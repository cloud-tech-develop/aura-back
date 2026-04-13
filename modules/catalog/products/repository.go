package products

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/domain"
)

type querier = db.Querier

type repository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, tenantSlug string, p *Product) error {
	query := fmt.Sprintf(`
		INSERT INTO "%s".product (sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query, p.SKU, p.Name, p.Description, p.CategoryID, p.BrandID,
		p.CostPrice, p.SalePrice, p.TaxRate, p.MinStock, p.CurrentStock, p.ImageURL, p.Status, p.EnterpriseID).
		Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error) {
	p := &Product{}
	query := fmt.Sprintf(`
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".product WHERE id = $1 AND deleted_at IS NULL`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Description, &p.CategoryID, &p.BrandID,
		&p.CostPrice, &p.SalePrice, &p.TaxRate, &p.MinStock, &p.CurrentStock,
		&p.ImageURL, &p.Status, &p.EnterpriseID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return p, nil
}

func (r *repository) GetBySKU(ctx context.Context, tenantSlug string, sku string, enterpriseID int64) (*Product, error) {
	p := &Product{}
	query := fmt.Sprintf(`
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".product WHERE sku = $1 AND enterprise_id = $2 AND deleted_at IS NULL`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query, sku, enterpriseID).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Description, &p.CategoryID, &p.BrandID,
		&p.CostPrice, &p.SalePrice, &p.TaxRate, &p.MinStock, &p.CurrentStock,
		&p.ImageURL, &p.Status, &p.EnterpriseID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}
	return p, nil
}

func (r *repository) List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Product, error) {
	query := fmt.Sprintf(`
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".product WHERE enterprise_id = $1 AND deleted_at IS NULL`, tenantSlug)

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d)", argPos, argPos)
		searchTerm := "%" + filters.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argPos += 2
	}

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argPos)
		args = append(args, *filters.CategoryID)
		argPos++
	}

	if filters.BrandID != nil {
		query += fmt.Sprintf(" AND brand_id = $%d", argPos)
		args = append(args, *filters.BrandID)
		argPos++
	}

	query += " ORDER BY name LIMIT $" + fmt.Sprintf("%d", argPos)
	args = append(args, filters.Limit)
	argPos++

	offset := (filters.Page - 1) * filters.Limit
	query += " OFFSET $" + fmt.Sprintf("%d", argPos)
	args = append(args, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var list []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Description, &p.CategoryID, &p.BrandID,
			&p.CostPrice, &p.SalePrice, &p.TaxRate, &p.MinStock, &p.CurrentStock,
			&p.ImageURL, &p.Status, &p.EnterpriseID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *repository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, first int64, rows int64, search string) (domain.PageResult, error) {
	// Build base WHERE clause - only active products
	baseWhere := `enterprise_id = $1 AND deleted_at IS NULL AND status = 'ACTIVE'`
	args := []interface{}{enterpriseID}
	argPos := 2

	// COUNT query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "%s".product WHERE `+baseWhere, tenantSlug)
	if search != "" {
		countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d)", argPos, argPos)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
		argPos += 2
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to count products: %w", err)
	}

	// SELECT query with pagination
	selectQuery := fmt.Sprintf(`
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".product WHERE `+baseWhere, tenantSlug)

	args = []interface{}{enterpriseID}
	argPos = 2

	if search != "" {
		selectQuery += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d)", argPos, argPos)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
		argPos += 2
	}

	selectQuery += " ORDER BY name LIMIT $" + fmt.Sprintf("%d", argPos)
	args = append(args, rows)
	argPos++

	offset := (first - 1) * rows
	selectQuery += " OFFSET $" + fmt.Sprintf("%d", argPos)
	args = append(args, offset)

	resultRows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to page products: %w", err)
	}
	defer resultRows.Close()

	var list []Product
	for resultRows.Next() {
		var p Product
		if err := resultRows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Description, &p.CategoryID, &p.BrandID,
			&p.CostPrice, &p.SalePrice, &p.TaxRate, &p.MinStock, &p.CurrentStock,
			&p.ImageURL, &p.Status, &p.EnterpriseID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return domain.PageResult{}, err
		}
		list = append(list, p)
	}

	// Calculate pagination
	page := first
	limit := rows
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

func (r *repository) Update(ctx context.Context, tenantSlug string, p *Product) error {
	query := fmt.Sprintf(`
		UPDATE "%s".product SET sku = $1, name = $2, description = $3, category_id = $4, brand_id = $5,
		cost_price = $6, sale_price = $7, tax_rate = $8, min_stock = $9, current_stock = $10, image_url = $11, status = $12, updated_at = NOW()
		WHERE id = $13 AND deleted_at IS NULL`, tenantSlug)

	_, err := r.db.ExecContext(ctx, query, p.SKU, p.Name, p.Description, p.CategoryID, p.BrandID,
		p.CostPrice, p.SalePrice, p.TaxRate, p.MinStock, p.CurrentStock, p.ImageURL, p.Status, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	query := fmt.Sprintf(`UPDATE "%s".product SET deleted_at = NOW() WHERE id = $1`, tenantSlug)
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
