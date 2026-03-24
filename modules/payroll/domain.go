package payroll

import (
	"context"
	"time"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
	EventEmployeeCreated   = "payroll.employee.created"
	EventEmployeeUpdated   = "payroll.employee.updated"
	EventSalarySet         = "payroll.salary.set"
	EventPayrollCalculated = "payroll.calculated"
	EventPayrollApproved   = "payroll.approved"
	EventPayrollPaid       = "payroll.paid"
)

const (
	StatusActive     = "ACTIVE"
	StatusInactive   = "INACTIVE"
	StatusTerminated = "TERMINATED"
)

const (
	EmploymentFullTime   = "FULL_TIME"
	EmploymentPartTime   = "PART_TIME"
	EmploymentContractor = "CONTRACTOR"
	EmploymentTemporary  = "TEMPORARY"
)

const (
	PeriodStatusOpen       = "OPEN"
	PeriodStatusProcessing = "PROCESSING"
	PeriodStatusApproved   = "APPROVED"
	PeriodStatusPaid       = "PAID"
	PeriodStatusClosed     = "CLOSED"
)

const (
	PayrollStatusDraft      = "DRAFT"
	PayrollStatusCalculated = "CALCULATED"
	PayrollStatusApproved   = "APPROVED"
	PayrollStatusPaid       = "PAID"
	PayrollStatusCancelled  = "CANCELLED"
)

// ─── Entities ─────────────────────────────────────────────────────────────────

// Employee represents an employee
type Employee struct {
	ID                int64      `json:"id"`
	EmployeeCode      string     `json:"employee_code"`
	ThirdPartyID      int64      `json:"third_party_id"`
	UserID            *int64     `json:"user_id,omitempty"`
	EmploymentType    string     `json:"employment_type"`
	Position          string     `json:"position"`
	Department        string     `json:"department"`
	HireDate          time.Time  `json:"hire_date"`
	TerminationDate   *time.Time `json:"termination_date,omitempty"`
	TerminationReason string     `json:"termination_reason"`
	BankName          string     `json:"bank_name"`
	BankAccount       string     `json:"bank_account"`
	BankCLABE         string     `json:"bank_clabe"`
	CURP              string     `json:"curp"`
	RFC               string     `json:"rfc"`
	NSS               string     `json:"nss"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}

// Salary represents an employee's salary record
type Salary struct {
	ID            int64      `json:"id"`
	EmployeeID    int64      `json:"employee_id"`
	BaseSalary    float64    `json:"base_salary"`
	DailySalary   float64    `json:"daily_salary"`
	SalaryType    string     `json:"salary_type"`
	EffectiveDate time.Time  `json:"effective_date"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	IsCurrent     bool       `json:"is_current"`
	Notes         string     `json:"notes"`
	CreatedAt     time.Time  `json:"created_at"`
}

// PayrollPeriod represents a payroll period
type PayrollPeriod struct {
	ID              int64      `json:"id"`
	PeriodType      string     `json:"period_type"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	PaymentDate     time.Time  `json:"payment_date"`
	Status          string     `json:"status"`
	TotalGross      float64    `json:"total_gross"`
	TotalDeductions float64    `json:"total_deductions"`
	TotalNet        float64    `json:"total_net"`
	EmployeeCount   int        `json:"employee_count"`
	ApprovedBy      *int64     `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

// Payroll represents a payroll record for an employee
type Payroll struct {
	ID              int64      `json:"id"`
	PayrollPeriodID int64      `json:"payroll_period_id"`
	EmployeeID      int64      `json:"employee_id"`
	SalaryID        *int64     `json:"salary_id,omitempty"`
	BaseAmount      float64    `json:"base_amount"`
	WorkedDays      float64    `json:"worked_days"`
	GrossAmount     float64    `json:"gross_amount"`
	TotalAdditions  float64    `json:"total_additions"`
	TotalDeductions float64    `json:"total_deductions"`
	NetAmount       float64    `json:"net_amount"`
	Status          string     `json:"status"`
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

// PayrollDetail represents a line item in payroll
type PayrollDetail struct {
	ID            int64     `json:"id"`
	PayrollID     int64     `json:"payroll_id"`
	ConceptType   string    `json:"concept_type"` // ADDITION or DEDUCTION
	ConceptCode   string    `json:"concept_code"`
	ConceptName   string    `json:"concept_name"`
	Quantity      float64   `json:"quantity"`
	Rate          float64   `json:"rate"`
	Amount        float64   `json:"amount"`
	ReferenceID   *int64    `json:"reference_id,omitempty"`
	ReferenceType *string   `json:"reference_type,omitempty"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

// DeductionType represents a deduction configuration
type DeductionType struct {
	ID              int64     `json:"id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	IsTax           bool      `json:"is_tax"`
	IsMandatory     bool      `json:"is_mandatory"`
	CalculationType string    `json:"calculation_type"`
	DefaultValue    *float64  `json:"default_value,omitempty"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

// AdditionType represents an addition configuration
type AdditionType struct {
	ID              int64     `json:"id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	IsTaxable       bool      `json:"is_taxable"`
	CalculationType string    `json:"calculation_type"`
	DefaultValue    *float64  `json:"default_value,omitempty"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

// EmployeeLoan represents a loan or advance
type EmployeeLoan struct {
	ID                int64      `json:"id"`
	EmployeeID        int64      `json:"employee_id"`
	LoanType          string     `json:"loan_type"`
	TotalAmount       float64    `json:"total_amount"`
	RemainingAmount   float64    `json:"remaining_amount"`
	InstallmentAmount float64    `json:"installment_amount"`
	InstallmentsTotal int        `json:"installments_total"`
	InstallmentsPaid  int        `json:"installments_paid"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           *time.Time `json:"end_date,omitempty"`
	Status            string     `json:"status"`
	Notes             string     `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}

// Overtime represents overtime record
type Overtime struct {
	ID          int64     `json:"id"`
	EmployeeID  int64     `json:"employee_id"`
	WorkDate    time.Time `json:"work_date"`
	Hours       float64   `json:"hours"`
	RateType    string    `json:"rate_type"`
	HourlyRate  float64   `json:"hourly_rate"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	ApprovedBy  *int64    `json:"approved_by,omitempty"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}

// Bonus represents a bonus record
type Bonus struct {
	ID         int64     `json:"id"`
	EmployeeID int64     `json:"employee_id"`
	BonusType  string    `json:"bonus_type"`
	Amount     float64   `json:"amount"`
	BonusDate  time.Time `json:"bonus_date"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
}

// PayrollPayment represents a payment record
type PayrollPayment struct {
	ID              int64     `json:"id"`
	PayrollID       int64     `json:"payroll_id"`
	EmployeeID      int64     `json:"employee_id"`
	PaymentMethod   string    `json:"payment_method"`
	Amount          float64   `json:"amount"`
	PaymentDate     time.Time `json:"payment_date"`
	ReferenceNumber *string   `json:"reference_number,omitempty"`
	BankReference   *string   `json:"bank_reference,omitempty"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
}

// PayrollSummary for reporting
type PayrollSummary struct {
	EmployeeID      int64   `json:"employee_id"`
	EmployeeCode    string  `json:"employee_code"`
	EmployeeName    string  `json:"employee_name"`
	GrossAmount     float64 `json:"gross_amount"`
	TotalAdditions  float64 `json:"total_additions"`
	TotalDeductions float64 `json:"total_deductions"`
	NetAmount       float64 `json:"net_amount"`
	Status          string  `json:"status"`
}

// TaxReportItem for tax reporting
type TaxReportItem struct {
	EmployeeID   int64   `json:"employee_id"`
	EmployeeCode string  `json:"employee_code"`
	EmployeeName string  `json:"employee_name"`
	RFC          string  `json:"rfc"`
	GrossAmount  float64 `json:"gross_amount"`
	ISRAmount    float64 `json:"isr_amount"`
	IMSSAmount   float64 `json:"imss_amount"`
	TotalTax     float64 `json:"total_tax"`
}

// DeductionReportItem for deduction reporting
type DeductionReportItem struct {
	DeductionCode string  `json:"deduction_code"`
	DeductionName string  `json:"deduction_name"`
	TotalAmount   float64 `json:"total_amount"`
	EmployeeCount int64   `json:"employee_count"`
}

// LeaveType represents a leave configuration
type LeaveType struct {
	ID               int64     `json:"id"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	IsPaid           bool      `json:"is_paid"`
	DaysAllowed      int       `json:"days_allowed"`
	RequiresApproval bool      `json:"requires_approval"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
}

// EmployeeLeave represents a leave request
type EmployeeLeave struct {
	ID          int64      `json:"id"`
	EmployeeID  int64      `json:"employee_id"`
	LeaveTypeID int64      `json:"leave_type_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Days        float64    `json:"days"`
	Status      string     `json:"status"`
	Reason      string     `json:"reason"`
	ApprovedBy  *int64     `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// LeaveBalance for employee leave balance
type LeaveBalance struct {
	LeaveTypeID   int64   `json:"leave_type_id"`
	LeaveTypeName string  `json:"leave_type_name"`
	DaysAllowed   int     `json:"days_allowed"`
	DaysUsed      float64 `json:"days_used"`
	DaysRemaining float64 `json:"days_remaining"`
}

// ─── Request DTOs ──────────────────────────────────────────────────────────────

// CreateEmployeeRequest for HU-PAY-001
type CreateEmployeeRequest struct {
	EmployeeCode   string `json:"employee_code" binding:"required"`
	ThirdPartyID   int64  `json:"third_party_id" binding:"required"`
	UserID         *int64 `json:"user_id"`
	EmploymentType string `json:"employment_type" binding:"required"`
	Position       string `json:"position"`
	Department     string `json:"department"`
	HireDate       string `json:"hire_date" binding:"required"`
	BankName       string `json:"bank_name"`
	BankAccount    string `json:"bank_account"`
	BankCLABE      string `json:"bank_clabe"`
	CURP           string `json:"curp"`
	RFC            string `json:"rfc"`
	NSS            string `json:"nss"`
}

// UpdateEmployeeRequest for HU-PAY-002
type UpdateEmployeeRequest struct {
	Position    string `json:"position"`
	Department  string `json:"department"`
	BankName    string `json:"bank_name"`
	BankAccount string `json:"bank_account"`
	BankCLABE   string `json:"bank_clabe"`
}

// TerminateEmployeeRequest for HU-PAY-003
type TerminateEmployeeRequest struct {
	TerminationDate   string `json:"termination_date" binding:"required"`
	TerminationReason string `json:"termination_reason"`
}

// SetSalaryRequest for HU-PAY-007
type SetSalaryRequest struct {
	BaseSalary    float64 `json:"base_salary" binding:"required,gt=0"`
	SalaryType    string  `json:"salary_type" binding:"required"`
	EffectiveDate string  `json:"effective_date" binding:"required"`
	Notes         string  `json:"notes"`
}

// CreatePeriodRequest for HU-PAY-011
type CreatePeriodRequest struct {
	PeriodType  string `json:"period_type" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	PaymentDate string `json:"payment_date" binding:"required"`
	Notes       string `json:"notes"`
}

// RegisterOvertimeRequest for HU-PAY-018
type RegisterOvertimeRequest struct {
	EmployeeID int64   `json:"employee_id" binding:"required"`
	WorkDate   string  `json:"work_date" binding:"required"`
	Hours      float64 `json:"hours" binding:"required,gt=0"`
	RateType   string  `json:"rate_type" binding:"required"`
	Notes      string  `json:"notes"`
}

// RegisterBonusRequest for HU-PAY-019
type RegisterBonusRequest struct {
	EmployeeID int64   `json:"employee_id" binding:"required"`
	BonusType  string  `json:"bonus_type" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	BonusDate  string  `json:"bonus_date" binding:"required"`
	Notes      string  `json:"notes"`
}

// CreateLoanRequest for HU-PAY-010
type CreateLoanRequest struct {
	EmployeeID        int64   `json:"employee_id" binding:"required"`
	LoanType          string  `json:"loan_type" binding:"required"`
	TotalAmount       float64 `json:"total_amount" binding:"required,gt=0"`
	InstallmentsTotal int     `json:"installments_total" binding:"required,gt=0"`
	StartDate         string  `json:"start_date" binding:"required"`
	Notes             string  `json:"notes"`
}

// PayrollPaymentRequest for HU-PAY-021
type PayrollPaymentRequest struct {
	PayrollIDs      []int64 `json:"payroll_ids" binding:"required,min=1"`
	PaymentMethod   string  `json:"payment_method" binding:"required"`
	PaymentDate     string  `json:"payment_date" binding:"required"`
	ReferenceNumber string  `json:"reference_number"`
}

// CreateLeaveTypeRequest for HU-PAY-030
type CreateLeaveTypeRequest struct {
	Code             string `json:"code" binding:"required"`
	Name             string `json:"name" binding:"required"`
	IsPaid           bool   `json:"is_paid"`
	DaysAllowed      int    `json:"days_allowed"`
	RequiresApproval bool   `json:"requires_approval"`
}

// CreateLeaveRequest for HU-PAY-031
type CreateLeaveRequest struct {
	EmployeeID  int64  `json:"employee_id" binding:"required"`
	LeaveTypeID int64  `json:"leave_type_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	Reason      string `json:"reason"`
}

// ─── Repository Interface ─────────────────────────────────────────────────────

type Repository interface {
	// Employee
	CreateEmployee(ctx context.Context, e *Employee) (int64, error)
	GetEmployeeByID(ctx context.Context, id int64) (*Employee, error)
	GetEmployeeByCode(ctx context.Context, code string) (*Employee, error)
	UpdateEmployee(ctx context.Context, e *Employee) error
	ListEmployees(ctx context.Context, department string, status string, page, limit int) ([]Employee, int64, error)

	// Salary
	CreateSalary(ctx context.Context, s *Salary) (int64, error)
	GetCurrentSalary(ctx context.Context, employeeID int64) (*Salary, error)
	GetSalaryHistory(ctx context.Context, employeeID int64) ([]Salary, error)
	UpdateSalary(ctx context.Context, s *Salary) error

	// DeductionType & AdditionType
	ListDeductionTypes(ctx context.Context, activeOnly bool) ([]DeductionType, error)
	ListAdditionTypes(ctx context.Context, activeOnly bool) ([]AdditionType, error)

	// PayrollPeriod
	CreatePeriod(ctx context.Context, p *PayrollPeriod) (int64, error)
	GetPeriodByID(ctx context.Context, id int64) (*PayrollPeriod, error)
	UpdatePeriod(ctx context.Context, p *PayrollPeriod) error
	ListPeriods(ctx context.Context, status string, page, limit int) ([]PayrollPeriod, int64, error)

	// Payroll
	CreatePayroll(ctx context.Context, p *Payroll) (int64, error)
	GetPayrollByID(ctx context.Context, id int64) (*Payroll, error)
	UpdatePayroll(ctx context.Context, p *Payroll) error
	ListPayrolls(ctx context.Context, periodID int64, status string, page, limit int) ([]Payroll, int64, error)
	GetPayrollByEmployeeAndPeriod(ctx context.Context, employeeID, periodID int64) (*Payroll, error)

	// PayrollDetail
	CreatePayrollDetail(ctx context.Context, d *PayrollDetail) error
	GetPayrollDetails(ctx context.Context, payrollID int64) ([]PayrollDetail, error)

	// EmployeeLoan
	CreateLoan(ctx context.Context, l *EmployeeLoan) (int64, error)
	GetLoanByID(ctx context.Context, id int64) (*EmployeeLoan, error)
	UpdateLoan(ctx context.Context, l *EmployeeLoan) error
	ListLoans(ctx context.Context, employeeID int64, status string) ([]EmployeeLoan, error)

	// Overtime
	CreateOvertime(ctx context.Context, o *Overtime) (int64, error)
	GetOvertimeByID(ctx context.Context, id int64) (*Overtime, error)
	ListOvertime(ctx context.Context, employeeID *int64, status string, startDate, endDate *time.Time) ([]Overtime, error)
	ApproveOvertime(ctx context.Context, id int64, approvedBy int64) error

	// Bonus
	CreateBonus(ctx context.Context, b *Bonus) (int64, error)
	ListBonuses(ctx context.Context, employeeID *int64, status string) ([]Bonus, error)

	// Payment
	CreatePayment(ctx context.Context, p *PayrollPayment) (int64, error)
	ListPayments(ctx context.Context, periodID *int64, employeeID *int64) ([]PayrollPayment, error)

	// Reports
	GetPayrollSummary(ctx context.Context, periodID int64) ([]PayrollSummary, error)
	GetPayrollTotals(ctx context.Context, periodID int64) (gross, deductions, net float64, count int, err error)
	GetTaxReport(ctx context.Context, periodID int64) ([]TaxReportItem, error)
	GetDeductionReport(ctx context.Context, periodID int64) ([]DeductionReportItem, error)
	GetEmployeeEarnings(ctx context.Context, employeeID int64, startDate, endDate *time.Time) ([]Payroll, error)

	// Leave Management
	CreateLeaveType(ctx context.Context, lt *LeaveType) (int64, error)
	GetLeaveTypeByID(ctx context.Context, id int64) (*LeaveType, error)
	ListLeaveTypes(ctx context.Context, activeOnly bool) ([]LeaveType, error)
	CreateLeave(ctx context.Context, el *EmployeeLeave) (int64, error)
	GetLeaveByID(ctx context.Context, id int64) (*EmployeeLeave, error)
	UpdateLeave(ctx context.Context, el *EmployeeLeave) error
	ListLeaves(ctx context.Context, employeeID *int64, status string, startDate, endDate *time.Time) ([]EmployeeLeave, error)
	GetLeaveBalance(ctx context.Context, employeeID int64, year int) ([]LeaveBalance, error)

	// Commissions Import
	GetApprovedCommissions(ctx context.Context, startDate, endDate *time.Time) ([]interface{}, error)
}

// ─── Service Interface ────────────────────────────────────────────────────────

type Service interface {
	// HU-PAY-001: Register Employee
	CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*Employee, error)

	// HU-PAY-002: Update Employee Profile
	UpdateEmployee(ctx context.Context, id int64, req UpdateEmployeeRequest) (*Employee, error)

	// HU-PAY-003: Deactivate Employee
	TerminateEmployee(ctx context.Context, id int64, req TerminateEmployeeRequest) (*Employee, error)

	// HU-PAY-004: View Employee List
	ListEmployees(ctx context.Context, department string, status string, page, limit int) ([]Employee, int64, error)

	// HU-PAY-007: Set Employee Salary
	SetSalary(ctx context.Context, employeeID int64, req SetSalaryRequest) (*Salary, error)

	// HU-PAY-011: Create Payroll Period
	CreatePeriod(ctx context.Context, req CreatePeriodRequest) (*PayrollPeriod, error)

	// HU-PAY-012: Close Payroll Period
	ClosePeriod(ctx context.Context, periodID int64) error

	// HU-PAY-014: Calculate Payroll
	CalculatePayroll(ctx context.Context, periodID int64) error

	// HU-PAY-016: Approve Payroll
	ApprovePayroll(ctx context.Context, periodID int64, approvedBy int64) error

	// HU-PAY-018: Register Overtime
	RegisterOvertime(ctx context.Context, req RegisterOvertimeRequest) (*Overtime, error)

	// HU-PAY-019: Register Bonus
	RegisterBonus(ctx context.Context, req RegisterBonusRequest) (*Bonus, error)

	// HU-PAY-010: Manage Employee Loans
	CreateLoan(ctx context.Context, req CreateLoanRequest) (*EmployeeLoan, error)

	// HU-PAY-013: Reopen Payroll Period
	ReopenPeriod(ctx context.Context, periodID int64, reason string) error

	// HU-PAY-015: Preview Payroll
	PreviewPayroll(ctx context.Context, periodID int64) ([]PayrollSummary, error)

	// HU-PAY-017: Reject Payroll
	RejectPayroll(ctx context.Context, periodID int64, reason string) error

	// HU-PAY-020: Import Commissions
	ImportCommissions(ctx context.Context, periodID int64) error

	// HU-PAY-021: Process Payment
	ProcessPayment(ctx context.Context, req PayrollPaymentRequest) error

	// HU-PAY-024: View Payment History
	ListPayments(ctx context.Context, periodID *int64, employeeID *int64) ([]PayrollPayment, error)

	// HU-PAY-026: Payroll Summary Report
	GetPayrollSummary(ctx context.Context, periodID int64) ([]PayrollSummary, error)

	// HU-PAY-027: Tax Report
	GetTaxReport(ctx context.Context, periodID int64) ([]TaxReportItem, error)

	// HU-PAY-028: Employee Earnings Report
	GetEmployeeEarnings(ctx context.Context, employeeID int64, startDate, endDate *time.Time) ([]Payroll, error)

	// HU-PAY-029: Deduction Report
	GetDeductionReport(ctx context.Context, periodID int64) ([]DeductionReportItem, error)

	// HU-PAY-030: Configure Leave Types
	CreateLeaveType(ctx context.Context, req CreateLeaveTypeRequest) (*LeaveType, error)
	ListLeaveTypes(ctx context.Context, activeOnly bool) ([]LeaveType, error)

	// HU-PAY-031: Register Leave Request
	CreateLeaveRequest(ctx context.Context, req CreateLeaveRequest) (*EmployeeLeave, error)

	// HU-PAY-032: Approve Leave
	ApproveLeave(ctx context.Context, leaveID int64, approvedBy int64, approved bool) (*EmployeeLeave, error)

	// HU-PAY-034: View Leave Balance
	GetLeaveBalance(ctx context.Context, employeeID int64, year int) ([]LeaveBalance, error)

	// Additional
	GetEmployeeByID(ctx context.Context, id int64) (*Employee, error)
	GetPeriodByID(ctx context.Context, id int64) (*PayrollPeriod, error)
	ListPeriods(ctx context.Context, status string, page, limit int) ([]PayrollPeriod, int64, error)
	ListPayrolls(ctx context.Context, periodID int64, status string, page, limit int) ([]Payroll, int64, error)
	GetPayrollByID(ctx context.Context, id int64) (*Payroll, error)
	ListDeductionTypes(ctx context.Context, activeOnly bool) ([]DeductionType, error)
	ListAdditionTypes(ctx context.Context, activeOnly bool) ([]AdditionType, error)
	GetPayrollDetails(ctx context.Context, payrollID int64) ([]PayrollDetail, error)
}
