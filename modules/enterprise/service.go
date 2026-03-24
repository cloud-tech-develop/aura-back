package enterprise

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// ErrPlanLimitReached returned when enterprise limit is reached
var ErrPlanLimitReached = errors.New("ha alcanzado el límite de empresas de su plan")

var validSlug = regexp.MustCompile(`^[a-z0-9_]+$`)

// Constants for validation
const (
	MinSlugLength     = 3
	MaxSlugLength     = 50
	MinPasswordLength = 8
)

type service struct {
	repo     Repository
	eventBus events.EventBus
	migrator Migrator
	rawDB    *sql.DB // used for schema dropping
}

func NewService(db *sql.DB, eventBus events.EventBus, migrator Migrator) Service {
	return &service{
		repo:     NewRepository(db),
		eventBus: eventBus,
		migrator: migrator,
		rawDB:    db,
	}
}

func (s *service) Create(ctx context.Context, e *Enterprise, passwordHash string) error {
	e.Slug = strings.ToLower(e.Slug)
	if e.SubDomain != "" {
		e.SubDomain = strings.ToLower(e.SubDomain)
	}
	if !validSlug.MatchString(e.Slug) {
		return fmt.Errorf("slug inválido: solo minúsculas, números y _")
	}

	// Validate slug length (HU-001)
	if len(e.Slug) < MinSlugLength || len(e.Slug) > MaxSlugLength {
		return fmt.Errorf("slug debe tener entre %d y %d caracteres", MinSlugLength, MaxSlugLength)
	}

	// Validate subdomain length if provided
	if e.SubDomain != "" && (len(e.SubDomain) < MinSlugLength || len(e.SubDomain) > MaxSlugLength) {
		return fmt.Errorf("el subdominio debe tener entre %d y %d caracteres", MinSlugLength, MaxSlugLength)
	}

	// Validate plan limit before creating enterprise (HU-008)
	if e.TenantID > 0 {
		plan, err := s.repo.GetPlanByEnterpriseID(ctx, e.TenantID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("verificando plan: %w", err)
		}
		if plan != nil && plan.MaxEnterprises != nil {
			currentCount, err := s.repo.CountEnterprisesByTenant(ctx, e.TenantID)
			if err != nil {
				return fmt.Errorf("contando empresas: %w", err)
			}
			if currentCount >= int64(*plan.MaxEnterprises) {
				return ErrPlanLimitReached
			}
		}
	}

	// Check if slug already exists
	existingSlug, err := s.repo.GetBySlug(ctx, e.Slug)
	if err == nil && existingSlug != nil {
		return fmt.Errorf("el slug %s ya está registrado", e.Slug)
	}
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("verificando slug: %w", err)
	}

	// Check if subdomain already exists
	if e.SubDomain != "" {
		existingSubDomain, err := s.repo.GetBySubDomain(ctx, e.SubDomain)
		if err == nil && existingSubDomain != nil {
			return fmt.Errorf("el subdominio %s ya está registrado", e.SubDomain)
		}
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("verificando subdominio: %w", err)
		}
	}

	// Check if email already exists before creating
	existing, err := s.repo.GetByEmail(ctx, e.Email)
	if err == nil && existing != nil {
		return fmt.Errorf("el email %s ya está registrado", e.Email)
	}
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("verificando email: %w", err)
	}

	// Check if email already exists in public.users (HU-002)
	emailExistsInUsers, err := s.repo.EmailExistsInUsers(ctx, e.Email)
	if err != nil {
		return fmt.Errorf("verificando email en usuarios: %w", err)
	}
	if emailExistsInUsers {
		return fmt.Errorf("el correo electrónico %s ya está registrado", e.Email)
	}

	// CreateEnterprise via migrator handles schema + tenant + user + third_party
	if s.migrator != nil {
		if err := s.migrator.RunMigrations(ctx, e, passwordHash); err != nil {
			return fmt.Errorf("crear esquema: %w", err)
		}
	}

	s.publish(NewCreatedEvent(e))
	return nil
}

func (s *service) GetBySlug(ctx context.Context, slug string) (*Enterprise, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *service) GetBySubDomain(ctx context.Context, subDomain string) (*Enterprise, error) {
	return s.repo.GetBySubDomain(ctx, subDomain)
}

func (s *service) GetByEmail(ctx context.Context, email vo.Email) (*Enterprise, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *service) List(ctx context.Context, params ListParams) (ListResult, error) {
	return s.repo.List(ctx, params)
}

func (s *service) Update(ctx context.Context, e *Enterprise) error {
	e.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, e); err != nil {
		return err
	}
	s.publish(NewUpdatedEvent(e))
	return nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	e, err := s.repo.GetBySlug(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	if s.rawDB != nil {
		s.rawDB.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %q CASCADE", e.Slug))
	}
	s.publish(NewDeletedEvent(e))
	return nil
}

func (s *service) publish(event events.Event) {
	if s.eventBus == nil {
		return
	}
	if err := s.eventBus.Publish(event); err != nil {
		fmt.Printf("[enterprise.Service] warn: publish failed: %v\n", err)
	}
}
