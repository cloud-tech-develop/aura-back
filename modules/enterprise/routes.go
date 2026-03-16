package enterprise

import "github.com/gin-gonic/gin"

// Register mounts all enterprise routes onto the given router group.
// Public routes come before the auth middleware.
func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Public — no auth required
	public.POST("/enterprises", h.Create)

	// Protected — behind AuthMiddleware
	protected.GET("/enterprises", h.List)
	protected.GET("/enterprises/:slug", h.GetBySlug)
	protected.PUT("/enterprises/:slug", h.Update)
	protected.PATCH("/enterprises/:slug/status", h.UpdateStatus)
}
