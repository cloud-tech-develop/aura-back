package categories

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

type service struct {
	repo     Repository
	eventBus events.EventBus
	syncBus  events.EventBus // Cross-server sync bus (RabbitMQ)
	syncMu   sync.RWMutex
	logger   *logging.LoggerHandler
	isOffline bool
}

func NewService(db *db.DB, eventBus events.EventBus) Service {
	logger := logging.NewLoggerHandler("sync")
	logger.Log("[Category Service] Initializing service")

	svc := &service{
		repo:     NewRepository(db),
		eventBus: eventBus,
		logger:   logger,
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
		s.logger.Log("[Category Service] Offline mode: RabbitMQ sync deferred until /offline/ping")
		return
	}

	// Online: wildcard binding receives offline events from all tenants.
	s.eventBus.Subscribe(EventCategoryOnlineCreated+".*", s)
	s.eventBus.Subscribe(EventCategoryOnlineUpdated+".*", s)
	s.eventBus.Subscribe(EventCategoryOnlineDeleted+".*", s)
	s.logger.Log("[Category Service] Online mode: Subscribed to offline category events (wildcard)")
}

// SetSyncBus sets the RabbitMQ event bus for cross-server sync publishing.
func (s *service) SetSyncBus(bus events.EventBus) {
	s.syncMu.Lock()
	defer s.syncMu.Unlock()
	s.syncBus = bus
	s.logger.Log("[Category Service] Sync bus updated to RabbitMQ")
}

// Handle implements events.EventHandler for RabbitMQ events
func (s *service) Handle(event events.Event) error {
	s.logger.Logf("[Category Service] Received event from RabbitMQ: %s", event.GetName())

	payload, ok := event.GetPayload().(map[string]interface{})
	if !ok {
		s.logger.Log("[Category Service] Invalid payload type in event")
		return fmt.Errorf("invalid payload type")
	}

	tenantSlug, ok := payload["tenant_slug"].(string)
	if !ok || tenantSlug == "" {
		s.logger.Log("[Category Service] No tenant_slug in event, skipping")
		return fmt.Errorf("tenant_slug is required")
	}

	s.logger.Logf("[Category Service] Processing %s event for tenant: %s", event.GetName(), tenantSlug)

	ctx := context.Background()

	switch event.GetName() {
	case EventCategoryOfflineCreated, EventCategoryOnlineCreated:
		return s.handleRemoteCreate(ctx, tenantSlug, payload)
	case EventCategoryOfflineUpdated, EventCategoryOnlineUpdated:
		return s.handleRemoteUpdate(ctx, tenantSlug, payload)
	case EventCategoryOfflineDeleted, EventCategoryOnlineDeleted:
		return s.handleRemoteDelete(ctx, tenantSlug, payload)
	}

	return nil
}

func (s *service) handleRemoteCreate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	category := categoryFromPayload(payload)
	s.logger.Logf("[Category Service] Remote create: %s", category.Name)

	existing, err := s.repo.GetByID(ctx, tenantSlug, category.ID)
	if err == nil {
		var eventTime time.Time
		if tStr, ok := payload["timestamp"].(string); ok {
			eventTime, _ = time.Parse(time.RFC3339, tStr)
		}

		if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
			s.logger.Logf("[Category Service] Conflict: local version is newer for %s, skipping", category.Name)
			return nil
		}

		category.CreatedAt = existing.CreatedAt
		if err := s.repo.Update(ctx, tenantSlug, category); err != nil {
			return fmt.Errorf("failed to update category: %w", err)
		}
		s.logger.Logf("[Category Service] Category updated from remote (resolved conflict): %s", category.Name)
		return nil
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking category existence: %w", err)
	}

	if err := s.repo.Create(ctx, tenantSlug, category); err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	s.logger.Logf("[Category Service] Category created from remote: %s", category.Name)
	return nil
}

func (s *service) handleRemoteUpdate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	category := categoryFromPayload(payload)
	s.logger.Logf("[Category Service] Remote update: %s", category.Name)

	existing, err := s.repo.GetByID(ctx, tenantSlug, category.ID)
	if err == sql.ErrNoRows {
		if err := s.repo.Create(ctx, tenantSlug, category); err != nil {
			return fmt.Errorf("failed to create category on remote update: %w", err)
		}
		s.logger.Logf("[Category Service] Category created from remote update: %s", category.Name)
		return nil
	} else if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}

	var eventTime time.Time
	if tStr, ok := payload["timestamp"].(string); ok {
		eventTime, _ = time.Parse(time.RFC3339, tStr)
	}

	if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
		s.logger.Logf("[Category Service] Conflict: local version is newer for %s, skipping update", category.Name)
		return nil
	}

	if err := s.repo.Update(ctx, tenantSlug, category); err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	s.logger.Logf("[Category Service] Category updated from remote: %s", category.Name)
	return nil
}

func (s *service) handleRemoteDelete(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	categoryID := int64(payload["category_id"].(float64))
	s.logger.Logf("[Category Service] Remote delete: ID %d", categoryID)

	existing, err := s.repo.GetByID(ctx, tenantSlug, categoryID)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}

	var eventTime time.Time
	if tStr, ok := payload["timestamp"].(string); ok {
		eventTime, _ = time.Parse(time.RFC3339, tStr)
	}

	if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
		s.logger.Logf("[Category Service] Conflict: local version is newer for ID %d, skipping delete", categoryID)
		return nil
	}

	if err := s.repo.Delete(ctx, tenantSlug, categoryID); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	s.logger.Logf("[Category Service] Category deleted from remote: ID %d", categoryID)
	return nil
}

func (s *service) Create(ctx context.Context, tenantSlug string, c *Category) error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if err := s.repo.Create(ctx, tenantSlug, c); err != nil {
		return err
	}

	// Publish sync event to RabbitMQ for cross-server sync
	if s.isOffline {
		s.publishSync(NewSyncCreatedEventFromOffline(tenantSlug, c))
	} else {
		s.publishSync(NewSyncCreatedEvent(tenantSlug, c))
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*Category, error) {
	return s.repo.GetByID(ctx, tenantSlug, id)
}

func (s *service) List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]domain.ListId, error) {
	return s.repo.List(ctx, tenantSlug, enterpriseID)
}

func (s *service) ListAll(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Category, error) {
	return s.repo.ListAll(ctx, tenantSlug, enterpriseID)
}

func (s *service) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.Page(ctx, tenantSlug, enterpriseID, page, limit, search, sort, order, params)
}

func (s *service) Update(ctx context.Context, tenantSlug string, id int64, c *Category) error {
	existing, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		return err
	}
	c.ID = existing.ID
	c.EnterpriseID = existing.EnterpriseID
	c.CreatedAt = existing.CreatedAt
	if err := s.repo.Update(ctx, tenantSlug, c); err != nil {
		return err
	}

	// Publish sync event to RabbitMQ for cross-server sync
	if s.isOffline {
		s.publishSync(NewSyncUpdatedEventFromOffline(tenantSlug, c))
	} else {
		s.publishSync(NewSyncUpdatedEvent(tenantSlug, c))
	}

	return nil
}

func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	// Get category for sync event
	c, _ := s.repo.GetByID(ctx, tenantSlug, id)

	if err := s.repo.Delete(ctx, tenantSlug, id); err != nil {
		return err
	}

	// Publish sync event to RabbitMQ for cross-server sync
	if c != nil {
		if s.isOffline {
			s.publishSync(NewSyncDeletedEventFromOffline(tenantSlug, c))
		} else {
			s.publishSync(NewSyncDeletedEvent(tenantSlug, c))
		}
	}

	return nil
}

func (s *service) Upsert(ctx context.Context, tenantSlug string, c *Category) error {
	return s.repo.Upsert(ctx, tenantSlug, c)
}

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
		s.logger.Logf("[Category Service] warn: sync publish failed: %v", err)
	}
}

func categoryFromPayload(payload map[string]interface{}) *Category {
	c := &Category{
		ID:             int64FromPayload(payload, "category_id"),
		Name:           strFromPayload(payload, "name"),
		DefaultTaxRate: floatFromPayload(payload, "default_tax_rate"),
		Active:         boolFromPayload(payload, "active"),
		EnterpriseID:   int64FromPayload(payload, "enterprise_id"),
	}
	if v, ok := payload["description"]; ok && v != nil {
		desc := strFromPayload(payload, "description")
		c.Description = &desc
	}
	if v, ok := payload["parent_id"]; ok && v != nil {
		id := int64FromPayload(payload, "parent_id")
		c.ParentID = &id
	}
	return c
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