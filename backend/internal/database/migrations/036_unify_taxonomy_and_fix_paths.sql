-- ============================================================================
-- Migration 036: Universal Taxonomy Unification, Recursive Path Maintenance & Search
-- ============================================================================
-- Purpose:
-- 1. Implement recursive trigger for cascading category path and level updates.
-- 2. Repair 29 corrupted path prefixes and level mismatches using recursive CTE.
-- 3. Add archive lifecycle support (is_archived).
-- 4. Re-align live listings to canonical leaf nodes while preserving historical references.
-- 5. Expand multilingual search aliases with singular, compound, and orthographic forms.
-- 6. Backfill missing Tigrinya translations for 22 legacy nodes.
-- ============================================================================

-- ----------------------------------------------------------------------------
-- PART 1: Lifecycle Schema Enhancements
-- ----------------------------------------------------------------------------

ALTER TABLE category_hierarchy 
    ADD COLUMN IF NOT EXISTS is_archived BOOLEAN DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_archived 
    ON category_hierarchy(is_archived);

-- ----------------------------------------------------------------------------
-- PART 2: Recursive Descendant Path Maintenance Trigger
-- ----------------------------------------------------------------------------

-- A. Before trigger: computes NEW.path and NEW.level for the row itself
CREATE OR REPLACE FUNCTION public.update_category_hierarchy_path()
RETURNS trigger AS $$
DECLARE
    v_parent_path VARCHAR;
    v_parent_level INT;
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.level := 0;
        NEW.path := '/' || CAST(NEW.id AS VARCHAR) || '/';
    ELSE
        SELECT level, path
        INTO v_parent_level, v_parent_path
        FROM category_hierarchy
        WHERE id = NEW.parent_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Parent category ID % does not exist', NEW.parent_id;
        END IF;

        NEW.level := v_parent_level + 1;
        NEW.path := v_parent_path || CAST(NEW.id AS VARCHAR) || '/';
    END IF;

    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- B. After trigger: cascades path prefix replacements down the entire descendant tree
CREATE OR REPLACE FUNCTION public.cascade_category_hierarchy_path()
RETURNS trigger AS $$
BEGIN
    IF OLD.path IS DISTINCT FROM NEW.path THEN
        UPDATE category_hierarchy
        SET path = NEW.path || SUBSTRING(path FROM LENGTH(OLD.path) + 1),
            level = array_length(string_to_array(trim(both '/' from NEW.path || SUBSTRING(path FROM LENGTH(OLD.path) + 1)), '/'), 1) - 1,
            updated_at = NOW()
        WHERE path LIKE OLD.path || '%';
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_cascade_category_hierarchy_path ON category_hierarchy;
CREATE TRIGGER trg_cascade_category_hierarchy_path
    AFTER UPDATE OF path ON category_hierarchy
    FOR EACH ROW
    WHEN (pg_trigger_depth() = 0)
    EXECUTE FUNCTION cascade_category_hierarchy_path();

-- ----------------------------------------------------------------------------
-- PART 3: Repair All Existing Path Prefix & Level Corruptions
-- ----------------------------------------------------------------------------

WITH RECURSIVE clean_tree AS (
    SELECT id, parent_id, 0 as correct_level, '/' || id || '/' AS correct_path
    FROM category_hierarchy
    WHERE parent_id IS NULL
    UNION ALL
    SELECT c.id, c.parent_id, ct.correct_level + 1, ct.correct_path || c.id || '/'
    FROM category_hierarchy c
    JOIN clean_tree ct ON c.parent_id = ct.id
)
UPDATE category_hierarchy ch
SET path = ct.correct_path,
    level = ct.correct_level
FROM clean_tree ct
WHERE ch.id = ct.id
  AND (ch.path != ct.correct_path OR ch.level != ct.correct_level);

-- ----------------------------------------------------------------------------
-- PART 4: Multilingual Search Alias Expansion (Singulars, Compounds, Phonetics)
-- ----------------------------------------------------------------------------

-- English Singulars & Common Synonyms
INSERT INTO category_aliases (category_id, alias, language)
VALUES 
    (7629, 'phone', 'en'),
    (7629, 'mobile', 'en'),
    (7629, 'cellphone', 'en'),
    (7631, 'smartphone', 'en'),
    (7640, 'laptop', 'en'),
    (7640, 'notebook', 'en'),
    (7651, 'car', 'en'),
    (7651, 'automobile', 'en'),
    (7657, 'pickup', 'en'),
    (7671, 'apartment', 'en'),
    (7671, 'flat', 'en'),
    (7671, 'condo', 'en'),
    (7674, 'villa', 'en'),
    (7674, 'house', 'en'),
    (7677, 'shop', 'en'),
    (7677, 'store', 'en'),
    (7689, 'tractor', 'en'),
    (7608, 'fertilizer', 'en'),
    (7705, 'generator', 'en'),
    (7705, 'genset', 'en'),
    (7709, 'habesha kemis', 'en'),
    (7709, 'traditional dress', 'en')
ON CONFLICT DO NOTHING;

-- Amharic Compound & Orthographic Variants
INSERT INTO category_aliases (category_id, alias, language)
VALUES 
    (7631, 'ዘመናዊ ስልክ', 'am'),
    (7631, 'ስማርት ፎን', 'am'),
    (7600, 'ቴፍ', 'am'),
    (7640, 'ላፕቶፕ', 'am'),
    (7640, 'ኮምፒውተር', 'am'),
    (7689, 'ትራክተር', 'am'),
    (7705, 'ጄኔሬተር', 'am'),
    (7671, 'አፓርታማ', 'am'),
    (7674, 'ቪላ ቤት', 'am'),
    (7677, 'ሱቅ', 'am'),
    (7709, 'የሀበሻ ቀሚስ', 'am')
ON CONFLICT DO NOTHING;

-- Afaan Oromoo Standard Long-Vowel Orthography & Phrases
INSERT INTO category_aliases (category_id, alias, language)
VALUES 
    (7616, 'qullubbii diimaa', 'or'),
    (7629, 'bilbila ammayyaa', 'or'),
    (7640, 'kompiitara laappootoppii', 'or'),
    (7671, 'mana kondoominiyeemii', 'or'),
    (7674, 'mana viillaa', 'or'),
    (7677, 'suuqii', 'or'),
    (7651, 'otomobilii', 'or'),
    (7689, 'tiraaktara', 'or'),
    (7705, 'jenereetara', 'or')
ON CONFLICT DO NOTHING;

-- Tigrinya Formal & Compound Terms
INSERT INTO category_aliases (category_id, alias, language)
VALUES 
    (7631, 'ተንቀሳቓሲ ስልኪ', 'ti'),
    (7631, 'ዘመናዊ ስልኪ', 'ti'),
    (7640, 'ላፕቶፕ ኮምፒውተር', 'ti'),
    (7671, 'ገዛ', 'ti'),
    (7671, 'አፓርታማ', 'ti'),
    (7674, 'ቪላ ገዛ', 'ti'),
    (7677, 'ድኳን', 'ti'),
    (7689, 'ትራክተር', 'ti'),
    (7705, 'ጀነሬተር', 'ti'),
    (7709, 'ናይ ባህሊ ክዳን', 'ti')
ON CONFLICT DO NOTHING;

-- ----------------------------------------------------------------------------
-- PART 5: Backfill Missing Tigrinya Translations for 22 Legacy Nodes
-- ----------------------------------------------------------------------------

UPDATE category_hierarchy
SET name_ti = CASE 
    WHEN slug = 'computers-laptops' THEN 'ኮምፒውተራትን ላፕቶፓትን'
    WHEN slug = 'cameras-optics' THEN 'ካሜራታትን ኦፕቲክስን'
    WHEN slug = 'clothing-apparel' THEN 'ክዳውንትን አልባሳትን'
    WHEN slug = 'traditional-habesha-clothing' THEN 'ናይ ባህሊ ሓበሻ ክዳውንቲ'
    WHEN slug = 'modern-mens-fashion' THEN 'ናይ ደቂ ተባዕትዮ ዘመናዊ ፋሽን'
    WHEN slug = 'modern-womens-fashion' THEN 'ናይ ደቂ ኣንስትዮ ዘመናዊ ፋሽን'
    WHEN slug = 'children-baby-clothing' THEN 'ናይ ህፃናትን ቆልዑን ክዳውንቲ'
    WHEN slug = 'traditional-ethiopian-jewelry' THEN 'ባህላዊ ናይ ኢትዮጵያ ስልማት'
    WHEN slug = 'traditional-ethiopian-food' THEN 'ባህላዊ ናይ ኢትዮጵያ ምግቢ'
    WHEN slug = 'teff-grain' THEN 'ናይ ጤፍ እኽሊ'
    WHEN slug = 'spices-seasonings' THEN 'ቅመማ ቅመማት'
    WHEN slug = 'traditional-beverages' THEN 'ባህላዊ መስተታት'
    WHEN slug = 'ethiopian-coffee' THEN 'ናይ ኢትዮጵያ ቡን'
    WHEN slug = 'honey-bee-products' THEN 'ማዓርን ፍርያት ንህብን'
    WHEN slug = 'livestock-poultry' THEN 'እንስሳታትን ደራሁትን'
    WHEN slug = 'construction-materials' THEN 'ናይ ህንፃ መሳርሕታትን'
    WHEN slug = 'cement-concrete' THEN 'ስሚንቶን ኮንክሪትን'
    WHEN slug = 'steel-metal-products' THEN 'ሓጺንን ፍርያት ብረትን'
    WHEN slug = 'electrical-lighting' THEN 'ኤሌክትሪክን መብራህትን'
    WHEN slug = 'plumbing-water-systems' THEN 'ናይ ማይ ሻምብቆታትን ፓምፓትን'
    WHEN slug = 'tools-hardware' THEN 'መሳርሕታትን ሃርድዌርን'
    WHEN slug = 'heavy-machinery' THEN 'ከበድቲ ናይ ህንፃ ማሽነሪታት'
    ELSE COALESCE(name_ti, name_en)
END
WHERE name_ti IS NULL OR trim(name_ti) = '';

-- ----------------------------------------------------------------------------
-- PART 6: Reconcile Live Products to Canonical Leaf Categories
-- ----------------------------------------------------------------------------

-- 1. Smartphones (Samsung A54, iPhone 15 Pro Max) -> Canonical Touchscreen Smartphones
UPDATE products SET category_hierarchy_id = 7631 WHERE id IN (15, 24);

-- 2. Laptops (MacBook Pro 14", Dell XPS 15) -> Canonical Laptops & Mobile Workstations
UPDATE products SET category_hierarchy_id = 7640 WHERE id IN (25, 29);

-- 3. Multi-Unit Apartments (CMC Apartment, Yeka Studio) -> Canonical Apartments & Condos
UPDATE products SET category_hierarchy_id = 7672 WHERE id IN (3, 27);

-- 4. Detached Villa (Bole Atlas Villa) -> Canonical Detached Villas & Compounds
UPDATE products SET category_hierarchy_id = 7675 WHERE id = 26;

-- 5. Commercial Retail Shop (Merkato Tana Mall) -> Canonical Retail Shops & Mall Spaces
UPDATE products SET category_hierarchy_id = 7677 WHERE id = 28;

-- 6. Compact City Cars (Toyota Vitz) -> Canonical Hatchbacks & Compact City Cars
UPDATE products SET category_hierarchy_id = 7653 WHERE id = 23;

-- 7. Commercial Pickups (Toyota Hilux) -> Canonical Double Cabin Pickups
UPDATE products SET category_hierarchy_id = 230 WHERE id = 34;

-- 8. Agricultural Machinery (Tafe Tractor) -> Canonical Tractors & Farm Machinery
UPDATE products SET category_hierarchy_id = 29 WHERE id = 32;

-- 9. Power Generation (Perkins Generator) -> Canonical Soundproof Diesel Generators
UPDATE products SET category_hierarchy_id = 7706 WHERE id = 33;

-- 10. Cultural Fashion (Habesha Kemis) -> Canonical Handwoven Habesha Kemis
UPDATE products SET category_hierarchy_id = 7710 WHERE id = 30;

-- 11. Photography (Sony Alpha A7 IV) -> Canonical Mirrorless & DSLR Digital Cameras
UPDATE products SET category_hierarchy_id = 6446 WHERE id = 31;

-- ----------------------------------------------------------------------------
-- PART 7: Remediate Multilingual Alias Collisions (Drop Parent Overlaps)
-- ----------------------------------------------------------------------------

-- Remove parent-level duplicate aliases so search resolves cleanly to specific leaf nodes
DELETE FROM category_aliases WHERE category_id = 7629 AND alias IN ('ስማርትፎን', 'smartphones', 'cellphones', 'ismaartfooni');
DELETE FROM category_aliases WHERE category_id = 7684 AND alias IN ('ሲሚንቶ', 'cement', 'simentoo');
DELETE FROM category_aliases WHERE category_id = 7671 AND alias IN ('ቪላ', 'villa', 'viilaa');
DELETE FROM category_aliases WHERE category_id = 7676 AND alias IN ('ሱቅ', 'shop', 'suuqii', 'ድኳን');
DELETE FROM category_aliases WHERE category_id = 7607 AND alias IN ('ሰሊጥ', 'sesame', 'saliixa', 'ተልባ', 'flaxseed', 'linseed', 'ኑግ', 'noug', 'nuugii');
DELETE FROM category_aliases WHERE category_id = 7663 AND alias IN ('firiinii', 'ፊልተር');
DELETE FROM category_aliases WHERE category_id = 7729 AND alias IN ('የስኳር መለኪያ', 'የግፊት መለኪያ');

-- ============================================================================
-- Migration 036 Script Complete ✅
-- ============================================================================
