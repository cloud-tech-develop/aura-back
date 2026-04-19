package presentations

import (
	"github.com/gin-gonic/gin"
)

func Register(protected gin.IRouter, h *Handler) {
	protected.GET("/presentations", h.List)
	protected.GET("/presentations/:id", h.GetByID)
	protected.POST("/presentations/page", h.Page)
	protected.PUT("/presentations/:id", h.Update)
	protected.DELETE("/presentations/:id", h.Delete)
}
