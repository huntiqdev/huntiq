"""
HUNTIQ Enterprise Dataset Ingestion & Streaming Pipeline.
Scalable multi-format ingestion engine capable of reading:
- Individual or directories of 100+ CSV / TSV files
- JSON & JSONL datasets
- High-throughput Parquet files (scaling to hundreds of thousands of rows)
Provides:
- Auto-discovery of text & label column schemas
- Unicode normalization and multilingual cleaning
- Memory-efficient chunked streaming
- Stratified train/test splitting and validation reports
"""

from __future__ import annotations

import glob
import json
import os
from collections.abc import Generator
from pathlib import Path
from typing import Any

import pandas as pd
from sklearn.model_selection import train_test_split

from app.languages.normalizer import normalizer

TEXT_COLUMN_CANDIDATES = [
    "text",
    "query",
    "sentence",
    "message",
    "utterance",
    "prompt",
    "input",
    "body",
    "title",
    "description",
    "content",
    "review",
    "statement",
]

LABEL_COLUMN_CANDIDATES = [
    "label",
    "intent",
    "category",
    "target",
    "class",
    "sentiment",
    "rating",
    "ground_truth",
    "domain",
    "type",
    "tag",
]


class IngestionReport:
    """Holds ingestion statistics and data quality metrics."""

    def __init__(self):
        self.total_rows_read: int = 0
        self.valid_rows: int = 0
        self.dropped_empty: int = 0
        self.duplicates_removed: int = 0
        self.files_processed: list[str] = []
        self.detected_text_column: str | None = None
        self.detected_label_column: str | None = None
        self.unique_classes: list[str] = []
        self.class_distribution: dict[str, int] = {}

    def to_dict(self) -> dict[str, Any]:
        return {
            "total_rows_read": self.total_rows_read,
            "valid_rows": self.valid_rows,
            "dropped_empty": self.dropped_empty,
            "duplicates_removed": self.duplicates_removed,
            "file_count": len(self.files_processed),
            "files": self.files_processed,
            "detected_text_column": self.detected_text_column,
            "detected_label_column": self.detected_label_column,
            "unique_classes": self.unique_classes,
            "class_distribution": self.class_distribution,
        }


class DatasetIngestionPipeline:
    """
    Universal Dataset Ingestion Pipeline for HUNTIQ AI.
    Handles from small 100-sample files up to hundreds of thousands of rows.
    """

    def __init__(self, chunk_size: int = 50000):
        self.chunk_size = chunk_size

    @staticmethod
    def detect_columns(columns: list[str]) -> tuple[str | None, str | None]:
        """Auto-detects the text column and label column from a list of header names."""
        lower_cols = {c.lower().strip(): c for c in columns}

        text_col = None
        # 1. Exact match for text
        for candidate in TEXT_COLUMN_CANDIDATES:
            if candidate in lower_cols:
                text_col = lower_cols[candidate]
                break

        # 2. Substring match for text
        if not text_col:
            for candidate in TEXT_COLUMN_CANDIDATES:
                for l_col, orig_col in lower_cols.items():
                    if candidate in l_col:
                        text_col = orig_col
                        break
                if text_col:
                    break

        label_col = None
        # 1. Exact match for label
        for candidate in LABEL_COLUMN_CANDIDATES:
            if candidate in lower_cols and lower_cols[candidate] != text_col:
                label_col = lower_cols[candidate]
                break

        # 2. Substring match for label
        if not label_col:
            for candidate in LABEL_COLUMN_CANDIDATES:
                for l_col, orig_col in lower_cols.items():
                    if candidate in l_col and orig_col != text_col:
                        label_col = orig_col
                        break
                if label_col:
                    break

        return text_col, label_col

    def _resolve_paths(self, source_path: str | Path) -> list[str]:
        """Resolves single file, directory, or glob pattern to list of absolute filepaths."""
        path_str = str(source_path)
        if os.path.isdir(path_str):
            # Check for CSVs, JSONs, Parquets inside directory
            patterns = ["*.csv", "*.tsv", "*.json", "*.jsonl", "*.parquet"]
            files: list[str] = []
            for p in patterns:
                files.extend(glob.glob(os.path.join(path_str, p)))
                files.extend(glob.glob(os.path.join(path_str, "**", p), recursive=True))
            return sorted(set(files))
        elif "*" in path_str or "?" in path_str:
            return sorted(glob.glob(path_str))
        elif os.path.exists(path_str):
            return [path_str]
        else:
            return []

    def load_dataset(
        self,
        source_path: str | Path,
        text_column: str | None = None,
        label_column: str | None = None,
        deduplicate: bool = True,
        normalize_geez: bool = True,
    ) -> tuple[pd.DataFrame, IngestionReport]:
        """
        Loads and standardizes datasets from single or multiple CSVs, JSON, or Parquet files.
        Returns:
            DataFrame with normalized columns ['text', 'label'] and an IngestionReport.
        """
        report = IngestionReport()
        files = self._resolve_paths(source_path)

        if not files:
            raise FileNotFoundError(
                f"No valid dataset files found at path: {source_path}"
            )

        dfs: list[pd.DataFrame] = []

        for fpath in files:
            ext = os.path.splitext(fpath)[1].lower()
            try:
                if ext in [".csv", ".tsv"]:
                    sep = "\t" if ext == ".tsv" else ","
                    # Read using pandas
                    df_chunk = pd.read_csv(fpath, sep=sep, low_memory=False)
                elif ext == ".jsonl":
                    df_chunk = pd.read_json(fpath, lines=True)
                elif ext == ".json":
                    with open(fpath, encoding="utf-8") as f:
                        data = json.load(f)
                    if isinstance(data, list):
                        if data and isinstance(data[0], list) and len(data[0]) >= 2:
                            # Handle [[text, label], ...]
                            df_chunk = pd.DataFrame(data, columns=["text", "label"])
                        else:
                            df_chunk = pd.DataFrame(data)
                    elif isinstance(data, dict):
                        # Handle dict with records
                        for key in ["samples", "data", "items", "dataset", "examples"]:
                            if key in data and isinstance(data[key], list):
                                df_chunk = pd.DataFrame(data[key])
                                break
                        else:
                            df_chunk = pd.DataFrame([data])
                    else:
                        continue
                elif ext == ".parquet":
                    df_chunk = pd.read_parquet(fpath)
                else:
                    continue

                if not df_chunk.empty:
                    report.files_processed.append(fpath)
                    report.total_rows_read += len(df_chunk)
                    dfs.append(df_chunk)
            except Exception as e:  # noqa: BLE001
                print(f"[DatasetIngestionPipeline] Warning: Skipping {fpath}: {e}")

        if not dfs:
            raise ValueError(f"Could not read any valid data rows from: {files}")

        merged_df = pd.concat(dfs, ignore_index=True)

        # Detect columns if not explicitly supplied
        detected_text, detected_label = self.detect_columns(merged_df.columns.tolist())
        final_text_col = text_column or detected_text
        final_label_col = label_column or detected_label

        if not final_text_col:
            raise ValueError(
                f"Could not automatically detect text column in columns: {merged_df.columns.tolist()}. "
                f"Please specify text_column explicitly."
            )

        report.detected_text_column = final_text_col
        report.detected_label_column = final_label_col

        # Standardize DataFrame
        standard_df = pd.DataFrame()
        standard_df["text"] = merged_df[final_text_col].astype(str)

        if final_label_col and final_label_col in merged_df.columns:
            standard_df["label"] = merged_df[final_label_col].astype(str)
        else:
            standard_df["label"] = "unlabeled"

        # Cleaning: Remove empty / whitespace strings
        initial_len = len(standard_df)
        standard_df = standard_df[standard_df["text"].str.strip() != ""]
        standard_df = standard_df[standard_df["text"].str.strip() != "nan"]
        report.dropped_empty = initial_len - len(standard_df)

        # Ge'ez Unicode Normalization
        if normalize_geez:
            standard_df["text"] = standard_df["text"].apply(
                lambda t: normalizer.normalize(t)
            )

        # Deduplication
        if deduplicate:
            before_dedup = len(standard_df)
            standard_df = standard_df.drop_duplicates(subset=["text", "label"])
            report.duplicates_removed = before_dedup - len(standard_df)

        report.valid_rows = len(standard_df)
        report.unique_classes = sorted(standard_df["label"].unique().tolist())
        report.class_distribution = standard_df["label"].value_counts().to_dict()

        return standard_df, report

    def stream_batches(
        self,
        df: pd.DataFrame,
        batch_size: int = 1000,
    ) -> Generator[tuple[list[str], list[str]], None, None]:
        """Memory-efficient batch generator for training very large datasets."""
        total = len(df)
        for i in range(0, total, batch_size):
            chunk = df.iloc[i : min(i + batch_size, total)]
            yield chunk["text"].tolist(), chunk["label"].tolist()

    @staticmethod
    def stratified_split(
        df: pd.DataFrame, test_size: float = 0.20, random_state: int = 42
    ) -> tuple[pd.DataFrame, pd.DataFrame]:
        """Splits DataFrame into train and test sets with stratified class balance."""
        df["label"].tolist()
        class_counts = df["label"].value_counts()

        # If any class has only 1 sample, fallback to unstratified split
        can_stratify = all(class_counts >= 2)

        stratify = df["label"] if can_stratify else None
        train_df, test_df = train_test_split(
            df, test_size=test_size, random_state=random_state, stratify=stratify
        )
        return train_df.reset_index(drop=True), test_df.reset_index(drop=True)

    @staticmethod
    def export_parquet(df: pd.DataFrame, output_path: str | Path) -> str:
        """Exports DataFrame to optimized Parquet format."""
        out = Path(output_path)
        out.parent.mkdir(parents=True, exist_ok=True)
        df.to_parquet(str(out), index=False, compression="snappy")
        return str(out)


# Global singleton instance
dataset_pipeline = DatasetIngestionPipeline()
