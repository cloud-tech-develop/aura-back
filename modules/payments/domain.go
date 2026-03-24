package payments

import (
	"context"
	"time"
)

// Payment represents a payment transaction
type Payment struct {
	ID                 int64       `json:"id"`
	PaymentType        string      `json:"payment_type"` // SALE, PURCHASE, ACCOUNT_RECEIVABLE, ACCOUNT_PAYABLE
	ReferenceID        int64       `json:"reference_id"`
	ReferenceType      string      `json:"reference_type"`
	PaymentMethod      string      `json:"payment_method"`
	Amount             float64     `json:"amount"`
	ReferenceNumber    string      `json:"reference_number"`
	BankName           string      `json:"bank_name"`
	CardType           string      `json:"card_type"`
	CardLastDigits     string      `json:"card_last_digits"`
	AuthorizationCode  string      `json:"authorization_code"`
	ChangeAmount       float64     `json:"change_amount"`
	CashDrawerID       *int64      `json:"cash_drawer_id,omitempty"`
	BranchID           int64       `json:"branch_id"`
	EnterpriseID       int64       `json:"enterprise_id"`
	UserID             int64       `json:"user_id"`
	Notes              string      `json:"notes"`
	Status             string      `json:"status"`
	CancelledAt        *time.Time  `json:"cancelled_at,omitempty"`
	CancelledBy        *int64      `json:"cancelled_by,omitempty"`
	CancellationReason string      `json:"cancellation_reason,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          *time.Time  `json:"updated_at,omitempty"`
}

// PaymentTransaction represents a payment transaction log
type PaymentTransaction struct {
	ID                int64      `json:"id"`
	PaymentID         int64      `json:"payment_id"`
	TransactionType   string     `json:"transaction_type"` // CHARGE, REFUND, CHARGEBACK
	Amount            float64    `json:"amount"`
	PreviousBalance   *float64   `json:"previous_balance,omitempty"`
	NewBalance        *float64   `json:"new_balance,omitempty"`
	ProcessorReference string    `json:"processor_reference,omitempty"`
	ProcessorResponse string     `json:"processor_response,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

// CashDrawer represents a cash register drawer
type CashDrawer struct {
	ID             int64       `json:"id"`
	UserID         int64       `json:"user_id"`
	BranchID       int64       `json:"branch_id"`
	EnterpriseID   int64       `json:"enterprise_id"`
	OpeningBalance float64     `json:"opening_balance"`
	ClosingBalance *float64    `json:"closing_balance,omitempty"`
	CashIn         float64     `json:"cash_in"`
	CashOut        float64     `json:"cash_out"`
	Status         string      `json:"status"`
	OpenedAt       time.Time   `json:"opened_at"`
	ClosedAt       *time.Time  `json:"closed_at,omitempty"`
	Notes          string      `json:"notes"`
}

// CashMovement represents a cash movement in drawer
type CashMovement struct {
	ID            int64     `json:"id"`
	CashDrawerID  int64     `json:"cash_drawer_id"`
	MovementType  string    `json:"movement_type"` // IN, OUT
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	UserID        int64     `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// Payment methods
const (
	MethodCash         = "CASH"
	MethodDebitCard    = "DEBIT_CARD"
	MethodCreditCard   = "CREDIT_CARD"
	MethodBankTransfer = "BANK_TRANSFER"
	MethodCredit       = "CREDIT"
	MethodVoucher      = "VOUCHER"
	MethodCheck        = "CHECK"
)

// Payment types
const (
	PaymentTypeSale              = "SALE"
	PaymentTypePurchase         = "PURCHASE"
	PaymentTypeAccountReceivable = "ACCOUNT_RECEIVABLE"
	PaymentTypeAccountPayable   = "ACCOUNT_PAYABLE"
)

// Transaction types
const (
	TransactionTypeCharge    = "CHARGE"
	TransactionTypeRefund    = "REFUND"
	TransactionTypeChargeback = "CHARGEBACK"
)

// Payment status
const (
	PaymentStatusCompleted = "COMPLETED"
	PaymentStatusPending   = "PENDING"
	PaymentStatusCancelled = "CANCELLED"
	PaymentStatusRefunded  = "REFUNDED"
)

// Cash drawer status
const (
	DrawerStatusOpen   = "OPEN"
	DrawerStatusClosed = "CLOSED"
)

// Cash movement types
const (
	CashMovementIn  = "IN"
	CashMovementOut = "OUT"
)

// PaymentFilters filters for listing payments
type PaymentFilters struct {
	ReferenceID   *int64
	PaymentMethod string
	Status        string
	StartDate     *time.Time
	EndDate       *time.Time
	Page          int
	Limit         int
}

// Repository interface for payment operations
type Repository interface {
	CreatePayment(ctx context.Context, p *Payment) error
	GetPaymentByID(ctx context.Context, id int64) (*Payment, error)
	GetPaymentsByReference(ctx context.Context, referenceType string, referenceID int64) ([]Payment, error)
	UpdatePaymentStatus(ctx context.Context, id int64, status string) error
	CancelPayment(ctx context.Context, id int64, cancelledBy int64, reason string) error
	ListPayments(ctx context.Context, enterpriseID int64, filters PaymentFilters) ([]Payment, error)
	
	CreateTransaction(ctx context.Context, tx *PaymentTransaction) error
	GetTransactionsByPayment(ctx context.Context, paymentID int64) ([]PaymentTransaction, error)
	
	CreateCashDrawer(ctx context.Context, drawer *CashDrawer) error
	GetCashDrawerByID(ctx context.Context, id int64) (*CashDrawer, error)
	GetOpenCashDrawer(ctx context.Context, userID, branchID int64) (*CashDrawer, error)
	UpdateCashDrawer(ctx context.Context, drawer *CashDrawer) error
	ListCashDrawers(ctx context.Context, enterpriseID int64, userID *int64, status string) ([]CashDrawer, error)
	
	CreateCashMovement(ctx context.Context, mov *CashMovement) error
	GetCashMovements(ctx context.Context, drawerID int64) ([]CashMovement, error)
}

// Service interface for payment business logic
type Service interface {
	ProcessPayment(ctx context.Context, p *Payment) error
	ProcessMultiplePayments(ctx context.Context, payments []Payment) error
	CalculateChange(ctx context.Context, amount float64, paymentMethod string) float64
	GetPayment(ctx context.Context, id int64) (*Payment, error)
	GetPaymentsByOrder(ctx context.Context, referenceType string, referenceID int64) ([]Payment, error)
	CancelPayment(ctx context.Context, id int64, cancelledBy int64, reason string) error
	ListPayments(ctx context.Context, enterpriseID int64, filters PaymentFilters) ([]Payment, error)
	
	OpenCashDrawer(ctx context.Context, drawer *CashDrawer) error
	CloseCashDrawer(ctx context.Context, drawerID int64, closingBalance float64, notes string) error
	GetCashDrawer(ctx context.Context, drawerID int64) (*CashDrawer, error)
	GetOpenDrawer(ctx context.Context, userID, branchID int64) (*CashDrawer, error)
	AddCashIn(ctx context.Context, drawerID int64, amount float64, description string) error
	AddCashOut(ctx context.Context, drawerID int64, amount float64, description string) error
	ListCashDrawers(ctx context.Context, enterpriseID int64, userID *int64, status string) ([]CashDrawer, error)
}
