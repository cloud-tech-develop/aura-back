package brands

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
	logger.Log("[Brand Service] Initializing service")

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
		s.logger.Log("[Brand Service] Offline mode: RabbitMQ sync deferred until /offline/ping")
		return
	}

	// Online: wildcard binding receives offline events from all tenants.
	s.eventBus.Subscribe(EventBrandOnlineCreated+".*", s)
	s.eventBus.Subscribe(EventBrandOnlineUpdated+".*", s)
	s.eventBus.Subscribe(EventBrandOnlineDeleted+".*", s)
	s.logger.Log("[Brand Service] Online mode: Subscribed to offline brand events (wildcard)")
}

// SetSyncBus sets the RabbitMQ event bus for cross-server sync publishing.
func (s *service) SetSyncBus(bus events.EventBus) {
	s.syncMu.Lock()
	defer s.syncMu.Unlock()
	s.syncBus = bus
	s.logger.Log("[Brand Service] Sync bus updated to RabbitMQ")
}

// Handle implements events.EventHandler for RabbitMQ events
func (s *service) Handle(event events.Event) error {
	s.logger.Logf("[Brand Service] Received event from RabbitMQ: %s", event.GetName())

	payload, ok := event.GetPayload().(map[string]interface{})
	if !ok {
		s.logger.Log("[Brand Service] Invalid payload type in event")
		return fmt.Errorf("invalid payload type")
	}

	tenantSlug, ok := payload["tenant_slug"].(string)
	if !ok || tenantSlug == "" {
		s.logger.Log("[Brand Service] No tenant_slug in event, skipping")
		return fmt.Errorf("tenant_slug is required")
	}

	s.logger.Logf("[Brand Service] Processing %s event for tenant: %s", event.GetName(), tenantSlug)

	ctx := context.Background()

	switch event.GetName() {
	case EventBrandOfflineCreated, EventBrandOnlineCreated:
		return s.handleRemoteCreate(ctx, tenantSlug, payload)
	case EventBrandOfflineUpdated, EventBrandOnlineUpdated:
		return s.handleRemoteUpdate(ctx, tenantSlug, payload)
	case EventBrandOfflineDeleted, EventBrandOnlineDeleted:
		return s.handleRemoteDelete(ctx, tenantSlug, payload)
	}

	return nil
}

func (s *service) handleRemoteCreate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	brand := brandFromPayload(payload)
	s.logger.Logf("[Brand Service] Remote create: %s", brand.Name)

	existing, err := s.repo.GetByID(ctx, tenantSlug, brand.ID)
	if err == nil {
		var eventTime time.Time
		if tStr, ok := payload["timestamp"].(string); ok {
			eventTime, _ = time.Parse(time.RFC3339, tStr)
		}

		if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
			s.logger.Logf("[Brand Service] Conflict: local version is newer for %s, skipping", brand.Name)
			return nil
		}

		brand.CreatedAt = existing.CreatedAt
		if err := s.repo.Update(ctx, tenantSlug, brand); err != nil {
			return fmt.Errorf("failed to update brand: %w", err)
		}
		s.logger.Logf("[Brand Service] Brand updated from remote (resolved conflict): %s", brand.Name)
		return nil
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking brand existence: %w", err)
	}

	if err := s.repo.Create(ctx, tenantSlug, brand); err != nil {
		return fmt.Errorf("failed to create brand: %w", err)
	}
	s.logger.Logf("[Brand Service] Brand created from remote: %s", brand.Name)
	return nil
}

func (s *service) handleRemoteUpdate(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	brand := brandFromPayload(payload)
	s.logger.Logf("[Brand Service] Remote update: %s", brand.Name)

	existing, err := s.repo.GetByID(ctx, tenantSlug, brand.ID)
	if err == sql.ErrNoRows {
		if err := s.repo.Create(ctx, tenantSlug, brand); err != nil {
			return fmt.Errorf("failed to create brand on remote update: %w", err)
		}
		s.logger.Logf("[Brand Service] Brand created from remote update: %s", brand.Name)
		return nil
	} else if err != nil {
		return fmt.Errorf("error fetching brand: %w", err)
	}

	var eventTime time.Time
	if tStr, ok := payload["timestamp"].(string); ok {
		eventTime, _ = time.Parse(time.RFC3339, tStr)
	}

	if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
		s.logger.Logf("[Brand Service] Conflict: local version is newer for %s, skipping update", brand.Name)
		return nil
	}

	if err := s.repo.Update(ctx, tenantSlug, brand); err != nil {
		return fmt.Errorf("failed to update brand: %w", err)
	}
	s.logger.Logf("[Brand Service] Brand updated from remote: %s", brand.Name)
	return nil
}

func (s *service) handleRemoteDelete(ctx context.Context, tenantSlug string, payload map[string]interface{}) error {
	brandID := int64(payload["brand_id"].(float64))
	s.logger.Logf("[Brand Service] Remote delete: ID %d", brandID)

	existing, err := s.repo.GetByID(ctx, tenantSlug, brandID)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("error fetching brand: %w", err)
	}

	var eventTime time.Time
	if tStr, ok := payload["timestamp"].(string); ok {
		eventTime, _ = time.Parse(time.RFC3339, tStr)
	}

	if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
		s.logger.Logf("[Brand Service] Conflict: local version is newer for ID %d, skipping delete", brandID)
		return nil
	}

	if err := s.repo.Delete(ctx, tenantSlug, brandID); err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}
	s.logger.Logf("[Brand Service] Brand deleted from remote: ID %d", brandID)
	return nil
}

func (s *service) Create(ctx context.Context, tenantSlug string, b *Brand) error {
	if b.Name == "" {
		return fmt.Errorf("name is required")
	}
	if err := s.repo.Create(ctx, tenantSlug, b); err != nil {
		return err
	}

	// Publish sync event to RabbitMQ for cross-server sync
	if s.isOffline {
		s.publishSync(NewSyncCreatedEventFromOffline(tenantSlug, b))
	} else {
		s.publishSync(NewSyncCreatedEvent(tenantSlug, b))
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*Brand, error) {
	return s.repo.GetByID(ctx, tenantSlug, id)
}

func (s *service) ListAll(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Brand, error) {
	return s.repo.ListAll(ctx, tenantSlug, enterpriseID)
}

func (s *service) List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]BrandList, error) {
	return s.repo.List(ctx, tenantSlug, enterpriseID)
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

func (s *service) Update(ctx context.Context, tenantSlug string, id int64, b *Brand) error {
	existing, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		return err
	}
	b.ID = existing.ID
	b.EnterpriseID = existing.EnterpriseID
	b.CreatedAt = existing.CreatedAt
	if err := s.repo.Update(ctx, tenantSlug, b); err != nil {
		return err
	}

	// Publish sync event to RabbitMQ for cross-server sync
	if s.isOffline {
		s.publishSync(NewSyncUpdatedEventFromOffline(tenantSlug, b))
	} else {
		s.publishSync(NewSyncUpdatedEvent(tenantSlug, b))
	}

	return nil
}

func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	// Get brand for sync event
	b, _ := s.repo.GetByID(ctx, tenantSlug, id)

	if err := s.repo.Delete(ctx, tenantSlug, id); err != nil {
		return err
	}

	// Publish sync event to RabbitMQ for cross-server sync
	if b != nil {
		if s.isOffline {
			s.publishSync(NewSyncDeletedEventFromOffline(tenantSlug, b))
		} else {
			s.publishSync(NewSyncDeletedEvent(tenantSlug, b))
		}
	}

	return nil
}

func (s *service) Upsert(ctx context.Context, tenantSlug string, b *Brand) error {
	return s.repo.Upsert(ctx, tenantSlug, b)
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
		s.logger.Logf("[Brand Service] warn: sync publish failed: %v", err)
	}
}

func brandFromPayload(payload map[string]interface{}) *Brand {
	return &Brand{
		ID:           int64FromPayload(payload, "brand_id"),
		Name:         strFromPayload(payload, "name"),
		Description:  strFromPayload(payload, "description"),
		Active:       boolFromPayload(payload, "active"),
		EnterpriseID: int64FromPayload(payload, "enterprise_id"),
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