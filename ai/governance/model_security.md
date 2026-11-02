# Model Supply Chain Security & Artifact Integrity

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Security protocols for downloading, verifying, and storing model weights and runtime dependencies.

---

## 2. Detailed Technical & Operational Rules
- Safe Serialization Formats: Only Safetensors and GGUF formats are permitted; raw Python pickle files are strictly prohibited.
- Private Network Isolation: Production AI service runs on 127.0.0.1:8000 with zero public internet ingress.
- Secret Authentication: Header X-HUNTIQ-AI-SECRET required on all incoming HTTP requests.

---

## 3. Enforcement & Verification Protocols
Automated vulnerability scanning with Bandit and pip-audit.

---

## 4. Cross-System References
- Sovereign Go Authority: backend/internal/
- Python AI Inference Runtime: ai/app/
- Authoritative Marketplace Database: PostgreSQL 16
