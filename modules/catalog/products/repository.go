package products

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/domain"
)

// querier defines the database query interface
// Supports both direct DB and transaction queries
type querier = db.Querier

// repository implements the Repository interface
// Handles all database operations for products
type repository struct {
	db        *db.DB
	isOffline bool
}

// NewRepository creates a new product repository instance
// db: database connection instance
func NewRepository(db *db.DB) Repository {
	return &repository{
		db:        db,
		isOffline: db.IsSQLite(),
	}
}

// Create inserts a new product into the database
func (r *repository) Create(ctx context.Context, tenantSlug string, p *Product) error {
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		INSERT INTO %sproduct (
			sku, barcode, name, description, category_id, brand_id, unit_id,
			product_type, active, visible_in_pos,
			cost_price, sale_price, price_2, price_3,
			iva_percentage, consumption_tax_value,
			current_stock, min_stock, max_stock,
			manages_inventory, manages_batches, manages_serial, allow_negative_stock,
			image_url, enterprise_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
		RETURNING id, created_at`, tenant)

	err := r.db.QueryRowContext(ctx, query,
		p.SKU, p.Barcode, p.Name, p.Description,
		p.CategoryID, p.BrandID, p.UnitID,
		p.ProductType, p.Active, p.VisibleInPOS,
		p.CostPrice, p.SalePrice, p.Price2, p.Price3,
		p.IVAPercentage, p.ConsumptionTax,
		p.CurrentStock, p.MinStock, p.MaxStock,
		p.ManagesInventory, p.ManagesBatches, p.ManagesSerial, p.AllowNegativeStock,
		p.ImageURL, p.EnterpriseID,
	).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetByID retrieves a product by its ID
func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error) {
	p := &Product{}
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		SELECT 
			id, sku, barcode, name, description, category_id, brand_id, unit_id,
			product_type, active, visible_in_pos,
			cost_price, sale_price, price_2, price_3,
			iva_percentage, consumption_tax_value,
			current_stock, min_stock, max_stock,
			manages_inventory, manages_batches, manages_serial, allow_negative_stock,
			image_url, enterprise_id,
			created_at, updated_at, deleted_at
		FROM %sproduct WHERE id = $1 AND deleted_at IS NULL`, tenant)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SKU, &p.Barcode, &p.Name, &p.Description,
		&p.CategoryID, &p.BrandID, &p.UnitID,
		&p.ProductType, &p.Active, &p.VisibleInPOS,
		&p.CostPrice, &p.SalePrice, &p.Price2, &p.Price3,
		&p.IVAPercentage, &p.ConsumptionTax,
		&p.CurrentStock, &p.MinStock, &p.MaxStock,
		&p.ManagesInventory, &p.ManagesBatches, &p.ManagesSerial, &p.AllowNegativeStock,
		&p.ImageURL, &p.EnterpriseID,
		&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return p, nil
}

// GetBySKU retrieves a product by its SKU code
func (r *repository) GetBySKU(ctx context.Context, tenantSlug string, sku string, enterpriseID int64) (*Product, error) {
	p := &Product{}
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		SELECT 
			id, sku, barcode, name, description, category_id, brand_id, unit_id,
			product_type, active, visible_in_pos,
			cost_price, sale_price, price_2, price_3,
			iva_percentage, consumption_tax_value,
			current_stock, min_stock, max_stock,
			manages_inventory, manages_batches, manages_serial, allow_negative_stock,
			image_url, enterprise_id,
			created_at, updated_at, deleted_at
		FROM %sproduct WHERE sku = $1 AND enterprise_id = $2 AND deleted_at IS NULL LIMIT 1`, tenant)

	err := r.db.QueryRowContext(ctx, query, sku, enterpriseID).Scan(
		&p.ID, &p.SKU, &p.Barcode, &p.Name, &p.Description,
		&p.CategoryID, &p.BrandID, &p.UnitID,
		&p.ProductType, &p.Active, &p.VisibleInPOS,
		&p.CostPrice, &p.SalePrice, &p.Price2, &p.Price3,
		&p.IVAPercentage, &p.ConsumptionTax,
		&p.CurrentStock, &p.MinStock, &p.MaxStock,
		&p.ManagesInventory, &p.ManagesBatches, &p.ManagesSerial, &p.AllowNegativeStock,
		&p.ImageURL, &p.EnterpriseID,
		&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}
	return p, nil
}

// GetByBarcode retrieves a product by its barcode
func (r *repository) GetByBarcode(ctx context.Context, tenantSlug string, barcode string, enterpriseID int64) (*Product, error) {
	p := &Product{}
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		SELECT 
			id, sku, barcode, name, description, category_id, brand_id, unit_id,
			product_type, active, visible_in_pos,
			cost_price, sale_price, price_2, price_3,
			iva_percentage, consumption_tax_value,
			current_stock, min_stock, max_stock,
			manages_inventory, manages_batches, manages_serial, allow_negative_stock,
			image_url, enterprise_id,
			created_at, updated_at, deleted_at
		FROM %sproduct WHERE barcode = $1 AND enterprise_id = $2 AND deleted_at IS NULL`, tenant)

	err := r.db.QueryRowContext(ctx, query, barcode, enterpriseID).Scan(
		&p.ID, &p.SKU, &p.Barcode, &p.Name, &p.Description,
		&p.CategoryID, &p.BrandID, &p.UnitID,
		&p.ProductType, &p.Active, &p.VisibleInPOS,
		&p.CostPrice, &p.SalePrice, &p.Price2, &p.Price3,
		&p.IVAPercentage, &p.ConsumptionTax,
		&p.CurrentStock, &p.MinStock, &p.MaxStock,
		&p.ManagesInventory, &p.ManagesBatches, &p.ManagesSerial, &p.AllowNegativeStock,
		&p.ImageURL, &p.EnterpriseID,
		&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get product by barcode: %w", err)
	}
	return p, nil
}

// List retrieves a list of products with filters
func (r *repository) List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Product, error) {
	// Prevents lib/pq connection state corruption when client cancels request
	ctx = context.WithoutCancel(ctx)

	tenant := r.db.SchemaPrefix(tenantSlug)
	baseWhere := fmt.Sprintf(`enterprise_id = %d AND deleted_at IS NULL AND active = true`, enterpriseID)

	// Apply search filter
	if filters.Search != "" {
		safeSearch := strings.ReplaceAll(filters.Search, "'", "''")
		baseWhere += fmt.Sprintf(" AND (name ILIKE '%%%s%%' OR sku ILIKE '%%%s%%' OR barcode ILIKE '%%%s%%')", safeSearch, safeSearch, safeSearch)
	}

	// Apply category filter
	if filters.CategoryID != nil {
		baseWhere += fmt.Sprintf(" AND category_id = %d", *filters.CategoryID)
	}

	// Apply brand filter
	if filters.BrandID != nil {
		baseWhere += fmt.Sprintf(" AND brand_id = %d", *filters.BrandID)
	}

	offset := (filters.Page - 1) * filters.Limit
	query := fmt.Sprintf(`
		SELECT 
			id, sku, barcode, name, description, category_id, brand_id, unit_id,
			product_type, active, visible_in_pos,
			cost_price, sale_price, price_2, price_3,
			iva_percentage, consumption_tax_value,
			current_stock, min_stock, max_stock,
			manages_inventory, manages_batches, manages_serial, allow_negative_stock,
			image_url, enterprise_id,
			created_at, updated_at, deleted_at
		FROM %sproduct WHERE `+baseWhere+` ORDER BY name LIMIT %d OFFSET %d`, tenant, filters.Limit, offset)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var list []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Barcode, &p.Name, &p.Description,
			&p.CategoryID, &p.BrandID, &p.UnitID,
			&p.ProductType, &p.Active, &p.VisibleInPOS,
			&p.CostPrice, &p.SalePrice, &p.Price2, &p.Price3,
			&p.IVAPercentage, &p.ConsumptionTax,
			&p.CurrentStock, &p.MinStock, &p.MaxStock,
			&p.ManagesInventory, &p.ManagesBatches, &p.ManagesSerial, &p.AllowNegativeStock,
			&p.ImageURL, &p.EnterpriseID,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// Page retrieves a paginated list of products with optional filters
func (r *repository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	// Prevents lib/pq connection state corruption when client cancels request
	ctx = context.WithoutCancel(ctx)

	tenant := r.db.SchemaPrefix(tenantSlug)
	// Build base WHERE clause
	baseWhere := fmt.Sprintf(`p.enterprise_id = %d AND p.deleted_at IS NULL`, enterpriseID)

	// Apply filters from params
	if params != nil {
		if categoryID, ok := params["category_id"]; ok && categoryID != nil {
			var catID int64
			switch v := categoryID.(type) {
			case float64:
				catID = int64(v)
			case int64:
				catID = v
			}
			baseWhere += fmt.Sprintf(" AND p.category_id = %d", catID)
		}
		if brandID, ok := params["brand_id"]; ok && brandID != nil {
			var bID int64
			switch v := brandID.(type) {
			case float64:
				bID = int64(v)
			case int64:
				bID = v
			}
			baseWhere += fmt.Sprintf(" AND p.brand_id = %d", bID)
		}
		if active, ok := params["active"]; ok && active != nil {
			if activeBool, isBool := active.(bool); isBool {
				baseWhere += fmt.Sprintf(" AND p.active = %v", activeBool)
			}
		}
		if visibleInPOS, ok := params["visible_in_pos"]; ok && visibleInPOS != nil {
			if visibleBool, isBool := visibleInPOS.(bool); isBool {
				baseWhere += fmt.Sprintf(" AND p.visible_in_pos = %v", visibleBool)
			}
		}
	}

	// Apply search filter
	searchCond := ""
	if search != "" {
		safeSearch := strings.ReplaceAll(search, "'", "''")
		searchCond = fmt.Sprintf(" AND (p.name ILIKE '%%%s%%' OR p.sku ILIKE '%%%s%%' OR p.barcode ILIKE '%%%s%%')", safeSearch, safeSearch, safeSearch)
	}

	// COUNT query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %sproduct AS p WHERE `+baseWhere+searchCond, tenant)

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to count products: %w", err)
	}

	// Validate sort column
	validSorts := map[string]string{
		"id":            "id",
		"name":          "name",
		"sku":           "sku",
		"barcode":       "barcode",
		"cost_price":    "cost_price",
		"sale_price":    "sale_price",
		"current_stock": "current_stock",
		"created_at":    "created_at",
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
		SELECT 
			p.id, p.sku, p.barcode, p.name, p.description, 
			p.category_id, c.name as category_name, 
			p.brand_id, b.name as brand_name, 
			p.unit_id, u.name as unit_name,
			p.product_type, p.active, p.visible_in_pos,
			p.cost_price, p.sale_price, p.price_2, p.price_3,
			p.iva_percentage, p.consumption_tax_value,
			p.current_stock, p.min_stock, p.max_stock,
			p.manages_inventory, p.manages_batches, p.manages_serial, p.allow_negative_stock,
			p.image_url, p.enterprise_id,
			p.created_at, p.updated_at, p.deleted_at
		FROM %sproduct AS p
		LEFT JOIN %scategory c ON p.category_id = c.id
		LEFT JOIN %sbrand b ON p.brand_id = b.id
		LEFT JOIN %sunit u ON p.unit_id = u.id
		WHERE `+baseWhere+searchCond+` ORDER BY p.%s %s LIMIT %d OFFSET %d`,
		tenant, tenant, tenant, tenant, sort, order, limit, offset)

	resultRows, err := r.db.QueryContext(ctx, selectQuery)

	if err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to page products: %w", err)
	}
	defer resultRows.Close()

	var list []Product
	for resultRows.Next() {
		var p Product
		if err := resultRows.Scan(
			&p.ID, &p.SKU, &p.Barcode, &p.Name, &p.Description,
			&p.CategoryID, &p.CategoryName,
			&p.BrandID, &p.BrandName,
			&p.UnitID, &p.UnitName,
			&p.ProductType, &p.Active, &p.VisibleInPOS,
			&p.CostPrice, &p.SalePrice, &p.Price2, &p.Price3,
			&p.IVAPercentage, &p.ConsumptionTax,
			&p.CurrentStock, &p.MinStock, &p.MaxStock,
			&p.ManagesInventory, &p.ManagesBatches, &p.ManagesSerial, &p.AllowNegativeStock,
			&p.ImageURL, &p.EnterpriseID,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
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

// Update updates an existing product in the database
func (r *repository) Update(ctx context.Context, tenantSlug string, p *Product) error {
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`
		UPDATE %sproduct SET 
			sku = $1, barcode = $2, name = $3, description = $4, 
			category_id = $5, brand_id = $6, unit_id = $7,
			product_type = $8, active = $9, visible_in_pos = $10,
			cost_price = $11, sale_price = $12, price_2 = $13, price_3 = $14,
			iva_percentage = $15, consumption_tax_value = $16,
			current_stock = $17, min_stock = $18, max_stock = $19,
			manages_inventory = $20, manages_batches = $21, manages_serial = $22, allow_negative_stock = $23,
			image_url = $24, updated_at = NOW()
		WHERE id = $25 AND deleted_at IS NULL`, tenant)

	_, err := r.db.ExecContext(ctx, query,
		p.SKU, p.Barcode, p.Name, p.Description,
		p.CategoryID, p.BrandID, p.UnitID,
		p.ProductType, p.Active, p.VisibleInPOS,
		p.CostPrice, p.SalePrice, p.Price2, p.Price3,
		p.IVAPercentage, p.ConsumptionTax,
		p.CurrentStock, p.MinStock, p.MaxStock,
		p.ManagesInventory, p.ManagesBatches, p.ManagesSerial, p.AllowNegativeStock,
		p.ImageURL, p.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

// Delete performs a soft delete of a product
func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	skuDel := "DEL-" + time.Now().Format("20060102150405")
	tenant := r.db.SchemaPrefix(tenantSlug)
	query := fmt.Sprintf(`UPDATE %sproduct SET deleted_at = NOW(), sku = $1 WHERE id = $2`, tenant)
	_, err := r.db.ExecContext(ctx, query, skuDel, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
