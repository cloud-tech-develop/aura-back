package categories

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
		Name           string  `json:"name" binding:"required"`
		ParentID       *int64  `json:"parent_id"`
		DefaultTaxRate float64 `json:"default_tax_rate"`
		Active         *bool   `json:"active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	category := &Category{
		Name:           req.Name,
		ParentID:       req.ParentID,
		DefaultTaxRate: req.DefaultTaxRate,
		Active:         active,
		EnterpriseID:   enterpriseID,
	}

	if err := h.svc.Create(c.Request.Context(), tenantSlug, category); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, category)
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

	list, err := h.svc.List(c.Request.Context(), tenantSlug, enterpriseID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
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

	category, err := h.svc.GetByID(c.Request.Context(), tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Categoría no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, category)
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
		Name           string  `json:"name"`
		ParentID       *int64  `json:"parent_id"`
		DefaultTaxRate float64 `json:"default_tax_rate"`
		Active         *bool   `json:"active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	category := &Category{
		Name:           req.Name,
		ParentID:       req.ParentID,
		DefaultTaxRate: req.DefaultTaxRate,
		Active:         active,
	}

	if err := h.svc.Update(c.Request.Context(), tenantSlug, id, category); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Categoría no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, category)
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
			response.NotFound(c, "Categoría no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	c.Status(204)
}
