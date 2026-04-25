package enterprise

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
)

type postgresRepository struct {
	q         db.Querier
	isOffline bool
}

func NewRepository(q db.Querier) Repository {
	return &postgresRepository{
		q:         q,
		isOffline: q.IsSQLite(),
	}
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
 
	prefix := r.q.SchemaPrefix("public")
	query := fmt.Sprintf(`
		INSERT INTO %senterprises 
		 (tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at) 
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) 
		 ON CONFLICT (slug) DO NOTHING`, prefix)
 
	_, err = r.q.ExecContext(ctx, query,
		e.TenantID, e.Name, e.CommercialName, e.Slug, e.SubDomain, e.Email,
		e.Document, e.DV, e.Phone, e.MunicipalityID, e.Municipality, e.Status,
		settingsJSON, e.CreatedAt, e.UpdatedAt,
	)
	return err
}
 
func (r *postgresRepository) GetBySlug(ctx context.Context, slug string) (*Enterprise, error) {
	var e Enterprise
	var settingsJSON []byte
	prefix := r.q.SchemaPrefix("public")
	query := fmt.Sprintf(`
		SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM %senterprises WHERE slug = $1 AND deleted_at IS NULL`, prefix)
 
	err := r.q.QueryRowContext(ctx, query, slug).Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
		&e.Email, &e.Document, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status,
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
	prefix := r.q.SchemaPrefix("public")
	query := fmt.Sprintf(`
		SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM %senterprises WHERE sub_domain = $1 AND status = 'ACTIVE' AND deleted_at IS NULL`, prefix)
 
	err := r.q.QueryRowContext(ctx, query, subDomain).Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
		&e.Email, &e.Document, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status,
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
	prefix := r.q.SchemaPrefix("public")
	query := fmt.Sprintf(`
		SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM %senterprises WHERE email = $1 AND deleted_at IS NULL`, prefix)
 
	err := r.q.QueryRowContext(ctx, query, email).Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
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
 
func (r *postgresRepository) List(ctx context.Context, params ListParams) (ListResult, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
 
	prefix := r.q.SchemaPrefix("public")
	// Build query with filters
	baseQuery := fmt.Sprintf("FROM %senterprises WHERE deleted_at IS NULL", prefix)
	var args []interface{}
	argIndex := 1
 
	if params.Status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, params.Status)
		argIndex++
	}
 
	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.q.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return ListResult{}, err
	}
 
	// Get paginated data
	offset := (params.Page - 1) * params.Limit
	dataQuery := fmt.Sprintf("SELECT id, name, commercial_name, slug, sub_domain, email, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", baseQuery, argIndex, argIndex+1)
	args = append(args, params.Limit, offset)
 
	rows, err := r.q.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return ListResult{}, err
	}
	defer rows.Close()
 
	var list []Enterprise
	for rows.Next() {
		var e Enterprise
		var settingsJSON []byte
		if err := rows.Scan(&e.ID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
			&e.Email, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality, &e.Status,
			&settingsJSON, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return ListResult{}, err
		}
		_ = json.Unmarshal(settingsJSON, &e.Settings)
		list = append(list, e)
	}
 
	totalPages := total / int64(params.Limit)
	if total%int64(params.Limit) > 0 {
		totalPages++
	}
 
	return ListResult{
		Data: list,
		Pagination: Pagination{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
 
func (r *postgresRepository) Update(ctx context.Context, e *Enterprise) error {
	e.UpdatedAt = time.Now()
	settingsJSON, err := json.Marshal(e.Settings)
	if err != nil || string(settingsJSON) == "null" {
		settingsJSON = []byte("{}")
	}
	prefix := r.q.SchemaPrefix("public")
	query := fmt.Sprintf(`
		UPDATE %senterprises 
		 SET name=$1, commercial_name=$2, sub_domain=$3, email=$4, dv=$5, phone=$6, 
		     municipality_id=$7, municipality=$8, status=$9, settings=$10, updated_at=$11 
		 WHERE id=$12 AND deleted_at IS NULL`, prefix)
 
	_, err = r.q.ExecContext(ctx, query,
		e.Name, e.CommercialName, e.SubDomain, e.Email, e.DV, e.Phone,
		e.MunicipalityID, e.Municipality, e.Status, settingsJSON, e.UpdatedAt, e.ID,
	)
	return err
}
 
func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	prefix := r.q.SchemaPrefix("public")
	query := fmt.Sprintf(`UPDATE %senterprises SET deleted_at = NOW() WHERE id = $1`, prefix)
	_, err := r.q.ExecContext(ctx, query, id)
	return err
}
 
// dropSchema is an internal helper used by service.
func (r *postgresRepository) dropSchema(slug string) {
	_, _ = r.q.ExecContext(context.Background(), fmt.Sprintf("DROP SCHEMA IF EXISTS %q CASCADE", slug))
}
 
// EmailExistsInUsers checks if email already exists in public.users (HU-002)
func (r *postgresRepository) EmailExistsInUsers(ctx context.Context, email vo.Email) (bool, error) {
	var exists bool
	prefix := r.q.SchemaPrefix("public")
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %susers WHERE email = $1)`, prefix)
	err := r.q.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}
 
// EnterpriseExistsByStatus checks if an enterprise exists with given status
func (r *postgresRepository) EnterpriseExistsByStatus(ctx context.Context, slug string, status string) (bool, error) {
	var exists bool
	prefix := r.q.SchemaPrefix("public")
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %senterprises WHERE slug = $1 AND status = $2 AND deleted_at IS NULL)`, prefix)
	err := r.q.QueryRowContext(ctx, query, slug, status).Scan(&exists)
	return exists, err
}
 
// GetPlanByEnterpriseID retrieves the plan for a given enterprise
func (r *postgresRepository) GetPlanByEnterpriseID(ctx context.Context, enterpriseID int64) (*Plan, error) {
	var p Plan
	prefix := r.q.SchemaPrefix("public")
	query := fmt.Sprintf(`
		SELECT id, enterprise_id, max_users, max_enterprises, trial_until, created_at, updated_at 
		 FROM %splans WHERE enterprise_id = $1 AND deleted_at IS NULL`, prefix)
	err := r.q.QueryRowContext(ctx, query, enterpriseID).Scan(&p.ID, &p.EnterpriseID, &p.MaxUsers, &p.MaxEnterprises, &p.TrialUntil, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
 
// GetPlansByEnterpriseID retrieves all plans for a given enterprise
func (r *postgresRepository) GetPlansByEnterpriseID(ctx context.Context, enterpriseID int64) ([]Plan, error) {
	prefix := r.q.SchemaPrefix("public")
	query := fmt.Sprintf(`
		SELECT id, enterprise_id, max_users, max_enterprises, trial_until, created_at, updated_at, deleted_at 
		 FROM %splans WHERE enterprise_id = $1 AND deleted_at IS NULL`, prefix)
 
	rows, err := r.q.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
 
	var plans []Plan
	for rows.Next() {
		var p Plan
		if err := rows.Scan(&p.ID, &p.EnterpriseID, &p.MaxUsers, &p.MaxEnterprises, &p.TrialUntil, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}
 
// CountEnterprisesByTenant counts the number of enterprises for a tenant
func (r *postgresRepository) CountEnterprisesByTenant(ctx context.Context, tenantID int64) (int64, error) {
	var count int64
	prefix := r.q.SchemaPrefix("public")
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %senterprises WHERE tenant_id = $1 AND deleted_at IS NULL`, prefix)
	err := r.q.QueryRowContext(ctx, query, tenantID).Scan(&count)
	return count, err
}
