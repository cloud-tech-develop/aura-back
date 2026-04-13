package categories

import (
	"context"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain"
)

// Category entity
type Category struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	DefaultTaxRate float64    `json:"default_tax_rate"`
	Active         bool       `json:"active"`
	ParentID       *int64     `json:"parent_id,omitempty"`
	EnterpriseID   int64      `json:"enterprise_id"`
	GlobalID       string     `json:"global_id"`
	SyncStatus     string     `json:"sync_status"`
	LastSyncedAt   *time.Time `json:"last_synced_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// Repository interface
type Repository interface {
	Create(ctx context.Context, tenantSlug string, c *Category) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Category, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Category, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, first int64, rows int64, search string) (domain.PageResult, error)
	Update(ctx context.Context, tenantSlug string, c *Category) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}

// Service interface
type Service interface {
	Create(ctx context.Context, tenantSlug string, c *Category) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Category, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Category, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, first int64, rows int64, search string) (domain.PageResult, error)
	Update(ctx context.Context, tenantSlug string, id int64, c *Category) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}
