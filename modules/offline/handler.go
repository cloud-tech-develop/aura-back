package offline

import (
	"os"

	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/gin-gonic/gin"
)

// Handler handles offline sync requests
type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Ping handles GET /offline/ping
// This endpoint only works in offline mode (SQLite)
// Gets the enterprise slug from the JWT token (set by AuthMiddleware)
func (h *Handler) Ping(c *gin.Context) {
	driver := os.Getenv("DATABASE_DRIVER")
	dsn := os.Getenv("DATABASE_URL")

	isOffline := driver == "sqlite" || dsn == ""
	if !isOffline {
		isOffline = dsn == ""
	}

	if !isOffline {
		response.Forbidden(c, "Endpoint solo disponible en modo offline")
		return
	}

	slug, ok := tenant.SlugFromContext(c)
	if !ok || slug == "" {
		response.BadRequest(c, "No se pudo obtener el slug del token")
		return
	}

	prodURL := os.Getenv("URL_PROD")
	if prodURL == "" {
		prodURL = "http://localhost:8081"
	}

	token := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			token = authHeader
		}
	}

	result, err := h.svc.SyncAllBySlug(c.Request.Context(), prodURL, token, slug)
	if err != nil {
		response.OK(c, "Error al sincronizar: "+err.Error())
		return
	}

	response.OK(c, gin.H{
		"slug":    slug,
		"source":  prodURL,
		"mode":    "offline",
		"result":  result,
		"message": "Sincronización completada",
	})
}

// SyncTenant handles GET /offline/sync-tenant
// Sincroniza los datos del tenant desde producción
func (h *Handler) SyncTenant(c *gin.Context) {
	driver := os.Getenv("DATABASE_DRIVER")
	dsn := os.Getenv("DATABASE_URL")

	isOffline := driver == "sqlite" || dsn == ""
	if !isOffline {
		isOffline = dsn == ""
	}

	if !isOffline {
		response.Forbidden(c, "Endpoint solo disponible en modo offline")
		return
	}

	slug, ok := tenant.SlugFromContext(c)
	if !ok || slug == "" {
		response.BadRequest(c, "No se pudo obtener el slug del token")
		return
	}

	prodURL := os.Getenv("URL_PROD")
	if prodURL == "" {
		prodURL = "http://localhost:8081"
	}

	token := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			token = authHeader
		}
	}

	result, err := h.svc.SyncTenantBySlug(c.Request.Context(), prodURL, token, slug)
	if err != nil {
		response.OK(c, "Error al sincronizar tenant: "+err.Error())
		return
	}

	response.OK(c, gin.H{
		"slug":    slug,
		"source":  prodURL,
		"mode":    "offline",
		"result":  result,
		"message": "Sincronización de tenant completada",
	})
}

// ListEnterprises handles GET /offline/enterprises
func (h *Handler) ListEnterprises(c *gin.Context) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" && os.Getenv("DATABASE_DRIVER") != "sqlite" {
		response.Forbidden(c, "Endpoint solo disponible en modo offline")
		return
	}

	enterprises, err := h.svc.GetLocalEnterprises(c.Request.Context())
	if err != nil {
		response.BadRequest(c, "Error al listar empresas: "+err.Error())
		return
	}

	response.OK(c, gin.H{
		"data":   enterprises,
		"total":  len(enterprises),
		"source": "local",
	})
}
