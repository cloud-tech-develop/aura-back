-- Create purchase_order table
CREATE TABLE IF NOT EXISTS purchase_order (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(50) NOT NULL,
    supplier_id BIGINT NOT NULL,
    branch_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    order_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expected_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PARTIAL', 'RECEIVED', 'CANCELLED')),
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    total DECIMAL(12,2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_purchase_order_number ON purchase_order(order_number);
CREATE INDEX idx_purchase_order_supplier ON purchase_order(supplier_id);
CREATE INDEX idx_purchase_order_status ON purchase_order(status);
CREATE INDEX idx_purchase_order_date ON purchase_order(order_date);

-- Create purchase_order_item table
CREATE TABLE IF NOT EXISTS purchase_order_item (
    id BIGSERIAL PRIMARY KEY,
    purchase_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    received_quantity DECIMAL(10,2) NOT NULL DEFAULT 0,
    unit_cost DECIMAL(12,2) NOT NULL,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL DEFAULT 0,
    line_total DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_po_item_order ON purchase_order_item(purchase_order_id);

-- Create purchase table (for received goods)
CREATE TABLE IF NOT EXISTS purchase (
    id BIGSERIAL PRIMARY KEY,
    purchase_number VARCHAR(50) NOT NULL,
    purchase_order_id BIGINT,
    supplier_id BIGINT NOT NULL,
    branch_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    purchase_date DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'COMPLETED' CHECK (status IN ('COMPLETED', 'PARTIAL', 'CANCELLED')),
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    total DECIMAL(12,2) NOT NULL DEFAULT 0,
    paid_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    pending_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_purchase_number ON purchase(purchase_number);
CREATE INDEX idx_purchase_supplier ON purchase(supplier_id);
CREATE INDEX idx_purchase_order_ref ON purchase(purchase_order_id);
CREATE INDEX idx_purchase_status ON purchase(status);
CREATE INDEX idx_purchase_date ON purchase(purchase_date);

-- Create purchase_item table
CREATE TABLE IF NOT EXISTS purchase_item (
    id BIGSERIAL PRIMARY KEY,
    purchase_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    unit_cost DECIMAL(12,2) NOT NULL,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL DEFAULT 0,
    line_total DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_purchase_item_purchase ON purchase_item(purchase_id);

-- Create purchase_payment table
CREATE TABLE IF NOT EXISTS purchase_payment (
    id BIGSERIAL PRIMARY KEY,
    purchase_id BIGINT NOT NULL,
    payment_method VARCHAR(30) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    reference_number VARCHAR(50),
    notes TEXT,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_purchase_payment_purchase ON purchase_payment(purchase_id);
CREATE INDEX idx_purchase_payment_date ON purchase_payment(created_at);

-- Add foreign key constraints
ALTER TABLE purchase_order ADD CONSTRAINT fk_po_supplier FOREIGN KEY (supplier_id) REFERENCES third_parties(id);
ALTER TABLE purchase_order ADD CONSTRAINT fk_po_branch FOREIGN KEY (branch_id) REFERENCES branch(id);
ALTER TABLE purchase_order ADD CONSTRAINT fk_po_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE purchase_order_item ADD CONSTRAINT fk_po_item_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_order(id) ON DELETE CASCADE;
ALTER TABLE purchase_order_item ADD CONSTRAINT fk_po_item_product FOREIGN KEY (product_id) REFERENCES product(id);
ALTER TABLE purchase ADD CONSTRAINT fk_purchase_supplier FOREIGN KEY (supplier_id) REFERENCES third_parties(id);
ALTER TABLE purchase ADD CONSTRAINT fk_purchase_branch FOREIGN KEY (branch_id) REFERENCES branch(id);
ALTER TABLE purchase ADD CONSTRAINT fk_purchase_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE purchase ADD CONSTRAINT fk_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_order(id);
ALTER TABLE purchase_item ADD CONSTRAINT fk_purchase_item_purchase FOREIGN KEY (purchase_id) REFERENCES purchase(id) ON DELETE CASCADE;
ALTER TABLE purchase_item ADD CONSTRAINT fk_purchase_item_product FOREIGN KEY (product_id) REFERENCES product(id);
ALTER TABLE purchase_payment ADD CONSTRAINT fk_pp_purchase FOREIGN KEY (purchase_id) REFERENCES purchase(id) ON DELETE CASCADE;
ALTER TABLE purchase_payment ADD CONSTRAINT fk_pp_user FOREIGN KEY (user_id) REFERENCES users(id);
