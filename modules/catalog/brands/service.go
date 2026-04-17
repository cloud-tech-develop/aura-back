package brands

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

func (s *service) Create(ctx context.Context, tenantSlug string, b *Brand) error {
	if b.Name == "" {
		return fmt.Errorf("name is required")
	}
	return s.repo.Create(ctx, tenantSlug, b)
}

func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*Brand, error) {
	return s.repo.GetByID(ctx, tenantSlug, id)
}

func (s *service) List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]domain.ListId, error) {
	return s.repo.List(ctx, tenantSlug, enterpriseID)
}

func (s *service) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.Page(ctx, tenantSlug, enterpriseID, page, limit, search, sort, order, params)
}

func (s *service) Update(ctx context.Context, tenantSlug string, id int64, b *Brand) error {
	existing, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		return err
	}
	b.ID = existing.ID
	b.EnterpriseID = existing.EnterpriseID
	b.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, tenantSlug, b)
}

func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	return s.repo.Delete(ctx, tenantSlug, id)
}
