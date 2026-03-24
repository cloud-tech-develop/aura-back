package purchases

import (
	"context"
	"time"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
	EventPurchaseCreated   = "purchase.created"
	EventPurchaseReceived  = "purchase.received"
	EventPurchasePaid      = "purchase.paid"
	EventPurchaseCancelled = "purchase.cancelled"
)

const (
	StatusPending   = "PENDING"
	StatusPartial   = "PARTIAL"
	StatusReceived  = "RECEIVED"
	StatusCancelled = "CANCELLED"
	StatusCompleted = "COMPLETED"
)

const (
	PaymentMethodCash     = "CASH"
	PaymentMethodCredit   = "CREDIT"
	PaymentMethodTransfer = "TRANSFER"
	PaymentMethodCheque   = "CHEQUE"
)

// ─── Entities ─────────────────────────────────────────────────────────────────

// PurchaseOrder represents a purchase order to a supplier
type PurchaseOrder struct {
	ID            int64      `json:"id"`
	OrderNumber   string     `json:"order_number"`
	SupplierID    int64      `json:"supplier_id"`
	BranchID      int64      `json:"branch_id"`
	UserID        int64      `json:"user_id"`
	OrderDate     time.Time  `json:"order_date"`
	ExpectedDate  *time.Time `json:"expected_date,omitempty"`
	Status        string     `json:"status"`
	Subtotal      float64    `json:"subtotal"`
	DiscountTotal float64    `json:"discount_total"`
	TaxTotal      float64    `json:"tax_total"`
	Total         float64    `json:"total"`
	Notes         string     `json:"notes"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

// PurchaseOrderItem represents an item in a purchase order
type PurchaseOrderItem struct {
	ID               int64     `json:"id"`
	PurchaseOrderID  int64     `json:"purchase_order_id"`
	ProductID        int64     `json:"product_id"`
	Quantity         float64   `json:"quantity"`
	ReceivedQuantity float64   `json:"received_quantity"`
	UnitCost         float64   `json:"unit_cost"`
	DiscountAmount   float64   `json:"discount_amount"`
	TaxRate          float64   `json:"tax_rate"`
	LineTotal        float64   `json:"line_total"`
	CreatedAt        time.Time `json:"created_at"`
}

// Purchase represents a completed purchase (goods receipt)
type Purchase struct {
	ID              int64      `json:"id"`
	PurchaseNumber  string     `json:"purchase_number"`
	PurchaseOrderID *int64     `json:"purchase_order_id,omitempty"`
	SupplierID      int64      `json:"supplier_id"`
	BranchID        int64      `json:"branch_id"`
	UserID          int64      `json:"user_id"`
	PurchaseDate    time.Time  `json:"purchase_date"`
	Status          string     `json:"status"`
	Subtotal        float64    `json:"subtotal"`
	DiscountTotal   float64    `json:"discount_total"`
	TaxTotal        float64    `json:"tax_total"`
	Total           float64    `json:"total"`
	PaidAmount      float64    `json:"paid_amount"`
	PendingAmount   float64    `json:"pending_amount"`
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// PurchaseItem represents an item in a purchase
type PurchaseItem struct {
	ID             int64     `json:"id"`
	PurchaseID     int64     `json:"purchase_id"`
	ProductID      int64     `json:"product_id"`
	Quantity       float64   `json:"quantity"`
	UnitCost       float64   `json:"unit_cost"`
	DiscountAmount float64   `json:"discount_amount"`
	TaxRate        float64   `json:"tax_rate"`
	LineTotal      float64   `json:"line_total"`
	CreatedAt      time.Time `json:"created_at"`
}

// PurchasePayment represents a payment for a purchase
type PurchasePayment struct {
	ID              int64     `json:"id"`
	PurchaseID      int64     `json:"purchase_id"`
	PaymentMethod   string    `json:"payment_method"`
	Amount          float64   `json:"amount"`
	ReferenceNumber *string   `json:"reference_number,omitempty"`
	Notes           string    `json:"notes"`
	UserID          int64     `json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`
}

// ─── Request DTOs ──────────────────────────────────────────────────────────────

// CreatePurchaseOrderRequest for HU-PUR-001
type CreatePurchaseOrderRequest struct {
	SupplierID   int64              `json:"supplier_id" binding:"required"`
	BranchID     int64              `json:"branch_id" binding:"required"`
	ExpectedDate *time.Time         `json:"expected_date"`
	Notes        string             `json:"notes"`
	Items        []OrderItemRequest `json:"items" binding:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID      int64   `json:"product_id" binding:"required"`
	Quantity       float64 `json:"quantity" binding:"required,gt=0"`
	UnitCost       float64 `json:"unit_cost" binding:"required,gt=0"`
	DiscountAmount float64 `json:"discount_amount"`
	TaxRate        float64 `json:"tax_rate"`
}

// ReceiveGoodsRequest for HU-PUR-002
type ReceiveGoodsRequest struct {
	PurchaseOrderID int64                `json:"purchase_order_id" binding:"required"`
	Items           []ReceiveItemRequest `json:"items" binding:"required,min=1"`
	Notes           string               `json:"notes"`
}

type ReceiveItemRequest struct {
	ProductID      int64   `json:"product_id" binding:"required"`
	Quantity       float64 `json:"quantity" binding:"required,gt=0"`
	UnitCost       float64 `json:"unit_cost" binding:"required,gt=0"`
	DiscountAmount float64 `json:"discount_amount"`
	TaxRate        float64 `json:"tax_rate"`
}

// RecordPaymentRequest for HU-PUR-003
type RecordPaymentRequest struct {
	PurchaseID      int64   `json:"purchase_id" binding:"required"`
	PaymentMethod   string  `json:"payment_method" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	ReferenceNumber string  `json:"reference_number"`
	Notes           string  `json:"notes"`
}

// CancelPurchaseRequest for HU-PUR-004
type CancelPurchaseRequest struct {
	Reason string `json:"reason"`
}

// SupplierSummary represents supplier account summary
type SupplierSummary struct {
	SupplierID     int64   `json:"supplier_id"`
	SupplierName   string  `json:"supplier_name"`
	TotalPurchases int64   `json:"total_purchases"`
	TotalAmount    float64 `json:"total_amount"`
	PaidAmount     float64 `json:"paid_amount"`
	PendingAmount  float64 `json:"pending_amount"`
}

// ─── Repository Interface ─────────────────────────────────────────────────────

type Repository interface {
	// Purchase Order
	CreatePurchaseOrder(ctx context.Context, po *PurchaseOrder) (int64, error)
	GetPurchaseOrderByID(ctx context.Context, id int64) (*PurchaseOrder, error)
	GetPurchaseOrderByNumber(ctx context.Context, orderNumber string) (*PurchaseOrder, error)
	UpdatePurchaseOrder(ctx context.Context, po *PurchaseOrder) error
	ListPurchaseOrders(ctx context.Context, supplierID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]PurchaseOrder, int64, error)

	// Purchase Order Items
	CreatePurchaseOrderItem(ctx context.Context, item *PurchaseOrderItem) error
	GetPurchaseOrderItems(ctx context.Context, orderID int64) ([]PurchaseOrderItem, error)
	UpdatePurchaseOrderItem(ctx context.Context, item *PurchaseOrderItem) error

	// Purchase
	CreatePurchase(ctx context.Context, p *Purchase) (int64, error)
	GetPurchaseByID(ctx context.Context, id int64) (*Purchase, error)
	GetPurchaseByNumber(ctx context.Context, purchaseNumber string) (*Purchase, error)
	UpdatePurchase(ctx context.Context, p *Purchase) error
	ListPurchases(ctx context.Context, supplierID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Purchase, int64, error)

	// Purchase Items
	CreatePurchaseItem(ctx context.Context, item *PurchaseItem) error
	GetPurchaseItems(ctx context.Context, purchaseID int64) ([]PurchaseItem, error)

	// Purchase Payments
	CreatePurchasePayment(ctx context.Context, payment *PurchasePayment) error
	GetPurchasePayments(ctx context.Context, purchaseID int64) ([]PurchasePayment, error)
	GetSupplierSummary(ctx context.Context, supplierID int64) (*SupplierSummary, error)
}

// ─── Service Interface ────────────────────────────────────────────────────────

type Service interface {
	// HU-PUR-001: Create Purchase Order
	CreatePurchaseOrder(ctx context.Context, userID int64, req CreatePurchaseOrderRequest) (*PurchaseOrder, error)

	// HU-PUR-002: Receive Goods
	ReceiveGoods(ctx context.Context, userID int64, req ReceiveGoodsRequest) (*Purchase, error)

	// HU-PUR-003: Record Purchase Payment
	RecordPayment(ctx context.Context, userID int64, req RecordPaymentRequest) (*PurchasePayment, error)

	// HU-PUR-004: Cancel Purchase
	CancelPurchase(ctx context.Context, purchaseID int64, reason string) error

	// HU-PUR-005: View Purchase History
	GetPurchaseHistory(ctx context.Context, supplierID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Purchase, int64, error)

	// HU-PUR-006: Supplier Account Summary
	GetSupplierSummary(ctx context.Context, supplierID int64) (*SupplierSummary, error)

	// Additional
	GetPurchaseOrderByID(ctx context.Context, id int64) (*PurchaseOrder, error)
	GetPurchaseByID(ctx context.Context, id int64) (*Purchase, error)
	ListPurchaseOrders(ctx context.Context, supplierID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]PurchaseOrder, int64, error)
}
