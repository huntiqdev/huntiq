-- ============================================================
-- 011_settings_and_limits.sql
-- CONSOLIDATED: platform_settings + settings_change_log
-- PURPOSE: All platform configuration in one place
-- NOTE: Uses ALTER TABLE for settings_change_log. For platform_settings,
--        001_schema created it with (key VARCHAR(100) PRIMARY KEY),
--        so we ALTER to add the additional columns.
-- ============================================================

-- ── Platform Settings ───────────────────────────────────────────────────────
-- 001_schema.sql created: (key VARCHAR(100) PRIMARY KEY)
-- We need to add all the other columns
ALTER TABLE platform_settings
    ADD COLUMN IF NOT EXISTS value       TEXT          NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS value_type  VARCHAR(20)   NOT NULL DEFAULT 'string',
    ADD COLUMN IF NOT EXISTS category    VARCHAR(50)   NOT NULL DEFAULT 'features',
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS is_public   BOOLEAN       NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS updated_by  INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Settings Change Log ─────────────────────────────────────────────────────
ALTER TABLE settings_change_log
    ADD COLUMN IF NOT EXISTS setting_key   VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS old_value     TEXT,
    ADD COLUMN IF NOT EXISTS new_value     TEXT          NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS changed_by    INT           NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS change_reason TEXT,
    ADD COLUMN IF NOT EXISTS ip_address    VARCHAR(50),
    ADD COLUMN IF NOT EXISTS created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Indexes ─────────────────────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_settings_category    ON platform_settings(category);
CREATE INDEX IF NOT EXISTS idx_settings_public      ON platform_settings(is_public) WHERE is_public = TRUE;
CREATE INDEX IF NOT EXISTS idx_settings_log_key     ON settings_change_log(setting_key, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_settings_log_admin   ON settings_change_log(changed_by, created_at DESC);

-- ── Trigger to log all setting changes ──────────────────────────────────────
CREATE OR REPLACE FUNCTION log_setting_change()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' AND OLD.value IS DISTINCT FROM NEW.value THEN
        INSERT INTO settings_change_log (setting_key, old_value, new_value, changed_by)
        VALUES (NEW.key, OLD.value, NEW.value, NEW.updated_by);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS track_setting_changes ON platform_settings;
CREATE TRIGGER track_setting_changes AFTER UPDATE ON platform_settings
    FOR EACH ROW EXECUTE FUNCTION log_setting_change();

-- ── Insert default settings ─────────────────────────────────────────────────
INSERT INTO platform_settings (key, value, value_type, category, description, is_public) VALUES
    -- Fee configuration
    ('broker_fee_rate', '2.0', 'number', 'payment', 'Default broker fee rate percentage', true),
    ('referral_commission_rate', '1.2', 'number', 'payment', 'Referral commission rate percentage', true),
    ('min_referral_withdrawal', '500', 'number', 'payment', 'Minimum referral commission withdrawal amount in ETB', true),
    ('max_referral_withdrawal', '100000', 'number', 'payment', 'Maximum referral commission withdrawal amount in ETB', true),
    ('normal_startup_fee', '1.0', 'number', 'payment', 'Listing startup fee percentage for normal sellers', true),
    ('gold_startup_fee', '0.5', 'number', 'payment', 'Listing startup fee percentage for gold sellers', true),
    ('premium_startup_fee', '0.25', 'number', 'payment', 'Listing startup fee percentage for premium sellers', true),
    ('normal_weekly_increment', '0.5', 'number', 'payment', 'Weekly increment percentage for normal sellers', true),
    ('gold_weekly_increment', '0.25', 'number', 'payment', 'Weekly increment percentage for gold sellers', true),
    ('premium_weekly_increment', '0.1', 'number', 'payment', 'Weekly increment percentage for premium sellers', true),
    ('listing_free_weeks', '1', 'number', 'payment', 'Number of free weeks before weekly increment is charged', true),

    -- Image limits
    ('normal_max_images', '2', 'number', 'limits', 'Max images allowed for normal tier listings', true),
    ('gold_max_images', '7', 'number', 'limits', 'Max images allowed for gold tier listings', true),
    ('premium_max_images', '20', 'number', 'limits', 'Max images allowed for premium tier listings', true),
    ('min_images_per_listing', '2', 'number', 'limits', 'Minimum images required per product listing', true),

    -- Tier thresholds & costs
    ('buyer_gold_cost', '5900', 'number', 'payment', 'Buyer Gold tier upgrade cost in ETB', true),
    ('buyer_premium_cost', '15900', 'number', 'payment', 'Buyer Premium tier upgrade cost in ETB', true),
    ('gold_threshold', '5900', 'number', 'payment', 'Gold upgrade price/threshold in ETB', true),
    ('premium_threshold', '15900', 'number', 'payment', 'Premium upgrade price/threshold in ETB', true),

    -- Upgrade pricing
    ('buyer_gold_upgrade_price', '5900', 'number', 'payment', 'Buyer Gold tier upgrade price in ETB', true),
    ('buyer_premium_upgrade_price', '15900', 'number', 'payment', 'Buyer Premium tier upgrade price in ETB', true),
    ('gold_upgrade_price', '500', 'number', 'payment', 'Seller Gold tier upgrade price in ETB', true),
    ('premium_upgrade_price', '1500', 'number', 'payment', 'Seller Premium tier upgrade price in ETB', true),

    -- Payment merchant numbers
    ('telebirr_merchant_phone', '0912345678', 'string', 'payment', 'Merchant phone number for Telebirr payments', true),
    ('safaricom_merchant_phone', '0712345678', 'string', 'payment', 'Merchant phone number for Safaricom/M-Pesa payments', true),
    ('cbe_account_number', '1000123456789', 'string', 'payment', 'Commercial Bank of Ethiopia account number', true),

    -- Registration & upgrade controls
    ('allow_normal_registration', 'true', 'boolean', 'features', 'Allow new normal tier user registrations', false),
    ('allow_buyer_registration', 'true', 'boolean', 'features', 'Allow new buyer registrations', false),
    ('allow_seller_upgrade_to_gold', 'true', 'boolean', 'features', 'Allow sellers to request Gold upgrades', false),
    ('allow_seller_upgrade_to_premium', 'true', 'boolean', 'features', 'Allow sellers to request Premium upgrades', false),
    ('allow_buyer_tier_upgrade', 'true', 'boolean', 'features', 'Allow buyers to request tier upgrades', false),

    -- OTP & auth
    ('otp_required', 'true', 'boolean', 'security', 'Whether OTP verification is required for login/registration', false),
    ('otp_expiry_minutes', '15', 'number', 'security', 'OTP code expiration time in minutes', false),
    ('otp_gateway', 'log', 'string', 'security', 'Gateway to send OTP (e.g. gateway_name, log)', true),

    -- Feature flags
    ('maintenance_mode', 'false', 'boolean', 'features', 'System maintenance/read-only mode', false),
    ('telebirr_payment_enabled', 'true', 'boolean', 'features', 'Enable/disable Telebirr payment method', false),
    ('buyer_tier_upgrade_enabled', 'true', 'boolean', 'features', 'Enable/disable buyer tier upgrade requests', false),
    ('ai_recommendations_enabled', 'true', 'boolean', 'features', 'Enable/disable AI-powered listings recommendations', false),

    -- Additional site/contact info
    ('company_phone', '0911000000', 'string', 'company', 'Official company support phone number', true),
    ('company_email', 'support@habeshamarket.com', 'string', 'company', 'Official company support email address', true),
    ('company_location', 'Addis Ababa, Ethiopia', 'string', 'company', 'Physical location address of headquarters', true),

    -- System Limits
    ('max_title_length', '200', 'number', 'limits', 'Maximum character length of listing titles', true),
    ('max_description_length', '5000', 'number', 'limits', 'Maximum character length of listing descriptions', true),
    ('normal_max_locations', '1', 'number', 'limits', 'Max locations a normal seller can assign', true),
    ('gold_max_locations', '3', 'number', 'limits', 'Max locations a gold seller can assign', true),
    ('premium_max_locations', '10', 'number', 'limits', 'Max locations a premium seller can assign', true),
    ('max_banners_buyer', '5', 'number', 'limits', 'Max banner uploads allowed for buyers', true),
    ('max_banners_seller', '5', 'number', 'limits', 'Max banner uploads allowed for sellers', true),
    ('max_image_size_mb', '20', 'number', 'limits', 'Maximum allowed image upload size in megabytes', true),
    ('currency_tolerance_percent', '0.1', 'number', 'payment', 'Allowed currency drift/tolerance percentage', true),

    -- Rate limits
    ('max_otp_per_hour', '5', 'number', 'security', 'Maximum OTP requests per hour per phone', false),
    ('max_login_attempts', '10', 'number', 'security', 'Maximum login attempts per hour', false)
ON CONFLICT (key) DO UPDATE SET
    value = EXCLUDED.value,
    value_type = EXCLUDED.value_type,
    category = EXCLUDED.category,
    description = EXCLUDED.description,
    is_public = EXCLUDED.is_public;
