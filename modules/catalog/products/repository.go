package products

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
)

type querier = db.Querier

type repository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, tenantSlug string, p *Product) error {
	query := `
		INSERT INTO product (sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at`

	q := r.db.WithSchema(r.db, tenantSlug)
	err := q.QueryRowContext(ctx, query, p.SKU, p.Name, p.Description, p.CategoryID, p.BrandID,
		p.CostPrice, p.SalePrice, p.TaxRate, p.MinStock, p.CurrentStock, p.ImageURL, p.Status, p.EnterpriseID).
		Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error) {
	p := &Product{}
	query := `
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM product WHERE id = $1 AND deleted_at IS NULL`

	q := r.db.WithSchema(r.db, tenantSlug)
	err := q.QueryRowContext(ctx, query, id).Scan(
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
	query := `
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM product WHERE sku = $1 AND enterprise_id = $2 AND deleted_at IS NULL`

	q := r.db.WithSchema(r.db, tenantSlug)
	err := q.QueryRowContext(ctx, query, sku, enterpriseID).Scan(
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
	query := `
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM product WHERE enterprise_id = $1 AND deleted_at IS NULL`

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

	q := r.db.WithSchema(r.db, tenantSlug)
	rows, err := q.QueryContext(ctx, query, args...)
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

func (r *repository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, first int64, rows int64, search string) ([]Product, error) {
	query := `
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM product WHERE enterprise_id = $1 AND deleted_at IS NULL`

	args := []interface{}{enterpriseID}
	argPos := 2

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d)", argPos, argPos)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
		argPos += 2
	}

	query += " ORDER BY name LIMIT $" + fmt.Sprintf("%d", argPos)
	args = append(args, rows)
	argPos++

	offset := (first - 1) * rows
	query += " OFFSET $" + fmt.Sprintf("%d", argPos)
	args = append(args, offset)

	q := r.db.WithSchema(r.db, tenantSlug)
	resultRows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
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
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *repository) Update(ctx context.Context, tenantSlug string, p *Product) error {
	query := `
		UPDATE product SET sku = $1, name = $2, description = $3, category_id = $4, brand_id = $5,
		cost_price = $6, sale_price = $7, tax_rate = $8, min_stock = $9, current_stock = $10, image_url = $11, status = $12, updated_at = NOW()
		WHERE id = $13 AND deleted_at IS NULL`

	q := r.db.WithSchema(r.db, tenantSlug)
	_, err := q.ExecContext(ctx, query, p.SKU, p.Name, p.Description, p.CategoryID, p.BrandID,
		p.CostPrice, p.SalePrice, p.TaxRate, p.MinStock, p.CurrentStock, p.ImageURL, p.Status, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	query := `UPDATE product SET deleted_at = NOW() WHERE id = $1`
	q := r.db.WithSchema(r.db, tenantSlug)
	_, err := q.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
