# HUNTIQ Conversational Intelligence Pipeline Specification

---

## 1. Execution Guarantees
- **P95 Latency Budget:** Total inference decision latency $\le 120\text{ms}$ on CPU runtime.
- **Fail-Safe Gateways:** If confidence falls below $0.70$ (Module 28), route to gentle clarification (Module 12) rather than hallucinating search filters.
- **Go Sovereign Boundary:** Python never emits SQL queries. Python emits strictly typed JSON tool proposals to Go on `127.0.0.1:8080/internal/ai/tools`.
