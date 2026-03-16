-- Tabla de Pagos
CREATE TABLE payment (
    id BIGSERIAL PRIMARY KEY,
    sales_order_id BIGINT NOT NULL REFERENCES sales_order(id),
    payment_method VARCHAR(50) NOT NULL CHECK (payment_method IN ('CASH', 'CARD', 'TRANSFER', 'CREDIT')),
    amount DECIMAL(12,2) NOT NULL,
    reference VARCHAR(100),
    card_type VARCHAR(20) CHECK (card_type IN ('CREDIT', 'DEBIT')),
    card_last_four VARCHAR(4),
    bank_name VARCHAR(100),
    authorization_code VARCHAR(50),
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    empresa_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT payment_order_fk FOREIGN KEY (sales_order_id) REFERENCES sales_order(id),
    CONSTRAINT payment_user_fk FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT payment_empresa_fk FOREIGN KEY (empresa_id) REFERENCES public.enterprises(id),
    CONSTRAINT payment_positive_amount CHECK (amount > 0)
);

CREATE INDEX idx_payment_order ON payment(sales_order_id);
CREATE INDEX idx_payment_user ON payment(user_id);
CREATE INDEX idx_payment_empresa ON payment(empresa_id);

-- Tabla de Caja Registradora
CREATE TABLE cash_drawer (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    empresa_id BIGINT NOT NULL,
    opening_balance DECIMAL(12,2) NOT NULL DEFAULT 0,
    closing_balance DECIMAL(12,2),
    cash_in DECIMAL(12,2) NOT NULL DEFAULT 0,
    cash_out DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED')),
    opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMPTZ,
    
    CONSTRAINT cash_drawer_user_fk FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT cash_drawer_empresa_fk FOREIGN KEY (empresa_id) REFERENCES public.enterprises(id)
);

CREATE INDEX idx_cash_drawer_user ON cash_drawer(user_id);
CREATE INDEX idx_cash_drawer_status ON cash_drawer(status);
