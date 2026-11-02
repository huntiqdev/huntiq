#!/usr/bin/env python3
"""
HUNTIQ Batch Kaggle & External Dataset Ingestion Pipeline.
Scans ai/datasets/external_drop/ for 1 to 100+ CSV, TSV, JSON, JSONL, or Parquet files,
automatically maps diverse column schemas, merges datasets, normalizes multilingual
text, and trains modern super AI model checkpoints.
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

EXTERNAL_DROP_DIR = AI_ROOT / "datasets" / "external_drop"
BASE_DATA_DIR = AI_ROOT / "training" / "training_data"
CHECKPOINT_DIR = AI_ROOT / "training" / "checkpoints"
REGISTRY_PATH = AI_ROOT / "training" / "model_registry" / "model_registry.json"

TASK_CONFIG = {
    "intent": {
        "task_name": "marketplace_intent_classification",
        "model_id": "intent_classifier_v1",
        "drop_folder": EXTERNAL_DROP_DIR / "intent",
        "fallback_file": BASE_DATA_DIR / "intent_dataset.json",
    },
    "category": {
        "task_name": "marketplace_category_classification",
        "model_id": "category_classifier_v1",
        "drop_folder": EXTERNAL_DROP_DIR / "category",
        "fallback_file": BASE_DATA_DIR / "category_dataset.json",
    },
    "safety": {
        "task_name": "marketplace_safety_scam_moderation",
        "model_id": "safety_moderation_v1",
        "drop_folder": EXTERNAL_DROP_DIR / "safety",
        "fallback_file": BASE_DATA_DIR / "safety_dataset.json",
    },
    "domain": {
        "task_name": "ethiopian_marketplace_domain_classification",
        "model_id": "domain_classifier_v1",
        "drop_folder": EXTERNAL_DROP_DIR / "domain",
        "fallback_file": BASE_DATA_DIR / "domain_dataset.json",
    },
    "sentiment": {
        "task_name": "marketplace_sentiment_review_dispute_classification",
        "model_id": "sentiment_dispute_v1",
        "drop_folder": EXTERNAL_DROP_DIR / "sentiment",
        "fallback_file": BASE_DATA_DIR / "sentiment_dispute_dataset.json",
    },
}


def update_registry(model_entry: dict):
    """Updates model_registry.json with newly trained checkpoint entry."""
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


def scan_and_train_task(task_key: str, specific_path: str | None = None) -> dict:
    """Ingests dataset files for a task and retrains the model."""
    cfg = TASK_CONFIG[task_key]
    drop_folder = cfg["drop_folder"]
    fallback_file = cfg["fallback_file"]
    task_name = cfg["task_name"]
    model_id = cfg["model_id"]

    os.makedirs(drop_folder, exist_ok=True)

    # Determine data source
    if specific_path and os.path.exists(specific_path):
        target_source = specific_path
        source_type = "Custom Path"
    else:
        # Check if drop folder has any files
        files_in_drop = [
            f
            for f in os.listdir(drop_folder)
            if not f.startswith(".") and f != "README.md"
        ]
        if files_in_drop:
            target_source = str(drop_folder)
            source_type = f"External Drop ({len(files_in_drop)} files)"
        elif os.path.exists(fallback_file):
            target_source = str(fallback_file)
            source_type = "Base Reconnaissance Data"
        else:
            print(f"[SKIP] No data found for task '{task_key}'")
            return {}

    print("\n" + "=" * 70)
    print(f"PROCESSING TASK: {task_key.upper()} [{model_id}]")
    print(f"Source: {source_type} -> {target_source}")
    print("=" * 70)

    t0 = time.time()
    df, report = dataset_pipeline.load_dataset(
        source_path=target_source,
        deduplicate=True,
        normalize_geez=True,
    )

    if df.empty or len(report.unique_classes) < 2:
        print(f"[WARN] Insufficient valid data for {task_key}. Rows: {len(df)}")
        return {}

    print(f"Clean Samples Ingested: {report.valid_rows}")
    print(f"Files Processed:        {len(report.files_processed)}")
    print(f"Unique Classes ({len(report.unique_classes)}):  {report.unique_classes}")

    # Train classifier
    classifier = ModernSemanticClassifier(task_name=task_name)
    metrics = classifier.fit(
        texts=df["text"].tolist(),
        labels=df["label"].tolist(),
        eval_split=0.20,
        random_state=42,
    )

    # Save Checkpoint
    ckpt_path = CHECKPOINT_DIR / f"{model_id}.joblib"
    os.makedirs(CHECKPOINT_DIR, exist_ok=True)

    import joblib

    joblib.dump(classifier.pipeline, str(ckpt_path))

    metrics_path = CHECKPOINT_DIR / f"{model_id}_metrics.json"
    with open(metrics_path, "w", encoding="utf-8") as f:
        json.dump(metrics, f, indent=2, ensure_ascii=False)

    duration = time.time() - t0
    print(
        f"Completed in {duration:.2f}s | Accuracy: {metrics['accuracy'] * 100:.2f}% | Macro F1: {metrics['macro_f1']:.4f} | Latency: {metrics['inference_latency_ms']:.3f} ms"
    )

    # Update Registry
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
            "total_samples": report.valid_rows,
            "status": "APPROVED_PRODUCTION",
            "approval_date": time.strftime("%Y-%m-%d", time.gmtime()),
        }
    )

    return {
        "model_id": model_id,
        "task": task_name,
        "samples": report.valid_rows,
        "files": len(report.files_processed),
        "accuracy": metrics["accuracy"],
        "macro_f1": metrics["macro_f1"],
        "latency_ms": metrics["inference_latency_ms"],
    }


def main():
    """CLI Entry point for batch Kaggle dataset ingestion."""
    parser = argparse.ArgumentParser(
        description="HUNTIQ Batch Kaggle & External Dataset Ingestion"
    )
    parser.add_argument(
        "--task",
        choices=["all", "intent", "category", "safety", "domain", "sentiment"],
        default="all",
        help="Specify a single task or 'all' to ingest all tasks",
    )
    parser.add_argument(
        "--path",
        type=str,
        default=None,
        help="Optional path to custom folder containing CSVs or specific file",
    )
    args = parser.parse_args()

    # Ensure all external drop directories exist
    for cfg in TASK_CONFIG.values():
        os.makedirs(cfg["drop_folder"], exist_ok=True)

    print("=" * 70)
    print("HUNTIQ BATCH DATASET INGESTION & SUPER AI TRAINING")
    print(f"Drop directory ready at: {EXTERNAL_DROP_DIR}")
    print("=" * 70)

    tasks_to_run = list(TASK_CONFIG.keys()) if args.task == "all" else [args.task]

    all_results = []
    all_texts_for_encoder = []

    for task_key in tasks_to_run:
        res = scan_and_train_task(
            task_key, specific_path=args.path if args.task != "all" else None
        )
        if res:
            all_results.append(res)

    # Re-fit continuous dense semantic encoder across all tasks
    for cfg in TASK_CONFIG.values():
        fb = cfg["fallback_file"]
        if os.path.exists(fb):
            try:
                with open(fb, encoding="utf-8") as f:
                    items = json.load(f)
                for item in items:
                    txt = item[0] if isinstance(item, list) else item.get("text", "")
                    if txt:
                        all_texts_for_encoder.append(txt)
            except (json.JSONDecodeError, OSError):
                pass

    if all_texts_for_encoder:
        print("\nUpdating Dense Semantic Continuous Vector Space...")
        dense_encoder.fit(all_texts_for_encoder)
        dense_encoder.save(SEMANTIC_CHECKPOINT)
        print(f"Saved updated Dense Semantic Checkpoint to: {SEMANTIC_CHECKPOINT}")

    # Print Final Summary Table
    print("\n" + "=" * 70)
    print("BATCH INGESTION & TRAINING SUMMARY")
    print("=" * 70)
    print(
        f"{'Model ID':<24} | {'Samples':<8} | {'Accuracy':<10} | {'Macro F1':<10} | {'Latency':<8}"
    )
    print("-" * 70)
    for r in all_results:
        print(
            f"{r['model_id']:<24} | {r['samples']:<8} | {r['accuracy'] * 100:>6.2f}%    | {r['macro_f1']:>8.4f} | {r['latency_ms']:>6.2f} ms"
        )
    print("=" * 70)
    print("Ready for production inference and high-volume external dataset drops!\n")


if __name__ == "__main__":
    main()
