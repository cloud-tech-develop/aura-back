package products

import (
	"context"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// PresentationRequest represents a presentation to be created with a product
type PresentationRequest struct {
	ID              *int64  `json:"id"`
	Name            string  `json:"name" binding:"required"`
	Factor          float64 `json:"factor" binding:"required"`
	Barcode         string  `json:"barcode"`
	SalePrice       float64 `json:"sale_price"`
	CostPrice       float64 `json:"cost_price"`
	DefaultPurchase bool    `json:"default_purchase"`
	DefaultSale     bool    `json:"default_sale"`
}

// Product entity
// Represents a product in the catalog with all inventory and pricing information
type Product struct {
	ID                 int64                 `json:"id"`
	SKU                string                `json:"sku"`
	Barcode            string                `json:"barcode"`
	Name               string                `json:"name" binding:"required"`
	Description        string                `json:"description"`
	CategoryID         *int64                `json:"category_id"`
	CategoryName       string                `json:"category_name"`
	BrandID            *int64                `json:"brand_id"`
	BrandName          string                `json:"brand_name"`
	UnitID             int64                 `json:"unit_measure_id" binding:"required"`
	UnitName           string                `json:"unit_name"`
	ProductType        string                `json:"product_type"`
	Active             bool                  `json:"active"`
	VisibleInPOS       bool                  `json:"visible_in_pos"`
	CostPrice          float64               `json:"cost_price" binding:"required"`
	SalePrice          float64               `json:"sale_price" binding:"required"`
	Price2             float64               `json:"price_2"`
	Price3             *float64              `json:"price_3"`
	IVAPercentage      float64               `json:"iva_percentage"`
	ConsumptionTax     float64               `json:"consumption_tax"`
	CurrentStock       int                   `json:"current_stock"`
	MinStock           int                   `json:"min_stock"`
	MaxStock           int                   `json:"max_stock"`
	ManagesInventory   bool                  `json:"manages_inventory"`
	ManagesBatches     bool                  `json:"manages_batches"`
	ManagesSerial      bool                  `json:"manages_serial"`
	AllowNegativeStock bool                  `json:"allow_negative_stock"`
	ImageURL           string                `json:"image_url"`
	EnterpriseID       int64                 `json:"enterprise_id"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          *time.Time            `json:"updated_at"`
	DeletedAt          *time.Time            `json:"deleted_at"`
	Presentations      []PresentationRequest `json:"presentations"`
}

// ValidProductTypes defines the allowed product type values
var ValidProductTypes = domain.ValidProductTypes

// IsValidProductType checks if the product type is valid
func IsValidProductType(productType string) bool {
	for _, t := range ValidProductTypes {
		if t == productType {
			return true
		}
	}
	return false
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
// Defines the data access layer for products
type Repository interface {
	Create(ctx context.Context, tenantSlug string, p *Product) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error)
	GetBySKU(ctx context.Context, tenantSlug string, sku string, enterpriseID int64) (*Product, error)
	GetByBarcode(ctx context.Context, tenantSlug string, barcode string, enterpriseID int64) (*Product, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Product, error)
	Update(ctx context.Context, tenantSlug string, p *Product) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}

// Service interface
// Defines the business logic layer for products
type Service interface {
	Create(ctx context.Context, tenantSlug string, p *Product) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error)
	GetBySKU(ctx context.Context, tenantSlug string, sku string, enterpriseID int64) (*Product, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Product, error)
	Update(ctx context.Context, tenantSlug string, id int64, p *Product) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}

// Domain Events
// Product lifecycle events for event-driven architecture
const (
	EventCreated = "product.created"
	EventUpdated = "product.updated"
	EventDeleted = "product.deleted"
)

// Event structures for product lifecycle
type CreatedEvent struct{ events.BaseEvent }
type UpdatedEvent struct{ events.BaseEvent }
type DeletedEvent struct{ events.BaseEvent }

// ToEventPayload converts product to event payload map
func (e *Product) ToEventPayload() map[string]interface{} {
	return map[string]interface{}{
		"id":          e.ID,
		"sku":         e.SKU,
		"barcode":     e.Barcode,
		"name":        e.Name,
		"description": e.Description,
		"active":      e.Active,
		"created_at":  e.CreatedAt,
		"updated_at":  e.UpdatedAt,
	}
}

// NewCreatedEvent creates a new product created event
func NewCreatedEvent(e *Product) CreatedEvent {
	return CreatedEvent{events.NewBaseEvent(EventCreated, e.ToEventPayload())}
}

// NewUpdatedEvent creates a new product updated event
func NewUpdatedEvent(e *Product) UpdatedEvent {
	return UpdatedEvent{events.NewBaseEvent(EventUpdated, e.ToEventPayload())}
}

// NewDeletedEvent creates a new product deleted event
func NewDeletedEvent(e *Product) DeletedEvent {
	return DeletedEvent{events.NewBaseEvent(EventDeleted, e.ToEventPayload())}
}
