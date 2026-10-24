# Model Card Template — HUNTIQ Machine Learning Governance

## 1. Model Details
- **Model Name:** [e.g., huntiq-query-intent-classifier]
- **Model Version:** [e.g., 1.0.0]
- **Task Type:** [e.g., Multi-Class Text Classification]
- **Artifact SHA-256 Checksum:** [e.g., e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
- **Artifact Format:** [e.g., ONNX / Joblib / SafeTensors]
- **Release Date:** [YYYY-MM-DD]
- **Developers / Authors:** [HUNTIQ ML Engineering Team]
- **License:** Proprietary (Internal HUNTIQ Use Only)

---

## 2. Intended Use
- **Primary Use Case:** [Detailed description of production inference task]
- **Primary Consumers:** HUNTIQ Go Backend (`backend/internal/service/search_service.go`)
- **Supported Languages:** English (`en`), Amharic (`am`), Afaan Oromo (`om`), Tigrinya (`ti`)
- **Out-of-Scope Uses:** Autonomous listing deletion, direct financial decision making, unreviewed bans.

---

## 3. Data & Provenance
- **Training Dataset Manifest:** `data-platform/ml/datasets/processed/manifest_v1.json`
- **Evaluation Dataset Manifest:** `data-platform/ml/datasets/evaluation/manifest_eval_v1.json`
- **Dataset Hash (SHA-256):** [sha256]
- **Record Count:** [Train: N, Validation: N, Test: N]
- **Data Split Methodology:** [Time-based cutoff / Stratified 80-10-10]
- **PII Scrubbing Verification:** Completed with `data_quality/validator.py` v1.0.0

---

## 4. Evaluation & Performance Benchmarks

### 4.1 Comparison against Non-AI Baseline
| Metric | Non-AI Keyword Baseline | Candidate Model | Relative Improvement |
| :--- | :--- | :--- | :--- |
| **Accuracy** | 62.4% | 88.7% | +26.3% |
| **Macro F1** | 0.58 | 0.86 | +0.28 |
| **Top-3 Accuracy** | 74.1% | 96.2% | +22.1% |
| **p95 Latency** | 2.1 ms | 14.5 ms | +12.4 ms (Within 500ms budget) |
| **Memory Footprint** | 5 MB | 45 MB | Acceptable |

---

## 5. Known Limitations & Failure Modes
- **Short Queries (1-2 characters):** Model confidence drops below 0.50. Go backend must fallback to direct prefix matching.
- **Transliterated Script:** Latin transliterations of Amharic/Tigrinya may yield lower confidence than native Ge'ez script.
- **Novel Brands:** Unseen brand names outside the training window default to general category classes.

---

## 6. Approvals & Deployment Sign-Off
- **ML Engineer Sign-Off:** [Name & Date]
- **Platform Architect Approval:** [Name & Date]
- **Security Lead Approval:** [Name & Date]
- **Deployment Status:** [Draft / Approved / Production / Deprecated]
