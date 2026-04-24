package offline

import (
	"os"

	"github.com/cloud-tech-develop/aura-back/shared/response"
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
func (h *Handler) Ping(c *gin.Context) {
	// Verify we are in offline mode (SQLite)
	driver := os.Getenv("DATABASE_DRIVER")
	dsn := os.Getenv("DATABASE_URL")

	isOffline := driver == "sqlite" || dsn == ""

	if !isOffline {
		// Check if DATABASE_URL is empty (offline mode in main.go)
		isOffline = dsn == ""
	}

	if !isOffline {
		response.Forbidden(c, "Endpoint solo disponible en modo offline")
		return
	}

	// Get production URL from environment
	prodURL := os.Getenv("URL_PROD")
	if prodURL == "" {
		prodURL = "http://localhost:8081" // fallback to default
	}

	// Get token from Authorization header
	token := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// Extract Bearer token
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			token = authHeader
		}
	}

	// Debug: log token presence
	if token == "" {
		response.BadRequest(c, "Token de autenticación requerido. Use header 'Authorization: Bearer <token>'")
		return
	}

	// Sync enterprises from production
	saved, err := h.svc.SyncEnterprises(c.Request.Context(), prodURL, token)
	if err != nil {
		response.BadRequest(c, "Error al sincronizar: "+err.Error())
		return
	}

	// Return sync result
	response.OK(c, gin.H{
		"synced":  saved,
		"source":  prodURL,
		"mode":    "offline",
		"message": "Sincronización completada",
	})
}

// ListEnterprises handles GET /offline/enterprises
// Returns all enterprises stored locally in SQLite
func (h *Handler) ListEnterprises(c *gin.Context) {
	// Verify we are in offline mode (SQLite)
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
