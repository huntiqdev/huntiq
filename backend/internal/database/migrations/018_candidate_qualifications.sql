-- ============================================================================
-- 018_candidate_qualifications.sql
-- Table: candidate_qualifications (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for candidate_qualifications
ALTER TABLE candidate_qualifications ADD COLUMN IF NOT EXISTS candidate_id integer NOT NULL;
ALTER TABLE candidate_qualifications ADD COLUMN IF NOT EXISTS qualification_type_id integer;
ALTER TABLE candidate_qualifications ADD COLUMN IF NOT EXISTS field_of_study character varying(200) NOT NULL;
ALTER TABLE candidate_qualifications ADD COLUMN IF NOT EXISTS institution character varying(200) NOT NULL;
ALTER TABLE candidate_qualifications ADD COLUMN IF NOT EXISTS graduation_year integer;
ALTER TABLE candidate_qualifications ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'candidate_qualifications' AND constraint_name = 'candidate_qualifications_pkey'
    ) THEN
        ALTER TABLE ONLY candidate_qualifications ADD CONSTRAINT candidate_qualifications_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_candidate_qualifications_cand ON candidate_qualifications USING btree (candidate_id);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'candidate_qualifications' AND constraint_name = 'candidate_qualifications_candidate_id_fkey'
    ) THEN
        ALTER TABLE ONLY candidate_qualifications ADD CONSTRAINT candidate_qualifications_candidate_id_fkey FOREIGN KEY (candidate_id) REFERENCES candidate_profiles(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'candidate_qualifications' AND constraint_name = 'candidate_qualifications_qualification_type_id_fkey'
    ) THEN
        ALTER TABLE ONLY candidate_qualifications ADD CONSTRAINT candidate_qualifications_qualification_type_id_fkey FOREIGN KEY (qualification_type_id) REFERENCES qualification_types(id) ON DELETE SET NULL;
    END IF;
END $$;
