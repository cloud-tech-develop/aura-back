package shrinkage

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Shrinkage
	protected.POST("/shrinkage", h.RegisterShrinkage)
	protected.GET("/shrinkage/:id", h.GetShrinkage)
	protected.GET("/shrinkage", h.ListShrinkages)
	protected.POST("/shrinkage/:id/authorize", h.AuthorizeShrinkage)
	protected.POST("/shrinkage/:id/cancel", h.CancelShrinkage)

	// Shrinkage Reasons
	protected.POST("/shrinkage/reasons", h.CreateReason)
	protected.GET("/shrinkage/reasons", h.ListReasons)

	// Reporting
	protected.GET("/shrinkage/report", h.GetShrinkageReport)
}
