-- ============================================================================
-- 083_deep_hierarchical_expansion_layers_3_4_5_6.sql
-- Universal Multi-Layer Deep Classification Architecture Expansion
-- Features:
--   1. Align placement semantic roles strictly to depths:
--      0=ROOT, 1=SUPERFAMILY, 2=CATEGORY, 3=SUBCATEGORY, 4=SPECIFICATION, 5=VARIANT, 6=LEAF
--   2. Seed specialized scheme hierarchies (SYS_TOOLS, SYS_ELEC, SYS_AGRI, SYS_VEH, SYS_PROP, SYS_SERV)
--   3. Deep-expand key Ethiopian marketplace domains through Layer 6 with 4-language naming
--   4. Re-calculate placement paths and depths idempotently
-- ============================================================================

BEGIN;

-- 1. Update existing placement roles according to depth across all schemes
UPDATE classification_placements SET role = 'ROOT' WHERE depth = 0 AND role <> 'ROOT';
UPDATE classification_placements SET role = 'SUPERFAMILY' WHERE depth = 1 AND role <> 'SUPERFAMILY';
UPDATE classification_placements SET role = 'CATEGORY' WHERE depth = 2 AND role <> 'CATEGORY';
UPDATE classification_placements SET role = 'SUBCATEGORY' WHERE depth = 3 AND role <> 'SUBCATEGORY';
UPDATE classification_placements SET role = 'SPECIFICATION' WHERE depth = 4 AND role <> 'SPECIFICATION';
UPDATE classification_placements SET role = 'VARIANT' WHERE depth = 5 AND role <> 'VARIANT';
UPDATE classification_placements SET role = 'LEAF' WHERE depth = 6 AND role <> 'LEAF';

-- 2. Ensure all 9 classification schemes are active and present
INSERT INTO classification_schemes (code, name_en, name_am, name_om, name_ti, description, is_system, status)
VALUES
('SYS_PROD', 'Universal Product Classification', 'ሁለንተናዊ የምርት ምደባ', 'Ramaddii Oomishaa Waliigalaa', 'ኩለመዳላዊ ምደባ ፍርያት', 'Master product marketplace classification scheme', true, 'ACTIVE'),
('SYS_IND', 'Industry Classification', 'የኢንዱስትሪ ምደባ', 'Ramaddii Industirii', 'ምደባ ኢንዱስትሪ', 'Economic sector and business vertical taxonomy', true, 'ACTIVE'),
('SYS_SERV', 'Services & Trades Classification', 'የአገልግሎትና የሙያ ምደባ', 'Ramaddii Tajaajilootaa', 'ምደባ ኣገልግሎታትን ሞያታትን', 'Service providers, contractors and professional trades', true, 'ACTIVE'),
('SYS_PROP', 'Real Estate & Property Classification', 'የሪል እስቴትና የንብረት ምደባ', 'Ramaddii Qabeenya fi Manneen', 'ምደባ ሪል እስቴትን ንብረትን', 'Residential, commercial, and land property classification', true, 'ACTIVE'),
('SYS_VEH', 'Automotive & Mobility Classification', 'የተሽከርካሪና የትራንስፖርት ምደባ', 'Ramaddii Konkolaattotaa', 'ምደባ ተሽከርከርትን ትራንስፖርትን', 'Vehicles, heavy machinery, and spare parts taxonomy', true, 'ACTIVE'),
('SYS_AGRI', 'Agriculture & Farming Classification', 'የግብርናና የእርሻ ምደባ', 'Ramaddii Qonnaa', 'ምደባ ሕርሻን ልምዓትን', 'Crops, livestock, farm machinery and agricultural commodities', true, 'ACTIVE'),
('SYS_JOBS', 'Workforce & Employment Classification', 'የስራና የቅጥር ምደባ', 'Ramaddii Hojii', 'ምደባ ስራሕን ስራሕተኛታትን', 'Occupations, qualifications, and freelance labor taxonomy', true, 'ACTIVE'),
('SYS_TOOLS', 'Tools & Machinery Classification', 'የመሳሪያዎችና ማሽነሪዎች ምደባ', 'Ramaddii Meeshaalee fi Maashinootaa', 'ምደባ መሳርሒታትን ማሽነሪታትን', 'Universal tools, power equipment, and machinery taxonomy', false, 'ACTIVE'),
('SYS_ELEC', 'Electronics Classification', 'የኤሌክትሮኒክስ ምደባ', 'Ramaddii Eleektirooniksii', 'ምደባ ኤሌክትሮኒክስ', 'Universal consumer electronics, devices, and computing hardware', false, 'ACTIVE')
ON CONFLICT (code) DO UPDATE SET
    name_en = EXCLUDED.name_en,
    name_am = EXCLUDED.name_am,
    name_om = EXCLUDED.name_om,
    name_ti = EXCLUDED.name_ti,
    status = 'ACTIVE';

-- 3. Deep-expand Core Market Domains down through Layers 3, 4, 5, and 6
DO $$
DECLARE
    v_scheme_id INT;
    v_root_pid  INT;
    v_parent_pid INT;
    v_curr_cid   INT;
    v_curr_pid   INT;
BEGIN
    -- Chain: agriculture under agriculture
    SELECT id INTO v_scheme_id FROM classification_schemes WHERE code = 'SYS_PROD';
    SELECT id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'agriculture' LIMIT 1;
    v_parent_pid := v_root_pid;

    -- Layer 1: SUPERFAMILY (Cereal & Grain Farming)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_GRAINS', 'Cereal & Grain Farming', 'Cereal & Grain Farming', 'የእህል ሰብል እርሻ', 'Qonna Midhaanii', 'ሕርሻ እኽሊ', 'Cereal grains, staple flours, and harvest crops', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUPERFAMILY', 'cereal-grain-farming', 1, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 2: CATEGORY (Teff Farming & Trade)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_TEFF', 'Teff Farming & Trade', 'Teff Farming & Trade', 'የጤፍ እርሻና ንግድ', 'Qonna fi Daldala Xaafii', 'ሕርሻን ንግድን ጣፍ', 'Eragrostis tef cultivation, grain varieties, and distribution', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'CATEGORY', 'teff-farming-trade', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 3: SUBCATEGORY (White Magna Teff)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_WHITE_MAGNA_TEFF', 'White Magna Teff', 'White Magna Teff', 'የነጭ ማኛ ጤፍ', 'Xaafii Adii Maanyaa', 'ጻዕዳ ማኛ ጣፍ', 'Premium high-altitude extra-white Magna teff grains', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUBCATEGORY', 'white-magna-teff', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 4: SPECIFICATION (Adaa Magna Grade-1 Sifted)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_ADAA_MAGNA_G1', 'Adaa Magna Grade-1 Sifted', 'Adaa Magna Grade-1 Sifted', 'የአደአ ማኛ አንደኛ ደረጃ የተነፋ', 'Xaafii Maanyaa Ada''aa Sadarkaa 1ffaa', 'ኣደኣ ማኛ ቀዳማይ ብርኪ ዝተነጸሀ', 'Traditional Adaa Bishoftu origin Grade-1 mechanical double-sifted Magna teff', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SPECIFICATION', 'adaa-magna-grade-1-sifted', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 5: VARIANT (Machine-Cleaned 50kg Polypropylene Bag)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_MAGNA_50KG_BAG', 'Machine-Cleaned 50kg Polypropylene Bag', 'Machine-Cleaned 50kg Polypropylene Bag', 'በማሽን የተጠራ 50 ኪ.ግ ጆንያ', 'Qulqullina Maashiniitiin Saamsame 50kg', 'ብማሽን ዝተነጸሀ 50 ኪ.ግ ጆንያ', 'Stone-free 99.8% pure machine-cleaned teff packed in heavy-duty 50kg bags', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'VARIANT', 'machine-cleaned-50kg-bag', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 6: LEAF (Export-Grade Certified Organic Adaa Magna 50kg Sealed Batch)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_EXPORT_CERT_ADAA_50KG', 'Export-Grade Certified Organic Adaa Magna 50kg Sealed Batch', 'Export-Grade Certified Organic Adaa Magna 50kg Sealed Batch', 'የውጭ ገበያ ማረጋገጫ ያለው አደአ ማኛ 50 ኪ.ግ የታሸገ', 'Xaafii Maanyaa Qulqulluu Export 50kg Saamsame', 'ሰደድ-ብቑዕ ማረጋገጺ ዘለዎ ኣደኣ ማኛ 50 ኪ.ግ ዝተዓሸገ', 'Certified organic Adaa highlands export batch with phytosanitary compliance and barcode traceability', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'LEAF', 'export-grade-certified-adaa-magna-50kg', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Chain: agriculture under agriculture
    SELECT id INTO v_scheme_id FROM classification_schemes WHERE code = 'SYS_PROD';
    SELECT id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'agriculture' LIMIT 1;
    v_parent_pid := v_root_pid;

    -- Layer 1: SUPERFAMILY (Coffee & Tea Production)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_COFFEE_TEA', 'Coffee & Tea Production', 'Coffee & Tea Production', 'የቡናና ሻይ ምርት', 'Oomisha Bunaa fi Shaayii', 'ፍርያት ቡንን ሻይን', 'Green and roasted Arabica coffee, highland tea, and spices', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUPERFAMILY', 'coffee-tea-production', 1, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 2: CATEGORY (Specialty Arabica Coffee)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_SPECIALTY_COFFEE', 'Specialty Arabica Coffee', 'Specialty Arabica Coffee', 'ስፔሻሊቲ አረቢካ ቡና', 'Buna Ispeeshaalitii Arabikaa', 'ስፔሻሊቲ ዓረቢካ ቡን', 'High altitude Ethiopian single-origin specialty Arabica green beans', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'CATEGORY', 'specialty-arabica-coffee', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 3: SUBCATEGORY (Yirgacheffe & Sidama Washed Coffee)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_YIRGACHEFFE_WASHED', 'Yirgacheffe & Sidama Washed Coffee', 'Yirgacheffe & Sidama Washed Coffee', 'የይርጋጨፌና ሲዳማ የታጠበ ቡና', 'Buna Yirgaacafee fi Sidaamaa Miiccame', 'ቡን ይርጋጨፌን ሲዳማን ዝተሓጸበ', 'Wet-processed washed specialty coffee with bright citrus and floral cup profile', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUBCATEGORY', 'yirgacheffe-washed-coffee', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 4: SPECIFICATION (Grade-1 Natural Sun-Dried Yirgacheffe G1)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_YIRGA_G1_NATURAL', 'Grade-1 Natural Sun-Dried Yirgacheffe G1', 'Grade-1 Natural Sun-Dried Yirgacheffe G1', 'አንደኛ ደረጃ በተፈጥሮ ፀሐይ የደረቀ ይርጋጨፌ G1', 'Buna Yirgaacafee G1 Aduun Gowfame', 'ቀዳማይ ብርኪ ብፀሓይ ዝነቐጸ ይርጋጨፌ G1', 'Bed-dried natural Grade-1 Yirgacheffe micro-lot with jasmine and bergamot tasting notes', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SPECIFICATION', 'grade-1-sun-dried-yirgacheffe-g1', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 5: VARIANT (Green Unroasted Export Beans 60kg Jute Bag with GrainPro)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_GREEN_COFFEE_60KG', 'Green Unroasted Export Beans 60kg Jute Bag with GrainPro', 'Green Unroasted Export Beans 60kg Jute Bag with GrainPro', 'ጥሬ ቡና 60 ኪ.ግ የጅውት ከረጢት በግሬይንፕሮ', 'Buna Dheedhii Export 60kg Saamunaa GrainPro Wajjin', 'ጥረ ሰደድ ቡን 60 ኪ.ግ ጁት ምስ ግሬይንፕሮ', 'Moisture-sealed 11.5% water content specialty green export beans in hermetic GrainPro liner and jute bag', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'VARIANT', 'green-unroasted-60kg-jute-grainpro', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 6: LEAF (Yirgacheffe Kochere Grade-1 Micro-Lot Q-Score 88+ Specialty 60kg)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_KOCHERE_G1_QSCORE88', 'Yirgacheffe Kochere Grade-1 Micro-Lot Q-Score 88+ Specialty 60kg', 'Yirgacheffe Kochere Grade-1 Micro-Lot Q-Score 88+ Specialty 60kg', 'ይርጋጨፌ ኮቸሬ G1 ልዩ ማይክሮ-ሎት Q-ነጥብ 88+ 60 ኪ.ግ', 'Yirgaacafee Kochere G1 Micro-Lot Q-Score 88+ 60kg', 'ይርጋጨፌ ኮቸሬ ቀዳማይ ብርኪ ማይክሮ-ሎት Q-Score 88+ 60 ኪ.ግ', 'ECX / SCA cupping score 88+ certified single-washing station Kochere Grade-1 green export parcel', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'LEAF', 'yirgacheffe-kochere-g1-qscore-88-plus-60kg', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Chain: electronics under electronics
    SELECT id INTO v_scheme_id FROM classification_schemes WHERE code = 'SYS_PROD';
    SELECT id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'electronics' LIMIT 1;
    v_parent_pid := v_root_pid;

    -- Layer 1: SUPERFAMILY (Smartphones & Mobile Devices)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ELEC_SMARTPHONES', 'Smartphones & Mobile Devices', 'Smartphones & Mobile Devices', 'ስማርት ስልኮችና ተንቀሳቃሽ መሣሪያዎች', 'Ismaartifoonii fi Meeshaalee Moobaayilaa', 'ስማርት ስልክታትን ተንቀሳቐስቲ ናውትን', 'Cellular phones, smartphones, tablets, and mobile hardware', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUPERFAMILY', 'smartphones-mobile-handsets', 1, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 2: CATEGORY (5G Flagship Smartphones)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ELEC_5G_SMARTPHONES', '5G Flagship Smartphones', '5G Flagship Smartphones', 'ባለ 5ጂ ከፍተኛ ደረጃ ስማርት ስልኮች', 'Ismaartifoonii 5G Sadarkaa Olaanaa', '5ጂ ፍላግሺፕ ስማርት ስልክታት', 'High performance 5G NR Android and iOS flagship smartphones', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'CATEGORY', '5g-flagship-smartphones', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 3: SUBCATEGORY (Dual SIM Global AMOLED 120Hz Smartphones)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ELEC_DUAL_SIM_AMOLED', 'Dual SIM Global AMOLED 120Hz Smartphones', 'Dual SIM Global AMOLED 120Hz Smartphones', 'ባለ ሁለት ሲም ግሎባል አሞሌድ 120Hz ስልኮች', 'Dual SIM AMOLED 120Hz Ismaartifoonii', 'ክልተ ሲም ኣሞሌድ 120Hz ስልክታት', 'Unlocked dual nano-SIM smartphones featuring 120Hz refresh dynamic AMOLED displays', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUBCATEGORY', 'dual-sim-amoled-120hz-phones', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 4: SPECIFICATION (256GB Storage with 8GB-12GB LPDDR5 RAM)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ELEC_256GB_12GB_RAM', '256GB Storage with 8GB-12GB LPDDR5 RAM', '256GB Storage with 8GB-12GB LPDDR5 RAM', '256 ጊጋባይት ማከማቻ ከ8-12 ጊባ ራም ጋር', 'Kuusaa 256GB fi RAM 8-12GB', '256GB መኽዘንን 8-12GB ራምን', 'Hardware specification with 256GB UFS 4.0 flash storage and 12GB high-speed memory', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SPECIFICATION', '256gb-storage-12gb-ram-phones', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 5: VARIANT (Fast Charging 67W-120W with 5000mAh Battery & OIS Camera)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ELEC_5000MAH_OIS_FASTCHARGE', 'Fast Charging 67W-120W with 5000mAh Battery & OIS Camera', 'Fast Charging 67W-120W with 5000mAh Battery & OIS Camera', 'ፈጣን ቻርጅ 67W-120W ከ5000mAh ባትሪና OIS ካሜራ ጋር', 'Chaarjiingii Saffisaa 67W-120W fi 5000mAh', 'ቅልጡፍ ቻርጅ 67W-120W ምስ 5000mAh ባትሪን OIS ካሜራን', 'Long endurance battery variant with optical image stabilized 200MP camera and ultra-fast charging', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'VARIANT', 'fast-charge-67w-5000mah-ois-camera', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 6: LEAF (Samsung Galaxy S24 Ultra 512GB Titanium Black PTA/Ethio-Telecom Verified)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ELEC_SAMSUNG_S24_ULTRA_512', 'Samsung Galaxy S24 Ultra 512GB Titanium Black PTA/Ethio-Telecom Verified', 'Samsung Galaxy S24 Ultra 512GB Titanium Black PTA/Ethio-Telecom Verified', 'ሳምሰንግ ጋላክሲ S24 አልትራ 512GB ታይታኒየም ኢትዮ-ቴሌኮም የተረጋገጠ', 'Samsung Galaxy S24 Ultra 512GB Itiyoo-Telekoom Kan Mirkanaa''e', 'ሳምሰንግ ጋላክሲ S24 ኣልትራ 512GB ታይታኒየም ኢትዮ-ቴሌኮም ዝተረጋገጸ', 'Snapdragon 8 Gen 3 flagship, Titanium armor chassis, 200MP Quad Telephoto, S-Pen included, IMEI registered in Ethiopia', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'LEAF', 'samsung-galaxy-s24-ultra-512gb-titanium', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Chain: transport-mobility under transport-mobility
    SELECT id INTO v_scheme_id FROM classification_schemes WHERE code = 'SYS_PROD';
    SELECT id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'transport-mobility' LIMIT 1;
    v_parent_pid := v_root_pid;

    -- Layer 1: SUPERFAMILY (Commercial Passenger Mobility & Small Vehicles)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_COMMERCIAL_MOBILITY', 'Commercial Passenger Mobility & Small Vehicles', 'Commercial Passenger Mobility & Small Vehicles', 'የመንገደኞችና አነስተኛ የንግድ ተሽከርካሪዎች', 'Konkolaattota Daldalaa fi Imaltootaa', 'ናይ ተሳፈርትን ንግዳውን መካይን', 'Intra-city passenger transport, minibuses, and three-wheeled taxis', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUPERFAMILY', 'commercial-passenger-mobility', 1, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 2: CATEGORY (Three-Wheelers & Auto-Rickshaws (Bajaj))
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_BAJAJ_THREE_WHEEL', 'Three-Wheelers & Auto-Rickshaws (Bajaj)', 'Three-Wheelers & Auto-Rickshaws (Bajaj)', 'ባለሦስት እግር ተሽከርካሪዎች (ባጃጅ)', 'Konkolaattota Miila Sadii (Baajaajii)', 'ሰለስተ እግሪ ዘለወን ተሽከርከርቲ (ባጃጅ)', 'Urban three-wheeled auto-rickshaws, passenger tricycles, and cargo bajajs', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'CATEGORY', 'three-wheelers-bajaj', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 3: SUBCATEGORY (4-Stroke Petrol Passenger Auto-Rickshaws)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_4STROKE_PETROL_BAJAJ', '4-Stroke Petrol Passenger Auto-Rickshaws', '4-Stroke Petrol Passenger Auto-Rickshaws', 'ባለ 4-ስትሮክ ነዳጅ የመንገደኛ ባጃጆች', 'Baajaajii Imaltootaa 4-Stroke Pirootoolii', '4-ስትሮክ ነዳዲ ተሳፈርቲ ባጃጅ', 'Fuel efficient single cylinder 4-stroke air-cooled petrol passenger transport autorickshaws', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUBCATEGORY', '4-stroke-petrol-passenger-rickshaws', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 4: SPECIFICATION (200cc-225cc Twin-Spark DTS-i Commercial Spec)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_200CC_DTSI_SPEC', '200cc-225cc Twin-Spark DTS-i Commercial Spec', '200cc-225cc Twin-Spark DTS-i Commercial Spec', '200cc-225cc መንታ ሻማ DTS-i የንግድ ዝርዝር', '200cc-225cc DTS-i Ispeesifikeshinii', '200cc-225cc ክልተ ሻማ DTS-i ንግዳዊ መግለጺ', 'High torque digital twin spark ignition engine configuration optimized for steep urban gradients', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SPECIFICATION', '200cc-225cc-twin-spark-dtsi-spec', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 5: VARIANT (Reinforced Chassis with Canvas Roof & Reverse Gear)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_REINFORCED_CHASSIS_CANVAS', 'Reinforced Chassis with Canvas Roof & Reverse Gear', 'Reinforced Chassis with Canvas Roof & Reverse Gear', 'የተጠናከረ ቻሲስ ከሸራ ጣሪያና የኋሊት ማርሽ ጋር', 'Chaasii Jabaa fi Garagalchaa', 'ዝተረረ ቻሲስ ምስ ናይ ሸራ ናሕስን ድሕሪት ማርሽን', 'Heavy-duty shock absorbers, hydraulic front drum brakes, high-tensile steel frame with all-weather tarpaulin', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'VARIANT', 'reinforced-chassis-canvas-roof-reverse', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 6: LEAF (Bajaj RE 205cc DTS-i Passenger 4-Stroke Complete with Tarpaulin & Spare Wheel)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_BAJAJ_RE_205CC_COMPLETE', 'Bajaj RE 205cc DTS-i Passenger 4-Stroke Complete with Tarpaulin & Spare Wheel', 'Bajaj RE 205cc DTS-i Passenger 4-Stroke Complete with Tarpaulin & Spare Wheel', 'ባጃጅ RE 205cc DTS-i ባለ 4-ስትሮክ ሙሉ መለዋወጫና መሸፈኛ ያለው', 'Baajaajii RE 205cc DTS-i Guutuu Tarpaulin fi Spare Wajjin', 'ባጃጅ RE 205cc DTS-i ምሉእ መለዋወጥን መሸፈንን ዘለዎ', 'Genuine Bajaj RE Compact 205cc DTS-i model, factory assembled, complete customs clearing papers, 4+1 passenger capacity', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'LEAF', 'bajaj-re-205cc-dtsi-passenger-complete', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Chain: energy-utilities under energy-utilities
    SELECT id INTO v_scheme_id FROM classification_schemes WHERE code = 'SYS_PROD';
    SELECT id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'energy-utilities' LIMIT 1;
    v_parent_pid := v_root_pid;

    -- Layer 1: SUPERFAMILY (Solar & Renewable Energy Systems)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ENG_SOLAR_RENEWABLE', 'Solar & Renewable Energy Systems', 'Solar & Renewable Energy Systems', 'የፀሐይና ታዳሽ ኃይል ሥርዓቶች', 'Sirna Annisaa Aduu fi Haaromfamu', 'ስርዓታት ፀሓይን ተሓዳሲ ሓይልን', 'Photovoltaic panels, solar charge controllers, battery banks, and inverters', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUPERFAMILY', 'solar-renewable-energy-systems', 1, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 2: CATEGORY (Off-Grid & Hybrid Solar Inverters)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ENG_HYBRID_INVERTERS', 'Off-Grid & Hybrid Solar Inverters', 'Off-Grid & Hybrid Solar Inverters', 'ኦፍ-ግሪድና ሃይብሪድ የፀሐይ ኢንቨርተሮች', 'Inveertara Aduu Off-Grid fi Haayibiriidii', 'ኦፍ-ግሪድን ሃይብሪድን ናይ ፀሓይ ኢንቨርተራት', 'Multifunction inverter-chargers combining inverter, solar charger, and utility battery charger', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'CATEGORY', 'off-grid-hybrid-solar-inverters', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 3: SUBCATEGORY (Pure Sine Wave MPPT Solar Hybrid Inverters)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ENG_PURE_SINE_MPPT', 'Pure Sine Wave MPPT Solar Hybrid Inverters', 'Pure Sine Wave MPPT Solar Hybrid Inverters', 'ንፁህ ሳይን ዌቭ MPPT የፀሐይ ሃይብሪድ ኢንቨርተሮች', 'Pure Sine Wave MPPT Inveertara Haayibiriidii', 'ንጹህ ሳይን ዌቭ MPPT ናይ ፀሓይ ሃይብሪድ ኢንቨርተር', 'High efficiency pure sine wave power inverters with built-in Maximum Power Point Tracking (MPPT)', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUBCATEGORY', 'pure-sine-wave-mppt-inverters', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 4: SPECIFICATION (48V 5kW-6kW Parallel Capability with Smart BMS Interface)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ENG_48V_5KW_SMART_BMS', '48V 5kW-6kW Parallel Capability with Smart BMS Interface', '48V 5kW-6kW Parallel Capability with Smart BMS Interface', '48V 5kW-6kW ከስማርት ቢኤምኤስ ግንኙነት ጋር', '48V 5kW-6kW Smart BMS Wajjin', '48V 5kW-6kW ምስ ስማርት ቢኤምኤስ', '48VDC system architecture with parallel operation up to 9 units and CAN/RS485 lithium communication', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SPECIFICATION', '48v-5kw-6kw-parallel-smart-bms', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 5: VARIANT (Lithium LiFePO4 Compatible High-Voltage PV 450V Input)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ENG_LIFEPO4_HIGH_PV_450V', 'Lithium LiFePO4 Compatible High-Voltage PV 450V Input', 'Lithium LiFePO4 Compatible High-Voltage PV 450V Input', 'ከሊቲየም LiFePO4 ባትሪ ጋር የሚሰራ 450V PV ግብዓት', 'LiFePO4 Liitiyemii fi PV 450V', 'ምስ ሊቲየም LiFePO4 ዝሰርሕ 450V ፒቪ ግብኣት', 'Wide 120V-450VDC PV input window allowing battery-less solar operation during daylight hours', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'VARIANT', 'lifepo4-compatible-high-pv-450v', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 6: LEAF (Growatt SPF 5000ES 5kW 48V Pure Sine Wave Hybrid Solar Inverter with 100A MPPT)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_ENG_GROWATT_5000ES_5KW', 'Growatt SPF 5000ES 5kW 48V Pure Sine Wave Hybrid Solar Inverter with 100A MPPT', 'Growatt SPF 5000ES 5kW 48V Pure Sine Wave Hybrid Solar Inverter with 100A MPPT', 'ግሮዋት SPF 5000ES 5kW 48V ሃይብሪድ የፀሐይ ኢንቨርተር ከ100A MPPT ጋር', 'Growatt SPF 5000ES 5kW Inveertara Aduu 100A MPPT', 'ግሮዋት SPF 5000ES 5kW 48V ሃይብሪድ ናይ ፀሓይ ኢንቨርተር ምስ 100A MPPT', 'Growatt 5000W 230VAC off-grid hybrid power inverter, integrated 100A MPPT charge controller, WiFi monitoring module', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'LEAF', 'growatt-spf-5000es-5kw-48v-hybrid-inverter', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Chain: construction-infrastructure under construction-infrastructure
    SELECT id INTO v_scheme_id FROM classification_schemes WHERE code = 'SYS_PROD';
    SELECT id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'construction-infrastructure' LIMIT 1;
    v_parent_pid := v_root_pid;

    -- Layer 1: SUPERFAMILY (Building Materials & Construction Supplies)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_CONST_BUILDING_MATERIALS', 'Building Materials & Construction Supplies', 'Building Materials & Construction Supplies', 'የግንባታ ዕቃዎችና ቁሳቁሶች', 'Meeshaalee Ijaarsaa fi Dhiheessii', 'ናይ ህንጸት ኣቑሑትን ቀረብን', 'Cement, rebar steel, aggregates, roofing sheets, timber, and masonry binders', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUPERFAMILY', 'building-materials-supplies', 1, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 2: CATEGORY (Hydraulic Cement & Binding Materials)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_CONST_HYDRAULIC_CEMENT', 'Hydraulic Cement & Binding Materials', 'Hydraulic Cement & Binding Materials', 'ሲሚንቶና ማያያዣ የግንባታ ቁሳቁሶች', 'Simintoo fi Meeshaalee Qabsiisan', 'ሲመንቶን መጣበቒ ናውቲ ህንጸትን', 'Portland cements, pozzolanic binders, lime, and concrete admixtures', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'CATEGORY', 'hydraulic-cement-binding-materials', 2, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 3: SUBCATEGORY (Portland Pozzolana & Ordinary Portland Cement (PPC & OPC))
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_CONST_PORTLAND_CEMENT', 'Portland Pozzolana & Ordinary Portland Cement (PPC & OPC)', 'Portland Pozzolana & Ordinary Portland Cement (PPC & OPC)', 'ፖርትላንድ ፖዞላና እና ኦርዲነሪ ሲሚንቶ (PPC እና OPC)', 'Simintoo Pooritlaandii (PPC fi OPC)', 'ፖርትላንድ ፖዞላናን ኦርዲነሪ ሲመንቶን (PPCን OPCን)', 'Standard building cements manufactured according to Ethiopian Quality and Standards Authority (QSAE) norms', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SUBCATEGORY', 'portland-pozzolana-ordinary-cement', 3, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 4: SPECIFICATION (High Early Strength OPC 42.5N Structural Grade)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_CONST_OPC_425N_STRUCTURAL', 'High Early Strength OPC 42.5N Structural Grade', 'High Early Strength OPC 42.5N Structural Grade', 'ከፍተኛ ፈጣን ጥንካሬ ያለው OPC 42.5N ለመዋቅራዊ ግንባታ', 'OPC 42.5N Ijaarsa Caasaatiif', 'ልዑል ቅልጡፍ ጥንካረ ዘለዎ OPC 42.5N ንመዋቕራዊ ህንጸት', '42.5 MPa compressive strength ordinary portland cement designed for multi-story load bearing columns and beams', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'SPECIFICATION', 'opc-425n-high-strength-structural-grade', 4, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 5: VARIANT (Moisture-Resistant 50kg Paper-Laminated Bag)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_CONST_50KG_MOISTURE_RESISTANT', 'Moisture-Resistant 50kg Paper-Laminated Bag', 'Moisture-Resistant 50kg Paper-Laminated Bag', 'እርጥበት መቋቋም የሚችል 50 ኪ.ግ የታሸገ ሲሚንቶ', 'Simintoo 50kg Jiidhinsa Danda''u', 'ርስሓት ዝፃወር 50 ኪ.ግ ዝተዓሸገ ሲመንቶ', '3-ply kraft paper bag with inner polythene moisture barrier film packed 50kg per unit', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'VARIANT', '50kg-bag-moisture-proof-laminated', 5, false, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

    -- Layer 6: LEAF (Dangote OPC 42.5N Grade Heavy Structural Cement 50kg Bag Certified QSAE)
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_CONST_DANGOTE_OPC_425N_50KG', 'Dangote OPC 42.5N Grade Heavy Structural Cement 50kg Bag Certified QSAE', 'Dangote OPC 42.5N Grade Heavy Structural Cement 50kg Bag Certified QSAE', 'ዳንጎቴ OPC 42.5N የኢትዮጵያ ጥራት ደረጃ የተረጋገጠ 50 ኪ.ግ', 'Dangote OPC 42.5N Simintoo 50kg QSAE Kan Mirkanaa''e', 'ዳንጎቴ OPC 42.5N ብቑዕ ህንጸት ሲመንቶ 50 ኪ.ግ ብደረጃ QSAE', 'Dangote Cement Ethiopia Mugher plant Grade 42.5N Portland Cement, tested 28-day compressive strength > 45MPa', 'PRODUCT', 'ACTIVE')
    ON CONFLICT (code) DO UPDATE SET
        canonical_name = EXCLUDED.canonical_name,
        name_en = EXCLUDED.name_en,
        name_am = EXCLUDED.name_am,
        name_om = EXCLUDED.name_om,
        name_ti = EXCLUDED.name_ti,
        description = EXCLUDED.description
    RETURNING id INTO v_curr_cid;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_curr_cid, v_parent_pid, 'LEAF', 'dangote-opc-425n-structural-cement-50kg', 6, true, true, 'PUBLISHED')
    ON CONFLICT (scheme_id, slug) DO UPDATE SET
        concept_id = EXCLUDED.concept_id,
        parent_placement_id = EXCLUDED.parent_placement_id,
        role = EXCLUDED.role,
        depth = EXCLUDED.depth,
        is_leaf = EXCLUDED.is_leaf
    RETURNING id INTO v_curr_pid;
    v_parent_pid := v_curr_pid;

END $$;

-- 4. Populate Dedicated Vertical Schemes (SYS_TOOLS, SYS_ELEC, SYS_AGRI, SYS_VEH)
DO $$
DECLARE
    v_tools_id INT; v_elec_id INT; v_agri_id INT; v_veh_id INT;
    v_prod_id  INT;
    v_root_cid INT;
    v_t_root_pid INT; v_e_root_pid INT; v_a_root_pid INT; v_v_root_pid INT;
BEGIN
    SELECT id INTO v_prod_id  FROM classification_schemes WHERE code = 'SYS_PROD';
    SELECT id INTO v_tools_id FROM classification_schemes WHERE code = 'SYS_TOOLS';
    SELECT id INTO v_elec_id  FROM classification_schemes WHERE code = 'SYS_ELEC';
    SELECT id INTO v_agri_id  FROM classification_schemes WHERE code = 'SYS_AGRI';
    SELECT id INTO v_veh_id   FROM classification_schemes WHERE code = 'SYS_VEH';

    -- 4.1 SYS_TOOLS Root & Mirror Placements
    SELECT concept_id INTO v_root_cid FROM classification_placements WHERE scheme_id = v_prod_id AND slug = 'industrial-manufacturing' AND role = 'ROOT' LIMIT 1;
    IF v_root_cid IS NOT NULL THEN
        INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
        VALUES (v_tools_id, v_root_cid, NULL, 'ROOT', 'tools-machinery-root', 0, false, true, 'PUBLISHED')
        ON CONFLICT (scheme_id, slug) DO UPDATE SET depth = 0 RETURNING id INTO v_t_root_pid;
    END IF;

    -- Mirror Power Tools & Hardware under SYS_TOOLS
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    SELECT v_tools_id, cp.concept_id, v_t_root_pid, cp.role, 'tools-' || cp.slug, cp.depth, cp.is_leaf, true, 'PUBLISHED'
    FROM classification_placements cp
    WHERE cp.scheme_id = v_prod_id AND (cp.slug LIKE '%tool%' OR cp.slug LIKE '%drill%' OR cp.slug LIKE '%hardware%')
      AND cp.depth > 0
    ON CONFLICT (scheme_id, slug) DO NOTHING;

    -- 4.2 SYS_ELEC Root & Mirror Placements
    SELECT concept_id INTO v_root_cid FROM classification_placements WHERE scheme_id = v_prod_id AND slug = 'electronics' AND role = 'ROOT' LIMIT 1;
    IF v_root_cid IS NOT NULL THEN
        INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
        VALUES (v_elec_id, v_root_cid, NULL, 'ROOT', 'electronics-root', 0, false, true, 'PUBLISHED')
        ON CONFLICT (scheme_id, slug) DO UPDATE SET depth = 0 RETURNING id INTO v_e_root_pid;
    END IF;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    SELECT v_elec_id, cp.concept_id, v_e_root_pid, cp.role, 'elec-' || cp.slug, cp.depth, cp.is_leaf, true, 'PUBLISHED'
    FROM classification_placements cp
    WHERE cp.scheme_id = v_prod_id AND (cp.slug LIKE '%elec%' OR cp.slug LIKE '%phone%' OR cp.slug LIKE '%camera%' OR cp.slug LIKE '%computer%')
      AND cp.depth > 0
    ON CONFLICT (scheme_id, slug) DO NOTHING;

    -- 4.3 SYS_AGRI Root & Mirror Placements
    SELECT concept_id INTO v_root_cid FROM classification_placements WHERE scheme_id = v_prod_id AND slug = 'agriculture' AND role = 'ROOT' LIMIT 1;
    IF v_root_cid IS NOT NULL THEN
        INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
        VALUES (v_agri_id, v_root_cid, NULL, 'ROOT', 'agriculture-root', 0, false, true, 'PUBLISHED')
        ON CONFLICT (scheme_id, slug) DO UPDATE SET depth = 0 RETURNING id INTO v_a_root_pid;
    END IF;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    SELECT v_agri_id, cp.concept_id, v_a_root_pid, cp.role, 'agri-' || cp.slug, cp.depth, cp.is_leaf, true, 'PUBLISHED'
    FROM classification_placements cp
    WHERE cp.scheme_id = v_prod_id AND (cp.slug LIKE '%crop%' OR cp.slug LIKE '%teff%' OR cp.slug LIKE '%grain%' OR cp.slug LIKE '%coffee%' OR cp.slug LIKE '%farm%')
      AND cp.depth > 0
    ON CONFLICT (scheme_id, slug) DO NOTHING;

    -- 4.4 SYS_VEH Root & Mirror Placements
    SELECT concept_id INTO v_root_cid FROM classification_placements WHERE scheme_id = v_prod_id AND slug = 'transport-mobility' AND role = 'ROOT' LIMIT 1;
    IF v_root_cid IS NOT NULL THEN
        INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
        VALUES (v_veh_id, v_root_cid, NULL, 'ROOT', 'transport-mobility-root', 0, false, true, 'PUBLISHED')
        ON CONFLICT (scheme_id, slug) DO UPDATE SET depth = 0 RETURNING id INTO v_v_root_pid;
    END IF;

    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    SELECT v_veh_id, cp.concept_id, v_v_root_pid, cp.role, 'veh-' || cp.slug, cp.depth, cp.is_leaf, true, 'PUBLISHED'
    FROM classification_placements cp
    WHERE cp.scheme_id = v_prod_id AND (cp.slug LIKE '%veh%' OR cp.slug LIKE '%car%' OR cp.slug LIKE '%truck%' OR cp.slug LIKE '%bajaj%' OR cp.slug LIKE '%wheel%')
      AND cp.depth > 0
    ON CONFLICT (scheme_id, slug) DO NOTHING;
END $$;

-- 5. Recalculate tree depth and canonical semantic roles across all placements
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
UPDATE classification_placements SET role = 'LEAF' WHERE depth = 6 AND role <> 'LEAF';

-- 6. Recalculate is_leaf flags across all placements efficiently
UPDATE classification_placements p
SET is_leaf = false
WHERE is_leaf = true
  AND EXISTS (SELECT 1 FROM classification_placements child WHERE child.parent_placement_id = p.id);

UPDATE classification_placements p
SET is_leaf = true
WHERE is_leaf = false
  AND NOT EXISTS (SELECT 1 FROM classification_placements child WHERE child.parent_placement_id = p.id);

COMMIT;
