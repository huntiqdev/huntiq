# Schema Validation Workflow

---

## Steps:
1. Execute `jsonschema.validate()` against module `schema.json`.
2. Check type constraints, string lengths, and required fields.
3. Move failed records to `quarantine/` with rejection reason.
