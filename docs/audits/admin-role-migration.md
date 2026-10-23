# HUNTIQ Enterprise Administrative Role Migration & Database Plan

**Document:** `docs/audits/admin-role-migration.md`  
**Date:** October 2026  
**Status:** Approved Migration Plan  
**Repository:** `/home/poven/GO PROJECT/sellbuy`  

---

## 1. Migration Purpose & Scope

The legacy HUNTIQ admin model divided administrative accounts into a coarse binary distinction: `super` (with unchecked bypass of all security checks) and `normal` (sub-admins managed by boolean permission rows).

This migration replaces the legacy `admin_role` ENUM with the enterprise taxonomy:
- `platform_owner`
- `platform_admin`
- `security_admin`
- `operations_admin`
- `moderator`
- `support_agent`

The migration is executed via Migration `029_enterprise_admin_roles_and_security.sql`.

---

## 2. Legacy-to-Professional Mapping Rules

| Legacy Role | Destination Role | Pre-Migration Verification Condition | Post-Migration Security Constraints |
|---|---|---|---|
| `super` | `platform_owner` | User ID must be verified as the primary system administrator (currently ID 36). | MFA bypass removed; passkey registration required. Immutable audit logging enabled. |
| `normal` | `operations_admin` | Default baseline for general operational staff. | Inherits existing `admin_permissions` mappings; cannot access executive or security consoles. |
| `normal` | `moderator` | For staff assigned solely catalog/content permissions (`pending`, `products`). | Scoped to `/staff/login`; restricted from viewing platform finance or settings. |
| `normal` | `support_agent` | For staff assigned dispute/customer assistance (`appeals`, `refunds`). | Scoped to customer assistance operations only. |

### 2.1 Explicit Safety Invariants

1. **No Accidental Promotion:** Ordinary marketplace users (`role = 'buyer'`, `role = 'seller'`) are untouched. Only accounts with `role = 'admin'` are inspected.
2. **Deterministic Identity Preservation:** Existing user IDs, email addresses, phone numbers, and creation timestamps remain unchanged.
3. **Audit History Preservation:** A dedicated audit table `admin_role_migration_audit` tracks every role change, timestamp, and operator.
4. **Idempotence & Replay Safety:** The migration script can be run multiple times without data corruption (`IF NOT EXISTS`, conditional schema updates).

---

## 3. Database Schema Changes (Migration 029)

### 3.1 Type Expansion & Table Alterations
```sql
-- 1. Add new enum values to admin_role
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'platform_owner';
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'platform_admin';
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'security_admin';
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'operations_admin';
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'moderator';
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'support_agent';

-- 2. Audit table for role transitions
CREATE TABLE IF NOT EXISTS admin_role_migration_audit (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    old_role VARCHAR(50) NOT NULL,
    new_role VARCHAR(50) NOT NULL,
    migrated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes TEXT
);

-- 3. Data backfill: Migrate ID 36 ('super') to 'platform_owner'
INSERT INTO admin_role_migration_audit (user_id, old_role, new_role, notes)
SELECT id, admin_role::varchar, 'platform_owner', 'Initial enterprise taxonomy migration'
FROM users
WHERE role = 'admin' AND admin_role = 'super'
ON CONFLICT DO NOTHING;

UPDATE users
SET admin_role = 'platform_owner'
WHERE role = 'admin' AND admin_role = 'super';

-- 4. Passkey credential storage
CREATE TABLE IF NOT EXISTS admin_passkeys (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    credential_id BYTEA NOT NULL UNIQUE,
    public_key BYTEA NOT NULL,
    attestation_type VARCHAR(50) NOT NULL DEFAULT 'none',
    aaguid BYTEA,
    sign_count BIGINT NOT NULL DEFAULT 0,
    clone_warning BOOLEAN NOT NULL DEFAULT FALSE,
    attachment VARCHAR(20) DEFAULT 'platform',
    name VARCHAR(100) NOT NULL DEFAULT 'Passkey',
    transports TEXT[] DEFAULT '{}',
    backup_eligible BOOLEAN DEFAULT FALSE,
    backup_state BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_admin_passkeys_user ON admin_passkeys(user_id) WHERE revoked_at IS NULL;

-- 5. Passkey auth challenge & ceremony session storage
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
CREATE INDEX IF NOT EXISTS idx_admin_auth_challenges_lookup ON admin_auth_challenges(challenge) WHERE used_at IS NULL;

-- 6. Server-side session tracking
CREATE TABLE IF NOT EXISTS admin_sessions (
    id SERIAL PRIMARY KEY,
    session_token_hash VARCHAR(64) NOT NULL UNIQUE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    portal VARCHAR(50) NOT NULL,
    role VARCHAR(30) NOT NULL,
    ip_address VARCHAR(50),
    user_agent TEXT,
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    idle_timeout_seconds INT NOT NULL DEFAULT 1800,
    absolute_expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    revocation_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_admin_sessions_token ON admin_sessions(session_token_hash) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_admin_sessions_user ON admin_sessions(user_id) WHERE revoked_at IS NULL;

-- 7. Progressive throttling & abuse monitoring
CREATE TABLE IF NOT EXISTS admin_auth_attempts (
    id SERIAL PRIMARY KEY,
    identifier VARCHAR(150) NOT NULL,
    ip_address VARCHAR(50) NOT NULL,
    portal VARCHAR(50) NOT NULL,
    success BOOLEAN NOT NULL DEFAULT FALSE,
    failure_reason VARCHAR(100),
    attempted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_admin_auth_attempts_ident ON admin_auth_attempts(identifier, attempted_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_auth_attempts_ip ON admin_auth_attempts(ip_address, attempted_at DESC);
```

---

## 4. Affected Files & Consumer Updates

| Layer | File Path | Migration Changes Required |
|---|---|---|
| **Go Models** | `backend/internal/models/user.go` | Define constants: `RolePlatformOwner`, `RolePlatformAdmin`, `RoleSecurityAdmin`, `RoleOperationsAdmin`, `RoleModerator`, `RoleSupportAgent`. |
| **Go JWT** | `backend/internal/security/jwt.go` | Keep `AdminRole` claim compatible; support new role values. |
| **Go Middleware** | `backend/internal/middleware/middleware.go` | Update `adminRoleRank` to support new roles: `platform_owner: 10`, `platform_admin: 8`, `security_admin: 8`, `operations_admin: 4`, `moderator: 2`, `support_agent: 2`. Backward-compatible fallback for legacy tokens. |
| **Go Handlers** | `backend/internal/handler/admin_handler.go` | Update `UpdateAdminRole`, `ListAdmins`, `CreateAdmin` to validate against new taxonomy. |
| **Go Services** | `backend/internal/service/admin_service.go` | Add portal-level and role-level authorization checks. |
| **Frontend Context** | `frontend/src/context/AuthContext.tsx` | Export `isPlatformOwner`, `isPlatformAdmin`, `isSecurityAdmin`, `isOperationsAdmin`. |
| **Frontend Layout** | `frontend/src/app/admin/layout.tsx` | Show current role title and scope navigation accordingly. |
| **Frontend Portals**| `frontend/src/app/(portals)` | Build dedicated login routes: `/executive-console/login`, `/platform-console/login`, `/security-console/login`, `/staff/login`. |

---

## 5. Rollback & Recovery Procedures

If a critical failure occurs during deployment:
1. **Database Rollback:**
   ```sql
   UPDATE users SET admin_role = 'super' WHERE admin_role = 'platform_owner';
   UPDATE users SET admin_role = 'normal' WHERE admin_role IN ('operations_admin', 'moderator', 'support_agent');
   ```
2. **Code Rollback:**
   Revert Go binary to previous release.
3. **Emergency Fallback Login:**
   Platform Owner emergency CLI command `/cmd/admin-recovery` can re-seed or verify passkeys without application downtime.

---

## 6. Post-Implementation Migration Verification

1. **Migration 029 Applied & Audited:** Verified in PostgreSQL that `admin_role_migration_audit` recorded transition of user ID 36 from `super` to `platform_owner`.
2. **Current Database State:** Confirmed single administrator record (ID 36) in `platform_owner` state with `is_active = true`.
3. **Session Invalidation Finding:** The audit identified that `UpdateAdminRole` updates `users.admin_role` in the database but does not invoke `adminAuthSvc.RevokeAllSessions`. Remediation: invoke session invalidation on role modification.

