-- ============================================================================
-- 011_admin_sessions.sql
-- Table: admin_sessions (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for admin_sessions
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS session_token_hash character varying(64) NOT NULL;
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS user_id integer NOT NULL;
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS portal character varying(50) NOT NULL;
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS role character varying(30) NOT NULL;
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS ip_address character varying(50);
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS user_agent text;
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS last_active_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS idle_timeout_seconds integer DEFAULT 1800 NOT NULL;
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS absolute_expires_at timestamp with time zone NOT NULL;
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS revoked_at timestamp with time zone;
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS revocation_reason text;
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_sessions' AND constraint_name = 'admin_sessions_pkey'
    ) THEN
        ALTER TABLE ONLY admin_sessions ADD CONSTRAINT admin_sessions_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_sessions' AND constraint_name = 'admin_sessions_session_token_hash_key'
    ) THEN
        ALTER TABLE ONLY admin_sessions ADD CONSTRAINT admin_sessions_session_token_hash_key UNIQUE (session_token_hash);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_admin_sessions_expiry ON admin_sessions USING btree (absolute_expires_at);

CREATE INDEX IF NOT EXISTS idx_admin_sessions_token ON admin_sessions USING btree (session_token_hash) WHERE (revoked_at IS NULL);

CREATE INDEX IF NOT EXISTS idx_admin_sessions_user ON admin_sessions USING btree (user_id) WHERE (revoked_at IS NULL);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_sessions' AND constraint_name = 'admin_sessions_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY admin_sessions ADD CONSTRAINT admin_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
