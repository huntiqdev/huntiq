-- ============================================================================
-- 017_candidate_profiles.sql
-- Table: candidate_profiles (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for candidate_profiles
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS user_id integer NOT NULL;
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS professional_title character varying(200) NOT NULL;
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS occupation_id integer;
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS headline text;
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS experience_years integer DEFAULT 0;
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS highest_qualification_id integer;
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS skills text[];
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS languages text[];
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS portfolio_url text;
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS city character varying(100);
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS is_profile_public boolean DEFAULT true;
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS is_seeking_work boolean DEFAULT true;
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();
ALTER TABLE candidate_profiles ADD COLUMN IF NOT EXISTS updated_at timestamp without time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'candidate_profiles' AND constraint_name = 'candidate_profiles_pkey'
    ) THEN
        ALTER TABLE ONLY candidate_profiles ADD CONSTRAINT candidate_profiles_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'candidate_profiles' AND constraint_name = 'candidate_profiles_user_id_key'
    ) THEN
        ALTER TABLE ONLY candidate_profiles ADD CONSTRAINT candidate_profiles_user_id_key UNIQUE (user_id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_candidate_profiles_occupation ON candidate_profiles USING btree (occupation_id);

CREATE INDEX IF NOT EXISTS idx_candidate_profiles_seeking ON candidate_profiles USING btree (is_seeking_work);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'candidate_profiles' AND constraint_name = 'candidate_profiles_highest_qualification_id_fkey'
    ) THEN
        ALTER TABLE ONLY candidate_profiles ADD CONSTRAINT candidate_profiles_highest_qualification_id_fkey FOREIGN KEY (highest_qualification_id) REFERENCES qualification_types(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'candidate_profiles' AND constraint_name = 'candidate_profiles_occupation_id_fkey'
    ) THEN
        ALTER TABLE ONLY candidate_profiles ADD CONSTRAINT candidate_profiles_occupation_id_fkey FOREIGN KEY (occupation_id) REFERENCES occupations(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'candidate_profiles' AND constraint_name = 'candidate_profiles_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY candidate_profiles ADD CONSTRAINT candidate_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
