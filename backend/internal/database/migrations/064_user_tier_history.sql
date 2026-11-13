-- ============================================================================
-- 064_user_tier_history.sql
-- Table: user_tier_history (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for user_tier_history
ALTER TABLE user_tier_history ADD COLUMN IF NOT EXISTS user_id integer DEFAULT 0 NOT NULL;
ALTER TABLE user_tier_history ADD COLUMN IF NOT EXISTS tier_type character varying(20) DEFAULT 'seller'::character varying NOT NULL;
ALTER TABLE user_tier_history ADD COLUMN IF NOT EXISTS old_tier character varying(20) DEFAULT 'normal'::character varying NOT NULL;
ALTER TABLE user_tier_history ADD COLUMN IF NOT EXISTS new_tier character varying(20) DEFAULT 'normal'::character varying NOT NULL;
ALTER TABLE user_tier_history ADD COLUMN IF NOT EXISTS reason character varying(100) DEFAULT 'system'::character varying NOT NULL;
ALTER TABLE user_tier_history ADD COLUMN IF NOT EXISTS changed_by integer;
ALTER TABLE user_tier_history ADD COLUMN IF NOT EXISTS upgrade_request_id integer;
ALTER TABLE user_tier_history ADD COLUMN IF NOT EXISTS duration_days integer;
ALTER TABLE user_tier_history ADD COLUMN IF NOT EXISTS expires_at timestamp with time zone;
ALTER TABLE user_tier_history ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_tier_history' AND constraint_name = 'user_tier_history_pkey'
    ) THEN
        ALTER TABLE ONLY user_tier_history ADD CONSTRAINT user_tier_history_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_tier_history_type ON user_tier_history USING btree (tier_type, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_tier_history_user ON user_tier_history USING btree (user_id, created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_tier_history' AND constraint_name = 'user_tier_history_changed_by_fkey'
    ) THEN
        ALTER TABLE ONLY user_tier_history ADD CONSTRAINT user_tier_history_changed_by_fkey FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_tier_history' AND constraint_name = 'user_tier_history_upgrade_request_id_fkey'
    ) THEN
        ALTER TABLE ONLY user_tier_history ADD CONSTRAINT user_tier_history_upgrade_request_id_fkey FOREIGN KEY (upgrade_request_id) REFERENCES upgrade_requests(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_tier_history' AND constraint_name = 'user_tier_history_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY user_tier_history ADD CONSTRAINT user_tier_history_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
