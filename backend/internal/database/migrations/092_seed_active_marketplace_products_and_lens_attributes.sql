-- ============================================================================
-- 092_seed_active_marketplace_products_and_lens_attributes.sql
-- Seeds Active Marketplace Listings and Connects Technical Attributes
--
-- Adds high-quality, verified listings across key categories:
--   1. Camera Lenses (Sony FE 24-70mm GM II, Canon RF 50mm STM, Sony FE 70-200mm GM II)
--   2. Digital Cameras (Sony Alpha 7 IV Mirrorless, Canon EOS 5D Mark IV DSLR)
--   3. Smartphones (Apple iPhone 15 Pro Max, Samsung Galaxy S24 Ultra)
--   4. Computers & Laptops (Apple MacBook Pro 14 M3 Pro)
--   5. Vehicles (Toyota Vitz 2018 Clean Condition)
--   6. Real Estate (Bole Atlas Furnished Apartment for Rent)
--   7. Tools & Hardware (Bosch GBH 2-26 DRE Rotary Hammer SDS-Plus)
--
-- Each product is linked with:
--   - Real high-resolution product images in product_images
--   - Accurate technical specifications in product_attributes
--   - Status 'approved', fee_paid true, valid 180-day expiration
-- ============================================================================

BEGIN;

DO $$
DECLARE
    v_seller_id INT := 17;
    v_loc_bole INT := 1;
    v_pid INT;
BEGIN
    -- Ensure test seller exists and is active
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = v_seller_id AND is_active = true) THEN
        SELECT id INTO v_seller_id FROM users WHERE role = 'seller' AND is_active = true LIMIT 1;
        IF v_seller_id IS NULL THEN
            SELECT id INTO v_seller_id FROM users WHERE is_active = true LIMIT 1;
        END IF;
        IF v_seller_id IS NULL THEN
            INSERT INTO users (name, phone, role, is_active, is_verified)
            VALUES ('HUNTIQ Verified Seller', '+251911999888', 'seller', true, true)
            RETURNING id INTO v_seller_id;
        END IF;
    END IF;

    -- Ensure Bole location exists
    IF NOT EXISTS (SELECT 1 FROM locations WHERE id = v_loc_bole) THEN
        SELECT id INTO v_loc_bole FROM locations LIMIT 1;
    END IF;

    -- ========================================================================
    -- 1. SONY FE 24-70MM F/2.8 GM II LENS (CATEGORY 1012880)
    -- ========================================================================
    IF NOT EXISTS (SELECT 1 FROM products WHERE title = 'Sony FE 24-70mm f/2.8 GM II Lens') THEN
        INSERT INTO products (
            seller_id, location_id, title, price, description, condition, status,
            listing_type, fee_paid, expires_at, category_hierarchy_id, view_count, contact_count
        ) VALUES (
            v_seller_id, v_loc_bole, 'Sony FE 24-70mm f/2.8 GM II Lens', 185000,
            'Mint condition Sony G Master II standard zoom lens. Ultra-sharp optics, fast quiet autofocus for video and photography. Includes original front/rear caps, lens hood, and padded pouch.',
            'like_new', 'approved', 'sale', true, NOW() + INTERVAL '180 days', 1012880, 48, 7
        ) RETURNING id INTO v_pid;

        -- Images
        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail)
        VALUES 
        (v_pid, 'https://images.unsplash.com/photo-1617005082133-548c4dd27f35?auto=format&fit=crop&w=800&q=80', 0, true),
        (v_pid, 'https://images.unsplash.com/photo-1516035069371-29a1b244cc32?auto=format&fit=crop&w=800&q=80', 1, false);

        -- Attributes
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Sony' FROM attribute_definitions WHERE slug = 'brand';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'FE 24-70mm f/2.8 GM II' FROM attribute_definitions WHERE slug = 'model';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Sony E-Mount' FROM attribute_definitions WHERE slug = 'lens-mount';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Standard Zoom' FROM attribute_definitions WHERE slug = 'lens-type';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, '24-70mm' FROM attribute_definitions WHERE slug = 'focal-length';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'f/2.8' FROM attribute_definitions WHERE slug = 'max-aperture';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Full Frame (35mm)' FROM attribute_definitions WHERE slug = 'camera-sensor-size';
        INSERT INTO product_attributes (product_id, attribute_id, value_boolean)
        SELECT v_pid, id, false FROM attribute_definitions WHERE slug = 'image-stabilization';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, '82mm' FROM attribute_definitions WHERE slug = 'filter-thread-size';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Clean Glass (No Scratch/Fungus/Haze)' FROM attribute_definitions WHERE slug = 'lens-optical-condition';
    END IF;

    -- ========================================================================
    -- 2. CANON RF 50MM F/1.8 STM COMPACT PRIME LENS (CATEGORY 1012887)
    -- ========================================================================
    IF NOT EXISTS (SELECT 1 FROM products WHERE title = 'Canon RF 50mm f/1.8 STM Compact Prime Lens') THEN
        INSERT INTO products (
            seller_id, location_id, title, price, description, condition, status,
            listing_type, fee_paid, expires_at, category_hierarchy_id, view_count, contact_count
        ) VALUES (
            v_seller_id, v_loc_bole, 'Canon RF 50mm f/1.8 STM Compact Prime Lens', 28500,
            'Brand new sealed Canon RF 50mm f/1.8 STM prime lens. Extremely lightweight walkaround lens with beautiful shallow depth of field and bokeh for Canon EOS R mirrorless cameras.',
            'new', 'approved', 'sale', true, NOW() + INTERVAL '180 days', 1012887, 62, 12
        ) RETURNING id INTO v_pid;

        -- Images
        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail)
        VALUES 
        (v_pid, 'https://images.unsplash.com/photo-1622434641406-a158123450f9?auto=format&fit=crop&w=800&q=80', 0, true),
        (v_pid, 'https://images.unsplash.com/photo-1502920917128-1aa500764cbd?auto=format&fit=crop&w=800&q=80', 1, false);

        -- Attributes
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Canon' FROM attribute_definitions WHERE slug = 'brand';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'RF 50mm f/1.8 STM' FROM attribute_definitions WHERE slug = 'model';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Canon RF' FROM attribute_definitions WHERE slug = 'lens-mount';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Prime Lens' FROM attribute_definitions WHERE slug = 'lens-type';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, '50mm' FROM attribute_definitions WHERE slug = 'focal-length';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'f/1.8' FROM attribute_definitions WHERE slug = 'max-aperture';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Full Frame (35mm)' FROM attribute_definitions WHERE slug = 'camera-sensor-size';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Clean Glass (No Scratch/Fungus/Haze)' FROM attribute_definitions WHERE slug = 'lens-optical-condition';
    END IF;

    -- ========================================================================
    -- 3. SONY FE 70-200MM F/2.8 GM OSS II TELEPHOTO LENS (CATEGORY 1012891)
    -- ========================================================================
    IF NOT EXISTS (SELECT 1 FROM products WHERE title = 'Sony FE 70-200mm f/2.8 GM OSS II Telephoto Zoom Lens') THEN
        INSERT INTO products (
            seller_id, location_id, title, price, description, condition, status,
            listing_type, fee_paid, expires_at, category_hierarchy_id, view_count, contact_count
        ) VALUES (
            v_seller_id, v_loc_bole, 'Sony FE 70-200mm f/2.8 GM OSS II Telephoto Zoom Lens', 245000,
            'Flagship telephoto zoom lens, 29% lighter than first generation. Built-in Optical SteadyShot, pristine condition glass, tripod collar, original hood and case included.',
            'like_new', 'approved', 'sale', true, NOW() + INTERVAL '180 days', 1012891, 71, 9
        ) RETURNING id INTO v_pid;

        -- Images
        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail)
        VALUES 
        (v_pid, 'https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?auto=format&fit=crop&w=800&q=80', 0, true),
        (v_pid, 'https://images.unsplash.com/photo-1512790182412-b19e6d62bc39?auto=format&fit=crop&w=800&q=80', 1, false);

        -- Attributes
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Sony' FROM attribute_definitions WHERE slug = 'brand';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'FE 70-200mm f/2.8 GM OSS II' FROM attribute_definitions WHERE slug = 'model';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Sony E-Mount' FROM attribute_definitions WHERE slug = 'lens-mount';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Telephoto Zoom' FROM attribute_definitions WHERE slug = 'lens-type';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, '70-200mm' FROM attribute_definitions WHERE slug = 'focal-length';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'f/2.8' FROM attribute_definitions WHERE slug = 'max-aperture';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Full Frame (35mm)' FROM attribute_definitions WHERE slug = 'camera-sensor-size';
        INSERT INTO product_attributes (product_id, attribute_id, value_boolean)
        SELECT v_pid, id, true FROM attribute_definitions WHERE slug = 'image-stabilization';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, '77mm' FROM attribute_definitions WHERE slug = 'filter-thread-size';
    END IF;

    -- ========================================================================
    -- 4. SONY ALPHA 7 IV MIRRORLESS CAMERA (CATEGORY 1012908)
    -- ========================================================================
    IF NOT EXISTS (SELECT 1 FROM products WHERE title = 'Sony Alpha 7 IV (A7 IV) Full Frame Mirrorless Camera Body') THEN
        INSERT INTO products (
            seller_id, location_id, title, price, description, condition, status,
            listing_type, fee_paid, expires_at, category_hierarchy_id, view_count, contact_count
        ) VALUES (
            v_seller_id, v_loc_bole, 'Sony Alpha 7 IV (A7 IV) Full Frame Mirrorless Camera Body', 195000,
            '33MP BSI CMOS full frame sensor, 4K 60p 10-bit video, S-Cinetone, 759-point AF. Very clean body, low shutter count 4,500 clicks. Includes original box, Sony battery, and strap.',
            'like_new', 'approved', 'sale', true, NOW() + INTERVAL '180 days', 1012908, 95, 14
        ) RETURNING id INTO v_pid;

        -- Images
        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail)
        VALUES 
        (v_pid, 'https://images.unsplash.com/photo-1516035069371-29a1b244cc32?auto=format&fit=crop&w=800&q=80', 0, true),
        (v_pid, 'https://images.unsplash.com/photo-1502920917128-1aa500764cbd?auto=format&fit=crop&w=800&q=80', 1, false);

        -- Attributes
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Sony' FROM attribute_definitions WHERE slug = 'brand';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Alpha 7 IV (ILCE-7M4)' FROM attribute_definitions WHERE slug = 'model';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Sony E-Mount' FROM attribute_definitions WHERE slug = 'lens-mount';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Full Frame (35mm)' FROM attribute_definitions WHERE slug = 'camera-sensor-size';
        INSERT INTO product_attributes (product_id, attribute_id, value_number)
        SELECT v_pid, id, 33 FROM attribute_definitions WHERE slug = 'camera-megapixels';
        INSERT INTO product_attributes (product_id, attribute_id, value_number)
        SELECT v_pid, id, 4500 FROM attribute_definitions WHERE slug = 'camera-shutter-count';
    END IF;

    -- ========================================================================
    -- 5. CANON EOS 5D MARK IV DSLR BODY (CATEGORY 1012915)
    -- ========================================================================
    IF NOT EXISTS (SELECT 1 FROM products WHERE title = 'Canon EOS 5D Mark IV Full Frame DSLR Camera Body') THEN
        INSERT INTO products (
            seller_id, location_id, title, price, description, condition, status,
            listing_type, fee_paid, expires_at, category_hierarchy_id, view_count, contact_count
        ) VALUES (
            v_seller_id, v_loc_bole, 'Canon EOS 5D Mark IV Full Frame DSLR Camera Body', 115000,
            '30.4MP full frame CMOS sensor, Dual Pixel CMOS AF, 4K video recording, Dual CF and SD card slots. Studio maintained, clean sensor, includes 2 genuine batteries and dual charger.',
            'good', 'approved', 'sale', true, NOW() + INTERVAL '180 days', 1012915, 54, 8
        ) RETURNING id INTO v_pid;

        -- Images
        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail)
        VALUES 
        (v_pid, 'https://images.unsplash.com/photo-1502920917128-1aa500764cbd?auto=format&fit=crop&w=800&q=80', 0, true),
        (v_pid, 'https://images.unsplash.com/photo-1516035069371-29a1b244cc32?auto=format&fit=crop&w=800&q=80', 1, false);

        -- Attributes
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Canon' FROM attribute_definitions WHERE slug = 'brand';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'EOS 5D Mark IV' FROM attribute_definitions WHERE slug = 'model';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Canon EF' FROM attribute_definitions WHERE slug = 'lens-mount';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Full Frame (35mm)' FROM attribute_definitions WHERE slug = 'camera-sensor-size';
        INSERT INTO product_attributes (product_id, attribute_id, value_number)
        SELECT v_pid, id, 30 FROM attribute_definitions WHERE slug = 'camera-megapixels';
        INSERT INTO product_attributes (product_id, attribute_id, value_number)
        SELECT v_pid, id, 18200 FROM attribute_definitions WHERE slug = 'camera-shutter-count';
    END IF;

    -- ========================================================================
    -- 6. APPLE IPHONE 15 PRO MAX 256GB (CATEGORY 9)
    -- ========================================================================
    IF NOT EXISTS (SELECT 1 FROM products WHERE title = 'Apple iPhone 15 Pro Max 256GB Natural Titanium') THEN
        INSERT INTO products (
            seller_id, location_id, title, price, description, condition, status,
            listing_type, fee_paid, expires_at, category_hierarchy_id, view_count, contact_count
        ) VALUES (
            v_seller_id, v_loc_bole, 'Apple iPhone 15 Pro Max 256GB Natural Titanium', 175000,
            'Original unlocked iPhone 15 Pro Max, 256GB, Natural Titanium. 98% battery health, flawless screen and titanium edges, includes original box and USB-C cable.',
            'like_new', 'approved', 'sale', true, NOW() + INTERVAL '180 days', 9, 142, 28
        ) RETURNING id INTO v_pid;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail)
        VALUES 
        (v_pid, 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?auto=format&fit=crop&w=800&q=80', 0, true),
        (v_pid, 'https://images.unsplash.com/photo-1510557880182-3d4d3cba35a5?auto=format&fit=crop&w=800&q=80', 1, false);

        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Apple' FROM attribute_definitions WHERE slug = 'brand';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'iPhone 15 Pro Max' FROM attribute_definitions WHERE slug = 'model';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, '256' FROM attribute_definitions WHERE slug = 'storage';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Natural Titanium' FROM attribute_definitions WHERE slug = 'color';
    END IF;

    -- ========================================================================
    -- 7. SAMSUNG GALAXY S24 ULTRA 512GB (CATEGORY 9)
    -- ========================================================================
    IF NOT EXISTS (SELECT 1 FROM products WHERE title = 'Samsung Galaxy S24 Ultra 512GB Titanium Black') THEN
        INSERT INTO products (
            seller_id, location_id, title, price, description, condition, status,
            listing_type, fee_paid, expires_at, category_hierarchy_id, view_count, contact_count
        ) VALUES (
            v_seller_id, v_loc_bole, 'Samsung Galaxy S24 Ultra 512GB Titanium Black', 168000,
            'Brand new sealed Samsung Galaxy S24 Ultra 5G with integrated S-Pen, Galaxy AI features, Snapdragon 8 Gen 3, 200MP camera system, 512GB storage.',
            'new', 'approved', 'sale', true, NOW() + INTERVAL '180 days', 9, 88, 16
        ) RETURNING id INTO v_pid;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail)
        VALUES 
        (v_pid, 'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?auto=format&fit=crop&w=800&q=80', 0, true),
        (v_pid, 'https://images.unsplash.com/photo-1580910051074-3eb694886505?auto=format&fit=crop&w=800&q=80', 1, false);

        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Samsung' FROM attribute_definitions WHERE slug = 'brand';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Galaxy S24 Ultra' FROM attribute_definitions WHERE slug = 'model';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, '512' FROM attribute_definitions WHERE slug = 'storage';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Titanium Black' FROM attribute_definitions WHERE slug = 'color';
    END IF;

    -- ========================================================================
    -- 8. APPLE MACBOOK PRO 14 M3 PRO (CATEGORY 10)
    -- ========================================================================
    IF NOT EXISTS (SELECT 1 FROM products WHERE title = 'Apple MacBook Pro 14-inch M3 Pro 18GB RAM 512GB SSD') THEN
        INSERT INTO products (
            seller_id, location_id, title, price, description, condition, status,
            listing_type, fee_paid, expires_at, category_hierarchy_id, view_count, contact_count
        ) VALUES (
            v_seller_id, v_loc_bole, 'Apple MacBook Pro 14-inch M3 Pro 18GB RAM 512GB SSD', 225000,
            'Space Black MacBook Pro 14 with M3 Pro chip (11-core CPU, 14-core GPU), 18GB unified memory, 512GB SSD, Liquid Retina XDR display with 120Hz ProMotion. Pristine battery 100%.',
            'like_new', 'approved', 'sale', true, NOW() + INTERVAL '180 days', 10, 115, 19
        ) RETURNING id INTO v_pid;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail)
        VALUES 
        (v_pid, 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?auto=format&fit=crop&w=800&q=80', 0, true),
        (v_pid, 'https://images.unsplash.com/photo-1611186871348-b1ce696e52c9?auto=format&fit=crop&w=800&q=80', 1, false);

        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Apple' FROM attribute_definitions WHERE slug = 'brand';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'MacBook Pro 14' FROM attribute_definitions WHERE slug = 'model';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Apple M3 Pro' FROM attribute_definitions WHERE slug = 'processor';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, '18' FROM attribute_definitions WHERE slug = 'ram';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, '512' FROM attribute_definitions WHERE slug = 'storage';
    END IF;

    -- ========================================================================
    -- 9. TOYOTA VITZ 2018 AUTOMATIC (CATEGORY 14)
    -- ========================================================================
    IF NOT EXISTS (SELECT 1 FROM products WHERE title = 'Toyota Vitz 2018 Clean Condition Automatic Transmission') THEN
        INSERT INTO products (
            seller_id, location_id, title, price, description, condition, status,
            listing_type, fee_paid, expires_at, category_hierarchy_id, view_count, contact_count
        ) VALUES (
            v_seller_id, v_loc_bole, 'Toyota Vitz 2018 Clean Condition Automatic Transmission', 1850000,
            'Very clean Toyota Vitz 2018 model, push start, reverse camera, clean interior, well-maintained city car in Addis Ababa. All service records available.',
            'good', 'approved', 'sale', true, NOW() + INTERVAL '180 days', 14, 210, 34
        ) RETURNING id INTO v_pid;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail)
        VALUES 
        (v_pid, 'https://images.unsplash.com/photo-1549399542-7e3f8b79c341?auto=format&fit=crop&w=800&q=80', 0, true),
        (v_pid, 'https://images.unsplash.com/photo-1552519507-da3b142c6e3d?auto=format&fit=crop&w=800&q=80', 1, false);

        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Toyota' FROM attribute_definitions WHERE slug = 'make';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Vitz' FROM attribute_definitions WHERE slug = 'vehicle-model';
        INSERT INTO product_attributes (product_id, attribute_id, value_number)
        SELECT v_pid, id, 2018 FROM attribute_definitions WHERE slug = 'year';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Automatic' FROM attribute_definitions WHERE slug = 'transmission';
        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Petrol' FROM attribute_definitions WHERE slug = 'fuel-type';
    END IF;

    -- ========================================================================
    -- 10. BOLE ATLAS 2-BEDROOM FURNISHED APARTMENT FOR RENT (CATEGORY 20)
    -- ========================================================================
    IF NOT EXISTS (SELECT 1 FROM products WHERE title = 'Modern 2-Bedroom Fully Furnished Apartment in Bole Atlas') THEN
        INSERT INTO products (
            seller_id, location_id, title, price, rental_price, rental_duration_value, rental_duration_unit,
            rental_price_per_month, description, condition, status,
            listing_type, fee_paid, expires_at, category_hierarchy_id, view_count, contact_count
        ) VALUES (
            v_seller_id, v_loc_bole, 'Modern 2-Bedroom Fully Furnished Apartment in Bole Atlas', 55000, 55000, 1, 'month'::rental_duration_unit,
            55000,
            'Spacious 2-bedroom furnished apartment with panoramic balcony view, modern fitted kitchen, 24/7 security, backup generator, dedicated underground parking.',
            'like_new', 'approved', 'rent', true, NOW() + INTERVAL '180 days', 20, 160, 22
        ) RETURNING id INTO v_pid;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail)
        VALUES 
        (v_pid, 'https://images.unsplash.com/photo-1522708323590-d24dbb6b0267?auto=format&fit=crop&w=800&q=80', 0, true),
        (v_pid, 'https://images.unsplash.com/photo-1502672260266-1c1ef2d93688?auto=format&fit=crop&w=800&q=80', 1, false);

        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'Apartment' FROM attribute_definitions WHERE slug = 'property-type';
        INSERT INTO product_attributes (product_id, attribute_id, value_number)
        SELECT v_pid, id, 2 FROM attribute_definitions WHERE slug = 'bedrooms';
        INSERT INTO product_attributes (product_id, attribute_id, value_number)
        SELECT v_pid, id, 2 FROM attribute_definitions WHERE slug = 'bathrooms';
        INSERT INTO product_attributes (product_id, attribute_id, value_boolean)
        SELECT v_pid, id, true FROM attribute_definitions WHERE slug = 'furnished';
        INSERT INTO product_attributes (product_id, attribute_id, value_boolean)
        SELECT v_pid, id, true FROM attribute_definitions WHERE slug = 'parking';
    END IF;

    -- ========================================================================
    -- 11. BOSCH GBH 2-26 DRE ROTARY HAMMER DRILL (CATEGORY 625225)
    -- ========================================================================
    IF NOT EXISTS (SELECT 1 FROM products WHERE title = 'Bosch GBH 2-26 DRE Professional Rotary Hammer Drill SDS-Plus 800W') THEN
        INSERT INTO products (
            seller_id, location_id, title, price, description, condition, status,
            listing_type, fee_paid, expires_at, category_hierarchy_id, view_count, contact_count
        ) VALUES (
            v_seller_id, v_loc_bole, 'Bosch GBH 2-26 DRE Professional Rotary Hammer Drill SDS-Plus 800W', 16500,
            'Genuine Bosch GBH 2-26 DRE heavy duty SDS-Plus rotary hammer drill. 800W motor, 2.7 Joules impact energy, forward/reverse rotation, carrying case included.',
            'new', 'approved', 'sale', true, NOW() + INTERVAL '180 days', 625225, 41, 6
        ) RETURNING id INTO v_pid;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail)
        VALUES 
        (v_pid, 'https://images.unsplash.com/photo-1504148455328-c376907d081c?auto=format&fit=crop&w=800&q=80', 0, true),
        (v_pid, 'https://images.unsplash.com/photo-1572981779307-38b8cabb2407?auto=format&fit=crop&w=800&q=80', 1, false);

        INSERT INTO product_attributes (product_id, attribute_id, value_text)
        SELECT v_pid, id, 'SDS-Plus' FROM attribute_definitions WHERE slug = 'attr-chuck-type';
        INSERT INTO product_attributes (product_id, attribute_id, value_number)
        SELECT v_pid, id, 2.7 FROM attribute_definitions WHERE slug = 'attr-impact-energy';
    END IF;

END $$;

COMMIT;
