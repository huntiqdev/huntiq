# HUNTIQ — Taxonomy Administration & Governance Architecture

**Project:** HUNTIQ Universal Marketplace  
**Document:** Administrative Controls, Authorization, Audit Logging & Cache Invalidation  
**Implementation:** Go Fiber Backend (`backend/internal/handler/category_handler.go`)  

---

## 1. Strict Role-Based Access Control (RBAC)

Canonical taxonomy management is an exclusive platform administrative privilege:
- **`super_admin` / `admin`:** Authorized to create categories, update definitions, reorder nodes, move hierarchy positions, merge duplicates, archive nodes, and manage aliases.
- **`moderator`:** Can view taxonomy nodes, review seller listings, and flag compliance violations; **CANNOT modify taxonomy structure**.
- **`seller` / `buyer`:** Strictly prohibited from modifying, creating, or altering category nodes.

Category mutation routes in `backend/internal/handler/category_handler.go` require authentication and admin permission checks:
```go
// Only authenticated admins can mutate categories
adminGroup := app.Group("/api/admin/categories", middleware.RequireAuth(), middleware.RequireRole("admin", "super_admin"))
adminGroup.Post("/", categoryHandler.CreateCategory)
adminGroup.Put("/:id", categoryHandler.UpdateCategory)
adminGroup.Delete("/:id", categoryHandler.DeleteCategory)
adminGroup.Put("/:id/reorder", categoryHandler.ReorderCategory)
adminGroup.Put("/:id/status", categoryHandler.SetCategoryStatus)
```

---

## 2. Immutable Administrative Audit Logging

Every administrative operation that alters taxonomy state records an immutable event in the PostgreSQL `admin_logs` table via `adminSvc.InsertAdminLog`:

```go
h := handlers.NewCategoryHandler(categoryHierarchySvc).
    WithCache(appCache).
    WithAdminLogger(func(ctx context.Context, adminID int, action, notes string) {
        adminSvc.InsertAdminLog(ctx, models.AdminLog{
            AdminID:   adminID,
            Action:    action,
            Notes:     notes,
            IPAddress: "",
        })
    })
```

### Audited Actions & Data Payloads
- `CATEGORY_CREATE`: Logs category ID, name, parent ID, slug, and creator admin ID.
- `CATEGORY_UPDATE`: Logs before-and-after name, definition, and attributes.
- `CATEGORY_MOVE`: Logs old parent ID, new parent ID, old path, new path, and descendant count.
- `CATEGORY_REORDER`: Logs category ID, old display order, and new display order.
- `CATEGORY_STATUS`: Logs activation, deactivation, or archival reason.
- `ALIAS_CREATE` / `ALIAS_DELETE`: Logs category ID, alias term, and language code.

---

## 3. Post-Commit Cache Invalidation

To guarantee that buyers and sellers immediately see updated category trees while preventing cache inconsistency during database rollback:

1. **Transaction Boundary:** Database mutations are committed first.
2. **Post-Commit Invalidation:** Caches are purged **only after** `tx.Commit()` succeeds.
3. **Invalidated Cache Keys:**
   - `category_tree_*` (Localized full trees in memory).
   - `category_nav_*` (Header navigation cache).
   - `category_facets_*` (Search filter facets).

```go
// Invalidate cache only after successful commit
if h.cache != nil {
    h.cache.Delete(ctx, "category_tree_en")
    h.cache.Delete(ctx, "category_tree_am")
    h.cache.Delete(ctx, "category_tree_or")
    h.cache.Delete(ctx, "category_tree_ti")
    h.cache.Delete(ctx, "category_hierarchy_roots")
}
```
