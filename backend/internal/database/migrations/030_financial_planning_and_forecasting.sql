-- ============================================================
-- 030_financial_planning_and_forecasting.sql
-- PURPOSE: Enterprise Financial Planning & Analysis (FP&A),
--          growth modeling, infrastructure capacity estimation,
--          Ethiopian statutory tax forecasting, and audit tracking.
-- ============================================================

-- 1. Financial Scenarios Table
CREATE TABLE IF NOT EXISTS financial_scenarios (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Financial Assumptions Table
CREATE TABLE IF NOT EXISTS financial_assumptions (
    id SERIAL PRIMARY KEY,
    scenario_id VARCHAR(50) NOT NULL REFERENCES financial_scenarios(id) ON DELETE CASCADE,
    section VARCHAR(50) NOT NULL,
    assumptions_json JSONB NOT NULL,
    updated_by INT REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_scenario_section UNIQUE (scenario_id, section)
);

CREATE INDEX IF NOT EXISTS idx_financial_assumptions_scenario ON financial_assumptions(scenario_id);

-- 3. Ethiopian Statutory Tax Rules Catalog
CREATE TABLE IF NOT EXISTS financial_tax_rules (
    id SERIAL PRIMARY KEY,
    jurisdiction VARCHAR(50) NOT NULL DEFAULT 'ET',
    rule_code VARCHAR(50) NOT NULL UNIQUE,
    rule_name VARCHAR(150) NOT NULL,
    category VARCHAR(50) NOT NULL,
    standard_rate NUMERIC(6,4),
    threshold_amount NUMERIC(15,2),
    effective_date DATE NOT NULL,
    source_citation TEXT NOT NULL,
    source_url TEXT,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    reviewer_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 4. Financial Actuals Table (Historical Accounting Benchmark)
CREATE TABLE IF NOT EXISTS financial_actuals (
    id SERIAL PRIMARY KEY,
    period_month VARCHAR(7) NOT NULL UNIQUE,
    currency VARCHAR(10) NOT NULL DEFAULT 'ETB',
    active_users INT NOT NULL DEFAULT 0,
    completed_transactions INT NOT NULL DEFAULT 0,
    gross_merchandise_volume NUMERIC(15,2) NOT NULL DEFAULT 0,
    commission_revenue NUMERIC(15,2) NOT NULL DEFAULT 0,
    listing_fee_revenue NUMERIC(15,2) NOT NULL DEFAULT 0,
    subscription_revenue NUMERIC(15,2) NOT NULL DEFAULT 0,
    other_revenue NUMERIC(15,2) NOT NULL DEFAULT 0,
    infra_cost_actual NUMERIC(15,2) NOT NULL DEFAULT 0,
    staffing_cost_actual NUMERIC(15,2) NOT NULL DEFAULT 0,
    marketing_cost_actual NUMERIC(15,2) NOT NULL DEFAULT 0,
    admin_cost_actual NUMERIC(15,2) NOT NULL DEFAULT 0,
    vat_collected_actual NUMERIC(15,2) NOT NULL DEFAULT 0,
    vat_paid_actual NUMERIC(15,2) NOT NULL DEFAULT 0,
    income_tax_paid_actual NUMERIC(15,2) NOT NULL DEFAULT 0,
    closing_cash_balance NUMERIC(15,2) NOT NULL DEFAULT 0,
    is_audited BOOLEAN NOT NULL DEFAULT FALSE,
    data_source VARCHAR(100) NOT NULL DEFAULT 'manual_entry',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_financial_actuals_period ON financial_actuals(period_month);

-- 5. Financial Audit Logs
CREATE TABLE IF NOT EXISTS financial_audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    scenario_id VARCHAR(50),
    action VARCHAR(100) NOT NULL,
    details JSONB NOT NULL,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_financial_audit_logs_created ON financial_audit_logs(created_at DESC);

-- 6. Seed Default Scenarios
INSERT INTO financial_scenarios (id, name, description)
VALUES 
    ('base', 'Base Case Planning', 'Realistic adoption growth, standard infrastructure deployment, benchmark monetization rates.'),
    ('conservative', 'Conservative Downside Case', 'Slower adoption, lower monetization conversion, higher churn, and protective cash runway modeling.'),
    ('high_growth', 'High Growth Upside Case', 'Aggressive marketplace expansion, rapid active user scaling, higher tier adoption, and capacity stress testing.')
ON CONFLICT (id) DO NOTHING;

-- 7. Seed Official Ethiopian Statutory Tax Rules
INSERT INTO financial_tax_rules (
    rule_code, rule_name, category, standard_rate, threshold_amount, effective_date, source_citation, source_url, is_verified, reviewer_notes
) VALUES
(
    'CIT_STANDARD',
    'Corporate Income Tax (Standard Rate)',
    'income_tax',
    0.3000,
    NULL,
    '2016-08-11',
    'Federal Income Tax Proclamation No. 979/2016, Article 19(1)',
    'https://www.mor.gov.et',
    FALSE,
    'Requires tax-professional confirmation. Standard 30% rate on taxable business income for body corporate.'
),
(
    'VAT_STANDARD',
    'Value Added Tax (Standard Rate)',
    'vat',
    0.1500,
    NULL,
    '2003-01-01',
    'Value Added Tax Proclamation No. 285/2002 & Proclamation No. 1341/2024, Article 7',
    'https://www.mor.gov.et',
    FALSE,
    'Requires tax-professional confirmation. 15% rate applies to taxable supplies of goods and digital services.'
),
(
    'VAT_THRESHOLD',
    'VAT Mandatory Registration Threshold',
    'vat',
    NULL,
    1000000.00,
    '2003-01-01',
    'Value Added Tax Proclamation No. 285/2002, Article 16',
    'https://www.mor.gov.et',
    FALSE,
    'Requires tax-professional confirmation. Annual taxable turnover of ETB 1,000,000 requires compulsory VAT registration.'
),
(
    'WHT_SERVICES',
    'Withholding Tax on Payments for Services',
    'withholding',
    0.0200,
    3000.00,
    '2016-08-11',
    'Federal Income Tax Proclamation No. 979/2016, Article 92(1)(a)',
    'https://www.mor.gov.et',
    FALSE,
    'Requires tax-professional confirmation. 2% withholding on domestic supply of services exceeding ETB 3,000 per transaction.'
),
(
    'WHT_GOODS',
    'Withholding Tax on Payments for Supply of Goods',
    'withholding',
    0.0200,
    10000.00,
    '2016-08-11',
    'Federal Income Tax Proclamation No. 979/2016, Article 92(1)(b)',
    'https://www.mor.gov.et',
    FALSE,
    'Requires tax-professional confirmation. 2% withholding on supply of goods in a single transaction exceeding ETB 10,000.'
),
(
    'PENSION_EMPLOYER',
    'Private Organization Employees Pension (Employer Share)',
    'payroll',
    0.1100,
    NULL,
    '2011-06-24',
    'Private Organization Employees Pension Proclamation No. 715/2011 & No. 1268/2022, Article 10',
    'https://www.poesa.gov.et',
    FALSE,
    'Requires tax-professional confirmation. Mandatory 11% employer contribution on employee basic salary.'
),
(
    'LOSS_CARRYFORWARD',
    'Business Loss Carryforward Relief',
    'deductions',
    NULL,
    5.00,
    '2016-08-11',
    'Federal Income Tax Proclamation No. 979/2016, Article 26',
    'https://www.mor.gov.et',
    FALSE,
    'Requires tax-professional confirmation. Tax losses may be carried forward for a maximum of 5 consecutive tax years.'
)
ON CONFLICT (rule_code) DO NOTHING;
