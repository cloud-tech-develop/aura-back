package offline

import "github.com/gin-gonic/gin"

// Register mounts offline sync routes onto the given router group.
// These routes only work in offline mode (SQLite)
func Register(router gin.IRouter, h *Handler) {
	// Offline sync endpoints
	router.GET("/offline/ping", h.Ping)
	router.GET("/offline/enterprises", h.ListEnterprises)
}
