package categories

import (
	"context"
	"time"
)

// Category entity
type Category struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	ParentID     *int64     `json:"parent_id,omitempty"`
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
	Create(ctx context.Context, c *Category) error
	GetByID(ctx context.Context, id int64) (*Category, error)
	List(ctx context.Context, enterpriseID int64) ([]Category, error)
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id int64) error
}

// Service interface
type Service interface {
	Create(ctx context.Context, c *Category) error
	GetByID(ctx context.Context, id int64) (*Category, error)
	List(ctx context.Context, enterpriseID int64) ([]Category, error)
	Update(ctx context.Context, id int64, c *Category) error
	Delete(ctx context.Context, id int64) error
}
