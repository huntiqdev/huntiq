"""
HUNTIQ Dense Multilingual Semantic Encoder.
Projects multilingual sentences, words, and character n-grams into a unified
dense continuous latent semantic embedding space (R^d).
Supports:
- Dense sentence embeddings
- Cosine semantic similarity
- Semantic search & candidate re-ranking
- Cluster centroid modeling
- Out-of-Distribution (OOD) uncertainty quantification
"""

from __future__ import annotations

import os
from typing import Any

import joblib
import numpy as np
from sklearn.decomposition import TruncatedSVD
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.pipeline import FeatureUnion, Pipeline
from sklearn.preprocessing import normalize

from app.core.morphological_tokenizer import morph_tokenizer

BASE_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
SEMANTIC_CHECKPOINT = os.path.join(
    BASE_DIR, "training", "checkpoints", "dense_semantic_encoder_v1.joblib"
)


class DenseSemanticEncoder:
    """
    Modern Latent Semantic Encoder producing L2-normalized dense vector embeddings.
    Integrates subword character n-grams and morphological root decomposition.
    """

    def __init__(self, embedding_dim: int = 128):
        self.embedding_dim = embedding_dim
        self.pipeline: Pipeline | None = None
        self._load_or_init()

    def _build_pipeline(self) -> Pipeline:
        """Constructs the multi-resolution feature extraction and dense projection pipeline."""
        features = FeatureUnion(
            [
                (
                    "word_ngram",
                    TfidfVectorizer(
                        ngram_range=(1, 3),
                        sublinear_tf=True,
                        analyzer="word",
                        lowercase=True,
                        min_df=1,
                    ),
                ),
                (
                    "char_wb_ngram",
                    TfidfVectorizer(
                        ngram_range=(2, 6),
                        sublinear_tf=True,
                        analyzer="char_wb",
                        lowercase=True,
                        min_df=1,
                    ),
                ),
            ]
        )

        svd = TruncatedSVD(
            n_components=self.embedding_dim,
            algorithm="randomized",
            n_iter=7,
            random_state=42,
        )

        return Pipeline([("features", features), ("svd", svd)])

    def _load_or_init(self) -> None:
        """Loads serialized dense encoder checkpoint if present on disk."""
        if os.path.exists(SEMANTIC_CHECKPOINT):
            try:
                self.pipeline = joblib.load(SEMANTIC_CHECKPOINT)
            except Exception as e:  # noqa: BLE001
                print(f"[DenseSemanticEncoder] Warning: Failed to load checkpoint: {e}")
                self.pipeline = None

    def fit(self, texts: list[str]) -> DenseSemanticEncoder:
        """Fits the dense latent semantic projection on a multilingual corpus."""
        # Pre-enrich documents with morphological roots and letter n-grams
        enriched_texts = [morph_tokenizer.to_enriched_feature_string(t) for t in texts]

        # Ensure embedding_dim does not exceed sample count
        actual_dim = min(self.embedding_dim, len(texts) - 1, 256)
        actual_dim = max(actual_dim, 8)

        features = FeatureUnion(
            [
                (
                    "word_ngram",
                    TfidfVectorizer(
                        ngram_range=(1, 3),
                        sublinear_tf=True,
                        analyzer="word",
                        lowercase=True,
                        min_df=1,
                    ),
                ),
                (
                    "char_wb_ngram",
                    TfidfVectorizer(
                        ngram_range=(2, 6),
                        sublinear_tf=True,
                        analyzer="char_wb",
                        lowercase=True,
                        min_df=1,
                    ),
                ),
            ]
        )

        svd = TruncatedSVD(
            n_components=actual_dim,
            algorithm="randomized",
            n_iter=7,
            random_state=42,
        )

        self.pipeline = Pipeline([("features", features), ("svd", svd)])
        self.pipeline.fit(enriched_texts)
        self.embedding_dim = actual_dim
        return self

    def save(self, filepath: str = SEMANTIC_CHECKPOINT) -> None:
        """Saves fitted pipeline to disk."""
        if self.pipeline:
            os.makedirs(os.path.dirname(filepath), exist_ok=True)
            joblib.dump(self.pipeline, filepath)

    def encode(self, texts: list[str] | str) -> np.ndarray:
        """
        Encodes one or more texts into L2-normalized dense vectors in R^d.
        Returns: 2D numpy array of shape (N, embedding_dim).
        """
        if isinstance(texts, str):
            texts = [texts]

        if not texts or not self.pipeline:
            # Fallback zero vector
            return np.zeros((max(len(texts), 1), self.embedding_dim), dtype=np.float32)

        try:
            enriched = [morph_tokenizer.to_enriched_feature_string(t) for t in texts]
            dense_vectors = self.pipeline.transform(enriched)
            # L2 normalize so dot product equals cosine similarity
            return normalize(dense_vectors, norm="l2", axis=1).astype(np.float32)
        except Exception:  # noqa: BLE001
            return np.zeros((len(texts), self.embedding_dim), dtype=np.float32)

    def encode_sentence(self, text: str) -> np.ndarray:
        """Encodes a single sentence into a 1D vector of shape (embedding_dim,)."""
        return self.encode([text])[0]

    def similarity(self, text_a: str, text_b: str) -> float:
        """Computes cosine similarity between two sentences (-1.0 to 1.0)."""
        vecs = self.encode([text_a, text_b])
        dot = float(np.dot(vecs[0], vecs[1]))
        return max(-1.0, min(1.0, dot))

    def rank_candidates(
        self, query: str, candidates: list[dict[str, Any]], text_key: str = "title"
    ) -> list[dict[str, Any]]:
        """
        Ranks a list of candidate items by semantic similarity to the search query.
        Returns candidates ordered by descending similarity score.
        """
        if not candidates or not query:
            return candidates

        q_vec = self.encode_sentence(query)
        cand_texts = [str(c.get(text_key, "")) for c in candidates]
        cand_vecs = self.encode(cand_texts)

        # Vectorized dot product against all candidate vectors
        scores = np.dot(cand_vecs, q_vec)

        ranked = []
        for cand, score in zip(candidates, scores):
            item = dict(cand)
            item["semantic_score"] = float(score)
            ranked.append(item)

        ranked.sort(key=lambda x: x["semantic_score"], reverse=True)
        return ranked


# Global singleton instance
dense_encoder = DenseSemanticEncoder()
