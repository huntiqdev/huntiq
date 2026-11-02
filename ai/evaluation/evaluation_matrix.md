# Master AI Capability Evaluation Matrix

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Central matrix mapping every AI module to its primary benchmark, metric formula, target threshold, and evaluation cadence.

---

## 2. Detailed Technical & Operational Rules
- Classifiers: Evaluated via Precision, Recall, and Macro-F1.
- Extractors: Evaluated via Exact Match (EM) and token-level F1.
- Dialogue Systems: Evaluated via Goal Completion Rate and Slot Error Rate.
- Latency: Evaluated via P50/P95 response time in milliseconds.

---

## 3. Enforcement & Verification Protocols
Updated alongside every release of .

---

## 4. Cross-System References
- Sovereign Go Authority: `backend/internal/`
- Python AI Inference Runtime: `ai/app/`
- Authoritative Marketplace Database: PostgreSQL 16
