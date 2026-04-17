package presentations

import (
	"github.com/gin-gonic/gin"
)

// Register registers the presentation routes
// public: public router (no auth)
// protected: protected router (with auth)
// h: handler instance
func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Public routes (no authentication required)
	public.GET("/presentations", h.List)
	public.GET("/presentations/:id", h.GetByID)

	// Protected routes (authentication required)
	protected.POST("/presentations/page", h.Page)
	protected.PUT("/presentations/:id", h.Update)
	protected.DELETE("/presentations/:id", h.Delete)
}
