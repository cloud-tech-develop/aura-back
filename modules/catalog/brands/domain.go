package brands

import (
	"context"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain"
)

// Brand entity
type Brand struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	EnterpriseID int64      `json:"enterprise_id"`
	GlobalID     string     `json:"global_id"`
	SyncStatus   string     `json:"sync_status"`
	LastSyncedAt *time.Time `json:"last_synced_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// Repository interface
type Repository interface {
	Create(ctx context.Context, tenantSlug string, b *Brand) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Brand, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Brand, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error)
	Update(ctx context.Context, tenantSlug string, b *Brand) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}

// Service interface
type Service interface {
	Create(ctx context.Context, tenantSlug string, b *Brand) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Brand, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Brand, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error)
	Update(ctx context.Context, tenantSlug string, id int64, b *Brand) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}
