# HUNTIQ Enterprise Administrative Role Taxonomy & Authorization Matrix

**Document:** `docs/audits/admin-authorization-matrix.md`  
**Date:** October 2026  
**Status:** Approved Specification & Architecture Design  
**Repository:** `/home/poven/GO PROJECT/sellbuy`  

---

## 1. Official Administrative Role Taxonomy

HUNTIQ replaces legacy labels (`super`, `normal`, `sub-admin`) with a professional, enterprise-grade taxonomy. These roles define internal application privileges under the principle of least privilege and strict separation of duties.

| Internal Role Identifier | Professional Display Name | Primary Responsibility | Dedicated Portal |
|---|---|---|---|
| `platform_owner` | Platform Owner | Ultimate platform governance, high-risk operational authorization, executive role delegation | `/executive-console/login` |
| `platform_admin` | Platform Administrator | System configuration, infrastructure maintenance, platform catalog and policy updates | `/platform-console/login` |
| `security_admin` | Security Administrator | Security event auditing, incident investigation, access reviews, session termination | `/security-console/login` |
| `operations_admin` | Operations Administrator | Routine operational workflows, batch job health, seller verification, tier adjustments | `/staff/login` |
| `moderator` | Content Moderator | Product listing inspection, abuse report triage, content moderation | `/staff/login` |
| `support_agent` | Customer Support Agent | User account assistance, dispute mediation, ticket handling within scoped permissions | `/staff/login` |

### 1.1 Separation of Corporate Job Titles from Application Privileges

Corporate designations (e.g., Founder, CEO, CTO, VP Engineering, CISO, Head of Trust & Safety) are organizational labels. They **never** automatically confer administrative database permissions or superuser status.
- Corporate titles are stored in organizational profile attributes, distinct from `admin_role`.
- A CTO cannot self-promote to `platform_owner`.
- A CISO or Security Administrator cannot modify platform fee structures or product catalogs.
- A Content Moderator cannot escalate privileges by modifying request payloads, frontend state, or JWT claims.
- The Platform Owner cannot bypass security audit logs.

---

## 2. Legacy-to-Professional Role Migration Mapping

| Legacy `admin_role` | Legacy Permissions State | New `admin_role` | Migration Rationale |
|---|---|---|---|
| `super` (Sole existing account: ID 36) | Unchecked superuser access across all endpoints | `platform_owner` | Preserves governance authority while binding the account to WebAuthn passkey requirements and immutable audit trails. |
| `normal` (Sub-admins with `admin_permissions`) | Configured via `admin_permissions` rows | `operations_admin`, `moderator`, or `support_agent` | Specific assignment based on assigned permission sets, preserving existing functional capabilities without privilege expansion. |

---

## 3. Comprehensive Granular Permission Matrix

The application utilizes 24 granular permission keys. Under the zero-trust model:
- `platform_owner` possesses platform governance and administrative authority, subject to mandatory passkeys and audit logging.
- `security_admin` has exclusive security investigation and audit visibility, but **no** authority to modify marketplace inventory or financial rules.
- Operational roles receive only the permissions explicitly granted to their duties.

| Permission Key | Domain Description | `platform_owner` | `platform_admin` | `security_admin` | `operations_admin` | `moderator` | `support_agent` |
|---|---|:---:|:---:|:---:|:---:|:---:|:---:|
| **Governance & Access** | | | | | | | |
| `governance:manage` | Grant/revoke administrative roles | **YES** | NO | NO | NO | NO | NO |
| `governance:policy` | Modify core authentication/security policy | **YES** (with reauth) | NO | Audited Review | NO | NO | NO |
| `admins:read` | View administrative staff directory | **YES** | **YES** | **YES** | NO | NO | NO |
| `admins:manage` | Provision/update staff accounts | **YES** | Scoped Staff | NO | NO | NO | NO |
| **Security & Auditing** | | | | | | | |
| `logs:view` | View immutable system and admin audit logs | **YES** | NO | **YES** | NO | NO | NO |
| `security:investigate`| View security incidents, rate-limit events | **YES** | NO | **YES** | NO | NO | NO |
| `security:revoke` | Terminate compromised sessions/keys | **YES** | NO | **YES** | NO | NO | NO |
| `otp_logs:view` | Inspect OTP authentication logs | **YES** | NO | **YES** | NO | NO | NO |
| **Platform & Settings** | | | | | | | |
| `settings:read` | Read system platform configuration | **YES** | **YES** | **YES** | Optional | NO | NO |
| `settings:write` | Modify platform settings & limits | **YES** | **YES** | NO | NO | NO | NO |
| `operations:monitor` | View background jobs & system health | **YES** | **YES** | **YES** | **YES** | NO | NO |
| **Catalog & Moderation** | | | | | | | |
| `pending:review` | Review & approve/reject listings | **YES** | **YES** | NO | **YES** | **YES** | NO |
| `products:manage` | Edit, flag, or delete active products | **YES** | **YES** | NO | **YES** | **YES** | NO |
| `categories:manage` | Manage category taxonomy & hierarchy | **YES** | **YES** | NO | Optional | NO | NO |
| `locations:manage` | Manage location taxonomy & geocoding | **YES** | **YES** | NO | Optional | NO | NO |
| `banners:manage` | Manage promotional banners | **YES** | **YES** | NO | Optional | NO | NO |
| **Trust, Safety & Users**| | | | | | | |
| `sellers:manage` | Activate/deactivate seller accounts | **YES** | **YES** | NO | **YES** | NO | Scoped |
| `seller_verification` | Review business/ID verification docs | **YES** | **YES** | NO | **YES** | NO | Scoped |
| `trust_safety` | Investigate fraud & abuse reports | **YES** | **YES** | **YES** | **YES** | **YES** | NO |
| `appeals:review` | Process user account suspension appeals | **YES** | **YES** | **YES** | **YES** | NO | Scoped |
| `refunds:process` | Review and adjudicate refund disputes | **YES** | **YES** | NO | **YES** | NO | Scoped |
| `tier_upgrades:manage`| Approve seller/buyer tier upgrades | **YES** | **YES** | NO | **YES** | NO | NO |
| **Finance & Revenue** | | | | | | | |
| `profit:view` | Inspect platform revenue & fee analytics | **YES** | **YES** | Audited | NO | NO | NO |
| `withdrawals:manage` | Approve and process referral payouts | **YES** | **YES** | Audited | **YES** | NO | NO |

---

## 4. Backend Authorization Architecture

### 4.1 Request Evaluation Pipeline
Every incoming administrative request to `api.Group("/admin", ...)` is evaluated through four discrete gates:

```
[ Incoming Request ]
         │
         ▼
[ 1. Transport & Session Gate ]
   - Validate TLS connection
   - Extract Session Token (HttpOnly Cookie or Bearer)
   - Verify Session ID in `admin_sessions` table (not revoked, within idle & max expiry)
         │
         ▼
[ 2. Account Status Gate ]
   - Query user identity: `users.is_active == TRUE` and `role == 'admin'`
   - Reject suspended or terminated principals immediately (fail-closed)
         │
         ▼
[ 3. Portal & Role Alignment Gate ]
   - Verify that the authenticated `admin_role` matches the requested endpoint / portal scope:
     * Executive endpoints require `platform_owner`
     * Platform console endpoints require `platform_owner` or `platform_admin`
     * Security console endpoints require `platform_owner` or `security_admin`
     * Staff endpoints require valid operational role with appropriate permission
         │
         ▼
[ 4. Granular Permission Gate ]
   - Execute `adminSvc.HasPermission(ctx, adminID, requiredPermission)`
   - Deny by default if permission record is absent or `enabled == FALSE`
         │
         ▼
[ 5. Sensitive Operation Reauthentication Gate ]
   - For high-consequence operations (role grants, credential revocation, payout approval),
     verify freshness of authentication (< 10 minutes since WebAuthn verification)
         │
         ▼
[ Authorized Handler Execution + Immutable Audit Log Entry ]
```

### 4.2 Sensitive Operations Requiring Reauthentication

1. Modifying administrator roles or granting `platform_owner` privileges.
2. Registering or revoking Platform Owner passkeys.
3. Modifying platform security configuration or authentication policies.
4. Exporting sensitive user data or full audit logs.
5. Performing irreversible system repairs or purging data.
6. Processing large financial withdrawals (> ETB 50,000).

---

## 5. Privilege Escalation Defenses

1. **No Client Role Injection:** The backend never trusts roles from request bodies, URL params, or untrusted claims. Role is loaded exclusively from verified session state.
2. **Strict Hierarchy Ceiling:** No administrator can grant a role higher than their own. Only `platform_owner` can grant or revoke `platform_admin` or `security_admin` roles.
3. **Audit Immunity:** `security_admin` cannot delete or alter audit logs (`admin_logs`). PostgreSQL rules and service constraints prevent log erasure.
4. **Self-Promotion Block:** An administrator cannot modify their own permissions or upgrade their own role.

---

## 6. Remediation: Elimination of Scalar Ranks & Explicit Domain Authorization (Task 3)

### 6.1 Root Cause
In legacy implementations, `adminRoleRank` assigned scalar integer weights to roles:
- `platform_owner` / `super`: 10
- `platform_admin`: 8
- `security_admin`: 8
- `operations_admin`: 4

Because `RequireAdminLevel(requiredRole)` checked `userRank >= requiredRank`, Platform Admins (`rank 8`) inherited access to Security Admin routes (`/admin/logs`, `/admin/otp-logs`, `/admin/activity-logs`), directly violating separation of duties. Concurrently, Platform Admins were locked out of platform settings (`/admin/settings/with-reason`) and categories which previously required `"super"` (rank 10).

### 6.2 Architectural Remediation
1. **Elimination of Privilege Inheritance:**
   - Deprecated `adminRoleRank` scalar comparisons.
   - Built explicit domain middlewares:
     * `RequirePlatformOwner`: strictly restricted to `platform_owner` (and backward-compatible `super`).
     * `RequireSecurityAdmin`: restricted to `security_admin` and `platform_owner`. Platform Admins are strictly forbidden.
     * `RequirePlatformAdmin`: restricted to `platform_admin` and `platform_owner`. Security Admins are strictly forbidden.
     * `RequireExplicitAdminRole(roles ...string)`: explicit multi-role gate.
     * `RequirePermission(permission)`: upgraded to evaluate `models.RoleDomainPermissions[adminRole][permission]` baseline maps before falling back to database `admin_permissions`.
2. **Explicit Domain Baseline Permissions (`models.RoleDomainPermissions`):**
   - `platform_admin`: categories, hierarchy, attributes, locations, banners, settings, health, pending review, user/seller inspection. Strictly excluded from security logs, OTP logs, and owner provisioning.
   - `security_admin`: audit logs, OTP logs, activity logs, trust & safety, session review. Strictly excluded from product catalog editing and governance.
   - `operations_admin`: pending reviews, sellers, seller verification, tier upgrades, refunds, appeals, operations health.
   - `moderator`: pending reviews, active listings, product flags, trust/safety reports.
   - `support_agent`: dispute mediation, appeals triage, buyer/seller assistance.
3. **Privilege Escalation Prevention:**
   - In `admin_service.go`, `CreateAdmin` explicitly prevents non-owners from provisioning administrative accounts, disallows creating `platform_owner` accounts via API, and automatically provisions new administrators with their domain's baseline permissions.
   - In `admin_handler.go` and `auth_handler.go`, all governance endpoints require `models.RolePlatformOwner` or `models.LegacyRoleSuper`.
4. **Live Role Resolution & Session Invalidation:**
   - Authorization is derived from the live database role in `users` and active server session in `admin_sessions`, guaranteeing immediate revocation upon demotion or session termination.

---

## 7. Verification and Test Evidence

### 7.1 Integration Test Suite: `TestExplicitAdminAuthorizationMatrix`
Executed via `go test -v -run TestExplicitAdminAuthorizationMatrix ./internal/middleware`:
- `Platform_Administrator_Boundary_Tests`: PASS
  * Accessing `/settings/with-reason`, `/settings/log`, `/categories`, `/system/health`: 200 OK.
  * Accessing Security Admin routes (`/logs`, `/otp-logs`, `/activity-logs`): **403 Forbidden**.
  * Accessing Platform Owner routes (`/admins`, `/admins/:id/role`, `/users/:id/ban`): **403 Forbidden**.
- `Security_Administrator_Boundary_Tests`: PASS
  * Accessing `/logs`, `/otp-logs`, `/activity-logs`: 200 OK.
  * Accessing Platform Admin routes (`/settings/with-reason`, `/settings/log`, `/categories`, `/attributes`): **403 Forbidden**.
  * Accessing Platform Owner routes (`/admins`, `/admins/:id/role`, `/users/:id/ban`): **403 Forbidden**.
- `Operations_Administrator_Boundary_Tests`: PASS
  * Accessing operational routes (`/operations/health`, `/verification/documents`, `/products/pending`, `/sellers`, `/upgrade-requests`): 200 OK.
  * Accessing Security logs (`/logs`, `/otp-logs`): **403 Forbidden**.
  * Accessing Platform governance (`/admins`, `/settings/with-reason`): **403 Forbidden**.
- `Moderator_&_Support_Agent_Boundary_Tests`: PASS
  * Accessing unauthorized admin functions: **403 Forbidden**.
  * Accessing scoped allowed domains: 200 OK.
- `Platform_Owner_Governance_Tests`: PASS
  * Accessing all governance, security, and platform routes: 200 OK.
  * Attempting request with revoked/missing session: **401 Unauthorized**.
  * Demoted/disabled administrator: **401/403**.
- `Unknown_Role_Rejection`: PASS
  * Custom/unrecognized roles safely rejected on all administrative endpoints: **403 Forbidden**.

