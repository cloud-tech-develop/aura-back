-- Tabla de Carrito
CREATE TABLE cart (
    id BIGSERIAL PRIMARY KEY,
    cart_code VARCHAR(50) NOT NULL,
    customer_id BIGINT,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    empresa_id BIGINT NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    total DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'SAVED', 'CONVERTED', 'EXPIRED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    
    CONSTRAINT cart_empresa_fk FOREIGN KEY (empresa_id) REFERENCES public.enterprises(id),
    CONSTRAINT cart_user_fk FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT cart_customer_fk FOREIGN KEY (customer_id) REFERENCES third_parties(id),
    CONSTRAINT cart_code_unique UNIQUE (empresa_id, cart_code)
);

CREATE INDEX idx_cart_empresa ON cart(empresa_id);
CREATE INDEX idx_cart_user ON cart(user_id);
CREATE INDEX idx_cart_status ON cart(status);

-- Tabla de Items del Carrito
CREATE TABLE cart_item (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES cart(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL,
    tax_amount DECIMAL(12,2) NOT NULL,
    total DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    
    CONSTRAINT cart_item_cart_fk FOREIGN KEY (cart_id) REFERENCES cart(id) ON DELETE CASCADE,
    CONSTRAINT cart_item_product_fk FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE INDEX idx_cart_item_cart ON cart_item(cart_id);
CREATE INDEX idx_cart_item_product ON cart_item(product_id);

-- Trigger para updated_at
CREATE TRIGGER update_cart_updated_at BEFORE UPDATE ON cart
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_cart_item_updated_at BEFORE UPDATE ON cart_item
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
