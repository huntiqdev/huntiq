# HUNTIQ Dataset Quality Standards & Acceptance Thresholds

---

## 1. Metric Thresholds
- **Schema Compliance:** 100% required. Zero non-compliant records permitted.
- **PII Detection Rate:** 0.00% tolerance. All phone numbers, national IDs, and bank account numbers must be pseudonymized.
- **Language Detection Confidence:** $\ge 0.95$ for mono-lingual records; multi-lingual records must have explicit script span tagging.
- **Homophone Consistency (Ge'ez):** 100% unified to standard Ethiopian Academy of Languages canonical orthography.
