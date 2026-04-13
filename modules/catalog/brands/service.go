package brands

import (
	"context"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
)

type service struct {
	repo Repository
}

func NewService(db db.Querier) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) Create(ctx context.Context, b *Brand) error {
	if b.Name == "" {
		return fmt.Errorf("name is required")
	}
	return s.repo.Create(ctx, b)
}

func (s *service) GetByID(ctx context.Context, id int64) (*Brand, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context, enterpriseID int64) ([]Brand, error) {
	return s.repo.List(ctx, enterpriseID)
}

func (s *service) Update(ctx context.Context, id int64, b *Brand) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	b.ID = existing.ID
	b.EnterpriseID = existing.EnterpriseID
	b.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, b)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
