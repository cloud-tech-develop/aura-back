package tenant

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strings"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
)

//go:embed migrations
var migrationsFS embed.FS

var validSlug = regexp.MustCompile(`^[a-z0-9_]+$`)

// isSQLite indicates if the database driver is SQLite
var isSQLite = false

// SetSQLiteMode sets the manager to use SQLite mode for query adaptation
func SetSQLiteMode(sqlite bool) {
	isSQLite = sqlite
}

// AdaptQuery adapts a PostgreSQL query to SQLite-compatible syntax if in SQLite mode
func AdaptQuery(query string) string {
	if !isSQLite {
		return query
	}
	newQuery := db.AdaptQuery(query, "sqlite")

	newQuery = strings.ReplaceAll(newQuery, "RETURNING ", "")
	newQuery = strings.ReplaceAll(newQuery, "COMMENT ON COLUMN", "-- COMMENT ON COLUMN")
	newQuery = strings.ReplaceAll(newQuery, "COMMENT ON ", "-- COMMENT ON ")
	newQuery = strings.ReplaceAll(newQuery, "ON CONFLICT (slug) DO UPDATE SET", "OR REPLACE")
	newQuery = strings.ReplaceAll(newQuery, "ADD COLUMN IF NOT EXISTS", "ADD COLUMN")

	newQuery = removeCheckConstraints(newQuery)
	newQuery = removeReferences(newQuery)
	newQuery = removeFunctionDefinitions(newQuery)
	newQuery = removeTriggerDefinitions(newQuery)
	newQuery = removeDropTriggers(newQuery)

	return newQuery
}

func removeCheckConstraints(query string) string {
	// Remove CHECK (...) patterns
	result := query
	for {
		start := strings.Index(result, "CHECK (")
		if start == -1 {
			break
		}
		// Find matching closing parenthesis
		depth := 0
		end := start + 7 // After "CHECK ("
		for end < len(result) {
			if result[end] == '(' {
				depth++
			} else if result[end] == ')' {
				if depth == 0 {
					break
				}
				depth--
			}
			end++
		}
		// Remove CHECK constraint including comma before it if present
		before := strings.TrimRight(result[:start], " ,")
		result = before + result[end+1:]
	}
	return result
}

// removeReferences removes REFERENCES constraints (SQLite handles FK differently)
func removeReferences(query string) string {
	result := query

	// Process each CREATE TABLE block separately
	lines := strings.Split(result, "\n")
	var cleanLines []string
	inConstraint := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Start of constraint block
		if strings.HasPrefix(trimmed, "CONSTRAINT") {
			inConstraint = true
			// Remove trailing comma from the previous line
			if len(cleanLines) > 0 {
				lastIdx := len(cleanLines) - 1
				lastLine := strings.TrimRight(cleanLines[lastIdx], " \t,")
				cleanLines[lastIdx] = lastLine
			}
			continue
		}

		// If we're in a constraint block
		if inConstraint {
			// Check if this line ends the constraint (has ) followed by ;)
			if strings.Contains(trimmed, ")") && strings.HasSuffix(trimmed, ");") {
				inConstraint = false
				cleanLines = append(cleanLines, ");")
			}
			continue
		}

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		cleanLines = append(cleanLines, line)
	}

	result = strings.Join(cleanLines, "\n")

	// Remove standalone REFERENCES clauses from column definitions
	re := regexp.MustCompile(`\s+REFERENCES\s+[\w.]+\([^)]+\)`)
	result = re.ReplaceAllString(result, "")

	// Remove trailing commas before closing parenthesis
	result = regexp.MustCompile(`,\s*(\n\s*\))`).ReplaceAllString(result, "$1")

	// Clean up multiple blank lines
	for strings.Contains(result, "\n\n\n") {
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}

	// Final cleanup: remove lines that are just commas
	lines = strings.Split(result, "\n")
	cleanLines = nil
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && trimmed != "," {
			cleanLines = append(cleanLines, line)
		}
	}
	result = strings.Join(cleanLines, "\n")

	return result
}

// removeFunctionDefinitions removes CREATE FUNCTION blocks (SQLite doesn't support plpgsql)
func removeFunctionDefinitions(query string) string {
	result := query
	for {
		start := strings.Index(result, "CREATE")
		if start == -1 {
			break
		}
		// Check if it's a FUNCTION definition
		if strings.HasPrefix(strings.ToUpper(result[start:]), "CREATE OR REPLACE FUNCTION") ||
			strings.HasPrefix(strings.ToUpper(result[start:]), "CREATE FUNCTION") {
			// Find the end ($$ language or semicolon)
			end := strings.Index(result[start:], "$$ language")
			if end == -1 {
				end = strings.Index(result[start:], ";")
			}
			if end == -1 {
				break
			}
			result = result[:start] + result[start+end+len("$$ language 'plpgsql'")+1:]
		} else {
			break
		}
	}
	return result
}

// removeTriggerDefinitions removes CREATE TRIGGER blocks
func removeTriggerDefinitions(query string) string {
	result := query
	for {
		triggerIdx := strings.Index(result, "CREATE TRIGGER")
		if triggerIdx == -1 {
			break
		}
		// Find end of trigger (semicolon)
		end := strings.Index(result[triggerIdx:], ";")
		if end == -1 {
			break
		}
		result = result[:triggerIdx] + result[triggerIdx+end+1:]
	}
	return result
}

// removeDropTriggers removes DROP TRIGGER statements
func removeDropTriggers(query string) string {
	result := query
	for {
		dropIdx := strings.Index(result, "DROP TRIGGER")
		if dropIdx == -1 {
			break
		}
		// Find end of statement (semicolon)
		end := strings.Index(result[dropIdx:], ";")
		if end == -1 {
			break
		}
		result = result[:dropIdx] + result[dropIdx+end+1:]
	}
	return result
}

type Enterprise struct {
	ID             int64
	TenantID       int64  `json:"tenant_id"`
	Name           string `json:"name"` // Razón social
	CommercialName string `json:"commercial_name"`
	Slug           string
	SubDomain      string   `json:"sub_domain"`
	Email          vo.Email `json:"email"`
	Document       string   `json:"document"`
	DV             string   `json:"dv"`
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

// MigratePublic applies online (Postgres) public schema migrations.
// Called at startup in online mode.
func (m *Manager) MigratePublic() error {
	if isSQLite {
		return fmt.Errorf("MigratePublic is for online (Postgres) mode only; use MigrateOffline() instead")
	}
	return m.RunMigrations("public", "public")
}

// MigrateOffline runs all offline (SQLite) migrations: public tables + tenant tables.
// Called at startup in offline mode.
func (m *Manager) MigrateOffline() error {
	if !isSQLite {
		return fmt.Errorf("MigrateOffline is for offline (SQLite) mode only")
	}

	if err := m.RunOfflineMigrations("offline/public"); err != nil {
		return fmt.Errorf("offline public migrations: %w", err)
	}
	if err := m.RunOfflineMigrations("offline/tenant"); err != nil {
		return fmt.Errorf("offline tenant migrations: %w", err)
	}
	return nil
}

// RunOfflineMigrations applies SQLite-native .sql files from the given offline subdirectory.
// These files are already SQLite-compatible — no adaptation needed.
func (m *Manager) RunOfflineMigrations(subPath string) error {
	fullPath := fmt.Sprintf("migrations/%s", subPath)

	conn, err := m.db.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("obtener conexión: %w", err)
	}
	defer conn.Close()

	// Create version tracking table for this subPath
	safeTableName := "schema_offline_" + strings.ReplaceAll(strings.ReplaceAll(subPath, "/", "_"), "-", "_")
	_, _ = conn.ExecContext(context.Background(), fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version INTEGER PRIMARY KEY,
			applied_at TEXT DEFAULT CURRENT_TIMESTAMP
		)`, safeTableName))

	// Get current max version
	var version int64
	_ = conn.QueryRowContext(context.Background(),
		fmt.Sprintf("SELECT COALESCE(MAX(version), 0) FROM %s", safeTableName),
	).Scan(&version)

	// Read and sort migration files
	subFS, err := fs.Sub(migrationsFS, fullPath)
	if err != nil {
		return fmt.Errorf("sub fs [%s]: %w", fullPath, err)
	}
	migrationFiles, err := fs.ReadDir(subFS, ".")
	if err != nil {
		return fmt.Errorf("leer directorio [%s]: %w", fullPath, err)
	}

	var filesToApply []string
	for _, file := range migrationFiles {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		var fileVersion int64
		fmt.Sscanf(name, "%d", &fileVersion)
		if fileVersion <= version {
			continue
		}
		filesToApply = append(filesToApply, name)
	}
	sort.Strings(filesToApply)

	for _, name := range filesToApply {
		var fileVersion int64
		fmt.Sscanf(name, "%d", &fileVersion)

		fmt.Printf("➜ aplicando migración offline [%s/%s]\n", subPath, name)
		content, err := migrationsFS.ReadFile(fullPath + "/" + name)
		if err != nil {
			return fmt.Errorf("leer archivo [%s]: %w", name, err)
		}

		// Execute each statement separately (SQLite requires it)
		sqls := splitSQLStatements(string(content))
		for _, stmt := range sqls {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := conn.ExecContext(context.Background(), stmt); err != nil {
				return fmt.Errorf("ejecutar migración [%s]: %w\nSQL: %s", name, err, stmt)
			}
		}

		_, err = conn.ExecContext(context.Background(),
			fmt.Sprintf("INSERT INTO %s (version) VALUES (?)", safeTableName), fileVersion)
		if err != nil {
			return fmt.Errorf("registrar versión [%s]: %w", name, err)
		}
		fmt.Printf("✓ migración offline [%s/%s] aplicada\n", subPath, name)
	}

	return nil
}

// splitSQLStatements splits a SQL script into individual statements by semicolon,
// stripping comment lines from each segment.
func splitSQLStatements(script string) []string {
	var stmts []string
	for _, segment := range strings.Split(script, ";") {
		// Remove comment lines and blank lines from the segment
		var cleanLines []string
		for _, line := range strings.Split(segment, "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "--") {
				continue
			}
			cleanLines = append(cleanLines, line)
		}
		s := strings.TrimSpace(strings.Join(cleanLines, "\n"))
		if s != "" {
			stmts = append(stmts, s)
		}
	}
	return stmts
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

	// Validate slug length (HU-001)
	if len(e.Slug) < 3 || len(e.Slug) > 50 {
		return fmt.Errorf("slug debe tener entre 3 y 50 caracteres")
	}

	// Check if email already exists before starting transaction
	var existingID int64
	err := m.db.QueryRowContext(ctx,
		`SELECT id FROM public.enterprises WHERE email = $1 AND deleted_at IS NULL`,
		e.Email,
	).Scan(&existingID)
	if err == nil {
		return fmt.Errorf("el email %s ya está registrado", e.Email)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("verificando email: %w", err)
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
		`INSERT INTO public.enterprises (tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		 ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
		 RETURNING id`,
		e.TenantID, e.Name, e.CommercialName, e.Slug, e.SubDomain, e.Email, e.Document, e.DV, e.Phone, e.MunicipalityID, e.Municipality,
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

	// 4.1 Assign ADMIN role to the initial user (HU-003)
	_, err = tx.ExecContext(ctx,
		`INSERT INTO public.user_roles (user_id, role_id) VALUES ($1, 1)
		 ON CONFLICT DO NOTHING`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("asignar rol ADMIN: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transacción: %w", err)
	}

	// 5. Run Tenant Migrations (creates third_parties and roles)
	if err := m.RunMigrations(e.Slug, "tenant"); err != nil {
		return err
	}

	// 6. Create Initial Third Party in tenant schema using a dedicated connection
	conn, err := m.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("obtener conexión: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, fmt.Sprintf("SET search_path TO %q", e.Slug)); err != nil {
		return fmt.Errorf("set search_path: %w", err)
	}

	tpQuery := `
		INSERT INTO third_parties (user_id, first_name, last_name, document_number, document_type, personal_email, tax_responsibility, is_employee) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = conn.ExecContext(ctx, tpQuery,
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
		`SELECT id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, created_at, updated_at 
		 FROM public.enterprises WHERE slug = $1`,
		slug,
	).Scan(&e.ID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain, &e.Email, &e.Document, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status, &e.CreatedAt, &e.UpdatedAt)
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
		Email: vo.Email(""),
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

// RunMigrations applies online (Postgres) migrations for the given schema and subPath.
// This is the ONLINE-only migration runner. For offline/SQLite, use RunOfflineMigrations.
func (m *Manager) RunMigrations(schema, subPath string) error {

	conn, err := m.db.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("obtener conexión: %w", err)
	}
	defer conn.Close()

	// SET search_path is PostgreSQL-specific, skip for SQLite
	if !isSQLite {
		if _, err := conn.ExecContext(context.Background(), fmt.Sprintf("SET search_path TO %q", schema)); err != nil {
			return fmt.Errorf("set search_path: %w", err)
		}
	}

	// Adapt schema_migrations table for SQLite
	tableName := "schema_migrations"
	if isSQLite && subPath == "tenant" {
		tableName = "schema_migrations_tenant"
	}

	schemaMigrationsTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		version BIGSERIAL PRIMARY KEY,
		dirty boolean NOT NULL DEFAULT false,
		created_at TIMESTAMPTZ DEFAULT NOW()
	)`, tableName)
	if isSQLite {
		schemaMigrationsTable = AdaptQuery(schemaMigrationsTable)
	}
	_, _ = conn.ExecContext(context.Background(), schemaMigrationsTable)

	fullPath := fmt.Sprintf("migrations/%s", subPath)
	subFS, err := fs.Sub(migrationsFS, fullPath)
	if err != nil {
		return fmt.Errorf("sub fs [%s]: %w", fullPath, err)
	}

	migrationFiles, err := fs.ReadDir(subFS, ".")
	if err != nil {
		return fmt.Errorf("leer directorio [%s]: %w", fullPath, err)
	}

	var version int64
	versionQuery := "SELECT COALESCE(MAX(version), 0) FROM schema_migrations"
	if isSQLite {
		versionQuery = AdaptQuery(versionQuery)
	}
	err = conn.QueryRowContext(context.Background(), versionQuery).Scan(&version)
	if err != nil && err != sql.ErrNoRows {
		version = 0
	}

	var filesToMigrate []string
	for _, file := range migrationFiles {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		var fileVersion int64
		fmt.Sscanf(name, "%d", &fileVersion)
		if fileVersion <= version {
			continue
		}
		if !strings.HasSuffix(name, ".up.sql") && !strings.HasSuffix(name, ".offline.sql") {
			continue
		}
		filesToMigrate = append(filesToMigrate, name)
	}

	sort.Strings(filesToMigrate)

	// Postgres mode: only .up.sql files
	var finalFiles []string
	for _, name := range filesToMigrate {
		if strings.HasSuffix(name, ".up.sql") {
			finalFiles = append(finalFiles, name)
		}
	}

	for _, name := range finalFiles {
		var fileVersion int64
		fmt.Sscanf(name, "%d", &fileVersion)

		fmt.Printf("➜ aplicando migración [%s/%s]\n", schema, name)
		content, err := migrationsFS.ReadFile(fullPath + "/" + name)
		if err != nil {
			return fmt.Errorf("leer archivo [%s]: %w", name, err)
		}

		sql := string(content)
		if isSQLite {
			sql = AdaptQuery(sql)
		}

		if _, err := conn.ExecContext(context.Background(), sql); err != nil {
			return fmt.Errorf("ejecutar migración [%s]: %w", name, err)
		}

		insertVersion := fmt.Sprintf("INSERT INTO %s (version, dirty) VALUES ($1, false)", tableName)
		if isSQLite {
			insertVersion = AdaptQuery(insertVersion)
		}
		_, err = conn.ExecContext(context.Background(), insertVersion, fileVersion)
		if err != nil {
			return fmt.Errorf("registrar versión [%s]: %w", name, err)
		}
		fmt.Printf("✓ migración [%s/%s] aplicada\n", schema, name)
	}

	fmt.Printf("✓ esquema [%s/%s] migrado correctamente\n", schema, subPath)
	return nil
}
