# HUNTIQ Grounded Knowledge Architecture

---

## 1. Architecture Overview
HUNTIQ's platform knowledge system bridges verified company documentation with real-time conversational retrieval.

```text
[Official Policy Markdown (ai/knowledge/)]
                    │
                    ▼
[Deterministic Chunking & Embedding Pipeline]
                    │
                    ▼
[Local Vector / Keyword Index (Ollama / In-Memory)]
                    │
                    ▼
[Context Injection via Strict Pass-Through Delimiters]
                    │
                    ▼
[Grounded Answer Generation (Module 03 & 30)]
```

## 2. Inviolable Truth Hierarchy
1. Live Database State (PostgreSQL via Go Backend): Sole authority on listing prices, stock, and status.
2. Official Policy Documents (`ai/knowledge/`): Sole authority on fees, terms, and rules.
3. Model Parametric Memory: Never used as an authority on facts.
