package commissions

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Commission Rules
	protected.POST("/commissions/rules", h.CreateRule)
	protected.PUT("/commissions/rules/:id", h.UpdateRule)
	protected.DELETE("/commissions/rules/:id", h.DeleteRule)
	protected.GET("/commissions/rules", h.ListRules)

	// Commission Calculation
	protected.POST("/commissions/calculate", h.CalculateCommissions)

	// Commission History
	protected.GET("/commissions", h.ListCommissions)
	protected.GET("/commissions/:id", h.GetCommission)

	// Commission Settlement
	protected.POST("/commissions/settle", h.SettleCommissions)

	// Reports
	protected.GET("/commissions/report/summary", h.GetCommissionReport)
	protected.GET("/commissions/report/totals", h.GetCommissionTotals)
}
