package api

import (
	"net/http"

	"github.com/cloud-tech-develop/aura-back/application/enterprise"
	domainenterprise "github.com/cloud-tech-develop/aura-back/domain/enterprise"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/gin-gonic/gin"
)

type EnterpriseHandler struct {
	service *enterprise.Service
}

func NewEnterpriseHandler(service *enterprise.Service) *EnterpriseHandler {
	return &EnterpriseHandler{service: service}
}

func (h *EnterpriseHandler) Create(c *gin.Context) {
	var req struct {
		Password       string `json:"password" binding:"required"`
		Name           string `json:"name" binding:"required"`
		CommercialName string `json:"commercial_name"`
		Slug           string `json:"slug" binding:"required"`
		SubDomain      string `json:"sub_domain"`
		Email          string `json:"email" binding:"required"`
		DV             string `json:"dv"`
		Phone          string `json:"phone"`
		MunicipalityID string `json:"municipality_id"`
		Municipality   string `json:"municipality"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos requeridos faltantes"})
		return
	}

	hash, err := tenant.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al encriptar password"})
		return
	}

	e := &domainenterprise.Enterprise{
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

	if err := h.service.Create(c.Request.Context(), e, hash); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "enterprise " + e.Slug + " creado"})
}
