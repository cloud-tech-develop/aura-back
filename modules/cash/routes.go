package cash

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Cash Drawer
	protected.GET("/cash/drawer/:branchID", h.GetDrawerByBranch)
	protected.POST("/cash/drawer", h.ConfigureDrawer)

	// Cash Shift
	protected.POST("/cash/shift/open", h.OpenShift)
	protected.POST("/cash/shift/:shiftID/close", h.CloseShift)
	protected.GET("/cash/shift/active", h.GetActiveShift)
	protected.GET("/cash/shift/:shiftID", h.GetShiftSummary)
	protected.GET("/cash/shifts", h.ListShifts)

	// Cash Movement
	protected.POST("/cash/movement", h.RecordMovement)

	// Reconciliation
	protected.POST("/cash/shift/:shiftID/reconcile", h.ReconcileShift)
}
