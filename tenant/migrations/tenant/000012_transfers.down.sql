-- Drop foreign key constraints first
ALTER TABLE transfer_item DROP CONSTRAINT IF EXISTS fk_transfer_item_product;
ALTER TABLE transfer_item DROP CONSTRAINT IF EXISTS fk_transfer_item_transfer;
ALTER TABLE transfer DROP CONSTRAINT IF EXISTS fk_transfer_user;
ALTER TABLE transfer DROP CONSTRAINT IF EXISTS fk_transfer_destination;
ALTER TABLE transfer DROP CONSTRAINT IF EXISTS fk_transfer_origin;

-- Drop tables
DROP TABLE IF EXISTS transfer_item;
DROP TABLE IF EXISTS transfer;
