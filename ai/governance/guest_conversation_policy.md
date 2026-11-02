# Guest Conversation Handling & Ephemeral Privacy Policy

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Defines privacy safeguards for unauthenticated visitors browsing HUNTIQ.

---

## 2. Detailed Technical & Operational Rules
- Zero Account Requirement: Guests can search and converse with the AI freely without logging in.
- Session Isolation: Each guest receives a random ephemeral session ID stored only in session cookies.
- Memory Flushed: Guest context is completely erased from Redis/memory after 30 minutes.

---

## 3. Enforcement & Verification Protocols
Enforced in ai/app/core/context_manager.py.

---

## 4. Cross-System References
- Sovereign Go Authority: backend/internal/
- Python AI Inference Runtime: ai/app/
- Authoritative Marketplace Database: PostgreSQL 16
