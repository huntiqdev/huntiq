-- Migration 071: Universal Entity Intelligence Engine
-- Author: Principal Knowledge-Graph Architect, Ontology Engineer, Multilingual NLP & Database Architect
-- Core Principle: ONE CANONICAL ENTITY IDENTITY -> ONE MASTER SEMANTIC DEFINITION -> MANY VALID RELATIONSHIPS -> MANY AUTOMATIC DISCOVERY PATHS -> NO UNNECESSARY DUPLICATION.

-- 1. Ensure required extensions exist for sub-millisecond offline string, trigram & phonetic matching
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;
CREATE EXTENSION IF NOT EXISTS unaccent;

-- 2. Core Authoritative Canonical Entity Registry
CREATE TABLE IF NOT EXISTS canonical_entities (
    id SERIAL PRIMARY KEY,
    canonical_name_en VARCHAR(255) NOT NULL,
    canonical_name_am VARCHAR(255),
    canonical_name_om VARCHAR(255),
    canonical_name_ti VARCHAR(255),
    canonical_slug VARCHAR(255) NOT NULL UNIQUE,
    entity_nature VARCHAR(50) NOT NULL CHECK (entity_nature IN ('MATERIAL', 'NON_MATERIAL')),
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN (
        'PRODUCT', 'TOOL_EQUIPMENT', 'VEHICLE', 'PROPERTY', 
        'SERVICE', 'PROFESSION', 'SKILL', 'RAW_MATERIAL', 'BUSINESS_OPPORTUNITY'
    )),
    short_definition_en TEXT NOT NULL,
    short_definition_am TEXT,
    short_definition_om TEXT,
    short_definition_ti TEXT,
    detailed_definition_en TEXT,
    detailed_definition_am TEXT,
    typical_examples JSONB DEFAULT '[]'::jsonb,
    typical_use_cases JSONB DEFAULT '[]'::jsonb,
    functional_capabilities JSONB DEFAULT '[]'::jsonb,
    relevant_attributes JSONB DEFAULT '{}'::jsonb,
    disambiguation_rules JSONB DEFAULT '{}'::jsonb,
    primary_category_id INT REFERENCES category_hierarchy(id) ON DELETE SET NULL,
    parent_entity_id INT REFERENCES canonical_entities(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'PENDING_REVIEW', 'ARCHIVED')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for high-speed canonical entity discovery
CREATE INDEX IF NOT EXISTS idx_canonical_entities_slug ON canonical_entities(canonical_slug);
CREATE INDEX IF NOT EXISTS idx_canonical_entities_type ON canonical_entities(entity_type);
CREATE INDEX IF NOT EXISTS idx_canonical_entities_nature ON canonical_entities(entity_nature);
CREATE INDEX IF NOT EXISTS idx_canonical_entities_primary_cat ON canonical_entities(primary_category_id);
CREATE INDEX IF NOT EXISTS idx_canonical_entities_parent ON canonical_entities(parent_entity_id);
CREATE INDEX IF NOT EXISTS idx_canonical_entities_name_en_trgm ON canonical_entities USING gin(canonical_name_en gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_canonical_entities_slug_trgm ON canonical_entities USING gin(canonical_slug gin_trgm_ops);

-- 3. Master Entity Aliases and Spelling / Phonetic Registry
CREATE TABLE IF NOT EXISTS entity_aliases (
    id SERIAL PRIMARY KEY,
    entity_id INT NOT NULL REFERENCES canonical_entities(id) ON DELETE CASCADE,
    alias_text VARCHAR(255) NOT NULL,
    normalized_alias VARCHAR(255) NOT NULL,
    language VARCHAR(10) NOT NULL DEFAULT 'en',
    alias_type VARCHAR(50) NOT NULL CHECK (alias_type IN (
        'EXACT_SYNONYM', 'COMMON_NAME', 'INFORMAL', 'ABBREVIATION', 
        'ACRONYM', 'REGIONAL', 'TYPO_CORRECTION', 'PHONETIC', 
        'TRANSLITERATION', 'BRAND_REFERENCE', 'MODEL_SPECIFIC'
    )),
    confidence_score NUMERIC(3,2) DEFAULT 1.00 CHECK (confidence_score >= 0.00 AND confidence_score <= 1.00),
    status VARCHAR(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'DEPRECATED')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_entity_alias_unique UNIQUE(entity_id, normalized_alias, language)
);

CREATE INDEX IF NOT EXISTS idx_entity_aliases_entity_id ON entity_aliases(entity_id);
CREATE INDEX IF NOT EXISTS idx_entity_aliases_norm ON entity_aliases(normalized_alias);
CREATE INDEX IF NOT EXISTS idx_entity_aliases_lang ON entity_aliases(language);
CREATE INDEX IF NOT EXISTS idx_entity_aliases_norm_trgm ON entity_aliases USING gin(normalized_alias gin_trgm_ops);

-- 4. Typed Entity Relationship Graph
CREATE TABLE IF NOT EXISTS entity_relationships (
    id SERIAL PRIMARY KEY,
    source_entity_id INT NOT NULL REFERENCES canonical_entities(id) ON DELETE CASCADE,
    target_entity_id INT REFERENCES canonical_entities(id) ON DELETE CASCADE,
    target_category_id INT REFERENCES category_hierarchy(id) ON DELETE CASCADE,
    relationship_type VARCHAR(50) NOT NULL CHECK (relationship_type IN (
        'IS_A', 'INSTANCE_OF', 'HAS_BRAND', 'HAS_FEATURE', 
        'USED_FOR', 'USED_IN_INDUSTRY', 'DISCOVERABLE_THROUGH', 
        'COMPLEMENTS', 'REQUIRES_SKILL', 'OFFERED_AS_SERVICE', 
        'SIMILAR_TO', 'CONTRASTED_WITH'
    )),
    confidence_score NUMERIC(3,2) DEFAULT 1.00 CHECK (confidence_score >= 0.00 AND confidence_score <= 1.00),
    metadata JSONB DEFAULT '{}'::jsonb,
    notes TEXT,
    status VARCHAR(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'DEPRECATED')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_relationship_target CHECK (
        (target_entity_id IS NOT NULL AND target_category_id IS NULL) OR
        (target_entity_id IS NULL AND target_category_id IS NOT NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_entity_rel_source ON entity_relationships(source_entity_id);
CREATE INDEX IF NOT EXISTS idx_entity_rel_target_entity ON entity_relationships(target_entity_id);
CREATE INDEX IF NOT EXISTS idx_entity_rel_target_cat ON entity_relationships(target_category_id);
CREATE INDEX IF NOT EXISTS idx_entity_rel_type ON entity_relationships(relationship_type);
CREATE UNIQUE INDEX IF NOT EXISTS uq_entity_rel_dedup ON entity_relationships (source_entity_id, COALESCE(target_entity_id, 0), COALESCE(target_category_id, 0), relationship_type);

-- 5. Link Products to Canonical Entity
ALTER TABLE products ADD COLUMN IF NOT EXISTS canonical_entity_id INT REFERENCES canonical_entities(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_products_canonical_entity_id ON products(canonical_entity_id);

-- 6. High-Performance Text Normalizer for Multilingual Marketplace Queries
CREATE OR REPLACE FUNCTION public.fn_normalize_entity_text(p_text text)
RETURNS text
LANGUAGE plpgsql
IMMUTABLE
AS $function$
DECLARE
    v_norm text;
BEGIN
    IF p_text IS NULL THEN
        RETURN '';
    END IF;
    -- Lowercase & trim
    v_norm := TRIM(LOWER(p_text));
    -- Collapse multiple whitespace sequences into single space
    v_norm := regexp_replace(v_norm, '\s+', ' ', 'g');
    -- Remove non-alphanumeric punctuation except Ethiopic syllabics (\u1200-\u137F)
    v_norm := regexp_replace(v_norm, '[^[:alnum:]\s\u1200-\u137F]', '', 'g');
    RETURN TRIM(v_norm);
END;
$function$;

-- 7. Universal Canonical Entity Resolver (Multi-stage Typo, Trigram, Phonetic & Levenshtein)
CREATE OR REPLACE FUNCTION public.resolve_canonical_entity(
    p_term text,
    p_lang text DEFAULT 'en',
    p_limit int DEFAULT 5
)
RETURNS TABLE (
    out_entity_id integer,
    out_canonical_name_en character varying,
    out_canonical_name_am character varying,
    out_canonical_slug character varying,
    out_entity_nature character varying,
    out_entity_type character varying,
    out_short_definition_en text,
    out_short_definition_am text,
    out_primary_category_id integer,
    out_score double precision,
    out_confidence_tier character varying,
    out_match_method character varying,
    out_suggested_correction character varying
)
LANGUAGE plpgsql
STABLE
AS $function$
DECLARE
    v_norm text;
    v_singular text;
    v_plural text;
    v_stem text;
    v_stem_alt text;
    v_dmeta text;
BEGIN
    v_norm := public.fn_normalize_entity_text(p_term);
    IF v_norm = '' THEN
        RETURN;
    END IF;

    IF p_limit <= 0 OR p_limit > 20 THEN
        p_limit := 5;
    END IF;

    -- Compute singular/plural variants for English / Latin queries
    v_singular := v_norm;
    v_plural := v_norm;
    IF v_norm LIKE '%ies' AND length(v_norm) > 4 THEN
        v_singular := substring(v_norm from 1 for length(v_norm) - 3) || 'y';
    ELSIF v_norm LIKE '%es' AND length(v_norm) > 3 THEN
        v_singular := substring(v_norm from 1 for length(v_norm) - 2);
    ELSIF v_norm LIKE '%s' AND length(v_norm) > 2 THEN
        v_singular := substring(v_norm from 1 for length(v_norm) - 1);
    ELSE
        v_plural := v_norm || 's';
    END IF;

    -- Compute stem variants (e.g. drilling -> drill, dril; cleaner -> clean)
    v_stem := v_norm;
    v_stem_alt := v_norm;
    IF v_norm LIKE '%ing' AND length(v_norm) > 5 THEN
        v_stem_alt := substring(v_norm from 1 for length(v_norm) - 3);
        IF substring(v_norm from length(v_norm)-4 for 1) = substring(v_norm from length(v_norm)-3 for 1) THEN
            v_stem := substring(v_norm from 1 for length(v_norm) - 4);
        ELSE
            v_stem := v_stem_alt;
        END IF;
    ELSIF v_norm LIKE '%ed' AND length(v_norm) > 4 THEN
        v_stem := substring(v_norm from 1 for length(v_norm) - 2);
    ELSIF v_norm LIKE '%ers' AND length(v_norm) > 5 THEN
        v_stem := substring(v_norm from 1 for length(v_norm) - 3);
    ELSIF v_norm LIKE '%er' AND length(v_norm) > 4 THEN
        v_stem := substring(v_norm from 1 for length(v_norm) - 2);
    END IF;

    -- Compute double metaphone for Latin terms where applicable
    BEGIN
        v_dmeta := dmetaphone(v_norm);
    EXCEPTION WHEN OTHERS THEN
        v_dmeta := '';
    END;

    RETURN QUERY
    WITH candidate_matches AS (
        -- Stage 1: Exact Match on Canonical Name or Slug (Score: 100.0, Tier: HIGH)
        SELECT 
            e.id AS entity_id,
            100.0::double precision AS match_score,
            'HIGH'::varchar AS conf_tier,
            'EXACT_CANONICAL'::varchar AS method,
            e.canonical_name_en AS correction
        FROM canonical_entities e
        WHERE e.status = 'ACTIVE'
          AND (
              LOWER(e.canonical_slug) IN (v_norm, v_singular, v_plural, v_stem, v_stem_alt)
              OR LOWER(e.canonical_name_en) IN (v_norm, v_singular, v_plural, v_stem, v_stem_alt)
              OR LOWER(COALESCE(e.canonical_name_am, '')) IN (v_norm, v_singular, v_plural)
              OR LOWER(COALESCE(e.canonical_name_om, '')) IN (v_norm, v_singular, v_plural)
              OR LOWER(COALESCE(e.canonical_name_ti, '')) IN (v_norm, v_singular, v_plural)
          )

        UNION ALL

        -- Stage 2: Exact Match on Entity Aliases (Score: 95.0 * alias confidence, Tier: HIGH)
        SELECT 
            a.entity_id,
            (95.0 * a.confidence_score)::double precision AS match_score,
            'HIGH'::varchar AS conf_tier,
            'EXACT_ALIAS'::varchar AS method,
            e.canonical_name_en AS correction
        FROM entity_aliases a
        JOIN canonical_entities e ON a.entity_id = e.id
        WHERE a.status = 'ACTIVE' AND e.status = 'ACTIVE'
          AND a.normalized_alias IN (v_norm, v_singular, v_plural, v_stem, v_stem_alt)

        UNION ALL

        -- Stage 3: High-Confidence Trigram Similarity on Aliases (similarity >= 0.60)
        SELECT 
            a.entity_id,
            (similarity(a.normalized_alias, v_norm) * 90.0)::double precision AS match_score,
            CASE WHEN similarity(a.normalized_alias, v_norm) >= 0.80 THEN 'HIGH' ELSE 'MEDIUM' END::varchar AS conf_tier,
            'TRIGRAM_SIMILARITY'::varchar AS method,
            e.canonical_name_en AS correction
        FROM entity_aliases a
        JOIN canonical_entities e ON a.entity_id = e.id
        WHERE a.status = 'ACTIVE' AND e.status = 'ACTIVE'
          AND (similarity(a.normalized_alias, v_norm) >= 0.60 OR similarity(a.normalized_alias, v_stem_alt) >= 0.60)

        UNION ALL

        -- Stage 4: Levenshtein Distance Match (length > 3, distance <= 2)
        SELECT 
            a.entity_id,
            ((1.0 - (levenshtein(v_norm, a.normalized_alias)::float / GREATEST(length(v_norm), length(a.normalized_alias)))) * 85.0)::double precision AS match_score,
            'MEDIUM'::varchar AS conf_tier,
            'LEVENSHTEIN_EDIT_DISTANCE'::varchar AS method,
            e.canonical_name_en AS correction
        FROM entity_aliases a
        JOIN canonical_entities e ON a.entity_id = e.id
        WHERE a.status = 'ACTIVE' AND e.status = 'ACTIVE'
          AND length(v_norm) >= 4
          AND length(a.normalized_alias) >= 4
          AND levenshtein(v_norm, a.normalized_alias) <= 2

        UNION ALL

        -- Stage 5: Phonetic / Double Metaphone Match (Latin only)
        SELECT 
            a.entity_id,
            75.0::double precision AS match_score,
            'MEDIUM'::varchar AS conf_tier,
            'PHONETIC_DMETAPHONE'::varchar AS method,
            e.canonical_name_en AS correction
        FROM entity_aliases a
        JOIN canonical_entities e ON a.entity_id = e.id
        WHERE a.status = 'ACTIVE' AND e.status = 'ACTIVE'
          AND v_dmeta <> ''
          AND dmetaphone(a.normalized_alias) = v_dmeta
          AND length(v_norm) >= 3

        UNION ALL

        -- Stage 6: Substring / Prefix Matching on Canonical Names & Slugs
        SELECT 
            e.id AS entity_id,
            60.0::double precision AS match_score,
            'LOW'::varchar AS conf_tier,
            'SUBSTRING'::varchar AS method,
            e.canonical_name_en AS correction
        FROM canonical_entities e
        WHERE e.status = 'ACTIVE'
          AND (
              e.canonical_name_en ILIKE v_norm || '%'
              OR e.canonical_slug ILIKE v_norm || '%'
              OR COALESCE(e.canonical_name_am, '') ILIKE v_norm || '%'
          )
    ),
    aggregated AS (
        SELECT 
            cm.entity_id,
            MAX(cm.match_score) AS best_score,
            (ARRAY_AGG(cm.conf_tier ORDER BY cm.match_score DESC))[1] AS best_tier,
            (ARRAY_AGG(cm.method ORDER BY cm.match_score DESC))[1] AS best_method,
            (ARRAY_AGG(cm.correction ORDER BY cm.match_score DESC))[1] AS best_correction
        FROM candidate_matches cm
        GROUP BY cm.entity_id
    )
    SELECT 
        e.id,
        e.canonical_name_en,
        e.canonical_name_am,
        e.canonical_slug,
        e.entity_nature,
        e.entity_type,
        e.short_definition_en,
        e.short_definition_am,
        e.primary_category_id,
        agg.best_score,
        agg.best_tier,
        agg.best_method,
        agg.best_correction
    FROM aggregated agg
    JOIN canonical_entities e ON agg.entity_id = e.id
    ORDER BY agg.best_score DESC, e.id ASC
    LIMIT p_limit;
END;
$function$;

-- 8. High-Performance Discovery Graph Path Extractor
CREATE OR REPLACE FUNCTION public.get_entity_discovery_paths(p_entity_id int)
RETURNS TABLE (
    out_path_type character varying,
    out_target_id integer,
    out_target_name character varying,
    out_relationship_type character varying,
    out_depth integer,
    out_confidence double precision
)
LANGUAGE plpgsql
STABLE
AS $function$
BEGIN
    RETURN QUERY
    -- 1. Primary Category Discovery Path
    SELECT 
        'PRIMARY_CATEGORY'::varchar AS path_type,
        c.id AS target_id,
        c.name_en AS target_name,
        'PRIMARY_TAXONOMY'::varchar AS relationship_type,
        c.level AS depth,
        1.0::double precision AS confidence
    FROM canonical_entities e
    JOIN category_hierarchy c ON e.primary_category_id = c.id
    WHERE e.id = p_entity_id

    UNION ALL

    -- 2. Direct Related Category Graph Paths (DISCOVERABLE_THROUGH, USED_FOR, USED_IN_INDUSTRY)
    SELECT 
        'RELATED_CATEGORY'::varchar AS path_type,
        c.id AS target_id,
        c.name_en AS target_name,
        r.relationship_type::varchar AS relationship_type,
        c.level AS depth,
        r.confidence_score::double precision AS confidence
    FROM entity_relationships r
    JOIN category_hierarchy c ON r.target_category_id = c.id
    WHERE r.source_entity_id = p_entity_id AND r.status = 'ACTIVE'

    UNION ALL

    -- 3. Related Canonical Entities (IS_A, INSTANCE_OF, COMPLEMENTS, REQUIRES_SKILL)
    SELECT 
        'RELATED_ENTITY'::varchar AS path_type,
        t.id AS target_id,
        t.canonical_name_en AS target_name,
        r.relationship_type::varchar AS relationship_type,
        1 AS depth,
        r.confidence_score::double precision AS confidence
    FROM entity_relationships r
    JOIN canonical_entities t ON r.target_entity_id = t.id
    WHERE r.source_entity_id = p_entity_id AND r.status = 'ACTIVE'

    ORDER BY path_type ASC, depth ASC, confidence DESC;
END;
$function$;

-- 9. Authoritative Master Seed Data: 17 Multi-Domain Canonical Entities
-- Spans Material & Non-Material across SYS_PROD, SYS_SRV, SYS_EMP, SYS_GIG, SYS_MED, SYS_OPP, SYS_PRM

INSERT INTO canonical_entities (
    id, canonical_name_en, canonical_name_am, canonical_name_om, canonical_name_ti,
    canonical_slug, entity_nature, entity_type,
    short_definition_en, short_definition_am, short_definition_om, short_definition_ti,
    detailed_definition_en, detailed_definition_am,
    typical_examples, typical_use_cases, functional_capabilities,
    relevant_attributes, disambiguation_rules, primary_category_id
) VALUES
-- 1. Smartphone (SYS_PROD)
(
    1, 'Smartphone', 'ስማርት ፎን', 'Bilbila Ammayyaa', 'ስማርት ቴሌፎን',
    'smartphone', 'MATERIAL', 'PRODUCT',
    'A portable electronic device providing advanced mobile telecommunications, computing, internet connectivity, and digital media processing.',
    'ተንቀሳቃሽ የስልክ፣ የመልዕክት፣ የኢንተርኔት እና የተራቀቀ ኮምፒውቲንግ አገልግሎት የሚሰጥ ዘመናዊ ኤሌክትሮኒክስ መሣሪያ።',
    'Meeshaa elektirooniksii socho''aa qunnamtii bilbilaa, intarneetii fi tajaajila ammayyaa kennu.',
    'ዝተራቐቐ ናይ ቴሌፎን፣ ኢንተርኔትን ኮምፒዩተርን ግልጋሎት ዝህብ ተንቀሳቓሲ መሳርሒ።',
    'A multifunctional consumer handheld device running a mobile operating system (e.g. Android, iOS) integrating touchscreen display, cellular modems, cameras, GPS, and third-party application execution.',
    'ተንቀሳቃሽ ኦፕሬቲንግ ሲስተም (አንድሮይድ፣ አይኦኤስ) የሚጠቀም፣ የንክኪ ስክሪን፣ ካሜራ፣ ጂፒኤስ እና የኢንተርኔት ግንኙነት ያለው የግል መገናኛ መሣሪያ።',
    '["Apple iPhone 15 Pro", "Samsung Galaxy S24 Ultra", "Tecno Camon 30", "Xiaomi Redmi Note 13"]'::jsonb,
    '["Voice and video calling", "Mobile banking and digital payments", "Web browsing and social networking", "High-definition photography and videography"]'::jsonb,
    '["Cellular voice/data communication", "Touchscreen interaction", "Biometric authentication (fingerprint/face recognition)", "High-speed application processing"]'::jsonb,
    '{"storage_capacities": ["64GB", "128GB", "256GB", "512GB", "1TB"], "ram_options": ["4GB", "6GB", "8GB", "12GB", "16GB"], "operating_systems": ["Android", "iOS"], "sim_support": ["Single SIM", "Dual SIM", "eSIM"]}'::jsonb,
    '{"must_not_confuse_with": ["Landline phone", "Walkie-talkie", "Standalone digital camera", "Tablet PC"]}'::jsonb,
    1425 -- Smartphones (smartphones-mobile-devices)
),

-- 2. Digital Camera (SYS_PROD)
(
    2, 'Digital Camera', 'ዲጂታል ካሜራ', 'Kaameraa Diijitaalaa', 'ዲጂታል ካሜራ',
    'camera', 'MATERIAL', 'PRODUCT',
    'An optical device engineered primarily for capturing high-resolution still photographs and motion video to digital storage.',
    'ፎቶግራፎችን እና ቪዲዮዎችን በከፍተኛ ጥራት በዲጂታል መንገድ ለመቅረጽ የተዘጋጀ የኦፕቲካል መሣሪያ።',
    'Meeshaa optikaalaa suuraa fi viidiyoo qulqullina olaanaadhaan waraabuuf oolu.',
    'ስእልን ቪድዮን ብብቑዕ ጽሬት ንምቕራጽ ዝተሰርሐ ናይ ኦፕቲክስ መሳርሒ።',
    'A dedicated optical recording instrument incorporating interchangeable lenses or optical zoom, an electronic image sensor (CMOS/CCD), manual exposure controls, and professional image signal processing.',
    'ሌንስ፣ ሴንሰር እና መቆጣጠሪያዎችን የያዘ፣ ለሙያዊም ሆነ ለግል ፎቶግራፊ እና ቪዲዮግራፊ የሚያገለግል መሣሪያ።',
    '["Canon EOS R6 Mark II", "Sony Alpha A7 IV", "Nikon Z8", "Fujifilm X-T5"]'::jsonb,
    '["Professional studio photography", "Cinematic video production", "Photojournalism", "Wildlife and sports recording"]'::jsonb,
    '["High dynamic range optical sensing", "Interchangeable lens mounting", "Raw image data recording", "Manual aperture/shutter/ISO controls"]'::jsonb,
    '{"sensor_format": ["Full-frame", "APS-C", "Micro Four Thirds", "Medium Format"], "megapixels": ["24MP", "33MP", "45MP", "61MP"], "video_resolution": ["4K 60fps", "8K 30fps"]}'::jsonb,
    '{"must_not_confuse_with": ["Smartphone camera module", "CCTV security camera", "Webcam"]}'::jsonb,
    12 -- Cameras (cameras)
),

-- 3. Walkie-Talkie (SYS_PROD)
(
    3, 'Walkie-Talkie Transceiver', 'ዋኪ ቶኪ መገናኛ ሬዲዮ', 'Raadiyoo Walkie-Talkie', 'ዋኪ ቶኪ ሬድዮ',
    'walkie-talkie', 'MATERIAL', 'PRODUCT',
    'A portable, handheld two-way radio transceiver designed for immediate push-to-talk bidirectional voice communications without cellular towers.',
    'የሞባይል ኔትወርክ ሳይፈልግ በቀጥታ የድምፅ ግንኙነት የሚያደርግ ተንቀሳቃሽ የሁለት አቅጣጫ ሬዲዮ መገናኛ መሣሪያ።',
    'Meeshaa raadiyoo kallattiin qunnamtii sagalee saffisaa gochuuf dandeessisu.',
    'ብዘይ ሞባይል ኔትዎርክ ብቐጥታ ናይ ድምጺ ርክብ ዝግበረሉ ተንቀሳቓሲ ሬድዮ መሳርሒ።',
    'A battery-powered RF transceiver operating across UHF/VHF frequencies utilizing a push-to-talk (PTT) switch for instant group communication in tactical, security, and industrial operations.',
    'የራዲዮ ፍሪኩዌንሲዎችን በመጠቀም በቡድን ፈጣን ንግግር ለማድረግ የሚያስችል የደህንነትና የኢንዱስትሪ መገናኛ መሣሪያ።',
    '["Baofeng UV-5R", "Motorola Solutions T800", "Kenwood ProTalk TK-3402U", "Retevis RT22"]'::jsonb,
    '["Building site and construction coordination", "Security team surveillance", "Event crowd management", "Remote agricultural communications"]'::jsonb,
    '["Push-to-talk RF transmission", "Channel selector switches", "Squelch noise filtering", "Long-range off-grid radio link"]'::jsonb,
    '{"frequency_bands": ["VHF 136-174MHz", "UHF 400-520MHz", "FRS/PMR446"], "range_km": ["3km", "5km", "10km"], "battery_capacity_mah": ["1500mAh", "2800mAh"]}'::jsonb,
    '{"must_not_confuse_with": ["Cellular smartphone", "Broadcast FM/AM radio receiver", "Satellite phone"]}'::jsonb,
    9 -- Phones (mobile-phones)
),

-- 4. Power Drill (SYS_PROD / Tools)
(
    4, 'Power Drill', 'የኤሌክትሪክ ድሪል / መሰርሰሪያ', 'Makiinaa Qotuu / Drill', 'ኤሌክትሪክ መሳርሒ መውቀሪ / ድሪል',
    'power-drill', 'MATERIAL', 'TOOL_EQUIPMENT',
    'A handheld rotating or percussive electromechanical power tool designed to bore holes in solid substrates and drive threaded fasteners.',
    'በእንጨት፣ በብረት እና በኮንክሪት ላይ ቀዳዳዎችን ለመሰርሰር እና ብሎኖችን ለማሰር የሚያገለግል የኤሌክትሪክ መሣሪያ።',
    'Meeshaa elektirikii sibiila, mukaa fi konkiriti irratti qaawwa qotuu fi booltii hidhuuf gargaaru.',
    'ኣብ እንጨይቲ፣ ሓጺንን ኮንክሪትን ጐድጓድ ንምውቃዕን ብሎን ንምእሳርን ዝጠቅም ናይ ኤሌክትሪክ መሳርሒ።',
    'An electromechanical workshop and construction tool driven by a brushed or brushless motor, with adjustable chuck, torque clutch, and optional hammering percussion for masonry boring.',
    'በባትሪ ወይም በገመድ ኤሌክትሪክ የሚሰራ፣ ፍጥነቱ እና ጉልበቱ የሚቆጣጠር፣ ለግንባታና የእንጨት/ብረት ስራ የሚውል ዘመናዊ መሣሪያ።',
    '["Bosch Professional GSB 18V-50", "DeWalt DCD791 20V MAX", "Makita HP333D Cordless Hammer Drill", "Ingco Cordless 20V Impact Drill"]'::jsonb,
    '["Drilling pilot holes in masonry and hardwood", "Driving structural screws and anchors", "Cabinetry and furniture fabrication", "Metal sheet perforation"]'::jsonb,
    '["Bi-directional variable speed rotation", "Torque limiter clutch mechanism", "Impact percussion pulse generator", "Keyless chuck clamping"]'::jsonb,
    '{"power_source": ["Cordless Battery (18V/20V)", "Corded Electric (220V)"], "chuck_size_mm": ["10mm", "13mm"], "motor_type": ["Brushless", "Brushed"], "impact_rate_bpm": ["0-25000", "0-30000"]}'::jsonb,
    '{"must_not_confuse_with": ["Water well drilling rig", "Dental clinic drill", "Industrial mining bore system", "Rotary hammer demolition breaker"]}'::jsonb,
    3623 -- Drills (drills-and-drivers-hardware)
),

-- 5. Passenger Car (SYS_PROD / Vehicles)
(
    5, 'Passenger Car', 'የግል መኪና', 'Konkolaataa Daangaa', 'ናይ ውልቂ መኪና',
    'passenger-car', 'MATERIAL', 'VEHICLE',
    'A wheeled motor vehicle primarily intended for transporting passengers, containing between one and eight passenger seating positions.',
    'ሰዎችን ከቦታ ወደ ቦታ ለማጓጓዝ የተሰራ አራት ጎማ ያለው የሞተር ተሽከርካሪ።',
    'Konkolaataa motora qabu kan namoota deddeebisuuf qophaa''e.',
    'ሰባት ንምጉዕዓዝ ዝተዳለወ ባለ 4 እግሪ ሞተር ተሽከርካሪ።',
    'A road vehicle propelled by an internal combustion engine, hybrid electric drivetrain, or battery electric system, fitted with seats for personal transportation on paved networks.',
    'በነዳጅ ወይም በኤሌክትሪክ የሚሰራ፣ ሹፌሩን ጨምሮ እስከ 5-8 ሰዎችን የሚያጓጉዝ የግል ወይም የቤተሰብ ተሽከርካሪ።',
    '["Toyota Corolla Sedan", "Hyundai Elantra", "Suzuki Dzire", "Honda Civic"]'::jsonb,
    '["Daily urban commuting", "Intercity personal transit", "Rideshare and taxi operation", "Family road transportation"]'::jsonb,
    '["Automated steering control", "Four-wheel braking with ABS", "Climate cabin regulation", "Passenger safety restraint cage"]'::jsonb,
    '{"fuel_types": ["Petrol", "Diesel", "Hybrid", "Electric"], "transmission": ["Manual", "Automatic", "CVT"], "seating_capacity": [4, 5, 7], "engine_displacement_cc": ["1000cc", "1300cc", "1600cc", "2000cc"]}'::jsonb,
    '{"must_not_confuse_with": ["Heavy commercial truck", "Passenger bus", "Motorcycle", "Agricultural tractor"]}'::jsonb,
    14 -- Cars (cars)
),

-- 6. Toyota Vitz (SYS_PROD / Vehicles / Specific Model)
(
    6, 'Toyota Vitz Subcompact', 'ቶዮታ ቪትስ መኪና', 'Konkolaataa Toyota Vitz', 'ቶዮታ ቪትስ',
    'toyota-vitz', 'MATERIAL', 'VEHICLE',
    'A subcompact B-segment three- or five-door hatchback automobile engineered and manufactured globally by Toyota Motor Corporation.',
    'በቶዮታ ኩባንያ የሚመረት አነስተኛ ነዳጅ ቆጣቢ የከተማ መኪና (ሃችባክ)።',
    'Konkolaataa qusataa gabaabaa kaampaanii Tooyootaatiin oomishamu.',
    'ብቶዮታ ዝስራሕ ቍጠባዊ ናይ ከተማ መኪና።',
    'A popular subcompact hatchback celebrated for exceptional fuel economy, reliability, maneuverability in congested cities, and broad parts availability throughout Ethiopia.',
    'በኢትዮጵያ ገበያ እጅግ ተወዳጅ የሆነ፣ ነዳጅ የሚቆጥብ እና በቀላሉ ለመንዳት ምቹ የሆነ የከተማ መኪና።',
    '["Toyota Vitz 2005 1.0L", "Toyota Vitz 2010 1.3L 4WD", "Toyota Vitz Jewela 2015", "Toyota Yaris Hatchback"]'::jsonb,
    '["City commuting in dense urban centers", "Fuel-efficient delivery and errands", "First-time vehicle ownership", "Urban ride-hailing services"]'::jsonb,
    '["Ultra-low fuel consumption", "Tight turning radius", "Electronic power steering", "Fold-down rear cargo capacity"]'::jsonb,
    '{"engine_displacement": ["996cc (1KR-FE)", "1298cc (2SZ-FE / 1NR-FE)", "1496cc (1NZ-FE)"], "transmission": ["Automatic (CVT)", "5-speed Manual"], "drive_layout": ["Front-Wheel Drive (FWD)", "All-Wheel Drive (4WD)"]}'::jsonb,
    '{"must_not_confuse_with": ["Toyota Corolla", "Toyota RAV4", "Suzuki Swift", "Toyota Hilux"]}'::jsonb,
    14 -- Cars (cars)
),

-- 7. Apartment Property (SYS_PRM / Property)
(
    7, 'Apartment Residence', 'አፓርታማ የመኖሪያ ቤት', 'Mana Jireenyaa Gamoo (Apartment)', 'ኣፓርታማ ናይ መኖርያ ገዛ',
    'apartment', 'MATERIAL', 'PROPERTY',
    'A self-contained residential housing unit located within a multi-unit condominium or multi-story apartment complex.',
    'በፎቅ ወይም በህንፃ ውስጥ የተሰራ የተሟላ የመኖሪያ ቤት ክፍል (ኮንዶሚኒየም ወይም ሪል እስቴት)።',
    'Mana jireenyaa gamoo qindaawaa keessatti argamu.',
    'ኣብ ደርቢ ህንጻ ውሽጢ ዝርከብ ምሉእ ናይ መኖርያ ገዛ።',
    'A private residential domicile situated within a shared real estate superstructure featuring independent utilities, private entry, kitchen, bath facilities, and shared common area amenities.',
    'ራሱን የቻለ ኩሽና፣ መኝታ፣ ሳሎን እና መታጠቢያ ቤት ያለው በህንፃ ውስጥ የሚገኝ የመኖሪያ ቤት።',
    '["2-Bedroom Bole condominium", "Luxury 3-bedroom furnished apartment CMC", "Studio apartment Kazanchis", "Penthouse apartment Sarbet"]'::jsonb,
    '["Primary family residence", "Long-term rental tenancy", "Short-term expat and tourist accommodation", "Real estate investment property"]'::jsonb,
    '["Independent living accommodation", "Secure building entry", "Elevator and standby generator access", "Reserved basement parking"]'::jsonb,
    '{"bedrooms": [1, 2, 3, 4], "bathrooms": [1, 2, 3], "floor_area_sqm": ["50sqm", "85sqm", "120sqm", "160sqm", "220sqm"], "furnishing": ["Furnished", "Semi-Furnished", "Unfurnished"], "tenure": ["For Sale", "For Rent"]}'::jsonb,
    '{"must_not_confuse_with": ["Commercial office space", "Standalone villa compound", "Warehouse depot", "Retail shop stall"]}'::jsonb,
    20 -- Apartments (apartments)
),

-- 8. Professional Plumber (SYS_EMP / Profession)
(
    8, 'Professional Plumber', 'የቧንቧ ባለሙያ / ቧንቧ ሠራተኛ', 'Ogeessa Ujummoo Bishaanii', 'ሞያተኛ ቧንቧ',
    'plumber', 'NON_MATERIAL', 'PROFESSION',
    'A skilled licensed tradesperson specializing in assembling, installing, inspecting, and repairing piping systems for potable water and drainage.',
    'የውሃ እና የፍሳሽ ቧንቧዎችን የመገጣጠም፣ የመዘርጋትና የመጠገን ሙያ ያለው የቴክኒክ ባለሙያ።',
    'Ogeessa ujummoolee bishaanii fi dhangala''aa diriirsuu fi suphuu danda''u.',
    'መስመር ቧንቧ ማይ ዝዝርግሕን ዝጽግንን ሞያተኛ።',
    'An artisan professional trained in sanitary plumbing codes, pipe fitting, water pump installation, solar water heater integration, sanitary ware installation, and hydraulic line diagnostic testing.',
    'በህንጻዎች እና በቤቶች ውስጥ የንጹህ ውሃ እና የፍሳሽ መስመሮችን የሚዘረጋ፣ የሚጠግንና የሚመረምር የሞያ ባለቤት።',
    '["Master Sanitary Plumber", "Residential Pipe Fitter", "Commercial Water Works Technician"]'::jsonb,
    '["New construction sanitary installation", "Emergency pipeline burst repair", "Pressure pump overhaul", "Bathroom fixture renovation"]'::jsonb,
    '["Hydraulic schematic interpretation", "PPR/PVC hot fusion pipe joining", "Leak pressure diagnostics", "Drainage gradient surveying"]'::jsonb,
    '{"experience_years": ["1-3 years", "3-5 years", "5-10 years", "10+ years"], "certifications": ["TVET Level IV Sanitary", "Master Craftsman License"], "employment_type": ["Full-time", "Contract", "On-demand Emergency"]}'::jsonb,
    '{"must_not_confuse_with": ["Electrician", "Civil mason", "HVAC technician", "Plumbing equipment supplier"]}'::jsonb,
    1244 -- Plumbing Services (plumbing-services)
),

-- 9. Plumbing Service (SYS_SRV / Service)
(
    9, 'Plumbing Repair & Installation Service', 'የቧንቧ ስራ እና ጥገና አገልግሎት', 'Tajaajila Ujummoo Bishaanii Suphuu', 'ግልጋሎት ጽገና መስመር ቧንቧ',
    'plumbing-service', 'NON_MATERIAL', 'SERVICE',
    'The professional execution of plumbing system design, maintenance, leak mitigation, sanitary ware fitting, and water infrastructure rehabilitation.',
    'የውሃ እና ፍሳሽ ቧንቧ ጥገና፣ አዳዲስ መስመሮችን የመዘርጋት እና የውሃ መሣሪያዎችን የመግጠም አገልግሎት።',
    'Tajaajila ujummoo diriirsuu, suphuu fi dhangala''aa to''achuu.',
    'ናይ ማይን ፍሳሽን መስመር ቧንቧ ናይ ምዝርጋሕን ምጽጋንን ግልጋሎት።',
    'Commercial or residential service contracts executed by qualified technicians to address pipe bursts, water tank installations, sewer unclogging, fixture replacement, and preventative hydraulic maintenance.',
    'ለቤትም ሆነ ለድርጅት የሚሰጥ የውሃ ቧንቧ፣ የቦይለር፣ የሲንክ እና የመታጠቢያ ቤት ዕቃዎች ጥገናና ተከላ አገልግሎት።',
    '["Emergency pipe leak repair", "Borehole pump connection service", "Bathroom remodeling plumbing", "Sewer line snake unclogging"]'::jsonb,
    '["Rapid response residential flood prevention", "Commercial building preventative maintenance", "Industrial water line installations"]'::jsonb,
    '["Pressurized hydrostatic testing", "High-pressure jet unclogging", "Electro-fusion pipe welding", "Acoustic leak detection"]'::jsonb,
    '{"pricing_model": ["Per hour", "Fixed quote per project", "Emergency flat rate"], "service_warranty": ["30 days", "90 days", "1 year"], "availability": ["24/7 On-call", "Standard business hours"]}'::jsonb,
    '{"must_not_confuse_with": ["Plumbing raw pipe sales", "Electrical rewiring service", "Water delivery trucking"]}'::jsonb,
    1244 -- Plumbing Services (plumbing-services)
),

-- 10. Software Engineer (SYS_EMP / Profession)
(
    10, 'Software Engineer', 'ሶፍትዌር ኢንጂነር', 'Injinera Sooftiweerii', 'ሶፍትዌር መሃንድስ',
    'software-engineer', 'NON_MATERIAL', 'PROFESSION',
    'A trained computer science professional who applies engineering principles to design, develop, test, deploy, and maintain software applications and systems.',
    'ኮምፒውተር ፕሮግራሞችን፣ ድረ-ገጾችን እና የሞባይል አፕሊኬሽኖችን የሚያለማና የሚገነባ የቴክኖሎጂ ባለሙያ።',
    'Ogeessa teeknoolojii sooftiweerii fi appilikeeshinii kalaqu fi hojjetu.',
    'ሶፍትዌርን አፕሊኬሽንን ዝሰርሕን ዝምዕብልን ሞያተኛ።',
    'A qualified technologist proficient in algorithm design, data structures, full-stack programming, cloud architecture, and database management, engineering resilient computing solutions.',
    'በተለያዩ የፕሮግራሚንግ ቋንቋዎች ኮድ በመጻፍ የንግድ፣ የሳይንስ እና የማህበራዊ ችግሮችን የሚፈቱ የሶፍትዌር ስርዓቶችን የሚገነባ ባለሙያ።',
    '["Senior Backend Go Engineer", "Frontend React/Next.js Architect", "Mobile Flutter Developer", "DevOps Cloud Engineer"]'::jsonb,
    '["Building scalable cloud commerce platforms", "Developing native mobile applications", "Architecting resilient REST & gRPC microservices", "Database schema and query performance tuning"]'::jsonb,
    '["Algorithmic problem solving", "Object-oriented & functional code authoring", "Distributed systems concurrency modeling", "Automated unit and integration testing"]'::jsonb,
    '{"primary_languages": ["Go", "TypeScript", "Python", "Rust", "Java", "C++"], "frameworks": ["Next.js", "React", "PostgreSQL", "Kubernetes", "Docker"], "experience_level": ["Junior", "Mid-Level", "Senior", "Principal/Lead"]}'::jsonb,
    '{"must_not_confuse_with": ["Computer hardware repair technician", "Data entry clerk", "Graphic designer", "IT helpdesk support"]}'::jsonb,
    1191 -- Software Development Services (software-development-services)
),

-- 11. Software Development Service (SYS_SRV / Service)
(
    11, 'Custom Software Development Service', 'የሶፍትዌር ማበልፀግ አገልግሎት', 'Tajaajila Sooftiweerii Hojjechuu', 'ግልጋሎት ምስራሕ ሶፍትዌር',
    'software-development', 'NON_MATERIAL', 'SERVICE',
    'A contracted technology service involving the turnkey analysis, design, programming, quality assurance, and hosting of custom software applications.',
    'ለድርጅቶች ወይም ለግለሰቦች የተዘጋጀ የኮምፒውተር ወይም የሞባይል ሶፍትዌር ሰርቶ የማስረከብ የቴክኖሎጂ አገልግሎት።',
    'Tajaajila sooftiweerii fi weebsaayitii daldalaaf hojjetanii kennuu.',
    'ናይ ንግዲ ትካላት ሶፍትዌርን መርበብ ሓበሬታን ናይ ምምዕባል ግልጋሎት።',
    'An agency or consultancy technical offering providing bespoke web applications, enterprise ERP/CRM development, fintech payments integration, mobile app creation, and maintenance agreements.',
    'የኢ-ኮሜርስ፣ የክፍያ ስርዓቶች፣ የሞባይል አፕ እና የድረ-ገጽ ስርዓቶችን የማበልጸግና የማስተዳደር ፕሮፌሽናል አገልግሎት።',
    '["Fintech payment gateway integration", "Custom inventory ERP development", "Cross-platform iOS/Android app building", "Enterprise cloud migration"]'::jsonb,
    '["Digital transformation of legacy commerce", "Automating business transaction pipelines", "Deploying scalable customer portals"]'::jsonb,
    '["Agile project sprint delivery", "Secure code auditing and vulnerability penetration", "Continuous integration and deployment pipelines", "High availability server maintenance"]'::jsonb,
    '{"delivery_model": ["Time and Materials", "Milestone-based Fixed Price", "Dedicated Engineering Team Staffing"], "sla_options": ["Standard 99.5%", "Mission-critical 99.99%"], "contract_scope": ["Full Turnkey", "Backend Only", "Mobile Only"]}'::jsonb,
    '{"must_not_confuse_with": ["Commercial off-the-shelf software license", "Hardware maintenance service", "Digital marketing campaign"]}'::jsonb,
    1191 -- Software Development Services (software-development-services)
),

-- 12. House Cleaning Service (SYS_SRV / Service)
(
    12, 'Residential House Cleaning Service', 'የቤት ጽዳት አገልግሎት', 'Tajaajila Mana Qulqulleessuu', 'ግልጋሎት ጽሬት ገዛ',
    'house-cleaning-service', 'NON_MATERIAL', 'SERVICE',
    'A domestic care service offering thorough sanitization, dusting, vacuuming, window washing, and general cleanliness maintenance for homes.',
    'ለመኖሪያ ቤቶች የሚሰጥ የወለል፣ የመስኮት፣ የኩሽና፣ የሳሎን እና የመታጠቢያ ቤት ጥልቅ የጽዳት አገልግሎት።',
    'Tajaajila mana jireenyaa fi mi''a mana keessaa qulqulleessuu.',
    'ናይ መኖርያ ኣባይቲ ምሉእ ናይ ምጽራይ ግልጋሎት።',
    'Professional custodial services employing trained personnel equipped with non-toxic sanitizing agents, industrial vacuums, and steam extractors to clean apartments, villas, and residential premises.',
    'ባለሙያዎችን በመላክ ቤትን ከነፍሳት የጸዳ፣ ንጹህና ያማረ የሚያደርግ ሙያዊ የጽዳት አገልግሎት።',
    '["Post-construction deep residential cleaning", "Move-in / move-out apartment sanitization", "Bi-weekly scheduled home maid cleaning", "Upholstery sofa steam cleaning"]'::jsonb,
    '["Restoring hygiene to post-renovation residences", "Preparing leased residences for new tenants", "Routine household sanitation maintenance"]'::jsonb,
    '["Deep allergen and dust extraction", "Bathroom lime-scale and mold eradication", "Kitchen grease degreasing", "Hardwood and tile polishing"]'::jsonb,
    '{"service_frequency": ["One-time deep clean", "Weekly", "Bi-weekly", "Monthly"], "crew_size": [1, 2, 4, 6], "cleaning_supplies": ["Supplied by agency", "Client provided"]}'::jsonb,
    '{"must_not_confuse_with": ["Industrial toxic waste cleanup", "Dry cleaning laundry service", "Car wash detailing"]}'::jsonb,
    889 -- Deep Cleaning Services (deep-cleaning-services)
),

-- 13. Portland Cement (SYS_PROD / Raw Material)
(
    13, 'Portland Cement', 'ፖርትላንድ ሲሚንቶ', 'Simintoo Boortilaandii', 'ሲሚንቶ',
    'portland-cement', 'MATERIAL', 'RAW_MATERIAL',
    'A hydraulic binder produced by pulverizing clinker, crucial for formulating structural concrete, mortar, and stucco in building construction.',
    'ለኮንክሪት፣ ለግንባታ እና ለማንኛውም ህንጻ ስራ ወሳኝ የሆነ የዱቄት የግንባታ ማጣበቂያ ጥሬ እቃ።',
    'Meeshaalee ijaarsaaf barbaachisaa ta''e kan simintoo dhaagaa irraa tolfamu.',
    'ንህንጻን ንኮንክሪትን ዝጠቅም መሰረታዊ ጥረ ኣቕሓ ሲሚንቶ።',
    'A manufactured mineral powder comprised of calcium silicates, aluminates, and ferrites that hydrates chemically upon water addition to set into stone-like compressive strength.',
    'አሸዋና ጠጠር ከውሃ ጋር በማዋሃድ ጠንካራ ኮንክሪት ለመስራት የሚያገለግል በፋብሪካ የተፈበረከ የሲሚንቶ ምርት።',
    '["Dangote OPC 42.5R Portland Cement", "Mugher PPC 32.5N Cement", "National Cement OPC 42.5N", "Derba MIDROC Cement 42.5R"]'::jsonb,
    '["Structural reinforced concrete pouring (slabs, columns, beams)", "Brick and block masonry mortar binding", "Foundation footing stabilization", "Precast concrete element manufacturing"]'::jsonb,
    '["High hydraulic compressive strength bonding", "Controlled setting and curing rates", "Moisture and weather resilience", "Cohesive mortar workability"]'::jsonb,
    '{"grade": ["OPC 42.5R", "OPC 42.5N", "PPC 32.5R", "PPC 32.5N"], "package_weight_kg": [50], "setting_time_initial_min": [60, 90], "compressive_strength_28day_mpa": ["32.5 MPa", "42.5 MPa", "52.5 MPa"]}'::jsonb,
    '{"must_not_confuse_with": ["Gypsum plaster", "Hydrated lime powder", "Pre-mixed ready concrete", "Ceramic tile grout"]}'::jsonb,
    1212 -- Cement (cement-binders-mortar)
),

-- 14. Agricultural Tractor (SYS_PROD / Tool Equipment / Agriculture)
(
    14, 'Agricultural Tractor', 'የእርሻ ትራክተር', 'Tiraaktara Qonnaa', 'ትራክተር ሕርሻ',
    'agricultural-tractor', 'MATERIAL', 'TOOL_EQUIPMENT',
    'A high-traction, low-speed heavy motorized machine engineered to pull, push, and power agricultural implements for soil cultivation and harvesting.',
    'ማሳን ለማረስ፣ ለማለስለስ፣ ምርት ለመሰብሰብ እና የግብርና እቃዎችን ለመጎተት የሚያገለግል ኃይለኛ የሞተር ማሽን።',
    'Makiinaa qonnaa kan lafa qotuuf, misoomsuuf fi oomisha walitti qabuuf oolu.',
    'መሬት ንምሕራስን ንምድላውን ዝጠቅም ናይ ሕርሻ ከቢድ ማሽን።',
    'A heavy diesel prime mover equipped with large rear pneumatic drive wheels or tracks, three-point hydraulic hitch, and power take-off (PTO) shaft to operate plows, seeders, and balers.',
    'ለዘመናዊ ግብርና ምርታማነት ወሳኝ የሆነ፣ ማረሻና መዝሪያዎችን በጉልበቱ በማንቀሳቀስ የእርሻ ስራን የሚያቀላጥፍ ተሽከርካሪ።',
    '["Massey Ferguson 375 75HP Tractor", "John Deere 5075E 4WD", "New Holland TD5.90", "Mahindra 6075 4WD Heavy Agri Tractor"]'::jsonb,
    '["Deep moldboard plowing and subsoiling", "Rotary tillage and seedbed preparation", "Heavy agricultural trailer haulage", "Operating irrigation pumps via PTO shaft"]'::jsonb,
    '["High torque diesel pulling power", "Three-point hydraulic implement suspension", "Rotary Power Take-Off (PTO) output", "Differential axle lock for muddy terrain"]'::jsonb,
    '{"horsepower_hp": ["50HP", "75HP", "90HP", "120HP", "150HP+"], "drive_system": ["2WD", "4WD / MFWD"], "pto_speed_rpm": ["540 RPM", "1000 RPM"], "hitch_category": ["Category 2", "Category 3"]}'::jsonb,
    '{"must_not_confuse_with": ["Excavator earthmover", "Road construction grader", "Passenger commercial vehicle", "Lawn mower"]}'::jsonb,
    317 -- Tractors (tractors)
),

-- 15. Solar Inverter (SYS_PROD / Energy)
(
    15, 'Solar Power Inverter', 'የፀሐይ ብርሃን ኃይል ኢንቨርተር', 'Inveertarii Annisaa Biiftuu', 'ናይ ፀሓይ ኃይሊ ኢንቨርተር',
    'solar-inverter', 'MATERIAL', 'PRODUCT',
    'An electrical power electronic apparatus that converts direct current (DC) electricity from solar panels into alternating current (AC) for grid or off-grid loads.',
    'ከፀሐይ ፓነል የሚገኘውን ዲሲ (DC) የኤሌክትሪክ ኃይል ወደ ኤሲ (AC) በመቀየር የቤትና የድርጅት እቃዎች እንዲሰሩ የሚያደርግ የኤሌክትሮኒክስ መሣሪያ።',
    'Meeshaa elektirikii annisaa aduu gara elektirikii manaa jijjiiru.',
    'ካብ ፀሓይ ዝርከብ ሓይሊ ኤሌክትሪክ ናብ ናይ ገዛ ኤሌክትሪክ ዝቕይር መሳርሒ።',
    'A solid-state power conversion converter with maximum power point tracking (MPPT), battery management, and bidirectional grid synchronization for solar installations.',
    'መብራት ሲጠፋ በባትሪ ወይም በፀሐይ ብርሃን ኃይል ቀጣይነት ያለው የኤሌክትሪክ አቅርቦት የሚሰጥ የቴክኖሎጂ እቃ።',
    '["Growatt SPF 5000 ES Off-Grid Inverter", "Victron MultiPlus-II 48V 5kVA", "Must Solar 3kW Hybrid Inverter", "Huawei SUN2000 10KTL Grid-Tie"]'::jsonb,
    '["Providing off-grid solar power to rural dwellings", "Commercial backup power during utility load shedding", "Industrial solar grid tie-in", "Solar powered agricultural water pumping"]'::jsonb,
    '["Pure sine wave AC voltage synthesis", "MPPT solar charge optimization", "Seamless UPS transfer switching (<10ms)", "Lithium and gel battery charge regulation"]'::jsonb,
    '{"rated_power_kva": ["1kVA", "3kVA", "5kVA", "10kVA", "15kVA"], "battery_voltage_vdc": ["12V", "24V", "48V"], "inverter_type": ["Pure Sine Wave Hybrid", "Off-Grid Standalone", "Grid-Tied Microinverter"]}'::jsonb,
    '{"must_not_confuse_with": ["Automotive DC-DC booster", "Pure generator engine", "Storage battery unit", "Solar PV panel module"]}'::jsonb,
    35 -- Tools & Hardware / Power Systems
),

-- 16. Public Infrastructure Tender (SYS_OPP / Business Opportunity)
(
    16, 'Public Infrastructure Construction Tender', 'የመንግስት የኮንስትራክሽን ጨረታ', 'Caalbaa Ijaarsaa Mootummaa', 'ናይ ህንጸት ጨረታ',
    'public-infrastructure-tender', 'NON_MATERIAL', 'BUSINESS_OPPORTUNITY',
    'A formal public procurement invitation issued by government or enterprise entities seeking qualified contractors for engineering works.',
    'በመንግስት ወይም በትላልቅ ድርጅቶች ለህንፃ፣ መንገድ ወይም ውሃ ስራ ብቁ የሆኑ ተቋራጮችን ለመቅጠር የሚወጣ ህጋዊ የጨረታ ማስታወቂያ።',
    'Waamicha dorgommii ijaarsa pirojektoota gurguddoo mootummaa.',
    'ናይ መንግስቲ ወይ ትካል ናይ ስራሕ ዕድልን ውዕልን ንምውሳድ ዝወጽእ ጨረታ።',
    'A structured, legally binding procurement opportunity governed by public procurement authority guidelines, detailing bill of quantities (BOQ), bidding guarantee, eligibility criteria, and deadlines.',
    'ተቋራጮችና ባለሙያዎች ተወዳድረው የገንዘብና የስራ ውል የሚያገኙበት ህጋዊ የንግድ እድል።',
    '["Federal Roads Authority Highway Asphalt Tender", "Regional Deep Water Well Drilling Bid", "Educational Campus Building Construction RFQ", "Municipal Substation Electrification Tender"]'::jsonb,
    '["Securing commercial general contractor revenue", "Supplying industrial materials to state projects", "Civil engineering consortium procurement"]'::jsonb,
    '["Legally protected bidding framework", "Pre-qualified bidder evaluation", "Standardized FIDIC dispute conditions", "Letter of Credit or Bank Guarantee backing"]'::jsonb,
    '{"procurement_method": ["National Competitive Bidding (NCB)", "International Competitive Bidding (ICB)", "Direct Procurement Request"], "bid_bond_etb": ["50,000 ETB", "200,000 ETB", "500,000 ETB+"], "contractor_grade": ["GC/BC Grade 1", "Grade 2", "Grade 3", "Specialized Drilling Contractor"]}'::jsonb,
    '{"must_not_confuse_with": ["Direct retail purchase", "Employment vacancy job posting", "Stock share issuance"]}'::jsonb,
    1020109 -- Deep Water Well Drilling Bids (SYS_OPP)
),

-- 17. Freelance Graphic Designer (SYS_GIG / Freelance Skill)
(
    17, 'Freelance Graphic Designer', 'ፍሪላንስ ግራፊክ ዲዛይነር', 'Dizaayinara Giriifiksii Hojii Dhuunfaa', 'ፍሪላንስ ናይ ግራፊክስ ሞያተኛ',
    'freelance-graphic-designer', 'NON_MATERIAL', 'SKILL',
    'An independent creative visual specialist who crafts brand identities, promotional collateral, packaging, and digital UI assets on a contract basis.',
    'አርማዎችን፣ የማስታወቂያ ፖስተሮችን፣ ማሸጊያዎችን እና የዲጂታል ዲዛይኖችን በግል ስራ ተቀጥሮ የሚሰራ የፈጠራ ባለሙያ።',
    'Ogeessa dhuunfaa fakkiiwwan, loogoo fi beeksisa daldalaa hojjetu.',
    'ናይ ምልክትን ማስታወቂያን ዲዛይን ብውልቂ ዝሰርሕ ናይ ስነ-ጥበብ ሞያተኛ።',
    'A visual artist proficient in vector illustration, typography, color theory, and digital rendering tools, delivering commercial graphic assets for marketing, branding, and publication.',
    'የኮምፒውተር የንድፍ መተግበሪያዎችን በመጠቀም ለንግድ ተቋማት ማራኪ የማስታወቂያ እና የማስተዋወቂያ ስራዎችን የሚያከናውን ባለሙያ።',
    '["Corporate Brand Identity Designer", "Social Media Content Graphic Artist", "Packaging and Label Print Designer", "Illustrator and Vector Character Artist"]'::jsonb,
    '["Designing corporate logos and stationery suites", "Creating high-converting social media advertising creatives", "Designing product packaging cartons and labels", "Formatting annual corporate reports and magazines"]'::jsonb,
    '["Vector artwork authoring", "Print-ready CMYK prepress calibration", "Digital RGB color profile management", "Visual narrative branding development"]'::jsonb,
    '{"software_proficiency": ["Adobe Photoshop", "Adobe Illustrator", "Figma", "Canva", "InDesign", "Blender 3D"], "engagement_terms": ["Hourly Gig", "Fixed Per-Design Milestone", "Monthly Retainer"], "turnaround_time": ["24 hours", "3 days", "7 days"]}'::jsonb,
    '{"must_not_confuse_with": ["Full-stack software developer", "Fine arts canvas painter", "Printing press machine operator"]}'::jsonb,
    1020099 -- Voiceover & Creative Gig Category (SYS_GIG)
)
ON CONFLICT (id) DO UPDATE SET
    canonical_name_en = EXCLUDED.canonical_name_en,
    canonical_name_am = EXCLUDED.canonical_name_am,
    canonical_name_om = EXCLUDED.canonical_name_om,
    canonical_name_ti = EXCLUDED.canonical_name_ti,
    canonical_slug = EXCLUDED.canonical_slug,
    entity_nature = EXCLUDED.entity_nature,
    entity_type = EXCLUDED.entity_type,
    short_definition_en = EXCLUDED.short_definition_en,
    short_definition_am = EXCLUDED.short_definition_am,
    detailed_definition_en = EXCLUDED.detailed_definition_en,
    detailed_definition_am = EXCLUDED.detailed_definition_am,
    typical_examples = EXCLUDED.typical_examples,
    typical_use_cases = EXCLUDED.typical_use_cases,
    functional_capabilities = EXCLUDED.functional_capabilities,
    relevant_attributes = EXCLUDED.relevant_attributes,
    disambiguation_rules = EXCLUDED.disambiguation_rules,
    primary_category_id = EXCLUDED.primary_category_id,
    updated_at = NOW();

-- Align sequence
SELECT setval('canonical_entities_id_seq', (SELECT MAX(id) FROM canonical_entities));

-- 10. Seed Rich Entity Aliases (Latin, Amharic Ge'ez, Afaan Oromo, Tigrinya, Typos, Misspellings, Transliterations)
INSERT INTO entity_aliases (entity_id, alias_text, normalized_alias, language, alias_type, confidence_score)
VALUES
    -- Smartphone (Entity 1)
    (1, 'smartphone', 'smartphone', 'en', 'EXACT_SYNONYM', 1.00),
    (1, 'smartphones', 'smartphones', 'en', 'EXACT_SYNONYM', 1.00),
    (1, 'phone', 'phone', 'en', 'COMMON_NAME', 0.95),
    (1, 'phones', 'phones', 'en', 'COMMON_NAME', 0.95),
    (1, 'mobile phone', 'mobile phone', 'en', 'COMMON_NAME', 0.98),
    (1, 'mobile', 'mobile', 'en', 'COMMON_NAME', 0.90),
    (1, 'cellphone', 'cellphone', 'en', 'COMMON_NAME', 0.95),
    (1, 'cell phone', 'cell phone', 'en', 'COMMON_NAME', 0.95),
    (1, 'phne', 'phne', 'en', 'TYPO_CORRECTION', 0.90),
    (1, 'fone', 'fone', 'en', 'PHONETIC', 0.88),
    (1, 'smartfone', 'smartfone', 'en', 'PHONETIC', 0.90),
    (1, 'smrtphone', 'smrtphone', 'en', 'TYPO_CORRECTION', 0.92),
    (1, 'iphone', 'iphone', 'en', 'BRAND_REFERENCE', 0.90),
    (1, 'samsung galaxy', 'samsung galaxy', 'en', 'BRAND_REFERENCE', 0.90),
    (1, 'ስልክ', 'ስልክ', 'am', 'COMMON_NAME', 1.00),
    (1, 'ስማርት ፎን', 'ስማርት ፎን', 'am', 'EXACT_SYNONYM', 1.00),
    (1, 'ሞባይል', 'ሞባይል', 'am', 'COMMON_NAME', 0.98),
    (1, 'ሞባይል ስልክ', 'ሞባይል ስልክ', 'am', 'COMMON_NAME', 1.00),
    (1, 'selk', 'selk', 'am', 'TRANSLITERATION', 0.90),
    (1, 'bilbila', 'bilbila', 'om', 'COMMON_NAME', 1.00),
    (1, 'bilbila socho''aa', 'bilbila sochoaa', 'om', 'COMMON_NAME', 0.95),
    (1, 'ቴሌፎን', 'ቴሌፎን', 'ti', 'COMMON_NAME', 1.00),

    -- Digital Camera (Entity 2)
    (2, 'camera', 'camera', 'en', 'EXACT_SYNONYM', 1.00),
    (2, 'cameras', 'cameras', 'en', 'EXACT_SYNONYM', 1.00),
    (2, 'camra', 'camra', 'en', 'TYPO_CORRECTION', 0.92),
    (2, 'digital camera', 'digital camera', 'en', 'EXACT_SYNONYM', 1.00),
    (2, 'dslr', 'dslr', 'en', 'ABBREVIATION', 0.95),
    (2, 'mirrorless camera', 'mirrorless camera', 'en', 'COMMON_NAME', 0.95),
    (2, 'photo camera', 'photo camera', 'en', 'COMMON_NAME', 0.95),
    (2, 'ካሜራ', 'ካሜራ', 'am', 'EXACT_SYNONYM', 1.00),
    (2, 'ዲጂታል ካሜራ', 'ዲጂታል ካሜራ', 'am', 'EXACT_SYNONYM', 1.00),
    (2, 'kamera', 'kamera', 'am', 'TRANSLITERATION', 0.95),
    (2, 'kaameraa', 'kaameraa', 'om', 'COMMON_NAME', 1.00),

    -- Walkie-Talkie (Entity 3)
    (3, 'walkie-talkie', 'walkietalkie', 'en', 'EXACT_SYNONYM', 1.00),
    (3, 'walkie talkie', 'walkie talkie', 'en', 'EXACT_SYNONYM', 1.00),
    (3, 'walki talki', 'walki talki', 'en', 'TYPO_CORRECTION', 0.92),
    (3, 'two way radio', 'two way radio', 'en', 'COMMON_NAME', 0.95),
    (3, 'radio transceiver', 'radio transceiver', 'en', 'COMMON_NAME', 0.95),
    (3, 'handheld radio', 'handheld radio', 'en', 'COMMON_NAME', 0.90),
    (3, 'ዋኪ ቶኪ', 'ዋኪ ቶኪ', 'am', 'EXACT_SYNONYM', 1.00),
    (3, 'የእጅ ሬዲዮ', 'የእጅ ሬዲዮ', 'am', 'COMMON_NAME', 0.95),
    (3, 'ራዲዮ', 'ራዲዮ', 'am', 'COMMON_NAME', 0.85),

    -- Power Drill (Entity 4)
    (4, 'drill', 'drill', 'en', 'COMMON_NAME', 0.98),
    (4, 'drills', 'drills', 'en', 'COMMON_NAME', 0.98),
    (4, 'dril', 'dril', 'en', 'TYPO_CORRECTION', 0.92),
    (4, 'power drill', 'power drill', 'en', 'EXACT_SYNONYM', 1.00),
    (4, 'electric drill', 'electric drill', 'en', 'COMMON_NAME', 0.98),
    (4, 'impact drill', 'impact drill', 'en', 'COMMON_NAME', 0.95),
    (4, 'cordless drill', 'cordless drill', 'en', 'COMMON_NAME', 0.95),
    (4, 'hammer drill', 'hammer drill', 'en', 'COMMON_NAME', 0.95),
    (4, 'መሰርሰሪያ', 'መሰርሰሪያ', 'am', 'EXACT_SYNONYM', 1.00),
    (4, 'ድሪል', 'ድሪል', 'am', 'EXACT_SYNONYM', 1.00),
    (4, 'የኤሌክትሪክ መሰርሰሪያ', 'የኤሌክትሪክ መሰርሰሪያ', 'am', 'EXACT_SYNONYM', 1.00),
    (4, 'meserseria', 'meserseria', 'am', 'TRANSLITERATION', 0.92),
    (4, 'makiinaa qotuu', 'makiinaa qotuu', 'om', 'COMMON_NAME', 0.95),

    -- Passenger Car (Entity 5)
    (5, 'car', 'car', 'en', 'COMMON_NAME', 1.00),
    (5, 'cars', 'cars', 'en', 'COMMON_NAME', 1.00),
    (5, 'automobile', 'automobile', 'en', 'EXACT_SYNONYM', 1.00),
    (5, 'passenger car', 'passenger car', 'en', 'EXACT_SYNONYM', 1.00),
    (5, 'vehicle', 'vehicle', 'en', 'COMMON_NAME', 0.90),
    (5, 'sedan', 'sedan', 'en', 'COMMON_NAME', 0.92),
    (5, 'mekena', 'mekena', 'am', 'TRANSLITERATION', 0.95),
    (5, 'mekina', 'mekina', 'am', 'TRANSLITERATION', 0.95),
    (5, 'መኪና', 'መኪና', 'am', 'EXACT_SYNONYM', 1.00),
    (5, 'ተሽከርካሪ', 'ተሽከርካሪ', 'am', 'COMMON_NAME', 0.92),
    (5, 'konkolaataa', 'konkolaataa', 'om', 'COMMON_NAME', 1.00),
    (5, 'ሞተር መኪና', 'ሞተር መኪና', 'ti', 'COMMON_NAME', 1.00),

    -- Toyota Vitz (Entity 6)
    (6, 'vitz', 'vitz', 'en', 'MODEL_SPECIFIC', 1.00),
    (6, 'vits', 'vits', 'en', 'TYPO_CORRECTION', 0.95),
    (6, 'toyota vitz', 'toyota vitz', 'en', 'EXACT_SYNONYM', 1.00),
    (6, 'toyota vits', 'toyota vits', 'en', 'TYPO_CORRECTION', 0.95),
    (6, 'vitz mekena', 'vitz mekena', 'am', 'TRANSLITERATION', 0.98),
    (6, 'yaris', 'yaris', 'en', 'EXACT_SYNONYM', 0.90),
    (6, 'ቪትስ', 'ቪትስ', 'am', 'MODEL_SPECIFIC', 1.00),
    (6, 'ቶዮታ ቪትስ', 'ቶዮታ ቪትስ', 'am', 'EXACT_SYNONYM', 1.00),
    (6, 'ቪትስ መኪና', 'ቪትስ መኪና', 'am', 'COMMON_NAME', 1.00),
    (6, 'ቪትዝ', 'ቪትዝ', 'am', 'PHONETIC', 0.95),

    -- Apartment (Entity 7)
    (7, 'apartment', 'apartment', 'en', 'EXACT_SYNONYM', 1.00),
    (7, 'apartments', 'apartments', 'en', 'EXACT_SYNONYM', 1.00),
    (7, 'apartama', 'apartama', 'am', 'TRANSLITERATION', 0.95),
    (7, 'condo', 'condo', 'en', 'INFORMAL', 0.95),
    (7, 'condominium', 'condominium', 'en', 'COMMON_NAME', 0.98),
    (7, 'flat', 'flat', 'en', 'COMMON_NAME', 0.92),
    (7, 'አፓርታማ', 'አፓርታማ', 'am', 'EXACT_SYNONYM', 1.00),
    (7, 'ኮንዶሚኒየም', 'ኮንዶሚኒየም', 'am', 'COMMON_NAME', 0.98),
    (7, 'ኮንዶ', 'ኮንዶ', 'am', 'INFORMAL', 0.92),
    (7, 'የመኖሪያ ቤት', 'የመኖሪያ ቤት', 'am', 'COMMON_NAME', 0.85),
    (7, 'mana gamoo', 'mana gamoo', 'om', 'COMMON_NAME', 0.95),

    -- Plumber (Entity 8)
    (8, 'plumber', 'plumber', 'en', 'EXACT_SYNONYM', 1.00),
    (8, 'plumbers', 'plumbers', 'en', 'EXACT_SYNONYM', 1.00),
    (8, 'pipe fitter', 'pipe fitter', 'en', 'COMMON_NAME', 0.95),
    (8, 'sanitary technician', 'sanitary technician', 'en', 'COMMON_NAME', 0.95),
    (8, 'የቧንቧ ባለሙያ', 'የቧንቧ ባለሙያ', 'am', 'EXACT_SYNONYM', 1.00),
    (8, 'የቧንቧ ሰራተኛ', 'የቧንቧ ሰራተኛ', 'am', 'COMMON_NAME', 0.98),
    (8, 'ቧንቧ ሠራተኛ', 'ቧንቧ ሠራተኛ', 'am', 'COMMON_NAME', 0.98),
    (8, 'ቦምቢስታ', 'ቦምቢስታ', 'am', 'INFORMAL', 0.92),
    (8, 'ogeessa ujummoo', 'ogeessa ujummoo', 'om', 'COMMON_NAME', 1.00),

    -- Plumbing Service (Entity 9)
    (9, 'plumbing service', 'plumbing service', 'en', 'EXACT_SYNONYM', 1.00),
    (9, 'plumbing repair', 'plumbing repair', 'en', 'COMMON_NAME', 0.98),
    (9, 'pipe repair', 'pipe repair', 'en', 'COMMON_NAME', 0.95),
    (9, 'leak fix', 'leak fix', 'en', 'INFORMAL', 0.90),
    (9, 'የቧንቧ ጥገና', 'የቧንቧ ጥገና', 'am', 'EXACT_SYNONYM', 1.00),
    (9, 'የቧንቧ ስራ', 'የቧንቧ ስራ', 'am', 'COMMON_NAME', 0.95),

    -- Software Engineer (Entity 10)
    (10, 'software engineer', 'software engineer', 'en', 'EXACT_SYNONYM', 1.00),
    (10, 'software developer', 'software developer', 'en', 'COMMON_NAME', 0.98),
    (10, 'programmer', 'programmer', 'en', 'COMMON_NAME', 0.95),
    (10, 'coder', 'coder', 'en', 'INFORMAL', 0.90),
    (10, 'backend developer', 'backend developer', 'en', 'COMMON_NAME', 0.95),
    (10, 'frontend developer', 'frontend developer', 'en', 'COMMON_NAME', 0.95),
    (10, 'ሶፍትዌር ኢንጂነር', 'ሶፍትዌር ኢንጂነር', 'am', 'EXACT_SYNONYM', 1.00),
    (10, 'ሶፍትዌር አበልጻጊ', 'ሶፍትዌር አበልጻጊ', 'am', 'COMMON_NAME', 0.98),
    (10, 'ፕሮግራመር', 'ፕሮግራመር', 'am', 'COMMON_NAME', 0.95),

    -- Software Development Service (Entity 11)
    (11, 'software development', 'software development', 'en', 'EXACT_SYNONYM', 1.00),
    (11, 'app development', 'app development', 'en', 'COMMON_NAME', 0.95),
    (11, 'web development', 'web development', 'en', 'COMMON_NAME', 0.95),
    (11, 'website building', 'website building', 'en', 'COMMON_NAME', 0.92),
    (11, 'የሶፍትዌር ማበልፀግ', 'የሶፍትዌር ማበልፀግ', 'am', 'EXACT_SYNONYM', 1.00),
    (11, 'የድረገፅ ስራ', 'የድረገፅ ስራ', 'am', 'COMMON_NAME', 0.92),

    -- House Cleaning Service (Entity 12)
    (12, 'house cleaning', 'house cleaning', 'en', 'EXACT_SYNONYM', 1.00),
    (12, 'home cleaning', 'home cleaning', 'en', 'COMMON_NAME', 0.98),
    (12, 'residential cleaning', 'residential cleaning', 'en', 'COMMON_NAME', 0.98),
    (12, 'maid service', 'maid service', 'en', 'COMMON_NAME', 0.92),
    (12, 'deep cleaning', 'deep cleaning', 'en', 'COMMON_NAME', 0.95),
    (12, 'የቤት ጽዳት', 'የቤት ጽዳት', 'am', 'EXACT_SYNONYM', 1.00),
    (12, 'የቤት ጽዳት አገልግሎት', 'የቤት ጽዳት አገልግሎት', 'am', 'EXACT_SYNONYM', 1.00),
    (12, 'tajaajila qulqullinaa', 'tajaajila qulqullinaa', 'om', 'COMMON_NAME', 0.95),

    -- Portland Cement (Entity 13)
    (13, 'cement', 'cement', 'en', 'COMMON_NAME', 0.98),
    (13, 'portland cement', 'portland cement', 'en', 'EXACT_SYNONYM', 1.00),
    (13, 'hydraulic cement', 'hydraulic cement', 'en', 'COMMON_NAME', 0.95),
    (13, 'opc cement', 'opc cement', 'en', 'ABBREVIATION', 0.95),
    (13, 'ppc cement', 'ppc cement', 'en', 'ABBREVIATION', 0.95),
    (13, 'ሲሚንቶ', 'ሲሚንቶ', 'am', 'EXACT_SYNONYM', 1.00),
    (13, 'ዳንጎቴ ሲሚንቶ', 'ዳንጎቴ ሲሚንቶ', 'am', 'BRAND_REFERENCE', 0.95),
    (13, 'ሙገር ሲሚንቶ', 'ሙገር ሲሚንቶ', 'am', 'BRAND_REFERENCE', 0.95),
    (13, 'simintoo', 'simintoo', 'om', 'COMMON_NAME', 1.00),

    -- Agricultural Tractor (Entity 14)
    (14, 'tractor', 'tractor', 'en', 'COMMON_NAME', 0.98),
    (14, 'tractors', 'tractors', 'en', 'COMMON_NAME', 0.98),
    (14, 'farm tractor', 'farm tractor', 'en', 'COMMON_NAME', 0.98),
    (14, 'agricultural tractor', 'agricultural tractor', 'en', 'EXACT_SYNONYM', 1.00),
    (14, 'ትራክተር', 'ትራክተር', 'am', 'EXACT_SYNONYM', 1.00),
    (14, 'የእርሻ ትራክተር', 'የእርሻ ትራክተር', 'am', 'EXACT_SYNONYM', 1.00),
    (14, 'tiraaktara', 'tiraaktara', 'om', 'COMMON_NAME', 1.00),

    -- Solar Inverter (Entity 15)
    (15, 'solar inverter', 'solar inverter', 'en', 'EXACT_SYNONYM', 1.00),
    (15, 'inverter', 'inverter', 'en', 'COMMON_NAME', 0.95),
    (15, 'inverters', 'inverters', 'en', 'COMMON_NAME', 0.95),
    (15, 'power inverter', 'power inverter', 'en', 'COMMON_NAME', 0.95),
    (15, 'hybrid inverter', 'hybrid inverter', 'en', 'COMMON_NAME', 0.95),
    (15, 'ኢንቨርተር', 'ኢንቨርተር', 'am', 'EXACT_SYNONYM', 1.00),
    (15, 'የፀሐይ ኢንቨርተር', 'የፀሐይ ኢንቨርተር', 'am', 'EXACT_SYNONYM', 1.00),

    -- Public Infrastructure Tender (Entity 16)
    (16, 'tender', 'tender', 'en', 'COMMON_NAME', 0.95),
    (16, 'tenders', 'tenders', 'en', 'COMMON_NAME', 0.95),
    (16, 'construction bid', 'construction bid', 'en', 'COMMON_NAME', 0.95),
    (16, 'government tender', 'government tender', 'en', 'COMMON_NAME', 0.98),
    (16, 'infrastructure tender', 'infrastructure tender', 'en', 'EXACT_SYNONYM', 1.00),
    (16, 'ጨረታ', 'ጨረታ', 'am', 'COMMON_NAME', 1.00),
    (16, 'የጨረታ ማስታወቂያ', 'የጨረታ ማስታወቂያ', 'am', 'EXACT_SYNONYM', 1.00),
    (16, 'caalbaa', 'caalbaa', 'om', 'COMMON_NAME', 1.00),

    -- Freelance Graphic Designer (Entity 17)
    (17, 'graphic designer', 'graphic designer', 'en', 'EXACT_SYNONYM', 1.00),
    (17, 'logo designer', 'logo designer', 'en', 'COMMON_NAME', 0.95),
    (17, 'visual designer', 'visual designer', 'en', 'COMMON_NAME', 0.92),
    (17, 'freelance designer', 'freelance designer', 'en', 'COMMON_NAME', 0.95),
    (17, 'ግራፊክ ዲዛይነር', 'ግራፊክ ዲዛይነር', 'am', 'EXACT_SYNONYM', 1.00),
    (17, 'ዲዛይነር', 'ዲዛይነር', 'am', 'COMMON_NAME', 0.90),
    (17, 'ሎጎ ሰሪ', 'ሎጎ ሰሪ', 'am', 'INFORMAL', 0.90)
ON CONFLICT (entity_id, normalized_alias, language) DO UPDATE SET
    alias_type = EXCLUDED.alias_type,
    confidence_score = EXCLUDED.confidence_score,
    status = 'ACTIVE';

-- 11. Seed Typed Entity Relationship Graph
-- Links Entities to Entities and Entities to Multi-Domain Categories
INSERT INTO entity_relationships (
    source_entity_id, target_entity_id, target_category_id,
    relationship_type, confidence_score, metadata, notes
) VALUES
    -- 1. Toyota Vitz (Entity 6) IS_A Passenger Car (Entity 5)
    (6, 5, NULL, 'IS_A', 1.00, '{"inheritance": "strict_subconcept"}'::jsonb, 'Toyota Vitz is a specific model instance of passenger car'),
    
    -- 2. Toyota Vitz (Entity 6) HAS_BRAND Toyota
    (6, NULL, 14, 'DISCOVERABLE_THROUGH', 1.00, '{"facet": "vehicle_cars"}'::jsonb, 'Discoverable in cars taxonomy'),

    -- 3. Power Drill (Entity 4) USED_FOR Construction and Carpentry
    (4, NULL, 3623, 'DISCOVERABLE_THROUGH', 1.00, '{"primary_domain": "hardware_tools"}'::jsonb, 'Direct hardware tools classification'),
    (4, NULL, 2498, 'DISCOVERABLE_THROUGH', 0.80, '{"secondary_domain": "drilling_equipment"}'::jsonb, 'Cross-domain discovery for industrial drilling equipment'),

    -- 4. Plumber (Entity 8) OFFERED_AS_SERVICE Plumbing Service (Entity 9)
    (8, 9, NULL, 'OFFERED_AS_SERVICE', 1.00, '{"profession_to_service": true}'::jsonb, 'A plumber delivers plumbing services'),
    
    -- 5. Plumbing Service (Entity 9) REQUIRES_SKILL Plumber (Entity 8)
    (9, 8, NULL, 'REQUIRES_SKILL', 1.00, '{"service_to_skill": true}'::jsonb, 'Plumbing service execution requires plumbing trade skills'),

    -- 6. Software Engineer (Entity 10) OFFERED_AS_SERVICE Software Development (Entity 11)
    (10, 11, NULL, 'OFFERED_AS_SERVICE', 1.00, '{"talent_to_service": true}'::jsonb, 'Software engineers provide custom software development services'),

    -- 7. Software Development (Entity 11) REQUIRES_SKILL Software Engineer (Entity 10)
    (11, 10, NULL, 'REQUIRES_SKILL', 1.00, '{"service_to_profession": true}'::jsonb, 'Software development projects require software engineering skills'),

    -- 8. Smartphone (Entity 1) SIMILAR_TO Walkie-Talkie (Entity 3) but contrasted
    (1, 3, NULL, 'SIMILAR_TO', 0.70, '{"domain": "handheld_communication"}'::jsonb, 'Both provide voice communication over radio frequencies'),
    (1, 3, NULL, 'CONTRASTED_WITH', 0.95, '{"key_difference": "Cellular IP vs simplex RF Push-to-Talk"}'::jsonb, 'Smartphone requires cellular carrier; walkie-talkie operates off-grid simplex'),

    -- 9. Smartphone (Entity 1) CONTRASTED_WITH Digital Camera (Entity 2)
    (1, 2, NULL, 'CONTRASTED_WITH', 0.90, '{"key_difference": "General mobile computer vs dedicated optical recording instrument"}'::jsonb, 'Smartphone has camera sensor, but Digital Camera has dedicated optical system and manual sensor controls'),

    -- 10. Agricultural Tractor (Entity 14) USED_IN_INDUSTRY Agriculture & Farming
    (14, NULL, 317, 'DISCOVERABLE_THROUGH', 1.00, '{"primary_domain": "agriculture_tractors"}'::jsonb, 'Direct primary category classification'),

    -- 11. Portland Cement (Entity 13) USED_IN_INDUSTRY Building Construction
    (13, NULL, 1212, 'DISCOVERABLE_THROUGH', 1.00, '{"primary_domain": "building_cement"}'::jsonb, 'Direct construction material category'),

    -- 12. Solar Inverter (Entity 15) COMPLEMENTS Renewable Solar Power Generation
    (15, NULL, 35, 'DISCOVERABLE_THROUGH', 0.95, '{"primary_domain": "tools_hardware_energy"}'::jsonb, 'Solar power electrical components'),

    -- 13. Public Infrastructure Tender (Entity 16) DISCOVERABLE_THROUGH Business Opportunities & Tenders
    (16, NULL, 1020109, 'DISCOVERABLE_THROUGH', 1.00, '{"primary_domain": "tenders_bids_opp"}'::jsonb, 'Commercial bids and public infrastructure contracts'),

    -- 14. Freelance Graphic Designer (Entity 17) DISCOVERABLE_THROUGH Gigs & Media
    (17, NULL, 1020099, 'DISCOVERABLE_THROUGH', 0.95, '{"primary_domain": "freelance_creative_gigs"}'::jsonb, 'Freelance creative design and brand identity services')
ON CONFLICT DO NOTHING;

-- 12. Extra Brand & Typo Aliases
INSERT INTO entity_aliases (entity_id, alias_text, normalized_alias, language, alias_type, confidence_score, status)
VALUES
    (1, 'samsung', 'samsung', 'en', 'BRAND_REFERENCE', 0.95, 'ACTIVE'),
    (1, 'sumsung', 'sumsung', 'en', 'TYPO_CORRECTION', 0.90, 'ACTIVE'),
    (1, 'samsong', 'samsong', 'en', 'TYPO_CORRECTION', 0.90, 'ACTIVE'),
    (1, 'phon', 'phon', 'en', 'TYPO_CORRECTION', 0.90, 'ACTIVE'),
    (6, 'vits', 'vits', 'en', 'TYPO_CORRECTION', 0.90, 'ACTIVE'),
    (4, 'dril', 'dril', 'en', 'TYPO_CORRECTION', 0.90, 'ACTIVE')
ON CONFLICT (entity_id, normalized_alias, language) DO UPDATE SET
    alias_type = EXCLUDED.alias_type,
    confidence_score = EXCLUDED.confidence_score,
    status = 'ACTIVE';

-- 13. Map Existing Database Products to Canonical Entities
UPDATE products SET canonical_entity_id = 1 WHERE canonical_entity_id IS NULL AND (title ILIKE '%phone%' OR title ILIKE '%samsung%' OR title ILIKE '%tecno%' OR title ILIKE '%iphone%');
UPDATE products SET canonical_entity_id = 6 WHERE canonical_entity_id IS NULL AND title ILIKE '%vitz%';
UPDATE products SET canonical_entity_id = 4 WHERE canonical_entity_id IS NULL AND title ILIKE '%drill%';

-- 14. Resilient Multi-Entity Brand Extraction Function
CREATE OR REPLACE FUNCTION extract_brands(p_query text)
RETURNS TABLE(brand_id integer, brand_name character varying, match_position integer)
LANGUAGE plpgsql
AS $$
DECLARE
    v_norm text;
BEGIN
    v_norm := LOWER(TRIM(COALESCE(p_query, '')));
    IF v_norm = '' THEN
        RETURN;
    END IF;

    RETURN QUERY
    SELECT 
        ea.id AS brand_id,
        ea.alias_text AS brand_name,
        POSITION(LOWER(ea.normalized_alias) IN v_norm)::INT AS match_position
    FROM entity_aliases ea
    WHERE ea.status = 'ACTIVE'
      AND (ea.alias_type IN ('BRAND_REFERENCE', 'MODEL_SPECIFIC') OR ea.normalized_alias IN ('samsung', 'apple', 'toyota', 'bosch', 'tecno', 'hyundai', 'lenovo', 'dell', 'hp'))
      AND (
          v_norm LIKE '%' || LOWER(ea.normalized_alias) || '%'
          OR LOWER(ea.normalized_alias) LIKE '%' || v_norm || '%'
      )
    ORDER BY LENGTH(ea.normalized_alias) DESC
    LIMIT 5;
END;
$$;
