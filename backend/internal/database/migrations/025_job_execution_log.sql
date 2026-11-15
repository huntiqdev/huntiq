-- ============================================================================
-- 025_job_execution_log.sql
-- Table: job_execution_log (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for job_execution_log
ALTER TABLE job_execution_log ADD COLUMN IF NOT EXISTS job_name character varying(100) NOT NULL;
ALTER TABLE job_execution_log ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'running'::character varying NOT NULL;
ALTER TABLE job_execution_log ADD COLUMN IF NOT EXISTS started_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE job_execution_log ADD COLUMN IF NOT EXISTS completed_at timestamp with time zone;
ALTER TABLE job_execution_log ADD COLUMN IF NOT EXISTS duration_ms integer;
ALTER TABLE job_execution_log ADD COLUMN IF NOT EXISTS records_processed integer DEFAULT 0;
ALTER TABLE job_execution_log ADD COLUMN IF NOT EXISTS error_message text;
ALTER TABLE job_execution_log ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

COMMENT ON TABLE job_execution_log IS 'Background job execution history for operations monitoring';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'job_execution_log' AND constraint_name = 'job_execution_log_pkey'
    ) THEN
        ALTER TABLE ONLY job_execution_log ADD CONSTRAINT job_execution_log_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_job_log_name_status ON job_execution_log USING btree (job_name, status, started_at DESC);

CREATE INDEX IF NOT EXISTS idx_job_log_started ON job_execution_log USING btree (started_at DESC);
