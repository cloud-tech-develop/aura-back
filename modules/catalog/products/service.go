package products

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/presentations"
	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/cloud-tech-develop/aura-back/shared/logging"
)

// service implements the Service interface
// Contains business logic for product management
type service struct {
	repo            Repository
	presentationSvc presentations.Service
	eventBus        events.EventBus // Local bus: audit logs, internal domain events
	syncBus         events.EventBus // Cross-server sync bus (RabbitMQ); set via SetSyncBus after /offline/ping
	syncMu          sync.RWMutex
	logger          *logging.LoggerHandler
	tenantSlug      string
	isOffline       bool
}

// NewService creates a new product service instance
// db: database connection instance
// eventBus: event bus for publishing domain events (RabbitMQ cuando USE_RABBITMQ=true)
func NewService(db *db.DB, eventBus events.EventBus, presSvc presentations.Service) Service {
	logger := logging.NewLoggerHandler("sync")
	logger.Log("[Product Service] Initializing service")

	// Crear el servicio
	svc := &service{
		repo:            NewRepository(db),
		presentationSvc: presSvc,
		eventBus:        eventBus,
		logger:          logger,
		isOffline:       db.IsSQLite(),
	}

	// Intentar suscribirse a los eventos de RabbitMQ si está disponible
	// Esto permite recibir eventos del otro servidor
	svc.subscribeToRabbitMQEvents()

	return svc
}

// subscribeToRabbitMQEvents registers cross-server sync subscriptions.
//
// Offline mode: subscriptions are deferred — the event bus is in-memory at startup
// and RabbitMQ is only available after /offline/ping calls SetSyncBus.
//
// Online mode: subscribes to wildcard routing keys so one queue receives events
// from all offline tenants (routing: product.online.created.*).
func (s *service) subscribeToRabbitMQEvents() {
	if s.eventBus == nil {
		return
	}

	if s.isOffline {
		// Subscriptions happen in SetSyncBus once the tenant slug is known.
		s.logger.Log("[Product Service] Offline mode: RabbitMQ sync deferred until /offline/ping")
		return
	}

	// Online: wildcard binding receives offline events from all tenants.
	// Routing key pattern: product.online.created.*
	s.eventBus.Subscribe(EventProductOnlineCreated+".*", s)
	s.eventBus.Subscribe(EventProductOnlineUpdated+".*", s)
	s.eventBus.Subscribe(EventProductOnlineDeleted+".*", s)
	s.logger.Log("[Product Service] Online mode: Subscribed to offline product events (wildcard)")
}

// SetSyncBus sets the RabbitMQ event bus used for cross-server sync publishing and subscribing.
// Called by offline.Service.ActivateRabbitMQ after /offline/ping resolves the tenant slug.
func (s *service) SetSyncBus(bus events.EventBus) {
	s.syncMu.Lock()
	defer s.syncMu.Unlock()
	s.syncBus = bus
	s.logger.Log("[Product Service] Sync bus updated to RabbitMQ")
}

// Handle implementa events.EventHandler para recibir eventos de RabbitMQ
// Cuando llega un evento del otro servidor, actualiza la base de datos local
func (s *service) Handle(event events.Event) error {
	s.logger.Logf("[Product Service] Received event from RabbitMQ: %s", event.GetName())

	payload, ok := event.GetPayload().(map[string]interface{})
	if !ok {
		s.logger.Log("[Product Service] Invalid payload type in event")
		return fmt.Errorf("invalid payload type")
	}

	// Obtener el tenant del evento
	tenantSlug, ok := payload["tenant_slug"].(string)
	if !ok || tenantSlug == "" {
		s.logger.Log("[Product Service] No tenant_slug in event, skipping")
		return fmt.Errorf("tenant_slug is required")
	}

	s.logger.Logf("[Product Service] Processing %s event for tenant: %s", event.GetName(), tenantSlug)

	switch event.GetName() {
	case EventProductOfflineCreated, EventProductOnlineCreated:
		return s.handleRemoteCreate(tenantSlug, payload)
	case EventProductOfflineUpdated, EventProductOnlineUpdated:
		return s.handleRemoteUpdate(tenantSlug, payload)
	case EventProductOfflineDeleted, EventProductOnlineDeleted:
		return s.handleRemoteDelete(tenantSlug, payload)
	}

	return nil
}

// Upsert upserts a product
// Validates business rules before persisting
func (s *service) Upsert(ctx context.Context, tenantSlug string, p Product) error {
	if err := s.repo.Upsert(ctx, tenantSlug, &p); err != nil {
		return err
	}
	return nil
}

// Create creates a new product in the catalog
// Validates business rules before persisting
func (s *service) Create(ctx context.Context, tenantSlug string, p *Product) error {
	logger := logging.NewLoggerHandler("logs")

	// Validate SKU is required
	if p.SKU == "" {
		return fmt.Errorf("sku is required")
	}

	// Validate name is required
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Validate unit measure is required
	if p.UnitID == 0 {
		return fmt.Errorf("unit_id is required")
	}

	// Validate product type is valid
	if p.ProductType == "" {
		p.ProductType = "STANDARD"
	} else if !IsValidProductType(p.ProductType) {
		return fmt.Errorf("invalid product type: %s", p.ProductType)
	}

	// Set default values for new product
	p.Active = true
	p.VisibleInPOS = true
	if p.ManagesInventory {
		p.CurrentStock = 0
	}

	// Validate pricing: cost and sale prices must be non-negative
	if p.CostPrice < 0 {
		return fmt.Errorf("cost_price cannot be negative")
	}
	if p.SalePrice < 0 {
		return fmt.Errorf("sale_price cannot be negative")
	}

	// Check barcode uniqueness if provided
	if p.Barcode != "" {
		_, err := s.repo.GetByBarcode(ctx, tenantSlug, p.Barcode, p.EnterpriseID)
		if err == nil {
			return fmt.Errorf("barcode %s already exists", p.Barcode)
		}
		if err != sql.ErrNoRows {
			return fmt.Errorf("error checking barcode: %w", err)
		}
	}

	// Create product in repository
	err := s.repo.Create(ctx, tenantSlug, p)
	if err != nil {
		return err
	}

	if p.ID == 0 && len(p.Presentations) > 0 {
		prod, _ := s.repo.GetBySKU(ctx, tenantSlug, p.SKU, p.EnterpriseID)
		if prod != nil {
			p.ID = prod.ID
		}
	}

	// Create presentations if provided
	if len(p.Presentations) > 0 {
		logger.Logf("[Product Service] Creating %d presentations for product ID: %d", len(p.Presentations), p.ID)

		// Convert Product.Presentations to presentations.PresentationRequest
		presRequests := make([]presentations.PresentationRequest, len(p.Presentations))
		for i, pres := range p.Presentations {
			presRequests[i] = presentations.PresentationRequest{
				ID:              pres.ID,
				Name:            pres.Name,
				Factor:          pres.Factor,
				Barcode:         pres.Barcode,
				SalePrice:       pres.SalePrice,
				CostPrice:       pres.CostPrice,
				DefaultPurchase: pres.DefaultPurchase,
				DefaultSale:     pres.DefaultSale,
			}
		}

		if err := s.presentationSvc.Create(ctx, tenantSlug, p.EnterpriseID, p.ID, presRequests); err != nil {
			logger.Logf("[Product Service] Failed to create presentations: %v", err)
			return fmt.Errorf("failed to create presentations: %w", err)
		}
	}

	// Publish domain event to RabbitMQ
	s.publish(NewCreatedEvent(p), tenantSlug)

	// Publish sync event to RabbitMQ for cross-server sync
	if s.isOffline {
		s.publishSync(NewSyncCreatedEventFromOffline(tenantSlug, p))
	} else {
		s.publishSync(NewSyncCreatedEvent(tenantSlug, p))
	}

	logger.Log("[Product Service] Product creation completed successfully")
	return nil
}

// GetByID retrieves a product by its ID
func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error) {
	product, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("error fetching product: %w", err)
	}
	return product, nil
}

// GetBySKU retrieves a product by its SKU code
func (s *service) GetBySKU(ctx context.Context, tenantSlug string, sku string, enterpriseID int64) (*Product, error) {
	product, err := s.repo.GetBySKU(ctx, tenantSlug, sku, enterpriseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("error fetching product by sku: %w", err)
	}
	return product, nil
}

// Page retrieves a paginated list of products
func (s *service) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.Page(ctx, tenantSlug, enterpriseID, page, limit, search, sort, order, params)
}

// List retrieves a list of products with filters
func (s *service) List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Product, error) {
	return s.repo.List(ctx, tenantSlug, enterpriseID)
}

// Update updates an existing product
// Validates business rules before persisting
func (s *service) Update(ctx context.Context, tenantSlug string, id int64, p *Product) error {
	logger := logging.NewLoggerHandler("logs")

	// Get existing product to validate and preserve values
	existing, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("error fetching product: %w", err)
	}

	// Preserve unchanged values from existing product
	if p.SKU == "" {
		p.SKU = existing.SKU
	} else {
		logger.Logf("[Product Service] SKU will be updated from %s to %s", existing.SKU, p.SKU)
	}
	if p.Barcode == "" {
		p.Barcode = existing.Barcode
	}
	if p.Name == "" {
		p.Name = existing.Name
	} else {
		logger.Logf("[Product Service] Name will be updated from %s to %s", existing.Name, p.Name)
	}
	if p.UnitID == 0 {
		p.UnitID = existing.UnitID
		logger.Log("[Product Service] Using existing unit measure")
	}
	if p.ProductType == "" {
		p.ProductType = existing.ProductType
	} else if !IsValidProductType(p.ProductType) {
		return fmt.Errorf("invalid product type: %s", p.ProductType)
	}

	// Preserve pricing if not specified
	if p.SalePrice == 0 {
		p.SalePrice = existing.SalePrice
	} else {
		logger.Logf("[Product Service] Sale price will be updated from %v to %v", existing.SalePrice, p.SalePrice)
	}
	if p.CostPrice == 0 {
		p.CostPrice = existing.CostPrice
	}

	// Set ID and EnterpriseID for update
	p.ID = id
	p.EnterpriseID = existing.EnterpriseID

	// Update in repository
	err = s.repo.Update(ctx, tenantSlug, p)
	if err != nil {
		return err
	}

	// Update presentations if provided
	if len(p.Presentations) > 0 {

		// Convert Product.Presentations to presentations.PresentationRequest
		presRequests := make([]presentations.PresentationRequest, len(p.Presentations))
		for i, pres := range p.Presentations {
			presRequests[i] = presentations.PresentationRequest{
				ID:              pres.ID,
				Name:            pres.Name,
				Factor:          pres.Factor,
				Barcode:         pres.Barcode,
				SalePrice:       pres.SalePrice,
				CostPrice:       pres.CostPrice,
				DefaultPurchase: pres.DefaultPurchase,
				DefaultSale:     pres.DefaultSale,
			}
		}

		if err := s.presentationSvc.Upsert(ctx, tenantSlug, p.EnterpriseID, id, presRequests); err != nil {
			return fmt.Errorf("failed to update presentations: %w", err)
		}
	}

	// Publish domain event to RabbitMQ
	s.publish(NewUpdatedEvent(p), tenantSlug)

	// Publish sync event to RabbitMQ for cross-server sync
	if s.isOffline {
		s.publishSync(NewSyncUpdatedEventFromOffline(tenantSlug, p))
	} else {
		s.publishSync(NewSyncUpdatedEvent(tenantSlug, p))
	}

	logger.Log("[Product Service] Product update completed successfully")
	return nil
}

// Delete performs a soft delete of a product
func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	logger := logging.NewLoggerHandler("logs")

	// Verify product exists
	_, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("error fetching product: %w", err)
	}

	// Soft delete in repository
	err = s.repo.Delete(ctx, tenantSlug, id)
	if err != nil {
		return err
	}

	// Publish domain event
	s.publish(NewDeletedEvent(&Product{ID: id}), tenantSlug)

	// Publish sync event to RabbitMQ for other servers
	if s.isOffline {
		s.publishSync(NewSyncDeletedEventFromOffline(tenantSlug, &Product{ID: id}))
	} else {
		s.publishSync(NewSyncDeletedEvent(tenantSlug, &Product{ID: id}))
	}

	logger.Log("[Product Service] Product deletion completed successfully")
	return nil
}

// publish publishes a domain event through the event bus (RabbitMQ)
func (s *service) publish(event events.Event, tenantSlug string) {
	if s.eventBus == nil {
		return
	}
	if err := s.eventBus.Publish(event); err != nil {
		s.logger.Logf("[Product Service] warn: publish failed: %v", err)
	}
}

// publishSync publishes a cross-server sync event.
// Uses the RabbitMQ sync bus if available (set after /offline/ping), falls back to local bus.
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
		s.logger.Logf("[Product Service] warn: sync publish failed: %v", err)
	}
}

// ─── Handlers para eventos remotos de RabbitMQ ───────────────────────────────────

// handleRemoteCreate crea un producto en la base de datos local desde un evento de RabbitMQ
func (s *service) handleRemoteCreate(tenantSlug string, payload map[string]interface{}) error {
	product := s.productFromPayload(payload)
	s.logger.Logf("[Product Service] Remote create: %s (%s)", product.Name, product.SKU)

	ctx := context.Background()

	// Intentar crear, si ya existe hacer upsert
	existing, err := s.repo.GetByID(ctx, tenantSlug, product.ID)
	if err == nil {
		// Check timestamp to resolve conflict
		var eventTime time.Time
		if tStr, ok := payload["timestamp"].(string); ok {
			eventTime, _ = time.Parse(time.RFC3339, tStr)
		}

		if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
			s.logger.Logf("[Product Service] Conflicto: la versión local es más reciente para %s, omitiendo", product.SKU)
			return nil
		}

		product.CreatedAt = existing.CreatedAt
		if err := s.repo.Update(ctx, tenantSlug, product); err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}
		s.logger.Logf("[Product Service] Product updated from remote (resolved conflict): %s", product.SKU)
		return nil
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking product existence: %w", err)
	}

	if err := s.repo.Create(ctx, tenantSlug, product); err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	s.logger.Logf("[Product Service] Product created from remote: %s", product.SKU)
	return nil
}

// handleRemoteUpdate actualiza un producto en la base de datos local desde un evento de RabbitMQ
func (s *service) handleRemoteUpdate(tenantSlug string, payload map[string]interface{}) error {
	product := s.productFromPayload(payload)
	s.logger.Logf("[Product Service] Remote update: %s (%s)", product.Name, product.SKU)

	ctx := context.Background()
	existing, err := s.repo.GetByID(ctx, tenantSlug, product.ID)
	if err == sql.ErrNoRows {
		// Create if it doesn't exist
		if err := s.repo.Create(ctx, tenantSlug, product); err != nil {
			return fmt.Errorf("failed to create product on remote update: %w", err)
		}
		s.logger.Logf("[Product Service] Product created from remote update: %s", product.SKU)
		return nil
	} else if err != nil {
		return fmt.Errorf("error fetching product: %w", err)
	}

	// Comprobar conflicto
	var eventTime time.Time
	if tStr, ok := payload["timestamp"].(string); ok {
		eventTime, _ = time.Parse(time.RFC3339, tStr)
	}

	if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
		s.logger.Logf("[Product Service] Conflicto: la versión local es más reciente para %s, omitiendo update", product.SKU)
		return nil
	}

	if err := s.repo.Update(ctx, tenantSlug, product); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	s.logger.Logf("[Product Service] Product updated from remote: %s", product.SKU)
	return nil
}

// handleRemoteDelete elimina un producto en la base de datos local desde un evento de RabbitMQ
func (s *service) handleRemoteDelete(tenantSlug string, payload map[string]interface{}) error {
	productID := int64(payload["product_id"].(float64))
	s.logger.Logf("[Product Service] Remote delete: ID %d", productID)

	ctx := context.Background()
	existing, err := s.repo.GetByID(ctx, tenantSlug, productID)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("error fetching product: %w", err)
	}

	var eventTime time.Time
	if tStr, ok := payload["timestamp"].(string); ok {
		eventTime, _ = time.Parse(time.RFC3339, tStr)
	}

	if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
		s.logger.Logf("[Product Service] Conflicto: la versión local es más reciente para ID %d, omitiendo delete", productID)
		return nil
	}

	if err := s.repo.Delete(ctx, tenantSlug, productID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	s.logger.Logf("[Product Service] Product deleted from remote: ID %d", productID)
	return nil
}

// productFromPayload converts a RabbitMQ event payload map to a Product.
// Uses safe type assertions to avoid panics on malformed messages.
func (s *service) productFromPayload(payload map[string]interface{}) *Product {
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
	if v, ok := payload["unit_measure_id"]; ok && v != nil {
		p.UnitID = int64AsFloat(payload, "unit_measure_id")
	}
	if v, ok := payload["price_3"]; ok && v != nil {
		p3 := floatVal(payload, "price_3")
		p.Price3 = &p3
	}

	return p
}

// ─── payload helpers ─────────────────────────────────────────────────────────

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

func intVal(m map[string]interface{}, k string) int {
	v, _ := m[k].(float64)
	return int(v)
}

func boolVal(m map[string]interface{}, k string) bool {
	v, _ := m[k].(bool)
	return v
}
