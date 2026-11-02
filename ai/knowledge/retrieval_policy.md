# Knowledge Retrieval & Grounded Citation Policy

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Specifies how the AI assistant queries, validates, and cites official platform knowledge documents.

---

## 2. Detailed Technical & Operational Rules
- Official Authority: Knowledge citations must reference verified documents in ai/knowledge/.
- Exact Quote Accuracy: Fee schedules, operating hours, and contact numbers must be quoted verbatim.
- Fallback on Uncertainty: If an inquiry cannot be answered by official documentation, route to human support.

---

## 3. Enforcement & Verification Protocols
Verified in ai/datasets/03_huntiq_platform_knowledge/evaluation/.

---

## 4. Cross-System References
- Sovereign Go Authority: backend/internal/
- Python AI Inference Runtime: ai/app/
- Authoritative Marketplace Database: PostgreSQL 16
