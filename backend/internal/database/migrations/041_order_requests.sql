-- ============================================================================
-- 041_order_requests.sql
-- Table: order_requests (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for order_requests
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS product_id integer DEFAULT 0 NOT NULL;
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS seller_id integer DEFAULT 0 NOT NULL;
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS buyer_name character varying(100) DEFAULT ''::character varying NOT NULL;
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS buyer_phone character varying(20) DEFAULT ''::character varying NOT NULL;
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS buyer_email character varying(150);
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS buyer_message text;
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS message text;
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS offer_price numeric(12,2);
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS status character varying(20) DEFAULT 'pending'::character varying NOT NULL;
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS seller_note text;
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS buyer_tier_at_order character varying(20) DEFAULT 'normal'::character varying;
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE order_requests ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone DEFAULT now() NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'order_requests' AND constraint_name = 'order_requests_pkey'
    ) THEN
        ALTER TABLE ONLY order_requests ADD CONSTRAINT order_requests_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_order_requests_buyer_phone ON order_requests USING btree (buyer_phone);

CREATE INDEX IF NOT EXISTS idx_order_requests_product ON order_requests USING btree (product_id);

CREATE INDEX IF NOT EXISTS idx_order_requests_seller ON order_requests USING btree (seller_id, status);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'order_requests' AND constraint_name = 'order_requests_product_id_fkey'
    ) THEN
        ALTER TABLE ONLY order_requests ADD CONSTRAINT order_requests_product_id_fkey FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'order_requests' AND constraint_name = 'order_requests_seller_id_fkey'
    ) THEN
        ALTER TABLE ONLY order_requests ADD CONSTRAINT order_requests_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
