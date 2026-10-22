-- ============================================================================
-- Migration 017: Dynamic Category Hierarchy System
-- ============================================================================
-- Purpose: Upgrade flat category system to support unlimited hierarchy depth
-- Date: October 21, 2026
-- Status: Safe migration - preserves existing categories table
-- ============================================================================

-- ============================================================================
-- PART 1: Create New Category Hierarchy Table
-- ============================================================================

CREATE TABLE IF NOT EXISTS category_hierarchy (
    id SERIAL PRIMARY KEY,
    parent_id INT REFERENCES category_hierarchy(id) ON DELETE RESTRICT,
    
    -- Multi-language support for Ethiopian marketplace
    name_en VARCHAR(200) NOT NULL,
    name_am VARCHAR(200),  -- Amharic name
    name_om VARCHAR(200),  -- Afaan Oromo name
    
    slug VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- Hierarchy metadata
    level INT NOT NULL DEFAULT 0,  -- 0 = root, 1 = child of root, etc.
    path VARCHAR(500),  -- Materialized path: /1/5/12/ for fast ancestor queries
    
    -- Display and ordering
    display_order INT DEFAULT 0,
    icon VARCHAR(100),  -- Optional icon name/class
    
    -- Status and settings
    is_active BOOLEAN DEFAULT true,
    is_leaf BOOLEAN DEFAULT false,  -- Can products be attached to this category?
    
    -- SEO
    meta_title VARCHAR(200),
    meta_description TEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT no_circular_reference CHECK (id != parent_id),
    CONSTRAINT valid_level CHECK (level >= 0 AND level < 20)
);

-- ============================================================================
-- PART 2: Indexes for Performance
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_parent ON category_hierarchy(parent_id);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_level ON category_hierarchy(level);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_slug ON category_hierarchy(slug);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_path ON category_hierarchy(path);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_active ON category_hierarchy(is_active);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_leaf ON category_hierarchy(is_leaf);
CREATE INDEX IF NOT EXISTS idx_category_hierarchy_display_order ON category_hierarchy(display_order);

-- Unique constraint: slug must be unique per level to avoid confusion
CREATE UNIQUE INDEX IF NOT EXISTS idx_category_hierarchy_slug_unique ON category_hierarchy(slug) WHERE is_active = true;

-- ============================================================================
-- PART 3: Migrate Existing Categories
-- ============================================================================

-- Insert root-level categories from existing categories table
INSERT INTO category_hierarchy (
    name_en, 
    slug, 
    level, 
    path, 
    is_active, 
    is_leaf,
    display_order,
    created_at
)
SELECT 
    COALESCE(name_en, name),
    slug,
    0,  -- Root level
    '/' || CAST(id AS VARCHAR) || '/',  -- Path: /1/, /2/, etc.
    COALESCE(is_active, true),
    false,  -- Root categories are not leaves initially
    ROW_NUMBER() OVER (ORDER BY id),
    NOW()
FROM categories c
WHERE c.id IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM category_hierarchy ch WHERE ch.slug = c.slug)
ORDER BY c.id;

-- ============================================================================
-- PART 4: Add Category Hierarchy Link to Products
-- ============================================================================

-- Add new column to products table (nullable for safe migration)
ALTER TABLE products ADD COLUMN IF NOT EXISTS category_hierarchy_id INT;

-- Add foreign key constraint safely
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_products_category_hierarchy'
    ) THEN
        ALTER TABLE products 
            ADD CONSTRAINT fk_products_category_hierarchy 
            FOREIGN KEY (category_hierarchy_id) 
            REFERENCES category_hierarchy(id) 
            ON DELETE RESTRICT;
    END IF;
END $$;

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_products_category_hierarchy ON products(category_hierarchy_id);

-- Migrate existing product category links
UPDATE products p
SET category_hierarchy_id = ch.id
FROM category_hierarchy ch, categories c
WHERE p.category_id = c.id
  AND ch.name_en = c.name
  AND ch.level = 0;

-- ============================================================================
-- PART 5: Parse Existing Subcategory/Item into Hierarchy
-- ============================================================================

-- Function to get or create category in hierarchy
CREATE OR REPLACE FUNCTION get_or_create_category(
    p_parent_id INT,
    p_name VARCHAR,
    p_slug VARCHAR,
    p_level INT
) RETURNS INT AS $$
DECLARE
    v_id INT;
    v_path VARCHAR;
BEGIN
    -- Check if category exists
    SELECT id INTO v_id
    FROM category_hierarchy
    WHERE parent_id = p_parent_id 
      AND name_en = p_name;
    
    IF v_id IS NOT NULL THEN
        RETURN v_id;
    END IF;
    
    -- Get parent path
    IF p_parent_id IS NULL THEN
        v_path := '/';
    ELSE
        SELECT path || CAST(p_parent_id AS VARCHAR) || '/' 
        INTO v_path
        FROM category_hierarchy
        WHERE id = p_parent_id;
    END IF;
    
    -- Create new category
    INSERT INTO category_hierarchy (
        parent_id, 
        name_en, 
        slug, 
        level, 
        path, 
        is_leaf
    ) VALUES (
        p_parent_id,
        p_name,
        p_slug,
        p_level,
        v_path,
        true  -- Leaf by default
    )
    RETURNING id INTO v_id;
    
    -- Mark parent as non-leaf
    IF p_parent_id IS NOT NULL THEN
        UPDATE category_hierarchy 
        SET is_leaf = false 
        WHERE id = p_parent_id;
    END IF;
    
    RETURN v_id;
END;
$$ LANGUAGE plpgsql;

-- Migrate subcategories (level 1)
DO $$
DECLARE
    r RECORD;
    v_subcategory_id INT;
    v_slug VARCHAR;
BEGIN
    FOR r IN (
        SELECT DISTINCT 
            p.category_id,
            p.subcategory,
            ch.id as parent_hierarchy_id
        FROM products p
        JOIN categories c ON c.id = p.category_id
        JOIN category_hierarchy ch ON ch.name_en = c.name AND ch.level = 0
        WHERE p.subcategory IS NOT NULL 
          AND TRIM(p.subcategory) != ''
    ) LOOP
        -- Generate slug from subcategory name
        v_slug := LOWER(REGEXP_REPLACE(TRIM(r.subcategory), '[^a-zA-Z0-9]+', '-', 'g'));
        v_slug := REGEXP_REPLACE(v_slug, '^-+|-+$', '', 'g');
        
        -- Create subcategory
        v_subcategory_id := get_or_create_category(
            r.parent_hierarchy_id,
            TRIM(r.subcategory),
            v_slug,
            1
        );
        
        -- Update products to point to subcategory
        UPDATE products 
        SET category_hierarchy_id = v_subcategory_id
        WHERE category_id = r.category_id 
          AND subcategory = r.subcategory
          AND category_hierarchy_id IS NULL;
    END LOOP;
END $$;

-- Migrate items (level 2)
DO $$
DECLARE
    r RECORD;
    v_item_id INT;
    v_slug VARCHAR;
BEGIN
    FOR r IN (
        SELECT DISTINCT 
            p.category_id,
            p.subcategory,
            p.item,
            ch.id as subcategory_hierarchy_id
        FROM products p
        JOIN categories c ON c.id = p.category_id
        JOIN category_hierarchy parent ON parent.name_en = c.name AND parent.level = 0
        JOIN category_hierarchy ch ON ch.parent_id = parent.id 
                                   AND ch.name_en = TRIM(p.subcategory) 
                                   AND ch.level = 1
        WHERE p.item IS NOT NULL 
          AND TRIM(p.item) != ''
          AND p.subcategory IS NOT NULL
    ) LOOP
        -- Generate slug from item name
        v_slug := LOWER(REGEXP_REPLACE(TRIM(r.item), '[^a-zA-Z0-9]+', '-', 'g'));
        v_slug := REGEXP_REPLACE(v_slug, '^-+|-+$', '', 'g');
        
        -- Create item category
        v_item_id := get_or_create_category(
            r.subcategory_hierarchy_id,
            TRIM(r.item),
            v_slug,
            2
        );
        
        -- Update products to point to item
        UPDATE products 
        SET category_hierarchy_id = v_item_id
        WHERE category_id = r.category_id 
          AND subcategory = r.subcategory
          AND item = r.item
          AND category_hierarchy_id IS NULL;
    END LOOP;
END $$;

-- Set remaining products to their root category
UPDATE products p
SET category_hierarchy_id = ch.id
FROM category_hierarchy ch, categories c
WHERE p.category_id = c.id
  AND ch.name_en = c.name
  AND ch.level = 0
  AND p.category_hierarchy_id IS NULL;

-- ============================================================================
-- PART 6: Helper Views
-- ============================================================================

-- View: Category with breadcrumb path
CREATE OR REPLACE VIEW category_breadcrumbs AS
WITH RECURSIVE category_path AS (
    -- Base case: root categories
    SELECT 
        id,
        parent_id,
        name_en,
        slug,
        level,
        path,
        ARRAY[name_en::TEXT] as breadcrumb_names,
        ARRAY[slug::TEXT] as breadcrumb_slugs,
        ARRAY[id] as breadcrumb_ids
    FROM category_hierarchy
    WHERE parent_id IS NULL
    
    UNION ALL
    
    -- Recursive case: build path upward
    SELECT 
        ch.id,
        ch.parent_id,
        ch.name_en,
        ch.slug,
        ch.level,
        ch.path,
        cp.breadcrumb_names || ch.name_en::TEXT,
        cp.breadcrumb_slugs || ch.slug::TEXT,
        cp.breadcrumb_ids || ch.id
    FROM category_hierarchy ch
    JOIN category_path cp ON ch.parent_id = cp.id
)
SELECT 
    id,
    name_en,
    slug,
    level,
    breadcrumb_names,
    breadcrumb_slugs,
    breadcrumb_ids,
    ARRAY_LENGTH(breadcrumb_names, 1) as depth
FROM category_path;

-- View: Category with product count
DROP VIEW IF EXISTS category_stats CASCADE;
CREATE VIEW category_stats AS
SELECT 
    ch.id,
    ch.id as category_id,
    ch.parent_id,
    ch.name_en,
    ch.slug,
    ch.level,
    ch.is_leaf,
    COUNT(p.id) as product_count,
    COUNT(p.id) FILTER (WHERE p.status = 'approved') as approved_count,
    COUNT(p.id) FILTER (WHERE p.status = 'pending') as pending_count
FROM category_hierarchy ch
LEFT JOIN products p ON p.category_hierarchy_id = ch.id
GROUP BY ch.id, ch.parent_id, ch.name_en, ch.slug, ch.level, ch.is_leaf;

-- ============================================================================
-- PART 7: Helper Functions
-- ============================================================================

-- Function: Get all ancestors of a category
CREATE OR REPLACE FUNCTION get_category_ancestors(p_category_id INT)
RETURNS TABLE (
    id INT,
    name_en VARCHAR,
    slug VARCHAR,
    level INT
) AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE ancestors AS (
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        WHERE ch.id = p_category_id
        
        UNION ALL
        
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        JOIN ancestors a ON ch.id = a.parent_id
    )
    SELECT a.id, a.name_en, a.slug, a.level
    FROM ancestors a
    ORDER BY a.level;
END;
$$ LANGUAGE plpgsql;

-- Function: Get all descendants of a category
CREATE OR REPLACE FUNCTION get_category_descendants(p_category_id INT)
RETURNS TABLE (
    id INT,
    name_en VARCHAR,
    slug VARCHAR,
    level INT
) AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE descendants AS (
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        WHERE ch.id = p_category_id
        
        UNION ALL
        
        SELECT 
            ch.id,
            ch.parent_id,
            ch.name_en,
            ch.slug,
            ch.level
        FROM category_hierarchy ch
        JOIN descendants d ON ch.parent_id = d.id
    )
    SELECT d.id, d.name_en, d.slug, d.level
    FROM descendants d
    ORDER BY d.level;
END;
$$ LANGUAGE plpgsql;

-- Function: Get category full path name
CREATE OR REPLACE FUNCTION get_category_full_name(p_category_id INT)
RETURNS VARCHAR AS $$
DECLARE
    v_full_name VARCHAR;
BEGIN
    SELECT ARRAY_TO_STRING(breadcrumb_names, ' > ')
    INTO v_full_name
    FROM category_breadcrumbs
    WHERE id = p_category_id;
    
    RETURN COALESCE(v_full_name, '');
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- PART 8: Update Triggers
-- ============================================================================

-- Trigger: Update path and level when parent changes
CREATE OR REPLACE FUNCTION update_category_hierarchy_path()
RETURNS TRIGGER AS $$
DECLARE
    v_parent_path VARCHAR;
    v_parent_level INT;
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.level := 0;
        NEW.path := '/' || CAST(NEW.id AS VARCHAR) || '/';
    ELSE
        SELECT level, path
        INTO v_parent_level, v_parent_path
        FROM category_hierarchy
        WHERE id = NEW.parent_id;
        
        NEW.level := v_parent_level + 1;
        NEW.path := v_parent_path || CAST(NEW.id AS VARCHAR) || '/';
    END IF;
    
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_category_hierarchy_path
    BEFORE INSERT OR UPDATE OF parent_id ON category_hierarchy
    FOR EACH ROW
    EXECUTE FUNCTION update_category_hierarchy_path();

-- ============================================================================
-- PART 9: Validation and Cleanup
-- ============================================================================

-- Ensure all products have a category_hierarchy_id
DO $$
DECLARE
    v_unlinked_count INT;
BEGIN
    SELECT COUNT(*) INTO v_unlinked_count
    FROM products
    WHERE category_hierarchy_id IS NULL;
    
    IF v_unlinked_count > 0 THEN
        RAISE WARNING 'Found % products without category_hierarchy_id. Please review.', v_unlinked_count;
    END IF;
END $$;

-- Drop the helper function (no longer needed after migration)
DROP FUNCTION IF EXISTS get_or_create_category(INT, VARCHAR, VARCHAR, INT);

-- ============================================================================
-- PART 10: Comments for Documentation
-- ============================================================================

COMMENT ON TABLE category_hierarchy IS 'Dynamic N-level category hierarchy supporting Ethiopian marketplace';
COMMENT ON COLUMN category_hierarchy.path IS 'Materialized path for fast ancestor/descendant queries. Format: /1/5/12/';
COMMENT ON COLUMN category_hierarchy.is_leaf IS 'True if products can be directly attached to this category';
COMMENT ON COLUMN category_hierarchy.level IS 'Depth in tree: 0=root, 1=child of root, etc.';

COMMENT ON VIEW category_breadcrumbs IS 'Category with full breadcrumb path for navigation';
COMMENT ON VIEW category_stats IS 'Category with product counts for analytics';

COMMENT ON FUNCTION get_category_ancestors(INT) IS 'Returns all ancestors of a category from root to parent';
COMMENT ON FUNCTION get_category_descendants(INT) IS 'Returns all descendants of a category including itself';
COMMENT ON FUNCTION get_category_full_name(INT) IS 'Returns full category name with breadcrumbs: Electronics > Phones > Samsung';

-- ============================================================================
-- Migration Complete ✅
-- ============================================================================

