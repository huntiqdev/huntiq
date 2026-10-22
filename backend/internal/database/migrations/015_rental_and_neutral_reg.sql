-- ============================================================
-- 015_rental_and_neutral_reg.sql
-- CONSOLIDATES: 025_rental + 028_seller_preference + 029_neutral_registration + 030_optimize_categories
-- PURPOSE: Complete marketplace evolution - rentals, neutral registration, consolidated categories
-- THIS IS THE FINAL MIGRATION - All modern platform features in one place
-- ============================================================

-- ══════════════════════════════════════════════════════════
-- PART A: RENTAL SYSTEM (from 025)
-- Extends products to support BOTH sale AND rental listings
-- ══════════════════════════════════════════════════════════

-- Create enums for listing types
DO $$ BEGIN
    CREATE TYPE listing_type AS ENUM ('sale', 'rent');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE rental_duration_unit AS ENUM ('hour', 'day', 'week', 'month', 'year');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- Extend products table with rental fields
DO $$
BEGIN
    -- listing_type: determines if product is for sale or rent
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='products' AND column_name='listing_type') THEN
        ALTER TABLE products ADD COLUMN listing_type listing_type NOT NULL DEFAULT 'sale';
    END IF;
    
    -- Rental-specific pricing fields
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='products' AND column_name='rental_price') THEN
        ALTER TABLE products 
            ADD COLUMN rental_price NUMERIC(12,2) CHECK (rental_price IS NULL OR (rental_price > 0 AND rental_price <= 999999999.99)),
            ADD COLUMN rental_duration_value INT CHECK (rental_duration_value IS NULL OR rental_duration_value > 0),
            ADD COLUMN rental_duration_unit rental_duration_unit,
            ADD COLUMN deposit_amount NUMERIC(12,2) DEFAULT 0 CHECK (deposit_amount >= 0);
    END IF;
    
    -- Rental availability dates
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='products' AND column_name='available_from') THEN
        ALTER TABLE products 
            ADD COLUMN available_from DATE,
            ADD COLUMN available_until DATE;
    END IF;
    
    -- Rental status tracking
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='products' AND column_name='is_currently_rented') THEN
        ALTER TABLE products 
            ADD COLUMN is_currently_rented BOOLEAN DEFAULT FALSE,
            ADD COLUMN current_tenant_id INT REFERENCES users(id) ON DELETE SET NULL,
            ADD COLUMN rented_until DATE;
    END IF;
END $$;

-- Add rental field constraints
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name='products' AND constraint_name='check_rental_fields'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT check_rental_fields CHECK (
            (listing_type = 'sale' AND rental_price IS NULL AND rental_duration_value IS NULL) 
            OR 
            (listing_type = 'rent' AND rental_price IS NOT NULL AND rental_duration_value IS NOT NULL)
        );
    END IF;
END $$;

-- Indexes for rental queries
CREATE INDEX IF NOT EXISTS idx_products_listing_type ON products (listing_type);
CREATE INDEX IF NOT EXISTS idx_products_rental_price ON products (rental_price) WHERE listing_type = 'rent';
CREATE INDEX IF NOT EXISTS idx_products_available_from ON products (available_from) WHERE listing_type = 'rent';
CREATE INDEX IF NOT EXISTS idx_products_is_rented ON products (is_currently_rented) WHERE listing_type = 'rent' AND is_currently_rented = FALSE;

COMMENT ON COLUMN products.listing_type IS 'Type of listing: sale (default) or rent';
COMMENT ON COLUMN products.rental_price IS 'Rental price per duration unit (only for rent listings)';
COMMENT ON COLUMN products.is_currently_rented IS 'Whether rental item is currently rented out';

-- ══════════════════════════════════════════════════════════
-- PART B: SELLER PREFERENCE (from 028)
-- Unified seller/landlord system - one role, filtered by preference
-- ══════════════════════════════════════════════════════════

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='seller_preference') THEN
        ALTER TABLE users 
            ADD COLUMN seller_preference VARCHAR(10) DEFAULT 'sale' 
            CHECK (seller_preference IN ('sale', 'rent', 'both'));
    END IF;
END $$;

-- Backfill existing sellers (backward compatibility)
UPDATE users 
SET seller_preference = 'sale' 
WHERE role = 'seller' AND seller_preference IS NULL;

-- Index for efficient filtering
CREATE INDEX IF NOT EXISTS idx_users_seller_preference 
    ON users(seller_preference) 
    WHERE role = 'seller';

COMMENT ON COLUMN users.seller_preference IS 
    'Seller listing preference: sale (sell products), rent (rent properties), both (can do either)';

-- ══════════════════════════════════════════════════════════
-- PART C: NEUTRAL REGISTRATION (from 029)
-- Users register WITHOUT choosing role, select demand after registration
-- ══════════════════════════════════════════════════════════

DO $$
BEGIN
    -- Add current_demand column (FIND or LIST)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='current_demand') THEN
        ALTER TABLE users ADD COLUMN current_demand VARCHAR(10);
        ALTER TABLE users ADD CONSTRAINT users_current_demand_check 
            CHECK (current_demand IN ('FIND', 'LIST') OR current_demand IS NULL);
    END IF;

    -- Add rental-specific columns (consolidated with Part A)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='products' AND column_name='rental_price_per_month') THEN
        ALTER TABLE products 
            ADD COLUMN rental_price_per_month NUMERIC(12,2),
            ADD COLUMN rental_min_duration_months INT,
            ADD COLUMN rental_max_duration_months INT,
            ADD COLUMN rental_deposit_amount NUMERIC(12,2),
            ADD COLUMN rental_terms TEXT,
            ADD COLUMN available_from_ts TIMESTAMPTZ;
    END IF;
END $$;

-- Constraints for rental fields (idempotent — wrapped in DO blocks)
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name='products' AND constraint_name='products_rental_price_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_rental_price_check
            CHECK (rental_price_per_month IS NULL OR (rental_price_per_month > 0 AND rental_price_per_month <= 999999999.99));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name='products' AND constraint_name='products_rental_min_duration_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_rental_min_duration_check
            CHECK (rental_min_duration_months IS NULL OR rental_min_duration_months >= 1);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name='products' AND constraint_name='products_rental_max_duration_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_rental_max_duration_check
            CHECK (rental_max_duration_months IS NULL OR rental_max_duration_months >= rental_min_duration_months);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name='products' AND constraint_name='products_rental_deposit_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_rental_deposit_check
            CHECK (rental_deposit_amount IS NULL OR rental_deposit_amount >= 0);
    END IF;
END $$;

-- Indexes for neutral registration flow
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_current_demand 
    ON users (current_demand) WHERE current_demand IS NOT NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_seller_listing_type 
    ON products (seller_id, listing_type);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_status_listing_type 
    ON products (status, listing_type) WHERE status = 'approved';

-- Materialized view for marketplace (optimized public queries)
CREATE MATERIALIZED VIEW IF NOT EXISTS marketplace_listings AS
SELECT 
    p.id, p.seller_id, p.title, p.price, p.rental_price_per_month,
    p.description, p.category_id, p.subcategory, p.item, p.condition,
    p.location_id, p.kebele, p.specific_area, p.latitude, p.longitude,
    p.listing_type, p.access_tier, p.seller_tier_snapshot, p.view_count,
    p.status, p.created_at, p.expires_at,
    u.name as seller_name, u.rating as seller_rating, u.business_name as seller_business_name
FROM products p
INNER JOIN users u ON p.seller_id = u.id
WHERE p.status = 'approved'
  AND (p.expires_at IS NULL OR p.expires_at > NOW())
  AND u.is_active = TRUE
  AND u.is_banned = FALSE;

CREATE INDEX IF NOT EXISTS idx_marketplace_listing_type ON marketplace_listings (listing_type);
CREATE INDEX IF NOT EXISTS idx_marketplace_category ON marketplace_listings (category_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_location ON marketplace_listings (location_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_marketplace_listings_id ON marketplace_listings (id);

COMMENT ON MATERIALIZED VIEW marketplace_listings IS 
    'Cached public marketplace listings (refresh via cron job every 5 minutes)';

-- ══════════════════════════════════════════════════════════
-- PART D: CATEGORY CONSOLIDATION (from 030)
-- Reduce categories from 23 to 10 core categories (56% reduction)
-- ══════════════════════════════════════════════════════════

-- Create category consolidation mapping
CREATE TEMP TABLE category_consolidation AS
SELECT 
    old_slug, 
    new_slug, 
    new_name,
    new_name_am
FROM (VALUES
    -- Core Categories (Keep these 10)
    ('electronics', 'electronics', 'Electronics', 'ኤሌክትሮኒክስ'),
    ('vehicles', 'vehicles', 'Vehicles', 'ተሽከርካሪዎች'),
    ('property', 'property', 'Property & Real Estate', 'ንብረት'),
    ('home_appliances', 'home', 'Home & Living', 'ቤት እና የኑሮ'),
    ('furniture', 'home', 'Home & Living', 'ቤት እና የኑሮ'),
    ('clothing', 'fashion', 'Fashion & Clothing', 'ፋሽን'),
    ('agriculture', 'agriculture', 'Agriculture & Farming', 'ግብርና'),
    ('animals', 'agriculture', 'Agriculture & Farming', 'ግብርና'),
    ('industry', 'business', 'Business & Industry', 'ንግድ'),
    ('construction', 'business', 'Business & Industry', 'ንግድ'),
    ('tools', 'business', 'Business & Industry', 'ንግድ'),
    ('services', 'services', 'Services & Jobs', 'አገልግሎቶች'),
    ('jobs', 'services', 'Services & Jobs', 'አገልግሎቶች'),
    ('food', 'lifestyle', 'Lifestyle & Wellness', 'የአኗኗር ዘይቤ'),
    ('health_beauty', 'lifestyle', 'Lifestyle & Wellness', 'የአኗኗር ዘይቤ'),
    ('sports', 'lifestyle', 'Lifestyle & Wellness', 'የአኗኗር ዘይቤ'),
    ('education', 'education', 'Education & Learning', 'ትምህርት'),
    ('entertainment', 'education', 'Education & Learning', 'ትምህርት'),
    ('business', 'business', 'Business & Industry', 'ንግድ'),
    ('luxury', 'fashion', 'Fashion & Clothing', 'ፋሽን'),
    ('misc', 'other', 'Other', 'ሌላ'),
    ('other', 'other', 'Other', 'ሌላ')
) AS mapping(old_slug, new_slug, new_name, new_name_am);

-- Create consolidated categories
INSERT INTO categories (name, slug, name_en, name_am, name_or, name_ti, sort_order, is_active)
VALUES 
    ('Electronics', 'electronics', 'Electronics', 'ኤሌክትሮኒክስ', 'Elektirooniksii', 'ኤሌክትሮኒክስ', 1, true),
    ('Vehicles', 'vehicles', 'Vehicles', 'ተሽከርካሪዎች', 'Konkolaataa', 'መኪናታት', 2, true),
    ('Property & Real Estate', 'property', 'Property & Real Estate', 'ንብረት', 'Qabeenya', 'ንብረት', 3, true),
    ('Home & Living', 'home', 'Home & Living', 'ቤት እና የኑሮ', 'Mana fi Jireeny', 'ቤትን ህይወት', 4, true),
    ('Fashion & Clothing', 'fashion', 'Fashion & Clothing', 'ፋሽን', 'Faashinii', 'ፋሽን', 5, true),
    ('Agriculture & Farming', 'agriculture', 'Agriculture & Farming', 'ግብርና', 'Qonna', 'ሓረስታይ', 6, true),
    ('Business & Industry', 'business', 'Business & Industry', 'ንግድ', 'Daldala', 'ንግዲ', 7, true),
    ('Services & Jobs', 'services', 'Services & Jobs', 'አገልግሎቶች', 'Tajaajila', 'ኣገልግሎት', 8, true),
    ('Lifestyle & Wellness', 'lifestyle', 'Lifestyle & Wellness', 'የአኗኗር ዘይቤ', 'Akkaataa Jireenyaa', 'ኣነባብራ', 9, true),
    ('Education & Learning', 'education', 'Education & Learning', 'ትምህርት', 'Barnoot', 'ትምህርቲ', 10, true),
    ('Other', 'other', 'Other', 'ሌላ', 'Kan biraa', 'ካልእ', 99, true)
ON CONFLICT (slug) DO UPDATE SET
    name = EXCLUDED.name,
    name_en = EXCLUDED.name_en,
    name_am = EXCLUDED.name_am,
    is_active = EXCLUDED.is_active;

-- Migrate products to consolidated categories
DO $$
DECLARE
    old_cat_id INT;
    new_cat_id INT;
    old_slug_var VARCHAR;
    new_slug_var VARCHAR;
BEGIN
    FOR old_slug_var, new_slug_var IN 
        SELECT DISTINCT old_slug, new_slug FROM category_consolidation WHERE old_slug != new_slug
    LOOP
        SELECT id INTO old_cat_id FROM categories WHERE slug = old_slug_var LIMIT 1;
        SELECT id INTO new_cat_id FROM categories WHERE slug = new_slug_var LIMIT 1;
        
        IF old_cat_id IS NOT NULL AND new_cat_id IS NOT NULL THEN
            UPDATE products 
            SET category_id = new_cat_id, 
                subcategory = COALESCE(subcategory, old_slug_var)
            WHERE category_id = old_cat_id;
            
            RAISE NOTICE 'Migrated category % to %', old_slug_var, new_slug_var;
        END IF;
    END LOOP;
END $$;

-- Deactivate old categories (preserve data integrity)
UPDATE categories 
SET is_active = FALSE 
WHERE slug NOT IN (
    'electronics', 'vehicles', 'property', 'home', 'fashion', 
    'agriculture', 'business', 'services', 'lifestyle', 'education', 'other'
);

-- Analytics table for category migration
CREATE TABLE IF NOT EXISTS category_migration_log (
    id SERIAL PRIMARY KEY,
    old_category_slug VARCHAR(100),
    new_category_slug VARCHAR(100),
    products_migrated INT DEFAULT 0,
    migrated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO category_migration_log (old_category_slug, new_category_slug, products_migrated)
SELECT 
    cc.old_slug,
    cc.new_slug,
    COUNT(p.id)
FROM category_consolidation cc
LEFT JOIN categories old_cat ON old_cat.slug = cc.old_slug
LEFT JOIN products p ON p.category_id = old_cat.id
WHERE cc.old_slug != cc.new_slug
GROUP BY cc.old_slug, cc.new_slug;

-- Optimize category queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_categories_active ON categories (is_active) WHERE is_active = TRUE;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_categories_slug_active ON categories (slug, is_active);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_category_status ON products (category_id, status);

-- Update statistics
ANALYZE categories;
ANALYZE products;

-- ══════════════════════════════════════════════════════════
-- MIGRATION COMPLETE - SUMMARY
-- ══════════════════════════════════════════════════════════

DO $$
DECLARE
    total_active_cats INT;
    total_inactive_cats INT;
    total_products INT;
    total_users INT;
    rental_products INT;
BEGIN
    SELECT COUNT(*) INTO total_active_cats FROM categories WHERE is_active = TRUE;
    SELECT COUNT(*) INTO total_inactive_cats FROM categories WHERE is_active = FALSE;
    SELECT COUNT(*) INTO total_products FROM products;
    SELECT COUNT(*) INTO total_users FROM users;
    SELECT COUNT(*) INTO rental_products FROM products WHERE listing_type = 'rent';
    
    RAISE NOTICE '========================================';
    RAISE NOTICE 'MIGRATION 015 COMPLETE!';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Active categories: %', total_active_cats;
    RAISE NOTICE 'Inactive categories: %', total_inactive_cats;
    RAISE NOTICE 'Total products: %', total_products;
    RAISE NOTICE 'Rental listings: %', rental_products;
    RAISE NOTICE 'Total users: %', total_users;
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Features: Rentals ✓ Neutral Reg ✓ Consolidated Categories ✓';
END $$;
