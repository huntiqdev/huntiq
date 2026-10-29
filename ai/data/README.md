# HUNTIQ Offline AI Data Repository

This directory contains offline data preparation guidelines and local datasets for fine-tuning evaluation.

## Invariant Rules
1. **Never commit raw PII:** Passwords, credit cards, bank accounts, or raw un-anonymized phone numbers must be scrubbed before saving data.
2. **Deterministic Train/Test Split:** Any offline fine-tuning dataset must strictly separate train, validation, and held-out test splits to prevent data leakage.
3. **PostgreSQL as Truth:** All marketplace entities and IDs must correspond to real schema types defined in the Go backend.
