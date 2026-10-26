-- Migration 034: Drop redundant duplicate secondary indexes (AUD-008)
--
-- Rationale:
-- 1. idx_category_hierarchy_slug and idx_category_hierarchy_slug_unique are duplicate/subsets of uq_category_hierarchy_slug.
-- 2. idx_attribute_definitions_slug is a non-unique duplicate of attribute_definitions_slug_key.
-- 3. idx_users_phone_role is a non-unique duplicate of users_phone_role_key.
-- 4. idx_brands_slug is a non-unique duplicate of brands_slug_key.
--
-- All lookups on these columns are backed with equal or higher selectivity by the corresponding UNIQUE constraint indexes.

DROP INDEX IF EXISTS idx_category_hierarchy_slug_unique;
DROP INDEX IF EXISTS idx_category_hierarchy_slug;
DROP INDEX IF EXISTS idx_attribute_definitions_slug;
DROP INDEX IF EXISTS idx_users_phone_role;
DROP INDEX IF EXISTS idx_brands_slug;
