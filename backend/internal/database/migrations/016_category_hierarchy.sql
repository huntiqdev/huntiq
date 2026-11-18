-- ============================================================================
-- 016_category_hierarchy.sql
-- Table: category_hierarchy (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for category_hierarchy
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS parent_id integer;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS name_en character varying(200) NOT NULL;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS name_am character varying(200);
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS name_om character varying(200);
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS slug character varying(200) NOT NULL;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS description text;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS level integer DEFAULT 0 NOT NULL;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS path character varying(500);
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS display_order integer DEFAULT 0;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS icon character varying(100);
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS is_leaf boolean DEFAULT false;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS meta_title character varying(200);
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS meta_description text;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS updated_at timestamp without time zone DEFAULT now();
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS name_ti character varying(200);
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS image_url text;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS is_archived boolean DEFAULT false;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS classification_system_code character varying(50) DEFAULT 'SYS_PROD'::character varying;
ALTER TABLE category_hierarchy ADD COLUMN IF NOT EXISTS i18n jsonb DEFAULT '{}'::jsonb NOT NULL;

-- Table Constraints for category_hierarchy
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_hierarchy' AND constraint_name = 'no_circular_reference'
    ) THEN
        ALTER TABLE category_hierarchy ADD CONSTRAINT no_circular_reference CHECK ((id <> parent_id));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_hierarchy' AND constraint_name = 'valid_level'
    ) THEN
        ALTER TABLE category_hierarchy ADD CONSTRAINT valid_level CHECK (((level >= 0) AND (level <= 3)));
    END IF;
END $$;

COMMENT ON TABLE category_hierarchy IS 'Dynamic N-level category hierarchy supporting Ethiopian marketplace';

COMMENT ON COLUMN category_hierarchy.level IS 'Depth in tree: 0=root, 1=child of root, etc.';

COMMENT ON COLUMN category_hierarchy.path IS 'Materialized path for fast ancestor/descendant queries. Format: /1/5/12/';

COMMENT ON COLUMN category_hierarchy.is_leaf IS 'True if products can be directly attached to this category';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_hierarchy' AND constraint_name = 'category_hierarchy_pkey'
    ) THEN
        ALTER TABLE ONLY category_hierarchy ADD CONSTRAINT category_hierarchy_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_hierarchy' AND constraint_name = 'uq_category_hierarchy_slug'
    ) THEN
        ALTER TABLE ONLY category_hierarchy ADD CONSTRAINT uq_category_hierarchy_slug UNIQUE (slug);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_active ON category_hierarchy USING btree (is_active);

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_archived ON category_hierarchy USING btree (is_archived);

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_classification_system ON category_hierarchy USING btree (classification_system_code);

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_display_order ON category_hierarchy USING btree (display_order);

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_image_url ON category_hierarchy USING btree (image_url) WHERE (image_url IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_leaf ON category_hierarchy USING btree (is_leaf);

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_level ON category_hierarchy USING btree (level);

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_parent ON category_hierarchy USING btree (parent_id);

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_path ON category_hierarchy USING btree (path);

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_path_ops ON category_hierarchy USING btree (path varchar_pattern_ops);

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_i18n ON category_hierarchy USING gin (i18n);

-- classification_systems FK pruned in consolidation

-- parent_id FK deferred to 057
