-- Remove active field from brand table
ALTER TABLE brand DROP COLUMN IF EXISTS active;
DROP INDEX IF EXISTS idx_brand_active;
