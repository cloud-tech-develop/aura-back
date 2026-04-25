package users

import "github.com/gin-gonic/gin"

// Register mounts all user routes onto the given router group.
func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Public — no auth required (for offline sync)
	public.GET("/users", h.ListByEnterpriseID)
	public.GET("/users-sync", h.ListByEnterpriseIDForSync) // For offline sync - includes password hashes
	public.GET("/user-roles", h.ListUserRolesByEnterpriseID)

	// Protected routes (require auth)
	protected.GET("/admin/users", h.List)
	protected.GET("/admin/users/:id", h.GetByID)
	protected.POST("/admin/users", h.Create)
	protected.PUT("/admin/users/:id", h.Update)
	protected.PATCH("/admin/users/:id/status", h.UpdateStatus)
	protected.PATCH("/admin/users/:id/roles", h.AssignRoles)
	protected.GET("/admin/roles", h.ListRoles)
}
