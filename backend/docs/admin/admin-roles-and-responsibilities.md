# HUNTIQ Admin Roles & Responsibilities Specification

## 1. Executive Summary & Architecture Overview

HUNTIQ operates a defense-in-depth, zero-trust administrative governance model combining role-based access control (RBAC), fine-grained capability checks, and WebAuthn (FIDO2) phishing-resistant hardware passkeys. The system strictly segregates operational tiers to enforce the principle of least privilege, preventing single-point credential compromise from compromising the marketplace ecosystem, banking integrations, or tenant records.

### Tier Architecture

| Governance Tier | System Label | Primary Authenticator | Session Lifetime | Target Consoles |
|---|---|---|---|---|
| **Tier-0** | Root Governance / Platform Owner | WebAuthn Hardware Passkey (User Verification Enforced) | 1 hour idle / 4 hours absolute | `/executive-console/login` |
| **Tier-1** | Operations & Platform Configuration | WebAuthn Hardware Passkey (PIN / Biometrics) | 2 hours idle / 8 hours absolute | `/platform-console/login` |
| **Tier-2** | Security & Compliance Oversight | WebAuthn Hardware Passkey (FIDO2 Token) | 2 hours idle / 8 hours absolute | `/security-console/login` |
| **Tier-3** | Operational Staff & Support | WebAuthn Passkey / Time-Limited Password + MFA | 4 hours idle / 12 hours absolute | `/staff/login` |

---

## 2. Verified Role Definitions & Technical Identifiers

### Investigation of Roles

Database migration `029_enterprise_admin_roles_and_security.sql`, Go models `models/passkey.go`, and backend middleware `backend/internal/middleware/middleware.go` define the actual schema and runtime enum values.

```sql
-- Migration 029 Role Enum Definition
CREATE TYPE admin_role_type AS ENUM (
    'platform_owner',   -- Legacy equivalent: 'super'
    'platform_admin',
    'security_admin',
    'operations_admin', -- Legacy equivalent: 'normal'
    'moderator',
    'support_agent'
);
```

#### Detailed Role Specifications

### 1. Platform Owner
- **Code & Database Identifier:** `platform_owner` (legacy alias: `super`).
- **Database Storage:** Stored in `users.admin_role = 'platform_owner'` (or `super` with `role = 'admin'`).
- **Purpose:** Ultimate sovereign governance of the HUNTIQ platform. Possesses unconstrained authority over root configurations, security policies, administrative provisioning, database maintenance flags, and banking gateways.
- **Allowed Pages & Routes:** Full access to all 25 administrative routes, including `/admin/users` (Admin Accounts), `/admin/settings`, `/admin/finance`, and `/admin/logs`.
- **Actual Backend Permissions:** Inherits all permissions. Specifically satisfies `RequirePlatformOwner` and `RequireSuperAdmin` middleware. Bypasses domain permission checks.
- **Capabilities:**
  - *Read:* All data, metrics, PII, financial ledgers, audit logs, and system states.
  - *Write/Manage:* Provision/deprovision admins, grant/revoke RBAC permissions, modify tax models, alter commission fees, configure fee schedules, execute manual refunds, and initiate platform maintenance locks.
- **Restricted Actions:** Self-deletion is prevented at the database and handler levels. Reauthentication with hardware passkey is mandatory for critical destructive operations.
- **Authentication Requirements:** WebAuthn Passkey with `userVerification: required` at `/executive-console/login`. Zero password fallback permitted.

---

### 2. Platform Administrator
- **Code & Database Identifier:** `platform_admin`.
- **Database Storage:** `users.admin_role = 'platform_admin'`.
- **Purpose:** Core platform operations, catalog structure governance, infrastructure configuration, category taxonomies, financial models, and operational monitoring.
- **Allowed Pages & Routes:** `/admin/dashboard`, `/admin/operations`, `/admin/products`, `/admin/categories`, `/admin/attributes`, `/admin/locations`, `/admin/search`, `/admin/finance`, `/admin/profit`, `/admin/settings`, `/admin/verification`, `/admin/trust-safety`, `/admin/refunds`, `/admin/appeals`, `/admin/referrals`, `/admin/upgrades`, `/admin/premium-users`.
- **Actual Backend Permissions:** Satisfies `RequirePlatformAdmin` and standard permission checks: `operations`, `categories`, `attributes`, `locations`, `products`, `finance`, `profit`, `settings`, `seller_verification`, `refunds`, `appeals`, `referrals`.
- **Capabilities:**
  - *Read:* Catalog, seller records, transactions, FP&A projections, profit analytics, system operations metrics.
  - *Write/Manage:* Add/edit categories, modify dynamic attribute schemas, update GPS boundaries, adjust platform operational settings, process seller tier upgrades, and approve referral withdrawals.
- **Restricted Actions:** Cannot create or manage other administrative accounts (`/admin/users` is restricted to Platform Owner). Cannot view raw security audit logs (`/admin/logs` requires Security Admin or Platform Owner).
- **Authentication Requirements:** WebAuthn Passkey at `/platform-console/login`.

---

### 3. Security Administrator
- **Code & Database Identifier:** `security_admin`.
- **Database Storage:** `users.admin_role = 'security_admin'`.
- **Purpose:** Platform cybersecurity, compliance, abuse investigation, session termination, tamper-resistant audit trail verification, and fraud response.
- **Allowed Pages & Routes:** `/admin/dashboard`, `/admin/operations`, `/admin/logs`, `/admin/trust-safety`.
- **Actual Backend Permissions:** Satisfies `RequireSecurityAdmin`, `logs`, `trust_safety`, and session termination endpoints.
- **Capabilities:**
  - *Read:* Security event logs, access attempts, admin session directories, abuse reports, IP firewall tables.
  - *Write/Manage:* Suspend fraudulent accounts, revoke active admin or seller sessions, block malicious IPs, manage report dispositions (investigating, resolved, dismissed).
- **Restricted Actions:** Cannot modify product prices, commission rates, catalog taxonomies, or execute financial refunds.
- **Authentication Requirements:** WebAuthn Passkey at `/security-console/login`.

---

### 4. Operations Administrator
- **Code & Database Identifier:** `operations_admin` (legacy alias: `normal`).
- **Database Storage:** `users.admin_role = 'operations_admin'` (or `normal`).
- **Purpose:** Day-to-day marketplace operational oversight, seller onboarding, tier verification, listing approvals, dispute resolution, and referral payout management.
- **Allowed Pages & Routes:** `/admin/dashboard`, `/admin/operations`, `/admin/products`, `/admin/categories`, `/admin/attributes`, `/admin/locations`, `/admin/search`, `/admin/verification`, `/admin/trust-safety`, `/admin/refunds`, `/admin/appeals`, `/admin/upgrades`, `/admin/premium-users`.
- **Actual Backend Permissions:** `products`, `categories`, `attributes`, `locations`, `seller_verification`, `trust_safety`, `refunds`, `appeals`, `seller_upgrades`.
- **Capabilities:**
  - *Read:* Listing queues, seller submissions, verification documents, tier upgrade payment slips.
  - *Write/Manage:* Approve/reject listings, verify seller documents, grant tier upgrades, process refunds, respond to appeals.
- **Restricted Actions:** Cannot access financial planning (`/admin/finance`), profit models (`/admin/profit`), security logs (`/admin/logs`), or admin user management (`/admin/users`).
- **Authentication Requirements:** WebAuthn Passkey or MFA-backed staff credentials at `/staff/login`.

---

### 5. Content Moderator
- **Code & Database Identifier:** `moderator`.
- **Database Storage:** `users.admin_role = 'moderator'`.
- **Purpose:** Marketplace content hygiene, catalog quality, spam elimination, and initial triage of reported listings and abusive comments.
- **Allowed Pages & Routes:** `/admin/dashboard`, `/admin/products`, `/admin/trust-safety`.
- **Actual Backend Permissions:** `pending` (listings awaiting moderation), `products` (listing review), `trust_safety` (review abuse reports).
- **Capabilities:**
  - *Read:* Pending product listings, listing media, user abuse reports.
  - *Write/Manage:* Approve or reject pending listings with moderation reason, mark abuse reports as reviewed.
- **Restricted Actions:** Cannot suspend user accounts directly (requires Security Admin or Platform Admin), cannot modify categories, cannot issue refunds, cannot view financials or settings.
- **Authentication Requirements:** WebAuthn Passkey or MFA-backed credentials at `/staff/login`.

---

### 6. Customer Support Agent
- **Code & Database Identifier:** `support_agent`.
- **Database Storage:** `users.admin_role = 'support_agent'`.
- **Purpose:** Customer inquiries, user assistance, ticket escalation, basic account recovery triage, and verification status checks.
- **Allowed Pages & Routes:** `/admin/dashboard`, `/admin/verification` (view-only), `/admin/appeals` (view/triage).
- **Actual Backend Permissions:** Read permissions for verification queues and appeal submissions.
- **Capabilities:**
  - *Read:* User verification status, seller profiles, public listings, submitted appeal text.
  - *Write/Manage:* Add notes to appeals or verification tickets for supervisor review.
- **Restricted Actions:** Cannot approve/reject verification documents, cannot suspend accounts, cannot view financial logs or settings.
- **Authentication Requirements:** WebAuthn Passkey or MFA-backed credentials at `/staff/login`.

---

### 7. Finance Officer / Finance Administrator (Architectural Analysis)
- **Status in Implementation:** Not implemented as a distinct database enum role.
- **Implementation Mechanism:** In HUNTIQ, financial governance (`/admin/finance` and `/admin/profit`) is governed by high-privilege domain permissions: `finance`, `finance:manage`, `profit`, and `analytics`.
- **Delegation Model:** The Platform Owner or Platform Admin holds these capabilities natively. A dedicated staff member assigned the `operations_admin` role can be granted the `finance` and `profit` permissions via the `admin_permissions` table to function as a Finance Officer without altering the database enum schema.

---

## 3. Role-to-Page Access Matrix

Legend:
- **Manage**: Full read, write, approve, edit, and configuration capabilities.
- **View**: Read-only access to records, analytics, or status queues.
- **Restricted**: Access denied by backend middleware and hidden in navigation.
- **Not verified**: Page/endpoint not implemented or requires unverified mock tokens.

| # | Page Route | Platform Owner | Platform Admin | Security Admin | Operations Admin | Content Moderator | Support Agent | Backend Enforcement |
|---|---|---|---|---|---|---|---|---|
| 1 | `/admin/login` | Manage (Directory) | Manage (Directory) | Manage (Directory) | Manage (Directory) | Manage (Directory) | Manage (Directory) | Public Router |
| 2 | `/executive-console/login` | Manage | Restricted | Restricted | Restricted | Restricted | Restricted | Backend WebAuthn Tier-0 Check |
| 3 | `/platform-console/login` | Manage | Manage | Restricted | Restricted | Restricted | Restricted | Backend WebAuthn Tier-1 Check |
| 4 | `/security-console/login` | Manage | Restricted | Manage | Restricted | Restricted | Restricted | Backend WebAuthn Tier-2 Check |
| 5 | `/staff/login` | Manage | Manage | Manage | Manage | Manage | Manage | Backend WebAuthn Tier-3 / MFA |
| 6 | `/admin` (Redirect) | Manage | Manage | Manage | Manage | Manage | Manage | Frontend & Layout Guard |
| 7 | `/admin/dashboard` | Manage | Manage | Manage | Manage | Manage | Manage | Auth Middleware (`RequireAdmin`) |
| 8 | `/admin/operations` | Manage | Manage | Manage | Manage | Restricted | Restricted | Permission `operations` |
| 9 | `/admin/products` | Manage | Manage | Restricted | Manage | Manage | Restricted | Permission `products` / `pending` |
| 10 | `/admin/categories` | Manage | Manage | Restricted | Manage | Restricted | Restricted | Permission `categories` |
| 11 | `/admin/attributes` | Manage | Manage | Restricted | Manage | Restricted | Restricted | Permission `categories` |
| 12 | `/admin/locations` | Manage | Manage | Restricted | Manage | Restricted | Restricted | Permission `locations` |
| 13 | `/admin/search` | Manage | Manage | Restricted | Manage | Restricted | Restricted | Permission `analytics` |
| 14 | `/admin/finance` | Manage | Manage | Restricted | Restricted | Restricted | Restricted | Permission `finance` / `profit` |
| 15 | `/admin/profit` | Manage | Manage | Restricted | Restricted | Restricted | Restricted | Permission `profit` |
| 16 | `/admin/users` | Manage | Restricted | Restricted | Restricted | Restricted | Restricted | `RequirePlatformOwner` |
| 17 | `/admin/settings` | Manage | Manage | Restricted | Restricted | Restricted | Restricted | Permission `settings` |
| 18 | `/admin/logs` | Manage | Restricted | Manage | Restricted | Restricted | Restricted | `RequireSecurityAdmin` / `logs` |
| 19 | `/admin/verification` | Manage | Manage | Restricted | Manage | Restricted | View | Permission `seller_verification` |
| 20 | `/admin/trust-safety` | Manage | Manage | Manage | Manage | Manage | Restricted | Permission `trust_safety` |
| 21 | `/admin/refunds` | Manage | Manage | Restricted | Manage | Restricted | Restricted | Permission `refunds` |
| 22 | `/admin/appeals` | Manage | Manage | Restricted | Manage | Restricted | View | Permission `appeals` |
| 23 | `/admin/referrals` | Manage | Manage | Restricted | Restricted | Restricted | Restricted | Permission `referrals` |
| 24 | `/admin/upgrades` | Manage | Manage | Restricted | Manage | Restricted | Restricted | Permission `seller_upgrades` |
| 25 | `/admin/premium-users` | Manage | Manage | Restricted | Restricted | Restricted | Restricted | Permission `premium_users` |

---

## 4. Calculated Metric Reconciliation

### 1. Total Verified Roles: **6 DB Enum Roles (+1 Delegated Domain Profile)**
- **Platform Owner** (`platform_owner` / `super`)
- **Platform Admin** (`platform_admin`)
- **Security Admin** (`security_admin`)
- **Operations Admin** (`operations_admin` / `normal`)
- **Moderator** (`moderator`)
- **Support Agent** (`support_agent`)
- *(Finance Officer is delegated via `finance` & `profit` permissions).*
- *Calculation Method:* Verified directly against PostgreSQL enum type `admin_role_type` in migration `029_enterprise_admin_roles_and_security.sql` and `models/passkey.go`.

### 2. Total Admin Pages: **25 Pages**
- **5 Dedicated Login Portals**:
  1. `/admin/login` (Directory router)
  2. `/executive-console/login` (Platform Owner WebAuthn)
  3. `/platform-console/login` (Platform Admin WebAuthn)
  4. `/security-console/login` (Security Admin WebAuthn)
  5. `/staff/login` (Enterprise Staff WebAuthn + Password/MFA)
- **20 Interior Admin Pages**:
  1. `/admin` (Redirect root)
  2. `/admin/dashboard`
  3. `/admin/operations`
  4. `/admin/products`
  5. `/admin/categories`
  6. `/admin/attributes`
  7. `/admin/locations`
  8. `/admin/search`
  9. `/admin/finance`
  10. `/admin/profit`
  11. `/admin/users`
  12. `/admin/settings`
  13. `/admin/logs`
  14. `/admin/verification`
  15. `/admin/trust-safety`
  16. `/admin/refunds`
  17. `/admin/appeals`
  18. `/admin/referrals`
  19. `/admin/upgrades`
  20. `/admin/premium-users`
- *Calculation Method:* Verified via filesystem glob of `frontend/src/app/**/page.tsx` for admin-scoped and console-scoped routes.

### 3. Total Login Portals: **5 Portals**
- Verified routes: `/admin/login`, `/executive-console/login`, `/platform-console/login`, `/security-console/login`, `/staff/login`.

### 4. Total Backend Endpoint Groups Without Dedicated UI: **8 Groups**
- Identified in backend router (`backend/cmd/api/main.go` and `admin_handler.go`):
  1. Sellers list & management (`/api/admin/sellers`)
  2. Global transactions ledger (`/api/admin/transactions`)
  3. System health diagnostics (`/api/admin/system/health`)
  4. Admin notification management (`/api/admin/notifications`)
  5. Promotional banners & carousels (`/api/admin/banners`)
  6. OTP security audit logs (`/api/admin/otp-logs`)
  7. Admin audit activity logs (`/api/admin/activity-logs`)
  8. Category search frequency & trends (`/api/admin/category-trends`)

### 5. Total Orphaned Admin Components: **2 Components**
1. `DocumentViewerModal.tsx`: Present in `frontend/src/components/admin/enterprise/DocumentViewerModal.tsx` for viewing business licenses and ID documents; verification page currently uses static thumbnail placeholders.
2. `MultiStepListingWizard.tsx`: Mounted in `AdminLayout.tsx` but modal trigger from header is rarely invoked in standard admin workflows.

### 6. Role/Page Combination Counts
Based on 6 enum roles across 25 administrative routes (150 total combinations):
- **Verified Manage Combinations:** 78
- **Verified View Combinations:** 2
- **Enforced Restricted Combinations:** 70
- **Problematic Combinations:** 0 (all routes now enforce cryptographic JWT/session validation on both frontend and backend).
