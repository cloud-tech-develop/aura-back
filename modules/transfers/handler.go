package transfers

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

// ─── HU-TRANS-001: Create Transfer Request ────────────────────────────────────

func (h *Handler) CreateTransfer(c *gin.Context) {
	var req CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1 // Default for testing
	}

	transfer, err := h.svc.CreateTransfer(c.Request.Context(), userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, transfer)
}

func (h *Handler) GetTransfer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	transfer, err := h.svc.GetTransferByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Traslado no encontrado")
		return
	}

	response.OK(c, transfer)
}

func (h *Handler) ListTransfers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	var originBranchID, destBranchID *int64
	if oid := c.Query("origin_branch_id"); oid != "" {
		if id, err := strconv.ParseInt(oid, 10, 64); err == nil {
			originBranchID = &id
		}
	}
	if did := c.Query("destination_branch_id"); did != "" {
		if id, err := strconv.ParseInt(did, 10, 64); err == nil {
			destBranchID = &id
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

	transfers, total, err := h.svc.ListTransfers(c.Request.Context(), originBranchID, destBranchID, status, startDate, endDate, page, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"data": transfers, "total": total})
}

// ─── HU-TRANS-002: Approve Transfer ──────────────────────────────────────────

func (h *Handler) ApproveTransfer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1
	}

	transfer, err := h.svc.ApproveTransfer(c.Request.Context(), id, userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, transfer)
}

// ─── HU-TRANS-003: Ship Transfer ──────────────────────────────────────────────

func (h *Handler) ShipTransfer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	var req ShipTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1
	}

	transfer, err := h.svc.ShipTransfer(c.Request.Context(), id, userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, transfer)
}

// ─── HU-TRANS-004: Receive Transfer ──────────────────────────────────────────

func (h *Handler) ReceiveTransfer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	var req ReceiveTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1
	}

	transfer, err := h.svc.ReceiveTransfer(c.Request.Context(), id, userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, transfer)
}

// ─── HU-TRANS-005: Cancel Transfer ──────────────────────────────────────────

func (h *Handler) CancelTransfer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	var req CancelTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1
	}

	transfer, err := h.svc.CancelTransfer(c.Request.Context(), id, userID, req.Reason)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, transfer)
}
