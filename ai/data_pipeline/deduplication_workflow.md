# Deduplication & Quality Control Workflow

---

## Steps:
1. Compute MD5 text hashes to remove exact duplicate utterances.
2. Compute MinHash LSH to detect near-duplicate queries (>90% Jaccard similarity).
3. Write unique records to `processed/deduplicated.jsonl`.
