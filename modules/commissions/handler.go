package commissions

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

// ─── HU-COMM-001: Configure Commission Rules ─────────────────────────────────

func (h *Handler) CreateRule(c *gin.Context) {
	var req CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rule, err := h.svc.CreateRule(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, rule)
}

func (h *Handler) UpdateRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	var req UpdateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rule, err := h.svc.UpdateRule(c.Request.Context(), id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, rule)
}

func (h *Handler) DeleteRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	if err := h.svc.DeleteRule(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Regla eliminada"})
}

func (h *Handler) ListRules(c *gin.Context) {
	activeOnly := c.Query("active") == "true"

	rules, err := h.svc.ListRules(c.Request.Context(), activeOnly)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, rules)
}

// ─── HU-COMM-002: Calculate Commissions on Sale ──────────────────────────────

func (h *Handler) CalculateCommissions(c *gin.Context) {
	var req CalculateCommissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	commissions, err := h.svc.CalculateCommissions(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, commissions)
}

// ─── HU-COMM-003: View Commission History ─────────────────────────────────────

func (h *Handler) ListCommissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	var employeeID *int64
	if eid := c.Query("employee_id"); eid != "" {
		if id, err := strconv.ParseInt(eid, 10, 64); err == nil {
			employeeID = &id
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

	commissions, total, err := h.svc.ListCommissions(c.Request.Context(), employeeID, status, startDate, endDate, page, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"data": commissions, "total": total})
}

func (h *Handler) GetCommission(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	commission, err := h.svc.GetCommissionByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Comisión no encontrada")
		return
	}

	response.OK(c, commission)
}

// ─── HU-COMM-004: Settle Commissions ─────────────────────────────────────────

func (h *Handler) SettleCommissions(c *gin.Context) {
	var req SettleCommissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1 // Default for testing
	}

	if err := h.svc.SettleCommissions(c.Request.Context(), userID, req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Comisiones liquidadas"})
}

// ─── HU-COMM-005: Commission Reports ─────────────────────────────────────────

func (h *Handler) GetCommissionReport(c *gin.Context) {
	var filter CommissionReportFilter

	if eid := c.Query("employee_id"); eid != "" {
		if id, err := strconv.ParseInt(eid, 10, 64); err == nil {
			filter.EmployeeID = &id
		}
	}
	if bid := c.Query("branch_id"); bid != "" {
		if id, err := strconv.ParseInt(bid, 10, 64); err == nil {
			filter.BranchID = &id
		}
	}
	filter.Status = c.Query("status")
	if sd := c.Query("start_date"); sd != "" {
		if t, err := time.Parse("2006-01-02", sd); err == nil {
			filter.StartDate = &t
		}
	}
	if ed := c.Query("end_date"); ed != "" {
		if t, err := time.Parse("2006-01-02", ed); err == nil {
			filter.EndDate = &t
		}
	}

	summary, err := h.svc.GetCommissionReport(c.Request.Context(), filter)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, summary)
}

func (h *Handler) GetCommissionTotals(c *gin.Context) {
	var employeeID *int64
	if eid := c.Query("employee_id"); eid != "" {
		if id, err := strconv.ParseInt(eid, 10, 64); err == nil {
			employeeID = &id
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

	totals, err := h.svc.GetCommissionTotals(c.Request.Context(), employeeID, startDate, endDate)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, totals)
}
