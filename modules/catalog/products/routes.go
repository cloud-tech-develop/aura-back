package products

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	protected.POST("/catalog/products", h.Create)
	protected.GET("/catalog/products", h.List)
	protected.GET("/catalog/products/:id", h.GetByID)
	protected.POST("/catalog/products/page", h.Page)
	protected.PUT("/catalog/products/:id", h.Update)
	protected.DELETE("/catalog/products/:id", h.Delete)
}
