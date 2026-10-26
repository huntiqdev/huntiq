# HUNTIQ AI — Model Evaluation Strategy

**Document Version:** 1.0.0  
**Status:** Proposal / Planning Document  
**Related:** [MODEL_LIFECYCLE.md](./MODEL_LIFECYCLE.md) · [TRAINING_PIPELINE.md](../training/TRAINING_PIPELINE.md) · [ai-testing.md](../testing/ai-testing.md)

---

## 1. Purpose

This document defines the evaluation strategy for HUNTIQ AI models. It covers metrics, baseline comparisons, offline evaluation, error analysis, latency, resource usage, and deployment acceptance criteria.

> **Important:** No trained models, measured metrics, or production benchmarks exist at the time of writing. All thresholds in this document are **proposals** that must be established from measured baseline results before being used as acceptance gates.

---

## 2. Evaluation Principles

1. **Measure first, then set thresholds.** Do not invent acceptance criteria before measuring the baseline. Measure the PostgreSQL FTS baseline on the held-out test set, then set the acceptance threshold as a percentage improvement above it.
2. **Never evaluate on training data.** All reported metrics must come from the fixed held-out test set (`ml/datasets/evaluation/`).
3. **Offline evaluation precedes any production consideration.** No model is deployed without offline evaluation results documented in the model card.
4. **Latency is a first-class metric.** A highly accurate model that exceeds the 500ms Go-to-Python budget is not deployable.
5. **Human error analysis is required.** Automated metrics are necessary but not sufficient — a sample of errors must be reviewed manually before approval.

---

## 3. The Non-AI Baseline

The existing HUNTIQ search system is the baseline that any AI model must demonstrably outperform:

| Component | Implementation |
|---|---|
| Full-text search | `plainto_tsquery('simple', ...)` on `search_vector` (GIN index) |
| Synonym expansion | `search_synonyms` table (migration 019) |
| Query normalization | `normalize_query()` PostgreSQL function (migration 019) |
| Brand detection | `extract_brands()` PostgreSQL function (migration 019) |
| Result logging | `search_analytics` table (migration 019) |

**How to measure the baseline:**

1. Take the fixed held-out test set from `ml/datasets/evaluation/`.
2. For each query in the test set, retrieve the top-k results from PostgreSQL using the existing search.
3. Compare the top result's category to the ground-truth intent/category label.
4. Record Accuracy@1, Accuracy@3, and p95 latency.

**The baseline results become the minimum bar.** A model is approved only if it exceeds this bar by a meaningful margin (proposed: ≥10% relative improvement on Accuracy@1).

---

## 4. Use-Case-Specific Evaluation Metrics

### 4.1 Query Intent Classification

**Task:** Given a user search query and locale, classify it into one of the authorized intent categories defined in `data-contracts/schemas/query_intent_contract.json`.

**Authorized intent labels:**
`buy_product`, `rent_property`, `find_service`, `job_search`, `vehicle_search`, `price_check`, `ambiguous`

**Metrics:**

| Metric | Formula | Rationale |
|---|---|---|
| **Accuracy** | Correct / Total | Overall correctness on the test set |
| **Macro F1** | Mean F1 across all intent classes | Handles class imbalance; ensures minority intents are not ignored |
| **Precision per class** | TP / (TP + FP) | Measures how often predicted intent is correct |
| **Recall per class** | TP / (TP + FN) | Measures how often the true intent is captured |
| **Confidence distribution** | Histogram of `confidence` scores | Reveals threshold calibration quality |
| **Low-confidence rate** | Fraction with confidence < 0.70 | Proxy for real-world fallback frequency |

**Proposed initial acceptance thresholds (to be revised after baseline measurement):**

| Metric | Minimum to Exceed Baseline |
|---|---|
| Accuracy | Baseline + ≥10% relative |
| Macro F1 | ≥ 0.70 |
| p95 latency | ≤ 50ms (inference only, excluding network) |
| Low-confidence rate | ≤ 20% on test set |

---

### 4.2 Category Suggestion

**Task:** Given a listing title and description, suggest the top-k canonical category IDs from the HUNTIQ taxonomy.

**Metrics:**

| Metric | Formula | Rationale |
|---|---|---|
| **Accuracy@1** | Top-1 match / Total | Primary metric — first suggestion is correct |
| **Accuracy@3** | Top-3 match / Total | Seller experience — correct category in top 3 |
| **MRR** | Mean Reciprocal Rank | Measures average position of the first correct suggestion |
| **Confidence calibration** | Reliability diagram (confidence vs. accuracy) | Ensures confidence scores reflect true probabilities |
| **p95 latency** | 95th percentile inference time | Must fit within 500ms budget |

**Proposed initial acceptance thresholds (to be revised after baseline measurement):**

| Metric | Minimum to Exceed Baseline |
|---|---|
| Accuracy@1 | Baseline + ≥15% relative |
| Accuracy@3 | ≥ 85% |
| p95 latency | ≤ 80ms (inference only) |

---

## 5. Offline Evaluation Procedure

1. **Load the fixed test set** from `ml/datasets/evaluation/`.
2. **Run inference** using the candidate model for every test record.
3. **Compute metrics** using `sklearn.metrics` or equivalent.
4. **Record results** in the model card (`ml/model_cards/<name>.md`).
5. **Compare against baseline** — document the relative improvement.
6. **Human error analysis** — manually review a random sample of 50 misclassified examples. Categorize error types (e.g., ambiguous query, unseen brand, transliteration).

---

## 6. Regression Testing

Before any new model version is deployed, it must be tested against the **exact same fixed test set** used to evaluate the previous version. This ensures:

- No accuracy regression on previously handled cases.
- Performance improvements are real, not artifacts of test set changes.

> **The test set must never be updated once it is created.** If the taxonomy changes and new categories need evaluation, create a new held-out set and document it as a separate version.

---

## 7. Latency and Resource Evaluation

Latency is measured end-to-end from the Go backend perspective:

| Measurement Point | Budget |
|---|---|
| Go → AI service (full round trip) | ≤ 500ms (enforced by Go timeout) |
| AI service inference only (model compute) | ≤ 50–100ms (leaves margin for I/O) |
| AI service cold start | ≤ 5s (acceptable for service restart) |

**Resource limits (proposed for container deployment):**

| Resource | Proposed Limit |
|---|---|
| Memory | ≤ 512MB RSS at idle; ≤ 1GB at peak inference |
| CPU | ≤ 1 vCPU for scikit-learn models; 2 vCPU for larger models |
| Disk (model artifact) | ≤ 200MB |

These limits are proposals. Actual limits must be established from measured resource usage during load testing.

---

## 8. Deployment Acceptance Criteria

A model version is approved for deployment only when all of the following are confirmed in the model card:

| Criterion | Required Evidence |
|---|---|
| Offline metrics exceed baseline | Measured results documented in model card |
| No regression on previous test set | Documented comparison in model card |
| Latency within budget | Measured under realistic load |
| PII validation passed | `DataQualityValidator` report with `valid: true` |
| Human error analysis completed | Error sample reviewed and categorized |
| Security lead sign-off | Name and date in model card |
| Platform architect sign-off | Name and date in model card |

---

## 9. Open Questions

These questions must be resolved through measurement before the thresholds above are finalized:

1. What is the current baseline Accuracy@1 for category suggestion using PostgreSQL FTS on the held-out test set?
2. What fraction of real user queries in `search_analytics` have `result_count = 0` (zero-result queries)?
3. What is the distribution of query lengths across languages (en, am, om, ti)?
4. What fraction of queries contain Ge'ez script vs. Latin transliteration?

These measurements will determine whether TF-IDF + Logistic Regression is sufficient or whether a more complex model is justified.
