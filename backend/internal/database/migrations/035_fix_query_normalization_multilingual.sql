-- ============================================================================
-- Migration 035: Fix Query Normalization for Multilingual & Synchronize Brands
-- ============================================================================
-- Purpose:
-- 1. Update normalize_query(p_query TEXT) to preserve Ethiopic Unicode characters
--    (\u1200-\u137F) and Afaan Oromoo hudhaa apostrophes instead of stripping them.
-- 2. Synchronize brand_dictionary from authoritative brands catalog (migration 031).
-- 3. Establish an automatic trigger ensuring future catalog additions remain synchronized.
-- ============================================================================

-- PART 1: Multilingual-safe query normalization
CREATE OR REPLACE FUNCTION normalize_query(p_query TEXT)
RETURNS TEXT AS $$
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
$$ LANGUAGE plpgsql IMMUTABLE;

-- PART 2: Synchronize brand_dictionary with 72 verified brands from brands catalog
INSERT INTO brand_dictionary (name, normalized_name, aliases, category_ids, is_active)
SELECT 
    b.name,
    normalize_query(b.name),
    COALESCE(
        jsonb_build_array(LOWER(b.slug)),
        '[]'::jsonb
    ),
    b.category_ids,
    b.is_active
FROM brands b
ON CONFLICT (name) DO UPDATE 
SET 
    normalized_name = EXCLUDED.normalized_name,
    category_ids = COALESCE(EXCLUDED.category_ids, brand_dictionary.category_ids),
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- PART 3: Create synchronization trigger for future brands insertions/updates
CREATE OR REPLACE FUNCTION sync_brand_to_dictionary()
RETURNS TRIGGER AS $$
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
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_sync_brand_to_dictionary ON brands;
CREATE TRIGGER trg_sync_brand_to_dictionary
AFTER INSERT OR UPDATE ON brands
FOR EACH ROW EXECUTE FUNCTION sync_brand_to_dictionary();
