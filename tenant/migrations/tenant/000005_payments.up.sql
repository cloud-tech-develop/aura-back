-- Tabla de Pagos
CREATE TABLE IF NOT EXISTS payment (
    id BIGSERIAL PRIMARY KEY,
    payment_type VARCHAR(20) NOT NULL DEFAULT 'SALE' CHECK (payment_type IN ('SALE', 'PURCHASE', 'ACCOUNT_RECEIVABLE', 'ACCOUNT_PAYABLE')),
    reference_id BIGINT NOT NULL,
    reference_type VARCHAR(50) NOT NULL,
    payment_method VARCHAR(30) NOT NULL CHECK (payment_method IN ('CASH', 'DEBIT_CARD', 'CREDIT_CARD', 'BANK_TRANSFER', 'CREDIT', 'VOUCHER', 'CHECK')),
    amount DECIMAL(12,2) NOT NULL,
    reference_number VARCHAR(100),
    bank_name VARCHAR(100),
    card_type VARCHAR(20) CHECK (card_type IN ('CREDIT', 'DEBIT')),
    card_last_digits VARCHAR(4),
    authorization_code VARCHAR(50),
    change_amount DECIMAL(12,2) DEFAULT 0,
    cash_drawer_id BIGINT REFERENCES cash_drawer(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    enterprise_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'COMPLETED' CHECK (status IN ('COMPLETED', 'PENDING', 'CANCELLED', 'REFUNDED')),
    cancelled_at TIMESTAMPTZ,
    cancelled_by BIGINT REFERENCES public.users(id),
    cancellation_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    
    CONSTRAINT payment_amount_check CHECK (amount > 0),
    CONSTRAINT payment_reference_fk FOREIGN KEY (reference_id) REFERENCES sales_order(id),
    CONSTRAINT payment_user_fk FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT payment_enterprise_fk FOREIGN KEY (enterprise_id) REFERENCES public.enterprises(id)
);

CREATE INDEX idx_payment_reference ON payment(reference_type, reference_id);
CREATE INDEX idx_payment_method ON payment(payment_method);
CREATE INDEX idx_payment_drawer ON payment(cash_drawer_id) WHERE cash_drawer_id IS NOT NULL;
CREATE INDEX idx_payment_branch ON payment(branch_id);
CREATE INDEX idx_payment_user ON payment(user_id);
CREATE INDEX idx_payment_created ON payment(created_at);
CREATE INDEX idx_payment_status ON payment(status);

-- Tabla de Transacciones de Pago
CREATE TABLE IF NOT EXISTS payment_transaction (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT NOT NULL REFERENCES payment(id),
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('CHARGE', 'REFUND', 'CHARGEBACK')),
    amount DECIMAL(12,2) NOT NULL,
    previous_balance DECIMAL(12,2),
    new_balance DECIMAL(12,2),
    processor_reference VARCHAR(100),
    processor_response TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT payment_tx_fk FOREIGN KEY (payment_id) REFERENCES payment(id)
);

CREATE INDEX idx_payment_tx_payment ON payment_transaction(payment_id);
CREATE INDEX idx_payment_tx_created ON payment_transaction(created_at);

-- Tabla de Caja Registradora
CREATE TABLE IF NOT EXISTS cash_drawer (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    enterprise_id BIGINT NOT NULL,
    opening_balance DECIMAL(12,2) NOT NULL DEFAULT 0,
    closing_balance DECIMAL(12,2),
    cash_in DECIMAL(12,2) NOT NULL DEFAULT 0,
    cash_out DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED')),
    opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMPTZ,
    notes TEXT,
    
    CONSTRAINT cash_drawer_user_fk FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT cash_drawer_enterprise_fk FOREIGN KEY (enterprise_id) REFERENCES public.enterprises(id)
);

CREATE INDEX idx_cash_drawer_user ON cash_drawer(user_id);
CREATE INDEX idx_cash_drawer_status ON cash_drawer(status);
CREATE INDEX idx_cash_drawer_branch ON cash_drawer(branch_id);

-- Tabla de Movimientos de Caja
CREATE TABLE IF NOT EXISTS cash_movement (
    id BIGSERIAL PRIMARY KEY,
    cash_drawer_id BIGINT NOT NULL REFERENCES cash_drawer(id),
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('IN', 'OUT')),
    amount DECIMAL(12,2) NOT NULL,
    description TEXT,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT cash_movement_amount_check CHECK (amount > 0),
    CONSTRAINT cash_movement_drawer_fk FOREIGN KEY (cash_drawer_id) REFERENCES cash_drawer(id)
);

CREATE INDEX idx_cash_movement_drawer ON cash_movement(cash_drawer_id);
CREATE INDEX idx_cash_movement_type ON cash_movement(movement_type);

-- Trigger para updated_at en payment
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_payment_updated_at BEFORE UPDATE ON payment
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
