package units

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Public for offline sync
	public.GET("/catalog/units", h.List)
	
	// Protected routes
	protected.POST("/catalog/units", h.Create)
	protected.GET("/catalog/units/:id", h.GetByID)
	protected.POST("/catalog/units/page", h.Page)
	protected.PUT("/catalog/units/:id", h.Update)
	protected.DELETE("/catalog/units/:id", h.Delete)
}
