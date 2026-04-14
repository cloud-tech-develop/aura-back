-- Table: product
CREATE TABLE IF NOT EXISTS product (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    category_id BIGINT,
    brand_id BIGINT,
    cost_price DECIMAL(12,2) NOT NULL,
    sale_price DECIMAL(12,2) NOT NULL,
    tax_rate DECIMAL(5,2) NOT NULL DEFAULT 19.00,
    min_stock INTEGER DEFAULT 0,
    current_stock INTEGER DEFAULT 0,
    image_url VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'DISCONTINUED')),
    enterprise_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT product_enterprise_fk FOREIGN KEY (enterprise_id) REFERENCES public.enterprises(id),
    CONSTRAINT product_category_fk FOREIGN KEY (category_id) REFERENCES category(id),
    CONSTRAINT product_brand_fk FOREIGN KEY (brand_id) REFERENCES brand(id),
    CONSTRAINT product_sku_unique UNIQUE (enterprise_id, sku),
    CONSTRAINT product_price_check CHECK (sale_price >= cost_price)
);

-- Indexes
CREATE INDEX idx_product_enterprise ON product(enterprise_id);
CREATE INDEX idx_product_category ON product(category_id);
CREATE INDEX idx_product_brand ON product(brand_id);
CREATE INDEX idx_product_sku ON product(sku);
CREATE INDEX idx_product_status ON product(status);
CREATE INDEX idx_product_deleted_at ON product(deleted_at) WHERE deleted_at IS NULL;

-- Trigger for updated_at
CREATE TRIGGER update_product_updated_at BEFORE UPDATE ON product
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON COLUMN product.sku IS 'Product SKU code';
COMMENT ON COLUMN product.name IS 'Product name';
COMMENT ON COLUMN product.description IS 'Product description';
COMMENT ON COLUMN product.category_id IS 'Category foreign key';
COMMENT ON COLUMN product.brand_id IS 'Brand foreign key';
COMMENT ON COLUMN product.cost_price IS 'Product cost price';
COMMENT ON COLUMN product.sale_price IS 'Product sale price';
COMMENT ON COLUMN product.tax_rate IS 'Product tax rate percentage';
COMMENT ON COLUMN product.min_stock IS 'Minimum stock threshold';
COMMENT ON COLUMN product.current_stock IS 'Current stock quantity';
COMMENT ON COLUMN product.image_url IS 'Product image URL';
COMMENT ON COLUMN product.status IS 'Product status: ACTIVE, INACTIVE, DISCONTINUED';