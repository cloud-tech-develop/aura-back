-- Trigger function for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- Table: category
CREATE TABLE IF NOT EXISTS category (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id BIGINT,
    default_tax_rate DECIMAL(5,2) DEFAULT 0.00,
    active BOOLEAN DEFAULT TRUE,
    enterprise_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT category_enterprise_fk FOREIGN KEY (enterprise_id) REFERENCES public.enterprises(id),
    CONSTRAINT category_parent_fk FOREIGN KEY (parent_id) REFERENCES category(id),
    CONSTRAINT category_name_unique UNIQUE (enterprise_id, name)
);

-- Indexes
CREATE INDEX idx_category_enterprise ON category(enterprise_id);
CREATE INDEX idx_category_parent ON category(parent_id);
CREATE INDEX idx_category_deleted_at ON category(deleted_at) WHERE deleted_at IS NULL;

-- Trigger for updated_at
CREATE TRIGGER update_category_updated_at BEFORE UPDATE ON category
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table: brand (with active field)
CREATE TABLE IF NOT EXISTS brand (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT TRUE,
    enterprise_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT brand_enterprise_fk FOREIGN KEY (enterprise_id) REFERENCES public.enterprises(id),
    CONSTRAINT brand_name_unique UNIQUE (enterprise_id, name)
);

-- Indexes
CREATE INDEX idx_brand_enterprise ON brand(enterprise_id);
CREATE INDEX idx_brand_active ON brand(active) WHERE active = TRUE;
CREATE INDEX idx_brand_deleted_at ON brand(deleted_at) WHERE deleted_at IS NULL;

-- Trigger for updated_at
CREATE TRIGGER update_brand_updated_at BEFORE UPDATE ON brand
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE category IS 'Product categories';
COMMENT ON COLUMN category.name IS 'Category name';
COMMENT ON COLUMN category.description IS 'Category description';
COMMENT ON COLUMN category.active IS 'Category active status';
COMMENT ON TABLE brand IS 'Product brands';
COMMENT ON COLUMN brand.name IS 'Brand name';
COMMENT ON COLUMN brand.description IS 'Brand description';
COMMENT ON COLUMN brand.active IS 'Brand active status';