package cart

import (
	"database/sql"
	"strconv"

	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{svc: NewService(db)}
}

// CreateCart - POST /carts
func (h *Handler) CreateCart(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	userID := c.GetInt64("user_id")

	if enterpriseID == 0 || userID == 0 {
		response.BadRequest(c, "enterprise_id or user_id not found")
		return
	}

	var req struct {
		BranchID   int64   `json:"branch_id" binding:"required"`
		CustomerID *int64  `json:"customer_id"`
		CartCode   string  `json:"cart_code"`
		CartType   string  `json:"cart_type"`
		Notes      string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cart := &Cart{
		CartCode:     req.CartCode,
		CartType:     req.CartType,
		CustomerID:   req.CustomerID,
		UserID:       userID,
		BranchID:     req.BranchID,
		EnterpriseID: enterpriseID,
		Notes:        req.Notes,
	}

	if err := h.svc.CreateCart(c.Request.Context(), cart); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, cart)
}

// ListCarts - GET /carts
func (h *Handler) ListCarts(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var branchID *int64
	if bid := c.Query("branch_id"); bid != "" {
		if id, err := strconv.ParseInt(bid, 10, 64); err == nil {
			branchID = &id
		}
	}

	filters := CartFilters{
		Page:      page,
		Limit:     limit,
		CartType:  c.Query("type"),
		Status:    c.Query("status"),
		BranchID:  branchID,
	}

	list, err := h.svc.ListCarts(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

// GetCart - GET /carts/:id
func (h *Handler) GetCart(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart ID")
		return
	}

	cart, err := h.svc.GetCart(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "cart not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, cart)
}

// GetCartByCode - GET /carts/code/:code
func (h *Handler) GetCartByCode(c *gin.Context) {
	code := c.Param("code")
	enterpriseID := c.GetInt64("enterprise_id")

	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	cart, err := h.svc.GetCartByCode(c.Request.Context(), code, enterpriseID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "cart not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, cart)
}

// AddItem - POST /carts/:id/items
func (h *Handler) AddItem(c *gin.Context) {
	idParam := c.Param("id")
	cartID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart ID")
		return
	}

	var req struct {
		ProductID        int64   `json:"product_id" binding:"required"`
		ProductVariantID *int64  `json:"product_variant_id"`
		Quantity         int     `json:"quantity" binding:"required,min=1"`
		UnitPrice        float64 `json:"unit_price" binding:"required"`
		DiscountType     string  `json:"discount_type"`
		DiscountValue    float64 `json:"discount_value"`
		TaxRate          float64 `json:"tax_rate" binding:"required"`
		Notes            string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item := &CartItem{
		ProductID:       req.ProductID,
		ProductVariantID: req.ProductVariantID,
		Quantity:        req.Quantity,
		UnitPrice:       req.UnitPrice,
		DiscountType:    req.DiscountType,
		DiscountValue:   req.DiscountValue,
		TaxRate:         req.TaxRate,
		Notes:           req.Notes,
	}

	if err := h.svc.AddItem(c.Request.Context(), cartID, item); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, item)
}

// UpdateItem - PUT /carts/:id/items/:itemId
func (h *Handler) UpdateItem(c *gin.Context) {
	cartIDParam := c.Param("id")
	cartID, err := strconv.ParseInt(cartIDParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart ID")
		return
	}

	itemIDParam := c.Param("itemId")
	itemID, err := strconv.ParseInt(itemIDParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid item ID")
		return
	}

	var req struct {
		Quantity      int     `json:"quantity" binding:"required,min=1"`
		UnitPrice     float64 `json:"unit_price"`
		DiscountType  string  `json:"discount_type"`
		DiscountValue float64 `json:"discount_value"`
		TaxRate       float64 `json:"tax_rate"`
		Notes         string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item := &CartItem{
		Quantity:      req.Quantity,
		UnitPrice:     req.UnitPrice,
		DiscountType:  req.DiscountType,
		DiscountValue: req.DiscountValue,
		TaxRate:       req.TaxRate,
		Notes:         req.Notes,
	}

	if err := h.svc.UpdateItem(c.Request.Context(), cartID, itemID, item); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "item updated"})
}

// RemoveItem - DELETE /carts/:id/items/:itemId
func (h *Handler) RemoveItem(c *gin.Context) {
	cartIDParam := c.Param("id")
	cartID, err := strconv.ParseInt(cartIDParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart ID")
		return
	}

	itemIDParam := c.Param("itemId")
	itemID, err := strconv.ParseInt(itemIDParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid item ID")
		return
	}

	if err := h.svc.RemoveItem(c.Request.Context(), cartID, itemID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.Status(204)
}

// ConvertToSale - POST /carts/:id/convert
func (h *Handler) ConvertToSale(c *gin.Context) {
	idParam := c.Param("id")
	cartID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart ID")
		return
	}

	cart, err := h.svc.ConvertToSale(c.Request.Context(), cartID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"message": "cart converted to sale",
		"cart":    cart,
	})
}

// ConvertToQuotation - POST /carts/:id/quotation
func (h *Handler) ConvertToQuotation(c *gin.Context) {
	idParam := c.Param("id")
	cartID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart ID")
		return
	}

	var req struct {
		ValidDays int `json:"valid_days"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.ValidDays = 30 // default 30 days
	}

	cart, err := h.svc.ConvertToQuotation(c.Request.Context(), cartID, req.ValidDays)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"message": "cart converted to quotation",
		"cart":    cart,
	})
}

// SetCustomer - PUT /carts/:id/customer
func (h *Handler) SetCustomer(c *gin.Context) {
	idParam := c.Param("id")
	cartID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart ID")
		return
	}

	var req struct {
		CustomerID *int64 `json:"customer_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.SetCustomer(c.Request.Context(), cartID, req.CustomerID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "customer set"})
}

// ApplyDiscount - POST /carts/:id/discount
func (h *Handler) ApplyDiscount(c *gin.Context) {
	idParam := c.Param("id")
	cartID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart ID")
		return
	}

	var req struct {
		DiscountType  string  `json:"discount_type" binding:"required"`
		DiscountValue float64 `json:"discount_value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ApplyDiscount(c.Request.Context(), cartID, req.DiscountType, req.DiscountValue); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "discount applied"})
}

// ApplyItemDiscount - POST /carts/:id/items/:itemId/discount
func (h *Handler) ApplyItemDiscount(c *gin.Context) {
	cartIDParam := c.Param("id")
	cartID, err := strconv.ParseInt(cartIDParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart ID")
		return
	}

	itemIDParam := c.Param("itemId")
	itemID, err := strconv.ParseInt(itemIDParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid item ID")
		return
	}

	var req struct {
		DiscountType  string  `json:"discount_type" binding:"required"`
		DiscountValue float64 `json:"discount_value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ApplyItemDiscount(c.Request.Context(), cartID, itemID, req.DiscountType, req.DiscountValue); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "item discount applied"})
}

// DeleteCart - DELETE /carts/:id
func (h *Handler) DeleteCart(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart ID")
		return
	}

	if err := h.svc.DeleteCart(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.Status(204)
}
