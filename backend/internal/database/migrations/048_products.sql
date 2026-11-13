-- ============================================================================
-- 048_products.sql
-- Table: products (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for products
ALTER TABLE products ADD COLUMN IF NOT EXISTS seller_id integer NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS location_id integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS title character varying(200) DEFAULT ''::character varying NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS price numeric(12,2) DEFAULT 0 NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS description text;
ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS subcategory character varying(100) DEFAULT ''::character varying;
ALTER TABLE products ADD COLUMN IF NOT EXISTS item character varying(100) DEFAULT ''::character varying;
ALTER TABLE products ADD COLUMN IF NOT EXISTS condition product_condition DEFAULT 'good'::product_condition NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS status product_status DEFAULT 'pending'::product_status NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS admin_notes text;
ALTER TABLE products ADD COLUMN IF NOT EXISTS reviewed_by integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS reviewed_at timestamp with time zone;
ALTER TABLE products ADD COLUMN IF NOT EXISTS view_count integer DEFAULT 0 NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS contact_count integer DEFAULT 0 NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS expires_at timestamp with time zone;
ALTER TABLE products ADD COLUMN IF NOT EXISTS sold_at timestamp with time zone;
ALTER TABLE products ADD COLUMN IF NOT EXISTS kebele character varying(100) DEFAULT NULL::character varying;
ALTER TABLE products ADD COLUMN IF NOT EXISTS specific_area text;
ALTER TABLE products ADD COLUMN IF NOT EXISTS latitude numeric(10,7) DEFAULT NULL::numeric;
ALTER TABLE products ADD COLUMN IF NOT EXISTS longitude numeric(10,7) DEFAULT NULL::numeric;
ALTER TABLE products ADD COLUMN IF NOT EXISTS listing_duration_days integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS startup_fee_percent numeric(5,2) DEFAULT NULL::numeric;
ALTER TABLE products ADD COLUMN IF NOT EXISTS weekly_increment_percent numeric(5,2) DEFAULT NULL::numeric;
ALTER TABLE products ADD COLUMN IF NOT EXISTS total_fee_percent numeric(5,2) DEFAULT NULL::numeric;
ALTER TABLE products ADD COLUMN IF NOT EXISTS total_fee_amount numeric(12,2) DEFAULT NULL::numeric;
ALTER TABLE products ADD COLUMN IF NOT EXISTS fee_paid boolean DEFAULT false NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS fee_payment_pending boolean DEFAULT false NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS fee_paid_at timestamp with time zone;
ALTER TABLE products ADD COLUMN IF NOT EXISTS seller_tier_snapshot character varying(10) DEFAULT 'normal'::character varying NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS access_tier character varying(10) DEFAULT NULL::character varying;
ALTER TABLE products ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS search_vector tsvector GENERATED ALWAYS AS (to_tsvector('simple'::regconfig, (((COALESCE(title, ''::character varying))::text || ' '::text) || COALESCE(description, ''::text)))) STORED;
ALTER TABLE products ADD COLUMN IF NOT EXISTS listing_type listing_type DEFAULT 'sale'::listing_type NOT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS rental_price numeric(12,2);
ALTER TABLE products ADD COLUMN IF NOT EXISTS rental_duration_value integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS rental_duration_unit rental_duration_unit;
ALTER TABLE products ADD COLUMN IF NOT EXISTS deposit_amount numeric(12,2) DEFAULT 0;
ALTER TABLE products ADD COLUMN IF NOT EXISTS available_from date;
ALTER TABLE products ADD COLUMN IF NOT EXISTS available_until date;
ALTER TABLE products ADD COLUMN IF NOT EXISTS is_currently_rented boolean DEFAULT false;
ALTER TABLE products ADD COLUMN IF NOT EXISTS current_tenant_id integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS rented_until date;
ALTER TABLE products ADD COLUMN IF NOT EXISTS rental_price_per_month numeric(12,2);
ALTER TABLE products ADD COLUMN IF NOT EXISTS rental_min_duration_months integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS rental_max_duration_months integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS rental_deposit_amount numeric(12,2);
ALTER TABLE products ADD COLUMN IF NOT EXISTS rental_terms text;
ALTER TABLE products ADD COLUMN IF NOT EXISTS available_from_ts timestamp with time zone;
ALTER TABLE products ADD COLUMN IF NOT EXISTS category_hierarchy_id integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS location_hierarchy_id integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS brand_id integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS product_model_id integer;
ALTER TABLE products ADD COLUMN IF NOT EXISTS custom_model_name character varying(200);

-- Table Constraints for products
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'check_rental_fields'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT check_rental_fields CHECK ((((listing_type = 'sale'::listing_type) AND (rental_price IS NULL) AND (rental_duration_value IS NULL)) OR ((listing_type = 'rent'::listing_type) AND (rental_price IS NOT NULL) AND (rental_duration_value IS NOT NULL))));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_access_tier_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_access_tier_check CHECK (((access_tier IS NULL) OR ((access_tier)::text = ANY ((ARRAY['gold'::character varying, 'premium'::character varying])::text[]))));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_deposit_amount_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_deposit_amount_check CHECK ((deposit_amount >= (0)::numeric));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_latitude_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_latitude_check CHECK (((latitude IS NULL) OR ((latitude >= ('-90'::integer)::numeric) AND (latitude <= (90)::numeric))));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_longitude_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_longitude_check CHECK (((longitude IS NULL) OR ((longitude >= ('-180'::integer)::numeric) AND (longitude <= (180)::numeric))));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_price_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_price_check CHECK (((price > (0)::numeric) AND (price <= 999999999.99)));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_rental_deposit_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_rental_deposit_check CHECK (((rental_deposit_amount IS NULL) OR (rental_deposit_amount >= (0)::numeric)));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_rental_duration_value_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_rental_duration_value_check CHECK (((rental_duration_value IS NULL) OR (rental_duration_value > 0)));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_rental_max_duration_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_rental_max_duration_check CHECK (((rental_max_duration_months IS NULL) OR (rental_max_duration_months >= rental_min_duration_months)));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_rental_min_duration_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_rental_min_duration_check CHECK (((rental_min_duration_months IS NULL) OR (rental_min_duration_months >= 1)));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_rental_price_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_rental_price_check CHECK (((rental_price IS NULL) OR ((rental_price > (0)::numeric) AND (rental_price <= 999999999.99))));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_seller_tier_snapshot_check'
    ) THEN
        ALTER TABLE products ADD CONSTRAINT products_seller_tier_snapshot_check CHECK (((seller_tier_snapshot)::text = ANY ((ARRAY['normal'::character varying, 'gold'::character varying, 'premium'::character varying])::text[])));
    END IF;
END $$;

COMMENT ON COLUMN products.listing_type IS 'Type of listing: sale (default) or rent';

COMMENT ON COLUMN products.rental_price IS 'Rental price per duration unit (only for rent listings)';

COMMENT ON COLUMN products.is_currently_rented IS 'Whether rental item is currently rented out';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_pkey'
    ) THEN
        ALTER TABLE ONLY products ADD CONSTRAINT products_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_products_access_tier ON products USING btree (access_tier) WHERE (access_tier IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_products_available_from ON products USING btree (available_from) WHERE (listing_type = 'rent'::listing_type);

CREATE INDEX IF NOT EXISTS idx_products_brand ON products USING btree (brand_id);

CREATE INDEX IF NOT EXISTS idx_products_category_hierarchy ON products USING btree (category_hierarchy_id);

CREATE INDEX IF NOT EXISTS idx_products_category_id ON products USING btree (category_id);

CREATE INDEX IF NOT EXISTS idx_products_category_status ON products USING btree (category_id, status);

CREATE INDEX IF NOT EXISTS idx_products_coords ON products USING btree (latitude, longitude) WHERE (latitude IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_products_description_trgm ON products USING gin (description gin_trgm_ops) WHERE (status = 'approved'::product_status);

CREATE INDEX IF NOT EXISTS idx_products_expires_at ON products USING btree (expires_at) WHERE (expires_at IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_products_fee_paid ON products USING btree (fee_paid) WHERE (fee_paid = false);

CREATE INDEX IF NOT EXISTS idx_products_fee_payment_pending ON products USING btree (fee_payment_pending) WHERE (fee_payment_pending = true);

CREATE INDEX IF NOT EXISTS idx_products_is_rented ON products USING btree (is_currently_rented) WHERE ((listing_type = 'rent'::listing_type) AND (is_currently_rented = false));

CREATE INDEX IF NOT EXISTS idx_products_item ON products USING btree (item);

CREATE INDEX IF NOT EXISTS idx_products_listing_type ON products USING btree (listing_type);

CREATE INDEX IF NOT EXISTS idx_products_location_hierarchy ON products USING btree (location_hierarchy_id);

CREATE INDEX IF NOT EXISTS idx_products_location_id ON products USING btree (location_id);

CREATE INDEX IF NOT EXISTS idx_products_model ON products USING btree (product_model_id);

CREATE INDEX IF NOT EXISTS idx_products_rental_price ON products USING btree (rental_price) WHERE (listing_type = 'rent'::listing_type);

CREATE INDEX IF NOT EXISTS idx_products_reviewed_by ON products USING btree (reviewed_by) WHERE (reviewed_by IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_products_search ON products USING gin (search_vector);

CREATE INDEX IF NOT EXISTS idx_products_seller_id ON products USING btree (seller_id);

CREATE INDEX IF NOT EXISTS idx_products_seller_listing_type ON products USING btree (seller_id, listing_type);

CREATE INDEX IF NOT EXISTS idx_products_status ON products USING btree (status);

CREATE INDEX IF NOT EXISTS idx_products_status_listing_type ON products USING btree (status, listing_type) WHERE (status = 'approved'::product_status);

CREATE INDEX IF NOT EXISTS idx_products_subcategory ON products USING btree (subcategory);

CREATE INDEX IF NOT EXISTS idx_products_title_trgm ON products USING gin (title gin_trgm_ops) WHERE (status = 'approved'::product_status);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'fk_products_category_hierarchy'
    ) THEN
        ALTER TABLE ONLY products ADD CONSTRAINT fk_products_category_hierarchy FOREIGN KEY (category_hierarchy_id) REFERENCES category_hierarchy(id) ON DELETE RESTRICT;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_brand_id_fkey'
    ) THEN
        ALTER TABLE ONLY products ADD CONSTRAINT products_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_category_id_fkey'
    ) THEN
        ALTER TABLE ONLY products ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_current_tenant_id_fkey'
    ) THEN
        ALTER TABLE ONLY products ADD CONSTRAINT products_current_tenant_id_fkey FOREIGN KEY (current_tenant_id) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_location_hierarchy_id_fkey'
    ) THEN
        ALTER TABLE ONLY products ADD CONSTRAINT products_location_hierarchy_id_fkey FOREIGN KEY (location_hierarchy_id) REFERENCES location_hierarchy(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_location_id_fkey'
    ) THEN
        ALTER TABLE ONLY products ADD CONSTRAINT products_location_id_fkey FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_product_model_id_fkey'
    ) THEN
        ALTER TABLE ONLY products ADD CONSTRAINT products_product_model_id_fkey FOREIGN KEY (product_model_id) REFERENCES product_models(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_reviewed_by_fkey'
    ) THEN
        ALTER TABLE ONLY products ADD CONSTRAINT products_reviewed_by_fkey FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'products' AND constraint_name = 'products_seller_id_fkey'
    ) THEN
        ALTER TABLE ONLY products ADD CONSTRAINT products_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
