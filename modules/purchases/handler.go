package purchases

import (
	"strconv"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// ─── HU-PUR-001: Create Purchase Order ───────────────────────────────────────

func (h *Handler) CreatePurchaseOrder(c *gin.Context) {
	var req CreatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1 // Default for testing
	}

	po, err := h.svc.CreatePurchaseOrder(c.Request.Context(), userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, po)
}

func (h *Handler) GetPurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	po, err := h.svc.GetPurchaseOrderByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Orden de compra no encontrada")
		return
	}

	response.OK(c, po)
}

func (h *Handler) ListPurchaseOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	var supplierID *int64
	if sid := c.Query("supplier_id"); sid != "" {
		if id, err := strconv.ParseInt(sid, 10, 64); err == nil {
			supplierID = &id
		}
	}

	var startDate, endDate *time.Time
	if sd := c.Query("start_date"); sd != "" {
		if t, err := time.Parse("2006-01-02", sd); err == nil {
			startDate = &t
		}
	}
	if ed := c.Query("end_date"); ed != "" {
		if t, err := time.Parse("2006-01-02", ed); err == nil {
			endDate = &t
		}
	}

	orders, total, err := h.svc.ListPurchaseOrders(c.Request.Context(), supplierID, status, startDate, endDate, page, limit)

	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"data": orders, "total": total})
}

// ─── HU-PUR-002: Receive Goods ────────────────────────────────────────────────

func (h *Handler) ReceiveGoods(c *gin.Context) {
	var req ReceiveGoodsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1
	}

	purchase, err := h.svc.ReceiveGoods(c.Request.Context(), userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, purchase)
}

// ─── HU-PUR-003: Record Purchase Payment ──────────────────────────────────────

func (h *Handler) RecordPayment(c *gin.Context) {
	var req RecordPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1
	}

	payment, err := h.svc.RecordPayment(c.Request.Context(), userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, payment)
}

// ─── HU-PUR-004: Cancel Purchase ──────────────────────────────────────────────

func (h *Handler) CancelPurchase(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	var req CancelPurchaseRequest
	_ = c.ShouldBindJSON(&req)

	if err := h.svc.CancelPurchase(c.Request.Context(), id, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Compra cancelada"})
}

// ─── HU-PUR-005: View Purchase History ─────────────────────────────────────────

func (h *Handler) GetPurchase(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	purchase, err := h.svc.GetPurchaseByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Compra no encontrada")
		return
	}

	response.OK(c, purchase)
}

func (h *Handler) ListPurchases(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	var supplierID *int64
	if sid := c.Query("supplier_id"); sid != "" {
		if id, err := strconv.ParseInt(sid, 10, 64); err == nil {
			supplierID = &id
		}
	}

	var startDate, endDate *time.Time
	if sd := c.Query("start_date"); sd != "" {
		if t, err := time.Parse("2006-01-02", sd); err == nil {
			startDate = &t
		}
	}
	if ed := c.Query("end_date"); ed != "" {
		if t, err := time.Parse("2006-01-02", ed); err == nil {
			endDate = &t
		}
	}

	purchases, total, err := h.svc.GetPurchaseHistory(c.Request.Context(), supplierID, status, startDate, endDate, page, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"data": purchases, "total": total})
}

// ─── HU-PUR-006: Supplier Account Summary ──────────────────────────────────────

func (h *Handler) GetSupplierSummary(c *gin.Context) {
	supplierID, err := strconv.ParseInt(c.Param("supplierID"), 10, 64)
	if err != nil {
		response.BadRequest(c, "supplier_id inválido")
		return
	}

	summary, err := h.svc.GetSupplierSummary(c.Request.Context(), supplierID)
	if err != nil {
		response.NotFound(c, "Resumen de proveedor no encontrado")
		return
	}

	response.OK(c, summary)
}
