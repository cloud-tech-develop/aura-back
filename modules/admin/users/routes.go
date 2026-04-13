package users

import "github.com/gin-gonic/gin"

// Register mounts all user routes onto the given router group.
func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Protected routes (require auth)
	protected.GET("/admin/users", h.List)
	protected.GET("/admin/users/:id", h.GetByID)
	protected.POST("/admin/users", h.Create)
	protected.PUT("/admin/users/:id", h.Update)
	protected.PATCH("/admin/users/:id/status", h.UpdateStatus)
	protected.PATCH("/admin/users/:id/roles", h.AssignRoles)
	protected.GET("/admin/roles", h.ListRoles)
}
