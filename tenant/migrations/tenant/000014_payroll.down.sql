-- Drop foreign key constraints first
ALTER TABLE payroll_payment DROP CONSTRAINT IF EXISTS fk_pp_employee;
ALTER TABLE payroll_payment DROP CONSTRAINT IF EXISTS fk_pp_payroll;
ALTER TABLE bonus DROP CONSTRAINT IF EXISTS fk_bonus_employee;
ALTER TABLE overtime DROP CONSTRAINT IF EXISTS fk_overtime_employee;
ALTER TABLE employee_loan DROP CONSTRAINT IF EXISTS fk_el_employee;
ALTER TABLE payroll_detail DROP CONSTRAINT IF EXISTS fk_pd_payroll;
ALTER TABLE payroll DROP CONSTRAINT IF EXISTS fk_payroll_salary;
ALTER TABLE payroll DROP CONSTRAINT IF EXISTS fk_payroll_employee;
ALTER TABLE payroll DROP CONSTRAINT IF EXISTS fk_payroll_period;
ALTER TABLE payroll_period DROP CONSTRAINT IF EXISTS fk_pp_approved_by;
ALTER TABLE salary DROP CONSTRAINT IF EXISTS fk_salary_employee;
ALTER TABLE employee DROP CONSTRAINT IF EXISTS fk_employee_third_party;

-- Drop tables
DROP TABLE IF EXISTS payroll_payment;
DROP TABLE IF EXISTS bonus;
DROP TABLE IF EXISTS overtime;
DROP TABLE IF EXISTS employee_loan;
DROP TABLE IF EXISTS payroll_detail;
DROP TABLE IF EXISTS payroll;
DROP TABLE IF EXISTS payroll_period;
DROP TABLE IF EXISTS addition_type;
DROP TABLE IF EXISTS deduction_type;
DROP TABLE IF EXISTS salary;
DROP TABLE IF EXISTS employee;
