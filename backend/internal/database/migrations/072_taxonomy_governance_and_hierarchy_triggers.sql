-- ============================================================================
-- 072_taxonomy_governance_and_hierarchy_triggers.sql
-- Materialized tree path cascades & hierarchy triggers
-- ============================================================================

CREATE OR REPLACE FUNCTION public.update_location_hierarchy_path()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
    parent_path VARCHAR(500);
    parent_level INT;
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.level := 0;
        NEW.path := '/' || NEW.id || '/';
    ELSE
        SELECT path, level INTO parent_path, parent_level
        FROM location_hierarchy
        WHERE id = NEW.parent_id;
        
        IF parent_path IS NULL THEN
            NEW.level := 1;
            NEW.path := '/' || NEW.parent_id || '/' || NEW.id || '/';
        ELSE
            NEW.level := parent_level + 1;
            NEW.path := parent_path || NEW.id || '/';
        END IF;
    END IF;
    
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.update_category_hierarchy_path()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
    v_parent_path VARCHAR;
    v_parent_level INT;
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.level := 0;
        NEW.path := '/' || CAST(NEW.id AS VARCHAR) || '/';
    ELSE
        SELECT level, path
        INTO v_parent_level, v_parent_path
        FROM category_hierarchy
        WHERE id = NEW.parent_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Parent category ID % does not exist', NEW.parent_id;
        END IF;

        NEW.level := v_parent_level + 1;
        NEW.path := v_parent_path || CAST(NEW.id AS VARCHAR) || '/';
    END IF;

    NEW.updated_at := NOW();
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION public.cascade_category_hierarchy_path()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    IF OLD.path IS DISTINCT FROM NEW.path THEN
        UPDATE category_hierarchy
        SET path = NEW.path || SUBSTRING(path FROM LENGTH(OLD.path) + 1),
            level = array_length(string_to_array(trim(both '/' from NEW.path || SUBSTRING(path FROM LENGTH(OLD.path) + 1)), '/'), 1) - 1,
            updated_at = NOW()
        WHERE path LIKE OLD.path || '%';
    END IF;
    RETURN NULL;
END;
$function$;

DROP TRIGGER IF EXISTS trg_location_path ON location_hierarchy;
CREATE TRIGGER trg_location_path BEFORE INSERT OR UPDATE OF parent_id, id ON location_hierarchy FOR EACH ROW EXECUTE FUNCTION update_location_hierarchy_path();

DROP TRIGGER IF EXISTS trg_update_category_hierarchy_path ON category_hierarchy;
CREATE TRIGGER trg_update_category_hierarchy_path BEFORE INSERT OR UPDATE OF parent_id, id ON category_hierarchy FOR EACH ROW EXECUTE FUNCTION update_category_hierarchy_path();

DROP TRIGGER IF EXISTS trg_cascade_category_hierarchy_path ON category_hierarchy;
CREATE TRIGGER trg_cascade_category_hierarchy_path AFTER UPDATE OF path ON category_hierarchy FOR EACH ROW EXECUTE FUNCTION cascade_category_hierarchy_path();
