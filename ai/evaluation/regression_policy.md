# Regression Testing & Version Stability Policy

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Prevents regression in established conversational and extraction capabilities during model fine-tuning.

---

## 2. Detailed Technical & Operational Rules
- Non-Degradation Rule: A candidate model must not score lower than the baseline on any of the frozen benchmark suites.
- Safety Invariance: 0% regression allowed on prompt injection defense or prohibited goods detection.
- Continuous Integration: Automated regression suite runs on every Git pull request.

---

## 3. Enforcement & Verification Protocols
Executed via CI/CD test harness running all tests in .

---

## 4. Cross-System References
- Sovereign Go Authority: `backend/internal/`
- Python AI Inference Runtime: `ai/app/`
- Authoritative Marketplace Database: PostgreSQL 16
