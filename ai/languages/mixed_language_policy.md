# Code-Switching & Mixed-Script Handling Policy

---

## 1. Multi-Script Tolerance
Sentences blending English and Ge'ez (e.g. "Toyota Vitz መኪና ዋጋው ስንት ነው?") must not trigger parser errors. The tokenizer must process Latin and Ge'ez script blocks gracefully within the same utterance.
