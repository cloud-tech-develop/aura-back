-- Drop category triggers
DROP TRIGGER IF EXISTS update_category_updated_at ON category;
DROP TRIGGER IF EXISTS update_brand_updated_at ON brand;

-- Drop tables
DROP TABLE IF EXISTS category CASCADE;
DROP TABLE IF EXISTS brand CASCADE;