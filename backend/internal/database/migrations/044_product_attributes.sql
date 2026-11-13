-- ============================================================================
-- 044_product_attributes.sql
-- Table: product_attributes (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for product_attributes
ALTER TABLE product_attributes ADD COLUMN IF NOT EXISTS product_id integer NOT NULL;
ALTER TABLE product_attributes ADD COLUMN IF NOT EXISTS attribute_id integer NOT NULL;
ALTER TABLE product_attributes ADD COLUMN IF NOT EXISTS value_text text;
ALTER TABLE product_attributes ADD COLUMN IF NOT EXISTS value_number numeric(15,2);
ALTER TABLE product_attributes ADD COLUMN IF NOT EXISTS value_boolean boolean;
ALTER TABLE product_attributes ADD COLUMN IF NOT EXISTS value_date date;
ALTER TABLE product_attributes ADD COLUMN IF NOT EXISTS value_json jsonb;
ALTER TABLE product_attributes ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();
ALTER TABLE product_attributes ADD COLUMN IF NOT EXISTS updated_at timestamp without time zone DEFAULT now();
ALTER TABLE product_attributes ADD COLUMN IF NOT EXISTS search_vector tsvector;

COMMENT ON TABLE product_attributes IS 'Stores actual product attribute values';

COMMENT ON COLUMN product_attributes.value_json IS 'For multiselect attributes (stored as JSON array)';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_attributes' AND constraint_name = 'product_attributes_pkey'
    ) THEN
        ALTER TABLE ONLY product_attributes ADD CONSTRAINT product_attributes_pkey PRIMARY KEY (id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_attributes' AND constraint_name = 'product_attributes_product_id_attribute_id_key'
    ) THEN
        ALTER TABLE ONLY product_attributes ADD CONSTRAINT product_attributes_product_id_attribute_id_key UNIQUE (product_id, attribute_id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_product_attributes_attribute ON product_attributes USING btree (attribute_id);

CREATE INDEX IF NOT EXISTS idx_product_attributes_product ON product_attributes USING btree (product_id);

CREATE INDEX IF NOT EXISTS idx_product_attributes_search ON product_attributes USING gin (search_vector);

CREATE INDEX IF NOT EXISTS idx_product_attributes_value_json ON product_attributes USING gin (value_json);

CREATE INDEX IF NOT EXISTS idx_product_attributes_value_number ON product_attributes USING btree (value_number);

CREATE INDEX IF NOT EXISTS idx_product_attributes_value_text ON product_attributes USING btree (value_text);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_attributes' AND constraint_name = 'product_attributes_attribute_id_fkey'
    ) THEN
        ALTER TABLE ONLY product_attributes ADD CONSTRAINT product_attributes_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES attribute_definitions(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'product_attributes' AND constraint_name = 'product_attributes_product_id_fkey'
    ) THEN
        ALTER TABLE ONLY product_attributes ADD CONSTRAINT product_attributes_product_id_fkey FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
    END IF;
END $$;
