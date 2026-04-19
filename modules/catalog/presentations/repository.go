package presentations

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
// Handles all database operations for presentations
type repository struct {
	db *db.DB
}

// NewRepository creates a new presentation repository instance
func NewRepository(db *db.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new presentation into the database
func (r *repository) Create(ctx context.Context, tenantSlug string, enterpriseID int64, p *Presentation) error {
	query := fmt.Sprintf(`
		INSERT INTO "%s".presentation (
			product_id, name, factor, barcode, cost_price, sale_price,
			default_purchase, default_sale, enterprise_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query,
		p.ProductID, p.Name, p.Factor, p.Barcode,
		p.CostPrice, p.SalePrice, p.DefaultPurchase, p.DefaultSale, enterpriseID,
	).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create presentation: %w", err)
	}
	return nil
}

// CreateMany inserts multiple presentations into the database
func (r *repository) CreateMany(ctx context.Context, tenantSlug string, enterpriseID int64, presentations []*Presentation) error {
	for _, p := range presentations {
		if err := r.Create(ctx, tenantSlug, enterpriseID, p); err != nil {
			return err
		}
	}
	return nil
}

// GetByID retrieves a presentation by its ID
func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Presentation, error) {
	p := &Presentation{}
	query := fmt.Sprintf(`
		SELECT 
			id, product_id, name, factor, barcode, cost_price, sale_price,
			default_purchase, default_sale, enterprise_id,
			created_at, updated_at, deleted_at
		FROM "%s".presentation WHERE id = $1 AND deleted_at IS NULL`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.ProductID, &p.Name, &p.Factor, &p.Barcode,
		&p.CostPrice, &p.SalePrice, &p.DefaultPurchase, &p.DefaultSale,
		&p.EnterpriseID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get presentation: %w", err)
	}
	return p, nil
}

// GetByProductID retrieves all presentations for a product
func (r *repository) GetByProductID(ctx context.Context, tenantSlug string, productID int64) ([]Presentation, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, product_id, name, factor, barcode, cost_price, sale_price,
			default_purchase, default_sale, enterprise_id,
			created_at, updated_at, deleted_at
		FROM "%s".presentation 
		WHERE product_id = $1 AND deleted_at IS NULL
		ORDER BY name`, tenantSlug)

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get presentations by product: %w", err)
	}
	defer rows.Close()

	var list []Presentation
	for rows.Next() {
		var p Presentation
		if err := rows.Scan(
			&p.ID, &p.ProductID, &p.Name, &p.Factor, &p.Barcode,
			&p.CostPrice, &p.SalePrice, &p.DefaultPurchase, &p.DefaultSale,
			&p.EnterpriseID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// List retrieves a list of presentations with filters
func (r *repository) List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Presentation, error) {
	ctx = context.WithoutCancel(ctx)

	baseWhere := fmt.Sprintf(`enterprise_id = %d AND deleted_at IS NULL`, enterpriseID)

	if filters.ProductID != nil {
		baseWhere += fmt.Sprintf(" AND product_id = %d", *filters.ProductID)
	}

	if filters.Search != "" {
		safeSearch := strings.ReplaceAll(filters.Search, "'", "''")
		baseWhere += fmt.Sprintf(" AND name ILIKE '%%%s%%'", safeSearch)
	}

	offset := (filters.Page - 1) * filters.Limit
	query := fmt.Sprintf(`
		SELECT 
			id, product_id, name, factor, barcode, cost_price, sale_price,
			default_purchase, default_sale, enterprise_id,
			created_at, updated_at, deleted_at
		FROM "%s".presentation WHERE `+baseWhere+` ORDER BY name LIMIT %d OFFSET %d`,
		tenantSlug, filters.Limit, offset)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list presentations: %w", err)
	}
	defer rows.Close()

	var list []Presentation
	for rows.Next() {
		var p Presentation
		if err := rows.Scan(
			&p.ID, &p.ProductID, &p.Name, &p.Factor, &p.Barcode,
			&p.CostPrice, &p.SalePrice, &p.DefaultPurchase, &p.DefaultSale,
			&p.EnterpriseID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// Page retrieves a paginated list of presentations
func (r *repository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	ctx = context.WithoutCancel(ctx)

	baseWhere := fmt.Sprintf(`enterprise_id = %d AND deleted_at IS NULL`, enterpriseID)

	// Apply filters from params
	if params != nil {
		if productID, ok := params["product_id"]; ok && productID != nil {
			var pID int64
			switch v := productID.(type) {
			case float64:
				pID = int64(v)
			case int64:
				pID = v
			}
			baseWhere += fmt.Sprintf(" AND product_id = %d", pID)
		}
	}

	// Apply search filter
	searchCond := ""
	if search != "" {
		safeSearch := strings.ReplaceAll(search, "'", "''")
		searchCond = fmt.Sprintf(" AND name ILIKE '%%%s%%'", safeSearch)
	}

	// COUNT query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "%s".presentation WHERE `+baseWhere+searchCond, tenantSlug)

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to count presentations: %w", err)
	}

	// Validate sort column
	validSorts := map[string]string{
		"id":         "id",
		"name":       "name",
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

	// SELECT query with pagination
	offset := (page - 1) * limit
	selectQuery := fmt.Sprintf(`
		SELECT 
			id, product_id, name, factor, barcode, cost_price, sale_price,
			default_purchase, default_sale, enterprise_id,
			created_at, updated_at, deleted_at
		FROM "%s".presentation 
		WHERE `+baseWhere+searchCond+` ORDER BY %s %s LIMIT %d OFFSET %d`,
		tenantSlug, sort, order, limit, offset)

	resultRows, err := r.db.QueryContext(ctx, selectQuery)
	if err != nil {
		return domain.PageResult{}, fmt.Errorf("failed to page presentations: %w", err)
	}
	defer resultRows.Close()

	var list []Presentation
	for resultRows.Next() {
		var p Presentation
		if err := resultRows.Scan(
			&p.ID, &p.ProductID, &p.Name, &p.Factor, &p.Barcode,
			&p.CostPrice, &p.SalePrice, &p.DefaultPurchase, &p.DefaultSale,
			&p.EnterpriseID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
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

// Update updates an existing presentation
func (r *repository) Update(ctx context.Context, tenantSlug string, p *Presentation) error {
	query := fmt.Sprintf(`
		UPDATE "%s".presentation SET 
			name = $1, factor = $2, barcode = $3, 
			cost_price = $4, sale_price = $5,
			default_purchase = $6, default_sale = $7, updated_at = NOW()
		WHERE id = $8 AND deleted_at IS NULL`, tenantSlug)

	_, err := r.db.ExecContext(ctx, query,
		p.Name, p.Factor, p.Barcode,
		p.CostPrice, p.SalePrice, p.DefaultPurchase, p.DefaultSale, p.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update presentation: %w", err)
	}
	return nil
}

// Delete performs a soft delete of a presentation
func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	ctx = context.WithoutCancel(ctx)
	nameProduct, err := r.GetByID(ctx, tenantSlug, id)
	if err != nil {
		return fmt.Errorf("failed to get presentation: %w", err)
	}
	nameDelete := fmt.Sprintf("%s_deleted_%s", nameProduct.Name, time.Now().Format("20060102150405"))
	query := fmt.Sprintf(`UPDATE "%s".presentation SET deleted_at = NOW(), name = $1 WHERE id = $2`, tenantSlug)
	_, err = r.db.ExecContext(ctx, query, nameDelete, id)
	if err != nil {
		return fmt.Errorf("failed to delete presentation: %w", err)
	}
	return nil
}
