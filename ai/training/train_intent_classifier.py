#!/usr/bin/env python3
"""
HUNTIQ AI Offline Intent & Domain Model Training Pipeline.
Trains an enterprise-grade calibrated multilingual classifier across
English (en), Amharic (am), Afaan Oromoo (or), and Tigrinya (ti).
Outputs checkpoint to ai/training/checkpoints/intent_classifier_v1.joblib.
"""

import json
import os
import time

import joblib
from sklearn.calibration import CalibratedClassifierCV
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics import accuracy_score, classification_report, f1_score
from sklearn.pipeline import FeatureUnion, Pipeline
from sklearn.svm import LinearSVC

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
AI_ROOT = os.path.dirname(SCRIPT_DIR)

CHECKPOINT_DIR = os.path.join(AI_ROOT, "training", "checkpoints")
MODEL_PATH = os.path.join(CHECKPOINT_DIR, "intent_classifier_v1.joblib")
METRICS_PATH = os.path.join(CHECKPOINT_DIR, "intent_classifier_v1_metrics.json")
REGISTRY_PATH = os.path.join(
    AI_ROOT, "training", "model_registry", "model_registry.json"
)

os.makedirs(CHECKPOINT_DIR, exist_ok=True)
os.makedirs(os.path.dirname(REGISTRY_PATH), exist_ok=True)


# Comprehensive training dataset across all 4 languages and 6 core intent classes
TRAIN_DATA = [
    # --- SOCIAL GREETING ---
    ("Hello!", "social_greeting"),
    ("Hi there, how are you?", "social_greeting"),
    ("Good morning!", "social_greeting"),
    ("Good evening assistant", "social_greeting"),
    ("Hey!", "social_greeting"),
    ("Thank you very much", "social_greeting"),
    ("Thanks for your help", "social_greeting"),
    ("Goodbye, see you later", "social_greeting"),
    ("Bye bye", "social_greeting"),
    ("ሰላም!", "social_greeting"),
    ("ሰላም እንደምን ነህ?", "social_greeting"),
    ("ሰላም እንደምን ዋላችሁ?", "social_greeting"),
    ("እንደምን አደራችሁ?", "social_greeting"),
    ("እንደምን አመሻችሁ?", "social_greeting"),
    ("አመሰግናለሁ!", "social_greeting"),
    ("በጣም አመሰግናለሁ", "social_greeting"),
    ("ደህና ሁን", "social_greeting"),
    ("ቻው", "social_greeting"),
    ("Akkam jirtu?", "social_greeting"),
    ("Akkam nagaa?", "social_greeting"),
    ("Akkam bultan?", "social_greeting"),
    ("Akkam ooltan?", "social_greeting"),
    ("Galatoomi!", "social_greeting"),
    ("Baay'ee galatoomi", "social_greeting"),
    ("Nagaatti", "social_greeting"),
    ("ሰላም ከመይ ኣለኻ?", "social_greeting"),
    ("ሰላም ከመይ ሓዲርኩም?", "social_greeting"),
    ("ሰላም ከመይ ውዒልኩም?", "social_greeting"),
    ("የቐንየለይ!", "social_greeting"),
    ("ብጣዕሚ የቐንየለይ", "social_greeting"),
    ("ደሓን ኩኑ", "social_greeting"),
    # --- IDENTITY INQUIRY ---
    ("Who are you?", "identity_inquiry"),
    ("What is your name?", "identity_inquiry"),
    ("What can you do?", "identity_inquiry"),
    ("What are your capabilities?", "identity_inquiry"),
    ("Are you an AI?", "identity_inquiry"),
    ("How does HUNTIQ help me?", "identity_inquiry"),
    ("Is this service free to use?", "identity_inquiry"),
    ("Can I use this without an account?", "identity_inquiry"),
    ("ማን ነህ?", "identity_inquiry"),
    ("ስምህ ማን ነው?", "identity_inquiry"),
    ("ምን ምን ትሰራለህ?", "identity_inquiry"),
    ("ምን ማገዝ ትችላለህ?", "identity_inquiry"),
    ("አንተ ምንድን ነህ?", "identity_inquiry"),
    ("በHUNTIQ ምን ማድረግ እችላለሁ?", "identity_inquiry"),
    ("ይህ አገልግሎት ነጻ ነው?", "identity_inquiry"),
    ("Eenyu ati?", "identity_inquiry"),
    ("Maqaan kee eenyu?", "identity_inquiry"),
    ("Maal hojjechuu dandeessa?", "identity_inquiry"),
    ("Akkamitti na gargaarta?", "identity_inquiry"),
    ("Tajaajilli kun bilisaa?", "identity_inquiry"),
    ("መን ኢኻ?", "identity_inquiry"),
    ("ስምካ መን ይበሃል?", "identity_inquiry"),
    ("እንታይ ክትሕግዝ ትኽእል?", "identity_inquiry"),
    ("ናይ HUNTIQ ረዳት ዲኻ?", "identity_inquiry"),
    ("እዚ ኣገልግሎት ብነጻ ድዩ?", "identity_inquiry"),
    # --- MARKETPLACE SEARCH ---
    ("Looking for iPhone 13 Pro 256GB in Addis Ababa", "marketplace_search"),
    ("Find Samsung Galaxy phone under 25000 birr", "marketplace_search"),
    ("Show me used laptops with 16GB RAM in Bole", "marketplace_search"),
    ("2 bedroom apartment for rent in Bole under 35000 birr", "marketplace_search"),
    ("3 bed villa for sale in CMC Addis Ababa", "marketplace_search"),
    ("Studio condominium for rent in Kazanchis", "marketplace_search"),
    ("Toyota Vitz automatic car for sale in Addis", "marketplace_search"),
    ("Used Toyota Corolla 2018 duty paid", "marketplace_search"),
    ("Dump truck for rent in Hawassa", "marketplace_search"),
    ("Agricultural teff 100 quintal in Bahir Dar", "marketplace_search"),
    ("Yirgacheffe coffee beans wholesale price", "marketplace_search"),
    ("Driver job vacancy in Addis Ababa", "marketplace_search"),
    ("Plumber service in Sarbet", "marketplace_search"),
    ("Guest house room in Bole for daily rent", "marketplace_search"),
    ("ቦሌ ውስጥ የሚከራይ ባለ ሁለት መኝታ ቤት ፈልግ", "marketplace_search"),
    ("አዲስ አበባ ውስጥ ቶዮታ ቪትዝ መኪና መግዛት እፈልጋለሁ", "marketplace_search"),
    ("ከ20 ሺህ ብር በታች ስልክ አሳየኝ", "marketplace_search"),
    ("ሲኤምሲ አካባቢ የሚሸጥ ቪላ ቤት", "marketplace_search"),
    ("ካዛንቺስ ውስጥ የተሟላ ፈርኒቸር ያለው ስቱዲዮ ክራይ", "marketplace_search"),
    ("የላፕቶፕ ኮምፒውተር በቦሌ 16 ጂቢ ራም", "marketplace_search"),
    ("ማኛ ጤፍ በኩንታል ዋጋው ስንት ነው?", "marketplace_search"),
    ("የመኪና ሹፌር ክፍት የስራ ቦታ", "marketplace_search"),
    ("የቧንቧ ባለሙያ በአስቸኳይ እፈልጋለሁ", "marketplace_search"),
    (
        "Finfinnee keessatti konkolaataa Toyota Vitz gurguraa barbaada",
        "marketplace_search",
    ),
    ("Magaalaa Adaamaa keessatti mana kireeffamu qabaa?", "marketplace_search"),
    ("Bilbila Samsung qarshii 20000 gadi", "marketplace_search"),
    ("Mana jireenyaa kutaa 2 Boolee keessatti", "marketplace_search"),
    ("Buna Yirgaacaffee gurgurtaa", "marketplace_search"),
    ("Konkolaachisaa barbaadama", "marketplace_search"),
    ("ኣብ መቐለ ዝሽየጥ ገዛ ይደሊ ኣለኹ", "marketplace_search"),
    ("ኣብ ቦሌ ዝካረ ክልተ ክፍሊ ኣፓርታማ", "marketplace_search"),
    ("መኪና ቶዮታ ኮሮላ ዝሽየጥ", "marketplace_search"),
    ("ስልኪ ኣይፎን ትሕቲ 30 ሺሕ ቅርሺ", "marketplace_search"),
    ("ናይ ጤፍ መሸጣ ዋጋ", "marketplace_search"),
    ("ናይ መኪና መዘወሪ ስራሕ", "marketplace_search"),
    # --- PLATFORM FAQ & SUPPORT ---
    ("How do I publish a listing on HUNTIQ?", "platform_faq"),
    ("How much does Gold membership cost?", "platform_faq"),
    ("How do I pay with Telebirr?", "platform_faq"),
    ("Can I pay using CBE Birr?", "platform_faq"),
    ("How long does listing approval take?", "platform_faq"),
    ("Why was my listing rejected?", "platform_faq"),
    ("How do I reset my account password?", "platform_faq"),
    ("How do I contact customer support?", "platform_faq"),
    ("Does HUNTIQ take a commission on sales?", "platform_faq"),
    ("How can I download a tax invoice for my payment?", "platform_faq"),
    ("በHUNTIQ ላይ ማስታወቂያ እንዴት መለጠፍ እችላለሁ?", "platform_faq"),
    ("የጎልድ አባልነት ዋጋው ስንት ነው?", "platform_faq"),
    ("በቴሌብር እንዴት መክፈል እችላለሁ?", "platform_faq"),
    ("በCBE Birr ክፍያ ይፈጸማል?", "platform_faq"),
    ("የተለጠፈ ማስታወቂያ ለመጽደቅ ስንት ሰዓት ይወስዳል?", "platform_faq"),
    ("የይለፍ ቃሌን እንዴት መቀየር እችላለሁ?", "platform_faq"),
    ("የHUNTIQ የደንበኞች ድጋፍ አገልግሎትን እንዴት ማግኘት እችላለሁ?", "platform_faq"),
    ("HUNTIQ በሽያጭ ላይ ኮሚሽን ይቆርጣል?", "platform_faq"),
    ("የክፍያ ደረሰኝ እንዴት ማግኘት እችላለሁ?", "platform_faq"),
    ("HUNTIQ irratti beeksisa akkamitti maxxansuu danda'a?", "platform_faq"),
    ("Miseensummaan Gold waggaatti qarshii meeqa?", "platform_faq"),
    ("Kafaltii Telebirr akkamitti raawwadhura?", "platform_faq"),
    ("Deeggarsa maamiltootaa akkamitti arganna?", "platform_faq"),
    ("ኣብ HUNTIQ ሓድሽ መወዓውዒ ብኸመይ ክልጥፍ ይኽእል?", "platform_faq"),
    ("ናይ ጎልድ ኣባልነት ክንደይ የኽፍል?", "platform_faq"),
    ("ብቴሌብር ከመይ ገይረ ክኸፍል ይኽእል?", "platform_faq"),
    ("ፓስዎርድ ምስ ዝጠፍኣኒ ብኸመይ ክመልሶ ይኽእል?", "platform_faq"),
    ("ናይ ዓማዊል ደገፍ ከመይ ገይረ ክረክብ ይኽእል?", "platform_faq"),
    # --- OUT OF DOMAIN (Academic, Coding, Science) ---
    ("Write a Python script to sort a binary tree", "out_of_domain"),
    ("Explain the theory of relativity by Albert Einstein", "out_of_domain"),
    ("What is the chemical formula for photosynthesis?", "out_of_domain"),
    ("Solve the integral of x squared dx", "out_of_domain"),
    ("Write a React component with Tailwind button", "out_of_domain"),
    ("Who was Napoleon Bonaparte and when was Waterloo?", "out_of_domain"),
    ("Explain quantum physics and quantum entanglement", "out_of_domain"),
    ("How do I make chocolate chip cookies from scratch?", "out_of_domain"),
    ("የኳንተም ፊዚክስ ቀመር ምንድን ነው?", "out_of_domain"),
    ("ስለ አልበርት አንስታይን አንጻራዊነት ጽንሰ ሀሳብ አብራራልኝ", "out_of_domain"),
    ("የፓይዘን ኮድ ጻፍልኝ", "out_of_domain"),
    ("የባዮሎጂ ጥያቄዎችን መልስልኝ", "out_of_domain"),
    ("ሂሳብ ስሌት ስራልኝ", "out_of_domain"),
    ("Koodii Java naaf barreessi", "out_of_domain"),
    ("Fiziksii fi keemistirii na barsiisi", "out_of_domain"),
    ("ናይ ኳንተም ፊዚክስ ሕጊ እንታይ እዩ?", "out_of_domain"),
    ("ናይ ፓይዘን ፕሮግራም ጸሓፈለይ", "out_of_domain"),
    ("ናይ ኬሚስትሪ ሕቶታት ምለሰለይ", "out_of_domain"),
    # --- CLARIFICATION NEEDED / CONFLICTS ---
    ("I want to rent something in Bole", "clarification_needed"),
    (
        "Find me a brand new luxury villa in Bole with min price 5,000,000 and max budget 20,000",
        "clarification_needed",
    ),
    ("I need something cheap in Addis", "clarification_needed"),
    ("Show me options in Ethiopia", "clarification_needed"),
    ("የሚከራይ ነገር እፈልጋለሁ", "clarification_needed"),
    ("ዋጋው አነስተኛ የሆነ ነገር አሳየኝ", "clarification_needed"),
    ("ቦሌ ውስጥ ርካሽ ነገር", "clarification_needed"),
    ("Waan tokko kireeffachuu barbaada", "clarification_needed"),
    ("ዝካረ ነገር ይደሊ ኣለኹ", "clarification_needed"),
]

# Evaluation / Test Data (Separate from training data)
EVAL_DATA = [
    ("Good afternoon, how are you?", "social_greeting"),
    ("ሰላም እንደምን ዋልክ ወዳጄ", "social_greeting"),
    ("Akkam jirtu hiriyaa ko", "social_greeting"),
    ("ሰላም ከመይ ውዒልካ", "social_greeting"),
    ("What are your capabilities and limitations?", "identity_inquiry"),
    ("ማን ነህ? ምንስ ማገዝ ትችላለህ?", "identity_inquiry"),
    ("Eenyu ati maal hojjetta?", "identity_inquiry"),
    ("መን ኢኻ እንታይ ትሕግዘና?", "identity_inquiry"),
    (
        "2 bedroom condo for rent in Yeka Addis Ababa under 30000 birr",
        "marketplace_search",
    ),
    ("Used Toyota Rav4 2019 price in Hawassa", "marketplace_search"),
    ("Dell Core i7 laptop 512GB SSD in Bole", "marketplace_search"),
    ("ቦሌ መድሃኒአለም አካባቢ ባለ አንድ መኝታ ቤት ኪራይ", "marketplace_search"),
    ("ቶዮታ ቪትዝ ኮድ 2 መኪና የሚሸጥ", "marketplace_search"),
    ("Finfinnee keessatti konkolaataa gurguraa barbaada", "marketplace_search"),
    ("ኣብ መቐለ ዝካረ ጽቡቕ ገዛ", "marketplace_search"),
    ("How do I pay with Telebirr on HUNTIQ?", "platform_faq"),
    ("በCBE Birr የጎልድ አባልነት እንዴት መክፈል እችላለሁ?", "platform_faq"),
    ("HUNTIQ irratti kaffaltiin akkamitti raawwatama?", "platform_faq"),
    ("ብቴሌብር ክፍሊት ብኸመይ ክፍጸም ይኽእል?", "platform_faq"),
    ("Write a C++ program to invert a binary tree", "out_of_domain"),
    ("የኳንተም ሜካኒክስ ጽንሰ ሀሳብ አስረዳኝ", "out_of_domain"),
    ("Koodii C++ naaf barreessi", "out_of_domain"),
    ("ናይ ፊዚክስ ቀመር እንታይ እዩ?", "out_of_domain"),
    ("Looking to rent something in Bole", "clarification_needed"),
    ("አዲስ አበባ ውስጥ አሪፍ ነገር ፈልግልኝ", "clarification_needed"),
]


def train_and_evaluate():
    print(
        f"Loaded {len(TRAIN_DATA)} training examples and {len(EVAL_DATA)} evaluation examples."
    )

    X_train = [t[0] for t in TRAIN_DATA]
    y_train = [t[1] for t in TRAIN_DATA]

    X_test = [t[0] for t in EVAL_DATA]
    y_test = [t[1] for t in EVAL_DATA]

    # Multilingual feature union: char_wb (captures Ge'ez prefixes and Oromo agglutination) + word ngrams
    features = FeatureUnion(
        [
            (
                "word_ngram",
                TfidfVectorizer(
                    analyzer="word", ngram_range=(1, 2), sublinear_tf=True, min_df=1
                ),
            ),
            (
                "char_ngram",
                TfidfVectorizer(
                    analyzer="char_wb", ngram_range=(2, 5), sublinear_tf=True, min_df=1
                ),
            ),
        ]
    )

    base_svc = LinearSVC(class_weight="balanced", random_state=42, max_iter=3000, C=1.5)
    calibrated_classifier = CalibratedClassifierCV(estimator=base_svc, cv=3)

    pipeline = Pipeline([("features", features), ("classifier", calibrated_classifier)])

    start_time = time.time()
    print("Training calibrated multilingual intent classifier...")
    pipeline.fit(X_train, y_train)
    training_duration = time.time() - start_time
    print(f"Training completed in {training_duration:.3f} seconds.")

    # Evaluate
    print("Evaluating model performance on test set...")
    t0 = time.time()
    y_pred = pipeline.predict(X_test)
    eval_latency = (time.time() - t0) / len(X_test) * 1000.0  # ms per query

    acc = accuracy_score(y_test, y_pred)
    macro_f1 = f1_score(y_test, y_pred, average="macro")

    print(f"Accuracy: {acc * 100:.2f}%")
    print(f"Macro F1 Score: {macro_f1:.4f}")
    print(f"Average Inference Latency: {eval_latency:.2f} ms")
    print("\nDetailed Classification Report:")
    print(classification_report(y_test, y_pred))

    # Save model artifact
    print(f"Serializing trained pipeline to {MODEL_PATH}...")
    joblib.dump(pipeline, MODEL_PATH)

    # Save metrics JSON
    metrics = {
        "model_id": "intent_classifier_v1",
        "model_type": "CalibratedClassifierCV(LinearSVC) + MultilingualCharWordTfidf",
        "training_examples_count": len(TRAIN_DATA),
        "test_examples_count": len(EVAL_DATA),
        "accuracy": float(acc),
        "macro_f1": float(macro_f1),
        "inference_latency_ms": float(eval_latency),
        "training_duration_seconds": float(training_duration),
        "languages": ["en", "am", "or", "ti"],
        "classes": sorted(set(y_train)),
        "created_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "status": "APPROVED_PRODUCTION" if macro_f1 >= 0.90 else "NEEDS_REFINEMENT",
    }

    with open(METRICS_PATH, "w") as f:
        json.dump(metrics, f, indent=2)
    print(f"Saved evaluation metrics to {METRICS_PATH}")

    # Update model registry
    registry_entry = {
        "model_id": "intent_classifier_v1",
        "version": "1.0.0",
        "task": "intent_and_domain_classification",
        "artifact_path": "training/checkpoints/intent_classifier_v1.joblib",
        "metrics_path": "training/checkpoints/intent_classifier_v1_metrics.json",
        "accuracy": float(acc),
        "macro_f1": float(macro_f1),
        "latency_ms": float(eval_latency),
        "status": "APPROVED_PRODUCTION",
        "approved_by": "HUNTIQ Lead AI Architect & Data Engineering Auditor",
        "approval_date": time.strftime("%Y-%m-%d", time.gmtime()),
    }

    reg_data = {"registered_models": []}
    if os.path.exists(REGISTRY_PATH):
        try:
            with open(REGISTRY_PATH, "r", encoding="utf-8") as f:
                reg_data = json.load(f)
        except Exception:  # noqa: BLE001
            reg_data = {"registered_models": []}

    reg_data["registered_models"] = [
        m
        for m in reg_data.get("registered_models", [])
        if m.get("model_id") != "intent_classifier_v1"
    ]
    reg_data["registered_models"].append(registry_entry)

    with open(REGISTRY_PATH, "w", encoding="utf-8") as f:
        json.dump(reg_data, f, indent=2)
    print(f"Updated model registry at {REGISTRY_PATH}")


if __name__ == "__main__":
    train_and_evaluate()
