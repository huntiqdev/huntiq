-- ============================================================================
-- 039_notifications.sql
-- Table: notifications (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for notifications
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS user_id integer DEFAULT 0 NOT NULL;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type character varying(50) DEFAULT 'system_message'::character varying NOT NULL;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS title character varying(200) DEFAULT ''::character varying NOT NULL;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS message text DEFAULT ''::text NOT NULL;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS link character varying(500);
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS is_read boolean DEFAULT false NOT NULL;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS is_admin boolean DEFAULT false NOT NULL;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS priority character varying(20) DEFAULT 'normal'::character varying;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'notifications' AND constraint_name = 'notifications_pkey'
    ) THEN
        ALTER TABLE ONLY notifications ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_notifications_admin ON notifications USING btree (is_admin, created_at DESC) WHERE (is_admin = true);

CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications USING btree (type);

CREATE INDEX IF NOT EXISTS idx_notifications_unread ON notifications USING btree (user_id, is_read) WHERE (is_read = false);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications USING btree (user_id, created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'notifications' AND constraint_name = 'notifications_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY notifications ADD CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
