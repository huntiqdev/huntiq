-- ============================================================
-- 025_profit_analytics_views.sql
-- PURPOSE: Create comprehensive profit analytics views required by
--          admin_repo.go (ProfitByTier, MonthlyProfitTrend,
--          WeeklyProfitTrend, AnnualProfitTrend, ProfitSummary).
-- ============================================================

-- 1. admin_profit_summary (all 12 columns expected by ProfitSummary)
DROP VIEW IF EXISTS admin_profit_summary CASCADE;
CREATE OR REPLACE VIEW admin_profit_summary AS
SELECT
    (SELECT COUNT(*)::INT FROM transactions WHERE status = 'completed') AS total_transactions,
    (SELECT COALESCE(SUM(amount), 0)::FLOAT8 FROM transactions WHERE status = 'completed') AS gross_volume,
    (SELECT COALESCE(SUM(broker_fee), 0)::FLOAT8 FROM transactions WHERE status = 'completed') AS broker_fee_revenue,
    2.0::FLOAT8 AS broker_fee_pct,
    (SELECT COALESCE(SUM(total_fee_amount), 0)::FLOAT8 FROM products WHERE fee_paid = TRUE AND status IN ('approved', 'expired', 'sold')) AS listing_fee_collected,
    (SELECT COALESCE(SUM(total_fee_amount), 0)::FLOAT8 FROM products WHERE fee_paid = FALSE AND status = 'pending') AS listing_fee_outstanding,
    (SELECT COALESCE(SUM(total_fee_amount), 0)::FLOAT8 FROM products) AS listing_fee_total,
    (SELECT COALESCE(SUM(amount), 0)::FLOAT8 FROM upgrade_requests WHERE status = 'approved' AND role = 'seller') AS seller_upgrade_revenue,
    (SELECT COALESCE(SUM(amount), 0)::FLOAT8 FROM upgrade_requests WHERE status = 'approved' AND role = 'buyer') AS buyer_upgrade_revenue,
    (SELECT COALESCE(SUM(amount), 0)::FLOAT8 FROM upgrade_requests WHERE status = 'approved') AS total_upgrade_revenue,
    ((SELECT COALESCE(SUM(broker_fee), 0)::FLOAT8 FROM transactions WHERE status = 'completed')
     + (SELECT COALESCE(SUM(total_fee_amount), 0)::FLOAT8 FROM products WHERE fee_paid = TRUE AND status IN ('approved', 'expired', 'sold'))
     + (SELECT COALESCE(SUM(amount), 0)::FLOAT8 FROM upgrade_requests WHERE status = 'approved')) AS total_platform_profit,
    '1.2'::VARCHAR AS referral_commission_rate;

-- 2. profit_by_seller_tier
CREATE OR REPLACE VIEW profit_by_seller_tier AS
WITH tiers AS (
    SELECT 'normal' AS tier
    UNION ALL SELECT 'gold'
    UNION ALL SELECT 'premium'
),
prod_stats AS (
    SELECT
        COALESCE(u.tier::text, 'normal') AS tier,
        COUNT(p.id) AS products,
        COALESCE(SUM(p.total_fee_amount) FILTER (WHERE p.fee_paid = TRUE), 0) AS listing_fees
    FROM users u
    LEFT JOIN products p ON p.seller_id = u.id AND p.status IN ('approved', 'expired', 'sold')
    WHERE u.role = 'seller'
    GROUP BY COALESCE(u.tier::text, 'normal')
),
tx_stats AS (
    SELECT
        COALESCE(u.tier::text, 'normal') AS tier,
        COALESCE(SUM(t.amount), 0) AS volume,
        COALESCE(SUM(t.broker_fee), 0) AS broker_fees
    FROM users u
    LEFT JOIN transactions t ON t.seller_id = u.id AND t.status = 'completed'
    WHERE u.role = 'seller'
    GROUP BY COALESCE(u.tier::text, 'normal')
),
upg_stats AS (
    SELECT
        COALESCE(ur.requested_tier, 'gold') AS tier,
        COALESCE(SUM(ur.amount), 0) AS upgrade_fees
    FROM upgrade_requests ur
    WHERE ur.status = 'approved' AND ur.role = 'seller'
    GROUP BY COALESCE(ur.requested_tier, 'gold')
)
SELECT
    t.tier,
    COALESCE(ps.products, 0)::INT AS products,
    COALESCE(ts.volume, 0)::FLOAT8 AS volume,
    COALESCE(ts.broker_fees, 0)::FLOAT8 AS broker_fees,
    COALESCE(ps.listing_fees, 0)::FLOAT8 AS listing_fees,
    COALESCE(us.upgrade_fees, 0)::FLOAT8 AS upgrade_fees,
    (COALESCE(ts.broker_fees, 0) + COALESCE(ps.listing_fees, 0) + COALESCE(us.upgrade_fees, 0))::FLOAT8 AS total_profit
FROM tiers t
LEFT JOIN prod_stats ps ON ps.tier = t.tier
LEFT JOIN tx_stats ts ON ts.tier = t.tier
LEFT JOIN upg_stats us ON us.tier = t.tier;

-- 3. monthly_profit_trend
CREATE OR REPLACE VIEW monthly_profit_trend AS
WITH months AS (
    SELECT to_char(d, 'YYYY-MM') AS month, date_trunc('month', d) AS m_start
    FROM generate_series(NOW() - INTERVAL '11 months', NOW(), INTERVAL '1 month') d
),
tx_m AS (
    SELECT
        to_char(created_at, 'YYYY-MM') AS month,
        COUNT(id) AS tx_count,
        COALESCE(SUM(amount), 0) AS volume,
        COALESCE(SUM(broker_fee), 0) AS broker_fees
    FROM transactions
    WHERE status = 'completed'
    GROUP BY to_char(created_at, 'YYYY-MM')
),
fee_m AS (
    SELECT
        to_char(created_at, 'YYYY-MM') AS month,
        COALESCE(SUM(total_fee_amount), 0) AS listing_fees
    FROM products
    WHERE fee_paid = TRUE AND status IN ('approved', 'expired', 'sold')
    GROUP BY to_char(created_at, 'YYYY-MM')
)
SELECT
    m.month,
    COALESCE(tx.tx_count, 0)::INT AS transactions,
    COALESCE(tx.volume, 0)::FLOAT8 AS volume,
    COALESCE(tx.broker_fees, 0)::FLOAT8 AS broker_fees,
    COALESCE(fm.listing_fees, 0)::FLOAT8 AS listing_fees,
    (COALESCE(tx.broker_fees, 0) + COALESCE(fm.listing_fees, 0))::FLOAT8 AS total_profit,
    CASE
        WHEN COALESCE(tx.volume, 0) > 0 THEN ROUND(((COALESCE(tx.broker_fees, 0) + COALESCE(fm.listing_fees, 0)) / tx.volume * 100)::numeric, 2)::FLOAT8
        ELSE 0.0
    END AS avg_fee_pct
FROM months m
LEFT JOIN tx_m tx ON tx.month = m.month
LEFT JOIN fee_m fm ON fm.month = m.month
ORDER BY m.m_start DESC;

-- 4. weekly_profit_trend
CREATE OR REPLACE VIEW weekly_profit_trend AS
WITH weeks AS (
    SELECT
        to_char(d, 'YYYY-"W"IW') AS week,
        date_trunc('week', d) AS week_start
    FROM generate_series(NOW() - INTERVAL '12 weeks', NOW(), INTERVAL '1 week') d
),
tx_w AS (
    SELECT
        to_char(created_at, 'YYYY-"W"IW') AS week,
        COUNT(id) AS tx_count,
        COALESCE(SUM(amount), 0) AS volume,
        COALESCE(SUM(broker_fee), 0) AS broker_fees
    FROM transactions
    WHERE status = 'completed'
    GROUP BY to_char(created_at, 'YYYY-"W"IW')
),
fee_w AS (
    SELECT
        to_char(created_at, 'YYYY-"W"IW') AS week,
        COALESCE(SUM(total_fee_amount), 0) AS listing_fees
    FROM products
    WHERE fee_paid = TRUE AND status IN ('approved', 'expired', 'sold')
    GROUP BY to_char(created_at, 'YYYY-"W"IW')
)
SELECT
    w.week,
    COALESCE(tx.volume, 0)::FLOAT8 AS volume,
    COALESCE(tx.broker_fees, 0)::FLOAT8 AS broker_fees,
    COALESCE(fw.listing_fees, 0)::FLOAT8 AS listing_fees,
    (COALESCE(tx.broker_fees, 0) + COALESCE(fw.listing_fees, 0))::FLOAT8 AS total_profit,
    COALESCE(tx.tx_count, 0)::INT AS transactions,
    CASE
        WHEN COALESCE(tx.volume, 0) > 0 THEN ROUND(((COALESCE(tx.broker_fees, 0) + COALESCE(fw.listing_fees, 0)) / tx.volume * 100)::numeric, 2)::FLOAT8
        ELSE 0.0
    END AS avg_fee_pct,
    w.week_start
FROM weeks w
LEFT JOIN tx_w tx ON tx.week = w.week
LEFT JOIN fee_w fw ON fw.week = w.week;

-- 5. annual_profit_trend
CREATE OR REPLACE VIEW annual_profit_trend AS
WITH years AS (
    SELECT EXTRACT(YEAR FROM d)::INT AS year
    FROM generate_series(NOW() - INTERVAL '4 years', NOW(), INTERVAL '1 year') d
),
tx_y AS (
    SELECT
        EXTRACT(YEAR FROM created_at)::INT AS year,
        COUNT(id) AS tx_count,
        COALESCE(SUM(amount), 0) AS volume,
        COALESCE(SUM(broker_fee), 0) AS broker_fees
    FROM transactions
    WHERE status = 'completed'
    GROUP BY EXTRACT(YEAR FROM created_at)::INT
),
fee_y AS (
    SELECT
        EXTRACT(YEAR FROM created_at)::INT AS year,
        COALESCE(SUM(total_fee_amount), 0) AS listing_fees
    FROM products
    WHERE fee_paid = TRUE AND status IN ('approved', 'expired', 'sold')
    GROUP BY EXTRACT(YEAR FROM created_at)::INT
)
SELECT
    y.year,
    COALESCE(ty.tx_count, 0)::INT AS transactions,
    COALESCE(ty.volume, 0)::FLOAT8 AS volume,
    COALESCE(ty.broker_fees, 0)::FLOAT8 AS broker_fees,
    COALESCE(fy.listing_fees, 0)::FLOAT8 AS listing_fees,
    (COALESCE(ty.broker_fees, 0) + COALESCE(fy.listing_fees, 0))::FLOAT8 AS total_profit,
    CASE
        WHEN COALESCE(ty.volume, 0) > 0 THEN ROUND(((COALESCE(ty.broker_fees, 0) + COALESCE(fy.listing_fees, 0)) / ty.volume * 100)::numeric, 2)::FLOAT8
        ELSE 0.0
    END AS avg_fee_pct
FROM years y
LEFT JOIN tx_y ty ON ty.year = y.year
LEFT JOIN fee_y fy ON fy.year = y.year;
