-- Create shrinkage_reason table
CREATE TABLE IF NOT EXISTS shrinkage_reason (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    requires_authorization BOOLEAN DEFAULT FALSE,
    authorization_threshold DECIMAL(12,2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create shrinkage table
CREATE TABLE IF NOT EXISTS shrinkage (
    id BIGSERIAL PRIMARY KEY,
    shrinkage_number VARCHAR(50) NOT NULL,
    branch_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    reason_id BIGINT NOT NULL,
    shrinkage_date DATE NOT NULL DEFAULT CURRENT_DATE,
    total_value DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED', 'CANCELLED')),
    notes TEXT,
    authorized_by BIGINT,
    authorized_at TIMESTAMPTZ,
    cancellation_reason TEXT,
    cancelled_by BIGINT,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_shrinkage_number ON shrinkage(shrinkage_number);
CREATE INDEX idx_shrinkage_branch ON shrinkage(branch_id);
CREATE INDEX idx_shrinkage_reason ON shrinkage(reason_id);
CREATE INDEX idx_shrinkage_status ON shrinkage(status);
CREATE INDEX idx_shrinkage_date ON shrinkage(shrinkage_date);

-- Create shrinkage_item table
CREATE TABLE IF NOT EXISTS shrinkage_item (
    id BIGSERIAL PRIMARY KEY,
    shrinkage_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    batch_number VARCHAR(50),
    serial_number VARCHAR(100),
    quantity DECIMAL(10,2) NOT NULL,
    unit_cost DECIMAL(12,2) NOT NULL,
    total_value DECIMAL(12,2) NOT NULL,
    reason_detail TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shrinkage_item_shrinkage ON shrinkage_item(shrinkage_id);

-- Add foreign key constraints
ALTER TABLE shrinkage ADD CONSTRAINT fk_shrinkage_branch FOREIGN KEY (branch_id) REFERENCES branch(id);
ALTER TABLE shrinkage ADD CONSTRAINT fk_shrinkage_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE shrinkage ADD CONSTRAINT fk_shrinkage_reason FOREIGN KEY (reason_id) REFERENCES shrinkage_reason(id);
ALTER TABLE shrinkage_item ADD CONSTRAINT fk_shrinkage_item_shrinkage FOREIGN KEY (shrinkage_id) REFERENCES shrinkage(id) ON DELETE CASCADE;
ALTER TABLE shrinkage_item ADD CONSTRAINT fk_shrinkage_item_product FOREIGN KEY (product_id) REFERENCES product(id);

-- Insert default shrinkage reasons
INSERT INTO shrinkage_reason (code, name, description, requires_authorization, authorization_threshold) VALUES
('DAMAGE', 'Daño', 'Producto dañado en almacén o transporte', FALSE, 100.00),
('EXPIRED', 'Vencimiento', 'Producto vencido sin posibilidad de venta', FALSE, NULL),
('THEFT', 'Robo', 'Pérdida por robo o hurto', TRUE, 50.00),
('LOSS', 'Extravío', 'Producto extraviado sin causa conocida', FALSE, NULL),
('QUALITY', 'Calidad', 'Producto que no cumple estándares de calidad', FALSE, NULL),
('SPOILAGE', 'Deterioro', 'Producto deteriorado por condiciones ambientales', FALSE, NULL),
('HANDLING', 'Manipulación', 'Daño por manipulación incorrecta', FALSE, 200.00),
('OTHER', 'Otro', 'Otras causas de merma', FALSE, NULL);
