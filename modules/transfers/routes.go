package transfers

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Transfers
	protected.POST("/transfers", h.CreateTransfer)
	protected.GET("/transfers/:id", h.GetTransfer)
	protected.GET("/transfers", h.ListTransfers)
	protected.POST("/transfers/:id/approve", h.ApproveTransfer)
	protected.POST("/transfers/:id/ship", h.ShipTransfer)
	protected.POST("/transfers/:id/receive", h.ReceiveTransfer)
	protected.POST("/transfers/:id/cancel", h.CancelTransfer)
}
