package shrinkage

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

// ─── HU-SHR-001: Register Shrinkage ──────────────────────────────────────────

func (h *Handler) RegisterShrinkage(c *gin.Context) {
	var req RegisterShrinkageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1 // Default for testing
	}

	shrinkage, err := h.svc.RegisterShrinkage(c.Request.Context(), userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, shrinkage)
}

func (h *Handler) GetShrinkage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	shrinkage, err := h.svc.GetShrinkageByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Merma no encontrada")
		return
	}

	response.OK(c, shrinkage)
}

func (h *Handler) ListShrinkages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	var branchID *int64
	if bid := c.Query("branch_id"); bid != "" {
		if id, err := strconv.ParseInt(bid, 10, 64); err == nil {
			branchID = &id
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

	shrinkages, total, err := h.svc.ListShrinkages(c.Request.Context(), branchID, status, startDate, endDate, page, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"data": shrinkages, "total": total})
}

// ─── HU-SHR-002: Configure Shrinkage Reasons ──────────────────────────────────

func (h *Handler) CreateReason(c *gin.Context) {
	var req CreateReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	reason, err := h.svc.CreateReason(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, reason)
}

func (h *Handler) ListReasons(c *gin.Context) {
	activeOnly := c.Query("active") == "true"

	reasons, err := h.svc.ListReasons(c.Request.Context(), activeOnly)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, reasons)
}

// ─── HU-SHR-003: Authorize High-Value Shrinkage ────────────────────────────────

func (h *Handler) AuthorizeShrinkage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	var req AuthorizeShrinkageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1
	}

	shrinkage, err := h.svc.AuthorizeShrinkage(c.Request.Context(), id, userID, req.Approved, req.Notes)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, shrinkage)
}

// ─── HU-SHR-004: View Shrinkage Report ─────────────────────────────────────────

func (h *Handler) GetShrinkageReport(c *gin.Context) {
	var branchID *int64
	if bid := c.Query("branch_id"); bid != "" {
		if id, err := strconv.ParseInt(bid, 10, 64); err == nil {
			branchID = &id
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

	report, err := h.svc.GetShrinkageReport(c.Request.Context(), branchID, startDate, endDate)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, report)
}

// ─── HU-SHR-005: Cancel Shrinkage ──────────────────────────────────────────────

func (h *Handler) CancelShrinkage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	var req CancelShrinkageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1
	}

	shrinkage, err := h.svc.CancelShrinkage(c.Request.Context(), id, userID, req.Reason)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, shrinkage)
}
