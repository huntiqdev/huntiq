-- ============================================================================
-- Migration 035: Ultra-Enterprise Global Taxonomy Architecture & Alias Search
-- ============================================================================
-- Purpose:
-- 1. Ensure pg_trgm extension is active for fast fuzzy & trigram alias searches.
-- 2. Create optimized indexes on category_aliases (trigram and language composite).
-- 3. Create function resolve_category_by_alias() to resolve multilingual search
--    terms (e.g., 'ስልክ', 'bilbila', 'onion', 'ቀይ ሽንኩርት') directly to canonical IDs.
-- 4. Create function get_category_descendants_by_term() for search query expansion.
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Index on category_aliases for fast exact and prefix matching
CREATE INDEX IF NOT EXISTS idx_category_aliases_lang_alias 
    ON category_aliases(language, alias);

-- Trigram index on category_aliases for typo-tolerant fuzzy matching
CREATE INDEX IF NOT EXISTS idx_category_aliases_trgm 
    ON category_aliases USING gin(alias gin_trgm_ops);

-- ============================================================================
-- FUNCTION: resolve_category_by_alias
-- Resolves any search token to its primary canonical category ID.
-- ============================================================================
CREATE OR REPLACE FUNCTION resolve_category_by_alias(p_term TEXT)
RETURNS INT AS $$
DECLARE
    v_cat_id INT;
    v_clean_term TEXT;
BEGIN
    v_clean_term := LOWER(TRIM(p_term));
    IF v_clean_term = '' THEN
        RETURN NULL;
    END IF;

    -- 1. Check exact match in category_aliases
    SELECT category_id INTO v_cat_id
    FROM category_aliases
    WHERE LOWER(alias) = v_clean_term
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 2. Check exact slug or name matches in category_hierarchy
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
    ORDER BY level DESC  -- prefer most specific leaf
    LIMIT 1;

    RETURN v_cat_id;
END;
$$ LANGUAGE plpgsql STABLE;

-- ============================================================================
-- FUNCTION: get_category_descendants_by_term
-- Returns all descendant category IDs for any search term resolving to a category.
-- ============================================================================
CREATE OR REPLACE FUNCTION get_category_descendants_by_term(p_term TEXT)
RETURNS TABLE(id INT) AS $$
DECLARE
    v_cat_id INT;
BEGIN
    v_cat_id := resolve_category_by_alias(p_term);
    IF v_cat_id IS NOT NULL THEN
        RETURN QUERY SELECT d.id FROM get_category_descendants(v_cat_id) d;
    END IF;
END;
$$ LANGUAGE plpgsql STABLE;

-- ============================================================================
-- Migration 035 Complete ✅
-- ============================================================================
