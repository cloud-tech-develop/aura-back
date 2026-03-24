-- Eliminar columna description de invoice_prefix
ALTER TABLE invoice_prefix DROP COLUMN IF EXISTS description;
