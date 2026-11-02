# HUNTIQ Dataset Versioning Specification

**Project:** HUNTIQ AI Data Platform  
**Standard:** Semantic Versioning for Data (Data-SemVer 2.0.0)

---

## 1. Version Format: `MAJOR.MINOR.PATCH`
- **MAJOR (Breaking Schema / Label Change):** Incremented when fields are added, renamed, or removed in `schema.json`, or when label definitions change.
- **MINOR (Batch Ingestion / Linguistic Expansion):** Incremented when new verified data batches are merged into `processed/` or when language splits expand without breaking schema contracts.
- **PATCH (Correction / Minor Typo Fix):** Incremented when existing record typos are repaired or metadata tags are updated.

## 2. Immutability Policy
Once a dataset version is published with a Git commit tag, training splits (`train.jsonl`, `validation.jsonl`, `test.jsonl`) are strictly immutable. Any modification requires a new version identifier.
