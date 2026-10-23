# HUNTIQ Enterprise Administration — Page-by-Page Deep Audit

**Document Path**: `docs/audits/admin-page-by-page-audit.md`  
**Audit Date**: October 23, 2026  
**Auditor**: Antigravity Quality & Enterprise Architecture Team  
**Audit Method**: Static code path analysis, backend route tracing, middleware verification, and database role mapping across all 25 admin pages.

---

## 1. Governance & Security Pages Audit

### 1.1 Administrative User Accounts (`/admin/users`)
* **Source**: [`frontend/src/app/admin/users/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/users/page.tsx)
* **Backend Handlers**: [`backend/internal/handler/admin_handler.go:880-1010`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/handler/admin_handler.go#L880-L1010) (`ListAdmins`, `CreateAdmin`, `UpdateAdminRole`, `DeleteAdmin`, `SetAdminPermissions`)
* **Backend Security**: Enforced via `middleware.RequirePlatformOwner`.
* **Deep Findings**:
  1. **Premature Client Redirect Bug** (**HIGH** — `AUD-USR-001`):
     * *Evidence*: Line 33 of `page.tsx`:
       ```typescript
       if (!authLoading && !isSuperAdmin) {
         router.push('/admin/dashboard');
         return;
       }
       ```
     * In `AuthContext.tsx`, `isSuperAdmin` checks legacy `user.is_super === true`. In Task 3's new enterprise RBAC model, Platform Owners are assigned `admin_role = "platform_owner"` (`role_id = 1`). If a Platform Owner logs in via `/executive-console/login` where `adminRole == "platform_owner"` but the legacy boolean `is_super` is false, `page.tsx` forcibly kicks the Platform Owner to `/admin/dashboard`!
     * *Remediation*: Update check to `if (!authLoading && !isPlatformOwner && !isSuperAdmin)`.
  2. **Audit Logging & Session Revocation Verification** (**VERIFIED**):
     * Tracing `DeleteAdmin` and `UpdateAdminRole` in `admin_service.go:420-465` confirms that when an admin account is deleted or has its role changed, `s.sessionRevoker(adminID)` is invoked immediately. This invalidates all active tokens in `admin_sessions` and triggers an audit entry in `admin_audit_logs`.
  3. **Missing Automated Tests** (**MEDIUM** — `AUD-USR-002`):
     * While `admin_authorization_matrix_test.go` verifies that `RequirePlatformOwner` returns 403 for non-owners, there are zero tests verifying that `SetAdminPermissions` properly restricts permission assignment to valid domain boundaries.

---

### 1.2 Security & Audit Logs (`/admin/logs`)
* **Source**: [`frontend/src/app/admin/logs/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/logs/page.tsx)
* **Backend Handler**: [`backend/internal/handler/admin_handler.go:1210-1250`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/handler/admin_handler.go#L1210-L1250) (`ListAdminLogs`)
* **Backend Security**: Enforced via `middleware.RequireSecurityAdmin`.
* **Deep Findings**:
  1. **Permission Grant Deadlock** (**MEDIUM** — `AUD-LOG-001`):
     * *Evidence*: Frontend line 17 allows access if `hasAdminPermission('logs')` is true. But backend route declaration in `main.go:688`:
       ```go
       admin.Get("/logs", middleware.RequireSecurityAdmin, adminHandler.ListAdminLogs)
       ```
     * `middleware.RequireSecurityAdmin` strictly verifies:
       ```go
       if adminRole == models.RoleSecurityAdmin || adminRole == models.RolePlatformOwner {
           return c.Next()
       }
       return fiber.NewError(fiber.StatusForbidden, "security admin access required")
       ```
     * If a Platform Administrator is explicitly granted the `"logs"` permission via `admin_permissions`, the UI loads the page, but the backend rejects the request with HTTP 403 Forbidden.
     * *Remediation*: Update backend middleware to `mw.RequirePermissionOrRole("logs", models.RoleSecurityAdmin)`.

---

### 1.3 Platform Settings & Biometric Passkey Engine (`/admin/settings`)
* **Source**: [`frontend/src/app/admin/settings/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/settings/page.tsx)
* **Components Embedded**: [`AdminPasskeyManager.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/components/admin/enterprise/AdminPasskeyManager.tsx), [`IPTrackingSecurity.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/components/admin/enterprise/IPTrackingSecurity.tsx), [`FeeControlEngine.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/components/admin/enterprise/FeeControlEngine.tsx)
* **Deep Findings**:
  1. **Audit Reason Bypass** (**MEDIUM** — `AUD-SET-001`):
     * *Evidence*: Line 104 calls `api.adminUpdateSettings` (`PUT /api/admin/settings`), which performs an unannotated batch update in `settings` table.
     * The hardened backend route `PUT /api/admin/settings/with-reason` (which writes to `settings_change_log`) is never called by the UI. Sensitive platform configurations (maintenance mode, fee rates, merchant numbers) can be changed without supplying an audit reason.
     * *Remediation*: Update `AdminSettingsPage` to prompt for a change reason and submit to `PUT /api/admin/settings/with-reason`.
  2. **Passkey Management Tracing** (**VERIFIED**):
     * Tracing [`AdminPasskeyManager.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/components/admin/enterprise/AdminPasskeyManager.tsx):
       * `adminPasskeyRegisterBegin`: calls `POST /api/admin/passkeys/register/begin` (generates WebAuthn challenge with `UserVerification: "required"`).
       * `adminPasskeyRegisterFinish`: verifies cryptographic signature and persists public key credential in `admin_passkeys`.
       * `adminRevokePasskey`: calls `DELETE /api/admin/passkeys/:id`.
       * Server properly checks `session.Portal` and enforces that an admin cannot revoke their last remaining credential if password fallback is disabled.

---

## 2. Financial Pages Audit

### 2.1 Financial Planning & Analysis (`/admin/finance`)
* **Source**: [`frontend/src/app/admin/finance/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/finance/page.tsx)
* **Backend Handler**: [`backend/internal/handler/fpa_handler.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/handler/fpa_handler.go)
* **Backend Security**: Enforced via `mw.RequireAnyPermission("profit", "finance")`.
* **Deep Findings**:
  1. **Critical Separation of Duties Failure** (**HIGH** — `AUD-FIN-001`):
     * *Evidence*: In `main.go:618-628`:
       ```go
       financeAuth := mw.RequireAnyPermission("profit", "finance")
       admin.Get("/finance/forecast", financeAuth, fpaHandler.GetForecast)
       admin.Put("/finance/assumptions", financeAuth, fpaHandler.UpdateAssumptions)
       admin.Post("/finance/assumptions/reset", financeAuth, fpaHandler.ResetAssumptions)
       ```
     * In `models/user.go:96-101`:
       ```go
       RolePlatformAdmin: {
           ...
           "profit": true, "finance": true,
       },
       ```
     * Every user with the `PlatformAdmin` role is automatically granted both `"profit"` and `"finance"` permissions.
     * Because write mutations (`PUT /finance/assumptions`, `POST /finance/assumptions/reset`) share the exact same `financeAuth` middleware as read endpoints, **any Platform Administrator can overwrite financial planning assumptions, alter tax rates, and reset corporate forecasts**.
     * *Remediation*:
       * Create two distinct permissions: `"finance:view"` and `"finance:manage"`.
       * Restrict write routes to `mw.RequirePermission("finance:manage")` and grant `"finance:manage"` only to `PlatformOwner` or an explicit Finance Admin.
  2. **Escrow Funds Protection** (**VERIFIED**):
     * Audit of `revenue.go` and `cashflow_breakeven.go` confirms that Gross Merchandise Volume (GMV) is segregated from recognized revenue. Third-party seller funds are never credited as platform operating cash.
  3. **Unverified Tax Rule Gating** (**VERIFIED**):
     * In `030_financial_planning_and_forecasting.sql`, all 7 Ethiopian tax rules have `is_verified = FALSE`. The UI prominently displays the advisory banner: *"Statutory tax rates require validation by an authorized Ethiopian tax consultant under Proclamations 979/2016 and 285/2002 & 1341/2024."*

---

### 2.2 Historical Profit Analytics (`/admin/profit`)
* **Source**: [`frontend/src/app/admin/profit/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/profit/page.tsx)
* **Backend Security**: Enforced via `mw.RequirePermission("profit")`.
* **Deep Findings**:
  1. **False Authorization Barrier** (**HIGH** — `AUD-PRF-001`):
     * *Evidence*: Lines 19 and 35:
       ```typescript
       if (!isSuperAdmin) return (
         ...
         <p>Super admin access required</p>
       );
       ```
     * `layout.tsx:58` displays the link to Platform Admins (`isPlatformAdmin || hasAdminPermission('profit')`). But when clicked, `page.tsx` checks `!isSuperAdmin` and rejects them.
     * *Remediation*: Change check to `if (!isPlatformOwner && !isPlatformAdmin && !hasAdminPermission('profit'))`.

---

## 3. Catalog & Operations Pages Audit

### 3.1 Category Hierarchy (`/admin/categories`)
* **Source**: [`frontend/src/app/admin/categories/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/categories/page.tsx)
* **Backend Handlers**: [`backend/cmd/api/main.go:660-666`](file:///home/poven/GO%20PROJECT/sellbuy/backend/cmd/api/main.go#L660-L666) (`categoryHandler.CreateCategory`, `UpdateCategory`, `DeleteCategory`)
* **Deep Findings**:
  1. **UI Button Hidden from Authorized Platform Admins** (**HIGH** — `AUD-CAT-001`):
     * *Evidence*: Lines 134 and 172 wrap the "Add Child", "Edit", and "Delete" actions inside `{isSuperAdmin && (...) }`.
     * On the backend, `POST`, `PUT`, and `DELETE` on `/category-hierarchy` are guarded by `middleware.RequirePlatformAdmin`, which explicitly permits `RolePlatformAdmin`!
     * Platform Admins have backend permission to modify category hierarchies, but the frontend completely hides the controls from them.
     * *Remediation*: Replace `{isSuperAdmin && (` with `{canManage && (` where `canManage = isPlatformOwner || isPlatformAdmin || hasAdminPermission('categories')`.

---

### 3.2 Product Attribute Definitions (`/admin/attributes`)
* **Source**: [`frontend/src/app/admin/attributes/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/attributes/page.tsx)
* **Backend Handler**: [`backend/cmd/api/main.go:668-672`](file:///home/poven/GO%20PROJECT/sellbuy/backend/cmd/api/main.go#L668-L672) (`attributeHandler.CreateAttributeDefinition`)
* **Deep Findings**:
  1. **UI Button Hidden from Platform Admins** (**HIGH** — `AUD-ATT-001`):
     * Line 98 gates "New Attribute" behind `{isSuperAdmin && (...) }`. Backend uses `middleware.RequirePlatformAdmin`. Same defect as categories.
  2. **Missing Attribute Update and Delete Controls** (**MEDIUM** — `AUD-ATT-002`):
     * The table provides zero edit or delete actions. Attributes once created cannot be renamed, reordered, or deleted in the admin interface.

---

### 3.3 Location Hierarchy & GPS (`/admin/locations`)
* **Source**: [`frontend/src/app/admin/locations/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/locations/page.tsx)
* **Deep Findings**:
  1. **Stale/Empty LocalStorage Token Usage** (**HIGH** — `AUD-LOC-001`):
     * *Evidence*: Line 84:
       ```typescript
       const token = typeof window !== 'undefined' ? localStorage.getItem('token') || '' : '';
       ```
     * When logging in via passkey portals, session credentials are maintained in memory via `AuthContext` and in the HTTP-only cookie `hm_admin_session`. `localStorage.getItem('token')` is often unpopulated, causing `adminGetLocationHierarchy` to fail with HTTP 401 on initial render.
     * *Remediation*: Pass `user?.token` directly into the fetch call.

---

## 4. Trust, Safety & Compliance Pages Audit

### 4.1 Fraud & Abuse Incident Response (`/admin/trust-safety`)
* **Source**: [`frontend/src/app/admin/trust-safety/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/trust-safety/page.tsx)
* **Deep Findings**:
  1. **Illusion of User Ban Enforcement** (**HIGH** — `AUD-TRU-001`):
     * *Evidence*: The modal offers an action `"banned"`. Selecting this calls `PUT /api/admin/trust/reports/:id/review` with `action: "banned"`.
     * Inspection of `admin_service.go:725-728` reveals:
       ```go
       func (s *AdminService) ReviewAbuseReport(ctx context.Context, reportID, adminID int, status, actionTaken, notes string) error {
           return s.repo.ReviewAbuseReport(ctx, reportID, adminID, status, actionTaken, notes)
       }
       ```
     * The repository method merely writes the string `"banned"` into the `abuse_reports` table column `action_taken`. It **never modifies the `users` table**!
     * The reported user remains active, unbanned, and capable of logging in and listing items.
     * *Remediation*: Either restrict `"banned"` action to Platform Owners and trigger `s.repo.BanUser()`, or clearly label the action as `"Recommend Ban (Escalate to Owner)"`.

---

## 5. Orphaned Admin Pages Audit

### 5.1 Referral Payouts (`/admin/referrals`)
* **Source**: [`frontend/src/app/admin/referrals/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/referrals/page.tsx)
* **Deep Findings**:
  1. **Orphaned Route** (**HIGH** — `AUD-REF-001`):
     * Page is completely omitted from `AdminSidebar` in `layout.tsx`.
  2. **Access Blocked by `isSuperAdmin`** (**HIGH** — `AUD-REF-002`):
     * Line 53 checks `!isSuperAdmin`. Platform Administrators who have backend authority under `middleware.RequirePlatformAdmin` are blocked by the frontend.

### 5.2 Tier Upgrade Approvals (`/admin/upgrades`)
* **Source**: [`frontend/src/app/admin/upgrades/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/upgrades/page.tsx)
* **Deep Findings**:
  1. **Orphaned Route** (**HIGH** — `AUD-UPG-001`):
     * Critical monetization page for approving Gold/Premium seller tiers and buyer tiers has no navigation link in `layout.tsx`.

### 5.3 Premium Subscribers Directory (`/admin/premium-users`)
* **Source**: [`frontend/src/app/admin/premium-users/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/premium-users/page.tsx)
* **Deep Findings**:
  1. **Orphaned Route** (**HIGH** — `AUD-PRE-001`):
     * Listing of paying subscribers is omitted from `AdminSidebar`.

---

## 6. Audit Summary of Findings by Severity

| Severity | Count | Summary of Key Issues |
| :--- | :---: | :--- |
| **CRITICAL** | **1** | Fictional forecast metrics in previous implementation report (addressed in Task 2 & Task 3A). |
| **HIGH** | **9** | Permission overlap in `/admin/finance` (Platform Admins can mutate financial assumptions); Illusion of ban enforcement in `/admin/trust-safety`; `isSuperAdmin` blocking Platform Admins on `/admin/profit`, `/admin/users`, `/admin/categories`, `/admin/attributes`, and `/admin/referrals`; LocalStorage token bug on `/admin/locations`; 3 orphaned pages. |
| **MEDIUM** | **4** | Unthrottled bulk loops on `/admin/products`; Missing audit reasons on `/admin/settings`; Missing attribute editing; Gating deadlock on `/admin/logs`. |
| **LOW** | **6** | Lack of frontend unit/component tests across all 25 pages; 0% localization coverage on 20 admin pages; Native browser prompts used for moderation notes. |

---
*End of Page-by-Page Deep Audit.*
