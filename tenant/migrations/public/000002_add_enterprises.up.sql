
CREATE TABLE IF NOT EXISTS public.enterprises (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL REFERENCES public.tenants(id),
    name            VARCHAR(255) NOT NULL,
    commercial_name VARCHAR(255),
    slug            VARCHAR(100) UNIQUE NOT NULL,
    sub_domain      VARCHAR(50),
    email           VARCHAR(150) NOT NULL UNIQUE,
    dv              VARCHAR(2),
    phone           VARCHAR(20),
    municipality_id VARCHAR(10),
    municipality    VARCHAR(255),
    status          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'SUSPENDED', 'DEBT', 'LEGAL_COLLECTION')),
    settings        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ DEFAULT NULL
);

COMMENT ON COLUMN public.enterprises.name IS 'Razón social';

CREATE INDEX idx_enterprises_slug ON public.enterprises(slug);
CREATE INDEX idx_enterprises_sub_domain ON public.enterprises(sub_domain);
CREATE INDEX idx_enterprises_email ON public.enterprises(email);
CREATE INDEX idx_enterprises_status ON public.enterprises(status);
