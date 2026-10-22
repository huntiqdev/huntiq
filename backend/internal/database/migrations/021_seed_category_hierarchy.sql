-- 021_seed_category_hierarchy.sql
-- PURPOSE: Populate standard level-1 subcategories and map attributes to categories.

-- PART 1: Insert standard subcategories (Level 1)
DO $$
DECLARE
    v_electronics_id INT;
    v_vehicles_id INT;
    v_property_id INT;
    v_agriculture_id INT;
    v_fashion_id INT;
    v_home_id INT;
    v_services_id INT;
    v_business_id INT;
BEGIN
    SELECT id INTO v_electronics_id FROM category_hierarchy WHERE slug = 'electronics' AND level = 0;
    SELECT id INTO v_vehicles_id FROM category_hierarchy WHERE slug = 'vehicles' AND level = 0;
    SELECT id INTO v_property_id FROM category_hierarchy WHERE slug = 'property' AND level = 0;
    SELECT id INTO v_agriculture_id FROM category_hierarchy WHERE slug = 'agriculture' AND level = 0;
    SELECT id INTO v_fashion_id FROM category_hierarchy WHERE slug = 'fashion' AND level = 0;
    SELECT id INTO v_home_id FROM category_hierarchy WHERE slug = 'home' AND level = 0;
    SELECT id INTO v_services_id FROM category_hierarchy WHERE slug = 'services' AND level = 0;
    SELECT id INTO v_business_id FROM category_hierarchy WHERE slug = 'business' AND level = 0;

    -- 1. Electronics Subcategories
    IF v_electronics_id IS NOT NULL THEN
        INSERT INTO category_hierarchy (parent_id, slug, name_en, name_am, name_om, is_leaf, display_order) VALUES
        (v_electronics_id, 'mobile-phones', 'Phones & Tablets', 'ስልኮችና ታብሌቶች', 'Bilbilootaa fi Taableetota', true, 1),
        (v_electronics_id, 'computers-laptops', 'Computers & Laptops', 'ኮምፒውተሮችና ላፕቶፖች', 'Kompitaroota', true, 2),
        (v_electronics_id, 'tv-audio', 'TV, Audio & Video', 'ቴሌቪዥንና ድምፅ', 'TV fi Sagalee', true, 3),
        (v_electronics_id, 'cameras', 'Cameras & Photography', 'ካሜራዎች', 'Kaameeraa', true, 4),
        (v_electronics_id, 'accessories', 'Accessories', 'ተጓዳኝ እቃዎች', 'Meeshaalee Dabalataa', true, 5)
        ON CONFLICT (slug) WHERE is_active = true DO NOTHING;
    END IF;

    -- 2. Vehicles Subcategories
    IF v_vehicles_id IS NOT NULL THEN
        INSERT INTO category_hierarchy (parent_id, slug, name_en, name_am, name_om, is_leaf, display_order) VALUES
        (v_vehicles_id, 'cars', 'Cars & Automobiles', 'መኪኖች', 'Konkolaattota', true, 1),
        (v_vehicles_id, 'trucks', 'Trucks & Heavy Vehicles', 'የጭነት መኪኖች', 'Gaanfa Konkolaataa', true, 2),
        (v_vehicles_id, 'motorcycles', 'Motorcycles & Bajaj', 'ሞተርሳይክልና ባጃጅ', 'Mootoraa fi Baajaajii', true, 3),
        (v_vehicles_id, 'spare-parts', 'Spare Parts', 'የመለዋወጫ እቃዎች', 'Qaama Bakka Buusaa', true, 4)
        ON CONFLICT (slug) WHERE is_active = true DO NOTHING;
    END IF;

    -- 3. Property Subcategories
    IF v_property_id IS NOT NULL THEN
        INSERT INTO category_hierarchy (parent_id, slug, name_en, name_am, name_om, is_leaf, display_order) VALUES
        (v_property_id, 'houses-sale', 'Houses for Sale', 'የሚሸጡ ቤቶች', 'Manneen Gurguraa', true, 1),
        (v_property_id, 'houses-rent', 'Houses for Rent', 'የሚከራዩ ቤቶች', 'Manneen Kiraa', true, 2),
        (v_property_id, 'apartments', 'Apartments & Condos', 'ኮንዶሚኒየም', 'Appaartaamaa', true, 3),
        (v_property_id, 'commercial-property', 'Commercial Property', 'የንግድ ቤቶች', 'Daldala', true, 4),
        (v_property_id, 'land-plots', 'Land & Plots', 'ቦታዎች', 'Lafa', true, 5)
        ON CONFLICT (slug) WHERE is_active = true DO NOTHING;
    END IF;

    -- 4. Agriculture Subcategories
    IF v_agriculture_id IS NOT NULL THEN
        INSERT INTO category_hierarchy (parent_id, slug, name_en, name_am, name_om, is_leaf, display_order) VALUES
        (v_agriculture_id, 'crops-coffee', 'Crops & Coffee', 'ሰብሎችና ቡና', 'Midhaan fi Buna', true, 1),
        (v_agriculture_id, 'livestock', 'Livestock & Cattle', 'የከብት እርባታ', 'Beeylada', true, 2),
        (v_agriculture_id, 'tractors-equipment', 'Tractors & Farm Machinery', 'ትራክተሮችና ማሽነሪዎች', 'Tiraaktaraa fi Maashinoota', true, 3)
        ON CONFLICT (slug) WHERE is_active = true DO NOTHING;
    END IF;

    -- 5. Fashion Subcategories
    IF v_fashion_id IS NOT NULL THEN
        INSERT INTO category_hierarchy (parent_id, slug, name_en, name_am, name_om, is_leaf, display_order) VALUES
        (v_fashion_id, 'mens-clothing', 'Men''s Fashion', 'የወንዶች ልብስ', 'Uffata Dhiiraa', true, 1),
        (v_fashion_id, 'womens-clothing', 'Women''s Fashion', 'የሴቶች ልብስ', 'Uffata Dubartootaa', true, 2),
        (v_fashion_id, 'traditional-habesha', 'Traditional Habesha Attire', 'የሀበሻ ባህላዊ ልብስ', 'Uffata Aadaa', true, 3),
        (v_fashion_id, 'shoes', 'Shoes & Footwear', 'ጫማዎች', 'Kophee', true, 4)
        ON CONFLICT (slug) WHERE is_active = true DO NOTHING;
    END IF;

    -- 6. Home Subcategories
    IF v_home_id IS NOT NULL THEN
        INSERT INTO category_hierarchy (parent_id, slug, name_en, name_am, name_om, is_leaf, display_order) VALUES
        (v_home_id, 'furniture', 'Furniture & Living', 'የቤት እቃዎችና ፈርኒቸር', 'Meeshaa Manoo', true, 1),
        (v_home_id, 'kitchen-appliances', 'Kitchen Appliances', 'የወጥ ቤት እቃዎች', 'Meeshaa Kushiinaa', true, 2)
        ON CONFLICT (slug) WHERE is_active = true DO NOTHING;
    END IF;

    -- 7. Services Subcategories
    IF v_services_id IS NOT NULL THEN
        INSERT INTO category_hierarchy (parent_id, slug, name_en, name_am, name_om, is_leaf, display_order) VALUES
        (v_services_id, 'professional-services', 'Professional Services', 'ሙያዊ አገልግሎቶች', 'Tajaajila Ogeessummaa', true, 1),
        (v_services_id, 'repair-maintenance', 'Repair & Maintenance', 'ጥገናና ዕድሳት', 'Suphaa', true, 2)
        ON CONFLICT (slug) WHERE is_active = true DO NOTHING;
    END IF;

    -- Mark root categories as non-leaf
    UPDATE category_hierarchy SET is_leaf = false WHERE level = 0;
END $$;

-- PART 2: Map Attributes to Categories (category_attributes)

-- Electronics: mobile-phones
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT ch.id, ad.id, CASE WHEN ad.slug IN ('brand', 'model') THEN true ELSE false END, ad.display_order
FROM category_hierarchy ch CROSS JOIN attribute_definitions ad
WHERE ch.slug = 'mobile-phones' AND ad.slug IN ('brand', 'model', 'storage', 'ram', 'screen-size', 'battery', 'color')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Electronics: computers-laptops
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT ch.id, ad.id, CASE WHEN ad.slug IN ('brand', 'model') THEN true ELSE false END, ad.display_order
FROM category_hierarchy ch CROSS JOIN attribute_definitions ad
WHERE ch.slug = 'computers-laptops' AND ad.slug IN ('brand', 'model', 'storage', 'ram', 'screen-size', 'processor')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Vehicles: cars
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT ch.id, ad.id, CASE WHEN ad.slug IN ('make', 'vehicle-model', 'year') THEN true ELSE false END, ad.display_order
FROM category_hierarchy ch CROSS JOIN attribute_definitions ad
WHERE ch.slug = 'cars' AND ad.slug IN ('make', 'vehicle-model', 'year', 'mileage', 'fuel-type', 'transmission', 'engine-size', 'doors')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Vehicles: trucks
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT ch.id, ad.id, CASE WHEN ad.slug IN ('make', 'vehicle-model', 'year') THEN true ELSE false END, ad.display_order
FROM category_hierarchy ch CROSS JOIN attribute_definitions ad
WHERE ch.slug = 'trucks' AND ad.slug IN ('make', 'vehicle-model', 'year', 'mileage', 'fuel-type', 'transmission')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Property: houses-sale, houses-rent, apartments
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT ch.id, ad.id, CASE WHEN ad.slug IN ('property-type', 'bedrooms', 'area-size') THEN true ELSE false END, ad.display_order
FROM category_hierarchy ch CROSS JOIN attribute_definitions ad
WHERE ch.slug IN ('houses-sale', 'houses-rent', 'apartments') AND ad.slug IN ('property-type', 'bedrooms', 'bathrooms', 'area-size', 'furnished', 'parking', 'floor')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Agriculture: tractors-equipment
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT ch.id, ad.id, CASE WHEN ad.slug IN ('brand', 'equipment-type') THEN true ELSE false END, ad.display_order
FROM category_hierarchy ch CROSS JOIN attribute_definitions ad
WHERE ch.slug = 'tractors-equipment' AND ad.slug IN ('brand', 'equipment-type', 'horsepower', 'working-hours', 'year')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Agriculture: livestock
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT ch.id, ad.id, CASE WHEN ad.slug IN ('livestock-type') THEN true ELSE false END, ad.display_order
FROM category_hierarchy ch CROSS JOIN attribute_definitions ad
WHERE ch.slug = 'livestock' AND ad.slug IN ('livestock-type', 'age')
ON CONFLICT (category_id, attribute_id) DO NOTHING;
