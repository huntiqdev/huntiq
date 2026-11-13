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
