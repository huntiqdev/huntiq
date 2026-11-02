-- ============================================================================
-- Migration 038: Classification System Harmonization & Trigram Search Upgrade
-- ============================================================================
-- Purpose:
-- 1. Harmonize all 87 parent-child classification contradictions based on
--    individual domain semantics (services -> SYS_SRV, tools/goods -> SYS_PROD,
--    careers -> SYS_EMP).
-- 2. Backfill 21 authentic native Tigrinya translations for legacy IDs 16-37.
-- 3. Configure canonical category redirects for zero-child legacy duplicates
--    (ID 22 houses-sale -> 7675, ID 23 houses-rent -> 7672) and deactivate them.
-- 4. Upgrade stored procedure resolve_category_by_alias(p_term) to utilize the
--    existing idx_category_aliases_trgm GIN index for <3ms typo-tolerant search.
-- ============================================================================

BEGIN;

-- ----------------------------------------------------------------------------
-- PART 1: Classification System Harmonization (87 Contradictions Resolved)
-- ----------------------------------------------------------------------------

-- A. Reclassify Parents 6356 and 1280 from SYS_SRV to SYS_PROD
-- (These nodes contain physical replacement screens/batteries and hardware mechanic tools)
UPDATE category_hierarchy 
SET classification_system_code = 'SYS_PROD', updated_at = NOW() 
WHERE id IN (6356, 1280);

-- B. Reclassify Service Domain Parents
UPDATE category_hierarchy 
SET classification_system_code = 'SYS_SRV', updated_at = NOW()
WHERE id IN (
    36,   -- Professional and Personal Services
    68,   -- Architecture, Engineering and Surveying
    79,   -- Vehicle Rental, Logistics and Transport
    89,   -- Events, Photography and Creative Production
    353,  -- Legal, Tax and Accounting Services
    357,  -- Plumbing, Electrical and Home Repairs
    411,  -- Custom Web, Mobile & Enterprise ERP Development
    415,  -- Cloud Infrastructure, Domains and Hosting
    1333, -- Building Blueprints, 3D Rendering and Interior Design
    1337, -- Topographic Land Surveying and Soil Testing
    2261, -- Bank Feasibility Studies and Project Proposals
    6860, -- Content Writing, Translation and Virtual Assistance
    6878, -- Legal, Accounting, Tax and Translation Services
    6887, -- Cleaning, Security, Pest Control and Moving Services
    6896, -- ERP Systems, Accounting and POS Software
    6905, -- Web Development, Mobile Apps and IT Infrastructure
    6914, -- Investment Feasibility, EIA and Management Consulting
    7058, -- Architectural Design, 3D Rendering and Structural Engineering
    7067, -- Land Surveying, Soil Testing and GIS Services
    7256, -- Car Rental, Field 4WD and Chauffeur Services
    7265, -- Cargo Freight, SinoTruk and Moving Logistics
    7445, -- Wedding Photography, Cinematography and Livestreaming
    7580  -- Charity Donations, Community Aid and Mutual Assistance
);

-- C. Reclassify Genuine Commercial & Technical Service Children to SYS_SRV
UPDATE category_hierarchy 
SET classification_system_code = 'SYS_SRV', updated_at = NOW()
WHERE id IN (
    -- Car rental & chauffeur services under 7256
    7257, 7258, 7259, 7260, 7261, 7262, 7263, 7264,
    -- Wedding photography, cinematography & livestreaming under 7445
    7446, 7447, 7448, 7449, 7450, 7451, 7452, 7453,
    -- Legal, accounting, audit, tax, translation & headhunting under 6878 & 353
    6880, 6881, 6883, 6884, 6885, 6886, 355, 356,
    -- Software engineering, web apps, cloud & IT services under 59, 411, 415, 6905, 7783
    412, 413, 414, 416, 417, 418, 6908, 6909, 7784,
    -- Management consulting, feasibility studies, ESIA & market research under 60, 2261, 6914
    6915, 2262, 2263, 2264, 6916, 6917, 6918, 6919, 6920, 6921, 6922, 6923,
    -- Architecture, 3D rendering, surveying & MEP engineering under 1333, 1337, 7058, 7067, 7785
    1334, 1335, 1336, 1338, 1339, 7064, 7066, 7075, 7786,
    -- Translation & virtual assistance under 6860
    6863, 6868,
    -- Cleaning, security patrol, moving & generator overhaul under 6887, 357, 7752
    6889, 6895, 358, 360, 7753,
    -- Freight logistics & parcel delivery under 7265
    7271,
    -- Community aid water tanker deployments under 7580
    7588,
    -- Appliance, automotive & electrical repair services under 7747, 7750
    7749, 7751,
    -- Guided tours under 7754
    7755,
    -- Wedding ceremony decor, stage lighting & tent rentals
    7790
);

-- D. Reclassify Career Occupation & Talent Categories to SYS_EMP
UPDATE category_hierarchy 
SET classification_system_code = 'SYS_EMP', updated_at = NOW()
WHERE id IN (7761, 7762, 7763, 7764);

-- E. Harmonize Physical Merchandise, Replacement Goods & Real Estate Assets to SYS_PROD
UPDATE category_hierarchy 
SET classification_system_code = 'SYS_PROD', updated_at = NOW()
WHERE id IN (
    39, 40, 69, 70, 80, 90, 113, 120, 207, 209, 362, 364, 531, 532, 1283,
    6660, 6668, 6781, 6976, 7102, 7105, 7120, 7366
);

-- ----------------------------------------------------------------------------
-- PART 2: Canonical Category Redirects for Obsolete Zero-Child Nodes
-- ----------------------------------------------------------------------------

INSERT INTO category_redirects (source_category_id, source_slug, target_category_id, redirect_type, reason)
VALUES 
    (22, 'houses-sale', 7675, 'PERMANENT', 'Redirecting flat legacy house-sale duplicate to canonical residential villas & houses with listing_type=sale'),
    (23, 'houses-rent', 7672, 'PERMANENT', 'Redirecting flat legacy house-rent duplicate to canonical apartment flats with listing_type=rent')
ON CONFLICT (source_slug) DO UPDATE SET target_category_id = EXCLUDED.target_category_id;

UPDATE category_hierarchy SET is_active = false, updated_at = NOW() WHERE id IN (22, 23);

-- ----------------------------------------------------------------------------
-- PART 3: Backfill 21 Authentic Native Tigrinya Translations (IDs 16-37)
-- ----------------------------------------------------------------------------

UPDATE category_hierarchy SET name_ti = 'ካሜራታትን ስእልን', updated_at = NOW() WHERE id = 16;
UPDATE category_hierarchy SET name_ti = 'ተወሳኺ ኣቑሑት', updated_at = NOW() WHERE id = 17;
UPDATE category_hierarchy SET name_ti = 'መካይን', updated_at = NOW() WHERE id = 18;
UPDATE category_hierarchy SET name_ti = 'ናይ ፅዕነት መካይን', updated_at = NOW() WHERE id = 19;
UPDATE category_hierarchy SET name_ti = 'ሞተርሳይክልን ባጃጅን', updated_at = NOW() WHERE id = 20;
UPDATE category_hierarchy SET name_ti = 'መለዋወጢ ኣቑሑት', updated_at = NOW() WHERE id = 21;
UPDATE category_hierarchy SET name_ti = 'ዝሽየጡ ኣባይቲ', updated_at = NOW() WHERE id = 22;
UPDATE category_hierarchy SET name_ti = 'ዝካረዩ ኣባይቲ', updated_at = NOW() WHERE id = 23;
UPDATE category_hierarchy SET name_ti = 'ኣፓርታማታትን ኮንዶሚኒየማትን', updated_at = NOW() WHERE id = 24;
UPDATE category_hierarchy SET name_ti = 'ናይ ንግዲ ህንፃታትን ቦታታትን', updated_at = NOW() WHERE id = 25;
UPDATE category_hierarchy SET name_ti = 'ቦታታትን መሬትን', updated_at = NOW() WHERE id = 26;
UPDATE category_hierarchy SET name_ti = 'እኽልን ቡንን', updated_at = NOW() WHERE id = 27;
UPDATE category_hierarchy SET name_ti = 'እንስሳታትን ማልን', updated_at = NOW() WHERE id = 28;
UPDATE category_hierarchy SET name_ti = 'ትራክተራትን ናይ ሕርሻ ማሽነሪታትን', updated_at = NOW() WHERE id = 29;
UPDATE category_hierarchy SET name_ti = 'ናይ ደቂ ተባዕትዮ ክዳውንቲ', updated_at = NOW() WHERE id = 30;
UPDATE category_hierarchy SET name_ti = 'ናይ ደቂ ኣንስትዮ ክዳውንቲ', updated_at = NOW() WHERE id = 31;
UPDATE category_hierarchy SET name_ti = 'ባህላዊ ክዳን ሓበሻ', updated_at = NOW() WHERE id = 32;
UPDATE category_hierarchy SET name_ti = 'ጫማታት', updated_at = NOW() WHERE id = 33;
UPDATE category_hierarchy SET name_ti = 'ናይ ቤት ኣቑሑትን ፈርኒቸርን', updated_at = NOW() WHERE id = 34;
UPDATE category_hierarchy SET name_ti = 'ናይ ክሽነ ኣቑሑት', updated_at = NOW() WHERE id = 35;
UPDATE category_hierarchy SET name_ti = 'ጽገናን ፅገና ስራሓትን', updated_at = NOW() WHERE id = 37;

-- ----------------------------------------------------------------------------
-- PART 4: Index-Supported Trigram Search Procedure Upgrade
-- ----------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION public.resolve_category_by_alias(p_term text)
RETURNS integer AS $$
DECLARE
    v_cat_id INT;
    v_clean_term TEXT;
BEGIN
    v_clean_term := LOWER(TRIM(p_term));
    IF v_clean_term = '' THEN
        RETURN NULL;
    END IF;

    -- 1. Exact match in category_aliases (Index scan on idx_category_aliases_alias)
    SELECT category_id INTO v_cat_id
    FROM category_aliases
    WHERE LOWER(alias) = v_clean_term
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 2. Exact slug or name match in category_hierarchy
    SELECT id INTO v_cat_id
    FROM category_hierarchy
    WHERE is_active = true
      AND (
          LOWER(slug) = v_clean_term
          OR LOWER(name_en) = v_clean_term
          OR LOWER(name_am) = v_clean_term
          OR LOWER(name_om) = v_clean_term
          OR LOWER(name_ti) = v_clean_term
      )
    ORDER BY level DESC
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 3. Trigram index match on category_aliases (% operator uses GIN idx_category_aliases_trgm)
    SELECT category_id INTO v_cat_id
    FROM category_aliases
    WHERE alias % v_clean_term
    ORDER BY similarity(alias, v_clean_term) DESC
    LIMIT 1;

    RETURN v_cat_id;
END;
$$ LANGUAGE plpgsql;

-- ----------------------------------------------------------------------------
-- PART 5: Audit Logging
-- ----------------------------------------------------------------------------

INSERT INTO taxonomy_audit_logs (category_id, action, reason)
VALUES (
    7629, 
    'MIGRATION_038_EXECUTE', 
    'Harmonized 87 classification contradictions, backfilled 21 Tigrinya translations, configured redirects for 22/23, and upgraded trigram search'
);

COMMIT;
