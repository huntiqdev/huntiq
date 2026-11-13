-- ============================================================================
-- 013_attribute_definitions.sql
-- Table: attribute_definitions (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for attribute_definitions
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS name_en character varying(100) NOT NULL;
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS name_am character varying(100);
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS name_om character varying(100);
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS slug character varying(100) NOT NULL;
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS description text;
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS data_type character varying(50) DEFAULT 'text'::character varying NOT NULL;
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS validation_rules jsonb;
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS is_required boolean DEFAULT false;
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS is_searchable boolean DEFAULT true;
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS is_filterable boolean DEFAULT true;
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS display_order integer DEFAULT 0;
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS unit character varying(50);
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS updated_at timestamp without time zone DEFAULT now();
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS name_ti character varying(200);

COMMENT ON TABLE attribute_definitions IS 'Global attribute templates (Brand, Model, etc.)';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'attribute_definitions' AND constraint_name = 'attribute_definitions_pkey'
    ) THEN
        ALTER TABLE ONLY attribute_definitions ADD CONSTRAINT attribute_definitions_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'attribute_definitions' AND constraint_name = 'attribute_definitions_slug_key'
    ) THEN
        ALTER TABLE ONLY attribute_definitions ADD CONSTRAINT attribute_definitions_slug_key UNIQUE (slug);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_attribute_definitions_data_type ON attribute_definitions USING btree (data_type);
