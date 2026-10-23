# HUNTIQ Enterprise Administrative Threat Model & Attack Surface Analysis

**Document:** `docs/audits/admin-login-threat-model.md`  
**Date:** October 2026  
**Methodology:** STRIDE (Spoofing, Tampering, Repudiation, Information Disclosure, Denial of Service, Elevation of Privilege) + NIST SP 800-63B / OWASP ASVS 4.0  
**Repository:** `/home/poven/GO PROJECT/sellbuy`  

---

## 1. Attack Surface Overview

The HUNTIQ administrative perimeter controls access to marketplace listings, seller financial accounts, transaction processing, platform configurations, and sensitive customer data across Ethiopia.

Administrative ingress points include:
1. HTTP Auth Endpoints: `/api/admin/login`, `/api/admin/login/verify`, `/api/admin/forgot-password`
2. Dedicated Frontend Portals: `/executive-console/login`, `/platform-console/login`, `/security-console/login`, `/staff/login`
3. Protected Administrative API Group: `/api/admin/*`
4. Session Management: Cookie `hm_token` and HTTP `Authorization: Bearer <token>`
5. WebAuthn Registration and Assertion Endpoints: `/api/admin/passkeys/*`

---

## 2. Comprehensive STRIDE Threat Analysis & Mitigations

| Threat ID | Threat Category | Threat Scenario | Existing Codebase State | Inherent Risk | Required Enterprise Mitigation | Verification Test |
|---|---|---|---|---|---|---|
| **THREAT-01** | **Spoofing** | Adversary sets up a look-alike phishing domain (`huntiq-admin.com`) to intercept credentials and relay OTPs in real-time. | Legacy email + OTP flow is fully vulnerable to reverse-proxy phishing (Evilginx). | **CRITICAL** | Implement FIDO2/WebAuthn passkeys. WebAuthn cryptographic assertions bind to the origin (`rpID`). Phishing proxies cannot obtain valid cryptographic signatures. | Test assertion failure against spoofed origin / invalid RP ID. |
| **THREAT-02** | **Elevation of Privilege** | Adversary compromises single password of `super` admin account and gains immediate total platform takeover. | `auth_service.go:504-511` explicitly bypasses OTP for `super` admin, returning full JWT immediately. | **CRITICAL** | Eliminate super admin MFA bypass. Require mandatory biometric/hardware passkey with user verification (`UV=required`) for Platform Owner. | Attempt password-only login for Platform Owner; verify strict denial without passkey assertion. |
| **THREAT-03** | **Denial of Service** | Attacker distributes brute-force requests across thousands of residential IPs to bypass IP-based rate limiting or exhaust resources. | In-memory `authLimiter` in `main.go:380` keys only by `c.IP()` (10 req/min). | **HIGH** | Implement centralized, multi-signal progressive throttling keyed by normalized `(Account, IP, Endpoint)`. Progressive delays (1s → 5s → 30s → temporary 5m restriction). Avoid permanent lockout to prevent attacker-induced DoS on real admins. | Simulate 50 requests across 10 IPs targeting one account; verify throttling activates globally. |
| **THREAT-04** | **Spoofing / Bypass** | Attacker sends arbitrary `X-Forwarded-For` header to rotate perceived IP and bypass IP-based rate limits. | `main.go:238` sets `ProxyHeader: fiber.HeaderXForwardedFor` without validating trusted upstream proxies. | **HIGH** | Only extract client IP from trusted reverse proxy CIDRs. Fallback to direct socket address if upstream is untrusted. | Send requests with forged `X-Forwarded-For: 8.8.8.8` and verify real socket IP is evaluated. |
| **THREAT-05** | **Information Disclosure** | Attacker probes `/api/admin/login` to enumerate valid administrator emails and identify who holds Platform Owner status. | Different error messages or response structures reveal if an account exists or has `super` status. | **MEDIUM** | Standardize public authentication responses to generic messages ("Invalid credentials or unauthorized portal"). Response times normalized via constant-time dummy checks. | Automated fuzzing comparing status codes, JSON keys, and timing for existing vs non-existent accounts. |
| **THREAT-06** | **Tampering / Hijacking** | Stolen JWT or cookie replay allows attacker persistent access for up to 24 hours even after admin termination. | Stateless JWTs with 24h lifetime (`jwt.go:90`). No token revocation list or database session tracking. | **HIGH** | Implement server-side `admin_sessions` table with cryptographic session tokens, idle timeout (15 min for Platform Owner, 30 min for Staff), and instant revocation on role change or logout. | Revoke session in backend and verify immediate 401 on next API request within 1 second. |
| **THREAT-07** | **Elevation of Privilege** | Sub-admin edits frontend form or crafts API call to grant themselves `super` role or unrestricted permissions. | `admin_handler.go:855` enforces `normal`, but API parameter tampering could occur if handler logic lacks strict check. | **HIGH** | Strict role hierarchy validation in Go service layer: an administrator cannot grant permissions higher than their assigned scope. Only `platform_owner` can grant administrative roles. | Submit PUT request as `operations_admin` attempting to grant `platform_owner`; verify 403 Forbidden. |
| **THREAT-08** | **Repudiation / Tampering** | Rogue administrator performs malicious database alteration and erases audit logs to cover tracks. | `admin_logs` table has foreign keys, but no cryptographic hash chaining or write-only constraints. | **HIGH** | Store `admin_logs` as append-only. No admin role (even `security_admin` or `platform_owner`) has API endpoints to delete or alter audit logs. | Attempt DELETE/PUT on audit log endpoints; verify strict denial. |
| **THREAT-09** | **Tampering / Spoofing** | Attacker hijacks WebAuthn enrollment ceremony by intercepting challenge or replaying previously used credentials. | WebAuthn not yet present. | **HIGH** | Server-side challenge generation with cryptographically secure random bytes (`crypto/rand`), single-use binding, 60-second expiration, origin and RP ID validation, and counter verification. | Replay passkey assertion; verify server rejects challenge reuse. |
| **THREAT-10** | **Elevation of Privilege** | Unprivileged staff access privileged portals or execute cross-portal administrative requests. | Existing frontend route guards check `isSuperAdmin` locally, but direct API access must be independently gated. | **HIGH** | Portal-to-role binding enforced at backend middleware: token and session record must match authorized portal level. Fail-closed on missing security context. | Attempt using staff session token against `/api/admin/system/*`; verify 403 Forbidden. |

---

## 3. Defense-in-Depth Architecture

To mitigate all identified threats, the system architecture establishes four defensive layers:

```
[ LAYER 1: Perimeter & Ingress ]
  ├── Trusted Reverse Proxy Validation (Anti-IP Spoofing)
  ├── Progressive Account & IP Rate Limiting (Anti-Stuffing)
  ├── Uniform Error Envelope (Anti-Enumeration)
  └── Dedicated Portal Routes (Isolation of Concerns)

[ LAYER 2: Strong Authentication Ceremony ]
  ├── WebAuthn / FIDO2 Passkeys (Phishing-Resistant Cryptography)
  ├── Mandatory User Verification (Biometric / PIN)
  ├── Strict Origin & RP ID Binding
  └── Time-Bounded, Single-Use Server Challenges

[ LAYER 3: Session & State Protection ]
  ├── Server-Side Session Tracking (`admin_sessions`)
  ├── Secure, HttpOnly, SameSite Cookies + Bearer Options
  ├── Configurable Idle & Absolute Expirations
  └── Instant Session Revocation on Security Events

[ LAYER 4: Zero-Trust Authorization & Audit ]
  ├── Least-Privilege Granular Permissions
  ├── Independent Backend Evaluation on Every Request
  ├── Sensitive Action Reauthentication Gate (< 10 min freshness)
  └── Immutable, Tamper-Evident Security Audit Logs
```

---

## 4. Threat Validation & Acceptance Testing Matrix

The implementation will be verified against the following test scenarios:
1. **Passkey Origin Mismatch Test:** Validate assertion rejection when origin != `APP_BASE_URL` / `FRONTEND_URL`.
2. **Replay Attack Test:** Validate that an assertion cannot be reused after initial challenge consumption.
3. **MFA Bypass Elimination Test:** Verify that Platform Owner and Platform Administrator cannot authenticate via password alone without a valid passkey assertion.
4. **Credential Stuffing Simulation:** Validate that 5 failed attempts within 2 minutes triggers progressive delay and account throttling without permanent lockout.
5. **Cross-Role Tampering Test:** Verify that `operations_admin` token cannot invoke `/api/admin/admins` or `/api/admin/settings`.
6. **Session Revocation Test:** Validate that calling session revocation immediately terminates access across all concurrent requests.
