-- ============================================================================
-- 067_core_functions_and_triggers.sql
-- Core platform procedures, maintenance routines, and automatic triggers
-- ============================================================================

CREATE OR REPLACE FUNCTION public.update_updated_at_column()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.cleanup_unverified_users()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
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
$function$;

CREATE OR REPLACE FUNCTION public.cleanup_expired_products()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
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
$function$;

CREATE OR REPLACE FUNCTION public.cleanup_expired_otps()
 RETURNS void
 LANGUAGE plpgsql
AS $function$
BEGIN
    DELETE FROM otp_codes WHERE status IN ('verified', 'expired') AND created_at < NOW() - INTERVAL '7 days';
    DELETE FROM otp_codes WHERE status = 'locked' AND created_at < NOW() - INTERVAL '24 hours';
END;
$function$;

CREATE OR REPLACE FUNCTION public.cleanup_old_activity_logs()
 RETURNS void
 LANGUAGE plpgsql
AS $function$
BEGIN
    DELETE FROM user_activity_logs WHERE created_at < NOW() - INTERVAL '6 months';
END;
$function$;

CREATE OR REPLACE FUNCTION public.cleanup_old_notifications()
 RETURNS void
 LANGUAGE plpgsql
AS $function$
BEGIN
    DELETE FROM notifications
    WHERE is_read = TRUE
      AND created_at < NOW() - INTERVAL '90 days';
END;
$function$;

CREATE OR REPLACE FUNCTION public.update_seller_rating()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
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
$function$;

CREATE OR REPLACE FUNCTION public.snapshot_seller_tier()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
    user_tier VARCHAR(20);
BEGIN
    -- Get current seller tier and snapshot it
    SELECT tier::text INTO user_tier FROM users WHERE id = NEW.seller_id;
    NEW.seller_tier_snapshot := COALESCE(user_tier, 'normal');

    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.create_pending_commission()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
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
$function$;

CREATE OR REPLACE FUNCTION public.auto_confirm_commissions()
 RETURNS void
 LANGUAGE plpgsql
AS $function$
BEGIN
    UPDATE referral_commissions
    SET status = 'confirmed', confirmed_at = NOW()
    WHERE status = 'pending'
      AND created_at < NOW() - INTERVAL '30 days';
END;
$function$;

CREATE OR REPLACE FUNCTION public.recalculate_all_ratings()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
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
$function$;

CREATE OR REPLACE FUNCTION public.record_tier_change()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    IF (OLD.tier IS DISTINCT FROM NEW.tier) OR (OLD.buyer_tier IS DISTINCT FROM NEW.buyer_tier) THEN
        INSERT INTO user_tier_history (user_id, tier_type, old_tier, new_tier, reason)
        VALUES (
            NEW.id,
            CASE WHEN OLD.tier IS DISTINCT FROM NEW.tier THEN 'seller' ELSE 'buyer' END,
            COALESCE(CASE WHEN OLD.tier IS DISTINCT FROM NEW.tier THEN OLD.tier::text ELSE OLD.buyer_tier::text END, 'normal'),
            COALESCE(CASE WHEN OLD.tier IS DISTINCT FROM NEW.tier THEN NEW.tier::text ELSE NEW.buyer_tier::text END, 'normal'),
            'system'
        );
    END IF;
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.log_setting_change()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    IF TG_OP = 'UPDATE' AND OLD.value IS DISTINCT FROM NEW.value THEN
        INSERT INTO settings_change_log (setting_key, old_value, new_value, changed_by)
        VALUES (NEW.key, OLD.value, NEW.value, NEW.updated_by);
    END IF;
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.refresh_all_materialized_views()
 RETURNS text
 LANGUAGE plpgsql
AS $function$
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
$function$;

CREATE OR REPLACE FUNCTION public.sync_otp_codes_columns()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    IF NEW.code != '' AND (NEW.code_encrypted IS NULL OR NEW.code_encrypted = '') THEN
        NEW.code_encrypted := NEW.code;
    ELSIF NEW.code_encrypted != '' AND (NEW.code IS NULL OR NEW.code = '') THEN
        NEW.code := NEW.code_encrypted;
    END IF;

    IF NEW.used = TRUE AND NEW.status = 'pending' THEN
        NEW.status := 'verified';
    ELSIF NEW.status = 'verified' AND NEW.used = FALSE THEN
        NEW.used := TRUE;
    END IF;

    IF NEW.locked_until IS NOT NULL AND NEW.locked_until > NOW() THEN
        NEW.status := 'locked';
    END IF;

    RETURN NEW;
END;
$function$;

-- Triggers
DROP TRIGGER IF EXISTS trg_update_seller_rating ON seller_reviews;
CREATE TRIGGER trg_update_seller_rating AFTER INSERT OR UPDATE ON seller_reviews FOR EACH ROW EXECUTE FUNCTION update_seller_rating();

DROP TRIGGER IF EXISTS trg_snapshot_seller_tier ON products;
CREATE TRIGGER trg_snapshot_seller_tier BEFORE INSERT ON products FOR EACH ROW EXECUTE FUNCTION snapshot_seller_tier();

DROP TRIGGER IF EXISTS trg_create_pending_commission ON transactions;
CREATE TRIGGER trg_create_pending_commission AFTER INSERT OR UPDATE ON transactions FOR EACH ROW EXECUTE FUNCTION create_pending_commission();

DROP TRIGGER IF EXISTS track_tier_changes ON users;
CREATE TRIGGER track_tier_changes AFTER UPDATE ON users FOR EACH ROW EXECUTE FUNCTION record_tier_change();

DROP TRIGGER IF EXISTS track_setting_changes ON platform_settings;
CREATE TRIGGER track_setting_changes AFTER UPDATE ON platform_settings FOR EACH ROW EXECUTE FUNCTION log_setting_change();

DROP TRIGGER IF EXISTS trg_sync_otp_codes ON otp_codes;
CREATE TRIGGER trg_sync_otp_codes BEFORE INSERT OR UPDATE ON otp_codes FOR EACH ROW EXECUTE FUNCTION sync_otp_codes_columns();

DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_order_requests_updated_at ON order_requests;
CREATE TRIGGER update_order_requests_updated_at BEFORE UPDATE ON order_requests FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_hero_banners_updated_at ON hero_banners;
CREATE TRIGGER update_hero_banners_updated_at BEFORE UPDATE ON hero_banners FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_referral_withdrawals_updated_at ON referral_withdrawals;
CREATE TRIGGER update_referral_withdrawals_updated_at BEFORE UPDATE ON referral_withdrawals FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_upgrade_requests_updated_at ON upgrade_requests;
CREATE TRIGGER update_upgrade_requests_updated_at BEFORE UPDATE ON upgrade_requests FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_refund_requests_updated_at ON refund_requests;
CREATE TRIGGER update_refund_requests_updated_at BEFORE UPDATE ON refund_requests FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_abuse_reports_updated_at ON abuse_reports;
CREATE TRIGGER update_abuse_reports_updated_at BEFORE UPDATE ON abuse_reports FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_appeals_updated_at ON account_appeals;
CREATE TRIGGER update_appeals_updated_at BEFORE UPDATE ON account_appeals FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_verification_docs_updated_at ON seller_verification_documents;
CREATE TRIGGER update_verification_docs_updated_at BEFORE UPDATE ON seller_verification_documents FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
