package payments

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Payment routes
	protected.POST("/payments", h.ProcessPayment)
	protected.POST("/payments/batch", h.ProcessMultiplePayments)
	protected.GET("/payments", h.ListPayments)
	protected.GET("/payments/:id", h.GetPayment)
	protected.GET("/payments/reference/:referenceType/:referenceId", h.GetPaymentsByOrder)
	protected.POST("/payments/:id/cancel", h.CancelPayment)

	// Cash drawer routes
	protected.POST("/cash-drawers", h.OpenCashDrawer)
	protected.GET("/cash-drawers", h.ListCashDrawers)
	protected.GET("/cash-drawers/open", h.GetOpenDrawer)
	protected.GET("/cash-drawers/:id", h.GetCashDrawer)
	protected.POST("/cash-drawers/:id/close", h.CloseCashDrawer)
	protected.POST("/cash-drawers/:id/cash-in", h.AddCashIn)
	protected.POST("/cash-drawers/:id/cash-out", h.AddCashOut)
}
