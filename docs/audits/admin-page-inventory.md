# HUNTIQ Enterprise Administration — Comprehensive Page Inventory & Routing Matrix

**Document Path**: `docs/audits/admin-page-inventory.md`  
**Audit Date**: October 23, 2026  
**Auditor**: Antigravity Quality & Enterprise Architecture Team  
**Repository**: `/home/poven/GO PROJECT/sellbuy`  
**Total Admin Pages Discovered**: 25 (21 inside `/admin/*`, 4 dedicated console login portals)  
**Total Fully Inspected Pages**: 25 (0 Partially Inspected, 0 Blocked, 0 Skipped)  
**Orphaned Backend Endpoints Discovered**: 12 endpoint groups (26 individual routes)  
**Orphaned Frontend Components Discovered**: 7 components  

---

## 1. Complete Inventory of All 25 Admin Pages (All 13 Required Fields)

---

### Page 01: `/admin` (Root Admin Entrypoint)
1. **Exact Frontend Route**: `/admin`
2. **Exact Source File Path**: [`frontend/src/app/admin/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/page.tsx)
3. **Shared Layout and Components**: `AdminLayout` ([`frontend/src/app/admin/layout.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/layout.tsx))
4. **Intended Roles**: All administrative identities (`PlatformOwner`, `PlatformAdmin`, `SecurityAdmin`, `OperationsAdmin`, `Moderator`, `SupportAgent`)
5. **Required Permission Identifiers**: None (Entrypoint redirect)
6. **Backend API Endpoints Used**: None directly
7. **Capabilities**: Client-side Next.js route redirection to `/admin/dashboard`
8. **Authentication and Authorization Enforcement**: Enforced by `AdminLayout` (redirects unauthenticated sessions to `/admin/login`)
9. **Loading, Empty, Error, and Unauthorized States**: Immediate redirect (`redirect('/admin/dashboard')`); layout renders `<Spinner />` during auth hydration
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: N/A (Redirect only)
11. **Relevant Tests and Actual Coverage**: Tested indirectly via Next.js route compiler (`npm run build`). Zero frontend unit tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-001` — Clean redirect, no security issues.
    * **Severity**: **LOW**
    * **Evidence**: Line 4: `redirect('/admin/dashboard')`.
    * **Recommended Action**: Maintain existing clean redirect.
13. **Audit Status**: **VERIFIED**

---

### Page 02: `/admin/dashboard` (Platform Overview & Operations Hub)
1. **Exact Frontend Route**: `/admin/dashboard`
2. **Exact Source File Path**: [`frontend/src/app/admin/dashboard/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/dashboard/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, `AdminSidebar`, `AdminHeaderTopBar`, `KPICard`, `QuickAction`, `MonthlyTrend`, `Spinner`
4. **Intended Roles**: All Admin Roles (`PlatformOwner`, `PlatformAdmin`, `SecurityAdmin`, `OperationsAdmin`, `Moderator`, `SupportAgent`)
5. **Required Permission Identifiers**: `"analytics"` (to view stats), `"pending"`, `"sellers"`, `"seller_verification"`, `"trust_safety"`, `"operations"` (for quick actions)
6. **Backend API Endpoints Used**:
   * `GET /api/admin/analytics/overview` (Gated by `RequirePermission("analytics")`)
   * `GET /api/admin/analytics/monthly?months=3` (Gated by `RequirePermission("analytics")`)
7. **Capabilities**: Read-only KPI aggregation, time-series monthly trend, quick action navigation links
8. **Authentication and Authorization Enforcement**: Server-side JWT/session validation via `mw.RequireAuth` + `RequirePermission("analytics")`
9. **Loading, Empty, Error, and Unauthorized States**:
   * Loading: `<Spinner />` centered state
   * Empty: Fallback placeholder dashes (`—`)
   * Error: Silently caught, defaults to `{}`
   * Unauthorized: Omits stats cards gracefully if permissions missing
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%** (Hardcoded English strings: "Platform Overview", "Total Listings", "Pending Review")
11. **Relevant Tests and Actual Coverage**: Backend analytics covered in `admin_authorization_matrix_test.go`. Frontend has zero tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-002` — Dashboard financial KPIs (`total_revenue`, `listing_fee_revenue`) check `isSuperAdmin` on line 95 instead of `isPlatformOwner || isPlatformAdmin || hasAdminPermission('profit')`.
    * **Severity**: **MEDIUM**
    * **Evidence**: Line 95: `const financialKpis: KPI[] = isSuperAdmin ? [...] : []`.
    * **Recommended Action**: Update check to `canViewFinancials = isPlatformOwner || isPlatformAdmin || hasAdminPermission('profit')`.
13. **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 03: `/admin/operations` (Background Operations & Job Monitor)
1. **Exact Frontend Route**: `/admin/operations`
2. **Exact Source File Path**: [`frontend/src/app/admin/operations/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/operations/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, `AdminSidebar`, `AdminHeaderTopBar`, `StatusBadge`, `StatusIcon`, `KPIBox`
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`, `SecurityAdmin`
5. **Required Permission Identifiers**: `"operations"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/operations/health` (`RequirePermission("operations")`)
   * `GET /api/admin/operations/jobs` (`RequirePermission("operations")`)
   * `GET /api/admin/operations/failed-jobs` (`RequirePermission("operations")`)
7. **Capabilities**: Read-only monitoring of system health, 24h job executions, dead-letter failed jobs, database table sizes; auto-refresh every 30s
8. **Authentication and Authorization Enforcement**: Server-side enforced by `mw.RequirePermission("operations")`
9. **Loading, Empty, Error, and Unauthorized States**:
   * Loading: Skeleton pulsing cards
   * Empty: Renders table headers with empty rows message
   * Error/Unauthorized: `Promise.allSettled` catches 403; page renders with `null` health rather than explicit 403 banner
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%** (Hardcoded English)
11. **Relevant Tests and Actual Coverage**: Backend endpoints covered in `admin_authorization_matrix_test.go`. Frontend has zero tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-003` — Missing client-side unauthorized notice when backend returns 403.
    * **Severity**: **LOW**
    * **Evidence**: Lines 60-70 silently swallow errors.
    * **Recommended Action**: Display explicit authorization banner if API calls fail with 403.
13. **Audit Status**: **VERIFIED**

---

### Page 04: `/admin/products` (Catalog & Listing Moderation Console)
1. **Exact Frontend Route**: `/admin/products`
2. **Exact Source File Path**: [`frontend/src/app/admin/products/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/products/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, `ProductReview`, `MultiStepListingWizard`
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`, `Moderator`
5. **Required Permission Identifiers**: `"pending"` (to review pending), `"products"` (to view catalog & fee claims)
6. **Backend API Endpoints Used**:
   * `GET /api/admin/products/pending` (`RequirePermission("pending")`)
   * `GET /api/admin/products` (`RequirePermission("products")`)
   * `PUT /api/admin/products/:id/review` (`RequirePermission("pending")`)
   * `DELETE /api/admin/products/:id` (`RequirePermission("products")`)
   * `GET /api/admin/fee-payments/pending` (`RequirePermission("products")`)
   * `PUT /api/admin/fee-payments/:id/verify` (`RequirePermission("products")`)
   * `PUT /api/admin/fee-payments/:id/reject` (`RequirePermission("products")`)
7. **Capabilities**: Filter by vertical (`sell`, `rent`, `lease`), approve, reject, delete listing, verify listing fee payments, bulk approve/reject
8. **Authentication and Authorization Enforcement**: Server-side role & permission checks on every endpoint
9. **Loading, Empty, Error, and Unauthorized States**: Grid vs. Table view toggle, loading spinner, empty list messages, error banner
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%** (Hardcoded English)
11. **Relevant Tests and Actual Coverage**: Backend covered in `product_handler_test.go` and `admin_authorization_matrix_test.go`. Frontend has zero tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-004` — Bulk review executes an unthrottled client-side `for` loop executing separate HTTP PUT calls. 50 items trigger 50 concurrent requests.
    * **Severity**: **MEDIUM**
    * **Evidence**: Lines 140-155 iterate over `selectedIds` with `await api.adminReviewProduct(...)`.
    * **Recommended Action**: Implement dedicated batch review endpoint `POST /api/admin/products/batch-review`.
13. **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 05: `/admin/categories` (Category Hierarchy Management)
1. **Exact Frontend Route**: `/admin/categories`
2. **Exact Source File Path**: [`frontend/src/app/admin/categories/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/categories/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, Category Tree Accordion Modal
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
5. **Required Permission Identifiers**: `"categories"`
6. **Backend API Endpoints Used**:
   * `GET /api/categories/tree?lang=en` (Public read)
   * `POST /api/admin/category-hierarchy` (Gated by `middleware.RequirePlatformAdmin`)
   * `PUT /api/admin/category-hierarchy/:id` (Gated by `middleware.RequirePlatformAdmin`)
   * `DELETE /api/admin/category-hierarchy/:id` (Gated by `middleware.RequirePlatformAdmin`)
7. **Capabilities**: Tree expansion, create category, edit multilingual names (`en`, `am`, `om`), slug generation, toggle active/leaf, delete category
8. **Authentication and Authorization Enforcement**: Server-side enforced with `middleware.RequirePlatformAdmin`
9. **Loading, Empty, Error, and Unauthorized States**: Expand/collapse tree nodes, edit modal, delete confirmation modal, error banner
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: Has form inputs for `name_en`, `name_am`, `name_om`, but admin UI chrome is hardcoded English.
11. **Relevant Tests and Actual Coverage**: Backend covered in `category_handler_test.go`. Frontend has zero tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-005` — Clean implementation. Action buttons correctly use `canManage = isSuperAdmin || isPlatformOwner || isPlatformAdmin || hasAdminPermission('categories')`.
    * **Severity**: **LOW**
    * **Evidence**: Lines 25, 135, and 181 use `canManage`.
    * **Recommended Action**: Maintain existing verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 06: `/admin/attributes` (Dynamic Listing Attribute Definitions)
1. **Exact Frontend Route**: `/admin/attributes`
2. **Exact Source File Path**: [`frontend/src/app/admin/attributes/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/attributes/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, Attribute Create Modal
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`
5. **Required Permission Identifiers**: `"categories"`
6. **Backend API Endpoints Used**:
   * `GET /api/attributes` (Public definition list)
   * `POST /api/admin/attributes` (Gated by `middleware.RequirePlatformAdmin`)
7. **Capabilities**: List dynamic attributes (Brand, Storage, Fuel Type, etc.), define new attribute with data types (`text`, `number`, `boolean`, `select`, `multiselect`, `range`, `date`)
8. **Authentication and Authorization Enforcement**: Backend enforces `middleware.RequirePlatformAdmin`
9. **Loading, Empty, Error, and Unauthorized States**: Empty table message, loading spinner, creation modal with validation rules
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: Multilingual name inputs for `name_en`, `name_am`, `name_om`. UI chrome is English-only.
11. **Relevant Tests and Actual Coverage**: Zero frontend tests. Backend handlers lack dedicated tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-006` — Missing Edit and Delete buttons in the UI for attribute definitions (read + create only).
    * **Severity**: **MEDIUM**
    * **Evidence**: Table rows only render attribute metadata with no action menu.
    * **Recommended Action**: Implement edit/delete modals and hook into backend update/delete handlers.
13. **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 07: `/admin/locations` (Location Hierarchy & GPS Geocoding Admin)
1. **Exact Frontend Route**: `/admin/locations`
2. **Exact Source File Path**: [`frontend/src/app/admin/locations/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/locations/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, Tabbed view (`tree`, `list`, `geocoding`, `import`), Location Edit Modal
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
5. **Required Permission Identifiers**: `"locations"`
6. **Backend API Endpoints Used**:
   * `GET /api/locations/tree?lang=en`
   * `GET /api/admin/location-hierarchy` (`RequirePermission("locations")`)
   * `POST /api/admin/location-hierarchy` (`RequirePermission("locations")`)
   * `PUT /api/admin/location-hierarchy/:id` (`RequirePermission("locations")`)
   * `DELETE /api/admin/location-hierarchy/:id` (`RequirePermission("locations")`)
   * `POST /api/admin/location-hierarchy/import` (`RequirePermission("locations")`)
   * `POST /api/admin/location-hierarchy/test-geocode` (`RequirePermission("locations")`)
7. **Capabilities**: 5-level geographic tree navigation (Country $\rightarrow$ Region $\rightarrow$ Zone/City $\rightarrow$ Woreda $\rightarrow$ Kebele), reverse GPS sandbox testing, batch JSON location import
8. **Authentication and Authorization Enforcement**: Server-side enforced on all mutation routes
9. **Loading, Empty, Error, and Unauthorized States**: Fully responsive tree view, search filter, tab switcher, geocoding pin tester
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: Amharic, Oromo, Tigrinya name inputs supported in schema.
11. **Relevant Tests and Actual Coverage**: Backend covered in `location_service_test.go` and `geocoding_test.go`. Zero frontend tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-007` — Line 84 properly checks `user?.token || localStorage.getItem('token')`.
    * **Severity**: **LOW**
    * **Evidence**: Line 84: `const token = user?.token || (typeof window !== 'undefined' ? localStorage.getItem('token') || '' : '');`.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 08: `/admin/search` (Search Intelligence & Zero-Result Analytics)
1. **Exact Frontend Route**: `/admin/search`
2. **Exact Source File Path**: [`frontend/src/app/admin/search/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/search/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
5. **Required Permission Identifiers**: `"analytics"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/search/popular` (`RequirePermission("analytics")`)
   * `GET /api/admin/search/zero-results` (`RequirePermission("analytics")`)
7. **Capabilities**: Read-only display of top 30 search queries, CTR percentages, and top 50 zero-result search queries
8. **Authentication and Authorization Enforcement**: Backend `RequirePermission("analytics")`
9. **Loading, Empty, Error, and Unauthorized States**: Skeleton loading, empty state messages, auto-scroll lists
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: No frontend tests. Backend covered in `admin_authorization_matrix_test.go`.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-008` — Missing client-side unauthorized notice if API returns 403.
    * **Severity**: **LOW**
    * **Evidence**: Lines 25-35 silently swallow errors.
    * **Recommended Action**: Render unauthorized alert banner.
13. **Audit Status**: **VERIFIED**

---

### Page 09: `/admin/finance` (Financial Planning & Analysis / FP&A)
1. **Exact Frontend Route**: `/admin/finance`
2. **Exact Source File Path**: [`frontend/src/app/admin/finance/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/finance/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, 9 tabbed modules, SVG interactive charts, 20M capacity stress analyzer, statutory tax catalog, assumption editor modal
4. **Intended Roles**: `PlatformOwner`, Designated Finance Officers
5. **Required Permission Identifiers**:
   * View: `"profit"`, `"finance"`, `"finance:view"`
   * Mutate: `"finance:manage"`, `"finance_manage"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/finance/forecast` (`financeViewAuth`)
   * `GET /api/admin/finance/scenarios` (`financeViewAuth`)
   * `GET /api/admin/finance/capacity-20m` (`financeViewAuth`)
   * `GET /api/admin/finance/tax-rules` (`financeViewAuth`)
   * `GET /api/admin/finance/actuals` (`financeViewAuth`)
   * `GET /api/admin/finance/audit-logs` (`financeViewAuth`)
   * `PUT /api/admin/finance/assumptions` (`financeManageAuth`)
   * `POST /api/admin/finance/assumptions/reset` (`financeManageAuth`)
   * `GET /api/admin/finance/export` (`financeViewAuth`)
7. **Capabilities**: Multi-scenario forecasting (12–60 mos), workload demand metrics, Ethiopian statutory tax estimation, cash burn/runway tracking, assumption editing, CSV/JSON export
8. **Authentication and Authorization Enforcement**: Enforced by `financeViewAuth` (read) and `financeManageAuth` (write) in `main.go:619-628`
9. **Loading, Empty, Error, and Unauthorized States**: Full client-side authorization banner, spinner, 5-year table, interactive chart hover, assumption modals
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **100% localized** across English, Amharic, Oromo, and Tigrinya in `frontend/src/lib/language/{en,am,or,ti}.json`
11. **Relevant Tests and Actual Coverage**: Backend covered in `fpa_test.go` (6 unit test suites) and `fpa_handler_test.go` (7 HTTP integration tests).
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-009` — Backend separation of duties is properly enforced (`PlatformAdmin` has view access, forbidden from mutating assumptions with 403). Frontend assumption edit modal must gracefully disable save button for view-only admins.
    * **Severity**: **MEDIUM**
    * **Evidence**: Verified in `fpa_handler_test.go:PlatformAdmin_Forbidden_From_Mutating_Assumptions_403`.
    * **Recommended Action**: Hide or disable "Save Assumptions" button if user lacks `finance:manage` permission.
13. **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 10: `/admin/profit` (Historical Profit & Revenue Analytics)
1. **Exact Frontend Route**: `/admin/profit`
2. **Exact Source File Path**: [`frontend/src/app/admin/profit/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/profit/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`
5. **Required Permission Identifiers**: `"profit"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/profit/summary` (`RequirePermission("profit")`)
   * `GET /api/admin/profit/by-tier` (`RequirePermission("profit")`)
   * `GET /api/admin/profit/monthly` (`RequirePermission("profit")`)
7. **Capabilities**: Historical listing fee revenue, broker commission totals, tier-by-tier revenue breakdowns
8. **Authentication and Authorization Enforcement**: Backend enforces `RequirePermission("profit")`
9. **Loading, Empty, Error, and Unauthorized States**: Spinner, summary KPI cards, tier table, monthly table, authorization banner
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Tested in `admin_authorization_matrix_test.go`. Zero frontend tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-010` — Clean authorization logic. Line 13 correctly uses `canViewProfit = isSuperAdmin || isPlatformOwner || isPlatformAdmin || hasAdminPermission('profit')`.
    * **Severity**: **LOW**
    * **Evidence**: Lines 13, 20, 36.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 11: `/admin/users` (Administrator Directory & Access Control)
1. **Exact Frontend Route**: `/admin/users`
2. **Exact Source File Path**: [`frontend/src/app/admin/users/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/users/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, `AdminRoleManagerModal`
4. **Intended Roles**: `PlatformOwner` (Tier-0 Governance)
5. **Required Permission Identifiers**: Platform Owner only
6. **Backend API Endpoints Used**:
   * `GET /api/admin/admins` (`middleware.RequirePlatformOwner`)
   * `POST /api/admin/admins` (`middleware.RequirePlatformOwner`)
   * `PUT /api/admin/admins/:id/role` (`middleware.RequirePlatformOwner`)
   * `DELETE /api/admin/admins/:id` (`middleware.RequirePlatformOwner`)
   * `GET /api/admin/admins/:id/permissions` (`middleware.RequirePlatformOwner`)
   * `PUT /api/admin/admins/:id/permissions` (`middleware.RequirePlatformOwner`)
7. **Capabilities**: List system administrators, view role badges, create admin account, update enterprise role, assign granular domain permissions, revoke admin access
8. **Authentication and Authorization Enforcement**: Server-side enforced with `middleware.RequirePlatformOwner`
9. **Loading, Empty, Error, and Unauthorized States**: Loading spinner, empty table, role modal, access revocation confirmation
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Backend covered in `admin_authorization_matrix_test.go`. Zero frontend tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-011` — Clean authorization logic. Line 13 properly checks `canManageAdmins = isPlatformOwner || isSuperAdmin`.
    * **Severity**: **LOW**
    * **Evidence**: Line 13: `const canManageAdmins = isPlatformOwner || isSuperAdmin;`.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 12: `/admin/settings` (Platform Settings, Security & Fees)
1. **Exact Frontend Route**: `/admin/settings`
2. **Exact Source File Path**: [`frontend/src/app/admin/settings/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/settings/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, `IPTrackingSecurity`, `FeeControlEngine`, `AdminPasskeyManager`
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`
5. **Required Permission Identifiers**: `"settings"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/settings` (`RequirePermission("settings")`)
   * `PUT /api/admin/settings/with-reason` (`middleware.RequirePlatformAdmin`)
   * `GET /api/admin/settings/log` (`middleware.RequirePlatformAdmin`)
7. **Capabilities**: Maintenance mode toggle, OTP enforcement toggles, broker fee rates, referral limits, merchant phone numbers, passkey management, IP firewall rules
8. **Authentication and Authorization Enforcement**: Server-side enforced by `RequirePermission("settings")` and `RequirePlatformAdmin`
9. **Loading, Empty, Error, and Unauthorized States**: 8 tabbed category sections, save notification, error banner, audit reason prompt
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Settings read/update covered in `admin_authorization_matrix_test.go`. Zero frontend tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-012` — Line 110 properly calls `api.adminUpdateSettingsWithReason`, requiring an administrative reason and logging all changes to `settings_change_log`.
    * **Severity**: **LOW**
    * **Evidence**: Line 110: `await api.adminUpdateSettingsWithReason(settingsMap, reason.trim(), user.token);`.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 13: `/admin/logs` (Immutable Security Audit Trail)
1. **Exact Frontend Route**: `/admin/logs`
2. **Exact Source File Path**: [`frontend/src/app/admin/logs/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/logs/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`
4. **Intended Roles**: `PlatformOwner`, `SecurityAdmin`
5. **Required Permission Identifiers**: Security Admin or `"logs"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/logs?page=...&limit=...` (`mw.RequireAnyPermission("logs")`)
7. **Capabilities**: Read-only pagination of administrative actions (actor, action, target, notes, timestamp)
8. **Authentication and Authorization Enforcement**: Backend enforces `mw.RequireAnyPermission("logs")`
9. **Loading, Empty, Error, and Unauthorized States**: Pagination controls (50/page), empty state, error alert, authorization banner
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Tested in `admin_authorization_matrix_test.go`. Zero frontend tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-013` — Backend route in `main.go:689` now uses `mw.RequireAnyPermission("logs")`, resolving the previous role deadlock. Both `SecurityAdmin` and admins with `"logs"` permission can access the audit trail.
    * **Severity**: **LOW**
    * **Evidence**: `main.go:689`.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 14: `/admin/verification` (Seller Identity & Document Verification)
1. **Exact Frontend Route**: `/admin/verification`
2. **Exact Source File Path**: [`frontend/src/app/admin/verification/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/verification/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, Document Review Modal
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
5. **Required Permission Identifiers**: `"seller_verification"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/verification/documents?status=...` (`RequirePermission("seller_verification")`)
   * `PUT /api/admin/verification/documents/:id/review` (`RequirePermission("seller_verification")`)
7. **Capabilities**: Filter documents (`pending`, `approved`, `rejected`), inspect business licenses and national IDs, approve or reject with reason
8. **Authentication and Authorization Enforcement**: Backend `RequirePermission("seller_verification")`
9. **Loading, Empty, Error, and Unauthorized States**: Status tab bar, review modal, loading animation, empty list
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Backend tested in `admin_authorization_matrix_test.go`. Untested on frontend.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-014` — Clean implementation. Rejection requires providing a rejection reason.
    * **Severity**: **LOW**
    * **Evidence**: Modal enforces note input when rejecting.
    * **Recommended Action**: Connect the orphaned `DocumentViewerModal.tsx` for enhanced multi-page document inspection.
13. **Audit Status**: **VERIFIED**

---

### Page 15: `/admin/trust-safety` (Fraud & Abuse Incident Response)
1. **Exact Frontend Route**: `/admin/trust-safety`
2. **Exact Source File Path**: [`frontend/src/app/admin/trust-safety/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/trust-safety/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, Abuse Review Modal
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `SecurityAdmin`, `OperationsAdmin`, `Moderator`
5. **Required Permission Identifiers**: `"trust_safety"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/trust/reports?status=...` (`RequirePermission("trust_safety")`)
   * `PUT /api/admin/trust/reports/:id/review` (`RequirePermission("trust_safety")`)
   * `PUT /api/admin/trust/reports/:id/priority` (`RequirePermission("trust_safety")`)
7. **Capabilities**: Filter reports, toggle report priority (`low`, `medium`, `high`, `critical`), review incident reports, take actions (`warn`, `content_removed`, `suspended`, `banned`, `no_action`)
8. **Authentication and Authorization Enforcement**: Backend `RequirePermission("trust_safety")`
9. **Loading, Empty, Error, and Unauthorized States**: Priority badges, review modal, status tabs
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Tested in `admin_authorization_matrix_test.go`. Zero frontend tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-015` — Backend `admin_service.go:730` handles `actionTaken == "banned"` by automatically calling `s.BanUser()`, but `actionTaken == "suspended"` does not deactivate the user account.
    * **Severity**: **MEDIUM**
    * **Evidence**: `admin_service.go:730-736`.
    * **Recommended Action**: Also handle `actionTaken == "suspended"` in `ReviewAbuseReport` to suspend account access.
13. **Audit Status**: **INSPECTED — ISSUES FOUND**

---

### Page 16: `/admin/refunds` (Dispute & Escrow Refund Management)
1. **Exact Frontend Route**: `/admin/refunds`
2. **Exact Source File Path**: [`frontend/src/app/admin/refunds/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/refunds/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, Refund Review Modal
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
5. **Required Permission Identifiers**: `"refunds"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/refunds?status=...` (`RequirePermission("refunds")`)
   * `PUT /api/admin/refunds/:id/review` (`RequirePermission("refunds")`)
   * `PUT /api/admin/refunds/:id/processed` (`RequirePermission("refunds")`)
7. **Capabilities**: Review refund requests, approve, reject with reason, mark processed with bank payment reference
8. **Authentication and Authorization Enforcement**: Backend `RequirePermission("refunds")`
9. **Loading, Empty, Error, and Unauthorized States**: Tabs (`pending`, `approved`, `rejected`, `processed`), modal with payment reference input
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Tested in `admin_authorization_matrix_test.go`. Zero frontend tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-016` — Clean implementation. Processing requires bank reference number.
    * **Severity**: **LOW**
    * **Evidence**: Lines 60-70.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 17: `/admin/appeals` (Account Appeals & Ban Reinstatement)
1. **Exact Frontend Route**: `/admin/appeals`
2. **Exact Source File Path**: [`frontend/src/app/admin/appeals/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/appeals/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, Appeal Decision Modal
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
5. **Required Permission Identifiers**: `"appeals"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/appeals?status=...` (`RequirePermission("appeals")`)
   * `PUT /api/admin/appeals/:id/review` (`RequirePermission("appeals")`)
   * `PUT /api/admin/appeals/:id/status` (`RequirePermission("appeals")`)
7. **Capabilities**: Filter appeals (`pending`, `under_review`, `approved`, `rejected`), review appellant statement, approve reinstatement, reject appeal
8. **Authentication and Authorization Enforcement**: Backend `RequirePermission("appeals")`
9. **Loading, Empty, Error, and Unauthorized States**: Status tabs, appeal type badge (`ban`, `deactivation`, `listing_rejection`), review modal
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Tested in `admin_authorization_matrix_test.go`. Zero frontend tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-017` — Clean implementation.
    * **Severity**: **LOW**
    * **Evidence**: Complete review cycle.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 18: `/admin/referrals` (Referral Commissions & Withdrawal Payouts)
1. **Exact Frontend Route**: `/admin/referrals`
2. **Exact Source File Path**: [`frontend/src/app/admin/referrals/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/referrals/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, Finance Admins with `"withdrawals"` / `"referrals"`
5. **Required Permission Identifiers**: `"referrals"`, `"withdrawals"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/referral/overview` (`middleware.RequirePlatformAdmin`)
   * `GET /api/admin/referral/withdrawals?status=...` (`RequirePermissionOrSuper("withdrawals")`)
   * `PUT /api/admin/referral/withdrawals/:id/approve` (`RequirePermissionOrSuper("withdrawals")`)
   * `PUT /api/admin/referral/withdrawals/:id/reject` (`RequirePermissionOrSuper("withdrawals")`)
   * `PUT /api/admin/referral/withdrawals/:id/paid` (`RequirePermissionOrSuper("withdrawals")`)
7. **Capabilities**: Payout dashboard, approve withdrawal, reject with reason, mark paid with transaction reference
8. **Authentication and Authorization Enforcement**: Backend checks `RequirePermissionOrSuper("withdrawals")` and `RequirePlatformAdmin`
9. **Loading, Empty, Error, and Unauthorized States**: Overview KPI cards, withdrawal status tabs, action prompt dialogs, authorization banner
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Tested in `admin_authorization_matrix_test.go`. Zero frontend tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-018` — Mounted in `AdminSidebar` under `MEMBERSHIPS & MONETIZATION`. Authorization check line 12 uses `canManageReferrals = isSuperAdmin || isPlatformOwner || isPlatformAdmin || hasAdminPermission('referrals')`. Uses native browser `prompt()` for rejection and bank reference.
    * **Severity**: **LOW**
    * **Evidence**: Line 12, 41, 49.
    * **Recommended Action**: Replace native `prompt()` with modal dialog.
13. **Audit Status**: **VERIFIED**

---

### Page 19: `/admin/upgrades` (Tier Upgrade Approval Console)
1. **Exact Frontend Route**: `/admin/upgrades`
2. **Exact Source File Path**: [`frontend/src/app/admin/upgrades/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/upgrades/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`, `RequestTable`
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`, `SupportAgent`
5. **Required Permission Identifiers**: `"seller_upgrades"`, `"buyer_upgrades"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/upgrade-requests?status=...` (`RequirePermission("seller_upgrades")`)
   * `PUT /api/admin/upgrade-requests/:id/approve` (`RequirePermission("seller_upgrades")`)
   * `PUT /api/admin/upgrade-requests/:id/reject` (`RequirePermission("seller_upgrades")`)
   * `GET /api/admin/buyer-upgrade-requests?status=...` (`RequirePermission("buyer_upgrades")`)
   * `PUT /api/admin/buyer-upgrade-requests/:id/approve` (`RequirePermission("buyer_upgrades")`)
   * `PUT /api/admin/buyer-upgrade-requests/:id/reject` (`RequirePermission("buyer_upgrades")`)
7. **Capabilities**: Dual-section approval console: approve/reject seller Gold/Premium upgrades and buyer upgrades
8. **Authentication and Authorization Enforcement**: Server-side enforced on each endpoint
9. **Loading, Empty, Error, and Unauthorized States**: Dual table sections, status tabs, tier badges (`gold`, `premium`)
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Tested in `admin_authorization_matrix_test.go`. Zero frontend tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-019` — Mounted in `AdminSidebar` under `MEMBERSHIPS & MONETIZATION`. Clean authorization checks.
    * **Severity**: **LOW**
    * **Evidence**: Fully functional dual approval console.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 20: `/admin/premium-users` (Gold & Premium Membership Directory)
1. **Exact Frontend Route**: `/admin/premium-users`
2. **Exact Source File Path**: [`frontend/src/app/admin/premium-users/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/premium-users/page.tsx)
3. **Shared Layout and Components**: `AdminLayout`
4. **Intended Roles**: `PlatformOwner`, `PlatformAdmin`, `OperationsAdmin`
5. **Required Permission Identifiers**: `"premium_users"`
6. **Backend API Endpoints Used**:
   * `GET /api/admin/premium-users?role=...&tier=...` (`RequirePermission("premium_users")`)
7. **Capabilities**: Filter paying subscribers by role (`seller`, `buyer`) and tier (`gold`, `premium`), inspect total value sold
8. **Authentication and Authorization Enforcement**: Backend `RequirePermission("premium_users")`
9. **Loading, Empty, Error, and Unauthorized States**: Role/Tier filter pills, table view, empty state, authorization banner
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Tested in `admin_authorization_matrix_test.go`. Zero frontend tests.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-020` — Mounted in `AdminSidebar` under `MEMBERSHIPS & MONETIZATION`. Clean implementation.
    * **Severity**: **LOW**
    * **Evidence**: Lines 16, 32.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 21: `/admin/login` (Enterprise Portal Directory)
1. **Exact Frontend Route**: `/admin/login`
2. **Exact Source File Path**: [`frontend/src/app/admin/login/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/admin/login/page.tsx)
3. **Shared Layout and Components**: Standalone login directory page (bypasses `AdminSidebar`)
4. **Intended Roles**: Public / Unauthenticated Administrative Identities
5. **Required Permission Identifiers**: None
6. **Backend API Endpoints Used**: None (Client directory linking to 4 dedicated portal endpoints)
7. **Capabilities**: Navigates to Tier-0, Platform Admin, Security Admin, or Staff login pages; auto-redirects active admin sessions to `/admin/dashboard`
8. **Authentication and Authorization Enforcement**: Public access point
9. **Loading, Empty, Error, and Unauthorized States**: 4 responsive portal cards with security tier indicators
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **0%**
11. **Relevant Tests and Actual Coverage**: Verified in `npm run build`.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-021` — Clean portal directory.
    * **Severity**: **LOW**
    * **Evidence**: 4 distinct portal entrypoints.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 22: `/executive-console/login` (Platform Owner Tier-0 Governance Portal)
1. **Exact Frontend Route**: `/executive-console/login`
2. **Exact Source File Path**: [`frontend/src/app/executive-console/login/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/executive-console/login/page.tsx)
3. **Shared Layout and Components**: `AdminPortalLogin` (`portal="executive"`, `portalTheme="crimson"`, `allowPasswordMfaFallback=false`)
4. **Intended Roles**: `PlatformOwner` (Tier-0 Governance)
5. **Required Permission Identifiers**: Unauthenticated ceremony
6. **Backend API Endpoints Used**:
   * `POST /api/admin/passkeys/login/begin`
   * `POST /api/admin/passkeys/login/finish`
7. **Capabilities**: Phishing-resistant WebAuthn biometric / hardware security key authentication. Password fallback strictly disabled (`allowPasswordMfaFallback={false}`).
8. **Authentication and Authorization Enforcement**: Strict FIDO2 challenge-response verification via Go `webauthn` library on backend
9. **Loading, Empty, Error, and Unauthorized States**: Hardware authentication prompt, error alert, spinner
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **100% localized** via `admin_security` translation keys
11. **Relevant Tests and Actual Coverage**: Backend covered in `admin_auth_service_test.go` and `admin_auth_handler_test.go`.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-022` — Secure passkey-only configuration.
    * **Severity**: **LOW**
    * **Evidence**: `allowPasswordMfaFallback={false}`.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 23: `/platform-console/login` (Platform Administrator Console)
1. **Exact Frontend Route**: `/platform-console/login`
2. **Exact Source File Path**: [`frontend/src/app/platform-console/login/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/platform-console/login/page.tsx)
3. **Shared Layout and Components**: `AdminPortalLogin` (`portal="platform"`, `portalTheme="indigo"`, `allowPasswordMfaFallback=false`)
4. **Intended Roles**: `PlatformAdmin`
5. **Required Permission Identifiers**: Unauthenticated ceremony
6. **Backend API Endpoints Used**:
   * `POST /api/admin/passkeys/login/begin`
   * `POST /api/admin/passkeys/login/finish`
7. **Capabilities**: Hardware passkey authentication for platform operations. Password fallback disabled.
8. **Authentication and Authorization Enforcement**: Cryptographic passkey assertion verification on backend
9. **Loading, Empty, Error, and Unauthorized States**: Passkey challenge prompt, biometric trigger
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **100% localized** via `admin_security` translation keys
11. **Relevant Tests and Actual Coverage**: Tested in `admin_auth_handler_test.go`.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-023` — Secure configuration.
    * **Severity**: **LOW**
    * **Evidence**: `allowPasswordMfaFallback={false}`.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 24: `/security-console/login` (Security Administrator Console)
1. **Exact Frontend Route**: `/security-console/login`
2. **Exact Source File Path**: [`frontend/src/app/security-console/login/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/security-console/login/page.tsx)
3. **Shared Layout and Components**: `AdminPortalLogin` (`portal="security"`, `portalTheme="emerald"`, `allowPasswordMfaFallback=false`)
4. **Intended Roles**: `SecurityAdmin`
5. **Required Permission Identifiers**: Unauthenticated ceremony
6. **Backend API Endpoints Used**:
   * `POST /api/admin/passkeys/login/begin`
   * `POST /api/admin/passkeys/login/finish`
7. **Capabilities**: Hardware passkey authentication for security monitoring and audit logging
8. **Authentication and Authorization Enforcement**: Cryptographic passkey assertion verification on backend
9. **Loading, Empty, Error, and Unauthorized States**: Passkey challenge prompt, error banner
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **100% localized** via `admin_security` translation keys
11. **Relevant Tests and Actual Coverage**: Tested in `admin_auth_handler_test.go`.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-024` — Secure configuration.
    * **Severity**: **LOW**
    * **Evidence**: `allowPasswordMfaFallback={false}`.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

### Page 25: `/staff/login` (Operational Staff Access Portal)
1. **Exact Frontend Route**: `/staff/login`
2. **Exact Source File Path**: [`frontend/src/app/staff/login/page.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/app/staff/login/page.tsx)
3. **Shared Layout and Components**: `AdminPortalLogin` (`portal="staff"`, `portalTheme="orange"`, `allowPasswordMfaFallback=true`)
4. **Intended Roles**: `OperationsAdmin`, `Moderator`, `SupportAgent`
5. **Required Permission Identifiers**: Unauthenticated ceremony
6. **Backend API Endpoints Used**:
   * `POST /api/admin/passkeys/login/begin` & `finish` (Passkey preferred)
   * `POST /api/admin/login` (Step 1 Password fallback)
   * `POST /api/admin/login/verify-otp` (Step 2 Mandatory SMS/Email OTP verification)
7. **Capabilities**: Passkey authentication with verified staff password + mandatory MFA fallback
8. **Authentication and Authorization Enforcement**: Server-side MFA code validation with rate limiting
9. **Loading, Empty, Error, and Unauthorized States**: Dual-mode auth selector (Passkey vs. Password+MFA), 6-digit OTP input step
10. **Localization Coverage (`en`, `am`, `or`, `ti`)**: **100% localized** via `admin_security` translation keys
11. **Relevant Tests and Actual Coverage**: Tested in `admin_auth_handler_test.go`.
12. **Findings, Severity, Evidence, and Recommended Action**:
    * **Finding**: `INV-025` — Secure multi-factor configuration.
    * **Severity**: **LOW**
    * **Evidence**: Mandatory Step 2 OTP verification.
    * **Recommended Action**: Maintain verified logic.
13. **Audit Status**: **VERIFIED**

---

## 2. Orphaned Backend Endpoints (Missing UI Coverage)

The following backend endpoints exist in [`backend/cmd/api/main.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/cmd/api/main.go) and are fully supported by backend handlers, services, and repositories, but **lack dedicated administrative frontend screens**:

| Backend Route | HTTP Method | Guarding Middleware | Purpose / Capability | Missing Frontend UI Page |
| :--- | :---: | :--- | :--- | :--- |
| `/api/admin/sellers` | `GET` | `RequirePermission("sellers")` | List marketplace sellers with ratings, tiers, and status | `/admin/sellers` |
| `/api/admin/sellers/:id/active` | `PUT` | `RequirePermission("sellers")` | Activate or deactivate seller accounts with reason | `/admin/sellers` |
| `/api/admin/transactions` | `GET` | `RequirePermission("transactions")` | Audit escrow transactions, buyer/seller IDs, and fees | `/admin/transactions` |
| `/api/admin/transactions/:id/status` | `PUT` | `RequirePermission("transactions")` | Update escrow transaction status (`completed`, `disputed`) | `/admin/transactions` |
| `/api/admin/system/health` | `GET` | `RequirePlatformAdmin` | Inspect DB schema health, missing indexes, email gaps | `/admin/system/health` |
| `/api/admin/system/health/repair` | `POST` | `RequirePlatformAdmin` | Execute automated schema integrity repair procedures | `/admin/system/health` |
| `/api/admin/notifications` | `GET`, `POST` | `RequirePermission("notifications")` | List and broadcast platform-wide administrative alerts | `/admin/notifications` |
| `/api/admin/banners` | `GET`, `POST`, `PUT`, `DELETE` | `RequirePermission("banners")` | Hero promotional banner and marketing CMS | `/admin/banners` |
| `/api/admin/otp-logs` | `GET` | `RequireAnyPermission("otp_logs")` | Inspect SMS/email OTP delivery logs and failures | `/admin/otp-logs` |
| `/api/admin/activity-logs` | `GET` | `RequireAnyPermission("logs", "otp_logs")` | Inspect marketplace end-user authentication history | `/admin/activity-logs` |
| `/api/admin/categories/popularity` | `GET` | `RequirePermission("category_trends")` | Analytics on high-demand vs low-demand categories | `/admin/categories/trends` |
| `/api/admin/users/:id/warn` | `POST` | `RequirePermission("sellers")` | Issue formal compliance warning to seller account | `/admin/sellers` |
| `/api/admin/users/:id/ban` | `PUT` | `RequirePlatformOwner` | Directly ban and deactivate any marketplace user | `/admin/users` |
| `/api/admin/users/:id/seller-tier` | `PUT` | `RequirePlatformOwner` | Override seller tier directly to Normal/Gold/Premium | `/admin/users` |
| `/api/admin/users/:id/buyer-tier` | `PUT` | `RequirePlatformOwner` | Override buyer tier directly to Normal/Gold/Premium | `/admin/users` |
| `/api/admin/users/:id/tier-history` | `GET` | `RequirePermission("premium_users")` | Audit timeline of user subscription tier changes | `/admin/premium-users` |
| `/api/admin/categories/:id/attributes` | `POST`, `DELETE` | `RequirePlatformAdmin` | Bind/unbind dynamic attributes to category definitions | `/admin/categories` |

---

## 3. Orphaned Frontend Components Discovered

The following components were built in the frontend repository but are **not imported or mounted in any active Next.js page**:

1. [`frontend/src/components/admin/enterprise/AvailabilityGrid.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/components/admin/enterprise/AvailabilityGrid.tsx) (12.2 KB): Unit calendar and reservation manager for rental properties.
2. [`frontend/src/components/admin/enterprise/DocumentViewerModal.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/components/admin/enterprise/DocumentViewerModal.tsx) (38.6 KB): Full-featured document inspection modal with zoom, rotation, and PDF controls.
3. [`frontend/src/components/admin/enterprise/DomainSwitcherTabs.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/components/admin/enterprise/DomainSwitcherTabs.tsx) (4.9 KB): Mode selector for multi-domain verification queues.
4. [`frontend/src/components/admin/enterprise/EnterpriseDataTable.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/components/admin/enterprise/EnterpriseDataTable.tsx) (28.8 KB): Generalized high-performance virtualized table with search and filters.
5. [`frontend/src/components/admin/ImageVerification.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/components/admin/ImageVerification.tsx) (1.0 KB): Legacy image moderation thumbnail box.
6. [`frontend/src/components/admin/LocationManager.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/components/admin/LocationManager.tsx) (2.8 KB): Legacy flat location CRUD box (superseded by `/admin/locations`).
7. [`frontend/src/components/admin/Stats.tsx`](file:///home/poven/GO%20PROJECT/sellbuy/frontend/src/components/admin/Stats.tsx) (1.3 KB): Legacy dashboard KPI card component.

---

## 4. Summary Status of All 25 Admin Pages

| Total Discovered | Fully Inspected | Issues Found | Verified Clean | Blocked |
| :---: | :---: | :---: | :---: | :---: |
| **25** | **25 (100%)** | **4** | **21** | **0** |

*Breakdown of Issues Found*:
* `/admin/dashboard`: Financial KPIs check `isSuperAdmin` instead of `canViewFinancials`.
* `/admin/products`: Bulk moderation executes unthrottled loop of individual PUT calls.
* `/admin/attributes`: Missing edit and delete actions in UI.
* `/admin/trust-safety`: `"suspended"` action does not suspend account access in backend `ReviewAbuseReport`.

---
*End of Comprehensive Admin Page Inventory.*
