# HUNTIQ AI Platform Testing and Quality Strategy

**Document Version:** 1.0.0  
**Target:** AI Engineers, Software Engineers, and QA  
**Status:** Approved Testing Strategy  

---

## 1. Testing Philosophy

The HUNTIQ AI Platform relies on **fast, deterministic, zero-external-dependency testing**:
- Unit and integration tests must run offline in under **5 seconds**.
- No production database, cloud GPU, or live third-party internet APIs are required to execute the test suite.
- All tests utilize synthetic fixtures.

---

## 2. Test Classification Matrix

| Layer | Directory | Purpose | Tools | Typical Execution Time |
| :--- | :--- | :--- | :--- | :--- |
| **Unit Tests** | `ai/tests/unit/` | Tests isolated modules: config parsing, error mapping, health response formatting, security token validation. | `pytest` | < 1 second |
| **Integration Tests** | `ai/tests/integration/` | Tests FastAPI application lifecycle, middleware, correlation ID tracking, and payload limits via ASGI client. | `pytest`, `httpx.AsyncClient` | < 2 seconds |
| **Contract Tests** | `ai/tests/contract/` | Verifies request and response schemas against Go client expectations and OpenAPI definitions. | `pytest`, `pydantic` | < 1 second |
| **Data Quality Tests** | `data-platform/tests/` | Verifies dataset manifests, detects schema drift, duplicate records, label inconsistencies, and PII leakage. | `pytest` | < 2 seconds |
| **Offline Evaluation** | `data-platform/ml/evaluation/` *(Future)* | Measures offline model metrics (Precision, Recall, F1, MRR, latency) on held-out test splits. | `scikit-learn` *(optional)* | Scheduled offline |

---

## 3. Running Tests Locally

### Run All AI Service Tests
```bash
cd ai
uv run pytest -v
```

### Run Tests with Coverage Report
```bash
cd ai
uv run pytest --cov=huntiq_ai --cov-report=term-missing
```

### Run Linting and Code Style Checks
```bash
cd ai
uv run ruff check .
uv run ruff format --check .
```

### Run Data Platform Contract Validation Tests
```bash
cd data-platform
uv run pytest tests/ -v
```

---

## 4. Test Invariants & Failure Criteria

A build or pull request must fail if:
1. Any unit or integration test fails.
2. Configuration permits insecure defaults in production mode (e.g. empty `HUNTIQ_AI_SECRET_KEY`).
3. Payload limit middleware fails to reject requests exceeding 64KB.
4. Health check exposes internal filesystem paths or secret keys.
5. Missing required translation or category ID in taxonomy manifest.
6. Data quality validator detects un-scrubbed Ethiopian phone numbers or email addresses in processed datasets.
