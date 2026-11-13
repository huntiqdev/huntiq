-- ============================================================================
-- 037_location_hierarchy.sql
-- Table: location_hierarchy (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for location_hierarchy
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS parent_id integer;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS location_type character varying(50) NOT NULL;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS name_en character varying(200) NOT NULL;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS name_am character varying(200);
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS name_om character varying(200);
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS name_ti character varying(200);
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS slug character varying(200) NOT NULL;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS admin_code character varying(50);
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS level integer DEFAULT 0 NOT NULL;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS path character varying(500);
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS display_order integer DEFAULT 0;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS is_leaf boolean DEFAULT false;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS is_deprecated boolean DEFAULT false;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS replaced_by_id integer;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS latitude numeric(10,7);
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS longitude numeric(10,7);
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS min_lat numeric(10,7);
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS max_lat numeric(10,7);
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS min_lng numeric(10,7);
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS max_lng numeric(10,7);
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS aliases text[] DEFAULT '{}'::text[];
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS source character varying(100) DEFAULT 'authoritative'::character varying;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS verification_status character varying(50) DEFAULT 'verified'::character varying;
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now();
ALTER TABLE location_hierarchy ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now();

-- Table Constraints for location_hierarchy
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'location_hierarchy' AND constraint_name = 'no_circular_location_ref'
    ) THEN
        ALTER TABLE location_hierarchy ADD CONSTRAINT no_circular_location_ref CHECK ((id <> parent_id));
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'location_hierarchy' AND constraint_name = 'valid_location_level'
    ) THEN
        ALTER TABLE location_hierarchy ADD CONSTRAINT valid_location_level CHECK (((level >= 0) AND (level < 10)));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'location_hierarchy' AND constraint_name = 'location_hierarchy_pkey'
    ) THEN
        ALTER TABLE ONLY location_hierarchy ADD CONSTRAINT location_hierarchy_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_loc_hier_active ON location_hierarchy USING btree (is_active);

CREATE INDEX IF NOT EXISTS idx_loc_hier_admin_code ON location_hierarchy USING btree (admin_code);

CREATE INDEX IF NOT EXISTS idx_loc_hier_bbox ON location_hierarchy USING btree (min_lat, max_lat, min_lng, max_lng) WHERE (min_lat IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_loc_hier_coords ON location_hierarchy USING btree (latitude, longitude) WHERE (latitude IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_loc_hier_level ON location_hierarchy USING btree (level);

CREATE INDEX IF NOT EXISTS idx_loc_hier_name_am_trgm ON location_hierarchy USING gin (name_am gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_loc_hier_name_en_trgm ON location_hierarchy USING gin (name_en gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_loc_hier_name_om_trgm ON location_hierarchy USING gin (name_om gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_loc_hier_name_ti_trgm ON location_hierarchy USING gin (name_ti gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_loc_hier_parent ON location_hierarchy USING btree (parent_id);

CREATE INDEX IF NOT EXISTS idx_loc_hier_path ON location_hierarchy USING btree (path);

CREATE INDEX IF NOT EXISTS idx_loc_hier_slug ON location_hierarchy USING btree (slug);

CREATE INDEX IF NOT EXISTS idx_loc_hier_type ON location_hierarchy USING btree (location_type);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'location_hierarchy' AND constraint_name = 'location_hierarchy_parent_id_fkey'
    ) THEN
        ALTER TABLE ONLY location_hierarchy ADD CONSTRAINT location_hierarchy_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES location_hierarchy(id) ON DELETE RESTRICT;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'location_hierarchy' AND constraint_name = 'location_hierarchy_replaced_by_id_fkey'
    ) THEN
        ALTER TABLE ONLY location_hierarchy ADD CONSTRAINT location_hierarchy_replaced_by_id_fkey FOREIGN KEY (replaced_by_id) REFERENCES location_hierarchy(id) ON DELETE SET NULL;
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS uq_loc_hier_slug ON location_hierarchy (slug);
