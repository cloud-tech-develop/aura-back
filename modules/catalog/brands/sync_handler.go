package brands

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

// SyncHandler handles brand synchronization events from RabbitMQ
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

// HandleWithContext handles an incoming brand sync event from RabbitMQ
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

	logger.Logf("[SyncHandler] Processing %s event for brand ID %d in tenant %s",
		action, int64(payload["brand_id"].(float64)), tenantSlug)

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
	b := BrandFromSyncPayload(payload)
	return h.repo.Create(ctx, tenantSlug, b)
}

func (h *SyncHandler) handleUpdate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	b := BrandFromSyncPayload(payload)
	return h.repo.Update(ctx, tenantSlug, b)
}

func (h *SyncHandler) handleDelete(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	brandID := int64(payload["brand_id"].(float64))
	return h.repo.Delete(ctx, tenantSlug, brandID)
}

// brandFromSyncPayload builds a Brand from a RabbitMQ sync payload map.
func BrandFromSyncPayload(payload map[string]interface{}) *Brand {
	return &Brand{
		ID:           Int64FromSyncPayload(payload, "brand_id"),
		Name:         StrFromSyncPayload(payload, "name"),
		Description:  StrFromSyncPayload(payload, "description"),
		Active:       BoolFromSyncPayload(payload, "active"),
		EnterpriseID: Int64FromSyncPayload(payload, "enterprise_id"),
	}
}

// BrandSyncFromOffline processes brands synced from offline to online
func (h *SyncHandler) BrandSyncFromOffline(ctx context.Context, tenantSlug string, brands []SyncBrandPayload) []SyncResult {
	logger := logging.NewLoggerHandler("sync")
	results := make([]SyncResult, 0, len(brands))

	for _, b := range brands {
		result := SyncResult{
			BrandID: b.BrandID,
			Name:    b.Name,
			Action:  b.Action,
		}

		switch b.Action {
		case "create":
			err := h.handleCreateFromOffline(ctx, tenantSlug, &b)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Brand created successfully"
			}
		case "update":
			err := h.handleUpdateFromOffline(ctx, tenantSlug, &b)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Brand updated successfully"
			}
		case "delete":
			err := h.handleDeleteFromOffline(ctx, tenantSlug, b.BrandID, b.Timestamp)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Brand deleted successfully"
			}
		default:
			result.Status = "error"
			result.Message = fmt.Sprintf("Unknown action: %s", b.Action)
		}

		logger.Logf("[SyncHandler] Sync result: %s - %s - %s",
			result.Name, result.Action, result.Message)
		results = append(results, result)
	}

	return results
}

func (h *SyncHandler) handleCreateFromOffline(ctx context.Context, tenantSlug string, p *SyncBrandPayload) error {
	brand := BrandFromSyncPayloadForCreate(p)

	existing, err := h.repo.GetByID(ctx, tenantSlug, p.BrandID)
	if err == nil {
		if isOnlineNewerBrand(existing.UpdatedAt, p.Timestamp) {
			return h.resolveConflict(ctx, tenantSlug, existing, brand, "create")
		}
		brand.CreatedAt = existing.CreatedAt
		return h.repo.Update(ctx, tenantSlug, brand)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking brand existence: %w", err)
	}

	return h.repo.Create(ctx, tenantSlug, brand)
}

func (h *SyncHandler) handleUpdateFromOffline(ctx context.Context, tenantSlug string, p *SyncBrandPayload) error {
	brand := BrandFromSyncPayloadForCreate(p)

	existing, err := h.repo.GetByID(ctx, tenantSlug, p.BrandID)
	if err == sql.ErrNoRows {
		return h.repo.Create(ctx, tenantSlug, brand)
	}
	if err != nil {
		return fmt.Errorf("error fetching brand: %w", err)
	}

	if isOnlineNewerBrand(existing.UpdatedAt, p.Timestamp) {
		return h.resolveConflict(ctx, tenantSlug, existing, brand, "update")
	}

	return h.repo.Update(ctx, tenantSlug, brand)
}

func (h *SyncHandler) handleDeleteFromOffline(ctx context.Context, tenantSlug string, brandID int64, offlineTimestamp time.Time) error {
	existing, err := h.repo.GetByID(ctx, tenantSlug, brandID)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error fetching brand: %w", err)
	}

	if isOnlineNewerBrand(existing.UpdatedAt, offlineTimestamp) {
		return fmt.Errorf("conflict: online has newer version, deletion rejected")
	}

	return h.repo.Delete(ctx, tenantSlug, brandID)
}

func (h *SyncHandler) resolveConflict(ctx context.Context, tenantSlug string, online, offline *Brand, action string) error {
	logger := logging.NewLoggerHandler("sync")
	logger.Logf("[SyncHandler] Conflict resolved in favor of online for brand ID %d (action: %s)",
		online.ID, action)
	return nil
}

func BrandFromSyncPayloadForCreate(p *SyncBrandPayload) *Brand {
	return &Brand{
		ID:           p.BrandID,
		Name:         p.Name,
		Description:  p.Description,
		Active:       p.Active,
		EnterpriseID: p.EnterpriseID,
	}
}

func convertToTimeBrand(dt *vo.DateTime) time.Time {
	if dt == nil {
		return time.Time{}
	}
	return time.Time(*dt)
}

func isOnlineNewerBrand(onlineTimestamp *vo.DateTime, offlineTimestamp time.Time) bool {
	if onlineTimestamp == nil {
		return false
	}
	return convertToTimeBrand(onlineTimestamp).After(offlineTimestamp)
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