-- ============================================================================
-- 063_unified_payment_and_fida_architecture.sql
-- Table: payment_orders, payment_provider_configs, and Fida SMS integration
-- Master FinTech Unified Payment Architecture (HUNTIQ Enterprise)
-- ============================================================================

-- 1. Payment Provider Configurations Table
CREATE TABLE IF NOT EXISTS payment_provider_configs (
    id SERIAL PRIMARY KEY,
    provider_code VARCHAR(50) NOT NULL UNIQUE,
    provider_name VARCHAR(100) NOT NULL,
    environment_mode VARCHAR(20) NOT NULL DEFAULT 'demo' CHECK (environment_mode IN ('demo', 'production')),
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    api_endpoint VARCHAR(255) NOT NULL DEFAULT '',
    app_id VARCHAR(100) NOT NULL DEFAULT '',
    app_key_encrypted TEXT NOT NULL DEFAULT '',
    short_code VARCHAR(50) NOT NULL DEFAULT '',
    notify_url VARCHAR(255) NOT NULL DEFAULT '',
    return_url VARCHAR(255) NOT NULL DEFAULT '',
    extra_config JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed Initial Payment Providers
INSERT INTO payment_provider_configs 
    (provider_code, provider_name, environment_mode, is_enabled, api_endpoint, short_code)
VALUES 
    ('demo_mock', 'HUNTIQ FinTech Demo Simulator', 'demo', true, 'internal://demo-payment', 'DEMO01'),
    ('telebirr', 'Ethio Telecom Telebirr SuperApp', 'demo', true, 'https://app.telebirr.et/payment', '100200'),
    ('ethio_telecom', 'Ethio Telecom TelePay / Direct Carrier Billing', 'demo', true, 'https://api.ethiotelecom.et/telepay/v1', 'ETB01')
ON CONFLICT (provider_code) DO NOTHING;

-- 2. Unified Payment Orders Table
CREATE TABLE IF NOT EXISTS payment_orders (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    environment VARCHAR(20) NOT NULL DEFAULT 'demo' CHECK (environment IN ('demo', 'production')),
    order_type VARCHAR(50) NOT NULL CHECK (order_type IN ('tier_upgrade', 'listing_fee', 'membership')),
    target_tier VARCHAR(20),
    user_role VARCHAR(20) NOT NULL DEFAULT 'seller',
    duration_days INTEGER NOT NULL DEFAULT 30,
    product_id INTEGER REFERENCES products(id) ON DELETE SET NULL,
    amount NUMERIC(12,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'ETB',
    title VARCHAR(200) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'CREATED' 
        CHECK (status IN ('CREATED', 'PENDING', 'PROCESSING', 'SUCCESSFUL', 'FAILED', 'CANCELLED', 'EXPIRED')),
    payer_phone VARCHAR(30),
    provider_tx_id VARCHAR(128),
    checkout_url TEXT,
    raw_callback JSONB,
    verification_metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_demo BOOLEAN NOT NULL DEFAULT true,
    otp_verified BOOLEAN NOT NULL DEFAULT false,
    otp_verified_at TIMESTAMPTZ,
    verified_at TIMESTAMPTZ,
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_orders_order_id ON payment_orders(order_id);
CREATE INDEX IF NOT EXISTS idx_payment_orders_user ON payment_orders(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_payment_orders_status ON payment_orders(status);
CREATE INDEX IF NOT EXISTS idx_payment_orders_provider ON payment_orders(provider);
CREATE INDEX IF NOT EXISTS idx_payment_orders_env ON payment_orders(environment, is_demo);
CREATE INDEX IF NOT EXISTS idx_payment_orders_created ON payment_orders(created_at DESC);

-- 3. Register Fida SMS Provider in sms_provider_configs
INSERT INTO sms_provider_configs 
    (provider_code, provider_name, api_endpoint, sender_id, is_active, is_default)
VALUES 
    ('fida', 'Fida National SMS & OTP Delivery Gateway', 'https://api.fida.et/v1/sms/otp', 'HUNTIQ', false, false)
ON CONFLICT (provider_code) DO NOTHING;

-- 4. Authoritative Platform Settings for Tier Upgrades & System Mode
INSERT INTO platform_settings (key, value, description)
VALUES 
    ('payment_operating_mode', 'demo', 'System-wide payment operating mode: demo or production'),
    ('fida_otp_required_for_upgrade', 'true', 'Require Fida SMS OTP verification prior to membership upgrade'),
    ('seller_gold_price_etb', '500', 'Authoritative monthly price for Seller Gold membership'),
    ('seller_premium_price_etb', '2000', 'Authoritative monthly price for Seller Premium membership'),
    ('buyer_gold_price_etb', '200', 'Authoritative monthly price for Buyer Gold membership'),
    ('buyer_premium_price_etb', '500', 'Authoritative monthly price for Buyer Premium membership')
ON CONFLICT (key) DO NOTHING;
