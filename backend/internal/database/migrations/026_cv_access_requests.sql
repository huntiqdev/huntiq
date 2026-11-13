-- ============================================================================
-- 026_cv_access_requests.sql
-- Table: cv_access_requests (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for cv_access_requests
ALTER TABLE cv_access_requests ADD COLUMN IF NOT EXISTS candidate_id integer NOT NULL;
ALTER TABLE cv_access_requests ADD COLUMN IF NOT EXISTS employer_id integer NOT NULL;
ALTER TABLE cv_access_requests ADD COLUMN IF NOT EXISTS job_posting_id integer;
ALTER TABLE cv_access_requests ADD COLUMN IF NOT EXISTS status character varying(50) DEFAULT 'pending'::character varying;
ALTER TABLE cv_access_requests ADD COLUMN IF NOT EXISTS request_note text;
ALTER TABLE cv_access_requests ADD COLUMN IF NOT EXISTS requested_at timestamp without time zone DEFAULT now();
ALTER TABLE cv_access_requests ADD COLUMN IF NOT EXISTS responded_at timestamp without time zone;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'cv_access_requests' AND constraint_name = 'cv_access_requests_pkey'
    ) THEN
        ALTER TABLE ONLY cv_access_requests ADD CONSTRAINT cv_access_requests_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_cv_access_candidate ON cv_access_requests USING btree (candidate_id);

CREATE INDEX IF NOT EXISTS idx_cv_access_employer ON cv_access_requests USING btree (employer_id);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'cv_access_requests' AND constraint_name = 'cv_access_requests_candidate_id_fkey'
    ) THEN
        ALTER TABLE ONLY cv_access_requests ADD CONSTRAINT cv_access_requests_candidate_id_fkey FOREIGN KEY (candidate_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'cv_access_requests' AND constraint_name = 'cv_access_requests_employer_id_fkey'
    ) THEN
        ALTER TABLE ONLY cv_access_requests ADD CONSTRAINT cv_access_requests_employer_id_fkey FOREIGN KEY (employer_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'cv_access_requests' AND constraint_name = 'cv_access_requests_job_posting_id_fkey'
    ) THEN
        ALTER TABLE ONLY cv_access_requests ADD CONSTRAINT cv_access_requests_job_posting_id_fkey FOREIGN KEY (job_posting_id) REFERENCES job_postings(id) ON DELETE SET NULL;
    END IF;
END $$;
