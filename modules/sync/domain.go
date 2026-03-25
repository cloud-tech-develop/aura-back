package sync

import (
	"context"
	"time"

	"github.com/cloud-tech-develop/aura-back/modules/products"
	"github.com/cloud-tech-develop/aura-back/modules/third-parties"
	"github.com/cloud-tech-develop/aura-back/modules/sales"
	"github.com/cloud-tech-develop/aura-back/modules/invoices"
)

// SyncBatch represents a bundle of data for synchronization
type SyncBatch struct {
	Categories   []products.Category   `json:"categories,omitempty"`
	Brands       []products.Brand      `json:"brands,omitempty"`
	Products     []products.Product    `json:"products,omitempty"`
	ThirdParties []thirdparties.ThirdParty `json:"third_parties,omitempty"`
	SalesOrders  []sales.SalesOrder    `json:"sales_orders,omitempty"`
	Invoices     []invoices.Invoice    `json:"invoices,omitempty"`
}

// SyncStats represents the result of a sync operation
type SyncStats struct {
	PulledCount int `json:"pulled_count"`
	PushedCount int `json:"pushed_count"`
	FailedCount int `json:"failed_count"`
	Duration    time.Duration `json:"duration"`
}

// Service defines the synchronization operations
type Service interface {
	// Pull fetches updates from the online server (Postgres) to local (SQLite)
	Pull(ctx context.Context, lastSync time.Time) (*SyncBatch, error)
	
	// Push sends local updates (SQLite) to the online server (Postgres)
	Push(ctx context.Context, batch *SyncBatch) (*SyncStats, error)
}
