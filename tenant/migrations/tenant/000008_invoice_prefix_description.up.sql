-- Añadir columna description a invoice_prefix
ALTER TABLE invoice_prefix ADD COLUMN IF NOT EXISTS description TEXT;
