package compositions

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
	db        *db.DB
	isOffline bool
}

func NewRepository(database *db.DB) Repository {
	return &repository{
		db:        database,
		isOffline: database.IsSQLite(),
	}
}

func (r *repository) Create(ctx context.Context, tenantSlug string, c *Composition) error {
	tenant := r.db.SchemaPrefix(tenantSlug)
	now := vo.Now()
	c.CreatedAt = now
	c.UpdatedAt = &now

	query := fmt.Sprintf(`
		INSERT INTO %sproduct_composition
			(parent_product_id, child_product_id, quantity, type, enterprise_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`, tenant)

	err := r.db.QueryRowContext(ctx, query,
		c.ParentProductID, c.ChildProductID, c.Quantity, c.Type, c.EnterpriseID,
		c.CreatedAt, c.UpdatedAt,
	).Scan(&c.ID)
	if err != nil {
		return fmt.Errorf("failed to create composition: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Composition, error) {
	tenant := r.db.SchemaPrefix(tenantSlug)
	joint := r.db.SchemaPrefix(tenantSlug)

	c := &Composition{}
	query := fmt.Sprintf(`
		SELECT
			pc.id, pc.parent_product_id, pc.child_product_id,
			COALESCE(p.name, '') AS parent_name,
			COALESCE(cp.name, '') AS child_name,
			pc.quantity, pc.type, pc.enterprise_id,
			pc.created_at, pc.updated_at
		FROM %sproduct_composition pc
		LEFT JOIN %sproduct p ON p.id = pc.parent_product_id
		LEFT JOIN %sproduct cp ON cp.id = pc.child_product_id
		WHERE pc.id = $1`, tenant, joint, joint)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.ParentProductID, &c.ChildProductID,
		&c.ParentName, &c.ChildName,
		&c.Quantity, &c.Type, &c.EnterpriseID,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get composition: %w", err)
	}
	return c, nil
}

func (r *repository) ListByParent(ctx context.Context, tenantSlug string, parentID int64) ([]Composition, error) {
	tenant := r.db.SchemaPrefix(tenantSlug)
	joint := r.db.SchemaPrefix(tenantSlug)

	query := fmt.Sprintf(`
		SELECT
			pc.id, pc.parent_product_id, pc.child_product_id,
			COALESCE(p.name, '') AS parent_name,
			COALESCE(cp.name, '') AS child_name,
			pc.quantity, pc.type, pc.enterprise_id,
			pc.created_at, pc.updated_at
		FROM %sproduct_composition pc
		LEFT JOIN %sproduct p ON p.id = pc.parent_product_id
		LEFT JOIN %sproduct cp ON cp.id = pc.child_product_id
		WHERE pc.parent_product_id = $1
		ORDER BY pc.id`, tenant, joint, joint)

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list compositions by parent: %w", err)
	}
	defer rows.Close()

	var list []Composition
	for rows.Next() {
		var c Composition
		if err := rows.Scan(
			&c.ID, &c.ParentProductID, &c.ChildProductID,
			&c.ParentName, &c.ChildName,
			&c.Quantity, &c.Type, &c.EnterpriseID,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *repository) ListAll(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Composition, error) {
	tenant := r.db.SchemaPrefix(tenantSlug)
	joint := r.db.SchemaPrefix(tenantSlug)

	query := fmt.Sprintf(`
		SELECT
			pc.id, pc.parent_product_id, pc.child_product_id,
			COALESCE(p.name, '') AS parent_name,
			COALESCE(cp.name, '') AS child_name,
			pc.quantity, pc.type, pc.enterprise_id,
			pc.created_at, pc.updated_at
		FROM %sproduct_composition pc
		LEFT JOIN %sproduct p ON p.id = pc.parent_product_id
		LEFT JOIN %sproduct cp ON cp.id = pc.child_product_id
		WHERE pc.enterprise_id = $1
		ORDER BY pc.id`, tenant, joint, joint)

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list all compositions: %w", err)
	}
	defer rows.Close()

	var list []Composition
	for rows.Next() {
		var c Composition
		if err := rows.Scan(
			&c.ID, &c.ParentProductID, &c.ChildProductID,
			&c.ParentName, &c.ChildName,
			&c.Quantity, &c.Type, &c.EnterpriseID,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *repository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, tipo string, sort string, order string) (domain.PageResult, error) {
	tenant := r.db.SchemaPrefix(tenantSlug)
	joint := r.db.SchemaPrefix(tenantSlug)

	baseWhere := fmt.Sprintf("pc.enterprise_id = %d", enterpriseID)

	// Apply type filter
	if tipo != "" {
		baseWhere += fmt.Sprintf(" AND pc.type = '%s'", strings.ReplaceAll(tipo, "'", "''"))
	}

	// Apply search filter
	searchCond := ""
	if search != "" {
		safeSearch := strings.ReplaceAll(search, "'", "''")
		searchCond = fmt.Sprintf(" AND (p.name ILIKE '%%%s%%' OR cp.name ILIKE '%%%s%%')", safeSearch, safeSearch)
	}

	// COUNT query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %sproduct_composition pc
		LEFT JOIN %sproduct p ON p.id = pc.parent_product_id
		LEFT JOIN %sproduct cp ON cp.id = pc.child_product_id
		WHERE `+baseWhere+searchCond, tenant, joint, joint)

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to count compositions: %w", err)
	}

	// Validate sort column
	validSorts := map[string]string{
		"id":              "pc.id",
		"quantity":        "pc.quantity",
		"type":            "pc.type",
		"parent_name":     "p.name",
		"child_name":      "cp.name",
		"created_at":      "pc.created_at",
	}
	if sortCol, ok := validSorts[sort]; ok {
		sort = sortCol
	} else {
		sort = "pc.id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// SELECT query with pagination
	offset := (page - 1) * limit
	selectQuery := fmt.Sprintf(`
		SELECT
			pc.id, pc.parent_product_id, pc.child_product_id,
			COALESCE(p.name, '') AS parent_name,
			COALESCE(cp.name, '') AS child_name,
			pc.quantity, pc.type, pc.enterprise_id,
			pc.created_at, pc.updated_at
		FROM %sproduct_composition pc
		LEFT JOIN %sproduct p ON p.id = pc.parent_product_id
		LEFT JOIN %sproduct cp ON cp.id = pc.child_product_id
		WHERE `+baseWhere+searchCond+` ORDER BY %s %s LIMIT %d OFFSET %d`,
		tenant, joint, joint, sort, order, limit, offset)

	rows, err := r.db.QueryContext(ctx, selectQuery)
	if err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to page compositions: %w", err)
	}
	defer rows.Close()

	var list []Composition
	for rows.Next() {
		var c Composition
		if err := rows.Scan(
			&c.ID, &c.ParentProductID, &c.ChildProductID,
			&c.ParentName, &c.ChildName,
			&c.Quantity, &c.Type, &c.EnterpriseID,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return domain.PageResult{}, err
		}
		list = append(list, c)
	}

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

func (r *repository) Update(ctx context.Context, tenantSlug string, c *Composition) error {
	tenant := r.db.SchemaPrefix(tenantSlug)

	query := fmt.Sprintf(`
		UPDATE %sproduct_composition SET
			quantity = $1, type = $2
		WHERE id = $3`, tenant)

	result, err := r.db.ExecContext(ctx, query, c.Quantity, c.Type, c.ID)
	if err != nil {
		return fmt.Errorf("failed to update composition: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	tenant := r.db.SchemaPrefix(tenantSlug)

	query := fmt.Sprintf(`DELETE FROM %sproduct_composition WHERE id = $1`, tenant)
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete composition: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) ExistsPair(ctx context.Context, tenantSlug string, parentID int64, childID int64) (bool, error) {
	tenant := r.db.SchemaPrefix(tenantSlug)

	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1 FROM %sproduct_composition
			WHERE parent_product_id = $1 AND child_product_id = $2
		)`, tenant)

	var exists bool
	err := r.db.QueryRowContext(ctx, query, parentID, childID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check composition existence: %w", err)
	}
	return exists, nil
}
