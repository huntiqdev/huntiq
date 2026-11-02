# Multi-Turn Dialogue & State Tracking Evaluation

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Assesses conversational memory, slot inheritance, and graceful context switching.

---

## 2. Detailed Technical & Operational Rules
- Slot Retention: User refinements (e.g. "make that cheaper") must preserve previous category and location slots with >= 98% accuracy.
- Domain Pivot Accuracy: When user changes topic ("forget houses, show me cars"), stale slots must be completely cleared.
- Follow-Up Coherence: Dialogue act responses must align with the active turn state.

---

## 3. Enforcement & Verification Protocols
Tested against multi-turn dialogue trees in .

---

## 4. Cross-System References
- Sovereign Go Authority: `backend/internal/`
- Python AI Inference Runtime: `ai/app/`
- Authoritative Marketplace Database: PostgreSQL 16
