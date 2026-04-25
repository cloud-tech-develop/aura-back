package thirdparties

import (
	"database/sql"
	"strconv"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(database *db.DB) *Handler {
	q := database.Wrap(database.DB)
	return &Handler{svc: NewService(q)}
}

// CreateThirdParty - POST /third-parties
func (h *Handler) CreateThirdParty(c *gin.Context) {
	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok || tenantSlug == "" {
		tenantSlug = c.Query("slug")
	}

	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	var req struct {
		UserID            *int64 `json:"user_id"`
		FirstName         string `json:"first_name" binding:"required"`
		LastName          string `json:"last_name" binding:"required"`
		DocumentNumber    string `json:"document_number" binding:"required"`
		DocumentType      string `json:"document_type" binding:"required"`
		PersonalEmail     string `json:"personal_email"`
		CommercialName    string `json:"commercial_name"`
		Address           string `json:"address"`
		Phone             string `json:"phone"`
		AdditionalEmail   string `json:"additional_email"`
		TaxResponsibility string `json:"tax_responsibility" binding:"required"`
		IsClient          bool   `json:"is_client"`
		IsProvider        bool   `json:"is_provider"`
		IsEmployee        bool   `json:"is_employee"`
		MunicipalityID    string `json:"municipality_id"`
		Municipality      string `json:"municipality"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tp := &ThirdParty{
		UserID:            req.UserID,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		DocumentNumber:    req.DocumentNumber,
		DocumentType:      req.DocumentType,
		PersonalEmail:     req.PersonalEmail,
		CommercialName:    req.CommercialName,
		Address:           req.Address,
		Phone:             req.Phone,
		AdditionalEmail:   req.AdditionalEmail,
		TaxResponsibility: req.TaxResponsibility,
		IsClient:          req.IsClient,
		IsProvider:        req.IsProvider,
		IsEmployee:        req.IsEmployee,
		MunicipalityID:    req.MunicipalityID,
		Municipality:      req.Municipality,
	}

	if err := h.svc.Create(c.Request.Context(), tenantSlug, tp); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, tp)
}

// ListThirdParties - GET /third-parties
// Supports enterprise_id from context (JWT) or query params (for offline sync)
func (h *Handler) ListThirdParties(c *gin.Context) {
	tenantSlug, hasSlug := tenant.SlugFromContext(c)
	if !hasSlug || tenantSlug == "" {
		tenantSlug = c.Query("slug")
	}

	enterpriseID := c.GetInt64("enterprise_id")
	
	// Try to get from query params for offline sync (no JWT context)
	if enterpriseID == 0 {
		idStr := c.Query("enterprise_id")
		if idStr != "" {
			var err error
			enterpriseID, err = strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				response.BadRequest(c, "enterprise_id inválido")
				return
			}
		}
	}

	if enterpriseID == 0 || tenantSlug == "" {
		response.BadRequest(c, "enterprise_id y slug son requeridos")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	typeFilter := c.Query("type")
	search := c.Query("search")

	filters := ThirdPartyFilters{
		Page:   page,
		Limit:  limit,
		Type:   typeFilter,
		Search: search,
	}

	list, err := h.svc.List(c.Request.Context(), tenantSlug, enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

// GetThirdParty - GET /third-parties/:id
func (h *Handler) GetThirdParty(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok || tenantSlug == "" {
		tenantSlug = c.Query("slug")
	}

	tp, err := h.svc.GetByID(c.Request.Context(), tenantSlug, id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Tercero no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, tp)
}

// GetThirdPartyByDocument - GET /third-parties/document/:documentNumber
func (h *Handler) GetThirdPartyByDocument(c *gin.Context) {
	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok || tenantSlug == "" {
		tenantSlug = c.Query("slug")
	}

	docNumber := c.Param("documentNumber")

	tp, err := h.svc.GetByDocument(c.Request.Context(), tenantSlug, docNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Tercero no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, tp)
}

// UpdateThirdParty - PUT /third-parties/:id
func (h *Handler) UpdateThirdParty(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	var req struct {
		UserID            *int64 `json:"user_id"`
		FirstName         string `json:"first_name"`
		LastName          string `json:"last_name"`
		DocumentType      string `json:"document_type"`
		PersonalEmail     string `json:"personal_email"`
		CommercialName    string `json:"commercial_name"`
		Address           string `json:"address"`
		Phone             string `json:"phone"`
		AdditionalEmail   string `json:"additional_email"`
		TaxResponsibility string `json:"tax_responsibility"`
		IsClient          *bool  `json:"is_client"`
		IsProvider        *bool  `json:"is_provider"`
		IsEmployee        *bool  `json:"is_employee"`
		MunicipalityID    string `json:"municipality_id"`
		Municipality      string `json:"municipality"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tp := &ThirdParty{
		UserID:            req.UserID,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		DocumentType:      req.DocumentType,
		PersonalEmail:     req.PersonalEmail,
		CommercialName:    req.CommercialName,
		Address:           req.Address,
		Phone:             req.Phone,
		AdditionalEmail:   req.AdditionalEmail,
		TaxResponsibility: req.TaxResponsibility,
		MunicipalityID:    req.MunicipalityID,
		Municipality:      req.Municipality,
	}

	// Handle boolean pointers
	if req.IsClient != nil {
		tp.IsClient = *req.IsClient
	}
	if req.IsProvider != nil {
		tp.IsProvider = *req.IsProvider
	}
	if req.IsEmployee != nil {
		tp.IsEmployee = *req.IsEmployee
	}

	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok || tenantSlug == "" {
		tenantSlug = c.Query("slug")
	}

	if err := h.svc.Update(c.Request.Context(), tenantSlug, id, tp); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Tercero no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	tp, _ = h.svc.GetByID(c.Request.Context(), tenantSlug, id)
	response.OK(c, tp)
}

// DeleteThirdParty - DELETE /third-parties/:id
func (h *Handler) DeleteThirdParty(c *gin.Context) {
	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok || tenantSlug == "" {
		tenantSlug = c.Query("slug")
	}

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), tenantSlug, id); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Tercero no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	c.Status(204)
}
