package products

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
)

type querier = db.Querier

type categoryRepository struct {
	db querier
}

func NewCategoryRepository(db querier) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, c *Category) error {
	query := `
		INSERT INTO category (name, description, parent_id, enterprise_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query, c.Name, c.Description, c.ParentID, c.EnterpriseID).
		Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id int64) (*Category, error) {
	c := &Category{}
	query := `
		SELECT id, name, description, parent_id, enterprise_id, created_at, updated_at, deleted_at
		FROM category WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.ParentID, &c.EnterpriseID,
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

func (r *categoryRepository) List(ctx context.Context, enterpriseID int64) ([]Category, error) {
	query := `
		SELECT id, name, description, parent_id, enterprise_id, created_at, updated_at, deleted_at
		FROM category WHERE enterprise_id = $1 AND deleted_at IS NULL
		ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.ParentID, &c.EnterpriseID,
			&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, c *Category) error {
	query := `
		UPDATE category SET name = $1, description = $2, parent_id = $3, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, c.Name, c.Description, c.ParentID, c.ID)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE category SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

// Brand Repository
type brandRepository struct {
	db querier
}

func NewBrandRepository(db querier) BrandRepository {
	return &brandRepository{db: db}
}

func (r *brandRepository) Create(ctx context.Context, b *Brand) error {
	query := `
		INSERT INTO brand (name, description, enterprise_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query, b.Name, b.Description, b.EnterpriseID).
		Scan(&b.ID, &b.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create brand: %w", err)
	}
	return nil
}

func (r *brandRepository) GetByID(ctx context.Context, id int64) (*Brand, error) {
	b := &Brand{}
	query := `
		SELECT id, name, description, enterprise_id, created_at, updated_at, deleted_at
		FROM brand WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.Name, &b.Description, &b.EnterpriseID,
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

func (r *brandRepository) List(ctx context.Context, enterpriseID int64) ([]Brand, error) {
	query := `
		SELECT id, name, description, enterprise_id, created_at, updated_at, deleted_at
		FROM brand WHERE enterprise_id = $1 AND deleted_at IS NULL
		ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brands: %w", err)
	}
	defer rows.Close()

	var brands []Brand
	for rows.Next() {
		var b Brand
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.EnterpriseID,
			&b.CreatedAt, &b.UpdatedAt, &b.DeletedAt); err != nil {
			return nil, err
		}
		brands = append(brands, b)
	}
	return brands, nil
}

func (r *brandRepository) Update(ctx context.Context, b *Brand) error {
	query := `
		UPDATE brand SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, b.Name, b.Description, b.ID)
	if err != nil {
		return fmt.Errorf("failed to update brand: %w", err)
	}
	return nil
}

func (r *brandRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE brand SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}
	return nil
}

// Product Repository
type productRepository struct {
	db querier
}

func NewProductRepository(db querier) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p *Product) error {
	query := `
		INSERT INTO product (sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query, p.SKU, p.Name, p.Description, p.CategoryID, p.BrandID,
		p.CostPrice, p.SalePrice, p.TaxRate, p.MinStock, p.CurrentStock, p.ImageURL, p.Status, p.EnterpriseID).
		Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id int64) (*Product, error) {
	p := &Product{}
	query := `
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM product WHERE id = $1 AND deleted_at IS NULL`

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

func (r *productRepository) GetBySKU(ctx context.Context, sku string, enterpriseID int64) (*Product, error) {
	p := &Product{}
	query := `
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM product WHERE sku = $1 AND enterprise_id = $2 AND deleted_at IS NULL`

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

func (r *productRepository) List(ctx context.Context, enterpriseID int64, filters ListFilters) ([]Product, error) {
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

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Description, &p.CategoryID, &p.BrandID,
			&p.CostPrice, &p.SalePrice, &p.TaxRate, &p.MinStock, &p.CurrentStock,
			&p.ImageURL, &p.Status, &p.EnterpriseID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) Update(ctx context.Context, p *Product) error {
	query := `
		UPDATE product SET sku = $1, name = $2, description = $3, category_id = $4, brand_id = $5,
		cost_price = $6, sale_price = $7, tax_rate = $8, min_stock = $9, current_stock = $10, image_url = $11, status = $12, updated_at = NOW()
		WHERE id = $13 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, p.SKU, p.Name, p.Description, p.CategoryID, p.BrandID,
		p.CostPrice, p.SalePrice, p.TaxRate, p.MinStock, p.CurrentStock, p.ImageURL, p.Status, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE product SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
