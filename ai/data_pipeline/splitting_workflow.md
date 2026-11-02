# Train / Validation / Test Splitting Workflow

---

## Rules:
- 80% Training (`training/train.jsonl`).
- 10% Validation (`training/validation.jsonl`).
- 10% Test (`training/test.jsonl`).
- Use deterministic random seed (seed = 42).
- Group by user/session ID to prevent data leakage between splits.
