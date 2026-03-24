package inventory

import (
	"context"
	"time"
)

// Inventory represents stock of a product at a branch
type Inventory struct {
	ID                int64      `json:"id"`
	ProductID         int64      `json:"product_id"`
	BranchID          int64      `json:"branch_id"`
	Quantity          int        `json:"quantity"`
	ReservedQuantity  int        `json:"reserved_quantity"`
	MinStock          int        `json:"min_stock"`
	MaxStock          *int       `json:"max_stock,omitempty"`
	Location          string     `json:"location"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}

// AvailableQuantity returns quantity - reserved
func (i *Inventory) AvailableQuantity() int {
	return i.Quantity - i.ReservedQuantity
}

// IsLowStock returns true if quantity is at or below min_stock
func (i *Inventory) IsLowStock() bool {
	return i.Quantity <= i.MinStock
}

// InventoryMovement represents a stock movement (Kardex entry)
type InventoryMovement struct {
	ID               int64      `json:"id"`
	InventoryID      int64      `json:"inventory_id"`
	MovementType     string     `json:"movement_type"` // ENTRY, EXIT, ADJUSTMENT
	MovementReason   string     `json:"movement_reason"`
	Quantity         int        `json:"quantity"`
	PreviousBalance  int        `json:"previous_balance"`
	NewBalance       int        `json:"new_balance"`
	ReferenceID      *int64     `json:"reference_id,omitempty"`
	ReferenceType    string     `json:"reference_type,omitempty"`
	BatchNumber      *string    `json:"batch_number,omitempty"`
	SerialNumber     *string    `json:"serial_number,omitempty"`
	ExpirationDate   *time.Time `json:"expiration_date,omitempty"`
	Notes            string     `json:"notes"`
	UserID           int64      `json:"user_id"`
	BranchID         int64      `json:"branch_id"`
	CreatedAt        time.Time  `json:"created_at"`
}

// MovementReason represents a type of movement
type MovementReason struct {
	ID                  int64     `json:"id"`
	Code                string    `json:"code"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	MovementType        string    `json:"movement_type"` // ENTRY, EXIT, ADJUSTMENT
	RequiresAuthorization bool    `json:"requires_authorization"`
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
}

// Movement types
const (
	MovementTypeEntry     = "ENTRY"
	MovementTypeExit      = "EXIT"
	MovementTypeAdjustment = "ADJUSTMENT"
)

// Movement reasons
const (
	MovementReasonSale          = "SALE"
	MovementReasonPurchase     = "PURCHASE"
	MovementReasonShrinkage    = "SHRINKAGE"
	MovementReasonTransferIn   = "TRANSFER_IN"
	MovementReasonTransferOut = "TRANSFER_OUT"
	MovementReasonAdjustment   = "ADJUSTMENT"
	MovementReasonReturn       = "RETURN"
	MovementReasonInitial      = "INITIAL"
	MovementReasonDamage      = "DAMAGE"
	MovementReasonTheft       = "THEFT"
	MovementReasonExpired     = "EXPIRED"
)

// InventoryFilters filters for listing inventory
type InventoryFilters struct {
	BranchID  *int64
	ProductID *int64
	LowStock  bool
	Page      int
	Limit     int
}

// MovementFilters filters for listing movements
type MovementFilters struct {
	InventoryID   *int64
	BranchID      *int64
	ProductID     *int64
	MovementType  string
	MovementReason string
	StartDate     *time.Time
	EndDate       *time.Time
	Page          int
	Limit         int
}

// Repository interface for inventory operations
type Repository interface {
	// Inventory CRUD
	CreateInventory(ctx context.Context, inv *Inventory) error
	GetInventory(ctx context.Context, productID, branchID int64) (*Inventory, error)
	UpdateInventory(ctx context.Context, inv *Inventory) error
	ListInventory(ctx context.Context, enterpriseID int64, filters InventoryFilters) ([]Inventory, error)
	GetLowStock(ctx context.Context, enterpriseID int64) ([]Inventory, error)
	GetInventoryByProduct(ctx context.Context, productID int64) ([]Inventory, error)

	// Movements
	CreateMovement(ctx context.Context, mov *InventoryMovement) error
	GetMovement(ctx context.Context, id int64) (*InventoryMovement, error)
	ListMovements(ctx context.Context, enterpriseID int64, filters MovementFilters) ([]InventoryMovement, error)
	GetProductKardex(ctx context.Context, productID int64, branchID int64) ([]InventoryMovement, error)

	// Movement Reasons
	ListReasons(ctx context.Context) ([]MovementReason, error)
	GetReason(ctx context.Context, code string) (*MovementReason, error)
}

// Service interface for inventory business logic
type Service interface {
	// Inventory
	GetInventory(ctx context.Context, productID, branchID int64) (*Inventory, error)
	ListInventory(ctx context.Context, enterpriseID int64, filters InventoryFilters) ([]Inventory, error)
	GetLowStock(ctx context.Context, enterpriseID int64) ([]Inventory, error)
	GetInventoryByProduct(ctx context.Context, productID int64) ([]Inventory, error)
	UpdateStock(ctx context.Context, productID, branchID int64, quantity int, reason, referenceType string, referenceID *int64, userID int64, notes string) (*Inventory, error)

	// Movements
	GetMovement(ctx context.Context, id int64) (*InventoryMovement, error)
	ListMovements(ctx context.Context, enterpriseID int64, filters MovementFilters) ([]InventoryMovement, error)
	GetProductKardex(ctx context.Context, productID int64, branchID int64) ([]InventoryMovement, error)

	// Reasons
	ListReasons(ctx context.Context) ([]MovementReason, error)
}
