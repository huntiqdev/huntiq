-- ============================================================
-- 029_enterprise_admin_roles_and_security.sql
-- PURPOSE: Professional enterprise administrator role taxonomy,
--          WebAuthn passkey credential persistence, server-side
--          session management, and progressive abuse defense.
-- ============================================================

-- 1. Extend admin_role enum with enterprise taxonomy
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'platform_owner';
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'platform_admin';
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'security_admin';
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'operations_admin';
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'moderator';
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'support_agent';

-- 2. Role Migration Audit Table
CREATE TABLE IF NOT EXISTS admin_role_migration_audit (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    old_role VARCHAR(50) NOT NULL,
    new_role VARCHAR(50) NOT NULL,
    migrated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes TEXT
);

-- 3. Data Migration: Safely transition existing legacy admin accounts
-- ID 36 (or any existing 'super' admin) -> 'platform_owner'
INSERT INTO admin_role_migration_audit (user_id, old_role, new_role, notes)
SELECT id, admin_role::varchar, 'platform_owner', 'Automated legacy super admin migration to platform_owner'
FROM users
WHERE role = 'admin' AND admin_role = 'super';

UPDATE users
SET admin_role = 'platform_owner', updated_at = NOW()
WHERE role = 'admin' AND admin_role = 'super';

-- Any existing 'normal' admin -> 'operations_admin'
INSERT INTO admin_role_migration_audit (user_id, old_role, new_role, notes)
SELECT id, admin_role::varchar, 'operations_admin', 'Automated legacy normal admin migration to operations_admin'
FROM users
WHERE role = 'admin' AND admin_role = 'normal';

UPDATE users
SET admin_role = 'operations_admin', updated_at = NOW()
WHERE role = 'admin' AND admin_role = 'normal';

-- 4. Passkey Credentials Table (WebAuthn / FIDO2)
CREATE TABLE IF NOT EXISTS admin_passkeys (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    credential_id BYTEA NOT NULL UNIQUE,
    public_key BYTEA NOT NULL,
    attestation_type VARCHAR(50) NOT NULL DEFAULT 'none',
    aaguid BYTEA,
    sign_count BIGINT NOT NULL DEFAULT 0,
    clone_warning BOOLEAN NOT NULL DEFAULT FALSE,
    attachment VARCHAR(20) DEFAULT 'platform', -- 'platform' or 'cross-platform'
    name VARCHAR(100) NOT NULL DEFAULT 'Passkey',
    transports TEXT[] DEFAULT '{}',
    backup_eligible BOOLEAN DEFAULT FALSE,
    backup_state BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    revocation_reason TEXT
);

CREATE INDEX IF NOT EXISTS idx_admin_passkeys_user ON admin_passkeys (user_id) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_admin_passkeys_cred_id ON admin_passkeys (credential_id);

-- 5. Passkey Challenges Table (Server-Controlled Single-Use Challenges)
CREATE TABLE IF NOT EXISTS admin_auth_challenges (
    id SERIAL PRIMARY KEY,
    challenge VARCHAR(255) NOT NULL UNIQUE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    portal VARCHAR(50) NOT NULL,
    ceremony_type VARCHAR(20) NOT NULL, -- 'registration', 'authentication', 'reauth'
    session_data JSONB NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_admin_auth_challenges_lookup ON admin_auth_challenges (challenge) WHERE used_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_admin_auth_challenges_user ON admin_auth_challenges (user_id, ceremony_type) WHERE used_at IS NULL;

-- 6. Server-Side Session Tracking Table (Zero-Trust Session Invalidation)
CREATE TABLE IF NOT EXISTS admin_sessions (
    id SERIAL PRIMARY KEY,
    session_token_hash VARCHAR(64) NOT NULL UNIQUE, -- SHA-256 hash of random session secret
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    portal VARCHAR(50) NOT NULL,
    role VARCHAR(30) NOT NULL,
    ip_address VARCHAR(50),
    user_agent TEXT,
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    idle_timeout_seconds INT NOT NULL DEFAULT 1800, -- 15m for executive, 30m for staff
    absolute_expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    revocation_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_admin_sessions_token ON admin_sessions (session_token_hash) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_admin_sessions_user ON admin_sessions (user_id) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_admin_sessions_expiry ON admin_sessions (absolute_expires_at);

-- 7. Progressive Throttling & Abuse Monitoring Table
CREATE TABLE IF NOT EXISTS admin_auth_attempts (
    id SERIAL PRIMARY KEY,
    identifier VARCHAR(150) NOT NULL,
    ip_address VARCHAR(50) NOT NULL,
    portal VARCHAR(50) NOT NULL,
    success BOOLEAN NOT NULL DEFAULT FALSE,
    failure_reason VARCHAR(100),
    attempted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_admin_auth_attempts_ident ON admin_auth_attempts (identifier, attempted_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_auth_attempts_ip ON admin_auth_attempts (ip_address, attempted_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_auth_attempts_time ON admin_auth_attempts (attempted_at DESC);

-- 8. Sensitive Action Reauthentication Table
CREATE TABLE IF NOT EXISTS admin_reauth_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    action_scope VARCHAR(100) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_admin_reauth_tokens_lookup ON admin_reauth_tokens (token_hash) WHERE used_at IS NULL;
