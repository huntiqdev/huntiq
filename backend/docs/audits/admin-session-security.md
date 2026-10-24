# HUNTIQ Administrative Session Security & Zero-Trust Protocol

**Document:** `docs/audits/admin-session-security.md`  
**Date:** October 2026  
**Status:** Architecture Specification  
**Repository:** `/home/poven/GO PROJECT/sellbuy`  

---

## 1. Zero-Trust Session Model

Under HUNTIQ's Zero-Trust architecture, no administrative session is implicitly trusted simply because a client holds a static cryptographic token or IP address. Every privileged request is authorized dynamically on the Go backend against active database state.

### 1.1 Dual Authentication Handshake
Administrative sessions support two concurrent, defense-in-depth transport modes:
1. **HttpOnly, Secure, SameSite Cookie (`hm_admin_session`):**
   - Inaccessible to client JavaScript (prevents Cross-Site Scripting [XSS] exfiltration).
   - `SameSite=Lax` (or `Strict` for critical administrative consoles) to block Cross-Site Request Forgery (CSRF).
   - `Secure=true` in production (enforcing HTTPS transmission).
2. **Authorization Header (`Bearer <session_token>`):**
   - Transmitted via standard header for headless CLI tools, mobile administrative apps, and programmatic integration.

---

## 2. Server-Side Session Tracking (`admin_sessions`)

Unlike legacy stateless JWTs, administrative sessions are recorded in `admin_sessions`:

| Field | Type | Purpose |
|---|---|---|
| `session_token_hash` | `VARCHAR(64)` | SHA-256 hash of high-entropy (256-bit) cryptographically random session secret. |
| `user_id` | `INT` | Foreign key referencing `users(id)`. |
| `portal` | `VARCHAR(50)` | Portal scope: `'executive'`, `'platform'`, `'security'`, `'staff'`. |
| `role` | `VARCHAR(30)` | Active administrative role identifier at time of issuance. |
| `ip_address` | `VARCHAR(50)` | Source IP address at session creation. |
| `user_agent` | `TEXT` | Client User-Agent string. |
| `last_active_at` | `TIMESTAMPTZ` | Timestamp updated on requests to enforce idle timeout. |
| `idle_timeout_seconds` | `INT` | Maximum permitted inactivity (Executive: 900s / 15m; Staff: 1800s / 30m). |
| `absolute_expires_at` | `TIMESTAMPTZ` | Hard absolute expiration (Executive: 4 hours; Staff: 12 hours). |
| `revoked_at` | `TIMESTAMPTZ` | Nullable timestamp set upon manual logout, password change, or security revocation. |
| `revocation_reason` | `TEXT` | Reason string (`"logout"`, `"admin_demoted"`, `"session_stale"`, `"security_revocation"`). |

---

## 3. Session Lifecycle & Invalidation Triggers

A session is invalidated immediately upon any of the following triggers:
1. **Explicit Logout:** `POST /api/admin/logout` marks session revoked in database and clears cookies.
2. **Idle Expiration:** Inactivity exceeding `idle_timeout_seconds` causes automatic rejection on the next request.
3. **Absolute Expiration:** Expiration of `absolute_expires_at` terminates the session regardless of activity.
4. **Account Suspension or Status Change:** If `users.is_active` becomes `FALSE`, all active sessions for that user ID are immediately rejected.
5. **Role Modification / Demotion:** When an administrator's role or permissions are modified, all concurrent sessions are revoked, forcing re-authentication.
6. **Credential Revocation:** Revocation of a passkey or password reset terminates all existing sessions.
7. **Security Administrator Revocation:** Security Administrators have authority to invoke `POST /api/admin/security/revoke-session` or `POST /api/admin/security/revoke-all-sessions/:userId`.

---

## 4. Post-Implementation Security Audit Finding: JWT Coexistence & Enforcement

During the post-implementation adversarial audit, a critical boundary gap was identified:
* **Finding:** `RequireAuth` middleware allowed 24-hour JWTs (`role == "admin"`) to authorize requests without validating against `admin_sessions`.
* **Impact:** Invalidation of a server-side session in `admin_sessions` did not terminate an active 24-hour JWT presented in the `Authorization: Bearer` header.
* **Remediation Requirement:** Privileged `/admin/*` operations must strictly require validation against `admin_sessions` or fresh database verification rather than accepting static JWT claims.

