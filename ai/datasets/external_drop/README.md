# HUNTIQ External Dataset Drop-In Platform

This directory is configured for high-volume batch ingestion from Kaggle, HuggingFace, or custom dataset exports.

## Folder Organization

You can drop single files, batches of 10, or 100+ CSV, TSV, JSON, JSONL, or Parquet files into any of these task folders:

```
ai/datasets/external_drop/
├── intent/          # Intent classification datasets (buy, sell, rent, price_inquiry, etc.)
├── category/        # Category classification (electronics, vehicles, real_estate, etc.)
├── safety/          # Scam detection, fraud, advance-fee, phishing, harassment
├── domain/          # In-domain marketplace vs. general out-of-domain queries
└── sentiment/       # Reviews, disputes, satisfaction, complaint ratings
```

## Supported Header Formats (Automatic Detection)

The ingestion pipeline automatically detects standard column variations without requiring manual formatting:

- **Text Column:** `text`, `query`, `sentence`, `prompt`, `comment`, `message`, `review`, `utterance`, `clean_text`
- **Label Column:** `label`, `intent`, `category`, `class`, `sentiment`, `target`, `tag`, `fraud`, `scam`

## How to Run Ingestion & Re-Train Models

Once you drop your CSVs, run:

```bash
cd ai
uv run python scripts/batch_ingest_kaggle.py
```

The script will:
1. Scan all folders and subdirectories
2. Read and merge up to hundreds of thousands of samples
3. Normalize Amharic (Fidel), Oromo, Tigrinya, and English characters
4. Train modern multi-resolution subword classifiers with letter-level n-grams and dense continuous vector embeddings
5. Save updated checkpoints into `ai/training/checkpoints/` and update `model_registry.json`.
