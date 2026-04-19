-- Table: product
-- Product master table for catalog management
CREATE TABLE IF NOT EXISTS product (
    id BIGSERIAL PRIMARY KEY,
    
    -- Basic identification fields
    sku VARCHAR(50) NOT NULL,                  -- Product SKU code (unique per enterprise)
    barcode VARCHAR(50),                      -- Product barcode for scanning
    name VARCHAR(200) NOT NULL,                -- Product name
    description TEXT,                        -- Product description
    
    -- Reference fields (foreign keys)
    category_id BIGINT,                       -- Category foreign key reference
    brand_id BIGINT,                          -- Brand foreign key reference
    unit_id BIGINT NOT NULL,                 -- Base unit of measure foreign key
    
    -- Product type classification
    product_type VARCHAR(20) NOT NULL DEFAULT 'STANDARD' CHECK (product_type IN ('STANDARD', 'SERVICE', 'KIT', 'WEIGHTABLE')),
    
    -- Status and visibility
    active BOOLEAN NOT NULL DEFAULT true,      -- Product active status
    visible_in_pos BOOLEAN NOT NULL DEFAULT true,  -- Visibility in POS interface
    
    -- Pricing fields
    cost_price DECIMAL(12,2) NOT NULL DEFAULT 0,    -- Product cost price (purchase price)
    sale_price DECIMAL(12,2) NOT NULL DEFAULT 0,    -- Product sale price (retail price)
    price_2 DECIMAL(12,2),                     -- Alternative price level 2 (wholesale)
    price_3 DECIMAL(12,2),                     -- Alternative price level 3 (special)
    
    -- Tax configuration
    iva_percentage DECIMAL(5,2) NOT NULL DEFAULT 19.00,   -- IVA tax percentage
    consumption_tax_value DECIMAL(12,2) NOT NULL DEFAULT 0, -- Consumption tax value
    
    -- Inventory control settings
    current_stock INTEGER NOT NULL DEFAULT 0,        -- Current inventory quantity
    min_stock INTEGER NOT NULL DEFAULT 0,          -- Minimum stock threshold for alerts
    max_stock INTEGER NOT NULL DEFAULT 0,          -- Maximum stock level for inventory limits
    
    -- Inventory management options
    manages_inventory BOOLEAN NOT NULL DEFAULT true,     -- Enable inventory tracking
    manages_batches BOOLEAN NOT NULL DEFAULT false,    -- Enable batch/lot tracking
    manages_serial BOOLEAN NOT NULL DEFAULT false,       -- Enable serial number tracking
    allow_negative_stock BOOLEAN NOT NULL DEFAULT false,  -- Allow negative stock (stock outs)
    
    -- Media and multimedia
    image_url VARCHAR(500),                    -- Product main image URL
    
    -- Enterprise ownership
    enterprise_id BIGINT NOT NULL,              -- Enterprise foreign key
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT product_enterprise_fk FOREIGN KEY (enterprise_id) REFERENCES public.enterprises(id),
    CONSTRAINT product_category_fk FOREIGN KEY (category_id) REFERENCES category(id),
    CONSTRAINT product_brand_fk FOREIGN KEY (brand_id) REFERENCES brand(id),
    CONSTRAINT product_unit_fk FOREIGN KEY (unit_id) REFERENCES unit(id),
    CONSTRAINT product_sku_unique UNIQUE (enterprise_id, sku),
    CONSTRAINT product_barcode_unique UNIQUE (enterprise_id, barcode)
);

-- Indexes for performance optimization
CREATE INDEX idx_product_enterprise ON product(enterprise_id);
CREATE INDEX idx_product_category ON product(category_id);
CREATE INDEX idx_product_brand ON product(brand_id);
CREATE INDEX idx_product_unit ON product(unit_id);
CREATE INDEX idx_product_sku ON product(sku);
CREATE INDEX idx_product_barcode ON product(barcode);
CREATE INDEX idx_product_product_type ON product(product_type);
CREATE INDEX idx_product_active ON product(active);
CREATE INDEX idx_product_visible_in_pos ON product(visible_in_pos);
CREATE INDEX idx_product_deleted_at ON product(deleted_at) WHERE deleted_at IS NULL;

-- Trigger for automatic updated_at timestamp
CREATE TRIGGER update_product_updated_at BEFORE UPDATE ON product
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Column comments (documentación en español)
COMMENT ON TABLE product IS 'Catálogo de productos - tabla principal del inventario de productos';
COMMENT ON COLUMN product.id IS 'Identificador único del producto';
COMMENT ON COLUMN product.sku IS 'Código SKU del producto (identificador interno)';
COMMENT ON COLUMN product.barcode IS 'Código de barras del producto para escaneo';
COMMENT ON COLUMN product.name IS 'Nombre comercial del producto';
COMMENT ON COLUMN product.description IS 'Descripción detallada del producto';
COMMENT ON COLUMN product.category_id IS 'Referencia a la categoría del producto';
COMMENT ON COLUMN product.brand_id IS 'Referencia a la marca del producto';
COMMENT ON COLUMN product.unit_id IS 'Referencia a la unidad de medida base';
COMMENT ON COLUMN product.product_type IS 'Tipo de producto: ESTANDAR, SERVICIO, COMBO, RECETA';
COMMENT ON COLUMN product.active IS 'Indica si el producto está activo';
COMMENT ON COLUMN product.visible_in_pos IS 'Indica si el producto es visible en el POS';
COMMENT ON COLUMN product.cost_price IS 'Precio de costo del producto (precio de compra)';
COMMENT ON COLUMN product.sale_price IS 'Precio de venta del producto (precio al público)';
COMMENT ON COLUMN product.price_2 IS 'Precio alternativo nivel 2 (mayoreo)';
COMMENT ON COLUMN product.price_3 IS 'Precio alternativo nivel 3 (especial)';
COMMENT ON COLUMN product.iva_percentage IS 'Porcentaje de IVA aplicado';
COMMENT ON COLUMN product.consumption_tax_value IS 'Valor de impuesto al consumo';
COMMENT ON COLUMN product.current_stock IS 'Cantidad actual en inventario';
COMMENT ON COLUMN product.min_stock IS 'Stock mínimo para alertas de reposición';
COMMENT ON COLUMN product.max_stock IS 'Stock máximo permitido';
COMMENT ON COLUMN product.manages_inventory IS 'Habilita control de inventario';
COMMENT ON COLUMN product.manages_batches IS 'Habilita manejo de lotes';
COMMENT ON COLUMN product.manages_serial IS 'Habilita manejo de números de serie';
COMMENT ON COLUMN product.allow_negative_stock IS ' Permite inventario negativo (ventas sin stock)';
COMMENT ON COLUMN product.image_url IS 'URL de la imagen del producto';
COMMENT ON COLUMN product.enterprise_id IS 'Identificador de la empresa (tenant)';
COMMENT ON COLUMN product.created_at IS 'Fecha de creación del registro';
COMMENT ON COLUMN product.updated_at IS 'Fecha de última modificación';
COMMENT ON COLUMN product.deleted_at IS 'Fecha de eliminación lógica';