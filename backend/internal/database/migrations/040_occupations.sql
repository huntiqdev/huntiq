-- ============================================================================
-- 040_occupations.sql
-- Table: occupations (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for occupations
ALTER TABLE occupations ADD COLUMN IF NOT EXISTS category_id integer;
ALTER TABLE occupations ADD COLUMN IF NOT EXISTS title_en character varying(200) NOT NULL;
ALTER TABLE occupations ADD COLUMN IF NOT EXISTS title_am character varying(200);
ALTER TABLE occupations ADD COLUMN IF NOT EXISTS title_om character varying(200);
ALTER TABLE occupations ADD COLUMN IF NOT EXISTS title_ti character varying(200);
ALTER TABLE occupations ADD COLUMN IF NOT EXISTS slug character varying(200) NOT NULL;
ALTER TABLE occupations ADD COLUMN IF NOT EXISTS code character varying(50);
ALTER TABLE occupations ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;
ALTER TABLE occupations ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'occupations' AND constraint_name = 'occupations_pkey'
    ) THEN
        ALTER TABLE ONLY occupations ADD CONSTRAINT occupations_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'occupations' AND constraint_name = 'occupations_slug_key'
    ) THEN
        ALTER TABLE ONLY occupations ADD CONSTRAINT occupations_slug_key UNIQUE (slug);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_occupations_category ON occupations USING btree (category_id);

CREATE INDEX IF NOT EXISTS idx_occupations_slug ON occupations USING btree (slug);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'occupations' AND constraint_name = 'occupations_category_id_fkey'
    ) THEN
        ALTER TABLE ONLY occupations ADD CONSTRAINT occupations_category_id_fkey FOREIGN KEY (category_id) REFERENCES category_hierarchy(id) ON DELETE SET NULL;
    END IF;
END $$;
