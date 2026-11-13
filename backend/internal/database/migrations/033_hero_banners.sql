-- ============================================================================
-- 033_hero_banners.sql
-- Table: hero_banners (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for hero_banners
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS banner_type character varying(20) DEFAULT 'seller'::character varying NOT NULL;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS media_type character varying(20) DEFAULT 'image'::character varying NOT NULL;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS media_url text DEFAULT ''::text NOT NULL;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS title character varying(200) DEFAULT ''::character varying NOT NULL;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS subtitle text;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS link_url text;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS promotion_text text;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS target_audience character varying(50) DEFAULT 'all'::character varying NOT NULL;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS sort_order integer DEFAULT 0 NOT NULL;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true NOT NULL;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS starts_at timestamp with time zone;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS ends_at timestamp with time zone;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS created_by integer;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE hero_banners ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'hero_banners' AND constraint_name = 'hero_banners_pkey'
    ) THEN
        ALTER TABLE ONLY hero_banners ADD CONSTRAINT hero_banners_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_hero_banners_active ON hero_banners USING btree (is_active, sort_order) WHERE (is_active = true);

CREATE INDEX IF NOT EXISTS idx_hero_banners_dates ON hero_banners USING btree (starts_at, ends_at) WHERE (is_active = true);

CREATE INDEX IF NOT EXISTS idx_hero_banners_type ON hero_banners USING btree (banner_type) WHERE (is_active = true);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'hero_banners' AND constraint_name = 'hero_banners_created_by_fkey'
    ) THEN
        ALTER TABLE ONLY hero_banners ADD CONSTRAINT hero_banners_created_by_fkey FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;
