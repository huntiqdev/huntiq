-- ============================================================================
-- 028_notifications.sql
-- Table: admin_notifications (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for admin_notifications
ALTER TABLE admin_notifications ADD COLUMN IF NOT EXISTS title character varying(200) DEFAULT ''::character varying NOT NULL;
ALTER TABLE admin_notifications ADD COLUMN IF NOT EXISTS message text DEFAULT ''::text NOT NULL;
ALTER TABLE admin_notifications ADD COLUMN IF NOT EXISTS priority character varying(20) DEFAULT 'normal'::character varying;
ALTER TABLE admin_notifications ADD COLUMN IF NOT EXISTS target_role character varying(50);
ALTER TABLE admin_notifications ADD COLUMN IF NOT EXISTS is_read boolean DEFAULT false NOT NULL;
ALTER TABLE admin_notifications ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE admin_notifications ADD COLUMN IF NOT EXISTS type character varying(50) DEFAULT 'info'::character varying NOT NULL;
ALTER TABLE admin_notifications ADD COLUMN IF NOT EXISTS read_by integer[] DEFAULT '{}'::integer[] NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_notifications' AND constraint_name = 'admin_notifications_pkey'
    ) THEN
        ALTER TABLE ONLY admin_notifications ADD CONSTRAINT admin_notifications_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_admin_notifications_priority ON admin_notifications USING btree (priority, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_admin_notifications_read_by ON admin_notifications USING gin (read_by);

CREATE INDEX IF NOT EXISTS idx_admin_notifications_type ON admin_notifications USING btree (type);
