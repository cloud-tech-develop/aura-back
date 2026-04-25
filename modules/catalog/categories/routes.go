package categories

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Public for offline sync
	public.GET("/catalog/categories", h.List)
	
	// Protected routes
	protected.POST("/catalog/categories", h.Create)
	protected.GET("/catalog/categories/:id", h.GetByID)
	protected.POST("/catalog/categories/page", h.Page)
	protected.PUT("/catalog/categories/:id", h.Update)
	protected.DELETE("/catalog/categories/:id", h.Delete)
}
