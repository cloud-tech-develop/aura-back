-- Remove sync fields from core tables

DROP INDEX IF EXISTS idx_invoice_prefix_global_id;
ALTER TABLE invoice_prefix DROP COLUMN IF EXISTS last_synced_at;
ALTER TABLE invoice_prefix DROP COLUMN IF EXISTS sync_status;
ALTER TABLE invoice_prefix DROP COLUMN IF EXISTS global_id;

DROP INDEX IF EXISTS idx_invoice_global_id;
ALTER TABLE invoice DROP COLUMN IF EXISTS last_synced_at;
ALTER TABLE invoice DROP COLUMN IF EXISTS sync_status;
ALTER TABLE invoice DROP COLUMN IF EXISTS global_id;

DROP INDEX IF EXISTS idx_third_parties_global_id;
ALTER TABLE third_parties DROP COLUMN IF EXISTS last_synced_at;
ALTER TABLE third_parties DROP COLUMN IF EXISTS sync_status;
ALTER TABLE third_parties DROP COLUMN IF EXISTS global_id;

DROP INDEX IF EXISTS idx_sales_order_global_id;
ALTER TABLE sales_order DROP COLUMN IF EXISTS last_synced_at;
ALTER TABLE sales_order DROP COLUMN IF EXISTS sync_status;
ALTER TABLE sales_order DROP COLUMN IF EXISTS global_id;

DROP INDEX IF EXISTS idx_product_global_id;
ALTER TABLE product DROP COLUMN IF EXISTS last_synced_at;
ALTER TABLE product DROP COLUMN IF EXISTS sync_status;
ALTER TABLE product DROP COLUMN IF EXISTS global_id;

DROP INDEX IF EXISTS idx_brand_global_id;
ALTER TABLE brand DROP COLUMN IF EXISTS last_synced_at;
ALTER TABLE brand DROP COLUMN IF EXISTS sync_status;
ALTER TABLE brand DROP COLUMN IF EXISTS global_id;

DROP INDEX IF EXISTS idx_category_global_id;
ALTER TABLE category DROP COLUMN IF EXISTS last_synced_at;
ALTER TABLE category DROP COLUMN IF EXISTS sync_status;
ALTER TABLE category DROP COLUMN IF EXISTS global_id;
