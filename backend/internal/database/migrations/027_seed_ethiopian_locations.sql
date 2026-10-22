-- ============================================================================
-- Migration 027: Authoritative Ethiopian Location Hierarchy Seed Dataset
-- ============================================================================
-- Purpose: Populate verified Ethiopian administrative structure:
--          - Country (Ethiopia)
--          - 12 Regions & 2 Chartered Cities (Addis Ababa, Dire Dawa)
--          - All 11 Addis Ababa Subcities with coordinates & bounding boxes
--          - Major metropolitan cities & regional zones
--          - Synchronize legacy locations table for backward compatibility
-- ============================================================================

DO $$
DECLARE
    v_country_id INT;
    v_addis_id INT;
    v_dire_id INT;
    v_oromia_id INT;
    v_amhara_id INT;
    v_tigray_id INT;
    v_sidama_id INT;
    v_somali_id INT;
    v_afar_id INT;
    v_bengum_id INT;
    v_gambela_id INT;
    v_harari_id INT;
    v_south_id INT;
    v_central_id INT;
    v_southwest_id INT;
    
    -- Subcity IDs
    v_bole_id INT;
    v_kirkos_id INT;
    v_yeka_id INT;
    v_arada_id INT;
    v_lideta_id INT;
    v_ak_id INT;
    v_gullele_id INT;
    v_kolfe_id INT;
    v_akaki_id INT;
    v_nifas_id INT;
    v_lemi_id INT;
BEGIN
    -- 1. Country: Ethiopia
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        NULL, 'country', 'Ethiopia', 'ኢትዮጵያ', 'Itoophiyaa', 'ኢትዮጵያ', 'ethiopia', 'ET',
        9.1450000, 40.4896730, 3.4000000, 14.9000000, 32.9900000, 47.9800000,
        ARRAY['Abyssinia', 'Federal Democratic Republic of Ethiopia'], 1
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_country_id;

    IF v_country_id IS NULL THEN
        SELECT id INTO v_country_id FROM location_hierarchy WHERE slug = 'ethiopia';
    END IF;

    -- 2. Chartered Cities & Regions (Level 1, parent = Ethiopia)
    
    -- Addis Ababa (Chartered City)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'city', 'Addis Ababa', 'አዲስ አበባ', 'Finfinnee', 'ኣዲስ ኣበባ', 'addis-ababa', 'ET-AA',
        9.0320000, 38.7468900, 8.8300000, 9.1000000, 38.6500000, 38.9000000,
        ARRAY['Finfinne', 'Addis Abeba', 'Sheger', 'AA'], 1
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_addis_id;
    IF v_addis_id IS NULL THEN SELECT id INTO v_addis_id FROM location_hierarchy WHERE slug = 'addis-ababa'; END IF;

    -- Dire Dawa (Chartered City)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'city', 'Dire Dawa', 'ድሬዳዋ', 'Dirree Dawaa', 'ድሬዳዋ', 'dire-dawa', 'ET-DD',
        9.5930000, 41.8660000, 9.4500000, 9.7500000, 41.6500000, 42.0500000,
        ARRAY['Dire Dawa City', 'Dirree Dhawaa'], 2
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_dire_id;
    IF v_dire_id IS NULL THEN SELECT id INTO v_dire_id FROM location_hierarchy WHERE slug = 'dire-dawa'; END IF;

    -- Oromia
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'Oromia', 'ኦሮሚያ', 'Oromiyaa', 'ኦሮሚያ', 'oromia', 'ET-OR',
        8.5410000, 39.2680000, 3.4000000, 10.4000000, 34.1000000, 42.9000000,
        ARRAY['Oromiya', 'Oromiyaa'], 3
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_oromia_id;
    IF v_oromia_id IS NULL THEN SELECT id INTO v_oromia_id FROM location_hierarchy WHERE slug = 'oromia'; END IF;

    -- Amhara
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'Amhara', 'አማራ', 'Amaara', 'ኣምሓራ', 'amhara', 'ET-AM',
        11.5900000, 37.3900000, 8.8000000, 13.8000000, 35.2000000, 40.3000000,
        ARRAY['Amhara Region', 'Amhara National Regional State'], 4
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_amhara_id;
    IF v_amhara_id IS NULL THEN SELECT id INTO v_amhara_id FROM location_hierarchy WHERE slug = 'amhara'; END IF;

    -- Tigray
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'Tigray', 'ትግራይ', 'Tigraay', 'ትግራይ', 'tigray', 'ET-TI',
        14.0300000, 38.3160000, 12.2000000, 14.9000000, 36.4500000, 40.0000000,
        ARRAY['Tigrai', 'Tigre'], 5
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_tigray_id;
    IF v_tigray_id IS NULL THEN SELECT id INTO v_tigray_id FROM location_hierarchy WHERE slug = 'tigray'; END IF;

    -- Sidama
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'Sidama', 'ሲዳማ', 'Sidaama', 'ሲዳማ', 'sidama', 'ET-SI',
        6.7800000, 38.4500000, 6.2000000, 7.3000000, 38.0000000, 39.2000000,
        ARRAY['Sidama Region', 'Sidaama'], 6
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_sidama_id;
    IF v_sidama_id IS NULL THEN SELECT id INTO v_sidama_id FROM location_hierarchy WHERE slug = 'sidama'; END IF;

    -- Somali
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'Somali', 'ሶማሌ', 'Somaalee', 'ሶማሌ', 'somali', 'ET-SO',
        7.1500000, 44.2500000, 4.0000000, 11.0000000, 40.5000000, 48.0000000,
        ARRAY['Somali Region', 'Ogaden'], 7
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_somali_id;
    IF v_somali_id IS NULL THEN SELECT id INTO v_somali_id FROM location_hierarchy WHERE slug = 'somali'; END IF;

    -- Afar
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'Afar', 'አፋር', 'Afaar', 'ዓፋር', 'afar', 'ET-AF',
        11.7500000, 40.9500000, 8.8000000, 14.5000000, 39.5000000, 42.5000000,
        ARRAY['Afar Region', 'Danakil'], 8
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_afar_id;
    IF v_afar_id IS NULL THEN SELECT id INTO v_afar_id FROM location_hierarchy WHERE slug = 'afar'; END IF;

    -- Benishangul-Gumuz
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'Benishangul-Gumuz', 'ቤንሻንጉል ጉሙዝ', 'Beniishaangul Gumuz', 'ቤኒሻንጉል ጉሙዝ', 'benishangul-gumuz', 'ET-BE',
        10.5000000, 35.5000000, 9.2000000, 12.0000000, 34.1000000, 37.0000000,
        ARRAY['Benishangul', 'Gumuz'], 9
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_bengum_id;
    IF v_bengum_id IS NULL THEN SELECT id INTO v_bengum_id FROM location_hierarchy WHERE slug = 'benishangul-gumuz'; END IF;

    -- Gambela
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'Gambela', 'ጋምቤላ', 'Gambellaa', 'ጋምቤላ', 'gambela', 'ET-GA',
        7.9000000, 34.5000000, 6.5000000, 8.8000000, 33.0000000, 35.5000000,
        ARRAY['Gambella', 'Gambella Region'], 10
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_gambela_id;
    IF v_gambela_id IS NULL THEN SELECT id INTO v_gambela_id FROM location_hierarchy WHERE slug = 'gambela'; END IF;

    -- Harari
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'Harari', 'ሐረሪ', 'Hararii', 'ሓረሪ', 'harari', 'ET-HA',
        9.3100000, 42.1300000, 9.1500000, 9.4500000, 42.0000000, 42.3000000,
        ARRAY['Harar Region', 'Harari People'], 11
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_harari_id;
    IF v_harari_id IS NULL THEN SELECT id INTO v_harari_id FROM location_hierarchy WHERE slug = 'harari'; END IF;

    -- South Ethiopia
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'South Ethiopia', 'ደቡብ ኢትዮጵያ', 'Kibba Itoophiyaa', 'ደቡብ ኢትዮጵያ', 'south-ethiopia', 'ET-SE',
        6.0000000, 37.5000000, 4.5000000, 7.5000000, 35.8000000, 38.5000000,
        ARRAY['Wolayita', 'Gamo', 'Gofa', 'South Region'], 12
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_south_id;
    IF v_south_id IS NULL THEN SELECT id INTO v_south_id FROM location_hierarchy WHERE slug = 'south-ethiopia'; END IF;

    -- Central Ethiopia
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'Central Ethiopia', 'ማዕከላዊ ኢትዮጵያ', 'Giddugala Itoophiyaa', 'ማእከላይ ኢትዮጵያ', 'central-ethiopia', 'ET-CE',
        7.8000000, 38.0000000, 7.0000000, 8.5000000, 37.0000000, 39.0000000,
        ARRAY['Gurage', 'Silte', 'Hadiya', 'Kambata'], 13
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_central_id;
    IF v_central_id IS NULL THEN SELECT id INTO v_central_id FROM location_hierarchy WHERE slug = 'central-ethiopia'; END IF;

    -- Southwest Ethiopia
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
    ) VALUES (
        v_country_id, 'region', 'Southwest Ethiopia', 'ደቡብ ምዕራብ ኢትዮጵያ', 'Kibba-Lixa Itoophiyaa', 'ደቡብ ምዕራብ ኢትዮጵያ', 'southwest-ethiopia', 'ET-SW',
        7.1500000, 36.2500000, 5.8000000, 8.2000000, 35.0000000, 37.2000000,
        ARRAY['Kaffa', 'Sheka', 'Bench Sheko', 'Dawro', 'West Omo'], 14
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_southwest_id;
    IF v_southwest_id IS NULL THEN SELECT id INTO v_southwest_id FROM location_hierarchy WHERE slug = 'southwest-ethiopia'; END IF;

    -- 3. All 11 Addis Ababa Subcities (Level 2, parent = Addis Ababa)
    -- Bole
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_addis_id, 'subcity', 'Bole', 'ቦሌ', 'Boolee', 'ቦሌ', 'addis-ababa-bole', 'AA-BO',
        8.9950000, 38.7900000, 8.9400000, 9.0400000, 38.7500000, 38.8600000,
        ARRAY['Bole Subcity', 'Bole Airport', 'Bole Medhanealem', 'Bole Atlas', 'Megenagna'], 1, true
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_bole_id;

    -- Kirkos
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_addis_id, 'subcity', 'Kirkos', 'ቂርቆስ', 'Qirqoos', 'ቂርቆስ', 'addis-ababa-kirkos', 'AA-KI',
        9.0080000, 38.7580000, 8.9900000, 9.0300000, 38.7400000, 38.7800000,
        ARRAY['Kirkos Subcity', 'Kazanchis', 'Meshualekia', 'Meskel Square', 'Olympia'], 2, true
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_kirkos_id;

    -- Yeka
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_addis_id, 'subcity', 'Yeka', 'የካ', 'Yakkkaa', 'የካ', 'addis-ababa-yeka', 'AA-YE',
        9.0420000, 38.8050000, 9.0100000, 9.0900000, 38.7700000, 38.8500000,
        ARRAY['Yeka Subcity', 'Kotebe', 'Ayat', 'CMC', 'Ferensay'], 3, true
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_yeka_id;

    -- Arada
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_addis_id, 'subcity', 'Arada', 'አራዳ', 'Aradaa', 'ኣራዳ', 'addis-ababa-arada', 'AA-AR',
        9.0350000, 38.7520000, 9.0200000, 9.0550000, 38.7350000, 38.7700000,
        ARRAY['Arada Subcity', 'Piazza', '4 Kilo', '6 Kilo', 'Amist Kilo'], 4, true
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_arada_id;

    -- Lideta
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_addis_id, 'subcity', 'Lideta', 'ልደታ', 'Lidetaa', 'ልደታ', 'addis-ababa-lideta', 'AA-LI',
        9.0120000, 38.7350000, 8.9950000, 9.0300000, 38.7150000, 38.7500000,
        ARRAY['Lideta Subcity', 'Abanet', 'Sengatera', 'Mexico'], 5, true
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_lideta_id;

    -- Addis Ketema
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_addis_id, 'subcity', 'Addis Ketema', 'አዲስ ከተማ', 'Magaalaa Haaraa', 'ኣዲስ ከተማ', 'addis-ababa-addis-ketema', 'AA-AK',
        9.0330000, 38.7300000, 9.0150000, 9.0500000, 38.7100000, 38.7450000,
        ARRAY['Merkato', 'Addis Ketema Subcity', 'Autobus Tera'], 6, true
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_ak_id;

    -- Gullele
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_addis_id, 'subcity', 'Gullele', 'ጉለሌ', 'Gullallee', 'ጉለሌ', 'addis-ababa-gullele', 'AA-GU',
        9.0680000, 38.7350000, 9.0400000, 9.1100000, 38.7000000, 38.7650000,
        ARRAY['Gullele Subcity', 'Shiro Meda', 'Entoto', 'Addisu Gebeya'], 7, true
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_gullele_id;

    -- Kolfe Keranio
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_addis_id, 'subcity', 'Kolfe Keranio', 'ኮልፌ ቀራኒዮ', 'Kolfi Qaraaniyoo', 'ኮልፌ ቀራንዮ', 'addis-ababa-kolfe-keranio', 'AA-KK',
        9.0200000, 38.7050000, 8.9750000, 9.0650000, 38.6650000, 38.7300000,
        ARRAY['Kolfe', 'Keranio', 'Tor Hailoch', 'Ayer Tena'], 8, true
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_kolfe_id;

    -- Akaki Kality
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_addis_id, 'subcity', 'Akaki Kality', 'አቃቂ ቃሊቲ', 'Aqaaqii Qaallittii', 'ኣቃቂ ቃሊቲ', 'addis-ababa-akaki-kality', 'AA-AQ',
        8.9100000, 38.7650000, 8.8400000, 8.9700000, 38.7100000, 38.8300000,
        ARRAY['Akaki', 'Kality', 'Tulu Dimtu', 'Koye Feche'], 9, true
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_akaki_id;

    -- Nifas Silk-Lafto
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_addis_id, 'subcity', 'Nifas Silk-Lafto', 'ንፋስ ስልክ ላፍቶ', 'Nifas Silk Laaftoo', 'ንፋስ ስልክ ላፍቶ', 'addis-ababa-nifas-silk-lafto', 'AA-NL',
        8.9700000, 38.7350000, 8.9250000, 9.0150000, 38.6900000, 38.7650000,
        ARRAY['Nefas Silk', 'Lafto', 'Jemo', 'Lebu', 'Mekanisa', 'Gotera'], 10, true
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_nifas_id;

    -- Lemi Kura
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_addis_id, 'subcity', 'Lemi Kura', 'ለሚ ኩራ', 'Lami Kuraa', 'ለሚ ኩራ', 'addis-ababa-lemi-kura', 'AA-LK',
        9.0150000, 38.8400000, 8.9700000, 9.0700000, 38.8000000, 38.9000000,
        ARRAY['Lemi Kura', 'Goro', 'Bole Arabsa', 'Summit'], 11, true
    ) ON CONFLICT DO NOTHING RETURNING id INTO v_lemi_id;

    -- 4. Key Metropolitan Cities & Major Regional Hubs
    -- Adama (Oromia)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_oromia_id, 'city', 'Adama', 'አዳማ', 'Adaamaa', 'ኣዳማ', 'oromia-adama', 'OR-ADA',
        8.5400000, 39.2700000, 8.4800000, 8.6000000, 39.2000000, 39.3400000,
        ARRAY['Nazret', 'Nazareth'], 1, true
    ) ON CONFLICT DO NOTHING;

    -- Bishoftu (Oromia)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_oromia_id, 'city', 'Bishoftu', 'ቢሾፍቱ', 'Bishooftuu', 'ቢሾፍቱ', 'oromia-bishoftu', 'OR-BIS',
        8.7500000, 38.9800000, 8.7000000, 8.8000000, 38.9200000, 39.0400000,
        ARRAY['Debre Zeyit', 'Debre Zeit'], 2, true
    ) ON CONFLICT DO NOTHING;

    -- Jimma (Oromia)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_oromia_id, 'city', 'Jimma', 'ጅማ', 'Jimmaa', 'ጅማ', 'oromia-jimma', 'OR-JIM',
        7.6700000, 36.8300000, 7.6200000, 7.7300000, 36.7800000, 36.9000000,
        ARRAY['Jima'], 3, true
    ) ON CONFLICT DO NOTHING;

    -- Shashemene (Oromia)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_oromia_id, 'city', 'Shashemene', 'ሻሸመኔ', 'Shaashamannee', 'ሻሸመኔ', 'oromia-shashemene', 'OR-SHA',
        7.2000000, 38.6000000, 7.1400000, 7.2500000, 38.5400000, 38.6500000,
        ARRAY['Shashamane'], 4, true
    ) ON CONFLICT DO NOTHING;

    -- Bahir Dar (Amhara)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_amhara_id, 'city', 'Bahir Dar', 'ባሕር ዳር', 'Baahir Daar', 'ባሕር ዳር', 'amhara-bahir-dar', 'AM-BD',
        11.5900000, 37.3900000, 11.5300000, 11.6500000, 37.3200000, 37.4500000,
        ARRAY['Bahr Dar'], 1, true
    ) ON CONFLICT DO NOTHING;

    -- Gondar (Amhara)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_amhara_id, 'city', 'Gondar', 'ጎንደር', 'Gondar', 'ጎንደር', 'amhara-gondar', 'AM-GD',
        12.6000000, 37.4600000, 12.5500000, 12.6600000, 37.4000000, 37.5200000,
        ARRAY['Gonder'], 2, true
    ) ON CONFLICT DO NOTHING;

    -- Dessie (Amhara)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_amhara_id, 'city', 'Dessie', 'ደሴ', 'Dasee', 'ደሴ', 'amhara-dessie', 'AM-DES',
        11.1300000, 39.6300000, 11.0800000, 11.1800000, 39.5800000, 39.6800000,
        ARRAY['Dese'], 3, true
    ) ON CONFLICT DO NOTHING;

    -- Mekelle (Tigray)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_tigray_id, 'city', 'Mekelle', 'መቀሌ', 'Maqalee', 'መቐለ', 'tigray-mekelle', 'TI-MK',
        13.4960000, 39.4750000, 13.4400000, 13.5600000, 39.4200000, 39.5400000,
        ARRAY['Makale', 'Meqele'], 1, true
    ) ON CONFLICT DO NOTHING;

    -- Hawassa (Sidama)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_sidama_id, 'city', 'Hawassa', 'ሀዋሳ', 'Hawaasaa', 'ሃዋሳ', 'sidama-hawassa', 'SI-HAW',
        7.0500000, 38.4800000, 6.9800000, 7.1200000, 38.4200000, 38.5400000,
        ARRAY['Awassa'], 1, true
    ) ON CONFLICT DO NOTHING;

    -- Jigjiga (Somali)
    INSERT INTO location_hierarchy (
        parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
        latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
    ) VALUES (
        v_somali_id, 'city', 'Jigjiga', 'ጅግጅጋ', 'Jigjigaa', 'ጅግጅጋ', 'somali-jigjiga', 'SO-JIG',
        9.3500000, 42.8000000, 9.3000000, 9.4000000, 42.7400000, 42.8600000,
        ARRAY['Jijiga'], 1, true
    ) ON CONFLICT DO NOTHING;

    -- 5. Backward Compatibility: Populate legacy `locations` table
    -- Sync Addis Ababa subcities as records in locations
    INSERT INTO locations (region, subcity, region_am, subcity_am, region_or, subcity_or, region_ti, subcity_ti)
    SELECT
        'Addis Ababa', l.name_en,
        'አዲስ አበባ', l.name_am,
        'Finfinnee', l.name_om,
        'ኣዲስ ኣበባ', l.name_ti
    FROM location_hierarchy l
    WHERE l.parent_id = v_addis_id AND l.location_type = 'subcity'
    ON CONFLICT DO NOTHING;

    -- Sync other major cities as records in locations
    INSERT INTO locations (region, zone, region_am, zone_am, region_or, zone_or, region_ti, zone_ti)
    SELECT
        p.name_en, l.name_en,
        p.name_am, l.name_am,
        p.name_om, l.name_om,
        p.name_ti, l.name_ti
    FROM location_hierarchy l
    JOIN location_hierarchy p ON p.id = l.parent_id
    WHERE l.location_type = 'city' AND p.location_type = 'region'
    ON CONFLICT DO NOTHING;

    -- Sync all regions into locations
    INSERT INTO locations (region, region_am, region_or, region_ti)
    SELECT
        l.name_en, l.name_am, l.name_om, l.name_ti
    FROM location_hierarchy l
    WHERE l.parent_id = v_country_id AND l.location_type IN ('region', 'city')
    ON CONFLICT DO NOTHING;

END $$;
