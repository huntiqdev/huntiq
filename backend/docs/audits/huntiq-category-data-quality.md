# HUNTIQ — Comprehensive Category Data Quality Report

**Audited System:** HUNTIQ Universal Marketplace (Ethiopia)  
**Authority:** PostgreSQL 18.4 (Debian runtime), Database `hulubedej_db`  
**Execution Timestamp:** 2026-10-24  
**Audit Standard:** Read-only schema and data inspection, constraint verification, and mathematical hierarchy reconciliation.

---

## 1. Executive Summary & Inventory Statistics

All metrics below were computed using direct, read-only SQL queries against the authoritative PostgreSQL table `public.category_hierarchy`.

| Metric | Measured Value | Validation Status | Notes |
| :--- | :--- | :--- | :--- |
| **Total Categories** | **2,380** | **PASS** | Canonical hierarchical taxonomy |
| **Active Categories** | **2,380** (100.0%) | **PASS** | `is_active = true` across all categories |
| **Inactive / Archived Categories** | **0** (0.0%) | **PASS** | Zero suppressed categories |
| **Root Categories (`parent_id IS NULL`)** | **69** | **PASS** | Level 0 root industry and domain pillars |
| **Branch Categories (`is_leaf = false`)** | **460** | **PASS** | Intermediate parent groupings |
| **Leaf Categories (`is_leaf = true`)** | **1,920** | **PASS** | Terminal selectable product nodes |
| **Minimum Depth (`min_level`)** | **0** | **PASS** | Root domains |
| **Maximum Depth (`max_level`)** | **3** | **PASS** | 4-tier tree (Domain → Subdomain → Family → Leaf) |
| **Legacy `categories` Table Count** | **11** | **LEGACY COMPAT** | Early flat category table retained for backward compat |

```sql
SELECT 
    COUNT(*) as total_categories,
    COUNT(*) FILTER (WHERE parent_id IS NULL) as root_categories,
    COUNT(*) FILTER (WHERE is_leaf = true) as leaf_categories,
    COUNT(*) FILTER (WHERE is_leaf = false) as branch_categories,
    COUNT(*) FILTER (WHERE is_active = true) as active_categories,
    COUNT(*) FILTER (WHERE is_active = false) as inactive_categories,
    MIN(level) as min_level,
    MAX(level) as max_level
FROM category_hierarchy;
```

---

## 2. Structural & Relational Hierarchy Analysis

### A. Integrity Invariant Checks
1. **Orphan Check:** `0` orphaned rows (`parent_id IS NOT NULL AND NOT EXISTS (SELECT 1 FROM category_hierarchy p WHERE p.id = c.parent_id)`). The foreign key constraint `category_hierarchy_parent_id_fkey` strictly enforces referential integrity on delete restrict.
2. **Self-Parenting Check:** `0` rows where `id = parent_id`. Prevented by check constraint `no_circular_reference (id <> parent_id)`.
3. **Level Monotonicity:** `0` level discrepancies. For every child node, `child.level = parent.level + 1` across all 2,380 categories.
4. **Root Depth:** All 69 root categories have `level = 0`. Check constraint `valid_level (level >= 0 AND level < 20)` is satisfied.

### B. Materialized Path (`path`) Inconsistencies
Materialized paths in `category_hierarchy.path` are formatted as `/root_id/sub_id/.../`.
- **Finding:** **2,289 non-root categories** exhibit a duplicated parent segment in their path (e.g. Category 466 under Category 46 has `path = '/46/46/'` instead of `'/46/466/'`).
- **Root Cause:** In [`backend/cmd/seed-taxonomy/main.go:139-144`](file:///home/poven/GO%20PROJECT/sellbuy/backend/cmd/seed-taxonomy/main.go#L139-L144), the batch seed script computed:
  ```go
  childPath := fmt.Sprintf("%s%d/", parentPath, parentID) // Bug: appended parentID instead of childID
  ```
  It then executed an explicit `UPDATE category_hierarchy SET path = $1 WHERE id = $2`. Because this update did not modify `parent_id`, the database trigger `trg_update_category_hierarchy_path` (which listens `BEFORE INSERT OR UPDATE OF parent_id`) was bypassed.
- **Impact:** Tree traversal via `get_category_ancestors` and `get_category_descendants` is **unaffected** because both use recursive CTEs based on `parent_id`. However, functions relying directly on path string parsing (such as `get_category_attributes_inherited`) must be guarded.
- **Remediation:** Apply an idempotent recursive CTE backfill update (detailed in Section 6).

### C. Branch vs. Leaf Inconsistencies
- Marked as leaf (`is_leaf = true`) but has children: **0 categories**.
- Marked as branch (`is_leaf = false`) but has no children: **5 categories**:
  - `ID 7`: "Business & Industry" (`level = 0`, root domain with no current sub-items)
  - `ID 9`: "Lifestyle & Wellness" (`level = 0`, root domain with no current sub-items)
  - `ID 10`: "Education & Learning" (`level = 0`, root domain with no current sub-items)
  - `ID 11`: "Other" (`level = 0`, general catch-all category)
  - `ID 6563`: "Cookware, Pots, Pans and Injera Mitad" (`level = 1`, parent `46`).
- **Explanation:** In migration `021_seed_category_hierarchy.sql`, an explicit statement `UPDATE category_hierarchy SET is_leaf = false WHERE level = 0;` was executed so root categories cannot have direct listings attached without drilling into sub-branches. Category 6563 is a leaf candidate that was marked `is_leaf = false`.

---

## 3. Localization & Media Completeness

| Translation / Asset Field | Populated Records | Missing Records | Completeness |
| :--- | :--- | :--- | :--- |
| **English (`name_en`)** | 2,380 | 0 | **100.0%** |
| **Amharic (`name_am`)** | 2,369 | **11** | **99.5%** |
| **Afaan Oromo (`name_om`)** | 2,369 | **11** | **99.5%** |
| **Tigrinya (`name_ti`)** | 2,347 | **33** | **98.6%** |
| **Slugs (`slug`)** | 2,380 | 0 | **100.0%** |
| **Icons (`icon`)** | 60 | 2,320 | **2.5%** (Root categories only) |
| **Image URLs (`image_url`)** | 0 | 2,380 | **0.0%** (Configured via migration 032) |

### Specific Missing Translations

#### 1. Root Categories Missing Amharic & Oromo (IDs 1 through 11)
During migration `017_category_hierarchy.sql:70-80`, the initial 11 root categories were copied from `categories` into `category_hierarchy`, but the `INSERT` column list omitted `name_am`, `name_or` (as `name_om`), and `name_ti`.
The translations already exist in the legacy `categories` table and can be cleanly synced:

| ID | English (`name_en`) | Legacy Amharic (`name_am`) | Legacy Oromo (`name_or`) | Legacy Tigrinya (`name_ti`) |
| :--- | :--- | :--- | :--- | :--- |
| 1 | Electronics | ኤሌክትሮኒክስ | Elektirooniksii | ኤሌክትሮኒክስ |
| 2 | Vehicles | ተሽከርካሪዎች | Konkolaataa | መኪናታት |
| 3 | Property & Real Estate | ንብረት | Qabeenya | ንብረት |
| 4 | Home & Living | ቤት እና የኑሮ | Mana fi Jireeny | ቤትን ህይወት |
| 5 | Fashion & Clothing | ፋሽን | Faashinii | ፋሽን |
| 6 | Agriculture & Farming | ግብርና | Qonna | ሓረስታይ |
| 7 | Business & Industry | ንግድ | Daldala | ንግዲ |
| 8 | Services & Jobs | አገልግሎቶች | Tajaajila | ኣገልግሎት |
| 9 | Lifestyle & Wellness | የአኗኗር ዘይቤ | Akkaataa Jireenyaa | ኣነባብራ |
| 10 | Education & Learning | ትምህርት | Barnoot | ትምህርቲ |
| 11 | Other | ሌላ | Kan biraa | ካልእ |

#### 2. Subcategories Missing Tigrinya (IDs 14 through 37)
Subcategories 14–37 were seeded in migration `021` before `name_ti` was introduced in migration `031`:
`IDs: 14, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 37`.

---

## 4. Slugs, Duplicates & Cross-Branch Concept Review

### Slug Integrity
- **Total Duplicate Slugs:** **0** (Enforced by PostgreSQL unique constraint `uq_category_hierarchy_slug`).
- **Slug Format:** 100% compliant with lowercase alphanumeric hyphenated pattern (`^[a-z0-9]+(-[a-z0-9]+)*$`).

### Cross-Branch Duplicate Concept Review
Out of 2,380 categories, exactly 4 concept names appear in more than one branch:

| Concept Name (Lower) | Occurrences | Category IDs | Levels | Parent IDs | Analysis & Recommendation |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `refrigerators and freezers` | 2 | `{163, 6509}` | `{1, 1}` | `{44, 44}` | **Duplicate sibling subcategories under Parent 44 ("Home Appliances")!** Category 163 (slug `refrigeration`) has 3 children; Category 6509 (slug `refrigerators-freezers`) has 8 children. **Action:** Human review to merge children into 6509 and deactivate 163. |
| `washing machines and dryers` | 2 | `{167, 6518}` | `{1, 1}` | `{44, 44}` | **Duplicate sibling subcategories under Parent 44 ("Home Appliances")!** Category 167 (slug `washers-dryers`) has 3 children; Category 6518 (slug `washing-machines-dryers`) has 8 children. **Action:** Human review to merge children into 6518 and deactivate 167. |
| `feature phones & keypad devices` | 2 | `{6302, 101}` | `{1, 2}` | `{38, 12}` | Category 6302 is under Root 38 ("Mobile Phones"), while Category 101 is under Subcategory 12 ("Phones"). Legitimate parallel taxonomy in separate domains. |
| `certified environmental & social impact assessment (esia) reports` | 2 | `{2264, 6916}` | `{2, 2}` | `{2261, 6914}` | Parallel professional consulting services in distinct sub-sectors. |

---

## 5. Product Category Referential Integrity

Direct query against `products` table:
```sql
SELECT id, title, category_id, category_hierarchy_id, status FROM products;
```
**Results:**
- **Product ID 3:** "Modern Furnished 2-Bedroom Apartment CMC"
  - `status`: `approved`
  - `category_id`: `3` (Valid: references "Property & Real Estate" in `categories`)
  - `category_hierarchy_id`: `NULL`
  - **Dangling Foreign Keys:** **0** (All product references point to valid records).
  - **Resolution:** Updated `product_repo.go` query includes compatibility fallback `(p.category_hierarchy_id IS NULL AND p.category_id = $N)` ensuring Product 3 matches category 3 filters without data loss.

---

## 6. Safe Remediation Proposals (SQL Scripts)

### Remediation Step 1: Sync Root Category Translations from Legacy Table
This query is non-destructive, safe, and restores the 11 missing root translations:
```sql
UPDATE category_hierarchy ch
SET name_am = c.name_am,
    name_om = c.name_or,
    name_ti = c.name_ti,
    updated_at = NOW()
FROM categories c
WHERE ch.id = c.id
  AND ch.level = 0
  AND (ch.name_am IS NULL OR ch.name_am = '');
```

### Remediation Step 2: Recompute Materialized Paths
This query recalculates the correct materialized paths (`/parent/child/`) without altering IDs, names, or parent relationships:
```sql
WITH RECURSIVE cat_paths AS (
    SELECT id, '/' || CAST(id AS VARCHAR) || '/' AS expected_path, 0 as expected_level
    FROM category_hierarchy
    WHERE parent_id IS NULL
    
    UNION ALL
    
    SELECT c.id, cp.expected_path || CAST(c.id AS VARCHAR) || '/', cp.expected_level + 1
    FROM category_hierarchy c
    JOIN cat_paths cp ON c.parent_id = cp.id
)
UPDATE category_hierarchy ch
SET path = cp.expected_path,
    level = cp.expected_level,
    updated_at = NOW()
FROM cat_paths cp
WHERE ch.id = cp.id
  AND (ch.path != cp.expected_path OR ch.path IS NULL);
```

### Remediation Step 3: Candidate Merge for Home Appliance Sibling Duplicates
1. Move children of category 163 to category 6509:
```sql
UPDATE category_hierarchy SET parent_id = 6509, updated_at = NOW() WHERE parent_id = 163;
```
2. Move children of category 167 to category 6518:
```sql
UPDATE category_hierarchy SET parent_id = 6518, updated_at = NOW() WHERE parent_id = 167;
```
3. Deactivate redundant parent categories:
```sql
UPDATE category_hierarchy SET is_active = false, updated_at = NOW() WHERE id IN (163, 167);
```
