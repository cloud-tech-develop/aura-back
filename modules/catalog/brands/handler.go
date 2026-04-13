package brands

import (
	"database/sql"
	"strconv"

	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	brand := &Brand{
		Name:         req.Name,
		Description:  req.Description,
		EnterpriseID: enterpriseID,
	}

	if err := h.svc.Create(c.Request.Context(), brand); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, brand)
}

func (h *Handler) List(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	list, err := h.svc.List(c.Request.Context(), enterpriseID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	brand, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Marca no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, brand)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	brand := &Brand{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.svc.Update(c.Request.Context(), id, brand); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Marca no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, brand)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Marca no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	c.Status(204)
}
