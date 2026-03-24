-- Drop foreign key constraints first
ALTER TABLE cash_movement DROP CONSTRAINT IF EXISTS fk_cash_movement_user;
ALTER TABLE cash_movement DROP CONSTRAINT IF EXISTS fk_cash_movement_shift;
ALTER TABLE cash_shift DROP CONSTRAINT IF EXISTS fk_cash_shift_branch;
ALTER TABLE cash_shift DROP CONSTRAINT IF EXISTS fk_cash_shift_user;
ALTER TABLE cash_shift DROP CONSTRAINT IF EXISTS fk_cash_shift_drawer;
ALTER TABLE cash_drawer DROP CONSTRAINT IF EXISTS fk_cash_drawer_branch;

-- Drop tables
DROP TABLE IF EXISTS cash_movement;
DROP TABLE IF EXISTS cash_shift;
DROP TABLE IF EXISTS cash_drawer;
