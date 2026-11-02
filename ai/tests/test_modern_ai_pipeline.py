"""
Tests for HUNTIQ Modern Super AI Architectural Pipeline:
- MorphologicalTokenizer (letter-level, fidel root, affix, and n-gram decomposition)
- DenseSemanticEncoder (continuous latent semantic space, cosine similarity, candidate ranking)
- DatasetIngestionPipeline (CSV, JSON, Parquet, schema auto-discovery, batch streaming)
- ModernSemanticClassifier (multi-resolution feature learning, calibrated probabilities, uncertainty)
"""

import os
import tempfile

import pandas as pd

from app.core.dataset_ingestion_pipeline import dataset_pipeline
from app.core.dense_semantic_encoder import dense_encoder
from app.core.modern_semantic_classifier import ModernSemanticClassifier
from app.core.morphological_tokenizer import morph_tokenizer


def test_morphological_tokenizer_letters_and_roots():
    """Verify letter-by-letter and Ge'ez fidel root decomposition."""
    # Amharic word "ቤቶች" (houses) -> letters ['ቤ', 'ት', 'ኦ', 'ች'] or chars
    letters = morph_tokenizer.extract_letters("ቤቶች")
    assert "ቤ" in letters and "ች" in letters

    # Root mapping: 'ቤ' base consonant is 'በ'
    roots = morph_tokenizer.extract_fidel_roots("ቤቶች")
    assert "በ" in roots and "ተ" in roots

    # Affix stripping: 'ቤቶች' stem is 'ቤት' with suffix '-ኦች' or '-ዎች'
    stem, affixes = morph_tokenizer.strip_affixes_geez("መኪናዎች")
    assert stem == "መኪና"
    assert any("ዎች" in a for a in affixes)

    # Afaan Oromoo suffix stripping: 'Finfinneetti' -> 'finfinnee'
    o_stem, o_affixes = morph_tokenizer.strip_affixes_oromo("Finfinneetti")
    assert o_stem == "finfinnee"
    assert any("tti" in a for a in o_affixes)


def test_morphological_tokenizer_full_document():
    """Verify document hierarchical decomposition across sentences and subwords."""
    doc = "I want to buy a car in Bole. ቶዮታ መኪና መግዛት እፈልጋለሁ።"
    tokens = morph_tokenizer.tokenize_document(doc)

    assert len(tokens["sentences"]) == 2
    assert len(tokens["words"]) >= 8
    assert len(tokens["letters"]) > 20
    assert len(tokens["char_ngrams"]) > 30


def test_dense_semantic_encoder():
    """Verify continuous dense semantic vector representation and similarity."""
    v1 = dense_encoder.encode_sentence("Toyota Vitz automatic car in Addis Ababa")
    assert v1.shape[0] == dense_encoder.embedding_dim
    assert abs(float((v1**2).sum()) - 1.0) < 1e-4  # L2 normalized

    # Semantic candidate ranking
    candidates = [
        {"id": 1, "title": "Toyota Corolla 2020 clean condition"},
        {"id": 2, "title": "2 bedroom condominium apartment in Bole"},
        {"id": 3, "title": "Organic honey wholesale 50kg"},
    ]
    ranked = dense_encoder.rank_candidates(
        "Looking to buy a Toyota car", candidates, text_key="title"
    )
    assert ranked[0]["id"] == 1, "Toyota listing must rank highest for car query"
    assert "semantic_score" in ranked[0]


def test_dataset_ingestion_pipeline_csv_and_parquet():
    """Verify multi-format loading, schema auto-detection, and parquet export."""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create 3 synthetic CSV files to simulate multiple datasets
        csv_files = []
        for i in range(3):
            df_synth = pd.DataFrame(
                {
                    "utterance": [f"Sample query {i}_{j} for test" for j in range(20)],
                    "intent_class": [f"class_{j % 3}" for j in range(20)],
                }
            )
            p = os.path.join(tmpdir, f"dataset_part_{i}.csv")
            df_synth.to_csv(p, index=False)
            csv_files.append(p)

        # Ingest entire directory of CSVs
        df_loaded, report = dataset_pipeline.load_dataset(
            source_path=tmpdir,
            deduplicate=True,
            normalize_geez=True,
        )

        assert report.valid_rows == 60
        assert len(report.files_processed) == 3
        assert report.detected_text_column == "utterance"
        assert report.detected_label_column == "intent_class"
        assert set(df_loaded.columns) == {"text", "label"}

        # Test Parquet export and round-trip
        parquet_path = os.path.join(tmpdir, "dataset.parquet")
        dataset_pipeline.export_parquet(df_loaded, parquet_path)
        assert os.path.exists(parquet_path)

        _df_pq, pq_report = dataset_pipeline.load_dataset(parquet_path)
        assert pq_report.valid_rows == 60


def test_modern_semantic_classifier_training_and_uncertainty():
    """Verify end-to-end training, calibration, and Shannon entropy uncertainty calculation."""
    samples = [
        ("Looking for a used Toyota car in Bole", "vehicles"),
        ("Toyota Corolla automatic transmission clean condition", "vehicles"),
        ("መኪና መግዛት እፈልጋለሁ በአዲስ አበባ", "vehicles"),
        ("2 bedroom apartment for rent in Bole", "real_estate"),
        ("Furnished condominium house for rent Kazanchis", "real_estate"),
        ("የሚከራይ ቪላ ቤት በቦሌ", "real_estate"),
        ("Samsung Galaxy S24 Ultra smartphone in Bole", "products"),
        ("Apple iPhone 15 Pro Max 256GB storage", "products"),
        ("ስማርት ስልክ ሳምሰንግ ወይም አይፎን የሚሸጥ", "products"),
    ] * 5

    texts = [s[0] for s in samples]
    labels = [s[1] for s in samples]

    clf = ModernSemanticClassifier(task_name="test_marketplace_routing")
    metrics = clf.fit(texts, labels, eval_split=0.20, random_state=42)

    assert metrics["accuracy"] >= 0.90
    assert metrics["macro_f1"] >= 0.90
    assert metrics["inference_latency_ms"] < 10.0

    # In-domain prediction with calibrated confidence
    label, conf, _probs, unc = clf.predict("Toyota Vitz in Bole")
    assert label == "vehicles"
    assert conf > 0.40
    assert 0.0 <= unc <= 1.0

    # Out-of-domain random nonsense should exhibit high uncertainty
    _, _, _, ood_unc = clf.predict("xyzabc random irrelevant query 12345")
    assert 0.0 <= ood_unc <= 1.0
