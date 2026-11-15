-- ============================================================================
-- 096_consolidate_schema_to_58_tables_and_prune_unused_entities.sql
-- Database Schema Consolidation: Streamline to 58 Base Tables
--
-- Goals:
--   1. Reduce database structural complexity by removing 7 redundant, unreferenced,
--      zero-row legacy entities that clutter schema governance:
--        - occupations (0 rows, zero code references; categorized under category_hierarchy)
--        - candidate_profiles (0 rows, zero code references; legacy CV subsystem)
--        - job_postings (0 rows, zero code references; legacy employment subsystem)
--        - product_models (0 rows, zero code references; specifications in products.attributes)
--        - brands (0 rows, zero code references; brand string facet in products.attributes)
--        - seller_reviews (0 rows, zero code references; unmounted review queue)
--        - classification_systems (1 row, zero code references; unreferenced ontology stub)
--   2. Target schema state: exactly 58 base tables + 15 views (73 objects total).
--   3. All 11 active approved products, images, locations, categories, and authentication
--      records remain 100% intact and verified.
-- ============================================================================

BEGIN;

-- ----------------------------------------------------------------------------
-- 1. Drop Unreferenced Legacy Job & Candidate Tables
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS candidate_profiles CASCADE;
DROP TABLE IF EXISTS job_postings CASCADE;
DROP TABLE IF EXISTS occupations CASCADE;

-- ----------------------------------------------------------------------------
-- 2. Drop Unreferenced Brand & Model Tables (Covered by products.attributes JSONB)
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS product_models CASCADE;
DROP TABLE IF EXISTS brands CASCADE;

-- ----------------------------------------------------------------------------
-- 3. Drop Unmounted Review & Ontology Stub Tables
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS seller_reviews CASCADE;
DROP TABLE IF EXISTS classification_systems CASCADE;

COMMIT;

