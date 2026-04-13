package products

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
)

type service struct {
	repo Repository
}

func NewService(db *db.DB) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) Create(ctx context.Context, tenantSlug string, p *Product) error {
	if p.SKU == "" {
		return fmt.Errorf("sku is required")
	}
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.SalePrice < p.CostPrice {
		return fmt.Errorf("sale price must be greater than or equal to cost price")
	}

	p.Status = "ACTIVE"
	p.CurrentStock = 0

	_, err := s.repo.GetBySKU(ctx, tenantSlug, p.SKU, p.EnterpriseID)
	if err == nil {
		return fmt.Errorf("sku %s already exists", p.SKU)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking sku: %w", err)
	}

	return s.repo.Create(ctx, tenantSlug, p)
}

func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error) {
	return s.repo.GetByID(ctx, tenantSlug, id)
}

func (s *service) Page(ctx context.Context, tenantSlug string, enterpriseID int64, first int64, rows int64, search string) ([]Product, error) {
	if first < 1 {
		first = 1
	}
	if rows < 1 {
		rows = 10
	}
	return s.repo.Page(ctx, tenantSlug, enterpriseID, first, rows, search)
}

func (s *service) List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Product, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 10
	}
	return s.repo.List(ctx, tenantSlug, enterpriseID, filters)
}

func (s *service) Update(ctx context.Context, tenantSlug string, id int64, p *Product) error {
	existing, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("error fetching product: %w", err)
	}

	if p.SKU == "" {
		p.SKU = existing.SKU
	}
	if p.Name == "" {
		p.Name = existing.Name
	}
	if p.SalePrice == 0 {
		p.SalePrice = existing.SalePrice
	}
	if p.CostPrice == 0 {
		p.CostPrice = existing.CostPrice
	}

	p.ID = id
	p.EnterpriseID = existing.EnterpriseID

	return s.repo.Update(ctx, tenantSlug, p)
}

func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	_, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("error fetching product: %w", err)
	}
	return s.repo.Delete(ctx, tenantSlug, id)
}
