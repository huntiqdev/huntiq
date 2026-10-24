# HUNTIQ Data Platform (`data-platform/`)

The offline data management, validation, labeling, and training preparation workspace for HUNTIQ.

---

## 1. Mission and Boundaries

- **Offline ML Lab:** Responsible for dataset contracts, privacy scrubbing, training pipelines, benchmark evaluations, and Model Cards.
- **Isolation:** Completely decoupled from the production inference serving runtime.
- **Data Invariance:** Never stores un-anonymized user data, plain passwords, or financial credentials in Git.
- **Canonical Taxonomy Alignment:** All category exports map deterministically to canonical IDs in PostgreSQL.

---

## 2. Directory Layout

```text
data-platform/
├── data-contracts/            # Pydantic & JSONSchema dataset and inference contracts
│   └── schemas/
├── taxonomy/                  # Versioned taxonomy snapshots and mappings
│   ├── schemas/
│   ├── manifests/
│   └── reports/
├── ml/                        # Machine learning workspace
│   ├── datasets/              # Versioned Parquet / JSONL files (gitignored)
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
└── data_quality/              # Quality gates, leakage detectors & validation rules
    └── validator.py
```

---

## 3. Dataset Lifecycle & Governance

Every dataset stored in `processed/` or `evaluation/` must have a corresponding `manifest.json` complying with `data-contracts/schemas/manifest_schema.json`.

```
Raw Export → Privacy Scrubbing → Validation Gate → Versioned Manifest → Deterministic Split → Model Training
```

### Mandatory Rules:
1. PII Redaction: Ethiopian phone numbers, email addresses, and names are scrubbed.
2. Leakage Prevention: Train and test splits must never share duplicate listings.
3. Provenance: Every record links back to an export timestamp and schema version.

---

## 4. Documentation

Platform-wide AI governance and testing specifications are in [`../docs/`](../docs/):
- [AI Data Security & PII Scrubbing Standards](../docs/security/ai-data-security.md)
- [AI Testing & Data Validation Gates](../docs/testing/ai-testing.md)
- [AI Architecture & Threat Model](../docs/architecture/ai-architecture.md)
