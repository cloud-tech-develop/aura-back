package categories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/domain"
)

type querier = db.Querier

type repository struct {
	db querier
}

func NewRepository(db querier) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, tenantSlug string, c *Category) error {
	query := fmt.Sprintf(`
		INSERT INTO "%s".category (name, description, parent_id, default_tax_rate, active, enterprise_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query, c.Name, c.Description, c.ParentID, c.DefaultTaxRate, c.Active, c.EnterpriseID).
		Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Category, error) {
	c := &Category{}
	query := fmt.Sprintf(`
		SELECT id, name, description, parent_id, default_tax_rate, active, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".category WHERE id = $1 AND deleted_at IS NULL`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.ParentID, &c.DefaultTaxRate, &c.Active, &c.EnterpriseID,
		&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return c, nil
}

func (r *repository) List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]domain.ListId, error) {
	// Prevents lib/pq connection state corruption when client cancels request (e.g., hot-reload)
	ctx = context.WithoutCancel(ctx)

	query := fmt.Sprintf(`
		SELECT id, name
		FROM "%s".category WHERE enterprise_id = %d AND deleted_at IS NULL AND active = true
		ORDER BY name`, tenantSlug, enterpriseID)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var list []domain.ListId
	for rows.Next() {
		var c domain.ListId
		if err := rows.Scan(&c.Id, &c.Name); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *repository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	// Prevents lib/pq connection state corruption when client cancels request (e.g., hot-reload)
	ctx = context.WithoutCancel(ctx)

	// Build base WHERE clause
	baseWhere := `enterprise_id = $1 AND deleted_at IS NULL`
	args := []interface{}{enterpriseID}
	argPos := 2

	// COUNT query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "%s".category WHERE `+baseWhere, tenantSlug)
	if search != "" {
		countQuery += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm)
		argPos++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to count categories: %w", err)
	}

	// Validate sort column (only allow safe columns)
	validSorts := map[string]string{
		"id":               "id",
		"name":             "name",
		"created_at":       "created_at",
		"default_tax_rate": "default_tax_rate",
	}
	if sortCol, ok := validSorts[sort]; ok {
		sort = sortCol
	} else {
		sort = "id"
	}
	// Validate order direction
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// SELECT query with pagination
	selectQuery := fmt.Sprintf(`
		SELECT id, name, description, parent_id, default_tax_rate, active, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".category WHERE `+baseWhere, tenantSlug)

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
		return domain.PageResult{}, fmt.Errorf("failed to page categories: %w", err)
	}
	defer resultRows.Close()

	var list []Category
	for resultRows.Next() {
		var c Category
		if err := resultRows.Scan(&c.ID, &c.Name, &c.Description, &c.ParentID, &c.DefaultTaxRate, &c.Active, &c.EnterpriseID,
			&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return domain.PageResult{}, err
		}
		list = append(list, c)
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

func (r *repository) Update(ctx context.Context, tenantSlug string, c *Category) error {
	query := fmt.Sprintf(`
		UPDATE "%s".category SET name = $1, description = $2, parent_id = $3, default_tax_rate = $4, active = $5, updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL`, tenantSlug)

	_, err := r.db.ExecContext(ctx, query, c.Name, c.Description, c.ParentID, c.DefaultTaxRate, c.Active, c.ID)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	query := fmt.Sprintf(`UPDATE "%s".category SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, tenantSlug)
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}
