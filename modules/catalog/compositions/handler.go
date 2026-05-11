package compositions

import (
	"database/sql"
	"net/http"
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

// Create handles POST /catalog/compositions
func (h *Handler) Create(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	enterpriseID := claims.EnterpriseID

	if tenantSlug == "" || enterpriseID == 0 {
		response.BadRequest(c, "tenant not found")
		return
	}

	var req struct {
		ParentProductID int64   `json:"parent_product_id" binding:"required"`
		ChildProductID  int64   `json:"child_product_id" binding:"required"`
		Quantity        float64 `json:"quantity"`
		Type            string  `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	composition := &Composition{
		ParentProductID: req.ParentProductID,
		ChildProductID:  req.ChildProductID,
		Quantity:        req.Quantity,
		Type:            req.Type,
		EnterpriseID:    enterpriseID,
	}

	if err := h.svc.Create(c.Request.Context(), tenantSlug, composition); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, composition)
}

// ListAll handles GET /catalog/compositions/all
func (h *Handler) ListAll(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	enterpriseID := claims.EnterpriseID

	if tenantSlug == "" || enterpriseID == 0 {
		response.BadRequest(c, "tenant not found")
		return
	}

	list, err := h.svc.ListAll(c.Request.Context(), tenantSlug, enterpriseID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

// List handles GET /catalog/compositions
func (h *Handler) List(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	enterpriseID := claims.EnterpriseID

	if tenantSlug == "" || enterpriseID == 0 {
		response.BadRequest(c, "tenant not found")
		return
	}

	// Optional filter by parent_product_id
	var parentID int64
	if pidStr := c.Query("parent_product_id"); pidStr != "" {
		if id, err := strconv.ParseInt(pidStr, 10, 64); err == nil {
			parentID = id
		}
	}

	if parentID > 0 {
		list, err := h.svc.ListByParent(c.Request.Context(), tenantSlug, parentID)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}
		response.OK(c, list)
		return
	}

	// No filter: return all for enterprise
	list, err := h.svc.ListAll(c.Request.Context(), tenantSlug, enterpriseID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, list)
}

// Page handles POST /catalog/compositions/page
func (h *Handler) Page(c *gin.Context) {
	claims, _ := tenant.ClaimsFromContext(c)
	tenantSlug := claims.Slug
	enterpriseID := claims.EnterpriseID

	if tenantSlug == "" || enterpriseID == 0 {
		response.BadRequest(c, "tenant not found")
		return
	}

	var req struct {
		Page   int64  `json:"page"`
		Limit  int64  `json:"limit"`
		Search string `json:"search"`
		Tipo   string `json:"tipo"`
		Sort   string `json:"sort"`
		Order  string `json:"order"`
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

	result, err := h.svc.Page(c.Request.Context(), tenantSlug, enterpriseID, req.Page, req.Limit, req.Search, req.Tipo, req.Sort, req.Order)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, result)
}

// GetByID handles GET /catalog/compositions/:id
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

	composition, err := h.svc.GetByID(c.Request.Context(), tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "composition not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, composition)
}

// Update handles PUT /catalog/compositions/:id
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
		Quantity float64 `json:"quantity"`
		Type     string  `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	composition := &Composition{
		Quantity: req.Quantity,
		Type:     req.Type,
	}

	if err := h.svc.Update(c.Request.Context(), tenantSlug, id, composition); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "composition not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, composition)
}

// Delete handles DELETE /catalog/compositions/:id
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
			response.NotFound(c, "composition not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
