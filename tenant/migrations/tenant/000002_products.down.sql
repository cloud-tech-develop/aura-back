DROP TRIGGER IF EXISTS update_product_updated_at ON product;
DROP TRIGGER IF EXISTS update_brand_updated_at ON brand;
DROP TRIGGER IF EXISTS update_category_updated_at ON category;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS product;
DROP TABLE IF EXISTS brand;
DROP TABLE IF EXISTS category;
