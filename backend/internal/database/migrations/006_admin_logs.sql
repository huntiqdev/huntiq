-- ============================================================================
-- 006_admin_logs.sql
-- Table: admin_logs (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for admin_logs
ALTER TABLE admin_logs ADD COLUMN IF NOT EXISTS admin_id integer DEFAULT 0 NOT NULL;
ALTER TABLE admin_logs ADD COLUMN IF NOT EXISTS action character varying(100) DEFAULT ''::character varying NOT NULL;
ALTER TABLE admin_logs ADD COLUMN IF NOT EXISTS target_type character varying(50);
ALTER TABLE admin_logs ADD COLUMN IF NOT EXISTS target_id integer;
ALTER TABLE admin_logs ADD COLUMN IF NOT EXISTS product_id integer;
ALTER TABLE admin_logs ADD COLUMN IF NOT EXISTS target_user_id integer;
ALTER TABLE admin_logs ADD COLUMN IF NOT EXISTS notes text;
ALTER TABLE admin_logs ADD COLUMN IF NOT EXISTS details text;
ALTER TABLE admin_logs ADD COLUMN IF NOT EXISTS ip_address character varying(50);
ALTER TABLE admin_logs ADD COLUMN IF NOT EXISTS user_agent text;
ALTER TABLE admin_logs ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_logs' AND constraint_name = 'admin_logs_pkey'
    ) THEN
        ALTER TABLE ONLY admin_logs ADD CONSTRAINT admin_logs_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_admin_logs_action ON admin_logs USING btree (action, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_admin_logs_admin ON admin_logs USING btree (admin_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_admin_logs_product ON admin_logs USING btree (product_id) WHERE (product_id IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_admin_logs_target ON admin_logs USING btree (target_type, target_id);

CREATE INDEX IF NOT EXISTS idx_admin_logs_target_user ON admin_logs USING btree (target_user_id) WHERE (target_user_id IS NOT NULL);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'admin_logs' AND constraint_name = 'admin_logs_admin_id_fkey'
    ) THEN
        ALTER TABLE ONLY admin_logs ADD CONSTRAINT admin_logs_admin_id_fkey FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
