-- HUNTIQ Taxonomy Path Repair Script (Idempotent)
-- Target: public.category_hierarchy
-- Purpose: Corrects materialized paths corrupted by parentID duplication defect (AUD-001)
-- Safety: Enclosed in transaction with DO-block verification.

BEGIN;

WITH RECURSIVE category_path_cte AS (
    -- Root categories (level 0)
    SELECT 
        id, 
        parent_id, 
        '/' || id || '/' AS expected_path
    FROM category_hierarchy
    WHERE parent_id IS NULL

    UNION ALL

    -- Child categories (level 1..N)
    SELECT 
        c.id, 
        c.parent_id, 
        p.expected_path || c.id || '/' AS expected_path
    FROM category_hierarchy c
    JOIN category_path_cte p ON c.parent_id = p.id
)
UPDATE category_hierarchy ch
SET path = cte.expected_path,
    updated_at = NOW()
FROM category_path_cte cte
WHERE ch.id = cte.id
  AND (ch.path IS DISTINCT FROM cte.expected_path);

-- Verification check
DO $$
DECLARE
    v_mismatch_count INT;
BEGIN
    WITH RECURSIVE check_tree AS (
        SELECT id, parent_id, '/' || id || '/' AS expected_path
        FROM category_hierarchy
        WHERE parent_id IS NULL
        UNION ALL
        SELECT c.id, c.parent_id, ct.expected_path || c.id || '/'
        FROM category_hierarchy c
        JOIN check_tree ct ON c.parent_id = ct.id
    )
    SELECT COUNT(*) INTO v_mismatch_count
    FROM category_hierarchy ch
    JOIN check_tree ct ON ch.id = ct.id
    WHERE ch.path != ct.expected_path;

    IF v_mismatch_count > 0 THEN
        RAISE EXCEPTION 'Path repair verification failed: % mismatches remain', v_mismatch_count;
    END IF;
    RAISE NOTICE 'Path repair verification successful: 0 mismatches remain.';
END $$;

COMMIT;
