-- Create commission_rule table
CREATE TABLE IF NOT EXISTS commission_rule (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    commission_type VARCHAR(30) NOT NULL CHECK (commission_type IN ('PERCENTAGE_SALE', 'PERCENTAGE_MARGIN', 'FIXED_AMOUNT')),
    employee_id BIGINT,
    product_id BIGINT,
    category_id BIGINT,
    value DECIMAL(12,2) NOT NULL,
    min_sale_amount DECIMAL(12,2) DEFAULT 0,
    start_date DATE,
    end_date DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_commission_rule_employee ON commission_rule(employee_id) WHERE employee_id IS NOT NULL;
CREATE INDEX idx_commission_rule_product ON commission_rule(product_id) WHERE product_id IS NOT NULL;
CREATE INDEX idx_commission_rule_category ON commission_rule(category_id) WHERE category_id IS NOT NULL;
CREATE INDEX idx_commission_rule_active ON commission_rule(is_active, start_date, end_date);

-- Create commission table
CREATE TABLE IF NOT EXISTS commission (
    id BIGSERIAL PRIMARY KEY,
    sales_order_id BIGINT NOT NULL,
    employee_id BIGINT NOT NULL,
    branch_id BIGINT NOT NULL,
    rule_id BIGINT,
    sale_amount DECIMAL(12,2) NOT NULL,
    profit_margin DECIMAL(12,2),
    commission_type VARCHAR(30) NOT NULL,
    commission_rate DECIMAL(12,2) NOT NULL,
    commission_amount DECIMAL(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'SETTLED', 'CANCELLED')),
    settled_at TIMESTAMPTZ,
    settled_by BIGINT,
    settlement_period VARCHAR(20),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_commission_employee ON commission(employee_id);
CREATE INDEX idx_commission_order ON commission(sales_order_id);
CREATE INDEX idx_commission_status ON commission(status);
CREATE INDEX idx_commission_created ON commission(created_at);
CREATE INDEX idx_commission_period ON commission(settlement_period) WHERE settlement_period IS NOT NULL;

-- Add foreign key constraints
ALTER TABLE commission_rule ADD CONSTRAINT fk_cr_employee FOREIGN KEY (employee_id) REFERENCES third_parties(id);
ALTER TABLE commission_rule ADD CONSTRAINT fk_cr_product FOREIGN KEY (product_id) REFERENCES product(id);
ALTER TABLE commission_rule ADD CONSTRAINT fk_cr_category FOREIGN KEY (category_id) REFERENCES category(id);
ALTER TABLE commission ADD CONSTRAINT fk_commission_order FOREIGN KEY (sales_order_id) REFERENCES sales_order(id);
ALTER TABLE commission ADD CONSTRAINT fk_commission_employee FOREIGN KEY (employee_id) REFERENCES third_parties(id);
ALTER TABLE commission ADD CONSTRAINT fk_commission_branch FOREIGN KEY (branch_id) REFERENCES branch(id);
ALTER TABLE commission ADD CONSTRAINT fk_commission_rule FOREIGN KEY (rule_id) REFERENCES commission_rule(id);
