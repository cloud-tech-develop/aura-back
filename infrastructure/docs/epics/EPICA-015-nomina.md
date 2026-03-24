# EPICA-015: Payroll Management

## 📌 General Information

- ID: EPICA-015
- State: Completed
- Priority: High
- Start Date: 2026-03-24
- Target Date: 2026-09-30
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Manage the complete payroll lifecycle for employees, including employee profiles, salary configuration, payroll calculation, deductions, payments, and reporting. This module ensures accurate and timely compensation while maintaining compliance with labor regulations.

**What problem does it solve?**

- Manual payroll calculation errors
- No centralized employee records
- Inconsistent salary payments
- Lack of payroll history
- Manual tax and deduction calculations
- No integration with commissions and attendance

**What value does it generate?**

- Automated payroll processing
- Accurate tax and deduction calculations
- Complete payroll audit trail
- Employee self-service payslips
- Integration with commissions and attendance
- Compliance with labor regulations

---

## 👥 Stakeholders

- End User: HR managers, accountants, store managers, employees
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Payroll Management module handles:

1. **Employee Management**: Registration, profiles, employment details
2. **Salary Configuration**: Base salary, allowances, benefits
3. **Payroll Periods**: Bi-weekly/monthly period management
4. **Deductions & Additions**: Taxes, insurance, loans, overtime, bonuses
5. **Payroll Processing**: Calculation, approval, payment
6. **Reporting**: Payslips, summaries, tax reports

---

## 📦 Scope

### Included:

- Employee registration and management
- Employment type configuration
- Salary structure management
- Payroll period management
- Automatic payroll calculation
- Tax calculations (ISR, IMSS, etc.)
- Deductions (loans, advances, insurance)
- Additions (overtime, bonuses, commissions)
- Payroll approval workflow
- Payment processing
- Payslip generation
- Payroll reports

### Not Included:

- Full HR management (recruitment, performance)
- Time tracking/attendance (separate module)
- Benefits administration
- Loan origination (only deduction tracking)

---

## 🧩 User Stories

### Employee Management

| HU         | Title                     | Description                                                | Priority | State      |
| ---------- | ------------------------- | ---------------------------------------------------------- | -------- | ---------- |
| HU-PAY-001 | Register Employee         | Create employee profile with personal and employment data  | High     | ✅ Implemented |
| HU-PAY-002 | Update Employee Profile   | Modify employee information                                | High     | ✅ Implemented |
| HU-PAY-003 | Deactivate Employee       | Mark employee as inactive (termination)                    | Medium   | ✅ Implemented |
| HU-PAY-004 | View Employee List        | List all employees with filters                            | High     | ✅ Implemented |
| HU-PAY-005 | Configure Employment Type | Define employment types (full-time, part-time, contractor) | Medium   | ✅ Implemented |

### Salary Configuration

| HU         | Title                      | Description                                      | Priority | State      |
| ---------- | -------------------------- | ------------------------------------------------ | -------- | ---------- |
| HU-PAY-006 | Configure Salary Structure | Define salary components (base, allowances)      | High     | ✅ Implemented |
| HU-PAY-007 | Set Employee Salary        | Assign salary to employee                        | High     | ✅ Implemented |
| HU-PAY-008 | Configure Deductions       | Define deduction types (taxes, insurance, loans) | High     | ✅ Implemented |
| HU-PAY-009 | Configure Additions        | Define addition types (overtime, bonuses)        | High     | ✅ Implemented |
| HU-PAY-010 | Manage Employee Loans      | Track employee loans and advances                | Medium   | ✅ Implemented |

### Payroll Periods

| HU         | Title                 | Description                                        | Priority | State      |
| ---------- | --------------------- | -------------------------------------------------- | -------- | ---------- |
| HU-PAY-011 | Create Payroll Period | Define bi-weekly/monthly periods                   | High     | ✅ Implemented |
| HU-PAY-012 | Close Payroll Period  | Lock period for processing                         | High     | ✅ Implemented |
| HU-PAY-013 | Reopen Payroll Period | Reopen period for corrections (with authorization) | Low      | ✅ Implemented |

### Payroll Processing

| HU         | Title              | Description                          | Priority | State      |
| ---------- | ------------------ | ------------------------------------ | -------- | ---------- |
| HU-PAY-014 | Calculate Payroll  | Process payroll for period           | High     | ✅ Implemented |
| HU-PAY-015 | Preview Payroll    | View payroll before approval         | High     | ✅ Implemented |
| HU-PAY-016 | Approve Payroll    | Authorize payroll for payment        | High     | ✅ Implemented |
| HU-PAY-017 | Reject Payroll     | Reject payroll with reasons          | Medium   | ✅ Implemented |
| HU-PAY-018 | Register Overtime  | Record employee overtime hours       | High     | ✅ Implemented |
| HU-PAY-019 | Register Bonus     | Add one-time bonuses to payroll      | Medium   | ✅ Implemented |
| HU-PAY-020 | Import Commissions | Import commissions from sales module | Medium   | ✅ Implemented |

### Payments

| HU         | Title                 | Description                                 | Priority | State      |
| ---------- | --------------------- | ------------------------------------------- | -------- | ---------- |
| HU-PAY-021 | Process Payment       | Execute payroll payment                     | High     | ✅ Implemented |
| HU-PAY-022 | Generate Bank File    | Create bank transfer file for bulk payments | Medium   | ✅ Implemented |
| HU-PAY-023 | Record Manual Payment | Register cash/check payments                | Medium   | ✅ Implemented |
| HU-PAY-024 | View Payment History  | List all payroll payments                   | High     | ✅ Implemented |

### Reporting

| HU         | Title                    | Description                      | Priority | State      |
| ---------- | ------------------------ | -------------------------------- | -------- | ---------- |
| HU-PAY-025 | Generate Payslip         | Create employee payslip PDF      | High     | ✅ Implemented |
| HU-PAY-026 | Payroll Summary Report   | View payroll summary by period   | High     | ✅ Implemented |
| HU-PAY-027 | Tax Report               | Generate tax deduction report    | High     | ✅ Implemented |
| HU-PAY-028 | Employee Earnings Report | Show employee earnings history   | Medium   | ✅ Implemented |
| HU-PAY-029 | Deduction Report         | Report of all deductions by type | Medium   | ✅ Implemented |

### Leave & Attendance Integration

| HU         | Title                      | Description                       | Priority | State      |
| ---------- | -------------------------- | --------------------------------- | -------- | ---------- |
| HU-PAY-030 | Configure Leave Types      | Define vacation, sick leave, etc. | Medium   | ✅ Implemented |
| HU-PAY-031 | Register Leave Request     | Employee requests time off        | Medium   | ✅ Implemented |
| HU-PAY-032 | Approve Leave              | Manager approves/rejects leave    | Medium   | ✅ Implemented |
| HU-PAY-033 | Calculate Leave Deductions | Deduct unpaid leave from payroll  | Medium   | ✅ Implemented |
| HU-PAY-034 | View Leave Balance         | Check employee vacation balance   | Medium   | ✅ Implemented |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

### Employee Rules

- Each employee has a unique employee code
- Employees can have multiple salary records (salary history)
- Employee deactivation requires termination date and reason
- Terminated employees cannot be included in new payroll periods

### Payroll Rules

- Payroll periods cannot overlap
- Payroll calculation requires period to be open
- Approved payroll cannot be modified without reopening
- All deductions and additions must be justified
- Commissions are imported from completed sales

### Payment Rules

- Payment can only be made for approved payroll
- Payment amount must match calculated net salary
- Bank transfers require valid bank account
- Cash payments require receipt confirmation

### Tax Rules

- Tax calculations follow current tax tables
- Social security contributions are calculated automatically
- Tax exemptions are applied per employee configuration
- Year-end tax adjustments are calculated automatically

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Database Schema (Tenant Schema)

**Table: employee**

```sql
CREATE TABLE employee (
    id BIGSERIAL PRIMARY KEY,
    employee_code VARCHAR(20) NOT NULL UNIQUE,
    third_party_id BIGINT NOT NULL REFERENCES third_parties(id),
    user_id BIGINT REFERENCES public.users(id),
    employment_type VARCHAR(20) NOT NULL CHECK (employment_type IN ('FULL_TIME', 'PART_TIME', 'CONTRACTOR', 'TEMPORARY')),
    position VARCHAR(100),
    department VARCHAR(100),
    hire_date DATE NOT NULL,
    termination_date DATE,
    termination_reason TEXT,
    bank_name VARCHAR(100),
    bank_account VARCHAR(50),
    bank_clabe VARCHAR(18),
    curp VARCHAR(18),
    rfc VARCHAR(13),
    nss VARCHAR(11),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'TERMINATED')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);
```

**Table: salary**

```sql
CREATE TABLE salary (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee(id),
    base_salary DECIMAL(12,2) NOT NULL,
    daily_salary DECIMAL(12,2) NOT NULL,
    salary_type VARCHAR(20) NOT NULL DEFAULT 'MONTHLY' CHECK (salary_type IN ('MONTHLY', 'BIWEEKLY', 'WEEKLY', 'HOURLY')),
    effective_date DATE NOT NULL,
    end_date DATE,
    is_current BOOLEAN DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Table: payroll_period**

```sql
CREATE TABLE payroll_period (
    id BIGSERIAL PRIMARY KEY,
    period_type VARCHAR(20) NOT NULL CHECK (period_type IN ('BIWEEKLY', 'MONTHLY', 'WEEKLY')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    payment_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'PROCESSING', 'APPROVED', 'PAID', 'CLOSED')),
    total_gross DECIMAL(14,2) DEFAULT 0,
    total_deductions DECIMAL(14,2) DEFAULT 0,
    total_net DECIMAL(14,2) DEFAULT 0,
    employee_count INTEGER DEFAULT 0,
    approved_by BIGINT REFERENCES public.users(id),
    approved_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,

    CONSTRAINT payroll_period_dates_check CHECK (end_date > start_date)
);
```

**Table: payroll**

```sql
CREATE TABLE payroll (
    id BIGSERIAL PRIMARY KEY,
    payroll_period_id BIGINT NOT NULL REFERENCES payroll_period(id),
    employee_id BIGINT NOT NULL REFERENCES employee(id),
    salary_id BIGINT REFERENCES salary(id),
    base_amount DECIMAL(12,2) NOT NULL,
    worked_days DECIMAL(5,2) NOT NULL DEFAULT 15,
    gross_amount DECIMAL(12,2) NOT NULL,
    total_additions DECIMAL(12,2) DEFAULT 0,
    total_deductions DECIMAL(12,2) DEFAULT 0,
    net_amount DECIMAL(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'CALCULATED', 'APPROVED', 'PAID', 'CANCELLED')),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);
```

**Table: payroll_detail**

```sql
CREATE TABLE payroll_detail (
    id BIGSERIAL PRIMARY KEY,
    payroll_id BIGINT NOT NULL REFERENCES payroll(id) ON DELETE CASCADE,
    concept_type VARCHAR(20) NOT NULL CHECK (concept_type IN ('ADDITION', 'DEDUCTION')),
    concept_code VARCHAR(20) NOT NULL,
    concept_name VARCHAR(100) NOT NULL,
    quantity DECIMAL(10,2) DEFAULT 1,
    rate DECIMAL(12,2) DEFAULT 0,
    amount DECIMAL(12,2) NOT NULL,
    reference_id BIGINT,
    reference_type VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Table: deduction_type**

```sql
CREATE TABLE deduction_type (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_tax BOOLEAN DEFAULT FALSE,
    is_mandatory BOOLEAN DEFAULT FALSE,
    calculation_type VARCHAR(20) NOT NULL DEFAULT 'PERCENTAGE' CHECK (calculation_type IN ('PERCENTAGE', 'FIXED', 'TABLE')),
    default_value DECIMAL(12,2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Table: addition_type**

```sql
CREATE TABLE addition_type (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_taxable BOOLEAN DEFAULT TRUE,
    calculation_type VARCHAR(20) NOT NULL DEFAULT 'FIXED' CHECK (calculation_type IN ('HOURLY', 'PERCENTAGE', 'FIXED')),
    default_value DECIMAL(12,2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Table: employee_loan**

```sql
CREATE TABLE employee_loan (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee(id),
    loan_type VARCHAR(20) NOT NULL CHECK (loan_type IN ('LOAN', 'ADVANCE', 'PURCHASE')),
    total_amount DECIMAL(12,2) NOT NULL,
    remaining_amount DECIMAL(12,2) NOT NULL,
    installment_amount DECIMAL(12,2) NOT NULL,
    installments_total INTEGER NOT NULL,
    installments_paid INTEGER DEFAULT 0,
    start_date DATE NOT NULL,
    end_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'COMPLETED', 'CANCELLED')),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);
```

**Table: leave_type**

```sql
CREATE TABLE leave_type (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    is_paid BOOLEAN DEFAULT TRUE,
    days_allowed INTEGER DEFAULT 0,
    requires_approval BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Table: employee_leave**

```sql
CREATE TABLE employee_leave (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee(id),
    leave_type_id BIGINT NOT NULL REFERENCES leave_type(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    days DECIMAL(5,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED', 'CANCELLED')),
    reason TEXT,
    approved_by BIGINT REFERENCES public.users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);
```

**Table: overtime**

```sql
CREATE TABLE overtime (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee(id),
    work_date DATE NOT NULL,
    hours DECIMAL(5,2) NOT NULL,
    rate_type VARCHAR(20) NOT NULL DEFAULT 'DOUBLE' CHECK (rate_type IN ('REGULAR', 'DOUBLE', 'TRIPLE')),
    hourly_rate DECIMAL(12,2) NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'PAID', 'CANCELLED')),
    approved_by BIGINT REFERENCES public.users(id),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Table: bonus**

```sql
CREATE TABLE bonus (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee(id),
    bonus_type VARCHAR(30) NOT NULL CHECK (bonus_type IN ('PERFORMANCE', 'HOLIDAY', 'TRANSPORT', 'FOOD', 'OTHER')),
    amount DECIMAL(12,2) NOT NULL,
    bonus_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'PAID', 'CANCELLED')),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Table: payroll_payment**

```sql
CREATE TABLE payroll_payment (
    id BIGSERIAL PRIMARY KEY,
    payroll_id BIGINT NOT NULL REFERENCES payroll(id),
    employee_id BIGINT NOT NULL REFERENCES employee(id),
    payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('BANK_TRANSFER', 'CASH', 'CHECK')),
    amount DECIMAL(12,2) NOT NULL,
    payment_date DATE NOT NULL,
    reference_number VARCHAR(50),
    bank_reference VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'COMPLETED', 'CANCELLED')),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

## 📊 Success Metrics

- Payroll processing time < 30 minutes for 100 employees
- Calculation accuracy 100%
- Payment error rate < 0.1%
- Employee satisfaction with payslips > 95%
- Tax compliance 100%

---

## 🚧 Risks

- Tax table changes require updates
- Bank file format variations
- Complex overtime calculations
- Multi-state tax variations
- Integration with external systems

---

## 📋 Implementation Priority

### Phase 1 (Critical)

1. HU-PAY-001: Register Employee
2. HU-PAY-006: Configure Salary Structure
3. HU-PAY-007: Set Employee Salary
4. HU-PAY-011: Create Payroll Period
5. HU-PAY-014: Calculate Payroll
6. HU-PAY-016: Approve Payroll
7. HU-PAY-025: Generate Payslip

### Phase 2 (Important)

1. HU-PAY-008: Configure Deductions
2. HU-PAY-009: Configure Additions
3. HU-PAY-018: Register Overtime
4. HU-PAY-020: Import Commissions
5. HU-PAY-021: Process Payment
6. HU-PAY-026: Payroll Summary Report

### Phase 3 (Enhancement)

1. HU-PAY-010: Manage Employee Loans
2. HU-PAY-022: Generate Bank File
3. HU-PAY-027: Tax Report
4. HU-PAY-030-034: Leave Management

### Phase 4 (Optional)

1. HU-PAY-013: Reopen Payroll Period
2. HU-PAY-019: Register Bonus
3. HU-PAY-028: Employee Earnings Report
4. HU-PAY-029: Deduction Report
