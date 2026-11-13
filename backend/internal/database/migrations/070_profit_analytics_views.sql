-- ============================================================================
-- 070_profit_analytics_views.sql
-- Financial and profit analytics reporting views
-- ============================================================================

CREATE OR REPLACE VIEW seller_profit_summary AS
SELECT seller_id,
    COALESCE(( SELECT sum(t.amount) AS sum
           FROM transactions t
          WHERE ((t.seller_id = p.seller_id) AND ((t.status)::text = 'completed'::text))), (0)::numeric) AS total_revenue,
    COALESCE(sum(total_fee_amount) FILTER (WHERE (fee_paid = true)), (0)::numeric) AS total_fees_paid,
    (COALESCE(( SELECT sum(t.amount) AS sum
           FROM transactions t
          WHERE ((t.seller_id = p.seller_id) AND ((t.status)::text = 'completed'::text))), (0)::numeric) - COALESCE(sum(total_fee_amount) FILTER (WHERE (fee_paid = true)), (0)::numeric)) AS net_profit
   FROM products p
  GROUP BY seller_id;

CREATE OR REPLACE VIEW admin_profit_summary AS
SELECT ( SELECT (count(*))::integer AS count
           FROM transactions
          WHERE ((transactions.status)::text = 'completed'::text)) AS total_transactions,
    ( SELECT (COALESCE(sum(transactions.amount), (0)::numeric))::double precision AS "coalesce"
           FROM transactions
          WHERE ((transactions.status)::text = 'completed'::text)) AS gross_volume,
    ( SELECT (COALESCE(sum(transactions.broker_fee), (0)::numeric))::double precision AS "coalesce"
           FROM transactions
          WHERE ((transactions.status)::text = 'completed'::text)) AS broker_fee_revenue,
    (2.0)::double precision AS broker_fee_pct,
    ( SELECT (COALESCE(sum(products.total_fee_amount), (0)::numeric))::double precision AS "coalesce"
           FROM products
          WHERE ((products.fee_paid = true) AND (products.status = ANY (ARRAY['approved'::product_status, 'expired'::product_status, 'sold'::product_status])))) AS listing_fee_collected,
    ( SELECT (COALESCE(sum(products.total_fee_amount), (0)::numeric))::double precision AS "coalesce"
           FROM products
          WHERE ((products.fee_paid = false) AND (products.status = 'pending'::product_status))) AS listing_fee_outstanding,
    ( SELECT (COALESCE(sum(products.total_fee_amount), (0)::numeric))::double precision AS "coalesce"
           FROM products) AS listing_fee_total,
    ( SELECT (COALESCE(sum(upgrade_requests.amount), (0)::numeric))::double precision AS "coalesce"
           FROM upgrade_requests
          WHERE (((upgrade_requests.status)::text = 'approved'::text) AND ((upgrade_requests.role)::text = 'seller'::text))) AS seller_upgrade_revenue,
    ( SELECT (COALESCE(sum(upgrade_requests.amount), (0)::numeric))::double precision AS "coalesce"
           FROM upgrade_requests
          WHERE (((upgrade_requests.status)::text = 'approved'::text) AND ((upgrade_requests.role)::text = 'buyer'::text))) AS buyer_upgrade_revenue,
    ( SELECT (COALESCE(sum(upgrade_requests.amount), (0)::numeric))::double precision AS "coalesce"
           FROM upgrade_requests
          WHERE ((upgrade_requests.status)::text = 'approved'::text)) AS total_upgrade_revenue,
    ((( SELECT (COALESCE(sum(transactions.broker_fee), (0)::numeric))::double precision AS "coalesce"
           FROM transactions
          WHERE ((transactions.status)::text = 'completed'::text)) + ( SELECT (COALESCE(sum(products.total_fee_amount), (0)::numeric))::double precision AS "coalesce"
           FROM products
          WHERE ((products.fee_paid = true) AND (products.status = ANY (ARRAY['approved'::product_status, 'expired'::product_status, 'sold'::product_status]))))) + ( SELECT (COALESCE(sum(upgrade_requests.amount), (0)::numeric))::double precision AS "coalesce"
           FROM upgrade_requests
          WHERE ((upgrade_requests.status)::text = 'approved'::text))) AS total_platform_profit,
    '1.2'::character varying AS referral_commission_rate;

CREATE OR REPLACE VIEW profit_by_seller_tier AS
WITH tiers AS (
         SELECT 'normal'::text AS tier
        UNION ALL
         SELECT 'gold'::text
        UNION ALL
         SELECT 'premium'::text
        ), prod_stats AS (
         SELECT COALESCE((u.tier)::text, 'normal'::text) AS tier,
            count(p.id) AS products,
            COALESCE(sum(p.total_fee_amount) FILTER (WHERE (p.fee_paid = true)), (0)::numeric) AS listing_fees
           FROM (users u
             LEFT JOIN products p ON (((p.seller_id = u.id) AND (p.status = ANY (ARRAY['approved'::product_status, 'expired'::product_status, 'sold'::product_status])))))
          WHERE ((u.role)::text = 'seller'::text)
          GROUP BY COALESCE((u.tier)::text, 'normal'::text)
        ), tx_stats AS (
         SELECT COALESCE((u.tier)::text, 'normal'::text) AS tier,
            COALESCE(sum(t_1.amount), (0)::numeric) AS volume,
            COALESCE(sum(t_1.broker_fee), (0)::numeric) AS broker_fees
           FROM (users u
             LEFT JOIN transactions t_1 ON (((t_1.seller_id = u.id) AND ((t_1.status)::text = 'completed'::text))))
          WHERE ((u.role)::text = 'seller'::text)
          GROUP BY COALESCE((u.tier)::text, 'normal'::text)
        ), upg_stats AS (
         SELECT COALESCE(ur.requested_tier, 'gold'::character varying) AS tier,
            COALESCE(sum(ur.amount), (0)::numeric) AS upgrade_fees
           FROM upgrade_requests ur
          WHERE (((ur.status)::text = 'approved'::text) AND ((ur.role)::text = 'seller'::text))
          GROUP BY COALESCE(ur.requested_tier, 'gold'::character varying)
        )
 SELECT t.tier,
    (COALESCE(ps.products, (0)::bigint))::integer AS products,
    (COALESCE(ts.volume, (0)::numeric))::double precision AS volume,
    (COALESCE(ts.broker_fees, (0)::numeric))::double precision AS broker_fees,
    (COALESCE(ps.listing_fees, (0)::numeric))::double precision AS listing_fees,
    (COALESCE(us.upgrade_fees, (0)::numeric))::double precision AS upgrade_fees,
    (((COALESCE(ts.broker_fees, (0)::numeric) + COALESCE(ps.listing_fees, (0)::numeric)) + COALESCE(us.upgrade_fees, (0)::numeric)))::double precision AS total_profit
   FROM (((tiers t
     LEFT JOIN prod_stats ps ON ((ps.tier = t.tier)))
     LEFT JOIN tx_stats ts ON ((ts.tier = t.tier)))
     LEFT JOIN upg_stats us ON (((us.tier)::text = t.tier)));

CREATE OR REPLACE VIEW monthly_profit_trend AS
WITH months AS (
         SELECT to_char(d.d, 'YYYY-MM'::text) AS month,
            date_trunc('month'::text, d.d) AS m_start
           FROM generate_series((now() - '11 mons'::interval), now(), '1 mon'::interval) d(d)
        ), tx_m AS (
         SELECT to_char(transactions.created_at, 'YYYY-MM'::text) AS month,
            count(transactions.id) AS tx_count,
            COALESCE(sum(transactions.amount), (0)::numeric) AS volume,
            COALESCE(sum(transactions.broker_fee), (0)::numeric) AS broker_fees
           FROM transactions
          WHERE ((transactions.status)::text = 'completed'::text)
          GROUP BY (to_char(transactions.created_at, 'YYYY-MM'::text))
        ), fee_m AS (
         SELECT to_char(products.created_at, 'YYYY-MM'::text) AS month,
            COALESCE(sum(products.total_fee_amount), (0)::numeric) AS listing_fees
           FROM products
          WHERE ((products.fee_paid = true) AND (products.status = ANY (ARRAY['approved'::product_status, 'expired'::product_status, 'sold'::product_status])))
          GROUP BY (to_char(products.created_at, 'YYYY-MM'::text))
        )
 SELECT m.month,
    (COALESCE(tx.tx_count, (0)::bigint))::integer AS transactions,
    (COALESCE(tx.volume, (0)::numeric))::double precision AS volume,
    (COALESCE(tx.broker_fees, (0)::numeric))::double precision AS broker_fees,
    (COALESCE(fm.listing_fees, (0)::numeric))::double precision AS listing_fees,
    ((COALESCE(tx.broker_fees, (0)::numeric) + COALESCE(fm.listing_fees, (0)::numeric)))::double precision AS total_profit,
        CASE
            WHEN (COALESCE(tx.volume, (0)::numeric) > (0)::numeric) THEN (round((((COALESCE(tx.broker_fees, (0)::numeric) + COALESCE(fm.listing_fees, (0)::numeric)) / tx.volume) * (100)::numeric), 2))::double precision
            ELSE (0.0)::double precision
        END AS avg_fee_pct
   FROM ((months m
     LEFT JOIN tx_m tx ON ((tx.month = m.month)))
     LEFT JOIN fee_m fm ON ((fm.month = m.month)))
  ORDER BY m.m_start DESC;

CREATE OR REPLACE VIEW weekly_profit_trend AS
WITH weeks AS (
         SELECT to_char(d.d, 'YYYY-"W"IW'::text) AS week,
            date_trunc('week'::text, d.d) AS week_start
           FROM generate_series((now() - '84 days'::interval), now(), '7 days'::interval) d(d)
        ), tx_w AS (
         SELECT to_char(transactions.created_at, 'YYYY-"W"IW'::text) AS week,
            count(transactions.id) AS tx_count,
            COALESCE(sum(transactions.amount), (0)::numeric) AS volume,
            COALESCE(sum(transactions.broker_fee), (0)::numeric) AS broker_fees
           FROM transactions
          WHERE ((transactions.status)::text = 'completed'::text)
          GROUP BY (to_char(transactions.created_at, 'YYYY-"W"IW'::text))
        ), fee_w AS (
         SELECT to_char(products.created_at, 'YYYY-"W"IW'::text) AS week,
            COALESCE(sum(products.total_fee_amount), (0)::numeric) AS listing_fees
           FROM products
          WHERE ((products.fee_paid = true) AND (products.status = ANY (ARRAY['approved'::product_status, 'expired'::product_status, 'sold'::product_status])))
          GROUP BY (to_char(products.created_at, 'YYYY-"W"IW'::text))
        )
 SELECT w.week,
    (COALESCE(tx.volume, (0)::numeric))::double precision AS volume,
    (COALESCE(tx.broker_fees, (0)::numeric))::double precision AS broker_fees,
    (COALESCE(fw.listing_fees, (0)::numeric))::double precision AS listing_fees,
    ((COALESCE(tx.broker_fees, (0)::numeric) + COALESCE(fw.listing_fees, (0)::numeric)))::double precision AS total_profit,
    (COALESCE(tx.tx_count, (0)::bigint))::integer AS transactions,
        CASE
            WHEN (COALESCE(tx.volume, (0)::numeric) > (0)::numeric) THEN (round((((COALESCE(tx.broker_fees, (0)::numeric) + COALESCE(fw.listing_fees, (0)::numeric)) / tx.volume) * (100)::numeric), 2))::double precision
            ELSE (0.0)::double precision
        END AS avg_fee_pct,
    w.week_start
   FROM ((weeks w
     LEFT JOIN tx_w tx ON ((tx.week = w.week)))
     LEFT JOIN fee_w fw ON ((fw.week = w.week)));

CREATE OR REPLACE VIEW annual_profit_trend AS
WITH years AS (
         SELECT (EXTRACT(year FROM d.d))::integer AS year
           FROM generate_series((now() - '4 years'::interval), now(), '1 year'::interval) d(d)
        ), tx_y AS (
         SELECT (EXTRACT(year FROM transactions.created_at))::integer AS year,
            count(transactions.id) AS tx_count,
            COALESCE(sum(transactions.amount), (0)::numeric) AS volume,
            COALESCE(sum(transactions.broker_fee), (0)::numeric) AS broker_fees
           FROM transactions
          WHERE ((transactions.status)::text = 'completed'::text)
          GROUP BY ((EXTRACT(year FROM transactions.created_at))::integer)
        ), fee_y AS (
         SELECT (EXTRACT(year FROM products.created_at))::integer AS year,
            COALESCE(sum(products.total_fee_amount), (0)::numeric) AS listing_fees
           FROM products
          WHERE ((products.fee_paid = true) AND (products.status = ANY (ARRAY['approved'::product_status, 'expired'::product_status, 'sold'::product_status])))
          GROUP BY ((EXTRACT(year FROM products.created_at))::integer)
        )
 SELECT y.year,
    (COALESCE(ty.tx_count, (0)::bigint))::integer AS transactions,
    (COALESCE(ty.volume, (0)::numeric))::double precision AS volume,
    (COALESCE(ty.broker_fees, (0)::numeric))::double precision AS broker_fees,
    (COALESCE(fy.listing_fees, (0)::numeric))::double precision AS listing_fees,
    ((COALESCE(ty.broker_fees, (0)::numeric) + COALESCE(fy.listing_fees, (0)::numeric)))::double precision AS total_profit,
        CASE
            WHEN (COALESCE(ty.volume, (0)::numeric) > (0)::numeric) THEN (round((((COALESCE(ty.broker_fees, (0)::numeric) + COALESCE(fy.listing_fees, (0)::numeric)) / ty.volume) * (100)::numeric), 2))::double precision
            ELSE (0.0)::double precision
        END AS avg_fee_pct
   FROM ((years y
     LEFT JOIN tx_y ty ON ((ty.year = y.year)))
     LEFT JOIN fee_y fy ON ((fy.year = y.year)));
