-- ============================================================================
-- 031_financial_scenarios.sql
-- Table: financial_scenarios (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for financial_scenarios
ALTER TABLE financial_scenarios ADD COLUMN IF NOT EXISTS name character varying(100) NOT NULL;
ALTER TABLE financial_scenarios ADD COLUMN IF NOT EXISTS description text;
ALTER TABLE financial_scenarios ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true NOT NULL;
ALTER TABLE financial_scenarios ADD COLUMN IF NOT EXISTS created_by integer;
ALTER TABLE financial_scenarios ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE financial_scenarios ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_scenarios' AND constraint_name = 'financial_scenarios_pkey'
    ) THEN
        ALTER TABLE ONLY financial_scenarios ADD CONSTRAINT financial_scenarios_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_scenarios' AND constraint_name = 'financial_scenarios_created_by_fkey'
    ) THEN
        ALTER TABLE ONLY financial_scenarios ADD CONSTRAINT financial_scenarios_created_by_fkey FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;
