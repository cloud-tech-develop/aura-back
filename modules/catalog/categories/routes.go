package categories

import "github.com/gin-gonic/gin"

func Register(protected gin.IRouter, h *Handler) {
	protected.GET("/catalog/categories/all", h.ListAll)
	protected.GET("/catalog/categories", h.List)
	protected.POST("/catalog/categories", h.Create)
	protected.GET("/catalog/categories/:id", h.GetByID)
	protected.POST("/catalog/categories/page", h.Page)
	protected.PUT("/catalog/categories/:id", h.Update)
	protected.DELETE("/catalog/categories/:id", h.Delete)
}
