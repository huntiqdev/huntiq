-- ============================================================================
-- 065_optimize_database_tables_and_cross_discovery.sql
-- Database Table Optimization: Prune legacy synthetic unused tables
-- Strengthen cross-category discovery indexes on classification_relationships
-- Dedicated single migration (HUNTIQ Enterprise Architecture)
-- ============================================================================

-- 1. Remove dangling foreign key column from products table
ALTER TABLE products DROP CONSTRAINT IF EXISTS products_canonical_entity_id_fkey;
ALTER TABLE products DROP COLUMN IF EXISTS canonical_entity_id;

-- 2. Safely drop unreferenced, empty synthetic/legacy tables
DROP TABLE IF EXISTS brand_dictionary CASCADE;
DROP TABLE IF EXISTS synonym_dictionary CASCADE;
DROP TABLE IF EXISTS category_redirects CASCADE;
DROP TABLE IF EXISTS collection_memberships CASCADE;
DROP TABLE IF EXISTS discovery_collections CASCADE;
DROP TABLE IF EXISTS entity_lifecycle_audit_logs CASCADE;
DROP TABLE IF EXISTS entity_relationships CASCADE;
DROP TABLE IF EXISTS canonical_entities CASCADE;
DROP TABLE IF EXISTS entity_types CASCADE;
DROP TABLE IF EXISTS merchant_inventory CASCADE;
DROP TABLE IF EXISTS invoices CASCADE;
DROP TABLE IF EXISTS service_appointments CASCADE;
DROP TABLE IF EXISTS merchant_analytics_events CASCADE;
DROP TABLE IF EXISTS referral_fraud_evaluations CASCADE;

-- 3. Ensure performance indexes for normalized cross-category discovery
CREATE INDEX IF NOT EXISTS idx_classif_rel_source ON classification_relationships (source_concept_id);
CREATE INDEX IF NOT EXISTS idx_classif_rel_target ON classification_relationships (target_concept_id);
CREATE INDEX IF NOT EXISTS idx_classif_rel_type ON classification_relationships (relationship_type_id);
CREATE INDEX IF NOT EXISTS idx_classif_rel_status ON classification_relationships (status);
