-- ============================================================================
-- 056_settings_change_log.sql
-- Table: settings_change_log (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for settings_change_log
ALTER TABLE settings_change_log ADD COLUMN IF NOT EXISTS setting_key character varying(100) DEFAULT ''::character varying NOT NULL;
ALTER TABLE settings_change_log ADD COLUMN IF NOT EXISTS old_value text;
ALTER TABLE settings_change_log ADD COLUMN IF NOT EXISTS new_value text DEFAULT ''::text NOT NULL;
ALTER TABLE settings_change_log ADD COLUMN IF NOT EXISTS changed_by integer DEFAULT 0 NOT NULL;
ALTER TABLE settings_change_log ADD COLUMN IF NOT EXISTS change_reason text;
ALTER TABLE settings_change_log ADD COLUMN IF NOT EXISTS ip_address character varying(50);
ALTER TABLE settings_change_log ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'settings_change_log' AND constraint_name = 'settings_change_log_pkey'
    ) THEN
        ALTER TABLE ONLY settings_change_log ADD CONSTRAINT settings_change_log_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_settings_log_admin ON settings_change_log USING btree (changed_by, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_settings_log_key ON settings_change_log USING btree (setting_key, created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'settings_change_log' AND constraint_name = 'settings_change_log_changed_by_fkey'
    ) THEN
        ALTER TABLE ONLY settings_change_log ADD CONSTRAINT settings_change_log_changed_by_fkey FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
