package cash

import (
	"context"
	"time"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
	EventShiftOpened    = "cash.shift.opened"
	EventShiftClosed    = "cash.shift.closed"
	EventMovement       = "cash.movement"
	EventReconciliation = "cash.reconciliation"
)

const (
	StatusOpen    = "OPEN"
	StatusClosed  = "CLOSED"
	StatusAudited = "AUDITED"
)

const (
	MovementIn  = "IN"
	MovementOut = "OUT"
)

const (
	ReasonSale       = "SALE"
	ReasonOpening    = "OPENING"
	ReasonClosing    = "CLOSING"
	ReasonExpense    = "EXPENSE"
	ReasonDrops      = "DROPS"
	ReasonWithdrawal = "WITHDRAWAL"
	ReasonAdjustment = "ADJUSTMENT"
	ReasonRefund     = "REFUND"
)

// ─── Entities ─────────────────────────────────────────────────────────────────

// CashDrawer represents a cash drawer configuration
type CashDrawer struct {
	ID        int64      `json:"id"`
	BranchID  int64      `json:"branch_id"`
	Name      string     `json:"name"`
	IsActive  bool       `json:"is_active"`
	MinFloat  float64    `json:"min_float"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// CashShift represents an open/closed cash shift
type CashShift struct {
	ID             int64      `json:"id"`
	CashDrawerID   int64      `json:"cash_drawer_id"`
	UserID         int64      `json:"user_id"`
	BranchID       int64      `json:"branch_id"`
	OpeningAmount  float64    `json:"opening_amount"`
	ClosingAmount  *float64   `json:"closing_amount,omitempty"`
	ExpectedAmount *float64   `json:"expected_amount,omitempty"`
	Difference     *float64   `json:"difference,omitempty"`
	OpeningNotes   string     `json:"opening_notes"`
	ClosingNotes   string     `json:"closing_notes"`
	Status         string     `json:"status"`
	OpenedAt       time.Time  `json:"opened_at"`
	ClosedAt       *time.Time `json:"closed_at,omitempty"`
	ClosedBy       *int64     `json:"closed_by,omitempty"`
	AuthorizedBy   *int64     `json:"authorized_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// CashMovement represents a cash entry/exit
type CashMovement struct {
	ID            int64     `json:"id"`
	ShiftID       int64     `json:"shift_id"`
	MovementType  string    `json:"movement_type"` // IN or OUT
	Reason        string    `json:"reason"`        // SALE, OPENING, etc.
	Amount        float64   `json:"amount"`
	ReferenceID   *int64    `json:"reference_id,omitempty"`
	ReferenceType *string   `json:"reference_type,omitempty"`
	Notes         string    `json:"notes"`
	UserID        int64     `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// ─── Request/Response DTOs ──────────────────────────────────────────────────

// ConfigureDrawerRequest is for HU-CASH-001
type ConfigureDrawerRequest struct {
	BranchID int64   `json:"branch_id" binding:"required"`
	Name     string  `json:"name"`
	MinFloat float64 `json:"min_float"`
}

// OpenShiftRequest is for HU-CASH-002
type OpenShiftRequest struct {
	CashDrawerID  int64   `json:"cash_drawer_id" binding:"required"`
	OpeningAmount float64 `json:"opening_amount" binding:"required"`
	OpeningNotes  string  `json:"opening_notes"`
}

// CloseShiftRequest is for HU-CASH-003
type CloseShiftRequest struct {
	ClosingAmount float64 `json:"closing_amount" binding:"required"`
	ClosingNotes  string  `json:"closing_notes"`
}

// RecordMovementRequest is for HU-CASH-004
type RecordMovementRequest struct {
	ShiftID       int64   `json:"shift_id" binding:"required"`
	MovementType  string  `json:"movement_type" binding:"required"` // IN or OUT
	Reason        string  `json:"reason" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	ReferenceID   *int64  `json:"reference_id"`
	ReferenceType *string `json:"reference_type"`
	Notes         string  `json:"notes"`
}

// ReconcileRequest is for HU-CASH-005
type ReconcileRequest struct {
	ExpectedAmount float64 `json:"expected_amount" binding:"required"`
	Notes          string  `json:"notes"`
}

// ShiftSummaryResponse is for HU-CASH-006
type ShiftSummaryResponse struct {
	Shift        CashShift      `json:"shift"`
	Movements    []CashMovement `json:"movements"`
	TotalIn      float64        `json:"total_in"`
	TotalOut     float64        `json:"total_out"`
	ExpectedCash float64        `json:"expected_cash"`
}

// ─── Repository Interface ─────────────────────────────────────────────────────

type Repository interface {
	// CashDrawer
	CreateDrawer(ctx context.Context, d *CashDrawer) error
	GetDrawerByBranch(ctx context.Context, branchID int64) (*CashDrawer, error)
	UpdateDrawer(ctx context.Context, d *CashDrawer) error

	// CashShift
	CreateShift(ctx context.Context, s *CashShift) error
	GetShiftByID(ctx context.Context, id int64) (*CashShift, error)
	GetActiveShiftByUser(ctx context.Context, userID int64) (*CashShift, error)
	UpdateShift(ctx context.Context, s *CashShift) error
	ListShifts(ctx context.Context, branchID int64, startDate, endDate *time.Time, status string, page, limit int) ([]CashShift, int64, error)

	// CashMovement
	CreateMovement(ctx context.Context, m *CashMovement) error
	GetMovementsByShift(ctx context.Context, shiftID int64) ([]CashMovement, error)
	GetMovementsSummary(ctx context.Context, shiftID int64) (totalIn, totalOut float64, err error)
}

// ─── Service Interface ────────────────────────────────────────────────────────

type Service interface {
	// HU-CASH-001: Configure Cash Drawer
	ConfigureDrawer(ctx context.Context, branchID int64, name string, minFloat float64) (*CashDrawer, error)

	// HU-CASH-002: Open Cash Shift
	OpenShift(ctx context.Context, userID, branchID, cashDrawerID int64, openingAmount float64, notes string) (*CashShift, error)

	// HU-CASH-003: Close Cash Shift
	CloseShift(ctx context.Context, shiftID int64, userID int64, closingAmount float64, notes string) (*CashShift, error)

	// HU-CASH-004: Record Cash Movement
	RecordMovement(ctx context.Context, userID int64, req RecordMovementRequest) (*CashMovement, error)

	// HU-CASH-005: Perform Cash Reconciliation
	ReconcileShift(ctx context.Context, shiftID int64, expectedAmount float64, notes string) (*CashShift, error)

	// HU-CASH-006: View Shift Summary
	GetShiftSummary(ctx context.Context, shiftID int64) (*ShiftSummaryResponse, error)

	// Additional
	GetActiveShift(ctx context.Context, userID int64) (*CashShift, error)
	GetDrawerByBranch(ctx context.Context, branchID int64) (*CashDrawer, error)
	ListShifts(ctx context.Context, branchID int64, startDate, endDate *time.Time, status string, page, limit int) ([]CashShift, int64, error)
}
