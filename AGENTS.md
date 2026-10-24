# HUNTIQ Repository Instructions for AI Agents & Developers

**Project:** HUNTIQ — Universal Marketplace Platform (Ethiopia)  
**Architecture:** Go Backend (Fiber v2) + PostgreSQL 16 + Next.js 15 (TypeScript) + Python AI Service (FastAPI)  

---

## 1. Service Ownership & Boundaries

```
huntiq/
├── backend/            # PILLAR 1: Sovereign Authority (Business logic, Auth, DB mutations, Categories, backend/docs)
├── frontend/           # PILLAR 2: Presentation Layer (Next.js 15 App Router, Tailwind CSS, Localized UI, frontend/docs)
└── ai/                 # PILLAR 3: AI Platform (FastAPI inference worker, offline data platform, ai/docs)
```

### Invariant Rules:
1. **The Go Backend is the Sole Sovereign Authority:**
   - All user authentication, authorization, listing lifecycle, payments, and database writes belong in `backend/`.
   - Never duplicate marketplace business logic in Python.
   - The Go backend never accepts unverified AI suggestions; all predictions are validated against strict business schemas and subject to immediate deterministic fallback.
2. **The Python AI Service is an Internal Worker:**
   - Deployed strictly on a private network (`127.0.0.1:8000`). Never expose port 8000 to the public internet.
   - Has zero database write credentials.
   - Authenticates incoming requests via header `X-HUNTIQ-AI-SECRET`.
3. **Data Platform is Offline:**
   - `ai/data-platform/` is strictly for data contracts, cleaning pipelines, offline training, and model cards.
   - Never commit raw datasets, user PII, or trained model binaries to Git.

---

## 2. Common Verification & Test Commands

### Go Backend
```bash
cd backend
go test ./...                       # Run all backend unit & integration tests
go run ./cmd/api                    # Start API server on :8080
```

### Next.js Frontend
```bash
cd frontend
npm run lint                        # Run ESLint checks
npm run build                       # Run production TypeScript build
npm run dev                         # Start dev server on :3000
```

### Python AI Service
```bash
cd ai
uv run pytest -v                    # Run all AI unit, integration, and contract tests
uv run ruff check .                 # Run Python linter
uv run ruff format --check .        # Check formatting
uv run uvicorn huntiq_ai.main:app   # Start local inference service on :8000
```

### Data Platform
```bash
cd ai/data-platform
uv run pytest tests/ -v             # Run data contract and quality validation tests
```

---

## 3. Strict Prohibitions

- **NO DESTRUCTIVE DB RESETS:** Never drop tables, reset the database, or delete active records.
- **NO FAKE AI:** Never create endpoints that return hard-coded random or fabricated predictions while claiming to be an AI model. If a model is not yet trained and approved, keep the endpoint explicitly unimplemented.
- **NO CREDENTIAL COMMITS:** Never commit passwords, tokens, API keys, or production `.env` files.
- **NO PII IN DATASETS:** Never collect passwords, credit cards, or raw un-anonymized phone numbers into `data-platform/`.
- **NO EXTERNAL DATA LEAKS:** Never transmit Ethiopian marketplace listings or queries to unvetted external public LLM APIs without explicit written approval.
- **PRESERVE WORKING CHANGES:** Never discard another developer's uncommitted work or run destructive `git reset --hard`.

---

## 4. How to Introduce a New AI Capability

1. **Step 1 — Define Contract:** Specify input/output schemas in `data-platform/data-contracts/schemas/`.
2. **Step 2 — Data & Baseline:** Prepare cleaned, anonymized dataset in `data-platform/ml/datasets/processed/` with a deterministic train/test split. Establish a non-AI rule-based baseline.
3. **Step 3 — Train & Evaluate:** Train candidate model offline, measure precision/recall/F1/latency against the baseline, and document results in `data-platform/ml/model_cards/`.
4. **Step 4 — Implement Inference:** Add typed route in `ai/src/huntiq_ai/api/routes/` with schema validation and authentication.
5. **Step 5 — Go Client Integration:** Update Go backend client with timeout, fallback, and validation.
6. **Step 6 — Tests:** Add unit, integration, and contract tests before deploying.
