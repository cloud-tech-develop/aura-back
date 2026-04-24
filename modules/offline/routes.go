package offline

import "github.com/gin-gonic/gin"

// Register mounts offline sync routes onto the given router group.
// These routes only work in offline mode (SQLite) and require authentication.
func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Offline sync endpoints (protected - requires token)
	protected.GET("/offline/ping", h.Ping)
	protected.GET("/offline/enterprises", h.ListEnterprises)
}
