-- ============================================================================
-- 091_camera_lenses_photography_taxonomy_and_attributes.sql
-- Complete Photography & Camera Lenses Taxonomy (Layers 3 to 6) & Dynamic Specifications
--
-- 1. Inserts comprehensive technical attribute definitions for photography & optics:
--    - lens-mount (Select: Sony E, Canon RF, Canon EF, Nikon Z, Nikon F, Fuji X, MFT, Leica L, PL)
--    - focal-length (Text: mm e.g. 50mm, 24-70mm)
--    - max-aperture (Select: f/1.2, f/1.4, f/1.8, f/2.0, f/2.8, f/3.5-5.6, f/4.0, f/4.5-6.3, f/5.6)
--    - lens-type (Select: Prime, Standard Zoom, Telephoto, Wide-Angle, Macro, Cinema, Fisheye)
--    - image-stabilization (Boolean: Optical IS / OIS / VR)
--    - filter-thread-size (Select: 49mm to 95mm)
--    - camera-sensor-size (Select: Full Frame 35mm, APS-C, Micro Four Thirds, Medium Format)
--    - camera-megapixels (Number: MP)
--    - camera-shutter-count (Number: count)
--    - lens-optical-condition (Select: Clean Glass, Minor Dust, Micro-Scratches, etc.)
--
-- 2. Links attributes in category_attributes and classification_attribute_rules
--    with is_inherited = true for Camera Lenses (1380), DSLRs (1381), and Mirrorless (1382).
--
-- 3. Purges dummy telecom template descendants under camera categories.
--
-- 4. Expands authentic photography taxonomy down through Layer 6:
--    - Prime & Portrait Lenses, Standard & Travel Zooms, Telephoto, Wide-Angle, Macro, Cinema
--    - Authentic leaf items: Sony, Canon, Nikon, Sigma, Tamron lenses and camera bodies.
--
-- 5. Seeds multilingual search aliases in category_aliases across en, am, om, ti.
-- ============================================================================

BEGIN;

-- ----------------------------------------------------------------------------
-- 1. Insert Photography & Optics Attribute Definitions
-- ----------------------------------------------------------------------------
INSERT INTO attribute_definitions (slug, name_en, name_am, name_om, name_ti, data_type, validation_rules, unit, is_required, is_searchable, is_filterable, display_order)
VALUES
('lens-mount', 'Lens Mount', 'የሌንስ ማያያዣ', 'Qabannoo Leensii', 'መተሓሓዚ ሌንስ', 'select', '{"options": ["Sony E-Mount", "Canon RF", "Canon EF", "Nikon Z", "Nikon F", "Fujifilm X", "Micro Four Thirds (MFT)", "Leica L-Mount", "Pentax K", "PL Mount"]}'::jsonb, NULL, true, true, true, 3),
('focal-length', 'Focal Length', 'የፎካል ርዝመት', 'Dheerina Fookaalii', 'ንውሓት ፎካል', 'text', '{"placeholder": "e.g. 50mm or 24-70mm"}'::jsonb, 'mm', true, true, true, 4),
('max-aperture', 'Maximum Aperture', 'ከፍተኛ አፐርቸር', 'Aparchara Ol''aanaa', 'ዝለዓለ ኣፐርቸር', 'select', '{"options": ["f/1.2", "f/1.4", "f/1.8", "f/2.0", "f/2.8", "f/3.5-5.6", "f/4.0", "f/4.5-6.3", "f/5.6"]}'::jsonb, NULL, true, true, true, 5),
('lens-type', 'Lens Type', 'የሌንስ ዓይነት', 'Gosa Leensii', 'ዓይነት ሌንስ', 'select', '{"options": ["Prime Lens", "Standard Zoom", "Telephoto Zoom", "Wide-Angle Zoom", "Macro Lens", "Cinema Lens", "Fisheye Lens", "Pancake Lens"]}'::jsonb, NULL, true, true, true, 6),
('image-stabilization', 'Image Stabilization', 'የምስል ማረጋጊያ (IS/OIS/VR)', 'Tasgabbooftuu Fakkii', 'መተኣማመኒ ስእሊ', 'boolean', '{"description": "Optical image stabilization built-in"}'::jsonb, NULL, false, true, true, 7),
('filter-thread-size', 'Filter Thread Size', 'የፊልተር ክር መጠን', 'Hammamtaa Filteerii', 'መጠን ፈትሊ ፊልተር', 'select', '{"options": ["49mm", "52mm", "55mm", "58mm", "62mm", "67mm", "72mm", "77mm", "82mm", "95mm", "Rear Drop-in"]}'::jsonb, 'mm', false, true, true, 8),
('camera-sensor-size', 'Sensor Format', 'የሴንሰር ቅርጸት', 'Boca Seensarii', 'ቅርጺ ሴንሰር', 'select', '{"options": ["Full Frame (35mm)", "APS-C / Crop Sensor", "Micro Four Thirds", "Medium Format", "1-Inch"]}'::jsonb, NULL, false, true, true, 9),
('camera-megapixels', 'Megapixels', 'ሜጋፒክሰል', 'Meegaapiksilii', 'ሜጋፒክሰል', 'number', '{"min": 1, "max": 200, "step": 1}'::jsonb, 'MP', false, true, true, 10),
('camera-shutter-count', 'Shutter Count', 'የሻተር ብዛት', 'Lakkoofsa Shutter', 'ቑፅሪ ሻተር', 'number', '{"min": 0, "max": 1000000, "step": 1}'::jsonb, NULL, false, false, true, 11),
('lens-optical-condition', 'Optical Condition', 'የኦፕቲካል ሁኔታ', 'Haala Optikaalaa', 'ኩነታት ኦፕቲካል', 'select', '{"options": ["Clean Glass (No Scratch/Fungus/Haze)", "Minor Dust (Does Not Affect Image)", "Surface Micro-Scratches", "Has Fungus / Cleaning Mark"]}'::jsonb, NULL, false, false, true, 12)
ON CONFLICT (slug) DO UPDATE SET
  name_en = EXCLUDED.name_en,
  name_am = EXCLUDED.name_am,
  name_om = EXCLUDED.name_om,
  name_ti = EXCLUDED.name_ti,
  data_type = EXCLUDED.data_type,
  validation_rules = EXCLUDED.validation_rules,
  unit = EXCLUDED.unit,
  updated_at = NOW();

-- ----------------------------------------------------------------------------
-- 2. Link Attributes in category_attributes (with is_inherited = true)
-- ----------------------------------------------------------------------------
-- Camera Lenses (ID 1380)
INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1380, id, true, true, 1 FROM attribute_definitions WHERE slug = 'brand'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 1;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1380, id, true, true, 2 FROM attribute_definitions WHERE slug = 'model'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 2;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1380, id, true, true, 3 FROM attribute_definitions WHERE slug = 'lens-mount'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 3;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1380, id, true, true, 4 FROM attribute_definitions WHERE slug = 'lens-type'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 4;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1380, id, true, true, 5 FROM attribute_definitions WHERE slug = 'focal-length'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 5;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1380, id, true, true, 6 FROM attribute_definitions WHERE slug = 'max-aperture'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 6;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1380, id, false, true, 7 FROM attribute_definitions WHERE slug = 'camera-sensor-size'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 7;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1380, id, false, true, 8 FROM attribute_definitions WHERE slug = 'image-stabilization'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 8;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1380, id, false, true, 9 FROM attribute_definitions WHERE slug = 'filter-thread-size'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 9;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1380, id, false, true, 10 FROM attribute_definitions WHERE slug = 'lens-optical-condition'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 10;

-- Mirrorless Cameras (ID 1382)
INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1382, id, true, true, 1 FROM attribute_definitions WHERE slug = 'brand'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 1;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1382, id, true, true, 2 FROM attribute_definitions WHERE slug = 'model'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 2;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1382, id, true, true, 3 FROM attribute_definitions WHERE slug = 'lens-mount'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 3;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1382, id, true, true, 4 FROM attribute_definitions WHERE slug = 'camera-sensor-size'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 4;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1382, id, true, true, 5 FROM attribute_definitions WHERE slug = 'camera-megapixels'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 5;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1382, id, false, true, 6 FROM attribute_definitions WHERE slug = 'camera-shutter-count'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 6;

-- DSLR Cameras (ID 1381)
INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1381, id, true, true, 1 FROM attribute_definitions WHERE slug = 'brand'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 1;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1381, id, true, true, 2 FROM attribute_definitions WHERE slug = 'model'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 2;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1381, id, true, true, 3 FROM attribute_definitions WHERE slug = 'lens-mount'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 3;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1381, id, true, true, 4 FROM attribute_definitions WHERE slug = 'camera-sensor-size'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 4;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1381, id, true, true, 5 FROM attribute_definitions WHERE slug = 'camera-megapixels'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 5;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1381, id, false, true, 6 FROM attribute_definitions WHERE slug = 'camera-shutter-count'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 6;

-- Action Cameras (ID 1379)
INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1379, id, true, true, 1 FROM attribute_definitions WHERE slug = 'brand'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 1;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1379, id, true, true, 2 FROM attribute_definitions WHERE slug = 'model'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 2;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1379, id, false, true, 3 FROM attribute_definitions WHERE slug = 'camera-megapixels'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 3;

-- Tripods & Mounts (ID 1383)
INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1383, id, true, true, 1 FROM attribute_definitions WHERE slug = 'brand'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 1;

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1383, id, true, true, 2 FROM attribute_definitions WHERE slug = 'model'
ON CONFLICT (category_id, attribute_id) DO UPDATE SET is_inherited = true, display_order = 2;

-- ----------------------------------------------------------------------------
-- 3. Link Classification Attribute Rules for Placement 1464 (camera-lenses)
-- ----------------------------------------------------------------------------
INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1464, id, true, true, 1 FROM attribute_definitions WHERE slug = 'brand'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1464, id, true, true, 2 FROM attribute_definitions WHERE slug = 'model'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1464, id, true, true, 3 FROM attribute_definitions WHERE slug = 'lens-mount'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1464, id, true, true, 4 FROM attribute_definitions WHERE slug = 'lens-type'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1464, id, true, true, 5 FROM attribute_definitions WHERE slug = 'focal-length'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1464, id, true, true, 6 FROM attribute_definitions WHERE slug = 'max-aperture'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1464, id, false, true, 7 FROM attribute_definitions WHERE slug = 'camera-sensor-size'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1464, id, false, true, 8 FROM attribute_definitions WHERE slug = 'image-stabilization'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1464, id, false, true, 9 FROM attribute_definitions WHERE slug = 'filter-thread-size'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1464, id, false, true, 10 FROM attribute_definitions WHERE slug = 'lens-optical-condition'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

-- Also for placement 1465 (dslr-cameras) and 1466 (mirrorless-cameras)
INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1466, id, true, true, 1 FROM attribute_definitions WHERE slug = 'brand'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1466, id, true, true, 2 FROM attribute_definitions WHERE slug = 'model'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1466, id, true, true, 3 FROM attribute_definitions WHERE slug = 'lens-mount'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1466, id, true, true, 4 FROM attribute_definitions WHERE slug = 'camera-sensor-size'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1466, id, true, true, 5 FROM attribute_definitions WHERE slug = 'camera-megapixels'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

INSERT INTO classification_attribute_rules (placement_id, attribute_id, is_required, is_inherited, display_order)
SELECT 1466, id, false, true, 6 FROM attribute_definitions WHERE slug = 'camera-shutter-count'
ON CONFLICT (placement_id, attribute_id) DO NOTHING;

-- ----------------------------------------------------------------------------
-- 4. Purge Dummy Template Descendants
-- ----------------------------------------------------------------------------
DELETE FROM category_hierarchy
WHERE (path LIKE '/1/1378/1379/%'
   OR path LIKE '/1/1378/1380/%'
   OR path LIKE '/1/1378/1381/%'
   OR path LIKE '/1/1378/1382/%'
   OR path LIKE '/1/1378/1383/%')
  AND id NOT IN (1379, 1380, 1381, 1382, 1383);

-- Update display order for primary photography categories so they appear first
UPDATE category_hierarchy SET display_order = 1 WHERE id = 1380; -- Camera Lenses
UPDATE category_hierarchy SET display_order = 2 WHERE id = 1382; -- Mirrorless Cameras
UPDATE category_hierarchy SET display_order = 3 WHERE id = 1381; -- DSLR Cameras
UPDATE category_hierarchy SET display_order = 4 WHERE id = 1379; -- Action Cameras
UPDATE category_hierarchy SET display_order = 5 WHERE id = 1383; -- Tripods & Mounts

-- ----------------------------------------------------------------------------
-- 5. Seed Authentic Photography & Lens Taxonomy (Layers 3, 4, 5, 6)
-- ----------------------------------------------------------------------------
DO $$
DECLARE
    -- Subcategory IDs (Level 3)
    v_l3_prime INT;
    v_l3_zoom INT;
    v_l3_tele INT;
    v_l3_wide INT;
    v_l3_macro INT;
    v_l3_cine INT;
    v_l3_adapt INT;

    -- Mirrorless Subcategory IDs (Level 3)
    v_l3_ml_ff INT;
    v_l3_ml_apsc INT;
    v_l3_ml_mft INT;

    -- DSLR Subcategory IDs (Level 3)
    v_l3_dslr_ff INT;
    v_l3_dslr_apsc INT;

    -- Spec IDs (Level 4)
    v_l4_fast_50 INT;
    v_l4_portrait_85 INT;
    v_l4_street_35 INT;
    v_l4_zoom_2470 INT;
    v_l4_zoom_travel INT;
    v_l4_zoom_kit INT;
    v_l4_tele_70200 INT;
    v_l4_tele_wildlife INT;
    v_l4_wide_1635 INT;
    v_l4_wide_primes INT;
    v_l4_macro_100 INT;
    v_l4_cine_primes INT;
    v_l4_smart_adapters INT;
    v_l4_filters INT;

    v_l4_ml_flagship INT;
    v_l4_ml_hybrid INT;
    v_l4_dslr_action INT;
    v_l4_dslr_studio INT;

    -- Variant IDs (Level 5)
    v_l5_sony_e INT;
    v_l5_canon_rf INT;
    v_l5_nikon_z INT;
    v_l5_tamron_e INT;
    v_l5_sigma_e INT;
    v_l5_fuji_x INT;

    v_l5_ml_sony INT;
    v_l5_ml_canon INT;
    v_l5_ml_nikon INT;
    v_l5_ml_fuji INT;
    v_l5_dslr_canon INT;
    v_l5_dslr_nikon INT;
BEGIN
    -- ========================================================================
    -- LEVEL 3: SUBCATEGORIES UNDER CAMERA LENSES (1380)
    -- ========================================================================
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (1380, 'Prime & Portrait Lenses', 'ፕራይም እና የቁም ፎቶ ሌንሶች', 'Leensii Dhaabbataa fi Suuraa Fuulaa', 'ፕራይምን ስእሊ ሰብን ሌንስታት', 'prime-portrait-lenses', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l3_prime;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (1380, 'Zoom & Standard Lenses', 'ዙም እና መደበኛ ሌንሶች', 'Leensii Zoomii fi Baramaa', 'ዙምን ስሩዕን ሌንስታት', 'zoom-standard-lenses', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l3_zoom;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (1380, 'Telephoto & Super-Telephoto Lenses', 'ቴሌፎቶ እና ሱፐር ቴሌፎቶ ሌንሶች', 'Leensii Faaggaa fi Telephoto', 'ቴሌፎቶን ሱፐር ቴሌፎቶን ሌንስታት', 'telephoto-super-telephoto-lenses', true, false, 3)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l3_tele;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (1380, 'Wide-Angle & Ultra-Wide Lenses', 'ሰፊ አንግል እና እጅግ ሰፊ ሌንሶች', 'Leensii Bal''ina Kallattii', 'ሰፊሕ አንግልን ኣዝዩ ሰፊሕን ሌንስታት', 'wide-angle-ultra-wide-lenses', true, false, 4)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l3_wide;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (1380, 'Macro & Close-Up Lenses', 'ማክሮ እና የቅርብ ፎቶ ሌንሶች', 'Leensii Maakroo fi Dhiyoo', 'ማክሮን ናይ ቀረባን ሌንስታት', 'macro-close-up-lenses', true, false, 5)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l3_macro;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (1380, 'Cinema & Video Lenses', 'ሲኒማ እና ቪዲዮ ቀረጻ ሌንሶች', 'Leensii Sinimaa fi Viidiyoo', 'ሲኒማን ቪድዮን ሌንስታት', 'cinema-video-lenses', true, false, 6)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l3_cine;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (1380, 'Lens Adapters, Filters & Converters', 'የሌንስ አስማሚዎች፣ ፊልተሮች እና መቀየሪያዎች', 'Adaapteroota fi Filteerii Leensii', 'አስማማዕቲ ሌንስን ፊልተርን', 'lens-adapters-filters-converters', true, false, 7)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l3_adapt;

    -- ========================================================================
    -- LEVEL 4: SPECIFICATIONS UNDER LENSES
    -- ========================================================================
    -- Under Prime:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_prime, 'Fast 50mm Standard Primes (f/1.2 - f/1.8)', 'ፈጣን 50 ሚሜ መደበኛ ፕራይም ሌንሶች (f/1.2 - f/1.8)', 'Leensii 50mm Iskeelii Saffisaa', 'ቅልጡፍ 50 ሚሜ ስሩዕ ፕራይም ሌንስ', 'fast-50mm-standard-primes', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_fast_50;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_prime, '85mm - 135mm Portrait Telephoto Primes', '85 ሚሜ - 135 ሚሜ የቁም ፎቶ ቴሌፎቶ ፕራይም', 'Leensii 85mm-135mm Suuraa Fuulaa', '85-135 ሚሜ ናይ ገጽ ስእሊ ፕራይም', '85mm-135mm-portrait-telephoto-primes', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_portrait_85;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_prime, '35mm Street & Documentary Primes', '35 ሚሜ የጎዳና እና ዘጋቢ ፊልም ፕራይም', 'Leensii 35mm Suuraa Daandii', '35 ሚሜ ጎደናን ዶክመንተሪን ፕራይም', '35mm-street-documentary-primes', true, false, 3)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_street_35;

    -- Under Zoom:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_zoom, '24-70mm Professional Constant f/2.8 Zooms', '24-70 ሚሜ ፕሮፌሽናል ቋሚ f/2.8 ዙም', 'Leensii 24-70mm Profeeshinaalaa f/2.8', '24-70 ሚሜ ሞያዊ f/2.8 ዙም', '24-70mm-professional-f28-zooms', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_zoom_2470;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_zoom, '24-105mm / 24-120mm All-in-One Travel Zooms', '24-105 ሚሜ የጉዞ ሁለገብ ዙም ሌንሶች', 'Leensii Imalaa 24-105mm Hunda-galeessa', '24-105 ሚሜ ናይ ጉዕዞ ኩለመዳይ ዙም', '24-105mm-all-in-one-travel-zooms', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_zoom_travel;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_zoom, '18-55mm Compact Kit & Everyday Zooms', '18-55 ሚሜ ኮምፓክት ኪት እና የዕለት ተዕለት ዙም', 'Leensii Guyyaa Guyyaa 18-55mm', '18-55 ሚሜ መዓልታዊ ኪት ዙም', '18-55mm-compact-kit-everyday-zooms', true, false, 3)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_zoom_kit;

    -- Under Telephoto:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_tele, '70-200mm Professional Fast Telephoto Zooms', '70-200 ሚሜ ፕሮፌሽናል ፈጣን ቴሌፎቶ ዙም', 'Leensii 70-200mm Telephoto Saffisaa', '70-200 ሚሜ ሞያዊ ቅልጡፍ ቴሌፎቶ', '70-200mm-professional-fast-telephoto-zooms', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_tele_70200;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_tele, '100-400mm / 200-600mm Sports & Wildlife Zooms', '100-400 ሚሜ / 200-600 ሚሜ የስፖርት እና የዱር እንስሳት ዙም', 'Leensii Ispoortii fi Bineensotaa 200-600mm', '200-600 ሚሜ ስፖርትን እንስሳታትን ዙም', '100-400mm-200-600mm-wildlife-zooms', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_tele_wildlife;

    -- Under Wide-Angle:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_wide, '14-24mm / 16-35mm Professional Wide Zooms', '14-24 ሚሜ / 16-35 ሚሜ ፕሮፌሽናል ሰፊ ዙም', 'Leensii Bal''aa Profeeshinaalaa 16-35mm', '16-35 ሚሜ ሞያዊ ሰፊሕ ዙም', '14-24mm-16-35mm-professional-wide-zooms', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_wide_1635;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_wide, '20mm - 24mm Astrophotography & Landscape Primes', '20 ሚሜ - 24 ሚሜ የከዋክብት እና የተፈጥሮ ገጽታ ፕራይም', 'Leensii Urjii fi Bakkee 20mm-24mm', '20-24 ሚሜ ከዋክብትን መልክዕ ምድርን ፕራይም', '20mm-24mm-astrophotography-landscape-primes', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_wide_primes;

    -- Under Macro:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_macro, '90mm - 105mm 1:1 True Magnification Macro', '90 ሚሜ - 105 ሚሜ 1:1 እውነተኛ ማክሮ ሌንሶች', 'Leensii 1:1 Maakroo Qulqulluu 90-105mm', '90-105 ሚሜ 1:1 ሓቀኛ ማክሮ ሌንስ', '90mm-105mm-1-1-true-macro', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_macro_100;

    -- Under Cinema:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_cine, 'Full-Frame T-Stop Cine Primes', 'ሙሉ ፍሬም T-Stop ሲኒማ ፕራይም ሌንሶች', 'Leensii Sinimaa T-Stop Guutuu', 'ሙሉእ ፍሬም T-Stop ሲኒማ ፕራይም', 'full-frame-t-stop-cine-primes', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_cine_primes;

    -- Under Adapters:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_adapt, 'Smart Electronic Autofocus Mount Adapters', 'ስማርት ኤሌክትሮኒክ አውቶፎከስ ማውንት አስማሚዎች', 'Adaapteroota Ootofookasii Ismaartii', 'ስማርት ኤሌክትሮኒክ ኣውቶፎከስ አስማማዕቲ', 'smart-electronic-autofocus-mount-adapters', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_smart_adapters;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_adapt, 'Circular Polarizing & Variable ND Filters', 'ፖላራይዚንግ እና ተለዋዋጭ ND ፊልተሮች', 'Filteerii ND fi Poolaaraayizingii', 'ፖላራይዚንግን ተቐያያሪ ND ፊልተርን', 'circular-polarizing-variable-nd-filters', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_filters;

    -- ========================================================================
    -- LEVEL 5: VARIANTS (CAMERA MOUNTS)
    -- ========================================================================
    -- Under 24-70mm Zoom:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_zoom_2470, 'Sony E-Mount Compatible Edition', 'የሶኒ ኢ-ማውንት ተስማሚ እትም', 'Maxxansa Sony E-Mount', 'ናይ ሶኒ E-Mount ሕታም', 'sony-e-mount-compatible-edition', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l5_sony_e;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_zoom_2470, 'Canon RF Mount Professional Edition', 'የካኖን RF ማውንት ፕሮፌሽናል እትም', 'Maxxansa Canon RF Mount', 'ናይ ካኖን RF Mount ሕታም', 'canon-rf-mount-professional-edition', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l5_canon_rf;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_zoom_2470, 'Nikon Z Mount Precision Edition', 'የኒኮን Z ማውንት ትክክለኛ እትም', 'Maxxansa Nikon Z Mount', 'ናይ ኒኮን Z Mount ሕታም', 'nikon-z-mount-precision-edition', true, false, 3)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l5_nikon_z;

    -- Under Fast 50mm:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_fast_50, 'Sony FE Prime Mount Edition', 'የሶኒ FE ፕራይም ማውንት እትም', 'Maxxansa Sony FE Prime', 'ናይ ሶኒ FE ፕራይም ሕታም', 'sony-fe-prime-mount-edition', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l5_sigma_e;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_fast_50, 'Canon RF Fast Prime Edition', 'የካኖን RF ፈጣን ፕራይም እትም', 'Maxxansa Canon RF Prime', 'ናይ ካኖን RF ፕራይም ሕታም', 'canon-rf-fast-prime-edition', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l5_tamron_e;

    -- ========================================================================
    -- LEVEL 6: AUTHENTIC LEAF PRODUCTS (CAMERA LENSES)
    -- ========================================================================
    -- Under Sony E-Mount (24-70mm):
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l5_sony_e, 'Sony FE 24-70mm f/2.8 GM II Lens', 'ሶኒ FE 24-70 ሚሜ f/2.8 GM II ሌንስ', 'Leensii Sony FE 24-70mm f/2.8 GM II', 'ሌንስ ሶኒ FE 24-70mm f/2.8 GM II', 'sony-fe-24-70mm-f28-gm-ii-lens', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l5_sony_e, 'Sigma 24-70mm f/2.8 DG DN Art Lens for Sony E', 'ሲግማ 24-70 ሚሜ f/2.8 DG DN Art ሌንስ ለሶኒ ኢ', 'Leensii Sigma 24-70mm f/2.8 DG DN Art', 'ሌንስ ሲግማ 24-70mm f/2.8 DG DN Art', 'sigma-24-70mm-f28-dg-dn-art-lens-sony-e', true, true, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l5_sony_e, 'Tamron 28-75mm f/2.8 Di III VXD G2 for Sony E', 'ታምሮን 28-75 ሚሜ f/2.8 Di III VXD G2 ሌንስ ለሶኒ ኢ', 'Leensii Tamron 28-75mm f/2.8 Di III VXD G2', 'ሌንስ ታምሮን 28-75mm f/2.8 Di III VXD G2', 'tamron-28-75mm-f28-di-iii-vxd-g2-sony-e', true, true, 3)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- Under Canon RF Mount (24-70mm):
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l5_canon_rf, 'Canon RF 24-70mm f/2.8L IS USM Lens', 'ካኖን RF 24-70 ሚሜ f/2.8L IS USM ሌንስ', 'Leensii Canon RF 24-70mm f/2.8L IS USM', 'ሌንስ ካኖን RF 24-70mm f/2.8L IS USM', 'canon-rf-24-70mm-f28l-is-usm-lens', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- Under Nikon Z Mount (24-70mm):
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l5_nikon_z, 'Nikon NIKKOR Z 24-70mm f/2.8 S Lens', 'ኒኮን ኒኮር Z 24-70 ሚሜ f/2.8 S ሌንስ', 'Leensii Nikon NIKKOR Z 24-70mm f/2.8 S', 'ሌንስ ኒኮን NIKKOR Z 24-70mm f/2.8 S', 'nikon-nikkor-z-24-70mm-f28-s-lens', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- Under Fast 50mm Primes:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l5_sigma_e, 'Sony FE 50mm f/1.4 GM Full Frame Prime Lens', 'ሶኒ FE 50 ሚሜ f/1.4 GM ሙሉ ፍሬም ፕራይም ሌንስ', 'Leensii Sony FE 50mm f/1.4 GM Prime', 'ሌንስ ሶኒ FE 50mm f/1.4 GM ፕራይም', 'sony-fe-50mm-f14-gm-prime-lens', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l5_sigma_e, 'Sony FE 50mm f/1.2 GM Full Frame Prime Lens', 'ሶኒ FE 50 ሚሜ f/1.2 GM ሙሉ ፍሬም ፕራይም ሌንስ', 'Leensii Sony FE 50mm f/1.2 GM Prime', 'ሌንስ ሶኒ FE 50mm f/1.2 GM ፕራይም', 'sony-fe-50mm-f12-gm-prime-lens', true, true, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l5_tamron_e, 'Canon RF 50mm f/1.8 STM Compact Prime Lens', 'ካኖን RF 50 ሚሜ f/1.8 STM ኮምፓክት ፕራይም ሌንስ', 'Leensii Canon RF 50mm f/1.8 STM Prime', 'ሌንስ ካኖን RF 50mm f/1.8 STM ፕራይም', 'canon-rf-50mm-f18-stm-compact-prime-lens', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- Under Portrait 85mm:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_portrait_85, 'Sony FE 85mm f/1.4 GM Portrait Lens', 'ሶኒ FE 85 ሚሜ f/1.4 GM የቁም ፎቶ ሌንስ', 'Leensii Sony FE 85mm f/1.4 GM', 'ሌንስ ሶኒ FE 85mm f/1.4 GM', 'sony-fe-85mm-f14-gm-portrait-lens', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_portrait_85, 'Canon RF 85mm f/1.2L USM Portrait Prime Lens', 'ካኖን RF 85 ሚሜ f/1.2L USM የቁም ፎቶ ሌንስ', 'Leensii Canon RF 85mm f/1.2L USM', 'ሌንስ ካኖን RF 85mm f/1.2L USM', 'canon-rf-85mm-f12l-usm-portrait-prime-lens', true, true, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- Under 70-200mm Telephoto:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_tele_70200, 'Sony FE 70-200mm f/2.8 GM OSS II Telephoto Lens', 'ሶኒ FE 70-200 ሚሜ f/2.8 GM OSS II ቴሌፎቶ ሌንስ', 'Leensii Sony FE 70-200mm f/2.8 GM OSS II', 'ሌንስ ሶኒ FE 70-200mm f/2.8 GM OSS II', 'sony-fe-70-200mm-f28-gm-oss-ii-telephoto-lens', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_tele_70200, 'Canon RF 70-200mm f/2.8L IS USM Telephoto Lens', 'ካኖን RF 70-200 ሚሜ f/2.8L IS USM ቴሌፎቶ ሌንስ', 'Leensii Canon RF 70-200mm f/2.8L IS USM', 'ሌንስ ካኖን RF 70-200mm f/2.8L IS USM', 'canon-rf-70-200mm-f28l-is-usm-telephoto-lens', true, true, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- Under Super-Telephoto 200-600mm:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_tele_wildlife, 'Sony FE 200-600mm f/5.6-6.3 G OSS Super-Telephoto Lens', 'ሶኒ FE 200-600 ሚሜ f/5.6-6.3 G OSS ሱፐር ቴሌፎቶ ሌንስ', 'Leensii Sony FE 200-600mm f/5.6-6.3 G OSS', 'ሌንስ ሶኒ FE 200-600mm f/5.6-6.3 G OSS', 'sony-fe-200-600mm-f56-63-g-oss-lens', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_tele_wildlife, 'Sigma 150-600mm f/5-6.3 DG DN OS Sports Lens for Sony E', 'ሲግማ 150-600 ሚሜ f/5-6.3 DG DN OS የስፖርት ሌንስ ለሶኒ', 'Leensii Sigma 150-600mm f/5-6.3 DG DN OS Sports', 'ሌንስ ሲግማ 150-600mm f/5-6.3 DG DN OS', 'sigma-150-600mm-f5-63-dg-dn-os-sports-lens', true, true, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- Under Wide Zoom 16-35mm:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_wide_1635, 'Sony FE 16-35mm f/2.8 GM II Wide-Angle Zoom Lens', 'ሶኒ FE 16-35 ሚሜ f/2.8 GM II ሰፊ ዙም ሌንስ', 'Leensii Sony FE 16-35mm f/2.8 GM II', 'ሌንስ ሶኒ FE 16-35mm f/2.8 GM II', 'sony-fe-16-35mm-f28-gm-ii-lens', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_wide_1635, 'Canon RF 15-35mm f/2.8L IS USM Ultra-Wide Zoom Lens', 'ካኖን RF 15-35 ሚሜ f/2.8L IS USM እጅግ ሰፊ ዙም ሌንስ', 'Leensii Canon RF 15-35mm f/2.8L IS USM', 'ሌንስ ካኖን RF 15-35mm f/2.8L IS USM', 'canon-rf-15-35mm-f28l-is-usm-lens', true, true, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- Under Macro 90-105mm:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_macro_100, 'Sony FE 90mm f/2.8 Macro G OSS Full Frame Lens', 'ሶኒ FE 90 ሚሜ f/2.8 Macro G OSS ማክሮ ሌንስ', 'Leensii Sony FE 90mm f/2.8 Macro G OSS', 'ሌንስ ሶኒ FE 90mm f/2.8 Macro G OSS', 'sony-fe-90mm-f28-macro-g-oss-lens', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_macro_100, 'Canon RF 100mm f/2.8L Macro IS USM Lens', 'ካኖን RF 100 ሚሜ f/2.8L Macro IS USM ማክሮ ሌንስ', 'Leensii Canon RF 100mm f/2.8L Macro IS USM', 'ሌንስ ካኖን RF 100mm f/2.8L Macro IS USM', 'canon-rf-100mm-f28l-macro-is-usm-lens', true, true, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_macro_100, 'Sigma 105mm f/2.8 DG DN Macro Art Lens for Sony E', 'ሲግማ 105 ሚሜ f/2.8 DG DN Macro Art ሌንስ ለሶኒ', 'Leensii Sigma 105mm f/2.8 DG DN Macro Art', 'ሌንስ ሲግማ 105mm f/2.8 DG DN Macro Art', 'sigma-105mm-f28-dg-dn-macro-art-lens', true, true, 3)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- Under Travel Zooms:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_zoom_travel, 'Sony FE 24-105mm f/4 G OSS Standard Zoom Lens', 'ሶኒ FE 24-105 ሚሜ f/4 G OSS የጉዞ ዙም ሌንስ', 'Leensii Sony FE 24-105mm f/4 G OSS', 'ሌንስ ሶኒ FE 24-105mm f/4 G OSS', 'sony-fe-24-105mm-f4-g-oss-lens', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_zoom_travel, 'Canon RF 24-105mm f/4L IS USM Standard Zoom Lens', 'ካኖን RF 24-105 ሚሜ f/4L IS USM የጉዞ ዙም ሌንስ', 'Leensii Canon RF 24-105mm f/4L IS USM', 'ሌንስ ካኖን RF 24-105mm f/4L IS USM', 'canon-rf-24-105mm-f4l-is-usm-lens', true, true, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- ========================================================================
    -- LEVEL 3, 4, 5, 6: MIRRORLESS CAMERAS (1382)
    -- ========================================================================
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (1382, 'Full Frame Mirrorless Cameras', 'ሙሉ ፍሬም ሚረርለስ ካሜራዎች', 'Kaameraa Miirarlesii Guutuu (Full Frame)', 'ሙሉእ ፍሬም ሚረርለስ ካሜራታት', 'full-frame-mirrorless-cameras', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l3_ml_ff;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (1382, 'APS-C / Crop Sensor Mirrorless Cameras', 'APS-C / ክሮፕ ሴንሰር ሚረርለስ ካሜራዎች', 'Kaameraa Miirarlesii APS-C', 'APS-C ሚረርለስ ካሜራታት', 'apsc-crop-sensor-mirrorless-cameras', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l3_ml_apsc;

    -- Under Full Frame Mirrorless (Level 4):
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_ml_ff, 'Hybrid Video & Stills Mirrorless Bodies', 'ድብልቅ ቪዲዮ እና ፎቶ ሚረርለስ ካሜራዎች', 'Kaameraa Miirarlesii Viidiyoo fi Suuraa', 'ድቅል ቪድዮን ፎቶን ሚረርለስ ካሜራ', 'hybrid-video-stills-mirrorless-bodies', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_ml_hybrid;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_ml_ff, 'Flagship High-Resolution Mirrorless Bodies', 'ከፍተኛ ጥራት ፍላግሺፕ ሚረርለስ ካሜራዎች', 'Kaameraa Miirarlesii Qulqullina Ol''aanaa', 'ዝለዓለ ፅሬት ፍላግሺፕ ሚረርለስ ካሜራ', 'flagship-high-resolution-mirrorless-bodies', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l4_ml_flagship;

    -- Under Hybrid (Level 5):
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_ml_hybrid, 'Sony Alpha Full Frame Series', 'የሶኒ አልፋ ሙሉ ፍሬም ተከታታይ', 'Waliigala Sony Alpha Full Frame', 'ናይ ሶኒ አልፋ ሙሉእ ፍሬም ተኸታታሊ', 'sony-alpha-full-frame-series', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l5_ml_sony;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_ml_hybrid, 'Canon EOS R Full Frame Series', 'የካኖን EOS R ሙሉ ፍሬም ተከታታይ', 'Waliigala Canon EOS R Full Frame', 'ናይ ካኖን EOS R ሙሉእ ፍሬም ተኸታታሊ', 'canon-eos-r-full-frame-series', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l5_ml_canon;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_ml_hybrid, 'Nikon Z Full Frame Series', 'የኒኮን Z ሙሉ ፍሬም ተከታታይ', 'Waliigala Nikon Z Full Frame', 'ናይ ኒኮን Z ሙሉእ ፍሬም ተኸታታሊ', 'nikon-z-full-frame-series', true, false, 3)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l5_ml_nikon;

    -- Level 6 Leaves (Mirrorless Bodies):
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l5_ml_sony, 'Sony Alpha 7 IV (A7 IV) Full Frame Mirrorless Camera Body', 'ሶኒ አልፋ 7 IV ሙሉ ፍሬም ሚረርለስ ካሜራ ቦዲ', 'Sony Alpha 7 IV (A7 IV) Kaameraa Miirarlesii', 'ሶኒ አልፋ 7 IV ሙሉእ ፍሬም ሚረርለስ ካሜራ', 'sony-alpha-7-iv-a7-iv-mirrorless-body', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_ml_flagship, 'Sony Alpha 7R V (A7R V) 61MP Full Frame Mirrorless Camera Body', 'ሶኒ አልፋ 7R V 61 ሜጋፒክሰል ሙሉ ፍሬም ካሜራ ቦዲ', 'Sony Alpha 7R V (A7R V) 61MP Kaameraa', 'ሶኒ አልፋ 7R V 61MP ሙሉእ ፍሬም ካሜራ', 'sony-alpha-7r-v-a7r-v-61mp-mirrorless-body', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l5_ml_canon, 'Canon EOS R6 Mark II Full Frame Mirrorless Camera Body', 'ካኖን EOS R6 Mark II ሙሉ ፍሬም ሚረርለስ ካሜራ ቦዲ', 'Canon EOS R6 Mark II Kaameraa Miirarlesii', 'ካኖን EOS R6 Mark II ሙሉእ ፍሬም ካሜራ', 'canon-eos-r6-mark-ii-mirrorless-body', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_ml_flagship, 'Canon EOS R5 45MP 8K Full Frame Mirrorless Camera Body', 'ካኖን EOS R5 45 ሜጋፒክሰል 8K ሙሉ ፍሬም ካሜራ ቦዲ', 'Canon EOS R5 45MP 8K Kaameraa Miirarlesii', 'ካኖን EOS R5 45MP 8K ሙሉእ ፍሬም ካሜራ', 'canon-eos-r5-45mp-8k-mirrorless-body', true, true, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l5_ml_nikon, 'Nikon Z6 III Full Frame Mirrorless Camera Body', 'ኒኮን Z6 III ሙሉ ፍሬም ሚረርለስ ካሜራ ቦዲ', 'Nikon Z6 III Kaameraa Miirarlesii', 'ኒኮን Z6 III ሙሉእ ፍሬም ካሜራ', 'nikon-z6-iii-mirrorless-body', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l4_ml_flagship, 'Nikon Z8 Flagship 45.7MP Full Frame Mirrorless Camera Body', 'ኒኮን Z8 ፍላግሺፕ 45.7 ሜጋፒክሰል ካሜራ ቦዲ', 'Nikon Z8 Kaameraa Miirarlesii 45.7MP', 'ኒኮን Z8 ፍላግሺፕ 45.7MP ካሜራ', 'nikon-z8-flagship-mirrorless-body', true, true, 3)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- Under APS-C:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_ml_apsc, 'Fujifilm X-T5 40MP APS-C Mirrorless Camera Body', 'ፉጂፊልም X-T5 40 ሜጋፒክሰል APS-C ሚረርለስ ካሜራ ቦዲ', 'Fujifilm X-T5 40MP Kaameraa Miirarlesii', 'ፉጂፊልም X-T5 40MP APS-C ሚረርለስ ካሜራ', 'fujifilm-x-t5-40mp-mirrorless-body', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    -- ========================================================================
    -- LEVEL 3, 4, 5, 6: DSLR CAMERAS (1381)
    -- ========================================================================
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (1381, 'Full Frame Professional DSLRs', 'ሙሉ ፍሬም ፕሮፌሽናል ዲኤስኤልአር ካሜራዎች', 'Kaameraa DSLR Profeeshinaalaa Guutuu', 'ሙሉእ ፍሬም ሞያዊ DSLR ካሜራታት', 'full-frame-professional-dslrs', true, false, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l3_dslr_ff;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (1381, 'APS-C Enthusiast DSLRs', 'APS-C ዲኤስኤልአር ካሜራዎች', 'Kaameraa DSLR APS-C', 'APS-C DSLR ካሜራታት', 'apsc-enthusiast-dslrs', true, false, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true RETURNING id INTO v_l3_dslr_apsc;

    -- Leaves for DSLR:
    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_dslr_ff, 'Canon EOS 5D Mark IV Full Frame DSLR Camera Body', 'ካኖን EOS 5D Mark IV ሙሉ ፍሬም ዲኤስኤልአር ካሜራ ቦዲ', 'Canon EOS 5D Mark IV Kaameraa DSLR', 'ካኖን EOS 5D Mark IV ሙሉእ ፍሬም DSLR ካሜራ', 'canon-eos-5d-mark-iv-dslr-body', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_dslr_ff, 'Nikon D850 Full Frame 45.7MP DSLR Camera Body', 'ኒኮን D850 ሙሉ ፍሬም 45.7 ሜጋፒክሰል ዲኤስኤልአር ካሜራ ቦዲ', 'Nikon D850 Kaameraa DSLR 45.7MP', 'ኒኮን D850 ሙሉእ ፍሬም 45.7MP DSLR ካሜራ', 'nikon-d850-full-frame-dslr-body', true, true, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_dslr_apsc, 'Canon EOS 90D APS-C DSLR Camera Body', 'ካኖን EOS 90D APS-C ዲኤስኤልአር ካሜራ ቦዲ', 'Canon EOS 90D Kaameraa DSLR', 'ካኖን EOS 90D APS-C DSLR ካሜራ', 'canon-eos-90d-dslr-body', true, true, 1)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

    INSERT INTO category_hierarchy (parent_id, name_en, name_am, name_om, name_ti, slug, is_active, is_leaf, display_order)
    VALUES (v_l3_dslr_apsc, 'Nikon D7500 DX-Format DSLR Camera Body', 'ኒኮን D7500 DX ፎርማት ዲኤስኤልአር ካሜራ ቦዲ', 'Nikon D7500 Kaameraa DSLR', 'ኒኮን D7500 DX DSLR ካሜራ', 'nikon-d7500-dslr-body', true, true, 2)
    ON CONFLICT (slug) DO UPDATE SET name_en = EXCLUDED.name_en, is_active = true, is_leaf = true;

END $$;

-- ----------------------------------------------------------------------------
-- 6. Multilingual Search Aliases for Camera Lenses & Photography
-- ----------------------------------------------------------------------------
INSERT INTO category_aliases (category_id, alias, language)
VALUES
-- Camera Lenses (1380)
(1380, 'camera lenses', 'en'),
(1380, 'camera lens', 'en'),
(1380, 'lenses', 'en'),
(1380, 'lens', 'en'),
(1380, 'photo lens', 'en'),
(1380, 'prime lens', 'en'),
(1380, 'zoom lens', 'en'),
(1380, 'telephoto lens', 'en'),
(1380, 'macro lens', 'en'),
(1380, 'wide angle lens', 'en'),
(1380, '24-70mm', 'en'),
(1380, '50mm', 'en'),
(1380, '70-200mm', 'en'),
(1380, 'sony lens', 'en'),
(1380, 'canon lens', 'en'),
(1380, 'nikon lens', 'en'),
(1380, 'sigma lens', 'en'),
(1380, 'tamron lens', 'en'),
(1380, 'የካሜራ ሌንስ', 'am'),
(1380, 'ሌንስ', 'am'),
(1380, 'ካሜራ ሌንሶች', 'am'),
(1380, 'ሌንሶች', 'am'),
(1380, 'ፕራይም ሌንስ', 'am'),
(1380, 'ዙም ሌንስ', 'am'),
(1380, 'ቴሌፎቶ ሌንስ', 'am'),
(1380, 'ማክሮ ሌንስ', 'am'),
(1380, 'leensii kaameraa', 'om'),
(1380, 'leensii', 'om'),
(1380, 'leensota', 'om'),
(1380, 'kaameraa leensii', 'om'),
(1380, 'ናይ ካሜራ ሌንስ', 'ti'),
(1380, 'ሌንስ', 'ti'),
(1380, 'ሌንስታት', 'ti'),
(1380, 'ካሜራ ሌንስ', 'ti'),

-- Mirrorless Cameras (1382)
(1382, 'mirrorless cameras', 'en'),
(1382, 'mirrorless camera', 'en'),
(1382, 'mirrorless', 'en'),
(1382, 'sony a7', 'en'),
(1382, 'canon r6', 'en'),
(1382, 'canon r5', 'en'),
(1382, 'nikon z6', 'en'),
(1382, 'fujifilm x-t5', 'en'),
(1382, 'ሚረርለስ ካሜራ', 'am'),
(1382, 'ካሜራ ሚረርለስ', 'am'),
(1382, 'kaameraa miirarles', 'om'),
(1382, 'ሚረርለስ ካሜራ', 'ti'),

-- DSLR Cameras (1381)
(1381, 'dslr cameras', 'en'),
(1381, 'dslr camera', 'en'),
(1381, 'dslr', 'en'),
(1381, 'canon 5d', 'en'),
(1381, 'nikon d850', 'en'),
(1381, 'ዲኤስኤልአር ካሜራ', 'am'),
(1381, 'ዲኤስኤልአር', 'am'),
(1381, 'kaameraa dslr', 'om'),
(1381, 'ዲኤስኤልአር ካሜራ', 'ti'),

-- Action Cameras (1379)
(1379, 'action cameras', 'en'),
(1379, 'action camera', 'en'),
(1379, 'gopro', 'en'),
(1379, 'ጎፕሮ', 'am'),
(1379, 'አክሽን ካሜራ', 'am'),
(1379, 'kaameraa goopiroo', 'om'),
(1379, 'ጎፕሮ ካሜራ', 'ti'),

-- Cameras & Photography (1378)
(1378, 'cameras', 'en'),
(1378, 'camera', 'en'),
(1378, 'photography', 'en'),
(1378, 'ካሜራ', 'am'),
(1378, 'ፎቶግራፍ', 'am'),
(1378, 'kaameraa', 'om'),
(1378, 'ካሜራ', 'ti')
ON CONFLICT (category_id, alias, language) DO NOTHING;

COMMIT;
