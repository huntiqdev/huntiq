# HUNTIQ Dataset Contribution Guide

**Audience:** Internal ML engineers, linguistic contractors, and marketplace domain experts.

---

## 1. Contribution Workflow
1. **Branching:** Create a feature branch named `data/modNN-description`.
2. **Schema Adherence:** Ensure all new records conform to `schema.json` in the target module.
3. **Linguistic Verification:** Amharic, Afaan Oromoo, and Tigrinya samples must be verified by a native speaker.
4. **Validation Test:** Run `uv run pytest tests/test_benchmark.py` before opening a pull request.
5. **No PII:** Never include real phone numbers, credit card numbers, or personal national IDs.
