-- ============================================================================
-- 057_schema_integrity_and_sequences.sql
-- Sequence Alignment, Foreign Keys & Integrity Guarantees
-- ============================================================================

BEGIN;

-- Synchronize sequences for seeded auto-increment primary keys
SELECT setval(pg_get_serial_sequence('category_hierarchy', 'id'), COALESCE(MAX(id), 1)) FROM category_hierarchy;
SELECT setval(pg_get_serial_sequence('locations', 'id'), COALESCE(MAX(id), 1)) FROM locations;
SELECT setval(pg_get_serial_sequence('location_hierarchy', 'id'), COALESCE(MAX(id), 1)) FROM location_hierarchy;
SELECT setval(pg_get_serial_sequence('users', 'id'), COALESCE(MAX(id), 1)) FROM users;
SELECT setval(pg_get_serial_sequence('products', 'id'), COALESCE(MAX(id), 1)) FROM products;
SELECT setval(pg_get_serial_sequence('product_images', 'id'), COALESCE(MAX(id), 1)) FROM product_images;
SELECT setval(pg_get_serial_sequence('product_attributes', 'id'), COALESCE(MAX(id), 1)) FROM product_attributes;

-- Re-assert category_id synchronization with category_hierarchy_id
UPDATE products SET category_id = category_hierarchy_id WHERE category_hierarchy_id IS NOT NULL;

-- Category hierarchy self-referential parent_id constraint
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'category_hierarchy' AND constraint_name = 'category_hierarchy_parent_id_fkey'
    ) THEN
        ALTER TABLE ONLY category_hierarchy ADD CONSTRAINT category_hierarchy_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES category_hierarchy(id) ON DELETE RESTRICT;
    END IF;
END $$;

COMMIT;
