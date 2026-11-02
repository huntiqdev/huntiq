# HUNTIQ Dataset Lifecycle Management

**Project:** HUNTIQ Universal Marketplace AI  
**Scope:** End-to-end dataset maturity stages from raw ingestion to model retirement.

---

## 1. Lifecycle Stages
Every dataset in the HUNTIQ ecosystem progresses through 6 deterministic stages:

1. **Stage 0: Initialized (Schema Defined)**
   - Schema and metadata manifests created.
   - Raw ingestion directories prepared. No training records permitted.
2. **Stage 1: Raw Ingested (Staging)**
   - Batches imported into `raw/` with verified provenance hashes and source licensing.
3. **Stage 2: Normalized & Deduplicated**
   - Text standardized through Ge'ez homophone unification, Oromo Qubee normalizer, and Latin transliteration cleaner.
   - Exact and semantic duplicate records quarantined.
4. **Stage 3: Annotated & Human Verified**
   - High-leverage samples reviewed by native Amharic, Oromo, Tigrinya, and English linguists.
   - Inter-annotator agreement score (Cohen's Kappa $\ge 0.85$) required.
5. **Stage 4: Partitioned & Leakage Audited**
   - Deterministic 80/10/10 split into `training/train.jsonl`, `validation.jsonl`, and `test.jsonl`.
   - Entity and session overlap tests executed.
6. **Stage 5: Production Approved (Active)**
   - Dataset tagged with immutable SemVer and registered in `dataset_registry.json`.
