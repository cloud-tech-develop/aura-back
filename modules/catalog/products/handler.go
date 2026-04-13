package products

import (
	"database/sql"
	"strconv"

	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok {
		response.BadRequest(c, "tenant no encontrado")
		return
	}

	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	var req struct {
		SKU         string  `json:"sku" binding:"required"`
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		CategoryID  *int64  `json:"category_id"`
		BrandID     *int64  `json:"brand_id"`
		CostPrice   float64 `json:"cost_price" binding:"required"`
		SalePrice   float64 `json:"sale_price" binding:"required"`
		TaxRate     float64 `json:"tax_rate"`
		MinStock    int     `json:"min_stock"`
		ImageURL    string  `json:"image_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	product := &Product{
		SKU:          req.SKU,
		Name:         req.Name,
		Description:  req.Description,
		CategoryID:   req.CategoryID,
		BrandID:      req.BrandID,
		CostPrice:    req.CostPrice,
		SalePrice:    req.SalePrice,
		TaxRate:      req.TaxRate,
		MinStock:     req.MinStock,
		ImageURL:     req.ImageURL,
		EnterpriseID: enterpriseID,
	}

	if err := h.svc.Create(c.Request.Context(), tenantSlug, product); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, product)
}

func (h *Handler) Page(c *gin.Context) {
	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok {
		response.BadRequest(c, "tenant no encontrado")
		return
	}

	var req struct {
		First  int64  `json:"first"`
		Rows   int64  `json:"rows"`
		Search string `json:"search"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	result, err := h.svc.Page(c.Request.Context(), tenantSlug, enterpriseID, req.First, req.Rows, req.Search)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, result)
}

func (h *Handler) List(c *gin.Context) {
	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok {
		response.BadRequest(c, "tenant no encontrado")
		return
	}

	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
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

func (h *Handler) GetByID(c *gin.Context) {
	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok {
		response.BadRequest(c, "tenant no encontrado")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	product, err := h.svc.GetByID(c.Request.Context(), tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Producto no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, product)
}

func (h *Handler) Update(c *gin.Context) {
	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok {
		response.BadRequest(c, "tenant no encontrado")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	var req struct {
		SKU         string  `json:"sku"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		CategoryID  *int64  `json:"category_id"`
		BrandID     *int64  `json:"brand_id"`
		CostPrice   float64 `json:"cost_price"`
		SalePrice   float64 `json:"sale_price"`
		TaxRate     float64 `json:"tax_rate"`
		MinStock    int     `json:"min_stock"`
		ImageURL    string  `json:"image_url"`
		Status      string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	product := &Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		BrandID:     req.BrandID,
		CostPrice:   req.CostPrice,
		SalePrice:   req.SalePrice,
		TaxRate:     req.TaxRate,
		MinStock:    req.MinStock,
		ImageURL:    req.ImageURL,
		Status:      req.Status,
	}

	if err := h.svc.Update(c.Request.Context(), tenantSlug, id, product); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Producto no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, product)
}

func (h *Handler) Delete(c *gin.Context) {
	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok {
		response.BadRequest(c, "tenant no encontrado")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), tenantSlug, id); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Producto no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	c.Status(204)
}
