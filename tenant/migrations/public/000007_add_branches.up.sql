-- Tabla de Sucursales (Branches)
CREATE TABLE branches (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(50),
    empresa_id BIGINT NOT NULL REFERENCES public.enterprises(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT branches_empresa_fk FOREIGN KEY (empresa_id) REFERENCES public.enterprises(id)
);

CREATE INDEX idx_branches_empresa ON branches(empresa_id);

-- Trigger para updated_at
CREATE TRIGGER update_branches_updated_at BEFORE UPDATE ON branches
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
