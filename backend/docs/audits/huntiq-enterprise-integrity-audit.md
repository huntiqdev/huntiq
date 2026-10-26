# HUNTIQ — Ultra-Deep Enterprise Audit: Taxonomy, Database Configuration, Frontend, Backend, and Production Integrity

**Document:** `backend/docs/audits/huntiq-enterprise-integrity-audit.md`  
**Auditor:** Principal Software Architect, Senior Go/PostgreSQL Engineer, Next.js/TypeScript Specialist, Security & QA Lead  
**Audit Date:** 2026-10-24  
**Target Environment:** HUNTIQ Universal Marketplace (Go + PostgreSQL 18.4 runtime + Next.js 15.1.7 + Python 3.14 / FastAPI)

---

## 1. Executive Summary

A comprehensive, evidence-driven audit was performed across all three architectural pillars of the HUNTIQ repository:
- **Pillar 1:** Go Backend Authority (`backend/`) and PostgreSQL Database (`hulubedej_db`)
- **Pillar 2:** Next.js 15 Presentation Layer (`frontend/`)
- **Pillar 3:** Python AI Inference Worker & Offline Data Platform (`ai/`)

Every subsystem was evaluated against concrete code execution, live HTTP endpoints, PostgreSQL read-only inspections, and production compilation suites.

### Key Audit Conclusions:
1. **Database as Sole Taxonomy Authority:** Confirmed. PostgreSQL `category_hierarchy` contains 2,380 active categories structured up to 4 levels deep. The database enforces hierarchy constraints (`no_circular_reference`, `valid_level`, `fk_products_category_hierarchy`).
2. **Frontend Category Data Flow:** Confirmed 100% database-driven. All production frontend components (`HomeHero`, `CategoryDrawer`, `FilterSidebar`, `products/page.tsx`, `SearchBar`, `ProductForm`, `CategoryTreeSelector`, `admin/categories/page.tsx`) fetch categories and trees dynamically from the Go backend. Legacy static category files (`categoryTree.ts`, `categories.ts`) are completely unimported by production UI components.
3. **Backend Validation Authority:** Confirmed. The Go backend authoritatively verifies category existence, active status, hierarchy descent (`IsDescendantOf`), seller tier authorization, listing fee status, and image URLs on all listing mutations.
4. **Targeted Defect Discovered & Repaired:**
   - **Attribute Null Scan Defect:** In [`attribute_repo.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/repository/attribute_repo.go#L125), `GetCategoryAttributes`, `GetAttributeDefinition`, `GetAttributeBySlug`, and `GetProductAttributes` crashed with HTTP 500 when scanning nullable `unit` or `description` columns into non-pointer strings. Fixed with explicit `COALESCE(unit, '')` queries and added 3 regression tests.
   - **Category Ancestors Query Defect:** Fixed query joining recursive function `get_category_ancestors($1)` to provide localized breadcrumb paths ordered by `level ASC`. Added 4 regression tests.
   - **Admin Notification Schema Mismatch:** Fixed with forward migration `033_admin_notifications_type.sql` adding `type` and `read_by`, and mapped `priority` semantics in Go. Added 2 regression tests.
   - **Nominatim User-Agent:** Corrected legacy SellBuy string to `HUNTIQ-Universal-Marketplace/1.0`. Added 1 regression test.
5. **Release Decision:** **READY WITH CONDITIONS** (Controlled Staging Deployment approved; production deployment subject to completing the 3 database data quality reconciliations documented herein).

---

## 2. Architecture & Data-Flow Dependency Map

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                              HUNTIQ UNIVERSAL ARCHITECTURE                             │
└────────────────────────────────────────────────────────────────────────────────────────┘

    [ Next.js 15 Frontend (Port 3000) ]
        │  ▲
        │  │ HTTP (REST + JWT Bearer)
        ▼  │
    [ Go Backend Sovereign Authority (Fiber v2 on Port 8080) ]
        │                  │                           ▲
        │ Database         │ HTTP Internal Secret      │ Offline Dataset
        │ Connection Pool  │ (X-HUNTIQ-AI-SECRET)      │ Snapshots
        │ (pgxpool)        ▼                           │
        │             [ Python AI Worker ]        [ AI Data Platform ]
        │             (127.0.0.1:8000)            (Offline Governance)
        ▼             (Zero DB Credentials)       (No Raw PII in Git)
    [ PostgreSQL 18.4 Authority ]
    (hulubedej_db: 33 Migrations, 2,380 Categories, Sovereign Business Rules)
```

### Detailed Subsystem Data Flows:

1. **Admin Category Management:**
   - Admin UI (`/admin/categories`) → `AdminAuthMiddleware` (verifies admin JWT + role/permission) → `CategoryHandler.Admin*` → `CategoryHierarchyService` (checks slug uniqueness, circular references) → `CategoryHierarchyRepo` → PostgreSQL `category_hierarchy`.
2. **Seller Listing Creation & Updates:**
   - Seller UI (`/seller/product/new`) → `CategoryTreeSelector` fetches `GET /api/categories/tree?lang=locale` → Dynamic attributes loaded via `GET /api/categories/:id/attributes` → Submission via `POST /api/seller/products` → `ProductHandler.SellerCreate` → `ProductService.CreateForSeller`:
     - Validates category existence in database.
     - Validates `is_active = true`.
     - Validates subcategory belongs to parent via `catHierarchyRepo.IsDescendantOf`.
     - Validates listing fee options and prohibited items.
     - Writes listing with `status = 'pending'` and `fee_paid = FALSE` into `products`.
3. **Public Catalog, Search & Filtering:**
   - Client (`/products`, `HomeHero`, `CategoryDrawer`, `SearchBar`) → `api.getCategoryTree(locale)` → `CategoryHandler.GetCategoryTree` → `CategoryHierarchyRepo.GetCategoryTree` (recursive CTE with localized Amharic/Oromo/Tigrinya translation resolution) → In-memory client cache.
   - Filter query (`/products?category_hierarchy_id=14`) → `ProductHandler.ListApproved` → `ProductRepo.ListApproved` → Parameterized SQL:
     ```sql
     AND (p.category_hierarchy_id IN (SELECT id FROM get_category_descendants($1))
          OR (p.category_hierarchy_id IS NULL AND p.category_id = $1))
     ```
   - Searches `search_vector @@ websearch_to_tsquery('simple', $query)` with brand synonym expansion.

---

## 3. Taxonomy Database Integrity Audit

Conducted read-only SQL audits against `category_hierarchy`:
- **Total Records:** 2,380 categories across 69 root domains (`min_level: 0`, `max_level: 3`).
- **Referential Integrity:** 0 orphan categories. Every child references a valid parent (`FOREIGN KEY (parent_id) REFERENCES category_hierarchy(id) ON DELETE RESTRICT`).
- **Cycle Prevention:** 0 circular references (`CHECK (id <> parent_id)` enforced).
- **Depth Consistency:** 100% of child categories satisfy `child.level = parent.level + 1`.
- **Materialized Path Finding:** 2,289 non-root categories have paths where the seed script repeated the parent ID instead of appending the child ID (e.g. `/46/46/` instead of `/46/466/`). While recursive CTE lookups ignore this path, a non-destructive path repair query has been formulated in the accompanying data quality report.
- **Selectable Leaves:** 1,920 terminal leaf categories (`is_leaf = true`). Root categories have `is_leaf = false` to require sellers to select specific subcategories.

---

## 4. Category Configuration Findings

| Finding | Category ID(s) | Description | Severity | Status / Remediation |
| :--- | :--- | :--- | :--- | :--- |
| **Missing Root Translations** | 1–11 | Root categories lack `name_am`, `name_om`, `name_ti` in `category_hierarchy` (migration 017 omitted them from the insert). | Medium | Translations exist in legacy `categories` table. Remediation SQL provided in data quality report. |
| **Missing Tigrinya Translations** | 14–37 | Subcategories 14–37 lack `name_ti` (seeded before migration 031 added the column). | Low | English/Amharic/Oromo populated. Tigrinya translations queued for review. |
| **Home Appliance Sibling Duplicates** | 163 & 6509, 167 & 6518 | Two duplicate sibling categories exist under Parent 44 ("Home Appliances"): "Refrigerators and Freezers" and "Washing Machines and Dryers". | Medium | Candidate merge script provided to consolidate children under 6509 and 6518. |
| **Icons on Subcategories** | Subcategories | Only the 60 root categories have icon identifiers (e.g. `Utensils`, `Laptop`); subcategories have empty string. | Low | Intended design (drawer and tree UI use folder/chevron glyphs for subcategories). |
| **Category Image URLs** | All | `image_url` column added in migration 032 is currently NULL for existing categories. | Low | System functions with SVG icon fallbacks until marketing CDN assets are uploaded. |

---

## 5. Frontend Hardcoded-Data Inventory

Every category-related constant or symbol across `frontend/src` was audited and classified according to the audit rubric:

| File & Line | Data Type | Production Reachable? | Classification | Status & Notes |
| :--- | :--- | :--- | :--- | :--- |
| `lib/constants/categoryTree.ts:13` | Static Object (`CATEGORY_TREE`) | **NO** | **Deprecated but unused** | Retained with `@deprecated` notice for legacy reference. 0 imports in production code. |
| `lib/constants/categories.ts:28` | Helper Func (`getCategorySlugs`) | **NO** | **Deprecated but unused** | Pure type helper for validators. 0 runtime category data injected into UI. |
| `lib/language/*.json:35-80` | JSON Translation Objects | **YES** | **Presentation-only constant** | Standard i18n UI strings for generic fallback labels. |
| `components/public/FilterSidebar.tsx:15` | Array (`CONDITIONS`) | **YES** | **Legitimate static business rule** | Defines `new` and `used` condition values matching PostgreSQL check constraints. |
| `components/public/HomeHero.tsx:44` | Array (`SLIDE_MODES`) | **YES** | **Presentation-only constant** | Marketing hero toggle for "Buy / Rent" modes; categories themselves loaded dynamically via `api.getRootCategories()`. |
| `components/public/SearchBar.tsx:71` | Dynamic State (`searchIndex`) | **YES** | **Canonical Database-Driven** | Populated entirely in memory from `api.getCategoryTree(locale)` across all 3 levels. |
| `components/public/CategoryDrawer.tsx:56` | Dynamic State (`tree`) | **YES** | **Canonical Database-Driven** | Populated entirely from `api.getCategoryTree(locale)`. |
| `components/seller/ProductForm.tsx:338` | Dynamic Selector Component | **YES** | **Canonical Database-Driven** | Uses `<CategoryTreeSelector />` which queries `api.getCategoryTree(language)`. |
| `app/products/page.tsx:121` | Dynamic State (`categories`) | **YES** | **Canonical Database-Driven** | Populated from `api.getCategoryTree(locale)` with multi-level slug/ID resolution. |
| `app/admin/categories/page.tsx:62` | Dynamic Tree Grid | **YES** | **Canonical Database-Driven** | Full administrative CRUD and status toggle loaded from Go backend. |

---

## 6. Backend / API Consistency Findings

### Discovered & Repaired: Attribute Repo Null Scan Defect
- **Location:** [`backend/internal/repository/attribute_repo.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/repository/attribute_repo.go)
- **Defect:** `GetCategoryAttributes` (called by `GET /api/categories/:id/attributes` and `ProductForm.tsx`) crashed with `HTTP 500` because 22 attribute definitions have `unit = NULL` in PostgreSQL. Scanning NULL into Go's non-pointer `string` fields failed.
- **Repair:**
  - Updated `GetCategoryAttributes` query to `SELECT attribute_id, name_en, slug, data_type, validation_rules, COALESCE(unit, ''), is_required, is_filterable, display_order FROM get_category_attributes($1)`.
  - Updated `GetAttributeDefinition`, `GetAttributeBySlug`, and `GetProductAttributes` to use `COALESCE(unit, '')` and `COALESCE(description, '')`.
  - Added regression test suite in [`attribute_repo_test.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/repository/attribute_repo_test.go).
- **Verification:** `curl -s "http://localhost:8080/api/categories/14/attributes"` returned HTTP 200 with 6 valid attributes (`Brand`, `Model`, `Storage`, `RAM`, `Screen Size`, `Processor`).

---

## 7. Search & Filtering Audit

- **Filter Integrity:**
  - Selecting a category in `FilterSidebar` triggers `api.getCategoryChildren(id)` to dynamically load subcategories.
  - When the user selects a new parent category, all incompatible child selections are cleared (`setSubcategory(''); setItem('');`).
  - Empty search filters return `[]` cleanly without leaking unapproved or cross-tenant products.
- **SQL Parameterization:**
  - All category and location filters in [`product_repo.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/repository/product_repo.go) use strictly parameterized `$N` placeholders. No user-controlled strings are concatenated into the SQL grammar.
- **Full-Text Search:**
  - Uses PostgreSQL `search_vector @@ websearch_to_tsquery('simple', $query)` indexed with GIN index `idx_products_search`.
  - Automatically extracts brand names using `extract_brands($query)` and boosts attribute matches.

---

## 8. Security Findings

- **Authentication & Authorization:**
  - Standard user endpoints protected via JWT middleware (`UserID` extracted from claims).
  - Admin endpoints enforce `AdminAuthMiddleware` with granular role checks (`SuperAdmin`, `PlatformOwner`, `PlatformAdmin`, `Moderator`, `OperationsAdmin`, `SecurityAdmin`, `SupportAgent`).
  - Category deletion blocks if active child subcategories or products reference the category.
- **AI Service Boundary Protection:**
  - Python AI service is strictly isolated to `127.0.0.1:8000`.
  - Requires pre-shared secret in header `X-HUNTIQ-AI-SECRET`. Rejects unauthorized requests with RFC 7807 problem details.
  - Python worker possesses zero PostgreSQL database write credentials.
- **Content Security & CORS:**
  - Security headers enforced: `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`, `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`.

---

## 9. Test Execution & Verification Matrix

All commands were directly executed in the environment. Passing unit tests are recorded alongside production compilation and live API checks:

| Subsystem | Command | Exit Code | Execution Time | Results |
| :--- | :--- | :--- | :--- | :--- |
| **Python AI Worker** | `cd ai && uv run pytest tests/ -v` | `0` | 0.27s | **15 passed** (contract, config, health, security) |
| **AI Data Platform** | `cd ai/data-platform && uv run pytest tests/ -v` | `0` | 0.10s | **10 passed** (validation, PII scrub, manifest) |
| **Go Backend Unit & Integration** | `cd backend && go test -count=1 ./...` | `0` | ~3.8s | **ALL PASSED** across all packages |
| **Go Category Ancestors Regressions** | `go test -v -run TestCategoryHierarchyRepo_GetBreadcrumbs ./internal/repository` | `0` | 0.12s | **4 passed** |
| **Go Admin Notification Regressions** | `go test -v -run TestAdminRepo.*Notification ./internal/repository` | `0` | 0.09s | **2 passed** |
| **Go Attribute Null-Scan Regressions** | `go test -v -run TestAttributeRepo ./internal/repository` | `0` | 0.08s | **3 passed** |
| **Go Geocoding User-Agent Regressions** | `go test -v ./internal/service/geocoding/...` | `0` | 0.01s | **5 passed** |
| **Go Production API Binary Build** | `cd backend && go build -o /tmp/huntiq-test-api ./cmd/api` | `0` | ~3.0s | **PASS** (Zero warnings, zero errors) |
| **Next.js TypeScript Check** | `cd frontend && npx tsc --noEmit` | `0` | ~7.0s | **PASS** (Zero type errors) |
| **Next.js ESLint Check** | `cd frontend && npm run lint` | `0` | ~5.0s | **PASS** (0 errors, standard warning rules) |
| **Next.js Production Build** | `cd frontend && npm run build` | `0` | ~45.0s | **PASS** (63/63 routes compiled) |

---

## 10. Files Changed & Root Causes

1. [`backend/internal/repository/category_hierarchy_repo.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/repository/category_hierarchy_repo.go): Fixed breadcrumbs query to join `get_category_ancestors($1)` and return localized rows ordered by `level ASC`. Initialized slice as `[]` for clean JSON empty arrays.
2. [`backend/internal/repository/admin_repo.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/repository/admin_repo.go): Mapped categorical notification types to urgency priorities and inserted both `type` and `priority`.
3. [`backend/internal/database/migrations/033_admin_notifications_type.sql`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/database/migrations/033_admin_notifications_type.sql): Forward migration adding `type` and `read_by` columns with indexes to `admin_notifications`.
4. [`backend/internal/database/db.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/database/db.go): Added idempotent schema guard for `admin_notifications` on server startup.
5. [`backend/internal/repository/attribute_repo.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/repository/attribute_repo.go): Added `COALESCE(unit, '')` and `COALESCE(description, '')` to eliminate pgx NULL scan errors across `GetCategoryAttributes`, `GetAttributeDefinition`, `GetAttributeBySlug`, and `GetProductAttributes`.
6. [`backend/internal/service/geocoding/geocoding.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/service/geocoding/geocoding.go): Replaced legacy SellBuy string with `HUNTIQ-Universal-Marketplace/1.0`.
7. [`backend/internal/repository/category_hierarchy_repo_test.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/repository/category_hierarchy_repo_test.go): Added 4 regression tests for category ancestors.
8. [`backend/internal/repository/admin_repo_test.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/repository/admin_repo_test.go): Added 2 regression tests for notification creation and read tracking.
9. [`backend/internal/repository/attribute_repo_test.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/repository/attribute_repo_test.go): Added 3 regression tests for attribute repository null safety.
10. [`backend/internal/service/geocoding/geocoding_test.go`](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/service/geocoding/geocoding_test.go): Added user-agent regression test.
11. [`backend/docs/audits/master-integration-audit-2026.md`](file:///home/poven/GO%20PROJECT/sellbuy/backend/docs/audits/master-integration-audit-2026.md): Updated audit text to accurately document fee verification on `products` table and clarify PostgreSQL 16+ baseline.
12. [`backend/docs/audits/huntiq-category-data-quality.md`](file:///home/poven/GO%20PROJECT/sellbuy/backend/docs/audits/huntiq-category-data-quality.md): Authored comprehensive category data quality report.

---

## 11. Prioritized Remediation Plan

### Immediate / Pre-Production (P1):
- Apply Remediation Step 1 in `huntiq-category-data-quality.md` to backfill Amharic, Oromo, and Tigrinya translations for Root Categories 1–11 from the legacy `categories` table.
- Apply Remediation Step 2 in `huntiq-category-data-quality.md` to recompute the materialized `path` values for the 2,289 categories with duplicate parent segments.

### Near-Term Staging (P2):
- Human review of duplicate sibling subcategories 163 vs 6509 and 167 vs 6518 under "Home Appliances" to merge children and deactivate redundant nodes.
- Backfill Tigrinya translations for subcategories 14–37.
- Upload high-resolution SVG or WebP category banner assets and populate `image_url` on root categories.

---

## 12. Release Recommendation

### **DECISION: READY WITH CONDITIONS**

**Justification:**
1. **Compilation & Test Suites:** 100% clean across Go, Python, Next.js, and TypeScript. All 63 Next.js routes bundle without errors. Full Go test suite passes (`go test -count=1 ./...`).
2. **Authority & Separation:** The Go backend remains the sovereign authority. Frontend category data is 100% dynamically loaded from PostgreSQL. The Python AI service is sandboxed without database write access.
3. **Reproducible Defects Resolved:** All runtime blockers (ancestors 500, admin notifications 500, category attributes 500) have been fixed with minimal, targeted changes and verified with regression tests.
4. **Conditions for Production Promotion:**
   - Execute the non-destructive translation sync and path recomputation scripts in staging.
   - Complete human review on the 2 duplicate subcategories under Home Appliances.
   - Perform end-to-end hardware passkey validation in a browser environment with authenticators.
