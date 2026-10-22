-- ============================================================
-- 009_referral_system.sql
-- CONSOLIDATED: referral_commissions + referral_withdrawals
-- PURPOSE: Complete referral/commission system
-- NOTE: Uses ALTER TABLE to extend shell tables from 001_schema.sql
-- ============================================================

-- ── Referral Commissions ────────────────────────────────────────────────────
ALTER TABLE referral_commissions
    ADD COLUMN IF NOT EXISTS referrer_id         INT           NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS referred_id         INT           NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS transaction_id      INT           REFERENCES transactions(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS commission_amount   NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (commission_amount >= 0),
    ADD COLUMN IF NOT EXISTS commission_percent  NUMERIC(5,2)  NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS commission_rate     NUMERIC(5,2)  DEFAULT 0,
    ADD COLUMN IF NOT EXISTS transaction_type    VARCHAR(50)   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS transaction_amount  NUMERIC(12,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS status              VARCHAR(20)   NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS pending_until       TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS confirmed_at        TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS paid_at             TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS notes               TEXT,
    ADD COLUMN IF NOT EXISTS created_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Referral Withdrawals ────────────────────────────────────────────────────
ALTER TABLE referral_withdrawals
    ADD COLUMN IF NOT EXISTS user_id             INT           NOT NULL DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS amount              NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (amount > 0),
    ADD COLUMN IF NOT EXISTS withdrawal_method   VARCHAR(50)   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS payment_method      VARCHAR(50),
    ADD COLUMN IF NOT EXISTS account_info        VARCHAR(200)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS payer_phone         VARCHAR(20),
    ADD COLUMN IF NOT EXISTS payment_phone       VARCHAR(20),
    ADD COLUMN IF NOT EXISTS reference           VARCHAR(200),
    ADD COLUMN IF NOT EXISTS status              VARCHAR(20)   NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS approved_by         INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS reviewed_by         INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS approved_at         TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS reviewed_at         TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS paid_at             TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS rejection_reason    TEXT,
    ADD COLUMN IF NOT EXISTS rejection_note      TEXT,
    ADD COLUMN IF NOT EXISTS admin_note          TEXT,
    ADD COLUMN IF NOT EXISTS notes               TEXT,
    ADD COLUMN IF NOT EXISTS created_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Indexes ─────────────────────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_referral_commissions_referrer ON referral_commissions(referrer_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_referral_commissions_referred ON referral_commissions(referred_id);
CREATE INDEX IF NOT EXISTS idx_referral_commissions_status   ON referral_commissions(status);
CREATE INDEX IF NOT EXISTS idx_referral_commissions_pending  ON referral_commissions(status, created_at) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_referral_withdrawals_user     ON referral_withdrawals(user_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_referral_withdrawals_status   ON referral_withdrawals(status, created_at DESC);

-- ── Triggers ────────────────────────────────────────────────────────────────
DROP TRIGGER IF EXISTS update_referral_withdrawals_updated_at ON referral_withdrawals;
CREATE TRIGGER update_referral_withdrawals_updated_at BEFORE UPDATE ON referral_withdrawals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Auto-confirm commissions after 30 days (pending → confirmed)
CREATE OR REPLACE FUNCTION auto_confirm_commissions()
RETURNS void AS $$
BEGIN
    UPDATE referral_commissions
    SET status = 'confirmed', confirmed_at = NOW()
    WHERE status = 'pending'
      AND created_at < NOW() - INTERVAL '30 days';
END;
$$ LANGUAGE plpgsql;

-- Unique constraint to prevent double-recording commission for the same transaction
ALTER TABLE referral_commissions
    DROP CONSTRAINT IF EXISTS uq_referral_commission_source;
ALTER TABLE referral_commissions
    ADD CONSTRAINT uq_referral_commission_source UNIQUE (referred_id, transaction_id);

-- Unique index to prevent multiple pending withdrawals per user
DROP INDEX IF EXISTS uq_referral_withdrawals_one_pending;
CREATE UNIQUE INDEX uq_referral_withdrawals_one_pending 
    ON referral_withdrawals (user_id) 
    WHERE status = 'pending';

-- View to calculate referral balances
CREATE OR REPLACE VIEW referral_balance_view AS
SELECT
    u.id AS user_id,
    COALESCE(SUM(rc.commission_amount), 0) AS total_earned,
    COALESCE((SELECT SUM(rw.amount) FROM referral_withdrawals rw WHERE rw.user_id = u.id AND rw.status IN ('approved', 'paid')), 0) AS total_withdrawn,
    COALESCE(SUM(rc.commission_amount), 0) -
    COALESCE((SELECT SUM(rw.amount) FROM referral_withdrawals rw WHERE rw.user_id = u.id AND rw.status IN ('approved', 'paid', 'pending')), 0) AS available_balance
FROM users u
LEFT JOIN referral_commissions rc ON rc.referrer_id = u.id AND rc.status = 'confirmed'
GROUP BY u.id;
