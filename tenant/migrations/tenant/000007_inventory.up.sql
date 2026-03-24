-- Tabla de Inventario por producto y sucursal
CREATE TABLE IF NOT EXISTS inventory (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES product(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER DEFAULT 0,
    min_stock INTEGER DEFAULT 0,
    max_stock INTEGER,
    location VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    
    CONSTRAINT inventory_product_branch_unique UNIQUE (product_id, branch_id),
    CONSTRAINT inventory_quantity_check CHECK (quantity >= 0),
    CONSTRAINT inventory_reserved_check CHECK (reserved_quantity >= 0)
);

CREATE INDEX idx_inventory_product ON inventory(product_id);
CREATE INDEX idx_inventory_branch ON inventory(branch_id);

-- Tabla de Movimientos de Inventario (Kardex)
CREATE TABLE IF NOT EXISTS inventory_movement (
    id BIGSERIAL PRIMARY KEY,
    inventory_id BIGINT NOT NULL REFERENCES inventory(id),
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('ENTRY', 'EXIT', 'ADJUSTMENT')),
    movement_reason VARCHAR(30) NOT NULL CHECK (movement_reason IN ('SALE', 'PURCHASE', 'SHRINKAGE', 'TRANSFER_IN', 'TRANSFER_OUT', 'ADJUSTMENT', 'RETURN', 'INITIAL', 'DAMAGE', 'THEFT', 'EXPIRED')),
    quantity INTEGER NOT NULL,
    previous_balance INTEGER NOT NULL,
    new_balance INTEGER NOT NULL,
    reference_id BIGINT,
    reference_type VARCHAR(50),
    batch_number VARCHAR(50),
    serial_number VARCHAR(100),
    expiration_date DATE,
    notes TEXT,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT inventory_movement_quantity_check CHECK (quantity > 0),
    CONSTRAINT inventory_movement_balance_check CHECK (new_balance >= 0)
);

CREATE INDEX idx_movement_inventory ON inventory_movement(inventory_id);
CREATE INDEX idx_movement_type ON inventory_movement(movement_type);
CREATE INDEX idx_movement_reason ON inventory_movement(movement_reason);
CREATE INDEX idx_movement_created ON inventory_movement(created_at);
CREATE INDEX idx_movement_reference ON inventory_movement(reference_type, reference_id);
CREATE INDEX idx_movement_batch ON inventory_movement(batch_number) WHERE batch_number IS NOT NULL;
CREATE INDEX idx_movement_serial ON inventory_movement(serial_number) WHERE serial_number IS NOT NULL;

-- Tabla de Razones de Movimiento
CREATE TABLE IF NOT EXISTS movement_reason (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('ENTRY', 'EXIT', 'ADJUSTMENT')),
    requires_authorization BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default movement reasons
INSERT INTO movement_reason (code, name, description, movement_type, requires_authorization) VALUES
('SALE', 'Venta', 'Salida por venta', 'EXIT', FALSE),
('PURCHASE', 'Compra', 'Entrada por compra', 'ENTRY', FALSE),
('SHRINKAGE', 'Merma', 'Salida por merma/daño', 'EXIT', TRUE),
('TRANSFER_IN', 'Transferencia Entrada', 'Entrada por transferencia', 'ENTRY', FALSE),
('TRANSFER_OUT', 'Transferencia Salida', 'Salida por transferencia', 'EXIT', FALSE),
('ADJUSTMENT', 'Ajuste', 'Ajuste de inventario', 'ADJUSTMENT', TRUE),
('RETURN', 'Devolución', 'Entrada por devolución', 'ENTRY', FALSE),
('INITIAL', 'Inventario Inicial', 'Inventario inicial', 'ENTRY', FALSE),
('DAMAGE', 'Daño', 'Salida por daño', 'EXIT', TRUE),
('THEFT', 'Hurto', 'Salida por hurto', 'EXIT', TRUE),
('EXPIRED', 'Vencido', 'Salida por vencimiento', 'EXIT', TRUE);

-- Trigger para updated_at en inventory
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_inventory_updated_at BEFORE UPDATE ON inventory
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
