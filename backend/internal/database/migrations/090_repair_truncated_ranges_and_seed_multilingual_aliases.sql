-- ============================================================================
-- Migration 090: Repair Truncated Range Specifications & Seed Multilingual Aliases
-- ============================================================================
-- 1. Restores the full specification titles for numeric range categories (800W, 100k,
--    115mm, 18V, 2000W, 300m², 500m²) so layers 4 & 5 display full technical specs.
-- 2. Creates unique index on category_aliases for idempotent alias population.
-- 3. Populates category_aliases from classification_aliases for all active categories.
-- 4. Seeds essential Ethiopian commercial aliases across all 4 languages
--    (en, am, om, ti) for instantaneous search discovery.
-- ============================================================================

-- 1. Repair Truncated Numeric Range Specification Names (Layers 4 & 5)
UPDATE category_hierarchy SET name_en = CASE 
    WHEN name_am LIKE '%ሞዴል 1%' OR slug LIKE '%_0' OR slug LIKE '%_0-1' THEN '800W - 1500W High Impact Energy 3-Mode Spec Model 1'
    ELSE '800W - 1500W High Impact Energy 3-Mode Spec Model 2'
END WHERE name_en = '800W';

UPDATE category_hierarchy SET name_en = CASE 
    WHEN name_am LIKE '%ሞዴል 1%' OR slug LIKE '%_0' OR slug LIKE '%_0-1' THEN '100k - 500k Verified Followers High Engagement Spec Model 1'
    ELSE '100k - 500k Verified Followers High Engagement Spec Model 2'
END WHERE name_en = '100k';

UPDATE category_hierarchy SET name_en = CASE 
    WHEN name_am LIKE '%ክፍል 1%' OR slug LIKE '%_0' OR slug LIKE '%_0-1' THEN '115mm - 230mm Heavy Angle Grinder Spec Part 1'
    ELSE '115mm - 230mm Heavy Angle Grinder Spec Part 2'
END WHERE name_en = '115mm';

UPDATE category_hierarchy SET name_en = CASE 
    WHEN name_am LIKE '%ክፍል 1%' OR slug LIKE '%_0' OR slug LIKE '%_0-1' THEN '18V - 20V Max Brushless Impact Wrench Spec Part 1'
    ELSE '18V - 20V Max Brushless Impact Wrench Spec Part 2'
END WHERE name_en = '18V';

UPDATE category_hierarchy SET 
    name_en = CASE 
        WHEN slug LIKE '%_0' OR slug LIKE '%_0-1' THEN '2000W - 5000W Heavy Industrial Hybrid Inverter Spec Part 1'
        ELSE '2000W - 5000W Heavy Industrial Hybrid Inverter Spec Part 2'
    END,
    name_am = CASE 
        WHEN slug LIKE '%_0' OR slug LIKE '%_0-1' THEN 'ከ2000 እስከ 5000 ዋት ከባድ የኢንዱስትሪ ኢንቨርተር ክፍል 1'
        ELSE 'ከ2000 እስከ 5000 ዋት ከባድ የኢንዱስትሪ ኢንቨርተር ክፍል 2'
    END
WHERE name_en = '2000W';

UPDATE category_hierarchy SET name_en = CASE 
    WHEN name_am LIKE '%አይነት 1%' OR slug LIKE '%_0' OR slug LIKE '%_0-1' THEN '300m² - 1000m² Prime Commercial Showroom Spec Model 1'
    ELSE '300m² - 1000m² Prime Commercial Showroom Spec Model 2'
END WHERE name_en = '300m²';

UPDATE category_hierarchy SET name_en = CASE 
    WHEN name_am LIKE '%ክፍል 1%' OR slug LIKE '%_0' OR slug LIKE '%_0-1' THEN '500m² - 2000m² High-Ceiling Clear Span Warehouse Spec Part 1'
    ELSE '500m² - 2000m² High-Ceiling Clear Span Warehouse Spec Part 2'
END WHERE name_en = '500m²';


-- 2. Create Unique Constraint Index on category_aliases
CREATE UNIQUE INDEX IF NOT EXISTS uq_category_aliases_cat_alias_lang 
ON category_aliases (category_id, alias, language);


-- 3. Populate category_aliases from classification_aliases
INSERT INTO category_aliases (category_id, alias, language)
SELECT DISTINCT cp.id, ca.alias, ca.lang
FROM classification_aliases ca
JOIN classification_placements cp ON ca.placement_id = cp.id
JOIN category_hierarchy c ON cp.id = c.id
WHERE c.is_active = true
ON CONFLICT (category_id, alias, language) DO NOTHING;


-- 4. Seed Essential High-Frequency Ethiopian Search Synonyms & Colloquial Aliases
-- Power Drills & Tools (2376)
INSERT INTO category_aliases (category_id, alias, language) VALUES
(2376, 'drill', 'en'), (2376, 'rotary drill', 'en'), (2376, 'hammer drill', 'en'), (2376, 'impact drill', 'en'),
(2376, 'ድሪል', 'am'), (2376, 'ቦረቦሬ', 'am'), (2376, 'የእጅ ድሪል', 'am'), (2376, 'መሰርሰሪያ', 'am'),
(2376, 'dirilii', 'om'), (2376, 'meeshaa qotinsaa', 'om'),
(2376, 'ድሪል', 'ti'), (2376, 'መሳርሒ ምኹዓት', 'ti')
ON CONFLICT (category_id, alias, language) DO NOTHING;

-- Agriculture & Teff (343, 331, 332, 4)
INSERT INTO category_aliases (category_id, alias, language) VALUES
(4, 'teff', 'en'), (4, 'wheat', 'en'), (4, 'barley', 'en'), (4, 'coffee', 'en'),
(4, 'ጤፍ', 'am'), (4, 'ስንዴ', 'am'), (4, 'ገብስ', 'am'), (4, 'ቡና', 'am'), (4, 'ማኛ', 'am'),
(4, 'xaafii', 'om'), (4, 'qamadii', 'om'), (4, 'garbuu', 'om'), (4, 'buna', 'om'),
(4, 'ጣፍ', 'ti'), (4, 'ስርናይ', 'ti'), (4, 'ስገም', 'ti'), (4, 'ቡን', 'ti')
ON CONFLICT (category_id, alias, language) DO NOTHING;

-- Smartphones & Mobile Phones (1425, 9)
INSERT INTO category_aliases (category_id, alias, language) VALUES
(1425, 'phone', 'en'), (1425, 'smartphone', 'en'), (1425, 'mobile', 'en'), (1425, 'cellphone', 'en'),
(1425, 'ስልክ', 'am'), (1425, 'ሞባይል', 'am'), (1425, 'ስማርትፎን', 'am'),
(1425, 'bilbila', 'om'), (1425, 'moobaayila', 'om'),
(1425, 'ሞባይል', 'ti'), (1425, 'ቴሌፎን', 'ti'),
(9, 'phone', 'en'), (9, 'smartphone', 'en'), (9, 'mobile', 'en'),
(9, 'ስልክ', 'am'), (9, 'ሞባይል', 'am'),
(9, 'bilbila', 'om'), (9, 'moobaayila', 'om'),
(9, 'ሞባይል', 'ti'), (9, 'ቴሌፎን', 'ti')
ON CONFLICT (category_id, alias, language) DO NOTHING;

-- Laptops & Computers (10, 1130)
INSERT INTO category_aliases (category_id, alias, language) VALUES
(10, 'laptop', 'en'), (10, 'computer', 'en'), (10, 'notebook', 'en'), (10, 'pc', 'en'),
(10, 'ላፕቶፕ', 'am'), (10, 'ኮምፒውተር', 'am'), (10, 'ዴስክቶፕ', 'am'),
(10, 'laaptooppii', 'om'), (10, 'kompiitara', 'om'),
(10, 'ላፕቶፕ', 'ti'), (10, 'ኮምፒተር', 'ti'),
(1130, 'laptop', 'en'), (1130, 'computer', 'en'), (1130, 'it equipment', 'en'),
(1130, 'ላፕቶፕ', 'am'), (1130, 'ኮምፒውተር', 'am'),
(1130, 'laaptooppii', 'om'), (1130, 'kompiitara', 'om')
ON CONFLICT (category_id, alias, language) DO NOTHING;

-- Vehicles & Automotive (3197, 3243)
INSERT INTO category_aliases (category_id, alias, language) VALUES
(3197, 'car', 'en'), (3197, 'vehicle', 'en'), (3197, 'automobile', 'en'), (3197, 'corolla', 'en'), (3197, 'toyota', 'en'), (3197, 'bajaj', 'en'),
(3197, 'መኪና', 'am'), (3197, 'ተሽከርካሪ', 'am'), (3197, 'ባጃጅ', 'am'), (3197, 'ኮሮላ', 'am'), (3197, 'ቶዮታ', 'am'),
(3197, 'konkolaataa', 'om'), (3197, 'bajaajii', 'om'),
(3197, 'መኪና', 'ti'), (3197, 'ባጃጅ', 'ti')
ON CONFLICT (category_id, alias, language) DO NOTHING;

-- Property & Real Estate (2758)
INSERT INTO category_aliases (category_id, alias, language) VALUES
(2758, 'house', 'en'), (2758, 'home', 'en'), (2758, 'apartment', 'en'), (2758, 'condominium', 'en'), (2758, 'villa', 'en'), (2758, 'rent', 'en'),
(2758, 'ቤት', 'am'), (2758, 'አፓርታማ', 'am'), (2758, 'ኮንዶሚኒየም', 'am'), (2758, 'ቪላ', 'am'), (2758, 'ኪራይ', 'am'), (2758, 'ሽያጭ', 'am'),
(2758, 'mana', 'om'), (2758, 'kondoominiyeemii', 'om'), (2758, 'kireeffannaa', 'om'),
(2758, 'ገዛ', 'ti'), (2758, 'ኮንዶሚኒየም', 'ti'), (2758, 'ክራይ', 'ti')
ON CONFLICT (category_id, alias, language) DO NOTHING;

-- Construction & Cement (1201, 1212)
INSERT INTO category_aliases (category_id, alias, language) VALUES
(1201, 'cement', 'en'), (1201, 'rebar', 'en'), (1201, 'iron', 'en'), (1201, 'construction', 'en'),
(1201, 'ሲሚንቶ', 'am'), (1201, 'ብረት', 'am'), (1201, 'ግንባታ', 'am'), (1201, 'ቆርቆሮ', 'am'),
(1201, 'siminto', 'om'), (1201, 'sibiila', 'om'),
(1201, 'ሲሚንቶ', 'ti'), (1201, 'ብረት', 'ti')
ON CONFLICT (category_id, alias, language) DO NOTHING;

