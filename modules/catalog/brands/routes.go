package brands

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	protected.POST("/catalog/brands", h.Create)
	protected.GET("/catalog/brands", h.List)
	protected.GET("/catalog/brands/:id", h.GetByID)
	protected.POST("/catalog/brands/page", h.Page)
	protected.PUT("/catalog/brands/:id", h.Update)
	protected.DELETE("/catalog/brands/:id", h.Delete)
}
