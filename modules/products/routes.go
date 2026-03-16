package products

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Category routes
	protected.POST("/categories", h.CreateCategory)
	protected.GET("/categories", h.ListCategories)
	protected.GET("/categories/:id", h.GetCategory)
	protected.PUT("/categories/:id", h.UpdateCategory)

	// Brand routes
	protected.POST("/brands", h.CreateBrand)
	protected.GET("/brands", h.ListBrands)
	protected.GET("/brands/:id", h.GetBrand)
	protected.PUT("/brands/:id", h.UpdateBrand)

	// Product routes
	protected.POST("/products", h.CreateProduct)
	protected.GET("/products", h.ListProducts)
	protected.GET("/products/:id", h.GetProduct)
	protected.PUT("/products/:id", h.UpdateProduct)
	protected.DELETE("/products/:id", h.DeleteProduct)
}
