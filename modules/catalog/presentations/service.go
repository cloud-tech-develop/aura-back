package presentations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/cloud-tech-develop/aura-back/shared/logging"
)

// service implements the Service interface
// Contains business logic for presentation management
type service struct {
	repo     Repository
	eventBus events.EventBus
}

// NewService creates a new presentation service instance
func NewService(db *db.DB, eventBus events.EventBus) Service {
	return &service{repo: NewRepository(db), eventBus: eventBus}
}

// Create creates multiple presentations for a product
func (s *service) Create(ctx context.Context, tenantSlug string, productID int64, presentations []PresentationRequest) error {
	logger := logging.NewLoggerHandler("logs")
	logger.Log(fmt.Sprintf("[Presentation Service] Creating %d presentations for product ID: %d", len(presentations), productID))

	// Validate product exists first by trying to get it
	if productID == 0 {
		logger.Log("[Presentation Service] Validation failed: product_id is required")
		return fmt.Errorf("product_id is required")
	}

	// Validate we have at least one presentation
	if len(presentations) == 0 {
		logger.Log("[Presentation Service] Validation failed: at least one presentation is required")
		return fmt.Errorf("at least one presentation is required")
	}

	// Convert request to entities
	entities := make([]*Presentation, len(presentations))
	for i, req := range presentations {
		logger.Logf("[Presentation Service] Validating presentation %d: name=%s", i+1, req.Name)

		// Validate required fields
		if req.Name == "" {
			logger.Logf("[Presentation Service] Validation failed: name is required for presentation %d", i+1)
			return fmt.Errorf("name is required for presentation %d", i+1)
		}
		if req.Factor == 0 {
			logger.Logf("[Presentation Service] Validation failed: factor is required for presentation %d", i+1)
			return fmt.Errorf("factor is required for presentation %d", i+1)
		}
		if req.SalePrice == 0 {
			logger.Logf("[Presentation Service] Validation failed: sale_price is required for presentation %d", i+1)
			return fmt.Errorf("sale_price is required for presentation %d", i+1)
		}
		if req.CostPrice == 0 {
			logger.Logf("[Presentation Service] Validation failed: cost_price is required for presentation %d", i+1)
			return fmt.Errorf("cost_price is required for presentation %d", i+1)
		}

		entities[i] = &Presentation{
			ProductID:       productID,
			Name:            req.Name,
			Factor:          req.Factor,
			Barcode:         req.Barcode,
			CostPrice:       req.CostPrice,
			SalePrice:       req.SalePrice,
			DefaultPurchase: req.DefaultPurchase,
			DefaultSale:     req.DefaultSale,
		}
	}

	// Create all presentations
	err := s.repo.CreateMany(ctx, tenantSlug, entities)
	if err != nil {
		logger.Logf("[Presentation Service] Repository create failed: %v", err)
		return err
	}

	logger.Logf("[Presentation Service] Created %d presentations successfully", len(presentations))
	return nil
}

// GetByID retrieves a presentation by its ID
func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*Presentation, error) {
	logger := logging.NewLoggerHandler("logs")
	logger.Logf("[Presentation Service] Fetching presentation by ID: %d", id)

	presentation, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logf("[Presentation Service] Presentation not found with ID: %d", id)
			return nil, sql.ErrNoRows
		}
		logger.Logf("[Presentation Service] Error fetching presentation: %v", err)
		return nil, fmt.Errorf("error fetching presentation: %w", err)
	}

	logger.Logf("[Presentation Service] Presentation found: ID=%d, Name=%s", presentation.ID, presentation.Name)
	return presentation, nil
}

// GetByProductID retrieves all presentations for a product
func (s *service) GetByProductID(ctx context.Context, tenantSlug string, productID int64) ([]Presentation, error) {
	logger := logging.NewLoggerHandler("logs")
	logger.Logf("[Presentation Service] Fetching presentations for product ID: %d", productID)

	if productID == 0 {
		return nil, fmt.Errorf("product_id is required")
	}

	presentations, err := s.repo.GetByProductID(ctx, tenantSlug, productID)
	if err != nil {
		logger.Logf("[Presentation Service] Error fetching presentations: %v", err)
		return nil, fmt.Errorf("error fetching presentations: %w", err)
	}

	logger.Logf("[Presentation Service] Found %d presentations", len(presentations))
	return presentations, nil
}

// Page retrieves a paginated list of presentations
func (s *service) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.Page(ctx, tenantSlug, enterpriseID, page, limit, search, sort, order, params)
}

// List retrieves a list of presentations with filters
func (s *service) List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Presentation, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 10
	}
	return s.repo.List(ctx, tenantSlug, enterpriseID, filters)
}

// Update updates an existing presentation
func (s *service) Update(ctx context.Context, tenantSlug string, id int64, p *Presentation) error {
	logger := logging.NewLoggerHandler("logs")
	logger.Logf("[Presentation Service] Starting presentation update for ID: %d", id)

	// Get existing to validate and preserve values
	existing, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logf("[Presentation Service] Presentation not found with ID: %d", id)
			return fmt.Errorf("presentation not found")
		}
		logger.Logf("[Presentation Service] Error fetching presentation: %v", err)
		return fmt.Errorf("error fetching presentation: %w", err)
	}

	// Preserve unchanged values
	if p.Name == "" {
		p.Name = existing.Name
	}
	if p.Factor == 0 {
		p.Factor = existing.Factor
	}
	if p.SalePrice == 0 {
		p.SalePrice = existing.SalePrice
	}
	if p.CostPrice == 0 {
		p.CostPrice = existing.CostPrice
	}

	p.ID = id
	p.ProductID = existing.ProductID
	p.EnterpriseID = existing.EnterpriseID

	logger.Logf("[Presentation Service] Updating presentation in repository")
	err = s.repo.Update(ctx, tenantSlug, p)
	if err != nil {
		logger.Logf("[Presentation Service] Repository update failed: %v", err)
		return err
	}

	logger.Logf("[Presentation Service] Presentation updated successfully")
	return nil
}

// Delete performs a soft delete of a presentation
func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	logger := logging.NewLoggerHandler("logs")
	logger.Logf("[Presentation Service] Starting presentation deletion for ID: %d", id)

	_, err := s.repo.GetByID(ctx, tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logf("[Presentation Service] Presentation not found with ID: %d", id)
			return fmt.Errorf("presentation not found")
		}
		logger.Logf("[Presentation Service] Error fetching presentation: %v", err)
		return fmt.Errorf("error fetching presentation: %w", err)
	}

	err = s.repo.Delete(ctx, tenantSlug, id)
	if err != nil {
		logger.Logf("[Presentation Service] Repository delete failed: %v", err)
		return err
	}

	logger.Logf("[Presentation Service] Presentation deleted successfully")
	return nil
}
