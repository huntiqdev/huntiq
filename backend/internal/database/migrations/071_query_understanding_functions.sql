-- ============================================================================
-- 071_query_understanding_functions.sql
-- Text normalization, brand extraction, synonym expansion, and search analytics views
-- ============================================================================

CREATE OR REPLACE FUNCTION public.normalize_query(p_query text)
 RETURNS text
 LANGUAGE plpgsql
 IMMUTABLE
AS $function$
DECLARE
    v_normalized TEXT;
BEGIN
    IF p_query IS NULL THEN
        RETURN '';
    END IF;

    -- 1. Lowercase
    v_normalized := LOWER(TRIM(p_query));
    
    -- 2. Preserve ASCII alphanumeric, Ethiopic Unicode (\u1200-\u137F), 
    --    apostrophes (Afaan Oromoo hudhaa), spaces, and hyphens.
    v_normalized := REGEXP_REPLACE(v_normalized, '[^a-z0-9\u1200-\u137F\s\''-]', '', 'g');
    
    -- 3. Consolidate multiple spaces
    v_normalized := REGEXP_REPLACE(v_normalized, '\s+', ' ', 'g');
    
    -- 4. Trim
    RETURN TRIM(v_normalized);
END;
$function$;

CREATE OR REPLACE FUNCTION public.extract_brands(p_query text)
 RETURNS TABLE(brand_id integer, brand_name character varying, match_position integer)
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    SELECT 
        bd.id,
        bd.name,
        POSITION(bd.normalized_name IN normalize_query(p_query))::INT AS match_position
    FROM brand_dictionary bd
    WHERE bd.is_active = true
      AND (
          normalize_query(p_query) LIKE '%' || bd.normalized_name || '%'
          OR EXISTS (
              SELECT 1 
              FROM jsonb_array_elements_text(bd.aliases) AS alias
              WHERE normalize_query(p_query) LIKE '%' || alias || '%'
          )
      )
    ORDER BY LENGTH(bd.normalized_name) DESC;
END;
$function$;

CREATE OR REPLACE FUNCTION public.expand_with_synonyms(p_query text)
 RETURNS text[]
 LANGUAGE plpgsql
AS $function$
DECLARE
    v_expanded TEXT[];
    v_word TEXT;
    v_synonyms JSONB;
BEGIN
    v_expanded := ARRAY[p_query];
    
    -- Split query into words
    FOR v_word IN 
        SELECT unnest(string_to_array(normalize_query(p_query), ' '))
    LOOP
        -- Find synonyms for this word
        SELECT synonyms INTO v_synonyms
        FROM synonym_dictionary
        WHERE term = v_word AND is_active = true
        LIMIT 1;
        
        IF v_synonyms IS NOT NULL THEN
            -- Add synonym variants
            v_expanded := v_expanded || ARRAY(
                SELECT jsonb_array_elements_text(v_synonyms)
            );
        END IF;
    END LOOP;
    
    RETURN v_expanded;
END;
$function$;

CREATE OR REPLACE FUNCTION public.update_brand_timestamp()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.sync_brand_to_dictionary()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    INSERT INTO brand_dictionary (name, normalized_name, aliases, category_ids, is_active)
    VALUES (
        NEW.name,
        normalize_query(NEW.name),
        jsonb_build_array(LOWER(NEW.slug)),
        NEW.category_ids,
        NEW.is_active
    )
    ON CONFLICT (name) DO UPDATE 
    SET 
        normalized_name = EXCLUDED.normalized_name,
        category_ids = EXCLUDED.category_ids,
        is_active = EXCLUDED.is_active,
        updated_at = NOW();
    RETURN NEW;
END;
$function$;

DROP TRIGGER IF EXISTS trg_update_brand_dictionary_timestamp ON brand_dictionary;
CREATE TRIGGER trg_update_brand_dictionary_timestamp BEFORE UPDATE ON brand_dictionary FOR EACH ROW EXECUTE FUNCTION update_brand_timestamp();

DROP TRIGGER IF EXISTS trg_sync_brand_to_dictionary ON brands;
CREATE TRIGGER trg_sync_brand_to_dictionary AFTER INSERT OR UPDATE ON brands FOR EACH ROW EXECUTE FUNCTION sync_brand_to_dictionary();

CREATE OR REPLACE VIEW popular_searches AS
SELECT normalized_query,
    count(*) AS search_count,
    avg(result_count) AS avg_results,
    count(clicked_product_id) AS clicks,
        CASE
            WHEN (count(*) > 0) THEN (((count(clicked_product_id))::double precision / (count(*))::double precision) * (100)::double precision)
            ELSE (0)::double precision
        END AS ctr_percent,
    max(created_at) AS last_searched
   FROM search_analytics
  WHERE (created_at >= (now() - '30 days'::interval))
  GROUP BY normalized_query
 HAVING (count(*) >= 3)
  ORDER BY (count(*)) DESC;

CREATE OR REPLACE VIEW zero_result_queries AS
SELECT query,
    normalized_query,
    count(*) AS occurrence_count,
    max(created_at) AS last_attempted
   FROM search_analytics
  WHERE ((result_count = 0) AND (created_at >= (now() - '30 days'::interval)))
  GROUP BY query, normalized_query
  ORDER BY (count(*)) DESC;
