# Cross-Dataset Duplicate Management & Partitioning

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
Guarantees that training examples are uniquely allocated and not redundantly duplicated across multiple module datasets.

---

## 2. Detailed Technical & Operational Rules
- Hash Partitioning: Exact utterance hashes are indexed centrally in dataset_registry.json.
- Multi-Task Exception Rule: If an utterance is intentionally used for both intent and entity training, it must be explicitly tagged as multi_task.
- Zero Leakage Guarantee: Ensures no evaluation example in one module appears in the training split of another module.

---

## 3. Enforcement & Verification Protocols
Scanned via automated duplicate detection pipelines in ai/data_pipeline/.

---

## 4. Cross-System References
- Sovereign Go Authority: backend/internal/
- Python AI Inference Runtime: ai/app/
- Authoritative Marketplace Database: PostgreSQL 16
