# HUNTIQ — Multilingual Category Search Architecture

**Project:** HUNTIQ Universal Marketplace  
**Document:** Search Indexing, Subtree Expansion & Multilingual Query Routing  
**Repository Engine:** `backend/internal/repository/product_repo.go`  

---

## 1. Search Query Pipeline

HUNTIQ implements a high-performance PostgreSQL search pipeline supporting English, Amharic, Afaan Oromoo, and Tigrinya:

```
User Search Query (e.g. "ስማርትፎን" or "smartphones")
                   │
                   ▼
┌────────────────────────────────────────────────────────┐
│ Phase 1: Alias & Canonical Category Resolution         │
│  - Query category_aliases for matching terms           │
│  - Match normalized canonical name_en, name_am, etc.   │
│  - Resolve to canonical CategoryID (e.g. 7631)         │
└──────────────────┬─────────────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────────────┐
│ Phase 2: Hierarchical Subtree Expansion                │
│  - Fetch target category path (e.g. /7629/7630/7631/)  │
│  - Query descendants via:                              │
│    WHERE p.category_hierarchy_id IN (                  │
│        SELECT id FROM category_hierarchy               │
│        WHERE path LIKE target_path || '%'              │
│    )                                                   │
└──────────────────┬─────────────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────────────┐
│ Phase 3: Orthogonal Attribute & Text Filtering         │
│  - Apply transaction_mode (e.g. for_sale, for_rent)    │
│  - Apply condition (e.g. new, used_good)               │
│  - Apply title/description full-text search            │
│  - Return ranked, paginated product listings           │
└────────────────────────────────────────────────────────┘
```

---

## 2. Materialized Path Subtree Indexing

To support arbitrary-depth queries with sub-millisecond latency:
- The `path` column is indexed using a B-Tree index:
  ```sql
  CREATE INDEX idx_category_hierarchy_path ON category_hierarchy (path);
  ```
- Subtree matching executes as an indexed range scan (`path LIKE '/7629/%'`), outperforming recursive joins and closure tables on read-heavy marketplace traffic.

---

## 3. Resolving Ambiguity & Polysemy

### A. Parent-Child Overlap Resolution
When an alias is searched (e.g., `phones`), the engine matches the most specific node (ID 7631). Broad browsing on the parent (ID 7629 "Consumer Electronics") naturally includes ID 7631 because ID 7631's path contains `/7629/`.

### B. Compound Term Disambiguation
For words with multiple meanings across domains:
- Fresh hot peppers: Indexed as `"ጥሬ ሚጥሚጣ"` / `"ጥሬ በርበሬ"` under Agriculture Produce (ID 7618).
- Dried berbere spice: Indexed as `"የተፈጨ በርበሬ"` / `"የበርበሬ ዱቄት"` under Traditional Spices (ID 7738).
- Photography DSLR: Indexed as `"የፎቶ ካሜራ"` (ID 6446).
- Security surveillance: Indexed as `"የደህንነት ካሜራ"` / `"ሲሲቲቪ"` (ID 7766).

Unqualified ambiguous terms surface results from both subtrees with ranking based on listing relevance and seller tier.
