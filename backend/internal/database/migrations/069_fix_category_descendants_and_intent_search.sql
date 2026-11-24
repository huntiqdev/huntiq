-- Migration 069: Fix Category Descendants Function & Hierarchical Intent Search
-- Author: Principal Database Architect & Taxonomy Engineer
-- Description: Corrects get_category_descendants_by_term() to query valid category_aliases
-- and category_hierarchy without legacy nonexistent tables, enabling instant category-wide
-- auto-discovery for all descendant products and subcategories.

CREATE OR REPLACE FUNCTION get_category_descendants_by_term(p_term text)
RETURNS TABLE(id integer, name_en character varying, slug character varying, level integer)
LANGUAGE plpgsql
AS $function$
DECLARE
    v_clean_term text;
BEGIN
    v_clean_term := TRIM(COALESCE(p_term, ''));
    IF v_clean_term = '' THEN
        RETURN;
    END IF;

    RETURN QUERY
    WITH RECURSIVE matched_roots AS (
        SELECT ch.id, ch.parent_id, ch.name_en, ch.slug, ch.level
        FROM category_hierarchy ch
        WHERE ch.is_active = true
          AND (
              ch.name_en ILIKE '%' || v_clean_term || '%'
              OR ch.name_am ILIKE '%' || v_clean_term || '%'
              OR ch.name_om ILIKE '%' || v_clean_term || '%'
              OR ch.name_ti ILIKE '%' || v_clean_term || '%'
              OR ch.slug ILIKE '%' || v_clean_term || '%'
              OR ch.id IN (
                  SELECT ca.category_id
                  FROM category_aliases ca
                  WHERE ca.alias ILIKE '%' || v_clean_term || '%'
                     OR v_clean_term ILIKE '%' || ca.alias || '%'
              )
          )
    ),
    descendants AS (
        SELECT m.id, m.parent_id, m.name_en, m.slug, m.level
        FROM matched_roots m

        UNION

        SELECT ch.id, ch.parent_id, ch.name_en, ch.slug, ch.level
        FROM category_hierarchy ch
        JOIN descendants d ON ch.parent_id = d.id
        WHERE ch.is_active = true
    )
    SELECT DISTINCT d.id, d.name_en, d.slug, d.level
    FROM descendants d
    ORDER BY d.level;
END;
$function$;
