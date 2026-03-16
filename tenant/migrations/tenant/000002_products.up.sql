-- Tabla de Categorías
CREATE TABLE category (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id BIGINT REFERENCES category(id),
    empresa_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT category_empresa_fk FOREIGN KEY (empresa_id) REFERENCES public.enterprises(id),
    CONSTRAINT category_parent_fk FOREIGN KEY (parent_id) REFERENCES category(id),
    CONSTRAINT category_name_unique UNIQUE (empresa_id, name)
);

CREATE INDEX idx_category_empresa ON category(empresa_id);
CREATE INDEX idx_category_parent ON category(parent_id);
CREATE INDEX idx_category_deleted_at ON category(deleted_at) WHERE deleted_at IS NULL;

-- Tabla de Marcas
CREATE TABLE brand (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    empresa_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT brand_empresa_fk FOREIGN KEY (empresa_id) REFERENCES public.enterprises(id),
    CONSTRAINT brand_name_unique UNIQUE (empresa_id, name)
);

CREATE INDEX idx_brand_empresa ON brand(empresa_id);
CREATE INDEX idx_brand_deleted_at ON brand(deleted_at) WHERE deleted_at IS NULL;

-- Tabla de Productos
CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    category_id BIGINT REFERENCES category(id),
    brand_id BIGINT REFERENCES brand(id),
    cost_price DECIMAL(12,2) NOT NULL,
    sale_price DECIMAL(12,2) NOT NULL,
    tax_rate DECIMAL(5,2) NOT NULL DEFAULT 19.00,
    min_stock INTEGER DEFAULT 0,
    current_stock INTEGER DEFAULT 0,
    image_url VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'DISCONTINUED')),
    empresa_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT product_empresa_fk FOREIGN KEY (empresa_id) REFERENCES public.enterprises(id),
    CONSTRAINT product_category_fk FOREIGN KEY (category_id) REFERENCES category(id),
    CONSTRAINT product_brand_fk FOREIGN KEY (brand_id) REFERENCES brand(id),
    CONSTRAINT product_sku_unique UNIQUE (empresa_id, sku),
    CONSTRAINT product_price_check CHECK (sale_price >= cost_price)
);

CREATE INDEX idx_product_empresa ON product(empresa_id);
CREATE INDEX idx_product_category ON product(category_id);
CREATE INDEX idx_product_brand ON product(brand_id);
CREATE INDEX idx_product_sku ON product(sku);
CREATE INDEX idx_product_status ON product(status);
CREATE INDEX idx_product_deleted_at ON product(deleted_at) WHERE deleted_at IS NULL;

-- Trigger para updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_category_updated_at BEFORE UPDATE ON category
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_brand_updated_at BEFORE UPDATE ON brand
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_product_updated_at BEFORE UPDATE ON product
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
