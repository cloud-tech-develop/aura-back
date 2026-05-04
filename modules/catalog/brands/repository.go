package brands

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
)

type querier = db.Querier

type repository struct {
	db        querier
	isOffline bool
}

func NewRepository(db querier) Repository {
	return &repository{
		db:        db,
		isOffline: db.IsSQLite(),
	}
}

func (r *repository) Create(ctx context.Context, tenantSlug string, b *Brand) error {
	tenant := r.db.SchemaPrefix(tenantSlug)

	now := vo.Now()
	b.CreatedAt = now
	b.UpdatedAt = &now

	cols := []string{
		"name", "description", "active", "enterprise_id",
		"created_at", "updated_at", "deleted_at",
	}
	args := brandArgs(b, r.isOffline)

	var query string
	placeholders := make([]string, len(cols))
	for i := range cols {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	if r.isOffline {
		if b.ID > 0 {
			cols = append(cols, "id")
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(cols)))
		}
		query = fmt.Sprintf(
			"INSERT INTO %sbrand (%s) VALUES (%s)",
			tenant,
			strings.Join(cols, ", "),
			strings.Join(placeholders, ", "),
		)
		_, err := r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to create brand offline: %w", err)
		}
		return nil
	} else {
		query = fmt.Sprintf(
			"INSERT INTO %sbrand (%s) VALUES (%s) RETURNING id",
			tenant,
			strings.Join(cols, ", "),
			strings.Join(placeholders, ", "),
		)
		err := r.db.QueryRowContext(ctx, query, args...).Scan(&b.ID)
		if err != nil {
			return fmt.Errorf("failed to create brand production: %w", err)
		}
		return nil
	}
}

// brandArgs returns the slice of arguments for INSERT queries.
// withID=true prepends b.ID as the first argument (used by SQLite always,
// and by Postgres when syncing an offline record that already has an ID).
func brandArgs(b *Brand, withID bool) []any {
	base := []any{
		b.Name, b.Description, b.Active, b.EnterpriseID,
		b.CreatedAt, b.UpdatedAt, nil,
	}
	if withID && b.ID > 0 {
		return append(base, b.ID)
	}
	return base
}

func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Brand, error) {
	b := &Brand{}
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		SELECT id, name, description, active, enterprise_id, created_at, updated_at, deleted_at
		FROM %sbrand WHERE id = $1 AND deleted_at IS NULL`, tenant)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.Name, &b.Description, &b.Active, &b.EnterpriseID,
		&b.CreatedAt, &b.UpdatedAt, &b.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get brand: %w", err)
	}
	return b, nil
}

func (r *repository) List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]BrandList, error) {
	// Prevents lib/pq connection state corruption when client cancels request (e.g., hot-reload)
	ctx = context.WithoutCancel(ctx)

	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		SELECT id, name 
		FROM %sbrand WHERE enterprise_id = $1
			AND deleted_at IS NULL
			AND active = true
		ORDER BY name;`, tenant)

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brands: %w", err)
	}
	defer rows.Close()

	var list []BrandList
	for rows.Next() {
		var b BrandList
		if err := rows.Scan(&b.ID, &b.Name); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *repository) Update(ctx context.Context, tenantSlug string, b *Brand) error {
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		UPDATE %sbrand SET name = $1, description = $2, active = $3, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL`, tenant)

	_, err := r.db.ExecContext(ctx, query, b.Name, b.Description, b.Active, b.ID)
	if err != nil {
		return fmt.Errorf("failed to update brand: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`UPDATE %sbrand SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, tenant)
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}
	return nil
}

func (r *repository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	// Prevents lib/pq connection state corruption when client cancels request (e.g., hot-reload)
	ctx = context.WithoutCancel(ctx)

	tenant := r.db.SchemaPrefix(tenantSlug)
	// Build base WHERE clause
	baseWhere := `enterprise_id = $1 AND deleted_at IS NULL`
	args := []interface{}{enterpriseID}
	argPos := 2

	// COUNT query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %sbrand WHERE `+baseWhere, tenant)
	if search != "" {
		countQuery += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm)
		argPos++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to count brands: %w", err)
	}

	// Validate sort column
	validSorts := map[string]string{
		"id":         "id",
		"name":       "name",
		"created_at": "created_at",
	}
	if sortCol, ok := validSorts[sort]; ok {
		sort = sortCol
	} else {
		sort = "id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// SELECT query with pagination
	selectQuery := fmt.Sprintf(`
		SELECT id, name, description, active, enterprise_id, created_at, updated_at, deleted_at
		FROM %sbrand WHERE `+baseWhere, tenant)

	args = []interface{}{enterpriseID}
	argPos = 2

	if search != "" {
		selectQuery += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm)
		argPos++
	}

	selectQuery += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d", sort, order, argPos)
	args = append(args, limit)
	argPos++

	offset := (page - 1) * limit
	selectQuery += fmt.Sprintf(" OFFSET $%d", argPos)
	args = append(args, offset)

	resultRows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to page brands: %w", err)
	}
	defer resultRows.Close()

	var list []Brand
	for resultRows.Next() {
		var b Brand
		if err := resultRows.Scan(&b.ID, &b.Name, &b.Description, &b.Active, &b.EnterpriseID,
			&b.CreatedAt, &b.UpdatedAt, &b.DeletedAt); err != nil {
			return domain.PageResult{}, err
		}
		list = append(list, b)
	}

	// Calculate pagination
	totalPages := (total + limit - 1) / limit
	if total == 0 {
		totalPages = 0
	}

	return domain.PageResult{
		Items:      list,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *repository) Upsert(ctx context.Context, tenantSlug string, b *Brand) error {
	now := vo.Now()
	b.UpdatedAt = &now

	exist, _ := r.GetByID(ctx, tenantSlug, b.ID)
	if exist != nil {
		return r.Update(ctx, tenantSlug, b)
	}
	return r.Create(ctx, tenantSlug, b)
}
