-- Tabla de Prefijos de Factura
CREATE TABLE invoice_prefix (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT NOT NULL,
    empresa_id BIGINT NOT NULL,
    prefix VARCHAR(10) NOT NULL,
    next_sequence BIGINT NOT NULL DEFAULT 1,
    description VARCHAR(100),
    
    CONSTRAINT invoice_prefix_branch_fk FOREIGN KEY (branch_id) REFERENCES public.branches(id),
    CONSTRAINT invoice_prefix_empresa_fk FOREIGN KEY (empresa_id) REFERENCES public.enterprises(id),
    CONSTRAINT invoice_prefix_unique UNIQUE (empresa_id, branch_id, prefix)
);

-- Tabla de Facturas
CREATE TABLE invoice (
    id BIGSERIAL PRIMARY KEY,
    invoice_number VARCHAR(50) NOT NULL,
    prefix VARCHAR(10) NOT NULL,
    sequence BIGINT NOT NULL,
    sales_order_id BIGINT NOT NULL REFERENCES sales_order(id),
    customer_id BIGINT REFERENCES public.third_parties(id),
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    branch_id BIGINT NOT NULL,
    empresa_id BIGINT NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL,
    discount DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_total DECIMAL(12,2) NOT NULL,
    total DECIMAL(12,2) NOT NULL,
    issue_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    due_date TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL DEFAULT 'ISSUED' CHECK (status IN ('ISSUED', 'PAID', 'CANCELLED', 'OVERDUE')),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT invoice_empresa_fk FOREIGN KEY (empresa_id) REFERENCES public.enterprises(id),
    CONSTRAINT invoice_branch_fk FOREIGN KEY (branch_id) REFERENCES public.branches(id),
    CONSTRAINT invoice_user_fk FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT invoice_sales_order_fk FOREIGN KEY (sales_order_id) REFERENCES sales_order(id),
    CONSTRAINT invoice_customer_fk FOREIGN KEY (customer_id) REFERENCES public.third_parties(id),
    CONSTRAINT invoice_number_unique UNIQUE (empresa_id, invoice_number),
    CONSTRAINT invoice_sequence_unique UNIQUE (empresa_id, branch_id, prefix, sequence)
);

CREATE INDEX idx_invoice_empresa ON invoice(empresa_id);
CREATE INDEX idx_invoice_branch ON invoice(branch_id);
CREATE INDEX idx_invoice_customer ON invoice(customer_id);
CREATE INDEX idx_invoice_number ON invoice(invoice_number);
CREATE INDEX idx_invoice_status ON invoice(status);
CREATE INDEX idx_invoice_deleted_at ON invoice(deleted_at) WHERE deleted_at IS NULL;

-- Tabla de Items de Factura
CREATE TABLE invoice_item (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES invoice(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    sales_order_item_id BIGINT REFERENCES sales_order_item(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL,
    tax_amount DECIMAL(12,2) NOT NULL,
    total DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT invoice_item_invoice_fk FOREIGN KEY (invoice_id) REFERENCES invoice(id) ON DELETE CASCADE,
    CONSTRAINT invoice_item_product_fk FOREIGN KEY (product_id) REFERENCES product(id),
    CONSTRAINT invoice_item_sales_order_item_fk FOREIGN KEY (sales_order_item_id) REFERENCES sales_order_item(id)
);

CREATE INDEX idx_invoice_item_invoice ON invoice_item(invoice_id);
CREATE INDEX idx_invoice_item_product ON invoice_item(product_id);

-- Trigger para updated_at
CREATE TRIGGER update_invoice_updated_at BEFORE UPDATE ON invoice
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
