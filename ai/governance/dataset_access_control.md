# Role-Based Access Control (RBAC) for AI Data Platform

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Restricts access to raw training data, model checkpoints, and annotation pipelines based on verified job role.

---

## 2. Detailed Technical & Operational Rules
- Annotators: Restricted to labeling interfaces without direct production database access.
- Data Engineers: Read/write access to ai/data_pipeline/ and ai/datasets/.
- Production Services: Read-only access to approved model weights on 127.0.0.1:8000.

---

## 3. Enforcement & Verification Protocols
Enforced via Linux filesystem permissions and Git branch protection rules.

---

## 4. Cross-System References
- Sovereign Go Authority: backend/internal/
- Python AI Inference Runtime: ai/app/
- Authoritative Marketplace Database: PostgreSQL 16
