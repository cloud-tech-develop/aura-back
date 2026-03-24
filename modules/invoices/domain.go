package invoices

import (
	"context"
	"time"
)

// InvoicePrefix represents a prefix for invoice numbering
type InvoicePrefix struct {
	ID              int64      `json:"id"`
	Prefix          string     `json:"prefix"`
	BranchID        int64      `json:"branch_id"`
	EnterpriseID    int64      `json:"enterprise_id"`
	CurrentNumber   int64      `json:"current_number"`
	ResolutionNumber string    `json:"resolution_number,omitempty"`
	ResolutionDate  *time.Time `json:"resolution_date,omitempty"`
	ValidFrom       *time.Time `json:"valid_from,omitempty"`
	ValidUntil      *time.Time `json:"valid_until,omitempty"`
	Description     string     `json:"description,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

// Invoice represents an invoice
type Invoice struct {
	ID              int64        `json:"id"`
	InvoiceNumber   string       `json:"invoice_number"`
	PrefixID        int64        `json:"prefix_id"`
	InvoiceType    string       `json:"invoice_type"` // SALE, CREDIT_NOTE, DEBIT_NOTE
	ReferenceID     *int64       `json:"reference_id,omitempty"`
	ReferenceType   string       `json:"reference_type,omitempty"`
	SalesOrderID    *int64       `json:"sales_order_id,omitempty"`
	CustomerID      int64        `json:"customer_id"`
	BranchID        int64        `json:"branch_id"`
	UserID          int64        `json:"user_id"`
	EnterpriseID    int64        `json:"enterprise_id"`
	InvoiceDate     time.Time    `json:"invoice_date"`
	DueDate         *time.Time   `json:"due_date,omitempty"`
	Subtotal        float64      `json:"subtotal"`
	DiscountTotal   float64      `json:"discount_total"`
	TaxExempt       float64      `json:"tax_exempt"`
	TaxableAmount   float64      `json:"taxable_amount"`
	Iva19           float64      `json:"iva_19"`
	Iva5            float64      `json:"iva_5"`
	Reteica         float64      `json:"reteica"`
	Retefuente      float64      `json:"retefuente"`
	ReteicaRate     float64      `json:"reteica_rate"`
	RetefuenteRate  float64      `json:"retefuente_rate"`
	Total           float64      `json:"total"`
	PaymentMethod   string       `json:"payment_method"`
	Status          string       `json:"status"`
	Notes           string       `json:"notes"`
	CancelledAt     *time.Time   `json:"cancelled_at,omitempty"`
	CancelledBy     *int64       `json:"cancelled_by,omitempty"`
	CancellationReason string    `json:"cancellation_reason,omitempty"`
	CreditNoteID    *int64       `json:"credit_note_id,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       *time.Time   `json:"updated_at,omitempty"`
	DeletedAt       *time.Time   `json:"deleted_at,omitempty"`
	Items           []InvoiceItem `json:"items,omitempty"`
	Customer        *ThirdParty  `json:"customer,omitempty"`
}

// InvoiceItem represents an item in an invoice
type InvoiceItem struct {
	ID             int64     `json:"id"`
	InvoiceID      int64     `json:"invoice_id"`
	ProductID      int64     `json:"product_id"`
	ProductName    string    `json:"product_name"`
	ProductSKU     string    `json:"product_sku,omitempty"`
	Quantity       float64   `json:"quantity"`
	UnitPrice      float64   `json:"unit_price"`
	DiscountAmount float64   `json:"discount_amount"`
	TaxRate        float64   `json:"tax_rate"`
	TaxAmount      float64   `json:"tax_amount"`
	LineTotal      float64   `json:"line_total"`
	CreatedAt      time.Time `json:"created_at"`
}

// InvoiceLog represents invoice audit log
type InvoiceLog struct {
	ID        int64     `json:"id"`
	InvoiceID int64     `json:"invoice_id"`
	Action    string    `json:"action"` // CREATED, ISSUED, SENT, VIEWED, CANCELLED
	UserID    int64     `json:"user_id"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

// ThirdParty represents customer info (simplified)
type ThirdParty struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	DocumentNumber string `json:"document_number"`
	DocumentType string `json:"document_type"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
}

// Invoice types
const (
	InvoiceTypeSale      = "SALE"
	InvoiceTypeCreditNote = "CREDIT_NOTE"
	InvoiceTypeDebitNote  = "DEBIT_NOTE"
)

// Invoice status
const (
	InvoiceStatusDraft     = "DRAFT"
	InvoiceStatusIssued    = "ISSUED"
	InvoiceStatusSent      = "SENT"
	InvoiceStatusViewed    = "VIEWED"
	InvoiceStatusCancelled = "CANCELLED"
)

// InvoiceFilters filters for listing invoices
type InvoiceFilters struct {
	Status       string
	InvoiceType  string
	BranchID     *int64
	CustomerID   *int64
	Page         int
	Limit        int
	StartDate    *time.Time
	EndDate      *time.Time
	Search       string
}

// Repository interface for invoice operations
type Repository interface {
	CreateInvoice(ctx context.Context, inv *Invoice) error
	GetInvoiceByID(ctx context.Context, id int64) (*Invoice, error)
	GetInvoiceByNumber(ctx context.Context, invoiceNumber string, enterpriseID int64) (*Invoice, error)
	UpdateInvoice(ctx context.Context, inv *Invoice) error
	UpdateInvoiceStatus(ctx context.Context, invoiceID int64, status string) error
	CancelInvoice(ctx context.Context, invoiceID int64, cancelledBy int64, reason string) error
	GetInvoices(ctx context.Context, enterpriseID int64, filters InvoiceFilters) ([]Invoice, error)
	CountInvoices(ctx context.Context, enterpriseID int64, filters InvoiceFilters) (int, error)
	
	CreateInvoiceItem(ctx context.Context, item *InvoiceItem) error
	GetInvoiceItems(ctx context.Context, invoiceID int64) ([]InvoiceItem, error)
	
	GetInvoicePrefix(ctx context.Context, branchID int64, prefix string) (*InvoicePrefix, error)
	GetInvoicePrefixByID(ctx context.Context, id int64) (*InvoicePrefix, error)
	CreateInvoicePrefix(ctx context.Context, prefix *InvoicePrefix) error
	UpdateInvoicePrefix(ctx context.Context, prefix *InvoicePrefix) error
	ListInvoicePrefixes(ctx context.Context, enterpriseID int64) ([]InvoicePrefix, error)
	
	CreateInvoiceLog(ctx context.Context, log *InvoiceLog) error
	GetInvoiceLogs(ctx context.Context, invoiceID int64) ([]InvoiceLog, error)
}

// Service interface for invoice business logic
type Service interface {
	GenerateInvoiceFromSale(ctx context.Context, salesOrderID int64, prefixID int64) (*Invoice, error)
	GenerateInvoice(ctx context.Context, inv *Invoice) error
	GetInvoice(ctx context.Context, id int64) (*Invoice, error)
	GetInvoiceByNumber(ctx context.Context, invoiceNumber string, enterpriseID int64) (*Invoice, error)
	GetInvoices(ctx context.Context, enterpriseID int64, filters InvoiceFilters) ([]Invoice, error)
	IssueInvoice(ctx context.Context, invoiceID int64) error
	CancelInvoice(ctx context.Context, invoiceID int64, reason string) error
	
	CreateInvoicePrefix(ctx context.Context, prefix *InvoicePrefix) error
	GetInvoicePrefixes(ctx context.Context, enterpriseID int64) ([]InvoicePrefix, error)
	GetInvoicePrefix(ctx context.Context, branchID int64, prefix string) (*InvoicePrefix, error)
	
	GetInvoiceLogs(ctx context.Context, invoiceID int64) ([]InvoiceLog, error)
}
