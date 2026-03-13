
ALTER TABLE public.tenants DROP COLUMN IF EXISTS enterprise_id;

DROP INDEX IF EXISTS idx_enterprises_status;
DROP INDEX IF EXISTS idx_enterprises_slug;
DROP TABLE IF EXISTS public.enterprises;

ALTER TABLE public.tenants ADD COLUMN IF NOT EXISTS name TEXT NOT NULL;
