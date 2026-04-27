package products

import (
	"database/sql"
	"strconv"

	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for products
// Converts JSON requests to domain entities and calls the service layer
type Handler struct {
	svc Service
}

// NewHandler creates a new product handler instance
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Create handles POST /products
// Creates a new product in the catalog
func (h *Handler) Create(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	enterpriseID := claims.EnterpriseID
	if tenantSlug == "" || enterpriseID == 0 {
		response.BadRequest(c, "tenant not found")
		return
	}

	// Request structure matching the provided JSON
	var req struct {
		// Basic fields
		SKU         string `json:"sku"`
		Barcode     string `json:"barcode"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`

		// Reference fields
		CategoryID *int64 `json:"category_id"`
		BrandID    *int64 `json:"brand_id"`
		UnitID     int64  `json:"unit_measure_id" binding:"required"`

		// Product type and status
		ProductType string `json:"product_type"`
		Active      *bool  `json:"active"`
		Status      string `json:"status"`

		// Pricing
		CostPrice float64  `json:"cost_price" binding:"required"`
		SalePrice float64  `json:"sale_price" binding:"required"`
		Price2    float64  `json:"price_2"`
		Price3    *float64 `json:"price_3"`

		// Taxes
		IVAPercentage  float64 `json:"iva_percentage"`
		ConsumptionTax float64 `json:"consumption_tax"`

		// Inventory controls
		ManagesInventory   *bool `json:"manages_inventory"`
		ManagesBatches     *bool `json:"manages_batches"`
		ManagesSerial      *bool `json:"manages_serial"`
		AllowNegativeStock *bool `json:"allow_negative_stock"`
		MinStock           *int  `json:"min_stock"`

		// Visibility
		VisibleInPOS *bool `json:"visible_in_pos"`

		// Presentations
		Presentations []PresentationRequest `json:"presentations"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Set defaults
	productType := req.ProductType
	if productType == "" {
		productType = "STANDARD"
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	visibleInPOS := true
	if req.VisibleInPOS != nil {
		visibleInPOS = *req.VisibleInPOS
	}

	managesInventory := true
	if req.ManagesInventory != nil {
		managesInventory = *req.ManagesInventory
	}

	ivaPercentage := req.IVAPercentage
	if ivaPercentage == 0 {
		ivaPercentage = 19.00
	}

	minStock := 0
	if req.MinStock != nil {
		minStock = *req.MinStock
	}

	product := &Product{
		SKU:                req.SKU,
		Barcode:            req.Barcode,
		Name:               req.Name,
		Description:        req.Description,
		ImageURL:           req.ImageURL,
		CategoryID:         req.CategoryID,
		BrandID:            req.BrandID,
		UnitID:             req.UnitID,
		ProductType:        productType,
		Active:             active,
		VisibleInPOS:       visibleInPOS,
		CostPrice:          req.CostPrice,
		SalePrice:          req.SalePrice,
		Price2:             req.Price2,
		Price3:             req.Price3,
		IVAPercentage:      ivaPercentage,
		ConsumptionTax:     req.ConsumptionTax,
		ManagesInventory:   managesInventory,
		ManagesBatches:     req.ManagesBatches != nil && *req.ManagesBatches,
		ManagesSerial:      req.ManagesSerial != nil && *req.ManagesSerial,
		AllowNegativeStock: req.AllowNegativeStock != nil && *req.AllowNegativeStock,
		MinStock:           minStock,
		EnterpriseID:       enterpriseID,
		Presentations:      req.Presentations,
	}

	if err := h.svc.Create(c.Request.Context(), tenantSlug, product); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, product)
}

// Page handles POST /products/page
// Returns a paginated list of products
func (h *Handler) Page(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	enterpriseID := claims.EnterpriseID
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

// List handles GET /products
// Returns a list of products with optional filters
func (h *Handler) List(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	enterpriseID := claims.EnterpriseID
	if tenantSlug == "" || enterpriseID == 0 {
		response.BadRequest(c, "tenant not found")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	var categoryID, brandID *int64
	if catIDStr := c.Query("category_id"); catIDStr != "" {
		if id, err := strconv.ParseInt(catIDStr, 10, 64); err == nil {
			categoryID = &id
		}
	}
	if brandIDStr := c.Query("brand_id"); brandIDStr != "" {
		if id, err := strconv.ParseInt(brandIDStr, 10, 64); err == nil {
			brandID = &id
		}
	}

	filters := ListFilters{
		Page:       page,
		Limit:      limit,
		Search:     search,
		CategoryID: categoryID,
		BrandID:    brandID,
	}

	list, err := h.svc.List(c.Request.Context(), tenantSlug, enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

// GetByID handles GET /products/:id
// Returns a single product by ID
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

	product, err := h.svc.GetByID(c.Request.Context(), tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "product not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, product)
}

// GetBySKU handles GET /catalog/products/exist/:sku
// Checks if a product exists by SKU
func (h *Handler) GetBySKU(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	if tenantSlug == "" {
		response.BadRequest(c, "tenant not found")
		return
	}

	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	sku := c.Param("sku")
	if sku == "" {
		response.BadRequest(c, "sku is required")
		return
	}

	product, err := h.svc.GetBySKU(c.Request.Context(), tenantSlug, sku, enterpriseID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.OK(c, gin.H{
				"exists": false,
				"sku":    sku,
			})
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"exists": true,
		"sku":    sku,
		"product": gin.H{
			"id":      product.ID,
			"name":    product.Name,
			"sku":     product.SKU,
			"barcode": product.Barcode,
		},
	})
}

// Update handles PUT /products/:id
// Updates an existing product
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

	// Request structure matching the provided JSON
	var req struct {
		// Basic fields
		SKU         string `json:"sku"`
		Name        string `json:"name"`
		Barcode     string `json:"barcode"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`

		// Reference fields
		CategoryID *int64 `json:"category_id"`
		BrandID    *int64 `json:"brand_id"`
		UnitID     int64  `json:"unit_measure_id"`

		// Product type and status
		ProductType string `json:"product_type"`
		Active      *bool  `json:"active"`
		Status      string `json:"status"`

		// Pricing
		CostPrice float64  `json:"cost_price"`
		SalePrice float64  `json:"sale_price"`
		Price2    float64  `json:"price_2"`
		Price3    *float64 `json:"price_3"`

		// Taxes
		IVAPercentage  float64 `json:"iva_percentage"`
		ConsumptionTax float64 `json:"consumption_tax"`

		// Inventory controls
		ManagesInventory   *bool `json:"manages_inventory"`
		ManagesBatches     *bool `json:"manages_batches"`
		ManagesSerial      *bool `json:"manages_serial"`
		AllowNegativeStock *bool `json:"allow_negative_stock"`
		MinStock           *int  `json:"min_stock"`

		// Visibility
		VisibleInPOS *bool `json:"visible_in_pos"`

		// Presentations
		Presentations []PresentationRequest `json:"presentations"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Handle boolean fields that could be omitted
	var active, visibleInPOS, managesInventory, managesBatches, managesSerial, allowNegative bool
	var activeSet, visibleSet, inventorySet, batchesSet, serialSet, negativeSet bool

	if req.Active != nil {
		active = *req.Active
		activeSet = true
	}
	if req.VisibleInPOS != nil {
		visibleInPOS = *req.VisibleInPOS
		visibleSet = true
	}
	if req.ManagesInventory != nil {
		managesInventory = *req.ManagesInventory
		inventorySet = true
	}
	if req.ManagesBatches != nil {
		managesBatches = *req.ManagesBatches
		batchesSet = true
	}
	if req.ManagesSerial != nil {
		managesSerial = *req.ManagesSerial
		serialSet = true
	}
	if req.AllowNegativeStock != nil {
		allowNegative = *req.AllowNegativeStock
		negativeSet = true
	}

	// Create product entity
	product := &Product{
		SKU:            req.SKU,
		Name:           req.Name,
		Barcode:        req.Barcode,
		Description:    req.Description,
		ImageURL:       req.ImageURL,
		CategoryID:     req.CategoryID,
		BrandID:        req.BrandID,
		UnitID:         req.UnitID,
		ProductType:    req.ProductType,
		CostPrice:      req.CostPrice,
		SalePrice:      req.SalePrice,
		Price2:         req.Price2,
		Price3:         req.Price3,
		IVAPercentage:  req.IVAPercentage,
		ConsumptionTax: req.ConsumptionTax,
	}

	// Handle active (using false as "not provided" indicator)
	if activeSet {
		product.Active = active
	}

	// Handle visibleInPOS
	if visibleSet {
		product.VisibleInPOS = visibleInPOS
	}

	// Handle inventory controls
	if inventorySet {
		product.ManagesInventory = managesInventory
	}
	if batchesSet {
		product.ManagesBatches = managesBatches
	}
	if serialSet {
		product.ManagesSerial = managesSerial
	}
	if negativeSet {
		product.AllowNegativeStock = allowNegative
	}

	// Handle MinStock
	if req.MinStock != nil {
		product.MinStock = *req.MinStock
	}

	// Handle Presentations
	product.Presentations = req.Presentations

	if err := h.svc.Update(c.Request.Context(), tenantSlug, id, product); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "product not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, product)
}

// Delete handles DELETE /products/:id
// Performs a soft delete of a product
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
			response.NotFound(c, "product not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	c.Status(204)
}
