-- ============================================================================
-- Migration 026: Configurable Multi-Level Location Hierarchy System
-- ============================================================================
-- Purpose: Provide an extensible, data-driven administrative location hierarchy
--          supporting Country, Region, Zone, City, Subcity, Woreda, and Kebele
--          with coordinates, bounding boxes, multilingual names, and aliases.
-- Status: Non-destructive — preserves existing locations and products tables.
-- ============================================================================

CREATE TABLE IF NOT EXISTS location_hierarchy (
    id SERIAL PRIMARY KEY,
    parent_id INT REFERENCES location_hierarchy(id) ON DELETE RESTRICT,
    
    -- Level type classification: 'country', 'region', 'zone', 'city', 'subcity', 'woreda', 'kebele'
    location_type VARCHAR(50) NOT NULL,
    
    -- Multilingual names (English, Amharic, Afaan Oromo, Tigrinya)
    name_en VARCHAR(200) NOT NULL,
    name_am VARCHAR(200),
    name_om VARCHAR(200),
    name_ti VARCHAR(200),
    
    slug VARCHAR(200) NOT NULL,
    admin_code VARCHAR(50), -- e.g. 'ET', 'ET-AA', 'ET-OR'
    
    -- Hierarchy metadata
    level INT NOT NULL DEFAULT 0, -- 0=country, 1=region, 2=zone/city, 3=woreda/subcity, 4=kebele
    path VARCHAR(500),            -- Materialized path: e.g. /1/2/5/ for fast descendant filtering
    display_order INT DEFAULT 0,
    
    -- Status flags
    is_active BOOLEAN DEFAULT true,
    is_leaf BOOLEAN DEFAULT false,
    is_deprecated BOOLEAN DEFAULT false,
    replaced_by_id INT REFERENCES location_hierarchy(id) ON DELETE SET NULL,
    
    -- Geographic centroids and bounding box for GPS reverse-geocoding resolution
    latitude NUMERIC(10,7),
    longitude NUMERIC(10,7),
    min_lat NUMERIC(10,7),
    max_lat NUMERIC(10,7),
    min_lng NUMERIC(10,7),
    max_lng NUMERIC(10,7),
    
    -- Searchable aliases & alternate spellings (e.g. ['Finfinnee', 'Addis Abeba', 'Sheger'])
    aliases TEXT[] DEFAULT '{}',
    
    -- Provenance & metadata
    source VARCHAR(100) DEFAULT 'authoritative',
    verification_status VARCHAR(50) DEFAULT 'verified',
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT no_circular_location_ref CHECK (id != parent_id),
    CONSTRAINT valid_location_level CHECK (level >= 0 AND level < 10)
);

-- Indexes for fast queries
CREATE INDEX IF NOT EXISTS idx_loc_hier_parent ON location_hierarchy(parent_id);
CREATE INDEX IF NOT EXISTS idx_loc_hier_type ON location_hierarchy(location_type);
CREATE INDEX IF NOT EXISTS idx_loc_hier_level ON location_hierarchy(level);
CREATE INDEX IF NOT EXISTS idx_loc_hier_path ON location_hierarchy(path);
CREATE INDEX IF NOT EXISTS idx_loc_hier_active ON location_hierarchy(is_active);
CREATE INDEX IF NOT EXISTS idx_loc_hier_slug ON location_hierarchy(slug);
CREATE INDEX IF NOT EXISTS idx_loc_hier_admin_code ON location_hierarchy(admin_code);
CREATE INDEX IF NOT EXISTS idx_loc_hier_coords ON location_hierarchy(latitude, longitude) WHERE latitude IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_loc_hier_bbox ON location_hierarchy(min_lat, max_lat, min_lng, max_lng) WHERE min_lat IS NOT NULL;

-- Trigram fuzzy search on localized names and aliases
CREATE INDEX IF NOT EXISTS idx_loc_hier_name_en_trgm ON location_hierarchy USING gin (name_en gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_loc_hier_name_am_trgm ON location_hierarchy USING gin (name_am gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_loc_hier_name_om_trgm ON location_hierarchy USING gin (name_om gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_loc_hier_name_ti_trgm ON location_hierarchy USING gin (name_ti gin_trgm_ops);

-- Materialized path update trigger
CREATE OR REPLACE FUNCTION update_location_hierarchy_path()
RETURNS TRIGGER AS $$
DECLARE
    parent_path VARCHAR(500);
    parent_level INT;
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.level := 0;
        NEW.path := '/' || NEW.id || '/';
    ELSE
        SELECT path, level INTO parent_path, parent_level
        FROM location_hierarchy
        WHERE id = NEW.parent_id;
        
        IF parent_path IS NULL THEN
            NEW.level := 1;
            NEW.path := '/' || NEW.parent_id || '/' || NEW.id || '/';
        ELSE
            NEW.level := parent_level + 1;
            NEW.path := parent_path || NEW.id || '/';
        END IF;
    END IF;
    
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_location_path ON location_hierarchy;
CREATE TRIGGER trg_location_path
BEFORE INSERT OR UPDATE OF parent_id, id ON location_hierarchy
FOR EACH ROW EXECUTE FUNCTION update_location_hierarchy_path();

-- Helper function: get ancestor IDs
CREATE OR REPLACE FUNCTION get_location_ancestors(loc_id INT)
RETURNS TABLE (id INT, location_type VARCHAR, name_en VARCHAR, name_am VARCHAR, name_om VARCHAR, name_ti VARCHAR, level INT) AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE ancestors AS (
        SELECT l.id, l.parent_id, l.location_type, l.name_en, l.name_am, l.name_om, l.name_ti, l.level
        FROM location_hierarchy l
        WHERE l.id = loc_id
        
        UNION ALL
        
        SELECT p.id, p.parent_id, p.location_type, p.name_en, p.name_am, p.name_om, p.name_ti, p.level
        FROM location_hierarchy p
        JOIN ancestors a ON a.parent_id = p.id
    )
    SELECT a.id, a.location_type, a.name_en, a.name_am, a.name_om, a.name_ti, a.level
    FROM ancestors a
    ORDER BY a.level ASC;
END;
$$ LANGUAGE plpgsql;

-- Helper function: get descendant IDs (hierarchical subtree)
CREATE OR REPLACE FUNCTION get_location_descendants(loc_id INT)
RETURNS TABLE (id INT) AS $$
DECLARE
    target_path VARCHAR(500);
BEGIN
    SELECT path INTO target_path FROM location_hierarchy WHERE location_hierarchy.id = loc_id;
    RETURN QUERY
    SELECT l.id
    FROM location_hierarchy l
    WHERE l.path LIKE target_path || '%';
END;
$$ LANGUAGE plpgsql;

-- Add location_hierarchy_id column to products table for rich hierarchy linking
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS location_hierarchy_id INT REFERENCES location_hierarchy(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_products_location_hierarchy ON products(location_hierarchy_id);
