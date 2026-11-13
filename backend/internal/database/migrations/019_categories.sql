-- ============================================================================
-- 019_categories.sql
-- Table: categories (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for categories
ALTER TABLE categories ADD COLUMN IF NOT EXISTS name character varying(100) DEFAULT ''::character varying NOT NULL;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS slug character varying(100) DEFAULT ''::character varying NOT NULL;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS name_en text;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS name_am text;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS name_or text;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS name_ti text;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS parent_id integer;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS sort_order integer DEFAULT 0 NOT NULL;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'categories' AND constraint_name = 'categories_name_key'
    ) THEN
        ALTER TABLE ONLY categories ADD CONSTRAINT categories_name_key UNIQUE (name);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'categories' AND constraint_name = 'categories_pkey'
    ) THEN
        ALTER TABLE ONLY categories ADD CONSTRAINT categories_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'categories' AND constraint_name = 'categories_slug_key'
    ) THEN
        ALTER TABLE ONLY categories ADD CONSTRAINT categories_slug_key UNIQUE (slug);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_categories_active ON categories USING btree (is_active) WHERE (is_active = true);

CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories USING btree (parent_id) WHERE (parent_id IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_categories_slug_active ON categories USING btree (slug, is_active);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'categories' AND constraint_name = 'categories_parent_id_fkey'
    ) THEN
        ALTER TABLE ONLY categories ADD CONSTRAINT categories_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE;
    END IF;
END $$;
