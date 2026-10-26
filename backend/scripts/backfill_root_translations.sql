-- HUNTIQ Root Category Translation Backfill (Idempotent)
-- Target: public.category_hierarchy (IDs 1..11)
-- Purpose: Backfills Amharic, Oromo, and Tigrinya translations from legacy categories table (AUD-002)
-- Safety: Only updates rows where destination translation is empty ('') and preserves existing translations.

BEGIN;

UPDATE category_hierarchy ch
SET 
    name_am = CASE 
        WHEN COALESCE(ch.name_am, '') = '' AND COALESCE(c.name_am, '') != '' THEN c.name_am 
        ELSE ch.name_am 
    END,
    name_om = CASE 
        WHEN COALESCE(ch.name_om, '') = '' AND COALESCE(c.name_or, '') != '' THEN c.name_or 
        ELSE ch.name_om 
    END,
    name_ti = CASE 
        WHEN COALESCE(ch.name_ti, '') = '' AND COALESCE(c.name_ti, '') != '' THEN c.name_ti 
        ELSE ch.name_ti 
    END,
    updated_at = NOW()
FROM categories c
WHERE ch.id = c.id
  AND ch.id BETWEEN 1 AND 11
  AND ch.level = 0
  AND (
      (COALESCE(ch.name_am, '') = '' AND COALESCE(c.name_am, '') != '') OR
      (COALESCE(ch.name_om, '') = '' AND COALESCE(c.name_or, '') != '') OR
      (COALESCE(ch.name_ti, '') = '' AND COALESCE(c.name_ti, '') != '')
  );

-- Verification
DO $$
DECLARE
    v_missing INT;
BEGIN
    SELECT COUNT(*) INTO v_missing
    FROM category_hierarchy
    WHERE id BETWEEN 1 AND 11
      AND (COALESCE(name_am, '') = '' OR COALESCE(name_om, '') = '' OR COALESCE(name_ti, '') = '');

    IF v_missing > 0 THEN
        RAISE EXCEPTION 'Translation backfill verification failed: % root categories still have blank translations', v_missing;
    END IF;
    RAISE NOTICE 'Translation backfill verification successful: all 11 root categories translated.';
END $$;

COMMIT;
