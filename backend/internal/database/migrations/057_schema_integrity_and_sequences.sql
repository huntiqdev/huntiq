-- ============================================================================
-- 057_schema_integrity_and_sequences.sql
-- Sequence Alignment, Foreign Keys & Integrity Guarantees
-- ============================================================================

BEGIN;

-- Synchronize sequences for seeded auto-increment primary keys
SELECT setval(pg_get_serial_sequence('category_hierarchy', 'id'), COALESCE(MAX(id), 1)) FROM category_hierarchy;
SELECT setval(pg_get_serial_sequence('locations', 'id'), COALESCE(MAX(id), 1)) FROM locations;
SELECT setval(pg_get_serial_sequence('location_hierarchy', 'id'), COALESCE(MAX(id), 1)) FROM location_hierarchy;
SELECT setval(pg_get_serial_sequence('users', 'id'), COALESCE(MAX(id), 1)) FROM users;
SELECT setval(pg_get_serial_sequence('products', 'id'), COALESCE(MAX(id), 1)) FROM products;
SELECT setval(pg_get_serial_sequence('product_images', 'id'), COALESCE(MAX(id), 1)) FROM product_images;
SELECT setval(pg_get_serial_sequence('product_attributes', 'id'), COALESCE(MAX(id), 1)) FROM product_attributes;

-- Re-assert category_id synchronization with category_hierarchy_id
UPDATE products SET category_id = category_hierarchy_id WHERE category_hierarchy_id IS NOT NULL;

-- Category hierarchy self-referential parent_id constraint
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_hierarchy' AND constraint_name = 'category_hierarchy_parent_id_fkey'
    ) THEN
        ALTER TABLE ONLY category_hierarchy ADD CONSTRAINT category_hierarchy_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES category_hierarchy(id) ON DELETE RESTRICT;
    END IF;
END $$;

-- ============================================================================
-- Enterprise SMS, Exact-Token Referral Rewards, Financial Ledger & Telebirr
-- ============================================================================

-- 1. SMS Provider Configurations
CREATE TABLE IF NOT EXISTS sms_provider_configs (
    id SERIAL PRIMARY KEY,
    provider_code VARCHAR(50) NOT NULL UNIQUE,
    provider_name VARCHAR(100) NOT NULL,
    api_endpoint VARCHAR(255) NOT NULL,
    sender_id VARCHAR(50) NOT NULL DEFAULT 'HUNTIQ',
    credentials JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT false,
    is_default BOOLEAN NOT NULL DEFAULT false,
    daily_limit INTEGER NOT NULL DEFAULT 50000,
    daily_sent INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO sms_provider_configs (provider_code, provider_name, api_endpoint, sender_id, is_active, is_default)
VALUES 
  ('mock_dev', 'HUNTIQ Dev Logger & Mock Gateway', 'internal://mock-sms', 'HUNTIQ-DEV', true, true),
  ('ethio_telecom', 'Ethio Telecom Enterprise SMS Gateway', 'https://api.ethiotelecom.et/sms/v1/send', 'HUNTIQ', false, false)
ON CONFLICT (provider_code) DO NOTHING;

-- 2. SMS Delivery Logs
CREATE TABLE IF NOT EXISTS sms_delivery_logs (
    id BIGSERIAL PRIMARY KEY,
    phone_masked VARCHAR(20) NOT NULL,
    purpose VARCHAR(50) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    provider_ref VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'queued',
    error_category VARCHAR(50),
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sms_logs_status ON sms_delivery_logs(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sms_logs_phone ON sms_delivery_logs(phone_masked, created_at DESC);

-- 3. Referral Tokens
CREATE TABLE IF NOT EXISTS referral_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(64) NOT NULL UNIQUE,
    referral_code VARCHAR(11) NOT NULL,
    click_count INTEGER NOT NULL DEFAULT 0,
    conversion_count INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_referral_tokens_user ON referral_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_referral_tokens_token ON referral_tokens(token);

-- 4. Referral Attributions
CREATE TABLE IF NOT EXISTS referral_attributions (
    id BIGSERIAL PRIMARY KEY,
    token_id INTEGER NOT NULL REFERENCES referral_tokens(id) ON DELETE CASCADE,
    inviter_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_hash VARCHAR(64) NOT NULL,
    user_agent_hash VARCHAR(64),
    landing_url VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'clicked',
    registered_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    attributed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '30 days')
);
CREATE INDEX IF NOT EXISTS idx_referral_attributions_token ON referral_attributions(token_id);
CREATE INDEX IF NOT EXISTS idx_referral_attributions_user ON referral_attributions(registered_user_id) WHERE registered_user_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_referral_attributions_ip ON referral_attributions(ip_hash, attributed_at DESC);

-- 5. Referral Reward Policies
CREATE TABLE IF NOT EXISTS referral_reward_policies (
    id SERIAL PRIMARY KEY,
    policy_name VARCHAR(100) NOT NULL,
    registration_reward_amount NUMERIC(12,2) NOT NULL DEFAULT 1.00,
    transaction_reward_percent NUMERIC(5,2) NOT NULL DEFAULT 5.00,
    min_qualifying_tx_amount NUMERIC(12,2) NOT NULL DEFAULT 50.00,
    max_reward_per_tx NUMERIC(12,2) NOT NULL DEFAULT 5000.00,
    max_lifetime_per_referred NUMERIC(12,2) NOT NULL DEFAULT 25000.00,
    max_rewards_per_inviter NUMERIC(12,2) NOT NULL DEFAULT 100000.00,
    reward_holding_days INTEGER NOT NULL DEFAULT 14,
    is_active BOOLEAN NOT NULL DEFAULT true,
    version INTEGER NOT NULL DEFAULT 1,
    effective_from TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    effective_to TIMESTAMPTZ,
    updated_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO referral_reward_policies (
    policy_name, registration_reward_amount, transaction_reward_percent,
    min_qualifying_tx_amount, max_reward_per_tx, reward_holding_days, is_active, version
) VALUES (
    'Standard Canonical Direct-Referral Policy 2026', 1.00, 5.00,
    50.00, 5000.00, 14, true, 1
) ON CONFLICT DO NOTHING;

-- 6. Referral Reward Policy Versions
CREATE TABLE IF NOT EXISTS referral_reward_policy_versions (
    id SERIAL PRIMARY KEY,
    policy_id INTEGER NOT NULL REFERENCES referral_reward_policies(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    registration_reward_amount NUMERIC(12,2) NOT NULL,
    transaction_reward_percent NUMERIC(5,2) NOT NULL,
    min_qualifying_tx_amount NUMERIC(12,2) NOT NULL,
    max_reward_per_tx NUMERIC(12,2) NOT NULL,
    reward_holding_days INTEGER NOT NULL,
    changed_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    change_reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 7. Wallet Accounts
CREATE TABLE IF NOT EXISTS wallet_accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    currency VARCHAR(10) NOT NULL DEFAULT 'ETB',
    available_balance NUMERIC(15,2) NOT NULL DEFAULT 0.00 CHECK (available_balance >= 0),
    pending_balance NUMERIC(15,2) NOT NULL DEFAULT 0.00 CHECK (pending_balance >= 0),
    reserved_balance NUMERIC(15,2) NOT NULL DEFAULT 0.00 CHECK (reserved_balance >= 0),
    total_withdrawn NUMERIC(15,2) NOT NULL DEFAULT 0.00 CHECK (total_withdrawn >= 0),
    lifetime_earned NUMERIC(15,2) NOT NULL DEFAULT 0.00 CHECK (lifetime_earned >= 0),
    is_locked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_wallet_accounts_user ON wallet_accounts(user_id);

-- 8. Immutable Double-Entry Ledger Transactions
CREATE TABLE IF NOT EXISTS ledger_transactions (
    id BIGSERIAL PRIMARY KEY,
    ledger_tx_id VARCHAR(64) NOT NULL UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    debit_account VARCHAR(50) NOT NULL,
    credit_account VARCHAR(50) NOT NULL,
    amount NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'ETB',
    entry_type VARCHAR(50) NOT NULL,
    reference_id VARCHAR(100),
    source_event VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'posted',
    idempotency_key VARCHAR(128) NOT NULL UNIQUE,
    notes TEXT,
    actor_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ledger_tx_user ON ledger_transactions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ledger_tx_entry_type ON ledger_transactions(entry_type);
CREATE INDEX IF NOT EXISTS idx_ledger_tx_idempotency ON ledger_transactions(idempotency_key);

-- 9. Telebirr Payment Orders
CREATE TABLE IF NOT EXISTS telebirr_orders (
    id SERIAL PRIMARY KEY,
    out_trade_no VARCHAR(64) NOT NULL UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(12,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'ETB',
    title VARCHAR(150) NOT NULL,
    order_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    prepay_id VARCHAR(128),
    transaction_no VARCHAR(128),
    payer_phone VARCHAR(20),
    raw_callback JSONB,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_telebirr_orders_out_trade ON telebirr_orders(out_trade_no);
CREATE INDEX IF NOT EXISTS idx_telebirr_orders_user ON telebirr_orders(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_telebirr_orders_status ON telebirr_orders(status);

-- 10. Referral Fraud Evaluations
CREATE TABLE IF NOT EXISTS referral_fraud_evaluations (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    inviter_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    risk_score INTEGER NOT NULL DEFAULT 0,
    decision VARCHAR(30) NOT NULL DEFAULT 'approved',
    signals JSONB NOT NULL DEFAULT '{}'::jsonb,
    investigated_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    investigation_notes TEXT,
    evaluated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_fraud_eval_inviter ON referral_fraud_evaluations(inviter_id);
CREATE INDEX IF NOT EXISTS idx_fraud_eval_decision ON referral_fraud_evaluations(decision);

COMMIT;
