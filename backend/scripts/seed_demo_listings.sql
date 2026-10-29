-- HUNTIQ Realistic Marketplace Demo Listings Seeder (Idempotent)
-- Author: HUNTIQ Engineering
-- Purpose: Seeds realistic verified Ethiopian marketplace listings across vehicles, electronics,
--          real estate rentals, laptops, and fashion with full translations, locations, and images.

BEGIN;

-- 1. Ensure existing products have proper location and category hierarchy IDs
UPDATE products 
SET location_id = 11, category_hierarchy_id = 7112, subcategory = 'Apartments', item = '2-Bedroom'
WHERE id = 3;

UPDATE products 
SET location_id = 1, category_hierarchy_id = 12, subcategory = 'Mobile Phones', item = 'Smartphone'
WHERE id = 15;

-- Add thumbnail for product 3 if missing
INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred, created_at, updated_at)
SELECT 3, 'https://images.unsplash.com/photo-1522708323590-d24dbb6b0267?w=800&auto=format&fit=crop&q=60', 1, true, false, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM product_images WHERE product_id = 3);

-- 2. Insert diverse verified listings (if not already seeded)
DO $$
DECLARE
    v_seller_id INT;
    v_prod_id INT;
BEGIN
    SELECT id INTO v_seller_id FROM users WHERE role = 'seller' ORDER BY id LIMIT 1;
    IF v_seller_id IS NULL THEN
        RAISE EXCEPTION 'No seller user found in users table';
    END IF;

    -- Listing 1: Toyota Vitz 2018 (Vehicles / Bole)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%Toyota Vitz 2018%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 1, 2, 74,
            'Toyota Vitz 2018 Automatic — Libre Verified', 1850000.00,
            'Clean Toyota Vitz 2018 model car, automatic transmission, original factory paint, 48,000 km, clean Libre title deed verified, single owner in Bole Addis Ababa. car automobile vehicle', 'good', 'approved',
            'sale', 0, 'gold', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1549399542-7e3f8b79c341?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;

    -- Listing 2: Apple iPhone 15 Pro Max 256GB (Electronics / Bole)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%iPhone 15 Pro Max 256GB%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 1, 1, 12,
            'Apple iPhone 15 Pro Max 256GB Natural Titanium', 145000.00,
            'Brand new sealed in box Apple iPhone 15 Pro Max 256GB Natural Titanium smartphone, factory unlocked, clean IMEI verified with Telecommunications Authority, full 1-year warranty. phone smartphone mobile', 'new', 'approved',
            'sale', 0, 'premium', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;

    -- Listing 3: Apple MacBook Pro 14" M3 Pro 18GB/512GB (Computers / Kirkos)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%MacBook Pro 14" M3 Pro%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 2, 1, 14,
            'Apple MacBook Pro 14" M3 Pro 18GB RAM 512GB SSD Space Black', 185000.00,
            'Latest Apple M3 Pro chip with 11-core CPU and 14-core GPU, 18GB unified memory, 512GB SSD storage, Liquid Retina XDR display, cycle count 12, immaculate condition with original 70W charger. laptop notebook computer', 'good', 'approved',
            'sale', 0, 'premium', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;

    -- Listing 4: Luxury 3-Bedroom Furnished Villa (Real Estate Rental / Bole)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%Luxury 3-Bedroom Furnished Villa%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, rental_price, rental_duration_value, rental_duration_unit,
            deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 1, 3, 7112,
            'Luxury 3-Bedroom Furnished Villa in Bole Atlas', 85000.00,
            'Expansive 3-bedroom, 3.5-bath luxury villa near Atlas Hotel, private garden, modern European kitchen, 24/7 dedicated generator, high-capacity water tanks, master ensuite with walk-in closet.', 'new', 'approved',
            'rent', 85000.00, 1, 'month',
            85000.00, 'premium', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1600596542815-ffad4c1539a9?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;

    -- Listing 5: Studio Apartment Furnished (Real Estate Rental / Yeka)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%Studio Apartment Furnished Yeka%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, rental_price, rental_duration_value, rental_duration_unit,
            deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 3, 3, 7112,
            'Cozy Studio Apartment Furnished Yeka Signal', 22000.00,
            'Fully furnished sunny studio apartment in Yeka Signal, high-speed WiFi ready, smart TV, hot shower, backup water supply, secure building with CCTV and parking space included.', 'good', 'approved',
            'rent', 22000.00, 1, 'month',
            22000.00, 'normal', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1502672260266-1c1ef2d93688?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;

    -- Listing 6: Commercial Retail Shop in Merkato (Commercial / Addis Ketema)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%Commercial Retail Shop in Merkato%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, rental_price, rental_duration_value, rental_duration_unit,
            deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 6, 3, 7112,
            'Prime Commercial Retail Shop in Merkato Tana Mall', 45000.00,
            'High foot-traffic 32 sqm commercial shop on ground floor near main entrance, ideal for electronics, clothing boutique, or cosmetics distribution, with 3-phase electricity.', 'good', 'approved',
            'rent', 45000.00, 1, 'month',
            45000.00, 'gold', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1555529669-e69e7aa0ba9a?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;

    -- Listing 7: Dell XPS 15 9530 Core i9 32GB (Computers / Bole)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%Dell XPS 15 9530%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 1, 1, 14,
            'Dell XPS 15 9530 Intel Core i9 32GB RAM RTX 4070', 195000.00,
            'Workstation powerhouse: 13th Gen Intel Core i9-13900H, 32GB DDR5 RAM, 1TB NVMe Gen4 SSD, NVIDIA GeForce RTX 4070 8GB, 3.5K OLED touchscreen display, 86Wh battery with 130W USB-C charger.', 'new', 'approved',
            'sale', 0, 'gold', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1593642632823-8f785ba67e45?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;

    -- Listing 8: Authentic Handwoven Habesha Kemis (Fashion / Kirkos)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%Handwoven Habesha Kemis%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 2, 5, 47,
            'Traditional Handwoven Habesha Kemis with Pure Silk Tilet (የሀበሻ ቀሚስ)', 18500.00,
            'Exquisite 100% pure Ethiopian cotton (Shemma) dress handcrafted by master weavers in Shiro Meda, featuring intricate golden and burgundy silk border embroidery (Tilet), includes matching Netela shawl.', 'new', 'approved',
            'sale', 0, 'gold', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1617627143750-d86bc21e42bb?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;

    -- Listing 9: Sony Alpha A7 IV Mirrorless Camera Body (Electronics / Arada)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%Sony Alpha A7 IV%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 4, 1, 41,
            'Sony Alpha A7 IV Full-Frame Mirrorless Camera Body', 160000.00,
            '33MP full-frame Exmor R CMOS sensor, 4K 60p 10-bit 4:2:2 video, 759-point fast hybrid AF with real-time Eye AF for humans, animals, and birds, shutter count only 3,200, includes 2 original Sony batteries and dual charger.', 'good', 'approved',
            'sale', 0, 'gold', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1516035069371-29a1b244cc32?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;

    -- Listing 10: Agricultural 4WD Farm Tractor 75HP (Agriculture / Oromia)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%Agricultural 4WD Farm Tractor%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 12, 6, 80,
            'Tafe 75HP 4WD Agricultural Farm Tractor with Disc Plough', 2150000.00,
            'Heavy-duty 75 horsepower 4-cylinder Simpson diesel engine, 4-wheel drive, dual clutch with power steering, heavy-duty rear hydraulic lift, accompanied by 3-disc reversible plough, inspected in Adama Oromia.', 'good', 'approved',
            'sale', 0, 'premium', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1592878904946-b3cd8ae243d0?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;

    -- Listing 11: Silent Diesel Industrial Generator 25kVA (Business / Lideta)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%Silent Diesel Industrial Generator%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 5, 7, 93,
            'Perkins 25kVA Soundproof Silent Diesel Generator with ATS', 680000.00,
            'Three-phase 380V/220V 25kVA standby silent acoustic canopy generator, genuine Perkins engine with Stamford alternator, Deep Sea electronic control panel and automatic transfer switch (ATS) included.', 'new', 'approved',
            'sale', 0, 'premium', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1581092160607-ee22621dd758?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;

    -- Listing 12: Toyota Hilux Double Cab 4x4 (Vehicles / Bole)
    IF NOT EXISTS (SELECT 1 FROM products WHERE title ILIKE '%Toyota Hilux Double Cab 4x4%') THEN
        INSERT INTO products (
            seller_id, location_id, category_id, category_hierarchy_id,
            title, price, description, condition, status,
            listing_type, deposit_amount, seller_tier_snapshot, fee_paid
        ) VALUES (
            v_seller_id, 1, 2, 74,
            'Toyota Hilux 2.8 D-4D Double Cab 4x4 Automatic 2021', 4950000.00,
            '2021 Toyota Hilux Invincible Double Cab, 2.8L turbo diesel engine, 4WD selectable drive, leather seats, 360-degree camera, canopy included, full maintenance history, Libre deed in hand in Bole. car automobile vehicle truck', 'good', 'approved',
            'sale', 0, 'premium', true
        ) RETURNING id INTO v_prod_id;

        INSERT INTO product_images (product_id, image_url, order_number, is_thumbnail, is_blurred)
        VALUES (v_prod_id, 'https://images.unsplash.com/photo-1559416523-140ddc3d238c?w=800&auto=format&fit=crop&q=60', 1, true, false);
    END IF;
END $$;

COMMIT;
