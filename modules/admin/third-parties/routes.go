package thirdparties

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Third Party routes
	protected.POST("/admin/third-parties", h.CreateThirdParty)
	protected.GET("/admin/third-parties", h.ListThirdParties)
	protected.GET("/admin/third-parties/:id", h.GetThirdParty)
	protected.GET("/admin/third-parties/document/:documentNumber", h.GetThirdPartyByDocument)
	protected.PUT("/admin/third-parties/:id", h.UpdateThirdParty)
	protected.DELETE("/admin/third-parties/:id", h.DeleteThirdParty)
}
