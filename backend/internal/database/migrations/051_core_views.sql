-- ============================================================================
-- 051_core_views.sql
-- Consolidated System & Analytic Views
-- ============================================================================

-- Backward-compatible categories view over category_hierarchy
CREATE OR REPLACE VIEW categories AS
SELECT id, name_en AS name, name_am, slug, icon, is_active, display_order, created_at, updated_at
FROM category_hierarchy;

-- ============================================================================
-- 068_core_views.sql
-- Core platform views and underlying helper functions
-- ============================================================================

CREATE OR REPLACE FUNCTION public.get_category_ancestors(p_category_id integer)
 RETURNS TABLE(id integer, name_en character varying, slug character varying, level integer)
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    WITH RECURSIVE ancestors AS (
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        WHERE ch.id = p_category_id
        
        UNION ALL
        
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        JOIN ancestors a ON ch.id = a.parent_id
    )
    SELECT a.id, a.name_en, a.slug, a.level
    FROM ancestors a
    ORDER BY a.level;
END;
$function$;

CREATE OR REPLACE FUNCTION public.get_category_descendants(p_category_id integer)
 RETURNS TABLE(id integer, name_en character varying, slug character varying, level integer)
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    WITH RECURSIVE descendants AS (
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        WHERE ch.id = p_category_id
        
        UNION ALL
        
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        JOIN descendants d ON ch.parent_id = d.id
    )
    SELECT d.id, d.name_en, d.slug, d.level
    FROM descendants d
    ORDER BY d.level;
END;
$function$;

CREATE OR REPLACE FUNCTION public.get_category_full_name(p_category_id integer)
 RETURNS character varying
 LANGUAGE plpgsql
AS $function$
DECLARE
    v_full_name VARCHAR;
BEGIN
    SELECT ARRAY_TO_STRING(breadcrumb_names, ' > ')
    INTO v_full_name
    FROM category_breadcrumbs
    WHERE id = p_category_id;
    
    RETURN COALESCE(v_full_name, '');
END;
$function$;

CREATE OR REPLACE FUNCTION public.get_category_attributes(p_category_id integer)
 RETURNS TABLE(attribute_id integer, name_en character varying, slug character varying, data_type character varying, validation_rules jsonb, unit character varying, is_required boolean, is_filterable boolean, display_order integer)
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    SELECT 
        attr_inh.attribute_id,
        attr_inh.name_en,
        attr_inh.slug,
        attr_inh.data_type,
        attr_inh.validation_rules,
        attr_inh.unit,
        attr_inh.is_required,
        attr_inh.is_filterable,
        attr_inh.display_order
    FROM get_category_attributes_inherited(p_category_id) attr_inh;
END;
$function$;

CREATE OR REPLACE FUNCTION public.update_attribute_timestamp()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.update_product_attribute_search_vector()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    IF NEW.value_text IS NOT NULL THEN
        NEW.search_vector := to_tsvector('simple', COALESCE(NEW.value_text, ''));
    END IF;
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.user_has_capability(p_user_id integer, p_capability character varying)
 RETURNS boolean
 LANGUAGE plpgsql
 STABLE
AS $function$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM user_capabilities
        WHERE user_id = p_user_id
          AND capability = p_capability
          AND status = 'active'
    );
END;
$function$;

CREATE OR REPLACE FUNCTION public.get_user_capabilities(p_user_id integer)
 RETURNS text[]
 LANGUAGE plpgsql
 STABLE
AS $function$
BEGIN
    RETURN ARRAY(
        SELECT capability FROM user_capabilities
        WHERE user_id = p_user_id AND status = 'active'
        ORDER BY capability ASC
    );
END;
$function$;

DROP TRIGGER IF EXISTS trg_update_attribute_definitions_timestamp ON attribute_definitions;
CREATE TRIGGER trg_update_attribute_definitions_timestamp BEFORE UPDATE ON attribute_definitions FOR EACH ROW EXECUTE FUNCTION update_attribute_timestamp();

DROP TRIGGER IF EXISTS trg_update_product_attribute_search ON product_attributes;
CREATE TRIGGER trg_update_product_attribute_search BEFORE UPDATE ON product_attributes FOR EACH ROW EXECUTE FUNCTION update_product_attribute_search_vector();

CREATE OR REPLACE VIEW referral_balance_view AS
SELECT u.id AS user_id,
    COALESCE(sum(rc.commission_amount), (0)::numeric) AS total_earned,
    COALESCE(( SELECT sum(rw.amount) AS sum
           FROM referral_withdrawals rw
          WHERE ((rw.user_id = u.id) AND ((rw.status)::text = ANY ((ARRAY['approved'::character varying, 'paid'::character varying])::text[])))), (0)::numeric) AS total_withdrawn,
    (COALESCE(sum(rc.commission_amount), (0)::numeric) - COALESCE(( SELECT sum(rw.amount) AS sum
           FROM referral_withdrawals rw
          WHERE ((rw.user_id = u.id) AND ((rw.status)::text = ANY ((ARRAY['approved'::character varying, 'paid'::character varying, 'pending'::character varying])::text[])))), (0)::numeric)) AS available_balance
   FROM (users u
     LEFT JOIN referral_commissions rc ON (((rc.referrer_id = u.id) AND ((rc.status)::text = 'confirmed'::text))))
  GROUP BY u.id;

CREATE OR REPLACE VIEW otp_log_view AS
SELECT id,
    COALESCE(phone, ''::character varying) AS phone,
    COALESCE(email, ''::character varying) AS email,
    purpose,
    user_type,
    used,
    attempts,
    expires_at,
    created_at,
    status,
    COALESCE(code_encrypted, code, ''::text) AS encrypted_code
   FROM otp_codes;

CREATE OR REPLACE VIEW category_breadcrumbs AS
WITH RECURSIVE category_path AS (
         SELECT category_hierarchy.id,
            category_hierarchy.parent_id,
            category_hierarchy.name_en,
            category_hierarchy.slug,
            category_hierarchy.level,
            category_hierarchy.path,
            ARRAY[(category_hierarchy.name_en)::text] AS breadcrumb_names,
            ARRAY[(category_hierarchy.slug)::text] AS breadcrumb_slugs,
            ARRAY[category_hierarchy.id] AS breadcrumb_ids
           FROM category_hierarchy
          WHERE (category_hierarchy.parent_id IS NULL)
        UNION ALL
         SELECT ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level,
            ch.path,
            (cp.breadcrumb_names || (ch.name_en)::text),
            (cp.breadcrumb_slugs || (ch.slug)::text),
            (cp.breadcrumb_ids || ch.id)
           FROM (category_hierarchy ch
             JOIN category_path cp ON ((ch.parent_id = cp.id)))
        )
 SELECT id,
    name_en,
    slug,
    level,
    breadcrumb_names,
    breadcrumb_slugs,
    breadcrumb_ids,
    array_length(breadcrumb_names, 1) AS depth
   FROM category_path;

CREATE OR REPLACE VIEW category_stats AS
SELECT ch.id,
    ch.id AS category_id,
    ch.parent_id,
    ch.name_en,
    ch.slug,
    ch.level,
    ch.is_leaf,
    count(p.id) AS product_count,
    count(p.id) FILTER (WHERE (p.status = 'approved'::product_status)) AS approved_count,
    count(p.id) FILTER (WHERE (p.status = 'pending'::product_status)) AS pending_count
   FROM (category_hierarchy ch
     LEFT JOIN products p ON ((p.category_hierarchy_id = ch.id)))
  GROUP BY ch.id, ch.parent_id, ch.name_en, ch.slug, ch.level, ch.is_leaf;

CREATE OR REPLACE VIEW category_applicable_attributes AS
SELECT ca.category_id,
    ch.name_en AS category_name,
    ch.slug AS category_slug,
    ad.id AS attribute_id,
    ad.name_en AS attribute_name,
    ad.slug AS attribute_slug,
    ad.data_type,
    ad.validation_rules,
    ad.unit,
    ca.is_required,
    ca.display_order,
    ad.is_filterable,
    ad.is_searchable
   FROM ((category_attributes ca
     JOIN category_hierarchy ch ON ((ch.id = ca.category_id)))
     JOIN attribute_definitions ad ON ((ad.id = ca.attribute_id)))
  ORDER BY ca.category_id, ca.display_order;

CREATE OR REPLACE VIEW user_capabilities_view AS
SELECT u.id AS user_id,
    u.name,
    u.phone,
    u.email,
    u.role,
    u.tier AS seller_tier,
    u.buyer_tier,
    u.admin_role,
    u.is_active,
    COALESCE(array_agg(uc.capability) FILTER (WHERE ((uc.status)::text = 'active'::text)), ARRAY['buyer'::character varying]) AS active_capabilities
   FROM (users u
     LEFT JOIN user_capabilities uc ON ((uc.user_id = u.id)))
  GROUP BY u.id;

-- ============================================================================
-- 069_operations_views.sql
-- Operations, monitoring, and job execution summary views
-- ============================================================================

CREATE OR REPLACE VIEW operations_health AS
SELECT ( SELECT count(*) AS count
           FROM job_execution_log
          WHERE (job_execution_log.started_at > (now() - '24:00:00'::interval))) AS total_jobs_24h,
    ( SELECT count(*) AS count
           FROM job_execution_log
          WHERE (((job_execution_log.status)::text = 'failed'::text) AND (job_execution_log.started_at > (now() - '24:00:00'::interval)))) AS failed_jobs_24h,
    ( SELECT count(*) AS count
           FROM job_execution_log
          WHERE (((job_execution_log.status)::text = 'running'::text) AND (job_execution_log.started_at < (now() - '01:00:00'::interval)))) AS stuck_jobs,
    ( SELECT count(*) AS count
           FROM users
          WHERE (users.is_active = false)) AS inactive_users,
    ( SELECT count(*) AS count
           FROM products
          WHERE ((products.status = 'pending'::product_status) AND (products.created_at < (now() - '7 days'::interval)))) AS stale_pending_products,
    ( SELECT count(*) AS count
           FROM otp_codes
          WHERE ((otp_codes.expires_at < now()) AND ((otp_codes.status)::text = 'pending'::text))) AS expired_otps,
    ( SELECT count(*) AS count
           FROM upgrade_requests
          WHERE (((upgrade_requests.status)::text = 'pending'::text) AND (upgrade_requests.created_at < (now() - '14 days'::interval)))) AS stale_upgrade_requests,
    ( SELECT pg_size_pretty(pg_database_size(current_database())) AS pg_size_pretty) AS database_size,
    ( SELECT pg_size_pretty(pg_total_relation_size('products'::regclass)) AS pg_size_pretty) AS products_table_size,
    ( SELECT pg_size_pretty(pg_total_relation_size('users'::regclass)) AS pg_size_pretty) AS users_table_size;

CREATE OR REPLACE VIEW job_status_summary AS
SELECT job_name,
    count(*) AS total_runs,
    count(*) FILTER (WHERE ((status)::text = 'success'::text)) AS successful_runs,
    count(*) FILTER (WHERE ((status)::text = 'failed'::text)) AS failed_runs,
    count(*) FILTER (WHERE ((status)::text = 'running'::text)) AS currently_running,
    max(started_at) AS last_run_at,
    max(
        CASE
            WHEN ((status)::text = 'success'::text) THEN started_at
            ELSE NULL::timestamp with time zone
        END) AS last_success_at,
    max(
        CASE
            WHEN ((status)::text = 'failed'::text) THEN started_at
            ELSE NULL::timestamp with time zone
        END) AS last_failure_at,
    avg(duration_ms) FILTER (WHERE ((status)::text = 'success'::text)) AS avg_duration_ms,
    sum(records_processed) AS total_records_processed
   FROM job_execution_log
  WHERE (started_at > (now() - '7 days'::interval))
  GROUP BY job_name
  ORDER BY (max(started_at)) DESC;

CREATE OR REPLACE VIEW recent_failed_jobs AS
SELECT id,
    job_name,
    started_at,
    completed_at,
    duration_ms,
    records_processed,
    error_message,
    created_at
   FROM job_execution_log
  WHERE (((status)::text = 'failed'::text) AND (started_at > (now() - '7 days'::interval)))
  ORDER BY started_at DESC
 LIMIT 100;

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
