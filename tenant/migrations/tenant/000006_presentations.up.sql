-- Table: presentation
-- Product presentations (variants/sizes like kilo, libra, unidad, etc.)

CREATE TABLE IF NOT EXISTS presentation (
    id BIGSERIAL PRIMARY KEY,
    
    -- Product reference
    product_id BIGINT NOT NULL,              -- Product foreign key
    
    -- Presentation details
    name VARCHAR(100) NOT NULL,               -- Presentation name (Kilo, Libra, etc.)
    factor DECIMAL(10,4) NOT NULL DEFAULT 1, -- Conversion factor to base unit
    barcode VARCHAR(50),                     -- Barcode for this presentation (optional)
    
    -- Pricing (can differ from product base price)
    cost_price DECIMAL(12,2) NOT NULL DEFAULT 0,  -- Cost price for this presentation
    sale_price DECIMAL(12,2) NOT NULL DEFAULT 0,  -- Sale price for this presentation
    
    -- Default flags
    default_purchase BOOLEAN NOT NULL DEFAULT false,  -- Default for purchase orders
    default_sale BOOLEAN NOT NULL DEFAULT false,      -- Default for sales/POS
    
    -- Enterprise ownership
    enterprise_id BIGINT NOT NULL,            -- Enterprise foreign key
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT presentation_enterprise_fk FOREIGN KEY (enterprise_id) REFERENCES public.enterprises(id),
    CONSTRAINT presentation_product_fk FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE,
    CONSTRAINT presentation_name_unique UNIQUE (enterprise_id, product_id, name)
);

-- Indexes for performance optimization
CREATE INDEX idx_presentation_enterprise ON presentation(enterprise_id);
CREATE INDEX idx_presentation_product ON presentation(product_id);
CREATE INDEX idx_presentation_barcode ON presentation(barcode);
CREATE INDEX idx_presentation_default_purchase ON presentation(default_purchase) WHERE default_purchase = true;
CREATE INDEX idx_presentation_default_sale ON presentation(default_sale) WHERE default_sale = true;
CREATE INDEX idx_presentation_deleted_at ON presentation(deleted_at) WHERE deleted_at IS NULL;

-- Trigger for automatic updated_at timestamp
CREATE TRIGGER update_presentation_updated_at BEFORE UPDATE ON presentation
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Column comments (documentación en español)
COMMENT ON TABLE presentation IS 'Presentaciones de productos - variantes/tamaños (kilo, libra, unidad, etc.)';
COMMENT ON COLUMN presentation.id IS 'Identificador único de la presentación';
COMMENT ON COLUMN presentation.product_id IS 'Referencia al producto padre';
COMMENT ON COLUMN presentation.name IS 'Nombre de la presentación (kilo, libra, unidad, etc.)';
COMMENT ON COLUMN presentation.factor IS 'Factor de conversión a la unidad base';
COMMENT ON COLUMN presentation.barcode IS 'Código de barras de esta presentación (opcional)';
COMMENT ON COLUMN presentation.cost_price IS 'Precio de costo de esta presentación';
COMMENT ON COLUMN presentation.sale_price IS 'Precio de venta de esta presentación';
COMMENT ON COLUMN presentation.default_purchase IS 'Indica si es la presentación por defecto para compras';
COMMENT ON COLUMN presentation.default_sale IS 'Indica si es la presentación por defecto para ventas/POS';
COMMENT ON COLUMN presentation.enterprise_id IS 'Identificador de la empresa (tenant)';
COMMENT ON COLUMN presentation.created_at IS 'Fecha de creación del registro';
COMMENT ON COLUMN presentation.updated_at IS 'Fecha de última modificación';
COMMENT ON COLUMN presentation.deleted_at IS 'Fecha de eliminación lógica';