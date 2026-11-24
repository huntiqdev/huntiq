-- Migration 070: Smart Category Intent & Disambiguation Resolution
-- Author: Principal Database Architect, Taxonomy Architect & Search Engineer
-- Description: Completely overhauls resolve_category_by_alias() without legacy table dependencies,
-- adds resolve_category_intent_multisource() for intelligent multi-branch disambiguation (e.g. 'drill'
-- in tools vs construction), and seeds high-frequency Latin/Ge'ez aliases for effortless user discovery.

-- 1. Upgrade resolve_category_by_alias with hierarchy level prioritization & stemming
CREATE OR REPLACE FUNCTION public.resolve_category_by_alias(p_term text)
RETURNS integer
LANGUAGE plpgsql
STABLE
AS $function$
DECLARE
    v_clean text;
    v_singular text;
    v_plural text;
    v_cat_id integer;
BEGIN
    IF p_term IS NULL THEN
        RETURN NULL;
    END IF;

    v_clean := TRIM(LOWER(p_term));
    IF v_clean = '' THEN
        RETURN NULL;
    END IF;

    -- Compute singular/plural variants for English words
    v_singular := v_clean;
    v_plural := v_clean;
    IF v_clean LIKE '%ies' AND length(v_clean) > 4 THEN
        v_singular := substring(v_clean from 1 for length(v_clean) - 3) || 'y';
    ELSIF v_clean LIKE '%es' AND length(v_clean) > 3 THEN
        v_singular := substring(v_clean from 1 for length(v_clean) - 2);
    ELSIF v_clean LIKE '%s' AND length(v_clean) > 2 THEN
        v_singular := substring(v_clean from 1 for length(v_clean) - 1);
    ELSE
        v_plural := v_clean || 's';
    END IF;

    -- 1. Exact match in category_hierarchy (slug or localized names), prioritizing top-level categories
    SELECT c.id INTO v_cat_id
    FROM category_hierarchy c
    WHERE c.is_active = true
      AND (
          LOWER(c.slug) IN (v_clean, v_singular, v_plural)
          OR LOWER(c.name_en) IN (v_clean, v_singular, v_plural)
          OR LOWER(COALESCE(c.name_am, '')) IN (v_clean, v_singular, v_plural)
          OR LOWER(COALESCE(c.name_om, '')) IN (v_clean, v_singular, v_plural)
          OR LOWER(COALESCE(c.name_ti, '')) IN (v_clean, v_singular, v_plural)
      )
    ORDER BY c.level ASC, c.display_order ASC, c.id ASC
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 2. Exact match in category_aliases joined with category_hierarchy, prioritizing lower level
    SELECT a.category_id INTO v_cat_id
    FROM category_aliases a
    JOIN category_hierarchy c ON a.category_id = c.id
    WHERE c.is_active = true
      AND LOWER(a.alias) IN (v_clean, v_singular, v_plural)
    ORDER BY c.level ASC, c.display_order ASC, c.id ASC
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 3. Prefix match in category_hierarchy
    SELECT c.id INTO v_cat_id
    FROM category_hierarchy c
    WHERE c.is_active = true
      AND (
          c.name_en ILIKE v_clean || '%'
          OR c.slug ILIKE v_clean || '%'
          OR COALESCE(c.name_am, '') LIKE p_term || '%'
      )
    ORDER BY c.level ASC, c.display_order ASC, c.id ASC
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 4. Prefix match in category_aliases
    SELECT a.category_id INTO v_cat_id
    FROM category_aliases a
    JOIN category_hierarchy c ON a.category_id = c.id
    WHERE c.is_active = true
      AND a.alias ILIKE v_clean || '%'
    ORDER BY c.level ASC, c.display_order ASC, c.id ASC
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 5. Substring / Word-boundary match in category_hierarchy
    SELECT c.id INTO v_cat_id
    FROM category_hierarchy c
    WHERE c.is_active = true
      AND (
          c.name_en ILIKE '% ' || v_clean || '%'
          OR c.name_en ILIKE '%' || v_clean || '%'
          OR COALESCE(c.name_am, '') LIKE '%' || p_term || '%'
      )
    ORDER BY c.level ASC, c.display_order ASC, c.id ASC
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 6. Substring match in category_aliases
    SELECT a.category_id INTO v_cat_id
    FROM category_aliases a
    JOIN category_hierarchy c ON a.category_id = c.id
    WHERE c.is_active = true
      AND a.alias ILIKE '%' || v_clean || '%'
    ORDER BY c.level ASC, c.display_order ASC, c.id ASC
    LIMIT 1;

    RETURN v_cat_id;
END;
$function$;


-- 2. Multi-Source Intent Disambiguation Function
-- When terms like 'drill' match multiple branches (Tools vs Industrial vs Mining),
-- this function returns all viable candidates ranked by relevancy, depth, and popularity.
DROP FUNCTION IF EXISTS public.resolve_category_intent_multisource(text, integer);
CREATE OR REPLACE FUNCTION public.resolve_category_intent_multisource(p_term text, p_limit int DEFAULT 5)
RETURNS TABLE (
    out_category_id integer,
    out_name_en character varying,
    out_name_am character varying,
    out_slug character varying,
    out_level integer,
    out_path text,
    out_score double precision,
    out_match_source text
)
LANGUAGE plpgsql
STABLE
AS $function$
DECLARE
    v_clean text;
    v_singular text;
    v_plural text;
    v_like text;
BEGIN
    v_clean := TRIM(LOWER(COALESCE(p_term, '')));
    IF v_clean = '' THEN
        RETURN;
    END IF;

    IF p_limit <= 0 OR p_limit > 20 THEN
        p_limit := 5;
    END IF;

    v_singular := v_clean;
    v_plural := v_clean;
    IF v_clean LIKE '%ies' AND length(v_clean) > 4 THEN
        v_singular := substring(v_clean from 1 for length(v_clean) - 3) || 'y';
    ELSIF v_clean LIKE '%es' AND length(v_clean) > 3 THEN
        v_singular := substring(v_clean from 1 for length(v_clean) - 2);
    ELSIF v_clean LIKE '%s' AND length(v_clean) > 2 THEN
        v_singular := substring(v_clean from 1 for length(v_clean) - 1);
    ELSE
        v_plural := v_clean || 's';
    END IF;
    v_like := '%' || v_clean || '%';

    RETURN QUERY
    WITH candidate_matches AS (
        -- Exact name matches (Score: 100)
        SELECT c.id, c.name_en, c.name_am, c.slug, c.level, c.path,
               100.0::double precision as score,
               'exact_category_name'::text as match_source
        FROM category_hierarchy c
        WHERE c.is_active = true
          AND (
              LOWER(c.name_en) IN (v_clean, v_singular, v_plural)
              OR LOWER(COALESCE(c.name_am, '')) IN (v_clean, v_singular, v_plural)
              OR LOWER(c.slug) IN (v_clean, v_singular, v_plural)
          )

        UNION ALL

        -- Exact alias matches (Score: 90)
        SELECT c.id, c.name_en, c.name_am, c.slug, c.level, c.path,
               90.0::double precision as score,
               'exact_alias'::text as match_source
        FROM category_aliases a
        JOIN category_hierarchy c ON a.category_id = c.id
        WHERE c.is_active = true
          AND LOWER(a.alias) IN (v_clean, v_singular, v_plural)

        UNION ALL

        -- Prefix matches (Score: 75)
        SELECT c.id, c.name_en, c.name_am, c.slug, c.level, c.path,
               75.0::double precision as score,
               'prefix_match'::text as match_source
        FROM category_hierarchy c
        WHERE c.is_active = true
          AND (c.name_en ILIKE v_clean || '%' OR c.slug ILIKE v_clean || '%')

        UNION ALL

        -- Substring matches in category names & aliases (Score: 50)
        SELECT c.id, c.name_en, c.name_am, c.slug, c.level, c.path,
               50.0::double precision as score,
               'substring_match'::text as match_source
        FROM category_hierarchy c
        WHERE c.is_active = true
          AND (c.name_en ILIKE v_like OR COALESCE(c.name_am, '') LIKE '%' || p_term || '%')
    ),
    deduped AS (
        SELECT cm.id, cm.name_en, cm.name_am, cm.slug, cm.level, cm.path,
               MIN(cm.match_source) as match_source,
               -- Adjust score: boost shallower levels (roots & category families) for broader discovery
               MAX(cm.score) - (cm.level * 2.0) as effective_score
        FROM candidate_matches cm
        GROUP BY cm.id, cm.name_en, cm.name_am, cm.slug, cm.level, cm.path
    )
    SELECT d.id, d.name_en, d.name_am, d.slug, d.level, d.path::text,
           d.effective_score, d.match_source
    FROM deduped d
    ORDER BY d.effective_score DESC, d.level ASC, d.id ASC
    LIMIT p_limit;
END;
$function$;


-- 3. Seed authoritative aliases for common colloquial terms, Ethiopian Latin/Amharic, and key categories
INSERT INTO category_aliases (category_id, alias, language)
VALUES
    -- Phones (Category 9 & 1425)
    (9, 'phone', 'en'),
    (9, 'phones', 'en'),
    (9, 'mobile phone', 'en'),
    (9, 'mobile phones', 'en'),
    (9, 'cellphone', 'en'),
    (9, 'cellphones', 'en'),
    (9, 'smartphone', 'en'),
    (9, 'smartphones', 'en'),
    (9, 'silk', 'en'),
    (9, 'ስልክ', 'am'),
    (9, 'ስልኮች', 'am'),
    (9, 'ሞባይል', 'am'),
    (9, 'ስማርትፎን', 'am'),
    (9, 'bilbila', 'om'),
    (9, 'ቴሌፎን', 'am'),

    -- Cars / Vehicles (Category 14 & 2)
    (14, 'car', 'en'),
    (14, 'cars', 'en'),
    (14, 'automobile', 'en'),
    (14, 'automobiles', 'en'),
    (14, 'auto', 'en'),
    (14, 'motor car', 'en'),
    (14, 'sedan', 'en'),
    (14, 'sedans', 'en'),
    (14, 'vitz', 'en'),
    (14, 'toyota', 'en'),
    (14, 'mekena', 'en'),
    (14, 'mekina', 'en'),
    (14, 'መኪና', 'am'),
    (14, 'መኪኖች', 'am'),
    (14, 'ተሽከርካሪ', 'am'),
    (14, 'konkolaataa', 'om'),

    -- Electronics (Category 1)
    (1, 'electronics', 'en'),
    (1, 'electronic', 'en'),
    (1, 'electronic devices', 'en'),
    (1, 'ኤሌክትሮኒክስ', 'am'),
    (1, 'ኤሌክትሮኒክ', 'am'),

    -- Construction (Category 1201)
    (1201, 'construction', 'en'),
    (1201, 'construction equipment', 'en'),
    (1201, 'building materials', 'en'),
    (1201, 'cement', 'en'),
    (1201, 'rebar', 'en'),
    (1201, 'ግንባታ', 'am'),
    (1201, 'የግንባታ እቃዎች', 'am'),
    (1201, 'ኮንስትራክሽን', 'am'),
    (1201, 'ijaarsa', 'om'),

    -- Drills & Power Tools (Category 3623 & 2376)
    (3623, 'drill', 'en'),
    (3623, 'drills', 'en'),
    (3623, 'cordless drill', 'en'),
    (3623, 'power drill', 'en'),
    (3623, 'መሰርሰሪያ', 'am'),
    (3623, 'ድሪል', 'am'),
    (2376, 'power drills', 'en'),
    (2376, 'industrial drill', 'en'),
    (2376, 'የኃይል መሰርሰሪያ', 'am'),

    -- Computers & Laptops (Category 10)
    (10, 'computer', 'en'),
    (10, 'computers', 'en'),
    (10, 'laptop', 'en'),
    (10, 'laptops', 'en'),
    (10, 'pc', 'en'),
    (10, 'desktop', 'en'),
    (10, 'ኮምፒውተር', 'am'),
    (10, 'ላፕቶፕ', 'am'),

    -- Real Estate / Property (Category 3)
    (3, 'house', 'en'),
    (3, 'houses', 'en'),
    (3, 'home', 'en'),
    (3, 'homes', 'en'),
    (3, 'apartment', 'en'),
    (3, 'apartments', 'en'),
    (3, 'condo', 'en'),
    (3, 'condominium', 'en'),
    (3, 'rent house', 'en'),
    (3, 'ቤት', 'am'),
    (3, 'ቤቶች', 'am'),
    (3, 'ኪራይ ቤት', 'am'),
    (3, 'ኮንዶሚኒየም', 'am')
ON CONFLICT (category_id, alias) DO NOTHING;
