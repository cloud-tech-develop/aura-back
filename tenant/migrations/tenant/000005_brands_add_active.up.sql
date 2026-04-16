-- Add active field to brand table
ALTER TABLE brand ADD COLUMN active BOOLEAN NOT NULL DEFAULT TRUE;
CREATE INDEX idx_brand_active ON brand(active) WHERE active = TRUE;
COMMENT ON COLUMN brand.active IS 'Brand status: true = active, false = inactive';
