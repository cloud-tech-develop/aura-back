package shrinkage

import (
	"context"
	"time"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
	EventShrinkageCreated    = "shrinkage.created"
	EventShrinkageAuthorized = "shrinkage.authorized"
	EventShrinkageCancelled  = "shrinkage.cancelled"
)

const (
	StatusPending   = "PENDING"
	StatusApproved  = "APPROVED"
	StatusRejected  = "REJECTED"
	StatusCancelled = "CANCELLED"
)

// ─── Entities ─────────────────────────────────────────────────────────────────

// ShrinkageReason represents a reason for inventory shrinkage
type ShrinkageReason struct {
	ID                     int64     `json:"id"`
	Code                   string    `json:"code"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	RequiresAuthorization  bool      `json:"requires_authorization"`
	AuthorizationThreshold *float64  `json:"authorization_threshold,omitempty"`
	IsActive               bool      `json:"is_active"`
	CreatedAt              time.Time `json:"created_at"`
}

// Shrinkage represents an inventory shrinkage record
type Shrinkage struct {
	ID                 int64      `json:"id"`
	ShrinkageNumber    string     `json:"shrinkage_number"`
	BranchID           int64      `json:"branch_id"`
	UserID             int64      `json:"user_id"`
	ReasonID           int64      `json:"reason_id"`
	ShrinkageDate      time.Time  `json:"shrinkage_date"`
	TotalValue         float64    `json:"total_value"`
	Status             string     `json:"status"`
	Notes              string     `json:"notes"`
	AuthorizedBy       *int64     `json:"authorized_by,omitempty"`
	AuthorizedAt       *time.Time `json:"authorized_at,omitempty"`
	CancellationReason string     `json:"cancellation_reason"`
	CancelledBy        *int64     `json:"cancelled_by,omitempty"`
	CancelledAt        *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

// ShrinkageItem represents an item in a shrinkage record
type ShrinkageItem struct {
	ID           int64     `json:"id"`
	ShrinkageID  int64     `json:"shrinkage_id"`
	ProductID    int64     `json:"product_id"`
	BatchNumber  *string   `json:"batch_number,omitempty"`
	SerialNumber *string   `json:"serial_number,omitempty"`
	Quantity     float64   `json:"quantity"`
	UnitCost     float64   `json:"unit_cost"`
	TotalValue   float64   `json:"total_value"`
	ReasonDetail string    `json:"reason_detail"`
	CreatedAt    time.Time `json:"created_at"`
}

// ─── Request DTOs ──────────────────────────────────────────────────────────────

// RegisterShrinkageRequest for HU-SHR-001
type RegisterShrinkageRequest struct {
	BranchID      int64                  `json:"branch_id" binding:"required"`
	ReasonID      int64                  `json:"reason_id" binding:"required"`
	ShrinkageDate string                 `json:"shrinkage_date"`
	Notes         string                 `json:"notes"`
	Items         []ShrinkageItemRequest `json:"items" binding:"required,min=1"`
}

type ShrinkageItemRequest struct {
	ProductID    int64   `json:"product_id" binding:"required"`
	BatchNumber  string  `json:"batch_number"`
	SerialNumber string  `json:"serial_number"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
	UnitCost     float64 `json:"unit_cost" binding:"required,gte=0"`
	ReasonDetail string  `json:"reason_detail"`
}

// CreateReasonRequest for HU-SHR-002
type CreateReasonRequest struct {
	Code                   string   `json:"code" binding:"required"`
	Name                   string   `json:"name" binding:"required"`
	Description            string   `json:"description"`
	RequiresAuthorization  bool     `json:"requires_authorization"`
	AuthorizationThreshold *float64 `json:"authorization_threshold"`
}

// AuthorizeShrinkageRequest for HU-SHR-003
type AuthorizeShrinkageRequest struct {
	Approved bool   `json:"approved"`
	Notes    string `json:"notes"`
}

// CancelShrinkageRequest for HU-SHR-005
type CancelShrinkageRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// ShrinkageReportItem for reporting
type ShrinkageReportItem struct {
	ReasonID      int64   `json:"reason_id"`
	ReasonName    string  `json:"reason_name"`
	TotalCount    int64   `json:"total_count"`
	TotalQuantity float64 `json:"total_quantity"`
	TotalValue    float64 `json:"total_value"`
}

// ─── Repository Interface ─────────────────────────────────────────────────────

type Repository interface {
	// ShrinkageReason
	CreateReason(ctx context.Context, reason *ShrinkageReason) error
	GetReasonByID(ctx context.Context, id int64) (*ShrinkageReason, error)
	GetReasonByCode(ctx context.Context, code string) (*ShrinkageReason, error)
	ListReasons(ctx context.Context, activeOnly bool) ([]ShrinkageReason, error)
	UpdateReason(ctx context.Context, reason *ShrinkageReason) error
	IsReasonUsed(ctx context.Context, reasonID int64) (bool, error)

	// Shrinkage
	CreateShrinkage(ctx context.Context, s *Shrinkage) (int64, error)
	GetShrinkageByID(ctx context.Context, id int64) (*Shrinkage, error)
	GetShrinkageByNumber(ctx context.Context, number string) (*Shrinkage, error)
	UpdateShrinkage(ctx context.Context, s *Shrinkage) error
	ListShrinkages(ctx context.Context, branchID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Shrinkage, int64, error)

	// ShrinkageItem
	CreateShrinkageItem(ctx context.Context, item *ShrinkageItem) error
	GetShrinkageItems(ctx context.Context, shrinkageID int64) ([]ShrinkageItem, error)

	// Reporting
	GetShrinkageReport(ctx context.Context, branchID *int64, startDate, endDate *time.Time) ([]ShrinkageReportItem, error)
}

// ─── Service Interface ────────────────────────────────────────────────────────

type Service interface {
	// HU-SHR-001: Register Shrinkage
	RegisterShrinkage(ctx context.Context, userID int64, req RegisterShrinkageRequest) (*Shrinkage, error)

	// HU-SHR-002: Configure Shrinkage Reasons
	CreateReason(ctx context.Context, req CreateReasonRequest) (*ShrinkageReason, error)
	ListReasons(ctx context.Context, activeOnly bool) ([]ShrinkageReason, error)

	// HU-SHR-003: Authorize High-Value Shrinkage
	AuthorizeShrinkage(ctx context.Context, shrinkageID int64, userID int64, approved bool, notes string) (*Shrinkage, error)

	// HU-SHR-004: View Shrinkage Report
	GetShrinkageReport(ctx context.Context, branchID *int64, startDate, endDate *time.Time) ([]ShrinkageReportItem, error)
	ListShrinkages(ctx context.Context, branchID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Shrinkage, int64, error)

	// HU-SHR-005: Cancel Shrinkage
	CancelShrinkage(ctx context.Context, shrinkageID int64, userID int64, reason string) (*Shrinkage, error)

	// Additional
	GetShrinkageByID(ctx context.Context, id int64) (*Shrinkage, error)
}
