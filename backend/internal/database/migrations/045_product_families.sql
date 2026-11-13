-- ============================================================================
-- 045_product_families.sql
-- Table: product_families (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for product_families
ALTER TABLE product_families ADD COLUMN IF NOT EXISTS brand_id integer NOT NULL;
ALTER TABLE product_families ADD COLUMN IF NOT EXISTS category_id integer;
ALTER TABLE product_families ADD COLUMN IF NOT EXISTS name character varying(200) NOT NULL;
ALTER TABLE product_families ADD COLUMN IF NOT EXISTS slug character varying(200) NOT NULL;
ALTER TABLE product_families ADD COLUMN IF NOT EXISTS series_name character varying(100);
ALTER TABLE product_families ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;
ALTER TABLE product_families ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();
ALTER TABLE product_families ADD COLUMN IF NOT EXISTS updated_at timestamp without time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_families' AND constraint_name = 'product_families_pkey'
    ) THEN
        ALTER TABLE ONLY product_families ADD CONSTRAINT product_families_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_families' AND constraint_name = 'uq_brand_family_slug'
    ) THEN
        ALTER TABLE ONLY product_families ADD CONSTRAINT uq_brand_family_slug UNIQUE (brand_id, slug);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_product_families_brand ON product_families USING btree (brand_id);

CREATE INDEX IF NOT EXISTS idx_product_families_category ON product_families USING btree (category_id);

CREATE INDEX IF NOT EXISTS idx_product_families_slug ON product_families USING btree (slug);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_families' AND constraint_name = 'product_families_brand_id_fkey'
    ) THEN
        ALTER TABLE ONLY product_families ADD CONSTRAINT product_families_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_families' AND constraint_name = 'product_families_category_id_fkey'
    ) THEN
        ALTER TABLE ONLY product_families ADD CONSTRAINT product_families_category_id_fkey FOREIGN KEY (category_id) REFERENCES category_hierarchy(id) ON DELETE SET NULL;
    END IF;
END $$;
