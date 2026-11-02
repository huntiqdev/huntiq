#!/usr/bin/env python3
"""
HUNTIQ Modern Super AI Universal Ingestion & Training Pipeline.
Supports:
- Ingesting single or 100+ CSV, TSV, JSON, JSONL, and Parquet files
- Letter-by-letter & morphological subword tokenization (Amharic, Oromo, Tigrinya, English)
- Dense continuous semantic latent space embedding
- Scalable training from small 100-sample sets to hundreds of thousands of samples
- Automated evaluation, accuracy benchmarking, and model registry synchronization
"""

from __future__ import annotations

import argparse
import json
import os
import sys
import time
from pathlib import Path

# Add project root to sys.path
SCRIPT_DIR = Path(__file__).resolve().parent
AI_ROOT = SCRIPT_DIR.parent
sys.path.insert(0, str(AI_ROOT))

from app.core.dataset_ingestion_pipeline import dataset_pipeline
from app.core.dense_semantic_encoder import SEMANTIC_CHECKPOINT, dense_encoder
from app.core.modern_semantic_classifier import ModernSemanticClassifier

DATA_DIR = AI_ROOT / "training" / "training_data"
CHECKPOINT_DIR = AI_ROOT / "training" / "checkpoints"
REGISTRY_PATH = AI_ROOT / "training" / "model_registry" / "model_registry.json"


def update_registry(model_entry: dict):
    """Updates model_registry.json with new trained model entry."""
    registry = {"registered_models": []}
    if os.path.exists(REGISTRY_PATH):
        try:
            with open(REGISTRY_PATH, encoding="utf-8") as f:
                registry = json.load(f)
        except Exception:  # noqa: BLE001
            registry = {"registered_models": []}

    models = [
        m
        for m in registry.get("registered_models", [])
        if m.get("model_id") != model_entry["model_id"]
    ]
    models.append(model_entry)
    registry["registered_models"] = models

    os.makedirs(os.path.dirname(REGISTRY_PATH), exist_ok=True)
    with open(REGISTRY_PATH, "w", encoding="utf-8") as f:
        json.dump(registry, f, indent=2, ensure_ascii=False)


def train_dense_semantic_encoder(all_corpus_texts: list[str]):
    """Fits and saves the dense continuous latent semantic encoder."""
    print("\n" + "=" * 70)
    print("TRAINING MULTILINGUAL DENSE SEMANTIC ENCODER (LATENT SPACE)")
    print("=" * 70)
    print(
        f"Total corpus size for semantic representation: {len(all_corpus_texts)} texts"
    )

    t0 = time.time()
    dense_encoder.fit(all_corpus_texts)
    dense_encoder.save(SEMANTIC_CHECKPOINT)
    duration = time.time() - t0

    print(
        f"Dense Semantic Encoder trained in {duration:.3f}s. Embedding Dim: {dense_encoder.embedding_dim}"
    )
    print(f"Saved checkpoint to: {SEMANTIC_CHECKPOINT}")

    # Test semantic similarity demonstration
    pairs = [
        ("Looking to buy a clean Toyota car in Bole", "ቶዮታ መኪና በቦሌ መግዛት እፈልጋለሁ"),
        ("iPhone 15 Pro Max 256GB phone", "ሳምሰንግ ጋላክሲ ስልክ የሚሸጥ"),
        (
            "Send 50% advance to Telebirr before seeing car",
            "Write a Python script for binary trees",
        ),
    ]
    print("\n--- Latent Semantic Cosine Similarity Test ---")
    for t1, t2 in pairs:
        sim = dense_encoder.similarity(t1, t2)
        print(f"Sim({t1[:35]}... <-> {t2[:35]}...): {sim:.4f}")


def train_single_task(
    task_name: str,
    model_id: str,
    source_path: str,
    text_col: str | None = None,
    label_col: str | None = None,
) -> tuple[float, float, float]:
    """Ingests dataset and trains a ModernSemanticClassifier."""
    print("\n" + "=" * 70)
    print(f"TRAINING MODERN SEMANTIC CLASSIFIER: {model_id} ({task_name})")
    print("=" * 70)

    df, report = dataset_pipeline.load_dataset(
        source_path=source_path,
        text_column=text_col,
        label_column=label_col,
        deduplicate=True,
        normalize_geez=True,
    )

    print(
        f"Ingested {report.valid_rows} clean samples from {len(report.files_processed)} files."
    )
    print(f"Classes ({len(report.unique_classes)}): {report.unique_classes}")

    classifier = ModernSemanticClassifier(task_name=task_name)
    metrics = classifier.fit(
        texts=df["text"].tolist(),
        labels=df["label"].tolist(),
        eval_split=0.20,
        random_state=42,
    )

    ckpt_path = CHECKPOINT_DIR / f"{model_id}.joblib"
    os.makedirs(CHECKPOINT_DIR, exist_ok=True)

    # Save pipeline in format compatible with joblib
    import joblib

    joblib.dump(classifier.pipeline, str(ckpt_path))

    metrics_path = CHECKPOINT_DIR / f"{model_id}_metrics.json"
    with open(metrics_path, "w", encoding="utf-8") as f:
        json.dump(metrics, f, indent=2, ensure_ascii=False)

    print(f"\n--- EVALUATION METRICS FOR {model_id} ---")
    print(f"Test Accuracy: {metrics['accuracy'] * 100:.2f}%")
    print(f"Macro F1:      {metrics['macro_f1']:.4f}")
    print(f"Latency:       {metrics['inference_latency_ms']:.3f} ms/sample")
    print(f"Saved Checkpoint: {ckpt_path}")

    # Register in Model Registry
    update_registry(
        {
            "model_id": model_id,
            "version": "2.0.0",
            "task": task_name,
            "artifact_path": f"training/checkpoints/{model_id}.joblib",
            "metrics_path": f"training/checkpoints/{model_id}_metrics.json",
            "accuracy": float(metrics["accuracy"]),
            "macro_f1": float(metrics["macro_f1"]),
            "latency_ms": float(metrics["inference_latency_ms"]),
            "status": "APPROVED_PRODUCTION",
            "approved_by": "HUNTIQ Modern Super AI Architectural Pipeline",
            "approval_date": time.strftime("%Y-%m-%d", time.gmtime()),
        }
    )

    return metrics["accuracy"], metrics["macro_f1"], metrics["inference_latency_ms"]


def run_full_pipeline(custom_source: str | None = None):
    """Executes end-to-end ingestion, semantic pre-training, and modern model training."""
    # 1. Collect all corpus texts for dense semantic encoder
    corpus_texts: list[str] = []

    tasks = [
        (
            "ethiopian_marketplace_domain_classification",
            "domain_classifier_v1",
            DATA_DIR / "domain_dataset.json",
        ),
        (
            "marketplace_safety_scam_moderation",
            "safety_moderation_v1",
            DATA_DIR / "safety_dataset.json",
        ),
        (
            "marketplace_category_classification",
            "category_classifier_v1",
            DATA_DIR / "category_dataset.json",
        ),
        (
            "marketplace_sentiment_review_dispute_classification",
            "sentiment_dispute_v1",
            DATA_DIR / "sentiment_dispute_dataset.json",
        ),
    ]

    for _, _, path in tasks:
        if os.path.exists(path):
            with open(path, encoding="utf-8") as f:
                data = json.load(f)
            for item in data:
                text = item[0] if isinstance(item, list) else item.get("text", "")
                if text:
                    corpus_texts.append(text)

    # Train dense semantic encoder
    if corpus_texts:
        train_dense_semantic_encoder(corpus_texts)

    # 2. Train each core modern classifier
    results = {}
    for task_name, model_id, path in tasks:
        acc, f1, lat = train_single_task(task_name, model_id, str(path))
        results[model_id] = {"accuracy": acc, "macro_f1": f1, "latency_ms": lat}

    # If custom dataset source provided (e.g. 100 CSVs or Parquet)
    if custom_source:
        print(f"\nProcessing Custom Dataset Source: {custom_source}")
        train_single_task(
            task_name="custom_dataset_classification",
            model_id="custom_model_v1",
            source_path=custom_source,
        )

    print("\n" + "=" * 70)
    print("ALL MODERN SUPER AI PIPELINES TRAINED & BENCHMARKED SUCCESSFULLY")
    print("=" * 70)
    for mid, r in results.items():
        print(
            f"- {mid:25s}: Accuracy: {r['accuracy'] * 100:.2f}% | Macro F1: {r['macro_f1']:.4f} | Latency: {r['latency_ms']:.3f} ms"
        )


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="HUNTIQ Modern Super AI Ingestion & Training CLI"
    )
    parser.add_argument(
        "--data-path",
        type=str,
        default=None,
        help="Path to single CSV, directory of CSVs, JSON, or Parquet files",
    )
    args = parser.parse_args()

    run_full_pipeline(custom_source=args.data_path)
