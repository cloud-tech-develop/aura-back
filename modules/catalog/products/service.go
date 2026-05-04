package products

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/presentations"
	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/cloud-tech-develop/aura-back/shared/logging"
)

// service implements the Service interface
// Contains business logic for product management
type service struct {
	repo            Repository
	presentationSvc presentations.Service
	eventBus        events.EventBus
}

// NewService creates a new product service instance
// db: database connection instance
// eventBus: event bus for publishing domain events
func NewService(db *db.DB, eventBus events.EventBus, presSvc presentations.Service) Service {
	return &service{repo: NewRepository(db), presentationSvc: presSvc, eventBus: eventBus}
}

// Upsert upserts a product
// Validates business rules before persisting
func (s *service) Upsert(ctx context.Context, tenantSlug string, p Product) error {
	if err := s.repo.Upsert(ctx, tenantSlug, &p); err != nil {
		return err
	}
	return nil
}

// Create creates a new product in the catalog
// Validates business rules before persisting
func (s *service) Create(ctx context.Context, tenantSlug string, p *Product) error {
	logger := logging.NewLoggerHandler("logs")

	// Validate SKU is required
	if p.SKU == "" {
		return fmt.Errorf("sku is required")
	}

	// Validate name is required
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Validate unit measure is required
	if p.UnitID == 0 {
		return fmt.Errorf("unit_id is required")
	}

	// Validate product type is valid
	if p.ProductType == "" {
		p.ProductType = "STANDARD"
	} else if !IsValidProductType(p.ProductType) {
		return fmt.Errorf("invalid product type: %s", p.ProductType)
	}

	// Set default values for new product
	p.Active = true
	p.VisibleInPOS = true
	if p.ManagesInventory {
		p.CurrentStock = 0
	}

	// Validate pricing: cost and sale prices must be non-negative
	if p.CostPrice < 0 {
		return fmt.Errorf("cost_price cannot be negative")
	}
	if p.SalePrice < 0 {
		return fmt.Errorf("sale_price cannot be negative")
	}

	// Check barcode uniqueness if provided
	if p.Barcode != "" {
		_, err := s.repo.GetByBarcode(ctx, tenantSlug, p.Barcode, p.EnterpriseID)
		if err == nil {
			return fmt.Errorf("barcode %s already exists", p.Barcode)
		}
		if err != sql.ErrNoRows {
			return fmt.Errorf("error checking barcode: %w", err)
		}
	}

	// Create product in repository
	err := s.repo.Create(ctx, tenantSlug, p)
	if err != nil {
		return err
	}

	if p.ID == 0 && len(p.Presentations) > 0 {
		prod, _ := s.repo.GetBySKU(ctx, tenantSlug, p.SKU, p.EnterpriseID)
		if prod != nil {
			p.ID = prod.ID
		}
	}

	// Create presentations if provided
	if len(p.Presentations) > 0 {
		logger.Logf("[Product Service] Creating %d presentations for product ID: %d", len(p.Presentations), p.ID)

		// Convert Product.Presentations to presentations.PresentationRequest
		presRequests := make([]presentations.PresentationRequest, len(p.Presentations))
		for i, pres := range p.Presentations {
			presRequests[i] = presentations.PresentationRequest{
				ID:              pres.ID,
				Name:            pres.Name,
				Factor:          pres.Factor,
				Barcode:         pres.Barcode,
				SalePrice:       pres.SalePrice,
				CostPrice:       pres.CostPrice,
				DefaultPurchase: pres.DefaultPurchase,
				DefaultSale:     pres.DefaultSale,
			}
		}

		if err := s.presentationSvc.Create(ctx, tenantSlug, p.EnterpriseID, p.ID, presRequests); err != nil {
			logger.Logf("[Product Service] Failed to create presentations: %v", err)
			return fmt.Errorf("failed to create presentations: %w", err)
		}
	}

	// Publish domain event
	s.publish(NewCreatedEvent(p))
	logger.Log("[Product Service] Product creation completed successfully")
	return nil
}

// GetByID retrieves a product by its ID
func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error) {
	product, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("error fetching product: %w", err)
	}
	return product, nil
}

// GetBySKU retrieves a product by its SKU code
func (s *service) GetBySKU(ctx context.Context, tenantSlug string, sku string, enterpriseID int64) (*Product, error) {
	product, err := s.repo.GetBySKU(ctx, tenantSlug, sku, enterpriseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("error fetching product by sku: %w", err)
	}
	return product, nil
}

// Page retrieves a paginated list of products
func (s *service) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.Page(ctx, tenantSlug, enterpriseID, page, limit, search, sort, order, params)
}

// List retrieves a list of products with filters
func (s *service) List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Product, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 10
	}
	return s.repo.List(ctx, tenantSlug, enterpriseID, filters)
}

// Update updates an existing product
// Validates business rules before persisting
func (s *service) Update(ctx context.Context, tenantSlug string, id int64, p *Product) error {
	logger := logging.NewLoggerHandler("logs")

	// Get existing product to validate and preserve values
	existing, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("error fetching product: %w", err)
	}

	// Preserve unchanged values from existing product
	if p.SKU == "" {
		p.SKU = existing.SKU
	} else {
		logger.Logf("[Product Service] SKU will be updated from %s to %s", existing.SKU, p.SKU)
	}
	if p.Barcode == "" {
		p.Barcode = existing.Barcode
	}
	if p.Name == "" {
		p.Name = existing.Name
	} else {
		logger.Logf("[Product Service] Name will be updated from %s to %s", existing.Name, p.Name)
	}
	if p.UnitID == 0 {
		p.UnitID = existing.UnitID
		logger.Log("[Product Service] Using existing unit measure")
	}
	if p.ProductType == "" {
		p.ProductType = existing.ProductType
	} else if !IsValidProductType(p.ProductType) {
		return fmt.Errorf("invalid product type: %s", p.ProductType)
	}

	// Preserve pricing if not specified
	if p.SalePrice == 0 {
		p.SalePrice = existing.SalePrice
	} else {
		logger.Logf("[Product Service] Sale price will be updated from %v to %v", existing.SalePrice, p.SalePrice)
	}
	if p.CostPrice == 0 {
		p.CostPrice = existing.CostPrice
	}

	// Set ID and EnterpriseID for update
	p.ID = id
	p.EnterpriseID = existing.EnterpriseID

	// Update in repository
	err = s.repo.Update(ctx, tenantSlug, p)
	if err != nil {
		return err
	}

	// Update presentations if provided
	if len(p.Presentations) > 0 {

		// Convert Product.Presentations to presentations.PresentationRequest
		presRequests := make([]presentations.PresentationRequest, len(p.Presentations))
		for i, pres := range p.Presentations {
			presRequests[i] = presentations.PresentationRequest{
				ID:              pres.ID,
				Name:            pres.Name,
				Factor:          pres.Factor,
				Barcode:         pres.Barcode,
				SalePrice:       pres.SalePrice,
				CostPrice:       pres.CostPrice,
				DefaultPurchase: pres.DefaultPurchase,
				DefaultSale:     pres.DefaultSale,
			}
		}

		if err := s.presentationSvc.Upsert(ctx, tenantSlug, p.EnterpriseID, id, presRequests); err != nil {
			return fmt.Errorf("failed to update presentations: %w", err)
		}
	}

	// Publish domain event
	s.publish(NewUpdatedEvent(p))
	logger.Log("[Product Service] Product update completed successfully")
	return nil
}

// Delete performs a soft delete of a product
func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	logger := logging.NewLoggerHandler("logs")

	// Verify product exists
	_, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("error fetching product: %w", err)
	}

	// Soft delete in repository
	err = s.repo.Delete(ctx, tenantSlug, id)
	if err != nil {
		return err
	}

	// Publish domain event
	s.publish(NewDeletedEvent(&Product{ID: id}))
	logger.Log("[Product Service] Product deletion completed successfully")
	return nil
}

// publish publishes a domain event through the event bus
func (s *service) publish(event events.Event) {
	if s.eventBus == nil {
		return
	}
	if err := s.eventBus.Publish(event); err != nil {
		fmt.Printf("[products.Service] warn: publish failed: %v\n", err)
	}
}
