package compositions

import (
	"context"

	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// Composition represents a product composition (KIT/RECETA)
type Composition struct {
	ID              int64        `json:"id"`
	ParentProductID int64        `json:"parent_product_id"`
	ChildProductID  int64        `json:"child_product_id"`
	ParentName      string       `json:"parent_name,omitempty"`
	ChildName       string       `json:"child_name,omitempty"`
	Quantity        float64      `json:"quantity"`
	Type            string       `json:"type"`
	EnterpriseID    int64        `json:"enterprise_id"`
	CreatedAt       vo.DateTime  `json:"created_at"`
	UpdatedAt       *vo.DateTime `json:"updated_at,omitempty"`
}

// Repository interface
type Repository interface {
	Create(ctx context.Context, tenantSlug string, c *Composition) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Composition, error)
	ListByParent(ctx context.Context, tenantSlug string, parentID int64) ([]Composition, error)
	ListAll(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Composition, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, tipo string, sort string, order string) (domain.PageResult, error)
	Update(ctx context.Context, tenantSlug string, c *Composition) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
	ExistsPair(ctx context.Context, tenantSlug string, parentID int64, childID int64) (bool, error)
}

// Service interface
type Service interface {
	Create(ctx context.Context, tenantSlug string, c *Composition) error
	GetByID(ctx context.Context, tenantSlug string, id int64) (*Composition, error)
	ListByParent(ctx context.Context, tenantSlug string, parentID int64) ([]Composition, error)
	ListAll(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Composition, error)
	Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, tipo string, sort string, order string) (domain.PageResult, error)
	Update(ctx context.Context, tenantSlug string, id int64, c *Composition) error
	Delete(ctx context.Context, tenantSlug string, id int64) error
}

// Domain Events
const (
	EventCreated = "composition.created"
	EventUpdated = "composition.updated"
	EventDeleted = "composition.deleted"
)

type CreatedEvent struct{ events.BaseEvent }
type UpdatedEvent struct{ events.BaseEvent }
type DeletedEvent struct{ events.BaseEvent }

func (e *Composition) ToEventPayload() map[string]interface{} {
	return map[string]interface{}{
		"id":                e.ID,
		"parent_product_id": e.ParentProductID,
		"child_product_id":  e.ChildProductID,
		"quantity":          e.Quantity,
		"type":              e.Type,
	}
}

func NewCreatedEvent(e *Composition) CreatedEvent {
	return CreatedEvent{events.NewBaseEvent(EventCreated, e.ToEventPayload())}
}

func NewUpdatedEvent(e *Composition) UpdatedEvent {
	return UpdatedEvent{events.NewBaseEvent(EventUpdated, e.ToEventPayload())}
}

func NewDeletedEvent(e *Composition) DeletedEvent {
	return DeletedEvent{events.NewBaseEvent(EventDeleted, e.ToEventPayload())}
}
