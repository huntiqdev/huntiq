-- ============================================================================
-- 080_taxonomy_optimization_and_deduplication.sql
-- 1. Remove 18 empty/orphaned early legacy placements with 0 children
-- 2. Correct multilingual translations for Tools, Machinery & Equipment
-- 3. Populate Amharic, Oromo, and Tigrinya translations for domain groupers
-- ============================================================================

BEGIN;

-- 1. Remove the 18 dead/empty legacy placements from early mock seeds (ids 11, 13-17, 23-25, 28-36)
-- These have 0 children, 0 products, and their genuine equivalents exist with full trees.
DELETE FROM classification_placements
WHERE id IN (
    11, 13,                -- tv-audio, accessories (under electronics)
    14, 15, 16, 17,        -- cars, trucks, motorcycles, spare-parts (under vehicles)
    23, 24, 25,            -- crops-coffee, livestock, tractors-equipment (under agriculture)
    28, 29,                -- traditional-habesha, shoes (under fashion)
    30, 31,                -- furniture, kitchen-appliances (under home)
    32, 33,                -- professional-services, repair-maintenance (under services)
    34, 35, 36             -- commercial-equipment, tools-hardware, office-supplies (under tools-machinery)
);

-- Also remove the 4 empty legacy top-level roots (vehicles, fashion, home, services)
-- Genuine equivalents: transport-mobility, clothing-apparel, furniture-furnishings/home-appliances, personal-household-services
DELETE FROM classification_placements
WHERE id IN (2, 5, 6, 7);

-- 2. Correct multilingual translations for "Tools, Machinery & Equipment" (Concept #8)
UPDATE classification_concepts
SET name_am = 'መሣሪያዎች፣ ማሽነሪዎችና መገልገያዎች',
    name_om = 'Meeshaalee, Maashinarii fi Qophiiwwan',
    name_ti = 'መሳሪያታት፣ ማሽነሪታትን ናውትን'
WHERE id = 8 OR code IN ('CONCEPT_TOOLS_MACHINERY', 'tools-machinery') OR canonical_name = 'Tools, Machinery & Equipment';

-- 3. Populate multilingual translations for cross-domain groupers under SYS_TOOLS and SYS_ELEC
UPDATE classification_concepts
SET name_am = REPLACE(name_en, ' – Tools & Equipment', ' – መሣሪያዎችና መገልገያዎች'),
    name_om = REPLACE(name_en, ' – Tools & Equipment', ' – Meeshaalee fi Qophiiwwan'),
    name_ti = REPLACE(name_en, ' – Tools & Equipment', ' – መሳሪያታትን ናውትን')
WHERE code LIKE 'TOOLS_DOMAIN_%' AND name_am IS NULL;

UPDATE classification_concepts
SET name_am = REPLACE(name_en, ' – Electronics & Components', ' – ኤሌክትሮኒክስና ክፍሎች'),
    name_om = REPLACE(name_en, ' – Electronics & Components', ' – Elektirooniksii fi Kutaalee'),
    name_ti = REPLACE(name_en, ' – Electronics & Components', ' – ኤሌክትሮኒክስን ኣካላትን')
WHERE code LIKE 'ELEC_DOMAIN_%' AND name_am IS NULL;

-- 4. Clean up any concepts orphaned by placement removal (using indexed NOT EXISTS)
DELETE FROM classification_concepts c
WHERE NOT EXISTS (SELECT 1 FROM classification_placements p WHERE p.concept_id = c.id)
  AND NOT EXISTS (SELECT 1 FROM classification_relationships r WHERE r.source_concept_id = c.id)
  AND NOT EXISTS (SELECT 1 FROM classification_relationships r WHERE r.target_concept_id = c.id);

-- 5. Recalculate is_leaf flags efficiently
UPDATE classification_placements p
SET is_leaf = false
WHERE is_leaf = true
  AND EXISTS (SELECT 1 FROM classification_placements child WHERE child.parent_placement_id = p.id);

UPDATE classification_placements p
SET is_leaf = true
WHERE is_leaf = false
  AND NOT EXISTS (SELECT 1 FROM classification_placements child WHERE child.parent_placement_id = p.id);

COMMIT;

