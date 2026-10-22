-- ============================================================
-- 013_reviews_and_warnings.sql
-- CONSOLIDATES: seller_reviews + user_warnings
-- PURPOSE: User moderation and feedback system
-- NOTE: Uses ALTER TABLE to extend shell tables from 001_schema.sql
-- ============================================================

-- ══════════════════════════════════════════════════════════
-- PART A: SELLER REVIEWS
-- ══════════════════════════════════════════════════════════

ALTER TABLE seller_reviews
    ADD COLUMN IF NOT EXISTS seller_id   INT          NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS reviewer_id INT          NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS order_id    INT          REFERENCES order_requests(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS rating      NUMERIC(2,1) NOT NULL DEFAULT 1 CHECK (rating >= 1 AND rating <= 5),
    ADD COLUMN IF NOT EXISTS comment     TEXT,
    ADD COLUMN IF NOT EXISTS created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW();

-- Ensure one review per order per reviewer
DO $ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name='seller_reviews' AND constraint_name='uq_review_per_order'
    ) THEN
        ALTER TABLE seller_reviews ADD CONSTRAINT uq_review_per_order UNIQUE (order_id, reviewer_id);
    END IF;
END $;

-- Indexes for efficient review queries
CREATE INDEX IF NOT EXISTS idx_seller_reviews_seller_id   ON seller_reviews (seller_id);
CREATE INDEX IF NOT EXISTS idx_seller_reviews_reviewer_id ON seller_reviews (reviewer_id);
CREATE INDEX IF NOT EXISTS idx_seller_reviews_order_id    ON seller_reviews (order_id) WHERE order_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_seller_reviews_rating      ON seller_reviews (rating);

COMMENT ON TABLE seller_reviews IS 'Buyer reviews for sellers after completed orders';
COMMENT ON COLUMN seller_reviews.rating IS 'Rating from 1.0 to 5.0 stars';

-- ══════════════════════════════════════════════════════════
-- PART B: USER WARNINGS
-- ══════════════════════════════════════════════════════════

ALTER TABLE user_warnings
    ADD COLUMN IF NOT EXISTS user_id    INT         REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS admin_id   INT         REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS reason     TEXT        NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS severity   VARCHAR(20) DEFAULT 'low' CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    ADD COLUMN IF NOT EXISTS is_active  BOOLEAN     DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Index for user warnings lookup
CREATE INDEX IF NOT EXISTS idx_user_warnings_user_id   ON user_warnings (user_id);
CREATE INDEX IF NOT EXISTS idx_user_warnings_admin_id  ON user_warnings (admin_id);
CREATE INDEX IF NOT EXISTS idx_user_warnings_severity  ON user_warnings (severity) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_user_warnings_active    ON user_warnings (is_active, user_id);

COMMENT ON TABLE user_warnings IS 'Admin warnings issued to users for policy violations';
COMMENT ON COLUMN user_warnings.severity IS 'Warning severity: low, medium, high, critical';
COMMENT ON COLUMN user_warnings.is_active IS 'Whether warning is still active (not expired or resolved)';

-- ══════════════════════════════════════════════════════════
-- PART C: BUYER TIER PAYMENTS
-- ══════════════════════════════════════════════════════════

ALTER TABLE buyer_tier_payments
    ADD COLUMN IF NOT EXISTS buyer_id        INT           NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS tier            VARCHAR(20)   NOT NULL DEFAULT 'gold',
    ADD COLUMN IF NOT EXISTS amount          NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reference       VARCHAR(200),
    ADD COLUMN IF NOT EXISTS payer_phone     VARCHAR(20),
    ADD COLUMN IF NOT EXISTS payment_method  VARCHAR(50),
    ADD COLUMN IF NOT EXISTS created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_buyer_tier_payments_buyer ON buyer_tier_payments(buyer_id, created_at DESC);

-- Verify critical tables exist
DO $
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'user_warnings'
    ) THEN
        RAISE EXCEPTION 'user_warnings table is missing — check 001_schema.sql';
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'seller_reviews'
    ) THEN
        RAISE EXCEPTION 'seller_reviews table is missing — check 001_schema.sql';
    END IF;
END $;

-- ══════════════════════════════════════════════════════════
-- PART D: ANALYTICS & AGGREGATE VIEWS
-- ══════════════════════════════════════════════════════════

-- Seller rating aggregate (cached view for performance)
CREATE MATERIALIZED VIEW IF NOT EXISTS seller_rating_stats AS
SELECT
    seller_id,
    COUNT(*) as total_reviews,
    ROUND(AVG(rating)::NUMERIC, 2) as avg_rating,
    COUNT(*) FILTER (WHERE rating >= 4) as positive_reviews,
    COUNT(*) FILTER (WHERE rating <= 2) as negative_reviews,
    MAX(created_at) as last_review_date
FROM seller_reviews
GROUP BY seller_id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_seller_rating_stats_seller_id ON seller_rating_stats (seller_id);

-- User warning count (for auto-ban logic)
CREATE MATERIALIZED VIEW IF NOT EXISTS user_warning_counts AS
SELECT
    user_id,
    COUNT(*) as total_warnings,
    COUNT(*) FILTER (WHERE severity = 'critical') as critical_count,
    COUNT(*) FILTER (WHERE severity = 'high') as high_count,
    COUNT(*) FILTER (WHERE is_active = TRUE) as active_warnings,
    MAX(created_at) as last_warning_date
FROM user_warnings
GROUP BY user_id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_warning_counts_user_id ON user_warning_counts (user_id);

COMMENT ON MATERIALIZED VIEW seller_rating_stats IS 'Cached seller rating statistics (refresh periodically)';
COMMENT ON MATERIALIZED VIEW user_warning_counts IS 'Cached user warning counts (refresh periodically)';
