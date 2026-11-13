-- ============================================================================
-- 020_category_aliases.sql
-- Table: category_aliases (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for category_aliases
ALTER TABLE category_aliases ADD COLUMN IF NOT EXISTS category_id integer NOT NULL;
ALTER TABLE category_aliases ADD COLUMN IF NOT EXISTS alias character varying(200) NOT NULL;
ALTER TABLE category_aliases ADD COLUMN IF NOT EXISTS language character varying(10) DEFAULT 'en'::character varying;
ALTER TABLE category_aliases ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_aliases' AND constraint_name = 'category_aliases_pkey'
    ) THEN
        ALTER TABLE ONLY category_aliases ADD CONSTRAINT category_aliases_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_category_aliases_alias ON category_aliases USING btree (alias);

CREATE INDEX IF NOT EXISTS idx_category_aliases_category ON category_aliases USING btree (category_id);

CREATE INDEX IF NOT EXISTS idx_category_aliases_lang_alias ON category_aliases USING btree (language, alias);

CREATE INDEX IF NOT EXISTS idx_category_aliases_trgm ON category_aliases USING gin (alias gin_trgm_ops);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_aliases' AND constraint_name = 'category_aliases_category_id_fkey'
    ) THEN
        ALTER TABLE ONLY category_aliases ADD CONSTRAINT category_aliases_category_id_fkey FOREIGN KEY (category_id) REFERENCES category_hierarchy(id) ON DELETE CASCADE;
    END IF;
END $$;
