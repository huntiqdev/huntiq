# HUNTIQ AI Platform Foundation — Repository Audit & Proposal Classification

**Date:** 2026-10-24  
**Author:** Principal AI Platform Architect & Systems Team  
**Status:** Approved for Implementation  

---

## 1. Executive Summary

This audit establishes the baseline for introducing an independent, production-grade AI platform foundation into the **HUNTIQ** marketplace repository. HUNTIQ is an Ethiopia-first universal marketplace built on a Go backend (Fiber v2), PostgreSQL 16 database, and Next.js 15 TypeScript frontend. 

Prior to this work, no AI or machine learning code existed in the repository. The objective is to establish a modular, secure, version-controlled, and testable AI foundation without altering existing business logic, rewriting existing applications, or introducing premature operational complexity.

---

## 2. Repository & Environment Audit

### 2.1 Git Status & Workspace Roots
- **Repository Root:** `/home/poven/GO PROJECT/sellbuy`
- **Current Branch:** `main` (tracked to `origin/main`)
- **Working Tree State:** Active work on taxonomy expansion (`backend/cmd/seed-taxonomy/domain_master_expander.go` and `catalog_domains_*.go`). All existing working-tree changes are strictly preserved.

### 2.2 Existing Systems Architecture
- **Backend (`backend/`):**
  - **Language & Runtime:** Go 1.23+
  - **Framework:** Fiber v2 (`github.com/gofiber/fiber/v2`)
  - **Database Driver:** `pgx/v5` (`github.com/jackc/pgx/v5/pgxpool`)
  - **Core Domains:** Authentication (JWT, OTP encryption), Marketplace Products/Listings, Categories & Attribute Hierarchies, Trust & Safety, Transactions & Ledger, Seller & Buyer Management, Referrals, Location Services.
  - **Test Suite:** Go standard testing with table-driven tests (`go test ./...`).
- **Frontend (`frontend/`):**
  - **Framework:** Next.js 15.1.7 App Router
  - **UI Library:** React 18.3.1, Tailwind CSS 3.4.17, Lucide React icons
  - **Language:** TypeScript 5.7.2
  - **Build & Quality:** `npm run build`, `npm run lint`.
- **Database (`PostgreSQL 16`):**
  - Schema migrations tracked in `backend/internal/database/migrations/` (migrations 001 to 031).
  - Search implementation: PostgreSQL Full-Text Search, B-Tree and GIN indexes, query logging in `search_queries`.
- **CI/CD & Deployment:**
  - GitHub Actions: `.github/workflows/webpack.yml` currently runs Node.js builds.
  - Runtime environment configuration via `.env` in `backend/` and `.env.local` in `frontend/`.

### 2.3 Python & Tooling Environment
- **System Python:** `/usr/bin/python3` (Python 3.14.8)
- **Environment Management:** Linux with PEP 668 managed environment.
- **Package & Project Manager:** `uv` 0.13.5 installed in `~/.local/bin/uv`. `uv` provides blazing-fast dependency resolution, standard `pyproject.toml` support, and deterministic `uv.lock` files.

---

## 3. Classification of Proposed Changes

Following strict platform governance, all proposed components are classified into three distinct categories:

| Component | Proposed Path | Classification | Justification |
| :--- | :--- | :--- | :--- |
| **Python AI Service Scaffold** | `ai/` | **Required now** | Standalone FastAPI service container with typed settings, structured logging, internal auth, and health checks. |
| **AI Project Metadata & Lock** | `ai/pyproject.toml`, `ai/uv.lock` | **Required now** | Pinned dependencies and deterministic reproducibility. |
| **Service Configuration & Security** | `ai/src/huntiq_ai/core/` | **Required now** | Typed Pydantic configuration, fail-fast validation, internal service authentication token validation, structured JSON logging. |
| **Health & Operations Endpoints** | `ai/src/huntiq_ai/api/routes/health.py` | **Required now** | Distinct `/health/live` and `/health/ready` endpoints with safe metadata. |
| **AI Automated Test Suite** | `ai/tests/` | **Required now** | Isolated pytest suite verifying config validation, request handling, error formatting, and security boundaries with synthetic fixtures. |
| **Data Platform Scaffold** | `data-platform/` | **Required now** | Directory structure for data contracts, dataset manifests, quality rules, and ML experiment lifecycles. |
| **Dataset Contract & Manifest Schemas** | `data-platform/data-contracts/` | **Required now** | JSONSchema/Pydantic schemas for dataset manifests, query intent contracts, and category suggestion contracts. |
| **Architecture & Governance Docs** | `docs/architecture/`, `docs/security/`, `docs/decisions/` | **Required now** | Complete architectural documentation, ADR-001, security threat model, and runbooks. |
| **Developer Guidelines** | `AGENTS.md`, `ai/AGENTS.md`, `data-platform/AGENTS.md` | **Required now** | Multi-level engineering and agent operating constraints. |
| **ML Training Optional Dependencies** | `ai/pyproject.toml` `[optional-dependencies]` | **Optional foundation** | Optional groups for `scikit-learn`, `pandas`, `pyarrow`, `joblib` so training is decoupled from lean inference. |
| **Containerization Definition** | `ai/Dockerfile`, `ai/.dockerignore` | **Optional foundation** | Multi-stage, non-root Docker container definition for production alignment. |
| **Predictive Inference Endpoints** | `ai/src/huntiq_ai/api/routes/predict.py` | **Future implementation** | Prediction endpoints remain unimplemented until specific task and baseline models are approved. |
| **Model Training & Pipelines** | `data-platform/ml/pipelines/` | **Future implementation** | Model training will be executed only after approved source data extraction and offline evaluation baselines. |
| **Vector DB / LLM Infrastructure** | External Services | **Future implementation** | Excluded from foundation; strictly justified on benchmark performance in later phases. |

---

## 4. Safety Guardrails & Invariants

1. **Zero Database Schema Changes:** No changes to existing PostgreSQL tables, migrations, or database state.
2. **Authoritative Backend Preservation:** The Go backend remains the single source of truth for all business rules, user identity, permissions, listing state, and financial transactions.
3. **Internal Network Isolation:** The AI service will never be exposed to the public internet; communication occurs strictly over private inter-service HTTP.
4. **Untrusted AI Output Rule:** Any future AI output received by the Go backend must be treated as untrusted suggestions, validated against business schemas, and subject to fallback.
5. **No Production Secrets or PII in Git:** Training datasets, private marketplace exports, raw user queries, and credentials are kept out of source control.
