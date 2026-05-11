CREATE TABLE IF NOT EXISTS product_composition (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_product_id INTEGER NOT NULL,
    child_product_id INTEGER NOT NULL,
    quantity REAL NOT NULL DEFAULT 1,
    type TEXT NOT NULL DEFAULT 'KIT',
    enterprise_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT
);

CREATE INDEX IF NOT EXISTS idx_composition_enterprise ON product_composition(enterprise_id);
CREATE INDEX IF NOT EXISTS idx_composition_parent ON product_composition(parent_product_id);
CREATE INDEX IF NOT EXISTS idx_composition_child ON product_composition(child_product_id);
