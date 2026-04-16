-- Table: unit (unidades de medida)
-- Unidades de medida para productos (kg, lb, unidad, etc.)

CREATE TABLE IF NOT EXISTS unit (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    abbreviation VARCHAR(20) NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    allow_decimals BOOLEAN DEFAULT TRUE,
    enterprise_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unit_enterprise_fk FOREIGN KEY (enterprise_id) REFERENCES public.enterprises(id),
    CONSTRAINT unit_name_unique UNIQUE (enterprise_id, name),
    CONSTRAINT unit_abbreviation_unique UNIQUE (enterprise_id, abbreviation)
);

-- Indexes
CREATE INDEX idx_unit_enterprise ON unit(enterprise_id);
CREATE INDEX idx_unit_deleted_at ON unit(deleted_at) WHERE deleted_at IS NULL;

-- Trigger for updated_at
CREATE TRIGGER update_unit_updated_at BEFORE UPDATE ON unit
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE unit IS 'Unidades de medida';
COMMENT ON COLUMN unit.name IS 'Nombre de la unidad de medida';
COMMENT ON COLUMN unit.abbreviation IS 'Abreviatura de la unidad (kg, lb, un, etc.)';
COMMENT ON COLUMN unit.active IS 'Estado activo de la unidad';
COMMENT ON COLUMN unit.allow_decimals IS 'Permite cantidades decimales';
