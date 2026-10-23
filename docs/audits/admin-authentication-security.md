# HUNTIQ Enterprise Administrative Authentication & Security Audit

**Document:** `docs/audits/admin-authentication-security.md`  
**Date:** October 2026  
**Auditor:** Antigravity AI Security Pair Programmer  
**Repository:** `/home/poven/GO PROJECT/sellbuy`  
**Target Architecture:** Go (Fiber v2) Backend, PostgreSQL 14+, Next.js 15 (Turbopack) + TypeScript  

---

## 1. Executive Summary

A comprehensive read-only security and architectural audit of the HUNTIQ administrative authentication and access control subsystem was conducted. The existing codebase possesses a functioning foundation—including bcrypt password hashing (cost 12), transactional database schema, granular permission tables (`admin_permissions`), and audit logging (`admin_logs`).

However, several critical enterprise security deficiencies exist:
1. **Critical Privilege Inversion:** The highest privilege tier (`super` admin) completely bypasses second-factor authentication (MFA), obtaining a 24-hour JWT with a single email/password step (`backend/internal/service/auth_service.go:504-511`), while lower-tier sub-admins require email OTP.
2. **Absence of Phishing-Resistant Biometric/Passkey Authentication:** No WebAuthn/FIDO2 framework is currently implemented. Authentication relies on shared secrets (passwords) and symmetric email OTPs vulnerable to phishing, intercept, and SIM/email compromise.
3. **In-Memory, IP-Only Rate Limiting:** The auth limiter in `backend/cmd/api/main.go` keys solely on `c.IP()` with an in-memory bucket (10 req/min). It provides zero defense against distributed credential stuffing across rotating IP pools and does not protect specific administrative accounts from targeted brute force.
4. **Header Spoofing Vulnerability:** Fiber is configured with `ProxyHeader: fiber.HeaderXForwardedFor` without validating trusted upstream reverse proxies, allowing clients to manipulate their perceived IP.
5. **Stateless JWTs Without Revocation:** JWTs have a 24-hour lifetime without an active server-side session registry, token blacklist, or session rotation mechanism. A compromised token cannot be immediately revoked without deactivating the account.
6. **Frontend Session Route Misplacement:** Next.js session route handlers were authored in `frontend/src/lib/api/route.ts` rather than the Next.js App Router standard `frontend/src/app/api/auth/session/route.ts`, causing cookie session synchronization issues.

---

## 2. Inventory of Existing Code, Routes, Tables, and Models

### 2.1 Backend Code Components

| File Path | Role & Purpose | Key Functions / Structs |
|---|---|---|
| `backend/cmd/api/main.go` | Route registration & middleware wiring | Routes: `/api/admin/login`, `/api/admin/login/verify`, `/api/admin/*` |
| `backend/internal/handler/auth_handler.go` | HTTP authentication controller | `AdminLogin`, `AdminLoginVerifyOTP`, `AdminForgotPassword` |
| `backend/internal/handler/admin_handler.go` | Admin feature handlers | `GetMe`, `GetMyPermissions`, `ListAdmins`, `CreateAdmin`, `UpdateAdminRole` |
| `backend/internal/service/auth_service.go` | Authentication domain logic | `AdminLoginStep1`, `AdminLoginStep2`, `SendEmailOTP` |
| `backend/internal/service/admin_service.go` | Admin & permission business logic | `HasPermission`, `SetAdminPermissions`, `InsertAdminLog` |
| `backend/internal/repository/user_repo.go` | User and credential persistence | `GetAdminByEmail`, `GetByID`, `UpdatePasswordHash` |
| `backend/internal/repository/admin_repo.go` | Admin entity and permissions store | `ListAdmins`, `CreateAdmin`, `GetAdminPermissions`, `SetAdminPermissions` |
| `backend/internal/middleware/middleware.go` | Auth & authorization filters | `RequireAuth`, `RequireRole`, `RequireAdminLevel`, `RequirePermission` |
| `backend/internal/security/jwt.go` | Token issuing and validation | `GenerateJWTFull`, `ParseJWT`, `Claims` struct |
| `backend/internal/security/password.go` | Password cryptography | `HashPassword` (bcrypt cost 12), `CheckPassword` |
| `backend/internal/security/otp_crypto.go` | OTP generation & hashing | `GenerateCryptoOTP`, `HashOTP` (SHA-256) |
| `backend/internal/models/user.go` | User & permission definitions | `User`, `AllAdminPermissions`, `DefaultAdminPermissions` |

### 2.2 PostgreSQL Schema & Tables

| Table / Type | Migration | Purpose & Key Columns |
|---|---|---|
| `users` | `001_schema.sql`, `002_users.sql` | `id`, `name`, `phone`, `email`, `password_hash`, `role` (`VARCHAR(20)`), `admin_role` (`admin_role` ENUM: `'normal'`, `'super'`), `is_active` |
| `admin_permissions` | `001_schema.sql`, `008_admin_system.sql` | `(admin_id, permission)` PK, `enabled BOOLEAN`, `updated_at`, `granted_by`, `granted_at` |
| `admin_logs` | `008_admin_system.sql` | `id`, `admin_id`, `action`, `target_type`, `target_id`, `product_id`, `target_user_id`, `notes`, `ip_address`, `created_at` |
| `otp_codes` | `012_security_and_otp.sql`, `023_otp_codes_columns.sql` | `id`, `phone`, `email`, `code_hash`, `purpose`, `expires_at`, `verified`, `attempts`, `created_at` |
| `platform_settings` | `001_schema.sql`, `011_settings_and_limits.sql` | System key-value settings (`maintenance_mode`, `otp_gateway`, etc.) |

### 2.3 Frontend Components (Next.js 15)

| File Path | Role & Purpose |
|---|---|
| `frontend/src/app/admin/login/page.tsx` | Single legacy admin login form handling Step 1 (email/password) and Step 2 (email OTP). |
| `frontend/src/app/admin/layout.tsx` | Admin navigation shell, checks `isSuperAdmin` and `hasAdminPermission`. |
| `frontend/src/context/AuthContext.tsx` | Client authentication provider, decodes JWT, polls permissions every 30s. |
| `frontend/src/lib/api/admin.ts` | API client wrapper for admin endpoints. |
| `frontend/src/lib/api/route.ts` | Cookie handler for `hm_token` (misplaced; needs App Router path). |

---

## 3. End-to-End Authentication Flow Analysis

### Current Step 1: Password Submission
1. Admin enters email and password at `/admin/login`.
2. Browser sends `POST /api/admin/login` with `{ email, password }`.
3. Fiber applies `authLimiter` (10 req/min per IP).
4. `AdminLogin` handler validates email format and password minimum length.
5. `auth_service.AdminLoginStep1` queries database via `userRepo.GetAdminByEmail`.
6. Checks:
   - `u.Role == "admin"`
   - `u.IsActive == true`
   - `security.CheckPassword(u.PasswordHash, req.Password)` using bcrypt.
7. **Privilege Branching:**
   - **If `u.AdminRole == "super"`:** Immediately generates a 24h JWT via `security.GenerateJWTFull`, updates `last_login_at`, and returns `{ token, is_super: true, message: "login successful" }`. **MFA is completely bypassed.**
   - **If `u.AdminRole == "normal"`:** Generates a 6-digit OTP, stores SHA-256 hash in `otp_codes` (expires in 5 minutes), sends email via SMTP (or console in dev), and returns `{ is_super: false, message: "OTP sent..." }`.

### Current Step 2: Normal Admin OTP Verification
1. Admin enters 6-digit code.
2. Browser sends `POST /api/admin/login/verify` with `{ email, code }`.
3. Fiber applies `otpLimiter` (5 req/min per IP).
4. `auth_service.AdminLoginStep2` verifies code against `otp_codes` table.
5. If valid, generates a 24h JWT, updates `last_login_at`, and returns `{ token }`.

### Current Authorization Enforcement
- Every admin request includes `Authorization: Bearer <token>`.
- `middleware.RequireAuth` parses JWT HMAC-SHA256 signature, validates `exp`, and checks `isActiveChecker` (cached for 10s).
- `middleware.RequireRole("admin")` validates that claims indicate `role == "admin"`.
- `middleware.RequirePermission(perm)`:
  - If `adminRole == "super"`, access is granted unconditionally.
  - If `adminRole == "normal"`, executes `adminSvc.HasPermission(adminID, perm)` against PostgreSQL.

---

## 4. Detailed Vulnerability & Gap Findings

### [SEC-VULN-01] Critical: Super Administrator MFA Bypass
- **File:** `backend/internal/service/auth_service.go:504-511`
- **Finding:** The most privileged administrative role possesses the weakest authentication posture. Super administrators are excused from second-factor verification, allowing account takeover via single credential leaks.
- **Remediation:** Remove bypass. Require mandatory user-verified passkey authentication for all privileged roles, with Platform Owner requiring strict phishing-resistant credentials.

### [SEC-VULN-02] High: In-Memory, IP-Only Rate Limiting
- **File:** `backend/cmd/api/main.go:370-387`
- **Finding:** The rate limiters are fiber in-memory instances keyed only by `c.IP()`. In a multi-replica cluster or behind rotating residential proxies, attackers can execute distributed credential stuffing. Furthermore, legitimate administrators behind shared NATs can be blocked by noisy neighbors.
- **Remediation:** Introduce a unified, observable throttling engine keyed by `(IP, Account, Endpoint)` with progressive backoff, lockouts, and persistent state.

### [SEC-VULN-03] High: Reverse Proxy IP Trust Spoofing
- **File:** `backend/cmd/api/main.go:238`
- **Finding:** `ProxyHeader: fiber.HeaderXForwardedFor` trusts `X-Forwarded-For` blindly without verifying whether the request came from a trusted upstream gateway (e.g., Cloudflare, Render load balancer, AWS ALB). Attackers can send `X-Forwarded-For: <random-ip>` to bypass the IP rate limiter.
- **Remediation:** Validate upstream proxy IPs and normalize source extraction.

### [SEC-VULN-04] Medium: Stateless 24-Hour JWTs Without Server-Side Revocation
- **File:** `backend/internal/security/jwt.go:90`, `backend/internal/middleware/middleware.go:70-92`
- **Finding:** Issued JWTs are valid for 24 hours. If an admin session is hijacked, an administrator is demoted, or authenticators are revoked, existing tokens remain valid until expiration unless the user is completely deactivated in `users.is_active`.
- **Remediation:** Implement a database-backed session registry with unique session IDs (`sid`), explicit idle timeouts (shorter for Platform Owner), and immediate revocation.

### [SEC-VULN-05] Medium: Single Administrative Portal for All Staff
- **File:** `frontend/src/app/admin/login/page.tsx`
- **Finding:** A single login page serves both super administrators and subordinate managers. There is no separation between high-assurance executive governance and everyday operational tasks.
- **Remediation:** Create dedicated privileged portals:
  - `/executive-console/login` (Platform Owner)
  - `/platform-console/login` (Platform Administrator)
  - `/security-console/login` (Security Administrator)
  - `/staff/login` (Operations, Support, Moderation)

### [SEC-VULN-06] Low: Next.js Session Route Placement Defect
- **File:** `frontend/src/lib/api/route.ts`
- **Finding:** The `POST/GET/DELETE` route handler managing `hm_token` is stored under `src/lib/api/route.ts` instead of `src/app/api/auth/session/route.ts`. In Next.js 15, route handlers outside `app/` are not registered as endpoints.
- **Remediation:** Move route handler to `frontend/src/app/api/auth/session/route.ts`.

---

## 5. Reusable Existing Controls

The application already has solid building blocks that must be preserved:
1. **Password Hashing:** `internal/security/password.go` uses `golang.org/x/crypto/bcrypt` with Cost 12.
2. **Permission Table Structure:** `admin_permissions` with `(admin_id, permission)` already supports granular capability flags.
3. **Audit Log Foundation:** `admin_logs` table already tracks `admin_id`, `action`, `target_type`, `target_id`, `ip_address`, and timestamps.
4. **Active Account Check:** `middleware.WithIsActiveChecker` already validates `is_active` in PostgreSQL with caching.
5. **Localization Infrastructure:** `frontend/src/lib/i18n.ts` provides fallback translation resolution across English, Amharic, Afaan Oromo, and Tigrinya.

---

## 6. Migration & Architecture Roadmap

1. **Role Taxonomy Expansion:** Migrate `admin_role` ENUM from `('normal', 'super')` to the professional taxonomy:
   - `platform_owner`
   - `platform_admin`
   - `security_admin`
   - `operations_admin`
   - `moderator`
   - `support_agent`
2. **WebAuthn Integration:** Add `github.com/go-webauthn/webauthn` for passkey ceremony handling, supporting platform authenticators (TouchID, FaceID, Windows Hello, Android) and FIDO2 hardware keys.
3. **Session & Security Tables:**
   - `admin_passkey_credentials`: Public keys, credential IDs, counters, device names.
   - `admin_sessions`: Server-tracked sessions with `session_id`, `user_id`, `portal`, `created_at`, `expires_at`, `revoked_at`.
   - `admin_auth_attempts`: Centralized brute-force tracking by normalized IP and account.
4. **Dedicated Frontend Portals:** Build distinct, accessible, localized login portals matching the HUNTIQ enterprise design system.

---

## 6. Post-Implementation Adversarial Audit Findings

Following initial implementation, an adversarial review identified the following findings documented in detail in [`admin-post-implementation-security-review.md`](file:///home/poven/GO%20PROJECT/sellbuy/docs/audits/admin-post-implementation-security-review.md):
1. **Passkey Endpoint Route Mismatch (High):** Frontend called `/admin/auth/passkey/*` while backend exposed `/admin/passkeys/*`, causing HTTP 404 on browser ceremonies.
2. **Stateless JWT Session Bypass (High):** Presentation of 24h JWTs in `Authorization: Bearer` bypassed server-side revocation in `admin_sessions`.
3. **Scalar Role Rank Privilege Inversion (Medium-High):** `Platform Administrator` shared rank 8 with `Security Administrator`, enabling cross-domain access to security logs.
4. **Rate-Limit DoS Vector (Medium):** `(identifier = $1 OR ip_address = $2)` in failed attempt query allowed deliberate remote account lockout.
5. **Re-Authentication Gap (Medium):** Authenticator registration and revocation lacked `RequireReauth` middleware.
6. **Exposed CLI Recovery Credential (Cleared):** Verified expired at `03:00:39` and permanently rejected by database challenge queries.

