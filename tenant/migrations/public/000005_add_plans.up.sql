CREATE TABLE IF NOT EXISTS public.plans (
    id              BIGSERIAL PRIMARY KEY,
    enterprise_id   BIGINT NOT NULL REFERENCES public.enterprises(id),
    max_users       INT DEFAULT NULL,
    trial_until     DATE DEFAULT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_plans_enterprise_id ON public.plans(enterprise_id);
