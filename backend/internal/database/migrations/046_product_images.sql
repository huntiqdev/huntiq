-- ============================================================================
-- 046_product_images.sql
-- Table: product_images (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for product_images
ALTER TABLE product_images ADD COLUMN IF NOT EXISTS product_id integer NOT NULL;
ALTER TABLE product_images ADD COLUMN IF NOT EXISTS image_url text DEFAULT ''::text NOT NULL;
ALTER TABLE product_images ADD COLUMN IF NOT EXISTS order_number smallint DEFAULT 0 NOT NULL;
ALTER TABLE product_images ADD COLUMN IF NOT EXISTS is_thumbnail boolean DEFAULT false NOT NULL;
ALTER TABLE product_images ADD COLUMN IF NOT EXISTS is_blurred boolean DEFAULT false NOT NULL;
ALTER TABLE product_images ADD COLUMN IF NOT EXISTS deleted_at timestamp with time zone;
ALTER TABLE product_images ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE product_images ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_images' AND constraint_name = 'product_images_pkey'
    ) THEN
        ALTER TABLE ONLY product_images ADD CONSTRAINT product_images_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_images' AND constraint_name = 'uq_product_image_order'
    ) THEN
        ALTER TABLE ONLY product_images ADD CONSTRAINT uq_product_image_order UNIQUE (product_id, order_number);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_product_images_blurred ON product_images USING btree (product_id) WHERE ((is_blurred = true) AND (deleted_at IS NULL));

CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images USING btree (product_id) WHERE (deleted_at IS NULL);

CREATE INDEX IF NOT EXISTS idx_product_images_thumbnail ON product_images USING btree (product_id) WHERE (is_thumbnail = true);

CREATE UNIQUE INDEX IF NOT EXISTS uq_product_thumbnail ON product_images USING btree (product_id) WHERE (is_thumbnail = true);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_images' AND constraint_name = 'product_images_product_id_fkey'
    ) THEN
        ALTER TABLE ONLY product_images ADD CONSTRAINT product_images_product_id_fkey FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
    END IF;
END $$;
