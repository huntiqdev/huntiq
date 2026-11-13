-- ============================================================================
-- 054_seller_reviews.sql
-- Table: seller_reviews (Attributes, Constraints & Indexes)
-- Master table creation is located in 001_schema.sql
-- Dedicated single-table migration (HUNTIQ 1-table-per-migration architecture)
-- ============================================================================

-- Table Columns for seller_reviews
ALTER TABLE seller_reviews ADD COLUMN IF NOT EXISTS seller_id integer DEFAULT 0 NOT NULL;
ALTER TABLE seller_reviews ADD COLUMN IF NOT EXISTS reviewer_id integer DEFAULT 0 NOT NULL;
ALTER TABLE seller_reviews ADD COLUMN IF NOT EXISTS order_id integer;
ALTER TABLE seller_reviews ADD COLUMN IF NOT EXISTS rating numeric(2,1) DEFAULT 1 NOT NULL;
ALTER TABLE seller_reviews ADD COLUMN IF NOT EXISTS comment text;
ALTER TABLE seller_reviews ADD COLUMN IF NOT EXISTS created_at timestamp with time zone DEFAULT now() NOT NULL;

-- Table Constraints for seller_reviews
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'seller_reviews' AND constraint_name = 'seller_reviews_rating_check'
    ) THEN
        ALTER TABLE seller_reviews ADD CONSTRAINT seller_reviews_rating_check CHECK (((rating >= (1)::numeric) AND (rating <= (5)::numeric)));
    END IF;
END $$;

COMMENT ON TABLE seller_reviews IS 'Buyer reviews for sellers after completed orders';

COMMENT ON COLUMN seller_reviews.rating IS 'Rating from 1.0 to 5.0 stars';

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'seller_reviews' AND constraint_name = 'seller_reviews_pkey'
    ) THEN
        ALTER TABLE ONLY seller_reviews ADD CONSTRAINT seller_reviews_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_seller_reviews_order_id ON seller_reviews USING btree (order_id) WHERE (order_id IS NOT NULL);

CREATE INDEX IF NOT EXISTS idx_seller_reviews_rating ON seller_reviews USING btree (rating);

CREATE INDEX IF NOT EXISTS idx_seller_reviews_reviewer_id ON seller_reviews USING btree (reviewer_id);

CREATE INDEX IF NOT EXISTS idx_seller_reviews_seller_id ON seller_reviews USING btree (seller_id);

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'seller_reviews' AND constraint_name = 'seller_reviews_order_id_fkey'
    ) THEN
        ALTER TABLE ONLY seller_reviews ADD CONSTRAINT seller_reviews_order_id_fkey FOREIGN KEY (order_id) REFERENCES order_requests(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'seller_reviews' AND constraint_name = 'seller_reviews_reviewer_id_fkey'
    ) THEN
        ALTER TABLE ONLY seller_reviews ADD CONSTRAINT seller_reviews_reviewer_id_fkey FOREIGN KEY (reviewer_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'seller_reviews' AND constraint_name = 'seller_reviews_seller_id_fkey'
    ) THEN
        ALTER TABLE ONLY seller_reviews ADD CONSTRAINT seller_reviews_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
