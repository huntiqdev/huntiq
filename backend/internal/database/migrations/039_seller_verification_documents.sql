-- ============================================================================
-- 039_seller_verification_documents.sql
-- Table: seller_verification_documents (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for seller_verification_documents
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS seller_id integer NOT NULL;
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS document_type character varying(50) NOT NULL;
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS file_url text NOT NULL;
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS file_name character varying(255);
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS file_size_bytes integer;
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'pending'::character varying NOT NULL;
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS reviewed_by integer;
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS reviewed_at timestamp with time zone;
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS reject_reason text;
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS notes text;
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE seller_verification_documents ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

COMMENT ON TABLE seller_verification_documents IS 'Seller identity and business verification documents';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'seller_verification_documents' AND constraint_name = 'seller_verification_documents_pkey'
    ) THEN
        ALTER TABLE ONLY seller_verification_documents ADD CONSTRAINT seller_verification_documents_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_verification_docs_reviewed ON seller_verification_documents USING btree (reviewed_by) WHERE (reviewed_by IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_verification_docs_seller ON seller_verification_documents USING btree (seller_id, status);

CREATE INDEX IF NOT EXISTS idx_verification_docs_status ON seller_verification_documents USING btree (status, created_at DESC);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'seller_verification_documents' AND constraint_name = 'seller_verification_documents_reviewed_by_fkey'
    ) THEN
        ALTER TABLE ONLY seller_verification_documents ADD CONSTRAINT seller_verification_documents_reviewed_by_fkey FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'seller_verification_documents' AND constraint_name = 'seller_verification_documents_seller_id_fkey'
    ) THEN
        ALTER TABLE ONLY seller_verification_documents ADD CONSTRAINT seller_verification_documents_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
