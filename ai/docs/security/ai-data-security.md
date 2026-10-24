# HUNTIQ AI Platform Data Security and Governance

**Document Version:** 1.0.0  
**Target:** Security Engineers, MLOps Engineers, and Data Custodians  
**Status:** Approved Security Architecture  

---

## 1. Threat Model & Security Principles

The HUNTIQ AI platform operates under a **Zero-Trust, Minimal-Privilege** framework:

1. **No External Attack Surface:** The AI service is strictly deployed on private loopback / VPC internal networks. It does not accept direct incoming traffic from public internet users.
2. **Untrusted Data Origin:** All input strings originating from users (search queries, listing titles, listing descriptions, candidate CV metadata) are treated as untrusted and potentially hostile (injection attacks, oversized payloads).
3. **Defense in Depth:** The Go backend validates input before sending it to Python; the Python service validates input before processing; the Go backend validates predictions before acting on them.

---

## 2. Sensitive Data Classification Matrix

| Data Category | Examples | Permitted in Production AI Inference? | Permitted in ML Training Datasets? | Permitted in Logs? |
| :--- | :--- | :--- | :--- | :--- |
| **Authentication Secrets** | Password hashes, JWT secrets, OTP seeds, API keys | **STRICTLY PROHIBITED** | **STRICTLY PROHIBITED** | **STRICTLY PROHIBITED** |
| **Financial / Payment** | Telebirr/CBE references, account numbers, card tokens | **STRICTLY PROHIBITED** | **STRICTLY PROHIBITED** | **STRICTLY PROHIBITED** |
| **Private Documents** | Uploaded PDF resumes, candidate ID cards, tax certificates | **STRICTLY PROHIBITED** | **STRICTLY PROHIBITED** | **STRICTLY PROHIBITED** |
| **Personal Identifiers** | Phone numbers, personal email addresses, national IDs | **STRICTLY PROHIBITED** | **STRICTLY PROHIBITED** (must be scrubbed/anonymized) | **STRICTLY PROHIBITED** |
| **Public Marketplace Text** | Product titles, category slugs, brand names, model numbers | Permitted | Permitted (in versioned Parquet/JSONL) | Truncated/Redacted only |
| **Search Queries** | User search strings | Permitted (sanitized) | Permitted (aggregated, scrubbed of PII) | Hashed or truncated only |

---

## 3. Data Sanitization & Scrubbing Pipeline

Before any dataset export from PostgreSQL is introduced to `data-platform/ml/datasets/raw/`, it must pass through an automated privacy filter:

1. **Regex Pattern Scrubbing:**
   - Ethiopian Phone Numbers (`/(?:\+251|0)9[0-9]{8}/g` -> `[PHONE_REDACTED]`)
   - Email Addresses (`/[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+/g` -> `[EMAIL_REDACTED]`)
   - Credit Cards / Account numbers
2. **Deterministic Pseudonymization:**
   - User IDs are hashed using a salted HMAC before being stored in analytics or interaction logs.
3. **No CV Content Leaks:**
   - Candidate CV documents uploaded to `cv_documents` table in PostgreSQL are private by default. They are never ingested into ML pipelines without explicit candidate opt-in and de-identification.

---

## 4. Model Artifact Security & Supply Chain Integrity

1. **Checksum Verification:**
   - Every trained model file (e.g. `.joblib`, `.onnx`) stored in `data-platform/ml/models/` must have a corresponding SHA-256 checksum recorded in its Model Card.
   - The serving application verifies the SHA-256 hash before loading weights into memory to prevent tampering.
2. **Safe Deserialization:**
   - Prefer standard formats like ONNX or native array weights over arbitrary untrusted pickle files.
   - Models are loaded only from a designated, read-only directory (`MODEL_DIR`).
3. **No Dynamic Model Downloading in Production:**
   - The production container image or deployment volume contains approved, verified model weights.
   - The service will never automatically download model files from unverified third-party internet URLs upon startup.

---

## 5. Non-Root Container Execution

The AI service Dockerfile enforces:
- Non-root user `appuser` (UID 10001).
- Read-only root filesystem wherever possible.
- Minimum required file permissions on model directories (`chmod 400 /app/models/*`).
- Dropped Linux capabilities (`cap_drop: ALL`).

---

## 6. Audit & Incident Response

- All security-relevant events (authentication failures, invalid tokens, payload rejections) are logged with severity `WARNING` or `ERROR`.
- If an authorization failure occurs, the response includes only an opaque error code (`UNAUTHORIZED`), never revealing whether the token was malformed, expired, or mismatched.
