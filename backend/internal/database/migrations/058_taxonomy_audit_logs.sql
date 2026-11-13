-- ============================================================================
-- 058_taxonomy_audit_logs.sql
-- Table: taxonomy_audit_logs (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for taxonomy_audit_logs
ALTER TABLE taxonomy_audit_logs ADD COLUMN IF NOT EXISTS admin_id integer;
ALTER TABLE taxonomy_audit_logs ADD COLUMN IF NOT EXISTS category_id integer;
ALTER TABLE taxonomy_audit_logs ADD COLUMN IF NOT EXISTS action character varying(50) NOT NULL;
ALTER TABLE taxonomy_audit_logs ADD COLUMN IF NOT EXISTS old_values jsonb;
ALTER TABLE taxonomy_audit_logs ADD COLUMN IF NOT EXISTS new_values jsonb;
ALTER TABLE taxonomy_audit_logs ADD COLUMN IF NOT EXISTS reason text;
ALTER TABLE taxonomy_audit_logs ADD COLUMN IF NOT EXISTS ip_address character varying(45);
ALTER TABLE taxonomy_audit_logs ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'taxonomy_audit_logs' AND constraint_name = 'taxonomy_audit_logs_pkey'
    ) THEN
        ALTER TABLE ONLY taxonomy_audit_logs ADD CONSTRAINT taxonomy_audit_logs_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_taxonomy_audit_category ON taxonomy_audit_logs USING btree (category_id);

CREATE INDEX IF NOT EXISTS idx_taxonomy_audit_created ON taxonomy_audit_logs USING btree (created_at);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'taxonomy_audit_logs' AND constraint_name = 'taxonomy_audit_logs_admin_id_fkey'
    ) THEN
        ALTER TABLE ONLY taxonomy_audit_logs ADD CONSTRAINT taxonomy_audit_logs_admin_id_fkey FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'taxonomy_audit_logs' AND constraint_name = 'taxonomy_audit_logs_category_id_fkey'
    ) THEN
        ALTER TABLE ONLY taxonomy_audit_logs ADD CONSTRAINT taxonomy_audit_logs_category_id_fkey FOREIGN KEY (category_id) REFERENCES category_hierarchy(id) ON DELETE SET NULL;
    END IF;
END $$;
