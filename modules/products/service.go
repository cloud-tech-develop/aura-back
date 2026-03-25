package products

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
)

type categoryService struct {
	repo CategoryRepository
}

func NewCategoryService(db db.Querier) CategoryService {
	return &categoryService{repo: NewCategoryRepository(db)}
}

func (s *categoryService) Create(ctx context.Context, c *Category) error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	c.EnterpriseID = ctx.Value("enterprise_id").(int64) // Assuming enterprise_id is set in context
	return s.repo.Create(ctx, c)
}

func (s *categoryService) GetByID(ctx context.Context, id int64) (*Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *categoryService) List(ctx context.Context, enterpriseID int64) ([]Category, error) {
	return s.repo.List(ctx, enterpriseID)
}

func (s *categoryService) Update(ctx context.Context, id int64, c *Category) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	c.ID = existing.ID
	c.EnterpriseID = existing.EnterpriseID
	c.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, c)
}

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

type brandService struct {
	repo BrandRepository
}

func NewBrandService(db db.Querier) BrandService {
	return &brandService{repo: NewBrandRepository(db)}
}

func (s *brandService) Create(ctx context.Context, b *Brand) error {
	if b.Name == "" {
		return fmt.Errorf("name is required")
	}
	b.EnterpriseID = ctx.Value("enterprise_id").(int64)
	return s.repo.Create(ctx, b)
}

func (s *brandService) GetByID(ctx context.Context, id int64) (*Brand, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *brandService) List(ctx context.Context, enterpriseID int64) ([]Brand, error) {
	return s.repo.List(ctx, enterpriseID)
}

func (s *brandService) Update(ctx context.Context, id int64, b *Brand) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	b.ID = existing.ID
	b.EnterpriseID = existing.EnterpriseID
	b.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, b)
}

func (s *brandService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

type productService struct {
	repo ProductRepository
}

func NewProductService(db db.Querier) ProductService {
	return &productService{repo: NewProductRepository(db)}
}

func (s *productService) Create(ctx context.Context, p *Product) error {
	if p.SKU == "" {
		return fmt.Errorf("sku is required")
	}
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.SalePrice < p.CostPrice {
		return fmt.Errorf("sale price must be greater than or equal to cost price")
	}

	p.EnterpriseID = ctx.Value("enterprise_id").(int64)
	p.Status = "ACTIVE"
	p.CurrentStock = 0 // Default stock

	// Check SKU uniqueness
	_, err := s.repo.GetBySKU(ctx, p.SKU, p.EnterpriseID)
	if err == nil {
		return fmt.Errorf("sku %s already exists", p.SKU)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking sku: %w", err)
	}

	return s.repo.Create(ctx, p)
}

func (s *productService) GetByID(ctx context.Context, id int64) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) List(ctx context.Context, enterpriseID int64, filters ListFilters) ([]Product, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 10
	}
	return s.repo.List(ctx, enterpriseID, filters)
}

func (s *productService) Update(ctx context.Context, id int64, p *Product) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if p.SalePrice < p.CostPrice {
		return fmt.Errorf("sale price must be greater than or equal to cost price")
	}

	// Check SKU uniqueness if changed
	if p.SKU != existing.SKU {
		_, err := s.repo.GetBySKU(ctx, p.SKU, existing.EnterpriseID)
		if err == nil {
			return fmt.Errorf("sku %s already exists", p.SKU)
		}
		if err != sql.ErrNoRows {
			return fmt.Errorf("error checking sku: %w", err)
		}
	}

	p.ID = existing.ID
	p.EnterpriseID = existing.EnterpriseID
	p.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, p)
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
