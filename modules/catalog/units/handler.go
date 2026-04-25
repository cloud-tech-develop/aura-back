package units

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
		Name          string `json:"name" binding:"required"`
		Abbreviation  string `json:"abbreviation" binding:"required"`
		Active        *bool  `json:"active"`
		AllowDecimals *bool  `json:"allow_decimals"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	allowDecimals := true
	if req.AllowDecimals != nil {
		allowDecimals = *req.AllowDecimals
	}

	unit := &Unit{
		Name:          req.Name,
		Abbreviation:  req.Abbreviation,
		Active:        active,
		AllowDecimals: allowDecimals,
		EnterpriseID:  enterpriseID,
	}

	if err := h.svc.Create(c.Request.Context(), tenantSlug, unit); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, unit)
}

func (h *Handler) List(c *gin.Context) {
	// Try to get slug from context (authenticated)
	tenantSlug, hasSlug := tenant.SlugFromContext(c)
	
	// If not in context, try to get from query param (public access for offline sync)
	if !hasSlug || tenantSlug == "" {
		tenantSlug = c.Query("slug")
	}
	
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		// Try from query param
		if idStr := c.Query("enterprise_id"); idStr != "" {
			var err error
			enterpriseID, err = strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				response.BadRequest(c, "enterprise_id inválido")
				return
			}
		}
	}
	
	// If no slug and no enterprise_id, require auth
	if tenantSlug == "" && enterpriseID == 0 && !hasSlug {
		response.BadRequest(c, "slug o enterprise_id es requerido")
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

	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	result, err := h.svc.Page(c.Request.Context(), tenantSlug, enterpriseID, req.Page, req.Limit, req.Search, req.Sort, req.Order, req.Params)
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

	unit, err := h.svc.GetByID(c.Request.Context(), tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Unidad no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, unit)
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
		Name          string `json:"name"`
		Abbreviation  string `json:"abbreviation"`
		Active        *bool  `json:"active"`
		AllowDecimals *bool  `json:"allow_decimals"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	unit := &Unit{
		Name:         req.Name,
		Abbreviation: req.Abbreviation,
	}

	if req.Active != nil {
		unit.Active = *req.Active
	}
	if req.AllowDecimals != nil {
		unit.AllowDecimals = *req.AllowDecimals
	}

	if err := h.svc.Update(c.Request.Context(), tenantSlug, id, unit); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Unidad no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, unit)
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
			response.NotFound(c, "Unidad no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	c.Status(204)
}
