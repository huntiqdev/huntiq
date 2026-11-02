# HUNTIQ Shared Intelligence Architecture

**Role:** Orchestrates the 30 conversational AI capabilities into a deterministic, high-performance marketplace reasoning pipeline.

---

## 1. Processing Pipeline
```text
User Utterance
      │
      ▼
[Reception & Normalization (Modules 14, 15)]
      │
      ▼
[Security & Safety Filter (Module 18)]
      │
      ▼
[Intent & Domain Classification (Modules 01, 02, 04, 05, 21)]
      │
      ▼
[Entity, Attribute & Constraint Extraction (Modules 06, 07, 08, 09, 10, 17, 23)]
      │
      ▼
[Dialogue State & Disambiguation (Modules 11, 12, 13, 22, 24, 25, 26, 27)]
      │
      ▼
[Authorized Tool Proposal to Go Backend (Module 19)]
      │
      ▼
[Grounded NLG Generation (Modules 03, 16, 20, 28, 29, 30)]
```
