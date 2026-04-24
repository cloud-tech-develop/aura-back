package enterprise

import "github.com/gin-gonic/gin"

// Register mounts all enterprise routes onto the given router group.
// Public routes come before the auth middleware.
func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Public — no auth required (for offline sync)
	public.POST("/enterprises", h.Create)
	public.GET("/enterprises/:slug", h.GetBySlug)

	// Protected — behind AuthMiddleware
	protected.GET("/enterprises", h.List)
	protected.PUT("/enterprises/:slug", h.Update)
	protected.PATCH("/enterprises/:slug/status", h.UpdateStatus)
}

// RegisterPublic mounts enterprise routes for public access (offline mode)
func RegisterPublic(public gin.IRouter, h *Handler) {
	public.POST("/enterprises", h.Create)
	public.GET("/enterprises/:slug", h.GetBySlug)
}
