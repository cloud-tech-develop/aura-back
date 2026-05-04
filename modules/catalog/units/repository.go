package units

import (
	"context"
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

func (r *repository) Upsert(ctx context.Context, tenantSlug string, u *Unit) error {
	now := vo.Now()
	u.UpdatedAt = &now

	exist, _ := r.GetByID(ctx, tenantSlug, u.ID)
	if exist != nil {
		return r.Update(ctx, tenantSlug, u)
	}

	return r.Create(ctx, tenantSlug, u)
}

func (r *repository) Create(ctx context.Context, tenantSlug string, u *Unit) error {
	tenant := r.db.SchemaPrefix(tenantSlug)

	now := vo.Now()
	u.CreatedAt = now
	u.UpdatedAt = &now

	cols := []string{
		"name", "abbreviation", "active", "allow_decimals", "enterprise_id",
		"created_at", "updated_at", "deleted_at",
	}
	args := unitArgs(u, r.isOffline)

	var query string
	placeholders := make([]string, len(cols))
	for i := range cols {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	if r.isOffline {
		if u.ID > 0 {
			cols = append(cols, "id")
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(cols)))
		}
		query = fmt.Sprintf(
			"INSERT INTO %sunit (%s) VALUES (%s)",
			tenant,
			strings.Join(cols, ", "),
			strings.Join(placeholders, ", "),
		)
		_, err := r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to create unit offline: %w", err)
		}
		return nil
	} else {
		query = fmt.Sprintf(
			"INSERT INTO %sunit (%s) VALUES (%s) RETURNING id",
			tenant,
			strings.Join(cols, ", "),
			strings.Join(placeholders, ", "),
		)
		err := r.db.QueryRowContext(ctx, query, args...).Scan(&u.ID)
		if err != nil {
			return fmt.Errorf("failed to create unit production: %w", err)
		}
		return nil
	}
}

// unitArgs returns the slice of arguments for INSERT queries.
// withID=true prepends u.ID as the first argument (used by SQLite always,
// and by Postgres when syncing an offline record that already has an ID).
func unitArgs(u *Unit, withID bool) []any {
	base := []any{
		u.Name, u.Abbreviation, u.Active, u.AllowDecimals, u.EnterpriseID,
		u.CreatedAt, u.UpdatedAt, nil,
	}
	if withID && u.ID > 0 {
		return append(base, u.ID)
	}
	return base
}

func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Unit, error) {
	u := &Unit{}
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		SELECT id, name, abbreviation, active, allow_decimals, enterprise_id, created_at, updated_at, deleted_at
		FROM %sunit WHERE id = $1 AND deleted_at IS NULL`, tenant)

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

	tenant := r.db.SchemaPrefix(tenantSlug)

	conditionActive := "true"
	if r.isOffline {
		conditionActive = "1"
	}

	query := fmt.Sprintf(`
		SELECT id, name, abbreviation
		FROM %sunit WHERE enterprise_id = $1 AND deleted_at IS NULL
		AND active = %s
		ORDER BY name`, tenant, conditionActive)

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

func (r *repository) ListAll(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Unit, error) {
	// Prevents lib/pq connection state corruption when client cancels request (e.g., hot-reload)
	ctx = context.WithoutCancel(ctx)

	tenant := r.db.SchemaPrefix(tenantSlug)

	conditionActive := "true"
	if r.isOffline {
		conditionActive = "1"
	}

	query := fmt.Sprintf(`
		SELECT id, name, abbreviation, active, allow_decimals, enterprise_id, created_at, updated_at, deleted_at
		FROM %sunit WHERE enterprise_id = $1 AND deleted_at IS NULL
		AND active = %s
		ORDER BY name`, tenant, conditionActive)

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list units: %w", err)
	}
	defer rows.Close()

	var list []Unit
	for rows.Next() {
		var u Unit
		if err := rows.Scan(&u.ID, &u.Name, &u.Abbreviation, &u.Active, &u.AllowDecimals, &u.EnterpriseID,
			&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
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
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %sunit WHERE `+baseWhere, tenant)
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
		FROM %sunit WHERE `+baseWhere, tenant)

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
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		UPDATE %sunit SET name = $1, abbreviation = $2, active = $3, allow_decimals = $4, updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL`, tenant)

	_, err := r.db.ExecContext(ctx, query, u.Name, u.Abbreviation, u.Active, u.AllowDecimals, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update unit: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`UPDATE %sunit SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, tenant)
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete unit: %w", err)
	}
	return nil
}
