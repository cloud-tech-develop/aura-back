package products

import (
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// SyncProductPayload contains all fields needed for synchronization between online and offline
// This payload is used for bidirectional sync via RabbitMQ
type SyncProductPayload struct {
	TenantSlug         string    `json:"tenant_slug"`
	ProductID          int64     `json:"product_id"`
	SKU                string    `json:"sku"`
	Barcode            string    `json:"barcode"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	CategoryID         *int64    `json:"category_id"`
	BrandID            *int64    `json:"brand_id"`
	UnitID             int64     `json:"unit_measure_id"`
	ProductType        string    `json:"product_type"`
	Active             bool      `json:"active"`
	VisibleInPOS       bool      `json:"visible_in_pos"`
	CostPrice          float64   `json:"cost_price"`
	SalePrice          float64   `json:"sale_price"`
	Price2             float64   `json:"price_2"`
	Price3             *float64  `json:"price_3"`
	IVAPercentage      float64   `json:"iva_percentage"`
	ConsumptionTax     float64   `json:"consumption_tax"`
	CurrentStock       int       `json:"current_stock"`
	MinStock           int       `json:"min_stock"`
	MaxStock           int       `json:"max_stock"`
	ManagesInventory   bool      `json:"manages_inventory"`
	ManagesBatches     bool      `json:"manages_batches"`
	ManagesSerial      bool      `json:"manages_serial"`
	AllowNegativeStock bool      `json:"allow_negative_stock"`
	ImageURL           string    `json:"image_url"`
	EnterpriseID       int64     `json:"enterprise_id"`
	Action             string    `json:"action"`    // "create", "update", "delete"
	Timestamp          time.Time `json:"timestamp"` // Timestamp of the change
	Source             string    `json:"source"`    // "online" or "offline"
}

// Sync constants for synchronization events
const (
	// Events for offline sync (online -> offline)
	EventProductOfflineCreated = "product.offline.created"
	EventProductOfflineUpdated = "product.offline.updated"
	EventProductOfflineDeleted = "product.offline.deleted"

	// Events for online sync (offline -> online)
	EventProductOnlineCreated = "product.online.created"
	EventProductOnlineUpdated = "product.online.updated"
	EventProductOnlineDeleted = "product.online.deleted"
)

// SyncEvent structures for synchronization
type SyncCreatedEvent struct{ events.BaseEvent }
type SyncUpdatedEvent struct{ events.BaseEvent }
type SyncDeletedEvent struct{ events.BaseEvent }

// ToSyncPayload converts Product to SyncProductPayload
func (p *Product) ToSyncPayload(action, source string) SyncProductPayload {
	return SyncProductPayload{
		TenantSlug:         "",
		ProductID:          p.ID,
		SKU:                p.SKU,
		Barcode:            p.Barcode,
		Name:               p.Name,
		Description:        p.Description,
		CategoryID:         p.CategoryID,
		BrandID:            p.BrandID,
		UnitID:             p.UnitID,
		ProductType:        p.ProductType,
		Active:             p.Active,
		VisibleInPOS:       p.VisibleInPOS,
		CostPrice:          p.CostPrice,
		SalePrice:          p.SalePrice,
		Price2:             p.Price2,
		Price3:             p.Price3,
		IVAPercentage:      p.IVAPercentage,
		ConsumptionTax:     p.ConsumptionTax,
		CurrentStock:       p.CurrentStock,
		MinStock:           p.MinStock,
		MaxStock:           p.MaxStock,
		ManagesInventory:   p.ManagesInventory,
		ManagesBatches:     p.ManagesBatches,
		ManagesSerial:      p.ManagesSerial,
		AllowNegativeStock: p.AllowNegativeStock,
		ImageURL:           p.ImageURL,
		EnterpriseID:       p.EnterpriseID,
		Action:             action,
		Timestamp:          time.Now(),
		Source:             source,
	}
}

// ToSyncPayloadWithTenant converts Product to SyncProductPayload with tenant slug
func (p *Product) ToSyncPayloadWithTenant(tenantSlug, action, source string) SyncProductPayload {
	pl := p.ToSyncPayload(action, source)
	pl.TenantSlug = tenantSlug
	return pl
}

// NewSyncCreatedEvent creates a new sync created event
func NewSyncCreatedEvent(tenantSlug string, p *Product) SyncCreatedEvent {
	return SyncCreatedEvent{events.NewBaseEvent(
		EventProductOfflineCreated,
		p.ToSyncPayloadWithTenant(tenantSlug, "create", "online"),
	)}
}

// NewSyncUpdatedEvent creates a new sync updated event
func NewSyncUpdatedEvent(tenantSlug string, p *Product) SyncUpdatedEvent {
	return SyncUpdatedEvent{events.NewBaseEvent(
		EventProductOfflineUpdated,
		p.ToSyncPayloadWithTenant(tenantSlug, "update", "online"),
	)}
}

// NewSyncDeletedEvent creates a new sync deleted event
func NewSyncDeletedEvent(tenantSlug string, p *Product) SyncDeletedEvent {
	return SyncDeletedEvent{events.NewBaseEvent(
		EventProductOfflineDeleted,
		p.ToSyncPayloadWithTenant(tenantSlug, "delete", "online"),
	)}
}

// NewSyncCreatedEventFromOffline creates a new sync created event from offline source
func NewSyncCreatedEventFromOffline(tenantSlug string, p *Product) SyncCreatedEvent {
	pl := p.ToSyncPayloadWithTenant(tenantSlug, "create", "offline")
	return SyncCreatedEvent{events.NewBaseEvent(EventProductOnlineCreated, pl)}
}

// NewSyncUpdatedEventFromOffline creates a new sync updated event from offline source
func NewSyncUpdatedEventFromOffline(tenantSlug string, p *Product) SyncUpdatedEvent {
	pl := p.ToSyncPayloadWithTenant(tenantSlug, "update", "offline")
	return SyncUpdatedEvent{events.NewBaseEvent(EventProductOnlineUpdated, pl)}
}

// NewSyncDeletedEventFromOffline creates a new sync deleted event from offline source
func NewSyncDeletedEventFromOffline(tenantSlug string, p *Product) SyncDeletedEvent {
	pl := p.ToSyncPayloadWithTenant(tenantSlug, "delete", "offline")
	return SyncDeletedEvent{events.NewBaseEvent(EventProductOnlineDeleted, pl)}
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	ProductID int64  `json:"product_id"`
	SKU       string `json:"sku"`
	Action    string `json:"action"` // "created", "updated", "deleted", "conflict"
	Status    string `json:"status"` // "success", "error", "conflict"
	Message   string `json:"message"`
}

// SyncRequest represents a sync request from offline to online
type SyncRequest struct {
	Products []SyncProductPayload `json:"products"`
}

// SyncResponse represents the response of a sync operation
type SyncResponse struct {
	Results  []SyncResult `json:"results"`
	SyncTime time.Time    `json:"sync_time"`
}
