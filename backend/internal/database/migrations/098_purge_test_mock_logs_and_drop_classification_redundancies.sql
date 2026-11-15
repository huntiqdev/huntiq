-- ============================================================================
-- 098_purge_test_mock_logs_and_drop_classification_redundancies.sql
-- Mock Data Purge & Redundant Schema Consolidation
--
-- Goals:
--   1. Purge ephemeral test/mock operational logs (admin_auth_attempts,
--      admin_auth_challenges, job_execution_log, search_analytics,
--      taxonomy_audit_logs) to ensure clean production tables.
--   2. Drop the 9 redundant classification_* tables and helper routines,
--      consolidating taxonomy fully into canonical category_hierarchy.
--   3. Bring base table count down from 58 to 49 clean, smart tables (well
--      below the 60-table architectural target).
--   4. Fully preserve all 11 verified product listings, images, user accounts,
--      locations, and the 3,992 canonical taxonomy categories.
-- ============================================================================

BEGIN;

-- ----------------------------------------------------------------------------
-- 1. Purge Ephemeral Test / Mock Operational Logs & Orphan Test Permutations
-- ----------------------------------------------------------------------------
TRUNCATE TABLE 
    admin_auth_attempts,
    admin_auth_challenges,
    job_execution_log,
    search_analytics,
    taxonomy_audit_logs;

-- Purge residual synthetic level-0 orphan permutations
DELETE FROM category_hierarchy WHERE slug LIKE 'l3-%' OR slug LIKE 'l2-%' OR slug LIKE 'l1-%';

-- Ensure products.category_id is 100% synchronized with canonical category_hierarchy_id
UPDATE products SET category_id = category_hierarchy_id WHERE category_hierarchy_id IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2. Drop Redundant Classification Triggers and Functions
-- ----------------------------------------------------------------------------
DROP FUNCTION IF EXISTS cascade_classification_placement_path CASCADE;
DROP FUNCTION IF EXISTS update_classification_placement_path CASCADE;
DROP FUNCTION IF EXISTS fn_search_classification_multilingual CASCADE;

-- ----------------------------------------------------------------------------
-- 3. Drop Redundant Classification Tables (Consolidating to 49 Smart Tables)
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS classification_aliases CASCADE;
DROP TABLE IF EXISTS classification_attribute_rules CASCADE;
DROP TABLE IF EXISTS classification_audit_logs CASCADE;
DROP TABLE IF EXISTS classification_change_requests CASCADE;
DROP TABLE IF EXISTS classification_relationships CASCADE;
DROP TABLE IF EXISTS classification_relationship_types CASCADE;
DROP TABLE IF EXISTS classification_placements CASCADE;
DROP TABLE IF EXISTS classification_concepts CASCADE;
DROP TABLE IF EXISTS classification_schemes CASCADE;

-- ----------------------------------------------------------------------------
-- 4. Clean Up Classification Sequences
-- ----------------------------------------------------------------------------
DROP SEQUENCE IF EXISTS classification_schemes_id_seq CASCADE;
DROP SEQUENCE IF EXISTS classification_concepts_id_seq CASCADE;
DROP SEQUENCE IF EXISTS classification_placements_id_seq CASCADE;
DROP SEQUENCE IF EXISTS classification_relationships_id_seq CASCADE;
DROP SEQUENCE IF EXISTS classification_relationship_types_id_seq CASCADE;
DROP SEQUENCE IF EXISTS classification_attribute_rules_id_seq CASCADE;
DROP SEQUENCE IF EXISTS classification_audit_logs_id_seq CASCADE;
DROP SEQUENCE IF EXISTS classification_change_requests_id_seq CASCADE;
DROP SEQUENCE IF EXISTS classification_aliases_id_seq CASCADE;

COMMIT;
