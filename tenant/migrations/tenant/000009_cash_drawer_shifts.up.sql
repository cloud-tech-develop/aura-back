-- Create cash_drawer table
CREATE TABLE IF NOT EXISTS cash_drawer (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL DEFAULT 'MAIN',
    is_active BOOLEAN DEFAULT TRUE,
    min_float DECIMAL(12,2) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_cash_drawer_branch ON cash_drawer(branch_id);

-- Create cash_shift table
CREATE TABLE IF NOT EXISTS cash_shift (
    id BIGSERIAL PRIMARY KEY,
    cash_drawer_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    branch_id BIGINT NOT NULL,
    opening_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    closing_amount DECIMAL(12,2),
    expected_amount DECIMAL(12,2),
    difference DECIMAL(12,2),
    opening_notes TEXT,
    closing_notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED', 'AUDITED')),
    opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMPTZ,
    closed_by BIGINT,
    authorized_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cash_shift_drawer ON cash_shift(cash_drawer_id);
CREATE INDEX idx_cash_shift_user ON cash_shift(user_id);
CREATE INDEX idx_cash_shift_status ON cash_shift(status);
CREATE INDEX idx_cash_shift_opened ON cash_shift(opened_at);

-- Create cash_movement table
CREATE TABLE IF NOT EXISTS cash_movement (
    id BIGSERIAL PRIMARY KEY,
    shift_id BIGINT NOT NULL,
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('IN', 'OUT')),
    reason VARCHAR(30) NOT NULL CHECK (reason IN ('SALE', 'OPENING', 'CLOSING', 'EXPENSE', 'DROPS', 'WITHDRAWAL', 'ADJUSTMENT', 'REFUND')),
    amount DECIMAL(12,2) NOT NULL,
    reference_id BIGINT,
    reference_type VARCHAR(50),
    notes TEXT,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT cash_movement_amount_check CHECK (amount > 0)
);

CREATE INDEX idx_cash_movement_shift ON cash_movement(shift_id);
CREATE INDEX idx_cash_movement_type ON cash_movement(movement_type);
CREATE INDEX idx_cash_movement_created ON cash_movement(created_at);

-- Add foreign key constraints
ALTER TABLE cash_drawer ADD CONSTRAINT fk_cash_drawer_branch FOREIGN KEY (branch_id) REFERENCES branch(id);
ALTER TABLE cash_shift ADD CONSTRAINT fk_cash_shift_drawer FOREIGN KEY (cash_drawer_id) REFERENCES cash_drawer(id);
ALTER TABLE cash_shift ADD CONSTRAINT fk_cash_shift_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE cash_shift ADD CONSTRAINT fk_cash_shift_branch FOREIGN KEY (branch_id) REFERENCES branch(id);
ALTER TABLE cash_movement ADD CONSTRAINT fk_cash_movement_shift FOREIGN KEY (shift_id) REFERENCES cash_shift(id);
ALTER TABLE cash_movement ADD CONSTRAINT fk_cash_movement_user FOREIGN KEY (user_id) REFERENCES users(id);
