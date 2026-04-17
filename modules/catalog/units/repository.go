package units

import (
	"context"
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

func (r *repository) Create(ctx context.Context, tenantSlug string, u *Unit) error {
	query := fmt.Sprintf(`
		INSERT INTO "%s".unit (name, abbreviation, active, allow_decimals, enterprise_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query, u.Name, u.Abbreviation, u.Active, u.AllowDecimals, u.EnterpriseID).
		Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create unit: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Unit, error) {
	u := &Unit{}
	query := fmt.Sprintf(`
		SELECT id, name, abbreviation, active, allow_decimals, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".unit WHERE id = $1 AND deleted_at IS NULL`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Abbreviation, &u.Active, &u.AllowDecimals, &u.EnterpriseID,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get unit: %w", err)
	}
	return u, nil
}

func (r *repository) List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]UnitList, error) {
	// Prevents lib/pq connection state corruption when client cancels request (e.g., hot-reload)
	ctx = context.WithoutCancel(ctx)

	query := fmt.Sprintf(`
		SELECT id, name, abbreviation
		FROM "%s".unit WHERE enterprise_id = $1 AND deleted_at IS NULL
		ORDER BY name`, tenantSlug)

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list units: %w", err)
	}
	defer rows.Close()

	var list []UnitList
	for rows.Next() {
		var u UnitList
		if err := rows.Scan(&u.Id, &u.Name, &u.Abbreviation); err != nil {
			return nil, err
		}
		list = append(list, u)
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
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "%s".unit WHERE `+baseWhere, tenantSlug)
	if search != "" {
		countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR abbreviation ILIKE $%d)", argPos, argPos)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm)
		argPos++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to count units: %w", err)
	}

	// Validate sort column
	validSorts := map[string]string{
		"id":           "id",
		"name":         "name",
		"abbreviation": "abbreviation",
		"active":       "active",
		"created_at":   "created_at",
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
		SELECT id, name, abbreviation, active, allow_decimals, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".unit WHERE `+baseWhere, tenantSlug)

	args = []interface{}{enterpriseID}
	argPos = 2

	if search != "" {
		selectQuery += fmt.Sprintf(" AND (name ILIKE $%d OR abbreviation ILIKE $%d)", argPos, argPos)
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
		return domain.PageResult{}, fmt.Errorf("failed to page units: %w", err)
	}
	defer resultRows.Close()

	var list []Unit
	for resultRows.Next() {
		var u Unit
		if err := resultRows.Scan(&u.ID, &u.Name, &u.Abbreviation, &u.Active, &u.AllowDecimals, &u.EnterpriseID,
			&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt); err != nil {
			return domain.PageResult{}, err
		}
		list = append(list, u)
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

func (r *repository) Update(ctx context.Context, tenantSlug string, u *Unit) error {
	query := fmt.Sprintf(`
		UPDATE "%s".unit SET name = $1, abbreviation = $2, active = $3, allow_decimals = $4, updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL`, tenantSlug)

	_, err := r.db.ExecContext(ctx, query, u.Name, u.Abbreviation, u.Active, u.AllowDecimals, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update unit: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	query := fmt.Sprintf(`UPDATE "%s".unit SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, tenantSlug)
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete unit: %w", err)
	}
	return nil
}
