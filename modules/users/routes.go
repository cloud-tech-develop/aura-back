package users

import "github.com/gin-gonic/gin"

// Register mounts all user routes onto the given router group.
func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Protected routes (require auth)
	protected.GET("/users", h.List)
	protected.GET("/users/:id", h.GetByID)
	protected.POST("/users", h.Create)
	protected.PUT("/users/:id", h.Update)
	protected.PATCH("/users/:id/status", h.UpdateStatus)
	protected.PATCH("/users/:id/roles", h.AssignRoles)
	protected.GET("/roles", h.ListRoles)
}
