-- Table: product (SQLite version for offline mode)
CREATE TABLE IF NOT EXISTS product (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    -- Basic identification fields
    sku VARCHAR(50) NOT NULL,
    barcode VARCHAR(50),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- Reference fields
    category_id INTEGER,
    brand_id INTEGER,
    unit_id INTEGER NOT NULL,
    
    -- Product type classification
    product_type VARCHAR(20) NOT NULL DEFAULT 'STANDARD',
    
    -- Status and visibility
    active INTEGER NOT NULL DEFAULT 1,
    visible_in_pos INTEGER NOT NULL DEFAULT 1,
    
    -- Pricing fields
    cost_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    sale_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    price_2 DECIMAL(12,2),
    price_3 DECIMAL(12,2),
    
    -- Tax configuration
    iva_percentage DECIMAL(5,2) NOT NULL DEFAULT 19.00,
    consumption_tax_value DECIMAL(12,2) NOT NULL DEFAULT 0,
    
    -- Inventory control settings
    current_stock INTEGER NOT NULL DEFAULT 0,
    min_stock INTEGER NOT NULL DEFAULT 0,
    max_stock INTEGER NOT NULL DEFAULT 0,
    
    -- Inventory management options
    manages_inventory INTEGER NOT NULL DEFAULT 1,
    manages_batches INTEGER NOT NULL DEFAULT 0,
    manages_serial INTEGER NOT NULL DEFAULT 0,
    allow_negative_stock INTEGER NOT NULL DEFAULT 0,
    
    -- Media
    image_url VARCHAR(500),
    
    -- Enterprise ownership
    enterprise_id INTEGER NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_product_enterprise ON product(enterprise_id);
CREATE INDEX idx_product_sku ON product(sku);
CREATE INDEX idx_product_barcode ON product(barcode);
CREATE INDEX idx_product_category ON product(category_id);
CREATE INDEX idx_product_brand ON product(brand_id);
CREATE INDEX idx_product_unit ON product(unit_id);
CREATE INDEX idx_product_active ON product(active);