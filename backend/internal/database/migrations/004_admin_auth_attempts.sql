-- ============================================================================
-- 004_admin_auth_attempts.sql
-- Table: admin_auth_attempts (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for admin_auth_attempts
ALTER TABLE admin_auth_attempts ADD COLUMN IF NOT EXISTS identifier character varying(150) NOT NULL;
ALTER TABLE admin_auth_attempts ADD COLUMN IF NOT EXISTS ip_address character varying(50) NOT NULL;
ALTER TABLE admin_auth_attempts ADD COLUMN IF NOT EXISTS portal character varying(50) NOT NULL;
ALTER TABLE admin_auth_attempts ADD COLUMN IF NOT EXISTS success boolean DEFAULT false NOT NULL;
ALTER TABLE admin_auth_attempts ADD COLUMN IF NOT EXISTS failure_reason character varying(100);
ALTER TABLE admin_auth_attempts ADD COLUMN IF NOT EXISTS attempted_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_auth_attempts' AND constraint_name = 'admin_auth_attempts_pkey'
    ) THEN
        ALTER TABLE ONLY admin_auth_attempts ADD CONSTRAINT admin_auth_attempts_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_admin_auth_attempts_ident ON admin_auth_attempts USING btree (identifier, attempted_at DESC);

CREATE INDEX IF NOT EXISTS idx_admin_auth_attempts_ip ON admin_auth_attempts USING btree (ip_address, attempted_at DESC);

CREATE INDEX IF NOT EXISTS idx_admin_auth_attempts_time ON admin_auth_attempts USING btree (attempted_at DESC);
