-- ============================================================================
-- 049_qualification_types.sql
-- Table: qualification_types (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for qualification_types
ALTER TABLE qualification_types ADD COLUMN IF NOT EXISTS name_en character varying(200) NOT NULL;
ALTER TABLE qualification_types ADD COLUMN IF NOT EXISTS name_am character varying(200);
ALTER TABLE qualification_types ADD COLUMN IF NOT EXISTS name_om character varying(200);
ALTER TABLE qualification_types ADD COLUMN IF NOT EXISTS name_ti character varying(200);
ALTER TABLE qualification_types ADD COLUMN IF NOT EXISTS level_order integer DEFAULT 0 NOT NULL;
ALTER TABLE qualification_types ADD COLUMN IF NOT EXISTS category character varying(50) DEFAULT 'degree'::character varying NOT NULL;
ALTER TABLE qualification_types ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'qualification_types' AND constraint_name = 'qualification_types_pkey'
    ) THEN
        ALTER TABLE ONLY qualification_types ADD CONSTRAINT qualification_types_pkey PRIMARY KEY (id);
    END IF;
END $$;
