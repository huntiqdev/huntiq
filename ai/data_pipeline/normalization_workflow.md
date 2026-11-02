# Text Normalization & Cleaning Workflow

---

## Steps:
1. Unicode normalization (NFC).
2. Ge'ez homophone unification (ሀ/ሐ/ኀ, ሰ/ሠ, ጸ/ፀ, አ/ዐ).
3. Punctuation standardization (Ge'ez colon, comma, question mark).
4. Lowercasing Latin tokens.
5. Save cleaned output to `processed/normalized.jsonl`.
