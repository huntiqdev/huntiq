# HUNTIQ Dataset Import Guide

**Target Audience:** Data Engineers, ML Engineers, and Annotation Leads.

---

## 1. Ingestion Workflow
1. **Source Registration:** Register provider, license, and checksum in `sources.json`.
2. **Raw Staging:** Place incoming JSONL, CSV, or Parquet files in `raw/external/` or `raw/original/`.
3. **Schema Ingestion:** Run `uv run python -m huntiq_ai.data_pipeline.ingest --module <module_id>`.
4. **Validation Check:** Verify that 100% of imported records conform to `schema.json`.
5. **Quarantine Handling:** Records failing validation are moved to `quarantine/` with an explicit failure code.
