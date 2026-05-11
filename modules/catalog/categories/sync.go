package categories

import (
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// SyncCategoryPayload contains all fields needed for synchronization between online and offline
type SyncCategoryPayload struct {
	TenantSlug      string    `json:"tenant_slug"`
	CategoryID      int64     `json:"category_id"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	DefaultTaxRate   float64   `json:"default_tax_rate"`
	Active          bool      `json:"active"`
	ParentID        *int64    `json:"parent_id,omitempty"`
	EnterpriseID    int64     `json:"enterprise_id"`
	Action          string    `json:"action"`    // "create", "update", "delete"
	Timestamp       time.Time `json:"timestamp"`  // Timestamp of the change
	Source          string    `json:"source"`     // "online" or "offline"
}

// Sync constants for synchronization events
const (
	// Events for offline sync (online -> offline)
	EventCategoryOfflineCreated = "category.offline.created"
	EventCategoryOfflineUpdated = "category.offline.updated"
	EventCategoryOfflineDeleted = "category.offline.deleted"

	// Events for online sync (offline -> online)
	EventCategoryOnlineCreated = "category.online.created"
	EventCategoryOnlineUpdated = "category.online.updated"
	EventCategoryOnlineDeleted = "category.online.deleted"
)

// SyncEvent structures for synchronization
type SyncCreatedEvent struct{ events.BaseEvent }
type SyncUpdatedEvent struct{ events.BaseEvent }
type SyncDeletedEvent struct{ events.BaseEvent }

// ToSyncPayload converts Category to SyncCategoryPayload
func (c *Category) ToSyncPayload(action, source string) SyncCategoryPayload {
	return SyncCategoryPayload{
		TenantSlug:    "",
		CategoryID:   c.ID,
		Name:         c.Name,
		Description:  c.Description,
		DefaultTaxRate: c.DefaultTaxRate,
		Active:       c.Active,
		ParentID:     c.ParentID,
		EnterpriseID: c.EnterpriseID,
		Action:       action,
		Timestamp:    time.Now(),
		Source:       source,
	}
}

// ToSyncPayloadWithTenant converts Category to SyncCategoryPayload with tenant slug
func (c *Category) ToSyncPayloadWithTenant(tenantSlug, action, source string) SyncCategoryPayload {
	pl := c.ToSyncPayload(action, source)
	pl.TenantSlug = tenantSlug
	return pl
}

// NewSyncCreatedEvent creates a new sync created event
func NewSyncCreatedEvent(tenantSlug string, c *Category) SyncCreatedEvent {
	return SyncCreatedEvent{events.NewBaseEvent(
		EventCategoryOfflineCreated,
		c.ToSyncPayloadWithTenant(tenantSlug, "create", "online"),
	)}
}

// NewSyncUpdatedEvent creates a new sync updated event
func NewSyncUpdatedEvent(tenantSlug string, c *Category) SyncUpdatedEvent {
	return SyncUpdatedEvent{events.NewBaseEvent(
		EventCategoryOfflineUpdated,
		c.ToSyncPayloadWithTenant(tenantSlug, "update", "online"),
	)}
}

// NewSyncDeletedEvent creates a new sync deleted event
func NewSyncDeletedEvent(tenantSlug string, c *Category) SyncDeletedEvent {
	return SyncDeletedEvent{events.NewBaseEvent(
		EventCategoryOfflineDeleted,
		c.ToSyncPayloadWithTenant(tenantSlug, "delete", "online"),
	)}
}

// NewSyncCreatedEventFromOffline creates a new sync created event from offline source
func NewSyncCreatedEventFromOffline(tenantSlug string, c *Category) SyncCreatedEvent {
	pl := c.ToSyncPayloadWithTenant(tenantSlug, "create", "offline")
	return SyncCreatedEvent{events.NewBaseEvent(EventCategoryOnlineCreated, pl)}
}

// NewSyncUpdatedEventFromOffline creates a new sync updated event from offline source
func NewSyncUpdatedEventFromOffline(tenantSlug string, c *Category) SyncUpdatedEvent {
	pl := c.ToSyncPayloadWithTenant(tenantSlug, "update", "offline")
	return SyncUpdatedEvent{events.NewBaseEvent(EventCategoryOnlineUpdated, pl)}
}

// NewSyncDeletedEventFromOffline creates a new sync deleted event from offline source
func NewSyncDeletedEventFromOffline(tenantSlug string, c *Category) SyncDeletedEvent {
	pl := c.ToSyncPayloadWithTenant(tenantSlug, "delete", "offline")
	return SyncDeletedEvent{events.NewBaseEvent(EventCategoryOnlineDeleted, pl)}
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	Action     string `json:"action"` // "created", "updated", "deleted", "conflict"
	Status     string `json:"status"` // "success", "error", "conflict"
	Message    string `json:"message"`
}