# HUNTIQ AI — Dataset Lifecycle

**Document Version:** 1.0.0  
**Status:** Proposal / Planning Document  
**Related:** [TRAINING_PIPELINE.md](./TRAINING_PIPELINE.md) · [ai-data-security.md](../security/ai-data-security.md)

---

## 1. Purpose

This document describes how datasets are created, governed, versioned, and retired for HUNTIQ AI training. It covers data provenance, authorization, privacy, deduplication, labeling, quality checks, and leakage prevention.

---

## 2. Implementation Status

| Component | Status | Evidence |
|---|---|---|
| Data contract schemas (JSONSchema) | ✅ Implemented | `data-contracts/schemas/query_intent_contract.json`, `category_suggestion_contract.json`, `manifest_schema.json` |
| Dataset manifest schema | ✅ Implemented | `data-contracts/schemas/manifest_schema.json` |
| PII scanner | ✅ Implemented | `data_quality/validator.py` — `scan_for_pii()` |
| Duplicate detector | ✅ Implemented | `data_quality/validator.py` — `detect_duplicates()` |
| Leakage detector | ✅ Implemented | `data_quality/validator.py` — `detect_leakage()` |
| Manifest validator | ✅ Implemented | `data_quality/validator.py` — `validate_manifest()` |
| Synthetic fixtures directory | ✅ Scaffolded | `ml/datasets/synthetic_fixtures/` (empty) |
| Dataset directory structure | ✅ Scaffolded | `raw/`, `interim/`, `processed/`, `evaluation/` (all empty) |
| Extraction/cleaning pipelines | ❌ Not implemented | `ml/pipelines/` does not exist |
| Labeling guidelines | ❌ Not implemented | `ml/labeling/` does not exist |
| Actual collected datasets | ❌ Do not exist | All dataset directories are empty |
| Version registry | ❌ Not implemented | No experiment tracker configured |

---

## 3. Data Authorization

All training data must be derived from sources explicitly authorized for this purpose:

| Source | Authorization Level | Authorized Fields |
|---|---|---|
| `products` table | Database export by data custodian | `title`, `description`, `category_id`, `subcategory`, `listing_type`, `condition` |
| `search_analytics` table (migration 019) | Database export by data custodian | `normalized_query`, `result_count`, `category_id` |
| `categories` / `category_hierarchy` tables | Public API (`GET /api/categories/tree`) | All taxonomy fields |
| Synthetic fixtures | Internal authorship — no export needed | Hand-crafted examples only |

**Prohibited:**
- `users` table (any column)
- `otp_codes`, `transactions`, `referral_commissions`, `upgrade_requests`
- Any column containing phone numbers, email addresses, passwords, payment references

---

## 4. Dataset Directory Structure and Stages

```
ml/datasets/
├── raw/                    # Unmodified authorized export from Go backend
├── interim/                # After PII scrubbing — before feature engineering
├── processed/              # Feature-engineered, split, validated — ready for training
├── evaluation/             # Fixed held-out test set — never used in training
└── synthetic_fixtures/     # Hand-crafted examples to bootstrap labeling
```

### Rules per stage:

| Stage | Rule |
|---|---|
| `raw/` | Never committed to Git. PII may still be present. No processing yet. |
| `interim/` | PII scrubbed and verified by validator. Still not split. Not committed to Git. |
| `processed/` | Split applied, manifest created, validator passed. Not committed to Git. |
| `evaluation/` | Fixed test set. Created once. Never updated after baseline measurement. Not committed to Git. |
| `synthetic_fixtures/` | Hand-crafted, no PII. May be committed to Git as examples only. |

> **All non-synthetic datasets are gitignored.** Only `manifest.json` files and `synthetic_fixtures/` may be committed.

---

## 5. Data Provenance

Every dataset in `processed/` or `evaluation/` must have a `manifest.json` conforming to `data-contracts/schemas/manifest_schema.json`.

Required manifest fields:

| Field | Description |
|---|---|
| `manifest_id` | Unique identifier (UUID) |
| `dataset_name` | Human-readable name (e.g., `huntiq_query_intent_v1`) |
| `version` | Semantic version (`1.0.0`) |
| `created_at` | ISO 8601 timestamp |
| `source_export_id` | Identifies the database export or script commit hash |
| `record_count` | Total records across all splits |
| `sha256_checksum` | SHA-256 of the processed JSONL file |
| `pii_scrubbed` | `true` — must be verified by validator |
| `splits` | Object with `train`, `val`, `test` counts summing to `record_count` |

---

## 6. Privacy and Sanitization

Scrubbing must be applied before data leaves `raw/` status. The `DataQualityValidator.scan_for_pii()` method in `data_quality/validator.py` detects:

- **Ethiopian phone numbers** — regex: `(?:(?:\+251\s?|0)[79]\d{2}[\s.-]?\d{3}[\s.-]?\d{3}|(?:\+251\s?|0)[79]\d{8})`
- **Email addresses** — standard RFC 5322 pattern
- **Credential keywords** — `password`, `secret`, `api_key`, `bearer <token>`

**Action on detection:** Replace with `[PHONE_REDACTED]`, `[EMAIL_REDACTED]`, `[SECRET_REDACTED]`. Log count replaced without revealing values.

---

## 7. Split Assignment and Leakage Prevention

**Required:** Use a **time-based cutoff**, not random shuffling, to assign `train`/`val`/`test` splits.

Rationale: Random splits on temporal data cause data leakage through autocorrelation — later queries tend to resemble earlier ones for trending items. A time-based cutoff simulates real deployment conditions.

Example cutoff policy:
- `train`: records created before `T - 60 days`
- `val`: records from `T - 60 days` to `T - 14 days`
- `test`: records from `T - 14 days` to `T`

The `DataQualityValidator.detect_leakage()` method verifies no query appears in both `train` and `val`/`test` after splitting.

---

## 8. Labeling

**Currently unimplemented.** Required before any supervised model can be trained:

- Label the `intent` field in `query_intent_contract.json` for collected queries.
- A labeling guideline document must be created at `ml/labeling/labeling_guidelines.md` (not yet written).
- Labels must be reviewed by at least one additional team member before entering `processed/`.
- Inter-annotator agreement should be measured and recorded in the dataset manifest.

**Synthetic fixtures** (`ml/datasets/synthetic_fixtures/`) may be used as seed examples to guide human labelers on ambiguous cases. These are hand-crafted and do not represent real user data.

---

## 9. Deduplication

The `DataQualityValidator.detect_duplicates()` method identifies identical query strings within the same dataset. Duplicates must be removed before the dataset enters `processed/`.

Near-duplicate detection (e.g., "samsung a16" vs "Samsung A16 phone") is not yet implemented. This is a future improvement.

---

## 10. Versioning and Reproducibility

Each dataset version is identified by:
1. The manifest `version` field (semantic version).
2. The `source_export_id` (identifies the exact database export script and commit hash).
3. The `sha256_checksum` of the processed file.

A new version must be created when:
- New records are added.
- Labels are corrected.
- The scrubbing or normalization pipeline changes.
- The split cutoff date changes.

Experiments must record which manifest version they used. This allows exact reproduction of any training run.

---

## 11. Retention Policy

- **`raw/` exports:** Retained for 30 days after `interim/` validation is complete, then deleted.
- **`interim/` files:** Retained until `processed/` version is validated and approved.
- **`processed/` and `evaluation/` datasets:** Retained for the lifetime of any model trained on them, plus 90 days after model deprecation.
- **`synthetic_fixtures/`:** Permanent (committed to Git, no PII).
- **Model artifacts:** See [MODEL_LIFECYCLE.md](../models/MODEL_LIFECYCLE.md).

---

## 12. How Existing Components Fit

| Component | Role in Dataset Lifecycle |
|---|---|
| `data-contracts/schemas/query_intent_contract.json` | Defines the schema each training record must conform to |
| `data-contracts/schemas/manifest_schema.json` | Defines what a valid dataset manifest looks like |
| `data_quality/validator.py` | Enforces all quality gates at Stages 3, 4, and 7 |
| `ml/datasets/synthetic_fixtures/` | Seed examples for labeling guidance |
| `taxonomy/schemas/taxonomy_export.schema.json` | Defines the schema for exporting the canonical category tree |
| `backend/internal/database/migrations/019_query_understanding.sql` | Source of `search_analytics` and `search_synonyms` tables used for training data extraction |
