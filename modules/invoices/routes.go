package invoices

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Invoice routes
	protected.POST("/invoices/generate", h.GenerateInvoiceFromSale)
	protected.POST("/invoices", h.CreateInvoice)
	protected.GET("/invoices", h.ListInvoices)
	protected.GET("/invoices/:id", h.GetInvoice)
	protected.GET("/invoices/number/:invoiceNumber", h.GetInvoiceByNumber)
	protected.POST("/invoices/:id/issue", h.IssueInvoice)
	protected.POST("/invoices/:id/cancel", h.CancelInvoice)
	protected.GET("/invoices/:id/logs", h.GetInvoiceLogs)

	// Invoice prefix routes
	protected.POST("/invoice-prefixes", h.CreateInvoicePrefix)
	protected.GET("/invoice-prefixes", h.ListInvoicePrefixes)
}
