-- ============================================================================
-- 023_seed_ethiopian_locations.sql
-- Authoritative Ethiopian Location Hierarchy Seed Dataset
-- Populates:
--   - Country (Ethiopia)
--   - 12 Regions & 2 Chartered Cities (Addis Ababa, Dire Dawa)
--   - All 11 Addis Ababa Subcities with coordinates & bounding boxes
--   - Major metropolitan cities & regional zones
--   - Synchronizes legacy locations table for backward compatibility
-- ============================================================================

-- 1. Country: Ethiopia
INSERT INTO location_hierarchy (
    parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
    latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
) VALUES (
    NULL, 'country', 'Ethiopia', 'ኢትዮጵያ', 'Itoophiyaa', 'ኢትዮጵያ', 'ethiopia', 'ET',
    9.1450000, 40.4896730, 3.4000000, 14.9000000, 32.9900000, 47.9800000,
    ARRAY['Abyssinia', 'Federal Democratic Republic of Ethiopia'], 1
) ON CONFLICT (slug) DO UPDATE SET
    name_en = EXCLUDED.name_en,
    name_am = EXCLUDED.name_am,
    name_om = EXCLUDED.name_om,
    name_ti = EXCLUDED.name_ti;

-- 2. Chartered Cities & Regions (parent = ethiopia)
INSERT INTO location_hierarchy (
    parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
    latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order
)
SELECT p.id, d.location_type, d.name_en, d.name_am, d.name_om, d.name_ti, d.slug, d.admin_code,
       d.latitude, d.longitude, d.min_lat, d.max_lat, d.min_lng, d.max_lng, d.aliases, d.display_order
FROM (SELECT id FROM location_hierarchy WHERE slug = 'ethiopia') p
CROSS JOIN (VALUES
    ('city', 'Addis Ababa', 'አዲስ አበባ', 'Finfinnee', 'ኣዲስ ኣበባ', 'addis-ababa', 'ET-AA', 9.0320000, 38.7468900, 8.8300000, 9.1000000, 38.6500000, 38.9000000, ARRAY['Finfinne', 'Addis Abeba', 'Sheger', 'AA'], 1),
    ('city', 'Dire Dawa', 'ድሬዳዋ', 'Dirree Dawaa', 'ድሬዳዋ', 'dire-dawa', 'ET-DD', 9.5930000, 41.8660000, 9.4500000, 9.7500000, 41.6500000, 42.0500000, ARRAY['Dire Dawa City', 'Dirree Dhawaa'], 2),
    ('region', 'Oromia', 'ኦሮሚያ', 'Oromiyaa', 'ኦሮሚያ', 'oromia', 'ET-OR', 8.5410000, 39.2680000, 3.4000000, 10.4000000, 34.1000000, 42.9000000, ARRAY['Oromiya', 'Oromiyaa'], 3),
    ('region', 'Amhara', 'አማራ', 'Amaara', 'ኣምሓራ', 'amhara', 'ET-AM', 11.5900000, 37.3900000, 8.8000000, 13.8000000, 35.2000000, 40.3000000, ARRAY['Amhara Region', 'Amhara National Regional State'], 4),
    ('region', 'Tigray', 'ትግራይ', 'Tigraay', 'ትግራይ', 'tigray', 'ET-TI', 14.0300000, 38.3160000, 12.2000000, 14.9000000, 36.4500000, 40.0000000, ARRAY['Tigrai', 'Tigre'], 5),
    ('region', 'Sidama', 'ሲዳማ', 'Sidaama', 'ሲዳማ', 'sidama', 'ET-SI', 6.7800000, 38.4500000, 6.2000000, 7.3000000, 38.0000000, 39.2000000, ARRAY['Sidama Region', 'Sidaama'], 6),
    ('region', 'Somali', 'ሶማሌ', 'Somaalee', 'ሶማሌ', 'somali', 'ET-SO', 7.1500000, 44.2500000, 4.0000000, 11.0000000, 40.5000000, 48.0000000, ARRAY['Somali Region', 'Ogaden'], 7),
    ('region', 'Afar', 'አፋር', 'Afaar', 'ዓፋር', 'afar', 'ET-AF', 11.7500000, 40.9500000, 8.8000000, 14.5000000, 39.5000000, 42.5000000, ARRAY['Afar Region', 'Danakil'], 8),
    ('region', 'Benishangul-Gumuz', 'ቤንሻንጉል ጉሙዝ', 'Beniishaangul Gumuz', 'ቤኒሻንጉል ጉሙዝ', 'benishangul-gumuz', 'ET-BE', 10.5000000, 35.5000000, 9.2000000, 12.0000000, 34.1000000, 37.0000000, ARRAY['Benishangul', 'Gumuz'], 9),
    ('region', 'Gambela', 'ጋምቤላ', 'Gambellaa', 'ጋምቤላ', 'gambela', 'ET-GA', 7.9000000, 34.5000000, 6.5000000, 8.8000000, 33.0000000, 35.5000000, ARRAY['Gambella', 'Gambella Region'], 10),
    ('region', 'Harari', 'ሐረሪ', 'Hararii', 'ሓረሪ', 'harari', 'ET-HA', 9.3100000, 42.1300000, 9.1500000, 9.4500000, 42.0000000, 42.3000000, ARRAY['Harar Region', 'Harari People'], 11),
    ('region', 'South Ethiopia', 'ደቡብ ኢትዮጵያ', 'Kibba Itoophiyaa', 'ደቡብ ኢትዮጵያ', 'south-ethiopia', 'ET-SE', 6.0000000, 37.5000000, 4.5000000, 7.5000000, 35.8000000, 38.5000000, ARRAY['Wolayita', 'Gamo', 'Gofa', 'South Region'], 12),
    ('region', 'Central Ethiopia', 'ማዕከላዊ ኢትዮጵያ', 'Giddugala Itoophiyaa', 'ማእከላይ ኢትዮጵያ', 'central-ethiopia', 'ET-CE', 7.8000000, 38.0000000, 7.0000000, 8.5000000, 37.0000000, 39.0000000, ARRAY['Gurage', 'Silte', 'Hadiya', 'Kambata'], 13),
    ('region', 'Southwest Ethiopia', 'ደቡብ ምዕራብ ኢትዮጵያ', 'Kibba-Lixa Itoophiyaa', 'ደቡብ ምዕራብ ኢትዮጵያ', 'southwest-ethiopia', 'ET-SW', 7.1500000, 36.2500000, 5.8000000, 8.2000000, 35.0000000, 37.2000000, ARRAY['Kaffa', 'Sheka', 'Bench Sheko', 'Dawro', 'West Omo'], 14)
) AS d(location_type, name_en, name_am, name_om, name_ti, slug, admin_code, latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order)
ON CONFLICT (slug) DO UPDATE SET
    name_en = EXCLUDED.name_en,
    name_am = EXCLUDED.name_am,
    name_om = EXCLUDED.name_om,
    name_ti = EXCLUDED.name_ti;

-- 3. All 11 Addis Ababa Subcities (parent = addis-ababa)
INSERT INTO location_hierarchy (
    parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
    latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
)
SELECT p.id, 'subcity', d.name_en, d.name_am, d.name_om, d.name_ti, d.slug, d.admin_code,
       d.latitude, d.longitude, d.min_lat, d.max_lat, d.min_lng, d.max_lng, d.aliases, d.display_order, true
FROM (SELECT id FROM location_hierarchy WHERE slug = 'addis-ababa') p
CROSS JOIN (VALUES
    ('Bole', 'ቦሌ', 'Boolee', 'ቦሌ', 'addis-ababa-bole', 'AA-BO', 8.9950000, 38.7900000, 8.9400000, 9.0400000, 38.7500000, 38.8600000, ARRAY['Bole Subcity', 'Bole Airport', 'Bole Medhanealem', 'Bole Atlas', 'Megenagna'], 1),
    ('Kirkos', 'ቂርቆስ', 'Qirqoos', 'ቂርቆስ', 'addis-ababa-kirkos', 'AA-KI', 9.0080000, 38.7580000, 8.9900000, 9.0300000, 38.7400000, 38.7800000, ARRAY['Kirkos Subcity', 'Kazanchis', 'Meshualekia', 'Meskel Square', 'Olympia'], 2),
    ('Yeka', 'የካ', 'Yakkkaa', 'የካ', 'addis-ababa-yeka', 'AA-YE', 9.0420000, 38.8050000, 9.0100000, 9.0900000, 38.7700000, 38.8500000, ARRAY['Yeka Subcity', 'Kotebe', 'Ayat', 'CMC', 'Ferensay'], 3),
    ('Arada', 'አራዳ', 'Aradaa', 'ኣራዳ', 'addis-ababa-arada', 'AA-AR', 9.0350000, 38.7520000, 9.0200000, 9.0550000, 38.7350000, 38.7700000, ARRAY['Arada Subcity', 'Piazza', '4 Kilo', '6 Kilo', 'Amist Kilo'], 4),
    ('Lideta', 'ልደታ', 'Lidetaa', 'ልደታ', 'addis-ababa-lideta', 'AA-LI', 9.0120000, 38.7350000, 8.9950000, 9.0300000, 38.7150000, 38.7500000, ARRAY['Lideta Subcity', 'Abanet', 'Sengatera', 'Mexico'], 5),
    ('Addis Ketema', 'አዲስ ከተማ', 'Magaalaa Haaraa', 'ኣዲስ ከተማ', 'addis-ababa-addis-ketema', 'AA-AK', 9.0330000, 38.7300000, 9.0150000, 9.0500000, 38.7100000, 38.7450000, ARRAY['Merkato', 'Addis Ketema Subcity', 'Autobus Tera'], 6),
    ('Gullele', 'ጉለሌ', 'Gullallee', 'ጉለሌ', 'addis-ababa-gullele', 'AA-GU', 9.0680000, 38.7350000, 9.0400000, 9.1100000, 38.7000000, 38.7650000, ARRAY['Gullele Subcity', 'Shiro Meda', 'Entoto', 'Addisu Gebeya'], 7),
    ('Kolfe Keranio', 'ኮልፌ ቀራኒዮ', 'Kolfi Qaraaniyoo', 'ኮልፌ ቀራንዮ', 'addis-ababa-kolfe-keranio', 'AA-KK', 9.0200000, 38.7050000, 8.9750000, 9.0650000, 38.6650000, 38.7300000, ARRAY['Kolfe', 'Keranio', 'Tor Hailoch', 'Ayer Tena'], 8),
    ('Akaki Kality', 'አቃቂ ቃሊቲ', 'Aqaaqii Qaallittii', 'ኣቃቂ ቃሊቲ', 'addis-ababa-akaki-kality', 'AA-AQ', 8.9100000, 38.7650000, 8.8400000, 8.9700000, 38.7100000, 38.8300000, ARRAY['Akaki', 'Kality', 'Tulu Dimtu', 'Koye Feche'], 9),
    ('Nifas Silk-Lafto', 'ንፋስ ስልክ ላፍቶ', 'Nifas Silk Laaftoo', 'ንፋስ ስልክ ላፍቶ', 'addis-ababa-nifas-silk-lafto', 'AA-NL', 8.9700000, 38.7350000, 8.9250000, 9.0150000, 38.6900000, 38.7650000, ARRAY['Nefas Silk', 'Lafto', 'Jemo', 'Lebu', 'Mekanisa', 'Gotera'], 10),
    ('Lemi Kura', 'ለሚ ኩራ', 'Lami Kuraa', 'ለሚ ኩራ', 'addis-ababa-lemi-kura', 'AA-LK', 9.0150000, 38.8400000, 8.9700000, 9.0700000, 38.8000000, 38.9000000, ARRAY['Lemi Kura', 'Goro', 'Bole Arabsa', 'Summit'], 11)
) AS d(name_en, name_am, name_om, name_ti, slug, admin_code, latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order)
ON CONFLICT (slug) DO UPDATE SET
    name_en = EXCLUDED.name_en,
    name_am = EXCLUDED.name_am,
    name_om = EXCLUDED.name_om,
    name_ti = EXCLUDED.name_ti;

-- 4. Key Metropolitan Cities & Regional Hubs (parent joined by parent_slug)
INSERT INTO location_hierarchy (
    parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
    latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order, is_leaf
)
SELECT p.id, 'city', d.name_en, d.name_am, d.name_om, d.name_ti, d.slug, d.admin_code,
       d.latitude, d.longitude, d.min_lat, d.max_lat, d.min_lng, d.max_lng, d.aliases, d.display_order, true
FROM (VALUES
    ('oromia', 'Adama', 'አዳማ', 'Adaamaa', 'ኣዳማ', 'oromia-adama', 'OR-ADA', 8.5400000, 39.2700000, 8.4800000, 8.6000000, 39.2000000, 39.3400000, ARRAY['Nazret', 'Nazareth'], 1),
    ('oromia', 'Bishoftu', 'ቢሾፍቱ', 'Bishooftuu', 'ቢሾፍቱ', 'oromia-bishoftu', 'OR-BIS', 8.7500000, 38.9800000, 8.7000000, 8.8000000, 38.9200000, 39.0400000, ARRAY['Debre Zeyit', 'Debre Zeit'], 2),
    ('oromia', 'Jimma', 'ጅማ', 'Jimmaa', 'ጅማ', 'oromia-jimma', 'OR-JIM', 7.6700000, 36.8300000, 7.6200000, 7.7300000, 36.7800000, 36.9000000, ARRAY['Jima'], 3),
    ('oromia', 'Shashemene', 'ሻሸመኔ', 'Shaashamannee', 'ሻሸመኔ', 'oromia-shashemene', 'OR-SHA', 7.2000000, 38.6000000, 7.1400000, 7.2500000, 38.5400000, 38.6500000, ARRAY['Shashamane'], 4),
    ('amhara', 'Bahir Dar', 'ባሕር ዳር', 'Baahir Daar', 'ባሕር ዳር', 'amhara-bahir-dar', 'AM-BD', 11.5900000, 37.3900000, 11.5300000, 11.6500000, 37.3200000, 37.4500000, ARRAY['Bahr Dar'], 1),
    ('amhara', 'Gondar', 'ጎንደር', 'Gondar', 'ጎንደር', 'amhara-gondar', 'AM-GD', 12.6000000, 37.4600000, 12.5500000, 12.6600000, 37.4000000, 37.5200000, ARRAY['Gonder'], 2),
    ('amhara', 'Dessie', 'ደሴ', 'Dasee', 'ደሴ', 'amhara-dessie', 'AM-DES', 11.1300000, 39.6300000, 11.0800000, 11.1800000, 39.5800000, 39.6800000, ARRAY['Dese'], 3),
    ('tigray', 'Mekelle', 'መቀሌ', 'Maqalee', 'መቐለ', 'tigray-mekelle', 'TI-MK', 13.4960000, 39.4750000, 13.4400000, 13.5600000, 39.4200000, 39.5400000, ARRAY['Makale', 'Meqele'], 1),
    ('sidama', 'Hawassa', 'ሀዋሳ', 'Hawaasaa', 'ሃዋሳ', 'sidama-hawassa', 'SI-HAW', 7.0500000, 38.4800000, 6.9800000, 7.1200000, 38.4200000, 38.5400000, ARRAY['Awassa'], 1),
    ('somali', 'Jigjiga', 'ጅግጅጋ', 'Jigjigaa', 'ጅግጅጋ', 'somali-jigjiga', 'SO-JIG', 9.3500000, 42.8000000, 9.3000000, 9.4000000, 42.7400000, 42.8600000, ARRAY['Jijiga'], 1)
) AS d(parent_slug, name_en, name_am, name_om, name_ti, slug, admin_code, latitude, longitude, min_lat, max_lat, min_lng, max_lng, aliases, display_order)
JOIN location_hierarchy p ON p.slug = d.parent_slug
ON CONFLICT (slug) DO UPDATE SET
    name_en = EXCLUDED.name_en,
    name_am = EXCLUDED.name_am,
    name_om = EXCLUDED.name_om,
    name_ti = EXCLUDED.name_ti;

-- 5. Backward Compatibility: Populate legacy locations table in batch
INSERT INTO locations (region, subcity, region_am, subcity_am, region_or, subcity_or, region_ti, subcity_ti)
SELECT
    'Addis Ababa', l.name_en,
    'አዲስ አበባ', l.name_am,
    'Finfinnee', l.name_om,
    'ኣዲስ ኣበባ', l.name_ti
FROM location_hierarchy l
JOIN location_hierarchy p ON p.id = l.parent_id
WHERE p.slug = 'addis-ababa' AND l.location_type = 'subcity';

INSERT INTO locations (region, zone, region_am, zone_am, region_or, zone_or, region_ti, zone_ti)
SELECT
    p.name_en, l.name_en,
    p.name_am, l.name_am,
    p.name_om, l.name_om,
    p.name_ti, l.name_ti
FROM location_hierarchy l
JOIN location_hierarchy p ON p.id = l.parent_id
WHERE l.location_type = 'city' AND p.location_type = 'region';

INSERT INTO locations (region, region_am, region_or, region_ti)
SELECT
    l.name_en, l.name_am, l.name_om, l.name_ti
FROM location_hierarchy l
JOIN location_hierarchy p ON p.id = l.parent_id
WHERE p.slug = 'ethiopia' AND l.location_type IN ('region', 'city');
