package tenant

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations
var migrationsFS embed.FS

var validSlug = regexp.MustCompile(`^[a-z0-9_]+$`)

type Enterprise struct {
	ID             int64
	TenantID       int64  `json:"tenant_id"`
	Name           string `json:"name"` // Razón social
	CommercialName string `json:"commercial_name"`
	Slug           string
	SubDomain      string `json:"sub_domain"`
	Email          string
	DV             string `json:"dv"`
	Phone          string
	MunicipalityID string `json:"municipality_id"`
	Municipality   string
	Status         string
	Settings       map[string]interface{}
	CreatedAt      string
	UpdatedAt      string
	DeletedAt      *string `json:"deleted_at,omitempty"`
}

type Tenant struct {
	ID        int64
	Name      string
	Slug      string
	CreatedAt string
	DeletedAt *string
}

type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{db: db}
}

// MigratePublic aplica migraciones al esquema public (tabla de enterprises y tenants)
func (m *Manager) MigratePublic() error {
	return m.RunMigrations("public", "public")
}

// CreateEnterprise creates a new enterprise with its schema and linked tenant
func (m *Manager) CreateEnterprise(ctx context.Context, e *Enterprise, passwordHash string) error {
	e.Slug = strings.ToLower(e.Slug)
	if e.SubDomain != "" {
		e.SubDomain = strings.ToLower(e.SubDomain)
	}
	if !validSlug.MatchString(e.Slug) {
		return fmt.Errorf("slug inválido: solo minúsculas, números y _")
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	// 1. Create Tenant
	var tenantID int64
	err = tx.QueryRowContext(ctx,
		`INSERT INTO public.tenants (name, slug) VALUES ($1, $2) 
		 ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name 
		 RETURNING id`,
		e.Name, e.Slug,
	).Scan(&tenantID)
	if err != nil {
		return fmt.Errorf("crear tenant: %w", err)
	}
	e.TenantID = tenantID

	// 2. Create Schema
	if _, err := tx.ExecContext(ctx,
		fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %q", e.Slug),
	); err != nil {
		return fmt.Errorf("crear esquema: %w", err)
	}

	// 3. Register Enterprise
	var enterpriseID int64
	err = tx.QueryRowContext(ctx,
		`INSERT INTO public.enterprises (tenant_id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		 ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
		 RETURNING id`,
		e.TenantID, e.Name, e.CommercialName, e.Slug, e.SubDomain, e.Email, e.DV, e.Phone, e.MunicipalityID, e.Municipality,
	).Scan(&enterpriseID)
	if err != nil {
		return fmt.Errorf("registrar enterprise: %w", err)
	}
	e.ID = enterpriseID

	// 4. Create Initial User in public.users
	var userID int64
	err = tx.QueryRowContext(ctx,
		`INSERT INTO public.users (enterprise_id, name, email, password_hash) 
		 VALUES ($1, $2, $3, $4) 
		 ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
		 RETURNING id`,
		e.ID, "Admin", e.Email, passwordHash,
	).Scan(&userID)
	if err != nil {
		return fmt.Errorf("crear usuario inicial en public: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transacción: %w", err)
	}

	// 5. Run Tenant Migrations (creates third_parties and roles)
	if err := m.RunMigrations(e.Slug, "tenant"); err != nil {
		return err
	}

	// 6. Create Initial Third Party in tenant schema
	tpQuery := fmt.Sprintf(`
		INSERT INTO %q.third_parties (user_id, first_name, last_name, document_number, document_type, personal_email, tax_responsibility, is_employee) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, e.Slug)

	_, err = m.db.ExecContext(ctx, tpQuery,
		userID, "Admin", e.Name, e.DV, "NIT", e.Email, "RESPONSIBLE", true,
	)
	if err != nil {
		return fmt.Errorf("crear tercero inicial: %w", err)
	}

	return nil
}

// GetEnterpriseBySlug retrieves an enterprise by slug
func (m *Manager) GetEnterpriseBySlug(ctx context.Context, slug string) (*Enterprise, error) {
	var e Enterprise
	err := m.db.QueryRowContext(ctx,
		`SELECT id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, created_at, updated_at 
		 FROM public.enterprises WHERE slug = $1`,
		slug,
	).Scan(&e.ID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain, &e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// GetEnterpriseBySubDomain retrieves an enterprise by sub_domain
func (m *Manager) GetEnterpriseBySubDomain(ctx context.Context, subDomain string) (*Enterprise, error) {
	var e Enterprise
	err := m.db.QueryRowContext(ctx,
		`SELECT id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, created_at, updated_at 
		 FROM public.enterprises WHERE sub_domain = $1 AND status = 'ACTIVE'`,
		subDomain,
	).Scan(&e.ID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain, &e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// ListEnterprises returns all enterprises
func (m *Manager) ListEnterprises(ctx context.Context) ([]Enterprise, error) {
	rows, err := m.db.QueryContext(ctx,
		`SELECT id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, created_at, updated_at 
		 FROM public.enterprises ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enterprises []Enterprise
	for rows.Next() {
		var e Enterprise
		if err := rows.Scan(&e.ID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain, &e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		enterprises = append(enterprises, e)
	}
	return enterprises, nil
}

// Create registers a new tenant (legacy, creates enterprise instead)
func (m *Manager) Create(ctx context.Context, name, slug string) error {
	return m.CreateEnterprise(ctx, &Enterprise{
		Name:  name,
		Slug:  slug,
		Email: "",
	}, "")
}

// MigrateAll applies migrations to all existing enterprises
func (m *Manager) MigrateAll(ctx context.Context) error {
	rows, err := m.db.QueryContext(ctx, `SELECT slug FROM public.enterprises`)
	if err != nil {
		return fmt.Errorf("listar enterprises: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return err
		}
		if err := m.RunMigrations(slug, "tenant"); err != nil {
			return fmt.Errorf("migrando [%s]: %w", slug, err)
		}
	}
	return nil
}

// RunMigrations usa golang-migrate apuntando al esquema indicado y subdirectorio de migraciones
func (m *Manager) RunMigrations(schema, subPath string) error {
	// Obtener una conexión dedicada para garantizar que el search_path sea correcto
	conn, err := m.db.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("obtener conexión: %w", err)
	}
	defer conn.Close()

	// Forzar el search_path explícitamente en la conexión
	if _, err := conn.ExecContext(context.Background(),
		fmt.Sprintf("SET search_path TO %q", schema),
	); err != nil {
		return fmt.Errorf("set search_path [%s]: %w", schema, err)
	}

	driver, err := postgres.WithInstance(m.db, &postgres.Config{
		SchemaName:      schema,
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return fmt.Errorf("driver migrate: %w", err)
	}

	// Apuntar al subdirectorio correcto dentro del embed
	fullPath := fmt.Sprintf("migrations/%s", subPath)
	subFS, err := fs.Sub(migrationsFS, fullPath)
	if err != nil {
		return fmt.Errorf("sub fs [%s]: %w", fullPath, err)
	}

	src, err := iofs.New(subFS, ".")
	if err != nil {
		return fmt.Errorf("fuente de migraciones [%s]: %w", fullPath, err)
	}

	mg, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("inicializar migrate: %w", err)
	}

	if err := mg.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Printf("✓ esquema [%s/%s] sin cambios\n", schema, subPath)
			return nil
		}
		if strings.Contains(err.Error(), "Dirty database version") {
			version, _, _ := mg.Version()
			if forceErr := mg.Force(int(version) - 1); forceErr != nil {
				return fmt.Errorf("forzar versión [%s/%s]: %w", schema, subPath, forceErr)
			}
			if err2 := mg.Up(); err2 != nil && err2 != migrate.ErrNoChange {
				return fmt.Errorf("re-aplicar migraciones [%s/%s]: %w", schema, subPath, err2)
			}
			fmt.Printf("✓ esquema [%s/%s] migrado después de reset dirty\n", schema, subPath)
			return nil
		}
		return fmt.Errorf("aplicar migraciones [%s/%s]: %w", schema, subPath, err)
	}

	fmt.Printf("✓ esquema [%s/%s] migrado correctamente\n", schema, subPath)
	return nil
}
