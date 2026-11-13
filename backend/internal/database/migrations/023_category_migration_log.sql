-- ============================================================================
-- 023_category_migration_log.sql
-- Table: category_migration_log (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for category_migration_log
ALTER TABLE category_migration_log ADD COLUMN IF NOT EXISTS old_category_slug character varying(100);
ALTER TABLE category_migration_log ADD COLUMN IF NOT EXISTS new_category_slug character varying(100);
ALTER TABLE category_migration_log ADD COLUMN IF NOT EXISTS products_migrated integer DEFAULT 0;
ALTER TABLE category_migration_log ADD COLUMN IF NOT EXISTS migrated_at timestamp with time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_migration_log' AND constraint_name = 'category_migration_log_pkey'
    ) THEN
        ALTER TABLE ONLY category_migration_log ADD CONSTRAINT category_migration_log_pkey PRIMARY KEY (id);
    END IF;
END $$;
