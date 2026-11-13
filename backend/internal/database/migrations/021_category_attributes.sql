-- ============================================================================
-- 021_category_attributes.sql
-- Table: category_attributes (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for category_attributes
ALTER TABLE category_attributes ADD COLUMN IF NOT EXISTS category_id integer NOT NULL;
ALTER TABLE category_attributes ADD COLUMN IF NOT EXISTS attribute_id integer NOT NULL;
ALTER TABLE category_attributes ADD COLUMN IF NOT EXISTS is_required boolean DEFAULT false;
ALTER TABLE category_attributes ADD COLUMN IF NOT EXISTS is_inherited boolean DEFAULT false;
ALTER TABLE category_attributes ADD COLUMN IF NOT EXISTS display_order integer DEFAULT 0;
ALTER TABLE category_attributes ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();

COMMENT ON TABLE category_attributes IS 'Maps attributes to categories';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_attributes' AND constraint_name = 'category_attributes_category_id_attribute_id_key'
    ) THEN
        ALTER TABLE ONLY category_attributes ADD CONSTRAINT category_attributes_category_id_attribute_id_key UNIQUE (category_id, attribute_id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_attributes' AND constraint_name = 'category_attributes_pkey'
    ) THEN
        ALTER TABLE ONLY category_attributes ADD CONSTRAINT category_attributes_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_category_attributes_attribute ON category_attributes USING btree (attribute_id);

CREATE INDEX IF NOT EXISTS idx_category_attributes_category ON category_attributes USING btree (category_id);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_attributes' AND constraint_name = 'category_attributes_attribute_id_fkey'
    ) THEN
        ALTER TABLE ONLY category_attributes ADD CONSTRAINT category_attributes_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES attribute_definitions(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_attributes' AND constraint_name = 'category_attributes_category_id_fkey'
    ) THEN
        ALTER TABLE ONLY category_attributes ADD CONSTRAINT category_attributes_category_id_fkey FOREIGN KEY (category_id) REFERENCES category_hierarchy(id) ON DELETE CASCADE;
    END IF;
END $$;
