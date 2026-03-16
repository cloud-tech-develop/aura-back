DROP TRIGGER IF EXISTS update_users_updated_at ON public.users;
DROP FUNCTION IF EXISTS update_updated_at_column();
ALTER TABLE public.users DROP COLUMN IF EXISTS updated_at;
