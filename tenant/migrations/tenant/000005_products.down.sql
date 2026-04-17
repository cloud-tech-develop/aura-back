-- Drop product trigger
DROP TRIGGER IF EXISTS update_product_updated_at ON product;

-- Drop table
DROP TABLE IF EXISTS product CASCADE;