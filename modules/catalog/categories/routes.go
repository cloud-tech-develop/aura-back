package categories

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	protected.POST("/catalog/categories", h.Create)
	protected.GET("/catalog/categories", h.List)
	protected.GET("/catalog/categories/:id", h.GetByID)
	protected.PUT("/catalog/categories/:id", h.Update)
	protected.DELETE("/catalog/categories/:id", h.Delete)
}
