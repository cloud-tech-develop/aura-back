package brands

import "github.com/gin-gonic/gin"

func Register(protected gin.IRouter, h *Handler) {
	// Public for offline sync
	protected.GET("/catalog/brands/all", h.ListAll)
	protected.GET("/catalog/brands", h.List)
	protected.POST("/catalog/brands", h.Create)
	protected.GET("/catalog/brands/:id", h.GetByID)
	protected.POST("/catalog/brands/page", h.Page)
	protected.PUT("/catalog/brands/:id", h.Update)
	protected.DELETE("/catalog/brands/:id", h.Delete)
}
