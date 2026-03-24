package sales

import (
	"context"
	"time"
)

// SalesOrder represents a sales order
type SalesOrder struct {
	ID           int64        `json:"id"`
	OrderNumber  string       `json:"order_number"`
	CustomerID   *int64       `json:"customer_id,omitempty"`
	UserID       int64        `json:"user_id"`
	BranchID     int64        `json:"branch_id"`
	EnterpriseID int64        `json:"enterprise_id"`
	Subtotal     float64      `json:"subtotal"`
	Discount     float64      `json:"discount"`
	TaxTotal     float64      `json:"tax_total"`
	Total        float64      `json:"total"`
	Status       string       `json:"status"`
	Notes        string       `json:"notes"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    *time.Time   `json:"updated_at,omitempty"`
	Items        []SalesOrderItem `json:"items,omitempty"`
}

// SalesOrderItem represents an item in a sales order
type SalesOrderItem struct {
	ID              int64      `json:"id"`
	SalesOrderID    int64      `json:"sales_order_id"`
	ProductID       int64      `json:"product_id"`
	Quantity        int        `json:"quantity"`
	UnitPrice       float64    `json:"unit_price"`
	DiscountPercent float64    `json:"discount_percent"`
	DiscountAmount  float64    `json:"discount_amount"`
	TaxRate         float64    `json:"tax_rate"`
	TaxAmount       float64    `json:"tax_amount"`
	Total           float64    `json:"total"`
	CreatedAt       time.Time  `json:"created_at"`
}

// Repository interface for sales order operations
type Repository interface {
	CreateOrder(ctx context.Context, order *SalesOrder) error
	GetOrderByID(ctx context.Context, id int64) (*SalesOrder, error)
	GetOrderByNumber(ctx context.Context, orderNumber string, enterpriseID int64) (*SalesOrder, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) error
	GetOrders(ctx context.Context, enterpriseID int64, filters OrderFilters) ([]SalesOrder, error)
	
	CreateOrderItem(ctx context.Context, item *SalesOrderItem) error
	GetOrderItems(ctx context.Context, orderID int64) ([]SalesOrderItem, error)
}

// Service interface for sales order business logic
type Service interface {
	CreateOrderFromCart(ctx context.Context, cartID int64) (*SalesOrder, error)
	GetOrder(ctx context.Context, id int64) (*SalesOrder, error)
	GetOrders(ctx context.Context, enterpriseID int64, filters OrderFilters) ([]SalesOrder, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) error
	CancelOrder(ctx context.Context, orderID int64) error
	CompleteOrder(ctx context.Context, orderID int64) error
}

// OrderFilters filters for listing orders
type OrderFilters struct {
	Status     string
	BranchID   *int64
	CustomerID *int64
	Page       int
	Limit      int
	StartDate  *time.Time
	EndDate    *time.Time
}

// Constants for order status
const (
	StatusPendingPayment = "PENDING_PAYMENT"
	StatusPaid           = "PAID"
	StatusCancelled      = "CANCELLED"
	StatusCompleted      = "COMPLETED"
)
