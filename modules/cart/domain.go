package cart

import (
	"context"
	"time"
)

// Cart represents a shopping cart or quotation
type Cart struct {
	ID            int64       `json:"id"`
	CartCode      string      `json:"cart_code"`
	CartType      string      `json:"cart_type"` // SALE, QUOTATION
	CustomerID    *int64      `json:"customer_id,omitempty"`
	UserID        int64       `json:"user_id"`
	BranchID      int64       `json:"branch_id"`
	EnterpriseID  int64       `json:"enterprise_id"`
	Subtotal      float64     `json:"subtotal"`
	Discount      float64     `json:"discount"`
	TaxTotal      float64     `json:"tax_total"`
	Total         float64     `json:"total"`
	Status        string      `json:"status"`
	Notes         string      `json:"notes"`
	ValidUntil    *time.Time  `json:"valid_until,omitempty"`
	ConvertedAt   *time.Time  `json:"converted_at,omitempty"`
	ReferenceID   *int64      `json:"reference_id,omitempty"`
	ReferenceType string      `json:"reference_type,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     *time.Time  `json:"updated_at,omitempty"`
	DeletedAt     *time.Time  `json:"deleted_at,omitempty"`
	Items         []CartItem  `json:"items,omitempty"`
}

// CartItem represents an item in the cart
type CartItem struct {
	ID                int64       `json:"id"`
	CartID            int64       `json:"cart_id"`
	ProductID         int64       `json:"product_id"`
	ProductVariantID  *int64      `json:"product_variant_id,omitempty"`
	Quantity          int         `json:"quantity"`
	UnitPrice         float64     `json:"unit_price"`
	DiscountType      string      `json:"discount_type"` // PERCENTAGE, FIXED
	DiscountValue     float64     `json:"discount_value"`
	DiscountAmount    float64     `json:"discount_amount"`
	TaxRate           float64     `json:"tax_rate"`
	TaxAmount         float64     `json:"tax_amount"`
	LineTotal         float64     `json:"line_total"`
	Notes             string      `json:"notes"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         *time.Time  `json:"updated_at,omitempty"`
}

// Cart status constants
const (
	CartStatusActive    = "ACTIVE"
	CartStatusSaved     = "SAVED"
	CartStatusConverted = "CONVERTED"
	CartStatusExpired   = "EXPIRED"
	CartStatusCancelled = "CANCELLED"
)

// Cart type constants
const (
	CartTypeSale       = "SALE"
	CartTypeQuotation  = "QUOTATION"
)

// Discount type constants
const (
	DiscountTypePercentage = "PERCENTAGE"
	DiscountTypeFixed      = "FIXED"
)

// CartFilters filters for listing carts
type CartFilters struct {
	CartType  string // "sale", "quotation"
	Status    string
	CustomerID *int64
	BranchID   *int64
	Page       int
	Limit      int
}

// Repository interface for cart operations
type Repository interface {
	CreateCart(ctx context.Context, cart *Cart) error
	GetCartByID(ctx context.Context, id int64) (*Cart, error)
	GetCartByCode(ctx context.Context, code string, enterpriseID int64) (*Cart, error)
	UpdateCart(ctx context.Context, cart *Cart) error
	UpdateCartStatus(ctx context.Context, cartID int64, status string) error
	DeleteCart(ctx context.Context, id int64) error
	ListCarts(ctx context.Context, enterpriseID int64, filters CartFilters) ([]Cart, error)
	
	AddItem(ctx context.Context, item *CartItem) error
	UpdateItem(ctx context.Context, item *CartItem) error
	RemoveItem(ctx context.Context, cartID, itemID int64) error
	GetItems(ctx context.Context, cartID int64) ([]CartItem, error)
	UpdateCartTotals(ctx context.Context, cartID int64) error
}

// Service interface for cart business logic
type Service interface {
	CreateCart(ctx context.Context, cart *Cart) error
	GetCart(ctx context.Context, id int64) (*Cart, error)
	GetCartByCode(ctx context.Context, code string, enterpriseID int64) (*Cart, error)
	UpdateCart(ctx context.Context, id int64, cart *Cart) error
	DeleteCart(ctx context.Context, id int64) error
	ListCarts(ctx context.Context, enterpriseID int64, filters CartFilters) ([]Cart, error)
	
	AddItem(ctx context.Context, cartID int64, item *CartItem) error
	UpdateItem(ctx context.Context, cartID, itemID int64, item *CartItem) error
	RemoveItem(ctx context.Context, cartID, itemID int64) error
	GetItems(ctx context.Context, cartID int64) ([]CartItem, error)
	
	ConvertToSale(ctx context.Context, cartID int64) (*Cart, error)
	ConvertToQuotation(ctx context.Context, cartID int64, validDays int) (*Cart, error)
	SetCustomer(ctx context.Context, cartID int64, customerID *int64) error
	ApplyDiscount(ctx context.Context, cartID int64, discountType string, discountValue float64) error
	ApplyItemDiscount(ctx context.Context, cartID, itemID int64, discountType string, discountValue float64) error
}
