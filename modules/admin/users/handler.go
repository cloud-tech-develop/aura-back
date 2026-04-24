package users

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
	"github.com/cloud-tech-develop/aura-back/shared/errors"
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

// ─── POST /users ────────────────────────────────────────────────────────

type createRequest struct {
	Name              string   `json:"name" binding:"required"`
	Email             vo.Email `json:"email" binding:"required"`
	Password          string   `json:"password" binding:"required"`
	FirstName         string   `json:"first_name"`
	LastName          string   `json:"last_name"`
	DocumentNumber    string   `json:"document_number"`
	DocumentType      string   `json:"document_type"`
	PersonalEmail     string   `json:"personal_email"`
	TaxResponsibility string   `json:"tax_responsibility"`
	IsEmployee        bool     `json:"is_employee"`
	Roles             []int64  `json:"roles"` // Role IDs to assign
}

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, errors.ErrFieldsRequired.Error()+".\n Error: "+err.Error())
		return
	}

	// Get tenant slug from context
	tenantSlug, ok := tenant.SlugFromContext(c)
	if !ok {
		response.BadRequest(c, "tenant no encontrado")
		return
	}

	// Get enterprise ID from context
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise ID no encontrado en contexto")
		return
	}

	// Validate password minimum length
	if len(req.Password) < 8 {
		response.BadRequest(c, "la contraseña debe tener al menos 8 caracteres")
		return
	}

	user := &User{
		EnterpriseID:      enterpriseID,
		Name:              req.Name,
		Email:             req.Email,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		DocumentNumber:    req.DocumentNumber,
		DocumentType:      req.DocumentType,
		PersonalEmail:     req.PersonalEmail,
		TaxResponsibility: req.TaxResponsibility,
		IsEmployee:        req.IsEmployee,
	}

	if err := h.svc.Create(c.Request.Context(), tenantSlug, user, req.Password); err != nil {
		if isConflictError(err) {
			response.Conflict(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	// Assign roles if provided
	if len(req.Roles) > 0 {
		roleLevel := c.GetInt("role_level")
		if err := h.svc.AssignRoles(c.Request.Context(), user.ID, req.Roles, roleLevel); err != nil {
			// Log error but don't fail the request? Or fail it?
			// For now, return error.
			response.BadRequest(c, fmt.Sprintf("usuario creado pero error al asignar roles: %v", err))
			return
		}
	}

	response.Created(c, user)
}

// ─── GET /users?enterprise_id=X (Public for offline sync) ─────────────────────

func (h *Handler) ListByEnterpriseID(c *gin.Context) {
	enterpriseID, _ := strconv.ParseInt(c.Query("enterprise_id"), 10, 64)
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id es requerido")
		return
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	users, err := h.svc.ListByEnterprise(c.Request.Context(), enterpriseID, page, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"data": users,
	})
}

// ─── GET /users ─────────────────────────────────────────────────────────

func (h *Handler) List(c *gin.Context) {
	// Get enterprise ID from context (current logged-in user's enterprise)
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise ID no encontrado en contexto")
		return
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, err := h.svc.ListByEnterprise(c.Request.Context(), enterpriseID, page, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, users)
}

// ─── GET /users/:id ─────────────────────────────────────────────────────

func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Usuario no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	// Security check: ensure user belongs to the current enterprise
	enterpriseID := c.GetInt64("enterprise_id")
	if user.EnterpriseID != enterpriseID {
		response.NotFound(c, "Usuario no encontrado") // Don't reveal existence
		return
	}

	response.OK(c, user)
}

// ─── PUT /users/:id ─────────────────────────────────────────────────────

type updateRequest struct {
	Name              string   `json:"name"`
	Email             vo.Email `json:"email"`
	FirstName         string   `json:"first_name"`
	LastName          string   `json:"last_name"`
	DocumentNumber    string   `json:"document_number"`
	DocumentType      string   `json:"document_type"`
	PersonalEmail     string   `json:"personal_email"`
	TaxResponsibility string   `json:"tax_responsibility"`
	IsEmployee        bool     `json:"is_employee"`
}

func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	// Get existing user
	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Usuario no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	// Security check
	enterpriseID := c.GetInt64("enterprise_id")
	if user.EnterpriseID != enterpriseID {
		response.NotFound(c, "Usuario no encontrado")
		return
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.DocumentNumber != "" {
		user.DocumentNumber = req.DocumentNumber
	}
	if req.DocumentType != "" {
		user.DocumentType = req.DocumentType
	}
	if req.PersonalEmail != "" {
		user.PersonalEmail = req.PersonalEmail
	}
	if req.TaxResponsibility != "" {
		user.TaxResponsibility = req.TaxResponsibility
	}
	// IsEmployee is bool, so we always update it if provided in JSON
	// But JSON omitempty might hide it. We assume false if not provided?
	// For simplicity, we update it.
	user.IsEmployee = req.IsEmployee

	if err := h.svc.Update(c.Request.Context(), user); err != nil {
		if isConflictError(err) {
			response.Conflict(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, user)
}

// ─── PATCH /users/:id/status ────────────────────────────────────────────

type statusRequest struct {
	Active bool `json:"active"`
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	// Security check: verify user exists and belongs to enterprise
	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Usuario no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	enterpriseID := c.GetInt64("enterprise_id")
	if user.EnterpriseID != enterpriseID {
		response.NotFound(c, "Usuario no encontrado")
		return
	}

	var req statusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateStatus(c.Request.Context(), id, req.Active); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Updated(c, gin.H{
		"id":     id,
		"active": req.Active,
	})
}

// ─── PATCH /users/:id/roles ─────────────────────────────────────────────

type rolesRequest struct {
	RoleIDs []int64 `json:"role_ids"`
}

func (h *Handler) AssignRoles(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	// Security check
	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Usuario no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	enterpriseID := c.GetInt64("enterprise_id")
	if user.EnterpriseID != enterpriseID {
		response.NotFound(c, "Usuario no encontrado")
		return
	}

	var req rolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get role level from JWT to validate permissions
	roleLevel := c.GetInt("role_level")

	if err := h.svc.AssignRoles(c.Request.Context(), id, req.RoleIDs, roleLevel); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"id":       id,
		"role_ids": req.RoleIDs,
	})
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

// ─── GET /roles ───────────────────────────────────────────────────────────

func (h *Handler) ListRoles(c *gin.Context) {
	// Get role level from JWT claims (set by middleware)
	roleLevel := c.GetInt("role_level")
	if roleLevel == 0 {
		// If not set, default to 0 (SUPERADMIN can see all)
		roleLevel = 0
	}

	roles, err := h.svc.ListRolesByMinLevel(c.Request.Context(), roleLevel)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, roles)
}

// ─── GET /user-roles?enterprise_id=X (Public for offline sync) ─────────────────────────────

func (h *Handler) ListUserRolesByEnterpriseID(c *gin.Context) {
	enterpriseID, _ := strconv.ParseInt(c.Query("enterprise_id"), 10, 64)
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id es requerido")
		return
	}

	userRoles, err := h.svc.ListUserRolesByEnterpriseID(c.Request.Context(), enterpriseID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"data": userRoles,
	})
}
