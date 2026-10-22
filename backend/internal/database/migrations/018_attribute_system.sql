-- ============================================================================
-- Migration 018: Dynamic Attribute System
-- ============================================================================
-- Enables category-specific product attributes (e.g., Brand for phones, Make for cars)
-- 
-- Features:
-- - Define attributes globally (e.g., "Brand", "Storage", "Year")
-- - Assign attributes to categories
-- - Store product attribute values (flexible schema)
-- - Support multiple data types (text, number, boolean, select, range)
-- - Validation rules per attribute
-- - Filterable and searchable attributes
-- ============================================================================

-- Transaction managed by migration runner

-- ============================================================================
-- PART 1: Attribute Definitions
-- ============================================================================

-- Table: attribute_definitions
-- Stores global attribute templates (e.g., "Brand", "Model", "Storage")
CREATE TABLE IF NOT EXISTS attribute_definitions (
    id SERIAL PRIMARY KEY,
    name_en VARCHAR(100) NOT NULL,
    name_am VARCHAR(100),
    name_om VARCHAR(100),
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    
    -- Data type: text, number, boolean, select, multiselect, range, date
    data_type VARCHAR(50) NOT NULL DEFAULT 'text',
    
    -- Validation rules (JSON)
    -- Examples:
    -- {"min": 0, "max": 1000} for numbers
    -- {"options": ["new", "like-new", "good", "fair"]} for select
    -- {"pattern": "^[0-9]{4}$"} for regex validation
    validation_rules JSONB,
    
    -- UI configuration
    is_required BOOLEAN DEFAULT false,
    is_searchable BOOLEAN DEFAULT true,      -- Included in search index
    is_filterable BOOLEAN DEFAULT true,      -- Shows as filter in UI
    display_order INT DEFAULT 0,
    
    -- Unit of measurement (e.g., "GB", "inches", "kg")
    unit VARCHAR(50),
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_attribute_definitions_slug ON attribute_definitions(slug);
CREATE INDEX IF NOT EXISTS idx_attribute_definitions_data_type ON attribute_definitions(data_type);

-- ============================================================================
-- PART 2: Category-Attribute Mapping
-- ============================================================================

-- Table: category_attributes
-- Maps attributes to categories (e.g., Electronics -> Brand, Model)
CREATE TABLE IF NOT EXISTS category_attributes (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES category_hierarchy(id) ON DELETE CASCADE,
    attribute_id INT NOT NULL REFERENCES attribute_definitions(id) ON DELETE CASCADE,
    
    -- Override settings per category
    is_required BOOLEAN DEFAULT false,
    is_inherited BOOLEAN DEFAULT false,      -- Inherited from parent category
    display_order INT DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(category_id, attribute_id)
);

CREATE INDEX IF NOT EXISTS idx_category_attributes_category ON category_attributes(category_id);
CREATE INDEX IF NOT EXISTS idx_category_attributes_attribute ON category_attributes(attribute_id);

-- ============================================================================
-- PART 3: Product Attribute Values
-- ============================================================================

-- Table: product_attributes
-- Stores actual product attribute values (flexible schema)
CREATE TABLE IF NOT EXISTS product_attributes (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    attribute_id INT NOT NULL REFERENCES attribute_definitions(id) ON DELETE CASCADE,
    
    -- Value storage (use appropriate column based on data_type)
    value_text TEXT,
    value_number NUMERIC(15, 2),
    value_boolean BOOLEAN,
    value_date DATE,
    value_json JSONB,                        -- For multiselect, complex values
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(product_id, attribute_id)
);

CREATE INDEX IF NOT EXISTS idx_product_attributes_product ON product_attributes(product_id);
CREATE INDEX IF NOT EXISTS idx_product_attributes_attribute ON product_attributes(attribute_id);
CREATE INDEX IF NOT EXISTS idx_product_attributes_value_text ON product_attributes(value_text);
CREATE INDEX IF NOT EXISTS idx_product_attributes_value_number ON product_attributes(value_number);
CREATE INDEX IF NOT EXISTS idx_product_attributes_value_json ON product_attributes USING gin(value_json);

-- ============================================================================
-- PART 4: Full-Text Search on Attributes
-- ============================================================================

-- Add search vector column to product_attributes
ALTER TABLE product_attributes ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Update search vector for text attributes
CREATE INDEX IF NOT EXISTS idx_product_attributes_search ON product_attributes USING gin(search_vector);

-- Trigger: Update search vector on insert/update
CREATE OR REPLACE FUNCTION update_product_attribute_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.value_text IS NOT NULL THEN
        NEW.search_vector := to_tsvector('simple', COALESCE(NEW.value_text, ''));
    END IF;
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_update_product_attribute_search ON product_attributes;
CREATE TRIGGER trg_update_product_attribute_search
    BEFORE INSERT OR UPDATE ON product_attributes
    FOR EACH ROW
    EXECUTE FUNCTION update_product_attribute_search_vector();

-- ============================================================================
-- PART 5: Seed Common Attributes
-- ============================================================================

-- Electronics attributes
INSERT INTO attribute_definitions (name_en, name_am, name_om, slug, data_type, is_filterable, display_order)
VALUES 
    ('Brand', 'ብራንድ', 'Maqaa Daldalaa', 'brand', 'text', true, 1),
    ('Model', 'ሞዴል', 'Moodela', 'model', 'text', true, 2),
    ('Storage Capacity', 'የማከማቻ አቅም', 'Dandeettii Kuusaa', 'storage', 'select', true, 3),
    ('RAM', 'RAM', 'RAM', 'ram', 'select', true, 4),
    ('Screen Size', 'የስክሪን መጠን', 'Guddinaa Iskiriinii', 'screen-size', 'number', true, 5),
    ('Processor', 'ፕሮሰሰር', 'Proosesera', 'processor', 'text', true, 6),
    ('Color', 'ቀለም', 'Halluu', 'color', 'select', true, 7),
    ('Battery Capacity', 'የባትሪ አቅም', 'Dandeettii Baatirii', 'battery', 'number', true, 8)
ON CONFLICT (slug) DO NOTHING;

-- Update validation rules for select attributes
UPDATE attribute_definitions 
SET validation_rules = '{"options": ["16GB", "32GB", "64GB", "128GB", "256GB", "512GB", "1TB"]}'
WHERE slug = 'storage';

UPDATE attribute_definitions 
SET validation_rules = '{"options": ["2GB", "3GB", "4GB", "6GB", "8GB", "12GB", "16GB"]}'
WHERE slug = 'ram';

UPDATE attribute_definitions 
SET validation_rules = '{"options": ["Black", "White", "Blue", "Red", "Green", "Gold", "Silver", "Gray"]}'
WHERE slug = 'color';

UPDATE attribute_definitions 
SET unit = 'inches'
WHERE slug = 'screen-size';

UPDATE attribute_definitions 
SET unit = 'mAh'
WHERE slug = 'battery';

-- Vehicle attributes
INSERT INTO attribute_definitions (name_en, name_am, name_om, slug, data_type, is_filterable, display_order)
VALUES 
    ('Make', 'አምራች', 'Oomishtuu', 'make', 'text', true, 1),
    ('Vehicle Model', 'ሞዴል', 'Moodela Konkolaataa', 'vehicle-model', 'text', true, 2),
    ('Year', 'ዓመት', 'Bara', 'year', 'number', true, 3),
    ('Mileage', 'የተጓዘ ርቀት', 'Fageenya Deemame', 'mileage', 'number', true, 4),
    ('Fuel Type', 'የነዳጅ አይነት', 'Gosa Boba''aa', 'fuel-type', 'select', true, 5),
    ('Transmission', 'ማስተላለፊያ', 'Dabarsaa', 'transmission', 'select', true, 6),
    ('Engine Size', 'የሞተር መጠን', 'Guddinaa Mootora', 'engine-size', 'number', true, 7),
    ('Number of Doors', 'የበር ብዛት', 'Lakkoofsa Balbala', 'doors', 'number', true, 8)
ON CONFLICT (slug) DO NOTHING;

UPDATE attribute_definitions 
SET validation_rules = '{"min": 1900, "max": 2030}'
WHERE slug = 'year';

UPDATE attribute_definitions 
SET unit = 'km', validation_rules = '{"min": 0, "max": 1000000}'
WHERE slug = 'mileage';

UPDATE attribute_definitions 
SET validation_rules = '{"options": ["Petrol", "Diesel", "Electric", "Hybrid", "CNG"]}'
WHERE slug = 'fuel-type';

UPDATE attribute_definitions 
SET validation_rules = '{"options": ["Manual", "Automatic", "Semi-Automatic"]}'
WHERE slug = 'transmission';

UPDATE attribute_definitions 
SET unit = 'L'
WHERE slug = 'engine-size';

-- Property attributes
INSERT INTO attribute_definitions (name_en, name_am, name_om, slug, data_type, is_filterable, display_order)
VALUES 
    ('Property Type', 'የንብረት አይነት', 'Gosa Qabeenya', 'property-type', 'select', true, 1),
    ('Bedrooms', 'መኝታ ክፍሎች', 'Kutaa Ciisichaa', 'bedrooms', 'number', true, 2),
    ('Bathrooms', 'መታጠቢያ ክፍሎች', 'Kutaa Dhiqannaa', 'bathrooms', 'number', true, 3),
    ('Area Size', 'የስፋት መጠን', 'Bal''ina', 'area-size', 'number', true, 4),
    ('Furnished', 'የተዘጋጀ', 'Qophaa''e', 'furnished', 'select', true, 5),
    ('Parking', 'የመኪና ማቆሚያ', 'Iddoo Kaa''umsaa', 'parking', 'boolean', true, 6),
    ('Floor Number', 'ደረጃ ቁጥር', 'Lakkoofsa Sadarkaa', 'floor', 'number', true, 7)
ON CONFLICT (slug) DO NOTHING;

UPDATE attribute_definitions 
SET validation_rules = '{"options": ["Apartment", "House", "Villa", "Studio", "Commercial", "Land"]}'
WHERE slug = 'property-type';

UPDATE attribute_definitions 
SET unit = 'm²', validation_rules = '{"min": 1, "max": 100000}'
WHERE slug = 'area-size';

UPDATE attribute_definitions 
SET validation_rules = '{"options": ["Fully Furnished", "Semi Furnished", "Unfurnished"]}'
WHERE slug = 'furnished';

-- Agriculture attributes
INSERT INTO attribute_definitions (name_en, name_am, name_om, slug, data_type, is_filterable, display_order)
VALUES 
    ('Equipment Type', 'የመሳሪያ አይነት', 'Gosa Meeshaa', 'equipment-type', 'select', true, 1),
    ('Horsepower', 'ሆርስፓወር', 'Humna Fardaa', 'horsepower', 'number', true, 2),
    ('Working Hours', 'የሥራ ሰዓት', 'Sa''aatii Hojii', 'working-hours', 'number', true, 3),
    ('Livestock Type', 'የእንስሳት አይነት', 'Gosa Beelada', 'livestock-type', 'select', true, 4),
    ('Age', 'ዕድሜ', 'Umurii', 'age', 'number', true, 5)
ON CONFLICT (slug) DO NOTHING;

UPDATE attribute_definitions 
SET validation_rules = '{"options": ["Tractor", "Plough", "Harvester", "Seeder", "Sprayer", "Other"]}'
WHERE slug = 'equipment-type';

UPDATE attribute_definitions 
SET unit = 'HP'
WHERE slug = 'horsepower';

UPDATE attribute_definitions 
SET unit = 'hours'
WHERE slug = 'working-hours';

UPDATE attribute_definitions 
SET validation_rules = '{"options": ["Cattle", "Sheep", "Goats", "Poultry", "Camels", "Donkeys"]}'
WHERE slug = 'livestock-type';

UPDATE attribute_definitions 
SET unit = 'years', validation_rules = '{"min": 0, "max": 50}'
WHERE slug = 'age';

-- ============================================================================
-- PART 6: Map Attributes to Categories
-- ============================================================================

-- Electronics (Mobile Phones)
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT 
    ch.id,
    ad.id,
    CASE WHEN ad.slug IN ('brand', 'model') THEN true ELSE false END,
    ad.display_order
FROM category_hierarchy ch
CROSS JOIN attribute_definitions ad
WHERE ch.slug = 'mobile-phones' 
  AND ad.slug IN ('brand', 'model', 'storage', 'ram', 'screen-size', 'battery', 'color')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Electronics (Computers & Laptops)
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT 
    ch.id,
    ad.id,
    CASE WHEN ad.slug IN ('brand', 'model') THEN true ELSE false END,
    ad.display_order
FROM category_hierarchy ch
CROSS JOIN attribute_definitions ad
WHERE ch.slug = 'computers-laptops' 
  AND ad.slug IN ('brand', 'model', 'storage', 'ram', 'screen-size', 'processor')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Vehicles (Cars)
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT 
    ch.id,
    ad.id,
    CASE WHEN ad.slug IN ('make', 'vehicle-model', 'year') THEN true ELSE false END,
    ad.display_order
FROM category_hierarchy ch
CROSS JOIN attribute_definitions ad
WHERE ch.slug = 'cars' 
  AND ad.slug IN ('make', 'vehicle-model', 'year', 'mileage', 'fuel-type', 'transmission', 'engine-size', 'doors')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Property (Houses)
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT 
    ch.id,
    ad.id,
    CASE WHEN ad.slug IN ('property-type', 'bedrooms', 'area-size') THEN true ELSE false END,
    ad.display_order
FROM category_hierarchy ch
CROSS JOIN attribute_definitions ad
WHERE ch.slug IN ('houses-sale', 'houses-rent') 
  AND ad.slug IN ('property-type', 'bedrooms', 'bathrooms', 'area-size', 'furnished', 'parking', 'floor')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Agriculture (Tractors)
INSERT INTO category_attributes (category_id, attribute_id, is_required, display_order)
SELECT 
    ch.id,
    ad.id,
    CASE WHEN ad.slug IN ('brand', 'equipment-type') THEN true ELSE false END,
    ad.display_order
FROM category_hierarchy ch
CROSS JOIN attribute_definitions ad
WHERE ch.slug = 'tractors-equipment' 
  AND ad.slug IN ('brand', 'equipment-type', 'horsepower', 'working-hours', 'year')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ============================================================================
-- PART 7: Helper Views
-- ============================================================================

-- View: Category with applicable attributes
CREATE OR REPLACE VIEW category_applicable_attributes AS
SELECT 
    ca.category_id,
    ch.name_en as category_name,
    ch.slug as category_slug,
    ad.id as attribute_id,
    ad.name_en as attribute_name,
    ad.slug as attribute_slug,
    ad.data_type,
    ad.validation_rules,
    ad.unit,
    ca.is_required,
    ca.display_order,
    ad.is_filterable,
    ad.is_searchable
FROM category_attributes ca
JOIN category_hierarchy ch ON ch.id = ca.category_id
JOIN attribute_definitions ad ON ad.id = ca.attribute_id
ORDER BY ca.category_id, ca.display_order;

-- ============================================================================
-- PART 8: Helper Functions
-- ============================================================================

-- Function: Get attributes for a category (including inherited)
CREATE OR REPLACE FUNCTION get_category_attributes(p_category_id INT)
RETURNS TABLE (
    attribute_id INT,
    name_en VARCHAR,
    slug VARCHAR,
    data_type VARCHAR,
    validation_rules JSONB,
    unit VARCHAR,
    is_required BOOLEAN,
    is_filterable BOOLEAN,
    display_order INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT
        ad.id,
        ad.name_en,
        ad.slug,
        ad.data_type,
        ad.validation_rules,
        ad.unit,
        ca.is_required,
        ad.is_filterable,
        ca.display_order
    FROM category_attributes ca
    JOIN attribute_definitions ad ON ad.id = ca.attribute_id
    WHERE ca.category_id = p_category_id
    ORDER BY ca.display_order;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- PART 9: Update Triggers
-- ============================================================================

-- Trigger: Update timestamp on attribute definition change
CREATE OR REPLACE FUNCTION update_attribute_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_update_attribute_definitions_timestamp ON attribute_definitions;
CREATE TRIGGER trg_update_attribute_definitions_timestamp
    BEFORE UPDATE ON attribute_definitions
    FOR EACH ROW
    EXECUTE FUNCTION update_attribute_timestamp();

-- ============================================================================
-- PART 10: Comments
-- ============================================================================

COMMENT ON TABLE attribute_definitions IS 'Global attribute templates (Brand, Model, etc.)';
COMMENT ON TABLE category_attributes IS 'Maps attributes to categories';
COMMENT ON TABLE product_attributes IS 'Stores actual product attribute values';
COMMENT ON COLUMN product_attributes.value_json IS 'For multiselect attributes (stored as JSON array)';
COMMENT ON FUNCTION get_category_attributes(INT) IS 'Returns all applicable attributes for a category';

-- ============================================================================
-- Migration 018 Complete ✅
-- ============================================================================
