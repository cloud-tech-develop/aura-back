package compositions

import "github.com/gin-gonic/gin"

func Register(protected gin.IRouter, h *Handler) {
	protected.GET("/catalog/compositions/all", h.ListAll)
	protected.GET("/catalog/compositions", h.List)
	protected.POST("/catalog/compositions", h.Create)
	protected.GET("/catalog/compositions/:id", h.GetByID)
	protected.POST("/catalog/compositions/page", h.Page)
	protected.PUT("/catalog/compositions/:id", h.Update)
	protected.DELETE("/catalog/compositions/:id", h.Delete)
}
