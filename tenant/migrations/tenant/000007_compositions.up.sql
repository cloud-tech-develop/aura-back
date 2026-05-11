-- Table: product_composition
-- Product compositions (KIT/RECETA) for composite products
CREATE TABLE IF NOT EXISTS product_composition (
    id BIGSERIAL PRIMARY KEY,

    -- Product references
    parent_product_id BIGINT NOT NULL,                -- Parent product ID (the composite product)
    child_product_id BIGINT NOT NULL,                 -- Child product ID (component product)

    -- Composition fields
    quantity DECIMAL(12,2) NOT NULL DEFAULT 1,        -- Quantity of child product needed
    type VARCHAR(20) NOT NULL DEFAULT 'KIT' CHECK (type IN ('KIT', 'RECETA')),  -- Composition type

    -- Enterprise ownership
    enterprise_id BIGINT NOT NULL,                    -- Enterprise foreign key

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT composition_no_self_ref CHECK (parent_product_id <> child_product_id),
    CONSTRAINT composition_unique_pair UNIQUE (enterprise_id, parent_product_id, child_product_id),
    CONSTRAINT composition_parent_fk FOREIGN KEY (parent_product_id) REFERENCES product(id),
    CONSTRAINT composition_child_fk FOREIGN KEY (child_product_id) REFERENCES product(id),
    CONSTRAINT composition_enterprise_fk FOREIGN KEY (enterprise_id) REFERENCES public.enterprises(id)
);

-- Indexes for performance optimization
CREATE INDEX idx_composition_enterprise ON product_composition(enterprise_id);
CREATE INDEX idx_composition_parent ON product_composition(parent_product_id);
CREATE INDEX idx_composition_child ON product_composition(child_product_id);

-- Trigger for automatic updated_at timestamp
CREATE TRIGGER update_composition_updated_at BEFORE UPDATE ON product_composition
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Column comments
COMMENT ON TABLE product_composition IS 'Composiciones de productos - relacion padre/hijo para productos KIT y RECETA';
COMMENT ON COLUMN product_composition.id IS 'Identificador unico de la composicion';
COMMENT ON COLUMN product_composition.parent_product_id IS 'ID del producto padre (producto compuesto)';
COMMENT ON COLUMN product_composition.child_product_id IS 'ID del producto hijo (componente)';
COMMENT ON COLUMN product_composition.quantity IS 'Cantidad del componente necesaria para una unidad del producto compuesto';
COMMENT ON COLUMN product_composition.type IS 'Tipo de composicion: KIT (empaquetado) o RECETA (manufacturado)';
COMMENT ON COLUMN product_composition.enterprise_id IS 'Identificador de la empresa (tenant)';
COMMENT ON COLUMN product_composition.created_at IS 'Fecha de creacion del registro';
COMMENT ON COLUMN product_composition.updated_at IS 'Fecha de ultima modificacion';
