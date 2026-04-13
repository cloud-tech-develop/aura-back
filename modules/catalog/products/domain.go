package products

import (
	"context"
	"time"
)

// Product entity
type Product struct {
	ID           int64      `json:"id"`
	SKU          string     `json:"sku"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	CategoryID   *int64     `json:"category_id,omitempty"`
	BrandID      *int64     `json:"brand_id,omitempty"`
	CostPrice    float64    `json:"cost_price"`
	SalePrice    float64    `json:"sale_price"`
	TaxRate      float64    `json:"tax_rate"`
	MinStock     int        `json:"min_stock"`
	CurrentStock int        `json:"current_stock"`
	ImageURL     string     `json:"image_url"`
	Status       string     `json:"status"`
	EnterpriseID int64      `json:"enterprise_id"`
	GlobalID     string     `json:"global_id"`
	SyncStatus   string     `json:"sync_status"`
	LastSyncedAt *time.Time `json:"last_synced_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// ListFilters for product queries
type ListFilters struct {
	Page       int
	Limit      int
	Search     string
	CategoryID *int64
	BrandID    *int64
}

// Repository interface
type Repository interface {
	Create(ctx context.Context, tenantSlug string, p *Product) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error)
	GetBySKU(ctx context.Context, tenantSlug string, sku string, enterpriseID int64) (*Product, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, first int64, rows int64, search string) ([]Product, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Product, error)
	Update(ctx context.Context, tenantSlug string, p *Product) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}

// Service interface
type Service interface {
	Create(ctx context.Context, tenantSlug string, p *Product) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, first int64, rows int64, search string) ([]Product, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Product, error)
	Update(ctx context.Context, tenantSlug string, id int64, p *Product) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}
