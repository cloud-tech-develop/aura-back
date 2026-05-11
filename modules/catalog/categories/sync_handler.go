package categories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/cloud-tech-develop/aura-back/shared/logging"
)

// Compile-time check: ensure SyncHandler implements events.EventHandler.
var _ events.EventHandler = (*SyncHandler)(nil)

// SyncHandler handles category synchronization events from RabbitMQ
type SyncHandler struct {
	repo     Repository
	eventBus events.EventBus
}

// NewSyncHandler creates a new sync handler instance
func NewSyncHandler(repo Repository, eventBus events.EventBus) *SyncHandler {
	return &SyncHandler{
		repo:     repo,
		eventBus: eventBus,
	}
}

// Handle implements events.EventHandler interface
func (h *SyncHandler) Handle(event events.Event) error {
	return h.HandleWithContext(context.Background(), event)
}

// HandleWithContext handles an incoming category sync event from RabbitMQ
func (h *SyncHandler) HandleWithContext(ctx context.Context, event events.Event) error {
	logger := logging.NewLoggerHandler("sync")

	payload, ok := event.GetPayload().(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid sync payload type")
	}

	tenantSlug, ok := payload["tenant_slug"].(string)
	if !ok || tenantSlug == "" {
		return fmt.Errorf("tenant_slug is required for sync")
	}

	action, ok := payload["action"].(string)
	if !ok {
		return fmt.Errorf("action is required for sync")
	}

	logger.Logf("[SyncHandler] Processing %s event for category ID %d in tenant %s",
		action, int64(payload["category_id"].(float64)), tenantSlug)

	switch action {
	case "create":
		return h.handleCreate(ctx, tenantSlug, payload)
	case "update":
		return h.handleUpdate(ctx, tenantSlug, payload)
	case "delete":
		return h.handleDelete(ctx, tenantSlug, payload)
	default:
		logger.Logf("[SyncHandler] Unknown action: %s", action)
		return fmt.Errorf("unknown sync action: %s", action)
	}
}

// EventHandlerFunc returns an events.EventHandler function for use with RabbitMQ subscriptions
func (h *SyncHandler) EventHandlerFunc() events.EventHandler {
	return h
}

func (h *SyncHandler) handleCreate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	c := CategoryFromSyncPayload(payload)
	return h.repo.Create(ctx, tenantSlug, c)
}

func (h *SyncHandler) handleUpdate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	c := CategoryFromSyncPayload(payload)
	return h.repo.Update(ctx, tenantSlug, c)
}

func (h *SyncHandler) handleDelete(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	categoryID := int64(payload["category_id"].(float64))
	return h.repo.Delete(ctx, tenantSlug, categoryID)
}

// categoryFromSyncPayload builds a Category from a RabbitMQ sync payload map.
func CategoryFromSyncPayload(payload map[string]interface{}) *Category {
	c := &Category{
		ID:             Int64FromSyncPayload(payload, "category_id"),
		Name:           StrFromSyncPayload(payload, "name"),
		DefaultTaxRate: FloatFromSyncPayload(payload, "default_tax_rate"),
		Active:         BoolFromSyncPayload(payload, "active"),
		EnterpriseID:   Int64FromSyncPayload(payload, "enterprise_id"),
	}
	if v, ok := payload["description"]; ok && v != nil {
		desc := StrFromSyncPayload(payload, "description")
		c.Description = &desc
	}
	if v, ok := payload["parent_id"]; ok && v != nil {
		id := Int64FromSyncPayload(payload, "parent_id")
		c.ParentID = &id
	}
	return c
}

// CategorySyncFromOffline processes categories synced from offline to online
func (h *SyncHandler) CategorySyncFromOffline(ctx context.Context, tenantSlug string, categories []SyncCategoryPayload) []SyncResult {
	logger := logging.NewLoggerHandler("sync")
	results := make([]SyncResult, 0, len(categories))

	for _, c := range categories {
		result := SyncResult{
			CategoryID: c.CategoryID,
			Name:       c.Name,
			Action:     c.Action,
		}

		switch c.Action {
		case "create":
			err := h.handleCreateFromOffline(ctx, tenantSlug, &c)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Category created successfully"
			}
		case "update":
			err := h.handleUpdateFromOffline(ctx, tenantSlug, &c)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Category updated successfully"
			}
		case "delete":
			err := h.handleDeleteFromOffline(ctx, tenantSlug, c.CategoryID, c.Timestamp)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Category deleted successfully"
			}
		default:
			result.Status = "error"
			result.Message = fmt.Sprintf("Unknown action: %s", c.Action)
		}

		logger.Logf("[SyncHandler] Sync result: %s - %s - %s",
			result.Name, result.Action, result.Message)
		results = append(results, result)
	}

	return results
}

func (h *SyncHandler) handleCreateFromOffline(ctx context.Context, tenantSlug string, p *SyncCategoryPayload) error {
	category := CategoryFromSyncPayloadForCreate(p)

	existing, err := h.repo.GetByID(ctx, tenantSlug, p.CategoryID)
	if err == nil {
		if isOnlineNewerCategory(existing.UpdatedAt, p.Timestamp) {
			return h.resolveConflict(ctx, tenantSlug, existing, category, "create")
		}
		category.CreatedAt = existing.CreatedAt
		return h.repo.Update(ctx, tenantSlug, category)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking category existence: %w", err)
	}

	return h.repo.Create(ctx, tenantSlug, category)
}

func (h *SyncHandler) handleUpdateFromOffline(ctx context.Context, tenantSlug string, p *SyncCategoryPayload) error {
	category := CategoryFromSyncPayloadForCreate(p)

	existing, err := h.repo.GetByID(ctx, tenantSlug, p.CategoryID)
	if err == sql.ErrNoRows {
		return h.repo.Create(ctx, tenantSlug, category)
	}
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}

	if isOnlineNewerCategory(existing.UpdatedAt, p.Timestamp) {
		return h.resolveConflict(ctx, tenantSlug, existing, category, "update")
	}

	return h.repo.Update(ctx, tenantSlug, category)
}

func (h *SyncHandler) handleDeleteFromOffline(ctx context.Context, tenantSlug string, categoryID int64, offlineTimestamp time.Time) error {
	existing, err := h.repo.GetByID(ctx, tenantSlug, categoryID)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}

	if isOnlineNewerCategory(existing.UpdatedAt, offlineTimestamp) {
		return fmt.Errorf("conflict: online has newer version, deletion rejected")
	}

	return h.repo.Delete(ctx, tenantSlug, categoryID)
}

func (h *SyncHandler) resolveConflict(ctx context.Context, tenantSlug string, online, offline *Category, action string) error {
	logger := logging.NewLoggerHandler("sync")
	logger.Logf("[SyncHandler] Conflict resolved in favor of online for category ID %d (action: %s)",
		online.ID, action)
	return nil
}

func CategoryFromSyncPayloadForCreate(p *SyncCategoryPayload) *Category {
	c := &Category{
		ID:              p.CategoryID,
		Name:            p.Name,
		Description:     p.Description,
		DefaultTaxRate:  p.DefaultTaxRate,
		Active:          p.Active,
		ParentID:        p.ParentID,
		EnterpriseID:    p.EnterpriseID,
	}
	return c
}

func convertToTimeCategory(dt *vo.DateTime) time.Time {
	if dt == nil {
		return time.Time{}
	}
	return time.Time(*dt)
}

func isOnlineNewerCategory(onlineTimestamp *vo.DateTime, offlineTimestamp time.Time) bool {
	if onlineTimestamp == nil {
		return false
	}
	return convertToTimeCategory(onlineTimestamp).After(offlineTimestamp)
}

// ─── payload helpers ─────────────────────────────────────────────────────────

func StrFromSyncPayload(m map[string]interface{}, k string) string {
	v, _ := m[k].(string)
	return v
}

func FloatFromSyncPayload(m map[string]interface{}, k string) float64 {
	v, _ := m[k].(float64)
	return v
}

func Int64FromSyncPayload(m map[string]interface{}, k string) int64 {
	v, _ := m[k].(float64)
	return int64(v)
}

func BoolFromSyncPayload(m map[string]interface{}, k string) bool {
	v, _ := m[k].(bool)
	return v
}