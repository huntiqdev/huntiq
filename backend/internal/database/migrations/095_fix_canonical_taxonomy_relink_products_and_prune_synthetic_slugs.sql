-- ============================================================================
-- 095_fix_canonical_taxonomy_relink_products_and_prune_synthetic_slugs.sql
-- Canonical Taxonomy Normalization, Product Relinking & Synthetic Slug Deactivation
--
-- Goals:
--   1. Deactivate all 621,850 synthetic permutation categories (slugs starting
--      with 'leaf-', 'var-', 'spec-', 'cat-', 'subcat-') so they never pollute
--      search results, faceted navigation, or category trees.
--   2. Activate core roots 1..8 and all 4,005 human-readable canonical categories.
--   3. Ensure missing canonical leaf categories exist with clean semantic slugs:
--        - apple-macbook-pro-14-m3-pro
--        - toyota-vitz-automatic-hatchback
--        - bole-atlas-furnished-2bed-apartment
--   4. Relink all 11 active products dynamically by their canonical category slug
--      fixing misassignments:
--        - Sony 70-200mm GM II -> sony-fe-70-200mm-f28-gm-oss-ii-telephoto-lens
--        - Bosch GBH 2-26 DRE Drill -> bosch-professional-gbh-2-26-dre-800w-drill
--        - Apple MacBook Pro 14 -> apple-macbook-pro-14-m3-pro
--        - Toyota Vitz 2018 -> toyota-vitz-automatic-hatchback
--        - Bole Atlas 2-Bed -> bole-atlas-furnished-2bed-apartment
--        - iPhone 15 Pro Max -> apple-iphone-15-pro-max-256gb-natural-titanium
--        - Samsung S24 Ultra -> samsung-galaxy-s24-ultra-5g-256gb
--   5. Synchronize and enrich products.attributes JSONB for all listings.
-- ============================================================================

BEGIN;

-- ----------------------------------------------------------------------------
-- 1. Deactivate All Synthetic Categories
-- ----------------------------------------------------------------------------
UPDATE category_hierarchy
SET is_active = false
WHERE slug LIKE 'leaf-%'
   OR slug LIKE 'var-%'
   OR slug LIKE 'spec-%'
   OR slug LIKE 'cat-%'
   OR slug LIKE 'subcat-%';

-- ----------------------------------------------------------------------------
-- 2. Activate Core Roots 1..8 and All Clean Canonical Hierarchy Branches
-- ----------------------------------------------------------------------------
UPDATE category_hierarchy
SET is_active = true
WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8);

-- Activate clean categories whose slugs do not have machine-generated prefixes
UPDATE category_hierarchy
SET is_active = true
WHERE slug NOT LIKE 'leaf-%'
  AND slug NOT LIKE 'var-%'
  AND slug NOT LIKE 'spec-%'
  AND slug NOT LIKE 'cat-%'
  AND slug NOT LIKE 'subcat-%'
  AND (
    level <= 2
    OR path ~ '^/(1|2|3|4|5|6|7|8)/'
  );

-- Specifically ensure tool and machinery subcategory chain is active
UPDATE category_hierarchy
SET is_active = true
WHERE id IN (8, 35, 3536, 3553, 3570, 3587, 3611);

-- ----------------------------------------------------------------------------
-- 3. Ensure Missing Canonical Leaf Categories Exist with Clean Semantic Slugs
-- ----------------------------------------------------------------------------

-- 3a. Apple MacBook Pro 14 M3 Pro (under Business & Professional Ultrabooks 3545)
DO $$
DECLARE
    v_parent_id INT := 3545;
    v_path VARCHAR(255);
    v_level INT;
BEGIN
    SELECT path, level INTO v_path, v_level FROM category_hierarchy WHERE id = v_parent_id;
    IF v_path IS NOT NULL AND NOT EXISTS (SELECT 1 FROM category_hierarchy WHERE slug = 'apple-macbook-pro-14-m3-pro') THEN
        INSERT INTO category_hierarchy (
            parent_id, name_en, name_am, name_om, name_ti, slug,
            level, path, is_leaf, is_active, display_order
        ) VALUES (
            v_parent_id,
            'Apple MacBook Pro 14 M3 Pro (18GB RAM, 512GB SSD)',
            'አፕል ማክቡክ ፕሮ 14 M3 ፕሮ',
            'Apple MacBook Pro 14 M3 Pro',
            'ኣፕል ማክቡክ ፕሮ 14 M3 ፕሮ',
            'apple-macbook-pro-14-m3-pro',
            v_level + 1,
            v_path, -- trigger or fixup will append id
            true, true, 1
        );
    END IF;
END $$;

-- 3b. Toyota Vitz 2018 Clean Condition Automatic (under Cars & Automobiles 14)
DO $$
DECLARE
    v_parent_id INT := 14;
    v_path VARCHAR(255);
    v_level INT;
BEGIN
    SELECT path, level INTO v_path, v_level FROM category_hierarchy WHERE id = v_parent_id;
    IF v_path IS NOT NULL AND NOT EXISTS (SELECT 1 FROM category_hierarchy WHERE slug = 'toyota-vitz-automatic-hatchback') THEN
        INSERT INTO category_hierarchy (
            parent_id, name_en, name_am, name_om, name_ti, slug,
            level, path, is_leaf, is_active, display_order
        ) VALUES (
            v_parent_id,
            'Toyota Vitz Automatic Hatchback',
            'ቶዮታ ቪትዝ አውቶማቲክ ሃችባክ',
            'Toyota Vitz Otomaatikii',
            'ቶዮታ ቪትዝ ኣውቶማቲክ ሃችባክ',
            'toyota-vitz-automatic-hatchback',
            v_level + 1,
            v_path,
            true, true, 1
        );
    END IF;
END $$;

-- 3c. Bole Atlas 2-Bedroom Furnished Apartment (under Two-Bedroom Furnished Apartments 3550)
DO $$
DECLARE
    v_parent_id INT := 3550;
    v_path VARCHAR(255);
    v_level INT;
BEGIN
    SELECT path, level INTO v_path, v_level FROM category_hierarchy WHERE id = v_parent_id;
    IF v_path IS NOT NULL AND NOT EXISTS (SELECT 1 FROM category_hierarchy WHERE slug = 'bole-atlas-furnished-2bed-apartment') THEN
        INSERT INTO category_hierarchy (
            parent_id, name_en, name_am, name_om, name_ti, slug,
            level, path, is_leaf, is_active, display_order
        ) VALUES (
            v_parent_id,
            'Bole Atlas Luxury 2-Bed Furnished Apartment',
            'ቦሌ አትላስ 2-መኝታ ሙሉ እቃ የተሟላለት አፓርታማ',
            'Bole Atlas Kutaa 2 kan Meeshaan Guutame',
            'ቦሌ ኣትላስ 2-መደቀሲ ምሉእ ንዋይ ዘለዎ ኣፓርታማ',
            'bole-atlas-furnished-2bed-apartment',
            v_level + 1,
            v_path,
            true, true, 1
        );
    END IF;
END $$;

-- Fix paths for any newly created categories
UPDATE category_hierarchy c
SET path = p.path || c.id || '/',
    level = p.level + 1,
    updated_at = NOW()
FROM category_hierarchy p
WHERE c.parent_id = p.id
  AND (c.path IS DISTINCT FROM (p.path || c.id || '/') OR c.level IS DISTINCT FROM (p.level + 1));

-- ----------------------------------------------------------------------------
-- 4. Dynamically Relink All Active Products to Canonical Slugs
-- ----------------------------------------------------------------------------
DO $$
DECLARE
    v_cat_id INT;
BEGIN
    -- 1. Sony FE 24-70mm GM II
    SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'sony-fe-24-70mm-f28-gm-ii-lens' LIMIT 1;
    IF v_cat_id IS NOT NULL THEN
        UPDATE products SET category_hierarchy_id = v_cat_id WHERE title ILIKE '%Sony FE 24-70mm%';
    END IF;

    -- 2. Canon RF 50mm f/1.8 STM
    SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'canon-rf-50mm-f18-stm-compact-prime-lens' LIMIT 1;
    IF v_cat_id IS NOT NULL THEN
        UPDATE products SET category_hierarchy_id = v_cat_id WHERE title ILIKE '%Canon RF 50mm%';
    END IF;

    -- 3. Sony FE 70-200mm GM OSS II (Fixes misassignment to Canon RF)
    SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'sony-fe-70-200mm-f28-gm-oss-ii-telephoto-lens' LIMIT 1;
    IF v_cat_id IS NOT NULL THEN
        UPDATE products SET category_hierarchy_id = v_cat_id WHERE title ILIKE '%Sony FE 70-200mm%';
    END IF;

    -- 4. Sony Alpha 7 IV Mirrorless Body
    SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'sony-alpha-7-iv-a7-iv-mirrorless-body' LIMIT 1;
    IF v_cat_id IS NOT NULL THEN
        UPDATE products SET category_hierarchy_id = v_cat_id WHERE title ILIKE '%Sony Alpha 7 IV%';
    END IF;

    -- 5. Canon EOS 5D Mark IV DSLR Body
    SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'canon-eos-5d-mark-iv-dslr-body' LIMIT 1;
    IF v_cat_id IS NOT NULL THEN
        UPDATE products SET category_hierarchy_id = v_cat_id WHERE title ILIKE '%Canon EOS 5D Mark IV%';
    END IF;

    -- 6. Apple iPhone 15 Pro Max
    SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'apple-iphone-15-pro-max-256gb-natural-titanium' LIMIT 1;
    IF v_cat_id IS NOT NULL THEN
        UPDATE products SET category_hierarchy_id = v_cat_id WHERE title ILIKE '%iPhone 15 Pro Max%';
    END IF;

    -- 7. Samsung Galaxy S24 Ultra
    SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'samsung-galaxy-s24-ultra-5g-256gb' LIMIT 1;
    IF v_cat_id IS NOT NULL THEN
        UPDATE products SET category_hierarchy_id = v_cat_id WHERE title ILIKE '%Samsung Galaxy S24 Ultra%';
    END IF;

    -- 8. Apple MacBook Pro 14
    SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'apple-macbook-pro-14-m3-pro' LIMIT 1;
    IF v_cat_id IS NOT NULL THEN
        UPDATE products SET category_hierarchy_id = v_cat_id WHERE title ILIKE '%MacBook Pro 14%';
    END IF;

    -- 9. Toyota Vitz 2018
    SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'toyota-vitz-automatic-hatchback' LIMIT 1;
    IF v_cat_id IS NOT NULL THEN
        UPDATE products SET category_hierarchy_id = v_cat_id WHERE title ILIKE '%Toyota Vitz%';
    END IF;

    -- 10. Modern 2-Bed Bole Apartment
    SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'bole-atlas-furnished-2bed-apartment' LIMIT 1;
    IF v_cat_id IS NULL THEN
        SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'bole-medhanialem-luxury-2bed-furnished-apt' LIMIT 1;
    END IF;
    IF v_cat_id IS NOT NULL THEN
        UPDATE products SET category_hierarchy_id = v_cat_id WHERE title ILIKE '%Furnished Apartment in Bole%';
    END IF;

    -- 11. Bosch GBH 2-26 DRE Rotary Hammer (Fixes synthetic 625225)
    SELECT id INTO v_cat_id FROM category_hierarchy WHERE slug = 'bosch-professional-gbh-2-26-dre-800w-drill' LIMIT 1;
    IF v_cat_id IS NOT NULL THEN
        UPDATE products SET category_hierarchy_id = v_cat_id WHERE title ILIKE '%Bosch GBH 2-26 DRE%';
    END IF;
END $$;

-- ----------------------------------------------------------------------------
-- 5. Enrich and Normalize Structured JSONB Attributes
-- ----------------------------------------------------------------------------
UPDATE products
SET attributes = attributes || '{
    "brand": "Bosch",
    "model": "GBH 2-26 DRE Professional",
    "tool_category": "Rotary Hammer",
    "chuck_type": "SDS-Plus",
    "power_watts": 800,
    "impact_energy_joules": 2.7,
    "voltage": "220V-240V",
    "corded": true
}'::jsonb
WHERE title ILIKE '%Bosch GBH 2-26 DRE%';

UPDATE products
SET attributes = attributes || '{
    "brand": "Apple",
    "model": "MacBook Pro 14",
    "processor": "Apple M3 Pro",
    "ram_gb": 18,
    "storage_gb": 512,
    "screen_size_inches": 14.2,
    "display_type": "Liquid Retina XDR",
    "color": "Space Black"
}'::jsonb
WHERE title ILIKE '%MacBook Pro 14%';

UPDATE products
SET attributes = attributes || '{
    "make": "Toyota",
    "model": "Vitz",
    "year": 2018,
    "fuel_type": "Petrol",
    "transmission": "Automatic",
    "body_type": "Hatchback",
    "origin": "Japan"
}'::jsonb
WHERE title ILIKE '%Toyota Vitz%';

UPDATE products
SET attributes = attributes || '{
    "property_type": "Apartment",
    "bedrooms": 2,
    "bathrooms": 2,
    "furnished": true,
    "floor_area_sqm": 120,
    "generator_backup": true,
    "elevator": true,
    "parking": true,
    "security": "24/7 Guard & CCTV"
}'::jsonb
WHERE title ILIKE '%Furnished Apartment in Bole%';

COMMIT;

