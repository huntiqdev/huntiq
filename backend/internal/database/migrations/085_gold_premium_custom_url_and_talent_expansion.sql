-- ============================================================================
-- 085_gold_premium_custom_url_and_talent_expansion.sql
-- Massive Multi-Layer Leaf Expansion, Non-Materialistic Talents & Gold/Premium Benefits
-- 
-- Features:
--   1. Add custom_url & social_links to products and users (Tier privilege for Gold/Premium)
--   2. Non-Materialistic Domains (Talents, CVs, Content Creators, Sports, Performing Arts)
--   3. Poly-Hierarchy Multi-Scheme Placements (Drill in Tools & Construction; Bajaj in Vehicles & Transport)
--   4. Deep Multi-Branch Expansion across Layers 3, 4, 5, and 6 with 4-Language Naming
--   5. Fully synchronized across classification_concepts, classification_placements, and category_hierarchy
-- ============================================================================

BEGIN;

-- 1. Add custom_url and social_links to products and users
ALTER TABLE products ADD COLUMN IF NOT EXISTS custom_url VARCHAR(500);
ALTER TABLE products ADD COLUMN IF NOT EXISTS social_links JSONB NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE users ADD COLUMN IF NOT EXISTS custom_url VARCHAR(500);
ALTER TABLE users ADD COLUMN IF NOT EXISTS social_links JSONB NOT NULL DEFAULT '{}'::jsonb;

CREATE INDEX IF NOT EXISTS idx_products_custom_url ON products(custom_url) WHERE custom_url IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_custom_url ON users(custom_url) WHERE custom_url IS NOT NULL;

-- 2. Multi-Branch & Talent Taxonomy Expansion
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

    -- --- Node: digital-creators-influencers (Depth 1, Role SUPERFAMILY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_WORK_DIGITAL_CREATORS', 'Digital Content Creators & Social Influencers', 'Digital Content Creators & Social Influencers', 'የዲጂታል ይዘት አዘጋጆችና የማህበራዊ ሚዲያ ተፅዕኖ ፈጣሪዎች', 'Uumtoota Qabiyyee Dijitaalaa fi Dhiibbaa Hawaasaa', 'ዲጂታል ትሕዝቶ ኣዳለውትን ናይ ማሕበራዊ ሚድያ ጽልዋ ፈጠርትን', 'Professional video content creators, comedy influencers, digital brand ambassadors, and podcasters', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'workforce-employment' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUPERFAMILY', 'digital-creators-influencers', 1, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'workforce-employment' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Digital Content Creators & Social Influencers', 'የዲጂታል ይዘት አዘጋጆችና የማህበራዊ ሚዲያ ተፅዕኖ ፈጣሪዎች', 'Uumtoota Qabiyyee Dijitaalaa fi Dhiibbaa Hawaasaa', 'ዲጂታል ትሕዝቶ ኣዳለውትን ናይ ማሕበራዊ ሚድያ ጽልዋ ፈጠርትን', 'digital-creators-influencers', 'Professional video content creators, comedy influencers, digital brand ambassadors, and podcasters', 1, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: sports-athletes-fitness-talents (Depth 1, Role SUPERFAMILY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SPORT_TALENTS_ATHLETES', 'Sports Athletes, Coaches & Fitness Talents', 'Sports Athletes, Coaches & Fitness Talents', 'ስፖርተኞች፣ አሰልጣኞችና የአካል ብቃት ባለሙያዎች', 'Atileetota Ispoortii, Leenjistoota fi Fayyummaa', 'ስፖርተኛታት፣ ኣሰልጠንትን ናይ ኣካል ብቕዓት ሞያተኛታትን', 'Athletic competitors, football talents, personal fitness coaches, and distance runners for hire and scouting', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'workforce-employment' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUPERFAMILY', 'sports-athletes-fitness-talents', 1, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'workforce-employment' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Sports Athletes, Coaches & Fitness Talents', 'ስፖርተኞች፣ አሰልጣኞችና የአካል ብቃት ባለሙያዎች', 'Atileetota Ispoortii, Leenjistoota fi Fayyummaa', 'ስፖርተኛታት፣ ኣሰልጠንትን ናይ ኣካል ብቕዓት ሞያተኛታትን', 'sports-athletes-fitness-talents', 'Athletic competitors, football talents, personal fitness coaches, and distance runners for hire and scouting', 1, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: social-video-content-creators (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_WORK_VIDEO_CREATORS', 'Short-Form Video & Comedy Content Creators', 'Short-Form Video & Comedy Content Creators', 'የአጭር ቪዲዮና አስቂኝ ይዘት ፈጣሪዎች', 'Uumtoota Viidiyoo Gabaabaa fi Qoosaa', 'ናይ ሓጺር ቪድዮን ቀልድን ትሕዝቶ ፈጠርቲ', 'TikTok, YouTube Shorts, and Instagram Reels creative producers specializing in comedy and viral storytelling', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'digital-creators-influencers' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'social-video-content-creators', 2, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'digital-creators-influencers' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Short-Form Video & Comedy Content Creators', 'የአጭር ቪዲዮና አስቂኝ ይዘት ፈጣሪዎች', 'Uumtoota Viidiyoo Gabaabaa fi Qoosaa', 'ናይ ሓጺር ቪድዮን ቀልድን ትሕዝቶ ፈጠርቲ', 'social-video-content-creators', 'TikTok, YouTube Shorts, and Instagram Reels creative producers specializing in comedy and viral storytelling', 2, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: football-soccer-players-scouting (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SPORT_FOOTBALL_SCOUTING', 'Football / Soccer Talents & Player Scouting CVs', 'Football / Soccer Talents & Player Scouting CVs', 'የእግር ኳስ ተጫዋቾች ተሰጥኦና የክለብ ምልመላ ሲቪ', 'Dandeettii Kubbaa Miilaa fi CV Taphattootaa', 'ናይ ኩዕሶ እግሪ ተጻወትቲ ተውህቦን ናይ ክለብ ምልመላ ሲቪን', 'Amateur and semi-professional football players seeking club trials, academy scouting, and seasonal contracts', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'sports-athletes-fitness-talents' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'football-soccer-players-scouting', 2, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'sports-athletes-fitness-talents' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Football / Soccer Talents & Player Scouting CVs', 'የእግር ኳስ ተጫዋቾች ተሰጥኦና የክለብ ምልመላ ሲቪ', 'Dandeettii Kubbaa Miilaa fi CV Taphattootaa', 'ናይ ኩዕሶ እግሪ ተጻወትቲ ተውህቦን ናይ ክለብ ምልመላ ሲቪን', 'football-soccer-players-scouting', 'Amateur and semi-professional football players seeking club trials, academy scouting, and seasonal contracts', 2, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: software-engineering-it-freelance (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_WORK_SOFTWARE_DEV', 'Software Development & Tech Freelancers', 'Software Development & Tech Freelancers', 'የሶፍትዌር አልሚዎችና የቴክኖሎጂ ፍሪላንሰሮች', 'Misoonsa Sooftiweerii fi Hojjattoota Teeknooloojii', 'ናይ ሶፍትዌር ምዕባለን ናይ ቴክኖሎጂ ፍሪላንሰራትን', 'Full-stack web developers, mobile app engineers, and cloud architects for contract and full-time hiring', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'workforce-employment' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'software-engineering-it-freelance', 2, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'workforce-employment' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Software Development & Tech Freelancers', 'የሶፍትዌር አልሚዎችና የቴክኖሎጂ ፍሪላንሰሮች', 'Misoonsa Sooftiweerii fi Hojjattoota Teeknooloojii', 'ናይ ሶፍትዌር ምዕባለን ናይ ቴክኖሎጂ ፍሪላንሰራትን', 'software-engineering-it-freelance', 'Full-stack web developers, mobile app engineers, and cloud architects for contract and full-time hiring', 2, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: drills-and-drivers-hardware (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TOOL_DRILLS_DRIVERS', 'Drills, Impact Drivers & Rotary Hammers', 'Drills, Impact Drivers & Rotary Hammers', 'ድሪሎች፣ ኢምፓክት ድራይቨሮችና ሮታሪ ሀመሮች', 'Boraasuu, Naannessituu fi Dhagaa Waraantoo', 'ድሪላት፣ ኢምፓክት ድራይቨራትን ሮታሪ ሃመራትን', 'Cordless and corded electric rotary drills, hammer drivers, and demolition hammers', 'EQUIPMENT', 'ACTIVE')
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
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'drills-and-drivers-hardware', 2, false, true, 'PUBLISHED')
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
    VALUES ('Drills, Impact Drivers & Rotary Hammers', 'ድሪሎች፣ ኢምፓክት ድራይቨሮችና ሮታሪ ሀመሮች', 'Boraasuu, Naannessituu fi Dhagaa Waraantoo', 'ድሪላት፣ ኢምፓክት ድራይቨራትን ሮታሪ ሃመራትን', 'drills-and-drivers-hardware', 'Cordless and corded electric rotary drills, hammer drivers, and demolition hammers', 2, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: ethiopian-traditional-spices (Depth 2, Role CATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_TRAD_SPICES', 'Ethiopian Traditional Spices & Seasonings', 'Ethiopian Traditional Spices & Seasonings', 'የኢትዮጵያ ባህላዊ ቅመማ ቅመሞችና ማጣፈጫዎች', 'Mi''eessituu fi Qoricha Nyaataa Aadaa Itoophiyaa', 'ናይ ኢትዮጵያ ባህላዊ ቅመማ ቅመማትን መመቀርታትን', 'Authentic organic Ethiopian spices: Berbere, Korarima, Mitmita, Ginger, Turmeric, Black Cumin (Tikur Azmud)', 'PRODUCT', 'ACTIVE')
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
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'CATEGORY', 'ethiopian-traditional-spices', 2, false, true, 'PUBLISHED')
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
    VALUES ('Ethiopian Traditional Spices & Seasonings', 'የኢትዮጵያ ባህላዊ ቅመማ ቅመሞችና ማጣፈጫዎች', 'Mi''eessituu fi Qoricha Nyaataa Aadaa Itoophiyaa', 'ናይ ኢትዮጵያ ባህላዊ ቅመማ ቅመማትን መመቀርታትን', 'ethiopian-traditional-spices', 'Authentic organic Ethiopian spices: Berbere, Korarima, Mitmita, Ginger, Turmeric, Black Cumin (Tikur Azmud)', 2, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: tiktok-comedy-skit-creators (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_WORK_TIKTOK_COMEDY', 'TikTok & Instagram Comedy Skit Creators', 'TikTok & Instagram Comedy Skit Creators', 'የቲክቶክና ኢንስታግራም አስቂኝ ድራማ ፈጣሪዎች', 'Uumtoota Diraamaa Qoosaa TikTok fi Instagram', 'ናይ ቲክቶክን ኢንስታግራምን ናይ ቀልዲ ድራማ ፈጠርቲ', 'Short-form comedic entertainers and sketch artists available for digital brand sponsorships and ads', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'social-video-content-creators' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'tiktok-comedy-skit-creators', 3, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'social-video-content-creators' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('TikTok & Instagram Comedy Skit Creators', 'የቲክቶክና ኢንስታግራም አስቂኝ ድራማ ፈጣሪዎች', 'Uumtoota Diraamaa Qoosaa TikTok fi Instagram', 'ናይ ቲክቶክን ኢንስታግራምን ናይ ቀልዲ ድራማ ፈጠርቲ', 'tiktok-comedy-skit-creators', 'Short-form comedic entertainers and sketch artists available for digital brand sponsorships and ads', 3, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: attacking-strikers-and-wingers (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SPORT_STRIKERS_WINGERS', 'Attacking Center Forwards & Wingers', 'Attacking Center Forwards & Wingers', 'አጥቂ አጥቂዎችና የመስመር ተጫዋቾች', 'Taphata Fuulduraa fi Sarara Cinaa', 'ኣጥቃዕቲ ተጻወትትን ናይ መስመር ተጻወትትን', 'Fast agile forwards, target strikers, and pace wingers with match statistics and highlight reels', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'football-soccer-players-scouting' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'attacking-strikers-and-wingers', 3, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'football-soccer-players-scouting' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Attacking Center Forwards & Wingers', 'አጥቂ አጥቂዎችና የመስመር ተጫዋቾች', 'Taphata Fuulduraa fi Sarara Cinaa', 'ኣጥቃዕቲ ተጻወትትን ናይ መስመር ተጻወትትን', 'attacking-strikers-and-wingers', 'Fast agile forwards, target strikers, and pace wingers with match statistics and highlight reels', 3, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: fullstack-web-mobile-developers (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_WORK_FULLSTACK_DEV', 'Full-Stack Web & Mobile App Engineers', 'Full-Stack Web & Mobile App Engineers', 'የሙሉ መስመር ድረ-ገፅና የሞባይል አፕሊኬሽን መሃንዲሶች', 'Injinariiwwan Weebsaayitii fi Mobaayilaa Guutuu', 'ናይ ድረ-ገጽን ሞባይል ኣፕሊኬሽንን መሃንድሳት', 'Engineers skilled in React, Next.js, Node.js, Go, Python, Flutter, and PostgreSQL databases', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'software-engineering-it-freelance' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'fullstack-web-mobile-developers', 3, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'software-engineering-it-freelance' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Full-Stack Web & Mobile App Engineers', 'የሙሉ መስመር ድረ-ገፅና የሞባይል አፕሊኬሽን መሃንዲሶች', 'Injinariiwwan Weebsaayitii fi Mobaayilaa Guutuu', 'ናይ ድረ-ገጽን ሞባይል ኣፕሊኬሽንን መሃንድሳት', 'fullstack-web-mobile-developers', 'Engineers skilled in React, Next.js, Node.js, Go, Python, Flutter, and PostgreSQL databases', 3, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: cordless-brushless-impact-drills (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TOOL_CORDLESS_BRUSHLESS', 'Cordless Brushless Lithium Impact Drills', 'Cordless Brushless Lithium Impact Drills', 'ያለ ሽቦ ብሩሽ-አልባ የሊቲየም ኢምፓክት ድሪሎች', 'Boraasuu Shibo-malee Liitiyeemii Cimaa', 'ሽቦ ኣልቦ ብሩሽ ዘይብሉ ናይ ሊቲየም ድሪላት', '20V/18V brushless motor cordless hammer drills with 2.0Ah/4.0Ah lithium-ion battery packs', 'EQUIPMENT', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'drills-and-drivers-hardware' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'cordless-brushless-impact-drills', 3, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'drills-and-drivers-hardware' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Cordless Brushless Lithium Impact Drills', 'ያለ ሽቦ ብሩሽ-አልባ የሊቲየም ኢምፓክት ድሪሎች', 'Boraasuu Shibo-malee Liitiyeemii Cimaa', 'ሽቦ ኣልቦ ብሩሽ ዘይብሉ ናይ ሊቲየም ድሪላት', 'cordless-brushless-impact-drills', '20V/18V brushless motor cordless hammer drills with 2.0Ah/4.0Ah lithium-ion battery packs', 3, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: berbere-chili-blend-farming (Depth 3, Role SUBCATEGORY) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_BERBERE_SPICE', 'Pure Sun-Dried Red Pepper & Gourmet Berbere Blend', 'Pure Sun-Dried Red Pepper & Gourmet Berbere Blend', 'ንፁህ በፀሐይ የደረቀ ቀይ በርበሬና የቅመም በርበሬ ድብልቅ', 'Barbaree Diimaa Aduun Goge fi Mi''eessituu', 'ንጹህ ብፀሓይ ዝነቐጸ ቀይሕ በርበረን ናይ ቅመም ድብልቅን', 'High-altitude organic red pepper pods sun-dried and traditionally blended with Korarima and sacred herbs', 'PRODUCT', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'ethiopian-traditional-spices' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SUBCATEGORY', 'berbere-chili-blend-farming', 3, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'ethiopian-traditional-spices' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Pure Sun-Dried Red Pepper & Gourmet Berbere Blend', 'ንፁህ በፀሐይ የደረቀ ቀይ በርበሬና የቅመም በርበሬ ድብልቅ', 'Barbaree Diimaa Aduun Goge fi Mi''eessituu', 'ንጹህ ብፀሓይ ዝነቐጸ ቀይሕ በርበረን ናይ ቅመም ድብልቅን', 'berbere-chili-blend-farming', 'High-altitude organic red pepper pods sun-dried and traditionally blended with Korarima and sacred herbs', 3, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: verified-creator-brand-promotions (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_WORK_CREATOR_PROMO_SPEC', '100k+ Followers Brand Sponsorship & Product Review Gigs', '100k+ Followers Brand Sponsorship & Product Review Gigs', 'ከ100ሺህ+ ተከታዮች በላይ የንግድ ማስታወቂያና የምርት ገለጻ አገልግሎት', 'Beeksisa Meeshaa fi Dhaabbilee Hordoftoota 100k+', 'ካብ 100ሽሕ+ ተኸተልቲ ዘለዎም ናይ ንግዲ ምልክታን ምልላይን', 'Sponsored promotional video campaigns, humorous product integration, and dedicated product unboxing', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'tiktok-comedy-skit-creators' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'verified-creator-brand-promotions', 4, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'tiktok-comedy-skit-creators' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('100k+ Followers Brand Sponsorship & Product Review Gigs', 'ከ100ሺህ+ ተከታዮች በላይ የንግድ ማስታወቂያና የምርት ገለጻ አገልግሎት', 'Beeksisa Meeshaa fi Dhaabbilee Hordoftoota 100k+', 'ካብ 100ሽሕ+ ተኸተልቲ ዘለዎም ናይ ንግዲ ምልክታን ምልላይን', 'verified-creator-brand-promotions', 'Sponsored promotional video campaigns, humorous product integration, and dedicated product unboxing', 4, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: under-23-licensed-club-strikers (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SPORT_U23_STRIKER_SPEC', 'Under-23 Regional League & Premier League Striker Profiles', 'Under-23 Regional League & Premier League Striker Profiles', 'ከ23 ዓመት በታች የክልልና የፕሪሚየር ሊግ አጥቂ መገለጫዎች', 'Taphattoota Fuulduraa Waggaa 23 Gadii Liigii Naannoo', 'ትሕቲ 23 ዓመት ናይ ዞባን ፕሪምየር ሊግን ኣጥቃዕቲ ፕሮፋይል', 'Registered competitive football strikers with official federation ID, physical metrics, and goals record', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'attacking-strikers-and-wingers' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'under-23-licensed-club-strikers', 4, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'attacking-strikers-and-wingers' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Under-23 Regional League & Premier League Striker Profiles', 'ከ23 ዓመት በታች የክልልና የፕሪሚየር ሊግ አጥቂ መገለጫዎች', 'Taphattoota Fuulduraa Waggaa 23 Gadii Liigii Naannoo', 'ትሕቲ 23 ዓመት ናይ ዞባን ፕሪምየር ሊግን ኣጥቃዕቲ ፕሮፋይል', 'under-23-licensed-club-strikers', 'Registered competitive football strikers with official federation ID, physical metrics, and goals record', 4, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: senior-react-golang-architects (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_WORK_REACT_GOLANG_SPEC', 'Senior Next.js & Golang Cloud Architecture Specialists', 'Senior Next.js & Golang Cloud Architecture Specialists', 'ከፍተኛ የኔክስት.ጄኤስ እና የጎላንግ ክላውድ ሲስተም ባለሙያዎች', 'Ogeeyyii Sadarkaa Olaanaa Next.js fi Golang', 'ላዕለዎት ናይ Next.jsን Golangን ክላውድ ሞያተኛታት', '5+ years experienced engineers capable of delivering high-concurrency microservices, fintech integrations, and portals', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'fullstack-web-mobile-developers' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'senior-react-golang-architects', 4, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'fullstack-web-mobile-developers' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Senior Next.js & Golang Cloud Architecture Specialists', 'ከፍተኛ የኔክስት.ጄኤስ እና የጎላንግ ክላውድ ሲስተም ባለሙያዎች', 'Ogeeyyii Sadarkaa Olaanaa Next.js fi Golang', 'ላዕለዎት ናይ Next.jsን Golangን ክላውድ ሞያተኛታት', 'senior-react-golang-architects', '5+ years experienced engineers capable of delivering high-concurrency microservices, fintech integrations, and portals', 4, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: 20v-heavy-torque-drill-kit (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TOOL_20V_TORQUE_SPEC', '20V Max 60Nm Heavy Torque Cordless Drill Kit', '20V Max 60Nm Heavy Torque Cordless Drill Kit', '20 ቮልት 60 ኒውተን ሜትር ከፍተኛ ጉልበት ያለው ድሪል ኪት', 'Kiitii Boraasuu 20V Humna Cimaa 60Nm', '20 ቮልት 60 ኒውተን ሜተር ላዕለዋይ ጉልበት ዘለዎ ድሪል ኪት', '20V dual speed (0-400 / 0-1500 RPM) hammer drill with 13mm metal keyless chuck and LED worklight', 'EQUIPMENT', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'cordless-brushless-impact-drills' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', '20v-heavy-torque-drill-kit', 4, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'cordless-brushless-impact-drills' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('20V Max 60Nm Heavy Torque Cordless Drill Kit', '20 ቮልት 60 ኒውተን ሜትር ከፍተኛ ጉልበት ያለው ድሪል ኪት', 'Kiitii Boraasuu 20V Humna Cimaa 60Nm', '20 ቮልት 60 ኒውተን ሜተር ላዕለዋይ ጉልበት ዘለዎ ድሪል ኪት', '20v-heavy-torque-drill-kit', '20V dual speed (0-400 / 0-1500 RPM) hammer drill with 13mm metal keyless chuck and LED worklight', 4, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: electric-green-bajaj-three-wheelers (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_ELECTRIC_BAJAJ_SPEC', 'Electric Battery Powered Green Passenger Bajaj', 'Electric Battery Powered Green Passenger Bajaj', 'በኤሌክትሪክ ባትሪ የሚሰራ ለአካባቢ ተስማሚ ባጃጅ', 'Baajaajii Elektiriikiitiin Hojjatu Magariisa', 'ብኤሌክትሪክ ባትሪ ዝሰርሕ ንከባቢ ምቹእ ባጃጅ', 'Zero emission electric passenger three-wheeler with 60V-72V lithium battery pack and regenerative braking', 'EQUIPMENT', 'ACTIVE')
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
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'electric-green-bajaj-three-wheelers', 4, false, true, 'PUBLISHED')
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
    VALUES ('Electric Battery Powered Green Passenger Bajaj', 'በኤሌክትሪክ ባትሪ የሚሰራ ለአካባቢ ተስማሚ ባጃጅ', 'Baajaajii Elektiriikiitiin Hojjatu Magariisa', 'ብኤሌክትሪክ ባትሪ ዝሰርሕ ንከባቢ ምቹእ ባጃጅ', 'electric-green-bajaj-three-wheelers', 'Zero emission electric passenger three-wheeler with 60V-72V lithium battery pack and regenerative braking', 4, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: mareko-fana-red-pepper-grade-1 (Depth 4, Role SPECIFICATION) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_MAREKO_FANA_SPEC', 'Mareko Fana Grade 1 High-Pungency Red Pepper', 'Mareko Fana Grade 1 High-Pungency Red Pepper', 'የማረቆ ፋና አንደኛ ደረጃ ከፍተኛ ቅመም ቀይ በርበሬ', 'Barbaree Maaraqoo Faanaa Sadarkaa 1ffaa', 'ማረቆ ፋና ቀዳማይ ደረጃ ብሉጽ ቀይሕ በርበረ', 'Renowned Mareko Fana cultivar prized for deep ruby-red color, robust aromatic heat, and high essential oil content', 'PRODUCT', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'berbere-chili-blend-farming' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'SPECIFICATION', 'mareko-fana-red-pepper-grade-1', 4, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'berbere-chili-blend-farming' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Mareko Fana Grade 1 High-Pungency Red Pepper', 'የማረቆ ፋና አንደኛ ደረጃ ከፍተኛ ቅመም ቀይ በርበሬ', 'Barbaree Maaraqoo Faanaa Sadarkaa 1ffaa', 'ማረቆ ፋና ቀዳማይ ደረጃ ብሉጽ ቀይሕ በርበረ', 'mareko-fana-red-pepper-grade-1', 'Renowned Mareko Fana cultivar prized for deep ruby-red color, robust aromatic heat, and high essential oil content', 4, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: single-skit-sponsored-video-production (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_WORK_SKIT_PRODUCTION_PACKAGE', 'Custom 60-Second Viral Comedy Skit Sponsorship Package', 'Custom 60-Second Viral Comedy Skit Sponsorship Package', 'የ60 ሰከንድ አስቂኝ የማስታወቂያ ቪዲዮ ዝግጅት ጥቅል', 'Paakeejii Viidiyoo Qoosaa Sekoondii 60', 'ናይ 60 ካልኢት ናይ ቀልዲ መወዓውዒ ቪድዮ ጥቅል', 'Full production including creative script, humorous acting, editing, brand tag, and posting to primary channel', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'verified-creator-brand-promotions' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'single-skit-sponsored-video-production', 5, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'verified-creator-brand-promotions' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Custom 60-Second Viral Comedy Skit Sponsorship Package', 'የ60 ሰከንድ አስቂኝ የማስታወቂያ ቪዲዮ ዝግጅት ጥቅል', 'Paakeejii Viidiyoo Qoosaa Sekoondii 60', 'ናይ 60 ካልኢት ናይ ቀልዲ መወዓውዒ ቪድዮ ጥቅል', 'single-skit-sponsored-video-production', 'Full production including creative script, humorous acting, editing, brand tag, and posting to primary channel', 5, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: verified-scouting-trial-dossier (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SPORT_SCOUTING_DOSSIER', 'Verified Scouting Match Footage & Trial CV Dossier', 'Verified Scouting Match Footage & Trial CV Dossier', 'የተረጋገጠ የጨዋታ ቪዲዮና የተጫዋች የሙከራ ሲቪ ማህደር', 'Viidiyoo Taphaa Qorannoo fi CV Taphatootaa Mirkanaa''e', 'ዝተረጋገጸ ናይ ጸወታ ቪድዮን ናይ ተጻዋታይ ፈተነ ሲቪን', 'Complete professional dossier: 10-minute HD match highlights link, sprint speed benchmark, and transfer availability', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'under-23-licensed-club-strikers' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'verified-scouting-trial-dossier', 5, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'under-23-licensed-club-strikers' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Verified Scouting Match Footage & Trial CV Dossier', 'የተረጋገጠ የጨዋታ ቪዲዮና የተጫዋች የሙከራ ሲቪ ማህደር', 'Viidiyoo Taphaa Qorannoo fi CV Taphatootaa Mirkanaa''e', 'ዝተረጋገጸ ናይ ጸወታ ቪድዮን ናይ ተጻዋታይ ፈተነ ሲቪን', 'verified-scouting-trial-dossier', 'Complete professional dossier: 10-minute HD match highlights link, sprint speed benchmark, and transfer availability', 5, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: turnkey-custom-app-development-service (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_WORK_TURNKEY_APP_DEV', 'End-to-End Custom Web & Mobile App Development Service', 'End-to-End Custom Web & Mobile App Development Service', 'ሙሉ የድረ-ገፅና የሞባይል አፕሊኬሽን የማልማት አገልግሎት ጥቅል', 'Tajaajila Misoonsa Weebsaayitii fi Mobaayilaa Guutuu', 'ሙሉእ ናይ ድረ-ገጽን ሞባይል ኣፕን ምዕባለ ኣገልግሎት ጥቅል', 'Complete lifecycle: UI/UX design, frontend, API backend, PostgreSQL database, Telebirr/Chapa payment gateway integration', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'senior-react-golang-architects' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'turnkey-custom-app-development-service', 5, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'senior-react-golang-architects' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('End-to-End Custom Web & Mobile App Development Service', 'ሙሉ የድረ-ገፅና የሞባይል አፕሊኬሽን የማልማት አገልግሎት ጥቅል', 'Tajaajila Misoonsa Weebsaayitii fi Mobaayilaa Guutuu', 'ሙሉእ ናይ ድረ-ገጽን ሞባይል ኣፕን ምዕባለ ኣገልግሎት ጥቅል', 'turnkey-custom-app-development-service', 'Complete lifecycle: UI/UX design, frontend, API backend, PostgreSQL database, Telebirr/Chapa payment gateway integration', 5, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: ingco-20v-lithium-drill-set-with-accessories (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TOOL_INGCO_20V_SET', 'INGCO 20V Lithium-Ion Impact Drill with 2 Batteries & 47-Piece Accessories', 'INGCO 20V Lithium-Ion Impact Drill with 2 Batteries & 47-Piece Accessories', 'ኢንግኮ 20ቮልት ሊቲየም ድሪል ከነ 2 ባትሪዎችና 47 መለዋወጫዎች', 'Ingco 20V Boraasuu Baatrii 2 fi Meeshaalee 47 Waliin', 'ኢንግኮ 20ቮልት ሊቲየም ድሪል ምስ 2 ባትሪን 47 መተካእታታትን', 'INGCO CIDLI200215 kit: 20V brushless impact drill, 2x 2.0Ah batteries, fast charger, canvas tool bag, bits set', 'EQUIPMENT', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '20v-heavy-torque-drill-kit' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', 'ingco-20v-lithium-drill-set-with-accessories', 5, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '20v-heavy-torque-drill-kit' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('INGCO 20V Lithium-Ion Impact Drill with 2 Batteries & 47-Piece Accessories', 'ኢንግኮ 20ቮልት ሊቲየም ድሪል ከነ 2 ባትሪዎችና 47 መለዋወጫዎች', 'Ingco 20V Boraasuu Baatrii 2 fi Meeshaalee 47 Waliin', 'ኢንግኮ 20ቮልት ሊቲየም ድሪል ምስ 2 ባትሪን 47 መተካእታታትን', 'ingco-20v-lithium-drill-set-with-accessories', 'INGCO CIDLI200215 kit: 20V brushless impact drill, 2x 2.0Ah batteries, fast charger, canvas tool bag, bits set', 5, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: 120km-range-lithium-electric-bajaj (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_120KM_E_BAJAJ', '120km Range Fast-Charging Lithium Electric Bajaj', '120km Range Fast-Charging Lithium Electric Bajaj', 'በአንድ ቻርጅ 120ኪሜ የሚጓዝ ፈጣን ቻርጀር ያለው የኤሌክትሪክ ባጃጅ', 'Baajaajii Elektiriikii Imala 120km Chaarjii Ariifataa', 'ብሓደ ቻርጅ 120ኪሜ ዝጓዓዝ ቅልጡፍ ቻርጀር ዘለዎ ባጃጅ', 'High density LiFePO4 battery, 3.5-hour fast recharge, digital LCD instrument cluster, solar roof panel assist', 'EQUIPMENT', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'electric-green-bajaj-three-wheelers' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', '120km-range-lithium-electric-bajaj', 5, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'electric-green-bajaj-three-wheelers' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('120km Range Fast-Charging Lithium Electric Bajaj', 'በአንድ ቻርጅ 120ኪሜ የሚጓዝ ፈጣን ቻርጀር ያለው የኤሌክትሪክ ባጃጅ', 'Baajaajii Elektiriikii Imala 120km Chaarjii Ariifataa', 'ብሓደ ቻርጅ 120ኪሜ ዝጓዓዝ ቅልጡፍ ቻርጀር ዘለዎ ባጃጅ', '120km-range-lithium-electric-bajaj', 'High density LiFePO4 battery, 3.5-hour fast recharge, digital LCD instrument cluster, solar roof panel assist', 5, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: 25kg-sack-sun-dried-mareko-fana-pods (Depth 5, Role VARIANT) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_MAREKO_25KG_SACK', '25kg Wholesale Burlap Sack Sun-Dried Mareko Fana Pods', '25kg Wholesale Burlap Sack Sun-Dried Mareko Fana Pods', '25 ኪሎ ግራም የጅምላ ጆንያ በፀሐይ የደረቀ የማረቆ ፋና ፍሬ', 'Kootaa 25kg Jumlaa Barbaree Maaraqoo Faanaa', '25 ኪሎ ግራም ጅምላ ጆንያ ብፀሓይ ዝነቐጸ ማረቆ ፋና ፍረ', 'Destemmed, sorted whole dry pods in ventilated burlap bags ready for gourmet spice grinding mills', 'PRODUCT', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'mareko-fana-red-pepper-grade-1' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'VARIANT', '25kg-sack-sun-dried-mareko-fana-pods', 5, false, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'mareko-fana-red-pepper-grade-1' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('25kg Wholesale Burlap Sack Sun-Dried Mareko Fana Pods', '25 ኪሎ ግራም የጅምላ ጆንያ በፀሐይ የደረቀ የማረቆ ፋና ፍሬ', 'Kootaa 25kg Jumlaa Barbaree Maaraqoo Faanaa', '25 ኪሎ ግራም ጅምላ ጆንያ ብፀሓይ ዝነቐጸ ማረቆ ፋና ፍረ', '25kg-sack-sun-dried-mareko-fana-pods', 'Destemmed, sorted whole dry pods in ventilated burlap bags ready for gourmet spice grinding mills', 5, v_ch_parent_id, true, false, 'SYS_PROD')
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

    -- --- Node: ethiopian-viral-comedian-sponsored-skit-package (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_WORK_ETHIO_COMEDIAN_PACKAGE', 'Ethiopian Viral Comedy Skit & Product Endorsement Campaign', 'Ethiopian Viral Comedy Skit & Product Endorsement Campaign', 'የኢትዮጵያ ታዋቂ ኮሜዲያን የማስታወቂያ ቪዲዮና የምርት ድጋፍ ዘመቻ', 'Diraamaa Qoosaa Beekamaa Itoophiyaa fi Beeksisa Oomishaa', 'ናይ ኢትዮጵያ ፍሉጥ ኮሜድያን ናይ መወዓውዒ ቪድዮን ናይ ፍርያት ደገፍን', 'Verified high-reach comedic influencer promotion with multi-platform posting on TikTok, Instagram & Telegram', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'single-skit-sponsored-video-production' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'ethiopian-viral-comedian-sponsored-skit-package', 6, true, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'single-skit-sponsored-video-production' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Ethiopian Viral Comedy Skit & Product Endorsement Campaign', 'የኢትዮጵያ ታዋቂ ኮሜዲያን የማስታወቂያ ቪዲዮና የምርት ድጋፍ ዘመቻ', 'Diraamaa Qoosaa Beekamaa Itoophiyaa fi Beeksisa Oomishaa', 'ናይ ኢትዮጵያ ፍሉጥ ኮሜድያን ናይ መወዓውዒ ቪድዮን ናይ ፍርያት ደገፍን', 'ethiopian-viral-comedian-sponsored-skit-package', 'Verified high-reach comedic influencer promotion with multi-platform posting on TikTok, Instagram & Telegram', 6, v_ch_parent_id, true, true, 'SYS_PROD')
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

    -- --- Node: professional-football-striker-talent-listing (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_SPORT_STRIKER_TALENT_COMPLETE', 'Professional Center Forward Football Striker Talent Listing', 'Professional Center Forward Football Striker Talent Listing', 'የፕሮፌሽናል ማዕከላዊ አጥቂ እግር ኳስ ተጫዋች የተሰጥኦ ምዝገባ', 'Galmee Dandeettii Taphataa Kubbaa Miilaa Fuulduraa', 'ናይ ፕሮፌሽናል ማእከላይ ኣጥቃዒ ኩዕሶ እግሪ ተውህቦ ምዝገባ', 'Verified club-ready football center forward, right/left-footed, aerial prowess, open for professional domestic and international trials', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'verified-scouting-trial-dossier' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'professional-football-striker-talent-listing', 6, true, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'verified-scouting-trial-dossier' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Professional Center Forward Football Striker Talent Listing', 'የፕሮፌሽናል ማዕከላዊ አጥቂ እግር ኳስ ተጫዋች የተሰጥኦ ምዝገባ', 'Galmee Dandeettii Taphataa Kubbaa Miilaa Fuulduraa', 'ናይ ፕሮፌሽናል ማእከላይ ኣጥቃዒ ኩዕሶ እግሪ ተውህቦ ምዝገባ', 'professional-football-striker-talent-listing', 'Verified club-ready football center forward, right/left-footed, aerial prowess, open for professional domestic and international trials', 6, v_ch_parent_id, true, true, 'SYS_PROD')
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

    -- --- Node: enterprise-fintech-ecommerce-portal-development (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_WORK_ENTERPRISE_DEV_COMPLETE', 'Enterprise E-Commerce & FinTech Portal Engineering Service', 'Enterprise E-Commerce & FinTech Portal Engineering Service', 'የድርጅት ደረጃ የኢ-ኮሜርስና የፋይናንስ ፖርታል ማልማት አገልግሎት', 'Tajaajila Misoonsa Poortaalii Daldala Dijitaalaa fi Faayinaansii', 'ናይ ትካል ደረጃ ናይ ኢ-ኮሜርስን ፋይናንስን ፖርታል ምዕባለ ኣገልግሎት', 'Production-grade scalable digital portal development with integrated automated payments, SMS notifications, and admin dashboard', 'SERVICE', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'turnkey-custom-app-development-service' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'enterprise-fintech-ecommerce-portal-development', 6, true, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'turnkey-custom-app-development-service' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Enterprise E-Commerce & FinTech Portal Engineering Service', 'የድርጅት ደረጃ የኢ-ኮሜርስና የፋይናንስ ፖርታል ማልማት አገልግሎት', 'Tajaajila Misoonsa Poortaalii Daldala Dijitaalaa fi Faayinaansii', 'ናይ ትካል ደረጃ ናይ ኢ-ኮሜርስን ፋይናንስን ፖርታል ምዕባለ ኣገልግሎት', 'enterprise-fintech-ecommerce-portal-development', 'Production-grade scalable digital portal development with integrated automated payments, SMS notifications, and admin dashboard', 6, v_ch_parent_id, true, true, 'SYS_PROD')
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

    -- --- Node: ingco-cidli200215-20v-cordless-impact-drill-complete (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_TOOL_INGCO_CIDLI200215_COMPLETE', 'INGCO CIDLI200215 20V Cordless Brushless Impact Drill Complete Kit', 'INGCO CIDLI200215 20V Cordless Brushless Impact Drill Complete Kit', 'ኢንግኮ CIDLI200215 20ቮልት ገመድ-አልባ ኢምፓክት ድሪል ሙሉ ኪት', 'Ingco CIDLI200215 20V Boraasuu Shibo-malee Guutuu', 'ኢንግኮ CIDLI200215 20ቮልት ሽቦ ኣልቦ ኢምፓክት ድሪል ሙሉእ', 'Complete contractor kit: 20V brushless impact drill, 2x lithium batteries, intelligent charger, carrying case', 'EQUIPMENT', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = 'ingco-20v-lithium-drill-set-with-accessories' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'ingco-cidli200215-20v-cordless-impact-drill-complete', 6, true, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = 'ingco-20v-lithium-drill-set-with-accessories' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('INGCO CIDLI200215 20V Cordless Brushless Impact Drill Complete Kit', 'ኢንግኮ CIDLI200215 20ቮልት ገመድ-አልባ ኢምፓክት ድሪል ሙሉ ኪት', 'Ingco CIDLI200215 20V Boraasuu Shibo-malee Guutuu', 'ኢንግኮ CIDLI200215 20ቮልት ሽቦ ኣልቦ ኢምፓክት ድሪል ሙሉእ', 'ingco-cidli200215-20v-cordless-impact-drill-complete', 'Complete contractor kit: 20V brushless impact drill, 2x lithium batteries, intelligent charger, carrying case', 6, v_ch_parent_id, true, true, 'SYS_PROD')
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

    -- --- Node: hulubedej-ev-rickshaw-eco-commuter-complete (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_VEH_EV_RICKSHAW_COMPLETE', 'Hulubedej Eco-Commuter 120km Electric Passenger Rickshaw Complete', 'Hulubedej Eco-Commuter 120km Electric Passenger Rickshaw Complete', 'ሁሉበደጅ ኢኮ-ኮምዩተር 120ኪሜ የኤሌክትሪክ መንገደኛ ባጃጅ ሙሉ ተሽከርካሪ', 'Hulubedej Eco-Commuter 120km Baajaajii Elektiriikii Guutuu', 'ሁሉበደጅ ኢኮ-ኮምዩተር 120ኪሜ ናይ ኤሌክትሪክ መጓዓዝያ ባጃጅ ሙሉእ', 'Complete electric passenger auto-rickshaw: 3kW motor, 72V 100Ah lithium battery, regenerative brakes, canopy cover', 'EQUIPMENT', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '120km-range-lithium-electric-bajaj' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'hulubedej-ev-rickshaw-eco-commuter-complete', 6, true, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '120km-range-lithium-electric-bajaj' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Hulubedej Eco-Commuter 120km Electric Passenger Rickshaw Complete', 'ሁሉበደጅ ኢኮ-ኮምዩተር 120ኪሜ የኤሌክትሪክ መንገደኛ ባጃጅ ሙሉ ተሽከርካሪ', 'Hulubedej Eco-Commuter 120km Baajaajii Elektiriikii Guutuu', 'ሁሉበደጅ ኢኮ-ኮምዩተር 120ኪሜ ናይ ኤሌክትሪክ መጓዓዝያ ባጃጅ ሙሉእ', 'hulubedej-ev-rickshaw-eco-commuter-complete', 'Complete electric passenger auto-rickshaw: 3kW motor, 72V 100Ah lithium battery, regenerative brakes, canopy cover', 6, v_ch_parent_id, true, true, 'SYS_PROD')
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

    -- --- Node: certified-mareko-fana-sun-dried-pepper-25kg (Depth 6, Role LEAF) ---
    -- 1. Concept
    INSERT INTO classification_concepts (code, canonical_name, name_en, name_am, name_om, name_ti, description, concept_type, status)
    VALUES ('CONCEPT_AGRI_MAREKO_25KG_CERTIFIED', 'Certified Organic Mareko Fana Whole Sun-Dried Red Pepper 25kg', 'Certified Organic Mareko Fana Whole Sun-Dried Red Pepper 25kg', 'የተረጋገጠ ኦርጋኒክ የማረቆ ፋና በፀሐይ የደረቀ ቀይ በርበሬ 25ኪግ', 'Barbaree Oorgaanikii Maaraqoo Faanaa 25kg Mirkanaa''e', 'ዝተረጋገጸ ኦርጋኒክ ማረቆ ፋና ብፀሓይ ዝነቐጸ ቀይሕ በርበረ 25ኪግ', 'Certified export-quality Mareko Fana dried pepper: zero aflatoxin, vibrant natural color, sorted and pesticide-free', 'PRODUCT', 'ACTIVE')
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
    SELECT id INTO v_parent_pid FROM classification_placements WHERE scheme_id = v_scheme_id AND slug = '25kg-sack-sun-dried-mareko-fana-pods' LIMIT 1;
    INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
    VALUES (v_scheme_id, v_cid, v_parent_pid, 'LEAF', 'certified-mareko-fana-sun-dried-pepper-25kg', 6, true, true, 'PUBLISHED')
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
    SELECT id INTO v_ch_parent_id FROM category_hierarchy WHERE slug = '25kg-sack-sun-dried-mareko-fana-pods' LIMIT 1;
    INSERT INTO category_hierarchy (name_en, name_am, name_om, name_ti, slug, description, level, parent_id, is_active, is_leaf, classification_system_code)
    VALUES ('Certified Organic Mareko Fana Whole Sun-Dried Red Pepper 25kg', 'የተረጋገጠ ኦርጋኒክ የማረቆ ፋና በፀሐይ የደረቀ ቀይ በርበሬ 25ኪግ', 'Barbaree Oorgaanikii Maaraqoo Faanaa 25kg Mirkanaa''e', 'ዝተረጋገጸ ኦርጋኒክ ማረቆ ፋና ብፀሓይ ዝነቐጸ ቀይሕ በርበረ 25ኪግ', 'certified-mareko-fana-sun-dried-pepper-25kg', 'Certified export-quality Mareko Fana dried pepper: zero aflatoxin, vibrant natural color, sorted and pesticide-free', 6, v_ch_parent_id, true, true, 'SYS_PROD')
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

-- 3. Poly-Hierarchy Relationships (Drill & Bajaj cross-scheme semantic links)
DO $$
DECLARE
    v_drill_cid INT;
    v_bajaj_cid INT;
    v_elec_serv_cid INT;
    v_const_serv_cid INT;
    v_trans_serv_cid INT;
BEGIN
    SELECT id INTO v_drill_cid FROM classification_concepts WHERE code = 'CONCEPT_TOOL_INGCO_CIDLI200215_COMPLETE' LIMIT 1;
    SELECT id INTO v_bajaj_cid FROM classification_concepts WHERE code = 'CONCEPT_VEH_BAJAJ_RE_205CC' LIMIT 1;
    SELECT id INTO v_elec_serv_cid FROM classification_concepts WHERE code = 'CONCEPT_SERV_FULL_WIRING_SERVICE' LIMIT 1;
    SELECT id INTO v_const_serv_cid FROM classification_concepts WHERE code = 'CONCEPT_TOOL_BOSCH_GBH_226_COMPLETE' LIMIT 1;
    SELECT id INTO v_trans_serv_cid FROM classification_concepts WHERE code = 'CONCEPT_VEH_ISUZU_NPR71_COMPLETE' LIMIT 1;

    IF v_drill_cid IS NOT NULL AND v_elec_serv_cid IS NOT NULL THEN
        INSERT INTO classification_relationships (source_concept_id, target_concept_id, relationship_type_code, status)
        VALUES (v_drill_cid, v_elec_serv_cid, 'USED_IN', 'ACTIVE')
        ON CONFLICT (source_concept_id, target_concept_id, relationship_type_code, context_scheme_id) DO NOTHING;
    END IF;

    IF v_drill_cid IS NOT NULL AND v_const_serv_cid IS NOT NULL THEN
        INSERT INTO classification_relationships (source_concept_id, target_concept_id, relationship_type_code, status)
        VALUES (v_drill_cid, v_const_serv_cid, 'COMPATIBLE_WITH', 'ACTIVE')
        ON CONFLICT (source_concept_id, target_concept_id, relationship_type_code, context_scheme_id) DO NOTHING;
    END IF;

    IF v_bajaj_cid IS NOT NULL AND v_trans_serv_cid IS NOT NULL THEN
        INSERT INTO classification_relationships (source_concept_id, target_concept_id, relationship_type_code, status)
        VALUES (v_bajaj_cid, v_trans_serv_cid, 'RELATED_TO', 'ACTIVE')
        ON CONFLICT (source_concept_id, target_concept_id, relationship_type_code, context_scheme_id) DO NOTHING;
    END IF;
END $$;

-- 4. Mirror non-materialistic placements into SYS_JOBS and SYS_SERV
DO $$
DECLARE
    v_prod_id INT;
    v_jobs_id INT;
    v_serv_id INT;
    v_root_pid INT;
BEGIN
    SELECT id INTO v_prod_id FROM classification_schemes WHERE code = 'SYS_PROD';
    SELECT id INTO v_jobs_id FROM classification_schemes WHERE code = 'SYS_JOBS';
    SELECT id INTO v_serv_id FROM classification_schemes WHERE code = 'SYS_SERV';

    IF v_jobs_id IS NOT NULL THEN
        -- Ensure root in SYS_JOBS
        SELECT id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_jobs_id AND depth = 0 LIMIT 1;
        IF v_root_pid IS NULL THEN
            SELECT concept_id INTO v_root_pid FROM classification_placements WHERE scheme_id = v_prod_id AND slug = 'workforce-employment' LIMIT 1;
            IF v_root_pid IS NOT NULL THEN
                INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
                VALUES (v_jobs_id, v_root_pid, NULL, 'ROOT', 'workforce-employment-root', 0, false, true, 'PUBLISHED')
                ON CONFLICT (scheme_id, slug) DO UPDATE SET depth = 0 RETURNING id INTO v_root_pid;
            END IF;
        END IF;

        IF v_root_pid IS NOT NULL THEN
            INSERT INTO classification_placements (scheme_id, concept_id, parent_placement_id, role, slug, depth, is_leaf, is_active, status)
            SELECT v_jobs_id, cp.concept_id, v_root_pid, cp.role, 'job-' || cp.slug, cp.depth, cp.is_leaf, true, 'PUBLISHED'
            FROM classification_placements cp
            WHERE cp.scheme_id = v_prod_id AND (cp.slug LIKE '%creator%' OR cp.slug LIKE '%talent%' OR cp.slug LIKE '%football%' OR cp.slug LIKE '%developer%' OR cp.slug LIKE '%app%')
              AND cp.depth > 0
            ON CONFLICT (scheme_id, slug) DO NOTHING;
        END IF;
    END IF;
END $$;

-- 5. Re-evaluate is_leaf and tree integrity efficiently
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
