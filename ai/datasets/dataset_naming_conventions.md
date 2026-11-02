# HUNTIQ Dataset Naming Conventions

---

## 1. Directory & File Standards
- **Module Directories:** Exactly `NN_dataset_name` where `NN` is a zero-padded integer from `01` to `30`.
- **JSON Schemas:** Named `schema.json` within module root.
- **Manifests:** Named `manifest.json`.
- **Split Files:** Strictly `train.jsonl`, `validation.jsonl`, `test.jsonl`.
- **Language Splits:** ISO-639-1 two-letter lowercase codes: `en.jsonl`, `am.jsonl`, `or.jsonl`, `ti.jsonl`.
