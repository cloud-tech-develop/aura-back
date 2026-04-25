package offline

import (
	"context"
	"time"
)

// ─── Entities ───────────────────────────────────────────────────────────────────

// Enterprise represents an enterprise synced from production
type Enterprise struct {
	ID             int64                  `json:"id"`
	TenantID       int64                  `json:"tenant_id"`
	Name          string                 `json:"name"`
	CommercialName string               `json:"commercial_name"`
	Slug          string                 `json:"slug"`
	SubDomain     string                 `json:"sub_domain"`
	Email         string                 `json:"email"`
	Document      string                 `json:"document"`
	DV            string                 `json:"dv"`
	Phone         string                 `json:"phone"`
	MunicipalityID string               `json:"municipality_id"`
	Municipality  string                 `json:"municipality"`
	Status        string                 `json:"status"`
	Settings      map[string]interface{} `json:"settings,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	DeletedAt     *time.Time             `json:"deleted_at,omitempty"`
}

// Plan represents a subscription plan for an enterprise
type Plan struct {
	ID            int64       `json:"id"`
	EnterpriseID int64       `json:"enterprise_id"`
	MaxUsers      *int        `json:"max_users,omitempty"`
	MaxEnterprises *int      `json:"max_enterprises,omitempty"`
	TrialUntil    *time.Time `json:"trial_until,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

// User represents a user from public.users
type User struct {
	ID            int64     `json:"id"`
	EnterpriseID int64     `json:"enterprise_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Active       bool      `json:"active"`
	PasswordHash string   `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// UserWithPassword represents a user with password hash (for offline sync)
type UserWithPassword struct {
	ID            int64     `json:"id"`
	EnterpriseID int64     `json:"enterprise_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// UserRole represents a user-role assignment from public.user_roles
type UserRole struct {
	ID       int64 `json:"id"`
	UserID  int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

// ─── Tenant Entities ───────────────────────────────────────────────────────

// ThirdParty represents a third party from tenant schema
type ThirdParty struct {
	ID                  int64     `json:"id"`
	UserID              int64     `json:"user_id"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	DocumentNumber     string    `json:"document_number"`
	DocumentType       string    `json:"document_type"`
	PersonalEmail      string    `json:"personal_email"`
	CommercialName   string    `json:"commercial_name"`
	Address           string    `json:"address"`
	Phone             string    `json:"phone"`
	AdditionalEmail   string    `json:"additional_email"`
	TaxResponsibility string    `json:"tax_responsibility"`
	IsClient          bool      `json:"is_client"`
	IsProvider       bool      `json:"is_provider"`
	IsEmployee       bool      `json:"is_employee"`
	MunicipalityID   string    `json:"municipality_id"`
	Municipality     string    `json:"municipality"`
	EnterpriseID     int64     `json:"enterprise_id"`
	CreatedAt        time.Time `json:"created_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// Category represents a product category
type Category struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Description   string    `json:"description"`
	ParentID      *int      `json:"parent_id"`
	DefaultTaxRate float64  `json:"default_tax_rate"`
	Active        bool      `json:"active"`
	EnterpriseID  int64     `json:"enterprise_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

// Brand represents a product brand
type Brand struct {
	ID           int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	EnterpriseID int64     `json:"enterprise_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// Unit represents a measurement unit
type Unit struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Abbreviation  string    `json:"abbreviation"`
	Active        bool      `json:"active"`
	AllowDecimals bool      `json:"allow_decimals"`
	EnterpriseID  int64     `json:"enterprise_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

// Product represents a product
type Product struct {
	ID                int64     `json:"id"`
	SKU               string    `json:"sku"`
	Barcode           string    `json:"barcode"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	CategoryID        *int64    `json:"category_id"`
	BrandID          *int64    `json:"brand_id"`
	UnitID           int64     `json:"unit_id"`
	ProductType       string    `json:"product_type"`
	Active           bool      `json:"active"`
	VisibleInPOS      bool      `json:"visible_in_pos"`
	CostPrice       float64   `json:"cost_price"`
	SalePrice       float64   `json:"sale_price"`
	Price2         *float64  `json:"price_2"`
	Price3         *float64  `json:"price_3"`
	IVAPercentage   float64   `json:"iva_percentage"`
	ConsumptionTax  float64   `json:"consumption_tax_value"`
	CurrentStock   int       `json:"current_stock"`
	MinStock       int       `json:"min_stock"`
	MaxStock       int       `json:"max_stock"`
	ManagesInventory bool   `json:"manages_inventory"`
	ManagesBatches  bool   `json:"manages_batches"`
	ManagesSerial    bool    `json:"manages_serial"`
	AllowNegativeStock bool   `json:"allow_negative_stock"`
	ImageURL       string    `json:"image_url"`
	EnterpriseID   int64     `json:"enterprise_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// Presentation represents a product presentation
type Presentation struct {
	ID              int64     `json:"id"`
	ProductID       int64     `json:"product_id"`
	Name           string    `json:"name"`
	Factor         float64   `json:"factor"`
	Barcode        string    `json:"barcode"`
	CostPrice      float64   `json:"cost_price"`
	SalePrice      float64   `json:"sale_price"`
	DefaultPurchase bool    `json:"default_purchase"`
	DefaultSale   bool      `json:"default_sale"`
	EnterpriseID   int64     `json:"enterprise_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	Enterprises   int      `json:"enterprises"`
	Plans       int      `json:"plans"`
	Users       int      `json:"users"`
	UserRoles   int      `json:"user_roles"`
	ThirdParties int      `json:"third_parties"`
	Categories  int      `json:"categories"`
	Brands     int      `json:"brands"`
	Units      int      `json:"units"`
	Products   int      `json:"products"`
	Presentations int     `json:"presentations"`
	Errors     []string `json:"errors,omitempty"`
}

// ─── Repository Interface ─────────────────────────────────────────────────────

type Repository interface {
	// Public schema operations
	UpsertEnterprise(ctx context.Context, e *Enterprise) error
	GetEnterpriseBySlug(ctx context.Context, slug string) (*Enterprise, error)
	ListEnterprises(ctx context.Context) ([]Enterprise, error)
	UpsertPlan(ctx context.Context, p *Plan) error
	UpsertUser(ctx context.Context, u *User) error
	UpsertUserRole(ctx context.Context, ur *UserRole) error

	// Tenant schema operations
	UpsertThirdParty(ctx context.Context, tp *ThirdParty) error
	UpsertCategory(ctx context.Context, c *Category) error
	UpsertBrand(ctx context.Context, b *Brand) error
	UpsertUnit(ctx context.Context, u *Unit) error
	UpsertProduct(ctx context.Context, p *Product) error
	UpsertPresentation(ctx context.Context, pr *Presentation) error
}

// ─── Service Interface ────────────────────────────────────────────────────────────

type Service interface {
	SyncAllBySlug(ctx context.Context, prodURL, token, slug string) (*SyncResult, error)
	SyncTenantBySlug(ctx context.Context, prodURL, token, slug string) (*SyncResult, error)
	GetLocalEnterprises(ctx context.Context) ([]Enterprise, error)
}