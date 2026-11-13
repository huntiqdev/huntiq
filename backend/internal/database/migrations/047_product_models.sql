-- ============================================================================
-- 047_product_models.sql
-- Table: product_models (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for product_models
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS family_id integer;
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS brand_id integer NOT NULL;
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS category_id integer;
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS name character varying(200) NOT NULL;
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS slug character varying(200) NOT NULL;
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS model_number character varying(100);
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS release_year integer;
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS default_specs jsonb;
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS is_verified boolean DEFAULT true;
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();
ALTER TABLE product_models ADD COLUMN IF NOT EXISTS updated_at timestamp without time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_models' AND constraint_name = 'product_models_pkey'
    ) THEN
        ALTER TABLE ONLY product_models ADD CONSTRAINT product_models_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_models' AND constraint_name = 'uq_brand_model_slug'
    ) THEN
        ALTER TABLE ONLY product_models ADD CONSTRAINT uq_brand_model_slug UNIQUE (brand_id, slug);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_product_models_brand ON product_models USING btree (brand_id);

CREATE INDEX IF NOT EXISTS idx_product_models_category ON product_models USING btree (category_id);

CREATE INDEX IF NOT EXISTS idx_product_models_family ON product_models USING btree (family_id);

CREATE INDEX IF NOT EXISTS idx_product_models_slug ON product_models USING btree (slug);

CREATE INDEX IF NOT EXISTS idx_product_models_specs ON product_models USING gin (default_specs);

CREATE INDEX IF NOT EXISTS idx_product_models_year ON product_models USING btree (release_year);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_models' AND constraint_name = 'product_models_brand_id_fkey'
    ) THEN
        ALTER TABLE ONLY product_models ADD CONSTRAINT product_models_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_models' AND constraint_name = 'product_models_category_id_fkey'
    ) THEN
        ALTER TABLE ONLY product_models ADD CONSTRAINT product_models_category_id_fkey FOREIGN KEY (category_id) REFERENCES category_hierarchy(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_models' AND constraint_name = 'product_models_family_id_fkey'
    ) THEN
        ALTER TABLE ONLY product_models ADD CONSTRAINT product_models_family_id_fkey FOREIGN KEY (family_id) REFERENCES product_families(id) ON DELETE SET NULL;
    END IF;
END $$;
