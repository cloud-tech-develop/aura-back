package categories

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

func (r *repository) Create(ctx context.Context, tenantSlug string, c *Category) error {
	tenant := r.db.SchemaPrefix(tenantSlug)

	now := vo.Now()
	c.CreatedAt = now
	c.UpdatedAt = &now

	cols := []string{
		"name", "description", "parent_id", "default_tax_rate", "active", "enterprise_id",
		"created_at", "updated_at", "deleted_at",
	}
	args := categoryArgs(c, r.isOffline)

	var query string
	placeholders := make([]string, len(cols))
	for i := range cols {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	if r.isOffline {
		if c.ID > 0 {
			cols = append(cols, "id")
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(cols)))
		}
		query = fmt.Sprintf(
			"INSERT INTO %scategory (%s) VALUES (%s)",
			tenant,
			strings.Join(cols, ", "),
			strings.Join(placeholders, ", "),
		)
		_, err := r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to create category offline: %w", err)
		}
		return nil
	} else {
		query = fmt.Sprintf(
			"INSERT INTO %scategory (%s) VALUES (%s) RETURNING id",
			tenant,
			strings.Join(cols, ", "),
			strings.Join(placeholders, ", "),
		)
		err := r.db.QueryRowContext(ctx, query, args...).Scan(&c.ID)
		if err != nil {
			return fmt.Errorf("failed to create category production: %w", err)
		}
		return nil
	}
}

// categoryArgs returns the slice of arguments for INSERT queries.
// withID=true prepends c.ID as the first argument (used by SQLite always,
// and by Postgres when syncing an offline record that already has an ID).
func categoryArgs(c *Category, withID bool) []any {
	base := []any{
		c.Name, c.Description, c.ParentID, c.DefaultTaxRate, c.Active, c.EnterpriseID,
		c.CreatedAt, c.UpdatedAt, nil,
	}
	if withID && c.ID > 0 {
		return append(base, c.ID)
	}
	return base
}

func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Category, error) {
	c := &Category{}
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		SELECT id, name, description, parent_id, default_tax_rate, active, enterprise_id, created_at, updated_at, deleted_at
		FROM %scategory WHERE id = $1 AND deleted_at IS NULL`, tenant)

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

	conditionActive := "true"
	if r.isOffline {
		conditionActive = "1"
	}

	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		SELECT id, name
		FROM %scategory WHERE enterprise_id = $1 AND deleted_at IS NULL AND active = %s
		ORDER BY name`, tenant, conditionActive)

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
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

	tenant := r.db.SchemaPrefix(tenantSlug)
	// Build base WHERE clause
	baseWhere := `enterprise_id = $1 AND deleted_at IS NULL`
	args := []interface{}{enterpriseID}
	argPos := 2

	// COUNT query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %scategory WHERE `+baseWhere, tenant)
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
		FROM %scategory WHERE `+baseWhere, tenant)

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
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		UPDATE %scategory SET name = $1, description = $2, parent_id = $3, default_tax_rate = $4, active = $5, updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL`, tenant)

	_, err := r.db.ExecContext(ctx, query, c.Name, c.Description, c.ParentID, c.DefaultTaxRate, c.Active, c.ID)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`UPDATE %scategory SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, tenant)
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

func (r *repository) Upsert(ctx context.Context, tenantSlug string, c *Category) error {

	now := vo.Now()
	c.UpdatedAt = &now

	exist, _ := r.GetByID(ctx, tenantSlug, c.ID)
	if exist != nil {
		return r.Update(ctx, tenantSlug, c)
	}
	return r.Create(ctx, tenantSlug, c)
}
