package payroll

import "github.com/gin-gonic/gin"

func Register(public gin.IRouter, protected gin.IRouter, h *Handler) {
	// Employees
	protected.POST("/payroll/employees", h.CreateEmployee)
	protected.GET("/payroll/employees/:id", h.GetEmployee)
	protected.PUT("/payroll/employees/:id", h.UpdateEmployee)
	protected.POST("/payroll/employees/:id/terminate", h.TerminateEmployee)
	protected.GET("/payroll/employees", h.ListEmployees)
	protected.POST("/payroll/employees/:id/salary", h.SetSalary)

	// Payroll Periods
	protected.POST("/payroll/periods", h.CreatePeriod)
	protected.GET("/payroll/periods/:id", h.GetPeriod)
	protected.GET("/payroll/periods", h.ListPeriods)
	protected.POST("/payroll/periods/:periodId/calculate", h.CalculatePayroll)
	protected.POST("/payroll/periods/:periodId/approve", h.ApprovePayroll)
	protected.POST("/payroll/periods/:periodId/close", h.ClosePeriod)
	protected.POST("/payroll/periods/:periodId/reopen", h.ReopenPeriod)
	protected.POST("/payroll/periods/:periodId/reject", h.RejectPayroll)

	// Payroll Records
	protected.GET("/payroll/periods/:periodId/payrolls", h.ListPayrolls)
	protected.GET("/payroll/periods/:periodId/preview", h.PreviewPayroll)
	protected.GET("/payroll/periods/:periodId/summary", h.GetPayrollSummary)
	protected.GET("/payroll/payslips/:payrollId", h.GetPayslip)

	// Overtime
	protected.POST("/payroll/overtime", h.RegisterOvertime)

	// Bonuses
	protected.POST("/payroll/bonuses", h.RegisterBonus)

	// Loans
	protected.POST("/payroll/loans", h.CreateLoan)

	// Payments
	protected.POST("/payroll/payments", h.ProcessPayment)
	protected.GET("/payroll/payments", h.ListPayments)

	// Commissions Import
	protected.POST("/payroll/periods/:periodId/import-commissions", h.ImportCommissions)

	// Reports
	protected.GET("/payroll/periods/:periodId/reports/tax", h.GetTaxReport)
	protected.GET("/payroll/periods/:periodId/reports/deductions", h.GetDeductionReport)
	protected.GET("/payroll/employees/:employeeId/earnings", h.GetEmployeeEarnings)

	// Types
	protected.GET("/payroll/deduction-types", h.ListDeductionTypes)
	protected.GET("/payroll/addition-types", h.ListAdditionTypes)

	// Leave Management
	protected.POST("/payroll/leave-types", h.CreateLeaveType)
	protected.GET("/payroll/leave-types", h.ListLeaveTypes)
	protected.POST("/payroll/leaves", h.CreateLeaveRequest)
	protected.POST("/payroll/leaves/:leaveId/approve", h.ApproveLeave)
	protected.GET("/payroll/employees/:employeeId/leave-balance", h.GetLeaveBalance)
}
