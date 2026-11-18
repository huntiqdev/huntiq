-- ============================================================================
-- 044_user_activity_logs.sql
-- Table: user_activity_logs (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for user_activity_logs
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS user_id integer;
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS event character varying(50) DEFAULT ''::character varying NOT NULL;
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS phone character varying(20);
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS email character varying(150);
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS user_type character varying(20);
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS ip_address character varying(50);
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS user_agent text;
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS success boolean DEFAULT true NOT NULL;
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS failure_reason text;
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS metadata jsonb;
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS meta text;
ALTER TABLE user_activity_logs ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_activity_logs' AND constraint_name = 'user_activity_logs_pkey'
    ) THEN
        ALTER TABLE ONLY user_activity_logs ADD CONSTRAINT user_activity_logs_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_activity_logs_event ON user_activity_logs USING btree (event, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_activity_logs_ip ON user_activity_logs USING btree (ip_address, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_activity_logs_phone ON user_activity_logs USING btree (phone, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_activity_logs_user ON user_activity_logs USING btree (user_id, created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_activity_logs' AND constraint_name = 'user_activity_logs_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY user_activity_logs ADD CONSTRAINT user_activity_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;
