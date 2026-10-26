# HUNTIQ — Phase 10: Final Database Repair Approval Check & Pre-Release Readiness Report

**Date:** October 2026  
**Auditor:** Principal Software Architect & Verification Lead  
**System:** HUNTIQ Universal Marketplace (Ethiopia)  
**Repository:** `huntiq` (`~/GO PROJECT/sellbuy`)  
**Target Database:** PostgreSQL 18.4 (Debian 18.4-1+b2), Database `hulubedej_db` on `localhost:5432` (User: `postgres`)  
**Baseline Verification Status:** Clean baseline, all working user changes preserved, 100% test pass rate across Go, TypeScript, and Python.

---

## 1. Verified Facts vs. Unverified Claims

| Subject | Verified Fact (with Concrete Evidence) | Rejected Unverified Claim / Fallacy |
| :--- | :--- | :--- |
| **Materialized Path Bug (`AUD-001`)** | Read-only CTE query proved exactly 2,289 rows in `category_hierarchy` have corrupted paths due to `parentID` duplication in `seed-taxonomy/main.go:140`. | *Claim that all browse flows were broken.* Recursive CTEs join on `parent_id`, so categories loaded normally, but prefix searches and attribute inheritance were affected. |
| **Root Category Translations (`AUD-002`)** | Read-only join confirmed Categories 1–11 in `category_hierarchy` have empty `name_am`, `name_om`, and `name_ti`, while legacy `categories` has them fully populated. | *Claim that data was permanently lost.* Translations were simply omitted during migration 017's `INSERT` column list. |
| **Dead Frontend Components (`AUD-012`/`AUD-013`)** | Codebase-wide ripgrep across `frontend/src` confirmed 0 callers for 8 components (`AvailabilityGrid`, `EnterpriseDataTable`, `LocationManager`, `Stats`, `ImageVerification`, `CategoryFilter`, `RentalListingForm`, `DomainSwitcherTabs`) and 2 obsolete constants (`categoryTree.ts`, `categories.ts`). | *Claim that deleting them might break dynamic routes.* Next.js 15 App Router discovers routes strictly from `app/**/page.tsx`; all 63 routes build cleanly without them. |
| **Search Analytics Wiring (`AUD-009`)** | `SearchRepo.LogSearch` was unwired in `main.go:70` (`_ = searchRepo`). Database check confirmed `search_analytics` had 0 rows. | *Claim that search analytics was working.* Admin popular search views were empty because the recording hook was never invoked. |
| **`SearchRepo.LogClick` Query Defect** | PostgreSQL syntax inspection discovered `UPDATE ... ORDER BY ... LIMIT 1` in `search_repo.go:52` was invalid SQL in PostgreSQL (MySQL syntax leak). Fixed via subquery `WHERE id = (SELECT id ... LIMIT 1)`. | *Claim that click logging was production-ready.* Executing `LogClick` on PostgreSQL failed with syntax error at or near `ORDER`. |
| **Redundant Indexes (`AUD-008`)** | `pg_constraint` query proved 0 constraints depend on the 5 candidate indexes, and `EXPLAIN` query plans confirmed the query planner already uses the primary UNIQUE constraint indexes. | *Claim that dropping redundant indexes harms query performance.* Unconditional UNIQUE indexes provide identical or superior index coverage with less write I/O. |
| **Duplicate Categories (`AUD-003`)** | Database query confirmed Categories 163 vs 6509 and 167 vs 6518 exist under Parent 44, and confirmed **0 active products** reference these 4 categories or their 22 child subcategories. | *Claim that duplicate categories can be deleted immediately.* Category deprecation/merging affects URLs and public facets, requiring explicit product-owner sign-off. |
| **Five Branch Categories without Children** | Catalog generator code analysis proved `domain_master_expander.go:26,71` passed un-prefixed `sub.Slug` to children of 6563, parenting them to sibling 1301. Root categories 7, 9, 10, 11 are Level 0 domain roots awaiting future subcategory seeds. | *Claim that the 5 categories are corrupt and should be deleted.* Categories 7, 9, 10, 11 are canonical domains; 6563 is a structural branch whose children were misrouted. |

---

## 2. Evidence of Test Execution & Exit Codes

All test commands below were executed during verification:

```bash
# 1. Full Backend Uncached Unit & Integration Test Suite
cd backend && go test -count=1 ./...
# Exit Code: 0 (PASS: 100% packages ok, including new search analytics tests)

# 2. Backend Static Analysis (Go Vet)
cd backend && go vet ./...
# Exit Code: 0 (PASS: 0 warnings, 0 vet issues)

# 3. Targeted Search Analytics Regression Tests
cd backend && go test -v ./internal/handler -run "TestProductHandler_SearchAnalytics"
# Exit Code: 0 (PASS: 4/4 test cases pass in 0.014s)

# 4. Frontend TypeScript Static Typecheck
cd frontend && npx tsc --noEmit
# Exit Code: 0 (PASS: 0 compile errors)

# 5. Frontend ESLint Verification
cd frontend && npm run lint
# Exit Code: 0 (PASS: 0 errors, standard Next.js entity/image warnings only)

# 6. Frontend Production Build (App Router)
cd frontend && npm run build
# Exit Code: 0 (PASS: Compiled successfully, generated 63/63 static & dynamic routes)

# 7. Python AI Worker Unit and Contract Tests
cd ai && uv run pytest -v
# Exit Code: 0 (PASS: 15/15 tests passed in 0.25s)

# 8. Python Data Platform Quality and Contract Tests
cd ai/data-platform && uv run pytest tests/ -v
# Exit Code: 0 (PASS: 10/10 tests passed in 0.06s)
```

---

## 3. Reconciled Master Findings Ledger (AUD-001 through AUD-017)

| Finding ID | Domain & Original Issue | Clarified Technical Definition & Scope | Status | Changed Files | Verification Tests | Production / DB Action Required? |
| :--- | :--- | :--- | :--- | :--- | :--- | :---: |
| **AUD-001** | Seed Taxonomy Path Defect | `seed-taxonomy/main.go:140` appended `parentID` instead of `childID`, generating corrupted paths (e.g. `/38/38/12/`) for 2,289 categories. | **Prepared for manual execution** | `backend/cmd/seed-taxonomy/main.go` | CTE read-only query identified 2,289 rows; dry-run executed with rollback confirmed 2,289 rows repaired. | **YES** (Execute `repair_category_paths.sql`) |
| **AUD-002** | Missing Root Translations | Migration 017 omitted `name_am`, `name_om`, `name_ti` for root categories 1–11 in `category_hierarchy`. Data exists in `categories`. | **Prepared for manual execution** | `backend/scripts/backfill_root_translations.sql` | Read-only join confirmed 11/11 empty in hierarchy; dry-run with rollback confirmed all 11 backfilled safely. | **YES** (Execute `backfill_root_translations.sql`) |
| **AUD-003** | Sibling Duplicate Categories | Under Parent 44 ("Home Appliances"), IDs 163 vs 6509 and 167 vs 6518 duplicate refrigerator and washing machine trees. | **Requires product-owner approval** | None | Database inspection confirmed 0 products reference these 4 categories or their 22 child nodes. | **YES** (Pending product owner sign-off) |
| **AUD-004** | Admin Category API Duality | Legacy administrative CRUD endpoints mutated flat `categories` table. Unused by production Next.js admin console. | **Fixed and verified** | `frontend/src/app/admin/categories/page.tsx`, `frontend/src/lib/api/admin.ts` | Grep confirmed 0 frontend callers to legacy flat admin APIs; primary admin page verified using `category-hierarchy`. | **NO** |
| **AUD-005** | Admin Location API Duality | Legacy administrative location CRUD operated on flat `locations` table. Unused by production admin console. | **Fixed and verified** | `frontend/src/app/admin/locations/page.tsx`, `frontend/src/lib/api/admin.ts` | Grep confirmed 0 frontend callers; admin location page verified using `location_hierarchy`. | **NO** |
| **AUD-006** | Universal Expansion Tables | Migration 031 created 10 recruitment/machinery tables with 0 application consumers. Invariants prohibit dropping tables. | **Deferred by design** | None | Schema inspection confirmed tables isolated with 0 active dependencies. | **NO** (Retained as Phase 2 extensions) |
| **AUD-007** | Disconnected Brand Catalogs | `brand_dictionary` (31 rows) used in `extract_brands(p_query)` for search matching, while `products.brand_id` references `brands` (72 rows). | **Requires product-owner approval** | None | Schema inspection proved search query matching only searches keyword tokens; foreign keys point to `brands`. | **YES** (Pending architectural decision) |
| **AUD-008** | Redundant Database Indexes | 5 secondary indexes duplicate existing UNIQUE constraints on `category_hierarchy`, `attribute_definitions`, `users`, and `brands`. | **Prepared for manual execution** | `backend/internal/database/migrations/034_drop_redundant_duplicate_indexes.sql` | `pg_constraint` confirmed 0 constraints depend on them; query plans verified planner uses UNIQUE indexes. | **YES** (Apply migration 034 during maintenance) |
| **AUD-009** | Unwired Search Analytics | `SearchRepo.LogSearch` was unwired in `main.go`. `search_analytics` table had 0 rows. | **Fixed and verified** | `backend/internal/handler/product_handler.go`, `backend/cmd/api/main.go`, `backend/internal/repository/search_repo.go`, `backend/internal/handler/product_handler_test.go` | Wired bounded worker (cap 500) with 3s timeout. Added 4 unit tests covering event dispatch, 0-results, failures, and drops. Fixed `LogClick` SQL syntax. | **NO** (Active on server startup) |
| **AUD-010** | Triplicate Image Attachment | Batch image loading copied across 3 repository methods. Public listings load thumbnail only; admin loads all images. | **Deferred by design** | None | Code inspection proved intentional semantic difference; premature consolidation risks payload bloat. | **NO** |
| **AUD-011** | Monolithic AdminRepo | `admin_repo.go` aggregates 9 domains in 2,643 lines. Pass-through methods exist in `admin_service.go`. | **Deferred by design** | `backend/internal/repository/admin_repo_test.go` | Added repository regression tests; high risk of regression without immediate functional benefit. | **NO** |
| **AUD-012** | Dead Frontend Components | 8 prototype components had 0 callers and contained hardcoded mock data (`AvailabilityGrid`, `EnterpriseDataTable`, etc.). | **Fixed and verified** | Deleted 8 component files in `frontend/src/components/` | Codebase grep confirmed 0 references; Next.js production build compiled 63/63 routes cleanly. | **NO** |
| **AUD-013** | Obsolete Static Constants | `categoryTree.ts` and `categories.ts` were static files marked `@deprecated` with 0 callers. | **Fixed and verified** | Deleted `categoryTree.ts` and `categories.ts` in `frontend/src/lib/constants/` | Grep confirmed 0 imports; production build succeeded. | **NO** |
| **AUD-014** | Missing Seller Attribute API | `ProductForm.tsx` bypassed client with inline `fetch` to save dynamic attributes. | **Fixed and verified** | `frontend/src/lib/api/seller.ts`, `frontend/src/components/seller/ProductForm.tsx` | Added typed `saveProductAttributes` to `sellerApi`; refactored `ProductForm.tsx`. | **NO** |
| **AUD-015** | Dual Category Selector UX | `ProductForm.tsx` rendered two category dropdowns for every listing (flat category + tree selector), causing dual state. | **Fixed and verified** | `frontend/src/components/seller/ProductForm.tsx` | Unified category selection to single `CategoryTreeSelector`; automatically populates hierarchy ID and legacy ID. | **NO** |
| **AUD-016** | Unwired Python AI Platform | AI worker is operational on `:8000` but has no model inference route. Invariant Rule 2 forbids fake AI endpoints. | **Deferred by design** | None | Passes all 15 unit/contract tests; candidates must be trained offline first per Section 4 of `AGENTS.md`. | **NO** |
| **AUD-017** | Inline DDL in Startup Bootstrap | `db.go:ensureSchema` executes idempotent DDL on startup duplicating migrations 029 and 033. | **Deferred by design** | None | Serves as safety net for bare developer environments without migration runners. | **NO** |

---

## 4. Exact Trigger, Constraint, and Index Impact Matrix

### 1. `backend/scripts/repair_category_paths.sql`
- **Target Table:** `public.category_hierarchy`
- **Columns Modified:** `path` (`VARCHAR(500)`), `updated_at` (`TIMESTAMP WITH TIME ZONE`).
- **Trigger Impact:**
  - `category_hierarchy` has exactly ONE trigger:
    ```sql
    CREATE TRIGGER trg_update_category_hierarchy_path 
    BEFORE INSERT OR UPDATE OF parent_id ON category_hierarchy 
    FOR EACH ROW EXECUTE FUNCTION update_category_hierarchy_path();
    ```
  - **Invariant Check:** Because this trigger explicitly specifies `UPDATE OF parent_id`, and `repair_category_paths.sql` only assigns `path` and `updated_at`, PostgreSQL skips this trigger entirely. The trigger will **NOT** fire, preventing recursive recalculations or unwanted side-effects.
- **Audit Mechanism Impact:**
  - Database-wide inspection confirmed **0 triggers** write to `admin_audit_logs` or audit tables. Admin auditing in HUNTIQ is strictly an application-layer concern (`AdminRepo.LogAdminAction`).
- **Generated Column Impact:**
  - Database schema query confirmed `category_hierarchy` has **0 generated columns** (only `products.search_vector` exists in `hulubedej_db`).
- **Constraint Impact:**
  - `category_hierarchy_pkey` (`id`): Untouched.
  - `uq_category_hierarchy_slug` (`slug`): Untouched.
  - `no_circular_reference` (`id <> parent_id`): Untouched.
  - `valid_level` (`level >= 0 AND level < 20`): Untouched.
  - `category_hierarchy_parent_id_fkey`: Untouched.
  - All foreign keys referencing `category_hierarchy` (`products`, `category_attributes`, `product_families`, etc.) remain 100% valid.
- **Index Impact:**
  - Updates B-Tree index `idx_category_hierarchy_path` with correct materialized path strings.

### 2. `backend/scripts/backfill_root_translations.sql`
- **Target Table:** `public.category_hierarchy` (destination), `public.categories` (source).
- **Columns Modified:** `name_am`, `name_om`, `name_ti`, `updated_at` on rows where `id BETWEEN 1 AND 11`.
- **Trigger & Audit Impact:** Neither table has triggers listening to name columns. 0 triggers fire.
- **Constraint Impact:** Target columns have no unique or foreign key constraints. 0 constraints affected.
- **Preservation & Overwrite Safety:**
  - Target columns currently contain `NULL` across all 33 cells (11 rows × 3 languages).
  - SQL uses `CASE WHEN COALESCE(ch.name_*, '') = '' ... ELSE ch.name_* END`, guaranteeing that any existing valid non-empty translation is preserved.

### 3. `backend/internal/database/migrations/034_drop_redundant_duplicate_indexes.sql`
- **Target Statements:**
  1. `DROP INDEX IF EXISTS idx_category_hierarchy_slug_unique;`
  2. `DROP INDEX IF EXISTS idx_category_hierarchy_slug;`
  3. `DROP INDEX IF EXISTS idx_attribute_definitions_slug;`
  4. `DROP INDEX IF EXISTS idx_users_phone_role;`
  5. `DROP INDEX IF EXISTS idx_brands_slug;`
- **Constraint Dependency Check:**
  - Query against `pg_constraint.conindid` confirmed **0 constraints** rely on any of these 5 indexes.
- **Preserved Semantics & Partial-Index Review:**
  - `idx_category_hierarchy_slug_unique` has partial predicate `WHERE (is_active = true)`.
  - Surviving index `uq_category_hierarchy_slug` is an **unconditional UNIQUE constraint** on `(slug)`.
  - Because no two rows (active or inactive) can share a slug under `uq_category_hierarchy_slug`, the partial unique index is a strict subset. Dropping it does not relax any business invariant.
  - Surviving unique constraint indexes (`uq_category_hierarchy_slug`, `attribute_definitions_slug_key`, `users_phone_role_key`, `brands_slug_key`) provide identical or superior index selectivity for all query plans.

---

## 5. Structural & Query Invariant Verification

### Path-Based Queries & Attribute Inheritance Verification
PostgreSQL function `get_category_attributes_inherited(p_category_id)` resolves category attributes by converting `category_hierarchy.path` into an array of ancestor IDs:
```sql
SELECT path INTO v_path FROM category_hierarchy WHERE id = p_category_id;
IF v_path IS NOT NULL AND v_path != '' THEN
    v_ancestor_ids := string_to_array(trim(both '/' from v_path), '/')::int[];
ELSE
    v_ancestor_ids := ARRAY[p_category_id];
END IF;
```

**Live Verification Results (Before vs. After Path Repair):**
1. **Subcategory Prefix Query Test (`Category 135` and its children `136..139`):**
   - Query: `SELECT COUNT(*) FROM category_hierarchy WHERE path LIKE '/39/135/%';`
   - **Before Repair (Current):** Returns **`0`** (corrupted path was `/39/39/135/`).
   - **After Repair (Simulated):** Returns **`5`** (matches parent 135 and all 4 child nodes 136, 137, 138, 139).
2. **Attribute Inheritance Test (`Category 138`):**
   - Query: `SELECT COUNT(*) FROM get_category_attributes_inherited(138);`
   - Returns all **8 inherited attributes** from domain 39 ("Computers, IT and Office Technology"), including Brand, Model, Condition, Processor, RAM, Storage, Screen Size, and OS.
3. **Index Scan Behavior under Collation `en_US.UTF-8`:**
   - Database collation is `en_US.UTF-8`. Exact equality lookups (`WHERE path = '/38/'`) perform an `Index Scan` using `idx_category_hierarchy_path`.
   - Sequential scan is selected by PostgreSQL for `LIKE` queries under 2,380 rows because total table footprint is only 294 pages (~2.3 MB), executing in 1.4ms.

---

## 6. Investigation of Branch Categories without Children (`is_leaf = false`)

Database verification confirmed exactly 5 categories are marked `is_leaf = false` but have 0 child rows in `category_hierarchy`:

```
  id  | parent_id |                name_en                |                 slug                 | level |  path   | is_leaf | is_active | display_order 
------+-----------+---------------------------------------+--------------------------------------+-------+---------+---------+-----------+---------------
    7 |           | Business & Industry                   | business                             |     0 | /7/     | f       | t         |             7
    9 |           | Lifestyle & Wellness                  | lifestyle                            |     0 | /9/     | f       | t         |             9
   10 |           | Education & Learning                  | education                            |     0 | /10/    | f       | t         |            10
   11 |           | Other                                 | other                                |     0 | /11/    | f       | t         |            11
 6563 |        46 | Cookware, Pots, Pans and Injera Mitad | kitchen-household-cookware-pots-pans |     1 | /46/46/ | f       | t         |             1
```

### Detailed Findings per Category:
1. **Categories 7, 9, 10, 11 (Level 0 Domain Roots):**
   - **Origin:** Seeded by migration `017_category_hierarchy.sql` from legacy `categories` table.
   - **Intended Role:** Broad Level 0 navigation domains. Marked `is_leaf = false` by marketplace design to force sellers to select specific subcategories.
   - **Status:** When the 60 universal taxonomy domains were generated, specialized domains (`business-services`, `commercial-supplies`, `health-medical`, `books-education`) were introduced. Categories 7, 9, 10, and 11 received 0 subcategories in `catalog_domains_*.go`.
   - **Impact:** 0 products, 0 attributes, and 0 models reference them.
   - **Action:** Retain as inactive or unpopulated domain placeholders; do NOT alter without product-owner sign-off.

2. **Category 6563 ("Cookware, Pots, Pans and Injera Mitad", Parent 46):**
   - **Origin:** Created during systematic catalog generation for `kitchen-household`.
   - **Intended Role:** Level 1 subcategory container for Ethiopian cookware (Injera mitad, clay disti, granite cookware, pressure cookers).
   - **Defect Root Cause Discovered in Code:**
     In `backend/cmd/seed-taxonomy/domain_master_expander.go`:
     ```go
     // Lines 25-27: Slug deduplication logic
     if seen[slug] {
         slug = fmt.Sprintf("%s-%s", parent, slug) // slug became "kitchen-household-cookware-pots-pans"
     }
     seen[slug] = true

     // Lines 70-71: Leaf addition loop passes un-prefixed sub.Slug!
     for leafIdx, leaf := range sub.Items {
         addNode(sub.Slug, leaf.Slug, leaf.NameEN, leaf.NameAM, leaf.NameOM, leaf.NameTI, true, leafIdx+1)
     }
     ```
     `addNode` received `sub.Slug` (`"cookware-pots-pans"`) as the parent slug instead of the deduplicated slug (`"kitchen-household-cookware-pots-pans"`).
     Consequently, child items 6564 through 6571 (including the Injera Mitad and Clay Stew Pots) resolved their parent to sibling subcategory `1301` ("cookware-pots-pans"), leaving Category 6563 with 0 children.
   - **Impact:** 0 products reference Category 6563.
   - **Action:** Category 6563 must **NOT** be changed automatically. When product owners resolve duplicate catalogs, children 6564–6571 can either remain under 1301 (with 6563 deactivated) or be re-parented to 6563.

---

## 7. Migration 034 Index Analysis & Safety

| Index to Drop | Table | Definition | Surviving Index & Constraint | Assessment |
| :--- | :--- | :--- | :--- | :--- |
| `idx_category_hierarchy_slug_unique` | `category_hierarchy` | `UNIQUE (slug) WHERE is_active = true` | `uq_category_hierarchy_slug` (`UNIQUE (slug)`) | **100% Redundant.** Surviving index is an unconditional UNIQUE constraint. |
| `idx_category_hierarchy_slug` | `category_hierarchy` | `INDEX (slug)` | `uq_category_hierarchy_slug` (`UNIQUE (slug)`) | **100% Redundant.** Non-unique B-tree identical in column and order to unique constraint index. |
| `idx_attribute_definitions_slug` | `attribute_definitions` | `INDEX (slug)` | `attribute_definitions_slug_key` (`UNIQUE (slug)`) | **100% Redundant.** Identical single-column B-tree. |
| `idx_users_phone_role` | `users` | `INDEX (phone, role)` | `users_phone_role_key` (`UNIQUE (phone, role)`) | **100% Redundant.** Composite index exactly duplicated by primary unique key. |
| `idx_brands_slug` | `brands` | `INDEX (slug)` | `brands_slug_key` (`UNIQUE (slug)`) | **100% Redundant.** Identical single-column B-tree. |

- **Total Disk Space Reclaimed:** ~576 kB.
- **I/O Savings:** Eliminates 5 duplicate B-tree updates on every category, user, brand, and attribute insert/update.
- **Constraint Dependency:** **0 constraints depend on these indexes.**

---

## 8. Root Translation Backfill Inspection

Read-only comparison of `category_hierarchy` (destination) vs `categories` (source) for IDs 1..11:

| ID | English Name (`name_en`) | Amharic Source (`categories.name_am`) | Oromo Source (`categories.name_or`) | Tigrinya Source (`categories.name_ti`) | Missing Destination Cells | Conflict Detected? |
| :-: | :--- | :--- | :--- | :--- | :-: | :-: |
| 1 | Electronics | ኤሌክትሮኒክስ | Elektirooniksii | ኤሌክትሮኒክስ | 3 (AM, OM, TI) | **None (0)** |
| 2 | Vehicles | ተሽከርካሪዎች | Konkolaataa | መኪናታት | 3 (AM, OM, TI) | **None (0)** |
| 3 | Property & Real Estate | ንብረት | Qabeenya | ንብረት | 3 (AM, OM, TI) | **None (0)** |
| 4 | Home & Living | ቤት እና የኑሮ | Mana fi Jireeny | ቤትን ህይወት | 3 (AM, OM, TI) | **None (0)** |
| 5 | Fashion & Clothing | ፋሽን | Faashinii | ፋሽን | 3 (AM, OM, TI) | **None (0)** |
| 6 | Agriculture & Farming | ግብርና | Qonna | ሓረስታይ | 3 (AM, OM, TI) | **None (0)** |
| 7 | Business & Industry | ንግድ | Daldala | ንግዲ | 3 (AM, OM, TI) | **None (0)** |
| 8 | Services & Jobs | አገልግሎቶች | Tajaajila | ኣገልግሎት | 3 (AM, OM, TI) | **None (0)** |
| 9 | Lifestyle & Wellness | የአኗኗር ዘይቤ | Akkaataa Jireenyaa | ኣነባብራ | 3 (AM, OM, TI) | **None (0)** |
| 10 | Education & Learning | ትምህርት | Barnoot | ትምህርቲ | 3 (AM, OM, TI) | **None (0)** |
| 11 | Other | ሌላ | Kan biraa | ካልእ | 3 (AM, OM, TI) | **None (0)** |

- **Total Missing Cells:** Exactly 33 cells across 11 rows.
- **Conflicts:** Exactly 0. No destination cell contains conflicting text.
- **Mapping Correctness:** `category_hierarchy.id` matches `categories.id` 1:1.
- **Verification Logic:** Uses `COALESCE(name_*, '') = ''` to guarantee NULL-safety and aborts transaction on any unpopulated root translation.

---

## 9. Search Analytics Operational & Privacy Review

### Concurrency & Performance Safety
- **Queue Bounds:** Enqueueing uses a buffered channel of capacity 500 (`searchChan`).
- **Zero-Latency Impact:** Non-blocking `select` drops events if the queue fills under burst traffic, guaranteeing **0ms latency overhead** on search requests.
- **Saturation Observability:** Added `atomic.AddUint64(&h.droppedCount, 1)` and `h.DroppedSearchEvents()` to monitor queue drop frequency.
- **Worker Timeout:** Dedicated background worker processes events with `context.WithTimeout(context.Background(), 3*time.Second)`.
- **Deduplication:** Only logs queries when `f.Page == 1`, preventing pagination requests from multiplying search metrics.
- **PostgreSQL Compliance:** Fixed non-standard MySQL syntax in `SearchRepo.LogClick` via subquery `WHERE id = (SELECT id ... ORDER BY created_at DESC LIMIT 1)`.

### Privacy & Data Minimization
- **Phone Number Protection:** Search queries from buyers may occasionally contain phone numbers. The normalized query removes whitespace and symbols. To prevent unscrubbed PII in analytics datasets, the data platform's PII scanner tests (`test_scan_for_pii_ethiopian_phone_numbers`) are integrated into offline pipelines.
- **IP Address Minimization:** `c.IP()` records client IP for security fraud analytics. For production privacy compliance, an IP truncation function (zeroing the last octet: `192.168.1.0` / IPv6 `/64`) is recommended prior to long-term archiving.

---

## 10. Remaining Product-Owner Decisions & Blockers

1. **AUD-003 (Duplicate Sibling Subcategories under Parent 44):**
   - Option A: Re-parent children of Category 163 to 6509 and Category 167 to 6518, then mark 163 and 167 as `is_active = false`.
   - Option B: Retain both branches with differentiated display names.
   - *Recommendation:* Option A (eliminates duplicate "Refrigerators and Freezers" and "Washing Machines and Dryers" facets).
2. **AUD-007 (Brand Catalog Duality):**
   - Option A: Create an automated PostgreSQL trigger/view synchronizing `brand_dictionary` from canonical `brands`.
   - Option B: Migrate `extract_brands(p_query)` to query `brands` directly with an added `aliases JSONB` column.
   - *Recommendation:* Option B (single source of truth).
3. **5 Leaf-Less Branch Categories (IDs 7, 9, 10, 11, and 6563):**
   - Categories 7, 9, 10, 11: Retain as empty Level 0 domain placeholders until specific subcategories are designed.
   - Category 6563: Deactivate (`is_active = false`) or re-parent children 6564–6571 from 1301 to 6563.

---

## 11. Independent GO / NO-GO Recommendations

| Script / Operation | Recommendation | Justification & Verification Evidence |
| :--- | :---: | :--- |
| **`backend/scripts/repair_category_paths.sql`** | **GO** | Fully idempotent, transactional, verified in dry-run with rollback. Corrects exactly 2,289 rows. Bypasses `trg_update_category_hierarchy_path` (parent_id not updated). Restores prefix query matching from 0 to 5 for test subtrees. 0 constraints violated. |
| **`backend/scripts/backfill_root_translations.sql`** | **GO** | Fully non-destructive, idempotent, transactional. Backfills exactly 11 rows (33 translation cells) from canonical legacy table. 0 conflicts. Preserves any existing text via `COALESCE` guardrail. |
| **`backend/internal/database/migrations/034_drop_redundant_duplicate_indexes.sql`** | **GO** | 0 constraints depend on candidate indexes. All 5 indexes are completely redundant with existing primary UNIQUE constraint indexes. Safe to execute during scheduled maintenance window to avoid lock contention. |

---

## 12. Final Manual Database Operations Approval Checklist

### Target Database & Environment
- **Host / Port:** `localhost:5432`
- **Database Name:** `hulubedej_db`
- **Database Engine:** PostgreSQL 18.4 (Debian 18.4-1+b2)
- **Role / User:** `postgres`

### Backup & Recovery Requirements
Prior to executing any change, an operator MUST execute:
```bash
# 1. Full binary dump of the database
pg_dump -U postgres -Fc hulubedej_db > /tmp/hulubedej_db_pre_repair_$(date +%Y%m%d_%H%M%S).dump

# 2. Fast point-in-time snapshot tables inside the database
psql -U postgres -d hulubedej_db -c "
CREATE TABLE category_hierarchy_backup_pre_repair AS SELECT * FROM category_hierarchy;
"
```

### Ordered Execution Commands
Execute the approved scripts in the following strict sequential order:

```bash
# Step 1: Repair Materialized Category Paths (AUD-001)
psql -U postgres -d hulubedej_db -f backend/scripts/repair_category_paths.sql
# Expected Output:
# UPDATE 2289
# NOTICE: Path repair verification successful: 0 mismatches remain.
# COMMIT

# Step 2: Backfill Root Category Translations (AUD-002)
psql -U postgres -d hulubedej_db -f backend/scripts/backfill_root_translations.sql
# Expected Output:
# UPDATE 11
# NOTICE: Translation backfill verification successful: all 11 root categories translated.
# COMMIT

# Step 3: Drop Redundant Secondary Indexes (AUD-008)
psql -U postgres -d hulubedej_db -f backend/internal/database/migrations/034_drop_redundant_duplicate_indexes.sql
# Expected Output:
# DROP INDEX (x5)
```

### Expected Affected-Row Counts
- `repair_category_paths.sql`: Exactly **2,289 rows** updated.
- `backfill_root_translations.sql`: Exactly **11 rows** updated.
- `034_drop_redundant_duplicate_indexes.sql`: Exactly **5 indexes dropped**, 0 rows modified.

### Post-Execution Validation Queries (Read-Only)
```sql
-- 1. Verify 0 corrupted paths remain in category hierarchy
WITH RECURSIVE expected_tree AS (
    SELECT id, parent_id, '/' || id || '/' AS expected_path
    FROM category_hierarchy WHERE parent_id IS NULL
    UNION ALL
    SELECT c.id, c.parent_id, et.expected_path || c.id || '/'
    FROM category_hierarchy c JOIN expected_tree et ON c.parent_id = et.id
)
SELECT 
    COUNT(*) AS total_categories,
    COUNT(*) FILTER (WHERE ch.path = et.expected_path) AS valid_paths,
    COUNT(*) FILTER (WHERE ch.path != et.expected_path OR ch.path IS NULL) AS corrupted_paths
FROM category_hierarchy ch
JOIN expected_tree et ON ch.id = et.id;
-- Expected: total_categories: 2380 | valid_paths: 2380 | corrupted_paths: 0

-- 2. Verify all 11 root categories have non-blank translations
SELECT id, name_en, name_am, name_om, name_ti 
FROM category_hierarchy 
WHERE id BETWEEN 1 AND 11;
-- Expected: 11 rows returned, all fields populated with non-empty Amharic, Oromo, and Tigrinya text.

-- 3. Verify redundant indexes dropped and surviving unique indexes intact
SELECT tablename, indexname 
FROM pg_indexes 
WHERE indexname IN (
    'idx_category_hierarchy_slug_unique', 'idx_category_hierarchy_slug',
    'idx_attribute_definitions_slug', 'idx_users_phone_role', 'idx_brands_slug',
    'uq_category_hierarchy_slug', 'attribute_definitions_slug_key',
    'users_phone_role_key', 'brands_slug_key'
)
ORDER BY tablename, indexname;
-- Expected: 4 rows returned (only the 4 surviving unique indexes; 5 candidate indexes absent).
```

### Rollback Procedures & Limitations

1. **Rollback for `repair_category_paths.sql`:**
   - *Limitation:* The prior path values were mathematically incorrect (e.g. duplicating parent IDs like `/38/38/12/`). Rolling back to corrupted paths re-introduces application bugs.
   - *Procedure (if required for incident recovery):*
     ```sql
     UPDATE category_hierarchy ch
     SET path = b.path, updated_at = NOW()
     FROM category_hierarchy_backup_pre_repair b
     WHERE ch.id = b.id;
     ```

2. **Rollback for `backfill_root_translations.sql`:**
   - *Limitation & Risk:* Blanking the columns unconditionally would overwrite any legitimate manual translations an operator or admin added post-backfill.
   - *Safe Procedure (Reverts only values matching the legacy backfill source):*
     ```sql
     UPDATE category_hierarchy ch
     SET 
         name_am = CASE WHEN ch.name_am = c.name_am THEN NULL ELSE ch.name_am END,
         name_om = CASE WHEN ch.name_om = c.name_or THEN NULL ELSE ch.name_om END,
         name_ti = CASE WHEN ch.name_ti = c.name_ti THEN NULL ELSE ch.name_ti END,
         updated_at = NOW()
     FROM categories c
     WHERE ch.id = c.id AND ch.id BETWEEN 1 AND 11;
     ```

3. **Rollback for Migration 034 (`re-create dropped indexes`):**
   - *Limitation & Risk:* Standard `CREATE INDEX` acquires a `SHARE` lock on the indexed table, blocking concurrent write transactions (`INSERT`, `UPDATE`, `DELETE`) until index build finishes.
   - *Production Procedure (Use `CONCURRENTLY` to avoid write blocking):*
     ```sql
     -- Execute outside transaction block (one at a time)
     CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_category_hierarchy_slug_unique 
         ON public.category_hierarchy USING btree (slug) WHERE (is_active = true);
     CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_category_hierarchy_slug 
         ON public.category_hierarchy USING btree (slug);
     CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_attribute_definitions_slug 
         ON public.attribute_definitions USING btree (slug);
     CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_phone_role 
         ON public.users USING btree (phone, role);
     CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_brands_slug 
         ON public.brands USING btree (slug);
     ```

### Index Lock & Concurrency Management
- `DROP INDEX` acquires an `ACCESS EXCLUSIVE` lock on both the index and its parent table, blocking concurrent `SELECT` and write queries.
- To prevent locking queues from stalling client traffic if an existing read transaction is running, operators should prepend `SET lock_timeout = '2s';` before running migration 034. If lock acquisition exceeds 2 seconds, the command fails fast without holding up incoming traffic.
