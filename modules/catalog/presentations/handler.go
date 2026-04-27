package presentations

import (
	"database/sql"
	"strconv"

	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for presentations
type Handler struct {
	svc Service
}

// NewHandler creates a new presentation handler instance
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Create handles POST /products/:id/presentations
// Creates multiple presentations for a product
func (h *Handler) Create(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	enterpriseID := claims.EnterpriseID

	if tenantSlug == "" || enterpriseID == 0 {
		response.BadRequest(c, "tenant not found")
		return
	}

	// Get product ID from URL parameter
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product ID")
		return
	}

	// Request structure for list of presentations
	var req PresentationListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Validate product ID is provided and valid
	if productID == 0 {
		response.BadRequest(c, "product_id is required")
		return
	}

	if err := h.svc.Create(c.Request.Context(), tenantSlug, enterpriseID, productID, req.Presentations); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, gin.H{"message": "Presentations created successfully", "count": len(req.Presentations)})
}

// UpsertArray handles PUT /products/:id/presentations
// Creates or updates presentations for a product based on ID presence
// Accepts direct array of presentations
func (h *Handler) UpsertArray(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	enterpriseID := claims.EnterpriseID

	if tenantSlug == "" || enterpriseID == 0 {
		response.BadRequest(c, "tenant not found")
		return
	}

	// Get product ID from URL parameter
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product ID")
		return
	}

	// Request structure for direct array of presentations
	var presentations []PresentationRequest
	if err := c.ShouldBindJSON(&presentations); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if len(presentations) == 0 {
		response.BadRequest(c, "at least one presentation is required")
		return
	}

	if err := h.svc.Upsert(c.Request.Context(), tenantSlug, enterpriseID, productID, presentations); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Presentations upserted successfully", "count": len(presentations)})
}

// Page handles POST /presentations/page
// Returns a paginated list of presentations
func (h *Handler) Page(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	enterpriseID := claims.EnterpriseID

	// Fallback: get from query params (for offline sync)
	if tenantSlug == "" || enterpriseID == 0 {
		if c.Query("slug") != "" {
			tenantSlug = c.Query("slug")
			if eid, err := strconv.ParseInt(c.Query("enterprise_id"), 10, 64); err == nil {
				enterpriseID = eid
			}
		}
	}

	if tenantSlug == "" || enterpriseID == 0 {
		response.BadRequest(c, "tenant not found")
		return
	}

	var req struct {
		Page   int64          `json:"page"`
		Limit  int64          `json:"limit"`
		Search string         `json:"search"`
		Sort   string         `json:"sort"`
		Order  string         `json:"order"`
		Params map[string]any `json:"params"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Apply defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Sort == "" {
		req.Sort = "id"
	}
	if req.Order == "" {
		req.Order = "asc"
	}
	if req.Params == nil {
		req.Params = make(map[string]any)
	}

	result, err := h.svc.Page(c.Request.Context(), tenantSlug, enterpriseID, req.Page, req.Limit, req.Search, req.Sort, req.Order, req.Params)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, result)
}

// List handles GET /presentations
// Returns a list of presentations, optionally filtered by product_id
func (h *Handler) List(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	enterpriseID := claims.EnterpriseID

	// Fallback: get from query params (for offline sync)
	if tenantSlug == "" || enterpriseID == 0 {
		if c.Query("slug") != "" {
			tenantSlug = c.Query("slug")
			if eid, err := strconv.ParseInt(c.Query("enterprise_id"), 10, 64); err == nil {
				enterpriseID = eid
			}
		}
	}

	if tenantSlug == "" || enterpriseID == 0 {
		response.BadRequest(c, "tenant not found")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	var productID *int64
	if pidStr := c.Query("product_id"); pidStr != "" {
		if id, err := strconv.ParseInt(pidStr, 10, 64); err == nil {
			productID = &id
		}
	}

	filters := ListFilters{
		Page:      page,
		Limit:     limit,
		Search:    search,
		ProductID: productID,
	}

	list, err := h.svc.List(c.Request.Context(), tenantSlug, enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

// GetByProductID handles GET /products/:id/presentations
// Returns all presentations for a specific product
func (h *Handler) GetByProductID(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug

	if tenantSlug == "" {
		response.BadRequest(c, "tenant not found")
		return
	}

	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product ID")
		return
	}

	presentations, err := h.svc.GetByProductID(c.Request.Context(), tenantSlug, productID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, presentations)
}

// GetByID handles GET /presentations/:id
// Returns a single presentation by ID
func (h *Handler) GetByID(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug

	if tenantSlug == "" {
		response.BadRequest(c, "tenant not found")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ID")
		return
	}

	presentation, err := h.svc.GetByID(c.Request.Context(), tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "presentation not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, presentation)
}

// Update handles PUT /presentations/:id
// Updates an existing presentation
func (h *Handler) Update(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug

	if tenantSlug == "" {
		response.BadRequest(c, "tenant not found")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ID")
		return
	}

	var req struct {
		Name            string  `json:"name"`
		Factor          float64 `json:"factor"`
		Barcode         string  `json:"barcode"`
		CostPrice       float64 `json:"cost_price"`
		SalePrice       float64 `json:"sale_price"`
		DefaultPurchase bool    `json:"default_purchase"`
		DefaultSale     bool    `json:"default_sale"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	presentation := &Presentation{
		Name:            req.Name,
		Factor:          req.Factor,
		Barcode:         req.Barcode,
		CostPrice:       req.CostPrice,
		SalePrice:       req.SalePrice,
		DefaultPurchase: req.DefaultPurchase,
		DefaultSale:     req.DefaultSale,
	}

	if err := h.svc.Update(c.Request.Context(), tenantSlug, id, presentation); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "presentation not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, presentation)
}

// Delete handles DELETE /presentations/:id
// Performs a soft delete of a presentation
func (h *Handler) Delete(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug

	if tenantSlug == "" {
		response.BadRequest(c, "tenant not found")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ID")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), tenantSlug, id); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "presentation not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	c.Status(204)
}
