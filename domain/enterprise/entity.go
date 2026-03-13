package enterprise

import (
	"time"

	"github.com/cloud-tech-develop/aura-back/domain/events"
)

const (
	EventEnterpriseCreated = "enterprise.created"
	EventEnterpriseUpdated = "enterprise.updated"
	EventEnterpriseDeleted = "enterprise.deleted"
)

type Enterprise struct {
	ID             int64
	TenantID       int64  `json:"tenant_id"`
	Name           string `json:"name"`
	CommercialName string `json:"commercial_name"`
	Slug           string
	SubDomain      string `json:"sub_domain"`
	Email          string
	DV             string `json:"dv"`
	Phone          string `json:"phone"`
	MunicipalityID string `json:"municipality_id"`
	Municipality   string `json:"municipality"`
	Status         string `json:"status"`
	Settings       map[string]interface{}
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

func (e *Enterprise) ToEventsPayload() map[string]interface{} {
	return map[string]interface{}{
		"id":              e.ID,
		"name":            e.Name,
		"commercial_name": e.CommercialName,
		"slug":            e.Slug,
		"sub_domain":      e.SubDomain,
		"email":           e.Email,
		"dv":              e.DV,
		"phone":           e.Phone,
		"municipality_id": e.MunicipalityID,
		"municipality":    e.Municipality,
		"status":          e.Status,
		"created_at":      e.CreatedAt,
		"updated_at":      e.UpdatedAt,
	}
}

type EnterpriseCreatedEvent struct {
	events.BaseEvent
}

func NewEnterpriseCreatedEvent(e *Enterprise) EnterpriseCreatedEvent {
	return EnterpriseCreatedEvent{
		BaseEvent: events.NewBaseEvent(EventEnterpriseCreated, e.ToEventsPayload()),
	}
}

type EnterpriseUpdatedEvent struct {
	events.BaseEvent
}

func NewEnterpriseUpdatedEvent(e *Enterprise) EnterpriseUpdatedEvent {
	return EnterpriseUpdatedEvent{
		BaseEvent: events.NewBaseEvent(EventEnterpriseUpdated, e.ToEventsPayload()),
	}
}

type EnterpriseDeletedEvent struct {
	events.BaseEvent
}

func NewEnterpriseDeletedEvent(e *Enterprise) EnterpriseDeletedEvent {
	return EnterpriseDeletedEvent{
		BaseEvent: events.NewBaseEvent(EventEnterpriseDeleted, e.ToEventsPayload()),
	}
}
