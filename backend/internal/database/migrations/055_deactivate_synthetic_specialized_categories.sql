-- ============================================================================
-- 099_deactivate_synthetic_specialized_categories.sql
-- PURPOSE: Deactivate all "Specialized & Custom" and "Specialty" synthetic
--          categories that were auto-generated and pollute the category tree.
--          These categories show as long, confusing names in the UI.
--
-- SAFE: Only deactivates (is_active = false). Does NOT delete any rows.
--       Existing listing references are preserved.
--       Reversible: set is_active = true to restore any category.
-- ============================================================================

BEGIN;

-- 1. Deactivate all "Specialized & Custom" synthetic categories
UPDATE category_hierarchy
SET
    is_active = false,
    is_archived = true,
    updated_at = NOW()
WHERE
    is_active = true
    AND (
        name_en ILIKE 'Specialized & Custom%'
        OR name_en ILIKE 'Specialized %'
        OR name_en ILIKE 'Specialty %'
        OR slug LIKE 'specialized-custom-%'
        OR slug LIKE 'specialized-%'
        OR slug = 'specialized-custom-miscellaneous'
    );

-- 2. Also deactivate categories with excessively long names (>60 chars) at root level (level=0)
--    These are old legacy synthetic roots from before the consolidation
UPDATE category_hierarchy
SET
    is_active = false,
    is_archived = true,
    updated_at = NOW()
WHERE
    is_active = true
    AND level = 0
    AND LENGTH(name_en) > 60;

-- 3. Report what was deactivated
DO $$
DECLARE
    deactivated_count INTEGER;
    remaining_active INTEGER;
    remaining_roots INTEGER;
BEGIN
    SELECT COUNT(*) INTO deactivated_count
    FROM category_hierarchy
    WHERE is_archived = true
      AND updated_at > NOW() - INTERVAL '10 seconds';

    SELECT COUNT(*) INTO remaining_active
    FROM category_hierarchy
    WHERE is_active = true;

    SELECT COUNT(*) INTO remaining_roots
    FROM category_hierarchy
    WHERE is_active = true AND level = 0;

    RAISE NOTICE 'Deactivated % synthetic categories', deactivated_count;
    RAISE NOTICE 'Remaining active categories: %', remaining_active;
    RAISE NOTICE 'Remaining active roots: %', remaining_roots;
END $$;

COMMIT;
