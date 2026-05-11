package compositions

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/domain"
)

type service struct {
	repo      Repository
	isOffline bool
}

func NewService(database *db.DB) Service {
	return &service{
		repo:      NewRepository(database),
		isOffline: database.IsSQLite(),
	}
}

func (s *service) Create(ctx context.Context, tenantSlug string, c *Composition) error {
	// Validate required fields
	if c.ParentProductID == 0 {
		return fmt.Errorf("parent_product_id is required")
	}
	if c.ChildProductID == 0 {
		return fmt.Errorf("child_product_id is required")
	}

	// Validate parent is not same as child
	if c.ParentProductID == c.ChildProductID {
		return fmt.Errorf("parent product cannot be the same as child product")
	}

	// Validate type
	if c.Type == "" {
		c.Type = "KIT"
	}
	if c.Type != "KIT" && c.Type != "RECETA" {
		return fmt.Errorf("invalid composition type: %s (must be KIT or RECETA)", c.Type)
	}

	// Validate quantity
	if c.Quantity <= 0 {
		c.Quantity = 1
	}

	// Check for duplicate
	exists, err := s.repo.ExistsPair(ctx, tenantSlug, c.ParentProductID, c.ChildProductID)
	if err != nil {
		return fmt.Errorf("error checking duplicate composition: %w", err)
	}
	if exists {
		return fmt.Errorf("composition already exists for this parent and child")
	}

	// Create
	if err := s.repo.Create(ctx, tenantSlug, c); err != nil {
		return err
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*Composition, error) {
	c, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("error fetching composition: %w", err)
	}
	return c, nil
}

func (s *service) ListByParent(ctx context.Context, tenantSlug string, parentID int64) ([]Composition, error) {
	return s.repo.ListByParent(ctx, tenantSlug, parentID)
}

func (s *service) ListAll(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Composition, error) {
	return s.repo.ListAll(ctx, tenantSlug, enterpriseID)
}

func (s *service) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, tipo string, sort string, order string) (domain.PageResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.Page(ctx, tenantSlug, enterpriseID, page, limit, search, tipo, sort, order)
}

func (s *service) Update(ctx context.Context, tenantSlug string, id int64, c *Composition) error {
	existing, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		return err
	}

	// Only update quantity and type (products are not changed)
	if c.Quantity > 0 {
		existing.Quantity = c.Quantity
	}
	if c.Type != "" {
		if c.Type != "KIT" && c.Type != "RECETA" {
			return fmt.Errorf("invalid composition type: %s (must be KIT or RECETA)", c.Type)
		}
		existing.Type = c.Type
	}

	if err := s.repo.Update(ctx, tenantSlug, existing); err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	// Verify existence
	_, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, tenantSlug, id)
}
