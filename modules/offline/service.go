package offline

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// service implements Service
type service struct {
	repo Repository
	http *http.Client
}

func NewService(db *sql.DB) Service {
	return &service{
		repo: NewRepository(db),
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SyncAllBySlug fetches enterprises from production filtered by the given slug,
// saves them to SQLite, then synchronizes plans in parallel goroutine if available
func (s *service) SyncAllBySlug(ctx context.Context, prodURL, token, slug string) (*SyncResult, error) {
	result := &SyncResult{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var asyncErrors []string

	// 1. Fetch and sync enterprises
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
	}

	// 2. Try to sync plan in background if enterprise has plan info
	if enterprise != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.syncPlan(ctx, prodURL, token, enterprise); err != nil {
				mu.Lock()
				asyncErrors = append(asyncErrors, fmt.Sprintf("plan: %v", err))
				mu.Unlock()
			}
		}()
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

// syncPlan tries to fetch and sync plan for the enterprise
// Note: This requires a dedicated endpoint which may not exist yet
func (s *service) syncPlan(ctx context.Context, prodURL, token string, enterprise *Enterprise) error {
	// Try to get plan info - this endpoint may not be public
	// For now, we skip this as there's no public /plans endpoint
	// The plan info would need to be included in the enterprise response or require a separate endpoint
	return nil
}

// GetLocalEnterprises returns all enterprises stored locally
func (s *service) GetLocalEnterprises(ctx context.Context) ([]Enterprise, error) {
	return s.repo.ListEnterprises(ctx)
}