# HUNTIQ Machine Learning Lifecycle & Training Framework

**Status:** Architecture Foundation & Phase 0 Definition  

---

## 1. Lifecycle Phases

The machine learning lifecycle for HUNTIQ consists of five distinct phases:

1. **Task Definition & Baseline (Phase 0 - Current):**
   - Formalize task (e.g. Query Intent Classification, Category Suggestion).
   - Define exact input/output contracts in `data-contracts/schemas/`.
   - Establish non-AI baseline heuristic (e.g., keyword dictionary matching).
2. **Data Curation & Anonymization (Phase 1):**
   - Extract sample from PostgreSQL via read-only script.
   - Run privacy scrubber (`data_quality/validator.py`).
   - Create deterministic 80/10/10 train/validation/test splits.
   - Generate signed `manifest.json`.
3. **Model Exploration & Benchmarking (Phase 2):**
   - Train lightweight CPU-friendly baseline (e.g., Logistic Regression / Naive Bayes / Random Forest with TF-IDF).
   - Evaluate against test split: measure Precision, Recall, Macro-F1, and Inference Latency (ms).
   - Require candidate model to beat the heuristic baseline before proceeding.
4. **Model Card & Approval (Phase 3):**
   - Document artifact checksum, dataset provenance, metrics, and known failure cases in `model_cards/`.
   - Security and Architecture review.
5. **Serving Integration (Phase 4):**
   - Package model artifact into `ai/models/`.
   - Expose behind internal authenticated route in `ai/src/huntiq_ai/api/routes/predict.py`.

---

## 2. Model Cards

Every deployed model artifact must have a Model Card matching `data-platform/ml/model_cards/MODEL_CARD_TEMPLATE.md`.
