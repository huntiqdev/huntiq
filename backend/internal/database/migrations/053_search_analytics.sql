-- ============================================================================
-- 053_search_analytics.sql
-- Table: search_analytics (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for search_analytics
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS query text NOT NULL;
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS normalized_query text;
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS result_count integer DEFAULT 0;
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS category_id integer;
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS filters_used jsonb;
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS clicked_product_id integer;
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS clicked_position integer;
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS user_id integer;
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS session_id character varying(100);
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS ip_address character varying(45);
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS user_agent text;
ALTER TABLE search_analytics ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();

COMMENT ON TABLE search_analytics IS 'Search query tracking and analytics';

CREATE SEQUENCE IF NOT EXISTS search_analytics_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'search_analytics' AND constraint_name = 'search_analytics_pkey'
    ) THEN
        ALTER TABLE ONLY search_analytics ADD CONSTRAINT search_analytics_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_search_analytics_clicked ON search_analytics USING btree (clicked_product_id);

CREATE INDEX IF NOT EXISTS idx_search_analytics_created ON search_analytics USING btree (created_at);

CREATE INDEX IF NOT EXISTS idx_search_analytics_query ON search_analytics USING btree (normalized_query);

CREATE INDEX IF NOT EXISTS idx_search_analytics_user ON search_analytics USING btree (user_id);
