-- ============================================================================
-- Migration 040: Products Trigram Full-Text Search Acceleration
-- ============================================================================
-- Purpose:
-- Accelerates substring and fuzzy search on products table (title, description)
-- using GIN trigram indexes for sub-5ms search during AI tool executions.
-- Compatible with PostgreSQL 16+ pg_trgm extension.
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- GIN trigram index for fast typo-tolerant substring search on product titles
CREATE INDEX IF NOT EXISTS idx_products_title_trgm
    ON products USING gin (title gin_trgm_ops)
    WHERE deleted_at IS NULL AND status = 'approved';

-- GIN trigram index for fast substring search on product descriptions
CREATE INDEX IF NOT EXISTS idx_products_description_trgm
    ON products USING gin (description gin_trgm_ops)
    WHERE deleted_at IS NULL AND status = 'approved';
