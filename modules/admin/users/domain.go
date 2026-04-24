package users

import (
	"context"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
)

const (
	EventCreated = "user_created"
	EventUpdated = "user_updated"
	EventDeleted = "user_deleted"
)

// User represents a user entity in the public.users table.
// It also carries data for the associated third_party in the tenant schema.
type User struct {
	ID           int64      `json:"id"`
	EnterpriseID int64      `json:"enterprise_id"`
	Name         string     `json:"name"` // Maps to public.users.name
	Email        vo.Email   `json:"email"`
	PasswordHash string     `json:"-"` // no visible en json
	Active       bool       `json:"active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"-"` // no visible en json
	DeletedAt    *time.Time `json:"-"` // no visible en json

	// Third party fields (not stored in public.users, but used for tenant.third_parties)
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	DocumentNumber    string `json:"document_number"`
	DocumentType      string `json:"document_type"`
	PersonalEmail     string `json:"personal_email"`
	TaxResponsibility string `json:"tax_responsibility"`
	IsEmployee        bool   `json:"is_employee"`
}

// Role represents a role entity in the public.roles table.
type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

// UserRole represents a user-role assignment in public.user_roles table
type UserRole struct {
	ID      int64 `json:"id"`
	UserID  int64 `json:"user_id"`
	RoleID  int64 `json:"role_id"`
}

// Repository defines the data access operations for users.
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email vo.Email) (*User, error)
	ListByEnterprise(ctx context.Context, enterpriseID int64, limit, offset int) ([]User, error)
	Update(ctx context.Context, user *User) error
	UpdateStatus(ctx context.Context, id int64, active bool) error
	ListRolesByMinLevel(ctx context.Context, minLevel int) ([]Role, error)
	ListUserRolesByEnterpriseID(ctx context.Context, enterpriseID int64) ([]UserRole, error)
}

// Service defines the business logic for users.
type Service interface {
	Create(ctx context.Context, tenantSlug string, user *User, password string) error
	GetByID(ctx context.Context, id int64) (*User, error)
	ListByEnterprise(ctx context.Context, enterpriseID int64, page, limit int) ([]User, error)
	Update(ctx context.Context, user *User) error
	UpdateStatus(ctx context.Context, id int64, active bool) error
	AssignRoles(ctx context.Context, userID int64, roleIDs []int64, minLevel int) error
	ListRolesByMinLevel(ctx context.Context, minLevel int) ([]Role, error)
	ListUserRolesByEnterpriseID(ctx context.Context, enterpriseID int64) ([]UserRole, error)
}
