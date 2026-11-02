# HUNTIQ Dataset Maintenance & Deprecation Guide

---

## 1. Routine Audits
- **Monthly License Scan:** Verify that all referenced public sources maintain permissive commercial terms.
- **Quarterly Drift Analysis:** Compare live marketplace query distributions against dataset distributions to identify newly emerging slang or product categories.
- **Deprecation Procedure:** Outdated records (e.g., obsolete model names or discontinued currency notes) are archived into `archive/v1/` rather than silently deleted.
