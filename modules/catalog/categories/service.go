package categories

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

func (s *service) Create(ctx context.Context, c *Category) error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	return s.repo.Create(ctx, c)
}

func (s *service) GetByID(ctx context.Context, id int64) (*Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context, enterpriseID int64) ([]Category, error) {
	return s.repo.List(ctx, enterpriseID)
}

func (s *service) Update(ctx context.Context, id int64, c *Category) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	c.ID = existing.ID
	c.EnterpriseID = existing.EnterpriseID
	c.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, c)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
