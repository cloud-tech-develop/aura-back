package products

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/cloud-tech-develop/aura-back/shared/logging"
)

type service struct {
	repo     Repository
	eventBus events.EventBus
}

func NewService(db *db.DB, eventBus events.EventBus) Service {
	return &service{repo: NewRepository(db), eventBus: eventBus}
}

func (s *service) Create(ctx context.Context, tenantSlug string, p *Product) error {
	logger := logging.NewLoggerHandler("logs")
	logger.Log(fmt.Sprintf("[Product Service] Starting product creation for SKU: %s", p.SKU))

	if p.SKU == "" {
		logger.Log("[Product Service] Validation failed: SKU is required")
		return fmt.Errorf("sku is required")
	}
	logger.Log("[Product Service] SKU validation passed")

	if p.Name == "" {
		logger.Log("[Product Service] Validation failed: Name is required")
		return fmt.Errorf("name is required")
	}
	logger.Log("[Product Service] Name validation passed")

	if p.SalePrice < p.CostPrice {
		logger.Logf("[Product Service] Validation failed: Sale price (%v) less than cost price (%v)", p.SalePrice, p.CostPrice)
		return fmt.Errorf("sale price must be greater than or equal to cost price")
	}
	logger.Logf("[Product Service] Price validation passed: cost=%v, sale=%v", p.CostPrice, p.SalePrice)

	p.Status = "ACTIVE"
	p.CurrentStock = 0
	logger.Log("[Product Service] Set default values: Status=ACTIVE, CurrentStock=0")

	logger.Logf("[Product Service] Checking if SKU %s already exists for enterprise %d", p.SKU, p.EnterpriseID)
	_, err := s.repo.GetBySKU(ctx, tenantSlug, p.SKU, p.EnterpriseID)
	if err == nil {
		logger.Logf("[Product Service] SKU %s already exists", p.SKU)
		return fmt.Errorf("sku %s already exists", p.SKU)
	}
	if err != sql.ErrNoRows {
		logger.Logf("[Product Service] Error checking SKU: %v", err)
		return fmt.Errorf("error checking sku: %w", err)
	}
	logger.Log("[Product Service] SKU uniqueness check passed")

	logger.Logf("[Product Service] Creating product in repository for tenant %s", tenantSlug)
	err = s.repo.Create(ctx, tenantSlug, p)
	if err != nil {
		logger.Logf("[Product Service] Repository create failed: %v", err)
		return err
	}
	logger.Logf("[Product Service] Product created successfully with ID: %d", p.ID)

	logger.Log("[Product Service] Publishing ProductCreated event")
	s.publish(NewCreatedEvent(p))
	logger.Log("[Product Service] Product creation completed successfully")
	return nil
}

func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error) {
	logger := logging.NewLoggerHandler("logs")
	logger.Logf("[Product Service] Fetching product by ID: %d for tenant %s", id, tenantSlug)

	product, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logf("[Product Service] Product not found with ID: %d", id)
			return nil, sql.ErrNoRows
		}
		logger.Logf("[Product Service] Error fetching product: %v", err)
		return nil, fmt.Errorf("error fetching product: %w", err)
	}
	logger.Logf("[Product Service] Product found: ID=%d, SKU=%s, Name=%s", product.ID, product.SKU, product.Name)
	return product, nil
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
	logger := logging.NewLoggerHandler("logs")
	logger.Logf("[Product Service] Starting product update for ID: %d", id)

	existing, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logf("[Product Service] Product not found with ID: %d", id)
			return fmt.Errorf("product not found")
		}
		logger.Logf("[Product Service] Error fetching product: %v", err)
		return fmt.Errorf("error fetching product: %w", err)
	}
	logger.Logf("[Product Service] Found existing product: ID=%d, SKU=%s", existing.ID, existing.SKU)

	if p.SKU == "" {
		p.SKU = existing.SKU
		logger.Log("[Product Service] Using existing SKU")
	} else {
		logger.Logf("[Product Service] SKU will be updated from %s to %s", existing.SKU, p.SKU)
	}
	if p.Name == "" {
		p.Name = existing.Name
		logger.Log("[Product Service] Using existing name")
	} else {
		logger.Logf("[Product Service] Name will be updated from %s to %s", existing.Name, p.Name)
	}
	if p.SalePrice == 0 {
		p.SalePrice = existing.SalePrice
		logger.Log("[Product Service] Using existing sale price")
	} else {
		logger.Logf("[Product Service] Sale price will be updated from %v to %v", existing.SalePrice, p.SalePrice)
	}
	if p.CostPrice == 0 {
		p.CostPrice = existing.CostPrice
		logger.Log("[Product Service] Using existing cost price")
	} else {
		logger.Logf("[Product Service] Cost price will be updated from %v to %v", existing.CostPrice, p.CostPrice)
	}

	p.ID = id
	p.EnterpriseID = existing.EnterpriseID
	logger.Logf("[Product Service] Set ID=%d and EnterpriseID=%d", id, existing.EnterpriseID)

	logger.Logf("[Product Service] Updating product in repository for tenant %s", tenantSlug)
	err = s.repo.Update(ctx, tenantSlug, p)
	if err != nil {
		logger.Logf("[Product Service] Repository update failed: %v", err)
		return err
	}
	logger.Logf("[Product Service] Product updated successfully with ID: %d", id)

	logger.Log("[Product Service] Publishing ProductUpdated event")
	s.publish(NewUpdatedEvent(p))
	logger.Log("[Product Service] Product update completed successfully")
	return nil
}

func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	logger := logging.NewLoggerHandler("logs")
	logger.Logf("[Product Service] Starting product deletion for ID: %d", id)

	_, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logf("[Product Service] Product not found with ID: %d", id)
			return fmt.Errorf("product not found")
		}
		logger.Logf("[Product Service] Error fetching product: %v", err)
		return fmt.Errorf("error fetching product: %w", err)
	}
	logger.Logf("[Product Service] Found product with ID: %d", id)

	logger.Logf("[Product Service] Deleting product from repository for tenant %s", tenantSlug)
	err = s.repo.Delete(ctx, tenantSlug, id)
	if err != nil {
		logger.Logf("[Product Service] Repository delete failed: %v", err)
		return err
	}
	logger.Logf("[Product Service] Product deleted successfully with ID: %d", id)

	logger.Log("[Product Service] Publishing ProductDeleted event")
	s.publish(NewDeletedEvent(&Product{ID: id}))
	logger.Log("[Product Service] Product deletion completed successfully")
	return nil
}

func (s *service) publish(event events.Event) {
	if s.eventBus == nil {
		return
	}
	if err := s.eventBus.Publish(event); err != nil {
		fmt.Printf("[products.Service] warn: publish failed: %v\n", err)
	}
}
