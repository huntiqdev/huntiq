-- ============================================================================
-- 061_classification_placements_instead_of_triggers.sql
-- Enable INSERT, UPDATE, and DELETE on classification_placements view
-- ============================================================================

CREATE OR REPLACE FUNCTION trg_classification_placements_insert()
RETURNS trigger AS $$
DECLARE
    v_scheme_code VARCHAR := 'SYS_PROD';
    v_parent_level INT := 0;
    v_parent_path VARCHAR := '';
    v_new_id INT;
    v_level INT := 0;
    v_path VARCHAR := '';
BEGIN
    IF NEW.scheme_id IS NOT NULL THEN
        SELECT code INTO v_scheme_code FROM classification_schemes WHERE id = NEW.scheme_id;
        IF v_scheme_code IS NULL THEN
            v_scheme_code := 'SYS_PROD';
        END IF;
    END IF;

    IF NEW.parent_placement_id IS NOT NULL AND NEW.parent_placement_id > 0 THEN
        SELECT level, path INTO v_parent_level, v_parent_path FROM category_hierarchy WHERE id = NEW.parent_placement_id;
        v_level := LEAST(3, COALESCE(v_parent_level, 0) + 1);
    END IF;

    INSERT INTO category_hierarchy (
        parent_id,
        name_en,
        slug,
        display_order,
        is_leaf,
        is_active,
        level,
        classification_system_code
    ) VALUES (
        NEW.parent_placement_id,
        COALESCE(NEW.display_name_override, 'Category ' || NEW.slug),
        NEW.slug,
        COALESCE(NEW.display_order, 0),
        COALESCE(NEW.is_leaf, false),
        COALESCE(NEW.is_active, true),
        v_level,
        v_scheme_code
    ) RETURNING id INTO v_new_id;

    IF v_parent_path IS NOT NULL AND v_parent_path != '' THEN
        v_path := v_parent_path || '/' || v_new_id;
    ELSE
        v_path := '/' || v_new_id;
    END IF;

    UPDATE category_hierarchy SET path = v_path WHERE id = v_new_id;

    NEW.id := v_new_id;
    NEW.path := v_path;
    NEW.depth := v_level;
    NEW.created_at := NOW();
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trg_classification_placements_update()
RETURNS trigger AS $$
DECLARE
    v_parent_level INT := 0;
    v_parent_path VARCHAR := '';
    v_level INT := 0;
    v_path VARCHAR := '';
BEGIN
    IF NEW.parent_placement_id IS NOT NULL AND NEW.parent_placement_id > 0 THEN
        SELECT level, path INTO v_parent_level, v_parent_path FROM category_hierarchy WHERE id = NEW.parent_placement_id;
        v_level := LEAST(3, COALESCE(v_parent_level, 0) + 1);
        IF v_parent_path IS NOT NULL AND v_parent_path != '' THEN
            v_path := v_parent_path || '/' || OLD.id;
        ELSE
            v_path := '/' || OLD.id;
        END IF;
    ELSE
        v_level := 0;
        v_path := '/' || OLD.id;
    END IF;

    UPDATE category_hierarchy
    SET parent_id = NEW.parent_placement_id,
        name_en = COALESCE(NEW.display_name_override, name_en),
        slug = COALESCE(NEW.slug, slug),
        display_order = COALESCE(NEW.display_order, display_order),
        is_leaf = COALESCE(NEW.is_leaf, is_leaf),
        is_active = COALESCE(NEW.is_active, is_active),
        level = v_level,
        path = v_path,
        updated_at = NOW()
    WHERE id = OLD.id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trg_classification_placements_delete()
RETURNS trigger AS $$
BEGIN
    UPDATE category_hierarchy SET is_active = false, is_archived = true, updated_at = NOW() WHERE id = OLD.id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_classification_placements_io_insert ON classification_placements;
CREATE TRIGGER trg_classification_placements_io_insert
INSTEAD OF INSERT ON classification_placements
FOR EACH ROW EXECUTE FUNCTION trg_classification_placements_insert();

DROP TRIGGER IF EXISTS trg_classification_placements_io_update ON classification_placements;
CREATE TRIGGER trg_classification_placements_io_update
INSTEAD OF UPDATE ON classification_placements
FOR EACH ROW EXECUTE FUNCTION trg_classification_placements_update();

DROP TRIGGER IF EXISTS trg_classification_placements_io_delete ON classification_placements;
CREATE TRIGGER trg_classification_placements_io_delete
INSTEAD OF DELETE ON classification_placements
FOR EACH ROW EXECUTE FUNCTION trg_classification_placements_delete();
