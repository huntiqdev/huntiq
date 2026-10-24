# HUNTIQ AI Platform Documentation Index

This directory contains all documentation for the HUNTIQ AI Platform, microservice architecture, integration contracts, data governance, security policies, and testing standards.

## 📚 Table of Contents

### 1. Architecture & Integration
- [AI Architecture & Threat Model](architecture/ai-architecture.md) — System component diagrams, trust boundaries, request lifecycles, and failure containment modes.
- [Go / Python Integration Contract](architecture/ai-integration.md) — HTTP REST contract, mutual authentication (`X-HUNTIQ-AI-SECRET`), correlation IDs, timeout thresholds, and Go deterministic fallback policies.

### 2. Architecture Decisions (ADR)
- [ADR-001: AI Platform Boundaries & Technology Stack](decisions/ADR-001-ai-platform-boundaries.md) — Rationale for Go backend sovereignty, stateless Python FastAPI internal workers, private network topology, and deferred vector database adoption.

### 3. Security & Governance
- [AI Data Security & Threat Modeling](security/ai-data-security.md) — Threat model, Ethiopian mobile number (`+2519...`, `09...`) and PII scrubbing standards, credential hygiene, and isolation from public internet.

### 4. Testing & Quality
- [AI Testing Strategy & Quality Gates](testing/ai-testing.md) — Unit testing, OpenAPI/RFC 7807 contract testing, data quality validation gates, and train/test leakage prevention.

### 5. Runbooks & Development
- [AI Local Development Runbook](runbooks/ai-local-development.md) — Step-by-step developer environment setup with `uv`, testing commands, Docker build instructions, and common troubleshooting tips.

### 6. Audits
- [AI Platform Foundation Pre-Audit](audits/ai-platform-foundation-audit.md) — Baseline inspection of the HUNTIQ repository, dependencies, database migrations, and existing search implementation prior to scaffold implementation.

---

## 🧭 Related Developer Guidelines
- [Root Repository Instructions](../../AGENTS.md)
- [AI Platform Guidelines (ai/AGENTS.md)](../AGENTS.md)
- [Data Platform Guidelines (ai/data-platform/AGENTS.md)](../data-platform/AGENTS.md)
