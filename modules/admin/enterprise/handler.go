package enterprise

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
	errs "github.com/cloud-tech-develop/aura-back/shared/errors"
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
	Document       vo.Document `json:"document"`
	DV             string      `json:"dv"`
	Phone          string      `json:"phone"`
	MunicipalityID string      `json:"municipality_id"`
	Municipality   string      `json:"municipality"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, errs.ErrFieldsRequired.Error()+".\n Error: "+err.Error())
		return
	}

	// Validate password minimum length (HU-001)
	if len(req.Password) < MinPasswordLength {
		response.BadRequest(c, "la contraseña debe tener al menos 8 caracteres")
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
		Document:       req.Document,
		DV:             req.DV,
		Phone:          req.Phone,
		MunicipalityID: req.MunicipalityID,
		Municipality:   req.Municipality,
	}

	if err := h.svc.Create(c.Request.Context(), e, hash); err != nil {
		// Return 403 Forbidden if plan limit reached (HU-008)
		if errors.Is(err, ErrPlanLimitReached) {
			response.Forbidden(c, err.Error())
			return
		}
		// Return 409 Conflict for duplicate entries
		if isConflictError(err) {
			response.Conflict(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, e)
}

// isConflictError returns true if the error is a conflict (duplicate)
func isConflictError(err error) bool {
	errStr := err.Error()
	return contains(errStr, "ya está registrado") ||
		contains(errStr, "ya existe") ||
		contains(errStr, "duplicate")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func formatInt64(n int64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte(n%10) + '0'
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// ─── GET /enterprises ─────────────────────────────────────────────────────────

func (h *Handler) List(c *gin.Context) {
	// Parse query parameters
	page := 1
	limit := 10
	status := ""

	if p := c.Query("page"); p != "" {
		if _, err := fmt.Sscanf(p, "%d", &page); err != nil {
			page = 1
		}
	}
	if l := c.Query("limit"); l != "" {
		if _, err := fmt.Sscanf(l, "%d", &limit); err != nil {
			limit = 10
		}
	}
	if s := c.Query("status"); s != "" {
		status = s
	}

	params := ListParams{
		Page:   page,
		Limit:  limit,
		Status: status,
	}

	result, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, result)
}

// ─── GET /enterprises/:slug ───────────────────────────────────────────────────

func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	ent, err := h.svc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		// Return 404 if not found (HU-005)
		if err == sql.ErrNoRows {
			response.NotFound(c, "Empresa no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, ent)
}

// ─── PUT /enterprises/:slug ───────────────────────────────────────────────────

type updateRequest struct {
	Name           string `json:"name"`
	CommercialName string `json:"commercial_name"`
	SubDomain      string `json:"sub_domain"`
	Phone          string `json:"phone"`
	MunicipalityID string `json:"municipality_id"`
	Municipality   string `json:"municipality"`
}

func (h *Handler) Update(c *gin.Context) {
	slug := c.Param("slug")

	// Get existing enterprise
	ent, err := h.svc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Empresa no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Update only provided fields
	if req.Name != "" {
		ent.Name = req.Name
	}
	if req.CommercialName != "" {
		ent.CommercialName = req.CommercialName
	}
	if req.SubDomain != "" {
		ent.SubDomain = req.SubDomain
	}
	if req.Phone != "" {
		ent.Phone = req.Phone
	}
	if req.MunicipalityID != "" {
		ent.MunicipalityID = req.MunicipalityID
	}
	if req.Municipality != "" {
		ent.Municipality = req.Municipality
	}

	if err := h.svc.Update(c.Request.Context(), ent); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, ent)
}

// ─── PATCH /enterprises/:slug/status ───────────────────────────────────────────

type statusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	slug := c.Param("slug")

	// Validate status
	validStatuses := map[string]bool{
		"ACTIVE":    true,
		"INACTIVE":  true,
		"SUSPENDED": true,
		"DEBT":      true,
	}

	var req statusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if !validStatuses[req.Status] {
		response.BadRequest(c, "estado inválido. Estados válidos: ACTIVE, INACTIVE, SUSPENDED, DEBT")
		return
	}

	// Get existing enterprise
	ent, err := h.svc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Empresa no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	ent.Status = req.Status

	if err := h.svc.Update(c.Request.Context(), ent); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"id":        ent.ID,
		"slug":      ent.Slug,
		"status":    ent.Status,
		"updatedAt": ent.UpdatedAt,
	})
}

// ─── GET /plans?enterprise_id=X ─────────────────────────────────────────────────────

func (h *Handler) GetPlans(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		// Try to get from query param
		enterpriseID, _ = strconv.ParseInt(c.Query("enterprise_id"), 10, 64)
	}
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id es requerido")
		return
	}

	plans, err := h.svc.GetPlansByEnterpriseID(c.Request.Context(), enterpriseID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, plans)
}
