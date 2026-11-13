-- ============================================================================
-- 030_financial_audit_logs.sql
-- Table: financial_audit_logs (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for financial_audit_logs
ALTER TABLE financial_audit_logs ADD COLUMN IF NOT EXISTS user_id integer;
ALTER TABLE financial_audit_logs ADD COLUMN IF NOT EXISTS scenario_id character varying(50);
ALTER TABLE financial_audit_logs ADD COLUMN IF NOT EXISTS action character varying(100) NOT NULL;
ALTER TABLE financial_audit_logs ADD COLUMN IF NOT EXISTS details jsonb NOT NULL;
ALTER TABLE financial_audit_logs ADD COLUMN IF NOT EXISTS ip_address character varying(45);
ALTER TABLE financial_audit_logs ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_audit_logs' AND constraint_name = 'financial_audit_logs_pkey'
    ) THEN
        ALTER TABLE ONLY financial_audit_logs ADD CONSTRAINT financial_audit_logs_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_financial_audit_logs_created ON financial_audit_logs USING btree (created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_audit_logs' AND constraint_name = 'financial_audit_logs_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY financial_audit_logs ADD CONSTRAINT financial_audit_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;
