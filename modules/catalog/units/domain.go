package units

import (
	"context"

	"github.com/cloud-tech-develop/aura-back/shared/domain"
)

// Unit entity
type Unit struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Abbreviation  string  `json:"abbreviation"`
	Active        bool    `json:"active"`
	AllowDecimals bool    `json:"allow_decimals"`
	EnterpriseID  int64   `json:"enterprise_id"`
	GlobalID      string  `json:"global_id"`
	SyncStatus    string  `json:"sync_status"`
	LastSyncedAt  *string `json:"last_synced_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at,omitempty"`
	DeletedAt     *string `json:"deleted_at,omitempty"`
}

// Repository interface
type Repository interface {
	Create(ctx context.Context, tenantSlug string, u *Unit) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Unit, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Unit, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error)
	Update(ctx context.Context, tenantSlug string, u *Unit) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}

// Service interface
type Service interface {
	Create(ctx context.Context, tenantSlug string, u *Unit) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Unit, error)
	List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Unit, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error)
	Update(ctx context.Context, tenantSlug string, id int64, u *Unit) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}
