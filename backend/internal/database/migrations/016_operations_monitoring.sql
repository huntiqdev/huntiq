-- ============================================================
-- 016_operations_monitoring.sql
-- PURPOSE: Operations monitoring and job tracking for Platform Operations Admin (Role 02)
-- ============================================================

-- ── Job Execution Log ─────────────────────────────────────────────────────────
-- Tracks background job runs for monitoring and troubleshooting
CREATE TABLE IF NOT EXISTS job_execution_log (
    id SERIAL PRIMARY KEY,
    job_name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'running', -- 'success', 'failed', 'running'
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    duration_ms INT,
    records_processed INT DEFAULT 0,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_job_log_name_status ON job_execution_log(job_name, status, started_at DESC);
CREATE INDEX IF NOT EXISTS idx_job_log_started ON job_execution_log(started_at DESC);

COMMENT ON TABLE job_execution_log IS 'Background job execution history for operations monitoring';

-- ── System Health Metrics ─────────────────────────────────────────────────────
-- Aggregated view of system health for operations dashboard
CREATE OR REPLACE VIEW operations_health AS
SELECT
    -- Job Status Summary
    (SELECT COUNT(*) FROM job_execution_log WHERE started_at > NOW() - INTERVAL '24 hours') AS total_jobs_24h,
    (SELECT COUNT(*) FROM job_execution_log WHERE status='failed' AND started_at > NOW() - INTERVAL '24 hours') AS failed_jobs_24h,
    (SELECT COUNT(*) FROM job_execution_log WHERE status='running' AND started_at < NOW() - INTERVAL '1 hour') AS stuck_jobs,
    
    -- Data Health
    (SELECT COUNT(*) FROM users WHERE is_active=FALSE) AS inactive_users,
    (SELECT COUNT(*) FROM products WHERE status='pending' AND created_at < NOW() - INTERVAL '7 days') AS stale_pending_products,
    (SELECT COUNT(*) FROM otp_codes WHERE expires_at < NOW() AND status='pending') AS expired_otps,
    (SELECT COUNT(*) FROM upgrade_requests WHERE status='pending' AND created_at < NOW() - INTERVAL '14 days') AS stale_upgrade_requests,
    
    -- Database Size
    (SELECT pg_size_pretty(pg_database_size(current_database()))) AS database_size,
    (SELECT pg_size_pretty(pg_total_relation_size('products'))) AS products_table_size,
    (SELECT pg_size_pretty(pg_total_relation_size('users'))) AS users_table_size;

COMMENT ON VIEW operations_health IS 'System health metrics for operations dashboard';

-- ── Job Status Summary View ───────────────────────────────────────────────────
-- Aggregates job execution stats by job name
CREATE OR REPLACE VIEW job_status_summary AS
SELECT
    job_name,
    COUNT(*) AS total_runs,
    COUNT(*) FILTER (WHERE status='success') AS successful_runs,
    COUNT(*) FILTER (WHERE status='failed') AS failed_runs,
    COUNT(*) FILTER (WHERE status='running') AS currently_running,
    MAX(started_at) AS last_run_at,
    MAX(CASE WHEN status='success' THEN started_at END) AS last_success_at,
    MAX(CASE WHEN status='failed' THEN started_at END) AS last_failure_at,
    AVG(duration_ms) FILTER (WHERE status='success') AS avg_duration_ms,
    SUM(records_processed) AS total_records_processed
FROM job_execution_log
WHERE started_at > NOW() - INTERVAL '7 days'
GROUP BY job_name
ORDER BY last_run_at DESC;

COMMENT ON VIEW job_status_summary IS 'Job execution statistics for the last 7 days';

-- ── Failed Jobs View ──────────────────────────────────────────────────────────
-- Recent failed jobs for troubleshooting
CREATE OR REPLACE VIEW recent_failed_jobs AS
SELECT
    id,
    job_name,
    started_at,
    completed_at,
    duration_ms,
    records_processed,
    error_message,
    created_at
FROM job_execution_log
WHERE status = 'failed'
  AND started_at > NOW() - INTERVAL '7 days'
ORDER BY started_at DESC
LIMIT 100;

COMMENT ON VIEW recent_failed_jobs IS 'Failed jobs in the last 7 days for troubleshooting';

-- ── Seller Verification Documents (for Role 08: Seller Verification Admin) ────
CREATE TABLE IF NOT EXISTS seller_verification_documents (
    id SERIAL PRIMARY KEY,
    seller_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL, -- 'national_id', 'business_license', 'trade_permit', 'tin_certificate'
    file_url TEXT NOT NULL,
    file_name VARCHAR(255),
    file_size_bytes INT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected'
    reviewed_by INT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    reject_reason TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_verification_docs_seller ON seller_verification_documents(seller_id, status);
CREATE INDEX IF NOT EXISTS idx_verification_docs_status ON seller_verification_documents(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_verification_docs_reviewed ON seller_verification_documents(reviewed_by) WHERE reviewed_by IS NOT NULL;

COMMENT ON TABLE seller_verification_documents IS 'Seller identity and business verification documents';

-- Trigger for document timestamps
DROP TRIGGER IF EXISTS update_verification_docs_updated_at ON seller_verification_documents;
CREATE TRIGGER update_verification_docs_updated_at BEFORE UPDATE ON seller_verification_documents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ── Abuse Reports (for Role 21: Fraud & Trust & Safety Admin) ────────────────
CREATE TABLE IF NOT EXISTS abuse_reports (
    id SERIAL PRIMARY KEY,
    reporter_id INT REFERENCES users(id) ON DELETE SET NULL,
    reporter_role VARCHAR(20), -- 'seller', 'buyer', 'admin', 'anonymous'
    reported_user_id INT REFERENCES users(id) ON DELETE CASCADE,
    reported_product_id INT REFERENCES products(id) ON DELETE CASCADE,
    reason VARCHAR(50) NOT NULL, -- 'fraud', 'spam', 'inappropriate', 'counterfeit', 'harassment', 'other'
    description TEXT NOT NULL,
    evidence_urls TEXT[], -- Array of screenshot/evidence URLs
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'investigating', 'resolved', 'dismissed'
    priority VARCHAR(20) NOT NULL DEFAULT 'normal', -- 'low', 'normal', 'high', 'critical'
    reviewed_by INT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    action_taken TEXT,
    resolution_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_abuse_reports_status ON abuse_reports(status, priority DESC, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_abuse_reports_reported_user ON abuse_reports(reported_user_id) WHERE reported_user_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_abuse_reports_reported_product ON abuse_reports(reported_product_id) WHERE reported_product_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_abuse_reports_reporter ON abuse_reports(reporter_id) WHERE reporter_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_abuse_reports_reviewer ON abuse_reports(reviewed_by) WHERE reviewed_by IS NOT NULL;

COMMENT ON TABLE abuse_reports IS 'User-reported abuse, fraud, and safety incidents';

-- Trigger for abuse report timestamps
DROP TRIGGER IF EXISTS update_abuse_reports_updated_at ON abuse_reports;
CREATE TRIGGER update_abuse_reports_updated_at BEFORE UPDATE ON abuse_reports
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ── Refund Requests (for Role 20: Refund & Dispute Admin) ────────────────────
CREATE TABLE IF NOT EXISTS refund_requests (
    id SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE SET NULL,
    requester_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    requester_role VARCHAR(20) NOT NULL, -- 'seller', 'buyer'
    amount NUMERIC(12,2) NOT NULL CHECK (amount > 0),
    reason VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    evidence_urls TEXT[], -- Array of evidence URLs
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected', 'processed'
    reviewed_by INT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    reject_reason TEXT,
    refund_reference VARCHAR(100), -- Payment gateway refund reference
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refund_requests_status ON refund_requests(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_refund_requests_requester ON refund_requests(requester_id, status);
CREATE INDEX IF NOT EXISTS idx_refund_requests_transaction ON refund_requests(transaction_id) WHERE transaction_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_refund_requests_reviewer ON refund_requests(reviewed_by) WHERE reviewed_by IS NOT NULL;

COMMENT ON TABLE refund_requests IS 'Refund and dispute requests for transactions';

-- Trigger for refund request timestamps
DROP TRIGGER IF EXISTS update_refund_requests_updated_at ON refund_requests;
CREATE TRIGGER update_refund_requests_updated_at BEFORE UPDATE ON refund_requests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ── Account Appeals (for Role 30: Appeals & Account Recovery Admin) ───────────
CREATE TABLE IF NOT EXISTS account_appeals (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    appeal_type VARCHAR(50) NOT NULL, -- 'ban', 'deactivation', 'listing_rejection', 'verification_rejection', 'warning'
    related_product_id INT REFERENCES products(id) ON DELETE SET NULL,
    related_document_id INT REFERENCES seller_verification_documents(id) ON DELETE SET NULL,
    reason TEXT NOT NULL,
    evidence_urls TEXT[],
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'under_review', 'approved', 'rejected'
    reviewed_by INT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    decision TEXT,
    action_taken TEXT, -- What was done (account reactivated, listing approved, etc.)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_appeals_status ON account_appeals(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_appeals_user ON account_appeals(user_id, status);
CREATE INDEX IF NOT EXISTS idx_appeals_type ON account_appeals(appeal_type, status);
CREATE INDEX IF NOT EXISTS idx_appeals_reviewer ON account_appeals(reviewed_by) WHERE reviewed_by IS NOT NULL;

COMMENT ON TABLE account_appeals IS 'User appeals for account actions and listing decisions';

-- Trigger for appeal timestamps
DROP TRIGGER IF EXISTS update_appeals_updated_at ON account_appeals;
CREATE TRIGGER update_appeals_updated_at BEFORE UPDATE ON account_appeals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ── Update admin_permissions with new permissions ──────────────────────────────
-- Add new permission keys to AllAdminPermissions in models/user.go:
-- 'operations', 'seller_verification', 'trust_safety', 'refunds', 'appeals'

COMMENT ON TABLE admin_permissions IS 'Granular admin permissions (now includes: operations, seller_verification, trust_safety, refunds, appeals)';

-- ── Grant default permissions to existing normal admins ───────────────────────
-- New permissions are off by default (sensitive operations)
-- Super admins can grant these via /admin/admins/:id/permissions

-- ══════════════════════════════════════════════════════════════════════════════
-- MIGRATION COMPLETE
-- ══════════════════════════════════════════════════════════════════════════════

DO $$
BEGIN
    RAISE NOTICE '========================================';
    RAISE NOTICE 'MIGRATION 016 COMPLETE!';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Added tables:';
    RAISE NOTICE '  - job_execution_log (operations monitoring)';
    RAISE NOTICE '  - seller_verification_documents (Role 08)';
    RAISE NOTICE '  - abuse_reports (Role 21)';
    RAISE NOTICE '  - refund_requests (Role 20)';
    RAISE NOTICE '  - account_appeals (Role 30)';
    RAISE NOTICE 'Added views:';
    RAISE NOTICE '  - operations_health';
    RAISE NOTICE '  - job_status_summary';
    RAISE NOTICE '  - recent_failed_jobs';
    RAISE NOTICE 'New admin permissions to add to models/user.go:';
    RAISE NOTICE '  - operations, seller_verification, trust_safety, refunds, appeals';
    RAISE NOTICE '========================================';
END $$;
