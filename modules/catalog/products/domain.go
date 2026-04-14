package products

import (
	"context"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/events"
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
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Product, error)
	Update(ctx context.Context, tenantSlug string, p *Product) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}

// Service interface
type Service interface {
	Create(ctx context.Context, tenantSlug string, p *Product) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Product, error)
	Update(ctx context.Context, tenantSlug string, id int64, p *Product) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}

// ─── Domain Events ────────────────────────────────────────────────────────
const (
	EventCreated = "product.created"
	EventUpdated = "product.updated"
	EventDeleted = "product.deleted"
)

type CreatedEvent struct{ events.BaseEvent }
type UpdatedEvent struct{ events.BaseEvent }
type DeletedEvent struct{ events.BaseEvent }

func (e *Product) ToEventPayload() map[string]interface{} {
	return map[string]interface{}{
		"id":          e.ID,
		"sku":         e.SKU,
		"name":        e.Name,
		"description": e.Description,
		"status":      e.Status,
		"created_at":  e.CreatedAt,
		"updated_at":  e.UpdatedAt,
	}
}

func NewCreatedEvent(e *Product) CreatedEvent {
	return CreatedEvent{events.NewBaseEvent(EventCreated, e.ToEventPayload())}
}
func NewUpdatedEvent(e *Product) UpdatedEvent {
	return UpdatedEvent{events.NewBaseEvent(EventUpdated, e.ToEventPayload())}
}
func NewDeletedEvent(e *Product) DeletedEvent {
	return DeletedEvent{events.NewBaseEvent(EventDeleted, e.ToEventPayload())}
}
