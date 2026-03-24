package cart

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Cart routes
	protected.POST("/carts", h.CreateCart)
	protected.GET("/carts", h.ListCarts)
	protected.GET("/carts/:id", h.GetCart)
	protected.GET("/carts/code/:code", h.GetCartByCode)
	protected.DELETE("/carts/:id", h.DeleteCart)

	// Cart item routes
	protected.POST("/carts/:id/items", h.AddItem)
	protected.PUT("/carts/:id/items/:itemId", h.UpdateItem)
	protected.DELETE("/carts/:id/items/:itemId", h.RemoveItem)
	protected.POST("/carts/:id/items/:itemId/discount", h.ApplyItemDiscount)

	// Conversion routes
	protected.POST("/carts/:id/convert", h.ConvertToSale)
	protected.POST("/carts/:id/quotation", h.ConvertToQuotation)

	// Cart operations
	protected.PUT("/carts/:id/customer", h.SetCustomer)
	protected.POST("/carts/:id/discount", h.ApplyDiscount)
}
