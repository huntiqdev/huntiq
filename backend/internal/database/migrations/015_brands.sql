-- ============================================================================
-- 015_brands.sql
-- Table: brands (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for brands
ALTER TABLE brands ADD COLUMN IF NOT EXISTS name character varying(200) NOT NULL;
ALTER TABLE brands ADD COLUMN IF NOT EXISTS slug character varying(200) NOT NULL;
ALTER TABLE brands ADD COLUMN IF NOT EXISTS category_ids integer[];
ALTER TABLE brands ADD COLUMN IF NOT EXISTS logo_url text;
ALTER TABLE brands ADD COLUMN IF NOT EXISTS country_of_origin character varying(100);
ALTER TABLE brands ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;
ALTER TABLE brands ADD COLUMN IF NOT EXISTS is_verified boolean DEFAULT true;
ALTER TABLE brands ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();
ALTER TABLE brands ADD COLUMN IF NOT EXISTS updated_at timestamp without time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'brands' AND constraint_name = 'brands_pkey'
    ) THEN
        ALTER TABLE ONLY brands ADD CONSTRAINT brands_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'brands' AND constraint_name = 'brands_slug_key'
    ) THEN
        ALTER TABLE ONLY brands ADD CONSTRAINT brands_slug_key UNIQUE (slug);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_brands_active ON brands USING btree (is_active);

CREATE INDEX IF NOT EXISTS idx_brands_category_ids ON brands USING gin (category_ids);
