package offline

import "github.com/gin-gonic/gin"

// Register mounts offline sync routes onto the given router group.
func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Offline sync endpoints (protected - requires token with slug in JWT)
	protected.GET("/offline/ping", h.Ping)
	protected.GET("/offline/enterprises", h.ListEnterprises)
}
