-- ============================================================================
-- Migration 039: Incorporate Category Redirects into Alias Resolution
-- ============================================================================
-- Purpose:
-- Fix alias resolution bypass where deactivated or obsolete slugs (e.g. 'houses-rent')
-- fall through to trigram similarity matching on category_aliases and incorrectly
-- match 'house' (ID 7675) instead of honoring the canonical redirect to ID 7672.
--
-- Resolution Order:
-- 1. Exact match in category_aliases
-- 2. Canonical redirect in category_redirects (e.g. houses-rent -> 7672)
-- 3. Exact slug or multilingual name match in category_hierarchy (is_active = true)
-- 4. Trigram index match on category_aliases (% operator with idx_category_aliases_trgm)
-- ============================================================================

BEGIN;

CREATE OR REPLACE FUNCTION public.resolve_category_by_alias(p_term text)
RETURNS integer AS $$
DECLARE
    v_cat_id INT;
    v_clean_term TEXT;
BEGIN
    v_clean_term := LOWER(TRIM(p_term));
    IF v_clean_term = '' THEN
        RETURN NULL;
    END IF;

    -- 1. Exact match in category_aliases (Index scan on idx_category_aliases_alias)
    SELECT category_id INTO v_cat_id
    FROM category_aliases
    WHERE LOWER(alias) = v_clean_term
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 2. Exact match in category_redirects (Canonical redirects for deactivated/legacy slugs)
    SELECT target_category_id INTO v_cat_id
    FROM category_redirects
    WHERE LOWER(source_slug) = v_clean_term
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 3. Exact slug or name match in category_hierarchy (Active categories only)
    SELECT id INTO v_cat_id
    FROM category_hierarchy
    WHERE is_active = true
      AND (
          LOWER(slug) = v_clean_term
          OR LOWER(name_en) = v_clean_term
          OR LOWER(name_am) = v_clean_term
          OR LOWER(name_om) = v_clean_term
          OR LOWER(name_ti) = v_clean_term
      )
    ORDER BY level DESC
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 4. Trigram index match on category_aliases (% operator uses GIN idx_category_aliases_trgm)
    SELECT category_id INTO v_cat_id
    FROM category_aliases
    WHERE alias % v_clean_term
    ORDER BY similarity(alias, v_clean_term) DESC
    LIMIT 1;

    RETURN v_cat_id;
END;
$$ LANGUAGE plpgsql;

-- Audit Logging
INSERT INTO taxonomy_audit_logs (category_id, action, reason)
VALUES (
    7672,
    'MIGRATION_039_EXECUTE',
    'Upgraded resolve_category_by_alias to prioritize category_redirects before trigram fuzzy matching'
);

COMMIT;
