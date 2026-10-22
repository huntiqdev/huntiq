-- ============================================================
-- 012_security_and_otp.sql
-- CONSOLIDATED: otp_codes + user_activity_logs + invitation_requests
-- PURPOSE: Security, authentication, and OTP system
-- NOTE: Uses ALTER TABLE to extend shell tables from 001_schema.sql
-- ============================================================

-- ── OTP Codes ───────────────────────────────────────────────────────────────
ALTER TABLE otp_codes
    ADD COLUMN IF NOT EXISTS phone          VARCHAR(20),
    ADD COLUMN IF NOT EXISTS email          VARCHAR(150),
    ADD COLUMN IF NOT EXISTS code_encrypted TEXT          NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS purpose        VARCHAR(50)   NOT NULL DEFAULT 'register',
    ADD COLUMN IF NOT EXISTS user_type      VARCHAR(20)   NOT NULL DEFAULT 'seller',
    ADD COLUMN IF NOT EXISTS attempts       INT           NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS max_attempts   INT           NOT NULL DEFAULT 5,
    ADD COLUMN IF NOT EXISTS status         VARCHAR(20)   NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS expires_at     TIMESTAMPTZ   NOT NULL DEFAULT (NOW() + INTERVAL '5 minutes'),
    ADD COLUMN IF NOT EXISTS verified_at    TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- Add contact check constraint
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name='otp_codes' AND constraint_name='otp_contact_check'
    ) THEN
        ALTER TABLE otp_codes ADD CONSTRAINT otp_contact_check
            CHECK (phone IS NOT NULL OR email IS NOT NULL);
    END IF;
END $$;

-- ── User Activity Logs ──────────────────────────────────────────────────────
ALTER TABLE user_activity_logs
    ADD COLUMN IF NOT EXISTS user_id        INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS event          VARCHAR(50)   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS phone          VARCHAR(20),
    ADD COLUMN IF NOT EXISTS email          VARCHAR(150),
    ADD COLUMN IF NOT EXISTS user_type      VARCHAR(20),
    ADD COLUMN IF NOT EXISTS ip_address     VARCHAR(50),
    ADD COLUMN IF NOT EXISTS user_agent     TEXT,
    ADD COLUMN IF NOT EXISTS success        BOOLEAN       NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS failure_reason TEXT,
    ADD COLUMN IF NOT EXISTS metadata       JSONB,
    ADD COLUMN IF NOT EXISTS meta           TEXT,
    ADD COLUMN IF NOT EXISTS created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Invitation Requests ─────────────────────────────────────────────────────
ALTER TABLE invitation_requests
    ADD COLUMN IF NOT EXISTS email           VARCHAR(150)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS phone           VARCHAR(20)   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS name            VARCHAR(100)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS message         TEXT,
    ADD COLUMN IF NOT EXISTS status          VARCHAR(20)   NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS invited_by      INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS approved_by     INT           REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS approved_at     TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS invitation_code VARCHAR(100)  UNIQUE,
    ADD COLUMN IF NOT EXISTS used_at         TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW();

-- ── Indexes ─────────────────────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_otp_phone         ON otp_codes(phone, status, expires_at DESC);
CREATE INDEX IF NOT EXISTS idx_otp_email         ON otp_codes(email, status, expires_at DESC);
CREATE INDEX IF NOT EXISTS idx_otp_status        ON otp_codes(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activity_logs_user  ON user_activity_logs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activity_logs_event ON user_activity_logs(event, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activity_logs_phone ON user_activity_logs(phone, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activity_logs_ip    ON user_activity_logs(ip_address, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_invitation_status   ON invitation_requests(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_invitation_code     ON invitation_requests(invitation_code) WHERE invitation_code IS NOT NULL;

-- ── Cleanup Functions ───────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION cleanup_expired_otps()
RETURNS void AS $$
BEGIN
    DELETE FROM otp_codes WHERE status IN ('verified', 'expired') AND created_at < NOW() - INTERVAL '7 days';
    DELETE FROM otp_codes WHERE status = 'locked' AND created_at < NOW() - INTERVAL '24 hours';
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION cleanup_old_activity_logs()
RETURNS void AS $$
BEGIN
    DELETE FROM user_activity_logs WHERE created_at < NOW() - INTERVAL '6 months';
END;
$$ LANGUAGE plpgsql;
