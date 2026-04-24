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
// Syncs enterprises matching the slug pattern and all related data (plans, users, user_roles)
func (h *Handler) Ping(c *gin.Context) {
	// Verify we are in offline mode (SQLite)
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

	// Get slug from JWT token (set by AuthMiddleware)
	slug, ok := tenant.SlugFromContext(c)
	if !ok || slug == "" {
		response.BadRequest(c, "No se pudo obtener el slug del token")
		return
	}

	// Get production URL from environment
	prodURL := os.Getenv("URL_PROD")
	if prodURL == "" {
		prodURL = "http://localhost:8081"
	}

	// Get token from Authorization header
	token := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			token = authHeader
		}
	}

	// Sync all data by slug
	result, err := h.svc.SyncAllBySlug(c.Request.Context(), prodURL, token, slug)
	if err != nil {
		response.BadRequest(c, "Error al sincronizar: "+err.Error())
		return
	}

	// Return sync result
	response.OK(c, gin.H{
		"slug":    slug,
		"source":  prodURL,
		"mode":    "offline",
		"result":  result,
		"message": "Sincronización completada",
	})
}

// ListEnterprises handles GET /offline/enterprises
// Returns all enterprises stored locally in SQLite
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