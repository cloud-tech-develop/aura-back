package presentations

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/cloud-tech-develop/aura-back/shared/logging"
)

// service implements the Service interface
// Contains business logic for presentation management
type service struct {
	repo      Repository
	eventBus  events.EventBus
	syncBus   events.EventBus // Cross-server sync bus (RabbitMQ)
	syncMu    sync.RWMutex
	logger    *logging.LoggerHandler
	isOffline bool
}

// NewService creates a new presentation service instance
func NewService(db *db.DB, eventBus events.EventBus) Service {
	logger := logging.NewLoggerHandler("sync")
	logger.Log("[Presentation Service] Initializing service")

	svc := &service{
		repo:      NewRepository(db),
		eventBus:  eventBus,
		logger:    logger,
		isOffline: db.IsSQLite(),
	}

	svc.subscribeToRabbitMQEvents()

	return svc
}

// subscribeToRabbitMQEvents registers cross-server sync subscriptions.
func (s *service) subscribeToRabbitMQEvents() {
	if s.eventBus == nil {
		return
	}

	if s.isOffline {
		s.logger.Log("[Presentation Service] Offline mode: RabbitMQ sync deferred until /offline/ping")
		return
	}

	// Online: wildcard binding receives offline events from all tenants.
	s.eventBus.Subscribe(EventPresentationOnlineCreated+".*", s)
	s.eventBus.Subscribe(EventPresentationOnlineUpdated+".*", s)
	s.eventBus.Subscribe(EventPresentationOnlineDeleted+".*", s)
	s.logger.Log("[Presentation Service] Online mode: Subscribed to offline presentation events (wildcard)")
}

// SetSyncBus sets the RabbitMQ event bus for cross-server sync publishing.
func (s *service) SetSyncBus(bus events.EventBus) {
	s.syncMu.Lock()
	defer s.syncMu.Unlock()
	s.syncBus = bus
	s.logger.Log("[Presentation Service] Sync bus updated to RabbitMQ")
}

// Handle implements events.EventHandler for RabbitMQ events
func (s *service) Handle(event events.Event) error {
	s.logger.Logf("[Presentation Service] Received event from RabbitMQ: %s", event.GetName())

	payload, ok := event.GetPayload().(map[string]interface{})
	if !ok {
		s.logger.Log("[Presentation Service] Invalid payload type in event")
		return fmt.Errorf("invalid payload type")
	}

	tenantSlug, ok := payload["tenant_slug"].(string)
	if !ok || tenantSlug == "" {
		s.logger.Log("[Presentation Service] No tenant_slug in event, skipping")
		return fmt.Errorf("tenant_slug is required")
	}

	s.logger.Logf("[Presentation Service] Processing %s event for tenant: %s", event.GetName(), tenantSlug)

	ctx := context.Background()

	switch event.GetName() {
	case EventPresentationOfflineCreated, EventPresentationOnlineCreated:
		return s.handleRemoteCreate(ctx, tenantSlug, payload)
	case EventPresentationOfflineUpdated, EventPresentationOnlineUpdated:
		return s.handleRemoteUpdate(ctx, tenantSlug, payload)
	case EventPresentationOfflineDeleted, EventPresentationOnlineDeleted:
		return s.handleRemoteDelete(ctx, tenantSlug, payload)
	}

	return nil
}

func (s *service) handleRemoteCreate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	p := presentationFromPayload(payload)
	s.logger.Logf("[Presentation Service] Remote create: ID %d", p.ID)

	existing, err := s.repo.GetByID(ctx, tenantSlug, p.ID)
	if err == nil {
		var eventTime time.Time
		if tStr, ok := payload["timestamp"].(string); ok {
			eventTime, _ = time.Parse(time.RFC3339, tStr)
		}

		if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
			s.logger.Logf("[Presentation Service] Conflict: local version is newer for presentation ID %d, skipping", p.ID)
			return nil
		}

		p.CreatedAt = existing.CreatedAt
		if err := s.repo.Update(ctx, tenantSlug, p); err != nil {
			return fmt.Errorf("failed to update presentation: %w", err)
		}
		s.logger.Logf("[Presentation Service] Presentation updated from remote (resolved conflict): ID %d", p.ID)
		return nil
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking presentation existence: %w", err)
	}

	if err := s.repo.Create(ctx, tenantSlug, p.EnterpriseID, p); err != nil {
		return fmt.Errorf("failed to create presentation: %w", err)
	}
	s.logger.Logf("[Presentation Service] Presentation created from remote: ID %d", p.ID)
	return nil
}

func (s *service) handleRemoteUpdate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	p := presentationFromPayload(payload)
	s.logger.Logf("[Presentation Service] Remote update: ID %d", p.ID)

	existing, err := s.repo.GetByID(ctx, tenantSlug, p.ID)
	if err == sql.ErrNoRows {
		if err := s.repo.Create(ctx, tenantSlug, p.EnterpriseID, p); err != nil {
			return fmt.Errorf("failed to create presentation on remote update: %w", err)
		}
		s.logger.Logf("[Presentation Service] Presentation created from remote update: ID %d", p.ID)
		return nil
	} else if err != nil {
		return fmt.Errorf("error fetching presentation: %w", err)
	}

	var eventTime time.Time
	if tStr, ok := payload["timestamp"].(string); ok {
		eventTime, _ = time.Parse(time.RFC3339, tStr)
	}

	if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
		s.logger.Logf("[Presentation Service] Conflict: local version is newer for ID %d, skipping update", p.ID)
		return nil
	}

	if err := s.repo.Update(ctx, tenantSlug, p); err != nil {
		return fmt.Errorf("failed to update presentation: %w", err)
	}
	s.logger.Logf("[Presentation Service] Presentation updated from remote: ID %d", p.ID)
	return nil
}

func (s *service) handleRemoteDelete(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	presentationID := int64(payload["presentation_id"].(float64))
	s.logger.Logf("[Presentation Service] Remote delete: ID %d", presentationID)

	existing, err := s.repo.GetByID(ctx, tenantSlug, presentationID)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("error fetching presentation: %w", err)
	}

	var eventTime time.Time
	if tStr, ok := payload["timestamp"].(string); ok {
		eventTime, _ = time.Parse(time.RFC3339, tStr)
	}

	if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
		s.logger.Logf("[Presentation Service] Conflict: local version is newer for ID %d, skipping delete", presentationID)
		return nil
	}

	if err := s.repo.Delete(ctx, tenantSlug, presentationID); err != nil {
		return fmt.Errorf("failed to delete presentation: %w", err)
	}
	s.logger.Logf("[Presentation Service] Presentation deleted from remote: ID %d", presentationID)
	return nil
}

// Create creates multiple presentations for a product
func (s *service) Create(ctx context.Context, tenantSlug string, enterpriseID int64, productID int64, presentations []PresentationRequest) error {
	logger := logging.NewLoggerHandler("logs")
	logger.Log(fmt.Sprintf("[Presentation Service] Creating %d presentations for product ID: %d", len(presentations), productID))

	// Validate product exists first by trying to get it
	if productID == 0 {
		logger.Log("[Presentation Service] Validation failed: product_id is required")
		return fmt.Errorf("product_id is required")
	}

	// Validate enterprise ID
	if enterpriseID == 0 {
		logger.Log("[Presentation Service] Validation failed: enterprise_id is required")
		return fmt.Errorf("enterprise_id is required")
	}

	// Validate we have at least one presentation
	if len(presentations) == 0 {
		logger.Log("[Presentation Service] Validation failed: at least one presentation is required")
		return fmt.Errorf("at least one presentation is required")
	}

	// Convert request to entities
	entities := make([]*Presentation, len(presentations))
	for i, req := range presentations {
		logger.Logf("[Presentation Service] Validating presentation %d: name=%s", i+1, req.Name)

		// Validate required fields
		if req.Name == "" {
			logger.Logf("[Presentation Service] Validation failed: name is required for presentation %d", i+1)
			return fmt.Errorf("name is required for presentation %d", i+1)
		}
		if req.Factor == 0 {
			logger.Logf("[Presentation Service] Validation failed: factor is required for presentation %d", i+1)
			return fmt.Errorf("factor is required for presentation %d", i+1)
		}

		entities[i] = &Presentation{
			ProductID:       productID,
			Name:            req.Name,
			Factor:          req.Factor,
			Barcode:         req.Barcode,
			CostPrice:       req.CostPrice,
			SalePrice:       req.SalePrice,
			DefaultPurchase: req.DefaultPurchase,
			DefaultSale:     req.DefaultSale,
			EnterpriseID:    enterpriseID,
		}
	}

	// Create all presentations
	err := s.repo.CreateMany(ctx, tenantSlug, enterpriseID, entities)
	if err != nil {
		logger.Logf("[Presentation Service] Repository create failed: %v", err)
		return err
	}

	// Publish sync events for each created presentation
	for _, p := range entities {
		if s.isOffline {
			s.publishSync(NewSyncCreatedEventFromOffline(tenantSlug, p))
		} else {
			s.publishSync(NewSyncCreatedEvent(tenantSlug, p))
		}
	}

	logger.Logf("[Presentation Service] Created %d presentations successfully", len(presentations))
	return nil
}

// Upsert creates or updates presentations for a product
// If PresentationRequest has ID, updates; otherwise creates new
func (s *service) Upsert(ctx context.Context, tenantSlug string, enterpriseID int64, productID int64, presentations []PresentationRequest) error {
	logger := logging.NewLoggerHandler("logs")
	logger.Log(fmt.Sprintf("[Presentation Service] Upserting %d presentations for product ID: %d", len(presentations), productID))

	// Validate product ID
	if productID == 0 {
		logger.Log("[Presentation Service] Validation failed: product_id is required")
		return fmt.Errorf("product_id is required")
	}

	// Validate enterprise ID
	if enterpriseID == 0 {
		logger.Log("[Presentation Service] Validation failed: enterprise_id is required")
		return fmt.Errorf("enterprise_id is required")
	}

	// Validate we have at least one presentation
	if len(presentations) == 0 {
		logger.Log("[Presentation Service] Validation failed: at least one presentation is required")
		return fmt.Errorf("at least one presentation is required")
	}

	// Get existing presentations for the product
	existingPresentations, err := s.repo.GetByProductID(ctx, tenantSlug, productID)
	if err != nil && err != sql.ErrNoRows {
		logger.Logf("[Presentation Service] Error fetching existing presentations: %v", err)
		return fmt.Errorf("error fetching existing presentations: %w", err)
	}

	// Build map of existing presentations by name (for new presentations, check if name exists)
	existingByName := make(map[string]int64)
	for _, p := range existingPresentations {
		existingByName[p.Name] = p.ID
	}

	// Separate presentations to create and to update
	var toCreate []*Presentation
	var toUpdate []PresentationRequest

	for i, req := range presentations {
		logger.Logf("[Presentation Service] Validating presentation %d: name=%s, id=%v", i+1, req.Name, req.ID)

		// Validate required fields
		if req.Name == "" {
			logger.Logf("[Presentation Service] Validation failed: name is required for presentation %d", i+1)
			return fmt.Errorf("name is required for presentation %d", i+1)
		}
		if req.Factor == 0 {
			logger.Logf("[Presentation Service] Validation failed: factor is required for presentation %d", i+1)
			return fmt.Errorf("factor is required for presentation %d", i+1)
		}

		if req.ID != nil && *req.ID > 0 {
			// Has ID - will update
			toUpdate = append(toUpdate, req)
		} else {
			// No ID - check if presentation with same name exists
			if existingID, exists := existingByName[req.Name]; exists {
				// Presentation exists by name - treat as update
				logger.Logf("[Presentation Service] Presentation with name '%s' exists with ID %d, updating", req.Name, existingID)
				req.ID = &existingID
				toUpdate = append(toUpdate, req)
			} else {
				// No existing presentation - create new
				toCreate = append(toCreate, &Presentation{
					ProductID:       productID,
					Name:            req.Name,
					Factor:          req.Factor,
					Barcode:         req.Barcode,
					CostPrice:       req.CostPrice,
					SalePrice:       req.SalePrice,
					DefaultPurchase: req.DefaultPurchase,
					DefaultSale:     req.DefaultSale,
					EnterpriseID:    enterpriseID,
				})
			}
		}
	}

	// Create new presentations
	if len(toCreate) > 0 {
		logger.Logf("[Presentation Service] Creating %d new presentations", len(toCreate))
		err := s.repo.CreateMany(ctx, tenantSlug, enterpriseID, toCreate)
		if err != nil {
			logger.Logf("[Presentation Service] Repository create failed: %v", err)
			return fmt.Errorf("failed to create presentations: %w", err)
		}
		// Publish sync events for created presentations
		for _, p := range toCreate {
			if s.isOffline {
				s.publishSync(NewSyncCreatedEventFromOffline(tenantSlug, p))
			} else {
				s.publishSync(NewSyncCreatedEvent(tenantSlug, p))
			}
		}
	}

	// Update existing presentations
	for _, req := range toUpdate {
		logger.Logf("[Presentation Service] Updating presentation ID: %d", *req.ID)

		// Get existing presentation
		existing, err := s.repo.GetByID(ctx, tenantSlug, *req.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Logf("[Presentation Service] Presentation not found with ID: %d", *req.ID)
				return fmt.Errorf("presentation not found: %d", *req.ID)
			}
			logger.Logf("[Presentation Service] Error fetching presentation: %v", err)
			return fmt.Errorf("error fetching presentation: %w", err)
		}

		// Update fields
		existing.Name = req.Name
		existing.Factor = req.Factor
		existing.Barcode = req.Barcode
		existing.CostPrice = req.CostPrice
		existing.SalePrice = req.SalePrice
		existing.DefaultPurchase = req.DefaultPurchase
		existing.DefaultSale = req.DefaultSale

		err = s.repo.Update(ctx, tenantSlug, existing)
		if err != nil {
			logger.Logf("[Presentation Service] Repository update failed: %v", err)
			return err
		}

		// Publish sync event for updated presentation
		if s.isOffline {
			s.publishSync(NewSyncUpdatedEventFromOffline(tenantSlug, existing))
		} else {
			s.publishSync(NewSyncUpdatedEvent(tenantSlug, existing))
		}
	}

	logger.Logf("[Presentation Service] Upsert completed successfully")
	return nil
}

// GetByID retrieves a presentation by its ID
func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*Presentation, error) {
	return s.repo.GetByID(ctx, tenantSlug, id)
}

// GetByProductID retrieves all presentations for a product
func (s *service) GetByProductID(ctx context.Context, tenantSlug string, productID int64) ([]Presentation, error) {
	return s.repo.GetByProductID(ctx, tenantSlug, productID)
}

// Page retrieves a paginated list of presentations
func (s *service) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.Page(ctx, tenantSlug, enterpriseID, page, limit, search, sort, order, params)
}

// List retrieves a list of presentations with product info
func (s *service) List(ctx context.Context, tenantSlug string, enterpriseID int64, productID int64) ([]PresentationWithProductInfo, error) {
	return s.repo.List(ctx, tenantSlug, enterpriseID, productID)
}

// Update updates an existing presentation
func (s *service) Update(ctx context.Context, tenantSlug string, id int64, p *Presentation) error {
	logger := logging.NewLoggerHandler("logs")
	logger.Logf("[Presentation Service] Starting presentation update for ID: %d", id)

	// Get existing presentation
	existing, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logf("[Presentation Service] Presentation not found with ID: %d", id)
			return fmt.Errorf("presentation not found")
		}
		logger.Logf("[Presentation Service] Error fetching presentation: %v", err)
		return fmt.Errorf("error fetching presentation: %w", err)
	}

	// Update fields
	existing.Name = p.Name
	existing.Factor = p.Factor
	existing.Barcode = p.Barcode
	existing.CostPrice = p.CostPrice
	existing.SalePrice = p.SalePrice
	existing.DefaultPurchase = p.DefaultPurchase
	existing.DefaultSale = p.DefaultSale

	err = s.repo.Update(ctx, tenantSlug, existing)
	if err != nil {
		logger.Logf("[Presentation Service] Repository update failed: %v", err)
		return err
	}

	// Publish sync event
	if s.isOffline {
		s.publishSync(NewSyncUpdatedEventFromOffline(tenantSlug, existing))
	} else {
		s.publishSync(NewSyncUpdatedEvent(tenantSlug, existing))
	}

	logger.Logf("[Presentation Service] Presentation updated successfully")
	return nil
}

// Delete performs a soft delete of a presentation
func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	logger := logging.NewLoggerHandler("logs")
	logger.Logf("[Presentation Service] Starting presentation deletion for ID: %d", id)

	presentation, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logf("[Presentation Service] Presentation not found with ID: %d", id)
			return fmt.Errorf("presentation not found")
		}
		logger.Logf("[Presentation Service] Error fetching presentation: %v", err)
		return fmt.Errorf("error fetching presentation: %w", err)
	}

	err = s.repo.Delete(ctx, tenantSlug, id)
	if err != nil {
		logger.Logf("[Presentation Service] Repository delete failed: %v", err)
		return err
	}

	// Publish sync event
	if s.isOffline {
		s.publishSync(NewSyncDeletedEventFromOffline(tenantSlug, presentation))
	} else {
		s.publishSync(NewSyncDeletedEvent(tenantSlug, presentation))
	}

	logger.Logf("[Presentation Service] Presentation deleted successfully")
	return nil
}

// ─── sync helpers ─────────────────────────────────────────────────────────

// publishSync publishes a cross-server sync event.
// Uses the RabbitMQ sync bus if available (set after /offline/ping), falls back to local eventBus.
func (s *service) publishSync(event events.Event) {
	s.syncMu.RLock()
	bus := s.syncBus
	s.syncMu.RUnlock()

	if bus == nil {
		bus = s.eventBus
	}
	if bus == nil {
		return
	}
	if err := bus.Publish(event); err != nil {
		s.logger.Logf("[Presentation Service] warn: sync publish failed: %v", err)
	}
}

func presentationFromPayload(payload map[string]interface{}) *Presentation {
	return &Presentation{
		ID:              int64FromPayload(payload, "presentation_id"),
		ProductID:       int64FromPayload(payload, "product_id"),
		Name:            strFromPayload(payload, "name"),
		Factor:          floatFromPayload(payload, "factor"),
		Barcode:         strFromPayload(payload, "barcode"),
		CostPrice:       floatFromPayload(payload, "cost_price"),
		SalePrice:       floatFromPayload(payload, "sale_price"),
		DefaultPurchase: boolFromPayload(payload, "default_purchase"),
		DefaultSale:     boolFromPayload(payload, "default_sale"),
		EnterpriseID:    int64FromPayload(payload, "enterprise_id"),
	}
}

func strFromPayload(m map[string]interface{}, k string) string {
	v, _ := m[k].(string)
	return v
}

func floatFromPayload(m map[string]interface{}, k string) float64 {
	v, _ := m[k].(float64)
	return v
}

func int64FromPayload(m map[string]interface{}, k string) int64 {
	v, _ := m[k].(float64)
	return int64(v)
}

func boolFromPayload(m map[string]interface{}, k string) bool {
	v, _ := m[k].(bool)
	return v
}