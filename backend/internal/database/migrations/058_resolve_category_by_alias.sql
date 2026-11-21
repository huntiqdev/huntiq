-- ============================================================================
-- 058_resolve_category_by_alias.sql
-- Missing Stored Function & Multilingual Core Category Aliases
-- ============================================================================

-- 1. Create or replace public.resolve_category_by_alias function
CREATE OR REPLACE FUNCTION public.resolve_category_by_alias(p_term text)
RETURNS integer
LANGUAGE plpgsql
STABLE
AS $$
DECLARE
    v_clean text;
    v_cat_id integer;
BEGIN
    IF p_term IS NULL THEN
        RETURN NULL;
    END IF;
    
    v_clean := TRIM(LOWER(p_term));
    IF v_clean = '' THEN
        RETURN NULL;
    END IF;

    -- 1. Exact match in category_aliases
    SELECT category_id INTO v_cat_id
    FROM category_aliases
    WHERE LOWER(alias) = v_clean
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 2. Exact match in category_hierarchy (slug or localized names)
    SELECT id INTO v_cat_id
    FROM category_hierarchy
    WHERE is_active = true
      AND (
          LOWER(slug) = v_clean
          OR LOWER(name_en) = v_clean
          OR LOWER(COALESCE(name_am, '')) = v_clean
          OR LOWER(COALESCE(name_om, '')) = v_clean
          OR LOWER(COALESCE(name_ti, '')) = v_clean
      )
    ORDER BY level ASC, display_order ASC
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 3. Check category_redirects
    SELECT destination_id INTO v_cat_id
    FROM category_redirects
    WHERE LOWER(source_slug) = v_clean
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 4. Trigram / Substring match in category_aliases
    SELECT category_id INTO v_cat_id
    FROM category_aliases
    WHERE alias ILIKE '%' || v_clean || '%'
    LIMIT 1;

    IF v_cat_id IS NOT NULL THEN
        RETURN v_cat_id;
    END IF;

    -- 5. Substring match in category_hierarchy name
    SELECT id INTO v_cat_id
    FROM category_hierarchy
    WHERE is_active = true
      AND (
          name_en ILIKE '%' || v_clean || '%'
          OR COALESCE(name_am, '') LIKE '%' || p_term || '%'
      )
    ORDER BY level ASC
    LIMIT 1;

    RETURN v_cat_id;
END;
$$;

-- 2. Seed High-Frequency Ethiopian Commercial Category Aliases
INSERT INTO category_aliases (category_id, alias, language)
VALUES
    -- Electronics & Phones (Category 21)
    (21, 'phone', 'en'),
    (21, 'phones', 'en'),
    (21, 'mobile', 'en'),
    (21, 'mobiles', 'en'),
    (21, 'smartphone', 'en'),
    (21, 'smartphones', 'en'),
    (21, 'cellphone', 'en'),
    (21, 'cellphones', 'en'),
    (21, 'tablet', 'en'),
    (21, 'tablets', 'en'),
    (21, 'ipad', 'en'),
    (21, 'iphone', 'en'),
    (21, 'samsung', 'en'),
    (21, 'ስልክ', 'am'),
    (21, 'ስልኮች', 'am'),
    (21, 'ስማርትፎን', 'am'),
    (21, 'ስማርትፎኖች', 'am'),
    (21, 'ሞባይል', 'am'),
    (21, 'ታብሌት', 'am'),
    (21, 'bilbila', 'or'),
    (21, 'bilbilootaa', 'or'),
    (21, 'ስልኪ', 'ti'),
    (21, 'ስልክታት', 'ti'),

    -- Vehicles & Automotive (Category 76 / 5)
    (76, 'car', 'en'),
    (76, 'cars', 'en'),
    (76, 'vehicle', 'en'),
    (76, 'vehicles', 'en'),
    (76, 'automobile', 'en'),
    (76, 'automobiles', 'en'),
    (76, 'motorcycle', 'en'),
    (76, 'motorcycles', 'en'),
    (76, 'bajaj', 'en'),
    (76, 'truck', 'en'),
    (76, 'trucks', 'en'),
    (76, 'መኪና', 'am'),
    (76, 'መኪኖች', 'am'),
    (76, 'ተሽከርካሪ', 'am'),
    (76, 'ተሽከርካሪዎች', 'am'),
    (76, 'ባጃጅ', 'am'),
    (76, 'ሞተርሳይክል', 'am'),
    (76, 'konkolaataa', 'or'),
    (76, 'መኪናታት', 'ti'),

    -- Property & Accommodation (Category 60 / 1)
    (60, 'property', 'en'),
    (60, 'house', 'en'),
    (60, 'houses', 'en'),
    (60, 'apartment', 'en'),
    (60, 'apartments', 'en'),
    (60, 'condo', 'en'),
    (60, 'condos', 'en'),
    (60, 'condominium', 'en'),
    (60, 'villa', 'en'),
    (60, 'villas', 'en'),
    (60, 'real estate', 'en'),
    (60, 'land', 'en'),
    (60, 'plot', 'en'),
    (60, 'ቤት', 'am'),
    (60, 'ቤቶች', 'am'),
    (60, 'ኮንዶሚኒየም', 'am'),
    (60, 'ቪላ', 'am'),
    (60, 'አፓርታማ', 'am'),
    (60, 'ቦታ', 'am'),
    (60, 'መሬት', 'am'),
    (60, 'mana', 'or'),
    (60, 'ገዛ', 'ti'),
    (60, 'ገዛውቲ', 'ti'),

    -- Computers & Laptops (Category 18)
    (18, 'laptop', 'en'),
    (18, 'computer', 'en'),
    (18, 'macbook', 'en'),
    (18, 'desktop', 'en'),
    (18, 'monitor', 'en'),
    (18, 'ኮምፒውተር', 'am'),
    (18, 'ላፕቶፕ', 'am'),
    (18, 'ዴስክቶፕ', 'am'),
    (18, 'kompiitara', 'or'),

    -- Clothing & Fashion (Category 15)
    (15, 'clothes', 'en'),
    (15, 'clothing', 'en'),
    (15, 'shoes', 'en'),
    (15, 'dress', 'en'),
    (15, 'suit', 'en'),
    (15, 'jacket', 'en'),
    (15, 'አልባሳት', 'am'),
    (15, 'ቀሚስ', 'am'),
    (15, 'የሀበሻ ቀሚስ', 'am'),
    (15, 'ጫማ', 'am'),
    (15, 'ጃኬት', 'am'),
    (15, 'ሱፍ', 'am'),
    (15, 'uffata', 'or'),
    (15, 'ክዳን', 'ti'),

    -- Books & Education (Category 9)
    (9, 'book', 'en'),
    (9, 'books', 'en'),
    (9, 'textbook', 'en'),
    (9, 'dictionary', 'en'),
    (9, 'መጽሐፍ', 'am'),
    (9, 'መጽሐፍት', 'am'),
    (9, 'መዝገበ ቃላት', 'am'),
    (9, 'ትምህርት', 'am'),
    (9, 'kitaaba', 'or'),
    (9, 'መጽሓፍ', 'ti'),

    -- Furniture & Home (Category 31)
    (31, 'furniture', 'en'),
    (31, 'sofa', 'en'),
    (31, 'bed', 'en'),
    (31, 'table', 'en'),
    (31, 'chair', 'en'),
    (31, 'wardrobe', 'en'),
    (31, 'ፈርኒቸር', 'am'),
    (31, 'ሶፋ', 'am'),
    (31, 'አልጋ', 'am'),
    (31, 'ወንበር', 'am'),
    (31, 'ጠረጴዛ', 'am'),
    (31, 'ቁምሳጥን', 'am'),
    (31, 'teessoo', 'or'),
    (31, 'ኮፍመቻ', 'ti'),

    -- Tools & Equipment (Category 71)
    (71, 'tool', 'en'),
    (71, 'tools', 'en'),
    (71, 'drill', 'en'),
    (71, 'grinder', 'en'),
    (71, 'generator', 'en'),
    (71, 'pump', 'en'),
    (71, 'መሳሪያ', 'am'),
    (71, 'መሳሪያዎች', 'am'),
    (71, 'ማሽነሪ', 'am'),
    (71, 'meeshaalee', 'or'),
    (71, 'መሳርሒ', 'ti')
ON CONFLICT DO NOTHING;
