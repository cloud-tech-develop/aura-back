package offline

import (
	"context"
	"time"
)

// ─── Entity ───────────────────────────────────────────────────────────────────

// Enterprise represents an enterprise synced from production
type Enterprise struct {
	ID             int64                  `json:"id"`
	TenantID       int64                  `json:"tenant_id"`
	Name           string                 `json:"name"`
	CommercialName string                 `json:"commercial_name"`
	Slug           string                 `json:"slug"`
	SubDomain      string                 `json:"sub_domain"`
	Email          string                 `json:"email"`
	Document       string                 `json:"document"`
	DV             string                 `json:"dv"`
	Phone          string                 `json:"phone"`
	MunicipalityID string                 `json:"municipality_id"`
	Municipality   string                 `json:"municipality"`
	Status         string                 `json:"status"`
	Settings       map[string]interface{} `json:"settings,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      *time.Time             `json:"deleted_at,omitempty"`
}

// ─── Repository Interface ─────────────────────────────────────────────────────

type Repository interface {
	Upsert(ctx context.Context, e *Enterprise) error
	GetBySlug(ctx context.Context, slug string) (*Enterprise, error)
	List(ctx context.Context) ([]Enterprise, error)
}

// ─── Service Interface ────────────────────────────────────────────────────────────

type Service interface {
	SyncEnterprises(ctx context.Context, prodURL string) (int, error)
	GetLocalEnterprises(ctx context.Context) ([]Enterprise, error)
}
