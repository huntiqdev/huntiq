-- ============================================================================
-- 036_job_postings.sql
-- Table: job_postings (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for job_postings
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS product_id integer NOT NULL;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS job_title character varying(250) NOT NULL;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS occupation_id integer;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS employer_name character varying(250) NOT NULL;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS workplace_type character varying(50) DEFAULT 'on_site'::character varying;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS employment_type character varying(50) DEFAULT 'full_time'::character varying;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS min_salary numeric(15,2);
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS max_salary numeric(15,2);
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS currency character varying(10) DEFAULT 'ETB'::character varying;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS education_level character varying(100);
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS experience_years integer DEFAULT 0;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS application_deadline date;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS application_instructions text;
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS skills_required text[];
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS updated_at timestamp without time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'job_postings' AND constraint_name = 'job_postings_pkey'
    ) THEN
        ALTER TABLE ONLY job_postings ADD CONSTRAINT job_postings_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'job_postings' AND constraint_name = 'job_postings_product_id_key'
    ) THEN
        ALTER TABLE ONLY job_postings ADD CONSTRAINT job_postings_product_id_key UNIQUE (product_id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_job_postings_deadline ON job_postings USING btree (application_deadline);

CREATE INDEX IF NOT EXISTS idx_job_postings_employment_type ON job_postings USING btree (employment_type);

CREATE INDEX IF NOT EXISTS idx_job_postings_occupation ON job_postings USING btree (occupation_id);

CREATE INDEX IF NOT EXISTS idx_job_postings_workplace_type ON job_postings USING btree (workplace_type);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'job_postings' AND constraint_name = 'job_postings_occupation_id_fkey'
    ) THEN
        ALTER TABLE ONLY job_postings ADD CONSTRAINT job_postings_occupation_id_fkey FOREIGN KEY (occupation_id) REFERENCES occupations(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'job_postings' AND constraint_name = 'job_postings_product_id_fkey'
    ) THEN
        ALTER TABLE ONLY job_postings ADD CONSTRAINT job_postings_product_id_fkey FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
    END IF;
END $$;
