package sales

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

// CreateOrder - POST /sales-orders
func (h *Handler) CreateOrder(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	userID := c.GetInt64("user_id")

	if enterpriseID == 0 || userID == 0 {
		response.BadRequest(c, "enterprise_id or user_id not found")
		return
	}

	var req struct {
		OrderNumber  string  `json:"order_number"`
		CustomerID   *int64  `json:"customer_id"`
		BranchID     int64   `json:"branch_id" binding:"required"`
		Subtotal     float64 `json:"subtotal" binding:"required"`
		Discount     float64 `json:"discount"`
		TaxTotal     float64 `json:"tax_total" binding:"required"`
		Total        float64 `json:"total" binding:"required"`
		Notes        string  `json:"notes"`
		Items        []struct {
			ProductID       int64   `json:"product_id" binding:"required"`
			Quantity        int     `json:"quantity" binding:"required,min=1"`
			UnitPrice       float64 `json:"unit_price" binding:"required"`
			DiscountPercent float64 `json:"discount_percent"`
			DiscountAmount  float64 `json:"discount_amount"`
			TaxRate         float64 `json:"tax_rate" binding:"required"`
			TaxAmount       float64 `json:"tax_amount" binding:"required"`
			Total           float64 `json:"total" binding:"required"`
		} `json:"items" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	order := &SalesOrder{
		OrderNumber:  req.OrderNumber,
		CustomerID:   req.CustomerID,
		UserID:       userID,
		BranchID:     req.BranchID,
		EnterpriseID: enterpriseID,
		Subtotal:     req.Subtotal,
		Discount:     req.Discount,
		TaxTotal:     req.TaxTotal,
		Total:        req.Total,
		Status:       StatusPendingPayment,
		Notes:        req.Notes,
	}

	// Get repository from service to create order and items
	repo := NewRepository(c.MustGet("db").(*sql.DB))

	if err := repo.CreateOrder(c.Request.Context(), order); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Create order items
	for _, itemReq := range req.Items {
		item := &SalesOrderItem{
			SalesOrderID:    order.ID,
			ProductID:       itemReq.ProductID,
			Quantity:        itemReq.Quantity,
			UnitPrice:       itemReq.UnitPrice,
			DiscountPercent: itemReq.DiscountPercent,
			DiscountAmount:  itemReq.DiscountAmount,
			TaxRate:         itemReq.TaxRate,
			TaxAmount:       itemReq.TaxAmount,
			Total:           itemReq.Total,
		}
		if err := repo.CreateOrderItem(c.Request.Context(), item); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
	}

	response.Created(c, order)
}

// GetOrder - GET /sales-orders/:id
func (h *Handler) GetOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order ID")
		return
	}

	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "order not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, order)
}

// ListOrders - GET /sales-orders
func (h *Handler) ListOrders(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	var branchID, customerID *int64
	if bidStr := c.Query("branch_id"); bidStr != "" {
		if id, err := strconv.ParseInt(bidStr, 10, 64); err == nil {
			branchID = &id
		}
	}
	if cidStr := c.Query("customer_id"); cidStr != "" {
		if id, err := strconv.ParseInt(cidStr, 10, 64); err == nil {
			customerID = &id
		}
	}

	filters := OrderFilters{
		Status:     status,
		BranchID:   branchID,
		CustomerID: customerID,
		Page:       page,
		Limit:      limit,
	}

	orders, err := h.svc.GetOrders(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, orders)
}

// UpdateOrderStatus - PUT /sales-orders/:id/status
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateOrderStatus(c.Request.Context(), id, req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "order status updated"})
}

// CancelOrder - POST /sales-orders/:id/cancel
func (h *Handler) CancelOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order ID")
		return
	}

	if err := h.svc.CancelOrder(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "order cancelled"})
}

// CompleteOrder - POST /sales-orders/:id/complete
func (h *Handler) CompleteOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order ID")
		return
	}

	if err := h.svc.CompleteOrder(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "order completed"})
}
