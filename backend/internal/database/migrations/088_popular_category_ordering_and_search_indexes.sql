-- ============================================================================
-- Migration 088: Popular Category Ordering, Deep Hierarchy Search & Mock Data Purge
-- ============================================================================
-- 1. Purges mock/sample products and demo sellers so marketplace begins clean.
-- 2. Establishes commercial popularity display_order for all 50 root economic sectors.
-- 3. Creates trigram & prefix search acceleration indexes on category_hierarchy
--    for sub-millisecond hierarchical breadcrumb matching.
-- ============================================================================

-- 1. Clean Mock / Sample Data
TRUNCATE TABLE products CASCADE;
DELETE FROM users WHERE phone IN ('+251911111111', '+251922222222', '+251933333333', '+251944444444');

-- 2. Popularity-Driven Root Category Order (Commercial & High-Demand First)
UPDATE category_hierarchy SET display_order = 1 WHERE id = 1 AND level = 0;    -- Electronics
UPDATE category_hierarchy SET display_order = 2 WHERE id = 3197 AND level = 0; -- Transport & Mobility (Vehicles / Automotive)
UPDATE category_hierarchy SET display_order = 3 WHERE id = 2758 AND level = 0; -- Property & Real Estate (Rent & Sale)
UPDATE category_hierarchy SET display_order = 4 WHERE id = 924 AND level = 0;  -- Clothing & Apparel
UPDATE category_hierarchy SET display_order = 5 WHERE id = 2212 AND level = 0; -- Home Appliances & Kitchenware
UPDATE category_hierarchy SET display_order = 6 WHERE id = 1130 AND level = 0; -- Computing & Information Technology
UPDATE category_hierarchy SET display_order = 7 WHERE id = 1814 AND level = 0; -- Food & Beverage
UPDATE category_hierarchy SET display_order = 8 WHERE id = 1201 AND level = 0; -- Construction & Infrastructure
UPDATE category_hierarchy SET display_order = 9 WHERE id = 2314 AND level = 0; -- Industrial Manufacturing
UPDATE category_hierarchy SET display_order = 10 WHERE id = 4 AND level = 0;  -- Agriculture & Food
UPDATE category_hierarchy SET display_order = 11 WHERE id = 2024 AND level = 0; -- Furniture & Home Furnishings
UPDATE category_hierarchy SET display_order = 12 WHERE id = 1932 AND level = 0; -- Footwear & Leather Goods
UPDATE category_hierarchy SET display_order = 13 WHERE id = 1651 AND level = 0; -- Fashion Accessories & Jewelry
UPDATE category_hierarchy SET display_order = 14 WHERE id = 690 AND level = 0;  -- Beauty & Personal Care
UPDATE category_hierarchy SET display_order = 15 WHERE id = 1452 AND level = 0; -- Energy & Utilities
UPDATE category_hierarchy SET display_order = 16 WHERE id = 3002 AND level = 0; -- Sports & Fitness
UPDATE category_hierarchy SET display_order = 17 WHERE id = 2121 AND level = 0; -- Healthcare & Medicine
UPDATE category_hierarchy SET display_order = 18 WHERE id = 2451 AND level = 0; -- Logistics & Supply Chain
UPDATE category_hierarchy SET display_order = 19 WHERE id = 972 AND level = 0;  -- Commerce & Trade
UPDATE category_hierarchy SET display_order = 20 WHERE id = 781 AND level = 0;  -- Business & Professional Services
UPDATE category_hierarchy SET display_order = 21 WHERE id = 2268 AND level = 0; -- Hospitality & Accommodation
UPDATE category_hierarchy SET display_order = 22 WHERE id = 3161 AND level = 0; -- Tourism & Travel
UPDATE category_hierarchy SET display_order = 23 WHERE id = 3396 AND level = 0; -- Workforce & Employment
UPDATE category_hierarchy SET display_order = 24 WHERE id = 1303 AND level = 0; -- Education & Training
UPDATE category_hierarchy SET display_order = 25 WHERE id = 3089 AND level = 0; -- Textiles & Fabrics
UPDATE category_hierarchy SET display_order = 26 WHERE id = 2579 AND level = 0; -- Office & Workplace Supplies
UPDATE category_hierarchy SET display_order = 27 WHERE id = 873 AND level = 0;  -- Cleaning & Facility Management
UPDATE category_hierarchy SET display_order = 28 WHERE id = 827 AND level = 0;  -- Chemicals & Raw Materials
UPDATE category_hierarchy SET display_order = 29 WHERE id = 2615 AND level = 0; -- Packaging & Printing
UPDATE category_hierarchy SET display_order = 30 WHERE id = 2497 AND level = 0; -- Mining & Extractive Industries
UPDATE category_hierarchy SET display_order = 31 WHERE id = 2966 AND level = 0; -- Security & Protective Services
UPDATE category_hierarchy SET display_order = 32 WHERE id = 1033 AND level = 0; -- Communications & Media
UPDATE category_hierarchy SET display_order = 33 WHERE id = 735 AND level = 0;  -- Books, Stationery & Publishing
UPDATE category_hierarchy SET display_order = 34 WHERE id = 608 AND level = 0;  -- Arts, Culture & Heritage
UPDATE category_hierarchy SET display_order = 35 WHERE id = 1513 AND level = 0; -- Entertainment & Recreation
UPDATE category_hierarchy SET display_order = 36 WHERE id = 1605 AND level = 0; -- Events & Exhibitions
UPDATE category_hierarchy SET display_order = 37 WHERE id = 562 AND level = 0;  -- Animal Care & Veterinary
UPDATE category_hierarchy SET display_order = 38 WHERE id = 1763 AND level = 0; -- Fishing & Aquaculture
UPDATE category_hierarchy SET display_order = 39 WHERE id = 1973 AND level = 0; -- Forestry & Timber
UPDATE category_hierarchy SET display_order = 40 WHERE id = 1554 AND level = 0; -- Environment & Conservation
UPDATE category_hierarchy SET display_order = 41 WHERE id = 3350 AND level = 0; -- Water & Sanitation
UPDATE category_hierarchy SET display_order = 42 WHERE id = 3309 AND level = 0; -- Waste Management & Recycling
UPDATE category_hierarchy SET display_order = 43 WHERE id = 1697 AND level = 0; -- Finance & Insurance
UPDATE category_hierarchy SET display_order = 44 WHERE id = 2405 AND level = 0; -- Legal & Regulatory Services
UPDATE category_hierarchy SET display_order = 45 WHERE id = 2661 AND level = 0; -- Personal & Household Services
UPDATE category_hierarchy SET display_order = 46 WHERE id = 1084 AND level = 0; -- Community & Social Development
UPDATE category_hierarchy SET display_order = 47 WHERE id = 2843 AND level = 0; -- Public Safety & Emergency Services
UPDATE category_hierarchy SET display_order = 48 WHERE id = 2075 AND level = 0; -- Government & Public Administration
UPDATE category_hierarchy SET display_order = 49 WHERE id = 2925 AND level = 0; -- Research & Scientific Services
UPDATE category_hierarchy SET display_order = 50 WHERE id = 2884 AND level = 0; -- Religion & Faith

-- 3. High-Speed Indexes for Hierarchical Search
CREATE INDEX IF NOT EXISTS idx_cat_hierarchy_name_en_trgm ON category_hierarchy USING gin (name_en gin_trgm_ops) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_cat_hierarchy_name_am_trgm ON category_hierarchy USING gin (name_am gin_trgm_ops) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_cat_hierarchy_slug_trgm ON category_hierarchy USING gin (slug gin_trgm_ops) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_cat_hierarchy_active_level ON category_hierarchy (is_active, level, display_order);
CREATE INDEX IF NOT EXISTS idx_cat_hierarchy_path_prefix ON category_hierarchy (path text_pattern_ops);

