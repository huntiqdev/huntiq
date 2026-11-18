-- ============================================================================
-- 005_admin_auth_challenges.sql
-- Table: admin_auth_challenges (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for admin_auth_challenges
ALTER TABLE admin_auth_challenges ADD COLUMN IF NOT EXISTS challenge character varying(255) NOT NULL;
ALTER TABLE admin_auth_challenges ADD COLUMN IF NOT EXISTS user_id integer;
ALTER TABLE admin_auth_challenges ADD COLUMN IF NOT EXISTS portal character varying(50) NOT NULL;
ALTER TABLE admin_auth_challenges ADD COLUMN IF NOT EXISTS ceremony_type character varying(20) NOT NULL;
ALTER TABLE admin_auth_challenges ADD COLUMN IF NOT EXISTS session_data jsonb NOT NULL;
ALTER TABLE admin_auth_challenges ADD COLUMN IF NOT EXISTS expires_at timestamp with time zone NOT NULL;
ALTER TABLE admin_auth_challenges ADD COLUMN IF NOT EXISTS used_at timestamp with time zone;
ALTER TABLE admin_auth_challenges ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_auth_challenges' AND constraint_name = 'admin_auth_challenges_challenge_key'
    ) THEN
        ALTER TABLE ONLY admin_auth_challenges ADD CONSTRAINT admin_auth_challenges_challenge_key UNIQUE (challenge);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_auth_challenges' AND constraint_name = 'admin_auth_challenges_pkey'
    ) THEN
        ALTER TABLE ONLY admin_auth_challenges ADD CONSTRAINT admin_auth_challenges_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_admin_auth_challenges_lookup ON admin_auth_challenges USING btree (challenge) WHERE (used_at IS NULL);

CREATE INDEX IF NOT EXISTS idx_admin_auth_challenges_user ON admin_auth_challenges USING btree (user_id, ceremony_type) WHERE (used_at IS NULL);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_auth_challenges' AND constraint_name = 'admin_auth_challenges_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY admin_auth_challenges ADD CONSTRAINT admin_auth_challenges_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
