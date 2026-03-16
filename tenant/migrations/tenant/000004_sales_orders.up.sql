-- Tabla de Órdenes de Venta
CREATE TABLE sales_order (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(50) NOT NULL,
    customer_id BIGINT REFERENCES third_parties(id),
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    empresa_id BIGINT NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL,
    discount DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_total DECIMAL(12,2) NOT NULL,
    total DECIMAL(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING_PAYMENT' CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED', 'COMPLETED')),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    
    CONSTRAINT sales_order_empresa_fk FOREIGN KEY (empresa_id) REFERENCES public.enterprises(id),
    CONSTRAINT sales_order_user_fk FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT sales_order_customer_fk FOREIGN KEY (customer_id) REFERENCES third_parties(id),
    CONSTRAINT sales_order_number_unique UNIQUE (empresa_id, branch_id, order_number)
);

CREATE INDEX idx_sales_order_empresa ON sales_order(empresa_id);
CREATE INDEX idx_sales_order_branch ON sales_order(branch_id);
CREATE INDEX idx_sales_order_customer ON sales_order(customer_id);
CREATE INDEX idx_sales_order_status ON sales_order(status);

-- Tabla de Items de Órdenes de Venta
CREATE TABLE sales_order_item (
    id BIGSERIAL PRIMARY KEY,
    sales_order_id BIGINT NOT NULL REFERENCES sales_order(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL,
    tax_amount DECIMAL(12,2) NOT NULL,
    total DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT sales_order_item_order_fk FOREIGN KEY (sales_order_id) REFERENCES sales_order(id) ON DELETE CASCADE,
    CONSTRAINT sales_order_item_product_fk FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE INDEX idx_sales_order_item_order ON sales_order_item(sales_order_id);
CREATE INDEX idx_sales_order_item_product ON sales_order_item(product_id);

-- Trigger para updated_at
CREATE TRIGGER update_sales_order_updated_at BEFORE UPDATE ON sales_order
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
