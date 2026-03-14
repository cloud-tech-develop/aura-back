package enterprise

import (
	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
	"github.com/cloud-tech-develop/aura-back/shared/errors"
	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc      Service
	tenantMg *tenant.Manager
}

func NewHandler(svc Service, tenantMg *tenant.Manager) *Handler {
	return &Handler{svc: svc, tenantMg: tenantMg}
}

// ─── POST /enterprises ────────────────────────────────────────────────────────

type createRequest struct {
	Password       string      `json:"password"        binding:"required"`
	Name           string      `json:"name"            binding:"required"`
	CommercialName string      `json:"commercial_name"`
	Slug           string      `json:"slug"            binding:"required"`
	SubDomain      string      `json:"sub_domain"`
	Email          vo.Email    `json:"email"           binding:"required"`
	DV             vo.Document `json:"dv"`
	Phone          string      `json:"phone"`
	MunicipalityID string      `json:"municipality_id"`
	Municipality   string      `json:"municipality"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, errors.ErrFieldsRequired.Error()+".\n Error: "+err.Error())
		return
	}

	hash, err := tenant.HashPassword(req.Password)
	if err != nil {
		response.InternalServerError(c, "error al encriptar password")
		return
	}

	e := &Enterprise{
		Name:           req.Name,
		CommercialName: req.CommercialName,
		Slug:           req.Slug,
		SubDomain:      req.SubDomain,
		Email:          req.Email,
		DV:             req.DV,
		Phone:          req.Phone,
		MunicipalityID: req.MunicipalityID,
		Municipality:   req.Municipality,
	}

	if err := h.svc.Create(c.Request.Context(), e, hash); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, e)
}

// ─── GET /enterprises ─────────────────────────────────────────────────────────

func (h *Handler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, list)
}

// ─── GET /enterprises/:slug ───────────────────────────────────────────────────

func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	ent, err := h.svc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		response.BadRequest(c, errors.ErrEnterpriseNotFound.Error())
		return
	}
	response.OK(c, ent)
}
