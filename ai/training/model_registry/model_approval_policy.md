# HUNTIQ Model Registry & Production Approval Policy

---

## 1. Approval Gate Protocol
A candidate fine-tuned model checkpoint may only be marked `STATUS: APPROVED_PRODUCTION` in `model_registry.json` when the following conditions are satisfied:

1. **Frozen Benchmark Evaluation:** Evaluated on all 3 golden test suites (`domain_benchmark.jsonl`, `conversational_benchmark.jsonl`, `entity_extraction.jsonl`).
2. **Macro F1 Score Threshold:** Macro F1 $\ge 0.92$ on intent classification and entity span extraction.
3. **Zero Security Regressions:** 100% defense against prompt injections in Module 18 evaluation suite.
4. **Latency Verification:** P95 response time $\le 120\text{ms}$ under 10 concurrent requests on target CPU inference instance.
5. **Dual Engineering Sign-Off:** Signed off by the AI Architect and Go Backend Lead.
