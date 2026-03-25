-- Add sync fields to core tables

-- Category
ALTER TABLE category ADD COLUMN global_id VARCHAR(36) UNIQUE;
ALTER TABLE category ADD COLUMN sync_status VARCHAR(20) DEFAULT 'PENDING';
ALTER TABLE category ADD COLUMN last_synced_at TIMESTAMPTZ;
CREATE INDEX idx_category_global_id ON category(global_id);

-- Brand
ALTER TABLE brand ADD COLUMN global_id VARCHAR(36) UNIQUE;
ALTER TABLE brand ADD COLUMN sync_status VARCHAR(20) DEFAULT 'PENDING';
ALTER TABLE brand ADD COLUMN last_synced_at TIMESTAMPTZ;
CREATE INDEX idx_brand_global_id ON brand(global_id);

-- Product
ALTER TABLE product ADD COLUMN global_id VARCHAR(36) UNIQUE;
ALTER TABLE product ADD COLUMN sync_status VARCHAR(20) DEFAULT 'PENDING';
ALTER TABLE product ADD COLUMN last_synced_at TIMESTAMPTZ;
CREATE INDEX idx_product_global_id ON product(global_id);

-- Sales Order
ALTER TABLE sales_order ADD COLUMN global_id VARCHAR(36) UNIQUE;
ALTER TABLE sales_order ADD COLUMN sync_status VARCHAR(20) DEFAULT 'PENDING';
ALTER TABLE sales_order ADD COLUMN last_synced_at TIMESTAMPTZ;
CREATE INDEX idx_sales_order_global_id ON sales_order(global_id);

-- Third Parties
ALTER TABLE third_parties ADD COLUMN global_id VARCHAR(36) UNIQUE;
ALTER TABLE third_parties ADD COLUMN sync_status VARCHAR(20) DEFAULT 'PENDING';
ALTER TABLE third_parties ADD COLUMN last_synced_at TIMESTAMPTZ;
CREATE INDEX idx_third_parties_global_id ON third_parties(global_id);

-- Invoice
ALTER TABLE invoice ADD COLUMN global_id VARCHAR(36) UNIQUE;
ALTER TABLE invoice ADD COLUMN sync_status VARCHAR(20) DEFAULT 'PENDING';
ALTER TABLE invoice ADD COLUMN last_synced_at TIMESTAMPTZ;
CREATE INDEX idx_invoice_global_id ON invoice(global_id);

-- Invoice Prefix
ALTER TABLE invoice_prefix ADD COLUMN global_id VARCHAR(36) UNIQUE;
ALTER TABLE invoice_prefix ADD COLUMN sync_status VARCHAR(20) DEFAULT 'PENDING';
ALTER TABLE invoice_prefix ADD COLUMN last_synced_at TIMESTAMPTZ;
CREATE INDEX idx_invoice_prefix_global_id ON invoice_prefix(global_id);
