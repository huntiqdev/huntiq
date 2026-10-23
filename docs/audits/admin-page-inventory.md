# HUNTIQ Enterprise Administration — Comprehensive Page Inventory & Routing Matrix

**Document Path**: `docs/audits/admin-page-inventory.md`  
**Audit Date**: October 23, 2026  
**Auditor**: Antigravity Quality & Enterprise Architecture Team  
**Total Admin Pages Discovered**: 25 (21 inside `/admin/*`, 4 dedicated console login portals)  
**Total Orphaned / Unlinked Pages Discovered**: 3 (`/admin/referrals`, `/admin/upgrades`, `/admin/premium-users`)  
**Audit Status**: 25 Pages Fully Inspected (0 Blocked, 0 Skipped)

---

## 1. Discovered Admin Pages Inventory & Inspection Matrix

### Page 01: `/admin` (Root Admin Entrypoint)
* **Frontend Route**: `/admin`
* **Source File**: [`frontend/src/app/admin/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/page.tsx)
* **Shared Layout / Components**: `AdminLayout` ([`frontend/src/app/admin/layout.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/layout.tsx))
* **Intended Roles**: All administrative identities (`PlatformOwner`, `PlatformAdmin`, `SecurityAdmin`, `OperationsAdmin`, `Moderator`, `SupportAgent`)
* **Required Permission Identifiers**: None (Entrypoint redirect)
* **Backend API Endpoints**: None
* **Capabilities**: Client-side Next.js route redirection to `/admin/dashboard`
* **Auth Enforcement**: Handled by `AdminLayout` (redirects unauthenticated sessions to `/admin/login`)
* **UI States**: Immediate server redirect (`redirect('/admin/dashboard')`)
* **Localization Coverage**: N/A (Redirect only)
* **Test Coverage**: Tested indirectly via Next.js route compiler (`npm run build`). No unit test.
* **Findings & Severity**: **LOW** (`INV-001`). Clean redirect, no security issues.
* **Audit Status**: **VERIFIED**

---

### Page 02: `/admin/dashboard` (Platform Overview & Operations Hub)
* **Frontend Route**: `/admin/dashboard`
* **Source File**: [`frontend/src/app/admin/dashboard/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/dashboard/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, `AdminSidebar`, `AdminHeaderTopBar`, `KPICard`, `QuickAction`, `MonthlyTrend`, `Spinner`
* **Intended Roles**: All Admin Roles
* **Required Permission Identifiers**: `"analytics"` (to view stats), `"pending"`, `"sellers"`, `"seller_verification"`, `"trust_safety"`, `"operations"` (for quick-action cards)
* **Backend API Endpoints**:
  * `GET /api/admin/analytics/overview` (Gated by `RequirePermission("analytics")`)
  * `GET /api/admin/analytics/monthly?months=3` (Gated by `RequirePermission("analytics")`)
* **Capabilities**: Read-only KPI aggregation, time-series monthly trend, quick navigation
* **Auth Enforcement**: Server-side JWT / session cookie validation via `mw.RequireAuth` + `RequirePermission("analytics")`
* **UI States**: 
  * Loading: `<Spinner />`
  * Empty: Fallback to placeholder dashes (`—`)
  * Error: Silently caught, defaults to `{}`
  * Unauthorized: Falls back to empty stats if `"analytics"` permission is missing
* **Localization Coverage**: 0% (All text hardcoded English: "Platform Overview", "Total Listings", "Pending Review")
* **Test Coverage**: No unit/component tests in `frontend/src`. Backend analytics covered in `admin_authorization_matrix_test.go`.
* **Findings & Severity**:
  * **MEDIUM** (`INV-002`): Financial KPIs (`total_revenue`, `listing_fee_revenue`) check legacy `isSuperAdmin` (`role_id == 1`) instead of `isPlatformOwner || isPlatformAdmin || hasAdminPermission('profit')`. Platform Admins are blocked from seeing high-level revenue figures on the dashboard.
  * **LOW** (`INV-003`): Missing localization strings.
* **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 03: `/admin/operations` (Background Operations & Job Monitor)
* **Frontend Route**: `/admin/operations`
* **Source File**: [`frontend/src/app/admin/operations/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/operations/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, `AdminSidebar`, `AdminHeaderTopBar`, `StatusBadge`, `StatusIcon`, `KPIBox`
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`, `SecurityAdmin`
* **Required Permission Identifiers**: `"operations"`
* **Backend API Endpoints**:
  * `GET /api/admin/operations/health` (Gated by `RequirePermission("operations")`)
  * `GET /api/admin/operations/jobs` (Gated by `RequirePermission("operations")`)
  * `GET /api/admin/operations/failed-jobs` (Gated by `RequirePermission("operations")`)
* **Capabilities**: Read-only monitoring of 24h job runs, failed jobs, table sizes, auto-refresh every 30s
* **Auth Enforcement**: Server-side enforced by `mw.RequirePermission("operations")`
* **UI States**:
  * Loading: Skeleton pulsing cards
  * Empty: Renders table headers with empty rows
  * Error / Unauthorized: `Promise.allSettled` catches 403; page renders with `null` health and empty jobs rather than displaying an explicit 403 Forbidden alert
* **Localization Coverage**: 0% (Hardcoded English)
* **Test Coverage**: No frontend test. Backend endpoints covered in `admin_authorization_matrix_test.go`.
* **Findings & Severity**: **LOW** (`INV-004`). Missing client-side unauthorized banner when 403 occurs.
* **Audit Status**: **VERIFIED**

---

### Page 04: `/admin/products` (Catalog & Listing Moderation Console)
* **Frontend Route**: `/admin/products`
* **Source File**: [`frontend/src/app/admin/products/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/products/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, `ProductReview`, `MultiStepListingWizard`
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`, `Moderator`
* **Required Permission Identifiers**: `"pending"` (to review pending), `"products"` (to view catalog & fee claims)
* **Backend API Endpoints**:
  * `GET /api/admin/products/pending` (`RequirePermission("pending")`)
  * `GET /api/admin/products` (`RequirePermission("products")`)
  * `PUT /api/admin/products/:id/review` (`RequirePermission("pending")`)
  * `DELETE /api/admin/products/:id` (`RequirePermission("products")`)
  * `GET /api/admin/fee-payments/pending` (`RequirePermission("products")`)
  * `PUT /api/admin/fee-payments/:id/verify` (`RequirePermission("products")`)
  * `PUT /api/admin/fee-payments/:id/reject` (`RequirePermission("products")`)
* **Capabilities**: List, filter by vertical (`sell`, `rent`, `lease`), approve, reject, delete listing, verify listing fees, bulk review
* **Auth Enforcement**: Server-side role & permission enforcement on every endpoint
* **UI States**: Grid vs. Table view toggle, loading spinner, empty list messages
* **Localization Coverage**: 0% (Hardcoded English)
* **Test Coverage**: Backend covered in `product_handler_test.go`. Frontend has zero tests.
* **Findings & Severity**:
  * **MEDIUM** (`INV-005`): Bulk review iterates over `selectedIds` in an unthrottled client-side `for` loop executing separate HTTP PUT calls. A bulk action on 50 listings fires 50 individual requests.
  * **LOW** (`INV-006`): Uses browser native `prompt()` for rejection notes.
* **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 05: `/admin/categories` (Category Hierarchy Management)
* **Frontend Route**: `/admin/categories`
* **Source File**: [`frontend/src/app/admin/categories/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/categories/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, Category Tree Accordion Modal
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
* **Required Permission Identifiers**: `"categories"`
* **Backend API Endpoints**:
  * `GET /api/categories/tree?lang=en` (Public read)
  * `POST /api/admin/category-hierarchy` (Gated by `middleware.RequirePlatformAdmin`)
  * `PUT /api/admin/category-hierarchy/:id` (Gated by `middleware.RequirePlatformAdmin`)
  * `DELETE /api/admin/category-hierarchy/:id` (Gated by `middleware.RequirePlatformAdmin`)
* **Capabilities**: Tree expansion, create category, edit multilingual names (`en`, `am`, `om`), slug generation, toggle active/leaf, delete
* **Auth Enforcement**: Server-side enforced with `middleware.RequirePlatformAdmin`
* **UI States**: Expand/collapse tree nodes, edit modal, delete confirmation modal, error banner
* **Localization Coverage**: Has inputs for `name_en`, `name_am`, `name_om`, but admin UI labels are English-only.
* **Test Coverage**: Untested on frontend.
* **Findings & Severity**:
  * **HIGH** (`INV-007`): **UI Capability Mismatch**: Lines 134 and 172 gate the "Add Child", "Edit", and "Delete" buttons behind `{isSuperAdmin && (`. A Platform Administrator (`role_id = 2`) has full backend authority (`middleware.RequirePlatformAdmin`) to manage categories, but the frontend hides all editing buttons from them!
* **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 06: `/admin/attributes` (Dynamic Listing Attribute Definitions)
* **Frontend Route**: `/admin/attributes`
* **Source File**: [`frontend/src/app/admin/attributes/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/attributes/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, Attribute Create Modal
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`
* **Required Permission Identifiers**: `"categories"`
* **Backend API Endpoints**:
  * `GET /api/attributes` (Public definition list)
  * `POST /api/admin/attributes` (Gated by `middleware.RequirePlatformAdmin`)
* **Capabilities**: List all dynamic attributes (Brand, Storage, Fuel Type, etc.), define new attribute with data types (`text`, `number`, `boolean`, `select`, `multiselect`, `range`, `date`)
* **Auth Enforcement**: Backend enforces `middleware.RequirePlatformAdmin`
* **UI States**: Empty table message, loading spinner, creation modal with validation rules
* **Localization Coverage**: Multilingual name inputs for `name_en`, `name_am`, `name_om`.
* **Test Coverage**: No frontend tests. Backend handlers lack dedicated tests.
* **Findings & Severity**:
  * **HIGH** (`INV-008`): "New Attribute" button is gated by `{isSuperAdmin && (` on line 98, blocking Platform Admins from creating attributes in the UI.
  * **MEDIUM** (`INV-009`): No Edit or Delete buttons exist in the UI for attribute definitions (read + create only).
* **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 07: `/admin/locations` (Location Hierarchy & GPS Geocoding Admin)
* **Frontend Route**: `/admin/locations`
* **Source File**: [`frontend/src/app/admin/locations/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/locations/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, Tabbed view (`tree`, `list`, `geocoding`, `import`), Location Edit Modal
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
* **Required Permission Identifiers**: `"locations"`
* **Backend API Endpoints**:
  * `GET /api/locations/tree?lang=en`
  * `GET /api/admin/location-hierarchy` (`RequirePermission("locations")`)
  * `POST /api/admin/location-hierarchy` (`RequirePermission("locations")`)
  * `PUT /api/admin/location-hierarchy/:id` (`RequirePermission("locations")`)
  * `DELETE /api/admin/location-hierarchy/:id` (`RequirePermission("locations")`)
  * `POST /api/admin/location-hierarchy/import` (`RequirePermission("locations")`)
  * `POST /api/admin/location-hierarchy/test-geocode` (`RequirePermission("locations")`)
* **Capabilities**: 5-level geographic tree navigation (Country $\rightarrow$ Region $\rightarrow$ Zone/City $\rightarrow$ Woreda $\rightarrow$ Kebele), reverse GPS sandbox testing, batch JSON location import
* **Auth Enforcement**: Server-side enforced on all mutation routes
* **UI States**: Fully responsive tree view, search filter, tab switcher, geocoding pin tester
* **Localization Coverage**: Amharic, Oromo, Tigrinya name inputs supported in schema.
* **Test Coverage**: Backend covered in `location_service_test.go` and `geocoding_test.go`. Zero frontend tests.
* **Findings & Severity**:
  * **HIGH** (`INV-010`): Line 84 extracts auth token via `localStorage.getItem('token')` instead of `user?.token` from `useAuth()`. Stale or missing `localStorage` values trigger 401 Unauthorized errors on initial load.
* **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 08: `/admin/search` (Search Intelligence & Zero-Result Analytics)
* **Frontend Route**: `/admin/search`
* **Source File**: [`frontend/src/app/admin/search/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/search/page.tsx)
* **Shared Layout / Components**: `AdminLayout`
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
* **Required Permission Identifiers**: `"analytics"`
* **Backend API Endpoints**:
  * `GET /api/admin/search/popular` (`RequirePermission("analytics")`)
  * `GET /api/admin/search/zero-results` (`RequirePermission("analytics")`)
* **Capabilities**: Read-only display of top 30 search queries, CTR percentages, and top 50 zero-result search queries
* **Auth Enforcement**: Backend `RequirePermission("analytics")`
* **UI States**: Skeleton loading, empty state messages, auto-scroll lists
* **Localization Coverage**: 0%
* **Test Coverage**: No automated tests.
* **Findings & Severity**: **LOW** (`INV-011`). Missing client-side unauthorized notice if API returns 403.
* **Audit Status**: **VERIFIED**

---

### Page 09: `/admin/finance` (Financial Planning & Analysis / FP&A)
* **Frontend Route**: `/admin/finance`
* **Source File**: [`frontend/src/app/admin/finance/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/finance/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, 9 tabbed modules, SVG charts, 20M capacity stress analyzer, statutory tax catalog, assumption editor modal
* **Intended Roles**: `PlatformOwner`, Authorized Finance Administrators
* **Required Permission Identifiers**: `"profit"`, `"finance"`
* **Backend API Endpoints**:
  * `GET /api/admin/finance/forecast` (`financeAuth`)
  * `GET /api/admin/finance/scenarios` (`financeAuth`)
  * `GET /api/admin/finance/capacity-20m` (`financeAuth`)
  * `GET /api/admin/finance/tax-rules` (`financeAuth`)
  * `GET /api/admin/finance/actuals` (`financeAuth`)
  * `GET /api/admin/finance/audit-logs` (`financeAuth`)
  * `PUT /api/admin/finance/assumptions` (`financeAuth`)
  * `POST /api/admin/finance/assumptions/reset` (`financeAuth`)
  * `GET /api/admin/finance/export` (`financeAuth`)
* **Capabilities**: Multi-scenario forecasting (12–60 mos), workload demand metrics, Ethiopian statutory tax estimation, cash burn/runway tracking, assumption editing, CSV/JSON export
* **Auth Enforcement**: Enforced by `mw.RequireAnyPermission("profit", "finance")`
* **UI States**: Full client-side authorization banner, spinner, 5-year table, interactive chart hover, assumption modals
* **Localization Coverage**: 100% localized across English, Amharic, Oromo, and Tigrinya in `frontend/src/lib/language/{en,am,or,ti}.json`
* **Test Coverage**: Backend covered in `fpa_test.go` (6 unit test suites) and `fpa_handler_test.go` (6 HTTP integration tests).
* **Findings & Severity**:
  * **HIGH** (`INV-012`): **Over-Privileged Platform Admins**: `RolePlatformAdmin` receives `"finance": true` by default in `models/user.go`. Both read and write endpoints share `financeAuth`, allowing infrastructure admins to modify tax rules and financial parameters without separation of duties.
* **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 10: `/admin/profit` (Historical Profit & Revenue Analytics)
* **Frontend Route**: `/admin/profit`
* **Source File**: [`frontend/src/app/admin/profit/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/profit/page.tsx)
* **Shared Layout / Components**: `AdminLayout`
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`
* **Required Permission Identifiers**: `"profit"`
* **Backend API Endpoints**:
  * `GET /api/admin/profit/summary` (`RequirePermission("profit")`)
  * `GET /api/admin/profit/by-tier` (`RequirePermission("profit")`)
  * `GET /api/admin/profit/monthly` (`RequirePermission("profit")`)
* **Capabilities**: Historical listing fee revenue, broker commission totals, tier-by-tier revenue breakdowns
* **Auth Enforcement**: Backend enforces `RequirePermission("profit")`
* **UI States**: Spinner, summary KPI cards, tier table, monthly table
* **Localization Coverage**: 0%
* **Test Coverage**: Tested in `admin_authorization_matrix_test.go`.
* **Findings & Severity**:
  * **HIGH** (`INV-013`): Line 19 and line 35 block users with `if (!isSuperAdmin)`. In `layout.tsx`, Platform Admins and users with `"profit"` permission are given the sidebar link, but navigating here displays "Super admin access required".
* **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 11: `/admin/users` (Administrator Directory & Access Control)
* **Frontend Route**: `/admin/users`
* **Source File**: [`frontend/src/app/admin/users/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/users/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, `AdminRoleManagerModal`
* **Intended Roles**: `PlatformOwner` (Tier-0 Governance)
* **Required Permission Identifiers**: Platform Owner only
* **Backend API Endpoints**:
  * `GET /api/admin/admins` (`middleware.RequirePlatformOwner`)
  * `POST /api/admin/admins` (`middleware.RequirePlatformOwner`)
  * `PUT /api/admin/admins/:id/role` (`middleware.RequirePlatformOwner`)
  * `DELETE /api/admin/admins/:id` (`middleware.RequirePlatformOwner`)
  * `GET /api/admin/admins/:id/permissions` (`middleware.RequirePlatformOwner`)
  * `PUT /api/admin/admins/:id/permissions` (`middleware.RequirePlatformOwner`)
* **Capabilities**: List system administrators, view role badges, create admin account, update enterprise role, assign granular domain permissions, revoke admin access
* **Auth Enforcement**: Server-side enforced with `middleware.RequirePlatformOwner`
* **UI States**: Loading spinner, empty table, role modal, access revocation confirmation
* **Localization Coverage**: 0%
* **Test Coverage**: Backend covered in `admin_authorization_matrix_test.go`.
* **Findings & Severity**:
  * **HIGH** (`INV-014`): Line 33 checks `!isSuperAdmin` instead of `!isPlatformOwner`. If a Platform Owner logs in without legacy `is_super == true` in their JWT, they are redirected away to `/admin/dashboard`.
* **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 12: `/admin/settings` (Platform Settings, Security & Fees)
* **Frontend Route**: `/admin/settings`
* **Source File**: [`frontend/src/app/admin/settings/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/settings/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, `IPTrackingSecurity`, `FeeControlEngine`, `AdminPasskeyManager`
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`
* **Required Permission Identifiers**: `"settings"`
* **Backend API Endpoints**:
  * `GET /api/admin/settings` (`RequirePermission("settings")`)
  * `PUT /api/admin/settings` (`RequirePermission("settings")`)
  * `PUT /api/admin/settings/with-reason` (`middleware.RequirePlatformAdmin`)
  * `GET /api/admin/settings/log` (`middleware.RequirePlatformAdmin`)
* **Capabilities**: Maintenance mode toggle, OTP enforcement toggles, broker fee rates, referral limits, merchant phone numbers, passkey management, IP firewall rules
* **Auth Enforcement**: Server-side enforced by `RequirePermission("settings")`
* **UI States**: 8 tabbed category sections, save notification, error banner
* **Localization Coverage**: 0%
* **Test Coverage**: Settings read/update covered in `admin_authorization_matrix_test.go`.
* **Findings & Severity**:
  * **MEDIUM** (`INV-015`): Line 104 uses `api.adminUpdateSettings` (`PUT /api/admin/settings`), completely bypassing `PUT /api/admin/settings/with-reason`. Changes submitted through this page fail to create audit records in `settings_change_log`.
* **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 13: `/admin/logs` (Immutable Security Audit Trail)
* **Frontend Route**: `/admin/logs`
* **Source File**: [`frontend/src/app/admin/logs/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/logs/page.tsx)
* **Shared Layout / Components**: `AdminLayout`
* **Intended Roles**: `PlatformOwner`, `SecurityAdmin`
* **Required Permission Identifiers**: Security Admin or `"logs"`
* **Backend API Endpoints**:
  * `GET /api/admin/logs?page=...&limit=...` (`middleware.RequireSecurityAdmin`)
* **Capabilities**: Read-only pagination of administrative actions (actor, action, target, notes, timestamp)
* **Auth Enforcement**: Backend enforces `middleware.RequireSecurityAdmin`
* **UI States**: Pagination controls (50/page), empty state, error alert
* **Localization Coverage**: 0%
* **Test Coverage**: Tested in `admin_authorization_matrix_test.go`.
* **Findings & Severity**:
  * **MEDIUM** (`INV-016`): Line 17 allows users with `hasAdminPermission('logs')` to see the UI. However, backend line 688 uses `middleware.RequireSecurityAdmin`, which strictly checks role (`RoleSecurityAdmin` or `RolePlatformOwner`) without consulting individual permission grants. Admins granted the `"logs"` permission receive HTTP 403.
* **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 14: `/admin/verification` (Seller Identity & Document Verification)
* **Frontend Route**: `/admin/verification`
* **Source File**: [`frontend/src/app/admin/verification/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/verification/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, Document Review Modal
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
* **Required Permission Identifiers**: `"seller_verification"`
* **Backend API Endpoints**:
  * `GET /api/admin/verification/documents?status=...` (`RequirePermission("seller_verification")`)
  * `PUT /api/admin/verification/documents/:id/review` (`RequirePermission("seller_verification")`)
* **Capabilities**: Filter documents (`pending`, `approved`, `rejected`), inspect business licenses and national IDs, approve or reject with reason
* **Auth Enforcement**: Backend `RequirePermission("seller_verification")`
* **UI States**: Status tab bar, review modal, loading animation, empty list
* **Localization Coverage**: 0%
* **Test Coverage**: Untested on frontend.
* **Findings & Severity**: **LOW** (`INV-017`). No client-side permission barrier on page component.
* **Audit Status**: **VERIFIED**

---

### Page 15: `/admin/trust-safety` (Fraud & Abuse Incident Response)
* **Frontend Route**: `/admin/trust-safety`
* **Source File**: [`frontend/src/app/admin/trust-safety/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/trust-safety/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, Abuse Review Modal
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `SecurityAdmin`, `OperationsAdmin`, `Moderator`
* **Required Permission Identifiers**: `"trust_safety"`
* **Backend API Endpoints**:
  * `GET /api/admin/trust/reports?status=...` (`RequirePermission("trust_safety")`)
  * `PUT /api/admin/trust/reports/:id/review` (`RequirePermission("trust_safety")`)
  * `PUT /api/admin/trust/reports/:id/priority` (`RequirePermission("trust_safety")`)
* **Capabilities**: Filter reports, toggle report priority (`low`, `medium`, `high`, `critical`), review incident reports
* **Auth Enforcement**: Backend `RequirePermission("trust_safety")`
* **UI States**: Priority badges, review modal, status tabs
* **Localization Coverage**: 0%
* **Test Coverage**: Tested in `admin_authorization_matrix_test.go`.
* **Findings & Severity**:
  * **HIGH** (`INV-018`): **Illusion of Enforcement**: The review modal allows moderators to select `action = "banned"` or `action = "suspended"`. However, the backend endpoint `ReviewAbuseReport` only updates the textual note on the abuse report record—it does **NOT** ban the target user. Real user bans require calling `PUT /api/admin/users/:id/ban`, which is restricted to `PlatformOwner`. The moderator believes the account was banned when it remains active.
* **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 16: `/admin/refunds` (Dispute & Escrow Refund Management)
* **Frontend Route**: `/admin/refunds`
* **Source File**: [`frontend/src/app/admin/refunds/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/refunds/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, Refund Review Modal
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
* **Required Permission Identifiers**: `"refunds"`
* **Backend API Endpoints**:
  * `GET /api/admin/refunds?status=...` (`RequirePermission("refunds")`)
  * `PUT /api/admin/refunds/:id/review` (`RequirePermission("refunds")`)
  * `PUT /api/admin/refunds/:id/processed` (`RequirePermission("refunds")`)
* **Capabilities**: Review refund requests, approve, reject with reason, mark processed with bank payment reference
* **Auth Enforcement**: Backend `RequirePermission("refunds")`
* **UI States**: Tabs (`pending`, `approved`, `rejected`, `processed`), modal with payment reference input
* **Localization Coverage**: 0%
* **Test Coverage**: Tested in `admin_authorization_matrix_test.go`.
* **Findings & Severity**: **LOW** (`INV-019`). Hardcoded ETB currency string.
* **Audit Status**: **VERIFIED**

---

### Page 17: `/admin/appeals` (Account Appeals & Ban Reinstatement)
* **Frontend Route**: `/admin/appeals`
* **Source File**: [`frontend/src/app/admin/appeals/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/appeals/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, Appeal Decision Modal
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
* **Required Permission Identifiers**: `"appeals"`
* **Backend API Endpoints**:
  * `GET /api/admin/appeals?status=...` (`RequirePermission("appeals")`)
  * `PUT /api/admin/appeals/:id/review` (`RequirePermission("appeals")`)
  * `PUT /api/admin/appeals/:id/status` (`RequirePermission("appeals")`)
* **Capabilities**: Filter appeals (`pending`, `under_review`, `approved`, `rejected`), review appellant statement, approve reinstatement, reject appeal
* **Auth Enforcement**: Backend `RequirePermission("appeals")`
* **UI States**: Status tabs, appeal type badge (`ban`, `deactivation`, `listing_rejection`), review modal
* **Localization Coverage**: 0%
* **Test Coverage**: Tested in `admin_authorization_matrix_test.go`.
* **Findings & Severity**: **LOW** (`INV-020`). Clean implementation.
* **Audit Status**: **VERIFIED**

---

### Page 18: `/admin/referrals` (Referral Commissions & Withdrawal Payouts)
* **Frontend Route**: `/admin/referrals`
* **Source File**: [`frontend/src/app/admin/referrals/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/referrals/page.tsx)
* **Shared Layout / Components**: `AdminLayout`
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, Finance Admins with `"withdrawals"`
* **Required Permission Identifiers**: `"withdrawals"`
* **Backend API Endpoints**:
  * `GET /api/admin/referral/overview` (`middleware.RequirePlatformAdmin`)
  * `GET /api/admin/referral/withdrawals?status=...` (`RequirePermissionOrSuper("withdrawals")`)
  * `PUT /api/admin/referral/withdrawals/:id/approve` (`RequirePermissionOrSuper("withdrawals")`)
  * `PUT /api/admin/referral/withdrawals/:id/reject` (`RequirePermissionOrSuper("withdrawals")`)
  * `PUT /api/admin/referral/withdrawals/:id/paid` (`RequirePermissionOrSuper("withdrawals")`)
* **Capabilities**: Payout dashboard, approve withdrawal, reject with reason, mark paid with transaction reference
* **Auth Enforcement**: Backend checks `RequirePermissionOrSuper("withdrawals")` and `RequirePlatformAdmin`
* **UI States**: Overview KPI cards, withdrawal status tabs, action prompt dialogs
* **Localization Coverage**: 0%
* **Test Coverage**: Tested in `admin_authorization_matrix_test.go`.
* **Findings & Severity**:
  * **HIGH** (`INV-021`): **Orphaned Admin Page**: `/admin/referrals` is **completely missing from `AdminSidebar` navigation** in `layout.tsx`. Admins cannot discover or reach this page without manually typing the URL.
  * **HIGH** (`INV-022`): Line 18 and line 53 check `!isSuperAdmin`. Platform Administrators who have authority under `middleware.RequirePlatformAdmin` are blocked by the frontend.
* **Audit Status**: **INSPECTED — ISSUES FOUND (ORPHANED)**

---

### Page 19: `/admin/upgrades` (Tier Upgrade Approval Console)
* **Frontend Route**: `/admin/upgrades`
* **Source File**: [`frontend/src/app/admin/upgrades/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/upgrades/page.tsx)
* **Shared Layout / Components**: `AdminLayout`, `RequestTable`
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`, `SupportAgent`
* **Required Permission Identifiers**: `"seller_upgrades"`, `"buyer_upgrades"`
* **Backend API Endpoints**:
  * `GET /api/admin/upgrade-requests?status=...` (`RequirePermission("seller_upgrades")`)
  * `PUT /api/admin/upgrade-requests/:id/approve` (`RequirePermission("seller_upgrades")`)
  * `PUT /api/admin/upgrade-requests/:id/reject` (`RequirePermission("seller_upgrades")`)
  * `GET /api/admin/buyer-upgrade-requests?status=...` (`RequirePermission("buyer_upgrades")`)
  * `PUT /api/admin/buyer-upgrade-requests/:id/approve` (`RequirePermission("buyer_upgrades")`)
  * `PUT /api/admin/buyer-upgrade-requests/:id/reject` (`RequirePermission("buyer_upgrades")`)
* **Capabilities**: Dual-section approval console: approve/reject seller Gold/Premium upgrades and buyer upgrades
* **Auth Enforcement**: Server-side enforced on each endpoint
* **UI States**: Dual table sections, status tabs, tier badges (`gold`, `premium`)
* **Localization Coverage**: 0%
* **Test Coverage**: Tested in `admin_authorization_matrix_test.go`.
* **Findings & Severity**:
  * **HIGH** (`INV-023`): **Orphaned Admin Page**: `/admin/upgrades` is **omitted from the sidebar navigation** in `layout.tsx`.
* **Audit Status**: **INSPECTED — ISSUES FOUND (ORPHANED)**

---

### Page 20: `/admin/premium-users` (Gold & Premium Membership Directory)
* **Frontend Route**: `/admin/premium-users`
* **Source File**: [`frontend/src/app/admin/premium-users/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/premium-users/page.tsx)
* **Shared Layout / Components**: `AdminLayout`
* **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
* **Required Permission Identifiers**: `"premium_users"`
* **Backend API Endpoints**:
  * `GET /api/admin/premium-users?role=...&tier=...` (`RequirePermission("premium_users")`)
* **Capabilities**: Filter paying subscribers by role (`seller`, `buyer`) and tier (`gold`, `premium`), inspect total value sold
* **Auth Enforcement**: Backend `RequirePermission("premium_users")`
* **UI States**: Role/Tier filter pills, table view, empty state
* **Localization Coverage**: 0%
* **Test Coverage**: Tested in `admin_authorization_matrix_test.go`.
* **Findings & Severity**:
  * **HIGH** (`INV-024`): **Orphaned Admin Page**: `/admin/premium-users` is **omitted from the sidebar navigation** in `layout.tsx`.
* **Audit Status**: **INSPECTED — ISSUES FOUND (ORPHANED)**

---

### Page 21: `/admin/login` (Enterprise Portal Directory)
* **Frontend Route**: `/admin/login`
* **Source File**: [`frontend/src/app/admin/login/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/login/page.tsx)
* **Shared Layout / Components**: Standalone login directory page (bypasses `AdminSidebar`)
* **Intended Roles**: Public / Unauthenticated Administrative Identities
* **Required Permission Identifiers**: None
* **Backend API Endpoints**: None (Client directory linking to 4 dedicated portal endpoints)
* **Capabilities**: Navigates to Tier-0, Platform Admin, Security Admin, or Staff login pages; auto-redirects active admin sessions to `/admin/dashboard`
* **Auth Enforcement**: Public access point
* **UI States**: 4 responsive portal cards with security tier indicators
* **Localization Coverage**: 0%
* **Test Coverage**: Verified in `npm run build`.
* **Findings & Severity**: **LOW** (`INV-025`). Clean portal directory.
* **Audit Status**: **VERIFIED**

---

### Page 22: `/executive-console/login` (Platform Owner Tier-0 Governance Portal)
* **Frontend Route**: `/executive-console/login`
* **Source File**: [`frontend/src/app/executive-console/login/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/executive-console/login/page.tsx)
* **Shared Layout / Components**: `AdminPortalLogin` (`portal="executive"`, `portalTheme="crimson"`, `allowPasswordMfaFallback=false`)
* **Intended Roles**: `PlatformOwner` (Tier-0 Governance)
* **Required Permission Identifiers**: Unauthenticated ceremony
* **Backend API Endpoints**:
  * `POST /api/admin/passkeys/login/begin`
  * `POST /api/admin/passkeys/login/finish`
* **Capabilities**: Phishing-resistant WebAuthn biometric / hardware security key authentication. Password fallback strictly disabled (`allowPasswordMfaFallback={false}`).
* **Auth Enforcement**: Strict FIDO2 challenge-response verification via Go `webauthn` library on backend
* **UI States**: Hardware authentication prompt, error alert, spinner
* **Localization Coverage**: 0%
* **Test Coverage**: Backend covered in `admin_auth_service_test.go` and `admin_auth_handler_test.go`.
* **Findings & Severity**: **LOW** (`INV-026`). Secure passkey-only configuration.
* **Audit Status**: **VERIFIED**

---

### Page 23: `/platform-console/login` (Platform Administrator Console)
* **Frontend Route**: `/platform-console/login`
* **Source File**: [`frontend/src/app/platform-console/login/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/platform-console/login/page.tsx)
* **Shared Layout / Components**: `AdminPortalLogin` (`portal="platform"`, `portalTheme="indigo"`, `allowPasswordMfaFallback=false`)
* **Intended Roles**: `PlatformAdmin`
* **Required Permission Identifiers**: Unauthenticated ceremony
* **Backend API Endpoints**:
  * `POST /api/admin/passkeys/login/begin`
  * `POST /api/admin/passkeys/login/finish`
* **Capabilities**: Hardware passkey authentication for platform operations. Password fallback disabled.
* **Auth Enforcement**: Cryptographic passkey assertion verification on backend
* **UI States**: Passkey challenge prompt, biometric trigger
* **Localization Coverage**: 0%
* **Test Coverage**: Tested in `admin_auth_handler_test.go`.
* **Findings & Severity**: **LOW** (`INV-027`). Secure configuration.
* **Audit Status**: **VERIFIED**

---

### Page 24: `/security-console/login` (Security Administrator Console)
* **Frontend Route**: `/security-console/login`
* **Source File**: [`frontend/src/app/security-console/login/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/security-console/login/page.tsx)
* **Shared Layout / Components**: `AdminPortalLogin` (`portal="security"`, `portalTheme="emerald"`, `allowPasswordMfaFallback=false`)
* **Intended Roles**: `SecurityAdmin`
* **Required Permission Identifiers**: Unauthenticated ceremony
* **Backend API Endpoints**:
  * `POST /api/admin/passkeys/login/begin`
  * `POST /api/admin/passkeys/login/finish`
* **Capabilities**: Hardware passkey authentication for security monitoring and audit logging
* **Auth Enforcement**: Cryptographic passkey assertion verification on backend
* **UI States**: Passkey challenge prompt, error banner
* **Localization Coverage**: 0%
* **Test Coverage**: Tested in `admin_auth_handler_test.go`.
* **Findings & Severity**: **LOW** (`INV-028`). Secure configuration.
* **Audit Status**: **VERIFIED**

---

### Page 25: `/staff/login` (Operational Staff Access Portal)
* **Frontend Route**: `/staff/login`
* **Source File**: [`frontend/src/app/staff/login/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/staff/login/page.tsx)
* **Shared Layout / Components**: `AdminPortalLogin` (`portal="staff"`, `portalTheme="orange"`, `allowPasswordMfaFallback=true`)
* **Intended Roles**: `OperationsAdmin`, `Moderator`, `SupportAgent`
* **Required Permission Identifiers**: Unauthenticated ceremony
* **Backend API Endpoints**:
  * `POST /api/admin/passkeys/login/begin` & `finish` (Passkey preferred)
  * `POST /api/admin/login` (Step 1 Password fallback)
  * `POST /api/admin/login/verify-otp` (Step 2 Mandatory SMS/Email OTP verification)
* **Capabilities**: Passkey authentication with verified staff password + mandatory MFA fallback
* **Auth Enforcement**: Server-side MFA code validation with rate limiting
* **UI States**: Dual-mode auth selector (Passkey vs. Password+MFA), 6-digit OTP input step
* **Localization Coverage**: 0%
* **Test Coverage**: Tested in `admin_auth_handler_test.go`.
* **Findings & Severity**: **LOW** (`INV-029`). Secure multi-factor configuration.
* **Audit Status**: **VERIFIED**

---

## 2. Summary Status of All 25 Admin Pages

| Total Discovered | Fully Inspected | Issues Found | Verified Clean | Blocked |
| :---: | :---: | :---: | :---: | :---: |
| **25** | **25 (100%)** | **11** | **14** | **0** |

---
*End of Admin Page Inventory.*
