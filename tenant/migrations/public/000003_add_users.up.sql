CREATE TABLE IF NOT EXISTS public.users (
    id             BIGSERIAL PRIMARY KEY,
    enterprise_id  BIGINT NOT NULL REFERENCES public.enterprises(id),
    email          TEXT UNIQUE NOT NULL,
    name           TEXT NOT NULL,
    password_hash  VARCHAR(255) NOT NULL,
    active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_users_email ON public.users(email);
CREATE INDEX idx_users_enterprise_id ON public.users(enterprise_id);
