-- 060_classification_attribute_rules_columns.sql
-- Fix missing columns on classification_attribute_rules required by GetEffectiveAttributes and SaveAttributeRule

ALTER TABLE classification_attribute_rules 
  ADD COLUMN IF NOT EXISTS concept_id INT,
  ADD COLUMN IF NOT EXISTS is_inherited BOOLEAN DEFAULT false,
  ADD COLUMN IF NOT EXISTS override_validation_rules JSONB,
  ADD COLUMN IF NOT EXISTS override_unit VARCHAR(50);

CREATE UNIQUE INDEX IF NOT EXISTS uq_classification_attr_rules ON classification_attribute_rules (placement_id, attribute_id);
