-- ============================================================
-- 020_seller_profit_view.sql
-- PURPOSE: Create the seller_profit_summary view required by
--          seller_repo.go's sqlSellerFinancialStats query.
--
-- Referenced columns:
--   total_revenue   — sum of all transaction amounts for a seller
--   net_profit      — revenue minus all confirmed listing fees paid
--   total_fees_paid — sum of total_fee_amount on fee_paid products
-- ============================================================

CREATE OR REPLACE VIEW seller_profit_summary AS
SELECT
    p.seller_id,
    -- Gross revenue: sum of completed transaction amounts
    COALESCE(
        (SELECT SUM(t.amount)
         FROM transactions t
         WHERE t.seller_id = p.seller_id
           AND t.status = 'completed'),
        0
    ) AS total_revenue,
    -- Total listing fees confirmed as paid by admin
    COALESCE(
        SUM(p.total_fee_amount) FILTER (WHERE p.fee_paid = TRUE),
        0
    ) AS total_fees_paid,
    -- Net profit = revenue - fees paid
    COALESCE(
        (SELECT SUM(t.amount)
         FROM transactions t
         WHERE t.seller_id = p.seller_id
           AND t.status = 'completed'),
        0
    ) - COALESCE(
        SUM(p.total_fee_amount) FILTER (WHERE p.fee_paid = TRUE),
        0
    ) AS net_profit
FROM products p
GROUP BY p.seller_id;

COMMENT ON VIEW seller_profit_summary IS
    'Per-seller financial summary: total_revenue, total_fees_paid, net_profit. '
    'Used by seller dashboard stats endpoint.';

-- Grant no special permissions — app uses a single DB user
-- Refresh is not needed (this is a plain view, not materialized)

-- ── Admin Profit Summary ─────────────────────────────────────────────────────
-- Referenced by admin_repo.go AnalyticsOverview:
--   SELECT COALESCE(total_platform_profit,0) FROM admin_profit_summary LIMIT 1
CREATE OR REPLACE VIEW admin_profit_summary AS
SELECT
    -- Total broker fees from completed transactions
    COALESCE((SELECT SUM(broker_fee) FROM transactions WHERE status = 'completed'), 0) AS total_broker_revenue,
    -- Total listing fees confirmed paid across the platform
    COALESCE((SELECT SUM(total_fee_amount) FROM products WHERE fee_paid = TRUE
              AND status IN ('approved', 'expired', 'sold')), 0) AS total_listing_fee_revenue,
    -- Combined platform profit
    COALESCE((SELECT SUM(broker_fee) FROM transactions WHERE status = 'completed'), 0)
    + COALESCE((SELECT SUM(total_fee_amount) FROM products WHERE fee_paid = TRUE
                AND status IN ('approved', 'expired', 'sold')), 0) AS total_platform_profit;

COMMENT ON VIEW admin_profit_summary IS
    'Platform-wide profit summary: broker fees + listing fees. '
    'Referenced by AnalyticsOverview in admin_repo.go.';
