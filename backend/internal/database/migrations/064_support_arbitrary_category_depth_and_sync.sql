-- ============================================================================
-- 064_support_arbitrary_category_depth_and_sync.sql
-- Relax category_hierarchy valid_level constraint to support arbitrary depth (>= 0)
-- Ensure materialized path and lineage indexes exist
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_hierarchy' AND constraint_name = 'valid_level'
    ) THEN
        ALTER TABLE category_hierarchy DROP CONSTRAINT valid_level;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_hierarchy' AND constraint_name = 'valid_level'
    ) THEN
        ALTER TABLE category_hierarchy ADD CONSTRAINT valid_level CHECK (level >= 0);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_path ON category_hierarchy USING btree (path);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_path_ops ON category_hierarchy USING btree (path varchar_pattern_ops);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_level ON category_hierarchy USING btree (level);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_parent ON category_hierarchy USING btree (parent_id);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_active_leaf ON category_hierarchy USING btree (is_active, is_leaf);

DELETE FROM category_aliases a USING category_aliases b
WHERE a.id < b.id AND a.category_id = b.category_id AND a.alias = b.alias;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_aliases' AND constraint_name = 'uq_category_aliases_cat_alias'
    ) THEN
        ALTER TABLE ONLY category_aliases ADD CONSTRAINT uq_category_aliases_cat_alias UNIQUE (category_id, alias);
    END IF;
END $$;

