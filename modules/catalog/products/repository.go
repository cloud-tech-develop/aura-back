package products

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	// Prevents lib/pq connection state corruption when client cancels request (e.g., hot-reload)
	ctx = context.WithoutCancel(ctx)

	baseWhere := fmt.Sprintf(`enterprise_id = %d AND deleted_at IS NULL`, enterpriseID)

	if filters.Search != "" {
		safeSearch := strings.ReplaceAll(filters.Search, "'", "''")
		baseWhere += fmt.Sprintf(" AND (name ILIKE '%%%s%%' OR sku ILIKE '%%%s%%')", safeSearch, safeSearch)
	}

	if filters.CategoryID != nil {
		baseWhere += fmt.Sprintf(" AND category_id = %d", *filters.CategoryID)
	}

	if filters.BrandID != nil {
		baseWhere += fmt.Sprintf(" AND brand_id = %d", *filters.BrandID)
	}

	offset := (filters.Page - 1) * filters.Limit
	query := fmt.Sprintf(`
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".product WHERE `+baseWhere+` ORDER BY name LIMIT %d OFFSET %d`, tenantSlug, filters.Limit, offset)

	rows, err := r.db.QueryContext(ctx, query)
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

func (r *repository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	// Prevents lib/pq connection state corruption when client cancels request (e.g., hot-reload)
	ctx = context.WithoutCancel(ctx)

	// Build base WHERE clause inlining variables to avoid PgBouncer statement pooling errors
	baseWhere := fmt.Sprintf(`enterprise_id = %d AND deleted_at IS NULL `, enterpriseID)

	// Apply params: category_id, brand_id, status
	if params != nil {
		if categoryID, ok := params["category_id"]; ok && categoryID != nil {
			var catID int64
			switch v := categoryID.(type) {
			case float64:
				catID = int64(v)
			case int64:
				catID = v
			}
			baseWhere += fmt.Sprintf(" AND category_id = %d", catID)
		}
		if brandID, ok := params["brand_id"]; ok && brandID != nil {
			var bID int64
			switch v := brandID.(type) {
			case float64:
				bID = int64(v)
			case int64:
				bID = v
			}
			baseWhere += fmt.Sprintf(" AND brand_id = %d", bID)
		}
		if status, ok := params["status"]; ok && status != nil {
			if statusStr, isStr := status.(string); isStr {
				safeStatus := strings.ReplaceAll(statusStr, "'", "''")
				baseWhere += fmt.Sprintf(" AND status = '%s'", safeStatus)
			}
		}
	}

	searchCond := ""
	if search != "" {
		safeSearch := strings.ReplaceAll(search, "'", "''")
		searchCond = fmt.Sprintf(" AND (name ILIKE '%%%s%%' OR sku ILIKE '%%%s%%')", safeSearch, safeSearch)
	}

	// COUNT query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "%s".product WHERE `+baseWhere+searchCond, tenantSlug)

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to count products: %w", err)
	}

	// Validate sort column
	validSorts := map[string]string{
		"id":         "id",
		"name":       "name",
		"sku":        "sku",
		"cost_price": "cost_price",
		"sale_price": "sale_price",
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

	offset := (page - 1) * limit
	selectQuery := fmt.Sprintf(`
		SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, created_at, updated_at, deleted_at
		FROM "%s".product WHERE `+baseWhere+searchCond+` ORDER BY %s %s LIMIT %d OFFSET %d`,
		tenantSlug, sort, order, limit, offset)

	resultRows, err := r.db.QueryContext(ctx, selectQuery)
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
