package categories

import (
	"context"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/domain"
)

type service struct {
	repo Repository
}

func NewService(db db.Querier) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) Create(ctx context.Context, tenantSlug string, c *Category) error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	return s.repo.Create(ctx, tenantSlug, c)
}

func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*Category, error) {
	return s.repo.GetByID(ctx, tenantSlug, id)
}

func (s *service) List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]Category, error) {
	return s.repo.List(ctx, tenantSlug, enterpriseID)
}

func (s *service) Page(ctx context.Context, tenantSlug string, enterpriseID int64, first int64, rows int64, search string) (domain.PageResult, error) {
	if first < 1 {
		first = 1
	}
	if rows < 1 {
		rows = 10
	}
	return s.repo.Page(ctx, tenantSlug, enterpriseID, first, rows, search)
}

func (s *service) Update(ctx context.Context, tenantSlug string, id int64, c *Category) error {
	existing, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		return err
	}
	c.ID = existing.ID
	c.EnterpriseID = existing.EnterpriseID
	c.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, tenantSlug, c)
}

func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	return s.repo.Delete(ctx, tenantSlug, id)
}
