-- Create transfer table
CREATE TABLE IF NOT EXISTS transfer (
    id BIGSERIAL PRIMARY KEY,
    transfer_number VARCHAR(50) NOT NULL,
    origin_branch_id BIGINT NOT NULL,
    destination_branch_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'SHIPPED', 'PARTIAL', 'RECEIVED', 'CANCELLED')),
    requested_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    shipped_date TIMESTAMPTZ,
    received_date TIMESTAMPTZ,
    notes TEXT,
    shipped_by BIGINT,
    received_by BIGINT,
    cancellation_reason TEXT,
    cancelled_by BIGINT,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_transfer_number ON transfer(transfer_number);
CREATE INDEX idx_transfer_origin ON transfer(origin_branch_id);
CREATE INDEX idx_transfer_destination ON transfer(destination_branch_id);
CREATE INDEX idx_transfer_status ON transfer(status);
CREATE INDEX idx_transfer_dates ON transfer(requested_date);

-- Create transfer_item table
CREATE TABLE IF NOT EXISTS transfer_item (
    id BIGSERIAL PRIMARY KEY,
    transfer_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    requested_quantity DECIMAL(10,2) NOT NULL,
    shipped_quantity DECIMAL(10,2) NOT NULL DEFAULT 0,
    received_quantity DECIMAL(10,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT transfer_item_quantity_check CHECK (requested_quantity > 0)
);

CREATE INDEX idx_transfer_item_transfer ON transfer_item(transfer_id);
CREATE INDEX idx_transfer_item_product ON transfer_item(product_id);

-- Add foreign key constraints
ALTER TABLE transfer ADD CONSTRAINT fk_transfer_origin FOREIGN KEY (origin_branch_id) REFERENCES branch(id);
ALTER TABLE transfer ADD CONSTRAINT fk_transfer_destination FOREIGN KEY (destination_branch_id) REFERENCES branch(id);
ALTER TABLE transfer ADD CONSTRAINT fk_transfer_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE transfer_item ADD CONSTRAINT fk_transfer_item_transfer FOREIGN KEY (transfer_id) REFERENCES transfer(id) ON DELETE CASCADE;
ALTER TABLE transfer_item ADD CONSTRAINT fk_transfer_item_product FOREIGN KEY (product_id) REFERENCES product(id);
