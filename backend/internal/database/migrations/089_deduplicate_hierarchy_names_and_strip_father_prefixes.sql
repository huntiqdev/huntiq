-- ============================================================================
-- Migration 089: Deduplicate Hierarchy Names & Strip Redundant Father Prefixes
-- ============================================================================
-- 1. Cleans Level 1 Superfamily names by removing redundant root domain titles
--    (e.g., under 'Agriculture & Food': 'Agricultural Infrastructure' -> 'Infrastructure & Storage',
--    'Agricultural Inputs' -> 'Inputs, Seeds & Fertilizers', etc.).
-- 2. Deduplicates Deep Layers (Level >= 2): strips compound parent prefixes
--    (e.g., 'Phones & Tablets - 5G Flagship & High-Performance Smartphones'
--    becomes '5G Flagship & High-Performance Smartphones'; and
--    'Agricultural Infrastructure - Bespoke & Custom Engineered Editions'
--    becomes 'Bespoke & Custom Engineered Editions').
-- 3. Drops and recreates trigram indexes for instantaneous index synchronization.
-- ============================================================================

-- 1. Superfamily De-duplication under Root Domains (Level 1)
-- Agriculture & Food (id 4)
UPDATE category_hierarchy SET 
    name_en = 'Infrastructure & Storage',
    name_am = 'መሠረተ ልማትና ማከማቻ',
    name_om = 'Bu''uuraalee fi Kuusaa',
    name_ti = 'መሰረተ ልምዓትን መኽዘንን'
WHERE id = 106 AND parent_id = 4;

UPDATE category_hierarchy SET 
    name_en = 'Inputs, Seeds & Fertilizers',
    name_am = 'ግብዓቶች፣ ዘሮችና ማዳበሪያዎች',
    name_om = 'Galteewwan, Sanyii fi Xaa''oo',
    name_ti = 'እታዎታት፣ ዘርኢታትን ድኽዒታትን'
WHERE id = 247 AND parent_id = 4;

UPDATE category_hierarchy SET 
    name_en = 'Machinery, Tools & Implements',
    name_am = 'ማሽነሪዎች፣ መሣሪያዎችና ዕቃዎች',
    name_om = 'Maashinoota fi Meeshaalee Qonnaa',
    name_ti = 'ማሽነሪታትን መሳርሒታትን'
WHERE id = 288 AND parent_id = 4;

UPDATE category_hierarchy SET 
    name_en = 'Processing, Milling & Packaging',
    name_am = 'ማቀነባበሪያ፣ ወፍጮና እሽግ',
    name_om = 'Qophii, Daakuu fi Saamsama',
    name_ti = 'ምድላው፣ መጥሓንን ምዕሻግን'
WHERE id = 131 AND parent_id = 4;

UPDATE category_hierarchy SET 
    name_en = 'Farming Services & Advisory',
    name_am = 'የእርሻ ድጋፍና የማማከር አገልግሎቶች',
    name_om = 'Tajaajila Qonnaa fi Gorsa',
    name_ti = 'ናይ ሕርሻ ኣገልግሎትን ምኽርን'
WHERE id = 179 AND parent_id = 4;

UPDATE category_hierarchy SET 
    name_en = 'Trade & Commodity Exchange',
    name_am = 'የሸቀጥ ንግድና የገበያ ልውውጥ',
    name_om = 'Daldala Meeshaalee fi Gabaa',
    name_ti = 'ንግዲ ሸቐጥን ምልውዋጥ ዕዳጋን'
WHERE id = 226 AND parent_id = 4;

UPDATE category_hierarchy SET 
    name_en = 'Enabling Services & Logistics',
    name_am = 'ተዛማጅ ሎጂስቲክስና ድጋፍ',
    name_om = 'Tajaajila Deggersaa fi Loojistiksii',
    name_ti = 'ደገፍቲ ኣገልግሎታትን ሎጂስቲክስን'
WHERE id = 79 AND parent_id = 4;

-- Electronics (id 1)
UPDATE category_hierarchy SET 
    name_en = 'Repair & Maintenance Services',
    name_am = 'የጥገናና እንክብካቤ አገልግሎት',
    name_om = 'Tajaajila Suphaa fi Qulqullinaa',
    name_ti = 'ናይ ጽገናን ሓልዮትን ኣገልግሎት'
WHERE id = 1394 AND parent_id = 1;

UPDATE category_hierarchy SET 
    name_en = 'Trade-In & Equipment Exchange',
    name_am = 'የመሳሪያዎች መለዋወጫና ልውውጥ',
    name_om = 'Jijjiirraa Meeshaalee fi Gabaa',
    name_ti = 'ምልውዋጥ መሳርሒታትን ዕዳጋን'
WHERE id = 1399 AND parent_id = 1;

UPDATE category_hierarchy SET 
    name_en = 'Warranty, Protection & Overhaul',
    name_am = 'ዋስትና፣ ከለላና ሙሉ ጥገና',
    name_om = 'Wabiifi Suphaa Guutuu',
    name_ti = 'ውሕስነትን ምሉእ ጽገናን'
WHERE id = 1404 AND parent_id = 1;

-- Transport & Mobility (id 3197)
UPDATE category_hierarchy SET 
    name_en = 'Infrastructure & Transit Services',
    name_am = 'መሠረተ ልማትና የትራንዚት አገልግሎት',
    name_om = 'Tajaajila Bu''uuraalee fi Tiraanzitii',
    name_ti = 'መሰረተ ልምዓትን ትራንዚትን'
WHERE id = 3262 AND parent_id = 3197;

UPDATE category_hierarchy SET 
    name_en = 'Batteries & Electrical Systems',
    name_am = 'ባትሪዎችና የኤሌክትሪክ ክፍሎች',
    name_om = 'Baatrii fi Sirna Elektiriikaa',
    name_ti = 'ባትሪታትን ናይ ኤሌክትሪክ ክፍልታትን'
WHERE id = 3272 AND parent_id = 3197;

UPDATE category_hierarchy SET 
    name_en = 'Auto Financing & Insurance',
    name_am = 'የመኪና ፋይናንስና ኢንሹራንስ',
    name_om = 'Faayinaansii fi Inshuraansii Konkolaataa',
    name_ti = 'ፋይናንስን ኢንሹራንስን መካይን'
WHERE id = 3277 AND parent_id = 3197;

UPDATE category_hierarchy SET 
    name_en = 'Parts, Spares & Components',
    name_am = 'መለዋወጫ እቃዎችና አካላት',
    name_om = 'Kutaalee Jijjiirraa fi Meeshaalee',
    name_ti = 'መተካእታታትን መለዋወጥን'
WHERE id = 3282 AND parent_id = 3197;

UPDATE category_hierarchy SET 
    name_en = 'Rental, Hire & Leasing',
    name_am = 'ኪራይና የሊዝ ውል',
    name_om = 'Kireeffannaa fi Liizii',
    name_ti = 'ክራይን ሊዝን'
WHERE id = 3288 AND parent_id = 3197;

UPDATE category_hierarchy SET 
    name_en = 'Garages, Repair & Servicing',
    name_am = 'ጋራዦች፣ ጥገናና ሰርቪስ',
    name_om = 'Gaaraajii, Suphaa fi Sarviisii',
    name_ti = 'ጋራጃት፣ ጽገናን ሰርቪስን'
WHERE id = 3293 AND parent_id = 3197;

UPDATE category_hierarchy SET 
    name_en = 'Towing & Breakdown Recovery',
    name_am = 'መጎተትና የአደጋ ጊዜ ድነት',
    name_om = 'Harkisa fi Baraarsa Balaa',
    name_ti = 'ምስሓብን ህጹጽ ረድኤትን'
WHERE id = 3299 AND parent_id = 3197;

UPDATE category_hierarchy SET 
    name_en = 'Car Wash & Auto Detailing',
    name_am = 'የመኪና ማጠቢያና ፖሊሽ',
    name_om = 'Dhiqqaafi Diiteeliingii Konkolaataa',
    name_ti = 'ምሕጻብን ምውቃብን መካይን'
WHERE id = 3304 AND parent_id = 3197;

-- Property & Real Estate (id 2758)
UPDATE category_hierarchy SET 
    name_en = 'Brokerage & Agency Services',
    name_am = 'የደላላና ወኪል አገልግሎቶች',
    name_om = 'Tajaajila Dilaalaa fi Bakka Bu''iinsaa',
    name_ti = 'ኣገልግሎት ደላሎን ወከልትን'
WHERE id = 2792 AND parent_id = 2758;

UPDATE category_hierarchy SET 
    name_en = 'Real Estate & Land Development',
    name_am = 'የሪል እስቴትና መሬት ልማት',
    name_om = 'Misooma Qabeenya fi Lafaa',
    name_ti = 'ልምዓት ሪል እስቴትን መሬትን'
WHERE id = 2797 AND parent_id = 2758;

UPDATE category_hierarchy SET 
    name_en = 'Title & Property Insurance',
    name_am = 'የካርታና የንብረት ኢንሹራንስ',
    name_om = 'Inshuraansii Kaartaa fi Qabeenyaa',
    name_ti = 'ኢንሹራንስ ካርታን ንብረትን'
WHERE id = 2802 AND parent_id = 2758;

UPDATE category_hierarchy SET 
    name_en = 'Legal, Title Deed & Notary',
    name_am = 'ህጋዊ፣ የባለቤትነት ካርታና ኖተሪ',
    name_om = 'Heera, Kaartaa Abbaa Qabeenyummaa',
    name_ti = 'ሕጋዊ ሰነዳትን ካርታን'
WHERE id = 2806 AND parent_id = 2758;

UPDATE category_hierarchy SET 
    name_en = 'Building Maintenance & Renovation',
    name_am = 'የህንፃ ጥገናና እድሳት',
    name_om = 'Suphaa fi Haaromsa Gamoo',
    name_ti = 'ጽገናን ሓድሽ ምግባርን ህንጻታት'
WHERE id = 2810 AND parent_id = 2758;

UPDATE category_hierarchy SET 
    name_en = 'Facility & Estate Management',
    name_am = 'የተቋማትና ግቢ አስተዳደር',
    name_om = 'Bulchiinsa Dhaabbilee fi Mooraa',
    name_ti = 'ምሕደራ ትካላትን ህንጻታትን'
WHERE id = 2814 AND parent_id = 2758;

UPDATE category_hierarchy SET 
    name_en = 'Valuation, Survey & Appraisal',
    name_am = 'የንብረት ግመታና ሰርቬይ',
    name_om = 'Tilmaama fi Sakatta''a Qabeenyaa',
    name_ti = 'ግምገማን መጽናዕትን ንብረት'
WHERE id = 2819 AND parent_id = 2758;

-- Construction & Infrastructure (id 1201)
UPDATE category_hierarchy SET 
    name_en = 'Heavy Equipment, Machinery & Cranes',
    name_am = 'ከባድ ማሽነሪዎችና ክሬኖች',
    name_om = 'Maashinoota Gurguddoo fi Kireenii',
    name_ti = 'ከበድቲ ማሽነሪታትን ክሬናትን'
WHERE id = 1225 AND parent_id = 1201;

UPDATE category_hierarchy SET 
    name_en = 'Hardware, Fasteners & Fixtures',
    name_am = 'የብረታ ብረት እቃዎችና ማሰሪያዎች',
    name_om = 'Meeshaalee Sibiilaa fi Hidhaa',
    name_ti = 'መሳርሒታት ብረትን መጣበቒታትን'
WHERE id = 1230 AND parent_id = 1201;

UPDATE category_hierarchy SET 
    name_en = 'Project Engineering & Supervision',
    name_am = 'የፕሮጀክት ምህንድስናና ቁጥጥር',
    name_om = 'Injinariingii fi To''annoo Pirojektii',
    name_ti = 'ምህንድስናን ቁጽጽርን ፕሮጀክት'
WHERE id = 1235 AND parent_id = 1201;

UPDATE category_hierarchy SET 
    name_en = 'Contractors & Trade Specialists',
    name_am = 'ኮንትራክተሮችና የሙያ ባለሙያዎች',
    name_om = 'Kentraaktaroota fi Ogeessota',
    name_ti = 'ተቛራጽትን ሰብ ሞያን'
WHERE id = 1240 AND parent_id = 1201;

-- Food & Beverage (id 1814)
UPDATE category_hierarchy SET 
    name_en = 'Delivery, Catering & Distribution',
    name_am = 'ምግብ ማድረስና የምግብ ዝግጅት አገልግሎት',
    name_om = 'Tajaajila Nyaata Dhiyeessuu fi Geessuu',
    name_ti = 'ምግቢ ምብጻሕን ምድላውን ኣገልግሎት'
WHERE id = 1845 AND parent_id = 1814;

UPDATE category_hierarchy SET 
    name_en = 'Commercial Processing & Packaging',
    name_am = 'የምግብ ማቀነባበሪያና ማሸጊያ ማሽኖች',
    name_om = 'Maashinoota Qophii fi Saamsama Nyaataa',
    name_ti = 'ማሽነሪታት ምድላውን ምዕሻግን ምግቢ'
WHERE id = 1850 AND parent_id = 1814;

UPDATE category_hierarchy SET 
    name_en = 'Laboratory Testing & Certification',
    name_am = 'የላቦራቶሪ ምርመራና የጥራት ማረጋገጫ',
    name_om = 'Qorannoo Laabooratorii fi Qulqullina',
    name_ti = 'ናይ ላቦራቶሪ መርመራን ውሕስነት ጽሬትን'
WHERE id = 1855 AND parent_id = 1814;

UPDATE category_hierarchy SET 
    name_en = 'Culinary Ingredients & Supplies',
    name_am = 'የምግብ ዝግጅት ግብዓቶችና ቅመማ ቅመሞች',
    name_om = 'Galteewwan Nyaataa fi Mi''eessituu',
    name_ti = 'እታዎታት ምግቢን ቀመማትን'
WHERE id = 1860 AND parent_id = 1814;

-- Industrial Manufacturing (id 2314)
UPDATE category_hierarchy SET 
    name_en = 'Compressors, Vacuum & Industrial Pumps',
    name_am = 'ኮምፕረሰሮችና የኢንዱስትሪ ፓምፖች',
    name_om = 'Kompireesara fi Paampii Industirii',
    name_ti = 'ኮምፕረሰራትን ናይ ኢንዱስትሪ ፓምፓታትን'
WHERE id = 2335 AND parent_id = 2314;

UPDATE category_hierarchy SET 
    name_en = 'Industrial Engineering & Automation',
    name_am = 'የኢንዱስትሪ ምህንድስናና አውቶሜሽን',
    name_om = 'Injinariingii fi Ootomeeshinii Industirii',
    name_ti = 'ምህንድስናን ኦቶሜሽንን ኢንዱስትሪ'
WHERE id = 2340 AND parent_id = 2314;

UPDATE category_hierarchy SET 
    name_en = 'Power Transmission, Gears & Motors',
    name_am = 'ጊርስ፣ ሞተሮችና ኃይል ማስተላለፊያ',
    name_om = 'Giirii, Mootoroota fi Daddabarsa Humnaa',
    name_ti = 'ጊርስ፣ ሞተራትን መመሓላለፊ ሓይልን'
WHERE id = 2345 AND parent_id = 2314;

UPDATE category_hierarchy SET 
    name_en = 'Industrial Safety, Workwear & PPE',
    name_am = 'የስራ ላይ ደህንነት አልባሳትና መጠበቂያዎች',
    name_om = 'Uffata Nageenya Hojii fi PPE',
    name_ti = 'ናይ ስራሕ ድሕንነት ክዳውንትን መከላከልን'
WHERE id = 2350 AND parent_id = 2314;

UPDATE category_hierarchy SET 
    name_en = 'Precision Testing & Metrology',
    name_am = 'ትክክለኛ የሙከራና የመለኪያ መሳሪያዎች',
    name_om = 'Meeshaalee Safaraa fi Qorannoo',
    name_ti = 'መሳርሒታት መለክዕን መርመራን'
WHERE id = 2355 AND parent_id = 2314;

-- Mining & Extractive Industries (id 2497)
UPDATE category_hierarchy SET 
    name_en = 'Extraction Machinery & Drilling Rigs',
    name_am = 'የቁፋሮ ማሽኖችና ድሪሎች',
    name_om = 'Maashinoota Qonnaa fi Boorii',
    name_ti = 'ማሽነሪታት ቁፋሮን መቆፈሪታትን'
WHERE id = 2518 AND parent_id = 2497;

UPDATE category_hierarchy SET 
    name_en = 'Underground & Surface Mine Safety',
    name_am = 'የከርሰ ምድርና የውጭ ደህንነት እቃዎች',
    name_om = 'Nageenya Qotinsa Lafaa fi Irraa',
    name_ti = 'ድሕንነት መዓድን ውሽጢ መሬትን ላዕልን'
WHERE id = 2523 AND parent_id = 2497;

UPDATE category_hierarchy SET 
    name_en = 'Geological Surveying & Exploration',
    name_am = 'የጂኦሎጂ ጥናትና ፍለጋ',
    name_om = 'Qorannoo Ji''ooloojii fi Sakatta''a',
    name_ti = 'ጂኦሎጂካዊ መጽናዕትን ምድላይን'
WHERE id = 2528 AND parent_id = 2497;

-- Office & Workplace Supplies (id 2579)
UPDATE category_hierarchy SET 
    name_en = 'Paper, Inks & Printing Consumables',
    name_am = 'ወረቀቶች፣ ቀለሞችና ማተሚያ ግብዓቶች',
    name_om = 'Waraqaa, Qalama fi Galtee Maxxansaa',
    name_ti = 'ወረቐት፣ ቐለማትን ናይ ሕትመት እታዎታትን'
WHERE id = 2585 AND parent_id = 2579;

UPDATE category_hierarchy SET 
    name_en = 'Printers, Shredders & Machines',
    name_am = 'ፕሪንተሮችና የቢሮ ማሽኖች',
    name_om = 'Piriintara fi Maashinoota Biiroo',
    name_ti = 'ፕሪንተራትን ናይ ቤት ጽሕፈት ማሽናትን'
WHERE id = 2590 AND parent_id = 2579;

UPDATE category_hierarchy SET 
    name_en = 'Desks, Ergonomic Chairs & Cabinets',
    name_am = 'ጠረጴዛዎች፣ ወንበሮችና ሰነድ መደርደሪያዎች',
    name_om = 'Meezii, Barcuma fi Saanduqa Galmee',
    name_ti = 'ጣውላታት፣ ኮፍመበሊታትን ናይ ሰነድ መኽዘን'
WHERE id = 2595 AND parent_id = 2579;

UPDATE category_hierarchy SET 
    name_en = 'Writing, Desk & Filing Supplies',
    name_am = 'የመፃፊያና የዴስክ እቃዎች',
    name_om = 'Meeshaalee Barreeffamaa fi Minjaalaa',
    name_ti = 'መሳርሒታት ጽሕፈትን ሰሌዳታትን'
WHERE id = 2600 AND parent_id = 2579;

-- Education & Training (id 1303)
UPDATE category_hierarchy SET 
    name_en = 'Academic Advisory & Admissions',
    name_am = 'የትምህርት ምደባና የማማከር አገልግሎት',
    name_om = 'Gorsa Barnootaa fi Seensaa',
    name_ti = 'ናይ ትምህርቲ ምኽርን ምዝገባን'
WHERE id = 1318 AND parent_id = 1303;

UPDATE category_hierarchy SET 
    name_en = 'Classroom, Lab & Learning Materials',
    name_am = 'የክፍል፣ የላብና የመማሪያ ቁሳቁሶች',
    name_om = 'Meeshaalee Kutaa fi Qorannoo',
    name_ti = 'ናይ ክፍሊ፣ ላቦራቶሪን መምሃሪ ናውትን'
WHERE id = 1322 AND parent_id = 1303;

-- Water & Sanitation (id 3350)
UPDATE category_hierarchy SET 
    name_en = 'Piping, Bowsers & Distribution',
    name_am = 'የቦቴ ውሃና የቧንቧ ስርጭት',
    name_om = 'Raabsa Bishaan Boozarii fi Ujummoo',
    name_ti = 'ስርጭት ቦቴ ማይን ሻምብቆታትን'
WHERE id = 3376 AND parent_id = 3350;

UPDATE category_hierarchy SET 
    name_en = 'Microbiological & Purity Testing',
    name_am = 'የውሃ ንፅህናና የላቦራቶሪ ምርመራ',
    name_om = 'Qorannoo Qulqullina Bishaanii',
    name_ti = 'መርመራ ጽሬት ማይን ላቦራቶሪን'
WHERE id = 3381 AND parent_id = 3350;

UPDATE category_hierarchy SET 
    name_en = 'Boreholes, Tanks & Treatment Plants',
    name_am = 'ቦረሆል፣ ታንከሮችና የማጣሪያ ተቋማት',
    name_om = 'Boorii, Taankii fi Dhaabbata Qulqulleessaa',
    name_ti = 'ዒላታት፣ ታንክታትን መጻረዪ ማይን'
WHERE id = 3386 AND parent_id = 3350;

UPDATE category_hierarchy SET 
    name_en = 'Filters, Reverse Osmosis & Chemicals',
    name_am = 'ማጣሪያዎችና የውሃ ኬሚካሎች',
    name_om = 'Fiiltara fi Keemikaala Bishaan Qulqulleessu',
    name_ti = 'መጻረዪታትን ከሚካላትን'
WHERE id = 3391 AND parent_id = 3350;


-- 2. Drop Trigram Indexes Temporarily for Lightning-Fast Batch Update
DROP INDEX IF EXISTS idx_cat_hierarchy_name_en_trgm;
DROP INDEX IF EXISTS idx_cat_hierarchy_name_am_trgm;
DROP INDEX IF EXISTS idx_category_hierarchy_name_am_trgm;
DROP INDEX IF EXISTS idx_category_hierarchy_name_om_trgm;
DROP INDEX IF EXISTS idx_category_hierarchy_name_ti_trgm;

-- 3. Strip Father Prefixes Across Deep Layers (Level >= 2)
UPDATE category_hierarchy
SET 
  name_en = CASE 
    WHEN name_en ~ '^\d+.*?\s+-\s+\d+' AND name_en LIKE '% - % - %' 
    THEN TRIM(split_part(name_en, ' - ', 3))
    ELSE TRIM(split_part(name_en, ' - ', 2))
  END,
  name_am = CASE 
    WHEN name_am ~ '^\d+.*?\s+-\s+\d+' AND name_am LIKE '% - % - %' 
    THEN TRIM(split_part(name_am, ' - ', 3))
    WHEN name_am LIKE '% - %'
    THEN TRIM(split_part(name_am, ' - ', 2))
    ELSE name_am 
  END,
  name_om = CASE 
    WHEN name_om ~ '^\d+.*?\s+-\s+\d+' AND name_om LIKE '% - % - %' 
    THEN TRIM(split_part(name_om, ' - ', 3))
    WHEN name_om LIKE '% - %'
    THEN TRIM(split_part(name_om, ' - ', 2))
    ELSE name_om 
  END,
  name_ti = CASE 
    WHEN name_ti ~ '^\d+.*?\s+-\s+\d+' AND name_ti LIKE '% - % - %' 
    THEN TRIM(split_part(name_ti, ' - ', 3))
    WHEN name_ti LIKE '% - %'
    THEN TRIM(split_part(name_ti, ' - ', 2))
    ELSE name_ti 
  END
WHERE level >= 2 AND name_en LIKE '% - %' AND is_active = true;

-- 4. Recreate Trigram Indexes for Sub-Millisecond Search Acceleration
CREATE INDEX IF NOT EXISTS idx_cat_hierarchy_name_en_trgm ON category_hierarchy USING gin (name_en gin_trgm_ops) WHERE (is_active = true);
CREATE INDEX IF NOT EXISTS idx_cat_hierarchy_name_am_trgm ON category_hierarchy USING gin (name_am gin_trgm_ops) WHERE (is_active = true);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_name_am_trgm ON category_hierarchy USING gin (name_am gin_trgm_ops) WHERE (name_am IS NOT NULL);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_name_om_trgm ON category_hierarchy USING gin (name_om gin_trgm_ops) WHERE (name_om IS NOT NULL);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_name_ti_trgm ON category_hierarchy USING gin (name_ti gin_trgm_ops) WHERE (name_ti IS NOT NULL);

