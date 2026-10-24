# HUNTIQ AI Service

The internal AI inference microservice for **HUNTIQ** (Ethiopia's Universal Marketplace).

---

## 1. Role and Boundaries

- **Stateless Prediction Worker:** Provides query classification, category suggestions, and duplicate detection for the authoritative Go backend.
- **Strict Boundary:** Does not write to PostgreSQL, does not manage user authentication, and is not exposed to the public internet.
- **Inter-Service Security:** All requests require the `X-HUNTIQ-AI-SECRET` header.

---

## 2. Prerequisites

- Python 3.12+ (tested on Python 3.14)
- `uv` package manager (`curl -LsSf https://astral.sh/uv/install.sh | bash`)

---

## 3. Local Development

### 1. Setup Environment
```bash
cp .env.example .env
uv sync
```

### 2. Run Test Suite
```bash
uv run pytest -v
```

### 3. Run Linting & Formatting
```bash
uv run ruff check .
uv run ruff format --check .
```

### 4. Start Local Server
```bash
uv run uvicorn huntiq_ai.main:app --host 127.0.0.1 --port 8000 --reload
```

---

## 4. Endpoints

| Endpoint | Method | Auth | Description |
| :--- | :--- | :--- | :--- |
| `/health/live` | GET | None | Liveness check confirming process responsiveness. |
| `/health/ready` | GET | None | Readiness check confirming configuration and models. |
| `/api/v1/predict/query-intent` | POST | `X-HUNTIQ-AI-SECRET` | *(Planned)* Disambiguates search queries into canonical category filters. |
| `/api/v1/predict/category-suggestion` | POST | `X-HUNTIQ-AI-SECRET` | *(Planned)* Suggests categories from title and description. |

---

## 5. Current Implementation Status

- [x] Typed settings with Pydantic (`core/config.py`).
- [x] Structured JSON logging (`core/logging.py`).
- [x] Inter-service secret verification (`core/security.py`).
- [x] Standard error handling (`core/errors.py`).
- [x] Liveness and readiness endpoints (`/health/live`, `/health/ready`).
- [x] Bounded request payloads (64KB limit).
- [x] Unit, integration, and contract tests.
- [ ] *Predictive model weights and inference pipelines (deferred to Phase 2 after offline baseline evaluation).*

---

## 6. Subsystems

- **Inference Runtime (`src/huntiq_ai`):** Stateless FastAPI microservice running on private network.
- **Offline Data Platform (`data-platform/`):** Dataset contracts, PII scrubbing, manifest validation, and model cards. See [`data-platform/README.md`](data-platform/README.md).

---

## 7. Architecture & Platform Documentation

All AI platform specifications, threat models, and runbooks are organized in [`docs/`](docs/):
- [AI Architecture & Threat Model](docs/architecture/ai-architecture.md)
- [Go / Python Integration Contract](docs/architecture/ai-integration.md)
- [ADR-001: Boundaries & Stack](docs/decisions/ADR-001-ai-platform-boundaries.md)
- [AI Data Security](docs/security/ai-data-security.md)
- [AI Testing Strategy](docs/testing/ai-testing.md)
- [Local Development Runbook](docs/runbooks/ai-local-development.md)
