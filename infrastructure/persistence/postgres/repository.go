package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/cloud-tech-develop/aura-back/domain/enterprise"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, e *enterprise.Enterprise) error {
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	if e.Status == "" {
		e.Status = "ACTIVE"
	}

	settingsJSON, err := json.Marshal(e.Settings)
	if err != nil || string(settingsJSON) == "null" {
		settingsJSON = []byte("{}")
	}

	_, err = r.db.ExecContext(ctx,
		`INSERT INTO public.enterprises (tenant_id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) ON CONFLICT (slug) DO NOTHING`,
		e.TenantID, e.Name, e.CommercialName, e.Slug, e.SubDomain, e.Email, e.DV, e.Phone, e.MunicipalityID, e.Municipality, e.Status, settingsJSON, e.CreatedAt, e.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetBySlug(ctx context.Context, slug string) (*enterprise.Enterprise, error) {
	var e enterprise.Enterprise
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM public.enterprises WHERE slug = $1`,
		slug,
	).Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain, &e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status, &settingsJSON, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if len(settingsJSON) > 0 {
		_ = json.Unmarshal(settingsJSON, &e.Settings)
	}
	return &e, nil
}

func (r *PostgresRepository) GetBySubDomain(ctx context.Context, subDomain string) (*enterprise.Enterprise, error) {
	var e enterprise.Enterprise
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM public.enterprises WHERE sub_domain = $1 AND status = 'ACTIVE'`,
		subDomain,
	).Scan(&e.ID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain, &e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status, &settingsJSON, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if len(settingsJSON) > 0 {
		_ = json.Unmarshal(settingsJSON, &e.Settings)
	}
	return &e, nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*enterprise.Enterprise, error) {
	var e enterprise.Enterprise
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM public.enterprises WHERE email = $1`,
		email,
	).Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain, &e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status, &settingsJSON, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if len(settingsJSON) > 0 {
		_ = json.Unmarshal(settingsJSON, &e.Settings)
	}
	return &e, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]enterprise.Enterprise, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM public.enterprises ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enterprises []enterprise.Enterprise
	for rows.Next() {
		var e enterprise.Enterprise
		var settingsJSON []byte
		if err := rows.Scan(&e.ID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain, &e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status, &settingsJSON, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		if len(settingsJSON) > 0 {
			_ = json.Unmarshal(settingsJSON, &e.Settings)
		}
		enterprises = append(enterprises, e)
	}
	return enterprises, nil
}

func (r *PostgresRepository) Update(ctx context.Context, e *enterprise.Enterprise) error {
	settingsJSON, err := json.Marshal(e.Settings)
	if err != nil || string(settingsJSON) == "null" {
		settingsJSON = []byte("{}")
	}
	_, err = r.db.ExecContext(ctx,
		`UPDATE public.enterprises SET name=$1, commercial_name=$2, sub_domain=$3, email=$4, dv=$5, phone=$6, municipality_id=$7, municipality=$8, status=$9, settings=$10, updated_at=$11 WHERE id=$12`,
		e.Name, e.CommercialName, e.SubDomain, e.Email, e.DV, e.Phone, e.MunicipalityID, e.Municipality, e.Status, settingsJSON, e.UpdatedAt, e.ID,
	)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(`DELETE FROM public.enterprises WHERE id = $1`, id)
	return err
}

func (r *PostgresRepository) Exec(query string, args ...interface{}) error {
	_, err := r.db.Exec(query, args...)
	return err
}
