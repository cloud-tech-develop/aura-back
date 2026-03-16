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
	ID           int64
	EnterpriseID int64
	Name         string // Maps to public.users.name
	Email        vo.Email
	PasswordHash string
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time

	// Third party fields (not stored in public.users, but used for tenant.third_parties)
	FirstName         string
	LastName          string
	DocumentNumber    string
	DocumentType      string
	PersonalEmail     string
	TaxResponsibility string
	IsEmployee        bool
}

// Repository defines the data access operations for users.
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email vo.Email) (*User, error)
	ListByEnterprise(ctx context.Context, enterpriseID int64, limit, offset int) ([]User, error)
	Update(ctx context.Context, user *User) error
	UpdateStatus(ctx context.Context, id int64, active bool) error
}

// Service defines the business logic for users.
type Service interface {
	Create(ctx context.Context, tenantSlug string, user *User, password string) error
	GetByID(ctx context.Context, id int64) (*User, error)
	ListByEnterprise(ctx context.Context, enterpriseID int64, page, limit int) ([]User, error)
	Update(ctx context.Context, user *User) error
	UpdateStatus(ctx context.Context, id int64, active bool) error
	AssignRoles(ctx context.Context, userID int64, roleIDs []int64) error
}
