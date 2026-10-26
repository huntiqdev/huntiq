-- ============================================================================
-- Migration 032: Category Image URL and Taxonomy Metadata
-- ============================================================================
ALTER TABLE category_hierarchy 
    ADD COLUMN IF NOT EXISTS image_url TEXT;

CREATE INDEX IF NOT EXISTS idx_category_hierarchy_image_url 
    ON category_hierarchy(image_url) WHERE image_url IS NOT NULL;
