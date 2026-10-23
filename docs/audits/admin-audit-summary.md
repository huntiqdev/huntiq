# HUNTIQ Enterprise Administration & FP&A — Audit Summary & Executive Report

**Document Path**: `docs/audits/admin-audit-summary.md`  
**Audit Date**: October 23, 2026  
**Auditor**: Antigravity FP&A & Enterprise Security Architecture Team  
**Scope**: 
1. Comprehensive page-by-page audit of all HUNTIQ administrative interfaces and backend endpoints.
2. Reconciliation of quarterly vs. monthly growth period compounding in [`backend/internal/fpa/growth.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/growth.go).
3. Verification of separation of duties, RBAC enforcement, session management, and localization parity.

---

## 1. Executive Summary

This audit conducted an exhaustive, read-only discovery and individual inspection of every administrative page, route, portal, layout, and API endpoint across the HUNTIQ (SellBuy) platform. Concurrently, it audited the Financial Planning & Analysis (FP&A) engine to resolve the discrepancy between quarterly business growth assumptions and monthly financial compounding.

### Key Metrics & Discovery Overview
* **Total Admin Pages Discovered**: **25**
  * **21 Internal Admin Pages** under [`frontend/src/app/admin/*`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin)
  * **4 Dedicated Console Portals** ([`/executive-console/login`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/executive-console/login/page.tsx), [`/platform-console/login`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/platform-console/login/page.tsx), [`/security-console/login`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/security-console/login/page.tsx), [`/staff/login`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/staff/login/page.tsx))
* **Inspection Status**: **25 Fully Inspected** (0 Partially Inspected, 0 Blocked, 0 Skipped).
* **Orphaned / Unlinked Pages**: **3 Pages** fully implemented in the file system but omitted from `AdminSidebar` navigation:
  * [`/admin/referrals`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/referrals/page.tsx)
  * [`/admin/upgrades`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/upgrades/page.tsx)
  * [`/admin/premium-users`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/premium-users/page.tsx)
* **Frontend Test Coverage**: **0%** (Zero `.test.tsx` / component tests across all 25 pages).
* **Localization Coverage**: Only **1 page** ([`/admin/finance`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/finance/page.tsx)) supports `en`, `am`, `or`, and `ti`. The remaining **20 internal pages are 100% hardcoded English**.

---

## 2. Growth Reconciliation Summary (Quarterly vs. Monthly)

### The Problem
Strategic planning intentions define a **3-month period as one financial quarter (4 quarters per year)**. However, [`backend/internal/fpa/growth.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/fpa/growth.go#L120-L136) compounds the early rate **every single month**:
$$\text{Current} = \text{Current} \times (1 + r)$$

If an executive inputs a **20% quarterly target** ($r_q = 0.20$), applying it monthly produces:
* **Year 1 Monthly Compounding**: $U_0 \times (1.20)^{12} \approx U_0 \times 8.9161$ (**+791.6% Annual Growth**)
* **Year 1 Intended Quarterly Compounding**: $U_0 \times (1.20)^4 \approx U_0 \times 2.0736$ (**+107.4% Annual Growth**)

This represents a **4.3-fold discrepancy** in active users in Year 1 alone.

### Month 1 to Month 60 Trajectory Comparison (Starting at 1,000 MAU)

| Milestone | Current `growth.go` (Monthly 20% Yr 1, then 2%/mo) | Equivalent Smooth Monthly ($r_m = 6.2659\%$, then 2%/mo) | True Discrete Quarterly (20% step/qtr, then 2%/mo) |
| :--- | :---: | :---: | :---: |
| **Month 1** | 1,000 | 1,000 | 1,000 |
| **Month 3 (Q1 end)** | 1,440 | 1,129 | 1,000 |
| **Month 4 (Q2 start)**| 1,728 | 1,200 | 1,200 |
| **Month 6 (Q2 end)** | 2,488 | 1,355 | 1,200 |
| **Month 12 (Q4 end)**| **7,430** | **1,951** | **1,728** |
| **Month 13 (Mature start)**| 7,579 | 1,990 | 2,074 |
| **Month 24 (Year 2)** | 9,451 | 2,482 | 2,586 |
| **Month 36 (Year 3)** | 11,986 | 3,148 | 3,280 |
| **Month 48 (Year 4)** | 15,201 | 3,992 | 4,160 |
| **Month 60 (Year 5)** | **19,222** | **5,063** | **5,276** |

### Mathematical Recommendation
Adopt **Equivalent Smooth Monthly Compounding** where the monthly rate $r_m$ is derived from the quarterly target $r_q$:
$$r_m = (1 + r_q)^{1/3} - 1$$
This reconciles strategic quarterly milestones with the granular monthly requirements of the financial engine (P&L, payroll, AWS server costs, and VAT returns) without creating artificial quarterly volatility.

---

## 3. High-Risk Findings & Vulnerabilities

### Critical & High Severity Issues

1. **Finance Separation of Duties Breakdown (`/admin/finance`) — HIGH (`AUD-FIN-001`)**
   * *Risk*: In `models/user.go`, `RolePlatformAdmin` receives baseline `"profit": true, "finance": true`. Both read (`GET /finance/forecast`) and write (`PUT /finance/assumptions`, `POST /finance/assumptions/reset`) routes share the exact same `financeAuth` middleware.
   * *Impact*: Any operational Platform Administrator can arbitrarily alter Ethiopian tax rates, corporate assumptions, or reset financial forecasts.
   * *Remediation*: Split into `"finance:view"` and `"finance:manage"`. Gating write endpoints behind `RequirePlatformOwner` or an explicit Chief Financial Officer permission.

2. **Illusion of User Ban Enforcement (`/admin/trust-safety`) — HIGH (`AUD-TRU-001`)**
   * *Risk*: When a moderator reviews an abuse report and chooses the action `"banned"`, `admin_service.go:ReviewAbuseReport` merely writes the string `"banned"` to the `abuse_reports` table.
   * *Impact*: The offending user account in the `users` table is **never suspended or deactivated**. They can continue to log in, post fraudulent listings, and message buyers.
   * *Remediation*: Enforce an automated database transaction calling `BanUser()` and revoking active sessions when `"banned"` is confirmed by an authorized admin.

3. **Legacy `isSuperAdmin` Frontend Access Barriers — HIGH (`AUD-PRF-001`, `AUD-CAT-001`, `AUD-ATT-001`, `AUD-USR-001`, `AUD-REF-001`)**
   * *Risk*: Five admin pages check legacy `isSuperAdmin` (`role_id == 1` or `user.is_super`) instead of role/permission checks:
     * [`/admin/categories`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/categories/page.tsx): "Add Child", "Edit", and "Delete" buttons are hidden from Platform Admins.
     * [`/admin/attributes`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/attributes/page.tsx): "New Attribute" button is hidden from Platform Admins.
     * [`/admin/profit`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/profit/page.tsx): Access completely blocked with *"Super admin access required"*.
     * [`/admin/referrals`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/referrals/page.tsx): Access completely blocked.
     * [`/admin/users`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/users/page.tsx): Prematurely redirects Platform Owners to `/admin/dashboard` if logging in via executive portal without legacy `is_super`.
   * *Remediation*: Standardize checks to `isPlatformOwner || isPlatformAdmin || hasAdminPermission(...)`.

4. **Stale/Empty LocalStorage Token Usage (`/admin/locations`) — HIGH (`AUD-LOC-001`)**
   * *Risk*: Line 84 reads `localStorage.getItem('token')` instead of auth context `user?.token`.
   * *Impact*: Admins authenticated via WebAuthn/passkey store tokens in memory/cookies, leading to HTTP 401 unauthenticated errors on location hierarchy load.

5. **Audit Reason Bypass (`/admin/settings`) — MEDIUM (`AUD-SET-001`)**
   * *Risk*: Saves platform configurations via `PUT /api/admin/settings` instead of the hardened `PUT /api/admin/settings/with-reason`.
   * *Impact*: Changes to transaction fee caps, maintenance mode, and merchant IDs leave `settings_change_log` unpopulated.

6. **Three Orphaned Administrative Pages — HIGH (`AUD-REF-001`, `AUD-UPG-001`, `AUD-PRE-001`)**
   * *Risk*: Referral payouts ([`/admin/referrals`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/referrals/page.tsx)), tier upgrade requests ([`/admin/upgrades`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/upgrades/page.tsx)), and premium subscriber management ([`/admin/premium-users`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/premium-users/page.tsx)) are missing from `AdminSidebar`.
   * *Impact*: Staff cannot easily navigate to monetization and payout workflows without manual URL typing.

---

## 4. Testing & Verification Limitations

* **Automated Backend Tests Passed**:
  * `go test -v ./internal/fpa/...` — All 11 FP&A test suites pass (run time: 0.057s).
  * `go test -v ./internal/middleware/admin_authorization_matrix_test.go` — Authorization matrix tests pass.
* **Frontend Test Deficit**:
  * **0 test files exist in `frontend/src`**.
  * No component tests, unit tests, or Playwright E2E tests exist for the admin console.
* **Browser Test Environment**:
  * Interactive headless browser testing is limited in this phase due to read-only constraints (avoiding database mutation or mock token pollution).
  * Recommended next step: Deploy a dedicated staging mock database with seeded admin accounts (`owner@sellbuy.et`, `admin@sellbuy.et`, `security@sellbuy.et`) to run Playwright E2E coverage across all 25 routes.

---

## 5. Prioritized Remediation Roadmap

```mermaid
flowchart TD
    subgraph Phase 1: Security & Access Hardening [Immediate P0]
        P1A[Split Finance View / Manage Permissions] --> P1B[Fix Ban Enforcement in Trust & Safety]
        P1B --> P1C[Replace Legacy isSuperAdmin UI Checks]
        P1C --> P1D[Fix LocalStorage Token in /admin/locations]
        P1D --> P1E[Enforce with-reason in /admin/settings]
    end

    subgraph Phase 2: Navigation & FP&A Growth [P1]
        P2A[Mount Orphaned Pages in AdminSidebar] --> P2B[Implement Smooth Quarterly Compounding in growth.go]
        P2B --> P2C[Expose Growth Period Selector in /admin/finance UI]
    end

    subgraph Phase 3: Quality & Localization [P2]
        P3A[Build Playwright E2E Test Suite for 25 Pages] --> P3B[Extend Localization to All Admin Screens]
        P3B --> P3C[Implement Attribute Edit/Delete & Dispute Console]
    end

    Phase 1 --> Phase 2 --> Phase 3
```

1. **Phase 1 (Immediate / P0 — Security & Integrity)**:
   * Isolate write access to FP&A financial assumptions.
   * Connect abuse report banning directly to `users` deactivation.
   * Purge all `isSuperAdmin` checks in favor of role/permission functions.
   * Fix token extraction in location hierarchy.
2. **Phase 2 (P1 — Navigation & Growth Model)**:
   * Add `/admin/referrals`, `/admin/upgrades`, and `/admin/premium-users` to `AdminSidebar`.
   * Add quarterly growth parameters (`early_quarterly_rate`, `mature_quarterly_rate`) to `types.go` and implement smooth derivation in `growth.go`.
3. **Phase 3 (P2 — Testing & Completeness)**:
   * Introduce Jest/Playwright tests for all admin screens.
   * Translate all 20 English-only admin interfaces into Amharic, Oromo, and Tigrinya.

---
*End of HUNTIQ Admin Audit Summary.*
