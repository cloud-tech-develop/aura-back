-- Drop presentation trigger
DROP TRIGGER IF EXISTS update_presentation_updated_at ON presentation;

-- Drop table
DROP TABLE IF EXISTS presentation CASCADE;