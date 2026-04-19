package presentations

import (
	"context"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// Presentation entity
// Represents a product presentation/variant (kilo, libra, unidad, etc.)
type Presentation struct {
	ID              int64      `json:"id"`
	ProductID       int64      `json:"product_id"`       // Product foreign key
	Name            string     `json:"name"`             // Presentation name (Kilo, Libra, etc.)
	Factor          float64    `json:"factor"`           // Conversion factor to base unit
	Barcode         string     `json:"barcode"`          // Barcode for this presentation (optional)
	CostPrice       float64    `json:"cost_price"`       // Cost price for this presentation
	SalePrice       float64    `json:"sale_price"`       // Sale price for this presentation
	DefaultPurchase bool       `json:"default_purchase"` // Default for purchase orders
	DefaultSale     bool       `json:"default_sale"`     // Default for sales/POS
	EnterpriseID    int64      `json:"enterprise_id"`    // Enterprise foreign key
	CreatedAt       time.Time  `json:"created_at"`       // Creation timestamp
	UpdatedAt       *time.Time `json:"updated_at"`       // Last update timestamp
	DeletedAt       *time.Time `json:"deleted_at"`       // Soft delete timestamp
}

// PresentationWithProductInfo contains presentation data with product info for Page response
type PresentationWithProductInfo struct {
	ID              int64   `json:"id"`
	ProductID       int64   `json:"product_id"`
	ProductName     string  `json:"product_name"`
	Name            string  `json:"name"`
	Factor          float64 `json:"factor"`
	Barcode         string  `json:"barcode"`
	CostPrice       float64 `json:"cost_price"`
	SalePrice       float64 `json:"sale_price"`
	DefaultPurchase bool    `json:"default_purchase"`
	DefaultSale     bool    `json:"default_sale"`
}

// ListFilters for presentation queries
type ListFilters struct {
	Page      int
	Limit     int
	Search    string
	ProductID *int64
}

// PresentationRequest represents the JSON structure for creating presentations
type PresentationRequest struct {
	Name            string  `json:"name" binding:"required"`
	Factor          float64 `json:"factor" binding:"required"`
	Barcode         string  `json:"barcode"`
	CostPrice       float64 `json:"cost_price" binding:"required"`
	SalePrice       float64 `json:"sale_price" binding:"required"`
	DefaultPurchase bool    `json:"default_purchase"`
	DefaultSale     bool    `json:"default_sale"`
}

// PresentationListRequest represents the JSON structure for creating multiple presentations
type PresentationListRequest struct {
	Presentations []PresentationRequest `json:"presentations" binding:"required"`
}

// Repository interface
// Defines the data access layer for presentations
type Repository interface {
	Create(ctx context.Context, tenantSlug string, enterpriseID int64, p *Presentation) error
	CreateMany(ctx context.Context, tenantSlug string, enterpriseID int64, presentations []*Presentation) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Presentation, error)
	GetByProductID(ctx context.Context, tenantSlug string, productID int64) ([]Presentation, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Presentation, error)
	Update(ctx context.Context, tenantSlug string, p *Presentation) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}

// Service interface
// Defines the business logic layer for presentations
type Service interface {
	Create(ctx context.Context, tenantSlug string, enterpriseID int64, productID int64, presentations []PresentationRequest) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Presentation, error)
	GetByProductID(ctx context.Context, tenantSlug string, productID int64) ([]Presentation, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Presentation, error)
	Update(ctx context.Context, tenantSlug string, id int64, p *Presentation) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}

// Domain Events
// Presentation lifecycle events for event-driven architecture
const (
	EventCreated = "presentation.created"
	EventUpdated = "presentation.updated"
	EventDeleted = "presentation.deleted"
)

// Event structures for presentation lifecycle
type CreatedEvent struct{ events.BaseEvent }
type UpdatedEvent struct{ events.BaseEvent }
type DeletedEvent struct{ events.BaseEvent }

// ToEventPayload converts presentation to event payload map
func (e *Presentation) ToEventPayload() map[string]interface{} {
	return map[string]interface{}{
		"id":         e.ID,
		"product_id": e.ProductID,
		"name":       e.Name,
		"factor":     e.Factor,
		"barcode":    e.Barcode,
		"sale_price": e.SalePrice,
		"created_at": e.CreatedAt,
	}
}

// NewCreatedEvent creates a new presentation created event
func NewCreatedEvent(e *Presentation) CreatedEvent {
	return CreatedEvent{events.NewBaseEvent(EventCreated, e.ToEventPayload())}
}

// NewUpdatedEvent creates a new presentation updated event
func NewUpdatedEvent(e *Presentation) UpdatedEvent {
	return UpdatedEvent{events.NewBaseEvent(EventUpdated, e.ToEventPayload())}
}

// NewDeletedEvent creates a new presentation deleted event
func NewDeletedEvent(e *Presentation) DeletedEvent {
	return DeletedEvent{events.NewBaseEvent(EventDeleted, e.ToEventPayload())}
}
