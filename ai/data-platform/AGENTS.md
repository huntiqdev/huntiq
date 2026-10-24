# HUNTIQ Data Platform (`ai/data-platform/`) Instructions for Agents & Developers

**Directory:** `/ai/data-platform`  
**Purpose:** Dataset Governance, Schemas, Manifests, Privacy Scrubbing, and Offline Model Evaluation  

---

## 1. Directory Structure

```text
ai/data-platform/
├── AGENTS.md                  # This file
├── README.md                  # Human developer documentation
├── data-contracts/
│   ├── README.md
│   └── schemas/
│       ├── manifest_schema.json
│       ├── query_intent_contract.json
│       └── category_suggestion_contract.json
├── taxonomy/
│   ├── README.md
│   ├── schemas/               # Canonical taxonomy export schemas
│   └── manifests/             # Versioned taxonomy snapshots
├── ml/
│   ├── README.md
│   ├── datasets/              # Strictly gitignored (except .gitkeep)
│   │   ├── raw/               # Raw exported snapshots
│   │   ├── interim/           # Scrubbed & anonymized data
│   │   ├── processed/         # Feature-engineered training splits
│   │   └── evaluation/        # Fixed held-out evaluation sets
│   ├── pipelines/             # Deterministic data extraction & cleaning scripts
│   ├── labeling/              # Human annotation guidelines & reviewed labels
│   ├── training/              # Offline model training code
│   ├── evaluation/            # Offline benchmarking & error analysis
│   ├── models/                # Trained artifacts (gitignored)
│   └── model_cards/           # Markdown model governance cards
├── data_quality/
│   ├── README.md
│   └── validator.py           # Integrity, leakage, and schema drift checks
└── tests/
    └── test_data_validation.py # Automated data contract verification
```

---

## 2. Dataset Governance Rules

1. **Never Commit Raw Data to Git:** The `.gitignore` at the repository root explicitly ignores `data-platform/ml/datasets/` and `data-platform/ml/models/`. Never bypass this rule.
2. **Mandatory Dataset Manifests:** Every dataset in `processed/` or `evaluation/` must have an associated `manifest.json` conforming to `data-contracts/schemas/manifest_schema.json`.
3. **Data Provenance:** Every manifest must specify:
   - Source table or export origin.
   - Exact query or extraction script commit hash.
   - Number of records.
   - SHA-256 checksum of the dataset file.
   - Anonymization and privacy filter version applied.
4. **No Data Leakage:**
   - Training and evaluation splits must never share identical or near-duplicate listings.
   - Split assignment must be deterministic (fixed random seed or time-based cutoff).
5. **Privacy Scrubbing:**
   - All phone numbers, email addresses, person names, and secret tokens must be scrubbed before data enters `interim/`.

---

## 3. Verification Commands

Run dataset contract and quality validation tests:
```bash
cd ai/data-platform
uv run pytest tests/ -v
```
