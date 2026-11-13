-- ============================================================================
-- 025_classification_systems.sql
-- Table: classification_systems (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for classification_systems
ALTER TABLE classification_systems ADD COLUMN IF NOT EXISTS name_en character varying(100) NOT NULL;
ALTER TABLE classification_systems ADD COLUMN IF NOT EXISTS name_am character varying(100) NOT NULL;
ALTER TABLE classification_systems ADD COLUMN IF NOT EXISTS name_om character varying(100) NOT NULL;
ALTER TABLE classification_systems ADD COLUMN IF NOT EXISTS name_ti character varying(100) NOT NULL;
ALTER TABLE classification_systems ADD COLUMN IF NOT EXISTS description text;
ALTER TABLE classification_systems ADD COLUMN IF NOT EXISTS icon character varying(100);
ALTER TABLE classification_systems ADD COLUMN IF NOT EXISTS display_order integer DEFAULT 0;
ALTER TABLE classification_systems ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT true;
ALTER TABLE classification_systems ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now();
ALTER TABLE classification_systems ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'classification_systems' AND constraint_name = 'classification_systems_pkey'
    ) THEN
        ALTER TABLE ONLY classification_systems ADD CONSTRAINT classification_systems_pkey PRIMARY KEY (code);
    END IF;
END $$;

-- Seed default product marketplace classification system
INSERT INTO classification_systems (code, name_en, name_am, name_om, name_ti, description, is_active)
VALUES ('SYS_PROD', 'Products', 'ምርቶች', 'Oomishaalee', 'ፍርያት', 'Standard marketplace products classification system', true)
ON CONFLICT (code) DO NOTHING;
