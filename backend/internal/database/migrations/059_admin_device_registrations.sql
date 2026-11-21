-- ============================================================================
-- 059_admin_device_registrations.sql
-- Multi-Computer Device Identity & Authorized Workstation Management
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

CREATE TABLE IF NOT EXISTS admin_device_registrations (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(64) NOT NULL UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_name VARCHAR(100) NOT NULL,
    browser VARCHAR(50),
    os VARCHAR(50),
    device_type VARCHAR(30) DEFAULT 'desktop' NOT NULL, -- desktop, laptop, mobile, workstation
    credential_id BYTEA, -- optional passkey credential bound to this device
    public_key BYTEA,
    status VARCHAR(20) DEFAULT 'active' NOT NULL, -- active, pending_verification, revoked, suspended
    is_trusted BOOLEAN DEFAULT TRUE NOT NULL,
    last_ip VARCHAR(50),
    last_active_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    registered_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    revoked_at TIMESTAMPTZ,
    revocation_reason TEXT
);

CREATE INDEX IF NOT EXISTS idx_device_reg_user ON admin_device_registrations(user_id) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_device_reg_id ON admin_device_registrations(device_id);

-- Link admin_sessions to authorized device registration
ALTER TABLE admin_sessions ADD COLUMN IF NOT EXISTS device_id VARCHAR(64);
CREATE INDEX IF NOT EXISTS idx_admin_sessions_device ON admin_sessions(device_id) WHERE revoked_at IS NULL;
