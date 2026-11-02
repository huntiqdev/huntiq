# Model Rollback & Downgrade Protocols

---

## Protocols:
If a newly deployed model exhibits higher hallucination rate or latency regression in production, Go backend immediately routes traffic to the previous approved checkpoint via `model_registry.json`.
