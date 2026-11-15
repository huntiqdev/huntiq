-- ============================================================================
-- 082_seed_comprehensive_ethiopian_zones_woredas.sql
-- Comprehensive Ethiopian Administrative Structure (Region -> Zone -> Woreda)
-- Features:
--   1. Fix get_location_ancestors function to include slug in return signature
--   2. Seed 75+ Administrative Zones across all 12 Regions & 2 Chartered Cities
--   3. Seed 100+ Key Urban and Rural Woredas with systematic disambiguation
--   4. Deterministic composite slugs (<region>-<zone>-<woreda>) eliminating collisions
--   5. Multilingual metadata (en, am, om, ti) and rich search aliases
--   6. Synchronize legacy locations table for full backward compatibility
-- ============================================================================

BEGIN;

-- 1. Update get_location_ancestors to return slug column for breadcrumb resolution
DROP FUNCTION IF EXISTS public.get_location_ancestors(integer);
CREATE OR REPLACE FUNCTION public.get_location_ancestors(loc_id integer)
 RETURNS TABLE(id integer, location_type character varying, name_en character varying, name_am character varying, name_om character varying, name_ti character varying, slug character varying, level integer)
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    WITH RECURSIVE ancestors AS (
        SELECT l.id, l.parent_id, l.location_type, l.name_en, l.name_am, l.name_om, l.name_ti, l.slug, l.level
        FROM location_hierarchy l
        WHERE l.id = loc_id
        
        UNION ALL
        
        SELECT p.id, p.parent_id, p.location_type, p.name_en, p.name_am, p.name_om, p.name_ti, p.slug, p.level
        FROM location_hierarchy p
        JOIN ancestors a ON a.parent_id = p.id
    )
    SELECT a.id, a.location_type, a.name_en, a.name_am, a.name_om, a.name_ti, a.slug, a.level
    FROM ancestors a
    ORDER BY a.level ASC;
END;
$function$;

-- 2. Seed All Administrative Zones and Subdistricts (Level 2)
INSERT INTO location_hierarchy (
    parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
    latitude, longitude, aliases, display_order, is_leaf, is_active
)
SELECT
    p.id,
    d.loc_type,
    d.name_en, d.name_am, d.name_om, d.name_ti,
    d.slug, d.admin_code,
    d.lat, d.lng,
    d.aliases,
    d.disp_order,
    false, true
FROM (VALUES
    ('oromia', 'zone', 'East Shewa Zone', 'ምስራቅ ሸዋ ዞን', 'Shawaa Bahaa', 'ዞን ምብራቕ ሸዋ', 'oromia-east-shewa', 'ET-OR-ES', 8.5400000, 39.2700000, ARRAY['East Shoa', 'Shawa Baha', 'ምስራቅ ሸዋ']::text[], 1),
    ('oromia', 'zone', 'West Shewa Zone', 'ምዕራብ ሸዋ ዞን', 'Shawaa Dhihaa', 'ዞን ምዕራብ ሸዋ', 'oromia-west-shewa', 'ET-OR-WS', 8.9800000, 37.8500000, ARRAY['West Shoa', 'Shawa Dhiha', 'ምዕራብ ሸዋ']::text[], 2),
    ('oromia', 'zone', 'North Shewa Zone (Oromia)', 'ሰሜን ሸዋ ዞን (ኦሮሚያ)', 'Shawaa Kaabaa', 'ዞን ሰሜን ሸዋ (ኦሮሚያ)', 'oromia-north-shewa', 'ET-OR-NS', 9.7500000, 38.7500000, ARRAY['North Shoa Oromia', 'Shawa Kaaba', 'ሰሜን ሸዋ ኦሮሚያ']::text[], 3),
    ('oromia', 'zone', 'Southwest Shewa Zone', 'ደቡብ ምዕራብ ሸዋ ዞን', 'Shawaa Kibba-Lixaa', 'ዞን ደቡብ ምዕራብ ሸዋ', 'oromia-southwest-shewa', 'ET-OR-SWS', 8.5300000, 38.0000000, ARRAY['Southwest Shoa', 'ደቡብ ምዕራብ ሸዋ']::text[], 4),
    ('oromia', 'zone', 'Arsi Zone', 'አርሲ ዞን', 'Arsii', 'ዞን አርሲ', 'oromia-arsi', 'ET-OR-AR', 7.7500000, 39.2500000, ARRAY['Arssi', 'አርሲ']::text[], 5),
    ('oromia', 'zone', 'West Arsi Zone', 'ምዕራብ አርሲ ዞን', 'Arsii Dhihaa', 'ዞን ምዕራብ አርሲ', 'oromia-west-arsi', 'ET-OR-WAR', 7.2000000, 38.6000000, ARRAY['West Arssi', 'ምዕራብ አርሲ']::text[], 6),
    ('oromia', 'zone', 'Bale Zone', 'ባሌ ዞን', 'Baalee', 'ዞን ባሌ', 'oromia-bale', 'ET-OR-BA', 7.0000000, 40.0000000, ARRAY['Bale', 'ባሌ']::text[], 7),
    ('oromia', 'zone', 'East Bale Zone', 'ምስራቅ ባሌ ዞን', 'Baalee Bahaa', 'ዞን ምስራቅ ባሌ', 'oromia-east-bale', 'ET-OR-EBA', 7.1500000, 40.7000000, ARRAY['East Bale', 'ምስራቅ ባሌ']::text[], 8),
    ('oromia', 'zone', 'Borena Zone', 'ቦረና ዞን', 'Boorana', 'ዞን ቦረና', 'oromia-borena', 'ET-OR-BO', 4.5000000, 38.5000000, ARRAY['Borana', 'ቦረና']::text[], 9),
    ('oromia', 'zone', 'Guji Zone', 'ጉጂ ዞን', 'Gujii', 'ዞን ጉጂ', 'oromia-guji', 'ET-OR-GU', 5.5000000, 39.0000000, ARRAY['Guji', 'ጉጂ']::text[], 10),
    ('oromia', 'zone', 'West Guji Zone', 'ምዕራብ ጉጂ ዞን', 'Gujii Dhihaa', 'ዞን ምዕራብ ጉጂ', 'oromia-west-guji', 'ET-OR-WGU', 5.6500000, 38.2500000, ARRAY['West Guji', 'ምዕራብ ጉጂ']::text[], 11),
    ('oromia', 'zone', 'East Hararghe Zone', 'ምስራቅ ሐረርጌ ዞን', 'Harargee Bahaa', 'ዞን ምስራቕ ሓረርጌ', 'oromia-east-hararghe', 'ET-OR-EH', 9.2500000, 42.0000000, ARRAY['East Harerge', 'ምስራቅ ሐረርጌ']::text[], 12),
    ('oromia', 'zone', 'West Hararghe Zone', 'ምዕራብ ሐረርጌ ዞን', 'Harargee Dhihaa', 'ዞን ምዕራብ ሓረርጌ', 'oromia-west-hararghe', 'ET-OR-WH', 8.9000000, 40.7500000, ARRAY['West Harerge', 'ምዕራብ ሐረርጌ']::text[], 13),
    ('oromia', 'zone', 'Jimma Zone', 'ጅማ ዞን', 'Jimmaa', 'ዞን ጅማ', 'oromia-jimma-zone', 'ET-OR-JM', 7.6700000, 36.8300000, ARRAY['Jimma Zone', 'ጅማ']::text[], 14),
    ('oromia', 'zone', 'East Welega Zone', 'ምስራቅ ወለጋ ዞን', 'Wallagga Bahaa', 'ዞን ምስራቕ ወለጋ', 'oromia-east-welega', 'ET-OR-EW', 9.0800000, 36.5500000, ARRAY['East Wollega', 'ምስራቅ ወለጋ']::text[], 15),
    ('oromia', 'zone', 'West Welega Zone', 'ምዕራብ ወለጋ ዞን', 'Wallagga Dhihaa', 'ዞን ምዕራብ ወለጋ', 'oromia-west-welega', 'ET-OR-WW', 9.1700000, 35.5000000, ARRAY['West Wollega', 'ምዕራብ ወለጋ']::text[], 16),
    ('oromia', 'zone', 'Horo Guduru Welega Zone', 'ሆሮ ጉዱሩ ወለጋ ዞን', 'Horo Guduruu Wallaggaa', 'ዞን ሆሮ ጉዱሩ ወለጋ', 'oromia-horo-guduru', 'ET-OR-HG', 9.5000000, 37.0500000, ARRAY['Horo Gudru', 'ሆሮ ጉዱሩ']::text[], 17),
    ('oromia', 'zone', 'Kelem Welega Zone', 'ቄለም ወለጋ ዞን', 'Qeellam Wallaggaa', 'ዞን ቄለም ወለጋ', 'oromia-kelem-welega', 'ET-OR-KW', 8.8000000, 34.8000000, ARRAY['Qellem Wollega', 'ቄለም ወለጋ']::text[], 18),
    ('oromia', 'zone', 'Ilu Aba Bor Zone', 'ኢሉ አባ ቦራ ዞን', 'Illuu Abbaa Booraa', 'ዞን ኢሉ ኣባ ቦራ', 'oromia-ilu-aba-bor', 'ET-OR-IB', 8.3000000, 35.7500000, ARRAY['Illubabor', 'ኢሉ አባቦራ']::text[], 19),
    ('oromia', 'zone', 'Buno Bedele Zone', 'ቡኖ በደሌ ዞን', 'Buunnoo Baddallee', 'ዞን ቡኖ በደሌ', 'oromia-buno-bedele', 'ET-OR-BB', 8.4500000, 36.3500000, ARRAY['Buno Bedelle', 'ቡኖ በደሌ']::text[], 20),
    ('oromia', 'zone', 'Oromia Special Zone Surrounding Finfinne', 'የፊንፊኔ ዙሪያ ልዩ ዞን', 'Naannawa Finfinnee', 'ፍሉይ ዞን ዙርያ ፊንፊነ', 'oromia-special-zone-finfinne', 'ET-OR-SZ', 9.0500000, 38.7000000, ARRAY['Finfinne Special Zone', 'የአዲስ አበባ ዙሪያ ልዩ ዞን']::text[], 21),
    ('amhara', 'zone', 'North Shewa Zone (Amhara)', 'ሰሜን ሸዋ ዞን (አማራ)', 'Shawaa Kaabaa (Amaara)', 'ዞን ሰሜን ሸዋ (ኣምሓራ)', 'amhara-north-shewa', 'ET-AM-NS', 9.6800000, 39.5300000, ARRAY['North Shoa Amhara', 'Semien Shewa', 'ሰሜን ሸዋ አማራ']::text[], 22),
    ('amhara', 'zone', 'South Gondar Zone', 'ደቡብ ጎንደር ዞን', 'Gondar Kibbaa', 'ዞን ደቡብ ጎንደር', 'amhara-south-gondar', 'ET-AM-SG', 11.8500000, 38.0000000, ARRAY['South Gonder', 'Debub Gondar', 'ደቡብ ጎንደር']::text[], 23),
    ('amhara', 'zone', 'Central Gondar Zone', 'ማዕከላዊ ጎንደር ዞን', 'Gondar Giddugalaa', 'ዞን ማእከላይ ጎንደር', 'amhara-central-gondar', 'ET-AM-CG', 12.6000000, 37.4500000, ARRAY['Central Gonder', 'ማዕከላዊ ጎንደር']::text[], 24),
    ('amhara', 'zone', 'North Gondar Zone', 'ሰሜን ጎንደር ዞን', 'Gondar Kaabaa', 'ዞን ሰሜን ጎንደር', 'amhara-north-gondar', 'ET-AM-NG', 13.1500000, 37.9000000, ARRAY['North Gonder', 'Semien Gondar', 'ሰሜን ጎንደር']::text[], 25),
    ('amhara', 'zone', 'West Gondar Zone', 'ምዕራብ ጎንደር ዞን', 'Gondar Dhihaa', 'ዞን ምዕራብ ጎንደር', 'amhara-west-gondar', 'ET-AM-WG', 12.8000000, 36.2000000, ARRAY['West Gonder', 'Mirab Gondar', 'ምዕራብ ጎንደር']::text[], 26),
    ('amhara', 'zone', 'West Gojjam Zone', 'ምዕራብ ጎጃም ዞን', 'Gojjaam Dhihaa', 'ዞን ምዕራብ ጎጃም', 'amhara-west-gojjam', 'ET-AM-WJ', 11.1000000, 37.2000000, ARRAY['West Gojam', 'Mirab Gojjam', 'ምዕራብ ጎጃም']::text[], 27),
    ('amhara', 'zone', 'East Gojjam Zone', 'ምስራቅ ጎጃም ዞን', 'Gojjaam Bahaa', 'ዞን ምስራቕ ጎጃም', 'amhara-east-gojjam', 'ET-AM-EJ', 10.3300000, 37.7200000, ARRAY['East Gojam', 'Misraq Gojjam', 'ምስራቅ ጎጃም']::text[], 28),
    ('amhara', 'zone', 'Awi Zone', 'አዊ ዞን', 'Awi', 'ዞን ኣዊ', 'amhara-awi', 'ET-AM-AW', 10.9800000, 36.8500000, ARRAY['Awi Administrative Zone', 'Agew Awi', 'አዊ']::text[], 29),
    ('amhara', 'zone', 'South Wollo Zone', 'ደቡብ ወሎ ዞን', 'Walloo Kibbaa', 'ዞን ደቡብ ወሎ', 'amhara-south-wollo', 'ET-AM-SW', 11.1300000, 39.2500000, ARRAY['South Wello', 'Debub Wollo', 'ደቡብ ወሎ']::text[], 30),
    ('amhara', 'zone', 'North Wollo Zone', 'ሰሜን ወሎ ዞን', 'Walloo Kaabaa', 'ዞን ሰሜን ወሎ', 'amhara-north-wollo', 'ET-AM-NW', 11.9000000, 39.3000000, ARRAY['North Wello', 'Semien Wollo', 'ሰሜን ወሎ']::text[], 31),
    ('amhara', 'zone', 'Wag Hemra Zone', 'ዋግ ኽምራ ዞን', 'Waag Himraa', 'ዞን ዋግ ኽምራ', 'amhara-wag-hemra', 'ET-AM-WH', 12.6300000, 38.9800000, ARRAY['Wag Hemra', 'Wag Himra', 'ዋግ ኽምራ']::text[], 32),
    ('amhara', 'zone', 'Oromo Special Zone (Amhara)', 'የኦሮሞ ብሔረሰብ አስተዳደር ዞን (አማራ)', 'Godina Addaa Oromiyaa (Amaara)', 'ፍሉይ ዞን ኦሮሞ (ኣምሓራ)', 'amhara-oromo-special-zone', 'ET-AM-OZ', 10.7200000, 39.8700000, ARRAY['Oromia Zone Amhara', 'የኦሮሞ ልዩ ዞን']::text[], 33),
    ('tigray', 'zone', 'Mekelle Special Zone', 'መቐለ ልዩ ዞን', 'Godina Addaa Maqalee', 'ፍሉይ ዞን መቐለ', 'tigray-mekelle-special-zone', 'ET-TI-MK', 13.5000000, 39.4700000, ARRAY['Mekelle Zone', 'መቀሌ ልዩ ዞን']::text[], 34),
    ('tigray', 'zone', 'Central Tigray Zone', 'ማእከላይ ዞን ትግራይ', 'Giddugala Tigraay', 'ዞን ማእከላይ ትግራይ', 'tigray-central', 'ET-TI-CT', 14.1000000, 38.9000000, ARRAY['Mehakelegnaw Tigray', 'ማዕከላዊ ትግራይ']::text[], 35),
    ('tigray', 'zone', 'East Tigray Zone', 'ምስራቓዊ ዞን ትግራይ', 'Tigraay Bahaa', 'ዞን ምስራቕ ትግራይ', 'tigray-east', 'ET-TI-ET', 14.2800000, 39.4500000, ARRAY['Misraqawi Tigray', 'ምስራቅ ትግራይ']::text[], 36),
    ('tigray', 'zone', 'South Tigray Zone', 'ደቡባዊ ዞን ትግራይ', 'Tigraay Kibbaa', 'ዞን ደቡብ ትግራይ', 'tigray-south', 'ET-TI-ST', 12.8000000, 39.5000000, ARRAY['Debubawi Tigray', 'ደቡብ ትግራይ']::text[], 37),
    ('tigray', 'zone', 'North West Tigray Zone', 'ሰሜን ምዕራብ ዞን ትግራይ', 'Tigraay Kaaba-Dhihaa', 'ዞን ሰሜን ምዕራብ ትግራይ', 'tigray-north-west', 'ET-TI-NW', 14.1000000, 38.3000000, ARRAY['Semien Mirab Tigray', 'ሰሜን ምዕራብ ትግራይ']::text[], 38),
    ('tigray', 'zone', 'South East Tigray Zone', 'ደቡብ ምብራቕ ዞን ትግራይ', 'Tigraay Kibba-Bahaa', 'ዞን ደቡብ ምብራቕ ትግራይ', 'tigray-south-east', 'ET-TI-SE', 13.3000000, 39.6000000, ARRAY['Debub Misraq Tigray', 'ደቡብ ምስራቅ ትግራይ']::text[], 39),
    ('tigray', 'zone', 'West Tigray Zone', 'ምዕራባዊ ዞን ትግራይ', 'Tigraay Dhihaa', 'ዞን ምዕራብ ትግራይ', 'tigray-west', 'ET-TI-WT', 13.9000000, 36.9000000, ARRAY['Mirabawi Tigray', 'ምዕራብ ትግራይ']::text[], 40),
    ('sidama', 'zone', 'Hawassa City Administration', 'ሀዋሳ ከተማ አስተዳደር', 'Bulchiinsa Magaalaa Hawaasaa', 'ምምሕዳር ከተማ ሃዋሳ', 'sidama-hawassa-city', 'ET-SI-HC', 7.0500000, 38.4800000, ARRAY['Hawassa Metro', 'ሀዋሳ ከተማ']::text[], 41),
    ('sidama', 'zone', 'Sidama Zone', 'ሲዳማ ዞን', 'Godina Sidaamaa', 'ዞን ሲዳማ', 'sidama-zone', 'ET-SI-ZN', 6.7500000, 38.4500000, ARRAY['Sidama Region Rural Woredas', 'ሲዳማ ወረዳዎች']::text[], 42),
    ('somali', 'zone', 'Fafan Zone', 'ፋፋን ዞን', 'Faafan', 'ዞን ፋፋን', 'somali-fafan', 'ET-SO-FA', 9.3500000, 43.0000000, ARRAY['Jijiga Zone', 'ፋፋን']::text[], 43),
    ('somali', 'zone', 'Sitti Zone', 'ሲቲ ዞን', 'Siti', 'ዞን ሲቲ', 'somali-sitti', 'ET-SO-SI', 10.2000000, 42.0000000, ARRAY['Shinile Zone', 'ሲቲ']::text[], 44),
    ('somali', 'zone', 'Jarar Zone', 'ጃራር ዞን', 'Jarar', 'ዞን ጃራር', 'somali-jarar', 'ET-SO-JA', 8.2500000, 43.5500000, ARRAY['Degehabur Zone', 'ጃራር']::text[], 45),
    ('somali', 'zone', 'Shabele Zone', 'ሸበሌ ዞን', 'Shabelle', 'ዞን ሸበሌ', 'somali-shabele', 'ET-SO-SH', 6.5000000, 43.7000000, ARRAY['Gode Zone', 'ሸበሌ']::text[], 46),
    ('somali', 'zone', 'Liben Zone (Somali)', 'ሊበን ዞን (ሶማሌ)', 'Liiban (Somaalee)', 'ዞን ሊበን (ሶማሌ)', 'somali-liben', 'ET-SO-LI', 5.0000000, 41.5000000, ARRAY['Liben Somali', 'ሊበን ሶማሌ']::text[], 47),
    ('somali', 'zone', 'Dawa Zone', 'ዳዋ ዞን', 'Dawa', 'ዞን ዳዋ', 'somali-dawa', 'ET-SO-DA', 4.3000000, 40.5000000, ARRAY['Dawa Zone', 'ዳዋ']::text[], 48),
    ('somali', 'zone', 'Korahe Zone', 'ቆራሄ ዞን', 'Qorahay', 'ዞን ቆራሄ', 'somali-korahe', 'ET-SO-KO', 6.5000000, 44.5000000, ARRAY['Korahey', 'ቆራሄ']::text[], 49),
    ('somali', 'zone', 'Doolo Zone', 'ዶሎ ዞን', 'Doollo', 'ዞን ዶሎ', 'somali-doolo', 'ET-SO-DO', 7.0000000, 46.5000000, ARRAY['Warder Zone', 'ዶሎ']::text[], 50),
    ('somali', 'zone', 'Erer Zone (Somali)', 'ኤረር ዞን (ሶማሌ)', 'Erer', 'ዞን ኤረር (ሶማሌ)', 'somali-erer', 'ET-SO-ER', 8.5000000, 42.5000000, ARRAY['Erer Somali', 'ኤረር']::text[], 51),
    ('somali', 'zone', 'Nogob Zone', 'ኖጎብ ዞን', 'Nogob', 'ዞን ኖጎብ', 'somali-nogob', 'ET-SO-NO', 7.8000000, 43.3000000, ARRAY['Fiq Zone', 'ኖጎብ']::text[], 52),
    ('somali', 'zone', 'Afder Zone', 'አፍዴር ዞን', 'Afder', 'ዞን ኣፍዴር', 'somali-afder', 'ET-SO-AF', 5.5000000, 43.0000000, ARRAY['Afder', 'አፍዴር']::text[], 53),
    ('afar', 'zone', 'Awsi Rasu (Zone 1)', 'አውሲ ረሱ (ዞን 1)', 'Awsi Rasu', 'ኣውሲ ረሱ (ዞን 1)', 'afar-awsi-rasu', 'ET-AF-Z1', 11.7500000, 41.0000000, ARRAY['Zone 1 Afar', 'Semera Zone', 'አውሲ ረሱ']::text[], 54),
    ('afar', 'zone', 'Kilbet Rasu (Zone 2)', 'ክልበት ረሱ (ዞን 2)', 'Kilbet Rasu', 'ክልበት ረሱ (ዞን 2)', 'afar-kilbet-rasu', 'ET-AF-Z2', 13.5000000, 40.3000000, ARRAY['Zone 2 Afar', 'ክልበት ረሱ']::text[], 55),
    ('afar', 'zone', 'Gabi Rasu (Zone 3)', 'ጋቢ ረሱ (ዞን 3)', 'Gabi Rasu', 'ጋቢ ረሱ (ዞን 3)', 'afar-gabi-rasu', 'ET-AF-Z3', 9.8000000, 40.3000000, ARRAY['Zone 3 Afar', 'Awash Zone', 'ጋቢ ረሱ']::text[], 56),
    ('afar', 'zone', 'Fanti Rasu (Zone 4)', 'ፈንቲ ረሱ (ዞን 4)', 'Fanti Rasu', 'ፈንቲ ረሱ (ዞን 4)', 'afar-fanti-rasu', 'ET-AF-Z4', 12.3000000, 40.0000000, ARRAY['Zone 4 Afar', 'ፈንቲ ረሱ']::text[], 57),
    ('afar', 'zone', 'Hari Rasu (Zone 5)', 'ሀሪ ረሱ (ዞን 5)', 'Hari Rasu', 'ሓሪ ረሱ (ዞን 5)', 'afar-hari-rasu', 'ET-AF-Z5', 10.5000000, 40.0000000, ARRAY['Zone 5 Afar', 'ሀሪ ረሱ']::text[], 58),
    ('benishangul-assosa', 'zone', 'Assosa Zone', 'አሶሳ ዞን', 'Assoosaa', 'ዞን ኣሶሳ', 'benishangul-assosa', 'ET-BE-AS', 10.0500000, 34.5300000, ARRAY['Asosa Zone', 'አሶሳ']::text[], 59),
    ('benishangul-metekel', 'zone', 'Metekel Zone', 'መተከል ዞን', 'Matakkali', 'ዞን መተከል', 'benishangul-metekel', 'ET-BE-ME', 11.0000000, 36.0000000, ARRAY['Metekel', 'መተከል']::text[], 60),
    ('benishangul-kamashi', 'zone', 'Kamashi Zone', 'ካማሺ ዞን', 'Kamaashii', 'ዞን ካማሺ', 'benishangul-kamashi', 'ET-BE-KA', 9.5000000, 35.8000000, ARRAY['Kamashi', 'ካማሺ']::text[], 61),
    ('benishangul-gumuz', 'zone', 'Mao-Komo Special Woreda', 'ማኦ-ኮሞ ልዩ ወረዳ', 'Mao-Komo', 'ፍሉይ ወረዳ ማኦ-ኮሞ', 'benishangul-mao-komo', 'ET-BE-MK', 9.4000000, 34.4000000, ARRAY['Mao Komo', 'ማኦ ኮሞ']::text[], 62),
    ('benishangul-gumuz', 'zone', 'Pawe Special Woreda', 'ፓዌ ልዩ ወረዳ', 'Paawee', 'ፍሉይ ወረዳ ፓዌ', 'benishangul-pawe', 'ET-BE-PW', 11.3000000, 36.4000000, ARRAY['Pawe', 'ፓዌ']::text[], 63),
    ('gambela', 'zone', 'Gambela Special Woreda', 'ጋምቤላ ልዩ ወረዳ', 'Gambellaa Addaa', 'ፍሉይ ወረዳ ጋምቤላ', 'gambela-special-woreda', 'ET-GA-GW', 8.2500000, 34.5800000, ARRAY['Gambela Town', 'ጋምቤላ ከተማ']::text[], 64),
    ('gambela', 'zone', 'Anuak Zone', 'አኝዋክ ዞን', 'Anywaa', 'ዞን ኣኝዋክ', 'gambela-anuak', 'ET-GA-AN', 7.8000000, 34.2000000, ARRAY['Anywaa Zone', 'አኝዋክ']::text[], 65),
    ('gambela', 'zone', 'Nuer Zone', 'ኑዌር ዞን', 'Nueer', 'ዞን ኑዌር', 'gambela-nuer', 'ET-GA-NU', 8.3000000, 33.8000000, ARRAY['Nuer Zone', 'ኑዌር']::text[], 66),
    ('gambela', 'zone', 'Majang Zone', 'ማጃንግ ዞን', 'Majaang', 'ዞን ማጃንግ', 'gambela-majang', 'ET-GA-MA', 7.3000000, 35.1000000, ARRAY['Majang Zone', 'ማጃንግ']::text[], 67),
    ('harari', 'zone', 'Harari Urban & Rural Woredas', 'ሐረሪ ከተማና ገጠር ወረዳዎች', 'Hararii Magaalaa fi Baadiyyaa', 'ወረዳታት ሓረሪ', 'harari-districts', 'ET-HA-DI', 9.3100000, 42.1300000, ARRAY['Harar Districts', 'ሐረሪ']::text[], 68),
    ('central-ethiopia', 'zone', 'Gurage Zone', 'ጉራጌ ዞን', 'Guraagee', 'ዞን ጉራጌ', 'central-ethiopia-gurage', 'ET-CE-GU', 8.2000000, 38.0000000, ARRAY['Gurage Zone', 'ጉራጌ']::text[], 69),
    ('central-ethiopia', 'zone', 'Silte Zone', 'ስልጤ ዞን', 'Sillxee', 'ዞን ስልጤ', 'central-ethiopia-silte', 'ET-CE-SI', 7.9000000, 38.2000000, ARRAY['Silti Zone', 'ስልጤ']::text[], 70),
    ('central-ethiopia', 'zone', 'Hadiya Zone', 'ሀዲያ ዞን', 'Haadiyyaa', 'ዞን ሃዲያ', 'central-ethiopia-hadiya', 'ET-CE-HA', 7.6000000, 37.7500000, ARRAY['Hadiya Zone', 'ሀዲያ']::text[], 71),
    ('central-ethiopia', 'zone', 'Halaba Zone', 'ሀላባ ዞን', 'Halaabaa', 'ዞን ሃላባ', 'central-ethiopia-halaba', 'ET-CE-HL', 7.3500000, 38.1000000, ARRAY['Alaba Zone', 'ሀላባ']::text[], 72),
    ('central-ethiopia', 'zone', 'Kembata Zone', 'ከምባታ ዞን', 'Kambaataa', 'ዞን ከምባታ', 'central-ethiopia-kembata', 'ET-CE-KM', 7.2500000, 37.8500000, ARRAY['Kembata Tembaro', 'ከምባታ']::text[], 73),
    ('central-ethiopia', 'zone', 'Tembaro Special Woreda', 'ተምባሮ ልዩ ወረዳ', 'Xambaaroo', 'ፍሉይ ወረዳ ተምባሮ', 'central-ethiopia-tembaro', 'ET-CE-TB', 7.2000000, 37.6000000, ARRAY['Tembaro', 'ተምባሮ']::text[], 74),
    ('central-ethiopia', 'zone', 'Yem Zone', 'የም ዞን', 'Yamm', 'ዞን የም', 'central-ethiopia-yem', 'ET-CE-YM', 7.9500000, 37.4000000, ARRAY['Yem Special Woreda', 'የም']::text[], 75),
    ('south-ethiopia', 'zone', 'Wolayita Zone', 'ወላይታ ዞን', 'Walaayittaa', 'ዞን ወላይታ', 'south-ethiopia-wolayita', 'ET-SE-WO', 6.8500000, 37.7500000, ARRAY['Wolaita Zone', 'ወላይታ']::text[], 76),
    ('south-ethiopia', 'zone', 'Gamo Zone', 'ጋሞ ዞን', 'Gaamoo', 'ዞን ጋሞ', 'south-ethiopia-gamo', 'ET-SE-GA', 6.1500000, 37.4500000, ARRAY['Gamo Zone', 'ጋሞ']::text[], 77),
    ('south-ethiopia', 'zone', 'Gofa Zone', 'ጎፋ ዞን', 'Gofaa', 'ዞን ጎፋ', 'south-ethiopia-gofa', 'ET-SE-GF', 6.3000000, 36.8500000, ARRAY['Gofa Zone', 'ጎፋ']::text[], 78),
    ('south-ethiopia', 'zone', 'South Omo Zone', 'ደቡብ ኦሞ ዞን', 'Oomoo Kibbaa', 'ዞን ደቡብ ኦሞ', 'south-ethiopia-south-omo', 'ET-SE-SO', 5.5000000, 36.5000000, ARRAY['South Omo', 'ደቡብ ኦሞ']::text[], 79),
    ('south-ethiopia', 'zone', 'Ari Zone', 'አሪ ዞን', 'Aarii', 'ዞን ኣሪ', 'south-ethiopia-ari', 'ET-SE-AR', 5.8000000, 36.6000000, ARRAY['Ari Zone', 'አሪ']::text[], 80),
    ('south-ethiopia', 'zone', 'Gedeo Zone', 'ጌዴኦ ዞን', 'Gidiyoo', 'ዞን ጌዴኦ', 'south-ethiopia-gedeo', 'ET-SE-GE', 6.1500000, 38.3000000, ARRAY['Gedeo Zone', 'ጌዴኦ']::text[], 81),
    ('south-ethiopia', 'zone', 'Konso Zone', 'ኮንሶ ዞን', 'Konsoo', 'ዞን ኮንሶ', 'south-ethiopia-konso', 'ET-SE-KO', 5.3500000, 37.4500000, ARRAY['Konso Special Zone', 'ኮንሶ']::text[], 82),
    ('south-ethiopia', 'zone', 'Basketo Zone', 'ባስኬቶ ዞን', 'Baaskeetoo', 'ዞን ባስኬቶ', 'south-ethiopia-basketo', 'ET-SE-BS', 6.2000000, 36.6000000, ARRAY['Basketo Special', 'ባስኬቶ']::text[], 83),
    ('south-ethiopia', 'zone', 'Burji Special Woreda', 'ቡርጂ ልዩ ወረዳ', 'Burjii', 'ፍሉይ ወረዳ ቡርጂ', 'south-ethiopia-burji', 'ET-SE-BJ', 5.4500000, 37.8500000, ARRAY['Burji', 'ቡርጂ']::text[], 84),
    ('south-ethiopia', 'zone', 'Amaro Special Woreda', 'አማሮ ልዩ ወረዳ', 'Amaarroo', 'ፍሉይ ወረዳ ኣማሮ', 'south-ethiopia-amaro', 'ET-SE-AM', 5.8500000, 37.9000000, ARRAY['Amaro', 'አማሮ']::text[], 85),
    ('south-ethiopia', 'zone', 'Alle Special Woreda', 'አሌ ልዩ ወረዳ', 'Allee', 'ፍሉይ ወረዳ ኣሌ', 'south-ethiopia-alle', 'ET-SE-AL', 5.5000000, 37.3000000, ARRAY['Alle', 'አሌ']::text[], 86),
    ('south-ethiopia', 'zone', 'Derashe Special Woreda', 'ደራሼ ልዩ ወረዳ', 'Diraashee', 'ፍሉይ ወረዳ ደራሼ', 'south-ethiopia-derashe', 'ET-SE-DR', 5.6000000, 37.4000000, ARRAY['Derashe', 'ደራሼ']::text[], 87),
    ('southwest-ethiopia', 'zone', 'Keffa Zone', 'ከፋ ዞን', 'Kafaa', 'ዞን ከፋ', 'southwest-ethiopia-keffa', 'ET-SW-KF', 7.2500000, 36.2500000, ARRAY['Kaffa Zone', 'ከፋ']::text[], 88),
    ('southwest-ethiopia', 'zone', 'Sheka Zone', 'ሼካ ዞን', 'Shakkichoo', 'ዞን ሼካ', 'southwest-ethiopia-sheka', 'ET-SW-SK', 7.3500000, 35.4500000, ARRAY['Sheka Zone', 'ሼካ']::text[], 89),
    ('southwest-ethiopia', 'zone', 'Bench Sheko Zone', 'ቤንች ሸኮ ዞን', 'Bench Shakkoo', 'ዞን ቤንች ሸኮ', 'southwest-ethiopia-bench-sheko', 'ET-SW-BS', 6.9000000, 35.5500000, ARRAY['Bench Maji', 'ቤንች ሸኮ']::text[], 90),
    ('southwest-ethiopia', 'zone', 'Dawro Zone', 'ዳውሮ ዞን', 'Daawroo', 'ዞን ዳውሮ', 'southwest-ethiopia-dawro', 'ET-SW-DW', 7.0000000, 37.1500000, ARRAY['Dawro Zone', 'ዳውሮ']::text[], 91),
    ('southwest-ethiopia', 'zone', 'West Omo Zone', 'ምዕራብ ኦሞ ዞን', 'Oomoo Dhihaa', 'ዞን ምዕራብ ኦሞ', 'southwest-ethiopia-west-omo', 'ET-SW-WO', 6.2000000, 35.7000000, ARRAY['West Omo', 'ምዕራብ ኦሞ']::text[], 92),
    ('southwest-ethiopia', 'zone', 'Konta Zone', 'ኮንታ ዞን', 'Konta', 'ዞን ኮንታ', 'southwest-ethiopia-konta', 'ET-SW-KT', 7.1000000, 36.8000000, ARRAY['Konta Special', 'ኮንታ']::text[], 93),
    ('dire-dawa', 'subcity', 'Dire Dawa Urban Subcity', 'ድሬዳዋ ከተማ ክፍለ ከተማ', 'Dirree Dawaa Magaalaa', 'ክፍለ ከተማ ድሬዳዋ', 'dire-dawa-urban', 'ET-DD-UR', 9.6000000, 41.8600000, ARRAY['Dire Dawa Urban', 'ድሬዳዋ ከተማ']::text[], 94),
    ('dire-dawa', 'subcity', 'Dire Dawa Rural District', 'ድሬዳዋ ገጠር አስተዳደር', 'Dirree Dawaa Baadiyyaa', 'ገጠር ምምሕዳር ድሬዳዋ', 'dire-dawa-rural', 'ET-DD-RU', 9.5500000, 41.8000000, ARRAY['Dire Dawa Rural Woredas', 'ድሬዳዋ ገጠር']::text[], 95)
) AS d(parent_slug, loc_type, name_en, name_am, name_om, name_ti, slug, admin_code, lat, lng, aliases, disp_order)
JOIN location_hierarchy p ON p.slug = d.parent_slug
ON CONFLICT (slug) DO UPDATE SET
    name_en = EXCLUDED.name_en,
    name_am = EXCLUDED.name_am,
    name_om = EXCLUDED.name_om,
    name_ti = EXCLUDED.name_ti,
    admin_code = EXCLUDED.admin_code,
    aliases = EXCLUDED.aliases;

-- 3. Seed Key Urban & Rural Woredas and City Administrations (Level 3)
INSERT INTO location_hierarchy (
    parent_id, location_type, name_en, name_am, name_om, name_ti, slug, admin_code,
    latitude, longitude, aliases, display_order, is_leaf, is_active
)
SELECT
    p.id,
    'woreda',
    d.name_en, d.name_am, d.name_om, d.name_ti,
    d.slug, d.admin_code,
    d.lat, d.lng,
    d.aliases,
    d.disp_order,
    true, true
FROM (VALUES
    ('addis-ababa-bole', 'Bole Woreda 01 (Atlas / Edna Mall)', 'ቦሌ ወረዳ 01 (አትላስ / ኤድና ሞል)', 'Boolee Woreda 01', 'ቦሌ ወረዳ 01', 'addis-ababa-bole-w01', 'AA-BO-01', 9.0000000, 38.7800000, ARRAY['Bole Atlas', 'Edna Mall', 'ቦሌ አትላስ']::text[], 1),
    ('addis-ababa-bole', 'Bole Woreda 03 (Medhanealem)', 'ቦሌ ወረዳ 03 (መድኃኔዓለም)', 'Boolee Woreda 03', 'ቦሌ ወረዳ 03', 'addis-ababa-bole-w03', 'AA-BO-03', 8.9900000, 38.7900000, ARRAY['Bole Medhanealem', 'Bole Brass', 'ቦሌ መድኃኔዓለም']::text[], 2),
    ('addis-ababa-bole', 'Bole Woreda 06 (Gerji / Imperial)', 'ቦሌ ወረዳ 06 (ገርጂ / ኢምፔሪያል)', 'Boolee Woreda 06', 'ቦሌ ወረዳ 06', 'addis-ababa-bole-w06', 'AA-BO-06', 8.9800000, 38.8100000, ARRAY['Gerji', 'Imperial', 'ገርጂ']::text[], 3),
    ('addis-ababa-bole', 'Bole Woreda 11 (Bole Arabsa / Ayat 2)', 'ቦሌ ወረዳ 11 (ቦሌ አራብሳ)', 'Boolee Woreda 11', 'ቦሌ ወረዳ 11', 'addis-ababa-bole-w11', 'AA-BO-11', 8.9700000, 38.8500000, ARRAY['Bole Arabsa', 'ቦሌ አራብሳ']::text[], 4),
    ('addis-ababa-kirkos', 'Kazanchis (Kirkos Woreda 01)', 'ካዛንቺስ (ቂርቆስ ወረዳ 01)', 'Kaazaanchis', 'ካዛንቺስ', 'addis-ababa-kirkos-kazanchis', 'AA-KI-01', 9.0200000, 38.7700000, ARRAY['Kazanchis', 'ካዛንቺስ', 'ECA']::text[], 5),
    ('addis-ababa-kirkos', 'Meskel Square / Stadium (Kirkos Woreda 02)', 'መስቀል አደባባይ / ስታዲየም', 'Addababaa Masqalaa', 'መስቀል ኣደባባይ', 'addis-ababa-kirkos-meskel-sq', 'AA-KI-02', 9.0100000, 38.7600000, ARRAY['Meskel Square', 'Stadium', 'መስቀል አደባባይ']::text[], 6),
    ('addis-ababa-yeka', 'CMC / Ayat (Yeka Woreda 08)', 'ሲኤምሲ / አያት (የካ ወረዳ 08)', 'CMC / Ayaat', 'ሲኤምሲ / ኣያት', 'addis-ababa-yeka-cmc-ayat', 'AA-YE-08', 9.0300000, 38.8600000, ARRAY['CMC', 'Ayat Real Estate', 'ሲኤምሲ']::text[], 7),
    ('addis-ababa-yeka', 'Kotebe (Yeka Woreda 03)', 'ኮተቤ (የካ ወረዳ 03)', 'Kotobee', 'ኮተቤ', 'addis-ababa-yeka-kotebe', 'AA-YE-03', 9.0400000, 38.8300000, ARRAY['Kotebe College', 'ኮተቤ']::text[], 8),
    ('addis-ababa-arada', 'Piazza (Arada Woreda 01)', 'ፒያሳ (አራዳ ወረዳ 01)', 'Piyaasaa', 'ፒያሳ', 'addis-ababa-arada-piazza', 'AA-AR-01', 9.0300000, 38.7500000, ARRAY['Piazza', 'Arada Giorgis', 'ፒያሳ']::text[], 9),
    ('addis-ababa-arada', '4 Kilo / 6 Kilo (Arada Woreda 02)', '4 ኪሎ / 6 ኪሎ', 'Arbaat Kilo / Ja''at Kilo', '4 ኪሎ / 6 ኪሎ', 'addis-ababa-arada-4kilo', 'AA-AR-02', 9.0400000, 38.7600000, ARRAY['Arat Kilo', 'Sidist Kilo', 'አራት ኪሎ']::text[], 10),
    ('addis-ababa-lideta', 'Mexico / Sengatera (Lideta Woreda 01)', 'ሜክሲኮ / ሰንጋተራ', 'Meeksikoo', 'ሜክሲኮ', 'addis-ababa-lideta-mexico', 'AA-LI-01', 9.0100000, 38.7400000, ARRAY['Mexico Square', 'Sengatera', 'ሜክሲኮ']::text[], 11),
    ('addis-ababa-addis-ketema', 'Merkato (Addis Ketema Woreda 01)', 'መርካቶ (አዲስ ከተማ ወረዳ 01)', 'MorkaaToo', 'መርካቶ', 'addis-ababa-addis-ketema-merkato', 'AA-AK-01', 9.0300000, 38.7300000, ARRAY['Merkato', 'Autobus Tera', 'መርካቶ']::text[], 12),
    ('addis-ababa-nifas-silk-lafto', 'Jemo 1-3 (Nifas Silk Woreda 01)', 'ጀሞ 1-3 (ንፋስ ስልክ)', 'Jemoo', 'ጀሞ', 'addis-ababa-nifas-silk-jemo', 'AA-NL-01', 8.9500000, 38.7100000, ARRAY['Jemo Condominium', 'ጀሞ']::text[], 13),
    ('addis-ababa-nifas-silk-lafto', 'Lebu / Mekanisa (Nifas Silk Woreda 02)', 'ለቡ / መካኒሳ', 'Labuu', 'ለቡ', 'addis-ababa-nifas-silk-lebu', 'AA-NL-02', 8.9700000, 38.7300000, ARRAY['Lebu Varnero', 'Mekanisa', 'ለቡ']::text[], 14),
    ('addis-ababa-kolfe-keranio', 'Ayer Tena / Tor Hailoch (Kolfe Woreda 01)', 'አየር ጤና / ጦር ኃይሎች', 'Qilleensa Fayyaa', 'ኣየር ጤና', 'addis-ababa-kolfe-ayer-tena', 'AA-KK-01', 9.0000000, 38.7000000, ARRAY['Ayer Tena', 'Tor Hailoch', 'አየር ጤና']::text[], 15),
    ('addis-ababa-akaki-kality', 'Kality / Tulu Dimtu (Akaki Woreda 01)', 'ቃሊቲ / ቱሉ ዲምቱ', 'Qaallittii', 'ቃሊቲ', 'addis-ababa-akaki-kality-w01', 'AA-AQ-01', 8.8900000, 38.7600000, ARRAY['Kality Industrial', 'Tulu Dimtu', 'ቃሊቲ']::text[], 16),
    ('addis-ababa-lemi-kura', 'Summit / Goro (Lemi Kura Woreda 01)', 'ሰሚት / ጎሮ (ለሚ ኩራ)', 'Saamiti / Goroo', 'ሰሚት / ጎሮ', 'addis-ababa-lemi-kura-summit', 'AA-LK-01', 9.0100000, 38.8600000, ARRAY['Summit Condominium', 'Goro', 'ሰሚት']::text[], 17),
    ('dire-dawa-urban', 'Sabian (Dire Dawa)', 'ሳቢያን (ድሬዳዋ)', 'Saabiyaan', 'ሳብያን', 'dire-dawa-urban-sabian', 'DD-UR-SA', 9.6000000, 41.8500000, ARRAY['Sabian', 'ሳቢያን']::text[], 18),
    ('dire-dawa-urban', 'Melka Jebdu (Dire Dawa)', 'መልካ ጀብዱ (ድሬዳዋ)', 'Malka Jabduu', 'መልካ ጀብዱ', 'dire-dawa-urban-melka-jebdu', 'DD-UR-MJ', 9.5700000, 41.8200000, ARRAY['Melka Jebdu', 'መልካ ጀብዱ']::text[], 19),
    ('dire-dawa-urban', 'Megala / Kebele 01-04 (Dire Dawa)', 'መጋላ (ድሬዳዋ)', 'Magaalaa Dirree', 'መጋላ', 'dire-dawa-urban-megala', 'DD-UR-MG', 9.5900000, 41.8600000, ARRAY['Megala', 'Ashawa', 'ድሬዳዋ መጋላ']::text[], 20),
    ('dire-dawa-rural', 'Gurgura Woreda', 'ጉርጉራ ወረዳ (ድሬዳዋ)', 'Gurgura', 'ጉርጉራ', 'dire-dawa-rural-gurgura', 'DD-RU-GU', 9.5500000, 41.7800000, ARRAY['Gurgura', 'ጉርጉራ']::text[], 21),
    ('dire-dawa-rural', 'Wahil Woreda', 'ዋሂል ወረዳ (ድሬዳዋ)', 'Waahil', 'ዋሂል', 'dire-dawa-rural-wahil', 'DD-RU-WA', 9.5000000, 41.7500000, ARRAY['Wahil', 'ዋሂል']::text[], 22),
    ('dire-dawa-rural', 'Jeldesa Woreda', 'ጀልዴሳ ወረዳ (ድሬዳዋ)', 'Jalddeessaa', 'ጀልዴሳ', 'dire-dawa-rural-jeldesa', 'DD-RU-JL', 9.6500000, 41.9500000, ARRAY['Jeldesa', 'ጀልዴሳ']::text[], 23),
    ('oromia-east-shewa', 'Adama Woreda & City', 'አዳማ ወረዳና ከተማ (ምስራቅ ሸዋ)', 'Aanaa fi Magaalaa Adaamaa', 'ኣዳማ (ምብራቕ ሸዋ)', 'oromia-east-shewa-adama', 'OR-ES-ADA', 8.5400000, 39.2700000, ARRAY['Nazret', 'Nazareth', 'Adama City', 'አዳማ']::text[], 24),
    ('oromia-east-shewa', 'Bishoftu City', 'ቢሾፍቱ ከተማ (ምስራቅ ሸዋ)', 'Magaalaa Bishooftuu', 'ቢሾፍቱ (ምብራቕ ሸዋ)', 'oromia-east-shewa-bishoftu', 'OR-ES-BIS', 8.7500000, 38.9800000, ARRAY['Debre Zeyit', 'Debre Zeit', 'Bishoftu', 'ቢሾፍቱ']::text[], 25),
    ('oromia-east-shewa', 'Mojo Woreda & City', 'ሞጆ ወረዳና ከተማ', 'Moozoo', 'ሞጆ', 'oromia-east-shewa-mojo', 'OR-ES-MOJ', 8.6000000, 39.1200000, ARRAY['Modjo', 'Mojo Dry Port', 'ሞጆ']::text[], 26),
    ('oromia-east-shewa', 'Batu / Ziway Woreda & City', 'ባቱ (ዝዋይ) ወረዳና ከተማ', 'Baatuu (Zway)', 'ባቱ (ዝዋይ)', 'oromia-east-shewa-batu', 'OR-ES-BAT', 7.9300000, 38.7200000, ARRAY['Batu', 'Ziway', 'Zway', 'ባቱ ዝዋይ']::text[], 27),
    ('oromia-east-shewa', 'Dugda Woreda', 'ዱግዳ ወረዳ (መቂ)', 'Aanaa Dugda', 'ዱግዳ', 'oromia-east-shewa-dugda', 'OR-ES-DUG', 8.1800000, 38.8200000, ARRAY['Dugda', 'Meki', 'መቂ']::text[], 28),
    ('oromia-east-shewa', 'Boset Woreda', 'ቦሰት ወረዳ (ወለንጪቲ)', 'Aanaa Boosat', 'ቦሰት', 'oromia-east-shewa-boset', 'OR-ES-BOS', 8.6800000, 39.5200000, ARRAY['Boset', 'Wolenchiti', 'ወለንጪቲ']::text[], 29),
    ('oromia-east-shewa', 'Bora Woreda', 'ቦራ ወረዳ', 'Aanaa Boraa', 'ቦራ', 'oromia-east-shewa-bora', 'OR-ES-BOR', 8.2500000, 38.9000000, ARRAY['Bora', 'Alem Tena', 'ቦራ']::text[], 30),
    ('oromia-east-shewa', 'Fentale Woreda', 'ፈንታሌ ወረዳ (መተሃራ)', 'Aanaa Fantaallee', 'ፈንታሌ', 'oromia-east-shewa-fentale', 'OR-ES-FEN', 8.8500000, 39.8800000, ARRAY['Fentale', 'Metehara', 'መተሃራ']::text[], 31),
    ('oromia-east-shewa', 'Lome Woreda', 'ሎሜ ወረዳ', 'Aanaa Lomee', 'ሎሜ', 'oromia-east-shewa-lome', 'OR-ES-LOM', 8.6500000, 39.1500000, ARRAY['Lome', 'ሎሜ']::text[], 32),
    ('oromia-east-shewa', 'Liben Chukala Woreda', 'ሊበን ጩቃላ ወረዳ', 'Aanaa Liiban Cuqqaalaa', 'ሊበን ጩቃላ', 'oromia-east-shewa-liben-chukala', 'OR-ES-LBC', 8.7000000, 38.8500000, ARRAY['Liben Chukala', 'Adulala', 'ሊበን ጩቃላ']::text[], 33),
    ('oromia-west-shewa', 'Ambo Woreda & City', 'አምቦ ወረዳና ከተማ (ምዕራብ ሸዋ)', 'Magaalaa fi Aanaa Amboo', 'ኣምቦ (ምዕራብ ሸዋ)', 'oromia-west-shewa-ambo', 'OR-WS-AMB', 8.9800000, 37.8500000, ARRAY['Ambo', 'Hagere Hiwot', 'አምቦ']::text[], 34),
    ('oromia-west-shewa', 'Holeta City', 'ሆለታ ከተማ (ምዕራብ ሸዋ)', 'Magaalaa Hoolotaa', 'ሆለታ (ምዕራብ ሸዋ)', 'oromia-west-shewa-holeta', 'OR-WS-HOL', 9.0500000, 38.5000000, ARRAY['Holeta Genet', 'ሆለታ']::text[], 35),
    ('oromia-west-shewa', 'Bako Tibe Woreda', 'ባኮ ጥቤ ወረዳ', 'Aanaa Baqqo Tiibbee', 'ባኮ ጥቤ', 'oromia-west-shewa-bako-tibe', 'OR-WS-BAK', 9.1300000, 37.0500000, ARRAY['Bako', 'ባኮ']::text[], 36),
    ('oromia-west-shewa', 'Ginchi / Dendi Woreda', 'ጊንጪ / ደንዲ ወረዳ', 'Aanaa Dandii / Giincii', 'ጊንጪ', 'oromia-west-shewa-ginchi-dendi', 'OR-WS-DEN', 9.0200000, 38.1300000, ARRAY['Ginchi', 'Dendi', 'ጊንጪ']::text[], 37),
    ('oromia-west-shewa', 'Jeldu Woreda', 'ጀልዱ ወረዳ', 'Aanaa Jalduu', 'ጀልዱ', 'oromia-west-shewa-jeldu', 'OR-WS-JEL', 9.2500000, 37.9000000, ARRAY['Jeldu', 'Gojjo', 'ጀልዱ']::text[], 38),
    ('oromia-west-shewa', 'Tikur Enchini Woreda', 'ጥቁር እንጪኒ ወረዳ', 'Aanaa Incinnii', 'ጥቁር እንጪኒ', 'oromia-west-shewa-tikur-enchini', 'OR-WS-TIK', 8.8000000, 37.6500000, ARRAY['Enchini', 'እንጪኒ']::text[], 39),
    ('oromia-west-shewa', 'Ejere Woreda', 'ኤጀሬ ወረዳ (አዲስ አለም)', 'Aanaa Eejjeree', 'ኤጀሬ', 'oromia-west-shewa-ejere', 'OR-WS-EJE', 9.0300000, 38.4000000, ARRAY['Addis Alem', 'Ejere', 'አዲስ አለም']::text[], 40),
    ('oromia-north-shewa', 'Fitche Woreda & City', 'ፍቼ ወረዳና ከተማ (ሰሜን ሸዋ፣ ኦሮሚያ)', 'Magaalaa fi Aanaa Fiichee', 'ፊቼ (ሰሜን ሸዋ፣ ኦሮሚያ)', 'oromia-north-shewa-fitche', 'OR-NS-FIT', 9.7700000, 38.7300000, ARRAY['Fitche', 'Fiche', 'ፍቼ ኦሮሚያ', 'Fiichee']::text[], 41),
    ('oromia-north-shewa', 'Degem Woreda', 'ደገም ወረዳ (ሰሜን ሸዋ፣ ኦሮሚያ)', 'Aanaa Daggam', 'ደገም', 'oromia-north-shewa-degem', 'OR-NS-DEG', 9.8500000, 38.6500000, ARRAY['Degem', 'Ali Doro', 'ደገም']::text[], 42),
    ('oromia-north-shewa', 'Gerar Jarso Woreda', 'ግራር ጃርሶ ወረዳ', 'Aanaa Giraar Jaarsoo', 'ግራር ጃርሶ', 'oromia-north-shewa-gerar-jarso', 'OR-NS-GJ', 9.7500000, 38.7000000, ARRAY['Gerar Jarso', 'ግራር ጃርሶ']::text[], 43),
    ('oromia-north-shewa', 'Kuyu Woreda', 'ኩዩ ወረዳ (ገርባ ጉራቻ)', 'Aanaa Kuyyuu', 'ኩዩ', 'oromia-north-shewa-kuyu', 'OR-NS-KUY', 9.5500000, 38.4500000, ARRAY['Gerba Guracha', 'Kuyu', 'ገርባ ጉራቻ']::text[], 44),
    ('oromia-north-shewa', 'Wara Jarso Woreda', 'ወረ ጃርሶ ወረዳ', 'Aanaa Waraa Jaarsoo', 'ወረ ጃርሶ', 'oromia-north-shewa-wara-jarso', 'OR-NS-WJ', 10.0500000, 38.3500000, ARRAY['Wara Jarso', 'ወረ ጃርሶ']::text[], 45),
    ('oromia-north-shewa', 'Yaya Gulele Woreda', 'ያያ ጉለሌ ወረዳ', 'Aanaa Yaayyaa Gullaallee', 'ያያ ጉለሌ', 'oromia-north-shewa-yaya-gulele', 'OR-NS-YG', 9.5500000, 38.7000000, ARRAY['Yaya Gulele', 'ያያ ጉለሌ']::text[], 46),
    ('oromia-north-shewa', 'Kimbibit Woreda', 'ክምቢቢት ወረዳ (ሸኖ)', 'Aanaa Kimbibit', 'ክምቢቢት', 'oromia-north-shewa-kimbibit', 'OR-NS-KIM', 9.3300000, 39.3000000, ARRAY['Sheno', 'Kimbibit', 'ሸኖ']::text[], 47),
    ('oromia-north-shewa', 'Aleltu Woreda', 'አለልቱ ወረዳ', 'Aanaa Alaltuu', 'አለልቱ', 'oromia-north-shewa-aleltu', 'OR-NS-ALE', 9.2000000, 39.1500000, ARRAY['Aleltu', 'አለልቱ']::text[], 48),
    ('oromia-north-shewa', 'Dera Woreda (Oromia)', 'ዴራ ወረዳ (ሰሜን ሸዋ፣ ኦሮሚያ)', 'Aanaa Deeraa (Oromiyaa)', 'ዴራ (ኦሮሚያ)', 'oromia-north-shewa-dera', 'OR-NS-DER', 10.1000000, 38.6500000, ARRAY['Dera Oromia', 'Gundo Meskel', 'ዴራ ኦሮሚያ']::text[], 49),
    ('oromia-southwest-shewa', 'Waliso Woreda & City', 'ወሊሶ ወረዳና ከተማ (ደቡብ ምዕራብ ሸዋ)', 'Magaalaa fi Aanaa Walisoo', 'ወሊሶ', 'oromia-southwest-shewa-waliso', 'OR-SWS-WAL', 8.5300000, 37.9700000, ARRAY['Woliso', 'Waliso', 'ወሊሶ']::text[], 50),
    ('oromia-southwest-shewa', 'Wonchi Woreda', 'ወንጪ ወረዳ', 'Aanaa Wancii', 'ወንጪ', 'oromia-southwest-shewa-wonchi', 'OR-SWS-WON', 8.6500000, 37.8800000, ARRAY['Wenchi Lake', 'Wonchi', 'ወንጪ']::text[], 51),
    ('oromia-southwest-shewa', 'Amaya Woreda', 'አመያ ወረዳ', 'Aanaa Amayyaa', 'አመያ', 'oromia-southwest-shewa-amaya', 'OR-SWS-AMA', 8.5500000, 37.5500000, ARRAY['Amaya', 'Gindo', 'አመያ']::text[], 52),
    ('oromia-southwest-shewa', 'Becho Woreda', 'በቾ ወረዳ (ቱሉ ቦሎ)', 'Aanaa Becho', 'በቾ', 'oromia-southwest-shewa-becho', 'OR-SWS-BEC', 8.6500000, 38.2000000, ARRAY['Tulu Bolo', 'Becho', 'ቱሉ ቦሎ']::text[], 53),
    ('oromia-southwest-shewa', 'Dawo Woreda', 'ዳዎ ወረዳ', 'Aanaa Daawoo', 'ዳዎ', 'oromia-southwest-shewa-dawo', 'OR-SWS-DAW', 8.7000000, 38.1000000, ARRAY['Dawo', 'Busa', 'ዳዎ']::text[], 54),
    ('oromia-arsi', 'Asella Woreda & City', 'አሰላ ወረዳና ከተማ (አርሲ)', 'Magaalaa fi Aanaa Asallaa', 'ኣሰላ (ኣርሲ)', 'oromia-arsi-asella', 'OR-AR-ASE', 7.9500000, 39.1200000, ARRAY['Asella', 'Asela', 'አሰላ አርሲ', 'Asallaa']::text[], 55),
    ('oromia-arsi', 'Tiyo Woreda', 'ጢዮ ወረዳ (አርሲ)', 'Aanaa Tiyyo', 'ጢዮ', 'oromia-arsi-tiyo', 'OR-AR-TIY', 7.9800000, 39.1500000, ARRAY['Tiyo', 'ጢዮ']::text[], 56),
    ('oromia-arsi', 'Bekoji Woreda', 'በቆጂ ወረዳ (አርሲ)', 'Aanaa Baqqojjii', 'በቆጂ', 'oromia-arsi-bekoji', 'OR-AR-BEK', 7.5300000, 39.2500000, ARRAY['Bekoji', 'በቆጂ']::text[], 57),
    ('oromia-arsi', 'Digelu Tijo Woreda', 'ዲገሉ ጢጆ ወረዳ (ሰጉሬ)', 'Aanaa Digaluu Xiijjo', 'ዲገሉ ጢጆ', 'oromia-arsi-digelu-tijo', 'OR-AR-DIG', 7.7500000, 39.1500000, ARRAY['Sagure', 'Digelu Tijo', 'ሰጉሬ']::text[], 58),
    ('oromia-arsi', 'Shirka Woreda', 'ሽርካ ወረዳ (ጎበሳ)', 'Aanaa Shirka', 'ሽርካ', 'oromia-arsi-shirka', 'OR-AR-SHI', 7.6000000, 39.5000000, ARRAY['Gobesa', 'Shirka', 'ሽርካ']::text[], 59),
    ('oromia-arsi', 'Munessa Woreda', 'ሙኔሳ ወረዳ (ቀርሳ)', 'Aanaa Muneessa', 'ሙኔሳ', 'oromia-arsi-munessa', 'OR-AR-MUN', 7.5500000, 38.9000000, ARRAY['Munessa', 'Kersa Arsi', 'ሙኔሳ']::text[], 60),
    ('oromia-arsi', 'Hetosa Woreda', 'ሄጦሳ ወረዳ (ኢተያ)', 'Aanaa Heexosa', 'ሄጦሳ', 'oromia-arsi-hetosa', 'OR-AR-HET', 8.1500000, 39.2500000, ARRAY['Iteya', 'Hetosa', 'ኢተያ']::text[], 61),
    ('oromia-arsi', 'Robe Woreda (Arsi)', 'ሮቤ ወረዳ (አርሲ)', 'Aanaa Roobee (Arsii)', 'ሮቤ (ኣርሲ)', 'oromia-arsi-robe', 'OR-AR-ROB', 7.8500000, 39.6300000, ARRAY['Robe Arsi', 'ሮቤ አርሲ']::text[], 62),
    ('oromia-west-arsi', 'Shashemene City & Woreda', 'ሻሸመኔ ከተማና ወረዳ (ምዕራብ አርሲ)', 'Magaalaa fi Aanaa Shaashamannee', 'ሻሸመኔ', 'oromia-west-arsi-shashemene', 'OR-WAR-SHA', 7.2000000, 38.6000000, ARRAY['Shashemene', 'Shashamane', 'ሻሸመኔ']::text[], 63),
    ('oromia-west-arsi', 'Arsi Negele Woreda & City', 'አርሲ ነገሌ ወረዳና ከተማ', 'Magaalaa fi Aanaa Arsii Nageellee', 'ኣርሲ ነገሌ', 'oromia-west-arsi-arsi-negele', 'OR-WAR-NEG', 7.3500000, 38.7000000, ARRAY['Arsi Negele', 'አርሲ ነገሌ']::text[], 64),
    ('oromia-west-arsi', 'Kofele Woreda', 'ኮፈሌ ወረዳ', 'Aanaa Kofalee', 'ኮፈሌ', 'oromia-west-arsi-kofele', 'OR-WAR-KOF', 7.0700000, 38.7800000, ARRAY['Kofele', 'ኮፈሌ']::text[], 65),
    ('oromia-west-arsi', 'Dodola Woreda & City', 'ዶዶላ ወረዳና ከተማ', 'Magaalaa fi Aanaa Doddolaa', 'ዶዶላ', 'oromia-west-arsi-dodola', 'OR-WAR-DOD', 6.9800000, 39.1800000, ARRAY['Dodola', 'ዶዶላ']::text[], 66),
    ('oromia-west-arsi', 'Adaba Woreda', 'አዳባ ወረዳ', 'Aanaa Adabbaa', 'አዳባ', 'oromia-west-arsi-adaba', 'OR-WAR-ADA', 7.0000000, 39.4000000, ARRAY['Adaba', 'አዳባ']::text[], 67),
    ('oromia-west-arsi', 'Nensebo Woreda', 'ነንሰቦ ወረዳ (ወርቃ)', 'Aanaa Nanseboo', 'ነንሰቦ', 'oromia-west-arsi-nensebo', 'OR-WAR-NEN', 6.6000000, 39.1000000, ARRAY['Nensebo', 'Werka', 'ነንሰቦ']::text[], 68),
    ('oromia-bale', 'Bale Robe City', 'ባሌ ሮቤ ከተማ', 'Magaalaa Roobee (Baalee)', 'ባሌ ሮቤ', 'oromia-bale-robe', 'OR-BA-ROB', 7.1200000, 40.0000000, ARRAY['Bale Robe', 'Robe Bale', 'ባሌ ሮቤ']::text[], 69),
    ('oromia-bale', 'Goba Woreda & City', 'ጎባ ወረዳና ከተማ', 'Magaalaa fi Aanaa Goba', 'ጎባ', 'oromia-bale-goba', 'OR-BA-GOB', 7.0000000, 39.9800000, ARRAY['Goba Bale', 'ጎባ']::text[], 70),
    ('oromia-bale', 'Sinana Woreda', 'ሲናና ወረዳ', 'Aanaa Sinaanaa', 'ሲናና', 'oromia-bale-sinana', 'OR-BA-SIN', 7.1000000, 40.2000000, ARRAY['Sinana', 'ሲናና']::text[], 71),
    ('oromia-bale', 'Agarfa Woreda', 'አጋርፋ ወረዳ', 'Aanaa Agarfaa', 'አጋርፋ', 'oromia-bale-agarfa', 'OR-BA-AGA', 7.2700000, 39.8200000, ARRAY['Agarfa', 'አጋርፋ']::text[], 72),
    ('oromia-bale', 'Dinsho Woreda', 'ዲንሾ ወረዳ (ባሌ ተራሮች)', 'Aanaa Dinnshoo', 'ዲንሾ', 'oromia-bale-dinsho', 'OR-BA-DIN', 7.1200000, 39.7800000, ARRAY['Dinsho', 'Bale Mountains', 'ዲንሾ']::text[], 73),
    ('oromia-bale', 'Delo Mena Woreda', 'ደሎ መና ወረዳ', 'Aanaa Dalloo Mannaa', 'ደሎ መና', 'oromia-bale-delo-mena', 'OR-BA-DEL', 6.4000000, 39.8500000, ARRAY['Delo Mena', 'Mena', 'ደሎ መና']::text[], 74),
    ('oromia-borena', 'Yabelo Woreda & City', 'ያቤሎ ወረዳና ከተማ (ቦረና)', 'Magaalaa fi Aanaa Yaaballoo', 'ያቤሎ (ቦረና)', 'oromia-borena-yabelo', 'OR-BO-YAB', 4.8800000, 38.0800000, ARRAY['Yabelo', 'Yabello', 'ያቤሎ']::text[], 75),
    ('oromia-borena', 'Moyale Woreda (Oromia)', 'ሞያሌ ወረዳ (ቦረና፣ ኦሮሚያ)', 'Aanaa Moyaalee (Boorana)', 'ሞያሌ (ኦሮሚያ)', 'oromia-borena-moyale', 'OR-BO-MOY', 3.5300000, 39.0500000, ARRAY['Moyale Oromia', 'Moyale Border', 'ሞያሌ ኦሮሚያ']::text[], 76),
    ('oromia-borena', 'Mega Woreda & City', 'ሜጋ ወረዳና ከተማ', 'Magaalaa fi Aanaa Meegaa', 'ሜጋ', 'oromia-borena-mega', 'OR-BO-MEG', 4.0500000, 38.3200000, ARRAY['Mega', 'ሜጋ']::text[], 77),
    ('oromia-borena', 'Teltele Woreda', 'ተልተሌ ወረዳ', 'Aanaa Taltaallee', 'ተልተሌ', 'oromia-borena-teltele', 'OR-BO-TEL', 5.0500000, 37.4000000, ARRAY['Teltele', 'ተልተሌ']::text[], 78),
    ('oromia-guji', 'Negele Borana City & Woreda', 'ነገሌ ቦረና ከተማና ወረዳ', 'Magaalaa Nageellee Booranaa', 'ነገሌ ቦረና', 'oromia-guji-negele-borana', 'OR-GU-NEG', 5.3300000, 39.5800000, ARRAY['Negele Borana', 'Negelle', 'ነገሌ']::text[], 79),
    ('oromia-guji', 'Adola Rede Woreda (ክብረ መንግሥት)', 'አዶላ ሬዴ ወረዳ (ክብረ መንግሥት)', 'Aanaa Adoolaa Reeddee', 'ኣዶላ ሬዴ', 'oromia-guji-adola-rede', 'OR-GU-ADO', 5.8800000, 38.9800000, ARRAY['Adola', 'Kibre Mengist', 'ክብረ መንግሥት']::text[], 80),
    ('oromia-guji', 'Shakiso Woreda & City', 'ሻኪሶ ወረዳና ከተማ (ወርቅ ማዕድን)', 'Magaalaa fi Aanaa Shaakkisoo', 'ሻኪሶ', 'oromia-guji-shakiso', 'OR-GU-SHA', 5.7500000, 38.9000000, ARRAY['Shakiso', 'Lega Dembi Gold', 'ሻኪሶ']::text[], 81),
    ('oromia-west-guji', 'Bule Hora Woreda & City', 'ቡሌ ሆራ ወረዳና ከተማ (ሀገረ ማርያም)', 'Magaalaa fi Aanaa Bulee Hoolaa', 'ቡሌ ሆራ', 'oromia-west-guji-bule-hora', 'OR-WGU-BUL', 5.6300000, 38.2300000, ARRAY['Bule Hora', 'Hagere Mariam Guji', 'ቡሌ ሆራ']::text[], 82),
    ('oromia-east-hararghe', 'Haramaya Woreda & City', 'ሐረማያ ወረዳና ከተማ', 'Magaalaa fi Aanaa Haramayaa', 'ሓረማያ', 'oromia-east-hararghe-haramaya', 'OR-EH-HAR', 9.4000000, 42.0200000, ARRAY['Haramaya University', 'Alemaya', 'ሐረማያ']::text[], 83),
    ('oromia-east-hararghe', 'Aweday City (የጫት ማዕከል)', 'አወዳይ ከተማ (የጫት ማዕከል)', 'Magaalaa Awaday', 'ኣወዳይ', 'oromia-east-hararghe-aweday', 'OR-EH-AWE', 9.3800000, 42.0600000, ARRAY['Aweday Khat Market', 'አወዳይ']::text[], 84),
    ('oromia-east-hararghe', 'Babile Woreda (Oromia)', 'ባቢሌ ወረዳ (ምስራቅ ሐረርጌ)', 'Aanaa Baabilee (Harargee)', 'ባቢሌ (ሓረርጌ)', 'oromia-east-hararghe-babile', 'OR-EH-BAB', 9.2200000, 42.3300000, ARRAY['Babile Elephant Sanctuary', 'ባቢሌ ኦሮሚያ']::text[], 85),
    ('oromia-east-hararghe', 'Kombolcha Woreda (East Hararghe)', 'ኮምቦልቻ ወረዳ (ምስራቅ ሐረርጌ)', 'Aanaa Qombolcha (Harargee)', 'ኮምቦልቻ (ሓረርጌ)', 'oromia-east-hararghe-kombolcha', 'OR-EH-KOM', 9.4300000, 42.1200000, ARRAY['Kombolcha Hararghe', 'ኮምቦልቻ ሐረርጌ']::text[], 86),
    ('oromia-west-hararghe', 'Chiro Woreda & City (አሰበ ተፈሪ)', 'ጭሮ ወረዳና ከተማ (አሰበ ተፈሪ)', 'Magaalaa fi Aanaa Ciroo', 'ጭሮ (ኣሰበ ተፈሪ)', 'oromia-west-hararghe-chiro', 'OR-WH-CHI', 9.0800000, 40.8700000, ARRAY['Chiro', 'Asebe Teferi', 'ጭሮ']::text[], 87),
    ('oromia-west-hararghe', 'Bedessa Woreda & City', 'በደሳ ወረዳና ከተማ', 'Magaalaa fi Aanaa Baddassaa', 'በደሳ', 'oromia-west-hararghe-bedessa', 'OR-WH-BED', 8.9000000, 40.7800000, ARRAY['Bedessa', 'በደሳ']::text[], 88),
    ('oromia-west-hararghe', 'Hirna City', 'ሂርና ከተማ', 'Magaalaa Hirna', 'ሂርና', 'oromia-west-hararghe-hirna', 'OR-WH-HIR', 9.2200000, 41.1000000, ARRAY['Hirna', 'ሂርና']::text[], 89),
    ('oromia-west-hararghe', 'Mieso Woreda (Oromia)', 'ሚኤሶ ወረዳ (ምዕራብ ሐረርጌ)', 'Aanaa Mi''eessoo (Harargee)', 'ሚኤሶ (ኦሮሚያ)', 'oromia-west-hararghe-mieso', 'OR-WH-MIE', 9.2300000, 40.7500000, ARRAY['Mieso Oromia', 'Mieso Railway', 'ሚኤሶ ኦሮሚያ']::text[], 90),
    ('oromia-jimma-zone', 'Jimma City (አባ ጅፋር)', 'ጅማ ከተማ (አባ ጅፋር)', 'Magaalaa Jimmaa', 'ጅማ ከተማ', 'oromia-jimma-city', 'OR-JM-CIT', 7.6700000, 36.8300000, ARRAY['Jimma Town', 'Aba Jifar', 'ጅማ ከተማ']::text[], 91),
    ('oromia-jimma-zone', 'Agaro Woreda & City', 'አጋሮ ወረዳና ከተማ (ቡና)', 'Magaalaa fi Aanaa Aggaaroo', 'ኣጋሮ', 'oromia-jimma-agaro', 'OR-JM-AGA', 7.8500000, 36.6500000, ARRAY['Agaro Coffee', 'አጋሮ']::text[], 92),
    ('oromia-jimma-zone', 'Gomma Woreda', 'ጎማ ወረዳ', 'Aanaa Gommaa', 'ጎማ', 'oromia-jimma-gomma', 'OR-JM-GOM', 7.8700000, 36.6300000, ARRAY['Gomma', 'ጎማ']::text[], 93),
    ('oromia-jimma-zone', 'Limmu Kosa Woreda (ሊሙ ቡና)', 'ሊሙ ኮሳ ወረዳ (ሊሙ ቡና)', 'Aanaa Limmuu Qoxxisaa', 'ሊሙ ኮሳ', 'oromia-jimma-limmu-kosa', 'OR-JM-LIM', 8.1700000, 36.9500000, ARRAY['Limmu Kosa', 'Limu Coffee', 'ሊሙ ኮሳ']::text[], 94),
    ('oromia-special-zone-finfinne', 'Burayu City', 'ቡራዩ ከተማ (ፊንፊኔ ዙሪያ)', 'Magaalaa Burraayyuu', 'ቡራዩ', 'oromia-finfinne-burayu', 'OR-SZ-BUR', 9.0700000, 38.6800000, ARRAY['Burayu', 'ቡራዩ']::text[], 95),
    ('oromia-special-zone-finfinne', 'Sululta City & Woreda', 'ሱሉልታ ከተማና ወረዳ', 'Magaalaa fi Aanaa Sulultaa', 'ሱሉልታ', 'oromia-finfinne-sululta', 'OR-SZ-SUL', 9.1800000, 38.7500000, ARRAY['Sululta', 'ሱሉልታ']::text[], 96),
    ('oromia-special-zone-finfinne', 'Sebeta City & Woreda', 'ሰበታ ከተማና ወረዳ', 'Magaalaa fi Aanaa Sabbataa', 'ሰበታ', 'oromia-finfinne-sebeta', 'OR-SZ-SEB', 8.9200000, 38.6200000, ARRAY['Sebeta', 'ሰበታ']::text[], 97),
    ('oromia-special-zone-finfinne', 'Legetafo Legedadi City', 'ለገጣፎ ለገዳዲ ከተማ', 'Magaalaa Laga Xaafoo Laga Daadhii', 'ለገጣፎ ለገዳዲ', 'oromia-finfinne-legetafo', 'OR-SZ-LEG', 9.0800000, 38.8800000, ARRAY['Legetafo', 'Legedadi', 'ለገጣፎ']::text[], 98),
    ('oromia-special-zone-finfinne', 'Dukem City', 'ዱከም ከተማ (ኢንዱስትሪ ፓርክ)', 'Magaalaa Dukam', 'ዱከም', 'oromia-finfinne-dukem', 'OR-SZ-DUK', 8.8000000, 38.9000000, ARRAY['Dukem Industrial', 'Eastern Industry Zone', 'ዱከም']::text[], 99),
    ('oromia-special-zone-finfinne', 'Sendafa Beke City & Woreda', 'ሰንዳፋ በኬ ከተማና ወረዳ', 'Magaalaa Sandafaa', 'ሰንዳፋ', 'oromia-finfinne-sendafa', 'OR-SZ-SEN', 9.1500000, 39.0200000, ARRAY['Sendafa', 'ሰንዳፋ']::text[], 100),
    ('oromia-special-zone-finfinne', 'Gelan City', 'ገላን ከተማ', 'Magaalaa Galaan', 'ገላን', 'oromia-finfinne-gelan', 'OR-SZ-GEL', 8.8500000, 38.8300000, ARRAY['Gelan', 'ገላን']::text[], 101),
    ('amhara-north-shewa', 'Debre Berhan Woreda & City', 'ደብረ ብርሃን ወረዳና ከተማ (ሰሜን ሸዋ፣ አማራ)', 'Magaalaa fi Aanaa Dabra Birhaan', 'ደብረ ብርሃን (ሰሜን ሸዋ፣ ኣምሓራ)', 'amhara-north-shewa-debre-berhan', 'AM-NS-DB', 9.6800000, 39.5300000, ARRAY['Debre Berhan City', 'ደብረ ብርሃን', 'DB']::text[], 102),
    ('amhara-north-shewa', 'Ankober Woreda', 'አንኮበር ወረዳ', 'Aanaa Ankobar', 'ኣንኮበር', 'amhara-north-shewa-ankober', 'AM-NS-ANK', 9.5800000, 39.7300000, ARRAY['Ankober', 'Historical Capital', 'አንኮበር']::text[], 103),
    ('amhara-north-shewa', 'Shewa Robit City', 'ሸዋ ሮቢት ከተማ', 'Magaalaa Shawaa Robiit', 'ሸዋ ሮቢት', 'amhara-north-shewa-shewa-robit', 'AM-NS-SR', 10.0000000, 39.9000000, ARRAY['Shewa Robit', 'Shoa Robit', 'ሸዋ ሮቢት']::text[], 104),
    ('amhara-north-shewa', 'Basona Werana Woreda', 'ባሶና ወራና ወረዳ', 'Aanaa Baasoonaa Waraanaa', 'ባሶና ወራና', 'amhara-north-shewa-basona-werana', 'AM-NS-BW', 9.7000000, 39.5000000, ARRAY['Basona Werana', 'ባሶና ወራና']::text[], 105),
    ('amhara-north-shewa', 'Angolalla Tera Woreda (ጫቻ)', 'አንጎለላ ጠራ ወረዳ (ጫቻ)', 'Aanaa Angolallaa Xarraa', 'ኣንጎለላ ጠራ', 'amhara-north-shewa-angolalla-tera', 'AM-NS-AT', 9.6000000, 39.3800000, ARRAY['Chacha', 'Angolalla Tera', 'ጫቻ']::text[], 106),
    ('amhara-north-shewa', 'Termaber Woreda (ደብረ ሲና)', 'ተርማበር ወረዳ (ደብረ ሲና)', 'Aanaa Tarmaabar', 'ተርማበር', 'amhara-north-shewa-termaber', 'AM-NS-TER', 9.8500000, 39.7500000, ARRAY['Debre Sina', 'Termaber Tunnel', 'ደብረ ሲና']::text[], 107),
    ('amhara-north-shewa', 'Menz Gera Woreda (መሃል ሜዳ)', 'መንዝ ጌራ ወረዳ (መሃል ሜዳ)', 'Aanaa Manz Geeraa', 'መንዝ ጌራ', 'amhara-north-shewa-menz-gera', 'AM-NS-MG', 10.3000000, 39.6700000, ARRAY['Mehal Meda', 'Menz Gera', 'መሃል ሜዳ']::text[], 108),
    ('amhara-south-gondar', 'Debre Tabor Woreda & City', 'ደብረ ታቦር ወረዳና ከተማ (ደቡብ ጎንደር)', 'Magaalaa fi Aanaa Dabra Taabor', 'ደብረ ታቦር', 'amhara-south-gondar-debre-tabor', 'AM-SG-DT', 11.8500000, 38.0200000, ARRAY['Debre Tabor', 'ደብረ ታቦር']::text[], 109),
    ('amhara-south-gondar', 'Fogera Woreda (ወረታ ሩዝ)', 'ፎገራ ወረዳ (ወረታ ሩዝ)', 'Aanaa Foggaraa', 'ፎገራ (ወረታ)', 'amhara-south-gondar-fogera', 'AM-SG-FOG', 11.9200000, 37.7000000, ARRAY['Woreta', 'Fogera Rice', 'ወረታ', 'ፎገራ']::text[], 110),
    ('amhara-south-gondar', 'Dera Woreda (Amhara)', 'ዴራ ወረዳ (ደቡብ ጎንደር፣ አማራ)', 'Aanaa Deeraa (Amaara)', 'ዴራ (ኣምሓራ)', 'amhara-south-gondar-dera', 'AM-SG-DER', 11.7000000, 37.6000000, ARRAY['Dera Amhara', 'Anbesame', 'አምበሳሜ', 'ዴራ አማራ']::text[], 111),
    ('amhara-south-gondar', 'Libo Kemkem Woreda (አዲስ ዘመን)', 'ሊቦ ከምከም ወረዳ (አዲስ ዘመን)', 'Aanaa Liiboo Kamkam', 'ሊቦ ከምከም', 'amhara-south-gondar-libo-kemkem', 'AM-SG-LIB', 12.1200000, 37.7800000, ARRAY['Addis Zemen', 'Libo Kemkem', 'አዲስ ዘመን']::text[], 112),
    ('amhara-south-gondar', 'Farta Woreda', 'ፋርታ ወረዳ', 'Aanaa Faartaa', 'ፋርታ', 'amhara-south-gondar-farta', 'AM-SG-FAR', 11.9000000, 38.0500000, ARRAY['Farta', 'Gasay', 'ፋርታ']::text[], 113),
    ('amhara-central-gondar', 'Gondar City (ፋሲል ግቢ)', 'ጎንደር ከተማ (ፋሲል ግቢ)', 'Magaalaa Gondar', 'ጎንደር ከተማ', 'amhara-central-gondar-city', 'AM-CG-CIT', 12.6000000, 37.4700000, ARRAY['Gondar Castle', 'Fasil Ghebbi', 'ጎንደር ከተማ']::text[], 114),
    ('amhara-central-gondar', 'Dembia Woreda (ኮልድባ)', 'ደምቢያ ወረዳ (ኮልድባ)', 'Aanaa Dambiyaa', 'ደምቢያ', 'amhara-central-gondar-dembia', 'AM-CG-DEM', 12.3500000, 37.3000000, ARRAY['Kola Diba', 'Dembia', 'ኮልድባ']::text[], 115),
    ('amhara-north-gondar', 'Debark Woreda & City (ስሜን ተራሮች)', 'ደባርቅ ወረዳና ከተማ (ስሜን ተራሮች)', 'Magaalaa fi Aanaa Dabaariq', 'ደባርቅ', 'amhara-north-gondar-debark', 'AM-NG-DEB', 13.1500000, 37.9000000, ARRAY['Debark', 'Simien Mountains National Park', 'ደባርቅ']::text[], 116),
    ('amhara-north-gondar', 'Dabat Woreda', 'ዳባት ወረዳ', 'Aanaa Daabaat', 'ዳባት', 'amhara-north-gondar-dabat', 'AM-NG-DAB', 12.9800000, 37.7700000, ARRAY['Dabat', 'ዳባት']::text[], 117),
    ('amhara-west-gojjam', 'Bahir Dar City (ጣና ሐይቅ)', 'ባሕር ዳር ከተማ (ጣና ሐይቅ)', 'Magaalaa Baahir Daar', 'ባሕር ዳር ከተማ', 'amhara-west-gojjam-bahir-dar-city', 'AM-WJ-BDC', 11.5900000, 37.3900000, ARRAY['Bahir Dar City', 'Lake Tana', 'ባሕር ዳር']::text[], 118),
    ('amhara-west-gojjam', 'Finote Selam Woreda & City', 'ፍኖተ ሰላም ወረዳና ከተማ', 'Magaalaa Finoota Salaam', 'ፍኖተ ሰላም', 'amhara-west-gojjam-finote-selam', 'AM-WJ-FS', 10.7000000, 37.2700000, ARRAY['Finote Selam', 'ፍኖተ ሰላም']::text[], 119),
    ('amhara-west-gojjam', 'Mecha Woreda (መርዓዊ)', 'ሜጫ ወረዳ (መርዓዊ)', 'Aanaa Meexaa', 'ሜጫ (መርዓዊ)', 'amhara-west-gojjam-mecha', 'AM-WJ-MEC', 11.4200000, 37.1500000, ARRAY['Merawi', 'Mecha', 'መርዓዊ']::text[], 120),
    ('amhara-west-gojjam', 'Bure Woreda (West Gojjam)', 'ቡሬ ወረዳ (ምዕራብ ጎጃም)', 'Aanaa Buree (Gojjaam)', 'ቡሬ (ጎጃም)', 'amhara-west-gojjam-bure', 'AM-WJ-BUR', 10.7000000, 37.0700000, ARRAY['Bure Gojjam', 'Bure Agro-Processing', 'ቡሬ ጎጃም']::text[], 121),
    ('amhara-east-gojjam', 'Debre Markos Woreda & City', 'ደብረ ማርቆስ ወረዳና ከተማ (ምስራቅ ጎጃም)', 'Magaalaa fi Aanaa Dabra Maarqos', 'ደብረ ማርቆስ', 'amhara-east-gojjam-debre-markos', 'AM-EJ-DM', 10.3300000, 37.7200000, ARRAY['Debre Markos', 'ደብረ ማርቆስ']::text[], 122),
    ('amhara-east-gojjam', 'Dejen Woreda (ዓባይ በረሃ)', 'ደጀን ወረዳ (ዓባይ በረሃ)', 'Aanaa Dajan', 'ደጀን', 'amhara-east-gojjam-dejen', 'AM-EJ-DEJ', 10.1700000, 38.1300000, ARRAY['Dejen', 'Abay Gorge Bridge', 'ደጀን']::text[], 123),
    ('amhara-east-gojjam', 'Enemay Woreda (ቢቸና)', 'እነማይ ወረዳ (ቢቸና)', 'Aanaa Innamaay', 'እነማይ (ቢቸና)', 'amhara-east-gojjam-enemay', 'AM-EJ-ENE', 10.4500000, 38.2000000, ARRAY['Bichena', 'Enemay', 'ቢቸና']::text[], 124),
    ('amhara-east-gojjam', 'Hulet Ej Enese Woreda (ሞጣ)', 'ሁለት እጁ እነሴ ወረዳ (ሞጣ)', 'Aanaa Huleet Ijjii Inasii', 'ሞጣ', 'amhara-east-gojjam-hulet-ej-enese', 'AM-EJ-HE', 11.0800000, 37.8800000, ARRAY['Motta', 'Mota', 'ሞጣ']::text[], 125),
    ('amhara-south-wollo', 'Dessie City (ቶሳ)', 'ደሴ ከተማ (ቶሳ)', 'Magaalaa Dasee', 'ደሴ ከተማ', 'amhara-south-wollo-dessie-city', 'AM-SW-DSC', 11.1300000, 39.6300000, ARRAY['Dessie Town', 'Tosa', 'ደሴ ከተማ']::text[], 126),
    ('amhara-south-wollo', 'Kombolcha City (ኢንዱስትሪ)', 'ኮምቦልቻ ከተማ (ደቡብ ወሎ)', 'Magaalaa Qombolchaa', 'ኮምቦልቻ (ደቡብ ወሎ)', 'amhara-south-wollo-kombolcha', 'AM-SW-KOM', 11.0800000, 39.7300000, ARRAY['Kombolcha Industrial Park', 'Kombolcha Dry Port', 'ኮምቦልቻ ደቡብ ወሎ']::text[], 127),
    ('amhara-south-wollo', 'Tehuledere Woreda (ሐይቅ እስጢፋኖስ)', 'ተሁለደሬ ወረዳ (ሐይቅ እስጢፋኖስ)', 'Aanaa Tahuuladaree', 'ተሁለደሬ (ሓይቅ)', 'amhara-south-wollo-tehuledere', 'AM-SW-TEH', 11.3000000, 39.6800000, ARRAY['Haiq', 'Lake Haiq', 'ሐይቅ']::text[], 128),
    ('amhara-north-wollo', 'Woldiya Woreda & City', 'ወልዲያ ወረዳና ከተማ (ሰሜን ወሎ)', 'Magaalaa fi Aanaa Waldiyyaa', 'ወልዲያ', 'amhara-north-wollo-woldiya', 'AM-NW-WOL', 11.8300000, 39.6000000, ARRAY['Woldia', 'ወልዲያ']::text[], 129),
    ('amhara-north-wollo', 'Lalibela Woreda & City (ውቅር አብያተ ክርስቲያናት)', 'ላሊበላ ወረዳና ከተማ (ውቅር አብያተ ክርስቲያናት)', 'Magaalaa Laallibalaa', 'ላሊበላ', 'amhara-north-wollo-lalibela', 'AM-NW-LAL', 12.0300000, 39.0500000, ARRAY['Lalibela Rock Hewn Churches', 'ላሊበላ', 'Roha']::text[], 130),
    ('amhara-north-wollo', 'Kobo Woreda (ራያ ቆቦ)', 'ቆቦ ወረዳ (ራያ ቆቦ)', 'Aanaa Qobboo', 'ራያ ቆቦ', 'amhara-north-wollo-kobo', 'AM-NW-KOB', 12.1500000, 39.6300000, ARRAY['Raya Kobo', 'Kobo Valley', 'ቆቦ']::text[], 131),
    ('tigray-central', 'Axum Woreda & City (ኣኽሱም ጽዮን)', 'አክሱም ወረዳና ከተማ', 'Magaalaa fi Aanaa Aksum', 'ከተማ ኣኽሱም ጽዮን', 'tigray-central-axum', 'TI-CT-AXU', 14.1200000, 38.7200000, ARRAY['Axum Obelisk', 'Axum Tsion', 'ኣክሱም']::text[], 132),
    ('tigray-central', 'Adwa Woreda & City (ዓድዋ ድል)', 'ዐድዋ ወረዳና ከተማ (የዐድዋ ድል)', 'Magaalaa fi Aanaa Adwaa', 'ከተማ ዓድዋ ዓወት', 'tigray-central-adwa', 'TI-CT-ADW', 14.1700000, 38.9000000, ARRAY['Adwa Victory', 'ዓድዋ']::text[], 133),
    ('tigray-central', 'Abi Adi / Kola Tembien (ዓቢይ ዓዲ)', 'አቢ አዲ / ቆላ ተምቤን', 'Magaalaa fi Aanaa Kollaa Tambeen', 'ዓቢይ ዓዲ / ቆላ ተምቤን', 'tigray-central-abi-adi', 'TI-CT-AA', 13.6200000, 39.0000000, ARRAY['Abi Adi', 'Kola Tembien', 'ዓቢይ ዓዲ']::text[], 134),
    ('tigray-east', 'Adigrat Woreda & City', 'ዓዲግራት ወረዳና ከተማ', 'Magaalaa fi Aanaa Adigraat', 'ከተማ ዓዲግራት', 'tigray-east-adigrat', 'TI-ET-ADI', 14.2800000, 39.4600000, ARRAY['Adigrat', 'ዓዲግራት']::text[], 135),
    ('tigray-east', 'Kilte Awlaelo Woreda (ውቅሮ)', 'ውቅሮ / ክልተ ኣውላዕሎ', 'Aanaa Qilxee Awlaaloo', 'ውቅሮ / ክልተ ኣውላዕሎ', 'tigray-east-wukro', 'TI-ET-WUK', 13.7800000, 39.6000000, ARRAY['Wukro', 'Kilte Awlaelo', 'ውቅሮ']::text[], 136),
    ('tigray-south', 'Maichew Woreda & City', 'ማይጨው ወረዳና ከተማ', 'Magaalaa fi Aanaa Maaycaw', 'ከተማ ማይጨው', 'tigray-south-maichew', 'TI-ST-MAI', 12.7800000, 39.5300000, ARRAY['Maychew', 'ማይጨው']::text[], 137),
    ('tigray-south', 'Alamata Woreda & City (ራያ)', 'አላማጣ ወረዳና ከተማ (ራያ)', 'Magaalaa fi Aanaa Alaamaaxaa', 'ከተማ ኣላማጣ (ራያ)', 'tigray-south-alamata', 'TI-ST-ALA', 12.4200000, 39.5300000, ARRAY['Alamata Raya', 'ኣላማጣ']::text[], 138),
    ('tigray-north-west', 'Shire / Inda Selassie Woreda & City', 'ሽረ እንዳስላሴ ወረዳና ከተማ', 'Magaalaa fi Aanaa Shiraa', 'ከተማ ሽረ እንዳስላሴ', 'tigray-north-west-shire', 'TI-NW-SHI', 14.1000000, 38.2800000, ARRAY['Shire Indaselassie', 'ሽረ']::text[], 139),
    ('sidama-zone', 'Aleta Wondo Woreda & Town', 'አለታ ወንዶ ወረዳና ከተማ', 'Alatta Wandoo', 'ዓለታ ወንዶ', 'sidama-aleta-wondo', 'SI-ZN-AW', 6.6000000, 38.4200000, ARRAY['Aleta Wondo Coffee', 'ዓለታ ወንዶ']::text[], 140),
    ('sidama-zone', 'Dale Woreda (ይርጋለም)', 'ዳሌ ወረዳ (ይርጋለም)', 'Daale (Yirgaalam)', 'ዳሌ (ይርጋለም)', 'sidama-dale', 'SI-ZN-DAL', 6.7500000, 38.4000000, ARRAY['Dale', 'Yirgalem Agro-Industry', 'ዳሌ']::text[], 141),
    ('sidama-zone', 'Wondo Genet Woreda (ፍልውሃ)', 'ወንዶ ገነት ወረዳ (ፍልውሃ)', 'Wando Gannat', 'ወንዶ ገነት', 'sidama-wondo-genet', 'SI-ZN-WG', 7.0200000, 38.6200000, ARRAY['Wondo Genet Hot Springs', 'ወንዶ ገነት']::text[], 142),
    ('sidama-zone', 'Bensa Woreda (ዳዬ)', 'ቤንሳ ወረዳ (ዳዬ ቡና)', 'Bansa (Daye)', 'ቤንሳ (ዳዬ)', 'sidama-bensa', 'SI-ZN-BEN', 6.5500000, 38.8000000, ARRAY['Bensa Daye Coffee', 'ቤንሳ']::text[], 143),
    ('somali-fafan', 'Jigjiga City (ቀበሌ 01-20)', 'ጅግጅጋ ከተማ አስተዳደር', 'Magaalaa Jigjigaa', 'ከተማ ጅግጅጋ', 'somali-fafan-jigjiga-city', 'SO-FA-JGC', 9.3500000, 42.8000000, ARRAY['Jijiga Metro', 'Jigjiga Town', 'ጅግጅጋ']::text[], 144),
    ('somali-fafan', 'Kebri Beyah Woreda', 'ቀብሪ በያህ ወረዳ', 'Aanaa Qabribayaah', 'ቀብሪ በያህ', 'somali-fafan-kebri-beyah', 'SO-FA-KB', 9.1000000, 43.1700000, ARRAY['Kebri Beyah', 'ቀብሪ በያህ']::text[], 145),
    ('somali-fafan', 'Awbare Woreda (ቴፍሬ በር)', 'አውባሬ ወረዳ (ቴፍሬ በር)', 'Aanaa Awbaarre', 'ኣውባሬ', 'somali-fafan-awbare', 'SO-FA-AWB', 9.7800000, 43.2200000, ARRAY['Awbare', 'Tefere Ber', 'አውባሬ']::text[], 146),
    ('somali-sitti', 'Shinile Woreda', 'ሽኒሌ ወረዳ', 'Aanaa Shiniile', 'ሽኒሌ', 'somali-sitti-shinile', 'SO-SI-SHI', 9.7000000, 41.8500000, ARRAY['Shinile', 'ሽኒሌ']::text[], 147),
    ('somali-shabele', 'Gode Woreda & City', 'ጎዴ ወረዳና ከተማ (ሸበሌ ወንዝ)', 'Magaalaa fi Aanaa Godey', 'ከተማ ጎዴ', 'somali-shabele-gode', 'SO-SH-GOD', 5.9500000, 43.5800000, ARRAY['Godey', 'Shabelle River', 'ጎዴ']::text[], 148),
    ('somali-dawa', 'Moyale Woreda (Somali)', 'ሞያሌ ወረዳ (ዳዋ፣ ሶማሌ)', 'Aanaa Moyaale (Somaalee)', 'ሞያሌ (ሶማሌ)', 'somali-dawa-moyale', 'SO-DA-MOY', 3.5400000, 39.0600000, ARRAY['Moyale Somali Border', 'ሞያሌ ሶማሌ']::text[], 149),
    ('afar-awsi-rasu', 'Semera-Logiya City', 'ሰመራ-ሎጊያ ከተማ (የክልሉ መዲና)', 'Magaalaa Samaraa-Logiyaa', 'ከተማ ሰመራ-ሎጊያ', 'afar-awsi-rasu-semera', 'AF-Z1-SEM', 11.7900000, 41.0100000, ARRAY['Semera Capital', 'Logiya Dry Port', 'ሰመራ']::text[], 150),
    ('afar-awsi-rasu', 'Asayita Woreda', 'አሳይታ ወረዳ', 'Aanaa Asayita', 'ኣሳይታ', 'afar-awsi-rasu-asayita', 'AF-Z1-ASA', 11.5700000, 41.4300000, ARRAY['Asaita', 'Aussa Sultanate', 'አሳይታ']::text[], 151),
    ('afar-gabi-rasu', 'Awash Sebat Kilo Town', 'አዋሽ ሰባት ኪሎ ከተማ', 'Magaalaa Hawaash', 'ከተማ ኣዋሽ ሰባት', 'afar-gabi-rasu-awash', 'AF-Z3-AWA', 8.9800000, 40.1700000, ARRAY['Awash 7 Kilo', 'Awash National Park', 'አዋሽ 7 ኪሎ']::text[], 152),
    ('benishangul-assosa', 'Assosa City & Woreda', 'አሶሳ ከተማና ወረዳ', 'Magaalaa fi Aanaa Assoosaa', 'ከተማ ኣሶሳ', 'benishangul-assosa-city', 'BE-AS-CIT', 10.0700000, 34.5300000, ARRAY['Assosa Town', 'አሶሳ']::text[], 153),
    ('benishangul-metekel', 'Gilgel Beles / Mandura (ህዳሴ ግድብ)', 'ግልገል በለስ (ታላቁ የህዳሴ ግድብ)', 'Magaalaa Gilgal Balas', 'ግልገል በለስ (ህዳሴ ግድብ)', 'benishangul-metekel-gilgel-beles', 'BE-ME-GB', 11.1500000, 36.3300000, ARRAY['GERD', 'Grand Ethiopian Renaissance Dam', 'Gilgel Beles', 'ግልገል በለስ']::text[], 154),
    ('gambela-special-woreda', 'Gambela Town & Zuria (ባሮ ወንዝ)', 'ጋምቤላ ከተማና ዙሪያ (ባሮ ወንዝ)', 'Magaalaa Gambellaa', 'ከተማ ጋምቤላ', 'gambela-town', 'GA-GW-CIT', 8.2500000, 34.5800000, ARRAY['Gambella Town', 'Baro River Port', 'ጋምቤላ']::text[], 155),
    ('central-ethiopia-gurage', 'Butajira Woreda & City', 'ቡታጅራ ወረዳና ከተማ (ጉራጌ)', 'Magaalaa fi Aanaa Butaajiraa', 'ቡታጅራ', 'central-ethiopia-gurage-butajira', 'CE-GU-BUT', 8.1200000, 38.3700000, ARRAY['Butajira', 'ቡታጅራ']::text[], 156),
    ('central-ethiopia-gurage', 'Wolkite Woreda & City', 'ወልቂጤ ወረዳና ከተማ (ጉራጌ)', 'Magaalaa fi Aanaa Walqixxee', 'ወልቂጤ', 'central-ethiopia-gurage-wolkite', 'CE-GU-WOL', 8.2800000, 37.7800000, ARRAY['Wolkite', 'ወልቂጤ']::text[], 157),
    ('central-ethiopia-silte', 'Worabe Woreda & City', 'ወራቤ ወረዳና ከተማ (ስልጤ)', 'Magaalaa fi Aanaa Waraabee', 'ወራቤ', 'central-ethiopia-silte-worabe', 'CE-SI-WOR', 7.8500000, 38.3300000, ARRAY['Worabe Capital', 'ወራቤ']::text[], 158),
    ('central-ethiopia-hadiya', 'Hosanna Woreda & City', 'ሆሳዕና ወረዳና ከተማ (ሀዲያ)', 'Magaalaa fi Aanaa Hosaa''inaa', 'ሆሳዕና', 'central-ethiopia-hadiya-hosanna', 'CE-HA-HOS', 7.5500000, 37.8500000, ARRAY['Hossana', 'Wachemo', 'ሆሳዕና']::text[], 159),
    ('south-ethiopia-wolayita', 'Wolayita Sodo City & Woreda', 'ወላይታ ሶዶ ከተማና ወረዳ', 'Magaalaa fi Aanaa SoddOo', 'ወላይታ ሶዶ', 'south-ethiopia-wolayita-sodo', 'SE-WO-SOD', 6.8500000, 37.7500000, ARRAY['Sodo Town', 'Wolaita Sodo', 'ወላይታ ሶዶ']::text[], 160),
    ('south-ethiopia-wolayita', 'Areka / Boloso Sore Woreda', 'አረካ / ቦሎሶ ሶሬ ወረዳ (ዝንጅብል)', 'Magaalaa Arekaa', 'ኣረካ', 'south-ethiopia-wolayita-areka', 'SE-WO-ARE', 7.0700000, 37.7000000, ARRAY['Areka', 'Boloso Sore Ginger', 'አረካ']::text[], 161),
    ('south-ethiopia-gamo', 'Arba Minch City & Woreda (አባያና ጫሞ)', 'አርባ ምንጭ ከተማና ወረዳ (አባያና ጫሞ ሐይቆች)', 'Magaalaa fi Aanaa Arbaa Minc', 'ኣርባ ምንጭ', 'south-ethiopia-gamo-arba-minch', 'SE-GA-AM', 6.0300000, 37.5500000, ARRAY['Arba Minch Forty Springs', 'Nechisar Park', 'አርባ ምንጭ']::text[], 162),
    ('south-ethiopia-gamo', 'Chencha Woreda (ዶርዜ ሸማ)', 'ጨንቻ ወረዳ (የዶርዜ ሸማቾች)', 'Aanaa Cencaa', 'ጨንቻ (ዶርዜ)', 'south-ethiopia-gamo-chencha', 'SE-GA-CHE', 6.2500000, 37.5700000, ARRAY['Chencha Dorze Weaving', 'ጨንቻ ዶርዜ']::text[], 163),
    ('south-ethiopia-south-omo', 'Jinka Woreda & City (ሙርሲ)', 'ጂንካ ወረዳና ከተማ (የሙርሲ ባህል)', 'Magaalaa Jiinkaa', 'ጂንካ', 'south-ethiopia-south-omo-jinka', 'SE-SO-JIN', 5.7800000, 36.5700000, ARRAY['Jinka Mursi Hamer Cultural Center', 'ጂንካ']::text[], 164),
    ('south-ethiopia-gedeo', 'Dilla Woreda & City', 'ዲላ ወረዳና ከተማ (ጌዴኦ)', 'Magaalaa Dillaa', 'ዲላ', 'south-ethiopia-gedeo-dilla', 'SE-GE-DIL', 6.4200000, 38.3200000, ARRAY['Dilla University', 'ዲላ']::text[], 165),
    ('south-ethiopia-gedeo', 'Yirgacheffe Woreda (ይርጋጨፌ ቡና)', 'ይርጋጨፌ ወረዳ (ዓለም አቀፍ ቡና)', 'Aanaa Yirgaacafee', 'ይርጋጨፌ', 'south-ethiopia-gedeo-yirgacheffe', 'SE-GE-YIR', 6.1700000, 38.2000000, ARRAY['Yirgacheffe Coffee Origin', 'ይርጋጨፌ']::text[], 166),
    ('southwest-ethiopia-keffa', 'Bonga Woreda & City (የአረቢካ ቡና መገኛ)', 'ቦንጋ ወረዳና ከተማ (የዓለም አረቢካ ቡና መገኛ)', 'Magaalaa fi Aanaa Boongaa', 'ቦንጋ (ቡና መበቆል)', 'southwest-ethiopia-keffa-bonga', 'SW-KF-BON', 7.2800000, 36.2300000, ARRAY['Bonga Wild Arabica Coffee Biosphere', 'ቦንጋ']::text[], 167),
    ('southwest-ethiopia-bench-sheko', 'Mizan Aman Woreda & City', 'ሚዛን አማን ወረዳና ከተማ', 'Magaalaa Miizaan Amaan', 'ሚዛን ኣማን', 'southwest-ethiopia-bench-mizan', 'SW-BS-MIZ', 6.9800000, 35.5800000, ARRAY['Mizan Teferi', 'Aman', 'ሚዛን አማን']::text[], 168),
    ('southwest-ethiopia-sheka', 'Tepi / Yeki Woreda (ቅመማ ቅመም)', 'ቴፒ / የኪ ወረዳ (ቅመማ ቅመም)', 'Magaalaa Xappii', 'ቴፒ', 'southwest-ethiopia-sheka-tepi', 'SW-SK-TEP', 7.2000000, 35.4300000, ARRAY['Tepi Spices Coffee', 'ቴፒ']::text[], 169)
) AS d(parent_slug, name_en, name_am, name_om, name_ti, slug, admin_code, lat, lng, aliases, disp_order)
JOIN location_hierarchy p ON p.slug = d.parent_slug
ON CONFLICT (slug) DO UPDATE SET
    name_en = EXCLUDED.name_en,
    name_am = EXCLUDED.name_am,
    name_om = EXCLUDED.name_om,
    name_ti = EXCLUDED.name_ti,
    admin_code = EXCLUDED.admin_code,
    aliases = EXCLUDED.aliases;

-- 4. Backward Compatibility: Synchronize legacy locations table
INSERT INTO locations (region, zone, woreda, region_am, zone_am, woreda_am, region_or, zone_or, woreda_or, region_ti, zone_ti, woreda_ti)
SELECT
    reg.name_en, zn.name_en, w.name_en,
    reg.name_am, zn.name_am, w.name_am,
    reg.name_om, zn.name_om, w.name_om,
    reg.name_ti, zn.name_ti, w.name_ti
FROM location_hierarchy w
JOIN location_hierarchy zn ON zn.id = w.parent_id
JOIN location_hierarchy reg ON reg.id = zn.parent_id
WHERE w.location_type = 'woreda'
  AND NOT EXISTS (
      SELECT 1 FROM locations l
      WHERE l.region = reg.name_en AND l.zone = zn.name_en AND l.woreda = w.name_en
  );

-- 5. For Addis Ababa and Dire Dawa subcities/woredas:
INSERT INTO locations (region, subcity, woreda, region_am, subcity_am, woreda_am, region_or, subcity_or, woreda_or, region_ti, subcity_ti, woreda_ti)
SELECT
    reg.name_en, zn.name_en, w.name_en,
    reg.name_am, zn.name_am, w.name_am,
    reg.name_om, zn.name_om, w.name_om,
    reg.name_ti, zn.name_ti, w.name_ti
FROM location_hierarchy w
JOIN location_hierarchy zn ON zn.id = w.parent_id
JOIN location_hierarchy reg ON reg.id = zn.parent_id
WHERE reg.location_type = 'city' AND w.location_type = 'woreda'
  AND NOT EXISTS (
      SELECT 1 FROM locations l
      WHERE l.region = reg.name_en AND l.subcity = zn.name_en AND l.woreda = w.name_en
  );

COMMIT;
