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
	Name           string                 `json:"name"`
	CommercialName string                 `json:"commercial_name"`
	Slug           string                 `json:"slug"`
	SubDomain      string                 `json:"sub_domain"`
	Email          string                 `json:"email"`
	Document       string                 `json:"document"`
	DV             string                 `json:"dv"`
	Phone          string                 `json:"phone"`
	MunicipalityID string                 `json:"municipality_id"`
	Municipality   string                 `json:"municipality"`
	Status         string                 `json:"status"`
	Settings       map[string]interface{} `json:"settings"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      *time.Time             `json:"deleted_at,omitempty"`
}

// Plan represents a subscription plan for an enterprise
type Plan struct {
	ID             int64      `json:"id"`
	EnterpriseID   int64      `json:"enterprise_id"`
	MaxUsers       *int       `json:"max_users,omitempty"`
	MaxEnterprises *int       `json:"max_enterprises,omitempty"`
	TrialUntil     *time.Time `json:"trial_until,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// User represents a user from public.users
type User struct {
	ID           int64      `json:"id"`
	EnterpriseID int64      `json:"enterprise_id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Active       bool       `json:"active"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// UserWithPassword represents a user with password hash (for offline sync)
type UserWithPassword struct {
	ID           int64      `json:"id"`
	EnterpriseID int64      `json:"enterprise_id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	Active       bool       `json:"active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// UserRole represents a user-role assignment from public.user_roles
type UserRole struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

// ─── Tenant Entities ───────────────────────────────────────────────────────

// ThirdParty represents a third party from tenant schema
type ThirdParty struct {
	ID                int64      `json:"id"`
	UserID            int64      `json:"user_id"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	DocumentNumber    string     `json:"document_number"`
	DocumentType      string     `json:"document_type"`
	PersonalEmail     string     `json:"personal_email"`
	CommercialName    string     `json:"commercial_name"`
	Address           string     `json:"address"`
	Phone             string     `json:"phone"`
	AdditionalEmail   string     `json:"additional_email"`
	TaxResponsibility string     `json:"tax_responsibility"`
	IsClient          bool       `json:"is_client"`
	IsProvider        bool       `json:"is_provider"`
	IsEmployee        bool       `json:"is_employee"`
	MunicipalityID    string     `json:"municipality_id"`
	Municipality      string     `json:"municipality"`
	EnterpriseID      int64      `json:"enterprise_id"`
	CreatedAt         time.Time  `json:"created_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	Enterprises   int      `json:"enterprises"`
	Plans         int      `json:"plans"`
	Users         int      `json:"users"`
	UserRoles     int      `json:"user_roles"`
	ThirdParties  int      `json:"third_parties"`
	Categories    int      `json:"categories"`
	Brands        int      `json:"brands"`
	Units         int      `json:"units"`
	Products      int      `json:"products"`
	Presentations int      `json:"presentations"`
	Errors        []string `json:"errors,omitempty"`
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
}

// ─── Service Interface ────────────────────────────────────────────────────────────

type Service interface {
	SyncAllBySlug(ctx context.Context, prodURL, token, slug string) (*SyncResult, error)
	SyncTenantBySlug(ctx context.Context, prodURL, token, slug string) (*SyncResult, error)
	GetLocalEnterprises(ctx context.Context) ([]Enterprise, error)
}
