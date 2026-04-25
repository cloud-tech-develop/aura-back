-- Table: unit (SQLite version for offline mode)
CREATE TABLE IF NOT EXISTS unit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    abbreviation VARCHAR(20) NOT NULL,
    active BOOLEAN DEFAULT 1,
    allow_decimals BOOLEAN DEFAULT 1,
    enterprise_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_unit_enterprise ON unit(enterprise_id);