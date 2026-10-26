# HUNTIQ AI — Training Pipeline

**Document Version:** 1.0.0  
**Status:** Proposal / Planning Document  
**Related:** [DATASET_LIFECYCLE.md](./DATASET_LIFECYCLE.md) · [ai-inference-architecture.md](../architecture/ai-inference-architecture.md) · [MODEL_LIFECYCLE.md](../models/MODEL_LIFECYCLE.md)

---

## 1. Purpose

This document describes the end-to-end workflow for developing, validating, and deploying a task-specific AI model for the HUNTIQ marketplace. It distinguishes between stages that are **currently implemented**, stages that are **partially implemented** (infrastructure exists, logic not yet written), and stages that are **planned proposals**.

No trained models, production deployments, or measured metrics exist at the time of writing. All metric thresholds are placeholders until measured from actual data.

---

## 2. Implementation Status Summary

| Stage | Status | Evidence |
|---|---|---|
| Data contracts and schemas | ✅ Implemented | `data-contracts/schemas/*.json` — JSONSchema contracts for `query_intent`, `category_suggestion`, `manifest` |
| Data quality validator | ✅ Implemented | `data_quality/validator.py` — PII scanning, duplicate detection, leakage detection, manifest validation |
| Data quality tests | ✅ Implemented | `data-platform/tests/` — pytest suite covering validator, manifest schema |
| Dataset directories | ✅ Scaffolded | `ml/datasets/{raw,interim,processed,evaluation,synthetic_fixtures}/` |
| Model card template | ✅ Implemented | `ml/model_cards/MODEL_CARD_TEMPLATE.md` |
| Taxonomy schema | ✅ Implemented | `taxonomy/schemas/taxonomy_export.schema.json` |
| Baseline (PostgreSQL FTS) | ✅ Implemented | `backend/internal/repository/search_repo.go`, migration 019 |
| FastAPI inference scaffold | ✅ Implemented | `src/huntiq_ai/` — health, security, config, logging, errors |
| Training pipelines | ❌ Not yet implemented | `ml/pipelines/` directory does not yet exist |
| Actual training datasets | ❌ Not yet collected | `ml/datasets/raw/`, `processed/`, `evaluation/` are empty |
| Model training code | ❌ Not yet implemented | `ml/training/` directory does not yet exist |
| Prediction endpoints | ❌ Not yet implemented | `api/routes/` has no `predict.py` |
| Experiment tracking | ❌ Not yet implemented | No tracking tool configured |
| Production model artifacts | ❌ Do not exist | `ml/models/` contains only `.gitkeep` |

---

## 3. End-to-End Training Workflow

### Stage 1 — Identify the AI Use Case

Before writing any code or collecting any data, define:

- **What problem does the model solve?** Example: "Given a marketplace search query in any of the four supported languages (en, am, om, ti), classify it into one of 7 intent categories so the Go backend can apply the correct search strategy."
- **What is the smallest useful version?** Start with the highest-frequency, easiest-to-label intent (e.g., `buy_product` vs. `rent_property`) before expanding to all 7 classes.
- **What is the acceptance criteria?** The model must outperform the deterministic PostgreSQL baseline on the held-out test set. See [MODEL_EVALUATION.md](../models/MODEL_EVALUATION.md).
- **What is the fallback?** The existing PostgreSQL full-text search (`search_repo.go`) continues to serve results if the model is absent or confidence is below threshold.

**Authorized use cases (from existing data contracts):**

| Use Case | Contract File | Status |
|---|---|---|
| Query intent classification | `query_intent_contract.json` | Contract defined, model not yet built |
| Category suggestion from title+description | `category_suggestion_contract.json` | Contract defined, model not yet built |

---

### Stage 2 — Obtain Authorized Training Data

**Requirement:** All training data must be obtained through authorized application interfaces. The Python AI service has no direct database write credentials and does not query PostgreSQL directly.

**Authorized data sources for HUNTIQ:**

1. **Search query logs** — Available via `search_analytics` table in PostgreSQL (migration 019). Access only through a controlled export approved by the data custodian.
2. **Listing titles and descriptions** — Available via read-only database export of `products` table (approved fields only: `title`, `description`, `category_id`, `subcategory`).
3. **Canonical taxonomy** — Available via `GET /api/categories/tree` endpoint (public, no auth required). The taxonomy export schema is defined at `taxonomy/schemas/taxonomy_export.schema.json`.
4. **Synthetic fixtures** — Hand-crafted examples stored in `ml/datasets/synthetic_fixtures/` to bootstrap labeling before real data is available.

**Prohibited data sources:**

- Raw phone numbers, email addresses, passwords, payment references — see [ai-data-security.md](../security/ai-data-security.md).
- Unreviewed scrapes from external websites.
- Any data not covered by HUNTIQ's data retention policy.

---

### Stage 3 — Validate Dataset Contracts and Data Quality

Before any further processing, run the data quality validator against the raw export:

```bash
cd ai/data-platform
uv run pytest tests/ -v
```

The `DataQualityValidator` in `data_quality/validator.py` enforces:

- **PII scanning:** Ethiopian phone patterns (`+2519...`, `09...`), email addresses, credential keywords.
- **Manifest validation:** All required fields present, `pii_scrubbed: true`, split sums match `record_count`.
- **Duplicate detection:** No identical query strings in the same dataset.
- **Leakage detection:** No query strings shared between `train` and `val`/`test` splits.

A dataset that fails any of these checks must not proceed to Stage 4.

---

### Stage 4 — Clean, Normalize, and Prepare Datasets

**Pipeline (proposed — not yet implemented in `ml/pipelines/`):**

1. **Scrub PII** — Apply regex substitutions from `data_quality/validator.py` to replace phone numbers with `[PHONE_REDACTED]`, emails with `[EMAIL_REDACTED]`.
2. **Normalize text** — Lowercase, strip leading/trailing whitespace, normalize Unicode (NFC form).
3. **Amharic/Tigrinya handling** — Preserve Ge'ez Unicode codepoints; do not apply Latin normalization.
4. **Deduplicate** — Remove exact duplicate `(query, intent)` pairs. Log the count removed.
5. **Assign split** — Use a **time-based cutoff** (not random) to assign `train`/`val`/`test` splits. This prevents data leakage from temporal autocorrelation.

**Output:** Save to `ml/datasets/processed/` as JSONL. Create a manifest file conforming to `data-contracts/schemas/manifest_schema.json`.

---

### Stage 5 — Establish a Non-AI Baseline

Before training any model, measure the performance of the existing deterministic system:

- **Current baseline:** PostgreSQL full-text search (`plainto_tsquery`) with `pg_trgm` extension, synonym expansion via `search_synonyms` table (migration 019).
- **Measure on the held-out test set:** Run the same test queries through the PostgreSQL search and record precision, recall, MRR, and latency.
- **Document in the model card:** The baseline result is the minimum bar a model must exceed before deployment approval.

> ⚠️ **No model should be deployed if it does not measurably outperform this baseline on the held-out test set.**

---

### Stage 6 — Train and Tune a Model

**Principles:**

- Start with the simplest model that could work (e.g., TF-IDF + Logistic Regression for intent classification).
- Add complexity only when a simpler approach is provably insufficient.
- All training code goes in `ml/training/` (not yet created).
- ML dependencies are in the optional group `[project.optional-dependencies] ml` in `pyproject.toml`.

**Candidate starting approaches by use case:**

| Use Case | Recommended First Approach | Rationale |
|---|---|---|
| Query intent (7 classes) | TF-IDF + Logistic Regression | Fast, interpretable, low memory. Good baseline for short text. |
| Category suggestion | FastText or TF-IDF + SVM | Handles multilingual short text well without GPU. |

---

### Stage 7 — Evaluate Against Predefined Metrics

Evaluation must be performed on the **fixed held-out test set only** — never on training or validation data.

See [MODEL_EVALUATION.md](../models/MODEL_EVALUATION.md) for the full evaluation strategy and metric definitions.

Key rules:
- Metrics must be measured, not estimated.
- The model must exceed the PostgreSQL baseline on the test set.
- Latency must be measured under realistic load conditions (not just single-request timing).

---

### Stage 8 — Record Experiments and Model Versions

**Currently unimplemented.** Required before production deployment:

- A lightweight experiment log (CSV or JSONL) recording: run ID, dataset manifest hash, hyperparameters, train/val/test metrics, and artifact path.
- Each candidate artifact stored in `ml/models/` with its SHA-256 checksum recorded in the model card.

---

### Stage 9 — Approve, Package, and Deploy

A model artifact may only be deployed when all of the following are true:

1. ✅ Model card is complete (`ml/model_cards/<model_name>.md`).
2. ✅ Evaluation metrics exceed the baseline on the held-out test set.
3. ✅ Security lead has verified no PII leaks in training data.
4. ✅ Platform architect has approved the inference endpoint schema.
5. ✅ SHA-256 checksum recorded and verified.
6. ✅ The FastAPI prediction endpoint (`api/routes/predict.py`) is implemented with authentication and fallback contract.

**Packaging:** The model artifact is embedded in the Docker image or mounted at the configured `MODEL_DIR` path. The service verifies the SHA-256 checksum at startup.

---

### Stage 10 — Monitor and Retrain

**Currently unimplemented.** Required once deployed:

- Track production confidence score distributions over time.
- Alert when the fraction of predictions below the 0.70 confidence threshold increases significantly (indicates distribution shift).
- Retrain when the test-set metrics on newly labeled data fall below the acceptance threshold.
- Never retrain automatically — human review required before a new model version is approved.

---

## 4. Engineering Principles

1. **Smallest useful capability first** — start with intent classification before category suggestion.
2. **Simple baselines before complex models** — TF-IDF + Logistic Regression must fail demonstrably before adding neural models.
3. **Deterministic business rules stay in Go** — the model produces a suggestion; Go decides what to do with it.
4. **Reproducibility** — fixed random seeds, time-based splits, versioned manifests.
5. **No GPU dependency in the foundation** — scikit-learn models must work on CPU.
6. **Fallback always available** — PostgreSQL FTS works independently of the AI service.
