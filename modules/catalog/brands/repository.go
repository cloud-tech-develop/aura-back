package brands

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

func (r *repository) Create(ctx context.Context, tenantSlug string, b *Brand) error {
	query := fmt.Sprintf(`
		INSERT INTO "%s".brand (name, description, active, enterprise_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query, b.Name, b.Description, b.Active, b.EnterpriseID).
		Scan(&b.ID, &b.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create brand: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Brand, error) {
	b := &Brand{}
	query := fmt.Sprintf(`
		SELECT id, name, description, active, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".brand WHERE id = $1 AND deleted_at IS NULL`, tenantSlug)

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

func (r *repository) List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Brand, error) {
	// Prevents lib/pq connection state corruption when client cancels request (e.g., hot-reload)
	ctx = context.WithoutCancel(ctx)

	query := fmt.Sprintf(`
		SELECT id, name, description, active, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".brand WHERE enterprise_id = $1 AND deleted_at IS NULL AND active = true
		ORDER BY name`, tenantSlug)

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brands: %w", err)
	}
	defer rows.Close()

	var list []Brand
	for rows.Next() {
		var b Brand
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Active, &b.EnterpriseID,
			&b.CreatedAt, &b.UpdatedAt, &b.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *repository) Update(ctx context.Context, tenantSlug string, b *Brand) error {
	query := fmt.Sprintf(`
		UPDATE "%s".brand SET name = $1, description = $2, active = $3, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL`, tenantSlug)

	_, err := r.db.ExecContext(ctx, query, b.Name, b.Description, b.Active, b.ID)
	if err != nil {
		return fmt.Errorf("failed to update brand: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	query := fmt.Sprintf(`UPDATE "%s".brand SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, tenantSlug)
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}
	return nil
}

func (r *repository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	// Prevents lib/pq connection state corruption when client cancels request (e.g., hot-reload)
	ctx = context.WithoutCancel(ctx)

	// Build base WHERE clause
	baseWhere := `enterprise_id = $1 AND deleted_at IS NULL`
	args := []interface{}{enterpriseID}
	argPos := 2

	// COUNT query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "%s".brand WHERE `+baseWhere, tenantSlug)
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
		FROM "%s".brand WHERE `+baseWhere, tenantSlug)

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
