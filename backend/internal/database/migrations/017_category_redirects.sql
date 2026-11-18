-- ============================================================================
-- 017_category_redirects.sql
-- Table: category_redirects (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for category_redirects
ALTER TABLE category_redirects ADD COLUMN IF NOT EXISTS source_category_id integer;
ALTER TABLE category_redirects ADD COLUMN IF NOT EXISTS source_slug character varying(255);
ALTER TABLE category_redirects ADD COLUMN IF NOT EXISTS target_category_id integer;
ALTER TABLE category_redirects ADD COLUMN IF NOT EXISTS old_slug character varying(255);
ALTER TABLE category_redirects ADD COLUMN IF NOT EXISTS new_slug character varying(255);
ALTER TABLE category_redirects ADD COLUMN IF NOT EXISTS redirect_type character varying(50) DEFAULT 'PERMANENT'::character varying NOT NULL;
ALTER TABLE category_redirects ADD COLUMN IF NOT EXISTS reason text;
ALTER TABLE category_redirects ADD COLUMN IF NOT EXISTS created_by integer;
ALTER TABLE category_redirects ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_redirects' AND constraint_name = 'category_redirects_pkey'
    ) THEN
        ALTER TABLE ONLY category_redirects ADD CONSTRAINT category_redirects_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_redirects' AND constraint_name = 'category_redirects_source_slug_key'
    ) THEN
        ALTER TABLE ONLY category_redirects ADD CONSTRAINT category_redirects_source_slug_key UNIQUE (source_slug);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_redirects' AND constraint_name = 'category_redirects_created_by_fkey'
    ) THEN
        ALTER TABLE ONLY category_redirects ADD CONSTRAINT category_redirects_created_by_fkey FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_redirects' AND constraint_name = 'category_redirects_target_category_id_fkey'
    ) THEN
        ALTER TABLE ONLY category_redirects ADD CONSTRAINT category_redirects_target_category_id_fkey FOREIGN KEY (target_category_id) REFERENCES category_hierarchy(id) ON DELETE CASCADE;
    END IF;
END $$;
