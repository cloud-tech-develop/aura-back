package presentations

import (
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// SyncPresentationPayload contains all fields needed for synchronization between online and offline
type SyncPresentationPayload struct {
	TenantSlug         string    `json:"tenant_slug"`
	PresentationID     int64     `json:"presentation_id"`
	ProductID          int64     `json:"product_id"`
	Name               string    `json:"name"`
	Factor             float64   `json:"factor"`
	Barcode            string    `json:"barcode"`
	CostPrice          float64   `json:"cost_price"`
	SalePrice          float64   `json:"sale_price"`
	DefaultPurchase    bool      `json:"default_purchase"`
	DefaultSale        bool      `json:"default_sale"`
	EnterpriseID       int64     `json:"enterprise_id"`
	Action             string    `json:"action"`    // "create", "update", "delete"
	Timestamp          time.Time `json:"timestamp"` // Timestamp of the change
	Source             string    `json:"source"`    // "online" or "offline"
}

// Sync constants for synchronization events
const (
	// Events for offline sync (online -> offline)
	EventPresentationOfflineCreated = "presentation.offline.created"
	EventPresentationOfflineUpdated = "presentation.offline.updated"
	EventPresentationOfflineDeleted = "presentation.offline.deleted"

	// Events for online sync (offline -> online)
	EventPresentationOnlineCreated = "presentation.online.created"
	EventPresentationOnlineUpdated = "presentation.online.updated"
	EventPresentationOnlineDeleted = "presentation.online.deleted"
)

// SyncEvent structures for synchronization
type SyncCreatedEvent struct{ events.BaseEvent }
type SyncUpdatedEvent struct{ events.BaseEvent }
type SyncDeletedEvent struct{ events.BaseEvent }

// ToSyncPayload converts Presentation to SyncPresentationPayload
func (p *Presentation) ToSyncPayload(action, source string) SyncPresentationPayload {
	return SyncPresentationPayload{
		TenantSlug:      "",
		PresentationID: p.ID,
		ProductID:      p.ProductID,
		Name:           p.Name,
		Factor:         p.Factor,
		Barcode:        p.Barcode,
		CostPrice:      p.CostPrice,
		SalePrice:      p.SalePrice,
		DefaultPurchase: p.DefaultPurchase,
		DefaultSale:    p.DefaultSale,
		EnterpriseID:   p.EnterpriseID,
		Action:         action,
		Timestamp:      time.Now(),
		Source:         source,
	}
}

// ToSyncPayloadWithTenant converts Presentation to SyncPresentationPayload with tenant slug
func (p *Presentation) ToSyncPayloadWithTenant(tenantSlug, action, source string) SyncPresentationPayload {
	pl := p.ToSyncPayload(action, source)
	pl.TenantSlug = tenantSlug
	return pl
}

// NewSyncCreatedEvent creates a new sync created event
func NewSyncCreatedEvent(tenantSlug string, p *Presentation) SyncCreatedEvent {
	return SyncCreatedEvent{events.NewBaseEvent(
		EventPresentationOfflineCreated,
		p.ToSyncPayloadWithTenant(tenantSlug, "create", "online"),
	)}
}

// NewSyncUpdatedEvent creates a new sync updated event
func NewSyncUpdatedEvent(tenantSlug string, p *Presentation) SyncUpdatedEvent {
	return SyncUpdatedEvent{events.NewBaseEvent(
		EventPresentationOfflineUpdated,
		p.ToSyncPayloadWithTenant(tenantSlug, "update", "online"),
	)}
}

// NewSyncDeletedEvent creates a new sync deleted event
func NewSyncDeletedEvent(tenantSlug string, p *Presentation) SyncDeletedEvent {
	return SyncDeletedEvent{events.NewBaseEvent(
		EventPresentationOfflineDeleted,
		p.ToSyncPayloadWithTenant(tenantSlug, "delete", "online"),
	)}
}

// NewSyncCreatedEventFromOffline creates a new sync created event from offline source
func NewSyncCreatedEventFromOffline(tenantSlug string, p *Presentation) SyncCreatedEvent {
	pl := p.ToSyncPayloadWithTenant(tenantSlug, "create", "offline")
	return SyncCreatedEvent{events.NewBaseEvent(EventPresentationOnlineCreated, pl)}
}

// NewSyncUpdatedEventFromOffline creates a new sync updated event from offline source
func NewSyncUpdatedEventFromOffline(tenantSlug string, p *Presentation) SyncUpdatedEvent {
	pl := p.ToSyncPayloadWithTenant(tenantSlug, "update", "offline")
	return SyncUpdatedEvent{events.NewBaseEvent(EventPresentationOnlineUpdated, pl)}
}

// NewSyncDeletedEventFromOffline creates a new sync deleted event from offline source
func NewSyncDeletedEventFromOffline(tenantSlug string, p *Presentation) SyncDeletedEvent {
	pl := p.ToSyncPayloadWithTenant(tenantSlug, "delete", "offline")
	return SyncDeletedEvent{events.NewBaseEvent(EventPresentationOnlineDeleted, pl)}
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	PresentationID int64  `json:"presentation_id"`
	Name           string `json:"name"`
	Action         string `json:"action"` // "created", "updated", "deleted", "conflict"
	Status         string `json:"status"` // "success", "error", "conflict"
	Message        string `json:"message"`
}