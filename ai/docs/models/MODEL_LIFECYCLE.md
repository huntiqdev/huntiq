# HUNTIQ AI — Model Lifecycle

**Document Version:** 1.0.0  
**Status:** Proposal / Planning Document  
**Related:** [MODEL_EVALUATION.md](./MODEL_EVALUATION.md) · [TRAINING_PIPELINE.md](../training/TRAINING_PIPELINE.md) · [ai-architecture.md](../architecture/ai-architecture.md)

---

## 1. Purpose

This document describes the full lifecycle of a HUNTIQ AI model from development through deprecation. It covers experiment tracking, artifact versioning, model cards, evaluation approval, deployment, rollback, monitoring, and retraining.

---

## 2. Implementation Status

| Component | Status | Evidence |
|---|---|---|
| Model card template | ✅ Implemented | `ml/model_cards/MODEL_CARD_TEMPLATE.md` |
| Model artifact directory | ✅ Scaffolded | `ml/models/.gitkeep` — empty, gitignored |
| FastAPI health endpoint (`models_loaded`) | ✅ Implemented | `schemas/health.py` — `models_loaded: false` field in readiness response |
| Experiment tracking system | ❌ Not implemented | No tool configured |
| Trained model artifacts | ❌ Do not exist | `ml/models/` is empty |
| Prediction endpoint | ❌ Not implemented | `api/routes/predict.py` does not exist |
| Deployment pipeline | ❌ Not implemented | No CI/CD stage configured for model deployment |
| Production monitoring | ❌ Not implemented | No dashboard or alerting configured |

---

## 3. Lifecycle Stages

```
Development → Experiment → Evaluation → Approval → Packaging → Deployment
     ↑                                                              ↓
     └──────────────── Monitoring → Retraining ←───────────────────┘
```

---

### Stage 1: Development

- Training code is written in `ml/training/` (not yet created).
- A candidate model is selected based on the use case requirements in [TRAINING_PIPELINE.md](../training/TRAINING_PIPELINE.md).
- A draft model card is created from `MODEL_CARD_TEMPLATE.md` before training begins.

---

### Stage 2: Experiment Tracking

**Currently unimplemented.** When implemented:

- Each training run is assigned a unique `run_id`.
- Recorded fields per run: `run_id`, `model_type`, `dataset_manifest_version`, `hyperparameters`, `train_metrics`, `val_metrics`, `test_metrics`, `artifact_path`, `artifact_sha256`, `duration_seconds`, `timestamp`.
- Experiments are stored in a lightweight JSONL log at `ml/experiments.jsonl` (gitignored).
- The best run per use case is promoted to a named artifact in `ml/models/`.

> **Note:** No experiment tracking tool (MLflow, W&B, etc.) is configured in the current repository. This is a deliberate decision to avoid infrastructure dependencies in the foundation phase. A lightweight JSONL log is sufficient for the first model.

---

### Stage 3: Artifact Versioning

Each promoted model artifact is named with the pattern:

```
ml/models/<use_case>_v<major>.<minor>.<patch>.<format>
```

Example:
```
ml/models/query_intent_v1.0.0.joblib
ml/models/query_intent_v1.0.0.sha256
```

The `.sha256` file contains the hex-encoded SHA-256 hash of the artifact. The production service verifies this hash at startup before loading.

> **Artifact format preference:** ONNX > Joblib (with `trust=False` or explicit allowlist) > other formats. Avoid arbitrary pickle files.

---

### Stage 4: Model Cards

Every artifact promoted to `ml/models/` must have a completed model card at:

```
ml/model_cards/<use_case>_v<version>.md
```

The card uses `MODEL_CARD_TEMPLATE.md` and must include:

- Measured evaluation metrics (not estimated — see [MODEL_EVALUATION.md](./MODEL_EVALUATION.md))
- Dataset manifest versions used for training and evaluation
- Known limitations and failure modes
- SHA-256 checksum of the artifact
- Sign-off from ML engineer, platform architect, and security lead

A model without a completed and signed card must not be deployed.

---

### Stage 5: Evaluation Approval

Before a model can be promoted to deployment:

1. Evaluation metrics must be measured on the **fixed held-out test set** — not on training or validation data.
2. Metrics must exceed the PostgreSQL FTS baseline (see [MODEL_EVALUATION.md](./MODEL_EVALUATION.md)).
3. Latency must be within the 500ms Go-to-Python budget with margin.
4. The model card must be reviewed and signed off.

> **Approval is a manual gate.** No automated promotion to production without explicit sign-off.

---

### Stage 6: Packaging

**Proposed packaging process (not yet implemented):**

1. Copy the approved artifact to the Docker build context or deployment volume.
2. Update `MODEL_DIR` path in the FastAPI configuration.
3. Verify the SHA-256 hash at container startup using `core/config.py`.
4. The readiness endpoint (`/health/ready`) transitions `models_loaded` from `false` to `true` only when the artifact is successfully loaded and verified.

---

### Stage 7: Deployment

**Currently unimplemented.** When implemented:

- The AI service is deployed as a sidecar or internal service alongside the Go backend.
- The Go backend `AI_SERVICE_URL` config points to the service.
- The Go backend calls `/health/ready` before making prediction requests.
- A circuit breaker or timeout (500ms) is enforced on every prediction call.
- Deployment does not replace the PostgreSQL FTS fallback — it augments it.

---

### Stage 8: Rollback

A rollback is required when:
- Production confidence scores fall significantly below the development test-set level.
- Error rate on prediction calls increases.
- A data quality or privacy incident is discovered post-deployment.

**Rollback procedure (proposed):**
1. Stop routing to the AI endpoint (disable `AI_SERVICE_URL` or set `AI_ENABLED=false` in Go config).
2. The Go backend automatically falls back to PostgreSQL FTS — no user impact.
3. Identify and resolve the issue.
4. Re-deploy the previous approved artifact.

---

### Stage 9: Monitoring

**Currently unimplemented.** When implemented, monitor:

| Metric | Alert Threshold | Action |
|---|---|---|
| Fraction of predictions below confidence 0.70 | > 15% over 24h | Investigate distribution shift; schedule retraining |
| p95 prediction latency | > 400ms | Review model complexity; optimize or rollback |
| AI service error rate | > 2% | Check service health; consider rollback |
| `models_loaded: false` in readiness | Any production instance | Alert immediately; restart or rollback |

---

### Stage 10: Retraining

Retraining is triggered when:
- Monitoring alerts indicate confidence degradation.
- A significant new product category is added to the taxonomy.
- A periodic review (proposed: every 90 days) shows test-set accuracy falling below the acceptance threshold.

**Retraining is never automatic.** Each retraining run follows the full pipeline from Stage 3 in [TRAINING_PIPELINE.md](../training/TRAINING_PIPELINE.md).

---

### Stage 11: Deprecation

A model is deprecated when:
- A newer version has been deployed and validated.
- The use case has been retired.

On deprecation:
- The model card is updated with `Deployment Status: Deprecated` and the deprecation date.
- The artifact is retained for 90 days, then deleted.
- The experiment log entry is retained indefinitely for audit purposes.
