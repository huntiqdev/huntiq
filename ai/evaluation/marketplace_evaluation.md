# Marketplace Domain Semantic Evaluation

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Evaluates contextual reasoning over Ethiopian commerce conventions, land terminology, and vehicle specs.

---

## 2. Detailed Technical & Operational Rules
- Currency Reasoning: Correctly understands Birr, ETB, "ሺህ" (thousands), and "ሚሊዮን" (millions).
- Land Unit Reasoning: Correctly parses "ካሬ" (square meters) and "ጋሻ" (gasha).
- Vehicle Trade Conventions: Correctly distinguishes commercial taxi permits, code 2/code 3 plates, and customs duty clearance.

---

## 3. Enforcement & Verification Protocols
Benchmarked against golden marketplace question-answering pairs across diverse Ethiopian product sectors.

---

## 4. Cross-System References
- Sovereign Go Authority: `backend/internal/`
- Python AI Inference Runtime: `ai/app/`
- Authoritative Marketplace Database: PostgreSQL 16
