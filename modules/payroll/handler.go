package payroll

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

// ─── HU-PAY-001: Register Employee ───────────────────────────────────────────

func (h *Handler) CreateEmployee(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	employee, err := h.svc.CreateEmployee(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, employee)
}

func (h *Handler) GetEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	employee, err := h.svc.GetEmployeeByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Empleado no encontrado")
		return
	}

	response.OK(c, employee)
}

func (h *Handler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	var req UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	employee, err := h.svc.UpdateEmployee(c.Request.Context(), id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, employee)
}

func (h *Handler) TerminateEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	var req TerminateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	employee, err := h.svc.TerminateEmployee(c.Request.Context(), id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, employee)
}

func (h *Handler) ListEmployees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	department := c.Query("department")
	status := c.Query("status")

	employees, total, err := h.svc.ListEmployees(c.Request.Context(), department, status, page, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"data": employees, "total": total})
}

// ─── HU-PAY-007: Set Employee Salary ──────────────────────────────────────────

func (h *Handler) SetSalary(c *gin.Context) {
	employeeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	var req SetSalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	salary, err := h.svc.SetSalary(c.Request.Context(), employeeID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, salary)
}

// ─── HU-PAY-011: Create Payroll Period ────────────────────────────────────────

func (h *Handler) CreatePeriod(c *gin.Context) {
	var req CreatePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	period, err := h.svc.CreatePeriod(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, period)
}

func (h *Handler) GetPeriod(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}

	period, err := h.svc.GetPeriodByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Período no encontrado")
		return
	}

	response.OK(c, period)
}

func (h *Handler) ListPeriods(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	periods, total, err := h.svc.ListPeriods(c.Request.Context(), status, page, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"data": periods, "total": total})
}

// ─── HU-PAY-014: Calculate Payroll ────────────────────────────────────────────

func (h *Handler) CalculatePayroll(c *gin.Context) {
	periodID, err := strconv.ParseInt(c.Param("periodId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "periodId inválido")
		return
	}

	if err := h.svc.CalculatePayroll(c.Request.Context(), periodID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Nómina calculada exitosamente"})
}

func (h *Handler) ListPayrolls(c *gin.Context) {
	periodID, err := strconv.ParseInt(c.Param("periodId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "periodId inválido")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	payrolls, total, err := h.svc.ListPayrolls(c.Request.Context(), periodID, status, page, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"data": payrolls, "total": total})
}

// ─── HU-PAY-016: Approve Payroll ──────────────────────────────────────────────

func (h *Handler) ApprovePayroll(c *gin.Context) {
	periodID, err := strconv.ParseInt(c.Param("periodId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "periodId inválido")
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1
	}

	if err := h.svc.ApprovePayroll(c.Request.Context(), periodID, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Nómina aprobada exitosamente"})
}

// ─── HU-PAY-012: Close Payroll Period ─────────────────────────────────────────

func (h *Handler) ClosePeriod(c *gin.Context) {
	periodID, err := strconv.ParseInt(c.Param("periodId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "periodId inválido")
		return
	}

	if err := h.svc.ClosePeriod(c.Request.Context(), periodID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Período cerrado exitosamente"})
}

// ─── HU-PAY-018: Register Overtime ────────────────────────────────────────────

func (h *Handler) RegisterOvertime(c *gin.Context) {
	var req RegisterOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	overtime, err := h.svc.RegisterOvertime(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, overtime)
}

// ─── HU-PAY-019: Register Bonus ───────────────────────────────────────────────

func (h *Handler) RegisterBonus(c *gin.Context) {
	var req RegisterBonusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	bonus, err := h.svc.RegisterBonus(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, bonus)
}

// ─── HU-PAY-010: Manage Employee Loans ─────────────────────────────────────────

func (h *Handler) CreateLoan(c *gin.Context) {
	var req CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	loan, err := h.svc.CreateLoan(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, loan)
}

// ─── HU-PAY-021: Process Payment ──────────────────────────────────────────────

func (h *Handler) ProcessPayment(c *gin.Context) {
	var req PayrollPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ProcessPayment(c.Request.Context(), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Pagos procesados exitosamente"})
}

// ─── HU-PAY-026: Payroll Summary Report ───────────────────────────────────────

func (h *Handler) GetPayrollSummary(c *gin.Context) {
	periodID, err := strconv.ParseInt(c.Param("periodId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "periodId inválido")
		return
	}

	summary, err := h.svc.GetPayrollSummary(c.Request.Context(), periodID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, summary)
}

// ─── Types ───────────────────────────────────────────────────────────────────

func (h *Handler) ListDeductionTypes(c *gin.Context) {
	activeOnly := c.Query("active") == "true"

	types, err := h.svc.ListDeductionTypes(c.Request.Context(), activeOnly)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, types)
}

func (h *Handler) ListAdditionTypes(c *gin.Context) {
	activeOnly := c.Query("active") == "true"

	types, err := h.svc.ListAdditionTypes(c.Request.Context(), activeOnly)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, types)
}

// ─── HU-PAY-013: Reopen Payroll Period ──────────────────────────────────────

func (h *Handler) ReopenPeriod(c *gin.Context) {
	periodID, err := strconv.ParseInt(c.Param("periodId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "periodId inválido")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ReopenPeriod(c.Request.Context(), periodID, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Período reabierto exitosamente"})
}

// ─── HU-PAY-015: Preview Payroll ────────────────────────────────────────────

func (h *Handler) PreviewPayroll(c *gin.Context) {
	periodID, err := strconv.ParseInt(c.Param("periodId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "periodId inválido")
		return
	}

	summary, err := h.svc.PreviewPayroll(c.Request.Context(), periodID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, summary)
}

// ─── HU-PAY-017: Reject Payroll ──────────────────────────────────────────────

func (h *Handler) RejectPayroll(c *gin.Context) {
	periodID, err := strconv.ParseInt(c.Param("periodId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "periodId inválido")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.RejectPayroll(c.Request.Context(), periodID, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Nómina rechazada"})
}

// ─── HU-PAY-020: Import Commissions ─────────────────────────────────────────

func (h *Handler) ImportCommissions(c *gin.Context) {
	periodID, err := strconv.ParseInt(c.Param("periodId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "periodId inválido")
		return
	}

	if err := h.svc.ImportCommissions(c.Request.Context(), periodID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Comisiones importadas"})
}

// ─── HU-PAY-024: View Payment History ────────────────────────────────────────

func (h *Handler) ListPayments(c *gin.Context) {
	var periodID, employeeID *int64

	if pid := c.Query("period_id"); pid != "" {
		if id, err := strconv.ParseInt(pid, 10, 64); err == nil {
			periodID = &id
		}
	}
	if eid := c.Query("employee_id"); eid != "" {
		if id, err := strconv.ParseInt(eid, 10, 64); err == nil {
			employeeID = &id
		}
	}

	payments, err := h.svc.ListPayments(c.Request.Context(), periodID, employeeID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, payments)
}

// ─── HU-PAY-025: Generate Payslip ────────────────────────────────────────────

func (h *Handler) GetPayslip(c *gin.Context) {
	payrollID, err := strconv.ParseInt(c.Param("payrollId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "payrollId inválido")
		return
	}

	payroll, err := h.svc.GetPayrollByID(c.Request.Context(), payrollID)
	if err != nil {
		response.NotFound(c, "Nómina no encontrada")
		return
	}

	details, _ := h.svc.GetPayrollDetails(c.Request.Context(), payrollID)

	response.OK(c, gin.H{
		"payroll": payroll,
		"details": details,
	})
}

// ─── HU-PAY-027: Tax Report ──────────────────────────────────────────────────

func (h *Handler) GetTaxReport(c *gin.Context) {
	periodID, err := strconv.ParseInt(c.Param("periodId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "periodId inválido")
		return
	}

	report, err := h.svc.GetTaxReport(c.Request.Context(), periodID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, report)
}

// ─── HU-PAY-028: Employee Earnings Report ────────────────────────────────────

func (h *Handler) GetEmployeeEarnings(c *gin.Context) {
	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "employeeId inválido")
		return
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

	earnings, err := h.svc.GetEmployeeEarnings(c.Request.Context(), employeeID, startDate, endDate)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, earnings)
}

// ─── HU-PAY-029: Deduction Report ────────────────────────────────────────────

func (h *Handler) GetDeductionReport(c *gin.Context) {
	periodID, err := strconv.ParseInt(c.Param("periodId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "periodId inválido")
		return
	}

	report, err := h.svc.GetDeductionReport(c.Request.Context(), periodID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, report)
}

// ─── HU-PAY-030: Configure Leave Types ──────────────────────────────────────

func (h *Handler) CreateLeaveType(c *gin.Context) {
	var req CreateLeaveTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	leaveType, err := h.svc.CreateLeaveType(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, leaveType)
}

func (h *Handler) ListLeaveTypes(c *gin.Context) {
	activeOnly := c.Query("active") == "true"

	types, err := h.svc.ListLeaveTypes(c.Request.Context(), activeOnly)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, types)
}

// ─── HU-PAY-031: Register Leave Request ──────────────────────────────────────

func (h *Handler) CreateLeaveRequest(c *gin.Context) {
	var req CreateLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	leave, err := h.svc.CreateLeaveRequest(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, leave)
}

// ─── HU-PAY-032: Approve Leave ──────────────────────────────────────────────

func (h *Handler) ApproveLeave(c *gin.Context) {
	leaveID, err := strconv.ParseInt(c.Param("leaveId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "leaveId inválido")
		return
	}

	var req struct {
		Approved bool `json:"approved"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		userID = 1
	}

	leave, err := h.svc.ApproveLeave(c.Request.Context(), leaveID, userID, req.Approved)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, leave)
}

// ─── HU-PAY-034: View Leave Balance ──────────────────────────────────────────

func (h *Handler) GetLeaveBalance(c *gin.Context) {
	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "employeeId inválido")
		return
	}

	year := time.Now().Year()
	if y := c.Query("year"); y != "" {
		if parsedYear, err := strconv.Atoi(y); err == nil {
			year = parsedYear
		}
	}

	balance, err := h.svc.GetLeaveBalance(c.Request.Context(), employeeID, year)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, balance)
}
