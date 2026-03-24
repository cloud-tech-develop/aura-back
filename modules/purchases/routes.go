package purchases

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Purchase Orders
	protected.POST("/purchases/orders", h.CreatePurchaseOrder)
	protected.GET("/purchases/orders/:id", h.GetPurchaseOrder)
	protected.GET("/purchases/orders", h.ListPurchaseOrders)

	// Goods Receipt
	protected.POST("/purchases/receive", h.ReceiveGoods)

	// Purchases
	protected.GET("/purchases/:id", h.GetPurchase)
	protected.GET("/purchases", h.ListPurchases)
	protected.POST("/purchases/:id/cancel", h.CancelPurchase)

	// Payments
	protected.POST("/purchases/payments", h.RecordPayment)

	// Supplier Summary
	protected.GET("/purchases/suppliers/:supplierID/summary", h.GetSupplierSummary)
}
