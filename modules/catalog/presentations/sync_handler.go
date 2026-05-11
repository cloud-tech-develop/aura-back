package presentations

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

// SyncHandler handles presentation synchronization events from RabbitMQ
type SyncHandler struct {
	repo Repository
}

// NewSyncHandler creates a new sync handler instance
func NewSyncHandler(repo Repository) *SyncHandler {
	return &SyncHandler{
		repo: repo,
	}
}

// Handle implements events.EventHandler interface
func (h *SyncHandler) Handle(event events.Event) error {
	ctx := context.Background()
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

	logger.Logf("[SyncHandler] Processing %s event for presentation ID %d in tenant %s",
		action, int64(payload["presentation_id"].(float64)), tenantSlug)

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
	p := presentationFromSyncPayload(payload)
	return h.repo.Create(ctx, tenantSlug, p.EnterpriseID, p)
}

func (h *SyncHandler) handleUpdate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	p := presentationFromSyncPayload(payload)
	return h.repo.Update(ctx, tenantSlug, p)
}

func (h *SyncHandler) handleDelete(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	presentationID := int64(payload["presentation_id"].(float64))
	return h.repo.Delete(ctx, tenantSlug, presentationID)
}

// presentationFromSyncPayload builds a Presentation from a RabbitMQ sync payload map.
func presentationFromSyncPayload(payload map[string]interface{}) *Presentation {
	return &Presentation{
		ID:              int64AsFloat(payload, "presentation_id"),
		ProductID:       int64AsFloat(payload, "product_id"),
		Name:            strVal(payload, "name"),
		Factor:          floatVal(payload, "factor"),
		Barcode:         strVal(payload, "barcode"),
		CostPrice:       floatVal(payload, "cost_price"),
		SalePrice:       floatVal(payload, "sale_price"),
		DefaultPurchase: boolVal(payload, "default_purchase"),
		DefaultSale:     boolVal(payload, "default_sale"),
		EnterpriseID:    int64AsFloat(payload, "enterprise_id"),
	}
}

// PresentationSyncFromOffline processes presentations synced from offline to online
func (h *SyncHandler) PresentationSyncFromOffline(ctx context.Context, tenantSlug string, presentations []SyncPresentationPayload) []SyncResult {
	logger := logging.NewLoggerHandler("sync")
	results := make([]SyncResult, 0, len(presentations))

	for _, p := range presentations {
		result := SyncResult{
			PresentationID: p.PresentationID,
			Name:          p.Name,
			Action:        p.Action,
		}

		switch p.Action {
		case "create":
			err := h.handleCreateFromOffline(ctx, tenantSlug, &p)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Presentation created successfully"
			}
		case "update":
			err := h.handleUpdateFromOffline(ctx, tenantSlug, &p)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Presentation updated successfully"
			}
		case "delete":
			err := h.handleDeleteFromOffline(ctx, tenantSlug, p.PresentationID, p.Timestamp)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Presentation deleted successfully"
			}
		default:
			result.Status = "error"
			result.Message = fmt.Sprintf("Unknown action: %s", p.Action)
		}

		logger.Logf("[SyncHandler] Sync result: %s - %s - %s",
			result.Name, result.Action, result.Message)
		results = append(results, result)
	}

	return results
}

func (h *SyncHandler) handleCreateFromOffline(ctx context.Context, tenantSlug string, p *SyncPresentationPayload) error {
	presentation := &Presentation{
		ID:              p.PresentationID,
		ProductID:       p.ProductID,
		Name:            p.Name,
		Factor:          p.Factor,
		Barcode:         p.Barcode,
		CostPrice:       p.CostPrice,
		SalePrice:       p.SalePrice,
		DefaultPurchase: p.DefaultPurchase,
		DefaultSale:     p.DefaultSale,
		EnterpriseID:    p.EnterpriseID,
	}

	// Check if presentation already exists
	existing, err := h.repo.GetByID(ctx, tenantSlug, p.PresentationID)
	if err == nil {
		// Presentation exists, check for conflict based on timestamp
		if isOnlineNewer(existing.UpdatedAt, p.Timestamp) {
			return h.resolveConflict(ctx, tenantSlug, existing, presentation, "create")
		}
		// Offline version is newer or equal, update it
		presentation.CreatedAt = existing.CreatedAt
		return h.repo.Update(ctx, tenantSlug, presentation)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking presentation existence: %w", err)
	}

	// New presentation, create it
	return h.repo.Create(ctx, tenantSlug, p.EnterpriseID, presentation)
}

func (h *SyncHandler) handleUpdateFromOffline(ctx context.Context, tenantSlug string, p *SyncPresentationPayload) error {
	presentation := &Presentation{
		ID:              p.PresentationID,
		ProductID:       p.ProductID,
		Name:            p.Name,
		Factor:          p.Factor,
		Barcode:         p.Barcode,
		CostPrice:       p.CostPrice,
		SalePrice:       p.SalePrice,
		DefaultPurchase: p.DefaultPurchase,
		DefaultSale:     p.DefaultSale,
		EnterpriseID:    p.EnterpriseID,
	}

	// Check if presentation exists
	existing, err := h.repo.GetByID(ctx, tenantSlug, p.PresentationID)
	if err == sql.ErrNoRows {
		// Presentation doesn't exist in online, create it
		return h.repo.Create(ctx, tenantSlug, p.EnterpriseID, presentation)
	}
	if err != nil {
		return fmt.Errorf("error fetching presentation: %w", err)
	}

	// Check for conflict based on timestamp
	if isOnlineNewer(existing.UpdatedAt, p.Timestamp) {
		return h.resolveConflict(ctx, tenantSlug, existing, presentation, "update")
	}

	// Offline version is newer or equal, update it
	return h.repo.Update(ctx, tenantSlug, presentation)
}

func (h *SyncHandler) handleDeleteFromOffline(ctx context.Context, tenantSlug string, presentationID int64, offlineTimestamp time.Time) error {
	existing, err := h.repo.GetByID(ctx, tenantSlug, presentationID)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error fetching presentation: %w", err)
	}

	if isOnlineNewer(existing.UpdatedAt, offlineTimestamp) {
		return fmt.Errorf("conflict: online has newer version, deletion rejected")
	}

	return h.repo.Delete(ctx, tenantSlug, presentationID)
}

func (h *SyncHandler) resolveConflict(ctx context.Context, tenantSlug string, online, offline *Presentation, action string) error {
	logger := logging.NewLoggerHandler("sync")
	logger.Logf("[SyncHandler] Conflict resolved in favor of online for presentation ID %d (action: %s)",
		online.ID, action)
	return nil
}

func convertToTime(dt *vo.DateTime) time.Time {
	if dt == nil {
		return time.Time{}
	}
	return time.Time(*dt)
}

func isOnlineNewer(onlineTimestamp *vo.DateTime, offlineTimestamp time.Time) bool {
	if onlineTimestamp == nil {
		return false
	}
	return convertToTime(onlineTimestamp).After(offlineTimestamp)
}

func strVal(m map[string]interface{}, k string) string {
	v, _ := m[k].(string)
	return v
}

func floatVal(m map[string]interface{}, k string) float64 {
	v, _ := m[k].(float64)
	return v
}

func int64AsFloat(m map[string]interface{}, k string) int64 {
	v, _ := m[k].(float64)
	return int64(v)
}

func boolVal(m map[string]interface{}, k string) bool {
	v, _ := m[k].(bool)
	return v
}