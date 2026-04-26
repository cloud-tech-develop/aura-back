-- Table: category (SQLite version for offline mode)
CREATE TABLE IF NOT EXISTS category (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id INTEGER,
    default_tax_rate DECIMAL(5,2) DEFAULT 0.00,
    active INTEGER DEFAULT 1,
    enterprise_id INTEGER NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_category_enterprise ON category(enterprise_id);
CREATE INDEX IF NOT EXISTS idx_category_parent ON category(parent_id);

-- Table: brand (SQLite version for offline mode)
CREATE TABLE IF NOT EXISTS brand (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    active INTEGER DEFAULT 1,
    enterprise_id INTEGER NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_brand_enterprise ON brand(enterprise_id);
