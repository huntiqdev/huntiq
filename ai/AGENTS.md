# HUNTIQ AI Service (`ai/`) Instructions for Agents & Developers

**Directory:** `/ai`  
**Runtime:** Python 3.12+ (tested on 3.14.8)  
**Package Manager:** `uv` with `pyproject.toml` and `uv.lock`  
**Web Framework:** FastAPI + Uvicorn  

---

## 1. Directory Structure

```text
ai/
├── AGENTS.md                  # This file
├── README.md                  # Human developer documentation
├── pyproject.toml             # Pinned project metadata & dependencies
├── uv.lock                    # Deterministic lockfile (MUST BE COMMITTED)
├── .python-version            # Pinned Python version
├── .env.example               # Safe template with placeholder values
├── Dockerfile                 # Minimal, non-root container
├── .dockerignore              # Ignore venv, caches, tests, data-platform, docs
│
├── data-platform/             # Offline dataset governance, contracts, validator
├── docs/                      # Dedicated AI architecture, ADRs, security, runbooks
│
├── src/
│   └── huntiq_ai/
│       ├── __init__.py
│       ├── main.py            # FastAPI application factory & lifecycle
│       ├── api/
│       │   ├── __init__.py
│       │   ├── router.py      # Aggregates sub-routers
│       │   └── routes/
│       │       ├── __init__.py
│       │       └── health.py  # /health/live and /health/ready
│       ├── core/
│       │   ├── __init__.py
│       │   ├── config.py      # Pydantic Settings
│       │   ├── logging.py     # Safe structured JSON logging
│       │   ├── errors.py      # Standard RFC 7807 error handlers
│       │   └── security.py    # Internal pre-shared secret validation
│       └── schemas/
│           ├── __init__.py
│           └── health.py      # Health response schemas
│
└── tests/
    ├── conftest.py            # Shared fixtures (TestClient, mock settings)
    ├── unit/
    │   ├── test_config.py     # Validates setting defaults & overrides
    │   ├── test_health.py     # Validates liveness & readiness logic
    │   └── test_security.py   # Validates token rejection and acceptance
    └── contract/
        └── test_contract.py   # Validates response schema adherence
```

---

## 2. Coding and Import Conventions

1. **Absolute Imports:** Always use absolute imports rooted at `huntiq_ai`:
   ```python
   # Correct:
   from huntiq_ai.core.config import get_settings
   from huntiq_ai.core.security import verify_internal_secret

   # Incorrect:
   from ..core.config import get_settings
   ```
2. **Typed Schemas:** All request and response models must inherit from `pydantic.BaseModel` with explicit field types and docstrings.
3. **No Global Mutable State:** Instantiate dependencies through FastAPI `Depends` or application state attached to `app.state`.
4. **Structured Logging:** Use `get_logger(__name__)`. Never use `print()` in production code. Never log authorization headers, passwords, or raw user PII.
5. **No Direct Database Access:** The AI service must not connect directly to PostgreSQL. All necessary attributes are received in the HTTP request payload from Go.
6. **No Mock AI Pretending to Be Real:** Never create a prediction endpoint that returns hard-coded fake predictions while claiming to be an AI model.

---

## 3. Dependency Management Commands

- **Install/Sync dependencies:** `uv sync`
- **Add a runtime dependency:** `uv add <package>`
- **Add a dev/test dependency:** `uv add --dev <package>`
- **Run test suite:** `uv run pytest -v`
- **Run lint check:** `uv run ruff check .`
- **Run code format:** `uv run ruff format .`
- **Start development server:** `uv run uvicorn huntiq_ai.main:app --reload --port 8000`
