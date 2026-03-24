-- Add max_enterprises column to plans table
ALTER TABLE public.plans ADD COLUMN IF NOT EXISTS max_enterprises INT DEFAULT NULL;

COMMENT ON COLUMN public.plans.max_enterprises IS 'Cantidad máxima de empresas permitidas en el plan';
