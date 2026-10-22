-- ============================================================
-- 003_locations.sql
-- PURPOSE : Ethiopian administrative hierarchy used to tag
--           product listings with a precise location.
--           Supports 4 languages: English, Amharic, Oromifa, Tigrinya.
-- HIERARCHY: region → zone → subcity → woreda
-- ============================================================
ALTER TABLE locations
    -- region     : top-level administrative region (e.g. 'Addis Ababa', 'Oromia')
    ADD COLUMN IF NOT EXISTS region     VARCHAR(100) NOT NULL DEFAULT '',
    -- zone       : zone within region (NULL for city administrations)
    ADD COLUMN IF NOT EXISTS zone       VARCHAR(100),
    -- subcity    : subcity within Addis Ababa or similar city admin
    ADD COLUMN IF NOT EXISTS subcity    VARCHAR(100),
    -- woreda     : smallest administrative unit
    ADD COLUMN IF NOT EXISTS woreda     VARCHAR(100),
    -- _am/_or/_ti: Amharic / Oromifa / Tigrinya translations for each level
    ADD COLUMN IF NOT EXISTS region_am  TEXT, ADD COLUMN IF NOT EXISTS region_or  TEXT, ADD COLUMN IF NOT EXISTS region_ti  TEXT,
    ADD COLUMN IF NOT EXISTS zone_am    TEXT, ADD COLUMN IF NOT EXISTS zone_or    TEXT, ADD COLUMN IF NOT EXISTS zone_ti    TEXT,
    ADD COLUMN IF NOT EXISTS subcity_am TEXT, ADD COLUMN IF NOT EXISTS subcity_or TEXT, ADD COLUMN IF NOT EXISTS subcity_ti TEXT,
    ADD COLUMN IF NOT EXISTS woreda_am  TEXT, ADD COLUMN IF NOT EXISTS woreda_or  TEXT, ADD COLUMN IF NOT EXISTS woreda_ti  TEXT,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- B-tree for exact match, GIN trigram for fuzzy search in all languages
CREATE INDEX IF NOT EXISTS idx_locations_region      ON locations (region);
CREATE INDEX IF NOT EXISTS idx_locations_subcity     ON locations (subcity);
CREATE INDEX IF NOT EXISTS idx_loc_region_trgm       ON locations USING gin (region     gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_loc_region_am_trgm    ON locations USING gin (region_am  gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_loc_zone_trgm         ON locations USING gin (zone       gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_loc_subcity_trgm      ON locations USING gin (subcity    gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_loc_subcity_am_trgm   ON locations USING gin (subcity_am gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_loc_woreda_trgm       ON locations USING gin (woreda     gin_trgm_ops);
