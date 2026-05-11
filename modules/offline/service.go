package offline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cloud-tech-develop/aura-back/infrastructure/messaging/rabbit"
	"github.com/cloud-tech-develop/aura-back/internal/db"
	catalogBrands "github.com/cloud-tech-develop/aura-back/modules/catalog/brands"
	catalogCategories "github.com/cloud-tech-develop/aura-back/modules/catalog/categories"
	catalogCompositions "github.com/cloud-tech-develop/aura-back/modules/catalog/compositions"
	catalogPresentations "github.com/cloud-tech-develop/aura-back/modules/catalog/presentations"
	catalogProducts "github.com/cloud-tech-develop/aura-back/modules/catalog/products"
	catalogUnits "github.com/cloud-tech-develop/aura-back/modules/catalog/units"
	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/cloud-tech-develop/aura-back/shared/logging"
	"github.com/cloud-tech-develop/aura-back/tenant"
)

// Events
const (
	EventEnterpriseSynced   = "offline.enterprise_synced"
	EventPlanSynced         = "offline.plan_synced"
	UserSynced              = "offline.user_synced"
	EventUserRoleSynced     = "offline.user_role_synced"
	EventThirdPartySynced   = "offline.third_party_synced"
	EventCategorySynced     = "offline.category_synced"
	EventBrandSynced        = "offline.brand_synced"
	EventUnitSynced         = "offline.unit_synced"
	EventProductSynced      = "offline.product_synced"
	EventPresentationSynced = "offline.presentation_synced"
	EventCompositionSynced  = "offline.composition_synced"
)

// EventPayload represents the sync event payload
type EventPayload struct {
	Table   string `json:"table"`
	Count   int    `json:"count"`
	Slug    string `json:"slug"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// service implements Service
type service struct {
	repo                 Repository
	http                 *http.Client
	eventBus             events.EventBus
	tenantMgr            *tenant.Manager
	logger                *logging.LoggerHandler
	productSvc           catalogProducts.Service
	presentationSvc      catalogPresentations.Service
	categorySvc          catalogCategories.Service
	brandSvc             catalogBrands.Service
	unitSvc              catalogUnits.Service
	compositionSvc       catalogCompositions.Service
	presentationRepo     catalogPresentations.Repository // Repo para sync handler
	mu                   sync.RWMutex
	rabbitActive         bool
	rabbitBus            *rabbit.RabbitMQEventBus // El event bus de RabbitMQ cuando se activa
}

func NewService(
	database *db.DB,
	eventBus events.EventBus,
	tenantMgr *tenant.Manager,
	productSvc catalogProducts.Service,
	presentationSvc catalogPresentations.Service,
	categorySvc catalogCategories.Service,
	brandSvc catalogBrands.Service,
	unitSvc catalogUnits.Service,
	compositionSvc catalogCompositions.Service,
	presentationRepo catalogPresentations.Repository,
) Service {
	return &service{
		repo: NewRepository(
			database,
			productSvc,
			presentationSvc,
			categorySvc,
			brandSvc,
			unitSvc,
		),
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
		eventBus:         eventBus,
		tenantMgr:        tenantMgr,
		logger:           logging.NewLoggerHandler("logs"),
		productSvc:       productSvc,
		presentationSvc:  presentationSvc,
		categorySvc:      categorySvc,
		brandSvc:         brandSvc,
		unitSvc:          unitSvc,
		compositionSvc:   compositionSvc,
		presentationRepo: presentationRepo,
		mu:               sync.RWMutex{},
		rabbitActive:     false,
		rabbitBus:        nil,
	}
}

// SyncTenantBySlug sincroniza los datos del tenant desde producción
func (s *service) SyncTenantBySlug(ctx context.Context, prodURL, token, slug string) (*SyncResult, error) {
	result := &SyncResult{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var asyncErrors []string

	enterprise, err := s.fetchEnterpriseBySlug(ctx, prodURL, token, slug)
	if err != nil {
		return nil, fmt.Errorf("fetch enterprise: %w", err)
	}
	if enterprise == nil {
		return nil, fmt.Errorf("empresa no encontrada: %s", slug)
	}

	// Run tenant offline migrations before saving data
	// (public tables are already applied at startup via MigrateOffline)
	if s.tenantMgr != nil {
		s.logger.Logf("[offline.Service] Running tenant offline migrations")
		if err := s.tenantMgr.RunOfflineMigrations("offline/tenant"); err != nil {
			s.logger.Logf("[offline.Service] warn: RunOfflineMigrations failed: %v", err)
			// Continue - tables might already exist
		} else {
			s.logger.Logf("[offline.Service] Tenant offline migrations completed")
		}
	}

	// Save enterprise locally
	if err := s.repo.UpsertEnterprise(ctx, enterprise); err != nil {
		return nil, fmt.Errorf("save enterprise locally: %w", err)
	}
	result.Enterprises = 1

	enterpriseID := enterprise.ID

	syncConfigs := []struct {
		name   string
		worker func() error
	}{
		{"plans", func() error { return s.syncPlans(ctx, prodURL, token, enterpriseID, result, &mu) }},
		{"users", func() error { return s.syncUsers(ctx, prodURL, token, enterpriseID, result, &mu) }},
		{"user_roles", func() error { return s.syncUserRoles(ctx, prodURL, token, enterpriseID, result, &mu) }},
		{"third_parties", func() error {
			return s.syncThirdParties(ctx, prodURL, token, enterprise.Slug, enterpriseID, result, &mu)
		}},
		{"categories", func() error { return s.syncCategories(ctx, prodURL, token, enterprise.Slug, enterpriseID, result, &mu) }},
		{"brands", func() error { return s.syncBrands(ctx, prodURL, token, enterprise.Slug, enterpriseID, result, &mu) }},
		{"units", func() error { return s.syncUnits(ctx, prodURL, token, enterprise.Slug, enterpriseID, result, &mu) }},
		{"products", func() error { return s.syncProducts(ctx, prodURL, token, enterprise.Slug, enterpriseID, result, &mu) }},
		{"presentations", func() error {
			return s.syncPresentations(ctx, prodURL, token, enterprise.Slug, enterpriseID, result, &mu)
		}},
		{"compositions", func() error {
			return s.syncCompositions(ctx, prodURL, token, enterprise.Slug, enterpriseID, result, &mu)
		}},
	}

	for _, cfg := range syncConfigs {
		wg.Add(1)
		go func(name string, worker func() error) {
			defer wg.Done()
			if err := worker(); err != nil {
				mu.Lock()
				asyncErrors = append(asyncErrors, fmt.Sprintf("%s: %v", name, err))
				mu.Unlock()
			}
		}(cfg.name, cfg.worker)
	}

	wg.Wait()

	mu.Lock()
	result.Errors = asyncErrors
	mu.Unlock()
	return result, nil
}

// fetchEnterpriseBySlug fetches a single enterprise by slug
func (s *service) fetchEnterpriseBySlug(ctx context.Context, prodURL, token, slug string) (*Enterprise, error) {
	url := prodURL + "/enterprises/" + slug

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	var apiResult struct {
		Data Enterprise `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		return nil, err
	}

	return &apiResult.Data, nil
}

// ─── Tenant Sync Methods ─────────────────────────────────────────────────────

func (s *service) syncPlans(ctx context.Context, prodURL, token string, enterpriseID int64, result *SyncResult, mu *sync.Mutex) error {
	url := fmt.Sprintf("%s/plans?enterprise_id=%d", prodURL, enterpriseID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	var apiResp struct {
		Data []Plan `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode plans: %w", err)
	}

	count := 0
	for i := range apiResp.Data {
		apiResp.Data[i].EnterpriseID = enterpriseID
		if err := s.repo.UpsertPlan(ctx, &apiResp.Data[i]); err != nil {
			continue
		}
		count++
	}

	mu.Lock()
	result.Plans = count
	mu.Unlock()
	return nil
}

func (s *service) syncUsers(ctx context.Context, prodURL, token string, enterpriseID int64, result *SyncResult, mu *sync.Mutex) error {
	url := fmt.Sprintf("%s/users-sync?enterprise_id=%d", prodURL, enterpriseID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	var apiResp struct {
		Data []UserWithPassword `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode users: %w", err)
	}

	count := 0
	for i := range apiResp.Data {
		u := &User{
			ID:           apiResp.Data[i].ID,
			EnterpriseID: enterpriseID,
			Name:         apiResp.Data[i].Name,
			Email:        apiResp.Data[i].Email,
			Active:       apiResp.Data[i].Active,
			PasswordHash: apiResp.Data[i].PasswordHash,
			CreatedAt:    apiResp.Data[i].CreatedAt,
			UpdatedAt:    apiResp.Data[i].UpdatedAt,
		}
		if err := s.repo.UpsertUser(ctx, u); err != nil {
			continue
		}
		count++
	}

	mu.Lock()
	result.Users = count
	mu.Unlock()
	return nil
}

func (s *service) syncUserRoles(ctx context.Context, prodURL, token string, enterpriseID int64, result *SyncResult, mu *sync.Mutex) error {
	url := fmt.Sprintf("%s/user-roles?enterprise_id=%d", prodURL, enterpriseID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	var apiResp struct {
		Data []UserRole `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode user_roles: %w", err)
	}

	count := 0
	for i := range apiResp.Data {
		if err := s.repo.UpsertUserRole(ctx, &apiResp.Data[i]); err != nil {
			continue
		}
		count++
	}

	mu.Lock()
	result.UserRoles = count
	mu.Unlock()
	return nil
}

func (s *service) syncThirdParties(ctx context.Context, prodURL, token string, slug string, enterpriseID int64, result *SyncResult, mu *sync.Mutex) error {
	url := fmt.Sprintf("%s/admin/third-parties?slug=%s&enterprise_id=%d&limit=1000", prodURL, slug, enterpriseID)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	// Decode using wrapper format: {"data": [...], "success": true}
	var apiResp struct {
		Data []ThirdParty `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode third_parties: %w", err)
	}

	count := 0
	for i := range apiResp.Data {
		apiResp.Data[i].EnterpriseID = enterpriseID
		if err := s.repo.UpsertThirdParty(ctx, &apiResp.Data[i]); err != nil {
			continue
		}
		count++
	}

	mu.Lock()
	result.ThirdParties = count
	mu.Unlock()
	if count > 0 {
		s.publishEvent(EventThirdPartySynced, count, slug, true, "")
	}
	return nil
}

func (s *service) syncCategories(ctx context.Context, prodURL, token string, slug string, enterpriseID int64, result *SyncResult, mu *sync.Mutex) error {
	url := prodURL + "/catalog/categories/all"

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	var apiResp struct {
		Data []catalogCategories.Category `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode categories: %w", err)
	}

	fmt.Println("Total categories: ", len(apiResp.Data))

	count := 0
	for i := range apiResp.Data {
		apiResp.Data[i].EnterpriseID = enterpriseID
		if err := s.categorySvc.Create(ctx, slug, &apiResp.Data[i]); err != nil {
			continue
		}
		count++
	}

	mu.Lock()
	result.Categories = count
	mu.Unlock()
	if count > 0 {
		s.publishEvent(EventCategorySynced, count, slug, true, "")
	}
	return nil
}

func (s *service) syncBrands(ctx context.Context, prodURL, token string, slug string, enterpriseID int64, result *SyncResult, mu *sync.Mutex) error {
	url := prodURL + "/catalog/brands/all"

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	var apiResp struct {
		Data []catalogBrands.Brand `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode brands: %w", err)
	}

	fmt.Println("Total Brands: ", len(apiResp.Data))

	count := 0
	for i := range apiResp.Data {
		apiResp.Data[i].EnterpriseID = enterpriseID
		if err := s.brandSvc.Upsert(ctx, slug, &apiResp.Data[i]); err != nil {
			continue
		}
		count++
	}

	mu.Lock()
	result.Brands = count
	mu.Unlock()
	if count > 0 {
		s.publishEvent(EventBrandSynced, count, slug, true, "")
	}
	return nil
}

func (s *service) syncUnits(ctx context.Context, prodURL, token string, slug string, enterpriseID int64, result *SyncResult, mu *sync.Mutex) error {
	url := prodURL + "/catalog/units/all"

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	var apiResp struct {
		Data []catalogUnits.Unit `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode units: %w", err)
	}

	fmt.Println("Total Units: ", len(apiResp.Data))

	count := 0
	for i := range apiResp.Data {
		apiResp.Data[i].EnterpriseID = enterpriseID
		if err := s.unitSvc.Upsert(ctx, slug, &apiResp.Data[i]); err != nil {
			continue
		}
		count++
	}

	mu.Lock()
	result.Units = count
	mu.Unlock()
	if count > 0 {
		s.publishEvent(EventUnitSynced, count, slug, true, "")
	}
	return nil
}

func (s *service) syncProducts(ctx context.Context, prodURL, token string, slug string, enterpriseID int64, result *SyncResult, mu *sync.Mutex) error {
	url := prodURL + "/catalog/products"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	var apiResp struct {
		Data []catalogProducts.Product `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode products: %w", err)
	}

	count := 0
	for _, p := range apiResp.Data {
		if err := s.productSvc.Upsert(ctx, slug, p); err != nil {
			continue
		}
		count++
	}

	mu.Lock()
	result.Products = count
	mu.Unlock()
	if count > 0 {
		s.publishEvent(EventProductSynced, count, slug, true, "")
	}
	return nil
}

func (s *service) syncPresentations(ctx context.Context, prodURL, token string, slug string, enterpriseID int64, result *SyncResult, mu *sync.Mutex) error {
	url := prodURL + "/catalog/presentations"

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	var apiResp struct {
		Data []catalogPresentations.Presentation `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode presentations: %w", err)
	}

	fmt.Println("Total presentations: ", len(apiResp.Data))

	// Si no hay presentaciones, no hay nada que sincronizar
	if len(apiResp.Data) == 0 {
		return nil
	}

	item := apiResp.Data[0]
	presentations := make([]catalogPresentations.PresentationRequest, len(apiResp.Data))
	for i, p := range apiResp.Data {
		presentations[i] = catalogPresentations.PresentationRequest{
			Name:            p.Name,
			Factor:          p.Factor,
			Barcode:         p.Barcode,
			CostPrice:       p.CostPrice,
			SalePrice:       p.SalePrice,
			DefaultPurchase: p.DefaultPurchase,
			DefaultSale:     p.DefaultSale,
		}
	}
	count := 0
	if err := s.presentationSvc.Upsert(ctx, slug, item.EnterpriseID, item.ProductID, presentations); err != nil {
		return err
	}
	count = len(presentations)

	mu.Lock()
	result.Presentations = count
	mu.Unlock()
	if count > 0 {
		s.publishEvent(EventPresentationSynced, count, slug, true, "")
	}
	return nil
}

func (s *service) syncCompositions(ctx context.Context, prodURL, token string, slug string, enterpriseID int64, result *SyncResult, mu *sync.Mutex) error {
	url := prodURL + "/catalog/compositions/all"

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	var apiResp struct {
		Data []catalogCompositions.Composition `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode compositions: %w", err)
	}

	count := 0
	for _, c := range apiResp.Data {
		c.EnterpriseID = enterpriseID
		if err := s.compositionSvc.Create(ctx, slug, &c); err != nil {
			continue
		}
		count++
	}

	mu.Lock()
	result.Compositions = count
	mu.Unlock()
	if count > 0 {
		s.publishEvent(EventCompositionSynced, count, slug, true, "")
	}
	return nil
}

// publishEvent publishes an event to the event bus
func (s *service) publishEvent(eventName string, count int, slug string, success bool, errMsg string) {
	if s.eventBus == nil {
		return
	}

	payload := EventPayload{
		Table:   eventName,
		Count:   count,
		Slug:    slug,
		Success: success,
		Error:   errMsg,
	}

	event := events.NewBaseEvent(eventName, payload)
	if err := s.eventBus.Publish(event); err != nil {
		fmt.Printf("[offline.Service] warn: publish failed: %v\n", err)
	}
}

// GetLocalEnterprises returns all enterprises stored locally
func (s *service) GetLocalEnterprises(ctx context.Context) ([]Enterprise, error) {
	return s.repo.ListEnterprises(ctx)
}

// GetActiveEnterprise returns the first (active) enterprise from SQLite
// Used for automatic sync on app startup when no JWT is available
func (s *service) GetActiveEnterprise(ctx context.Context) (*Enterprise, error) {
	enterprises, err := s.repo.ListEnterprises(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing enterprises: %w", err)
	}

	if len(enterprises) == 0 {
		return nil, fmt.Errorf("no enterprises found in local database")
	}

	// Return the first enterprise as the active one
	// In a multi-tenant scenario, you might want to select based on status or user preference
	return &enterprises[0], nil
}

// SyncAllBySlug synchronizes all tenant data from production
func (s *service) SyncAllBySlug(ctx context.Context, prodURL, token, slug string) (*SyncResult, error) {
	return s.SyncTenantBySlug(ctx, prodURL, token, slug)
}

// syncActivator is implemented by services that support swapping their sync bus at runtime.
type syncActivator interface {
	SetSyncBus(bus events.EventBus)
}

// ActivateRabbitMQ wires RabbitMQ for a specific tenant after /offline/ping resolves the slug.
//
// Sequence:
//  1. Create a tenant-scoped RabbitMQ bus (routing key suffix = slug).
//  2. Subscribe catalog service handlers to the "online→offline" events BEFORE Start().
//  3. Start consumers.
//  4. Tell catalog services to publish their sync events through RabbitMQ instead of memory.
func (s *service) ActivateRabbitMQ(ctx context.Context, slug string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.rabbitActive {
		s.logger.Logf("[offline.Service] RabbitMQ already active for tenant: %s", slug)
		return nil
	}

	s.logger.Logf("[offline.Service] Activating RabbitMQ for tenant: %s", slug)

	rb, err := rabbit.NewRabbitMQEventBusWithTenant(slug)
	if err != nil {
		return fmt.Errorf("failed to create RabbitMQ event bus: %w", err)
	}

	// Register catalog event handlers BEFORE Start() so consumer goroutines
	// are launched for these events when Start() is called.
	// Products: service.Handle() processes create/update/delete from online.
	if handler, ok := s.productSvc.(events.EventHandler); ok {
		rb.Subscribe(catalogProducts.EventProductOfflineCreated, handler)
		rb.Subscribe(catalogProducts.EventProductOfflineUpdated, handler)
		rb.Subscribe(catalogProducts.EventProductOfflineDeleted, handler)
		s.logger.Logf("[offline.Service] Product handler subscribed to offline events for tenant: %s", slug)
	}

	// Categories sync handler
	if categoryHandler, ok := s.categorySvc.(events.EventHandler); ok {
		rb.Subscribe(catalogCategories.EventCategoryOfflineCreated, categoryHandler)
		rb.Subscribe(catalogCategories.EventCategoryOfflineUpdated, categoryHandler)
		rb.Subscribe(catalogCategories.EventCategoryOfflineDeleted, categoryHandler)
		s.logger.Logf("[offline.Service] Category handler subscribed to offline events for tenant: %s", slug)
	}

	// Brands sync handler
	if brandHandler, ok := s.brandSvc.(events.EventHandler); ok {
		rb.Subscribe(catalogBrands.EventBrandOfflineCreated, brandHandler)
		rb.Subscribe(catalogBrands.EventBrandOfflineUpdated, brandHandler)
		rb.Subscribe(catalogBrands.EventBrandOfflineDeleted, brandHandler)
		s.logger.Logf("[offline.Service] Brand handler subscribed to offline events for tenant: %s", slug)
	}

	// Presentations sync handler - presentations use a separate SyncHandler
	presSyncHandler := catalogPresentations.NewSyncHandler(s.presentationRepo)
	rb.Subscribe(catalogPresentations.EventPresentationOfflineCreated, presSyncHandler)
	rb.Subscribe(catalogPresentations.EventPresentationOfflineUpdated, presSyncHandler)
	rb.Subscribe(catalogPresentations.EventPresentationOfflineDeleted, presSyncHandler)
	s.logger.Logf("[offline.Service] Presentation handler subscribed to offline events for tenant: %s", slug)

	if err := rb.Start(); err != nil {
		return fmt.Errorf("failed to start RabbitMQ event bus: %w", err)
	}

	// Update catalog services so they publish sync events to RabbitMQ (not the memory bus).
	// The RabbitMQ bus has b.tenant = slug, so routing keys automatically include the slug.
	if activator, ok := s.productSvc.(syncActivator); ok {
		activator.SetSyncBus(rb)
		s.logger.Logf("[offline.Service] Product service sync bus set to RabbitMQ for tenant: %s", slug)
	}

	s.rabbitBus = rb
	s.rabbitActive = true
	s.logger.Logf("[offline.Service] RabbitMQ activated for tenant: %s", slug)
	return nil
}
