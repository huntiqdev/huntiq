# Data Ingestion Pipeline Workflow

---

## Steps:
1. Ingest raw batches into `ai/data_pipeline/import_jobs/`.
2. Compute SHA-256 hash and record provenance metadata.
3. Validate format against `source_schema.json`.
4. Route valid files to module `raw/` staging.
