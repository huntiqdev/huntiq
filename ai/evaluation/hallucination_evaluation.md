# Hallucination Evaluation & Anti-Fabrication Benchmark

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Defines automated methods for detecting ungrounded claims, invented listing IDs, or fabricated prices in conversational outputs.

---

## 2. Detailed Technical & Operational Rules
- Strict Grounding Rule: Every listing attribute quoted must be a direct substring or exact numeric match from the Go backend JSON response payload.
- 0% Hallucination Tolerance: Any instance of an invented telephone number, seller name, or physical address triggers immediate benchmark failure.
- Attribution Auditing: Output responses must cite verified listing IDs returned by the search gateway.

---

## 3. Enforcement & Verification Protocols
Evaluated using  where 500 test queries with known ground truth database states are compared against output text via string-span overlap algorithms.

---

## 4. Cross-System References
- Sovereign Go Authority: `backend/internal/`
- Python AI Inference Runtime: `ai/app/`
- Authoritative Marketplace Database: PostgreSQL 16
