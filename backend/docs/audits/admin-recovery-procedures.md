# HUNTIQ Enterprise Administrative Recovery & Credential Management Protocol

**Document:** `docs/audits/admin-recovery-procedures.md`  
**Date:** October 2026  
**Status:** Operational Protocol & Security Standard  
**Repository:** `/home/poven/GO PROJECT/sellbuy`  

---

## 1. Principles of Privileged Recovery

Recovery for administrative accounts must never undermine primary authentication security:
1. **No Silent Downgrades:** A failed passkey authentication must **never** automatically fall back to password-only access.
2. **No Email-Only Privileged Reset:** Email verification alone cannot reset credentials for `platform_owner`, `platform_admin`, or `security_admin`.
3. **No Hidden Backdoors:** No secret maintenance accounts, developer bypass flags, or undocumented master passwords exist in the system.
4. **Mandatory Audit Logging:** Every recovery request, authorization step, and credential rotation is permanently recorded in `admin_logs`.

---

## 2. Emergency Platform Owner Recovery Protocol

If the sole registered Platform Owner loses all enrolled authenticators:

### 2.1 Authorized Recovery Custodians
Recovery requires direct cryptographic or terminal authorization by designated platform custodians:
- Platform Host Infrastructure Custodian (access to server host / secure terminal)
- Primary Database Administrator

### 2.2 CLI-Based Recovery Tool
A secure, offline-first command-line utility is provided:
`backend/cmd/admin-recovery`

Execution Requirements:
```bash
# Must be executed directly on the host or private bastion terminal with DB credentials:
cd /home/poven/GO\ PROJECT/sellbuy/backend
go run ./cmd/admin-recovery --action=generate-bootstrap-token --user-id=36
```
1. Generates a cryptographically random, high-entropy 256-bit emergency bootstrap token.
2. The token is hashed with SHA-256 and stored in `admin_auth_challenges` with a strict 15-minute lifespan.
3. The custodian uses the single-use token on the local executive console (`/executive-console/login?bootstrap=<token>`) to immediately enroll a replacement hardware passkey.
4. All previous sessions, passkeys, and tokens are revoked simultaneously.
5. An immutable security audit log entry is written with action `"platform_owner_emergency_recovery"`.

---

## 3. Staff Credential Recovery Protocol

For operational staff (`operations_admin`, `moderator`, `support_agent`):
1. **Request Submission:** Staff member requests credential reset through `/staff/login?recovery=true`.
2. **Administrative Approval:** A `platform_admin` or `platform_owner` reviews the request in the admin console.
3. **Identity Verification:** Secondary verification (out-of-band phone/voice verification with authorized manager).
4. **One-Time Enrollment Link:** Administrator issues a time-bounded (1 hour) enrollment link.
5. **Mandatory Password & Passkey Enrollment:** Upon opening the link, staff member sets a new strong password and registers their authenticator.

---

## 4. Post-Implementation Security Audit: Credential Lifecycle Verification

During the adversarial post-implementation audit:
1. **Plaintext Protection Verified:** The offline CLI tool (`cmd/admin-recovery`) was verified to hash recovery tokens with SHA-256 (`tokenHash := security.HashToken(token)`) before database insertion. Plaintext secrets are never persisted or logged.
2. **Deterministic Expiration Verified:** The test recovery challenge issued at `02:45:39` was deterministically observed to transition to `is_expired = true` upon reaching its 15-minute TTL (`03:00:39`).
3. **Zero Redemption Surface:** The audit verified that no web API endpoint consumes `ceremony_type = 'reauth'`, preventing any network-based redemption of CLI recovery tokens.

