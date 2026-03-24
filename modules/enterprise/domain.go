package enterprise

import (
	"context"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
	EventCreated = "enterprise.created"
	EventUpdated = "enterprise.updated"
	EventDeleted = "enterprise.deleted"
)

// ─── Plan Entity ────────────────────────────────────────────────────────────────

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

// ─── Entity ───────────────────────────────────────────────────────────────────

type Enterprise struct {
	ID             int64                  `json:"id"`
	TenantID       int64                  `json:"tenant_id"`
	Name           string                 `json:"name"`
	CommercialName string                 `json:"commercial_name"`
	Slug           string                 `json:"slug"`
	SubDomain      string                 `json:"sub_domain"`
	Email          vo.Email               `json:"email"`
	Document       vo.Document            `json:"document"`
	DV             string                 `json:"dv"`
	Phone          string                 `json:"phone"`
	MunicipalityID string                 `json:"municipality_id"`
	Municipality   string                 `json:"municipality"`
	Status         string                 `json:"status"`
	Settings       map[string]interface{} `json:"settings,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      *time.Time             `json:"deleted_at,omitempty"`
}

func (e *Enterprise) ToEventPayload() map[string]interface{} {
	return map[string]interface{}{
		"id":              e.ID,
		"name":            e.Name,
		"commercial_name": e.CommercialName,
		"slug":            e.Slug,
		"sub_domain":      e.SubDomain,
		"email":           e.Email.String(),
		"status":          e.Status,
		"created_at":      e.CreatedAt,
		"updated_at":      e.UpdatedAt,
	}
}

// ─── Repository Interface ─────────────────────────────────────────────────────

// ListParams holds pagination and filter parameters
type ListParams struct {
	Page   int
	Limit  int
	Status string
}

// ListResult holds the paginated result
type ListResult struct {
	Data       []Enterprise `json:"data"`
	Pagination Pagination   `json:"pagination"`
}

// Pagination holds pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}

type Repository interface {
	Create(ctx context.Context, e *Enterprise) error
	GetBySlug(ctx context.Context, slug string) (*Enterprise, error)
	GetBySubDomain(ctx context.Context, subDomain string) (*Enterprise, error)
	GetByEmail(ctx context.Context, email vo.Email) (*Enterprise, error)
	EmailExistsInUsers(ctx context.Context, email vo.Email) (bool, error)
	List(ctx context.Context, params ListParams) (ListResult, error)
	Update(ctx context.Context, e *Enterprise) error
	Delete(ctx context.Context, id int64) error
	GetPlanByEnterpriseID(ctx context.Context, enterpriseID int64) (*Plan, error)
	CountEnterprisesByTenant(ctx context.Context, tenantID int64) (int64, error)
}

// ─── Service Interface ────────────────────────────────────────────────────────

type Service interface {
	Create(ctx context.Context, e *Enterprise, passwordHash string) error
	GetBySlug(ctx context.Context, slug string) (*Enterprise, error)
	GetBySubDomain(ctx context.Context, subDomain string) (*Enterprise, error)
	GetByEmail(ctx context.Context, email vo.Email) (*Enterprise, error)
	List(ctx context.Context, params ListParams) (ListResult, error)
	Update(ctx context.Context, e *Enterprise) error
	Delete(ctx context.Context, id int64) error
}

// ─── Migrator Interface ───────────────────────────────────────────────────────

type Migrator interface {
	RunMigrations(ctx context.Context, e *Enterprise, passwordHash string) error
}

// ─── Domain Events ────────────────────────────────────────────────────────────

type CreatedEvent struct{ events.BaseEvent }
type UpdatedEvent struct{ events.BaseEvent }
type DeletedEvent struct{ events.BaseEvent }

func NewCreatedEvent(e *Enterprise) CreatedEvent {
	return CreatedEvent{events.NewBaseEvent(EventCreated, e.ToEventPayload())}
}
func NewUpdatedEvent(e *Enterprise) UpdatedEvent {
	return UpdatedEvent{events.NewBaseEvent(EventUpdated, e.ToEventPayload())}
}
func NewDeletedEvent(e *Enterprise) DeletedEvent {
	return DeletedEvent{events.NewBaseEvent(EventDeleted, e.ToEventPayload())}
}
