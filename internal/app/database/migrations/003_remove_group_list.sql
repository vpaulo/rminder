-- Remove unused group_list table and related columns/indexes

DROP INDEX IF EXISTS idx_group_list_position;
DROP INDEX IF EXISTS idx_list_group_id;
DROP TABLE IF EXISTS group_list;

ALTER TABLE list DROP COLUMN group_id;
ALTER TABLE persistence DROP COLUMN group_id;
