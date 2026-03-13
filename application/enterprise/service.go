package enterprise

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cloud-tech-develop/aura-back/domain/enterprise"
	"github.com/cloud-tech-develop/aura-back/domain/events"
	"github.com/cloud-tech-develop/aura-back/infrastructure/persistence/postgres"
)

var validSlug = regexp.MustCompile(`^[a-z0-9_]+$`)

type Service struct {
	repo     enterprise.Repository
	eventBus events.EventBus
	migrator Migrator
}

type Migrator interface {
	RunMigrations(e *enterprise.Enterprise, passwordHash string) error
}

func NewService(repo enterprise.Repository, eventBus events.EventBus, migrator Migrator) *Service {
	return &Service{
		repo:     repo,
		eventBus: eventBus,
		migrator: migrator,
	}
}

func (s *Service) Create(ctx context.Context, e *enterprise.Enterprise, passwordHash string) error {
	e.Slug = strings.ToLower(e.Slug)
	if e.SubDomain != "" {
		e.SubDomain = strings.ToLower(e.SubDomain)
	}
	if !validSlug.MatchString(e.Slug) {
		return fmt.Errorf("slug inválido: solo minúsculas, números y _")
	}

	schemaCreated, err := s.createSchema(ctx, e, passwordHash)
	if err != nil {
		return fmt.Errorf("crear esquema: %w", err)
	}


	if err := s.repo.Create(ctx, e); err != nil {
		if schemaCreated {
			s.dropSchema(ctx, e.Slug)
		}
		return fmt.Errorf("registrar enterprise: %w", err)
	}

	if s.eventBus != nil {
		event := enterprise.NewEnterpriseCreatedEvent(e)
		if err := s.eventBus.Publish(event); err != nil {
			fmt.Printf("[Service] Warning: failed to publish event: %v\n", err)
		}
	}

	return nil
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*enterprise.Enterprise, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *Service) GetBySubDomain(ctx context.Context, subDomain string) (*enterprise.Enterprise, error) {
	return s.repo.GetBySubDomain(ctx, subDomain)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*enterprise.Enterprise, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) List(ctx context.Context) ([]enterprise.Enterprise, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, e *enterprise.Enterprise) error {
	e.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, e); err != nil {
		return err
	}

	if s.eventBus != nil {
		event := enterprise.NewEnterpriseUpdatedEvent(e)
		if err := s.eventBus.Publish(event); err != nil {
			fmt.Printf("[Service] Warning: failed to publish event: %v\n", err)
		}
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	e, err := s.repo.GetBySlug(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.dropSchema(ctx, e.Slug)

	if s.eventBus != nil {
		event := enterprise.NewEnterpriseDeletedEvent(e)
		if err := s.eventBus.Publish(event); err != nil {
			fmt.Printf("[Service] Warning: failed to publish event: %v\n", err)
		}
	}

	return nil
}

func (s *Service) createSchema(ctx context.Context, e *enterprise.Enterprise, passwordHash string) (bool, error) {
	if s.migrator == nil {
		return false, nil
	}
	if err := s.migrator.RunMigrations(e, passwordHash); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) dropSchema(ctx context.Context, slug string) {
	if s.repo == nil {
		return
	}
	db, ok := s.repo.(*postgres.PostgresRepository)
	if !ok {
		return
	}
	db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", slug))
}

func NewServiceWithDB(db *sql.DB, eventBus events.EventBus, migrator Migrator) *Service {
	repo := postgres.NewPostgresRepository(db)
	return NewService(repo, eventBus, migrator)
}
