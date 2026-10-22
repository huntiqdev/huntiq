-- ============================================================
-- 005_products.sql
-- PURPOSE : Product listings + product_images child table.
--           product_images has no meaning without its parent product.
--
-- ACCESS TIER LOGIC:
--   access_tier NULL    → listing visible to all buyers
--   access_tier 'gold'  → buyer needs gold access pass to see
--                         seller phone, location, all images
--   access_tier 'premium'→ buyer needs premium access pass to see
--                         seller phone, location, unblurred images
--
-- IMAGE LIMITS BY SELLER TIER (from platform_settings):
--   normal  → max 2 images
--   gold    → max 7 images  (all visible to gold+ buyers)
--   premium → max 20 images (seller marks individual images as blurred)
-- ============================================================
ALTER TABLE products
    -- ── Ownership & Classification ────────────────────────────
    -- seller_id    : FK to the seller who listed this product
    ADD COLUMN IF NOT EXISTS seller_id                INT               NOT NULL REFERENCES users(id)      ON DELETE CASCADE,
    -- location_id  : FK to Ethiopian location (region/zone/subcity/woreda)
    ADD COLUMN IF NOT EXISTS location_id              INT               REFERENCES locations(id)            ON DELETE SET NULL,
    -- title        : listing headline shown on search cards
    ADD COLUMN IF NOT EXISTS title                    VARCHAR(200)      NOT NULL DEFAULT '',
    -- price        : asking price in ETB (must be > 0)
    -- price: > 0 and capped at 999,999,999.99 ETB — prevents absurd price DoS
    ADD COLUMN IF NOT EXISTS price                    NUMERIC(12,2)     NOT NULL DEFAULT 0 CHECK (price > 0 AND price <= 999999999.99),
    -- description  : full listing description (shown on detail page)
    ADD COLUMN IF NOT EXISTS description              TEXT,
    -- category_id  : FK to top-level category
    ADD COLUMN IF NOT EXISTS category_id              INT               REFERENCES categories(id)           ON DELETE SET NULL,
    -- subcategory  : free-text subcategory label within the category
    ADD COLUMN IF NOT EXISTS subcategory              VARCHAR(100)      DEFAULT '',
    -- item         : specific item name within subcategory
    ADD COLUMN IF NOT EXISTS item                     VARCHAR(100)      DEFAULT '',
    -- condition    : physical condition of the item
    ADD COLUMN IF NOT EXISTS condition                product_condition NOT NULL DEFAULT 'good',

    -- ── Status & Review ───────────────────────────────────────
    -- status       : listing lifecycle — pending→approved→sold/expired/rejected
    ADD COLUMN IF NOT EXISTS status                   product_status    NOT NULL DEFAULT 'pending',
    -- admin_notes  : admin review notes (rejection reason, fee payment notes)
    ADD COLUMN IF NOT EXISTS admin_notes              TEXT,
    -- reviewed_by  : which admin approved or rejected this listing
    ADD COLUMN IF NOT EXISTS reviewed_by              INT               REFERENCES users(id)                ON DELETE SET NULL,
    -- reviewed_at  : when the admin reviewed this listing
    ADD COLUMN IF NOT EXISTS reviewed_at              TIMESTAMPTZ,

    -- ── Analytics ─────────────────────────────────────────────
    -- view_count   : incremented each time listing detail page is viewed
    ADD COLUMN IF NOT EXISTS view_count               INT               NOT NULL DEFAULT 0,
    -- contact_count: incremented each time a buyer submits an order_request
    ADD COLUMN IF NOT EXISTS contact_count            INT               NOT NULL DEFAULT 0,

    -- ── Lifecycle Timestamps ──────────────────────────────────
    -- expires_at   : when the listing expires (set on approval based on duration)
    ADD COLUMN IF NOT EXISTS expires_at               TIMESTAMPTZ       DEFAULT NULL,
    -- sold_at      : auto-set by trigger when status changes to 'sold'
    ADD COLUMN IF NOT EXISTS sold_at                  TIMESTAMPTZ       DEFAULT NULL,

    -- ── Precise Location ──────────────────────────────────────
    -- kebele       : smallest sub-woreda unit (optional)
    ADD COLUMN IF NOT EXISTS kebele                   VARCHAR(100)      DEFAULT NULL,
    -- specific_area: free-text landmark or area description
    ADD COLUMN IF NOT EXISTS specific_area            TEXT              DEFAULT NULL,
    -- latitude/longitude: GPS coordinates for map display
    ADD COLUMN IF NOT EXISTS latitude                 NUMERIC(10,7)     DEFAULT NULL CHECK (latitude  IS NULL OR (latitude  >= -90  AND latitude  <= 90)),
    ADD COLUMN IF NOT EXISTS longitude                NUMERIC(10,7)     DEFAULT NULL CHECK (longitude IS NULL OR (longitude >= -180 AND longitude <= 180)),

    -- ── Listing Fee ───────────────────────────────────────────
    -- listing_duration_days   : how many days the listing runs
    ADD COLUMN IF NOT EXISTS listing_duration_days    INT               DEFAULT NULL,
    -- startup_fee_percent     : initial fee % charged at listing approval
    ADD COLUMN IF NOT EXISTS startup_fee_percent      NUMERIC(5,2)      DEFAULT NULL,
    -- weekly_increment_percent: additional % added each week after free weeks
    ADD COLUMN IF NOT EXISTS weekly_increment_percent NUMERIC(5,2)      DEFAULT NULL,
    -- total_fee_percent       : final computed fee % (startup + increments)
    ADD COLUMN IF NOT EXISTS total_fee_percent        NUMERIC(5,2)      DEFAULT NULL,
    -- total_fee_amount        : final fee in ETB = price × total_fee_percent / 100
    ADD COLUMN IF NOT EXISTS total_fee_amount         NUMERIC(12,2)     DEFAULT NULL,
    -- fee_paid                : TRUE = admin confirmed payment received
    ADD COLUMN IF NOT EXISTS fee_paid                 BOOLEAN           NOT NULL DEFAULT FALSE,
    -- fee_payment_pending     : TRUE = seller submitted payment claim, awaiting admin verification
    --   Flow: seller submits → fee_payment_pending=TRUE
    --         admin verifies → fee_paid=TRUE, fee_payment_pending=FALSE
    ADD COLUMN IF NOT EXISTS fee_payment_pending      BOOLEAN           NOT NULL DEFAULT FALSE,
    -- fee_paid_at             : when admin verified the fee payment (for revenue attribution)
    ADD COLUMN IF NOT EXISTS fee_paid_at              TIMESTAMPTZ       DEFAULT NULL,

    -- ── Tier Snapshot (denormalized) ──────────────────────────
    -- seller_tier_snapshot: seller's tier at listing time — kept in sync by DB trigger
    --   Used for filtering/display without joining users table
    -- seller_tier_snapshot: CHECK mirrors seller_tier enum — prevents arbitrary string injection
    ADD COLUMN IF NOT EXISTS seller_tier_snapshot     VARCHAR(10)       NOT NULL DEFAULT 'normal' CHECK (seller_tier_snapshot IN ('normal','gold','premium')),
    -- access_tier : what buyer access pass is needed to see full listing details
    --   NULL      = open listing, no pass needed
    --   'gold'    = buyer needs gold access pass
    --   'premium' = buyer needs premium access pass
    --   Set automatically from seller_tier_snapshot by trigger
    ADD COLUMN IF NOT EXISTS access_tier              VARCHAR(10)       DEFAULT NULL CHECK (access_tier IS NULL OR access_tier IN ('gold','premium')),

    -- ── Full-text Search ──────────────────────────────────────
    -- search_vector: auto-generated tsvector from title + description
    --   Used by GIN index for fast full-text search queries
    ADD COLUMN IF NOT EXISTS created_at               TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at               TIMESTAMPTZ       NOT NULL DEFAULT NOW();

-- Generated column — must be separate ALTER (PostgreSQL requirement)
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS search_vector tsvector GENERATED ALWAYS AS (
        to_tsvector('simple', coalesce(title,'') || ' ' || coalesce(description,''))
    ) STORED;

-- ── Indexes ──────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_products_seller_id           ON products (seller_id);
CREATE INDEX IF NOT EXISTS idx_products_status              ON products (status);
CREATE INDEX IF NOT EXISTS idx_products_category_id         ON products (category_id);
CREATE INDEX IF NOT EXISTS idx_products_location_id         ON products (location_id);
CREATE INDEX IF NOT EXISTS idx_products_subcategory         ON products (subcategory);
CREATE INDEX IF NOT EXISTS idx_products_item                ON products (item);
CREATE INDEX IF NOT EXISTS idx_products_expires_at          ON products (expires_at)          WHERE expires_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_products_fee_paid            ON products (fee_paid)             WHERE fee_paid = FALSE;
CREATE INDEX IF NOT EXISTS idx_products_fee_payment_pending ON products (fee_payment_pending)  WHERE fee_payment_pending = TRUE;
CREATE INDEX IF NOT EXISTS idx_products_access_tier         ON products (access_tier)          WHERE access_tier IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_products_search              ON products USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS idx_products_coords              ON products (latitude, longitude)  WHERE latitude IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_products_reviewed_by         ON products (reviewed_by)          WHERE reviewed_by IS NOT NULL;

-- ============================================================
-- product_images (child of products — no meaning without parent)
-- PURPOSE : Ordered images for a product listing.
--
-- THUMBNAIL LOGIC:
--   is_thumbnail = TRUE  → shown on search card / listing preview
--                          exactly ONE per product, always visible
--                          NEVER blurred regardless of seller tier
--   is_thumbnail = FALSE → shown only on detail page
--
-- BLUR LOGIC (premium sellers only):
--   is_blurred = FALSE → visible to all buyers
--   is_blurred = TRUE  → blurred for buyers without premium access pass
--                        premium buyers see full image
--                        thumbnail is NEVER blurred (enforced in Go)
--
-- DELETE/UPDATE (premium sellers only):
--   deleted_at NULL → active image, included in all queries
--   deleted_at SET  → soft-deleted, excluded from all queries
--   updated_at      → tracks when image URL was replaced
-- ============================================================
ALTER TABLE product_images
    -- product_id   : FK to parent product — cascade delete removes all images
    ADD COLUMN IF NOT EXISTS product_id   INT         NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    -- image_url    : relative path or CDN URL to the image file
    ADD COLUMN IF NOT EXISTS image_url    TEXT        NOT NULL DEFAULT '',
    -- order_number : display order (0 = first shown, thumbnail is typically order 0)
    ADD COLUMN IF NOT EXISTS order_number SMALLINT    NOT NULL DEFAULT 0,
    -- is_thumbnail : TRUE = this image is the search card thumbnail
    --               Exactly one TRUE per product — enforced by partial unique index
    ADD COLUMN IF NOT EXISTS is_thumbnail BOOLEAN     NOT NULL DEFAULT FALSE,
    -- is_blurred   : TRUE = premium seller marked this image as gated
    --               Non-premium buyers see blur overlay
    --               Only meaningful when seller tier = 'premium'
    ADD COLUMN IF NOT EXISTS is_blurred   BOOLEAN     NOT NULL DEFAULT FALSE,
    -- deleted_at   : soft delete timestamp — NULL = active, SET = deleted
    ADD COLUMN IF NOT EXISTS deleted_at   TIMESTAMPTZ DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- updated_at   : set when image URL is replaced (premium seller update)
    ADD COLUMN IF NOT EXISTS updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Each product has exactly one thumbnail (partial unique index — only one TRUE per product)
CREATE UNIQUE INDEX IF NOT EXISTS uq_product_thumbnail
    ON product_images (product_id) WHERE is_thumbnail = TRUE;
-- Each order_number is unique per product
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name='product_images' AND constraint_name='uq_product_image_order'
    ) THEN
        ALTER TABLE product_images ADD CONSTRAINT uq_product_image_order UNIQUE (product_id, order_number);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images (product_id)    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_product_images_thumbnail  ON product_images (product_id)    WHERE is_thumbnail = TRUE;
CREATE INDEX IF NOT EXISTS idx_product_images_blurred    ON product_images (product_id)    WHERE is_blurred = TRUE AND deleted_at IS NULL;
