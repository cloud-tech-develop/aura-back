package payroll

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"
)

var (
	ErrEmployeeNotFound    = errors.New("empleado no encontrado")
	ErrSalaryNotFound      = errors.New("salario no configurado")
	ErrPeriodNotFound      = errors.New("período no encontrado")
	ErrPeriodNotOpen       = errors.New("período no está abierto")
	ErrPeriodNotProcessing = errors.New("período no está en procesamiento")
	ErrPayrollNotFound     = errors.New("registro de nómina no encontrado")
	ErrEmployeeCodeExists  = errors.New("código de empleado ya existe")
	ErrLoanNotFound        = errors.New("préstamo no encontrado")
)

type service struct {
	repo Repository
}

func NewService(db *sql.DB) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) generateEmployeeCode() string {
	return fmt.Sprintf("EMP-%d", time.Now().UnixNano()%100000)
}

// ─── HU-PAY-001: Register Employee ───────────────────────────────────────────

func (s *service) CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*Employee, error) {
	// Check if code exists
	_, err := s.repo.GetEmployeeByCode(ctx, req.EmployeeCode)
	if err == nil {
		return nil, ErrEmployeeCodeExists
	}

	hireDate, err := time.Parse("2006-01-02", req.HireDate)
	if err != nil {
		return nil, fmt.Errorf("fecha de contratación inválida: %w", err)
	}

	employee := &Employee{
		EmployeeCode:   req.EmployeeCode,
		ThirdPartyID:   req.ThirdPartyID,
		UserID:         req.UserID,
		EmploymentType: req.EmploymentType,
		Position:       req.Position,
		Department:     req.Department,
		HireDate:       hireDate,
		BankName:       req.BankName,
		BankAccount:    req.BankAccount,
		BankCLABE:      req.BankCLABE,
		CURP:           req.CURP,
		RFC:            req.RFC,
		NSS:            req.NSS,
	}

	id, err := s.repo.CreateEmployee(ctx, employee)
	if err != nil {
		return nil, fmt.Errorf("creando empleado: %w", err)
	}
	employee.ID = id

	return employee, nil
}

// ─── HU-PAY-002: Update Employee Profile ──────────────────────────────────────

func (s *service) UpdateEmployee(ctx context.Context, id int64, req UpdateEmployeeRequest) (*Employee, error) {
	employee, err := s.repo.GetEmployeeByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("obteniendo empleado: %w", err)
	}

	if req.Position != "" {
		employee.Position = req.Position
	}
	if req.Department != "" {
		employee.Department = req.Department
	}
	if req.BankName != "" {
		employee.BankName = req.BankName
	}
	if req.BankAccount != "" {
		employee.BankAccount = req.BankAccount
	}
	if req.BankCLABE != "" {
		employee.BankCLABE = req.BankCLABE
	}

	if err := s.repo.UpdateEmployee(ctx, employee); err != nil {
		return nil, fmt.Errorf("actualizando empleado: %w", err)
	}

	return employee, nil
}

// ─── HU-PAY-003: Deactivate Employee ─────────────────────────────────────────

func (s *service) TerminateEmployee(ctx context.Context, id int64, req TerminateEmployeeRequest) (*Employee, error) {
	employee, err := s.repo.GetEmployeeByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("obteniendo empleado: %w", err)
	}

	termDate, err := time.Parse("2006-01-02", req.TerminationDate)
	if err != nil {
		return nil, fmt.Errorf("fecha de baja inválida: %w", err)
	}

	employee.TerminationDate = &termDate
	employee.TerminationReason = req.TerminationReason
	employee.Status = StatusTerminated

	if err := s.repo.UpdateEmployee(ctx, employee); err != nil {
		return nil, fmt.Errorf("dando de baja empleado: %w", err)
	}

	return employee, nil
}

// ─── HU-PAY-004: View Employee List ───────────────────────────────────────────

func (s *service) ListEmployees(ctx context.Context, department string, status string, page, limit int) ([]Employee, int64, error) {
	return s.repo.ListEmployees(ctx, department, status, page, limit)
}

// ─── HU-PAY-007: Set Employee Salary ──────────────────────────────────────────

func (s *service) SetSalary(ctx context.Context, employeeID int64, req SetSalaryRequest) (*Salary, error) {
	effectiveDate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		return nil, fmt.Errorf("fecha efectiva inválida: %w", err)
	}

	// Deactivate current salary
	currentSalary, err := s.repo.GetCurrentSalary(ctx, employeeID)
	if err == nil && currentSalary != nil {
		now := time.Now()
		currentSalary.EndDate = &now
		if err := s.repo.UpdateSalary(ctx, currentSalary); err != nil {
			return nil, fmt.Errorf("desactivando salario anterior: %w", err)
		}
	}

	// Calculate daily salary
	dailySalary := req.BaseSalary / 30

	salary := &Salary{
		EmployeeID:    employeeID,
		BaseSalary:    req.BaseSalary,
		DailySalary:   dailySalary,
		SalaryType:    req.SalaryType,
		EffectiveDate: effectiveDate,
		Notes:         req.Notes,
	}

	id, err := s.repo.CreateSalary(ctx, salary)
	if err != nil {
		return nil, fmt.Errorf("creando salario: %w", err)
	}
	salary.ID = id

	return salary, nil
}

// ─── HU-PAY-011: Create Payroll Period ────────────────────────────────────────

func (s *service) CreatePeriod(ctx context.Context, req CreatePeriodRequest) (*PayrollPeriod, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("fecha inicio inválida: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("fecha fin inválida: %w", err)
	}
	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		return nil, fmt.Errorf("fecha pago inválida: %w", err)
	}

	if endDate.Before(startDate) {
		return nil, errors.New("fecha fin debe ser mayor a fecha inicio")
	}

	period := &PayrollPeriod{
		PeriodType:  req.PeriodType,
		StartDate:   startDate,
		EndDate:     endDate,
		PaymentDate: paymentDate,
		Notes:       req.Notes,
	}

	id, err := s.repo.CreatePeriod(ctx, period)
	if err != nil {
		return nil, fmt.Errorf("creando período: %w", err)
	}
	period.ID = id

	return period, nil
}

// ─── HU-PAY-012: Close Payroll Period ─────────────────────────────────────────

func (s *service) ClosePeriod(ctx context.Context, periodID int64) error {
	period, err := s.repo.GetPeriodByID(ctx, periodID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPeriodNotFound
		}
		return fmt.Errorf("obteniendo período: %w", err)
	}

	if period.Status != PeriodStatusApproved {
		return errors.New("período debe estar aprobado para cerrar")
	}

	period.Status = PeriodStatusClosed
	if err := s.repo.UpdatePeriod(ctx, period); err != nil {
		return fmt.Errorf("cerrando período: %w", err)
	}

	return nil
}

// ─── HU-PAY-014: Calculate Payroll ────────────────────────────────────────────

func (s *service) CalculatePayroll(ctx context.Context, periodID int64) error {
	period, err := s.repo.GetPeriodByID(ctx, periodID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPeriodNotFound
		}
		return fmt.Errorf("obteniendo período: %w", err)
	}

	if period.Status != PeriodStatusOpen {
		return ErrPeriodNotOpen
	}

	// Set to processing
	period.Status = PeriodStatusProcessing
	if err := s.repo.UpdatePeriod(ctx, period); err != nil {
		return fmt.Errorf("actualizando período: %w", err)
	}

	// Get all active employees
	employees, _, err := s.repo.ListEmployees(ctx, "", StatusActive, 1, 1000)
	if err != nil {
		return fmt.Errorf("obteniendo empleados: %w", err)
	}

	for _, emp := range employees {
		// Get current salary
		salary, err := s.repo.GetCurrentSalary(ctx, emp.ID)
		if err != nil {
			continue // Skip employees without salary
		}

		// Check if payroll already exists
		existing, _ := s.repo.GetPayrollByEmployeeAndPeriod(ctx, emp.ID, periodID)
		if existing != nil {
			continue // Skip already calculated
		}

		// Calculate worked days
		workedDays := 15.0 // Default bi-weekly
		if period.PeriodType == "MONTHLY" {
			workedDays = 30.0
		}

		// Calculate base amount
		baseAmount := salary.DailySalary * workedDays

		// Calculate additions (overtime, bonuses)
		totalAdditions := 0.0
		overtimes, _ := s.repo.ListOvertime(ctx, &emp.ID, "APPROVED", &period.StartDate, &period.EndDate)
		for _, ot := range overtimes {
			totalAdditions += ot.TotalAmount
		}

		bonuses, _ := s.repo.ListBonuses(ctx, &emp.ID, "APPROVED")
		for _, b := range bonuses {
			if !b.BonusDate.Before(period.StartDate) && !b.BonusDate.After(period.EndDate) {
				totalAdditions += b.Amount
			}
		}

		// Calculate gross amount
		grossAmount := baseAmount + totalAdditions

		// Calculate deductions
		totalDeductions := 0.0

		// ISR (simplified 10% for demo)
		isr := grossAmount * 0.10
		totalDeductions += isr

		// IMSS (2.75%)
		imss := grossAmount * 0.0275
		totalDeductions += imss

		// Loan deductions
		loans, _ := s.repo.ListLoans(ctx, emp.ID, "ACTIVE")
		for _, loan := range loans {
			if loan.InstallmentsPaid < loan.InstallmentsTotal {
				installment := loan.InstallmentAmount
				if installment > loan.RemainingAmount {
					installment = loan.RemainingAmount
				}
				totalDeductions += installment
			}
		}

		// Calculate net amount
		netAmount := grossAmount - totalDeductions
		netAmount = math.Round(netAmount*100) / 100

		// Create payroll record
		payroll := &Payroll{
			PayrollPeriodID: periodID,
			EmployeeID:      emp.ID,
			SalaryID:        &salary.ID,
			BaseAmount:      baseAmount,
			WorkedDays:      workedDays,
			GrossAmount:     grossAmount,
			TotalAdditions:  totalAdditions,
			TotalDeductions: totalDeductions,
			NetAmount:       netAmount,
			Status:          PayrollStatusCalculated,
		}

		payrollID, err := s.repo.CreatePayroll(ctx, payroll)
		if err != nil {
			continue
		}

		// Create payroll details
		// Base salary detail
		s.repo.CreatePayrollDetail(ctx, &PayrollDetail{
			PayrollID:   payrollID,
			ConceptType: "ADDITION",
			ConceptCode: "SALARY",
			ConceptName: "Salario Base",
			Quantity:    workedDays,
			Rate:        salary.DailySalary,
			Amount:      baseAmount,
		})

		// Additions
		for _, ot := range overtimes {
			s.repo.CreatePayrollDetail(ctx, &PayrollDetail{
				PayrollID:     payrollID,
				ConceptType:   "ADDITION",
				ConceptCode:   "OVERTIME",
				ConceptName:   "Horas Extra",
				Quantity:      ot.Hours,
				Rate:          ot.HourlyRate,
				Amount:        ot.TotalAmount,
				ReferenceID:   &ot.ID,
				ReferenceType: stringPtr("overtime"),
			})
		}

		// Deductions
		s.repo.CreatePayrollDetail(ctx, &PayrollDetail{
			PayrollID:   payrollID,
			ConceptType: "DEDUCTION",
			ConceptCode: "ISR",
			ConceptName: "ISR",
			Rate:        10.0,
			Amount:      isr,
		})

		s.repo.CreatePayrollDetail(ctx, &PayrollDetail{
			PayrollID:   payrollID,
			ConceptType: "DEDUCTION",
			ConceptCode: "IMSS",
			ConceptName: "IMSS",
			Rate:        2.75,
			Amount:      imss,
		})

		// Loan deductions
		for _, loan := range loans {
			if loan.InstallmentsPaid < loan.InstallmentsTotal {
				installment := loan.InstallmentAmount
				if installment > loan.RemainingAmount {
					installment = loan.RemainingAmount
				}
				s.repo.CreatePayrollDetail(ctx, &PayrollDetail{
					PayrollID:     payrollID,
					ConceptType:   "DEDUCTION",
					ConceptCode:   loan.LoanType,
					ConceptName:   "Préstamo",
					Amount:        installment,
					ReferenceID:   &loan.ID,
					ReferenceType: stringPtr("loan"),
				})
			}
		}
	}

	// Update period totals
	gross, deductions, net, count, _ := s.repo.GetPayrollTotals(ctx, periodID)
	period.TotalGross = gross
	period.TotalDeductions = deductions
	period.TotalNet = net
	period.EmployeeCount = count
	s.repo.UpdatePeriod(ctx, period)

	return nil
}

func stringPtr(s string) *string {
	return &s
}

// ─── HU-PAY-016: Approve Payroll ──────────────────────────────────────────────

func (s *service) ApprovePayroll(ctx context.Context, periodID int64, approvedBy int64) error {
	period, err := s.repo.GetPeriodByID(ctx, periodID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPeriodNotFound
		}
		return fmt.Errorf("obteniendo período: %w", err)
	}

	if period.Status != PeriodStatusProcessing {
		return ErrPeriodNotProcessing
	}

	// Update all payrolls to approved
	payrolls, _, err := s.repo.ListPayrolls(ctx, periodID, PayrollStatusCalculated, 1, 1000)
	if err != nil {
		return fmt.Errorf("obteniendo nóminas: %w", err)
	}

	for _, p := range payrolls {
		p.Status = PayrollStatusApproved
		s.repo.UpdatePayroll(ctx, &p)
	}

	// Update period
	now := time.Now()
	period.Status = PeriodStatusApproved
	period.ApprovedBy = &approvedBy
	period.ApprovedAt = &now
	if err := s.repo.UpdatePeriod(ctx, period); err != nil {
		return fmt.Errorf("aprobando período: %w", err)
	}

	return nil
}

// ─── HU-PAY-018: Register Overtime ────────────────────────────────────────────

func (s *service) RegisterOvertime(ctx context.Context, req RegisterOvertimeRequest) (*Overtime, error) {
	workDate, err := time.Parse("2006-01-02", req.WorkDate)
	if err != nil {
		return nil, fmt.Errorf("fecha inválida: %w", err)
	}

	// Get employee's daily salary to calculate hourly rate
	salary, err := s.repo.GetCurrentSalary(ctx, req.EmployeeID)
	if err != nil {
		return nil, ErrSalaryNotFound
	}

	hourlyRate := salary.DailySalary / 8
	rateMultiplier := 2.0
	if req.RateType == "TRIPLE" {
		rateMultiplier = 3.0
	}

	totalAmount := req.Hours * hourlyRate * rateMultiplier

	overtime := &Overtime{
		EmployeeID:  req.EmployeeID,
		WorkDate:    workDate,
		Hours:       req.Hours,
		RateType:    req.RateType,
		HourlyRate:  hourlyRate,
		TotalAmount: totalAmount,
		Notes:       req.Notes,
	}

	id, err := s.repo.CreateOvertime(ctx, overtime)
	if err != nil {
		return nil, fmt.Errorf("registrando horas extra: %w", err)
	}
	overtime.ID = id

	return overtime, nil
}

// ─── HU-PAY-019: Register Bonus ───────────────────────────────────────────────

func (s *service) RegisterBonus(ctx context.Context, req RegisterBonusRequest) (*Bonus, error) {
	bonusDate, err := time.Parse("2006-01-02", req.BonusDate)
	if err != nil {
		return nil, fmt.Errorf("fecha inválida: %w", err)
	}

	bonus := &Bonus{
		EmployeeID: req.EmployeeID,
		BonusType:  req.BonusType,
		Amount:     req.Amount,
		BonusDate:  bonusDate,
		Notes:      req.Notes,
	}

	id, err := s.repo.CreateBonus(ctx, bonus)
	if err != nil {
		return nil, fmt.Errorf("registrando bono: %w", err)
	}
	bonus.ID = id

	return bonus, nil
}

// ─── HU-PAY-010: Manage Employee Loans ─────────────────────────────────────────

func (s *service) CreateLoan(ctx context.Context, req CreateLoanRequest) (*EmployeeLoan, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("fecha inválida: %w", err)
	}

	loan := &EmployeeLoan{
		EmployeeID:        req.EmployeeID,
		LoanType:          req.LoanType,
		TotalAmount:       req.TotalAmount,
		InstallmentsTotal: req.InstallmentsTotal,
		StartDate:         startDate,
		Notes:             req.Notes,
	}

	id, err := s.repo.CreateLoan(ctx, loan)
	if err != nil {
		return nil, fmt.Errorf("creando préstamo: %w", err)
	}
	loan.ID = id

	return loan, nil
}

// ─── HU-PAY-021: Process Payment ──────────────────────────────────────────────

func (s *service) ProcessPayment(ctx context.Context, req PayrollPaymentRequest) error {
	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		return fmt.Errorf("fecha de pago inválida: %w", err)
	}

	for _, payrollID := range req.PayrollIDs {
		payroll, err := s.repo.GetPayrollByID(ctx, payrollID)
		if err != nil {
			continue
		}

		if payroll.Status != PayrollStatusApproved {
			continue
		}

		// Create payment
		payment := &PayrollPayment{
			PayrollID:       payrollID,
			EmployeeID:      payroll.EmployeeID,
			PaymentMethod:   req.PaymentMethod,
			Amount:          payroll.NetAmount,
			PaymentDate:     paymentDate,
			ReferenceNumber: &req.ReferenceNumber,
		}

		s.repo.CreatePayment(ctx, payment)

		// Update payroll status
		payroll.Status = PayrollStatusPaid
		s.repo.UpdatePayroll(ctx, payroll)

		// Update loan installments
		loans, _ := s.repo.ListLoans(ctx, payroll.EmployeeID, "ACTIVE")
		for _, loan := range loans {
			if loan.InstallmentsPaid < loan.InstallmentsTotal {
				loan.InstallmentsPaid++
				loan.RemainingAmount -= loan.InstallmentAmount
				if loan.RemainingAmount <= 0 {
					loan.Status = "COMPLETED"
				}
				s.repo.UpdateLoan(ctx, &loan)
			}
		}
	}

	return nil
}

// ─── HU-PAY-026: Payroll Summary Report ───────────────────────────────────────

func (s *service) GetPayrollSummary(ctx context.Context, periodID int64) ([]PayrollSummary, error) {
	return s.repo.GetPayrollSummary(ctx, periodID)
}

// ─── Additional ──────────────────────────────────────────────────────────────

func (s *service) GetEmployeeByID(ctx context.Context, id int64) (*Employee, error) {
	return s.repo.GetEmployeeByID(ctx, id)
}

func (s *service) GetPeriodByID(ctx context.Context, id int64) (*PayrollPeriod, error) {
	return s.repo.GetPeriodByID(ctx, id)
}

func (s *service) ListPeriods(ctx context.Context, status string, page, limit int) ([]PayrollPeriod, int64, error) {
	return s.repo.ListPeriods(ctx, status, page, limit)
}

func (s *service) ListPayrolls(ctx context.Context, periodID int64, status string, page, limit int) ([]Payroll, int64, error) {
	return s.repo.ListPayrolls(ctx, periodID, status, page, limit)
}

func (s *service) GetPayrollByID(ctx context.Context, id int64) (*Payroll, error) {
	return s.repo.GetPayrollByID(ctx, id)
}

func (s *service) ListDeductionTypes(ctx context.Context, activeOnly bool) ([]DeductionType, error) {
	return s.repo.ListDeductionTypes(ctx, activeOnly)
}

func (s *service) ListAdditionTypes(ctx context.Context, activeOnly bool) ([]AdditionType, error) {
	return s.repo.ListAdditionTypes(ctx, activeOnly)
}

func (s *service) GetPayrollDetails(ctx context.Context, payrollID int64) ([]PayrollDetail, error) {
	return s.repo.GetPayrollDetails(ctx, payrollID)
}

// ─── HU-PAY-013: Reopen Payroll Period ──────────────────────────────────────

func (s *service) ReopenPeriod(ctx context.Context, periodID int64, reason string) error {
	period, err := s.repo.GetPeriodByID(ctx, periodID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPeriodNotFound
		}
		return fmt.Errorf("obteniendo período: %w", err)
	}

	if period.Status != PeriodStatusApproved {
		return errors.New("solo se pueden reabrir períodos aprobados")
	}

	// Reset payrolls to calculated
	payrolls, _, _ := s.repo.ListPayrolls(ctx, periodID, PayrollStatusApproved, 1, 1000)
	for _, p := range payrolls {
		p.Status = PayrollStatusCalculated
		s.repo.UpdatePayroll(ctx, &p)
	}

	period.Status = PeriodStatusProcessing
	period.Notes = reason
	if err := s.repo.UpdatePeriod(ctx, period); err != nil {
		return fmt.Errorf("reabriendo período: %w", err)
	}

	return nil
}

// ─── HU-PAY-015: Preview Payroll ────────────────────────────────────────────

func (s *service) PreviewPayroll(ctx context.Context, periodID int64) ([]PayrollSummary, error) {
	return s.repo.GetPayrollSummary(ctx, periodID)
}

// ─── HU-PAY-017: Reject Payroll ──────────────────────────────────────────────

func (s *service) RejectPayroll(ctx context.Context, periodID int64, reason string) error {
	period, err := s.repo.GetPeriodByID(ctx, periodID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPeriodNotFound
		}
		return fmt.Errorf("obteniendo período: %w", err)
	}

	if period.Status != PeriodStatusProcessing {
		return errors.New("solo se pueden rechazar períodos en procesamiento")
	}

	// Delete all calculated payrolls
	payrolls, _, _ := s.repo.ListPayrolls(ctx, periodID, PayrollStatusCalculated, 1, 1000)
	for _, p := range payrolls {
		p.Status = PayrollStatusCancelled
		s.repo.UpdatePayroll(ctx, &p)
	}

	period.Status = PeriodStatusOpen
	period.Notes = "Rechazado: " + reason
	if err := s.repo.UpdatePeriod(ctx, period); err != nil {
		return fmt.Errorf("rechazando nómina: %w", err)
	}

	return nil
}

// ─── HU-PAY-020: Import Commissions ─────────────────────────────────────────

func (s *service) ImportCommissions(ctx context.Context, periodID int64) error {
	period, err := s.repo.GetPeriodByID(ctx, periodID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPeriodNotFound
		}
		return fmt.Errorf("obteniendo período: %w", err)
	}

	// Get all payrolls for the period
	payrolls, _, err := s.repo.ListPayrolls(ctx, periodID, "", 1, 1000)
	if err != nil {
		return fmt.Errorf("obteniendo nóminas: %w", err)
	}

	// For each payroll, import commissions (simplified)
	for _, p := range payrolls {
		// In a real implementation, this would query the commissions module
		// For now, we'll add a placeholder commission detail
		detail := &PayrollDetail{
			PayrollID:   p.ID,
			ConceptType: "ADDITION",
			ConceptCode: "COMMISSION",
			ConceptName: "Comisión por ventas",
			Amount:      0, // Would be calculated from commissions module
		}
		s.repo.CreatePayrollDetail(ctx, detail)
	}

	_ = period // Use period to avoid unused variable
	return nil
}

// ─── HU-PAY-024: View Payment History ────────────────────────────────────────

func (s *service) ListPayments(ctx context.Context, periodID *int64, employeeID *int64) ([]PayrollPayment, error) {
	return s.repo.ListPayments(ctx, periodID, employeeID)
}

// ─── HU-PAY-027: Tax Report ──────────────────────────────────────────────────

func (s *service) GetTaxReport(ctx context.Context, periodID int64) ([]TaxReportItem, error) {
	return s.repo.GetTaxReport(ctx, periodID)
}

// ─── HU-PAY-028: Employee Earnings Report ────────────────────────────────────

func (s *service) GetEmployeeEarnings(ctx context.Context, employeeID int64, startDate, endDate *time.Time) ([]Payroll, error) {
	return s.repo.GetEmployeeEarnings(ctx, employeeID, startDate, endDate)
}

// ─── HU-PAY-029: Deduction Report ────────────────────────────────────────────

func (s *service) GetDeductionReport(ctx context.Context, periodID int64) ([]DeductionReportItem, error) {
	return s.repo.GetDeductionReport(ctx, periodID)
}

// ─── HU-PAY-030: Configure Leave Types ────────────────────────────────────────

func (s *service) CreateLeaveType(ctx context.Context, req CreateLeaveTypeRequest) (*LeaveType, error) {
	lt := &LeaveType{
		Code:             req.Code,
		Name:             req.Name,
		IsPaid:           req.IsPaid,
		DaysAllowed:      req.DaysAllowed,
		RequiresApproval: req.RequiresApproval,
	}

	id, err := s.repo.CreateLeaveType(ctx, lt)
	if err != nil {
		return nil, fmt.Errorf("creando tipo de permiso: %w", err)
	}
	lt.ID = id

	return lt, nil
}

func (s *service) ListLeaveTypes(ctx context.Context, activeOnly bool) ([]LeaveType, error) {
	return s.repo.ListLeaveTypes(ctx, activeOnly)
}

// ─── HU-PAY-031: Register Leave Request ──────────────────────────────────────

func (s *service) CreateLeaveRequest(ctx context.Context, req CreateLeaveRequest) (*EmployeeLeave, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("fecha inicio inválida: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("fecha fin inválida: %w", err)
	}

	// Calculate days
	days := endDate.Sub(startDate).Hours()/24 + 1

	leave := &EmployeeLeave{
		EmployeeID:  req.EmployeeID,
		LeaveTypeID: req.LeaveTypeID,
		StartDate:   startDate,
		EndDate:     endDate,
		Days:        days,
		Reason:      req.Reason,
	}

	// Check if leave type requires approval
	leaveType, err := s.repo.GetLeaveTypeByID(ctx, req.LeaveTypeID)
	if err == nil && !leaveType.RequiresApproval {
		leave.Status = "APPROVED"
		now := time.Now()
		leave.ApprovedAt = &now
	}

	id, err := s.repo.CreateLeave(ctx, leave)
	if err != nil {
		return nil, fmt.Errorf("creando solicitud de permiso: %w", err)
	}
	leave.ID = id

	return leave, nil
}

// ─── HU-PAY-032: Approve Leave ──────────────────────────────────────────────

func (s *service) ApproveLeave(ctx context.Context, leaveID int64, approvedBy int64, approved bool) (*EmployeeLeave, error) {
	leave, err := s.repo.GetLeaveByID(ctx, leaveID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("solicitud de permiso no encontrada")
		}
		return nil, fmt.Errorf("obteniendo permiso: %w", err)
	}

	if leave.Status != "PENDING" {
		return nil, errors.New("solicitud ya procesada")
	}

	now := time.Now()
	leave.ApprovedBy = &approvedBy
	leave.ApprovedAt = &now
	if approved {
		leave.Status = "APPROVED"
	} else {
		leave.Status = "REJECTED"
	}

	if err := s.repo.UpdateLeave(ctx, leave); err != nil {
		return nil, fmt.Errorf("actualizando permiso: %w", err)
	}

	return leave, nil
}

// ─── HU-PAY-034: View Leave Balance ──────────────────────────────────────────

func (s *service) GetLeaveBalance(ctx context.Context, employeeID int64, year int) ([]LeaveBalance, error) {
	return s.repo.GetLeaveBalance(ctx, employeeID, year)
}
