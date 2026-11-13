-- ============================================================================
-- 038_locations.sql
-- Table: locations (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for locations
ALTER TABLE locations ADD COLUMN IF NOT EXISTS region character varying(100) DEFAULT ''::character varying NOT NULL;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS zone character varying(100);
ALTER TABLE locations ADD COLUMN IF NOT EXISTS subcity character varying(100);
ALTER TABLE locations ADD COLUMN IF NOT EXISTS woreda character varying(100);
ALTER TABLE locations ADD COLUMN IF NOT EXISTS region_am text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS region_or text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS region_ti text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS zone_am text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS zone_or text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS zone_ti text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS subcity_am text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS subcity_or text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS subcity_ti text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS woreda_am text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS woreda_or text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS woreda_ti text;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'locations' AND constraint_name = 'locations_pkey'
    ) THEN
        ALTER TABLE ONLY locations ADD CONSTRAINT locations_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_loc_region_am_trgm ON locations USING gin (region_am gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_loc_region_trgm ON locations USING gin (region gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_loc_subcity_am_trgm ON locations USING gin (subcity_am gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_loc_subcity_trgm ON locations USING gin (subcity gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_loc_woreda_trgm ON locations USING gin (woreda gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_loc_zone_trgm ON locations USING gin (zone gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_locations_region ON locations USING btree (region);

CREATE INDEX IF NOT EXISTS idx_locations_subcity ON locations USING btree (subcity);
