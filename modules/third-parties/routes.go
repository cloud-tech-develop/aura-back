package thirdparties

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Third Party routes
	protected.POST("/third-parties", h.CreateThirdParty)
	protected.GET("/third-parties", h.ListThirdParties)
	protected.GET("/third-parties/:id", h.GetThirdParty)
	protected.GET("/third-parties/document/:documentNumber", h.GetThirdPartyByDocument)
	protected.PUT("/third-parties/:id", h.UpdateThirdParty)
	protected.DELETE("/third-parties/:id", h.DeleteThirdParty)
}
