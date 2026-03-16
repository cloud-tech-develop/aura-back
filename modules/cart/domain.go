package cart

import (
	"context"
	"time"
)

type Cart struct {
	ID         int64      `json:"id"`
	CartCode   string     `json:"cart_code"`
	CustomerID *int64     `json:"customer_id,omitempty"`
	UserID     int64      `json:"user_id"`
	BranchID   int64      `json:"branch_id"`
	EmpresaID  int64      `json:"empresa_id"`
	Subtotal   float64    `json:"subtotal"`
	Discount   float64    `json:"discount"`
	TaxTotal   float64    `json:"tax_total"`
	Total      float64    `json:"total"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

type CartItem struct {
	ID              int64   `json:"id"`
	CartID          int64   `json:"cart_id"`
	ProductID       int64   `json:"product_id"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	DiscountPercent float64 `json:"discount_percent"`
	DiscountAmount  float64 `json:"discount_amount"`
	TaxRate         float64 `json:"tax_rate"`
	TaxAmount       float64 `json:"tax_amount"`
	Total           float64 `json:"total"`
}

type Repository interface {
	CreateCart(ctx context.Context, cart *Cart) error
	GetCartByID(ctx context.Context, id int64) (*Cart, error)
	AddItem(ctx context.Context, item *CartItem) error
	UpdateItem(ctx context.Context, item *CartItem) error
	RemoveItem(ctx context.Context, cartID, itemID int64) error
	GetItems(ctx context.Context, cartID int64) ([]CartItem, error)
	UpdateCartTotals(ctx context.Context, cartID int64) error
}

type Service interface {
	CreateCart(ctx context.Context, cart *Cart) error
	GetCart(ctx context.Context, id int64) (*Cart, error)
	AddItem(ctx context.Context, cartID int64, item *CartItem) error
	UpdateItem(ctx context.Context, cartID, itemID int64, quantity int) error
	RemoveItem(ctx context.Context, cartID, itemID int64) error
	ConvertToSale(ctx context.Context, cartID int64) error // HU-SALES-002
}
