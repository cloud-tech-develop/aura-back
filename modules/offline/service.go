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

// SyncEnterprises fetches enterprises from production and saves them to local SQLite
func (s *service) SyncEnterprises(ctx context.Context, prodURL, token string) (int, error) {
	// Fetch enterprises from production
	url := prodURL + "/enterprises"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("crear petición: %w", err)
	}

	// Add Authorization header if token provided
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("conectar a producción: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("respuesta inválida de producción: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Data []Enterprise `json:"data"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decodificar respuesta: %w", err)
	}

	// Save each enterprise that doesn't exist locally
	saved := 0
	for _, ent := range result.Data.Data {
		// Enterprise doesn't exist, insert it
		fmt.Println("insertando empresa", ent.Name)
		fmt.Println(ent)
		if err := s.repo.Upsert(ctx, &ent); err != nil {
			return saved, fmt.Errorf("guardar empresa %s: %w", ent.Name, err)
		}
		saved++
	}

	return saved, nil
}

// GetLocalEnterprises returns all enterprises stored locally
func (s *service) GetLocalEnterprises(ctx context.Context) ([]Enterprise, error) {
	return s.repo.List(ctx)
}
