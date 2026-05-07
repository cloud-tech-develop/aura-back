package products

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

// SyncHandler handles product synchronization events from RabbitMQ
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

// HandleProductSync handles an incoming product sync event from RabbitMQ
// This is called when online publishes an event that offline subscribes to
func (h *SyncHandler) HandleProductSync(ctx context.Context, event events.Event) error {
	logger := logging.NewLoggerHandler("sync")

	payload, ok := event.GetPayload().(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid sync payload type")
	}

	// Extract tenant slug
	tenantSlug, ok := payload["tenant_slug"].(string)
	if !ok || tenantSlug == "" {
		return fmt.Errorf("tenant_slug is required for sync")
	}

	action, ok := payload["action"].(string)
	if !ok {
		return fmt.Errorf("action is required for sync")
	}

	logger.Logf("[SyncHandler] Processing %s event for product ID %d in tenant %s",
		action, int64(payload["product_id"].(float64)), tenantSlug)

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

func (h *SyncHandler) handleCreate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	product := productFromSyncPayload(payload)
	return h.repo.Create(ctx, tenantSlug, product)
}

func (h *SyncHandler) handleUpdate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	product := productFromSyncPayload(payload)
	return h.repo.Update(ctx, tenantSlug, product)
}

func (h *SyncHandler) handleDelete(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	productID := int64(payload["product_id"].(float64))
	return h.repo.Delete(ctx, tenantSlug, productID)
}

// productFromSyncPayload builds a Product from a RabbitMQ sync payload map.
// Uses the same safe helpers as service.productFromPayload.
func productFromSyncPayload(payload map[string]interface{}) *Product {
	p := &Product{
		ID:                 int64AsFloat(payload, "product_id"),
		SKU:                strVal(payload, "sku"),
		Barcode:            strVal(payload, "barcode"),
		Name:               strVal(payload, "name"),
		Description:        strVal(payload, "description"),
		ProductType:        strVal(payload, "product_type"),
		Active:             boolVal(payload, "active"),
		VisibleInPOS:       boolVal(payload, "visible_in_pos"),
		CostPrice:          floatVal(payload, "cost_price"),
		SalePrice:          floatVal(payload, "sale_price"),
		Price2:             floatVal(payload, "price_2"),
		IVAPercentage:      floatVal(payload, "iva_percentage"),
		ConsumptionTax:     floatVal(payload, "consumption_tax"),
		ManagesInventory:   boolVal(payload, "manages_inventory"),
		ManagesBatches:     boolVal(payload, "manages_batches"),
		ManagesSerial:      boolVal(payload, "manages_serial"),
		AllowNegativeStock: boolVal(payload, "allow_negative_stock"),
		ImageURL:           strVal(payload, "image_url"),
		EnterpriseID:       int64AsFloat(payload, "enterprise_id"),
		UnitID:             int64AsFloat(payload, "unit_measure_id"),
		MinStock:           intVal(payload, "min_stock"),
		MaxStock:           intVal(payload, "max_stock"),
		CurrentStock:       intVal(payload, "current_stock"),
	}
	if v, ok := payload["category_id"]; ok && v != nil {
		id := int64AsFloat(payload, "category_id")
		p.CategoryID = &id
	}
	if v, ok := payload["brand_id"]; ok && v != nil {
		id := int64AsFloat(payload, "brand_id")
		p.BrandID = &id
	}
	if v, ok := payload["price_3"]; ok && v != nil {
		p3 := floatVal(payload, "price_3")
		p.Price3 = &p3
	}
	return p
}

// Handle implements events.EventHandler interface
func (h *SyncHandler) Handle(event events.Event) error {
	return h.HandleProductSync(context.Background(), event)
}

// SyncHandlerFunc returns an events.EventHandler function for use with RabbitMQ subscriptions
func (h *SyncHandler) SyncHandlerFunc() events.EventHandler {
	return h
}

// ProductSyncFromOffline processes products synced from offline to online
// Returns sync results for each product
func (h *SyncHandler) ProductSyncFromOffline(ctx context.Context, tenantSlug string, products []SyncProductPayload) []SyncResult {
	logger := logging.NewLoggerHandler("sync")
	results := make([]SyncResult, 0, len(products))

	for _, p := range products {
		result := SyncResult{
			ProductID: p.ProductID,
			SKU:       p.SKU,
			Action:    p.Action,
		}

		switch p.Action {
		case "create":
			err := h.handleCreateFromOffline(ctx, tenantSlug, &p)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Product created successfully"
			}
		case "update":
			err := h.handleUpdateFromOffline(ctx, tenantSlug, &p)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Product updated successfully"
			}
		case "delete":
			err := h.handleDeleteFromOffline(ctx, tenantSlug, p.ProductID, p.Timestamp)
			if err != nil {
				result.Status = "error"
				result.Message = err.Error()
			} else {
				result.Status = "success"
				result.Message = "Product deleted successfully"
			}
		default:
			result.Status = "error"
			result.Message = fmt.Sprintf("Unknown action: %s", p.Action)
		}

		logger.Logf("[SyncHandler] Sync result: %s - %s - %s",
			result.SKU, result.Action, result.Message)
		results = append(results, result)
	}

	return results
}

func (h *SyncHandler) handleCreateFromOffline(ctx context.Context, tenantSlug string, p *SyncProductPayload) error {
	product := &Product{
		ID:                 p.ProductID,
		SKU:                p.SKU,
		Barcode:            p.Barcode,
		Name:               p.Name,
		Description:        p.Description,
		ProductType:        p.ProductType,
		Active:             p.Active,
		VisibleInPOS:       p.VisibleInPOS,
		CostPrice:          p.CostPrice,
		SalePrice:          p.SalePrice,
		Price2:             p.Price2,
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
		CategoryID:         p.CategoryID,
		BrandID:            p.BrandID,
		UnitID:             p.UnitID,
		Price3:             p.Price3,
	}

	// Check if product already exists
	existing, err := h.repo.GetByID(ctx, tenantSlug, p.ProductID)
	if err == nil {
		// Product exists, check for conflict based on timestamp
		if isOnlineNewer(existing.UpdatedAt, p.Timestamp) {
			// Online has newer version, conflict detected
			return h.resolveConflict(ctx, tenantSlug, existing, product, "create")
		}
		// Offline version is newer or equal, update it
		product.CreatedAt = existing.CreatedAt
		return h.repo.Update(ctx, tenantSlug, product)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking product existence: %w", err)
	}

	// New product, create it
	return h.repo.Create(ctx, tenantSlug, product)
}

func (h *SyncHandler) handleUpdateFromOffline(ctx context.Context, tenantSlug string, p *SyncProductPayload) error {
	product := &Product{
		ID:                 p.ProductID,
		SKU:                p.SKU,
		Barcode:            p.Barcode,
		Name:               p.Name,
		Description:        p.Description,
		ProductType:        p.ProductType,
		Active:             p.Active,
		VisibleInPOS:       p.VisibleInPOS,
		CostPrice:          p.CostPrice,
		SalePrice:          p.SalePrice,
		Price2:             p.Price2,
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
		CategoryID:         p.CategoryID,
		BrandID:            p.BrandID,
		UnitID:             p.UnitID,
		Price3:             p.Price3,
	}

	// Check if product exists
	existing, err := h.repo.GetByID(ctx, tenantSlug, p.ProductID)
	if err == sql.ErrNoRows {
		// Product doesn't exist in online, create it
		return h.repo.Create(ctx, tenantSlug, product)
	}
	if err != nil {
		return fmt.Errorf("error fetching product: %w", err)
	}

	// Check for conflict based on timestamp
	if isOnlineNewer(existing.UpdatedAt, p.Timestamp) {
		// Online has newer version, conflict detected
		return h.resolveConflict(ctx, tenantSlug, existing, product, "update")
	}

	// Offline version is newer or equal, update it
	return h.repo.Update(ctx, tenantSlug, product)
}

func (h *SyncHandler) handleDeleteFromOffline(ctx context.Context, tenantSlug string, productID int64, offlineTimestamp time.Time) error {
	// Check if product exists
	existing, err := h.repo.GetByID(ctx, tenantSlug, productID)
	if err == sql.ErrNoRows {
		// Product already deleted in online
		return nil
	}
	if err != nil {
		return fmt.Errorf("error fetching product: %w", err)
	}

	// Check for conflict based on timestamp
	if isOnlineNewer(existing.UpdatedAt, offlineTimestamp) {
		// Online has newer version, don't delete
		return fmt.Errorf("conflict: online has newer version, deletion rejected")
	}

	// Safe to delete
	return h.repo.Delete(ctx, tenantSlug, productID)
}

// resolveConflict resolves a sync conflict based on timestamps
// Currently uses the "last write wins" strategy based on updated_at
// TODO: Add sync_conflicts table for audit logging
func (h *SyncHandler) resolveConflict(ctx context.Context, tenantSlug string, online, offline *Product, action string) error {
	logger := logging.NewLoggerHandler("sync")

	// Strategy: Last write wins (online wins if online timestamp > offline timestamp)
	// This is already handled by checking timestamps in the caller
	// This method is for future audit logging

	logger.Logf("[SyncHandler] Conflict resolved in favor of online for product ID %d (action: %s)",
		online.ID, action)

	// TODO: Insert into sync_conflicts table for audit
	// This requires creating the table first via migration

	return nil
}

// InsertSyncConflict inserts a conflict record for audit purposes
// This is a placeholder for future implementation
func (h *SyncHandler) InsertSyncConflict(ctx context.Context, tenantSlug string, online, offline *Product, action, resolution string) error {
	// TODO: Implement when sync_conflicts table is created
	return nil
}

// convertToTime converts *vo.DateTime to time.Time for comparison
func convertToTime(dt *vo.DateTime) time.Time {
	if dt == nil {
		return time.Time{}
	}
	return time.Time(*dt)
}

// isOnlineNewer checks if online product has a newer timestamp than offline
func isOnlineNewer(onlineTimestamp *vo.DateTime, offlineTimestamp time.Time) bool {
	if onlineTimestamp == nil {
		return false
	}
	return convertToTime(onlineTimestamp).After(offlineTimestamp)
}
