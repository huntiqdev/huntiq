-- ============================================================================
-- 002_abuse_reports.sql
-- Table: abuse_reports (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for abuse_reports
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS reporter_id integer;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS reporter_role character varying(20);
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS reported_user_id integer;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS reported_product_id integer;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS reason character varying(50) NOT NULL;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS description text NOT NULL;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS evidence_urls text[];
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'pending'::character varying NOT NULL;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS priority character varying(20) DEFAULT 'normal'::character varying NOT NULL;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS reviewed_by integer;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS reviewed_at timestamp with time zone;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS action_taken text;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS resolution_notes text;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

COMMENT ON TABLE abuse_reports IS 'User-reported abuse, fraud, and safety incidents';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'abuse_reports' AND constraint_name = 'abuse_reports_pkey'
    ) THEN
        ALTER TABLE ONLY abuse_reports ADD CONSTRAINT abuse_reports_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_abuse_reports_reported_product ON abuse_reports USING btree (reported_product_id) WHERE (reported_product_id IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_abuse_reports_reported_user ON abuse_reports USING btree (reported_user_id) WHERE (reported_user_id IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_abuse_reports_reporter ON abuse_reports USING btree (reporter_id) WHERE (reporter_id IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_abuse_reports_reviewer ON abuse_reports USING btree (reviewed_by) WHERE (reviewed_by IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_abuse_reports_status ON abuse_reports USING btree (status, priority DESC, created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'abuse_reports' AND constraint_name = 'abuse_reports_reported_product_id_fkey'
    ) THEN
        ALTER TABLE ONLY abuse_reports ADD CONSTRAINT abuse_reports_reported_product_id_fkey FOREIGN KEY (reported_product_id) REFERENCES products(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'abuse_reports' AND constraint_name = 'abuse_reports_reported_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY abuse_reports ADD CONSTRAINT abuse_reports_reported_user_id_fkey FOREIGN KEY (reported_user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'abuse_reports' AND constraint_name = 'abuse_reports_reporter_id_fkey'
    ) THEN
        ALTER TABLE ONLY abuse_reports ADD CONSTRAINT abuse_reports_reporter_id_fkey FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'abuse_reports' AND constraint_name = 'abuse_reports_reviewed_by_fkey'
    ) THEN
        ALTER TABLE ONLY abuse_reports ADD CONSTRAINT abuse_reports_reviewed_by_fkey FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;
