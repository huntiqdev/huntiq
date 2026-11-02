# Human Evaluation Rubrics & Protocol

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Defines qualitative assessment criteria for native speaker linguists and domain specialists.

---

## 2. Detailed Technical & Operational Rules
- Naturalness Score (1-5): Grammatical correctness and natural Ethiopian phrasing.
- Grounding Score (1-5): Faithfulness to the underlying marketplace listings.
- Cultural Respect Score (1-5): Proper formal/casual register and polite greeting protocol.
- Inter-Annotator Agreement: Minimum Fleiss' Kappa >= 0.80 required.

---

## 3. Enforcement & Verification Protocols
Quarterly evaluations conducted on random samples of anonymized production interactions.

---

## 4. Cross-System References
- Sovereign Go Authority: `backend/internal/`
- Python AI Inference Runtime: `ai/app/`
- Authoritative Marketplace Database: PostgreSQL 16
