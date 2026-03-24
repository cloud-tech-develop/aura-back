-- Create employee table
CREATE TABLE IF NOT EXISTS employee (
    id BIGSERIAL PRIMARY KEY,
    employee_code VARCHAR(20) NOT NULL UNIQUE,
    third_party_id BIGINT NOT NULL,
    user_id BIGINT,
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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_employee_code ON employee(employee_code);
CREATE INDEX idx_employee_third_party ON employee(third_party_id);
CREATE INDEX idx_employee_status ON employee(status);
CREATE INDEX idx_employee_department ON employee(department);

-- Create salary table
CREATE TABLE IF NOT EXISTS salary (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL,
    base_salary DECIMAL(12,2) NOT NULL,
    daily_salary DECIMAL(12,2) NOT NULL,
    salary_type VARCHAR(20) NOT NULL DEFAULT 'MONTHLY' CHECK (salary_type IN ('MONTHLY', 'BIWEEKLY', 'WEEKLY', 'HOURLY')),
    effective_date DATE NOT NULL,
    end_date DATE,
    is_current BOOLEAN DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_salary_employee ON salary(employee_id);
CREATE INDEX idx_salary_current ON salary(is_current) WHERE is_current = TRUE;

-- Create deduction_type table
CREATE TABLE IF NOT EXISTS deduction_type (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_tax BOOLEAN DEFAULT FALSE,
    is_mandatory BOOLEAN DEFAULT FALSE,
    calculation_type VARCHAR(20) NOT NULL DEFAULT 'PERCENTAGE' CHECK (calculation_type IN ('PERCENTAGE', 'FIXED', 'TABLE')),
    default_value DECIMAL(12,2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create addition_type table
CREATE TABLE IF NOT EXISTS addition_type (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_taxable BOOLEAN DEFAULT TRUE,
    calculation_type VARCHAR(20) NOT NULL DEFAULT 'FIXED' CHECK (calculation_type IN ('HOURLY', 'PERCENTAGE', 'FIXED')),
    default_value DECIMAL(12,2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create payroll_period table
CREATE TABLE IF NOT EXISTS payroll_period (
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
    approved_by BIGINT,
    approved_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    
    CONSTRAINT payroll_period_dates_check CHECK (end_date > start_date)
);

CREATE INDEX idx_payroll_period_dates ON payroll_period(start_date, end_date);
CREATE INDEX idx_payroll_period_status ON payroll_period(status);

-- Create payroll table
CREATE TABLE IF NOT EXISTS payroll (
    id BIGSERIAL PRIMARY KEY,
    payroll_period_id BIGINT NOT NULL,
    employee_id BIGINT NOT NULL,
    salary_id BIGINT,
    base_amount DECIMAL(12,2) NOT NULL,
    worked_days DECIMAL(5,2) NOT NULL DEFAULT 15,
    gross_amount DECIMAL(12,2) NOT NULL,
    total_additions DECIMAL(12,2) DEFAULT 0,
    total_deductions DECIMAL(12,2) DEFAULT 0,
    net_amount DECIMAL(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'CALCULATED', 'APPROVED', 'PAID', 'CANCELLED')),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_payroll_period ON payroll(payroll_period_id);
CREATE INDEX idx_payroll_employee ON payroll(employee_id);
CREATE INDEX idx_payroll_status ON payroll(status);

-- Create payroll_detail table
CREATE TABLE IF NOT EXISTS payroll_detail (
    id BIGSERIAL PRIMARY KEY,
    payroll_id BIGINT NOT NULL,
    concept_type VARCHAR(20) NOT NULL CHECK (concept_type IN ('ADDITION', 'DEDUCTION')),
    concept_code VARCHAR(20) NOT NULL,
    concept_name VARCHAR(100) NOT NULL,
    quantity DECIMAL(10,2) DEFAULT 1,
    rate DECIMAL(12,2) DEFAULT 0,
    amount DECIMAL(12,2) NOT NULL,
    reference_id BIGINT,
    reference_type VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payroll_detail_payroll ON payroll_detail(payroll_id);

-- Create employee_loan table
CREATE TABLE IF NOT EXISTS employee_loan (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL,
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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_employee_loan_employee ON employee_loan(employee_id);
CREATE INDEX idx_employee_loan_status ON employee_loan(status);

-- Create overtime table
CREATE TABLE IF NOT EXISTS overtime (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL,
    work_date DATE NOT NULL,
    hours DECIMAL(5,2) NOT NULL,
    rate_type VARCHAR(20) NOT NULL DEFAULT 'DOUBLE' CHECK (rate_type IN ('REGULAR', 'DOUBLE', 'TRIPLE')),
    hourly_rate DECIMAL(12,2) NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'PAID', 'CANCELLED')),
    approved_by BIGINT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_overtime_employee ON overtime(employee_id);
CREATE INDEX idx_overtime_status ON overtime(status);
CREATE INDEX idx_overtime_date ON overtime(work_date);

-- Create bonus table
CREATE TABLE IF NOT EXISTS bonus (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL,
    bonus_type VARCHAR(30) NOT NULL CHECK (bonus_type IN ('PERFORMANCE', 'HOLIDAY', 'TRANSPORT', 'FOOD', 'OTHER')),
    amount DECIMAL(12,2) NOT NULL,
    bonus_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'PAID', 'CANCELLED')),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bonus_employee ON bonus(employee_id);
CREATE INDEX idx_bonus_status ON bonus(status);

-- Create payroll_payment table
CREATE TABLE IF NOT EXISTS payroll_payment (
    id BIGSERIAL PRIMARY KEY,
    payroll_id BIGINT NOT NULL,
    employee_id BIGINT NOT NULL,
    payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('BANK_TRANSFER', 'CASH', 'CHECK')),
    amount DECIMAL(12,2) NOT NULL,
    payment_date DATE NOT NULL,
    reference_number VARCHAR(50),
    bank_reference VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'COMPLETED', 'CANCELLED')),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payroll_payment_payroll ON payroll_payment(payroll_id);
CREATE INDEX idx_payroll_payment_employee ON payroll_payment(employee_id);
CREATE INDEX idx_payroll_payment_date ON payroll_payment(payment_date);

-- Add foreign key constraints
ALTER TABLE employee ADD CONSTRAINT fk_employee_third_party FOREIGN KEY (third_party_id) REFERENCES third_parties(id);
ALTER TABLE salary ADD CONSTRAINT fk_salary_employee FOREIGN KEY (employee_id) REFERENCES employee(id);
ALTER TABLE payroll_period ADD CONSTRAINT fk_pp_approved_by FOREIGN KEY (approved_by) REFERENCES users(id);
ALTER TABLE payroll ADD CONSTRAINT fk_payroll_period FOREIGN KEY (payroll_period_id) REFERENCES payroll_period(id);
ALTER TABLE payroll ADD CONSTRAINT fk_payroll_employee FOREIGN KEY (employee_id) REFERENCES employee(id);
ALTER TABLE payroll ADD CONSTRAINT fk_payroll_salary FOREIGN KEY (salary_id) REFERENCES salary(id);
ALTER TABLE payroll_detail ADD CONSTRAINT fk_pd_payroll FOREIGN KEY (payroll_id) REFERENCES payroll(id) ON DELETE CASCADE;
ALTER TABLE employee_loan ADD CONSTRAINT fk_el_employee FOREIGN KEY (employee_id) REFERENCES employee(id);
ALTER TABLE overtime ADD CONSTRAINT fk_overtime_employee FOREIGN KEY (employee_id) REFERENCES employee(id);
ALTER TABLE bonus ADD CONSTRAINT fk_bonus_employee FOREIGN KEY (employee_id) REFERENCES employee(id);
ALTER TABLE payroll_payment ADD CONSTRAINT fk_pp_payroll FOREIGN KEY (payroll_id) REFERENCES payroll(id);
ALTER TABLE payroll_payment ADD CONSTRAINT fk_pp_employee FOREIGN KEY (employee_id) REFERENCES employee(id);

-- Insert default deduction types
INSERT INTO deduction_type (code, name, description, is_tax, is_mandatory, calculation_type, default_value) VALUES
('ISR', 'ISR', 'Impuesto Sobre la Renta', TRUE, TRUE, 'TABLE', NULL),
('IMSS', 'IMSS', 'Seguro Social', TRUE, TRUE, 'PERCENTAGE', 2.75),
('INFONAVIT', 'INFONAVIT', 'Crédito Infonavit', FALSE, FALSE, 'PERCENTAGE', 5.0),
('PENSION', 'Pensión Alimenticia', 'Pensión alimenticia', FALSE, FALSE, 'PERCENTAGE', NULL),
('LOAN', 'Préstamo', 'Descuento de préstamo', FALSE, FALSE, 'FIXED', NULL),
('ADVANCE', 'Anticipo', 'Descuento de anticipo', FALSE, FALSE, 'FIXED', NULL),
('OTHER', 'Otro', 'Otras deducciones', FALSE, FALSE, 'FIXED', NULL);

-- Insert default addition types
INSERT INTO addition_type (code, name, description, is_taxable, calculation_type, default_value) VALUES
('OVERTIME', 'Horas Extra', 'Pago de horas extra', TRUE, 'HOURLY', NULL),
('BONUS', 'Bono', 'Bono de desempeño', TRUE, 'FIXED', NULL),
('COMMISSION', 'Comisión', 'Comisión por ventas', TRUE, 'FIXED', NULL),
('TRANSPORT', 'Transporte', 'Ayuda de transporte', FALSE, 'FIXED', 500.00),
('FOOD', 'Alimentación', 'Vale de alimentación', FALSE, 'FIXED', 600.00),
('HOLIDAY', 'Prima Vacacional', 'Prima vacacional', TRUE, 'PERCENTAGE', 25.0),
('CHRISTMAS', 'Aguinaldo', 'Gratificación navideña', TRUE, 'FIXED', NULL),
('OTHER', 'Otro', 'Otras percepciones', TRUE, 'FIXED', NULL);
