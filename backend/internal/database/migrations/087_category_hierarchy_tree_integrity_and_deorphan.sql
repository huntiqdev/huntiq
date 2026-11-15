-- ============================================================================
-- Migration 087: Category Hierarchy Tree Integrity & De-orphanization
-- ============================================================================
-- PURPOSE:
-- Deactivate all orphaned category subtrees whose parent or ancestor is inactive.
-- Prevents inactive branches from leaking thousands of leaf specifications as
-- top-level categories in UI navigation, search indexers, and drawers.
-- ============================================================================

DO $$
BEGIN
    -- 1. Cascade inactive status across all levels of category_hierarchy (level 1 down to level 7)
    UPDATE category_hierarchy 
    SET is_active = false, updated_at = NOW() 
    WHERE level = 1 
      AND parent_id NOT IN (SELECT id FROM category_hierarchy WHERE level = 0 AND is_active = true) 
      AND is_active = true;

    UPDATE category_hierarchy 
    SET is_active = false, updated_at = NOW() 
    WHERE level = 2 
      AND parent_id NOT IN (SELECT id FROM category_hierarchy WHERE level = 1 AND is_active = true) 
      AND is_active = true;

    UPDATE category_hierarchy 
    SET is_active = false, updated_at = NOW() 
    WHERE level = 3 
      AND parent_id NOT IN (SELECT id FROM category_hierarchy WHERE level = 2 AND is_active = true) 
      AND is_active = true;

    UPDATE category_hierarchy 
    SET is_active = false, updated_at = NOW() 
    WHERE level = 4 
      AND parent_id NOT IN (SELECT id FROM category_hierarchy WHERE level = 3 AND is_active = true) 
      AND is_active = true;

    UPDATE category_hierarchy 
    SET is_active = false, updated_at = NOW() 
    WHERE level = 5 
      AND parent_id NOT IN (SELECT id FROM category_hierarchy WHERE level = 4 AND is_active = true) 
      AND is_active = true;

    UPDATE category_hierarchy 
    SET is_active = false, updated_at = NOW() 
    WHERE level = 6 
      AND parent_id NOT IN (SELECT id FROM category_hierarchy WHERE level = 5 AND is_active = true) 
      AND is_active = true;

    UPDATE category_hierarchy 
    SET is_active = false, updated_at = NOW() 
    WHERE level = 7 
      AND parent_id NOT IN (SELECT id FROM category_hierarchy WHERE level = 6 AND is_active = true) 
      AND is_active = true;

    -- 2. Deactivate orphaned classification_placements under inactive parent placements
    WITH RECURSIVE inact_placements AS (
        SELECT cp.id 
        FROM classification_placements cp 
        JOIN classification_placements p ON cp.parent_placement_id = p.id 
        WHERE cp.is_active = true AND p.is_active = false
        UNION ALL
        SELECT ch.id 
        FROM classification_placements ch 
        JOIN inact_placements i ON ch.parent_placement_id = i.id 
        WHERE ch.is_active = true
    )
    UPDATE classification_placements 
    SET is_active = false, updated_at = NOW() 
    WHERE id IN (SELECT id FROM inact_placements);

END $$;

