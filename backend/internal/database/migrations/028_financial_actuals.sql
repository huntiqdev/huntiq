-- ============================================================================
-- 028_financial_actuals.sql
-- Table: financial_actuals (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for financial_actuals
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS period_month character varying(7) NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS currency character varying(10) DEFAULT 'ETB'::character varying NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS active_users integer DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS completed_transactions integer DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS gross_merchandise_volume numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS commission_revenue numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS listing_fee_revenue numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS subscription_revenue numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS other_revenue numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS infra_cost_actual numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS staffing_cost_actual numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS marketing_cost_actual numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS admin_cost_actual numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS vat_collected_actual numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS vat_paid_actual numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS income_tax_paid_actual numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS closing_cash_balance numeric(15,2) DEFAULT 0 NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS is_audited boolean DEFAULT false NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS data_source character varying(100) DEFAULT 'manual_entry'::character varying NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS notes text;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE financial_actuals ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_actuals' AND constraint_name = 'financial_actuals_period_month_key'
    ) THEN
        ALTER TABLE ONLY financial_actuals ADD CONSTRAINT financial_actuals_period_month_key UNIQUE (period_month);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_actuals' AND constraint_name = 'financial_actuals_pkey'
    ) THEN
        ALTER TABLE ONLY financial_actuals ADD CONSTRAINT financial_actuals_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_financial_actuals_period ON financial_actuals USING btree (period_month);
