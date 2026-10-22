-- ============================================================
-- 004_categories.sql
-- PURPOSE : Hierarchical product categories with multilingual names.
--           Supports 2 levels: parent category → subcategory.
--           parent_id NULL = top-level category.
-- LANGUAGES: English (name_en), Amharic (name_am),
--            Oromifa (name_or), Tigrinya (name_ti)
-- ============================================================
ALTER TABLE categories
    -- name      : canonical English name used as internal key
    ADD COLUMN IF NOT EXISTS name       VARCHAR(100) NOT NULL DEFAULT '' UNIQUE,
    -- slug      : URL-safe identifier e.g. 'home_appliances'
    ADD COLUMN IF NOT EXISTS slug       VARCHAR(100) NOT NULL DEFAULT '' UNIQUE,
    -- name_en/am/or/ti: display names per language
    ADD COLUMN IF NOT EXISTS name_en    TEXT,
    ADD COLUMN IF NOT EXISTS name_am    TEXT,
    ADD COLUMN IF NOT EXISTS name_or    TEXT,
    ADD COLUMN IF NOT EXISTS name_ti    TEXT,
    -- parent_id : NULL = top-level. Set = subcategory of parent.
    --             ON DELETE CASCADE: deleting parent removes all subcategories
    ADD COLUMN IF NOT EXISTS parent_id  INT     REFERENCES categories(id) ON DELETE CASCADE,
    -- sort_order: display order within same parent level (lower = shown first)
    ADD COLUMN IF NOT EXISTS sort_order INT     NOT NULL DEFAULT 0,
    -- is_active : FALSE = hidden from frontend, still exists in DB
    ADD COLUMN IF NOT EXISTS is_active  BOOLEAN NOT NULL DEFAULT TRUE;

CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories (parent_id) WHERE parent_id IS NOT NULL;
