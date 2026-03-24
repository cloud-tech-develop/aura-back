-- Drop foreign key constraints first
ALTER TABLE purchase_payment DROP CONSTRAINT IF EXISTS fk_pp_user;
ALTER TABLE purchase_payment DROP CONSTRAINT IF EXISTS fk_pp_purchase;
ALTER TABLE purchase_item DROP CONSTRAINT IF EXISTS fk_purchase_item_product;
ALTER TABLE purchase_item DROP CONSTRAINT IF EXISTS fk_purchase_item_purchase;
ALTER TABLE purchase DROP CONSTRAINT IF EXISTS fk_purchase_order;
ALTER TABLE purchase DROP CONSTRAINT IF EXISTS fk_purchase_user;
ALTER TABLE purchase DROP CONSTRAINT IF EXISTS fk_purchase_branch;
ALTER TABLE purchase DROP CONSTRAINT IF EXISTS fk_purchase_supplier;
ALTER TABLE purchase_order_item DROP CONSTRAINT IF EXISTS fk_po_item_product;
ALTER TABLE purchase_order_item DROP CONSTRAINT IF EXISTS fk_po_item_order;
ALTER TABLE purchase_order DROP CONSTRAINT IF EXISTS fk_po_user;
ALTER TABLE purchase_order DROP CONSTRAINT IF EXISTS fk_po_branch;
ALTER TABLE purchase_order DROP CONSTRAINT IF EXISTS fk_po_supplier;

-- Drop tables
DROP TABLE IF EXISTS purchase_payment;
DROP TABLE IF EXISTS purchase_item;
DROP TABLE IF EXISTS purchase;
DROP TABLE IF EXISTS purchase_order_item;
DROP TABLE IF EXISTS purchase_order;
