-- ============================================================================
-- Migration 019: Query Understanding & Search Intelligence
-- ============================================================================
-- Enables intelligent search with brand/model recognition, synonyms, and analytics
-- 
-- Features:
-- - Brand dictionary for manufacturer recognition
-- - Synonym dictionary for term expansion
-- - Search analytics for query tracking
-- - Query normalization
-- ============================================================================

-- Transaction managed by migration runner

-- ============================================================================
-- PART 1: Brand Dictionary
-- ============================================================================

-- Table: brand_dictionary
-- Stores known brands/manufacturers for recognition in search queries
CREATE TABLE IF NOT EXISTS brand_dictionary (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL UNIQUE,
    normalized_name VARCHAR(200) NOT NULL,  -- Lowercase, no special chars
    aliases JSONB,                          -- Alternative spellings/misspellings
    category_ids INT[],                     -- Applicable categories
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_brand_dictionary_normalized ON brand_dictionary(normalized_name);
CREATE INDEX IF NOT EXISTS idx_brand_dictionary_category_ids ON brand_dictionary USING gin(category_ids);

-- Seed common brands
INSERT INTO brand_dictionary (name, normalized_name, aliases, category_ids)
VALUES 
    -- Electronics brands
    ('Samsung', 'samsung', '["samung", "samsng", "samsungg"]', '{1}'),
    ('Apple', 'apple', '["aple", "appl", "iphone"]', '{1}'),
    ('Huawei', 'huawei', '["huwei", "huawai", "huaweii"]', '{1}'),
    ('Xiaomi', 'xiaomi', '["shaomi", "xiomi", "redmi"]', '{1}'),
    ('Oppo', 'oppo', '["opo", "oppoo"]', '{1}'),
    ('Vivo', 'vivo', '["vivo", "vivoo"]', '{1}'),
    ('Tecno', 'tecno', '["techno", "teknoo"]', '{1}'),
    ('Infinix', 'infinix', '["inifnix", "infinix"]', '{1}'),
    ('Nokia', 'nokia', '["nokiaa", "noki"]', '{1}'),
    ('Realme', 'realme', '["realmi", "realmee"]', '{1}'),
    ('Dell', 'dell', '["del", "dell"]', '{1}'),
    ('HP', 'hp', '["hewlett packard", "hp"]', '{1}'),
    ('Lenovo', 'lenovo', '["lenovo", "lenova"]', '{1}'),
    ('Asus', 'asus', '["asuss", "azus"]', '{1}'),
    ('Acer', 'acer', '["acer", "aser"]', '{1}'),
    
    -- Vehicle brands
    ('Toyota', 'toyota', '["toyata", "toyotta"]', '{2}'),
    ('Nissan', 'nissan', '["nisan", "nissann"]', '{2}'),
    ('Honda', 'honda', '["hond", "hondaa"]', '{2}'),
    ('Hyundai', 'hyundai', '["hundai", "hyunday"]', '{2}'),
    ('Kia', 'kia', '["kia", "kiaa"]', '{2}'),
    ('Suzuki', 'suzuki', '["suzki", "suzukii"]', '{2}'),
    ('Isuzu', 'isuzu', '["izuzu", "issuzu"]', '{2}'),
    ('Mercedes', 'mercedes', '["mercedes benz", "mercedesbenz", "benz"]', '{2}'),
    ('BMW', 'bmw', '["bmw", "beamer"]', '{2}'),
    ('Volkswagen', 'volkswagen', '["vw", "volkswagon", "volkswagen"]', '{2}'),
    ('Ford', 'ford', '["ford", "fordd"]', '{2}'),
    ('Chevrolet', 'chevrolet', '["chevy", "chevrolett"]', '{2}'),
    
    -- Agriculture brands
    ('John Deere', 'johndeere', '["john deer", "johndeere", "deere"]', '{6}'),
    ('Massey Ferguson', 'masseyferguson', '["massey", "ferguson", "mf"]', '{6}'),
    ('New Holland', 'newholland', '["newholland", "nh"]', '{6}'),
    ('Kubota', 'kubota', '["kubota", "kubotta"]', '{6}')
ON CONFLICT (name) DO NOTHING;

-- ============================================================================
-- PART 2: Synonym Dictionary
-- ============================================================================

-- Table: synonym_dictionary
-- Stores term synonyms for query expansion
CREATE TABLE IF NOT EXISTS synonym_dictionary (
    id SERIAL PRIMARY KEY,
    term VARCHAR(200) NOT NULL,
    synonyms JSONB NOT NULL,                -- Array of equivalent terms
    context VARCHAR(100),                   -- Optional: category/domain
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_synonym_dictionary_term ON synonym_dictionary(term);

-- Seed common synonyms
INSERT INTO synonym_dictionary (term, synonyms, context)
VALUES 
    -- Electronics
    ('phone', '["mobile", "cell phone", "smartphone", "cellular", "handset"]', 'electronics'),
    ('mobile', '["phone", "cell phone", "smartphone"]', 'electronics'),
    ('laptop', '["notebook", "computer", "pc"]', 'electronics'),
    ('computer', '["laptop", "desktop", "pc"]', 'electronics'),
    ('tv', '["television", "telly", "screen"]', 'electronics'),
    
    -- Property
    ('house', '["home", "residence", "property", "dwelling"]', 'property'),
    ('apartment', '["flat", "condo", "unit"]', 'property'),
    ('rent', '["rental", "lease", "hire"]', 'property'),
    ('buy', '["purchase", "sale", "buying"]', 'property'),
    
    -- Vehicles
    ('car', '["vehicle", "automobile", "auto"]', 'vehicles'),
    ('bike', '["motorcycle", "motorbike", "cycle"]', 'vehicles'),
    ('truck', '["lorry", "van"]', 'vehicles'),
    
    -- Agriculture
    ('tractor', '["farm equipment", "machinery"]', 'agriculture'),
    ('plough', '["plow", "ploughing"]', 'agriculture'),
    ('ox', '["oxen", "cattle"]', 'agriculture'),
    
    -- General
    ('new', '["brand new", "unused", "fresh"]', null),
    ('used', '["second hand", "pre-owned", "secondhand"]', null),
    ('cheap', '["affordable", "inexpensive", "budget"]', null),
    ('expensive', '["costly", "pricey", "premium"]', null);

-- ============================================================================
-- PART 3: Search Analytics
-- ============================================================================

-- Table: search_analytics
-- Tracks search queries for analysis and improvement
CREATE TABLE IF NOT EXISTS search_analytics (
    id BIGSERIAL PRIMARY KEY,
    query TEXT NOT NULL,
    normalized_query TEXT,
    result_count INT DEFAULT 0,
    category_id INT,
    filters_used JSONB,
    clicked_product_id INT,
    clicked_position INT,
    user_id INT,
    session_id VARCHAR(100),
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_search_analytics_query ON search_analytics(normalized_query);
CREATE INDEX IF NOT EXISTS idx_search_analytics_created ON search_analytics(created_at);
CREATE INDEX IF NOT EXISTS idx_search_analytics_user ON search_analytics(user_id);
CREATE INDEX IF NOT EXISTS idx_search_analytics_clicked ON search_analytics(clicked_product_id);

-- ============================================================================
-- PART 4: Query Normalization Function
-- ============================================================================

-- Function: Normalize query for matching
CREATE OR REPLACE FUNCTION normalize_query(p_query TEXT)
RETURNS TEXT AS $$
DECLARE
    v_normalized TEXT;
BEGIN
    -- Lowercase
    v_normalized := LOWER(TRIM(p_query));
    
    -- Remove special characters except spaces and hyphens
    v_normalized := REGEXP_REPLACE(v_normalized, '[^a-z0-9\s-]', '', 'g');
    
    -- Replace multiple spaces with single space
    v_normalized := REGEXP_REPLACE(v_normalized, '\s+', ' ', 'g');
    
    -- Trim
    v_normalized := TRIM(v_normalized);
    
    RETURN v_normalized;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- ============================================================================
-- PART 5: Brand Detection Function
-- ============================================================================

-- Function: Extract brands from query
CREATE OR REPLACE FUNCTION extract_brands(p_query TEXT)
RETURNS TABLE (
    brand_id INT,
    brand_name VARCHAR,
    match_position INT
) AS $$
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
$$ LANGUAGE plpgsql;

-- ============================================================================
-- PART 6: Synonym Expansion Function
-- ============================================================================

-- Function: Expand query with synonyms
CREATE OR REPLACE FUNCTION expand_with_synonyms(p_query TEXT)
RETURNS TEXT[] AS $$
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
$$ LANGUAGE plpgsql;

-- ============================================================================
-- PART 7: Popular Searches View
-- ============================================================================

-- View: Popular search queries
CREATE OR REPLACE VIEW popular_searches AS
SELECT 
    normalized_query,
    COUNT(*) as search_count,
    AVG(result_count) as avg_results,
    COUNT(clicked_product_id) as clicks,
    CASE 
        WHEN COUNT(*) > 0 THEN 
            (COUNT(clicked_product_id)::FLOAT / COUNT(*)::FLOAT) * 100
        ELSE 0
    END as ctr_percent,
    MAX(created_at) as last_searched
FROM search_analytics
WHERE created_at >= NOW() - INTERVAL '30 days'
GROUP BY normalized_query
HAVING COUNT(*) >= 3
ORDER BY search_count DESC;

-- ============================================================================
-- PART 8: Zero Result Queries View
-- ============================================================================

-- View: Queries that returned no results
CREATE OR REPLACE VIEW zero_result_queries AS
SELECT 
    query,
    normalized_query,
    COUNT(*) as occurrence_count,
    MAX(created_at) as last_attempted
FROM search_analytics
WHERE result_count = 0
  AND created_at >= NOW() - INTERVAL '30 days'
GROUP BY query, normalized_query
ORDER BY occurrence_count DESC;

-- ============================================================================
-- PART 9: Update Triggers
-- ============================================================================

-- Trigger: Update timestamp on brand dictionary change
CREATE OR REPLACE FUNCTION update_brand_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_update_brand_dictionary_timestamp ON brand_dictionary;
CREATE TRIGGER trg_update_brand_dictionary_timestamp
    BEFORE UPDATE ON brand_dictionary
    FOR EACH ROW
    EXECUTE FUNCTION update_brand_timestamp();

-- ============================================================================
-- PART 10: Comments
-- ============================================================================

COMMENT ON TABLE brand_dictionary IS 'Known brands/manufacturers for search recognition';
COMMENT ON TABLE synonym_dictionary IS 'Term synonyms for query expansion';
COMMENT ON TABLE search_analytics IS 'Search query tracking and analytics';
COMMENT ON FUNCTION normalize_query(TEXT) IS 'Normalizes search query for consistent matching';
COMMENT ON FUNCTION extract_brands(TEXT) IS 'Detects brands mentioned in search query';
COMMENT ON FUNCTION expand_with_synonyms(TEXT) IS 'Expands query with synonym variations';

-- ============================================================================
-- Migration 019 Complete ✅
-- ============================================================================
