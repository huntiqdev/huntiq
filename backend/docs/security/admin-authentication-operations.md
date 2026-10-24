# HUNTIQ Enterprise Administrative Authentication Operations Manual

**Document:** `docs/security/admin-authentication-operations.md`  
**Date:** October 2026  
**Status:** Operational Manual  
**Classification:** Enterprise Internal  

---

## 1. Daily Operations & Portal Directory

Administrative operations are partitioned across four specialized ingress portals:

| Portal | Route | Intended Users | Primary Daily Actions |
|---|---|---|---|
| **Executive Console** | `/executive-console/login` | Platform Owner | Executive governance, role provisioning, security policy approval |
| **Platform Console** | `/platform-console/login` | Platform Administrators | Category & attribute management, location hierarchies, platform settings |
| **Security Console** | `/security-console/login` | Security Administrators | Audit log review, incident investigation, session termination, passkey revocations |
| **Staff Portal** | `/staff/login` | Operations, Moderators, Support | Product review, seller verification, dispute mediation, ticket handling |

---

## 2. Authenticator Management Procedures

### 2.1 Enrolling a Backup Passkey
To prevent lockout, every privileged administrator is strongly encouraged to enroll at least two authenticators (e.g., Primary Laptop Touch ID / Windows Hello + Backup YubiKey).

1. Log in to the authorized console.
2. Navigate to **Account Security → Passkeys & Authenticators**.
3. Re-authenticate via current passkey (reauth freshness requirement: < 10 minutes).
4. Click **Register New Authenticator**.
5. Insert security key or verify secondary device.
6. Provide a descriptive label (e.g., `"Work MacBook TouchID"`, `"Backup YubiKey 5C"`).

### 2.2 Revoking a Lost or Compromised Passkey
1. Navigate to **Account Security → Passkeys & Authenticators**.
2. Locate the compromised credential.
3. Click **Revoke Passkey**.
4. Confirm re-authentication.
5. All sessions associated with that credential are automatically terminated.
6. If the administrator is attempting to revoke their last remaining authenticator, the system requires registering a replacement before the revocation can proceed.
