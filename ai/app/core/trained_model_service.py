"""
HUNTIQ AI Runtime Service for Trained Scikit-Learn Models.
Loads trained checkpoints:
- intent_classifier_v1: Multilingual Intent & Domain Classification
- category_classifier_v1: Marketplace 8-Taxonomy Category Routing
- safety_moderation_v1: Fraud, Scam, Phishing & Harassment Detection
- domain_classifier_v1: Ethiopian Marketplace vs General Domain Boundary Guard
to provide high-speed (<1ms) inference with calibrated probabilities.
"""

import os

import joblib

BASE_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
CHECKPOINT_DIR = os.path.join(BASE_DIR, "training", "checkpoints")

INTENT_MODEL_PATH = os.path.join(CHECKPOINT_DIR, "intent_classifier_v1.joblib")
CATEGORY_MODEL_PATH = os.path.join(CHECKPOINT_DIR, "category_classifier_v1.joblib")
SAFETY_MODEL_PATH = os.path.join(CHECKPOINT_DIR, "safety_moderation_v1.joblib")
DOMAIN_MODEL_PATH = os.path.join(CHECKPOINT_DIR, "domain_classifier_v1.joblib")
SENTIMENT_MODEL_PATH = os.path.join(CHECKPOINT_DIR, "sentiment_dispute_v1.joblib")


class TrainedModelService:
    def __init__(self):
        self.intent_pipeline = None
        self.category_pipeline = None
        self.safety_pipeline = None
        self.domain_pipeline = None
        self.sentiment_pipeline = None
        self._load_models()

    def _load_models(self):
        """Loads serialized pipelines if available on disk."""
        if os.path.exists(INTENT_MODEL_PATH):
            try:
                self.intent_pipeline = joblib.load(INTENT_MODEL_PATH)
            except Exception as e:  # noqa: BLE001
                print(
                    f"[TrainedModelService] Warning: Could not load intent model: {e}"
                )

        if os.path.exists(CATEGORY_MODEL_PATH):
            try:
                self.category_pipeline = joblib.load(CATEGORY_MODEL_PATH)
            except Exception as e:  # noqa: BLE001
                print(
                    f"[TrainedModelService] Warning: Could not load category model: {e}"
                )

        if os.path.exists(SAFETY_MODEL_PATH):
            try:
                self.safety_pipeline = joblib.load(SAFETY_MODEL_PATH)
            except Exception as e:  # noqa: BLE001
                print(
                    f"[TrainedModelService] Warning: Could not load safety model: {e}"
                )

        if os.path.exists(DOMAIN_MODEL_PATH):
            try:
                self.domain_pipeline = joblib.load(DOMAIN_MODEL_PATH)
            except Exception as e:  # noqa: BLE001
                print(
                    f"[TrainedModelService] Warning: Could not load domain model: {e}"
                )

        if os.path.exists(SENTIMENT_MODEL_PATH):
            try:
                self.sentiment_pipeline = joblib.load(SENTIMENT_MODEL_PATH)
            except Exception as e:  # noqa: BLE001
                print(
                    f"[TrainedModelService] Warning: Could not load sentiment model: {e}"
                )

    def is_intent_model_ready(self) -> bool:
        return self.intent_pipeline is not None

    def is_category_model_ready(self) -> bool:
        return self.category_pipeline is not None

    def is_safety_model_ready(self) -> bool:
        return self.safety_pipeline is not None

    def is_domain_model_ready(self) -> bool:
        return self.domain_pipeline is not None

    def is_sentiment_model_ready(self) -> bool:
        return self.sentiment_pipeline is not None

    def predict_intent(self, text: str) -> tuple[str | None, float, dict[str, float]]:
        """
        Predicts intent label with calibrated probability.
        Returns: (best_label, confidence, all_probabilities)
        """
        if not self.intent_pipeline:
            return None, 0.0, {}

        try:
            probs = self.intent_pipeline.predict_proba([text])[0]
            classes = self.intent_pipeline.classes_
            prob_dict = {cls: float(prob) for cls, prob in zip(classes, probs)}
            best_idx = probs.argmax()
            best_label = classes[best_idx]
            confidence = float(probs[best_idx])
            return best_label, confidence, prob_dict
        except Exception:  # noqa: BLE001
            return None, 0.0, {}

    def predict_category(self, text: str) -> tuple[str | None, float, dict[str, float]]:
        """
        Predicts marketplace category with calibrated probability.
        Returns: (best_category, confidence, all_probabilities)
        """
        if not self.category_pipeline:
            return None, 0.0, {}

        try:
            probs = self.category_pipeline.predict_proba([text])[0]
            classes = self.category_pipeline.classes_
            prob_dict = {cls: float(prob) for cls, prob in zip(classes, probs)}
            best_idx = probs.argmax()
            best_category = classes[best_idx]
            confidence = float(probs[best_idx])
            return best_category, confidence, prob_dict
        except Exception:  # noqa: BLE001
            return None, 0.0, {}

    def predict_safety(self, text: str) -> tuple[str | None, float, dict[str, float]]:
        """
        Predicts safety/moderation label with calibrated probability.
        Returns: (safety_class, confidence, all_probabilities)
        safety_class is one of: safe, scam_advance, phishing_harvesting, prohibited_contraband, harassment_abuse
        """
        if not self.safety_pipeline:
            return "safe", 1.0, {"safe": 1.0}

        try:
            probs = self.safety_pipeline.predict_proba([text])[0]
            classes = self.safety_pipeline.classes_
            prob_dict = {cls: float(prob) for cls, prob in zip(classes, probs)}
            best_idx = probs.argmax()
            best_class = classes[best_idx]
            confidence = float(probs[best_idx])
            return best_class, confidence, prob_dict
        except Exception:  # noqa: BLE001
            return "safe", 1.0, {"safe": 1.0}

    def predict_domain(self, text: str) -> tuple[str | None, float, dict[str, float]]:
        """
        Predicts domain boundary: ethiopian_marketplace vs general_out_of_domain.
        Returns: (domain_label, confidence, all_probabilities)
        """
        if not self.domain_pipeline:
            return None, 0.0, {}

        try:
            probs = self.domain_pipeline.predict_proba([text])[0]
            classes = self.domain_pipeline.classes_
            prob_dict = {cls: float(prob) for cls, prob in zip(classes, probs)}
            best_idx = probs.argmax()
            best_domain = classes[best_idx]
            confidence = float(probs[best_idx])
            return best_domain, confidence, prob_dict
        except Exception:  # noqa: BLE001
            return None, 0.0, {}

    def predict_sentiment_dispute(
        self, text: str
    ) -> tuple[str | None, float, dict[str, float]]:
        """
        Predicts sentiment and dispute risk level with calibrated probability.
        Classes: positive_review, neutral_inquiry, negative_dissatisfied, transaction_dispute
        Returns: (best_class, confidence, all_probabilities)
        """
        if not self.sentiment_pipeline:
            return "neutral_inquiry", 1.0, {"neutral_inquiry": 1.0}

        try:
            probs = self.sentiment_pipeline.predict_proba([text])[0]
            classes = self.sentiment_pipeline.classes_
            prob_dict = {cls: float(prob) for cls, prob in zip(classes, probs)}
            best_idx = probs.argmax()
            best_class = classes[best_idx]
            confidence = float(probs[best_idx])
            return best_class, confidence, prob_dict
        except Exception:  # noqa: BLE001
            return "neutral_inquiry", 1.0, {"neutral_inquiry": 1.0}


# Global singleton instance
trained_models = TrainedModelService()
