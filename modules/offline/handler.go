package offline

import (
	"os"

	"github.com/cloud-tech-develop/aura-back/shared/logging"
	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/gin-gonic/gin"
)

// Handler handles offline sync requests
type Handler struct {
	svc    Service
	logger *logging.LoggerHandler
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc, logger: logging.NewLoggerHandler("logs")}
}

// Ping handles GET /offline/ping
// This endpoint only works in offline mode (SQLite)
// Gets the enterprise slug from JWT token OR from local SQLite database
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

	var slug string
	var token string
	var syncSource string // "token" o "local_db"

	// Try to get slug from JWT token first
	slugFromToken, ok := tenant.SlugFromContext(c)
	if ok && slugFromToken != "" {
		slug = slugFromToken
		syncSource = "token"
	} else {
		// Fallback: get slug from local SQLite enterprises table
		enterprise, err := h.svc.GetActiveEnterprise(c.Request.Context())
		if err != nil {
			response.BadRequest(c, "No hay enterprise configurada. Primero sincronice con /offline/sync-tenant")
			return
		}
		slug = enterprise.Slug
		syncSource = "local_db"

		h.logger.Logf("[offline.Handler] No se encontró slug en token, usando slug local: %s", slug)
	}

	prodURL := os.Getenv("URL_PROD")
	if prodURL == "" {
		prodURL = "http://localhost:8081"
	}

	// Get token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			token = authHeader
		}
	}

	h.logger.Logf("[offline.Handler] Iniciando sincronizacion para slug: %s (fuente: %s)", slug, syncSource)

	result, err := h.svc.SyncAllBySlug(c.Request.Context(), prodURL, token, slug)
	if err != nil {
		response.OK(c, "Error al sincronizar: "+err.Error())
		return
	}

	// After successful sync, activate RabbitMQ with the tenant slug
	// This will subscribe to tenant-specific events
	if err := h.svc.ActivateRabbitMQ(c.Request.Context(), slug); err != nil {
		h.logger.Logf("[offline.Handler] Warning: Failed to activate RabbitMQ: %v", err)
		// Don't fail the response - sync was successful
	}

	response.OK(c, gin.H{
		"slug":        slug,
		"source":      prodURL,
		"mode":        "offline",
		"sync_source": syncSource,
		"result":      result,
		"message":     "Sincronización completada",
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
