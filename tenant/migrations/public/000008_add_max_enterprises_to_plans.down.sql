-- Remove max_enterprises column from plans table
ALTER TABLE public.plans DROP COLUMN IF EXISTS max_enterprises;
