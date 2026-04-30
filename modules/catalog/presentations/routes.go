package presentations

import (
	"github.com/gin-gonic/gin"
)

func Register(protected gin.IRouter, h *Handler) {
	protected.GET("/catalog/presentations", h.List)
	protected.GET("/catalog/presentations/:id", h.GetByID)
	protected.POST("/catalog/presentations/page", h.Page)
	protected.PUT("/catalog/presentations/:id", h.Update)
	protected.DELETE("/catalog/presentations/:id", h.Delete)
}
