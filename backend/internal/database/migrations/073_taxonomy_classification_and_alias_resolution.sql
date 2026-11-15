-- ============================================================================
-- 073_taxonomy_classification_and_alias_resolution.sql
-- Taxonomy resolution, inheritance, and geographic traversal procedures
-- ============================================================================

CREATE OR REPLACE FUNCTION public.get_category_attributes_inherited(p_category_id integer)
 RETURNS TABLE(attribute_id integer, name_en character varying, name_am character varying, name_om character varying, name_ti character varying, slug character varying, data_type character varying, validation_rules jsonb, unit character varying, is_required boolean, is_filterable boolean, is_searchable boolean, is_inherited boolean, display_order integer)
 LANGUAGE plpgsql
AS $function$
DECLARE
    v_path VARCHAR;
    v_ancestor_ids INT[];
BEGIN
    -- Get path of target category
    SELECT path INTO v_path
    FROM category_hierarchy
    WHERE id = p_category_id;
    
    IF v_path IS NOT NULL AND v_path != '' THEN
        -- Convert '/1/5/12/' into integer array [1, 5, 12]
        v_ancestor_ids := string_to_array(trim(both '/' from v_path), '/')::int[];
    ELSE
        v_ancestor_ids := ARRAY[p_category_id];
    END IF;

    RETURN QUERY
    WITH candidate_attributes AS (
        SELECT 
            ad.id AS attribute_id,
            ad.name_en,
            ad.name_am,
            ad.name_om,
            ad.name_ti,
            ad.slug,
            ad.data_type,
            ad.validation_rules,
            ad.unit,
            ca.is_required,
            ad.is_filterable,
            ad.is_searchable,
            (ca.category_id != p_category_id) AS is_inherited,
            ca.display_order,
            ca.category_id AS assigned_category_id,
            ROW_NUMBER() OVER (
                PARTITION BY ad.id 
                ORDER BY 
                    -- Prefer exact category assignment over ancestor assignments
                    CASE WHEN ca.category_id = p_category_id THEN 0 ELSE 1 END,
                    ca.display_order ASC
            ) AS rn
        FROM category_attributes ca
        JOIN attribute_definitions ad ON ad.id = ca.attribute_id
        WHERE ca.category_id = ANY(v_ancestor_ids) 
           OR ca.category_id = p_category_id
    )
    SELECT 
        ca.attribute_id,
        ca.name_en,
        ca.name_am,
        ca.name_om,
        ca.name_ti,
        ca.slug,
        ca.data_type,
        ca.validation_rules,
        ca.unit,
        ca.is_required,
        ca.is_filterable,
        ca.is_searchable,
        ca.is_inherited,
        ca.display_order
    FROM candidate_attributes ca
    WHERE ca.rn = 1
    ORDER BY ca.display_order ASC, ca.attribute_id ASC;
END;
$function$;

CREATE OR REPLACE FUNCTION public.resolve_category_by_alias(p_term text)
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
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
$function$;

CREATE OR REPLACE FUNCTION public.get_category_descendants_by_term(p_term text)
 RETURNS TABLE(id integer)
 LANGUAGE plpgsql
 STABLE
AS $function$
DECLARE
    v_cat_id INT;
BEGIN
    v_cat_id := resolve_category_by_alias(p_term);
    IF v_cat_id IS NOT NULL THEN
        RETURN QUERY SELECT d.id FROM get_category_descendants(v_cat_id) d;
    END IF;
END;
$function$;

DROP FUNCTION IF EXISTS public.get_location_ancestors(integer);
CREATE OR REPLACE FUNCTION public.get_location_ancestors(loc_id integer)
 RETURNS TABLE(id integer, location_type character varying, name_en character varying, name_am character varying, name_om character varying, name_ti character varying, level integer)
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    WITH RECURSIVE ancestors AS (
        SELECT l.id, l.parent_id, l.location_type, l.name_en, l.name_am, l.name_om, l.name_ti, l.level
        FROM location_hierarchy l
        WHERE l.id = loc_id
        
        UNION ALL
        
        SELECT p.id, p.parent_id, p.location_type, p.name_en, p.name_am, p.name_om, p.name_ti, p.level
        FROM location_hierarchy p
        JOIN ancestors a ON a.parent_id = p.id
    )
    SELECT a.id, a.location_type, a.name_en, a.name_am, a.name_om, a.name_ti, a.level
    FROM ancestors a
    ORDER BY a.level ASC;
END;
$function$;

DROP FUNCTION IF EXISTS public.get_location_descendants(integer);
CREATE OR REPLACE FUNCTION public.get_location_descendants(loc_id integer)
 RETURNS TABLE(id integer)
 LANGUAGE plpgsql
AS $function$
DECLARE
    target_path VARCHAR(500);
BEGIN
    SELECT path INTO target_path FROM location_hierarchy WHERE location_hierarchy.id = loc_id;
    RETURN QUERY
    SELECT l.id
    FROM location_hierarchy l
    WHERE l.path LIKE target_path || '%';
END;
$function$;
