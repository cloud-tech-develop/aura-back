package cash

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

// ─── HU-CASH-001: Configure Cash Drawer ─────────────────────────────────────

func (h *Handler) ConfigureDrawer(c *gin.Context) {
	var req ConfigureDrawerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	drawer, err := h.svc.ConfigureDrawer(c.Request.Context(), req.BranchID, req.Name, req.MinFloat)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, drawer)
}

func (h *Handler) GetDrawerByBranch(c *gin.Context) {
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "branch_id inválido")
		return
	}

	drawer, err := h.svc.GetDrawerByBranch(c.Request.Context(), branchID)
	if err != nil {
		response.NotFound(c, "Caja no encontrada")
		return
	}

	response.OK(c, drawer)
}

// ─── HU-CASH-002: Open Cash Shift ───────────────────────────────────────────

func (h *Handler) OpenShift(c *gin.Context) {
	var req OpenShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	branchID := c.GetInt64("branch_id")

	if userID == 0 {
		response.Unauthorized(c, "user_id no encontrado")
		return
	}

	shift, err := h.svc.OpenShift(c.Request.Context(), userID, branchID, req.CashDrawerID, req.OpeningAmount, req.OpeningNotes)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, shift)
}

// ─── HU-CASH-003: Close Cash Shift ───────────────────────────────────────────

func (h *Handler) CloseShift(c *gin.Context) {
	shiftIDStr := c.Param("shiftID")
	shiftID, err := strconv.ParseInt(shiftIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "shift_id inválido")
		return
	}

	var req CloseShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")

	shift, err := h.svc.CloseShift(c.Request.Context(), shiftID, userID, req.ClosingAmount, req.ClosingNotes)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, shift)
}

// ─── HU-CASH-004: Record Cash Movement ───────────────────────────────────────

func (h *Handler) RecordMovement(c *gin.Context) {
	var req RecordMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Unauthorized(c, "user_id no encontrado")
		return
	}

	movement, err := h.svc.RecordMovement(c.Request.Context(), userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, movement)
}

// ─── HU-CASH-005: Perform Cash Reconciliation ────────────────────────────────

func (h *Handler) ReconcileShift(c *gin.Context) {
	shiftIDStr := c.Param("shiftID")
	shiftID, err := strconv.ParseInt(shiftIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "shift_id inválido")
		return
	}

	var req ReconcileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	shift, err := h.svc.ReconcileShift(c.Request.Context(), shiftID, req.ExpectedAmount, req.Notes)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, shift)
}

// ─── HU-CASH-006: View Shift Summary ──────────────────────────────────────────

func (h *Handler) GetShiftSummary(c *gin.Context) {
	shiftIDStr := c.Param("shiftID")
	shiftID, err := strconv.ParseInt(shiftIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "shift_id inválido")
		return
	}

	summary, err := h.svc.GetShiftSummary(c.Request.Context(), shiftID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, summary)
}

// ─── Additional Endpoints ────────────────────────────────────────────────────

func (h *Handler) GetActiveShift(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Unauthorized(c, "user_id no encontrado")
		return
	}

	shift, err := h.svc.GetActiveShift(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "No hay turno activo")
		return
	}

	response.OK(c, shift)
}

func (h *Handler) ListShifts(c *gin.Context) {
	branchID := c.GetInt64("branch_id")
	if branchID == 0 {
		response.BadRequest(c, "branch_id no encontrado")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

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

	shifts, total, err := h.svc.ListShifts(c.Request.Context(), branchID, startDate, endDate, status, page, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"data":  shifts,
		"total": total,
	})
}
