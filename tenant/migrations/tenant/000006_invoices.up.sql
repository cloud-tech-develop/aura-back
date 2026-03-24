-- Tabla de Prefijos de Factura
CREATE TABLE IF NOT EXISTS invoice_prefix (
    id BIGSERIAL PRIMARY KEY,
    prefix VARCHAR(10) NOT NULL,
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    enterprise_id BIGINT NOT NULL REFERENCES public.enterprises(id),
    current_number INTEGER NOT NULL DEFAULT 0,
    resolution_number VARCHAR(50),
    resolution_date DATE,
    valid_from DATE,
    valid_until DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    
    CONSTRAINT invoice_prefix_unique UNIQUE (enterprise_id, branch_id, prefix)
);

CREATE INDEX idx_prefix_branch ON invoice_prefix(branch_id);
CREATE INDEX idx_prefix_enterprise ON invoice_prefix(enterprise_id);

-- Tabla de Facturas
CREATE TABLE IF NOT EXISTS invoice (
    id BIGSERIAL PRIMARY KEY,
    invoice_number VARCHAR(50) NOT NULL,
    prefix_id BIGINT NOT NULL REFERENCES invoice_prefix(id),
    invoice_type VARCHAR(20) NOT NULL DEFAULT 'SALE' CHECK (invoice_type IN ('SALE', 'CREDIT_NOTE', 'DEBIT_NOTE')),
    reference_id BIGINT,
    reference_type VARCHAR(50),
    sales_order_id BIGINT REFERENCES sales_order(id),
    customer_id BIGINT NOT NULL REFERENCES third_parties(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    enterprise_id BIGINT NOT NULL REFERENCES public.enterprises(id),
    invoice_date DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date DATE,
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_exempt DECIMAL(12,2) NOT NULL DEFAULT 0,
    taxable_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    iva_19 DECIMAL(12,2) NOT NULL DEFAULT 0,
    iva_5 DECIMAL(12,2) NOT NULL DEFAULT 0,
    reteica DECIMAL(12,2) NOT NULL DEFAULT 0,
    retefuente DECIMAL(12,2) NOT NULL DEFAULT 0,
    reteica_rate DECIMAL(5,2) DEFAULT 0,
    retefuente_rate DECIMAL(5,2) DEFAULT 0,
    total DECIMAL(12,2) NOT NULL DEFAULT 0,
    payment_method VARCHAR(30) DEFAULT 'CASH',
    notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'ISSUED', 'SENT', 'VIEWED', 'CANCELLED')),
    cancelled_at TIMESTAMPTZ,
    cancelled_by BIGINT REFERENCES public.users(id),
    cancellation_reason TEXT,
    credit_note_id BIGINT REFERENCES invoice(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT invoice_number_unique UNIQUE (enterprise_id, invoice_number)
);

CREATE INDEX idx_invoice_number ON invoice(invoice_number);
CREATE INDEX idx_invoice_prefix ON invoice(prefix_id);
CREATE INDEX idx_invoice_branch ON invoice(branch_id);
CREATE INDEX idx_invoice_customer ON invoice(customer_id);
CREATE INDEX idx_invoice_status ON invoice(status);
CREATE INDEX idx_invoice_date ON invoice(invoice_date);
CREATE INDEX idx_invoice_sales_order ON invoice(sales_order_id);
CREATE INDEX idx_invoice_type ON invoice(invoice_type);
CREATE INDEX idx_invoice_deleted ON invoice(deleted_at) WHERE deleted_at IS NULL;

-- Tabla de Items de Factura
CREATE TABLE IF NOT EXISTS invoice_item (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES invoice(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    product_name VARCHAR(200) NOT NULL,
    product_sku VARCHAR(50),
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL DEFAULT 19.00,
    tax_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    line_total DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT invoice_item_invoice_fk FOREIGN KEY (invoice_id) REFERENCES invoice(id) ON DELETE CASCADE,
    CONSTRAINT invoice_item_product_fk FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE INDEX idx_invoice_item_invoice ON invoice_item(invoice_id);
CREATE INDEX idx_invoice_item_product ON invoice_item(product_id);

-- Tabla de Logs de Factura (Audit Trail)
CREATE TABLE IF NOT EXISTS invoice_log (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES invoice(id),
    action VARCHAR(20) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    details TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoice_log_invoice ON invoice_log(invoice_id);
CREATE INDEX idx_invoice_log_created ON invoice_log(created_at);

-- Trigger para updated_at en invoice_prefix
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_invoice_prefix_updated_at BEFORE UPDATE ON invoice_prefix
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_invoice_updated_at BEFORE UPDATE ON invoice
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
