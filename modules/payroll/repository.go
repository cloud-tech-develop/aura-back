package payroll

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// ─── Employee ─────────────────────────────────────────────────────────────────

func (r *repository) CreateEmployee(ctx context.Context, e *Employee) (int64, error) {
	e.CreatedAt = time.Now()
	e.Status = StatusActive

	return e.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO employee (employee_code, third_party_id, user_id, employment_type, position, department,
			hire_date, bank_name, bank_account, bank_clabe, curp, rfc, nss, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING id`,
		e.EmployeeCode, e.ThirdPartyID, e.UserID, e.EmploymentType, e.Position, e.Department,
		e.HireDate, e.BankName, e.BankAccount, e.BankCLABE, e.CURP, e.RFC, e.NSS, e.Status, e.CreatedAt,
	).Scan(&e.ID)
}

func (r *repository) GetEmployeeByID(ctx context.Context, id int64) (*Employee, error) {
	var e Employee
	var termDate, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, employee_code, third_party_id, user_id, employment_type, position, department,
			hire_date, termination_date, termination_reason, bank_name, bank_account, bank_clabe, curp, rfc, nss, status, created_at, updated_at
		 FROM employee WHERE id = $1`, id,
	).Scan(&e.ID, &e.EmployeeCode, &e.ThirdPartyID, &e.UserID, &e.EmploymentType, &e.Position, &e.Department,
		&e.HireDate, &termDate, &e.TerminationReason, &e.BankName, &e.BankAccount, &e.BankCLABE, &e.CURP, &e.RFC, &e.NSS, &e.Status, &e.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if termDate.Valid {
		e.TerminationDate = &termDate.Time
	}
	if updatedAt.Valid {
		e.UpdatedAt = &updatedAt.Time
	}
	return &e, nil
}

func (r *repository) GetEmployeeByCode(ctx context.Context, code string) (*Employee, error) {
	var e Employee
	var termDate, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, employee_code, third_party_id, user_id, employment_type, position, department,
			hire_date, termination_date, termination_reason, bank_name, bank_account, bank_clabe, curp, rfc, nss, status, created_at, updated_at
		 FROM employee WHERE employee_code = $1`, code,
	).Scan(&e.ID, &e.EmployeeCode, &e.ThirdPartyID, &e.UserID, &e.EmploymentType, &e.Position, &e.Department,
		&e.HireDate, &termDate, &e.TerminationReason, &e.BankName, &e.BankAccount, &e.BankCLABE, &e.CURP, &e.RFC, &e.NSS, &e.Status, &e.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if termDate.Valid {
		e.TerminationDate = &termDate.Time
	}
	if updatedAt.Valid {
		e.UpdatedAt = &updatedAt.Time
	}
	return &e, nil
}

func (r *repository) UpdateEmployee(ctx context.Context, e *Employee) error {
	now := time.Now()
	e.UpdatedAt = &now
	_, err := r.db.ExecContext(ctx,
		`UPDATE employee SET position=$1, department=$2, bank_name=$3, bank_account=$4, bank_clabe=$5, updated_at=$6 WHERE id=$7`,
		e.Position, e.Department, e.BankName, e.BankAccount, e.BankCLABE, now, e.ID)
	return err
}

func (r *repository) ListEmployees(ctx context.Context, department string, status string, page, limit int) ([]Employee, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	baseQuery := "FROM employee WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if department != "" {
		baseQuery += fmt.Sprintf(" AND department = $%d", argIndex)
		args = append(args, department)
		argIndex++
	}
	if status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) "+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf("SELECT id, employee_code, third_party_id, user_id, employment_type, position, department, hire_date, termination_date, termination_reason, bank_name, bank_account, bank_clabe, curp, rfc, nss, status, created_at, updated_at %s ORDER BY employee_code LIMIT $%d OFFSET $%d", baseQuery, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		var termDate, updatedAt sql.NullTime
		if err := rows.Scan(&e.ID, &e.EmployeeCode, &e.ThirdPartyID, &e.UserID, &e.EmploymentType, &e.Position, &e.Department,
			&e.HireDate, &termDate, &e.TerminationReason, &e.BankName, &e.BankAccount, &e.BankCLABE, &e.CURP, &e.RFC, &e.NSS, &e.Status, &e.CreatedAt, &updatedAt); err != nil {
			return nil, 0, err
		}
		if termDate.Valid {
			e.TerminationDate = &termDate.Time
		}
		if updatedAt.Valid {
			e.UpdatedAt = &updatedAt.Time
		}
		employees = append(employees, e)
	}
	return employees, total, nil
}

// ─── Salary ───────────────────────────────────────────────────────────────────

func (r *repository) CreateSalary(ctx context.Context, s *Salary) (int64, error) {
	s.CreatedAt = time.Now()
	s.IsCurrent = true

	return s.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO salary (employee_id, base_salary, daily_salary, salary_type, effective_date, is_current, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		s.EmployeeID, s.BaseSalary, s.DailySalary, s.SalaryType, s.EffectiveDate, s.IsCurrent, s.Notes, s.CreatedAt,
	).Scan(&s.ID)
}

func (r *repository) GetCurrentSalary(ctx context.Context, employeeID int64) (*Salary, error) {
	var s Salary
	var endDate sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, employee_id, base_salary, daily_salary, salary_type, effective_date, end_date, is_current, notes, created_at
		 FROM salary WHERE employee_id = $1 AND is_current = TRUE`, employeeID,
	).Scan(&s.ID, &s.EmployeeID, &s.BaseSalary, &s.DailySalary, &s.SalaryType, &s.EffectiveDate, &endDate, &s.IsCurrent, &s.Notes, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	if endDate.Valid {
		s.EndDate = &endDate.Time
	}
	return &s, nil
}

func (r *repository) GetSalaryHistory(ctx context.Context, employeeID int64) ([]Salary, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, employee_id, base_salary, daily_salary, salary_type, effective_date, end_date, is_current, notes, created_at
		 FROM salary WHERE employee_id = $1 ORDER BY effective_date DESC`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var salaries []Salary
	for rows.Next() {
		var s Salary
		var endDate sql.NullTime
		if err := rows.Scan(&s.ID, &s.EmployeeID, &s.BaseSalary, &s.DailySalary, &s.SalaryType, &s.EffectiveDate, &endDate, &s.IsCurrent, &s.Notes, &s.CreatedAt); err != nil {
			return nil, err
		}
		if endDate.Valid {
			s.EndDate = &endDate.Time
		}
		salaries = append(salaries, s)
	}
	return salaries, nil
}

func (r *repository) UpdateSalary(ctx context.Context, s *Salary) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE salary SET end_date = $1, is_current = FALSE WHERE id = $2`,
		s.EndDate, s.ID)
	return err
}

// ─── DeductionType & AdditionType ──────────────────────────────────────────────

func (r *repository) ListDeductionTypes(ctx context.Context, activeOnly bool) ([]DeductionType, error) {
	query := `SELECT id, code, name, description, is_tax, is_mandatory, calculation_type, default_value, is_active, created_at FROM deduction_type`
	if activeOnly {
		query += " WHERE is_active = TRUE"
	}
	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []DeductionType
	for rows.Next() {
		var dt DeductionType
		var defVal sql.NullFloat64
		if err := rows.Scan(&dt.ID, &dt.Code, &dt.Name, &dt.Description, &dt.IsTax, &dt.IsMandatory, &dt.CalculationType, &defVal, &dt.IsActive, &dt.CreatedAt); err != nil {
			return nil, err
		}
		if defVal.Valid {
			dt.DefaultValue = &defVal.Float64
		}
		types = append(types, dt)
	}
	return types, nil
}

func (r *repository) ListAdditionTypes(ctx context.Context, activeOnly bool) ([]AdditionType, error) {
	query := `SELECT id, code, name, description, is_taxable, calculation_type, default_value, is_active, created_at FROM addition_type`
	if activeOnly {
		query += " WHERE is_active = TRUE"
	}
	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []AdditionType
	for rows.Next() {
		var at AdditionType
		var defVal sql.NullFloat64
		if err := rows.Scan(&at.ID, &at.Code, &at.Name, &at.Description, &at.IsTaxable, &at.CalculationType, &defVal, &at.IsActive, &at.CreatedAt); err != nil {
			return nil, err
		}
		if defVal.Valid {
			at.DefaultValue = &defVal.Float64
		}
		types = append(types, at)
	}
	return types, nil
}

// ─── PayrollPeriod ────────────────────────────────────────────────────────────

func (r *repository) CreatePeriod(ctx context.Context, p *PayrollPeriod) (int64, error) {
	p.CreatedAt = time.Now()
	p.Status = PeriodStatusOpen

	return p.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO payroll_period (period_type, start_date, end_date, payment_date, status, created_at, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		p.PeriodType, p.StartDate, p.EndDate, p.PaymentDate, p.Status, p.CreatedAt, p.Notes,
	).Scan(&p.ID)
}

func (r *repository) GetPeriodByID(ctx context.Context, id int64) (*PayrollPeriod, error) {
	var p PayrollPeriod
	var approvedAt, updatedAt sql.NullTime
	var approvedBy sql.NullInt64

	err := r.db.QueryRowContext(ctx,
		`SELECT id, period_type, start_date, end_date, payment_date, status, total_gross, total_deductions, total_net, employee_count, approved_by, approved_at, notes, created_at, updated_at
		 FROM payroll_period WHERE id = $1`, id,
	).Scan(&p.ID, &p.PeriodType, &p.StartDate, &p.EndDate, &p.PaymentDate, &p.Status, &p.TotalGross, &p.TotalDeductions, &p.TotalNet, &p.EmployeeCount, &approvedBy, &approvedAt, &p.Notes, &p.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if approvedBy.Valid {
		p.ApprovedBy = &approvedBy.Int64
	}
	if approvedAt.Valid {
		p.ApprovedAt = &approvedAt.Time
	}
	if updatedAt.Valid {
		p.UpdatedAt = &updatedAt.Time
	}
	return &p, nil
}

func (r *repository) UpdatePeriod(ctx context.Context, p *PayrollPeriod) error {
	now := time.Now()
	p.UpdatedAt = &now
	_, err := r.db.ExecContext(ctx,
		`UPDATE payroll_period SET status=$1, total_gross=$2, total_deductions=$3, total_net=$4, employee_count=$5, approved_by=$6, approved_at=$7, updated_at=$8 WHERE id=$9`,
		p.Status, p.TotalGross, p.TotalDeductions, p.TotalNet, p.EmployeeCount, p.ApprovedBy, p.ApprovedAt, now, p.ID)
	return err
}

func (r *repository) ListPeriods(ctx context.Context, status string, page, limit int) ([]PayrollPeriod, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	baseQuery := "FROM payroll_period WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) "+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf("SELECT id, period_type, start_date, end_date, payment_date, status, total_gross, total_deductions, total_net, employee_count, approved_by, approved_at, notes, created_at, updated_at %s ORDER BY start_date DESC LIMIT $%d OFFSET $%d", baseQuery, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var periods []PayrollPeriod
	for rows.Next() {
		var p PayrollPeriod
		var approvedAt, updatedAt sql.NullTime
		var approvedBy sql.NullInt64
		if err := rows.Scan(&p.ID, &p.PeriodType, &p.StartDate, &p.EndDate, &p.PaymentDate, &p.Status, &p.TotalGross, &p.TotalDeductions, &p.TotalNet, &p.EmployeeCount, &approvedBy, &approvedAt, &p.Notes, &p.CreatedAt, &updatedAt); err != nil {
			return nil, 0, err
		}
		if approvedBy.Valid {
			p.ApprovedBy = &approvedBy.Int64
		}
		if approvedAt.Valid {
			p.ApprovedAt = &approvedAt.Time
		}
		if updatedAt.Valid {
			p.UpdatedAt = &updatedAt.Time
		}
		periods = append(periods, p)
	}
	return periods, total, nil
}

// ─── Payroll ──────────────────────────────────────────────────────────────────

func (r *repository) CreatePayroll(ctx context.Context, p *Payroll) (int64, error) {
	p.CreatedAt = time.Now()
	p.Status = PayrollStatusDraft

	return p.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO payroll (payroll_period_id, employee_id, salary_id, base_amount, worked_days, gross_amount, total_additions, total_deductions, net_amount, status, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id`,
		p.PayrollPeriodID, p.EmployeeID, p.SalaryID, p.BaseAmount, p.WorkedDays, p.GrossAmount, p.TotalAdditions, p.TotalDeductions, p.NetAmount, p.Status, p.Notes, p.CreatedAt,
	).Scan(&p.ID)
}

func (r *repository) GetPayrollByID(ctx context.Context, id int64) (*Payroll, error) {
	var p Payroll
	var salaryID sql.NullInt64

	err := r.db.QueryRowContext(ctx,
		`SELECT id, payroll_period_id, employee_id, salary_id, base_amount, worked_days, gross_amount, total_additions, total_deductions, net_amount, status, notes, created_at, updated_at
		 FROM payroll WHERE id = $1`, id,
	).Scan(&p.ID, &p.PayrollPeriodID, &p.EmployeeID, &salaryID, &p.BaseAmount, &p.WorkedDays, &p.GrossAmount, &p.TotalAdditions, &p.TotalDeductions, &p.NetAmount, &p.Status, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if salaryID.Valid {
		p.SalaryID = &salaryID.Int64
	}
	return &p, nil
}

func (r *repository) UpdatePayroll(ctx context.Context, p *Payroll) error {
	now := time.Now()
	p.UpdatedAt = &now
	_, err := r.db.ExecContext(ctx,
		`UPDATE payroll SET gross_amount=$1, total_additions=$2, total_deductions=$3, net_amount=$4, status=$5, notes=$6, updated_at=$7 WHERE id=$8`,
		p.GrossAmount, p.TotalAdditions, p.TotalDeductions, p.NetAmount, p.Status, p.Notes, now, p.ID)
	return err
}

func (r *repository) ListPayrolls(ctx context.Context, periodID int64, status string, page, limit int) ([]Payroll, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	baseQuery := "FROM payroll WHERE payroll_period_id = $1"
	args := []interface{}{periodID}
	argIndex := 2

	if status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) "+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf("SELECT id, payroll_period_id, employee_id, salary_id, base_amount, worked_days, gross_amount, total_additions, total_deductions, net_amount, status, notes, created_at, updated_at %s ORDER BY employee_id LIMIT $%d OFFSET $%d", baseQuery, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var payrolls []Payroll
	for rows.Next() {
		var p Payroll
		var salaryID sql.NullInt64
		if err := rows.Scan(&p.ID, &p.PayrollPeriodID, &p.EmployeeID, &salaryID, &p.BaseAmount, &p.WorkedDays, &p.GrossAmount, &p.TotalAdditions, &p.TotalDeductions, &p.NetAmount, &p.Status, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if salaryID.Valid {
			p.SalaryID = &salaryID.Int64
		}
		payrolls = append(payrolls, p)
	}
	return payrolls, total, nil
}

func (r *repository) GetPayrollByEmployeeAndPeriod(ctx context.Context, employeeID, periodID int64) (*Payroll, error) {
	var p Payroll
	var salaryID sql.NullInt64

	err := r.db.QueryRowContext(ctx,
		`SELECT id, payroll_period_id, employee_id, salary_id, base_amount, worked_days, gross_amount, total_additions, total_deductions, net_amount, status, notes, created_at, updated_at
		 FROM payroll WHERE employee_id = $1 AND payroll_period_id = $2`, employeeID, periodID,
	).Scan(&p.ID, &p.PayrollPeriodID, &p.EmployeeID, &salaryID, &p.BaseAmount, &p.WorkedDays, &p.GrossAmount, &p.TotalAdditions, &p.TotalDeductions, &p.NetAmount, &p.Status, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if salaryID.Valid {
		p.SalaryID = &salaryID.Int64
	}
	return &p, nil
}

// ─── PayrollDetail ────────────────────────────────────────────────────────────

func (r *repository) CreatePayrollDetail(ctx context.Context, d *PayrollDetail) error {
	d.CreatedAt = time.Now()
	return r.db.QueryRowContext(ctx,
		`INSERT INTO payroll_detail (payroll_id, concept_type, concept_code, concept_name, quantity, rate, amount, reference_id, reference_type, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id`,
		d.PayrollID, d.ConceptType, d.ConceptCode, d.ConceptName, d.Quantity, d.Rate, d.Amount, d.ReferenceID, d.ReferenceType, d.Notes, d.CreatedAt,
	).Scan(&d.ID)
}

func (r *repository) GetPayrollDetails(ctx context.Context, payrollID int64) ([]PayrollDetail, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, payroll_id, concept_type, concept_code, concept_name, quantity, rate, amount, reference_id, reference_type, notes, created_at
		 FROM payroll_detail WHERE payroll_id = $1 ORDER BY concept_type, concept_code`, payrollID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []PayrollDetail
	for rows.Next() {
		var d PayrollDetail
		var refID sql.NullInt64
		var refType sql.NullString
		if err := rows.Scan(&d.ID, &d.PayrollID, &d.ConceptType, &d.ConceptCode, &d.ConceptName, &d.Quantity, &d.Rate, &d.Amount, &refID, &refType, &d.Notes, &d.CreatedAt); err != nil {
			return nil, err
		}
		if refID.Valid {
			d.ReferenceID = &refID.Int64
		}
		if refType.Valid {
			d.ReferenceType = &refType.String
		}
		details = append(details, d)
	}
	return details, nil
}

// ─── EmployeeLoan ─────────────────────────────────────────────────────────────

func (r *repository) CreateLoan(ctx context.Context, l *EmployeeLoan) (int64, error) {
	l.CreatedAt = time.Now()
	l.Status = "ACTIVE"
	l.RemainingAmount = l.TotalAmount
	l.InstallmentAmount = l.TotalAmount / float64(l.InstallmentsTotal)

	return l.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO employee_loan (employee_id, loan_type, total_amount, remaining_amount, installment_amount, installments_total, start_date, status, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		l.EmployeeID, l.LoanType, l.TotalAmount, l.RemainingAmount, l.InstallmentAmount, l.InstallmentsTotal, l.StartDate, l.Status, l.Notes, l.CreatedAt,
	).Scan(&l.ID)
}

func (r *repository) GetLoanByID(ctx context.Context, id int64) (*EmployeeLoan, error) {
	var l EmployeeLoan
	var endDate, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, employee_id, loan_type, total_amount, remaining_amount, installment_amount, installments_total, installments_paid, start_date, end_date, status, notes, created_at, updated_at
		 FROM employee_loan WHERE id = $1`, id,
	).Scan(&l.ID, &l.EmployeeID, &l.LoanType, &l.TotalAmount, &l.RemainingAmount, &l.InstallmentAmount, &l.InstallmentsTotal, &l.InstallmentsPaid, &l.StartDate, &endDate, &l.Status, &l.Notes, &l.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if endDate.Valid {
		l.EndDate = &endDate.Time
	}
	if updatedAt.Valid {
		l.UpdatedAt = &updatedAt.Time
	}
	return &l, nil
}

func (r *repository) UpdateLoan(ctx context.Context, l *EmployeeLoan) error {
	now := time.Now()
	l.UpdatedAt = &now
	_, err := r.db.ExecContext(ctx,
		`UPDATE employee_loan SET remaining_amount=$1, installments_paid=$2, status=$3, updated_at=$4 WHERE id=$5`,
		l.RemainingAmount, l.InstallmentsPaid, l.Status, now, l.ID)
	return err
}

func (r *repository) ListLoans(ctx context.Context, employeeID int64, status string) ([]EmployeeLoan, error) {
	query := `SELECT id, employee_id, loan_type, total_amount, remaining_amount, installment_amount, installments_total, installments_paid, start_date, end_date, status, notes, created_at, updated_at
		 FROM employee_loan WHERE employee_id = $1`
	args := []interface{}{employeeID}

	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []EmployeeLoan
	for rows.Next() {
		var l EmployeeLoan
		var endDate, updatedAt sql.NullTime
		if err := rows.Scan(&l.ID, &l.EmployeeID, &l.LoanType, &l.TotalAmount, &l.RemainingAmount, &l.InstallmentAmount, &l.InstallmentsTotal, &l.InstallmentsPaid, &l.StartDate, &endDate, &l.Status, &l.Notes, &l.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		if endDate.Valid {
			l.EndDate = &endDate.Time
		}
		if updatedAt.Valid {
			l.UpdatedAt = &updatedAt.Time
		}
		loans = append(loans, l)
	}
	return loans, nil
}

// ─── Overtime ─────────────────────────────────────────────────────────────────

func (r *repository) CreateOvertime(ctx context.Context, o *Overtime) (int64, error) {
	o.CreatedAt = time.Now()
	o.Status = "PENDING"

	return o.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO overtime (employee_id, work_date, hours, rate_type, hourly_rate, total_amount, status, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		o.EmployeeID, o.WorkDate, o.Hours, o.RateType, o.HourlyRate, o.TotalAmount, o.Status, o.Notes, o.CreatedAt,
	).Scan(&o.ID)
}

func (r *repository) GetOvertimeByID(ctx context.Context, id int64) (*Overtime, error) {
	var o Overtime
	var approvedBy sql.NullInt64

	err := r.db.QueryRowContext(ctx,
		`SELECT id, employee_id, work_date, hours, rate_type, hourly_rate, total_amount, status, approved_by, notes, created_at
		 FROM overtime WHERE id = $1`, id,
	).Scan(&o.ID, &o.EmployeeID, &o.WorkDate, &o.Hours, &o.RateType, &o.HourlyRate, &o.TotalAmount, &o.Status, &approvedBy, &o.Notes, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	if approvedBy.Valid {
		o.ApprovedBy = &approvedBy.Int64
	}
	return &o, nil
}

func (r *repository) ListOvertime(ctx context.Context, employeeID *int64, status string, startDate, endDate *time.Time) ([]Overtime, error) {
	query := `SELECT id, employee_id, work_date, hours, rate_type, hourly_rate, total_amount, status, approved_by, notes, created_at FROM overtime WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if employeeID != nil {
		query += fmt.Sprintf(" AND employee_id = $%d", argIndex)
		args = append(args, *employeeID)
		argIndex++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if startDate != nil {
		query += fmt.Sprintf(" AND work_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		query += fmt.Sprintf(" AND work_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}
	query += " ORDER BY work_date DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overtimes []Overtime
	for rows.Next() {
		var o Overtime
		var approvedBy sql.NullInt64
		if err := rows.Scan(&o.ID, &o.EmployeeID, &o.WorkDate, &o.Hours, &o.RateType, &o.HourlyRate, &o.TotalAmount, &o.Status, &approvedBy, &o.Notes, &o.CreatedAt); err != nil {
			return nil, err
		}
		if approvedBy.Valid {
			o.ApprovedBy = &approvedBy.Int64
		}
		overtimes = append(overtimes, o)
	}
	return overtimes, nil
}

func (r *repository) ApproveOvertime(ctx context.Context, id int64, approvedBy int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE overtime SET status = 'APPROVED', approved_by = $1 WHERE id = $2 AND status = 'PENDING'`,
		approvedBy, id)
	return err
}

// ─── Bonus ────────────────────────────────────────────────────────────────────

func (r *repository) CreateBonus(ctx context.Context, b *Bonus) (int64, error) {
	b.CreatedAt = time.Now()
	b.Status = "PENDING"

	return b.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO bonus (employee_id, bonus_type, amount, bonus_date, status, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		b.EmployeeID, b.BonusType, b.Amount, b.BonusDate, b.Status, b.Notes, b.CreatedAt,
	).Scan(&b.ID)
}

func (r *repository) ListBonuses(ctx context.Context, employeeID *int64, status string) ([]Bonus, error) {
	query := `SELECT id, employee_id, bonus_type, amount, bonus_date, status, notes, created_at FROM bonus WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if employeeID != nil {
		query += fmt.Sprintf(" AND employee_id = $%d", argIndex)
		args = append(args, *employeeID)
		argIndex++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	query += " ORDER BY bonus_date DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bonuses []Bonus
	for rows.Next() {
		var b Bonus
		if err := rows.Scan(&b.ID, &b.EmployeeID, &b.BonusType, &b.Amount, &b.BonusDate, &b.Status, &b.Notes, &b.CreatedAt); err != nil {
			return nil, err
		}
		bonuses = append(bonuses, b)
	}
	return bonuses, nil
}

// ─── Payment ──────────────────────────────────────────────────────────────────

func (r *repository) CreatePayment(ctx context.Context, p *PayrollPayment) (int64, error) {
	p.CreatedAt = time.Now()
	p.Status = "COMPLETED"

	return p.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO payroll_payment (payroll_id, employee_id, payment_method, amount, payment_date, reference_number, status, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		p.PayrollID, p.EmployeeID, p.PaymentMethod, p.Amount, p.PaymentDate, p.ReferenceNumber, p.Status, p.Notes, p.CreatedAt,
	).Scan(&p.ID)
}

func (r *repository) ListPayments(ctx context.Context, periodID *int64, employeeID *int64) ([]PayrollPayment, error) {
	query := `SELECT pp.id, pp.payroll_id, pp.employee_id, pp.payment_method, pp.amount, pp.payment_date, pp.reference_number, pp.bank_reference, pp.status, pp.notes, pp.created_at
		 FROM payroll_payment pp
		 JOIN payroll p ON p.id = pp.payroll_id
		 WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if periodID != nil {
		query += fmt.Sprintf(" AND p.payroll_period_id = $%d", argIndex)
		args = append(args, *periodID)
		argIndex++
	}
	if employeeID != nil {
		query += fmt.Sprintf(" AND pp.employee_id = $%d", argIndex)
		args = append(args, *employeeID)
		argIndex++
	}
	query += " ORDER BY pp.payment_date DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []PayrollPayment
	for rows.Next() {
		var p PayrollPayment
		var refNum, bankRef sql.NullString
		if err := rows.Scan(&p.ID, &p.PayrollID, &p.EmployeeID, &p.PaymentMethod, &p.Amount, &p.PaymentDate, &refNum, &bankRef, &p.Status, &p.Notes, &p.CreatedAt); err != nil {
			return nil, err
		}
		if refNum.Valid {
			p.ReferenceNumber = &refNum.String
		}
		if bankRef.Valid {
			p.BankReference = &bankRef.String
		}
		payments = append(payments, p)
	}
	return payments, nil
}

// ─── Reports ──────────────────────────────────────────────────────────────────

func (r *repository) GetPayrollSummary(ctx context.Context, periodID int64) ([]PayrollSummary, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT p.employee_id, e.employee_code, tp.name as employee_name,
			p.gross_amount, p.total_additions, p.total_deductions, p.net_amount, p.status
		 FROM payroll p
		 JOIN employee e ON e.id = p.employee_id
		 JOIN third_parties tp ON tp.id = e.third_party_id
		 WHERE p.payroll_period_id = $1
		 ORDER BY e.employee_code`, periodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []PayrollSummary
	for rows.Next() {
		var s PayrollSummary
		if err := rows.Scan(&s.EmployeeID, &s.EmployeeCode, &s.EmployeeName, &s.GrossAmount, &s.TotalAdditions, &s.TotalDeductions, &s.NetAmount, &s.Status); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}

func (r *repository) GetPayrollTotals(ctx context.Context, periodID int64) (gross, deductions, net float64, count int, err error) {
	err = r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(gross_amount),0), COALESCE(SUM(total_deductions),0), COALESCE(SUM(net_amount),0), COUNT(*)
		 FROM payroll WHERE payroll_period_id = $1`,
		periodID,
	).Scan(&gross, &deductions, &net, &count)
	return
}

func (r *repository) GetTaxReport(ctx context.Context, periodID int64) ([]TaxReportItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT p.employee_id, e.employee_code, tp.name as employee_name, COALESCE(e.rfc, ''),
			p.gross_amount,
			COALESCE((SELECT SUM(amount) FROM payroll_detail WHERE payroll_id = p.id AND concept_code = 'ISR'), 0) as isr_amount,
			COALESCE((SELECT SUM(amount) FROM payroll_detail WHERE payroll_id = p.id AND concept_code = 'IMSS'), 0) as imss_amount
		 FROM payroll p
		 JOIN employee e ON e.id = p.employee_id
		 JOIN third_parties tp ON tp.id = e.third_party_id
		 WHERE p.payroll_period_id = $1
		 ORDER BY e.employee_code`, periodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var report []TaxReportItem
	for rows.Next() {
		var item TaxReportItem
		if err := rows.Scan(&item.EmployeeID, &item.EmployeeCode, &item.EmployeeName, &item.RFC, &item.GrossAmount, &item.ISRAmount, &item.IMSSAmount); err != nil {
			return nil, err
		}
		item.TotalTax = item.ISRAmount + item.IMSSAmount
		report = append(report, item)
	}
	return report, nil
}

func (r *repository) GetDeductionReport(ctx context.Context, periodID int64) ([]DeductionReportItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT pd.concept_code, pd.concept_name, SUM(pd.amount) as total_amount, COUNT(DISTINCT pd.payroll_id) as employee_count
		 FROM payroll_detail pd
		 JOIN payroll p ON p.id = pd.payroll_id
		 WHERE p.payroll_period_id = $1 AND pd.concept_type = 'DEDUCTION'
		 GROUP BY pd.concept_code, pd.concept_name
		 ORDER BY total_amount DESC`, periodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var report []DeductionReportItem
	for rows.Next() {
		var item DeductionReportItem
		if err := rows.Scan(&item.DeductionCode, &item.DeductionName, &item.TotalAmount, &item.EmployeeCount); err != nil {
			return nil, err
		}
		report = append(report, item)
	}
	return report, nil
}

func (r *repository) GetEmployeeEarnings(ctx context.Context, employeeID int64, startDate, endDate *time.Time) ([]Payroll, error) {
	query := `SELECT p.id, p.payroll_period_id, p.employee_id, p.salary_id, p.base_amount, p.worked_days, 
			p.gross_amount, p.total_additions, p.total_deductions, p.net_amount, p.status, p.notes, p.created_at, p.updated_at
		 FROM payroll p
		 JOIN payroll_period pp ON pp.id = p.payroll_period_id
		 WHERE p.employee_id = $1`
	args := []interface{}{employeeID}
	argIndex := 2

	if startDate != nil {
		query += fmt.Sprintf(" AND pp.start_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		query += fmt.Sprintf(" AND pp.end_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}
	query += " ORDER BY pp.start_date DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payrolls []Payroll
	for rows.Next() {
		var p Payroll
		var salaryID sql.NullInt64
		if err := rows.Scan(&p.ID, &p.PayrollPeriodID, &p.EmployeeID, &salaryID, &p.BaseAmount, &p.WorkedDays,
			&p.GrossAmount, &p.TotalAdditions, &p.TotalDeductions, &p.NetAmount, &p.Status, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if salaryID.Valid {
			p.SalaryID = &salaryID.Int64
		}
		payrolls = append(payrolls, p)
	}
	return payrolls, nil
}

// ─── Leave Management ──────────────────────────────────────────────────────────

func (r *repository) CreateLeaveType(ctx context.Context, lt *LeaveType) (int64, error) {
	lt.CreatedAt = time.Now()
	lt.IsActive = true

	return lt.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO leave_type (code, name, is_paid, days_allowed, requires_approval, is_active, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		lt.Code, lt.Name, lt.IsPaid, lt.DaysAllowed, lt.RequiresApproval, lt.IsActive, lt.CreatedAt,
	).Scan(&lt.ID)
}

func (r *repository) GetLeaveTypeByID(ctx context.Context, id int64) (*LeaveType, error) {
	var lt LeaveType

	err := r.db.QueryRowContext(ctx,
		`SELECT id, code, name, is_paid, days_allowed, requires_approval, is_active, created_at
		 FROM leave_type WHERE id = $1`, id,
	).Scan(&lt.ID, &lt.Code, &lt.Name, &lt.IsPaid, &lt.DaysAllowed, &lt.RequiresApproval, &lt.IsActive, &lt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &lt, nil
}

func (r *repository) ListLeaveTypes(ctx context.Context, activeOnly bool) ([]LeaveType, error) {
	query := `SELECT id, code, name, is_paid, days_allowed, requires_approval, is_active, created_at FROM leave_type`
	if activeOnly {
		query += " WHERE is_active = TRUE"
	}
	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []LeaveType
	for rows.Next() {
		var lt LeaveType
		if err := rows.Scan(&lt.ID, &lt.Code, &lt.Name, &lt.IsPaid, &lt.DaysAllowed, &lt.RequiresApproval, &lt.IsActive, &lt.CreatedAt); err != nil {
			return nil, err
		}
		types = append(types, lt)
	}
	return types, nil
}

func (r *repository) CreateLeave(ctx context.Context, el *EmployeeLeave) (int64, error) {
	el.CreatedAt = time.Now()
	el.Status = "PENDING"

	return el.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO employee_leave (employee_id, leave_type_id, start_date, end_date, days, status, reason, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		el.EmployeeID, el.LeaveTypeID, el.StartDate, el.EndDate, el.Days, el.Status, el.Reason, el.CreatedAt,
	).Scan(&el.ID)
}

func (r *repository) GetLeaveByID(ctx context.Context, id int64) (*EmployeeLeave, error) {
	var el EmployeeLeave
	var approvedAt, updatedAt sql.NullTime
	var approvedBy sql.NullInt64

	err := r.db.QueryRowContext(ctx,
		`SELECT id, employee_id, leave_type_id, start_date, end_date, days, status, reason, approved_by, approved_at, created_at, updated_at
		 FROM employee_leave WHERE id = $1`, id,
	).Scan(&el.ID, &el.EmployeeID, &el.LeaveTypeID, &el.StartDate, &el.EndDate, &el.Days, &el.Status, &el.Reason, &approvedBy, &approvedAt, &el.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if approvedBy.Valid {
		el.ApprovedBy = &approvedBy.Int64
	}
	if approvedAt.Valid {
		el.ApprovedAt = &approvedAt.Time
	}
	if updatedAt.Valid {
		el.UpdatedAt = &updatedAt.Time
	}
	return &el, nil
}

func (r *repository) UpdateLeave(ctx context.Context, el *EmployeeLeave) error {
	now := time.Now()
	el.UpdatedAt = &now
	_, err := r.db.ExecContext(ctx,
		`UPDATE employee_leave SET status=$1, approved_by=$2, approved_at=$3, updated_at=$4 WHERE id=$5`,
		el.Status, el.ApprovedBy, el.ApprovedAt, now, el.ID)
	return err
}

func (r *repository) ListLeaves(ctx context.Context, employeeID *int64, status string, startDate, endDate *time.Time) ([]EmployeeLeave, error) {
	query := `SELECT id, employee_id, leave_type_id, start_date, end_date, days, status, reason, approved_by, approved_at, created_at, updated_at FROM employee_leave WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if employeeID != nil {
		query += fmt.Sprintf(" AND employee_id = $%d", argIndex)
		args = append(args, *employeeID)
		argIndex++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if startDate != nil {
		query += fmt.Sprintf(" AND start_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		query += fmt.Sprintf(" AND end_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}
	query += " ORDER BY start_date DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaves []EmployeeLeave
	for rows.Next() {
		var el EmployeeLeave
		var approvedAt, updatedAt sql.NullTime
		var approvedBy sql.NullInt64
		if err := rows.Scan(&el.ID, &el.EmployeeID, &el.LeaveTypeID, &el.StartDate, &el.EndDate, &el.Days, &el.Status, &el.Reason, &approvedBy, &approvedAt, &el.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		if approvedBy.Valid {
			el.ApprovedBy = &approvedBy.Int64
		}
		if approvedAt.Valid {
			el.ApprovedAt = &approvedAt.Time
		}
		if updatedAt.Valid {
			el.UpdatedAt = &updatedAt.Time
		}
		leaves = append(leaves, el)
	}
	return leaves, nil
}

func (r *repository) GetLeaveBalance(ctx context.Context, employeeID int64, year int) ([]LeaveBalance, error) {
	startOfYear := fmt.Sprintf("%d-01-01", year)
	endOfYear := fmt.Sprintf("%d-12-31", year)

	rows, err := r.db.QueryContext(ctx,
		`SELECT lt.id, lt.name, lt.days_allowed,
			COALESCE((SELECT SUM(el.days) FROM employee_leave el 
				WHERE el.leave_type_id = lt.id AND el.employee_id = $1 
				AND el.status = 'APPROVED' 
				AND el.start_date >= $2 AND el.end_date <= $3), 0) as days_used
		 FROM leave_type lt
		 WHERE lt.is_active = TRUE
		 ORDER BY lt.name`, employeeID, startOfYear, endOfYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []LeaveBalance
	for rows.Next() {
		var b LeaveBalance
		if err := rows.Scan(&b.LeaveTypeID, &b.LeaveTypeName, &b.DaysAllowed, &b.DaysUsed); err != nil {
			return nil, err
		}
		b.DaysRemaining = float64(b.DaysAllowed) - b.DaysUsed
		balances = append(balances, b)
	}
	return balances, nil
}

// ─── Commissions Import ──────────────────────────────────────────────────────

func (r *repository) GetApprovedCommissions(ctx context.Context, startDate, endDate *time.Time) ([]interface{}, error) {
	// This would connect to the commissions module
	// For now, return empty slice
	return []interface{}{}, nil
}
