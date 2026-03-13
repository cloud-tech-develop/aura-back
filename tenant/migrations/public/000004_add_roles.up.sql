CREATE TABLE IF NOT EXISTS public.roles (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(100),
    level       INT NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS public.user_roles (
    user_id     BIGINT NOT NULL REFERENCES public.users(id),
    role_id     BIGINT NOT NULL REFERENCES public.roles(id),
    PRIMARY KEY (user_id, role_id)
);

-- Seeder for roles
INSERT INTO public.roles (name, description, level) VALUES 
('SUPERADMIN', 'Super admin', 0),
('ADMIN', 'Administrator', 1),
('SUPERVISOR', 'Supervisor', 2),
('USER', 'Standard user', 3),
('SELLER', 'Sales and customers access', 3),
('CASHIER', 'Cashier', 3),
('ACCOUNTANT', 'Accountant', 3)
ON CONFLICT (name) DO NOTHING;
