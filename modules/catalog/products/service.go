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

// Create creates a new product in the catalog
// Validates business rules before persisting
func (s *service) Create(ctx context.Context, tenantSlug string, p *Product) error {
	logger := logging.NewLoggerHandler("logs")
	logger.Log(fmt.Sprintf("[Product Service] Starting product creation for SKU: %s", p.SKU))

	// Validate SKU is required
	if p.SKU == "" {
		logger.Log("[Product Service] Validation failed: SKU is required")
		return fmt.Errorf("sku is required")
	}
	logger.Log("[Product Service] SKU validation passed")

	// Validate name is required
	if p.Name == "" {
		logger.Log("[Product Service] Validation failed: Name is required")
		return fmt.Errorf("name is required")
	}
	logger.Log("[Product Service] Name validation passed")

	// Validate unit measure is required
	if p.UnitID == 0 {
		logger.Log("[Product Service] Validation failed: UnitID is required")
		return fmt.Errorf("unit_id is required")
	}
	logger.Log("[Product Service] UnitMeasureID validation passed")

	// Validate product type is valid
	if p.ProductType == "" {
		p.ProductType = "ESTANDAR"
	} else if !IsValidProductType(p.ProductType) {
		logger.Logf("[Product Service] Validation failed: invalid product type '%s'", p.ProductType)
		return fmt.Errorf("invalid product type: %s", p.ProductType)
	}
	logger.Logf("[Product Service] Product type validation passed: %s", p.ProductType)

	// Set default values for new product
	p.Active = true
	p.VisibleInPOS = true
	if p.ManagesInventory {
		p.CurrentStock = 0
	}
	logger.Log("[Product Service] Set default values: Active=true, VisibleInPOS=true, CurrentStock=0")

	// Validate pricing: cost and sale prices must be non-negative
	if p.CostPrice < 0 {
		logger.Logf("[Product Service] Validation failed: cost_price (%v) cannot be negative", p.CostPrice)
		return fmt.Errorf("cost_price cannot be negative")
	}
	if p.SalePrice < 0 {
		logger.Logf("[Product Service] Validation failed: sale_price (%v) cannot be negative", p.SalePrice)
		return fmt.Errorf("sale_price cannot be negative")
	}
	logger.Logf("[Product Service] Price validation passed: cost=%v, sale=%v", p.CostPrice, p.SalePrice)

	// Check SKU uniqueness within enterprise
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

	// Check barcode uniqueness if provided
	if p.Barcode != "" {
		logger.Logf("[Product Service] Checking if barcode %s already exists for enterprise %d", p.Barcode, p.EnterpriseID)
		_, err := s.repo.GetByBarcode(ctx, tenantSlug, p.Barcode, p.EnterpriseID)
		if err == nil {
			logger.Logf("[Product Service] Barcode %s already exists", p.Barcode)
			return fmt.Errorf("barcode %s already exists", p.Barcode)
		}
		if err != sql.ErrNoRows {
			logger.Logf("[Product Service] Error checking barcode: %v", err)
			return fmt.Errorf("error checking barcode: %w", err)
		}
		logger.Log("[Product Service] Barcode uniqueness check passed")
	}

	// Create product in repository
	logger.Logf("[Product Service] Creating product in repository for tenant %s", tenantSlug)
	err = s.repo.Create(ctx, tenantSlug, p)
	if err != nil {
		logger.Logf("[Product Service] Repository create failed: %v", err)
		return err
	}
	logger.Logf("[Product Service] Product created successfully with ID: %d", p.ID)

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
		logger.Logf("[Product Service] Presentations created successfully")
	}

	// Publish domain event
	logger.Log("[Product Service] Publishing ProductCreated event")
	s.publish(NewCreatedEvent(p))
	logger.Log("[Product Service] Product creation completed successfully")
	return nil
}

// GetByID retrieves a product by its ID
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

// GetBySKU retrieves a product by its SKU code
func (s *service) GetBySKU(ctx context.Context, tenantSlug string, sku string, enterpriseID int64) (*Product, error) {
	logger := logging.NewLoggerHandler("logs")
	logger.Logf("[Product Service] Fetching product by SKU: %s for enterprise: %d", sku, enterpriseID)

	product, err := s.repo.GetBySKU(ctx, tenantSlug, sku, enterpriseID)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logf("[Product Service] Product not found with SKU: %s", sku)
			return nil, sql.ErrNoRows
		}
		logger.Logf("[Product Service] Error fetching product by SKU: %v", err)
		return nil, fmt.Errorf("error fetching product by sku: %w", err)
	}
	logger.Logf("[Product Service] Product found: ID=%d, SKU=%s, Name=%s", product.ID, product.SKU, product.Name)
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
	logger.Logf("[Product Service] Starting product update for ID: %d", id)

	// Get existing product to validate and preserve values
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

	// Preserve unchanged values from existing product
	if p.SKU == "" {
		p.SKU = existing.SKU
		logger.Log("[Product Service] Using existing SKU")
	} else {
		logger.Logf("[Product Service] SKU will be updated from %s to %s", existing.SKU, p.SKU)
	}
	if p.Barcode == "" {
		p.Barcode = existing.Barcode
		logger.Log("[Product Service] Using existing barcode")
	}
	if p.Name == "" {
		p.Name = existing.Name
		logger.Log("[Product Service] Using existing name")
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
		logger.Logf("[Product Service] Invalid product type: %s", p.ProductType)
		return fmt.Errorf("invalid product type: %s", p.ProductType)
	}

	// Preserve pricing if not specified
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

	// Set ID and EnterpriseID for update
	p.ID = id
	p.EnterpriseID = existing.EnterpriseID
	logger.Logf("[Product Service] Set ID=%d and EnterpriseID=%d", id, existing.EnterpriseID)

	// Update in repository
	logger.Logf("[Product Service] Updating product in repository for tenant %s", tenantSlug)
	err = s.repo.Update(ctx, tenantSlug, p)
	if err != nil {
		logger.Logf("[Product Service] Repository update failed: %v", err)
		return err
	}
	logger.Logf("[Product Service] Product updated successfully with ID: %d", id)

	// Update presentations if provided
	if len(p.Presentations) > 0 {
		logger.Logf("[Product Service] Updating %d presentations for product ID: %d", len(p.Presentations), id)

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
			logger.Logf("[Product Service] Failed to update presentations: %v", err)
			return fmt.Errorf("failed to update presentations: %w", err)
		}
		logger.Logf("[Product Service] Presentations updated successfully")
	}

	// Publish domain event
	logger.Log("[Product Service] Publishing ProductUpdated event")
	s.publish(NewUpdatedEvent(p))
	logger.Log("[Product Service] Product update completed successfully")
	return nil
}

// Delete performs a soft delete of a product
func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	logger := logging.NewLoggerHandler("logs")
	logger.Logf("[Product Service] Starting product deletion for ID: %d", id)

	// Verify product exists
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

	// Soft delete in repository
	logger.Logf("[Product Service] Deleting product from repository for tenant %s", tenantSlug)
	err = s.repo.Delete(ctx, tenantSlug, id)
	if err != nil {
		logger.Logf("[Product Service] Repository delete failed: %v", err)
		return err
	}
	logger.Logf("[Product Service] Product deleted successfully with ID: %d", id)

	// Publish domain event
	logger.Log("[Product Service] Publishing ProductDeleted event")
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
