# HUNTIQ Enterprise Security Auditing, Monitoring & Alerting

**Document:** `docs/audits/admin-security-monitoring.md`  
**Date:** October 2026  
**Status:** Monitoring Specification & Security Operations Guide  
**Repository:** `/home/poven/GO PROJECT/sellbuy`  

---

## 1. Security Event Taxonomy

All security-relevant actions within the administrative perimeter produce structured audit entries in `admin_logs` and generate real-time monitoring events.

| Event Code | Event Category | Trigger Description | Severity | Logged Attributes |
|---|---|---|---|---|
| `AUTH_PASSKEY_SUCCESS` | Authentication | Successful passkey assertion verification | Informational | `admin_id`, `portal`, `ip_address`, `credential_id` |
| `AUTH_PASSKEY_FAILURE` | Authentication | Failed passkey signature, expired challenge, bad RP | Warning | `identifier`, `portal`, `ip_address`, `failure_reason` |
| `AUTH_RATE_LIMIT_TRIP` | Abuse Detection | Rate limit threshold exceeded (5+ failed attempts) | Warning | `identifier`, `ip_address`, `portal`, `attempt_count` |
| `AUTH_CLONE_WARNING` | Security Alert | Passkey signature counter rollback detected | **CRITICAL** | `admin_id`, `credential_id`, `sign_count`, `stored_count` |
| `SESSION_REVOKED` | Session Security | Session explicitly revoked by security admin or user | Notice | `admin_id`, `session_id`, `revoked_by`, `reason` |
| `ROLE_ASSIGNED` | Authorization | Privileged role granted or updated | High | `actor_admin_id`, `target_user_id`, `old_role`, `new_role` |
| `PASSKEY_ENROLLED` | Authenticator | New passkey registered for admin account | High | `admin_id`, `credential_id`, `authenticator_name` |
| `PASSKEY_REVOKED` | Authenticator | Passkey revoked or deregistered | High | `admin_id`, `credential_id`, `reason` |
| `CONFIG_MUTATION` | Platform Governance | System setting or policy changed | High | `admin_id`, `key`, `old_value`, `new_value` |
| `EMERGENCY_RECOVERY` | Disaster Recovery | Platform Owner emergency bootstrap invoked | **CRITICAL** | `target_user_id`, `token_id`, `ip_address` |

---

## 2. Redaction & Privacy Mandate

To prevent credential leaks and respect user privacy, the following attributes are **strictly prohibited** from ever being written to logs, consoles, or external metrics:
1. Plaintext passwords or password fragments.
2. Unhashed MFA/OTP codes.
3. Raw biometric data or biometric vectors.
4. Passkey private keys or seed phrases.
5. Plaintext session tokens or bootstrap secrets.
6. Full credit card numbers or sensitive payment keys.

---

## 3. High-Priority Alerting Rules

Automated alerts are raised to designated security channels upon detecting:
1. **Passkey Clone / Replay Detection (`AUTH_CLONE_WARNING`):** Immediate alert. Indicates private key extraction or device cloning.
2. **Repeated Privileged Login Failures:** More than 5 consecutive authentication failures on `/executive-console/login` within 10 minutes.
3. **Unexpected Role Promotion:** Any assignment of `platform_owner`, `platform_admin`, or `security_admin`.
4. **Emergency Bootstrap Invocation:** Any generation or usage of emergency bootstrap tokens.
5. **Off-Hours Sensitive Configuration Modification:** System settings or payment gateway modifications during non-business hours.
