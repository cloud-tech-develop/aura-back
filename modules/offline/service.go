package offline

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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

// SyncEnterpriseBySlug fetches a single enterprise by slug from production and saves it to local SQLite
func (s *service) SyncEnterpriseBySlug(ctx context.Context, prodURL, token, slug string) (*Enterprise, error) {
	// Fetch enterprise by slug from production
	url := prodURL + "/enterprises/" + slug

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("crear petición: %w", err)
	}

	// Add Authorization header if token provided
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("conectar a producción: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("empresa no encontrada: %s", slug)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("respuesta inválida de producción: %d", resp.StatusCode)
	}

	var result struct {
		Data Enterprise `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decodificar respuesta: %w", err)
	}

	// Save enterprise to local SQLite
	if err := s.repo.Upsert(ctx, &result.Data); err != nil {
		return nil, fmt.Errorf("guardar empresa %s: %w", result.Data.Name, err)
	}

	return &result.Data, nil
}

// GetLocalEnterprises returns all enterprises stored locally
func (s *service) GetLocalEnterprises(ctx context.Context) ([]Enterprise, error) {
	return s.repo.List(ctx)
}
