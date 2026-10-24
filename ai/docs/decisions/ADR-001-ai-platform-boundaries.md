# ADR-001: AI Platform Boundaries, Technologies, and Integration Model

**Date:** 2026-10-24  
**Status:** Accepted  
**Deciders:** Principal AI Platform Architect, Senior Go Systems Architect, Application Security Lead  

---

## 1. Context and Problem Statement

HUNTIQ is a marketplace platform with an established Go Fiber backend, PostgreSQL 16 database, and Next.js frontend. As marketplace discovery, categorization, and query understanding requirements grow, future machine learning and predictive capabilities will be needed (e.g., query intent classification, category suggestions, duplicate detection, ranking assistance).

However, introducing AI directly into the monolithic Go backend or running Python scripts ad-hoc creates tight coupling, unpredictable dependency trees, security hazards, and operational debt. We need an architectural decision that defines how the AI capability will be organized, isolated, secured, and integrated with the existing platform.

---

## 2. Decision

We establish an **independent Python AI microservice (`ai/`)** and an **isolated offline Data Platform (`data-platform/`)**, adhering to the following strict boundaries:

### 2.1 Service Boundaries & Authority
1. **The Go Backend Remains the Sovereign Authority:**
   - Owns user authentication, authorization, role enforcement, and session lifecycles.
   - Owns all database mutations, transactions, financial ledgers, and marketplace listing state.
   - Owns canonical taxonomy and category definitions in PostgreSQL.
   - Serves as the gateway for all client requests; frontend clients never communicate directly with the AI service.
2. **The Python AI Service (`ai/`) is a Stateless Prediction Worker:**
   - Deployed strictly on an internal private network (not publicly accessible).
   - Has zero direct database write permissions.
   - Authenticates inter-service requests using a pre-shared secret header (`X-HUNTIQ-AI-SECRET`).
   - Returns typed suggestions and predictions with confidence scores.
   - Does not enforce marketplace business rules or modify listings.
3. **The Data Platform (`data-platform/`) is Offline and Isolated:**
   - Dedicated exclusively to dataset ingestion, anonymization, labeling, training, offline evaluation, and model artifact governance.
   - Completely decoupled from the production inference serving runtime.
   - Never contains real production secrets, user passwords, or un-anonymized PII in source control.

### 2.2 Technology Stack
- **Python Version:** Python 3.12+ (tested on system Python 3.14).
- **Package Manager:** `uv` with `pyproject.toml` and committed `uv.lock`.
- **API Framework:** `FastAPI` + `Uvicorn` for lean, asynchronous, typed HTTP serving.
- **Data Validation & Settings:** `Pydantic v2` and `pydantic-settings`.
- **Testing & Quality:** `pytest`, `httpx` (async test client), `pytest-cov`, and `ruff` (linter/formatter).
- **ML Training Dependencies:** Kept in an optional dependency group (`[project.optional-dependencies] ml`) containing `scikit-learn`, `pandas`, `pyarrow`, and `joblib`.

### 2.3 Failure Mode & Fault Tolerance
- **Untrusted Output Policy:** The Go backend treats every response from the AI service as an untrusted recommendation.
- **Fail-Open / Fallback Behavior:** If the AI service times out (default 500ms timeout) or returns an error, the Go backend falls back immediately to existing deterministic logic (e.g., standard PostgreSQL ILIKE / Full-Text Search). The marketplace user never experiences an outage due to an AI failure.

---

## 3. Considered Alternatives

### Alternative A: Embedding ML/AI Models Directly in Go (e.g., Go ONNX Runtime / Go bindings)
- *Pros:* Single binary deployment, no inter-service network hop.
- *Cons:* Extremely difficult for ML engineers to train, iterate, and evaluate models; limited library ecosystem in Go; potential memory bloat and CGO compilation complexity.
- *Verdict:* Rejected for core model exploration and training. (Lightweight ONNX runtime in Go may be re-evaluated in Phase 4 if latency mandates).

### Alternative B: Direct Database Access for Python Service
- *Pros:* Python service could query database directly for feature extraction.
- *Cons:* Bypasses Go domain layer, couples database schema changes to AI code, risks race conditions and unauthorized reads of sensitive data.
- *Verdict:* Rejected. The Go backend must prepare and pass necessary contextual features in the HTTP request payload.

### Alternative C: Third-Party Cloud LLM APIs (OpenAI, Anthropic) via Frontend or Direct Calling
- *Pros:* Fast initial prototype.
- *Cons:* Uncontrolled operational costs, vendor lock-in, data sovereignty violations (sending Ethiopian marketplace data abroad), high network latency, potential hallucination in canonical category assignment.
- *Verdict:* Prohibited for the foundational core. Specialized domain models must be trained or fine-tuned on local taxonomy data.

---

## 4. Consequences

### Positive
- Developers can immediately understand service ownership without ambiguity.
- Python dependencies are strictly pinned and isolated via `uv`.
- Go backend security, database integrity, and core marketplace transactions remain 100% safeguarded.
- Future ML developers have a standardized, pre-configured framework for training, evaluating, and packaging models.
- Tests can be run locally in seconds without GPU infrastructure or production credentials.

### Negative / Trade-offs
- An additional service to monitor, containerize, and configure in local development and production.
- Inter-service network latency (~2ms to 10ms) between Go backend and Python service.
