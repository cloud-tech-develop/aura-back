package inventory

import (
	"database/sql"
	"strconv"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(database *db.DB) *Handler {
	q := database.Wrap(database.DB)
	return &Handler{svc: NewService(q)}
}

// GetInventory - GET /inventory/:productId/:branchId
func (h *Handler) GetInventory(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID de producto inválido")
		return
	}

	branchID, err := strconv.ParseInt(c.Param("branchId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID de sucursal inválido")
		return
	}

	inv, err := h.svc.GetInventory(c.Request.Context(), productID, branchID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Inventario no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, inv)
}

// ListInventory - GET /inventory
func (h *Handler) ListInventory(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var branchID, productID *int64
	if bid := c.Query("branch_id"); bid != "" {
		if id, err := strconv.ParseInt(bid, 10, 64); err == nil {
			branchID = &id
		}
	}
	if pid := c.Query("product_id"); pid != "" {
		if id, err := strconv.ParseInt(pid, 10, 64); err == nil {
			productID = &id
		}
	}

	lowStock := c.Query("low_stock") == "true"

	filters := InventoryFilters{
		Page:      page,
		Limit:     limit,
		BranchID:  branchID,
		ProductID: productID,
		LowStock:  lowStock,
	}

	list, err := h.svc.ListInventory(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

// GetLowStock - GET /inventory/low-stock
func (h *Handler) GetLowStock(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	list, err := h.svc.GetLowStock(c.Request.Context(), enterpriseID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

// GetInventoryByProduct - GET /inventory/product/:productId
func (h *Handler) GetInventoryByProduct(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID de producto inválido")
		return
	}

	list, err := h.svc.GetInventoryByProduct(c.Request.Context(), productID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

// UpdateStock - POST /inventory/movements
func (h *Handler) UpdateStock(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	userID := c.GetInt64("user_id")

	if enterpriseID == 0 || userID == 0 {
		response.BadRequest(c, "enterprise_id or user_id not found")
		return
	}

	var req struct {
		ProductID     int64   `json:"product_id" binding:"required"`
		BranchID      int64   `json:"branch_id" binding:"required"`
		Quantity      int     `json:"quantity" binding:"required"`
		Reason        string  `json:"reason" binding:"required"`
		ReferenceID   *int64  `json:"reference_id"`
		ReferenceType string  `json:"reference_type"`
		Notes         string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	inv, err := h.svc.UpdateStock(
		c.Request.Context(),
		req.ProductID,
		req.BranchID,
		req.Quantity,
		req.Reason,
		req.ReferenceType,
		req.ReferenceID,
		userID,
		req.Notes,
	)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, inv)
}

// GetMovement - GET /movements/:id
func (h *Handler) GetMovement(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	mov, err := h.svc.GetMovement(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Movimiento no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, mov)
}

// ListMovements - GET /movements
func (h *Handler) ListMovements(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var branchID, productID, inventoryID *int64
	if bid := c.Query("branch_id"); bid != "" {
		if id, err := strconv.ParseInt(bid, 10, 64); err == nil {
			branchID = &id
		}
	}
	if pid := c.Query("product_id"); pid != "" {
		if id, err := strconv.ParseInt(pid, 10, 64); err == nil {
			productID = &id
		}
	}
	if iid := c.Query("inventory_id"); iid != "" {
		if id, err := strconv.ParseInt(iid, 10, 64); err == nil {
			inventoryID = &id
		}
	}

	filters := MovementFilters{
		Page:           page,
		Limit:          limit,
		BranchID:       branchID,
		ProductID:      productID,
		InventoryID:    inventoryID,
		MovementType:   c.Query("movement_type"),
		MovementReason: c.Query("reason"),
	}

	list, err := h.svc.ListMovements(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

// GetProductKardex - GET /inventory/kardex/:productId/:branchId
func (h *Handler) GetProductKardex(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID de producto inválido")
		return
	}

	branchID, err := strconv.ParseInt(c.Param("branchId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID de sucursal inválido")
		return
	}

	kardex, err := h.svc.GetProductKardex(c.Request.Context(), productID, branchID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, kardex)
}

// ListReasons - GET /movement-reasons
func (h *Handler) ListReasons(c *gin.Context) {
	list, err := h.svc.ListReasons(c.Request.Context())
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}
