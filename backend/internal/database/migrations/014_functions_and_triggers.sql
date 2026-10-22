-- ============================================================
-- 014_functions_and_triggers.sql
-- CONSOLIDATES: cleanup functions + stored procedures + triggers
-- PURPOSE: Database functions, triggers, and maintenance procedures
-- FIXES: Corrected column references and tier values to match actual schema
-- ============================================================

-- ══════════════════════════════════════════════════════════
-- PART A: CLEANUP FUNCTIONS
-- ══════════════════════════════════════════════════════════

-- Function to cleanup unverified users after 24 hours
CREATE OR REPLACE FUNCTION cleanup_unverified_users() RETURNS INT AS $$
DECLARE
    deleted_count INT;
BEGIN
    WITH deleted AS (
        DELETE FROM users
        WHERE otp_verified = false
          AND created_at < NOW() - INTERVAL '24 hours'
        RETURNING id
    )
    SELECT count(*) INTO deleted_count FROM deleted;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_unverified_users IS 'Deletes unverified users older than 24 hours, returns count';

-- Function to cleanup expired products
CREATE OR REPLACE FUNCTION cleanup_expired_products() RETURNS INT AS $$
DECLARE
    updated_count INT;
BEGIN
    WITH updated AS (
        UPDATE products
        SET status = 'expired'
        WHERE status = 'approved'
          AND expires_at IS NOT NULL
          AND expires_at < NOW()
        RETURNING id
    )
    SELECT count(*) INTO updated_count FROM updated;
    RETURN updated_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_expired_products IS 'Marks expired products as expired, returns count';

-- ══════════════════════════════════════════════════════════
-- PART B: RATING UPDATE TRIGGERS
-- ══════════════════════════════════════════════════════════

-- Function to update user rating after review
CREATE OR REPLACE FUNCTION update_seller_rating() RETURNS TRIGGER AS $$
BEGIN
    -- Update seller average rating
    UPDATE users
    SET rating = (
        SELECT ROUND(AVG(rating)::NUMERIC, 2)
        FROM seller_reviews
        WHERE seller_id = NEW.seller_id
    )
    WHERE id = NEW.seller_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update seller rating on new review
DROP TRIGGER IF EXISTS trg_update_seller_rating ON seller_reviews;
CREATE TRIGGER trg_update_seller_rating
    AFTER INSERT OR UPDATE ON seller_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_seller_rating();

COMMENT ON FUNCTION update_seller_rating IS 'Updates seller average rating when reviews are added/updated';

-- ══════════════════════════════════════════════════════════
-- PART C: PRODUCT EXPIRATION TRIGGER
-- FIXED: Uses actual seller_tier enum values (normal/gold/premium)
--        instead of phantom tiers (platinum/silver/bronze)
-- ══════════════════════════════════════════════════════════

-- Function to set product expiration date based on seller tier
CREATE OR REPLACE FUNCTION set_product_expiration() RETURNS TRIGGER AS $$
DECLARE
    user_tier VARCHAR(20);
    expiry_days INT;
BEGIN
    -- Get seller tier
    SELECT tier::text INTO user_tier FROM users WHERE id = NEW.seller_id;

    -- Set expiration based on tier (if not manually set)
    IF NEW.expires_at IS NULL THEN
        CASE user_tier
            WHEN 'premium' THEN expiry_days := 90;
            WHEN 'gold' THEN expiry_days := 60;
            WHEN 'normal' THEN expiry_days := 30;
            ELSE expiry_days := 30;
        END CASE;

        NEW.expires_at := NOW() + (expiry_days || ' days')::INTERVAL;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to set expiration on product creation
DROP TRIGGER IF EXISTS trg_set_product_expiration ON products;
CREATE TRIGGER trg_set_product_expiration
    BEFORE INSERT ON products
    FOR EACH ROW
    EXECUTE FUNCTION set_product_expiration();

COMMENT ON FUNCTION set_product_expiration IS 'Automatically sets product expiration date based on seller tier';

-- ══════════════════════════════════════════════════════════
-- PART D: TIER SNAPSHOT TRIGGER
-- FIXED: Default is 'normal' not 'bronze' (matches CHECK constraint)
-- ══════════════════════════════════════════════════════════

-- Function to snapshot seller tier when creating products
CREATE OR REPLACE FUNCTION snapshot_seller_tier() RETURNS TRIGGER AS $$
DECLARE
    user_tier VARCHAR(20);
BEGIN
    -- Get current seller tier and snapshot it
    SELECT tier::text INTO user_tier FROM users WHERE id = NEW.seller_id;
    NEW.seller_tier_snapshot := COALESCE(user_tier, 'normal');

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to snapshot tier on product creation
DROP TRIGGER IF EXISTS trg_snapshot_seller_tier ON products;
CREATE TRIGGER trg_snapshot_seller_tier
    BEFORE INSERT ON products
    FOR EACH ROW
    EXECUTE FUNCTION snapshot_seller_tier();

COMMENT ON FUNCTION snapshot_seller_tier IS 'Snapshots seller tier at product creation time';

-- ══════════════════════════════════════════════════════════
-- PART E: REFERRAL COMMISSION TRIGGER
-- FIXED: Uses correct column names:
--   - 'invited_by' not 'referred_by' (users table)
--   - 'user_id' not 'buyer_id' (transactions table)
--   - 'referred_id' not 'referred_user_id' (referral_commissions table)
-- ══════════════════════════════════════════════════════════

-- Function to create pending commission on transaction
CREATE OR REPLACE FUNCTION create_pending_commission() RETURNS TRIGGER AS $$
DECLARE
    referrer_id_val INT;
    commission_rate NUMERIC;
    commission_amt NUMERIC;
BEGIN
    -- Only process for completed transactions
    IF NEW.status = 'completed' THEN
        -- Get user's referrer (uses 'invited_by' column, not 'referred_by')
        SELECT invited_by INTO referrer_id_val
        FROM users
        WHERE id = NEW.user_id AND invited_by IS NOT NULL;

        IF referrer_id_val IS NOT NULL THEN
            -- Get commission rate from settings (default 5%)
            commission_rate := 0.05; -- Can be fetched from platform_settings
            commission_amt := NEW.amount * commission_rate;

            -- Create pending commission (uses 'referred_id' column)
            INSERT INTO referral_commissions (
                referrer_id,
                referred_id,
                transaction_id,
                commission_amount,
                commission_percent,
                transaction_type,
                status,
                created_at
            )
            VALUES (
                referrer_id_val,
                NEW.user_id,
                NEW.id,
                commission_amt,
                commission_rate * 100,
                COALESCE(NEW.transaction_type, 'listing_fee'),
                'pending',
                NOW()
            )
            ON CONFLICT DO NOTHING;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to create commission on completed transaction
DROP TRIGGER IF EXISTS trg_create_pending_commission ON transactions;
CREATE TRIGGER trg_create_pending_commission
    AFTER INSERT OR UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION create_pending_commission();

COMMENT ON FUNCTION create_pending_commission IS 'Creates pending referral commission when transaction completes';

-- ══════════════════════════════════════════════════════════
-- PART F: MAINTENANCE HELPER FUNCTIONS
-- ══════════════════════════════════════════════════════════

-- Function to recalculate all seller ratings (maintenance)
CREATE OR REPLACE FUNCTION recalculate_all_ratings() RETURNS INT AS $$
DECLARE
    updated_count INT;
BEGIN
    WITH rating_updates AS (
        SELECT
            seller_id,
            ROUND(AVG(rating)::NUMERIC, 2) as new_rating
        FROM seller_reviews
        GROUP BY seller_id
    )
    UPDATE users u
    SET rating = ru.new_rating
    FROM rating_updates ru
    WHERE u.id = ru.seller_id
      AND u.role = 'seller'
    RETURNING u.id INTO updated_count;

    RETURN updated_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION recalculate_all_ratings IS 'Recalculates all seller ratings (use for maintenance)';

-- Function to refresh all materialized views
CREATE OR REPLACE FUNCTION refresh_all_materialized_views() RETURNS TEXT AS $$
DECLARE
    refresh_log TEXT := '';
BEGIN
    -- Refresh seller rating stats
    REFRESH MATERIALIZED VIEW CONCURRENTLY seller_rating_stats;
    refresh_log := refresh_log || 'seller_rating_stats refreshed; ';

    -- Refresh user warning counts
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_warning_counts;
    refresh_log := refresh_log || 'user_warning_counts refreshed; ';

    -- Refresh marketplace listings (created in 015)
    BEGIN
        REFRESH MATERIALIZED VIEW CONCURRENTLY marketplace_listings;
        refresh_log := refresh_log || 'marketplace_listings refreshed; ';
    EXCEPTION
        WHEN undefined_table THEN
            refresh_log := refresh_log || 'marketplace_listings not yet created; ';
    END;

    RETURN refresh_log || 'All views refreshed successfully';
EXCEPTION
    WHEN OTHERS THEN
        RETURN 'Error refreshing views: ' || SQLERRM;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION refresh_all_materialized_views IS 'Refreshes all materialized views (use in cron job)';

-- ══════════════════════════════════════════════════════════
-- PART G: SCHEDULED MAINTENANCE RECOMMENDATIONS
-- ══════════════════════════════════════════════════════════

-- Add these to your cron job or background worker:
-- 1. Run daily: SELECT cleanup_unverified_users();
-- 2. Run hourly: SELECT cleanup_expired_products();
-- 3. Run every 5 minutes: SELECT refresh_all_materialized_views();
-- 4. Run weekly: VACUUM ANALYZE;

COMMENT ON SCHEMA public IS 'Scheduled maintenance functions configured in migration 014';
