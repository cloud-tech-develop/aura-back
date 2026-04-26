-- Offline/SQLite equivalent of public schema bootstrap
CREATE TABLE IF NOT EXISTS tenants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS enterprises (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tenant_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    commercial_name VARCHAR(255),
    slug VARCHAR(100) UNIQUE NOT NULL,
    sub_domain VARCHAR(50),
    email VARCHAR(150) NOT NULL UNIQUE,
    document VARCHAR(50),
    dv VARCHAR(2),
    phone VARCHAR(20),
    municipality_id VARCHAR(10),
    municipality VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    settings TEXT DEFAULT '{}',
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_enterprises_slug ON enterprises(slug);
CREATE INDEX IF NOT EXISTS idx_enterprises_email ON enterprises(email);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    enterprise_id INTEGER NOT NULL,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_enterprise_id ON users(enterprise_id);

CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(100),
    level INTEGER NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, role_id)
);

-- Seeder de roles
INSERT OR IGNORE INTO roles (name, description, level) VALUES
('SUPERADMIN', 'Super admin', 0),
('ADMIN', 'Administrator', 1),
('SUPERVISOR', 'Supervisor', 2),
('USER', 'Standard user', 3),
('SELLER', 'Sales and customers access', 3),
('CASHIER', 'Cashier', 3),
('ACCOUNTANT', 'Accountant', 3);

CREATE TABLE IF NOT EXISTS plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    enterprise_id INTEGER NOT NULL,
    max_users INTEGER DEFAULT NULL,
    max_enterprises INTEGER DEFAULT NULL,
    trial_until TEXT DEFAULT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_plans_enterprise_id ON plans(enterprise_id);

CREATE TABLE IF NOT EXISTS branches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(50),
    enterprise_id INTEGER NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_branches_enterprise ON branches(enterprise_id);
