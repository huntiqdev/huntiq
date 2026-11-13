-- ============================================================================
-- 014_brand_dictionary.sql
-- Table: brand_dictionary (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for brand_dictionary
ALTER TABLE brand_dictionary ADD COLUMN IF NOT EXISTS name character varying(200) NOT NULL;
ALTER TABLE brand_dictionary ADD COLUMN IF NOT EXISTS normalized_name character varying(200) NOT NULL;
ALTER TABLE brand_dictionary ADD COLUMN IF NOT EXISTS aliases jsonb;
ALTER TABLE brand_dictionary ADD COLUMN IF NOT EXISTS category_ids integer[];
ALTER TABLE brand_dictionary ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;
ALTER TABLE brand_dictionary ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();
ALTER TABLE brand_dictionary ADD COLUMN IF NOT EXISTS updated_at timestamp without time zone DEFAULT now();

COMMENT ON TABLE brand_dictionary IS 'Known brands/manufacturers for search recognition';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'brand_dictionary' AND constraint_name = 'brand_dictionary_name_key'
    ) THEN
        ALTER TABLE ONLY brand_dictionary ADD CONSTRAINT brand_dictionary_name_key UNIQUE (name);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'brand_dictionary' AND constraint_name = 'brand_dictionary_pkey'
    ) THEN
        ALTER TABLE ONLY brand_dictionary ADD CONSTRAINT brand_dictionary_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_brand_dictionary_category_ids ON brand_dictionary USING gin (category_ids);

CREATE INDEX IF NOT EXISTS idx_brand_dictionary_normalized ON brand_dictionary USING btree (normalized_name);
