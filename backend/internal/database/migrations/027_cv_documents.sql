-- ============================================================================
-- 027_cv_documents.sql
-- Table: cv_documents (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for cv_documents
ALTER TABLE cv_documents ADD COLUMN IF NOT EXISTS user_id integer NOT NULL;
ALTER TABLE cv_documents ADD COLUMN IF NOT EXISTS file_url text NOT NULL;
ALTER TABLE cv_documents ADD COLUMN IF NOT EXISTS file_name character varying(255) NOT NULL;
ALTER TABLE cv_documents ADD COLUMN IF NOT EXISTS file_size integer NOT NULL;
ALTER TABLE cv_documents ADD COLUMN IF NOT EXISTS mime_type character varying(100) NOT NULL;
ALTER TABLE cv_documents ADD COLUMN IF NOT EXISTS is_private boolean DEFAULT true;
ALTER TABLE cv_documents ADD COLUMN IF NOT EXISTS download_count integer DEFAULT 0;
ALTER TABLE cv_documents ADD COLUMN IF NOT EXISTS created_at timestamp without time zone DEFAULT now();
ALTER TABLE cv_documents ADD COLUMN IF NOT EXISTS updated_at timestamp without time zone DEFAULT now();

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'cv_documents' AND constraint_name = 'cv_documents_pkey'
    ) THEN
        ALTER TABLE ONLY cv_documents ADD CONSTRAINT cv_documents_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_cv_documents_user ON cv_documents USING btree (user_id);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'cv_documents' AND constraint_name = 'cv_documents_user_id_fkey'
    ) THEN
        ALTER TABLE ONLY cv_documents ADD CONSTRAINT cv_documents_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
