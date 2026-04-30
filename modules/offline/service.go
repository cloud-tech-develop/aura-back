package offline

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/products"
	catalogProducts "github.com/cloud-tech-develop/aura-back/modules/catalog/products"
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
	repo       Repository
	http       *http.Client
	eventBus   events.EventBus
	tenantMgr  *tenant.Manager
	logger     *logging.LoggerHandler
	productSvc catalogProducts.Service
}

func NewService(database *db.DB, eventBus events.EventBus, tenantMgr *tenant.Manager, productSvc catalogProducts.Service) Service {
	return &service{
		repo: NewRepository(database, productSvc),
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
		eventBus:   eventBus,
		tenantMgr:  tenantMgr,
		logger:     logging.NewLoggerHandler("logs"),
		productSvc: productSvc,
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
	url := fmt.Sprintf("%s/catalog/categories/page?slug=%s&enterprise_id=%d", prodURL, slug, enterpriseID)

	body, _ := json.Marshal(map[string]interface{}{
		"limit": 1000,
		"page":  1,
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		Data struct {
			Items []Category `json:"items"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode categories: %w", err)
	}

	count := 0
	for i := range apiResp.Data.Items {
		apiResp.Data.Items[i].EnterpriseID = enterpriseID
		if err := s.repo.UpsertCategory(ctx, &apiResp.Data.Items[i]); err != nil {
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
	url := fmt.Sprintf("%s/catalog/brands/page?slug=%s&enterprise_id=%d", prodURL, slug, enterpriseID)

	body, _ := json.Marshal(map[string]interface{}{
		"limit": 1000,
		"page":  1,
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		Data struct {
			Items []Brand `json:"items"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode brands: %w", err)
	}

	count := 0
	for i := range apiResp.Data.Items {
		apiResp.Data.Items[i].EnterpriseID = enterpriseID
		if err := s.repo.UpsertBrand(ctx, &apiResp.Data.Items[i]); err != nil {
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
	url := fmt.Sprintf("%s/catalog/units/page?slug=%s&enterprise_id=%d", prodURL, slug, enterpriseID)

	body, _ := json.Marshal(map[string]interface{}{
		"limit": 1000,
		"page":  1,
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		Data struct {
			Items []Unit `json:"items"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode units: %w", err)
	}

	count := 0
	for i := range apiResp.Data.Items {
		apiResp.Data.Items[i].EnterpriseID = enterpriseID
		if err := s.repo.UpsertUnit(ctx, &apiResp.Data.Items[i]); err != nil {
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
		Data []products.Product `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode products: %w", err)
	}

	count := 0
	for _, p := range apiResp.Data {
		fmt.Println("product: ", p.Name, p.ID)
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
		Data []Presentation `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode presentations: %w", err)
	}

	fmt.Println("presentations:", apiResp.Data)

	count := 0
	for i := range apiResp.Data {
		apiResp.Data[i].EnterpriseID = enterpriseID
		if err := s.repo.UpsertPresentation(ctx, &apiResp.Data[i]); err != nil {
			continue
		}
		count++
	}

	mu.Lock()
	result.Presentations = count
	mu.Unlock()
	if count > 0 {
		s.publishEvent(EventPresentationSynced, count, slug, true, "")
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

// SyncAllBySlug synchronizes all tenant data from production
func (s *service) SyncAllBySlug(ctx context.Context, prodURL, token, slug string) (*SyncResult, error) {
	return s.SyncTenantBySlug(ctx, prodURL, token, slug)
}
