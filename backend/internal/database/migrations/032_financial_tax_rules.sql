-- ============================================================================
-- 032_financial_tax_rules.sql
-- Table: financial_tax_rules (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for financial_tax_rules
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS jurisdiction character varying(50) DEFAULT 'ET'::character varying NOT NULL;
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS rule_code character varying(50) NOT NULL;
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS rule_name character varying(150) NOT NULL;
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS category character varying(50) NOT NULL;
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS standard_rate numeric(6,4);
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS threshold_amount numeric(15,2);
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS effective_date date NOT NULL;
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS source_citation text NOT NULL;
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS source_url text;
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS is_verified boolean DEFAULT false NOT NULL;
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS reviewer_notes text;
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE financial_tax_rules ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_tax_rules' AND constraint_name = 'financial_tax_rules_pkey'
    ) THEN
        ALTER TABLE ONLY financial_tax_rules ADD CONSTRAINT financial_tax_rules_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'financial_tax_rules' AND constraint_name = 'financial_tax_rules_rule_code_key'
    ) THEN
        ALTER TABLE ONLY financial_tax_rules ADD CONSTRAINT financial_tax_rules_rule_code_key UNIQUE (rule_code);
    END IF;
END $$;
