-- Table: presentation (SQLite version for offline mode)
CREATE TABLE IF NOT EXISTS presentation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    -- Product reference
    product_id INTEGER NOT NULL,

    -- Presentation details
    name VARCHAR(100) NOT NULL,
    factor DECIMAL(10,4) NOT NULL DEFAULT 1,
    barcode VARCHAR(50),

    -- Pricing
    cost_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    sale_price DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Default flags
    default_purchase INTEGER NOT NULL DEFAULT 0,
    default_sale INTEGER NOT NULL DEFAULT 0,

    -- Enterprise ownership
    enterprise_id INTEGER NOT NULL,

    -- Timestamps
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_presentation_enterprise ON presentation(enterprise_id);
CREATE INDEX IF NOT EXISTS idx_presentation_product ON presentation(product_id);
CREATE INDEX IF NOT EXISTS idx_presentation_barcode ON presentation(barcode);
