-- ============================================================================
-- 057_synonym_dictionary.sql
-- Table: synonym_dictionary (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for synonym_dictionary
ALTER TABLE synonym_dictionary ADD COLUMN IF NOT EXISTS term character varying(200) NOT NULL;
ALTER TABLE synonym_dictionary ADD COLUMN IF NOT EXISTS synonyms jsonb NOT NULL;
ALTER TABLE synonym_dictionary ADD COLUMN IF NOT EXISTS context character varying(100);
ALTER TABLE synonym_dictionary ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;
ALTER TABLE synonym_dictionary ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();

COMMENT ON TABLE synonym_dictionary IS 'Term synonyms for query expansion';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'synonym_dictionary' AND constraint_name = 'synonym_dictionary_pkey'
    ) THEN
        ALTER TABLE ONLY synonym_dictionary ADD CONSTRAINT synonym_dictionary_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_synonym_dictionary_term ON synonym_dictionary USING btree (term);
