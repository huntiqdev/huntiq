-- ============================================================================
-- Migration 094: Consolidate Database Schema to ~60 Tables and JSONB Attributes
-- ============================================================================
-- 1. Adds native `attributes JSONB` to `products` with high-performance GIN indexing
--    and backfills all specifications from `product_attributes`.
-- 2. Retires 10 dead zero-row legacy tables (brand_dictionary, synonym_dictionary,
--    product_families, cv_documents, candidate_qualifications, qualification_types,
--    cv_access_requests, category_migration_log, admin_role_migration_audit).
-- 3. Drops legacy `categories` table and replaces with a zero-cost backward-compatible
--    VIEW over `category_hierarchy(level = 0)`.
-- 4. Cleans up orphaned foreign key constraints.
-- ============================================================================

BEGIN;

-- 1. Add attributes JSONB column to products and create GIN index
ALTER TABLE products ADD COLUMN IF NOT EXISTS attributes JSONB NOT NULL DEFAULT '{}'::jsonb;

-- Backfill technical specifications from product_attributes into products.attributes
WITH agg AS (
    SELECT 
        pa.product_id,
        jsonb_object_agg(
            ad.slug, 
            COALESCE(
                pa.value_json,
                to_jsonb(pa.value_text),
                to_jsonb(pa.value_number),
                to_jsonb(pa.value_boolean),
                to_jsonb(pa.value_date)
            )
        ) AS attr_json
    FROM product_attributes pa
    JOIN attribute_definitions ad ON ad.id = pa.attribute_id
    GROUP BY pa.product_id
)
UPDATE products p
SET attributes = p.attributes || agg.attr_json,
    updated_at = NOW()
FROM agg
WHERE p.id = agg.product_id;

-- Create high-performance GIN index for JSONB path querying
CREATE INDEX IF NOT EXISTS idx_products_attributes_gin ON products USING gin (attributes jsonb_path_ops);

-- 2. Remove legacy constraints that point to deprecated tables
ALTER TABLE products DROP CONSTRAINT IF EXISTS products_category_id_fkey;
ALTER TABLE candidate_profiles DROP CONSTRAINT IF EXISTS candidate_profiles_highest_qualification_id_fkey;
ALTER TABLE product_models DROP CONSTRAINT IF EXISTS product_models_family_id_fkey;

-- 3. Drop dead / empty legacy tables with zero active references
DROP TABLE IF EXISTS brand_dictionary CASCADE;
DROP TABLE IF EXISTS synonym_dictionary CASCADE;
DROP TABLE IF EXISTS category_migration_log CASCADE;
DROP TABLE IF EXISTS admin_role_migration_audit CASCADE;
DROP TABLE IF EXISTS product_families CASCADE;
DROP TABLE IF EXISTS cv_access_requests CASCADE;
DROP TABLE IF EXISTS cv_documents CASCADE;
DROP TABLE IF EXISTS candidate_qualifications CASCADE;
DROP TABLE IF EXISTS qualification_types CASCADE;

-- Drop obsolete empty views
DROP VIEW IF EXISTS category_applicable_attributes CASCADE;
DROP VIEW IF EXISTS category_breadcrumbs CASCADE;
DROP VIEW IF EXISTS category_stats CASCADE;

-- 4. Replace legacy categories table with transparent view over category_hierarchy
DROP TABLE IF EXISTS categories CASCADE;

CREATE OR REPLACE VIEW categories AS
SELECT id, name_en AS name, name_am, slug, icon, is_active, display_order, created_at, updated_at
FROM category_hierarchy
WHERE level = 0;

-- 5. Add flexible JSONB metadata to candidate_profiles for extensible resume data
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb;

COMMIT;

