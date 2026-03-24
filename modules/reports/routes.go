package reports

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	protected.GET("/reports/sales-summary", h.GetSalesSummary)
	protected.GET("/reports/product-sales", h.GetProductSales)
	protected.GET("/reports/payment-methods", h.GetPaymentMethodBreakdown)
	protected.GET("/reports/daily-sales", h.GetDailySales)
	protected.GET("/reports/top-customers", h.GetTopCustomers)

	protected.GET("/reports/sales", h.GetSalesByPeriod)
	protected.GET("/reports/sales/products", h.GetSalesByProduct)
	protected.GET("/reports/sales/employees", h.GetSalesByEmployee)
	protected.GET("/reports/inventory", h.GetInventoryStatus)
	protected.GET("/reports/inventory/movements", h.GetMovementHistory)

	protected.POST("/reports/:reportType/export/pdf", h.ExportToPDF)
	protected.POST("/reports/:reportType/export/excel", h.ExportToExcel)
}
