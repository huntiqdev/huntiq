-- ============================================================================
-- Migration 093: Reconcile Taxonomy Leaves, Repair Materialized Paths & Isolate Orphans
-- ============================================================================
-- 1. Automates leaf flag reconciliation across category_hierarchy and
--    classification_placements so zero-child nodes become selectable leaves (is_leaf = true)
--    and parent containers are marked non-leaves (is_leaf = false).
-- 2. Fully repairs materialized paths and levels across all branches (including
--    Bajaj Cargo/Passenger and all descendants) to ensure valid radix backtracking.
-- 3. Safely archives orphaned synthetic root categories.
-- ============================================================================

BEGIN;

-- 1. Reconcile is_leaf in category_hierarchy
-- Set is_leaf = true for all active nodes with 0 active children (resolves zombie container dead-ends)
UPDATE category_hierarchy c
SET is_leaf = true, updated_at = NOW()
WHERE c.is_active = true 
  AND c.is_leaf = false
  AND NOT EXISTS (
    SELECT 1 FROM category_hierarchy child 
    WHERE child.parent_id = c.id AND child.is_active = true
  );

-- Set is_leaf = false for all active nodes that DO have active children
UPDATE category_hierarchy c
SET is_leaf = false, updated_at = NOW()
WHERE c.is_active = true 
  AND c.is_leaf = true
  AND EXISTS (
    SELECT 1 FROM category_hierarchy child 
    WHERE child.parent_id = c.id AND child.is_active = true
  );

-- 2. Reconcile is_leaf in classification_placements
UPDATE classification_placements cp
SET is_leaf = true, updated_at = NOW()
WHERE cp.is_active = true 
  AND cp.is_leaf = false
  AND NOT EXISTS (
    SELECT 1 FROM classification_placements child 
    WHERE child.parent_placement_id = cp.id AND child.is_active = true
  );

UPDATE classification_placements cp
SET is_leaf = false, updated_at = NOW()
WHERE cp.is_active = true 
  AND cp.is_leaf = true
  AND EXISTS (
    SELECT 1 FROM classification_placements child 
    WHERE child.parent_placement_id = cp.id AND child.is_active = true
  );

-- 3. Repair parent, level, and radix path for Three-Wheelers / Bajaj branch
UPDATE category_hierarchy
SET parent_id = 3254,
    level = 3,
    path = '/2/16/3254/3255/',
    updated_at = NOW()
WHERE id = 3255;

UPDATE category_hierarchy
SET parent_id = 3254,
    level = 3,
    path = '/2/16/3254/3256/',
    updated_at = NOW()
WHERE id = 3256;

-- Replace old prefix /3197/3254/ with canonical parent prefix /2/16/3254/ for all Bajaj descendants
UPDATE category_hierarchy
SET path = REPLACE(path, '/3197/3254/', '/2/16/3254/'),
    updated_at = NOW()
WHERE path LIKE '/3197/3254/%';

-- 4. Recompute paths and levels layer by layer from roots down to layer 7 to ensure 100% mathematical consistency
DO $$
DECLARE
    v_updated INT := 1;
    v_iter INT := 0;
BEGIN
    WHILE v_updated > 0 AND v_iter < 8 LOOP
        v_iter := v_iter + 1;
        UPDATE category_hierarchy c
        SET path = p.path || c.id || '/',
            level = p.level + 1,
            updated_at = NOW()
        FROM category_hierarchy p
        WHERE c.parent_id = p.id
          AND (c.path IS DISTINCT FROM (p.path || c.id || '/') OR c.level IS DISTINCT FROM (p.level + 1));
        GET DIAGNOSTICS v_updated = ROW_COUNT;
    END LOOP;
END $$;

-- 5. Isolate and archive inactive synthetic orphaned level 0 nodes
UPDATE category_hierarchy
SET is_archived = true, updated_at = NOW()
WHERE parent_id IS NULL AND is_active = false AND (is_archived = false OR is_archived IS NULL);

COMMIT;

