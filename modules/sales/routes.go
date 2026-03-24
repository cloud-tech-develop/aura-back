package sales

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Sales order routes
	protected.POST("/sales-orders", h.CreateOrder)
	protected.GET("/sales-orders", h.ListOrders)
	protected.GET("/sales-orders/:id", h.GetOrder)
	protected.PUT("/sales-orders/:id/status", h.UpdateOrderStatus)
	protected.POST("/sales-orders/:id/cancel", h.CancelOrder)
	protected.POST("/sales-orders/:id/complete", h.CompleteOrder)
}
