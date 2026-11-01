# HUNTIQ — Taxonomy Migration Readiness & Safety Assurance Report

**Project:** HUNTIQ Universal Marketplace  
**Document:** Production Data Safety, Migration Pre-flight & Rollback Plan  
**Migration 036 Status:** ON HOLD (Awaiting Explicit Owner Approval)  

---

## 1. Production Data Preservation Assurance

In strict adherence to **Section 1 (Non-Negotiable Safety Rules)**:
1. **Migration 036 has NOT been applied to production database `hulubedej_db`.**
2. **Zero production categories have been modified, deleted, or overwritten.**
3. **Zero live product foreign-key references have been altered.**
4. **Zero legacy category records or aliases have been deleted.**
5. **No mock or synthetic data has been written to production tables.**

### Current Database State Verification
```sql
SELECT count(*) FROM category_hierarchy; -- Exactly 2,575 rows
SELECT count(*) FROM products;           -- Exactly 14 rows
SELECT count(*) FROM category_aliases;   -- Exactly 1,038 rows
```

---

## 2. Dry-Run Execution Evidence (Isolated Transaction Block)

The migration script [036_unify_taxonomy_and_fix_paths.sql](file:///home/poven/GO%20PROJECT/sellbuy/backend/internal/database/migrations/036_unify_taxonomy_and_fix_paths.sql) was dry-run validated inside a transactional block (`BEGIN; ... ROLLBACK;`).

### Dry-Run Statistics
- **Transaction Outcome:** Successfully executed all DDL, trigger replacements, CTE path repairs, and product re-alignments.
- **Path Repairs:** 29 corrupted legacy paths repaired down to 0.
- **Tigrinya Names Backfilled:** 22 missing translations populated.
- **Products Re-aligned:** 14 listings successfully re-targeted to canonical leaf categories.
- **Rollback:** Transaction rolled back completely. Zero persistent modifications made to the database.

---

## 3. Production Execution Checklist (When Approved)

When the owner grants explicit written approval to execute Migration 036, the following sequence must be observed:

```
[ ] Step 1: Pre-Migration Full Database Snapshot
    pg_dump -U postgres -d hulubedej_db -F c -b -v -f /backup/hulubedej_db_pre_036_$(date +%s).dump

[ ] Step 2: Verify Snapshot Integrity
    pg_restore --list /backup/hulubedej_db_pre_036_*.dump | wc -l

[ ] Step 3: Put Application in Brief Maintenance Mode
    System displays friendly maintenance notice to users.

[ ] Step 4: Execute Migration 036 in Single Transaction
    psql -U postgres -d hulubedej_db -f backend/internal/database/migrations/036_unify_taxonomy_and_fix_paths.sql

[ ] Step 5: Post-Migration Automated Sanity Check
    - Verify 0 orphaned listings: SELECT count(*) FROM products WHERE category_hierarchy_id NOT IN (SELECT id FROM category_hierarchy);
    - Verify 0 path mismatches: SELECT count(*) FROM category_hierarchy WHERE level > 0 AND path NOT LIKE '/' || parent_id || '/%';
    - Verify 14 listings active and accessible on web frontend.

[ ] Step 6: Flush In-Memory Application Caches
    Invalidate category_tree_* redis/in-memory keys via admin handler.

[ ] Step 7: Lift Maintenance Mode & Resume Full Operations
```

---

## 4. Rollback and Recovery Procedure

If any anomaly occurs during execution:
1. **Automatic Transaction Rollback:** Migration 036 is wrapped in a single `BEGIN ... COMMIT;` block. If any statement encounters an error, PostgreSQL rolls back the entire transaction automatically.
2. **Manual Restore (Disaster Recovery):**
   ```bash
   pg_restore -U postgres -d hulubedej_db --clean --if-exists /backup/hulubedej_db_pre_036_*.dump
   ```
