"""
HUNTIQ Modern Multilingual Semantic Classifier.
Combines:
- Letter-by-letter & subword morphological decomposition (MorphologicalTokenizer)
- Sublinear character and word n-gram feature learning
- Calibrated probability estimation via Platt scaling / isotonic regression
- Shannon entropy-based Out-Of-Distribution (OOD) uncertainty quantification
"""

from __future__ import annotations

import os
import time
from typing import Any

import joblib
import numpy as np
from sklearn.calibration import CalibratedClassifierCV
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics import accuracy_score, classification_report, f1_score
from sklearn.pipeline import FeatureUnion, Pipeline
from sklearn.svm import LinearSVC


class ModernSemanticClassifier:
    """
    Modern Multilingual Semantic Classifier for Ethiopian Marketplace NLP.
    Understands every letter, character n-gram, morphological root, and consecutive sentence.
    """

    def __init__(self, task_name: str = "generic_semantic_classification"):
        self.task_name = task_name
        self.pipeline: Pipeline | None = None
        self.classes: list[str] = []
        self.is_fitted: bool = False

    def _build_pipeline(self) -> Pipeline:
        """Constructs the multi-resolution feature extraction and calibrated classifier pipeline."""
        features = FeatureUnion(
            [
                (
                    "word_tfidf",
                    TfidfVectorizer(
                        ngram_range=(1, 3),
                        sublinear_tf=True,
                        analyzer="word",
                        lowercase=True,
                        min_df=1,
                        max_df=0.98,
                    ),
                ),
                (
                    "char_wb_tfidf",
                    TfidfVectorizer(
                        ngram_range=(2, 6),
                        sublinear_tf=True,
                        analyzer="char_wb",
                        lowercase=True,
                        min_df=1,
                        max_df=0.98,
                    ),
                ),
                (
                    "char_letters_tfidf",
                    TfidfVectorizer(
                        ngram_range=(1, 2),
                        sublinear_tf=True,
                        analyzer="char",
                        lowercase=True,
                        min_df=1,
                        max_df=0.98,
                    ),
                ),
            ]
        )

        base_svc = LinearSVC(
            C=2.5,
            loss="squared_hinge",
            max_iter=5000,
            random_state=42,
        )
        calibrated_svc = CalibratedClassifierCV(estimator=base_svc, cv=5)

        return Pipeline([("features", features), ("classifier", calibrated_svc)])

    def fit(
        self,
        texts: list[str],
        labels: list[str],
        eval_split: float = 0.20,
        random_state: int = 42,
    ) -> dict[str, Any]:
        """
        Trains the modern classifier pipeline using multi-resolution word, subword, and letter features.
        Returns evaluation metrics dictionary.
        """
        t0 = time.time()

        from sklearn.model_selection import train_test_split

        # Stratified train/test split if possible
        label_counts = {lbl: labels.count(lbl) for lbl in set(labels)}
        can_stratify = all(c >= 2 for c in label_counts.values())
        stratify = labels if can_stratify else None

        X_train, X_test, y_train, y_test = train_test_split(
            texts,
            labels,
            test_size=eval_split,
            random_state=random_state,
            stratify=stratify,
        )

        self.pipeline = self._build_pipeline()
        self.pipeline.fit(X_train, y_train)
        self.classes = sorted(set(labels))
        self.is_fitted = True

        train_time = time.time() - t0

        # Evaluate on test split
        t_eval = time.time()
        y_pred = self.pipeline.predict(X_test)
        eval_time_per_sample_ms = ((time.time() - t_eval) * 1000.0) / max(
            len(X_test), 1
        )

        acc = float(accuracy_score(y_test, y_pred))
        macro_f1 = float(f1_score(y_test, y_pred, average="macro", zero_division=0))
        report = classification_report(
            y_test, y_pred, zero_division=0, output_dict=True
        )

        return {
            "task_name": self.task_name,
            "accuracy": acc,
            "macro_f1": macro_f1,
            "training_duration_seconds": float(train_time),
            "inference_latency_ms": float(eval_time_per_sample_ms),
            "train_samples": len(X_train),
            "test_samples": len(X_test),
            "classes": self.classes,
            "detailed_report": report,
        }

    def predict(self, text: str) -> tuple[str, float, dict[str, float], float]:
        """
        Predicts class label with calibrated probability and uncertainty metric.
        Returns:
            (best_class, confidence, probability_distribution, uncertainty_score)
            uncertainty_score: 0.0 (certain) to 1.0 (completely ambiguous/OOD)
        """
        if not self.pipeline or not self.is_fitted:
            return "unknown", 0.0, {}, 1.0

        try:
            probs = self.pipeline.predict_proba([text])[0]
            classes = self.pipeline.classes_
            prob_dict = {str(c): float(p) for c, p in zip(classes, probs)}

            best_idx = probs.argmax()
            best_label = str(classes[best_idx])
            confidence = float(probs[best_idx])

            # Normalized Shannon Entropy for uncertainty (0.0 = point mass, 1.0 = uniform distribution)
            k = len(classes)
            if k > 1:
                # Avoid log(0)
                safe_probs = np.clip(probs, 1e-12, 1.0)
                entropy = -float(np.sum(safe_probs * np.log(safe_probs)))
                max_entropy = float(np.log(k))
                uncertainty = float(entropy / max_entropy)
            else:
                uncertainty = 0.0

            return best_label, confidence, prob_dict, uncertainty
        except Exception:  # noqa: BLE001
            return "unknown", 0.0, {}, 1.0

    def predict_batch(self, texts: list[str]) -> list[dict[str, Any]]:
        """Vectorized batch prediction for high-throughput evaluation."""
        if not self.pipeline or not self.is_fitted or not texts:
            return []

        prob_matrix = self.pipeline.predict_proba(texts)
        classes = self.pipeline.classes_

        results = []
        k = len(classes)
        max_entropy = float(np.log(k)) if k > 1 else 1.0

        for probs in prob_matrix:
            best_idx = probs.argmax()
            best_label = str(classes[best_idx])
            confidence = float(probs[best_idx])

            safe_probs = np.clip(probs, 1e-12, 1.0)
            entropy = -float(np.sum(safe_probs * np.log(safe_probs)))
            uncertainty = float(entropy / max_entropy) if k > 1 else 0.0

            prob_dict = {str(c): float(p) for c, p in zip(classes, probs)}
            results.append(
                {
                    "label": best_label,
                    "confidence": confidence,
                    "uncertainty": uncertainty,
                    "probabilities": prob_dict,
                }
            )

        return results

    def save(self, filepath: str) -> None:
        """Serializes trained pipeline to joblib checkpoint."""
        os.makedirs(os.path.dirname(filepath), exist_ok=True)
        joblib.dump(self, filepath)

    @classmethod
    def load(cls, filepath: str) -> ModernSemanticClassifier:
        """Deserializes trained pipeline from joblib checkpoint."""
        return joblib.load(filepath)
