package enterprise

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
)

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, e *Enterprise) error {
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
		`INSERT INTO public.enterprises 
		 (tenant_id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at) 
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) 
		 ON CONFLICT (slug) DO NOTHING`,
		e.TenantID, e.Name, e.CommercialName, e.Slug, e.SubDomain, e.Email,
		e.DV, e.Phone, e.MunicipalityID, e.Municipality, e.Status,
		settingsJSON, e.CreatedAt, e.UpdatedAt,
	)
	return err
}

func (r *postgresRepository) GetBySlug(ctx context.Context, slug string) (*Enterprise, error) {
	var e Enterprise
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM public.enterprises WHERE slug = $1 AND deleted_at IS NULL`,
		slug,
	).Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
		&e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status,
		&settingsJSON, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(settingsJSON, &e.Settings)
	return &e, nil
}

func (r *postgresRepository) GetBySubDomain(ctx context.Context, subDomain string) (*Enterprise, error) {
	var e Enterprise
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM public.enterprises WHERE sub_domain = $1 AND status = 'ACTIVE' AND deleted_at IS NULL`,
		subDomain,
	).Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
		&e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status,
		&settingsJSON, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(settingsJSON, &e.Settings)
	return &e, nil
}

func (r *postgresRepository) GetByEmail(ctx context.Context, email vo.Email) (*Enterprise, error) {
	var e Enterprise
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM public.enterprises WHERE email = $1 AND deleted_at IS NULL`,
		email,
	).Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
		&e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status,
		&settingsJSON, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	_ = json.Unmarshal(settingsJSON, &e.Settings)
	return &e, nil
}

func (r *postgresRepository) List(ctx context.Context) ([]Enterprise, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM public.enterprises WHERE deleted_at IS NULL ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Enterprise
	for rows.Next() {
		var e Enterprise
		var settingsJSON []byte
		if err := rows.Scan(&e.ID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
			&e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status,
			&settingsJSON, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(settingsJSON, &e.Settings)
		list = append(list, e)
	}
	return list, nil
}

func (r *postgresRepository) Update(ctx context.Context, e *Enterprise) error {
	e.UpdatedAt = time.Now()
	settingsJSON, err := json.Marshal(e.Settings)
	if err != nil || string(settingsJSON) == "null" {
		settingsJSON = []byte("{}")
	}
	_, err = r.db.ExecContext(ctx,
		`UPDATE public.enterprises 
		 SET name=$1, commercial_name=$2, sub_domain=$3, email=$4, dv=$5, phone=$6, 
		     municipality_id=$7, municipality=$8, status=$9, settings=$10, updated_at=$11 
		 WHERE id=$12 AND deleted_at IS NULL`,
		e.Name, e.CommercialName, e.SubDomain, e.Email, e.DV, e.Phone,
		e.MunicipalityID, e.Municipality, e.Status, settingsJSON, e.UpdatedAt, e.ID,
	)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE public.enterprises SET deleted_at = NOW() WHERE id = $1`, id,
	)
	return err
}

// dropSchema is an internal helper used by service.
func (r *postgresRepository) dropSchema(slug string) {
	r.db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %q CASCADE", slug))
}
