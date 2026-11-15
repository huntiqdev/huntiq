-- ============================================================================
-- 084_comprehensive_hierarchical_expansion_layers_3_to_6.sql
-- Complete Hierarchical Taxonomy Expansion Across Layers 3, 4, 5, and 6
-- 
-- Expands all essential marketplace sectors directly in PostgreSQL database:
--   - Electronics & Computing (Smartphones, Budget Phones, Laptops, 4K TVs, Inverters)
--   - Vehicles & Transport (Bajaj 3-Wheelers, Sedans, Commercial Box Trucks, Tires, Oils)
--   - Agriculture & Commodities (White Magna Teff, Specialty Coffee, High-Yield Dairy)
--   - Construction & Hardware (Rotary Hammer Drills, Hydraulic Cement PPC 32.5R)
--   - Real Estate & Property (Bole/Kazanchis Furnished Apartments, Condos, Villas)
--   - Fashion & Traditional Attire (Handwoven Tilfi Habesha Kemis Sets)
--   - Home & Kitchen (Smart Digital Injera Mitad 58cm, Inverter Refrigerators)
--   - Renewable Energy (Off-Grid Hybrid Inverters 5kW 48V)
--   - Services & Trades (Licensed Residential Electrical Wiring & DB Setup)
-- 
-- Fully localized in 4 languages: English (en), Amharic (am), Oromo (om), Tigrinya (ti)
-- Synchronized across classification_concepts, classification_placements, and category_hierarchy
-- ============================================================================

BEGIN;

DO $$
DECLARE
    v_scheme_id INT;
    v_parent_pid INT;
    v_cid INT;
    v_pid INT;
    v_ch_parent_id INT;
BEGIN
    SELECT id INTO v_scheme_id FROM classification_schemes WHERE code = 'SYS_PROD';
    IF v_scheme_id IS NULL THEN
        RAISE EXCEPTION 'SYS_PROD scheme not found';
    END IF;

    -- --- Node: smartphones-touchscreen (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_SMARTPHONES', 'Smartphones & Touchscreen Devices', 'Smartphones & Touchscreen Devices', 'ስማርት ስልኮችና የንኪኪ መሳሪያዎች', 'Ismaartifoonota fi Meeshaalee Tuqaa', 'ስማርት ፎናትን ናይ ንክኪ መሳርሕታትን', 'Android and iOS smartphones with touchscreens and cellular connectivity', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'mobile-phones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'smartphones-touchscreen', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'mobile-phones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Smartphones & Touchscreen Devices', 'ስማርት ስልኮችና የንኪኪ መሳሪያዎች', 'Ismaartifoonota fi Meeshaalee Tuqaa', 'ስማርት ፎናትን ናይ ንክኪ መሳርሕታትን', 'smartphones-touchscreen', 'Android and iOS smartphones with touchscreens and cellular connectivity', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: feature-button-phones (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_FEATURE_BUTTON', 'Feature & Button Keypad Phones', 'Feature & Button Keypad Phones', 'የቁልፍ እና የባህላዊ ኪፓድ ስልኮች', 'Bilbiloota Qabduu fi Kiibaardii Dhaabbataa', 'ናይ ቁልፊን መደበኛ ኪፓድን ዘለዎም ፎናት', 'Durable physical keypad mobile handsets with ultra-long standby batteries', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'mobile-phones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'feature-button-phones', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'mobile-phones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Feature & Button Keypad Phones', 'የቁልፍ እና የባህላዊ ኪፓድ ስልኮች', 'Bilbiloota Qabduu fi Kiibaardii Dhaabbataa', 'ናይ ቁልፊን መደበኛ ኪፓድን ዘለዎም ፎናት', 'feature-button-phones', 'Durable physical keypad mobile handsets with ultra-long standby batteries', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: portable-laptops-notebooks (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_COMP_PORTABLE_LAPTOPS', 'Laptops & Portable Notebooks', 'Laptops & Portable Notebooks', 'ተንቀሳቃሽ ላፕቶፖችና ኖትቡኮች', 'Laappitooppota fi Meeshaalee Kompiitaraa', 'ተንቀሳቐስቲ ላፕቶፓትን ኖትቡካትን', 'Portable clamshell personal computers, ultrabooks, and workstations', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'computers-laptops' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'portable-laptops-notebooks', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'computers-laptops' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Laptops & Portable Notebooks', 'ተንቀሳቃሽ ላፕቶፖችና ኖትቡኮች', 'Laappitooppota fi Meeshaalee Kompiitaraa', 'ተንቀሳቐስቲ ላፕቶፓትን ኖትቡካትን', 'portable-laptops-notebooks', 'Portable clamshell personal computers, ultrabooks, and workstations', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: smart-televisions-displays (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TV_SMART_TELEVISIONS', 'Smart Televisions & Flat Displays', 'Smart Televisions & Flat Displays', 'ስማርት ቴሌቪዥኖችና ጠፍጣፋ ስክሪኖች', 'Televiziyinii Ismaartii fi Iskiriinii Bal''aa', 'ስማርት ቴሌቪዥናትን ጸፍጻፍ ስክሪናትን', 'Smart LED, QLED, and OLED flat screen televisions with WiFi and streaming apps', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'tv-audio' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'smart-televisions-displays', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'tv-audio' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Smart Televisions & Flat Displays', 'ስማርት ቴሌቪዥኖችና ጠፍጣፋ ስክሪኖች', 'Televiziyinii Ismaartii fi Iskiriinii Bal''aa', 'ስማርት ቴሌቪዥናትን ጸፍጻፍ ስክሪናትን', 'smart-televisions-displays', 'Smart LED, QLED, and OLED flat screen televisions with WiFi and streaming apps', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: passenger-sedans-compacts (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_PASSENGER_SEDANS', 'Passenger Sedans & Compact Automobiles', 'Passenger Sedans & Compact Automobiles', 'የቤት አውቶሞቢሎችና ኮምፓክት መኪኖች', 'Konkolaattota Maatii fi Sedaanota', 'ናይ ገዛ መኻይንን ኮምፓክት ሴዳናትን', 'Four-door passenger sedans and compact hatchback automobiles', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'cars' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'passenger-sedans-compacts', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'cars' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Passenger Sedans & Compact Automobiles', 'የቤት አውቶሞቢሎችና ኮምፓክት መኪኖች', 'Konkolaattota Maatii fi Sedaanota', 'ናይ ገዛ መኻይንን ኮምፓክት ሴዳናትን', 'passenger-sedans-compacts', 'Four-door passenger sedans and compact hatchback automobiles', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: three-wheelers-bajaj (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_THREE_WHEELERS_BAJAJ', 'Three-Wheelers & Bajaj Auto Rickshaws', 'Three-Wheelers & Bajaj Auto Rickshaws', 'ባለ ሶስት እግር ባጃጅና አውቶ ሪክሾዎች', 'Baajaajii fi Rikkishoo Miila Sadii', 'ባዓል ሰለስተ እግሪ ባጃጅን ሪክሾን', 'Three-wheeled motorized passenger and cargo auto-rickshaws (Bajaj, TVS, Piaggio)', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'motorcycles' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'three-wheelers-bajaj', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'motorcycles' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Three-Wheelers & Bajaj Auto Rickshaws', 'ባለ ሶስት እግር ባጃጅና አውቶ ሪክሾዎች', 'Baajaajii fi Rikkishoo Miila Sadii', 'ባዓል ሰለስተ እግሪ ባጃጅን ሪክሾን', 'three-wheelers-bajaj', 'Three-wheeled motorized passenger and cargo auto-rickshaws (Bajaj, TVS, Piaggio)', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: medium-commercial-cargo-trucks (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_MEDIUM_CARGO_TRUCKS', 'Medium Commercial Cargo Trucks', 'Medium Commercial Cargo Trucks', 'መካከለኛ የንግድ ጭነት መኪኖች', 'Konkolaattota Fe''umsaa Giddu-galeessaa', 'ማእከላይ ናይ ንግዲ ፅዕነት መኻይን', 'Medium-duty diesel trucks (3.5 ton to 10 ton) for inter-city freight and commercial logistics', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'trucks' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'medium-commercial-cargo-trucks', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'trucks' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Medium Commercial Cargo Trucks', 'መካከለኛ የንግድ ጭነት መኪኖች', 'Konkolaattota Fe''umsaa Giddu-galeessaa', 'ማእከላይ ናይ ንግዲ ፅዕነት መኻይን', 'medium-commercial-cargo-trucks', 'Medium-duty diesel trucks (3.5 ton to 10 ton) for inter-city freight and commercial logistics', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: furnished-apartments-for-rent (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PROP_FURNISHED_APTS_RENT', 'Furnished Apartments for Rent', 'Furnished Apartments for Rent', 'የተሟላላቸው የሚከራዩ አፓርታማዎች', 'Appaartamaalee Meeshaan Guutaman Kireeffaman', 'ዕቑር ንብረት ዘለዎም ዝካረዩ ኣፓርታማታት', 'Fully furnished and serviced residential apartments available for medium to long-term rent', 'PROPERTY', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'apartments' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'furnished-apartments-for-rent', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'apartments' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Furnished Apartments for Rent', 'የተሟላላቸው የሚከራዩ አፓርታማዎች', 'Appaartamaalee Meeshaan Guutaman Kireeffaman', 'ዕቑር ንብረት ዘለዎም ዝካረዩ ኣፓርታማታት', 'furnished-apartments-for-rent', 'Fully furnished and serviced residential apartments available for medium to long-term rent', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: staple-grains-pulses (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_STAPLE_GRAINS', 'Staple Grains & Harvest Pulses', 'Staple Grains & Harvest Pulses', 'ዋና ዋና እህሎችና የሰብል ጥራጥሬዎች', 'Midhaan Nyaata Bu''uuraa fi Oomishaalee', 'ቀንዲ እኽልታትን ናይ ስርናይ ጥራጥሬታትን', 'Staple cereal grains, teff varieties, wheat, barley, pulses, and oilseeds', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'crops-coffee' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'staple-grains-pulses', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'crops-coffee' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Staple Grains & Harvest Pulses', 'ዋና ዋና እህሎችና የሰብል ጥራጥሬዎች', 'Midhaan Nyaata Bu''uuraa fi Oomishaalee', 'ቀንዲ እኽልታትን ናይ ስርናይ ጥራጥሬታትን', 'staple-grains-pulses', 'Staple cereal grains, teff varieties, wheat, barley, pulses, and oilseeds', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: power-tools-machinery (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TOOL_POWER_TOOLS', 'Power Tools & Electric Machinery', 'Power Tools & Electric Machinery', 'የኤሌክትሪክ ሃይል መሳሪያዎችና ማሽነሪዎች', 'Meeshaalee Humna Elektiriikii fi Maashinoota', 'ናይ ኤሌክትሪክ ሓይሊ መሳርሕታትን ማሽነሪታትን', 'Heavy-duty electric rotary hammers, angle grinders, circular saws, and workshop power equipment', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'tools-hardware' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'power-tools-machinery', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'tools-hardware' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Power Tools & Electric Machinery', 'የኤሌክትሪክ ሃይል መሳሪያዎችና ማሽነሪዎች', 'Meeshaalee Humna Elektiriikii fi Maashinoota', 'ናይ ኤሌክትሪክ ሓይሊ መሳርሕታትን ማሽነሪታትን', 'power-tools-machinery', 'Heavy-duty electric rotary hammers, angle grinders, circular saws, and workshop power equipment', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: hydraulic-structural-cement (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_MAT_STRUCTURAL_CEMENT', 'Hydraulic & Structural Building Cement', 'Hydraulic & Structural Building Cement', 'የህንፃ ግንባታና መዋቅራዊ ሲሚንቶ', 'Siminto Qulqullinaa fi Ijaarsaa', 'ናይ ህንጻ ህንጸትን መዋቕራዊን ሲሚንቶ', 'Portland Pozzolana Cement (PPC) and Ordinary Portland Cement (OPC) for structural casting and masonry', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'commercial-equipment' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'hydraulic-structural-cement', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'commercial-equipment' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Hydraulic & Structural Building Cement', 'የህንፃ ግንባታና መዋቅራዊ ሲሚንቶ', 'Siminto Qulqullinaa fi Ijaarsaa', 'ናይ ህንጻ ህንጸትን መዋቕራዊን ሲሚንቶ', 'hydraulic-structural-cement', 'Portland Pozzolana Cement (PPC) and Ordinary Portland Cement (OPC) for structural casting and masonry', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: women-traditional-habesha-kemis (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_FASH_HABESHA_KEMIS', 'Women''s Traditional Habesha Kemis & Dresses', 'Women''s Traditional Habesha Kemis & Dresses', 'የሴቶች ባህላዊ የሀበሻ ቀሚሶችና አልባሳት', 'Uffata Aadaa Dubartootaa Qamisii Habashaa', 'ናይ ደቂ ኣንስትዮ ባህላዊ ናይ ሓበሻ ቀሚሳትን ክዳውንትን', 'Handcrafted women''s traditional Ethiopian dresses made from handspun cotton (Shema) with intricate Tilfi embroidery', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'traditional-habesha' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'women-traditional-habesha-kemis', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'traditional-habesha' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Women''s Traditional Habesha Kemis & Dresses', 'የሴቶች ባህላዊ የሀበሻ ቀሚሶችና አልባሳት', 'Uffata Aadaa Dubartootaa Qamisii Habashaa', 'ናይ ደቂ ኣንስትዮ ባህላዊ ናይ ሓበሻ ቀሚሳትን ክዳውንትን', 'women-traditional-habesha-kemis', 'Handcrafted women''s traditional Ethiopian dresses made from handspun cotton (Shema) with intricate Tilfi embroidery', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: electric-injera-mitad-baking (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_HOME_INJERA_MITAD', 'Electric Injera Mitad & Traditional Baking Equipment', 'Electric Injera Mitad & Traditional Baking Equipment', 'የእንጀራ ኤሌክትሪክ ምጣድና መጋገሪያ መሳሪያዎች', 'Eelee Buddeenaa Elektiriikii fi Meeshaalee Tolchuu', 'ናይ እንጀራ ኤሌክትሪክ ምጣድን መጋገያ መሳርሕታትን', 'Electric clay and Teflon-coated griddles (Mitad) engineered for baking authentic Ethiopian Injera and flatbreads', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'kitchen-appliances' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'electric-injera-mitad-baking', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'kitchen-appliances' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Electric Injera Mitad & Traditional Baking Equipment', 'የእንጀራ ኤሌክትሪክ ምጣድና መጋገሪያ መሳሪያዎች', 'Eelee Buddeenaa Elektiriikii fi Meeshaalee Tolchuu', 'ናይ እንጀራ ኤሌክትሪክ ምጣድን መጋገያ መሳርሕታትን', 'electric-injera-mitad-baking', 'Electric clay and Teflon-coated griddles (Mitad) engineered for baking authentic Ethiopian Injera and flatbreads', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: solar-renewable-energy-systems (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SOLAR_RENEWABLE_SYSTEMS', 'Solar & Renewable Energy Systems', 'Solar & Renewable Energy Systems', 'የፀሐይ ብርሃንና ታዳሽ ሃይል ስርዓቶች', 'Sirna Humna Aduu fi Madaalawaa', 'ናይ ፀሓይ ብርሃንን ተሓዳሲ ሓይልን ስርዓታት', 'Solar PV panels, hybrid inverters, solar charge controllers, and deep-cycle lithium storage batteries', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'commercial-equipment' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'solar-renewable-energy-systems', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'commercial-equipment' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Solar & Renewable Energy Systems', 'የፀሐይ ብርሃንና ታዳሽ ሃይል ስርዓቶች', 'Sirna Humna Aduu fi Madaalawaa', 'ናይ ፀሓይ ብርሃንን ተሓዳሲ ሓይልን ስርዓታት', 'solar-renewable-energy-systems', 'Solar PV panels, hybrid inverters, solar charge controllers, and deep-cycle lithium storage batteries', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: electrical-installation-repair-trades (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SERV_ELECTRICAL_TRADES', 'Electrical Installation & Repair Trades', 'Electrical Installation & Repair Trades', 'የኤሌክትሪክ መስመር ዝርጋታና ጥገና ሙያ አገልግሎቶች', 'Tajaajila Dhaabbii fi Suphaa Sarara Elektiriikii', 'ናይ ኤሌክትሪክ መስመር ዝርግሓን ጽገናን ሞያዊ ኣገልግሎታት', 'Professional residential and commercial electrical contracting, circuit breaker panel wiring, and safety audits', 'SERVICE', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'repair-maintenance' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'electrical-installation-repair-trades', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'repair-maintenance' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Electrical Installation & Repair Trades', 'የኤሌክትሪክ መስመር ዝርጋታና ጥገና ሙያ አገልግሎቶች', 'Tajaajila Dhaabbii fi Suphaa Sarara Elektiriikii', 'ናይ ኤሌክትሪክ መስመር ዝርግሓን ጽገናን ሞያዊ ኣገልግሎታት', 'electrical-installation-repair-trades', 'Professional residential and commercial electrical contracting, circuit breaker panel wiring, and safety audits', 2, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: flagship-5g-smartphones (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_FLAGSHIP_5G', 'Flagship 5G Smartphones', 'Flagship 5G Smartphones', 'ዋና ከፍተኛ ደረጃ 5ጂ ስማርት ስልኮች', 'Ismaartifoonota Sadarkaa Olaanaa 5G', 'ላዕለዋይ ደረጃ 5ጂ ስማርት ፎናት', 'Premium flagship smartphones with 5G connectivity and top processors', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'smartphones-touchscreen' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'flagship-5g-smartphones', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'smartphones-touchscreen' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Flagship 5G Smartphones', 'ዋና ከፍተኛ ደረጃ 5ጂ ስማርት ስልኮች', 'Ismaartifoonota Sadarkaa Olaanaa 5G', 'ላዕለዋይ ደረጃ 5ጂ ስማርት ፎናት', 'flagship-5g-smartphones', 'Premium flagship smartphones with 5G connectivity and top processors', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: budget-entry-smartphones (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_BUDGET_ENTRY', 'Budget & Entry-Level Smartphones', 'Budget & Entry-Level Smartphones', 'ተመጣጣኝ ዋጋ ያላቸው ስማርት ስልኮች', 'Ismaartifoonota Gatii Madaalawaa', 'ማዕከላይን ተመጣጣንን ዋጋ ዘለዎም ስማርት ፎናት', 'Affordable 4G/LTE smartphones optimized for everyday connectivity and social media', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'smartphones-touchscreen' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'budget-entry-smartphones', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'smartphones-touchscreen' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Budget & Entry-Level Smartphones', 'ተመጣጣኝ ዋጋ ያላቸው ስማርት ስልኮች', 'Ismaartifoonota Gatii Madaalawaa', 'ማዕከላይን ተመጣጣንን ዋጋ ዘለዎም ስማርት ፎናት', 'budget-entry-smartphones', 'Affordable 4G/LTE smartphones optimized for everyday connectivity and social media', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: heavy-duty-keypad-phones (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_HEAVY_KEYPAD', 'Heavy-Duty Battery Keypad Phones', 'Heavy-Duty Battery Keypad Phones', 'ረጅም ሰዓት የሚቆይ ባትሪ ያላቸው የቁልፍ ስልኮች', 'Bilbiloota Qabduu Baatrii Cimaa Qaban', 'ንነዊሕ ሰዓት ዝጸንሕ ባትሪ ዘለዎም ናይ ቁልፊ ፎናት', 'Rugged feature phones engineered with multi-day standby batteries and loud speakers', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'feature-button-phones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'heavy-duty-keypad-phones', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'feature-button-phones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Heavy-Duty Battery Keypad Phones', 'ረጅም ሰዓት የሚቆይ ባትሪ ያላቸው የቁልፍ ስልኮች', 'Bilbiloota Qabduu Baatrii Cimaa Qaban', 'ንነዊሕ ሰዓት ዝጸንሕ ባትሪ ዘለዎም ናይ ቁልፊ ፎናት', 'heavy-duty-keypad-phones', 'Rugged feature phones engineered with multi-day standby batteries and loud speakers', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: business-work-ultrabooks (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_COMP_BUSINESS_ULTRABOOKS', 'Business & Professional Ultrabooks', 'Business & Professional Ultrabooks', 'የስራና የሙያተኞች አልትራቡክ ላፕቶፖች', 'Ultrabooks Hojii fi Ogeeyyii', 'ናይ ስራሕን ሙያተኛታትን ኣልትራቡክ ላፕቶፓት', 'Lightweight, enterprise-security ultrabooks designed for corporate productivity', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'portable-laptops-notebooks' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'business-work-ultrabooks', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'portable-laptops-notebooks' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Business & Professional Ultrabooks', 'የስራና የሙያተኞች አልትራቡክ ላፕቶፖች', 'Ultrabooks Hojii fi Ogeeyyii', 'ናይ ስራሕን ሙያተኛታትን ኣልትራቡክ ላፕቶፓት', 'business-work-ultrabooks', 'Lightweight, enterprise-security ultrabooks designed for corporate productivity', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: 4k-ultra-hd-smart-tvs (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TV_4K_ULTRA_HD', '4K Ultra HD Smart TVs', '4K Ultra HD Smart TVs', '4ኬ አልትራ ኤችዲ ስማርት ቴሌቪዥኖች', 'Televiziyinii Ismaartii 4K Ultra HD', '4ኬ ኣልትራ ኤችዲ ስማርት ቴሌቪዥናት', '3840x2160 resolution 4K UHD smart televisions with HDR10 support', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'smart-televisions-displays' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', '4k-ultra-hd-smart-tvs', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'smart-televisions-displays' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('4K Ultra HD Smart TVs', '4ኬ አልትራ ኤችዲ ስማርት ቴሌቪዥኖች', 'Televiziyinii Ismaartii 4K Ultra HD', '4ኬ ኣልትራ ኤችዲ ስማርት ቴሌቪዥናት', '4k-ultra-hd-smart-tvs', '3840x2160 resolution 4K UHD smart televisions with HDR10 support', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: economy-city-sedans (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_ECONOMY_SEDANS', 'Economy City Sedans & Hatchbacks', 'Economy City Sedans & Hatchbacks', 'ተመጣጣኝ የከተማ አውቶሞቢሎችና ሀችባኮች', 'Sedaanota Magaalaa fi Haachibaakii Qusataa', 'ተመጣጣኒ ዋጋ ናይ ከተማ መኻይንን ሃችባክን', 'Fuel-efficient compact sedans widely used for private family commute and ride-hailing', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'passenger-sedans-compacts' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'economy-city-sedans', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'passenger-sedans-compacts' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Economy City Sedans & Hatchbacks', 'ተመጣጣኝ የከተማ አውቶሞቢሎችና ሀችባኮች', 'Sedaanota Magaalaa fi Haachibaakii Qusataa', 'ተመጣጣኒ ዋጋ ናይ ከተማ መኻይንን ሃችባክን', 'economy-city-sedans', 'Fuel-efficient compact sedans widely used for private family commute and ride-hailing', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: commercial-passenger-bajaj (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_PASSENGER_BAJAJ', 'Commercial Passenger Bajaj Rickshaws', 'Commercial Passenger Bajaj Rickshaws', 'የሰው ማመላለሻ ንግድ ባጃጅ ሪክሾዎች', 'Baajaajii Namoota Geessitu Daldalaa', 'ናይ ሰብ መጓዓዝያ ባጃጅ ሪክሾታት', '3-passenger motorized auto-rickshaws standard for urban and regional public transport', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'three-wheelers-bajaj' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'commercial-passenger-bajaj', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'three-wheelers-bajaj' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Commercial Passenger Bajaj Rickshaws', 'የሰው ማመላለሻ ንግድ ባጃጅ ሪክሾዎች', 'Baajaajii Namoota Geessitu Daldalaa', 'ናይ ሰብ መጓዓዝያ ባጃጅ ሪክሾታት', 'commercial-passenger-bajaj', '3-passenger motorized auto-rickshaws standard for urban and regional public transport', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: isuzu-box-and-cargo-trucks (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_ISUZU_CARGO_TRUCKS', 'Isuzu NPR / FSR Diesel Cargo Trucks', 'Isuzu NPR / FSR Diesel Cargo Trucks', 'የአይሱዙ ኤንፒአር እና ኤፍኤስአር ዲዝል የጭነት መኪኖች', 'Konkolaattota Isuzu NPR fi FSR Fe''umsaa', 'ናይ ኣይሱዙ ኤንፒኣርን ኤፍኤስኣርን ናይ ፅዕነት መኻይን', 'Standard Isuzu diesel cargo trucks dominating regional market distribution', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'medium-commercial-cargo-trucks' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'isuzu-box-and-cargo-trucks', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'medium-commercial-cargo-trucks' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Isuzu NPR / FSR Diesel Cargo Trucks', 'የአይሱዙ ኤንፒአር እና ኤፍኤስአር ዲዝል የጭነት መኪኖች', 'Konkolaattota Isuzu NPR fi FSR Fe''umsaa', 'ናይ ኣይሱዙ ኤንፒኣርን ኤፍኤስኣርን ናይ ፅዕነት መኻይን', 'isuzu-box-and-cargo-trucks', 'Standard Isuzu diesel cargo trucks dominating regional market distribution', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: two-bedroom-furnished-apartments (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PROP_2BED_FURNISHED_APTS', '2-Bedroom Furnished Apartments', '2-Bedroom Furnished Apartments', 'ባለ ሁለት መኝታ የተሟላላቸው አፓርታማዎች', 'Appaartamaalee Kutaa Ciisichaa 2 Meeshaan Guutaman', 'ባዓል 2 መደቀሲ ዕቑር ንብረት ዘለዎም ኣፓርታማታት', 'Modern 2-bedroom, 2-bathroom furnished apartments with living room and fitted kitchen', 'PROPERTY', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'furnished-apartments-for-rent' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'two-bedroom-furnished-apartments', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'furnished-apartments-for-rent' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('2-Bedroom Furnished Apartments', 'ባለ ሁለት መኝታ የተሟላላቸው አፓርታማዎች', 'Appaartamaalee Kutaa Ciisichaa 2 Meeshaan Guutaman', 'ባዓል 2 መደቀሲ ዕቑር ንብረት ዘለዎም ኣፓርታማታት', 'two-bedroom-furnished-apartments', 'Modern 2-bedroom, 2-bathroom furnished apartments with living room and fitted kitchen', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: teff-grain-farming (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_TEFF_GRAINS', 'Teff Grain Harvest Varieties', 'Teff Grain Harvest Varieties', 'የጤፍ እህል ዝርያዎችና ንግድ', 'Sanyii Midhaan Xaafii fi Oomisha', 'ናይ ጣፍ እኽሊ ዓይነታትን ንግድን', 'Cultivated Eragrostis tef cereal varieties (Magna, White, Mixed Sergegna, Tikur/Brown)', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'staple-grains-pulses' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'teff-grain-farming', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'staple-grains-pulses' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Teff Grain Harvest Varieties', 'የጤፍ እህል ዝርያዎችና ንግድ', 'Sanyii Midhaan Xaafii fi Oomisha', 'ናይ ጣፍ እኽሊ ዓይነታትን ንግድን', 'teff-grain-farming', 'Cultivated Eragrostis tef cereal varieties (Magna, White, Mixed Sergegna, Tikur/Brown)', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: ethiopian-specialty-coffee (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_SPECIALTY_COFFEE', 'Ethiopian Specialty Arabica Coffee Beans', 'Ethiopian Specialty Arabica Coffee Beans', 'የኢትዮጵያ ስፔሻሊቲ አረቢካ ቡና', 'Buna Ispeeshaalitii Arabikaa Itoophiyaa', 'ናይ ኢትዮጵያ ስፔሻሊቲ ኣረቢካ ቡን', 'High altitude Ethiopian single-origin specialty Arabica green and roasted coffee beans', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'staple-grains-pulses' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'ethiopian-specialty-coffee', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'staple-grains-pulses' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Ethiopian Specialty Arabica Coffee Beans', 'የኢትዮጵያ ስፔሻሊቲ አረቢካ ቡና', 'Buna Ispeeshaalitii Arabikaa Itoophiyaa', 'ናይ ኢትዮጵያ ስፔሻሊቲ ኣረቢካ ቡን', 'ethiopian-specialty-coffee', 'High altitude Ethiopian single-origin specialty Arabica green and roasted coffee beans', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: rotary-hammer-demolition-drills (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TOOL_ROTARY_HAMMERS', 'Rotary Hammer & Demolition Drills', 'Rotary Hammer & Demolition Drills', 'ሮታሪ ሀመር እና ሰባሪ ድሪሎች', 'Boraasuu fi Waraantoo Dhagaa Cimaa', 'ሮታሪ ሃመርን ሰባሪ ድሪላትን', 'Electropneumatic rotary hammer drills for drilling reinforced concrete, granite, and masonry', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'power-tools-machinery' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'rotary-hammer-demolition-drills', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'power-tools-machinery' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Rotary Hammer & Demolition Drills', 'ሮታሪ ሀመር እና ሰባሪ ድሪሎች', 'Boraasuu fi Waraantoo Dhagaa Cimaa', 'ሮታሪ ሃመርን ሰባሪ ድሪላትን', 'rotary-hammer-demolition-drills', 'Electropneumatic rotary hammer drills for drilling reinforced concrete, granite, and masonry', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: portland-pozzolana-ordinary-cement (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_MAT_PPC_OPC_CEMENT', 'Portland Pozzolana & Ordinary Portland Cement', 'Portland Pozzolana & Ordinary Portland Cement', 'ፖርትላንድ ፖዞላና እና ኦርዲናሪ ፖርትላንድ ሲሚንቶ', 'Simintoo Poortilaandii fi Pozoolaanaa', 'ፖርትላንድ ፖዞላናን ኦርዲናሪ ፖርትላንድን ሲሚንቶ', 'General purpose and high-strength construction cements complying with Ethiopian Standard ES EN 197-1', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'hydraulic-structural-cement' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'portland-pozzolana-ordinary-cement', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'hydraulic-structural-cement' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Portland Pozzolana & Ordinary Portland Cement', 'ፖርትላንድ ፖዞላና እና ኦርዲናሪ ፖርትላንድ ሲሚንቶ', 'Simintoo Poortilaandii fi Pozoolaanaa', 'ፖርትላንድ ፖዞላናን ኦርዲናሪ ፖርትላንድን ሲሚንቶ', 'portland-pozzolana-ordinary-cement', 'General purpose and high-strength construction cements complying with Ethiopian Standard ES EN 197-1', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: handwoven-tilfi-habesha-kemis (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_FASH_HANDWOVEN_TILFI', 'Handwoven Tilfi & Tibeb Embroidered Habesha Kemis', 'Handwoven Tilfi & Tibeb Embroidered Habesha Kemis', 'በእጅ የተሸመነ የጥልፍና የጥበብ የሀበሻ ቀሚስ', 'Qamisii Habashaa Harka-Dhaawame Tilibii Qabu', 'ብኢድ ዝተሸለመ ናይ ጥልፍን ጥበብን ናይ ሓበሻ ቀሚስ', 'Exquisite handspun cotton gowns decorated with colorful geometric Tibeb borders and necklines', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'women-traditional-habesha-kemis' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'handwoven-tilfi-habesha-kemis', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'women-traditional-habesha-kemis' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Handwoven Tilfi & Tibeb Embroidered Habesha Kemis', 'በእጅ የተሸመነ የጥልፍና የጥበብ የሀበሻ ቀሚስ', 'Qamisii Habashaa Harka-Dhaawame Tilibii Qabu', 'ብኢድ ዝተሸለመ ናይ ጥልፍን ጥበብን ናይ ሓበሻ ቀሚስ', 'handwoven-tilfi-habesha-kemis', 'Exquisite handspun cotton gowns decorated with colorful geometric Tibeb borders and necklines', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: energy-saving-electric-injera-mitad (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_HOME_ENERGY_MITAD', 'Energy Saving & Thermostat Electric Injera Mitad', 'Energy Saving & Thermostat Electric Injera Mitad', 'ሀይል ቆጣቢና ቴርሞስታት ያለው የኤሌክትሪክ እንጀራ ምጣድ', 'Eelee Elektiriikii Qusataa fi To''annaa Hoo''aa Qabu', 'ሓይሊ ቆጣብን ቴርሞስታትን ዘለዎ ናይ ኤሌክትሪክ ምጣድ', 'Advanced thermal insulated electric baking plates designed for fast heat-up and reduced electricity consumption', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'electric-injera-mitad-baking' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'energy-saving-electric-injera-mitad', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'electric-injera-mitad-baking' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Energy Saving & Thermostat Electric Injera Mitad', 'ሀይል ቆጣቢና ቴርሞስታት ያለው የኤሌክትሪክ እንጀራ ምጣድ', 'Eelee Elektiriikii Qusataa fi To''annaa Hoo''aa Qabu', 'ሓይሊ ቆጣብን ቴርሞስታትን ዘለዎ ናይ ኤሌክትሪክ ምጣድ', 'energy-saving-electric-injera-mitad', 'Advanced thermal insulated electric baking plates designed for fast heat-up and reduced electricity consumption', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: off-grid-hybrid-solar-inverters (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SOLAR_HYBRID_INVERTERS', 'Off-Grid & Hybrid Solar Inverters', 'Off-Grid & Hybrid Solar Inverters', 'የፀሐይ ሃይል ኦፍ-ግሪድ እና ሃይብሪድ ኢንቨርተሮች', 'Inveertara Humna Aduu Haayibriidii', 'ናይ ፀሓይ ሓይሊ ኦፍ-ግሪድን ሃይብሪድን ኢንቨርተራት', 'Smart solar inverters integrating MPPT solar charge controllers, AC battery chargers, and grid bypass', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'solar-renewable-energy-systems' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'off-grid-hybrid-solar-inverters', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'solar-renewable-energy-systems' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Off-Grid & Hybrid Solar Inverters', 'የፀሐይ ሃይል ኦፍ-ግሪድ እና ሃይብሪድ ኢንቨርተሮች', 'Inveertara Humna Aduu Haayibriidii', 'ናይ ፀሓይ ሓይሊ ኦፍ-ግሪድን ሃይብሪድን ኢንቨርተራት', 'off-grid-hybrid-solar-inverters', 'Smart solar inverters integrating MPPT solar charge controllers, AC battery chargers, and grid bypass', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: residential-commercial-electrical-wiring (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SERV_RESIDENTIAL_WIRING', 'Residential & Commercial Building Electrical Wiring', 'Residential & Commercial Building Electrical Wiring', 'የመኖሪያና የንግድ ህንፃዎች የኤሌክትሪክ መስመር ዝርጋታ', 'Sarara Elektiriikii Manneen Jireenyaa fi Daldalaa', 'ናይ መነበሪን ንግዲን ህንጻታት ናይ ኤሌክትሪክ መስመር ዝርግሓ', 'Complete electrical rough-in and trim-out wiring, conduit pipe laying, and lighting fixture installation', 'SERVICE', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'electrical-installation-repair-trades' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'residential-commercial-electrical-wiring', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'electrical-installation-repair-trades' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Residential & Commercial Building Electrical Wiring', 'የመኖሪያና የንግድ ህንፃዎች የኤሌክትሪክ መስመር ዝርጋታ', 'Sarara Elektiriikii Manneen Jireenyaa fi Daldalaa', 'ናይ መነበሪን ንግዲን ህንጻታት ናይ ኤሌክትሪክ መስመር ዝርግሓ', 'residential-commercial-electrical-wiring', 'Complete electrical rough-in and trim-out wiring, conduit pipe laying, and lighting fixture installation', 3, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: dual-sim-flagship-amoled (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_DUAL_SIM_AMOLED', 'Dual SIM Dynamic AMOLED 120Hz Smartphones', 'Dual SIM Dynamic AMOLED 120Hz Smartphones', 'ባለ ሁለት ሲም ዳይናሚክ አሞሌድ 120ኸርዝ ስልኮች', 'Siimii Lama Dynamic AMOLED 120Hz Ismaartifoonota', 'ክልተ ሲም ዳይናሚክ ኣሞሌድ 120ኸርዝ ፎናት', 'Dual SIM flagship smartphones with 120Hz high refresh AMOLED screens', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'flagship-5g-smartphones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'dual-sim-flagship-amoled', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'flagship-5g-smartphones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Dual SIM Dynamic AMOLED 120Hz Smartphones', 'ባለ ሁለት ሲም ዳይናሚክ አሞሌድ 120ኸርዝ ስልኮች', 'Siimii Lama Dynamic AMOLED 120Hz Ismaartifoonota', 'ክልተ ሲም ዳይናሚክ ኣሞሌድ 120ኸርዝ ፎናት', 'dual-sim-flagship-amoled', 'Dual SIM flagship smartphones with 120Hz high refresh AMOLED screens', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: 4g-dual-sim-budget-smartphones (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_4G_DUAL_SIM_BUDGET', '4G Dual SIM 5000mAh Battery Budget Smartphones', '4G Dual SIM 5000mAh Battery Budget Smartphones', '4ጂ ባለሁለት ሲም 5000ሚሊ አምፔር ባትሪ ያላቸው ስልኮች', '4G Siimii Lama Baatrii 5000mAh Ismaartifoonota', '4ጂ ክልተ ሲም 5000ሚሊ ኣምፔር ባትሪ ዘለዎም ፎናት', 'Budget smartphones featuring 4G LTE, dual SIM standby, and high capacity 5000mAh battery', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'budget-entry-smartphones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', '4g-dual-sim-budget-smartphones', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'budget-entry-smartphones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('4G Dual SIM 5000mAh Battery Budget Smartphones', '4ጂ ባለሁለት ሲም 5000ሚሊ አምፔር ባትሪ ያላቸው ስልኮች', '4G Siimii Lama Baatrii 5000mAh Ismaartifoonota', '4ጂ ክልተ ሲም 5000ሚሊ ኣምፔር ባትሪ ዘለዎም ፎናት', '4g-dual-sim-budget-smartphones', 'Budget smartphones featuring 4G LTE, dual SIM standby, and high capacity 5000mAh battery', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: wireless-fm-torch-feature-phones (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_FM_TORCH_SPEC', 'Wireless FM Radio & High-Beam Torch Feature Phones', 'Wireless FM Radio & High-Beam Torch Feature Phones', 'ያለ ሽቦ ኤፍኤም ሬዲዮ እና የብርሃን ባትሪ ያላቸው የቁልፍ ስልኮች', 'Raadiyoo FM fi Turchii Ifa Cimaa Bilbiloota Qabduu', 'ሽቦ ኣልቦ ኤፍኤም ሬድዮን ናይ ብርሃን ባትሪን ዘለዎም ፎናት', 'Feature phones with integrated wireless FM radio antennas and high-beam dual LED flashlights', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'heavy-duty-keypad-phones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'wireless-fm-torch-feature-phones', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'heavy-duty-keypad-phones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Wireless FM Radio & High-Beam Torch Feature Phones', 'ያለ ሽቦ ኤፍኤም ሬዲዮ እና የብርሃን ባትሪ ያላቸው የቁልፍ ስልኮች', 'Raadiyoo FM fi Turchii Ifa Cimaa Bilbiloota Qabduu', 'ሽቦ ኣልቦ ኤፍኤም ሬድዮን ናይ ብርሃን ባትሪን ዘለዎም ፎናት', 'wireless-fm-torch-feature-phones', 'Feature phones with integrated wireless FM radio antennas and high-beam dual LED flashlights', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: intel-core-i7-business-laptops (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_COMP_I7_BUSINESS_SPEC', 'Intel Core i7 16GB RAM Enterprise Laptops', 'Intel Core i7 16GB RAM Enterprise Laptops', 'ኢንቴል ኮር አይ7 16ጂቢ ራም ያላቸው የድርጅት ላፕቶፖች', 'Intel Core i7 16GB RAM Kompiitara Dhaabbilee', 'ኢንቴል ኮር ኣይ7 16ጂቢ ራም ዘለዎም ናይ ትካል ላፕቶፓት', 'High-performance enterprise laptops with 11th-13th Gen Intel Core i7 and 16GB RAM', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'business-work-ultrabooks' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'intel-core-i7-business-laptops', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'business-work-ultrabooks' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Intel Core i7 16GB RAM Enterprise Laptops', 'ኢንቴል ኮር አይ7 16ጂቢ ራም ያላቸው የድርጅት ላፕቶፖች', 'Intel Core i7 16GB RAM Kompiitara Dhaabbilee', 'ኢንቴል ኮር ኣይ7 16ጂቢ ራም ዘለዎም ናይ ትካል ላፕቶፓት', 'intel-core-i7-business-laptops', 'High-performance enterprise laptops with 11th-13th Gen Intel Core i7 and 16GB RAM', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: 55-inch-4k-smart-tvs (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TV_55INCH_4K_SPEC', '55-Inch 4K UHD Smart Televisions', '55-Inch 4K UHD Smart Televisions', '55 ኢንች 4ኬ አልትራ ኤችዲ ስማርት ቴሌቪዥኖች', 'Televiziyinii Ismaartii 55-Inchi 4K UHD', '55 ኢንች 4ኬ ኣልትራ ኤችዲ ስማርት ቴሌቪዥናት', '55-inch diagonal screen size 4K UHD smart televisions for living rooms', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '4k-ultra-hd-smart-tvs' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', '55-inch-4k-smart-tvs', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '4k-ultra-hd-smart-tvs' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('55-Inch 4K UHD Smart Televisions', '55 ኢንች 4ኬ አልትራ ኤችዲ ስማርት ቴሌቪዥኖች', 'Televiziyinii Ismaartii 55-Inchi 4K UHD', '55 ኢንች 4ኬ ኣልትራ ኤችዲ ስማርት ቴሌቪዥናት', '55-inch-4k-smart-tvs', '55-inch diagonal screen size 4K UHD smart televisions for living rooms', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: 1200cc-1600cc-automatic-sedans (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_1200_1600_AUTO_SPEC', '1.2L-1.6L Petrol Automatic Sedans', '1.2L-1.6L Petrol Automatic Sedans', '1.2-1.6 ሊትር ቤንዚን አውቶማቲክ መኪኖች', 'Sedaanota Peetiroolaa Otomaatiksii 1.2L-1.6L', '1.2-1.6 ሊትሮ ቤንዚን ኣውቶማቲክ መኻይን', 'Front-wheel drive petrol passenger cars with 1.2L to 1.6L 4-cylinder engines and automatic transmission', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'economy-city-sedans' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', '1200cc-1600cc-automatic-sedans', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'economy-city-sedans' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('1.2L-1.6L Petrol Automatic Sedans', '1.2-1.6 ሊትር ቤንዚን አውቶማቲክ መኪኖች', 'Sedaanota Peetiroolaa Otomaatiksii 1.2L-1.6L', '1.2-1.6 ሊትሮ ቤንዚን ኣውቶማቲክ መኻይን', '1200cc-1600cc-automatic-sedans', 'Front-wheel drive petrol passenger cars with 1.2L to 1.6L 4-cylinder engines and automatic transmission', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: 4-stroke-twin-spark-petrol-bajaj (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_4STROKE_200CC_BAJAJ', '4-Stroke 200cc-225cc Twin Spark Petrol Bajaj', '4-Stroke 200cc-225cc Twin Spark Petrol Bajaj', 'ባለ 4 ስትሮክ 200-225ሲሲ ባለ ሁለት ሻማ ቤንዚን ባጃጅ', 'Baajaajii Peetiroolaa 4-Stroke 200cc-225cc', 'ባዓል 4 ስትሮክ 200-225ሲሲ ቤንዚን ባጃጅ', 'High-torque 4-stroke single-cylinder petrol engine with DTS-i twin spark plug technology', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'commercial-passenger-bajaj' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', '4-stroke-twin-spark-petrol-bajaj', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'commercial-passenger-bajaj' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('4-Stroke 200cc-225cc Twin Spark Petrol Bajaj', 'ባለ 4 ስትሮክ 200-225ሲሲ ባለ ሁለት ሻማ ቤንዚን ባጃጅ', 'Baajaajii Peetiroolaa 4-Stroke 200cc-225cc', 'ባዓል 4 ስትሮክ 200-225ሲሲ ቤንዚን ባጃጅ', '4-stroke-twin-spark-petrol-bajaj', 'High-torque 4-stroke single-cylinder petrol engine with DTS-i twin spark plug technology', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: isuzu-npr-fsr-diesel-trucks (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_ISUZU_NPR_SPEC', 'Isuzu NPR 4.6L-5.2L Turbo Diesel 5-Ton Trucks', 'Isuzu NPR 4.6L-5.2L Turbo Diesel 5-Ton Trucks', 'አይሱዙ ኤንፒአር 4.6-5.2 ሊትር ቱርቦ ዲዝል 5-ቶን መኪኖች', 'Isuzu NPR 4.6L-5.2L Turbo Diizilii Toonii 5', 'ኣይሱዙ ኤንፒኣር 4.6-5.2 ሊትሮ ቱርቦ ዲዝል 5-ቶን መኻይን', 'Isuzu NPR series turbo diesel 5-speed manual 5,000kg payload commercial truck', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'isuzu-box-and-cargo-trucks' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'isuzu-npr-fsr-diesel-trucks', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'isuzu-box-and-cargo-trucks' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Isuzu NPR 4.6L-5.2L Turbo Diesel 5-Ton Trucks', 'አይሱዙ ኤንፒአር 4.6-5.2 ሊትር ቱርቦ ዲዝል 5-ቶን መኪኖች', 'Isuzu NPR 4.6L-5.2L Turbo Diizilii Toonii 5', 'ኣይሱዙ ኤንፒኣር 4.6-5.2 ሊትሮ ቱርቦ ዲዝል 5-ቶን መኻይን', 'isuzu-npr-fsr-diesel-trucks', 'Isuzu NPR series turbo diesel 5-speed manual 5,000kg payload commercial truck', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: bole-kazanchis-luxury-2bed-apartments (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PROP_BOLE_2BED_SPEC', 'Bole & Kazanchis Luxury 2-Bedroom Furnished Apartments', 'Bole & Kazanchis Luxury 2-Bedroom Furnished Apartments', 'ቦሌና ካዛንቺስ የቅንጦት ባለ 2 መኝታ አፓርታማዎች', 'Appaartamaalee Qananii Boolee fi Kaazaanichiis Kutaa 2', 'ቦሌን ካዛንቺስን ቅንጦት ባዓል 2 መደቀሲ ኣፓርታማታት', 'Prime residential locations in Bole and Kazanchis diplomatic corridors', 'PROPERTY', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'two-bedroom-furnished-apartments' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'bole-kazanchis-luxury-2bed-apartments', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'two-bedroom-furnished-apartments' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Bole & Kazanchis Luxury 2-Bedroom Furnished Apartments', 'ቦሌና ካዛንቺስ የቅንጦት ባለ 2 መኝታ አፓርታማዎች', 'Appaartamaalee Qananii Boolee fi Kaazaanichiis Kutaa 2', 'ቦሌን ካዛንቺስን ቅንጦት ባዓል 2 መደቀሲ ኣፓርታማታት', 'bole-kazanchis-luxury-2bed-apartments', 'Prime residential locations in Bole and Kazanchis diplomatic corridors', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: white-magna-teff-varieties (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_WHITE_MAGNA_SPEC', 'White Magna Extra-Fine Teff Varieties', 'White Magna Extra-Fine Teff Varieties', 'የነጭ ማኛ ምርጥ የጤፍ ዝርያዎች', 'Xaafii Adii Maanyaa Sadarkaa Ol-aanaa', 'ጻዕዳ ማኛ ብሉጽ ናይ ጣፍ ዓይነታት', 'Premium highland white Magna teff known for extra brightness, soft texture, and high market value', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'teff-grain-farming' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'white-magna-teff-varieties', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'teff-grain-farming' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('White Magna Extra-Fine Teff Varieties', 'የነጭ ማኛ ምርጥ የጤፍ ዝርያዎች', 'Xaafii Adii Maanyaa Sadarkaa Ol-aanaa', 'ጻዕዳ ማኛ ብሉጽ ናይ ጣፍ ዓይነታት', 'white-magna-teff-varieties', 'Premium highland white Magna teff known for extra brightness, soft texture, and high market value', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: yirgacheffe-specialty-arabica (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_YIRGACHEFFE_SPEC', 'Yirgacheffe High-Altitude Specialty Arabica', 'Yirgacheffe High-Altitude Specialty Arabica', 'የይርጋጨፌ ከፍተኛ ቦታ ስፔሻሊቲ አረቢካ ቡና', 'Buna Yirgaacaffee Arabikaa Qulqulluu', 'ናይ ይርጋጨፌ ላዕለዋይ ቦታ ስፔሻሊቲ ኣረቢካ ቡን', 'Grown at 1,900m-2,200m altitude with signature floral jasmine notes and vibrant citrus acidity', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'ethiopian-specialty-coffee' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'yirgacheffe-specialty-arabica', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'ethiopian-specialty-coffee' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Yirgacheffe High-Altitude Specialty Arabica', 'የይርጋጨፌ ከፍተኛ ቦታ ስፔሻሊቲ አረቢካ ቡና', 'Buna Yirgaacaffee Arabikaa Qulqulluu', 'ናይ ይርጋጨፌ ላዕለዋይ ቦታ ስፔሻሊቲ ኣረቢካ ቡን', 'yirgacheffe-specialty-arabica', 'Grown at 1,900m-2,200m altitude with signature floral jasmine notes and vibrant citrus acidity', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: sds-plus-rotary-hammer-drills-800w (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TOOL_SDS_PLUS_800W', 'SDS-Plus 800W-1000W Rotary Hammer Drills', 'SDS-Plus 800W-1000W Rotary Hammer Drills', 'ኤስዲኤስ ፕላስ 800-1000ዋት ሮታሪ ሀመር ድሪሎች', 'Boraasuu SDS-Plus 800W-1000W Cimaa', 'ኤስዲኤስ ፕላስ 800-1000ዋት ሮታሪ ሃመር ድሪላት', 'SDS-Plus quick-change chuck 800W electric hammer drills delivering 2.7 to 3.2 Joules impact energy', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'rotary-hammer-demolition-drills' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'sds-plus-rotary-hammer-drills-800w', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'rotary-hammer-demolition-drills' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('SDS-Plus 800W-1000W Rotary Hammer Drills', 'ኤስዲኤስ ፕላስ 800-1000ዋት ሮታሪ ሀመር ድሪሎች', 'Boraasuu SDS-Plus 800W-1000W Cimaa', 'ኤስዲኤስ ፕላስ 800-1000ዋት ሮታሪ ሃመር ድሪላት', 'sds-plus-rotary-hammer-drills-800w', 'SDS-Plus quick-change chuck 800W electric hammer drills delivering 2.7 to 3.2 Joules impact energy', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: ppc-325r-structural-masonry-cement (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_MAT_PPC_325R_SPEC', 'Portland Pozzolana Cement PPC 32.5R Grade', 'Portland Pozzolana Cement PPC 32.5R Grade', 'ፖርትላንድ ፖዞላና ሲሚንቶ ፒፒሲ 32.5አር ደረጃ', 'Simintoo Poortilaandii Pozoolaanaa Sadarkaa PPC 32.5R', 'ፖርትላንድ ፖዞላና ሲሚንቶ ፒፒሲ 32.5ኣር ደረጃ', 'PPC 32.5R composite cement for masonry, plastering, reinforced concrete slabs, and brickwork', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'portland-pozzolana-ordinary-cement' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'ppc-325r-structural-masonry-cement', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'portland-pozzolana-ordinary-cement' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Portland Pozzolana Cement PPC 32.5R Grade', 'ፖርትላንድ ፖዞላና ሲሚንቶ ፒፒሲ 32.5አር ደረጃ', 'Simintoo Poortilaandii Pozoolaanaa Sadarkaa PPC 32.5R', 'ፖርትላንድ ፖዞላና ሲሚንቶ ፒፒሲ 32.5ኣር ደረጃ', 'ppc-325r-structural-masonry-cement', 'PPC 32.5R composite cement for masonry, plastering, reinforced concrete slabs, and brickwork', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: pure-cotton-gold-tilfi-habesha-kemis (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_FASH_GOLD_TILFI_SPEC', 'Pure Handspun Cotton Gold & Silver Thread Tilfi Kemis', 'Pure Handspun Cotton Gold & Silver Thread Tilfi Kemis', 'ከንፁህ ጥጥ የወርቅና የብር ክር ጥልፍ የሀበሻ ቀሚስ', 'Qamisii Jirbii Qulqulluu Kirrii Warqee Tilibame', 'ካብ ንጹህ ጥጥ ናይ ወርቅን ብሩርን ፈትሊ ጥልፊ ዘለዎ ቀሚስ', '100% genuine Ethiopian handspun cotton (YeShema) woven with metallic gold and silver threads', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'handwoven-tilfi-habesha-kemis' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'pure-cotton-gold-tilfi-habesha-kemis', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'handwoven-tilfi-habesha-kemis' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Pure Handspun Cotton Gold & Silver Thread Tilfi Kemis', 'ከንፁህ ጥጥ የወርቅና የብር ክር ጥልፍ የሀበሻ ቀሚስ', 'Qamisii Jirbii Qulqulluu Kirrii Warqee Tilibame', 'ካብ ንጹህ ጥጥ ናይ ወርቅን ብሩርን ፈትሊ ጥልፊ ዘለዎ ቀሚስ', 'pure-cotton-gold-tilfi-habesha-kemis', '100% genuine Ethiopian handspun cotton (YeShema) woven with metallic gold and silver threads', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: digital-thermostat-non-stick-injera-mitad (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_HOME_DIGITAL_MITAD_SPEC', 'Digital Thermostat 58cm Non-Stick Injera Mitad', 'Digital Thermostat 58cm Non-Stick Injera Mitad', 'ዲጂታል ቴርሞስታት 58ሳሜ የማያሲዝ የእንጀራ ምጣድ', 'Eelee 58cm Iskiriinii fi To''annaa Dijitaalaa Qabu', 'ዲጂታል ቴርሞስታት 58ሳሜ ዘይለግብ ናይ እንጀራ ምጣድ', '58cm large diameter baking surface with automatic temperature sensor, stainless steel cover lid, and heat dial', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'energy-saving-electric-injera-mitad' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'digital-thermostat-non-stick-injera-mitad', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'energy-saving-electric-injera-mitad' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Digital Thermostat 58cm Non-Stick Injera Mitad', 'ዲጂታል ቴርሞስታት 58ሳሜ የማያሲዝ የእንጀራ ምጣድ', 'Eelee 58cm Iskiriinii fi To''annaa Dijitaalaa Qabu', 'ዲጂታል ቴርሞስታት 58ሳሜ ዘይለግብ ናይ እንጀራ ምጣድ', 'digital-thermostat-non-stick-injera-mitad', '58cm large diameter baking surface with automatic temperature sensor, stainless steel cover lid, and heat dial', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: pure-sine-wave-mppt-inverters (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SOLAR_MPPT_SPEC', 'Pure Sine Wave MPPT High-PV Inverters', 'Pure Sine Wave MPPT High-PV Inverters', 'ንፁህ ሳይን ዌቭ ኤምፒፒቲ ከፍተኛ ፒቪ ኢንቨርተሮች', 'Inveertara Saayin Weev Qulqulluu MPPT', 'ንጹህ ሳይን ዌቭ ኤምፒፒቲ ላዕለዋይ ፒቪ ኢንቨርተራት', 'Pure sine wave 230V AC output with high input PV voltage range (120V-450V DC) and MPPT tracking', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'off-grid-hybrid-solar-inverters' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'pure-sine-wave-mppt-inverters', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'off-grid-hybrid-solar-inverters' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Pure Sine Wave MPPT High-PV Inverters', 'ንፁህ ሳይን ዌቭ ኤምፒፒቲ ከፍተኛ ፒቪ ኢንቨርተሮች', 'Inveertara Saayin Weev Qulqulluu MPPT', 'ንጹህ ሳይን ዌቭ ኤምፒፒቲ ላዕለዋይ ፒቪ ኢንቨርተራት', 'pure-sine-wave-mppt-inverters', 'Pure sine wave 230V AC output with high input PV voltage range (120V-450V DC) and MPPT tracking', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: conduit-distribution-board-installation (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SERV_CONDUIT_DB_SPEC', 'Conduit Routing & 3-Phase Distribution Board Installation', 'Conduit Routing & 3-Phase Distribution Board Installation', 'የኮንዲዩት ቧንቧ ማለፊያና ባለ 3-ፌዝ ዲቢ ቦርድ ተከላ', 'Ujummoo Elektiriikii fi Gabatee Qoodinsaa 3-Phase Dhaabuu', 'ናይ ኮንዲዩት መስመርን ባዓል 3-ፌዝ ዲቢ ቦርድ ምትካልን', 'Wall chasing, conduit pipe laying, cable pulling (2.5mm²/4mm²/6mm²), and multi-breaker main panel installation', 'SERVICE', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'residential-commercial-electrical-wiring' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'conduit-distribution-board-installation', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'residential-commercial-electrical-wiring' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Conduit Routing & 3-Phase Distribution Board Installation', 'የኮንዲዩት ቧንቧ ማለፊያና ባለ 3-ፌዝ ዲቢ ቦርድ ተከላ', 'Ujummoo Elektiriikii fi Gabatee Qoodinsaa 3-Phase Dhaabuu', 'ናይ ኮንዲዩት መስመርን ባዓል 3-ፌዝ ዲቢ ቦርድ ምትካልን', 'conduit-distribution-board-installation', 'Wall chasing, conduit pipe laying, cable pulling (2.5mm²/4mm²/6mm²), and multi-breaker main panel installation', 4, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: 256gb-storage-12gb-ram-smartphones (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_256GB_12GB_RAM', '256GB Storage & 12GB RAM Flagship Phones', '256GB Storage & 12GB RAM Flagship Phones', '256ጂቢ ማከማቻና 12ጂቢ ራም ያላቸው ስልኮች', 'Kuusaa 256GB fi 12GB RAM Ismaartifoonota', '256ጂቢ መኽዘንን 12ጂቢ ራምን ዘለዎም ፎናት', '256GB internal storage paired with 12GB high-speed LPDDR5 RAM', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'dual-sim-flagship-amoled' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', '256gb-storage-12gb-ram-smartphones', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'dual-sim-flagship-amoled' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('256GB Storage & 12GB RAM Flagship Phones', '256ጂቢ ማከማቻና 12ጂቢ ራም ያላቸው ስልኮች', 'Kuusaa 256GB fi 12GB RAM Ismaartifoonota', '256ጂቢ መኽዘንን 12ጂቢ ራምን ዘለዎም ፎናት', '256gb-storage-12gb-ram-smartphones', '256GB internal storage paired with 12GB high-speed LPDDR5 RAM', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: 64gb-storage-4gb-ram-smartphones (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_64GB_4GB_RAM', '64GB Storage & 4GB RAM Budget Smartphones', '64GB Storage & 4GB RAM Budget Smartphones', '64ጂቢ ማከማቻና 4ጂቢ ራም ያላቸው ስማርት ስልኮች', 'Kuusaa 64GB fi 4GB RAM Ismaartifoonota', '64ጂቢ መኽዘንን 4ጂቢ ራምን ዘለዎም ስማርት ፎናት', '64GB internal flash storage with 4GB RAM and expandable microSD slot', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '4g-dual-sim-budget-smartphones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', '64gb-storage-4gb-ram-smartphones', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '4g-dual-sim-budget-smartphones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('64GB Storage & 4GB RAM Budget Smartphones', '64ጂቢ ማከማቻና 4ጂቢ ራም ያላቸው ስማርት ስልኮች', 'Kuusaa 64GB fi 4GB RAM Ismaartifoonota', '64ጂቢ መኽዘንን 4ጂቢ ራምን ዘለዎም ስማርት ፎናት', '64gb-storage-4gb-ram-smartphones', '64GB internal flash storage with 4GB RAM and expandable microSD slot', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: dual-sim-2500mah-keypad-phones (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_DUAL_SIM_2500MAH', 'Dual SIM 2500mAh+ Extended Battery Keypad Phones', 'Dual SIM 2500mAh+ Extended Battery Keypad Phones', 'ባለ ሁለት ሲም 2500+ ሚሊ አምፔር ረጅም ባትሪ የቁልፍ ስልኮች', 'Siimii Lama Baatrii 2500mAh+ Bilbiloota Qabduu', 'ክልተ ሲም 2500+ ሚሊ ኣምፔር ነዊሕ ባትሪ ዘለዎም ፎናት', 'Dual SIM feature phones with 2500mAh to 3000mAh removable batteries', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'wireless-fm-torch-feature-phones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'dual-sim-2500mah-keypad-phones', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'wireless-fm-torch-feature-phones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Dual SIM 2500mAh+ Extended Battery Keypad Phones', 'ባለ ሁለት ሲም 2500+ ሚሊ አምፔር ረጅም ባትሪ የቁልፍ ስልኮች', 'Siimii Lama Baatrii 2500mAh+ Bilbiloota Qabduu', 'ክልተ ሲም 2500+ ሚሊ ኣምፔር ነዊሕ ባትሪ ዘለዎም ፎናት', 'dual-sim-2500mah-keypad-phones', 'Dual SIM feature phones with 2500mAh to 3000mAh removable batteries', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: 16gb-ram-512gb-ssd-14inch-ultrabook (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_COMP_16GB_512GB_14INCH', '14-Inch Full HD 16GB RAM 512GB NVMe SSD Ultrabook', '14-Inch Full HD 16GB RAM 512GB NVMe SSD Ultrabook', '14 ኢንች ፉል ኤችዲ 16ጂቢ ራም 512ጂቢ ኤስኤስዲ አልትራቡክ', '14-Inchi Full HD 16GB RAM 512GB SSD Ultrabook', '14 ኢንች ፉል ኤችዲ 16ጂቢ ራም 512ጂቢ ኤስኤስዲ ኣልትራቡክ', '14.0-inch anti-glare FHD screen with fast 512GB NVMe M.2 Solid State Drive', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'intel-core-i7-business-laptops' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', '16gb-ram-512gb-ssd-14inch-ultrabook', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'intel-core-i7-business-laptops' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('14-Inch Full HD 16GB RAM 512GB NVMe SSD Ultrabook', '14 ኢንች ፉል ኤችዲ 16ጂቢ ራም 512ጂቢ ኤስኤስዲ አልትራቡክ', '14-Inchi Full HD 16GB RAM 512GB SSD Ultrabook', '14 ኢንች ፉል ኤችዲ 16ጂቢ ራም 512ጂቢ ኤስኤስዲ ኣልትራቡክ', '16gb-ram-512gb-ssd-14inch-ultrabook', '14.0-inch anti-glare FHD screen with fast 512GB NVMe M.2 Solid State Drive', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: google-tv-hdr10-55-inch (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TV_55INCH_GOOGLE_TV', '55-Inch Frameless Google TV & Tizen 4K HDR', '55-Inch Frameless Google TV & Tizen 4K HDR', '55 ኢንች ጠርዝ-አልባ ጎግል ቲቪ እና ታይዘን 4ኬ ኤችዲአር', '55-Inchi Daangaa-malee Google TV fi Tizen 4K HDR', '55 ኢንች ደረት ዘይብሉ ጎግል ቲቪን ታይዘንን 4ኬ', 'Bezel-less 55-inch smart TV with Google TV OS, voice remote, and Dolby Audio', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '55-inch-4k-smart-tvs' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'google-tv-hdr10-55-inch', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '55-inch-4k-smart-tvs' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('55-Inch Frameless Google TV & Tizen 4K HDR', '55 ኢንች ጠርዝ-አልባ ጎግል ቲቪ እና ታይዘን 4ኬ ኤችዲአር', '55-Inchi Daangaa-malee Google TV fi Tizen 4K HDR', '55 ኢንች ደረት ዘይብሉ ጎግል ቲቪን ታይዘንን 4ኬ', 'google-tv-hdr10-55-inch', 'Bezel-less 55-inch smart TV with Google TV OS, voice remote, and Dolby Audio', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: toyota-sedans-petrol-automatic (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_TOYOTA_SEDANS_AUTO', 'Toyota Petrol Automatic Sedans', 'Toyota Petrol Automatic Sedans', 'ቶዮታ ቤንዚን አውቶማቲክ አውቶሞቢሎች', 'Tooyootaa Peetiroolaa Otomaatiksii Sedaanota', 'ቶዮታ ቤንዚን ኣውቶማቲክ ሴዳን መኻይን', 'Toyota dependable automatic passenger sedans with high resale value and readily available parts', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '1200cc-1600cc-automatic-sedans' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'toyota-sedans-petrol-automatic', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '1200cc-1600cc-automatic-sedans' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Toyota Petrol Automatic Sedans', 'ቶዮታ ቤንዚን አውቶማቲክ አውቶሞቢሎች', 'Tooyootaa Peetiroolaa Otomaatiksii Sedaanota', 'ቶዮታ ቤንዚን ኣውቶማቲክ ሴዳን መኻይን', 'toyota-sedans-petrol-automatic', 'Toyota dependable automatic passenger sedans with high resale value and readily available parts', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: canvas-roof-manual-transmission-bajaj (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_CANVAS_MANUAL_BAJAJ', 'Canvas Roof 4-Speed Manual Reverse Bajaj', 'Canvas Roof 4-Speed Manual Reverse Bajaj', 'የሸራ ጣሪያ ባለ 4 ማርሽ ከነ ሪቨርስ ባጃጅ', 'Baajaajii Baangaa Sharaa Maarshii 4 Duubatti Deebisu', 'ናይ ሸራ ናሕሲ 4 ማርሽ ምስ ንድሕሪት መመለሲ ባጃጅ', 'Standard weather-proof canvas canopy body, 4 forward gears plus reverse hand lever', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '4-stroke-twin-spark-petrol-bajaj' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'canvas-roof-manual-transmission-bajaj', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '4-stroke-twin-spark-petrol-bajaj' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Canvas Roof 4-Speed Manual Reverse Bajaj', 'የሸራ ጣሪያ ባለ 4 ማርሽ ከነ ሪቨርስ ባጃጅ', 'Baajaajii Baangaa Sharaa Maarshii 4 Duubatti Deebisu', 'ናይ ሸራ ናሕሲ 4 ማርሽ ምስ ንድሕሪት መመለሲ ባጃጅ', 'canvas-roof-manual-transmission-bajaj', 'Standard weather-proof canvas canopy body, 4 forward gears plus reverse hand lever', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: isuzu-npr-dry-cargo-box-truck (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_ISUZU_NPR_BOX_BODY', 'Isuzu NPR Heavy Aluminum Dry Cargo Box Body', 'Isuzu NPR Heavy Aluminum Dry Cargo Box Body', 'አይሱዙ ኤንፒአር የአሉሚኒየም ዝግ የደረቅ ጭነት ሣጥን', 'Isuzu NPR Saanduqa Alumiiniyemii Cufaa', 'ኣይሱዙ ኤንፒኣር ናይ ኣሉሚንየም ዕፁው ሳንዱቕ መኪና', 'Isuzu NPR with weatherproof reinforced corrugated aluminum closed box and rear double doors', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'isuzu-npr-fsr-diesel-trucks' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'isuzu-npr-dry-cargo-box-truck', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'isuzu-npr-fsr-diesel-trucks' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Isuzu NPR Heavy Aluminum Dry Cargo Box Body', 'አይሱዙ ኤንፒአር የአሉሚኒየም ዝግ የደረቅ ጭነት ሣጥን', 'Isuzu NPR Saanduqa Alumiiniyemii Cufaa', 'ኣይሱዙ ኤንፒኣር ናይ ኣሉሚንየም ዕፁው ሳንዱቕ መኪና', 'isuzu-npr-dry-cargo-box-truck', 'Isuzu NPR with weatherproof reinforced corrugated aluminum closed box and rear double doors', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: full-generator-wifi-elevator-2bed-bole (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PROP_BOLE_AMENITIES_2BED', 'Full Backup Generator, Elevator & Security 2-Bed Bole', 'Full Backup Generator, Elevator & Security 2-Bed Bole', 'ሙሉ ጄነሬተር፣ ሊፍትና ጥበቃ ያለው ባለ 2 መኝታ ቦሌ', 'Jeeneereetara Guutuu, Wiifii fi Liiftii Kutaa 2 Boolee', 'ሙሉእ ጀነሬተር፣ ሊፍትን ድሕንነትን ዘለዎ ባዓል 2 መደቀሲ ቦሌ', 'Serviced apartment with 24/7 automatic standby power generator, modern elevator, and dedicated parking', 'PROPERTY', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'bole-kazanchis-luxury-2bed-apartments' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'full-generator-wifi-elevator-2bed-bole', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'bole-kazanchis-luxury-2bed-apartments' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Full Backup Generator, Elevator & Security 2-Bed Bole', 'ሙሉ ጄነሬተር፣ ሊፍትና ጥበቃ ያለው ባለ 2 መኝታ ቦሌ', 'Jeeneereetara Guutuu, Wiifii fi Liiftii Kutaa 2 Boolee', 'ሙሉእ ጀነሬተር፣ ሊፍትን ድሕንነትን ዘለዎ ባዓል 2 መደቀሲ ቦሌ', 'full-generator-wifi-elevator-2bed-bole', 'Serviced apartment with 24/7 automatic standby power generator, modern elevator, and dedicated parking', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: adaa-magna-grade-1-sifted-teff (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_ADAA_MAGNA_G1', 'Adaa Magna Grade 1 Machine Sifted Teff', 'Adaa Magna Grade 1 Machine Sifted Teff', 'የአደአ ማኛ አንደኛ ደረጃ በማሽን የተነፋ ጤፍ', 'Xaafii Ada''aa Maanyaa Sadarkaa 1ffaa Qulqullaa''e', 'ኣደኣ ማኛ ቀዳማይ ደረጃ ብማሽን ዝተነፈ ጣፍ', 'Adaa/Bishoftu certified Grade 1 white Magna teff, machine destoned and sifted without impurities', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'white-magna-teff-varieties' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'adaa-magna-grade-1-sifted-teff', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'white-magna-teff-varieties' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Adaa Magna Grade 1 Machine Sifted Teff', 'የአደአ ማኛ አንደኛ ደረጃ በማሽን የተነፋ ጤፍ', 'Xaafii Ada''aa Maanyaa Sadarkaa 1ffaa Qulqullaa''e', 'ኣደኣ ማኛ ቀዳማይ ደረጃ ብማሽን ዝተነፈ ጣፍ', 'adaa-magna-grade-1-sifted-teff', 'Adaa/Bishoftu certified Grade 1 white Magna teff, machine destoned and sifted without impurities', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: grade-1-sun-dried-yirgacheffe-green (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_YIRG_G1_WASHED', 'Grade 1 Wet-Washed Sun-Dried Yirgacheffe Green Coffee', 'Grade 1 Wet-Washed Sun-Dried Yirgacheffe Green Coffee', 'አንደኛ ደረጃ በውሃ የታጠበና በፀሐይ የደረቀ ይርጋጨፌ ቡና', 'Buna Yirgaacaffee Sadarkaa 1ffaa Dhiqame Aduun Goge', 'ቀዳማይ ደረጃ ብማይ ዝተሓጸበን ብፀሓይ ዝነቐጸን ይርጋጨፌ ቡን', 'Specialty Grade 1 washed green coffee with cup score 87+ on raised African drying beds', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'yirgacheffe-specialty-arabica' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'grade-1-sun-dried-yirgacheffe-green', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'yirgacheffe-specialty-arabica' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Grade 1 Wet-Washed Sun-Dried Yirgacheffe Green Coffee', 'አንደኛ ደረጃ በውሃ የታጠበና በፀሐይ የደረቀ ይርጋጨፌ ቡና', 'Buna Yirgaacaffee Sadarkaa 1ffaa Dhiqame Aduun Goge', 'ቀዳማይ ደረጃ ብማይ ዝተሓጸበን ብፀሓይ ዝነቐጸን ይርጋጨፌ ቡን', 'grade-1-sun-dried-yirgacheffe-green', 'Specialty Grade 1 washed green coffee with cup score 87+ on raised African drying beds', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: bosch-gbh-2-26-dre-rotary-hammer-kit (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TOOL_BOSCH_GBH_KIT', 'Bosch GBH 2-26 DRE Rotary Hammer with Chisel & Drill Bits Kit', 'Bosch GBH 2-26 DRE Rotary Hammer with Chisel & Drill Bits Kit', 'ቦሽ 2-26 ሮታሪ ሀመር ከነ ቺዝልና ቢቶች ኪት', 'Bosh GBH 2-26 DRE Boraasuu Qottoo fi Qabduu Waliin', 'ቦሽ 2-26 ሮታሪ ሃመር ምስ ቺዝልን ቢታትን ኪት', 'Bosch professional GBH 2-26 DRE 800W rotary hammer complete with carry case, SDS chisels, and bits', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'sds-plus-rotary-hammer-drills-800w' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'bosch-gbh-2-26-dre-rotary-hammer-kit', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'sds-plus-rotary-hammer-drills-800w' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Bosch GBH 2-26 DRE Rotary Hammer with Chisel & Drill Bits Kit', 'ቦሽ 2-26 ሮታሪ ሀመር ከነ ቺዝልና ቢቶች ኪት', 'Bosh GBH 2-26 DRE Boraasuu Qottoo fi Qabduu Waliin', 'ቦሽ 2-26 ሮታሪ ሃመር ምስ ቺዝልን ቢታትን ኪት', 'bosch-gbh-2-26-dre-rotary-hammer-kit', 'Bosch professional GBH 2-26 DRE 800W rotary hammer complete with carry case, SDS chisels, and bits', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: 50kg-bag-moisture-proof-ppc-cement (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_MAT_PPC_50KG_BAG', '50kg Laminated Moisture-Proof Bag PPC Cement', '50kg Laminated Moisture-Proof Bag PPC Cement', '50ኪግ እርጥበት መከላከያ ጆንያ ፒፒሲ ሲሚንቶ', 'Kootaa 50kg Jiidhinsa Qolatu Simintoo PPC', '50ኪግ ርጥበት መኸላኸሊ ጆንያ ፒፒሲ ሲሚንቶ', '50kg factory packaged multi-ply kraft paper or laminated polypropylene moisture-resistant sacks', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'ppc-325r-structural-masonry-cement' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', '50kg-bag-moisture-proof-ppc-cement', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'ppc-325r-structural-masonry-cement' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('50kg Laminated Moisture-Proof Bag PPC Cement', '50ኪግ እርጥበት መከላከያ ጆንያ ፒፒሲ ሲሚንቶ', 'Kootaa 50kg Jiidhinsa Qolatu Simintoo PPC', '50ኪግ ርጥበት መኸላኸሊ ጆንያ ፒፒሲ ሲሚንቶ', '50kg-bag-moisture-proof-ppc-cement', '50kg factory packaged multi-ply kraft paper or laminated polypropylene moisture-resistant sacks', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: full-set-habesha-kemis-with-matching-netela (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_FASH_KEMIS_NETELA_SET', 'Full Set Habesha Kemis with Matching Tilfi Netela Shawl', 'Full Set Habesha Kemis with Matching Tilfi Netela Shawl', 'ሙሉ የሀበሻ ቀሚስ ከተዛማጅ ጥልፍ ነጠላ ጋር', 'Qamisii Habashaa Guutuu fi Naxalaa Wal-Simu', 'ሙሉእ ናይ ሓበሻ ቀሚስ ምስ ተመሳሳሊ ጥልፊ ነጸላ', 'Two-piece ceremonial outfit comprising the full-length dress and a matching double-layered gauze shawl (Netela)', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'pure-cotton-gold-tilfi-habesha-kemis' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'full-set-habesha-kemis-with-matching-netela', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'pure-cotton-gold-tilfi-habesha-kemis' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Full Set Habesha Kemis with Matching Tilfi Netela Shawl', 'ሙሉ የሀበሻ ቀሚስ ከተዛማጅ ጥልፍ ነጠላ ጋር', 'Qamisii Habashaa Guutuu fi Naxalaa Wal-Simu', 'ሙሉእ ናይ ሓበሻ ቀሚስ ምስ ተመሳሳሊ ጥልፊ ነጸላ', 'full-set-habesha-kemis-with-matching-netela', 'Two-piece ceremonial outfit comprising the full-length dress and a matching double-layered gauze shawl (Netela)', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: 58cm-diameter-heavy-duty-injera-mitad (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_HOME_58CM_3800W_MITAD', '58cm Diameter 3800W Dual Element Heavy Duty Injera Mitad', '58cm Diameter 3800W Dual Element Heavy Duty Injera Mitad', '58 ሳንቲሜትር 3800ዋት ባለ ሁለት ማሞቂያ ከባድ እንጀራ ምጣድ', 'Eelee 58cm 3800W Hoo''isaa Dachaa Qabu Cimaa', '58 ሳንቲሜትር 3800ዋት ክልተ መሞቒ ዘለዎ ከቢድ ምጣድ', 'Heavy-duty electric griddle with reinforced 3800W spiral elements, safety ground, and ceramic base', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'digital-thermostat-non-stick-injera-mitad' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', '58cm-diameter-heavy-duty-injera-mitad', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'digital-thermostat-non-stick-injera-mitad' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('58cm Diameter 3800W Dual Element Heavy Duty Injera Mitad', '58 ሳንቲሜትር 3800ዋት ባለ ሁለት ማሞቂያ ከባድ እንጀራ ምጣድ', 'Eelee 58cm 3800W Hoo''isaa Dachaa Qabu Cimaa', '58 ሳንቲሜትር 3800ዋት ክልተ መሞቒ ዘለዎ ከቢድ ምጣድ', '58cm-diameter-heavy-duty-injera-mitad', 'Heavy-duty electric griddle with reinforced 3800W spiral elements, safety ground, and ceramic base', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: 48v-5kw-parallel-smart-bms-inverters (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SOLAR_48V_5KW_INVERTER', '48V 5kW-6kW Parallel LiFePO4 Smart Inverters', '48V 5kW-6kW Parallel LiFePO4 Smart Inverters', '48ቮልት 5ኪሎዋት ሊቲየም ባትሪ የሚቀበል ስማርት ኢንቨርተር', 'Inveertara 48V 5kW LiFePO4 Ismaartii', '48ቮልት 5ኪሎዋት ሊቲየም ባትሪ ዝቕበል ስማርት ኢንቨርተር', '48V DC battery system, 5000W continuous output, CAN/RS485 LiFePO4 battery communication, parallel capability', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'pure-sine-wave-mppt-inverters' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', '48v-5kw-parallel-smart-bms-inverters', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'pure-sine-wave-mppt-inverters' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('48V 5kW-6kW Parallel LiFePO4 Smart Inverters', '48ቮልት 5ኪሎዋት ሊቲየም ባትሪ የሚቀበል ስማርት ኢንቨርተር', 'Inveertara 48V 5kW LiFePO4 Ismaartii', '48ቮልት 5ኪሎዋት ሊቲየም ባትሪ ዝቕበል ስማርት ኢንቨርተር', '48v-5kw-parallel-smart-bms-inverters', '48V DC battery system, 5000W continuous output, CAN/RS485 LiFePO4 battery communication, parallel capability', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: certified-master-electrician-turnkey-service (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SERV_MASTER_ELECTRICIAN', 'Certified Master Electrician Complete Turnkey Wiring Package', 'Certified Master Electrician Complete Turnkey Wiring Package', 'የተመሰከረለት ዋና ባለሙያ ሙሉ የኤሌክትሪክ ስራ ጥቅል', 'Ogeessa Elektiriikii Eeyyamame Guutuu Hojii Sararaa', 'ዝተመስከረሉ ዋና ሞያተኛ ሙሉእ ናይ ኤሌክትሪክ ስራሕ ፓኬጅ', 'Licensed master electrician end-to-end service: load calculation, earth grounding, circuit testing, and commissioning', 'SERVICE', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'conduit-distribution-board-installation' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'certified-master-electrician-turnkey-service', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'conduit-distribution-board-installation' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Certified Master Electrician Complete Turnkey Wiring Package', 'የተመሰከረለት ዋና ባለሙያ ሙሉ የኤሌክትሪክ ስራ ጥቅል', 'Ogeessa Elektiriikii Eeyyamame Guutuu Hojii Sararaa', 'ዝተመስከረሉ ዋና ሞያተኛ ሙሉእ ናይ ኤሌክትሪክ ስራሕ ፓኬጅ', 'certified-master-electrician-turnkey-service', 'Licensed master electrician end-to-end service: load calculation, earth grounding, circuit testing, and commissioning', 5, v_ch_parent_id, true, false, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: samsung-galaxy-s24-ultra-5g-256gb (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_SAMSUNG_S24_ULTRA_256GB', 'Samsung Galaxy S24 Ultra 5G 256GB Titanium', 'Samsung Galaxy S24 Ultra 5G 256GB Titanium', 'ሳምሰንግ ጋላክሲ ኤስ24 አልትራ 5ጂ 256ጂቢ ታይታኒየም', 'Saamsangi Gaalaaksii S24 Ultra 5G 256GB Taayitaaniyem', 'ሳምሰንግ ጋላክሲ ኤስ24 ኣልትራ 5ጂ 256ጂቢ ታይታኒየም', 'Samsung Galaxy S24 Ultra 5G with S-Pen, 200MP camera, 256GB storage, and titanium frame', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '256gb-storage-12gb-ram-smartphones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'samsung-galaxy-s24-ultra-5g-256gb', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '256gb-storage-12gb-ram-smartphones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Samsung Galaxy S24 Ultra 5G 256GB Titanium', 'ሳምሰንግ ጋላክሲ ኤስ24 አልትራ 5ጂ 256ጂቢ ታይታኒየም', 'Saamsangi Gaalaaksii S24 Ultra 5G 256GB Taayitaaniyem', 'ሳምሰንግ ጋላክሲ ኤስ24 ኣልትራ 5ጂ 256ጂቢ ታይታኒየም', 'samsung-galaxy-s24-ultra-5g-256gb', 'Samsung Galaxy S24 Ultra 5G with S-Pen, 200MP camera, 256GB storage, and titanium frame', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: apple-iphone-15-pro-max-256gb-natural-titanium (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_IPHONE_15_PRO_MAX_256GB', 'Apple iPhone 15 Pro Max 256GB Natural Titanium', 'Apple iPhone 15 Pro Max 256GB Natural Titanium', 'አፕል አይፎን 15 ፕሮ ማክስ 256ጂቢ ናቹራል ታይታኒየም', 'Aappil Ayifoon 15 Pro Max 256GB Taayitaaniyem Uumamaa', 'ኣፕል ኣይፎን 15 ፕሮ ማክስ 256ጂቢ ናቹራል ታይታኒየም', 'Apple iPhone 15 Pro Max 256GB with A17 Pro chip, 5x telephoto camera, and USB-C', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '256gb-storage-12gb-ram-smartphones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'apple-iphone-15-pro-max-256gb-natural-titanium', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '256gb-storage-12gb-ram-smartphones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Apple iPhone 15 Pro Max 256GB Natural Titanium', 'አፕል አይፎን 15 ፕሮ ማክስ 256ጂቢ ናቹራል ታይታኒየም', 'Aappil Ayifoon 15 Pro Max 256GB Taayitaaniyem Uumamaa', 'ኣፕል ኣይፎን 15 ፕሮ ማክስ 256ጂቢ ናቹራል ታይታኒየም', 'apple-iphone-15-pro-max-256gb-natural-titanium', 'Apple iPhone 15 Pro Max 256GB with A17 Pro chip, 5x telephoto camera, and USB-C', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: tecno-spark-20-64gb-4gb-dual-sim (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_TECNO_SPARK_20', 'Tecno Spark 20 64GB Dual SIM 4G Smartphone', 'Tecno Spark 20 64GB Dual SIM 4G Smartphone', 'ቴክኖ ስፓርክ 20 64ጂቢ ባለሁለት ሲም 4ጂ ስማርት ስልክ', 'Teknoo Spark 20 64GB Siimii Lama 4G Ismaartifoonii', 'ቴክኖ ስፓርክ 20 64ጂቢ ክልተ ሲም 4ጂ ስማርት ፎን', 'Tecno Spark 20 with 64GB storage, 4GB RAM, 50MP main camera, and dual stereo speakers', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '64gb-storage-4gb-ram-smartphones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'tecno-spark-20-64gb-4gb-dual-sim', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '64gb-storage-4gb-ram-smartphones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Tecno Spark 20 64GB Dual SIM 4G Smartphone', 'ቴክኖ ስፓርክ 20 64ጂቢ ባለሁለት ሲም 4ጂ ስማርት ስልክ', 'Teknoo Spark 20 64GB Siimii Lama 4G Ismaartifoonii', 'ቴክኖ ስፓርክ 20 64ጂቢ ክልተ ሲም 4ጂ ስማርት ፎን', 'tecno-spark-20-64gb-4gb-dual-sim', 'Tecno Spark 20 with 64GB storage, 4GB RAM, 50MP main camera, and dual stereo speakers', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: infinix-hot-40i-64gb-4gb-dual-sim (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_INFINIX_HOT_40I', 'Infinix Hot 40i 64GB Dual SIM 4G Smartphone', 'Infinix Hot 40i 64GB Dual SIM 4G Smartphone', 'ኢንፊኒክስ ሆት 40አይ 64ጂቢ ባለሁለት ሲም 4ጂ ስማርት ስልክ', 'Infiniks Hot 40i 64GB Siimii Lama 4G Ismaartifoonii', 'ኢንፊኒክስ ሆት 40ኣይ 64ጂቢ ክልተ ሲም 4ጂ ስማርት ፎን', 'Infinix Hot 40i with 90Hz punch-hole display, 50MP camera, 5000mAh battery, and 18W fast charging', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '64gb-storage-4gb-ram-smartphones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'infinix-hot-40i-64gb-4gb-dual-sim', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '64gb-storage-4gb-ram-smartphones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Infinix Hot 40i 64GB Dual SIM 4G Smartphone', 'ኢንፊኒክስ ሆት 40አይ 64ጂቢ ባለሁለት ሲም 4ጂ ስማርት ስልክ', 'Infiniks Hot 40i 64GB Siimii Lama 4G Ismaartifoonii', 'ኢንፊኒክስ ሆት 40ኣይ 64ጂቢ ክልተ ሲም 4ጂ ስማርት ፎን', 'infinix-hot-40i-64gb-4gb-dual-sim', 'Infinix Hot 40i with 90Hz punch-hole display, 50MP camera, 5000mAh battery, and 18W fast charging', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: itel-it2160-dual-sim-torch-fm (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_ITEL_IT2160', 'itel it2160 Dual SIM Wireless FM Keypad Phone', 'itel it2160 Dual SIM Wireless FM Keypad Phone', 'አይቴል 2160 ባለ ሁለት ሲም ያለ ሽቦ ሬዲዮ የቁልፍ ስልክ', 'Aayitel 2160 Siimii Lama Raadiyoo FM Bilbila Qabduu', 'ኣይቴል 2160 ክልተ ሲም ሽቦ ኣልቦ ሬድዮ ናይ ቁልፊ ፎን', 'itel it2160 1.77-inch display, dual SIM, built-in bright torch, and wireless FM radio', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'dual-sim-2500mah-keypad-phones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'itel-it2160-dual-sim-torch-fm', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'dual-sim-2500mah-keypad-phones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('itel it2160 Dual SIM Wireless FM Keypad Phone', 'አይቴል 2160 ባለ ሁለት ሲም ያለ ሽቦ ሬዲዮ የቁልፍ ስልክ', 'Aayitel 2160 Siimii Lama Raadiyoo FM Bilbila Qabduu', 'ኣይቴል 2160 ክልተ ሲም ሽቦ ኣልቦ ሬድዮ ናይ ቁልፊ ፎን', 'itel-it2160-dual-sim-torch-fm', 'itel it2160 1.77-inch display, dual SIM, built-in bright torch, and wireless FM radio', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: nokia-105-africa-edition-dual-sim (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PHONE_NOKIA_105', 'Nokia 105 African Edition Dual SIM Keypad Phone', 'Nokia 105 African Edition Dual SIM Keypad Phone', 'ኖኪያ 105 አፍሪካን ኤዲሽን ባለ ሁለት ሲም የቁልፍ ስልክ', 'Nookiyaa 105 Afirikaa Edishinii Siimii Lama Bilbila', 'ኖኪያ 105 ኣፍሪካን ኤዲሽን ክልተ ሲም ናይ ቁልፊ ፎን', 'Nokia 105 durable island-keymat phone with tough polycarbonate body, FM radio, and long battery', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'dual-sim-2500mah-keypad-phones' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'nokia-105-africa-edition-dual-sim', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'dual-sim-2500mah-keypad-phones' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Nokia 105 African Edition Dual SIM Keypad Phone', 'ኖኪያ 105 አፍሪካን ኤዲሽን ባለ ሁለት ሲም የቁልፍ ስልክ', 'Nookiyaa 105 Afirikaa Edishinii Siimii Lama Bilbila', 'ኖኪያ 105 ኣፍሪካን ኤዲሽን ክልተ ሲም ናይ ቁልፊ ፎን', 'nokia-105-africa-edition-dual-sim', 'Nokia 105 durable island-keymat phone with tough polycarbonate body, FM radio, and long battery', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: hp-elitebook-840-g8-core-i7-16gb-512gb (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_COMP_HP_ELITEBOOK_840_G8', 'HP EliteBook 840 G8 Intel Core i7 16GB RAM 512GB SSD', 'HP EliteBook 840 G8 Intel Core i7 16GB RAM 512GB SSD', 'ኤችፒ ኤሊትቡክ 840 ጂ8 ኢንቴል ኮር አይ7 16ጂቢ 512ጂቢ', 'HP EliteBook 840 G8 Intel Core i7 16GB 512GB SSD', 'ኤችፒ ኤሊትቡክ 840 ጂ8 ኢንቴል ኮር ኣይ7 16ጂቢ 512ጂቢ', 'HP EliteBook 840 G8 aluminum chassis, Intel Core i7-1165G7, 16GB DDR4, 512GB NVMe SSD, backlit keyboard', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '16gb-ram-512gb-ssd-14inch-ultrabook' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'hp-elitebook-840-g8-core-i7-16gb-512gb', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '16gb-ram-512gb-ssd-14inch-ultrabook' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('HP EliteBook 840 G8 Intel Core i7 16GB RAM 512GB SSD', 'ኤችፒ ኤሊትቡክ 840 ጂ8 ኢንቴል ኮር አይ7 16ጂቢ 512ጂቢ', 'HP EliteBook 840 G8 Intel Core i7 16GB 512GB SSD', 'ኤችፒ ኤሊትቡክ 840 ጂ8 ኢንቴል ኮር ኣይ7 16ጂቢ 512ጂቢ', 'hp-elitebook-840-g8-core-i7-16gb-512gb', 'HP EliteBook 840 G8 aluminum chassis, Intel Core i7-1165G7, 16GB DDR4, 512GB NVMe SSD, backlit keyboard', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: lenovo-thinkpad-t14-core-i7-16gb-512gb (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_COMP_LENOVO_THINKPAD_T14', 'Lenovo ThinkPad T14 Intel Core i7 16GB RAM 512GB SSD', 'Lenovo ThinkPad T14 Intel Core i7 16GB RAM 512GB SSD', 'ሌኖቮ ቲንክፓድ ቲ14 ኢንቴል ኮር አይ7 16ጂቢ 512ጂቢ', 'Lenovo ThinkPad T14 Intel Core i7 16GB 512GB SSD', 'ሌኖቮ ቲንክፓድ ቲ14 ኢንቴል ኮር ኣይ7 16ጂቢ 512ጂቢ', 'Lenovo ThinkPad T14 military-grade durability, Core i7, 16GB RAM, 512GB SSD, TrackPoint navigation', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '16gb-ram-512gb-ssd-14inch-ultrabook' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'lenovo-thinkpad-t14-core-i7-16gb-512gb', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '16gb-ram-512gb-ssd-14inch-ultrabook' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Lenovo ThinkPad T14 Intel Core i7 16GB RAM 512GB SSD', 'ሌኖቮ ቲንክፓድ ቲ14 ኢንቴል ኮር አይ7 16ጂቢ 512ጂቢ', 'Lenovo ThinkPad T14 Intel Core i7 16GB 512GB SSD', 'ሌኖቮ ቲንክፓድ ቲ14 ኢንቴል ኮር ኣይ7 16ጂቢ 512ጂቢ', 'lenovo-thinkpad-t14-core-i7-16gb-512gb', 'Lenovo ThinkPad T14 military-grade durability, Core i7, 16GB RAM, 512GB SSD, TrackPoint navigation', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: samsung-55-inch-crystal-uhd-4k-smart-tv (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TV_SAMSUNG_55_CRYSTAL', 'Samsung 55-Inch Crystal UHD 4K Smart TV CU7000', 'Samsung 55-Inch Crystal UHD 4K Smart TV CU7000', 'ሳምሰንግ 55 ኢንች ክሪስታል 4ኬ ስማርት ቴሌቪዥን ሲዩ7000', 'Saamsangi 55-Inchi Crystal UHD 4K Smart TV CU7000', 'ሳምሰንግ 55 ኢንች ክሪስታል 4ኬ ስማርት ቴሌቪዥን ሲዩ7000', 'Samsung 55-Inch Crystal UHD 4K Processor, PurColor, Smart Hub, HDR10+, and Object Tracking Sound', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'google-tv-hdr10-55-inch' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'samsung-55-inch-crystal-uhd-4k-smart-tv', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'google-tv-hdr10-55-inch' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Samsung 55-Inch Crystal UHD 4K Smart TV CU7000', 'ሳምሰንግ 55 ኢንች ክሪስታል 4ኬ ስማርት ቴሌቪዥን ሲዩ7000', 'Saamsangi 55-Inchi Crystal UHD 4K Smart TV CU7000', 'ሳምሰንግ 55 ኢንች ክሪስታል 4ኬ ስማርት ቴሌቪዥን ሲዩ7000', 'samsung-55-inch-crystal-uhd-4k-smart-tv', 'Samsung 55-Inch Crystal UHD 4K Processor, PurColor, Smart Hub, HDR10+, and Object Tracking Sound', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: hisense-55-inch-4k-uhd-smart-google-tv (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TV_HISENSE_55_GOOGLE', 'Hisense 55-Inch 4K UHD Smart Google TV A6K Series', 'Hisense 55-Inch 4K UHD Smart Google TV A6K Series', 'ሀይሴንስ 55 ኢንች 4ኬ ስማርት ጎግል ቴሌቪዥን ኤ6ኬ', 'Haayisensi 55-Inchi 4K UHD Smart Google TV A6K', 'ሃይሴንስ 55 ኢንች 4ኬ ስማርት ጎግል ቴሌቪዥን ኤ6ኬ', 'Hisense 55-inch 4K Google TV with Dolby Vision, DTS Virtual:X sound, and Chromecast built-in', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'google-tv-hdr10-55-inch' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'hisense-55-inch-4k-uhd-smart-google-tv', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'google-tv-hdr10-55-inch' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Hisense 55-Inch 4K UHD Smart Google TV A6K Series', 'ሀይሴንስ 55 ኢንች 4ኬ ስማርት ጎግል ቴሌቪዥን ኤ6ኬ', 'Haayisensi 55-Inchi 4K UHD Smart Google TV A6K', 'ሃይሴንስ 55 ኢንች 4ኬ ስማርት ጎግል ቴሌቪዥን ኤ6ኬ', 'hisense-55-inch-4k-uhd-smart-google-tv', 'Hisense 55-inch 4K Google TV with Dolby Vision, DTS Virtual:X sound, and Chromecast built-in', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: toyota-corolla-16-automatic-sedan (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_TOYOTA_COROLLA_16', 'Toyota Corolla 1.6L Dual VVT-i Automatic Sedan', 'Toyota Corolla 1.6L Dual VVT-i Automatic Sedan', 'ቶዮታ ኮሮላ 1.6 ሊትር ቪቪቲ-አይ አውቶማቲክ ሴዳን', 'Tooyootaa Koroolaa 1.6L Otomaatiksii Sedaan', 'ቶዮታ ኮሮላ 1.6 ሊትሮ ኣውቶማቲክ ሴዳን', 'Toyota Corolla 1.6L Dual VVT-i petrol engine, 4-speed automatic, air conditioning, ABS, and power windows', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'toyota-sedans-petrol-automatic' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'toyota-corolla-16-automatic-sedan', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'toyota-sedans-petrol-automatic' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Toyota Corolla 1.6L Dual VVT-i Automatic Sedan', 'ቶዮታ ኮሮላ 1.6 ሊትር ቪቪቲ-አይ አውቶማቲክ ሴዳን', 'Tooyootaa Koroolaa 1.6L Otomaatiksii Sedaan', 'ቶዮታ ኮሮላ 1.6 ሊትሮ ኣውቶማቲክ ሴዳን', 'toyota-corolla-16-automatic-sedan', 'Toyota Corolla 1.6L Dual VVT-i petrol engine, 4-speed automatic, air conditioning, ABS, and power windows', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: suzuki-dzire-12-automatic-sedan (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_SUZUKI_DZIRE_12', 'Suzuki Dzire 1.2L DualJet AGS Automatic Sedan', 'Suzuki Dzire 1.2L DualJet AGS Automatic Sedan', 'ሱዙኪ ዲዛየር 1.2 ሊትር ዱአልጄት አውቶማቲክ ሴዳን', 'Suzukii Dizaayer 1.2L DualJet Otomaatiksii Sedaan', 'ሱዙኪ ዲዛየር 1.2 ሊትሮ ዱኣልጄት ኣውቶማቲክ ሴዳን', 'Suzuki Dzire 1.2L DualJet ultra-low fuel consumption, Auto Gear Shift (AGS), touchscreen infotainment', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'toyota-sedans-petrol-automatic' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'suzuki-dzire-12-automatic-sedan', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'toyota-sedans-petrol-automatic' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Suzuki Dzire 1.2L DualJet AGS Automatic Sedan', 'ሱዙኪ ዲዛየር 1.2 ሊትር ዱአልጄት አውቶማቲክ ሴዳን', 'Suzukii Dizaayer 1.2L DualJet Otomaatiksii Sedaan', 'ሱዙኪ ዲዛየር 1.2 ሊትሮ ዱኣልጄት ኣውቶማቲክ ሴዳን', 'suzuki-dzire-12-automatic-sedan', 'Suzuki Dzire 1.2L DualJet ultra-low fuel consumption, Auto Gear Shift (AGS), touchscreen infotainment', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: bajaj-re-205cc-dtsi-passenger-rickshaw (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_BAJAJ_RE_205CC', 'Bajaj RE 205cc DTS-i 4-Stroke Passenger Rickshaw', 'Bajaj RE 205cc DTS-i 4-Stroke Passenger Rickshaw', 'ባጃጅ አርኢ 205ሲሲ ዲቲኤስ-አይ ባለ 4 ስትሮክ ባጃጅ', 'Baajaajii RE 205cc DTS-i 4-Stroke Rikkishoo', 'ባጃጅ ኣርኢ 205ሲሲ ዲቲኤስ-ኣይ 4 ስትሮክ ባጃጅ', 'Bajaj RE 205cc DTS-i petrol engine, 4 forward 1 reverse, reinforced chassis, complete with spare wheel', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'canvas-roof-manual-transmission-bajaj' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'bajaj-re-205cc-dtsi-passenger-rickshaw', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'canvas-roof-manual-transmission-bajaj' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Bajaj RE 205cc DTS-i 4-Stroke Passenger Rickshaw', 'ባጃጅ አርኢ 205ሲሲ ዲቲኤስ-አይ ባለ 4 ስትሮክ ባጃጅ', 'Baajaajii RE 205cc DTS-i 4-Stroke Rikkishoo', 'ባጃጅ ኣርኢ 205ሲሲ ዲቲኤስ-ኣይ 4 ስትሮክ ባጃጅ', 'bajaj-re-205cc-dtsi-passenger-rickshaw', 'Bajaj RE 205cc DTS-i petrol engine, 4 forward 1 reverse, reinforced chassis, complete with spare wheel', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: tvs-king-deluxe-200cc-passenger-rickshaw (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_TVS_KING_DELUXE', 'TVS King Deluxe 200cc 4-Stroke Passenger Rickshaw', 'TVS King Deluxe 200cc 4-Stroke Passenger Rickshaw', 'ቲቪኤስ ኪንግ ዴሉክስ 200ሲሲ ባለ 4 ስትሮክ ባጃጅ', 'TVS King Deluxe 200cc 4-Stroke Rikkishoo', 'ቲቪኤስ ኪንግ ዴሉክስ 200ሲሲ 4 ስትሮክ ባጃጅ', 'TVS King Deluxe 200cc single cylinder 4-stroke engine, dual halogen headlights, wide passenger seating', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'canvas-roof-manual-transmission-bajaj' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'tvs-king-deluxe-200cc-passenger-rickshaw', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'canvas-roof-manual-transmission-bajaj' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('TVS King Deluxe 200cc 4-Stroke Passenger Rickshaw', 'ቲቪኤስ ኪንግ ዴሉክስ 200ሲሲ ባለ 4 ስትሮክ ባጃጅ', 'TVS King Deluxe 200cc 4-Stroke Rikkishoo', 'ቲቪኤስ ኪንግ ዴሉክስ 200ሲሲ 4 ስትሮክ ባጃጅ', 'tvs-king-deluxe-200cc-passenger-rickshaw', 'TVS King Deluxe 200cc single cylinder 4-stroke engine, dual halogen headlights, wide passenger seating', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: isuzu-npr-71-diesel-dry-cargo-box-truck (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_ISUZU_NPR71_COMPLETE', 'Isuzu NPR 71 4HG1 Diesel Dry Cargo Box Truck Complete', 'Isuzu NPR 71 4HG1 Diesel Dry Cargo Box Truck Complete', 'አይሱዙ ኤንፒአር 71 4ኤችጂ1 ዲዝል የደረቅ ጭነት ሳጥን መኪና', 'Isuzu NPR 71 4HG1 Diizilii Saanduqa Cufaa Guutuu', 'ኣይሱዙ ኤንፒኣር 71 4ኤችጂ1 ዲዝል ዕፁው ሳንዱቕ መኪና ሙሉእ', 'Isuzu NPR 71 4HG1 4.6L diesel engine, 5-ton payload, heavy aluminum enclosed box, ready for long-haul cargo', 'EQUIPMENT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'isuzu-npr-dry-cargo-box-truck' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'isuzu-npr-71-diesel-dry-cargo-box-truck', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'isuzu-npr-dry-cargo-box-truck' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Isuzu NPR 71 4HG1 Diesel Dry Cargo Box Truck Complete', 'አይሱዙ ኤንፒአር 71 4ኤችጂ1 ዲዝል የደረቅ ጭነት ሳጥን መኪና', 'Isuzu NPR 71 4HG1 Diizilii Saanduqa Cufaa Guutuu', 'ኣይሱዙ ኤንፒኣር 71 4ኤችጂ1 ዲዝል ዕፁው ሳንዱቕ መኪና ሙሉእ', 'isuzu-npr-71-diesel-dry-cargo-box-truck', 'Isuzu NPR 71 4HG1 4.6L diesel engine, 5-ton payload, heavy aluminum enclosed box, ready for long-haul cargo', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: bole-medhanialem-luxury-2bed-furnished-apt (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_PROP_BOLE_MEDHANIALEM_2BED', 'Bole Medhanialem 120m² Luxury 2-Bed Furnished Apartment', 'Bole Medhanialem 120m² Luxury 2-Bed Furnished Apartment', 'ቦሌ መድኃኔዓለም 120 ካሬ የቅንጦት ባለ 2 መኝታ አፓርታማ', 'Boolee Madhaani’aalam 120m² Appaartamaa Qananii Kutaa 2', 'ቦሌ መድኃኔዓለም 120 ካሬ ቅንጦት ባዓል 2 መደቀሲ ኣፓርታማ', '120m² interior layout, 2 master en-suite bedrooms, open kitchen with oven, standby generator, Bole Medhanialem', 'PROPERTY', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'full-generator-wifi-elevator-2bed-bole' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'bole-medhanialem-luxury-2bed-furnished-apt', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'full-generator-wifi-elevator-2bed-bole' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Bole Medhanialem 120m² Luxury 2-Bed Furnished Apartment', 'ቦሌ መድኃኔዓለም 120 ካሬ የቅንጦት ባለ 2 መኝታ አፓርታማ', 'Boolee Madhaani’aalam 120m² Appaartamaa Qananii Kutaa 2', 'ቦሌ መድኃኔዓለም 120 ካሬ ቅንጦት ባዓል 2 መደቀሲ ኣፓርታማ', 'bole-medhanialem-luxury-2bed-furnished-apt', '120m² interior layout, 2 master en-suite bedrooms, open kitchen with oven, standby generator, Bole Medhanialem', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: adaa-magna-grade-1-certified-50kg-bag (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_ADAA_MAGNA_50KG', 'Adaa Magna Grade 1 Cleaned Teff 50kg Laminated Bag', 'Adaa Magna Grade 1 Cleaned Teff 50kg Laminated Bag', 'አደአ ማኛ አንደኛ ደረጃ የተጣራ ጤፍ 50ኪግ ጆንያ', 'Ada''aa Maanyaa Sadarkaa 1ffaa Qulqulluu Kootaa 50kg', 'ኣደኣ ማኛ ቀዳማይ ደረጃ ዝተጻረየ ጣፍ 50ኪግ ጆንያ', 'Adaa Magna Grade 1 50kg polypropylene woven moisture-resistant sack, machine destoned, 100% pure', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'adaa-magna-grade-1-sifted-teff' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'adaa-magna-grade-1-certified-50kg-bag', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'adaa-magna-grade-1-sifted-teff' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Adaa Magna Grade 1 Cleaned Teff 50kg Laminated Bag', 'አደአ ማኛ አንደኛ ደረጃ የተጣራ ጤፍ 50ኪግ ጆንያ', 'Ada''aa Maanyaa Sadarkaa 1ffaa Qulqulluu Kootaa 50kg', 'ኣደኣ ማኛ ቀዳማይ ደረጃ ዝተጻረየ ጣፍ 50ኪግ ጆንያ', 'adaa-magna-grade-1-certified-50kg-bag', 'Adaa Magna Grade 1 50kg polypropylene woven moisture-resistant sack, machine destoned, 100% pure', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: yirgacheffe-kochere-g1-q88-green-60kg (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_YIRG_KOCHERE_60KG', 'Yirgacheffe Kochere Grade 1 Green Coffee 60kg Jute Bag', 'Yirgacheffe Kochere Grade 1 Green Coffee 60kg Jute Bag', 'ይርጋጨፌ ኮቸሬ አንደኛ ደረጃ ጥሬ ቡና 60ኪግ ጆንያ', 'Buna Yirgaacaffee Kocarree Sadarkaa 1ffaa 60kg', 'ይርጋጨፌ ኮቸሬ ቀዳማይ ደረጃ ጥረ ቡን 60ኪግ ጆንያ', 'Yirgacheffe Kochere micro-lot, Q-Score 88.5, GrainPro lined 60kg export jute sack, fresh crop', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'grade-1-sun-dried-yirgacheffe-green' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'yirgacheffe-kochere-g1-q88-green-60kg', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'grade-1-sun-dried-yirgacheffe-green' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Yirgacheffe Kochere Grade 1 Green Coffee 60kg Jute Bag', 'ይርጋጨፌ ኮቸሬ አንደኛ ደረጃ ጥሬ ቡና 60ኪግ ጆንያ', 'Buna Yirgaacaffee Kocarree Sadarkaa 1ffaa 60kg', 'ይርጋጨፌ ኮቸሬ ቀዳማይ ደረጃ ጥረ ቡን 60ኪግ ጆንያ', 'yirgacheffe-kochere-g1-q88-green-60kg', 'Yirgacheffe Kochere micro-lot, Q-Score 88.5, GrainPro lined 60kg export jute sack, fresh crop', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: bosch-professional-gbh-2-26-dre-800w-drill (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TOOL_BOSCH_GBH_226_COMPLETE', 'Bosch Professional GBH 2-26 DRE 800W Rotary Hammer Drill', 'Bosch Professional GBH 2-26 DRE 800W Rotary Hammer Drill', 'ቦሽ ፕሮፌሽናል ጂቢኤች 2-26 ዲአርኢ 800ዋት ሮታሪ ድሪል', 'Bosh Piroofeeshinaal GBH 2-26 DRE 800W Boraasuu', 'ቦሽ ፕሮፌሽናል ጂቢኤች 2-26 ዲኣርኢ 800ዋት ሮታሪ ድሪል', 'Genuine Bosch Professional GBH 2-26 DRE: 800W motor, 2.7 Joules, rotation stop for chiselling, 220V', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'bosch-gbh-2-26-dre-rotary-hammer-kit' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'bosch-professional-gbh-2-26-dre-800w-drill', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'bosch-gbh-2-26-dre-rotary-hammer-kit' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Bosch Professional GBH 2-26 DRE 800W Rotary Hammer Drill', 'ቦሽ ፕሮፌሽናል ጂቢኤች 2-26 ዲአርኢ 800ዋት ሮታሪ ድሪል', 'Bosh Piroofeeshinaal GBH 2-26 DRE 800W Boraasuu', 'ቦሽ ፕሮፌሽናል ጂቢኤች 2-26 ዲኣርኢ 800ዋት ሮታሪ ድሪል', 'bosch-professional-gbh-2-26-dre-800w-drill', 'Genuine Bosch Professional GBH 2-26 DRE: 800W motor, 2.7 Joules, rotation stop for chiselling, 220V', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: dangote-ppc-325r-structural-cement-50kg (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_MAT_DANGOTE_PPC_50KG', 'Dangote PPC 32.5R Portland Pozzolana Cement 50kg Bag', 'Dangote PPC 32.5R Portland Pozzolana Cement 50kg Bag', 'ዳንጎቴ ፒፒሲ 32.5አር ፖርትላንድ ፖዞላና ሲሚንቶ 50ኪግ', 'Dangootee PPC 32.5R Simintoo Poortilaandii 50kg', 'ዳንጎቴ ፒፒሲ 32.5ኣር ፖርትላንድ ፖዞላና ሲሚንቶ 50ኪግ', 'Dangote Cement Ethiopia PPC 32.5R high performance blended hydraulic cement, 50kg bag', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '50kg-bag-moisture-proof-ppc-cement' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'dangote-ppc-325r-structural-cement-50kg', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '50kg-bag-moisture-proof-ppc-cement' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Dangote PPC 32.5R Portland Pozzolana Cement 50kg Bag', 'ዳንጎቴ ፒፒሲ 32.5አር ፖርትላንድ ፖዞላና ሲሚንቶ 50ኪግ', 'Dangootee PPC 32.5R Simintoo Poortilaandii 50kg', 'ዳንጎቴ ፒፒሲ 32.5ኣር ፖርትላንድ ፖዞላና ሲሚንቶ 50ኪግ', 'dangote-ppc-325r-structural-cement-50kg', 'Dangote Cement Ethiopia PPC 32.5R high performance blended hydraulic cement, 50kg bag', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: mugher-ppc-325r-structural-cement-50kg (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_MAT_MUGHER_PPC_50KG', 'Mugher PPC 32.5R Portland Pozzolana Cement 50kg Bag', 'Mugher PPC 32.5R Portland Pozzolana Cement 50kg Bag', 'ሙገር ፒፒሲ 32.5አር ፖርትላንድ ፖዞላና ሲሚንቶ 50ኪግ', 'Muger PPC 32.5R Simintoo Poortilaandii 50kg', 'ሙገር ፒፒሲ 32.5ኣር ፖርትላንድ ፖዞላና ሲሚንቶ 50ኪግ', 'Mugher Cement Enterprise PPC 32.5R time-tested structural hydraulic cement, 50kg bag', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '50kg-bag-moisture-proof-ppc-cement' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'mugher-ppc-325r-structural-cement-50kg', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '50kg-bag-moisture-proof-ppc-cement' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Mugher PPC 32.5R Portland Pozzolana Cement 50kg Bag', 'ሙገር ፒፒሲ 32.5አር ፖርትላንድ ፖዞላና ሲሚንቶ 50ኪግ', 'Muger PPC 32.5R Simintoo Poortilaandii 50kg', 'ሙገር ፒፒሲ 32.5ኣር ፖርትላንድ ፖዞላና ሲሚንቶ 50ኪግ', 'mugher-ppc-325r-structural-cement-50kg', 'Mugher Cement Enterprise PPC 32.5R time-tested structural hydraulic cement, 50kg bag', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: luxury-shewa-handwoven-kemis-gold-tilfi-set (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_FASH_SHEWA_KEMIS_COMPLETE', 'Luxury Shewa Royal Handwoven Kemis with Gold Tilfi & Netela', 'Luxury Shewa Royal Handwoven Kemis with Gold Tilfi & Netela', 'ሮያል የሸዋ እጅ ጥልፍ የሀበሻ ቀሚስ ከነጠላ ጋር ሙሉ ስብስብ', 'Qamisii Aadaa Shawaa Mootii Warqeen Tilibame Naxalaa Waliin', 'ሮያል ናይ ሸዋ ኢድ ጥልፊ ናይ ሓበሻ ቀሚስ ምስ ነጸላ ሙሉእ', 'Master artisan crafted Shewa royal style Kemis with fine hand-spun cotton, dense gold thread embroidery, and matching Netela', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'full-set-habesha-kemis-with-matching-netela' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'luxury-shewa-handwoven-kemis-gold-tilfi-set', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'full-set-habesha-kemis-with-matching-netela' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Luxury Shewa Royal Handwoven Kemis with Gold Tilfi & Netela', 'ሮያል የሸዋ እጅ ጥልፍ የሀበሻ ቀሚስ ከነጠላ ጋር ሙሉ ስብስብ', 'Qamisii Aadaa Shawaa Mootii Warqeen Tilibame Naxalaa Waliin', 'ሮያል ናይ ሸዋ ኢድ ጥልፊ ናይ ሓበሻ ቀሚስ ምስ ነጸላ ሙሉእ', 'luxury-shewa-handwoven-kemis-gold-tilfi-set', 'Master artisan crafted Shewa royal style Kemis with fine hand-spun cotton, dense gold thread embroidery, and matching Netela', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: hulubedej-smart-digital-injera-mitad-pro-58cm (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_HOME_HULUBEDEJ_MITAD_PRO', 'Hulubedej Smart Digital Thermostat Injera Mitad Pro 58cm', 'Hulubedej Smart Digital Thermostat Injera Mitad Pro 58cm', 'ሁሉበደጅ ስማርት ዲጂታል ቴርሞስታት የእንጀራ ምጣድ ፕሮ 58ሳሜ', 'Hulubedej Ismaartii Eelee Buddeenaa Dijitaala Pro 58cm', 'ሁሉበደጅ ስማርት ዲጂታል ቴርሞስታት ናይ እንጀራ ምጣድ ፕሮ 58ሳሜ', 'Hulubedej Pro electric injera mitad with precision digital heat governor, 58cm non-stick ceramic baking disc, 3800W 220V', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '58cm-diameter-heavy-duty-injera-mitad' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'hulubedej-smart-digital-injera-mitad-pro-58cm', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '58cm-diameter-heavy-duty-injera-mitad' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Hulubedej Smart Digital Thermostat Injera Mitad Pro 58cm', 'ሁሉበደጅ ስማርት ዲጂታል ቴርሞስታት የእንጀራ ምጣድ ፕሮ 58ሳሜ', 'Hulubedej Ismaartii Eelee Buddeenaa Dijitaala Pro 58cm', 'ሁሉበደጅ ስማርት ዲጂታል ቴርሞስታት ናይ እንጀራ ምጣድ ፕሮ 58ሳሜ', 'hulubedej-smart-digital-injera-mitad-pro-58cm', 'Hulubedej Pro electric injera mitad with precision digital heat governor, 58cm non-stick ceramic baking disc, 3800W 220V', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: growatt-spf-5000es-5kw-48v-hybrid-inverter (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SOLAR_GROWATT_5000ES', 'Growatt SPF 5000ES 5kW 48V Off-Grid Hybrid Inverter', 'Growatt SPF 5000ES 5kW 48V Off-Grid Hybrid Inverter', 'ግሮዋት ኤስፒኤፍ 5000ኢኤስ 5ኪሎዋት 48ቮልት ሃይብሪድ ኢንቨርተር', 'Growatt SPF 5000ES 5kW 48V Inveertara Haayibriidii', 'ግሮዋት ኤስፒኤፍ 5000ኢኤስ 5ኪሎዋት 48ቮልት ሃይብሪድ ኢንቨርተር', 'Growatt SPF 5000ES: 5kW 230V output, 450V PV input, 100A MPPT, works with or without battery, WiFi monitoring', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '48v-5kw-parallel-smart-bms-inverters' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'growatt-spf-5000es-5kw-48v-hybrid-inverter', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '48v-5kw-parallel-smart-bms-inverters' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Growatt SPF 5000ES 5kW 48V Off-Grid Hybrid Inverter', 'ግሮዋት ኤስፒኤፍ 5000ኢኤስ 5ኪሎዋት 48ቮልት ሃይብሪድ ኢንቨርተር', 'Growatt SPF 5000ES 5kW 48V Inveertara Haayibriidii', 'ግሮዋት ኤስፒኤፍ 5000ኢኤስ 5ኪሎዋት 48ቮልት ሃይብሪድ ኢንቨርተር', 'growatt-spf-5000es-5kw-48v-hybrid-inverter', 'Growatt SPF 5000ES: 5kW 230V output, 450V PV input, 100A MPPT, works with or without battery, WiFi monitoring', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

    -- --- Node: residential-full-building-wiring-and-breaker-setup (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SERV_FULL_WIRING_SERVICE', 'Turnkey Residential Building Electrical Wiring & Breaker Setup', 'Turnkey Residential Building Electrical Wiring & Breaker Setup', 'ሙሉ የመኖሪያ ቤት ኤሌክትሪክ መስመር ዝርጋታና ቆጣሪ ተከላ አገልግሎት', 'Tajaajila Guutuu Sarara Elektiriikii Mana Jireenyaa fi Qoorannoo', 'ሙሉእ ናይ መነበሪ ገዛ ኤሌክትሪክ መስመር ዝርግሓን ቆጻሪ ምትካልን ኣገልግሎት', 'Complete electrical contracting service for G+0 to G+2 residences including materials inspection, wiring, breakers, earth rod, and certificate', 'SERVICE', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        concept_type = EXCLUDED.concept_type,
        status = 'ACTIVE'
    RETURNING id INTO v_cid;

    -- 2. Placement in SYS_PROD
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'certified-master-electrician-turnkey-service' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'residential-full-building-wiring-and-breaker-setup', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf,
        is_active = true,
        status = 'PUBLISHED'
    RETURNING id INTO v_pid;

    -- 3. Category Hierarchy Sync
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'certified-master-electrician-turnkey-service' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Turnkey Residential Building Electrical Wiring & Breaker Setup', 'ሙሉ የመኖሪያ ቤት ኤሌክትሪክ መስመር ዝርጋታና ቆጣሪ ተከላ አገልግሎት', 'Tajaajila Guutuu Sarara Elektiriikii Mana Jireenyaa fi Qoorannoo', 'ሙሉእ ናይ መነበሪ ገዛ ኤሌክትሪክ መስመር ዝርግሓን ቆጻሪ ምትካልን ኣገልግሎት', 'residential-full-building-wiring-and-breaker-setup', 'Complete electrical contracting service for G+0 to G+2 residences including materials inspection, wiring, breakers, earth rod, and certificate', 6, v_ch_parent_id, true, true, 'SYS_PROD')
    ON CONFLICT (slug) DO UPDATE SET
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description,
        level = EXCLUDED.level,
        parent_id = EXCLUDED.parent_id,
        is_active = true,
        is_leaf = EXCLUDED.is_leaf;

END $$;

-- 4. Mirror vertical domain placements into specialized schemes
DO $$
DECLARE
    v_prod_id INT;
    v_elec_id INT;
    v_veh_id  INT;
    v_agri_id INT;
    v_tool_id INT;
    v_prop_id INT;
    v_serv_id INT;
    v_root_pid INT;
BEGIN
    SELECT id INTO v_prod_id FROM classification_schemes WHERE code = 'SYS_PROD';
    SELECT id INTO v_elec_id FROM classification_schemes WHERE code = 'SYS_ELEC';
    SELECT id INTO v_veh_id  FROM classification_schemes WHERE code = 'SYS_VEH';
    SELECT id INTO v_agri_id FROM classification_schemes WHERE code = 'SYS_AGRI';
    SELECT id INTO v_tool_id FROM classification_schemes WHERE code = 'SYS_TOOLS';
    SELECT id INTO v_prop_id FROM classification_schemes WHERE code = 'SYS_PROP';
    SELECT id INTO v_serv_id FROM classification_schemes WHERE code = 'SYS_SERV';

    -- Electronics mirror
    IF v_elec_id IS NOT NULL THEN
        SELECT id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_elec_id AND depth = 0 LIMIT 1;
        IF v_root_pid IS NOT NULL THEN
            INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
            SELECT v_elec_id, cp.concept_id, v_root_pid, cp.role, 'elec-' || cp.slug, cp.depth, cp.is_leaf, true, 'PUBLISHED'
            FROM classification_placements cp
            WHERE cp.scheme_id = v_prod_id AND (cp.slug LIKE '%phone%' OR cp.slug LIKE '%laptop%' OR cp.slug LIKE '%tv%' OR cp.slug LIKE '%samsung%' OR cp.slug LIKE '%apple%')
              AND cp.depth > 0
            ON CONFLICT (scheme_id, slug) DO NOTHING;
        END IF;
    END IF;

    -- Vehicles mirror
    IF v_veh_id IS NOT NULL THEN
        SELECT id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_veh_id AND depth = 0 LIMIT 1;
        IF v_root_pid IS NOT NULL THEN
            INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
            SELECT v_veh_id, cp.concept_id, v_root_pid, cp.role, 'veh-' || cp.slug, cp.depth, cp.is_leaf, true, 'PUBLISHED'
            FROM classification_placements cp
            WHERE cp.scheme_id = v_prod_id AND (cp.slug LIKE '%corolla%' OR cp.slug LIKE '%dzire%' OR cp.slug LIKE '%bajaj%' OR cp.slug LIKE '%isuzu%' OR cp.slug LIKE '%sedan%' OR cp.slug LIKE '%car%')
              AND cp.depth > 0
            ON CONFLICT (scheme_id, slug) DO NOTHING;
        END IF;
    END IF;

    -- Agriculture mirror
    IF v_agri_id IS NOT NULL THEN
        SELECT id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_agri_id AND depth = 0 LIMIT 1;
        IF v_root_pid IS NOT NULL THEN
            INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
            SELECT v_agri_id, cp.concept_id, v_root_pid, cp.role, 'agri-' || cp.slug, cp.depth, cp.is_leaf, true, 'PUBLISHED'
            FROM classification_placements cp
            WHERE cp.scheme_id = v_prod_id AND (cp.slug LIKE '%teff%' OR cp.slug LIKE '%coffee%' OR cp.slug LIKE '%grain%' OR cp.slug LIKE '%dairy%')
              AND cp.depth > 0
            ON CONFLICT (scheme_id, slug) DO NOTHING;
        END IF;
    END IF;

    -- Tools mirror
    IF v_tool_id IS NOT NULL THEN
        SELECT id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_tool_id AND depth = 0 LIMIT 1;
        IF v_root_pid IS NOT NULL THEN
            INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
            SELECT v_tool_id, cp.concept_id, v_root_pid, cp.role, 'tool-' || cp.slug, cp.depth, cp.is_leaf, true, 'PUBLISHED'
            FROM classification_placements cp
            WHERE cp.scheme_id = v_prod_id AND (cp.slug LIKE '%drill%' OR cp.slug LIKE '%bosch%' OR cp.slug LIKE '%tool%' OR cp.slug LIKE '%hammer%')
              AND cp.depth > 0
            ON CONFLICT (scheme_id, slug) DO NOTHING;
        END IF;
    END IF;
END $$;

-- 5. Re-evaluate is_leaf status accurately across all placements and categories efficiently
UPDATE classification_placements p
SET is_leaf = false
WHERE is_leaf = true
  AND EXISTS (SELECT 1 FROM classification_placements child WHERE child.parent_placement_id = p.id);

UPDATE classification_placements p
SET is_leaf = true
WHERE is_leaf = false
  AND NOT EXISTS (SELECT 1 FROM classification_placements child WHERE child.parent_placement_id = p.id);

UPDATE category_hierarchy ch
SET is_leaf = false
WHERE is_leaf = true
  AND EXISTS (SELECT 1 FROM category_hierarchy child WHERE child.parent_id = ch.id);

UPDATE category_hierarchy ch
SET is_leaf = true
WHERE is_leaf = false
  AND NOT EXISTS (SELECT 1 FROM category_hierarchy child WHERE child.parent_id = ch.id);

-- 6. Ensure depths and roles match strict tree hierarchy
UPDATE classification_placements p
SET depth = t.lvl
FROM (
    WITH RECURSIVE tree AS (
        SELECT id, 0 AS lvl FROM classification_placements WHERE parent_placement_id IS NULL
        UNION ALL
        SELECT p.id, t.lvl + 1 FROM classification_placements p JOIN tree t ON p.parent_placement_id = t.id
    )
    SELECT id, lvl FROM tree
) t
WHERE p.id = t.id AND p.depth <> t.lvl;

UPDATE classification_placements SET role = 'ROOT' WHERE depth = 0 AND role <> 'ROOT';
UPDATE classification_placements SET role = 'SUPERFAMILY' WHERE depth = 1 AND role <> 'SUPERFAMILY';
UPDATE classification_placements SET role = 'CATEGORY' WHERE depth = 2 AND role <> 'CATEGORY';
UPDATE classification_placements SET role = 'SUBCATEGORY' WHERE depth = 3 AND role <> 'SUBCATEGORY';
UPDATE classification_placements SET role = 'SPECIFICATION' WHERE depth = 4 AND role <> 'SPECIFICATION';
UPDATE classification_placements SET role = 'VARIANT' WHERE depth = 5 AND role <> 'VARIANT';
UPDATE classification_placements SET role = 'LEAF' WHERE depth >= 6 AND role <> 'LEAF';

COMMIT;
