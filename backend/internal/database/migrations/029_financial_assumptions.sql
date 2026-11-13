-- ============================================================================
-- 029_financial_assumptions.sql
-- Table: financial_assumptions (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for financial_assumptions
ALTER TABLE financial_assumptions ADD COLUMN IF NOT EXISTS scenario_id character varying(50) NOT NULL;
ALTER TABLE financial_assumptions ADD COLUMN IF NOT EXISTS section character varying(50) NOT NULL;
ALTER TABLE financial_assumptions ADD COLUMN IF NOT EXISTS assumptions_json jsonb NOT NULL;
ALTER TABLE financial_assumptions ADD COLUMN IF NOT EXISTS updated_by integer;
ALTER TABLE financial_assumptions ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_assumptions' AND constraint_name = 'financial_assumptions_pkey'
    ) THEN
        ALTER TABLE ONLY financial_assumptions ADD CONSTRAINT financial_assumptions_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_assumptions' AND constraint_name = 'uq_scenario_section'
    ) THEN
        ALTER TABLE ONLY financial_assumptions ADD CONSTRAINT uq_scenario_section UNIQUE (scenario_id, section);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_financial_assumptions_scenario ON financial_assumptions USING btree (scenario_id);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_assumptions' AND constraint_name = 'financial_assumptions_scenario_id_fkey'
    ) THEN
        ALTER TABLE ONLY financial_assumptions ADD CONSTRAINT financial_assumptions_scenario_id_fkey FOREIGN KEY (scenario_id) REFERENCES financial_scenarios(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_assumptions' AND constraint_name = 'financial_assumptions_updated_by_fkey'
    ) THEN
        ALTER TABLE ONLY financial_assumptions ADD CONSTRAINT financial_assumptions_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;
