-- Rollback inventory tables
DROP TRIGGER IF EXISTS update_inventory_updated_at ON inventory;
DROP TABLE IF EXISTS inventory_movement;
DROP TABLE IF EXISTS movement_reason;
DROP TABLE IF EXISTS inventory;
