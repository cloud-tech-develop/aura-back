package brands

import (
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// SyncBrandPayload contains all fields needed for synchronization between online and offline
type SyncBrandPayload struct {
	TenantSlug   string    `json:"tenant_slug"`
	BrandID      int64     `json:"brand_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Active       bool      `json:"active"`
	EnterpriseID int64     `json:"enterprise_id"`
	Action       string    `json:"action"`    // "create", "update", "delete"
	Timestamp    time.Time `json:"timestamp"` // Timestamp of the change
	Source       string    `json:"source"`    // "online" or "offline"
}

// Sync constants for synchronization events
const (
	// Events for offline sync (online -> offline)
	EventBrandOfflineCreated = "brand.offline.created"
	EventBrandOfflineUpdated = "brand.offline.updated"
	EventBrandOfflineDeleted = "brand.offline.deleted"

	// Events for online sync (offline -> online)
	EventBrandOnlineCreated = "brand.online.created"
	EventBrandOnlineUpdated = "brand.online.updated"
	EventBrandOnlineDeleted = "brand.online.deleted"
)

// SyncEvent structures for synchronization
type SyncCreatedEvent struct{ events.BaseEvent }
type SyncUpdatedEvent struct{ events.BaseEvent }
type SyncDeletedEvent struct{ events.BaseEvent }

// ToSyncPayload converts Brand to SyncBrandPayload
func (b *Brand) ToSyncPayload(action, source string) SyncBrandPayload {
	return SyncBrandPayload{
		TenantSlug:   "",
		BrandID:     b.ID,
		Name:        b.Name,
		Description: b.Description,
		Active:      b.Active,
		EnterpriseID: b.EnterpriseID,
		Action:      action,
		Timestamp:   time.Now(),
		Source:      source,
	}
}

// ToSyncPayloadWithTenant converts Brand to SyncBrandPayload with tenant slug
func (b *Brand) ToSyncPayloadWithTenant(tenantSlug, action, source string) SyncBrandPayload {
	pl := b.ToSyncPayload(action, source)
	pl.TenantSlug = tenantSlug
	return pl
}

// NewSyncCreatedEvent creates a new sync created event
func NewSyncCreatedEvent(tenantSlug string, b *Brand) SyncCreatedEvent {
	return SyncCreatedEvent{events.NewBaseEvent(
		EventBrandOfflineCreated,
		b.ToSyncPayloadWithTenant(tenantSlug, "create", "online"),
	)}
}

// NewSyncUpdatedEvent creates a new sync updated event
func NewSyncUpdatedEvent(tenantSlug string, b *Brand) SyncUpdatedEvent {
	return SyncUpdatedEvent{events.NewBaseEvent(
		EventBrandOfflineUpdated,
		b.ToSyncPayloadWithTenant(tenantSlug, "update", "online"),
	)}
}

// NewSyncDeletedEvent creates a new sync deleted event
func NewSyncDeletedEvent(tenantSlug string, b *Brand) SyncDeletedEvent {
	return SyncDeletedEvent{events.NewBaseEvent(
		EventBrandOfflineDeleted,
		b.ToSyncPayloadWithTenant(tenantSlug, "delete", "online"),
	)}
}

// NewSyncCreatedEventFromOffline creates a new sync created event from offline source
func NewSyncCreatedEventFromOffline(tenantSlug string, b *Brand) SyncCreatedEvent {
	pl := b.ToSyncPayloadWithTenant(tenantSlug, "create", "offline")
	return SyncCreatedEvent{events.NewBaseEvent(EventBrandOnlineCreated, pl)}
}

// NewSyncUpdatedEventFromOffline creates a new sync updated event from offline source
func NewSyncUpdatedEventFromOffline(tenantSlug string, b *Brand) SyncUpdatedEvent {
	pl := b.ToSyncPayloadWithTenant(tenantSlug, "update", "offline")
	return SyncUpdatedEvent{events.NewBaseEvent(EventBrandOnlineUpdated, pl)}
}

// NewSyncDeletedEventFromOffline creates a new sync deleted event from offline source
func NewSyncDeletedEventFromOffline(tenantSlug string, b *Brand) SyncDeletedEvent {
	pl := b.ToSyncPayloadWithTenant(tenantSlug, "delete", "offline")
	return SyncDeletedEvent{events.NewBaseEvent(EventBrandOnlineDeleted, pl)}
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	BrandID int64  `json:"brand_id"`
	Name    string `json:"name"`
	Action  string `json:"action"` // "created", "updated", "deleted", "conflict"
	Status  string `json:"status"` // "success", "error", "conflict"
	Message string `json:"message"`
}