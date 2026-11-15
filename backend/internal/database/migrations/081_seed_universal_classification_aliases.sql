-- ============================================================================
-- 081_seed_universal_classification_aliases.sql
-- Ingest and map multilingual synonyms, transliterations, and alternate names
-- into classification_aliases for multi-lingual search & discovery.
-- ============================================================================

BEGIN;

-- 1. Populate classification_aliases from category_aliases where matching concept exists
INSERT INTO classification_aliases (concept_id, placement_id, alias, lang, alias_type, is_active, created_at)
SELECT DISTINCT
    cp.concept_id,
    cp.id AS placement_id,
    ca.alias,
    CASE 
        WHEN ca.language IN ('or', 'om') THEN 'om'
        WHEN ca.language = 'am' THEN 'am'
        WHEN ca.language = 'ti' THEN 'ti'
        ELSE 'en'
    END AS lang,
    'SYNONYM' AS alias_type,
    true AS is_active,
    NOW() AS created_at
FROM category_aliases ca
JOIN category_hierarchy ch ON ch.id = ca.category_id
JOIN classification_placements cp ON cp.slug = ch.slug
WHERE NOT EXISTS (
    SELECT 1 FROM classification_aliases cla 
    WHERE cla.concept_id = cp.concept_id 
      AND cla.alias = ca.alias 
      AND cla.lang = ca.language
);

-- 2. Seed High-Impact Everyday Aliases for Key Sectors & Hardware
-- Bajaj / Rickshaw
INSERT INTO classification_aliases (concept_id, placement_id, alias, lang, alias_type)
SELECT cp.concept_id, cp.id, a.alias, a.lang, a.type
FROM classification_placements cp
CROSS JOIN (VALUES 
    ('bajaj', 'en', 'SYNONYM'),
    ('tuk tuk', 'en', 'ALTERNATIVE_NAME'),
    ('auto rickshaw', 'en', 'SYNONYM'),
    ('ባጃጅ', 'am', 'SYNONYM'),
    ('baajaaj', 'om', 'SYNONYM'),
    ('ባጃጅ', 'ti', 'SYNONYM')
) AS a(alias, lang, type)
WHERE cp.slug LIKE '%bajaj%' OR cp.slug LIKE '%three-wheelers%'
ON CONFLICT DO NOTHING;

-- Power Drills & Rotaries
INSERT INTO classification_aliases (concept_id, placement_id, alias, lang, alias_type)
SELECT cp.concept_id, cp.id, a.alias, a.lang, a.type
FROM classification_placements cp
CROSS JOIN (VALUES 
    ('drill', 'en', 'SYNONYM'),
    ('rotary drill', 'en', 'SYNONYM'),
    ('ድሪል', 'am', 'TRANSLITERATION'),
    ('ቦረቦሬ', 'am', 'SYNONYM'),
    ('dirilii', 'om', 'TRANSLITERATION'),
    ('ድሪል', 'ti', 'TRANSLITERATION')
) AS a(alias, lang, type)
WHERE cp.slug LIKE '%drill%'
ON CONFLICT DO NOTHING;

-- Mobile Phones & Smartphones
INSERT INTO classification_aliases (concept_id, placement_id, alias, lang, alias_type)
SELECT cp.concept_id, cp.id, a.alias, a.lang, a.type
FROM classification_placements cp
CROSS JOIN (VALUES 
    ('smartphone', 'en', 'SYNONYM'),
    ('cell phone', 'en', 'SYNONYM'),
    ('mobile', 'en', 'SYNONYM'),
    ('ስልክ', 'am', 'SYNONYM'),
    ('ሞባይል', 'am', 'TRANSLITERATION'),
    ('bilbila', 'om', 'SYNONYM'),
    ('ሞባይል', 'ti', 'TRANSLITERATION')
) AS a(alias, lang, type)
WHERE cp.slug IN ('mobile-phones', 'smartphones-mobile-devices')
ON CONFLICT DO NOTHING;

-- Laptops & Computers
INSERT INTO classification_aliases (concept_id, placement_id, alias, lang, alias_type)
SELECT cp.concept_id, cp.id, a.alias, a.lang, a.type
FROM classification_placements cp
CROSS JOIN (VALUES 
    ('laptop', 'en', 'SYNONYM'),
    ('notebook pc', 'en', 'SYNONYM'),
    ('ኮምፒውተር', 'am', 'SYNONYM'),
    ('ላፕቶፕ', 'am', 'TRANSLITERATION'),
    ('kompiitara', 'om', 'SYNONYM'),
    ('ኮምፒተር', 'ti', 'SYNONYM')
) AS a(alias, lang, type)
WHERE cp.slug IN ('computers-laptops', 'laptops-notebooks')
ON CONFLICT DO NOTHING;

-- Tractors & Farm Machinery
INSERT INTO classification_aliases (concept_id, placement_id, alias, lang, alias_type)
SELECT cp.concept_id, cp.id, a.alias, a.lang, a.type
FROM classification_placements cp
CROSS JOIN (VALUES 
    ('tractor', 'en', 'SYNONYM'),
    ('farm machine', 'en', 'SYNONYM'),
    ('ትራክተር', 'am', 'TRANSLITERATION'),
    ('tiraaktara', 'om', 'TRANSLITERATION'),
    ('ትራክተር', 'ti', 'TRANSLITERATION')
) AS a(alias, lang, type)
WHERE cp.slug LIKE '%tractor%'
ON CONFLICT DO NOTHING;

COMMIT;

