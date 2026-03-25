package thirdparties

import (
	"context"
	"time"
)

// ThirdParty represents a third party entity (client, supplier, or employee)
type ThirdParty struct {
	ID                int64      `json:"id"`
	UserID            *int64     `json:"user_id,omitempty"`
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
	GlobalID          string     `json:"global_id"`
	SyncStatus        string     `json:"sync_status"`
	LastSyncedAt      *time.Time `json:"last_synced_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

// Document types
const (
	DocumentTypeCC       = "CC"       // Cédula de Ciudadanía
	DocumentTypeCE       = "CE"       // Cédula de Extranjería
	DocumentTypeNIT      = "NIT"      // Número de Identificación Tributaria
	DocumentTypePASSPORT = "PASSPORT" // Pasaporte
	DocumentTypeRUT      = "RUT"      // Registro Único Tributario
)

// Tax responsibilities
const (
	TaxRespResponsible   = "RESPONSIBLE"
	TaxRespNotResponsible = "NOT-RESPONSIBLE"
)

// ThirdPartyFilters filters for listing third parties
type ThirdPartyFilters struct {
	Type      string // "client", "provider", "employee"
	Search    string
	Status    string // "active", "inactive", "all"
	Page      int
	Limit     int
}

// Repository interface for third party operations
type Repository interface {
	Create(ctx context.Context, tp *ThirdParty) error
	GetByID(ctx context.Context, id int64) (*ThirdParty, error)
	GetByDocument(ctx context.Context, docNumber string) (*ThirdParty, error)
	Update(ctx context.Context, id int64, tp *ThirdParty) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, enterpriseID int64, filters ThirdPartyFilters) ([]ThirdParty, error)
	Count(ctx context.Context, enterpriseID int64, filters ThirdPartyFilters) (int, error)
}

// Service interface for third party business logic
type Service interface {
	Create(ctx context.Context, tp *ThirdParty) error
	GetByID(ctx context.Context, id int64) (*ThirdParty, error)
	GetByDocument(ctx context.Context, docNumber string) (*ThirdParty, error)
	Update(ctx context.Context, id int64, tp *ThirdParty) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, enterpriseID int64, filters ThirdPartyFilters) ([]ThirdParty, error)
}
