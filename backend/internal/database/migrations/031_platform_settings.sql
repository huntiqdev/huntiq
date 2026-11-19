-- ============================================================================
-- 031_platform_settings.sql
-- Table: platform_settings (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for platform_settings
ALTER TABLE platform_settings ADD COLUMN IF NOT EXISTS key character varying(100) UNIQUE;
ALTER TABLE platform_settings ADD COLUMN IF NOT EXISTS value text DEFAULT ''::text NOT NULL;
ALTER TABLE platform_settings ADD COLUMN IF NOT EXISTS value_type character varying(20) DEFAULT 'string'::character varying NOT NULL;
ALTER TABLE platform_settings ADD COLUMN IF NOT EXISTS category character varying(50) DEFAULT 'features'::character varying NOT NULL;
ALTER TABLE platform_settings ADD COLUMN IF NOT EXISTS description text;
ALTER TABLE platform_settings ADD COLUMN IF NOT EXISTS is_public boolean DEFAULT false NOT NULL;
ALTER TABLE platform_settings ADD COLUMN IF NOT EXISTS updated_by integer;
ALTER TABLE platform_settings ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE platform_settings ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'platform_settings' AND constraint_name = 'platform_settings_pkey'
    ) THEN
        ALTER TABLE ONLY platform_settings ADD CONSTRAINT platform_settings_pkey PRIMARY KEY (key);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_settings_category ON platform_settings USING btree (category);

CREATE INDEX IF NOT EXISTS idx_settings_public ON platform_settings USING btree (is_public) WHERE (is_public = true);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'platform_settings' AND constraint_name = 'platform_settings_updated_by_fkey'
    ) THEN
        ALTER TABLE ONLY platform_settings ADD CONSTRAINT platform_settings_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;
