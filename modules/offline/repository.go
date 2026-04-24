package offline

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// Upsert inserts an enterprise if it doesn't exist (no overwrite)
func (r *repository) Upsert(ctx context.Context, e *Enterprise) error {
	// Check if enterprise already exists
	existing, err := r.GetBySlug(ctx, e.Slug)
	if err == nil && existing != nil {
		// Enterprise exists, don't overwrite
		return nil
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Set defaults
	if e.Status == "" {
		e.Status = "ACTIVE"
	}
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	settingsJSON, err := json.Marshal(e.Settings)
	if err != nil || string(settingsJSON) == "null" {
		settingsJSON = []byte("{}")
	}

	_, err = r.db.ExecContext(ctx,
		`INSERT INTO enterprises 
		 (id, tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at, deleted_at) 
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.TenantID, e.Name, e.CommercialName, e.Slug, e.SubDomain,
		e.Email, e.Document, e.DV, e.Phone, e.MunicipalityID, e.Municipality,
		e.Status, settingsJSON, e.CreatedAt, e.UpdatedAt, nil,
	)
	return err
}

// GetBySlug retrieves an enterprise by its slug
func (r *repository) GetBySlug(ctx context.Context, slug string) (*Enterprise, error) {
	var e Enterprise
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM enterprises WHERE slug = ? AND deleted_at IS NULL`,
		slug,
	).Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
		&e.Email, &e.Document, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality,
		&e.Status, &settingsJSON, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(settingsJSON, &e.Settings)
	return &e, nil
}

// List returns all enterprises stored locally
func (r *repository) List(ctx context.Context) ([]Enterprise, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM enterprises WHERE deleted_at IS NULL ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Enterprise
	for rows.Next() {
		var e Enterprise
		var settingsJSON []byte
		if err := rows.Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
			&e.Email, &e.Document, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality,
			&e.Status, &settingsJSON, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(settingsJSON, &e.Settings)
		list = append(list, e)
	}

	return list, nil
}
