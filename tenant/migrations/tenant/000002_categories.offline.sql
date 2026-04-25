-- Table: category (SQLite version for offline mode)
CREATE TABLE IF NOT EXISTS category (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id INTEGER,
    default_tax_rate DECIMAL(5,2) DEFAULT 0.00,
    active BOOLEAN DEFAULT 1,
    enterprise_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_category_enterprise ON category(enterprise_id);
CREATE INDEX idx_category_parent ON category(parent_id);

-- Table: brand (SQLite version for offline mode)
CREATE TABLE IF NOT EXISTS brand (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT 1,
    enterprise_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_brand_enterprise ON brand(enterprise_id);