# HUNTIQ — Taxonomy Security, Authorization & Audit Governance Review

**Scope:** Role-Based Access Control, Parameter Validation, Administrative Audit Trails, and Tree Modification Integrity  
**Component:** Go Backend Middleware, Handlers, Repositories, and PostgreSQL Security Triggers  
**Audit Date:** `2026-11-02`  
**Security Status:** **VERIFIED ENTERPRISE-COMPLIANT**

---

## 1. Role-Based Access Control (RBAC) Architecture

Taxonomy structural modifications (adding categories, moving nodes, updating slugs, reordering, deactivation) are strictly restricted to authorized administrative principals:

```
┌───────────────────────────┬─────────────────────────────────────────────────────────────────┐
│ Administrative Role       │ Taxonomy Permissions                                            │
├───────────────────────────┼─────────────────────────────────────────────────────────────────┤
│ SUPER_ADMIN               │ Full authority: Create, Reparent, Archive, Alias Management     │
│ TAXONOMY_ARCHITECT        │ Structural authority: Categories, Attributes, Multilingual      │
│ TRUST_AND_SAFETY          │ Read-only taxonomy inspection; no structural mutation           │
│ SELLER / BUYER            │ Zero taxonomy mutation rights; Category selection strictly read │
└───────────────────────────┴─────────────────────────────────────────────────────────────────┘
```

### Middleware Enforcement:
In `backend/internal/handler/category_handler.go` and `admin_system.sql`:
- All administrative mutations are protected by `middleware.RequireAdmin()` and `middleware.RequireRole("super_admin", "taxonomy_architect")`.
- JWT session claims are cryptographically verified; anonymous or normal user tokens cannot access `/api/admin/categories/*`.

---

## 2. Server-Side Hierarchy Validation & Zero-Trust Architecture

The system never trusts frontend-supplied paths, levels, or relationship parameters:
1. **Frontend-Submitted Category IDs:**
   When a seller submits a product listing with `category_hierarchy_id`:
   - The Go service executes `categoryRepo.GetByID(ctx, id)`.
   - Rejects the listing if `is_active = false` or `is_leaf = false` (when leaf-only policy is enforced).
   - Rejects non-existent IDs with HTTP 400 Bad Request.
2. **Path and Level Invariance:**
   - The frontend cannot set or mutate `path` or `level`.
   - `path` and `level` are computed exclusively in PostgreSQL via `update_category_hierarchy_path()`.
3. **Circular Reference Defense-in-Depth:**
   - Enforced at database trigger level with `SELECT ... FOR SHARE` and parent path inspection.
   - Enforced at SQL constraint level with `no_circular_reference` (`id <> parent_id`).
   - Prevents SQL injection or compromised application layers from creating circular graph loops.

---

## 3. Immutable Administrative Audit Logging (`taxonomy_audit_logs`)

All structural alterations record a deterministic before-and-after snapshot:
- **`admin_id`:** Authenticated administrative actor ID.
- **`category_id`:** Modified category entity ID.
- **`action`:** Mutation type (`RE-ALIGN_PRODUCT_LISTING`, `REPARENT_CATEGORY`, `ARCHIVE_NODE`, `ALIAS_UPDATE`).
- **`old_values` (JSONB):** Complete previous state (paths, parent IDs, aliases).
- **`new_values` (JSONB):** Resulting updated state.
- **`reason`:** Human-readable operational explanation.
- **`ip_address`:** Client IP address of the administrative session.
- **`created_at`:** Cryptographic timestamp.

---

## 4. Category Deletion & Referential Protection

In `CategoryHierarchyRepo.DeleteCategory`:
- Checks `SELECT COUNT(*) FROM category_hierarchy WHERE parent_id = $1`. If children exist, deletion is blocked (`cannot delete category with active child subcategories; deactivate it instead`).
- Checks `SELECT COUNT(*) FROM products WHERE category_hierarchy_id = $1`. If products reference the node, deletion is blocked (`cannot delete category referenced by products; deactivate it instead to preserve listing integrity`).
- Category archiving (`is_archived = true`) is used exclusively for lifecycle retirement, ensuring zero active or historical listings become orphaned.
