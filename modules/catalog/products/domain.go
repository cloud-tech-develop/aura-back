package products

import (
	"context"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// Product entity
// Represents a product in the catalog with all inventory and pricing information
type Product struct {
	ID                 int64      `json:"id"`
	SKU                string     `json:"sku"`                  // Product SKU code (unique per enterprise)
	Barcode            string     `json:"barcode"`              // Product barcode for scanning
	Name               string     `json:"name"`                 // Product name
	Description        string     `json:"description"`          // Product description
	CategoryID         *int64     `json:"category_id"`          // Category foreign key reference
	BrandID            *int64     `json:"brand_id"`             // Brand foreign key reference
	UnitID             int64      `json:"unit_id"`              // Base unit of measure foreign key
	ProductType        string     `json:"product_type"`         // Product type: ESTANDAR, SERVICIO, COMBO, RECETA
	Active             bool       `json:"active"`               // Product active status
	VisibleInPOS       bool       `json:"visible_in_pos"`       // Visibility in POS interface
	CostPrice          float64    `json:"cost_price"`           // Product cost price (purchase price)
	SalePrice          float64    `json:"sale_price"`           // Product sale price (retail price)
	Price2             float64    `json:"price_2"`              // Alternative price level 2 (wholesale)
	Price3             *float64   `json:"price_3"`              // Alternative price level 3 (special)
	IVAPercentage      float64    `json:"iva_percentage"`       // IVA tax percentage
	ConsumptionTax     float64    `json:"consumption_tax"`      // Consumption tax percentage
	CurrentStock       int        `json:"current_stock"`        // Current inventory quantity
	MinStock           int        `json:"min_stock"`            // Minimum stock threshold for alerts
	MaxStock           int        `json:"max_stock"`            // Maximum stock level for inventory limits
	ManagesInventory   bool       `json:"manages_inventory"`    // Enable inventory tracking
	ManagesBatches     bool       `json:"manages_batches"`      // Enable batch/lot tracking
	ManagesSerial      bool       `json:"manages_serial"`       // Enable serial number tracking
	AllowNegativeStock bool       `json:"allow_negative_stock"` // Allow negative stock
	ImageURL           string     `json:"image_url"`            // Product image URL
	EnterpriseID       int64      `json:"enterprise_id"`        // Enterprise foreign key
	CreatedAt          time.Time  `json:"created_at"`           // Creation timestamp
	UpdatedAt          *time.Time `json:"updated_at"`           // Last update timestamp
	DeletedAt          *time.Time `json:"deleted_at"`           // Soft delete timestamp
}

// ValidProductTypes defines the allowed product type values
//
//	."ESTANDAR", "KIT", "PESABLE",    "SERVICIO"
var ValidProductTypes = []string{"ESTANDAR", "SERVICIO", "COMBO", "RECETA"}

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
