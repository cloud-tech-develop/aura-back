-- Table: third_parties (SQLite version for offline mode)
CREATE TABLE IF NOT EXISTS third_parties (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    document_number VARCHAR(50) NOT NULL,
    document_type VARCHAR(20) NOT NULL,
    personal_email VARCHAR(150),
    commercial_name VARCHAR(255),
    address VARCHAR(255),
    phone VARCHAR(20),
    additional_email VARCHAR(150),
    tax_responsibility VARCHAR(20) NOT NULL DEFAULT 'RESPONSIBLE',
    is_client BOOLEAN NOT NULL DEFAULT 0,
    is_provider BOOLEAN NOT NULL DEFAULT 0,
    is_employee BOOLEAN NOT NULL DEFAULT 0,
    municipality_id VARCHAR(10),
    municipality VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_third_parties_document ON third_parties(document_number);
CREATE INDEX idx_third_parties_user ON third_parties(user_id);