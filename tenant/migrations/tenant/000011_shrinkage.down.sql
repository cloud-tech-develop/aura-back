-- Drop foreign key constraints first
ALTER TABLE shrinkage_item DROP CONSTRAINT IF EXISTS fk_shrinkage_item_product;
ALTER TABLE shrinkage_item DROP CONSTRAINT IF EXISTS fk_shrinkage_item_shrinkage;
ALTER TABLE shrinkage DROP CONSTRAINT IF EXISTS fk_shrinkage_reason;
ALTER TABLE shrinkage DROP CONSTRAINT IF EXISTS fk_shrinkage_user;
ALTER TABLE shrinkage DROP CONSTRAINT IF EXISTS fk_shrinkage_branch;

-- Drop tables
DROP TABLE IF EXISTS shrinkage_item;
DROP TABLE IF EXISTS shrinkage;
DROP TABLE IF EXISTS shrinkage_reason;
