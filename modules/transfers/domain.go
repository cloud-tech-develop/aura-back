package transfers

import (
	"context"
	"time"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
	EventTransferCreated   = "transfer.created"
	EventTransferApproved  = "transfer.approved"
	EventTransferShipped   = "transfer.shipped"
	EventTransferReceived  = "transfer.received"
	EventTransferCancelled = "transfer.cancelled"
)

const (
	StatusPending   = "PENDING"
	StatusApproved  = "APPROVED"
	StatusShipped   = "SHIPPED"
	StatusPartial   = "PARTIAL"
	StatusReceived  = "RECEIVED"
	StatusCancelled = "CANCELLED"
)

// ─── Entities ─────────────────────────────────────────────────────────────────

// Transfer represents an inventory transfer between branches
type Transfer struct {
	ID                  int64      `json:"id"`
	TransferNumber      string     `json:"transfer_number"`
	OriginBranchID      int64      `json:"origin_branch_id"`
	DestinationBranchID int64      `json:"destination_branch_id"`
	UserID              int64      `json:"user_id"`
	Status              string     `json:"status"`
	RequestedDate       time.Time  `json:"requested_date"`
	ShippedDate         *time.Time `json:"shipped_date,omitempty"`
	ReceivedDate        *time.Time `json:"received_date,omitempty"`
	Notes               string     `json:"notes"`
	ShippedBy           *int64     `json:"shipped_by,omitempty"`
	ReceivedBy          *int64     `json:"received_by,omitempty"`
	CancellationReason  string     `json:"cancellation_reason"`
	CancelledBy         *int64     `json:"cancelled_by,omitempty"`
	CancelledAt         *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
}

// TransferItem represents an item in a transfer
type TransferItem struct {
	ID                int64     `json:"id"`
	TransferID        int64     `json:"transfer_id"`
	ProductID         int64     `json:"product_id"`
	RequestedQuantity float64   `json:"requested_quantity"`
	ShippedQuantity   float64   `json:"shipped_quantity"`
	ReceivedQuantity  float64   `json:"received_quantity"`
	CreatedAt         time.Time `json:"created_at"`
}

// ─── Request DTOs ──────────────────────────────────────────────────────────────

// CreateTransferRequest for HU-TRANS-001
type CreateTransferRequest struct {
	OriginBranchID      int64                 `json:"origin_branch_id" binding:"required"`
	DestinationBranchID int64                 `json:"destination_branch_id" binding:"required"`
	Notes               string                `json:"notes"`
	Items               []TransferItemRequest `json:"items" binding:"required,min=1"`
}

type TransferItemRequest struct {
	ProductID         int64   `json:"product_id" binding:"required"`
	RequestedQuantity float64 `json:"requested_quantity" binding:"required,gt=0"`
}

// ShipTransferRequest for HU-TRANS-003
type ShipTransferRequest struct {
	Items []ShipItemRequest `json:"items" binding:"required,min=1"`
	Notes string            `json:"notes"`
}

type ShipItemRequest struct {
	ProductID       int64   `json:"product_id" binding:"required"`
	ShippedQuantity float64 `json:"shipped_quantity" binding:"required,gt=0"`
}

// ReceiveTransferRequest for HU-TRANS-004
type ReceiveTransferRequest struct {
	Items []ReceiveItemRequest `json:"items" binding:"required,min=1"`
	Notes string               `json:"notes"`
}

type ReceiveItemRequest struct {
	ProductID        int64   `json:"product_id" binding:"required"`
	ReceivedQuantity float64 `json:"received_quantity" binding:"required,gt=0"`
}

// CancelTransferRequest for HU-TRANS-005
type CancelTransferRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// ─── Repository Interface ─────────────────────────────────────────────────────

type Repository interface {
	// Transfer
	CreateTransfer(ctx context.Context, t *Transfer) (int64, error)
	GetTransferByID(ctx context.Context, id int64) (*Transfer, error)
	GetTransferByNumber(ctx context.Context, number string) (*Transfer, error)
	UpdateTransfer(ctx context.Context, t *Transfer) error
	ListTransfers(ctx context.Context, originBranchID, destBranchID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Transfer, int64, error)

	// TransferItem
	CreateTransferItem(ctx context.Context, item *TransferItem) error
	GetTransferItems(ctx context.Context, transferID int64) ([]TransferItem, error)
	UpdateTransferItem(ctx context.Context, item *TransferItem) error
}

// ─── Service Interface ────────────────────────────────────────────────────────

type Service interface {
	// HU-TRANS-001: Create Transfer Request
	CreateTransfer(ctx context.Context, userID int64, req CreateTransferRequest) (*Transfer, error)

	// HU-TRANS-002: Approve Transfer
	ApproveTransfer(ctx context.Context, transferID int64, userID int64) (*Transfer, error)

	// HU-TRANS-003: Ship Transfer
	ShipTransfer(ctx context.Context, transferID int64, userID int64, req ShipTransferRequest) (*Transfer, error)

	// HU-TRANS-004: Receive Transfer
	ReceiveTransfer(ctx context.Context, transferID int64, userID int64, req ReceiveTransferRequest) (*Transfer, error)

	// HU-TRANS-005: Cancel Transfer
	CancelTransfer(ctx context.Context, transferID int64, userID int64, reason string) (*Transfer, error)

	// HU-TRANS-006: View Transfer History
	ListTransfers(ctx context.Context, originBranchID, destBranchID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Transfer, int64, error)

	// Additional
	GetTransferByID(ctx context.Context, id int64) (*Transfer, error)
}
