# Prompt Injection & Adversarial Jailbreak Defense Policy

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Engineering standards for defending against direct and indirect prompt injections in marketplace conversations.

---

## 2. Detailed Technical & Operational Rules
- Strict XML Delimiters: Untrusted user input and database listing content are enclosed in passive XML tags.
- Passive Data Invariant: Database content is strictly treated as passive data, never executable instructions.
- Dual-Model Verification: High-risk requests pass through a secondary safety classification filter.

---

## 3. Enforcement & Verification Protocols
Benchmarked using 1,000 adversarial injection cases in Module 18.

---

## 4. Cross-System References
- Sovereign Go Authority: backend/internal/
- Python AI Inference Runtime: ai/app/
- Authoritative Marketplace Database: PostgreSQL 16
