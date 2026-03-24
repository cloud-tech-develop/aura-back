-- Drop foreign key constraints first
ALTER TABLE commission DROP CONSTRAINT IF EXISTS fk_commission_rule;
ALTER TABLE commission DROP CONSTRAINT IF EXISTS fk_commission_branch;
ALTER TABLE commission DROP CONSTRAINT IF EXISTS fk_commission_employee;
ALTER TABLE commission DROP CONSTRAINT IF EXISTS fk_commission_order;
ALTER TABLE commission_rule DROP CONSTRAINT IF EXISTS fk_cr_category;
ALTER TABLE commission_rule DROP CONSTRAINT IF EXISTS fk_cr_product;
ALTER TABLE commission_rule DROP CONSTRAINT IF EXISTS fk_cr_employee;

-- Drop tables
DROP TABLE IF EXISTS commission;
DROP TABLE IF EXISTS commission_rule;
