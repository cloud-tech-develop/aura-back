package offline

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// Events
const (
	EventEnterpriseSynced = "offline.enterprise_synced"
	EventPlanSynced       = "offline.plan_synced"
	UserSynced            = "offline.user_synced"
	EventUserRoleSynced   = "offline.user_role_synced"
)

// EventPayload represents the sync event payload
type EventPayload struct {
	Table   string      `json:"table"`
	Count   int         `json:"count"`
	Slug    string      `json:"slug"`
	Success bool       `json:"success"`
	Error   string      `json:"error,omitempty"`
}

// service implements Service
type service struct {
	repo    Repository
	http    *http.Client
	eventBus events.EventBus
}

func NewService(db *sql.DB, eventBus events.EventBus) Service {
	return &service{
		repo: NewRepository(db),
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
		eventBus: eventBus,
	}
}

// SyncAllBySlug fetches enterprise by slug from production,
// saves it to SQLite, then synchronizes plans, users, and user_roles in parallel goroutines
func (s *service) SyncAllBySlug(ctx context.Context, prodURL, token, slug string) (*SyncResult, error) {
	result := &SyncResult{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var asyncErrors []string

	// 1. Fetch and sync enterprise first
	enterprise, err := s.fetchEnterpriseBySlug(ctx, prodURL, token, slug)
	if err != nil {
		return nil, fmt.Errorf("fetch enterprise: %w", err)
	}

	if enterprise == nil {
		return nil, fmt.Errorf("empresa no encontrada: %s", slug)
	}

	if err := s.repo.UpsertEnterprise(ctx, enterprise); err != nil {
		mu.Lock()
		asyncErrors = append(asyncErrors, fmt.Sprintf("enterprise %s: %v", enterprise.Name, err))
		mu.Unlock()
	} else {
		result.Enterprises++
		// Publish event
		s.publishEvent(EventEnterpriseSynced, 1, slug, true, "")
	}

	// 2. Schedule parallel sync for plans, users, user_roles using enterprise ID
	enterpriseID := enterprise.ID
	syncConfigs := []struct {
		name    string
		worker func() error
	}{
		{
			"plans",
			func() error {
				return s.syncPlans(ctx, prodURL, token, enterpriseID, result)
			},
		},
		{
			"users",
			func() error {
				return s.syncUsers(ctx, prodURL, token, enterpriseID, result)
			},
		},
		{
			"user_roles",
			func() error {
				return s.syncUserRoles(ctx, prodURL, token, enterpriseID, result)
			},
		},
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

	result.Errors = asyncErrors
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

// syncPlans fetches and saves plans for the given enterprise
func (s *service) syncPlans(ctx context.Context, prodURL, token string, enterpriseID int64, result *SyncResult) error {
	url := fmt.Sprintf("%s/plans?enterprise_id=%d", prodURL, enterpriseID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

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

	var apiResult struct {
		Data struct {
			Data []Plan `json:"data"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		return err
	}

	count := 0
	for _, plan := range apiResult.Data.Data {
		plan.EnterpriseID = enterpriseID
		if err := s.repo.UpsertPlan(ctx, &plan); err != nil {
			continue
		}
		count++
	}

	result.Plans = count
	if count > 0 {
		s.publishEvent(EventPlanSynced, count, "", true, "")
	}

	return nil
}

// syncUsers fetches and saves users for the given enterprise
func (s *service) syncUsers(ctx context.Context, prodURL, token string, enterpriseID int64, result *SyncResult) error {
	url := fmt.Sprintf("%s/users?enterprise_id=%d", prodURL, enterpriseID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

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

	var apiResult struct {
		Data struct {
			Data []User `json:"data"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		return err
	}

	count := 0
	for _, user := range apiResult.Data.Data {
		user.EnterpriseID = enterpriseID
		if err := s.repo.UpsertUser(ctx, &user); err != nil {
			continue
		}
		count++
	}

	result.Users = count
	if count > 0 {
		s.publishEvent(UserSynced, count, "", true, "")
	}

	return nil
}

// syncUserRoles fetches and saves user_roles for the given enterprise
func (s *service) syncUserRoles(ctx context.Context, prodURL, token string, enterpriseID int64, result *SyncResult) error {
	url := fmt.Sprintf("%s/user-roles?enterprise_id=%d", prodURL, enterpriseID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

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

	var apiResult struct {
		Data struct {
			Data []UserRole `json:"data"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		return err
	}

	count := 0
	for _, ur := range apiResult.Data.Data {
		if err := s.repo.UpsertUserRole(ctx, &ur); err != nil {
			continue
		}
		count++
	}

	result.UserRoles = count
	if count > 0 {
		s.publishEvent(EventUserRoleSynced, count, "", true, "")
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