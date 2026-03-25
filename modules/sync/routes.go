package sync

import (
	"github.com/gin-gonic/gin"
)

func Register(public, protected gin.IRouter, h *Handler) {
	sync := protected.Group("/sync")
	{
		sync.GET("/pull", h.Pull)
		sync.POST("/push", h.Push)
	}
}
