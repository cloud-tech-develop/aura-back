package inventory

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Inventory routes
	protected.GET("/inventory", h.ListInventory)
	protected.GET("/inventory/low-stock", h.GetLowStock)
	protected.GET("/inventory/:productId/:branchId", h.GetInventory)
	protected.GET("/inventory/product/:productId", h.GetInventoryByProduct)
	protected.GET("/inventory/kardex/:productId/:branchId", h.GetProductKardex)

	// Movement routes
	protected.POST("/inventory/movements", h.UpdateStock)
	protected.GET("/movements", h.ListMovements)
	protected.GET("/movements/:id", h.GetMovement)

	// Movement reasons
	protected.GET("/movement-reasons", h.ListReasons)
}
