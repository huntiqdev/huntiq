-- ============================================================================
-- 065_user_warnings.sql
-- Table: user_warnings (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for user_warnings
ALTER TABLE user_warnings ADD COLUMN IF NOT EXISTS user_id integer;
ALTER TABLE user_warnings ADD COLUMN IF NOT EXISTS admin_id integer;
ALTER TABLE user_warnings ADD COLUMN IF NOT EXISTS reason text DEFAULT ''::text NOT NULL;
ALTER TABLE user_warnings ADD COLUMN IF NOT EXISTS severity character varying(20) DEFAULT 'low'::character varying;
ALTER TABLE user_warnings ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;
ALTER TABLE user_warnings ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

-- Table Constraints for user_warnings
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_warnings' AND constraint_name = 'user_warnings_severity_check'
    ) THEN
        ALTER TABLE user_warnings ADD CONSTRAINT user_warnings_severity_check CHECK (((severity)::text = ANY ((ARRAY['low'::character varying, 'medium'::character varying, 'high'::character varying, 'critical'::character varying])::text[])));
    END IF;
END $$;

COMMENT ON TABLE user_warnings IS 'Admin warnings issued to users for policy violations';

COMMENT ON COLUMN user_warnings.severity IS 'Warning severity: low, medium, high, critical';

COMMENT ON COLUMN user_warnings.is_active IS 'Whether warning is still active (not expired or resolved)';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_warnings' AND constraint_name = 'user_warnings_pkey'
    ) THEN
        ALTER TABLE ONLY user_warnings ADD CONSTRAINT user_warnings_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_user_warnings_active ON user_warnings USING btree (is_active, user_id);

CREATE INDEX IF NOT EXISTS idx_user_warnings_admin_id ON user_warnings USING btree (admin_id);

CREATE INDEX IF NOT EXISTS idx_user_warnings_severity ON user_warnings USING btree (severity) WHERE (is_active = true);

CREATE INDEX IF NOT EXISTS idx_user_warnings_user_id ON user_warnings USING btree (user_id);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_warnings' AND constraint_name = 'user_warnings_admin_id_fkey'
    ) THEN
        ALTER TABLE ONLY user_warnings ADD CONSTRAINT user_warnings_admin_id_fkey FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'user_warnings' AND constraint_name = 'user_warnings_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY user_warnings ADD CONSTRAINT user_warnings_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
