from data_quality.validator import DataQualityValidator


def test_synthetic_dataset_full_validation_passes(
    sample_dataset: list[dict], sample_manifest: dict
):
    """Clean synthetic dataset passes full quality suite with 0 violations."""
    report = DataQualityValidator.validate_dataset(
        records=sample_dataset,
        manifest=sample_manifest,
        text_field="query",
        split_field="split",
    )
    assert report.valid is True
    assert report.total_samples == len(sample_dataset)
    assert len(report.pii_violations) == 0
    assert len(report.duplicate_violations) == 0
    assert len(report.leakage_violations) == 0
    assert len(report.manifest_errors) == 0


def test_scan_for_pii_ethiopian_phone_numbers():
    """Detects various Ethiopian phone number formats."""
    # Test formats: 0911234567, +251 911 234 567, 0712345678
    test_cases = [
        "call me at 0911234567 for details",
        "contact +251 911 234 567 via telegram",
        "call seller on 0712345678",
    ]
    for text in test_cases:
        violations = DataQualityValidator.scan_for_pii(text)
        assert len(violations) > 0
        assert any("phone number" in v for v in violations)


def test_scan_for_pii_email_address():
    """Detects email addresses in text."""
    text = "reach out to seller at seller@gmail.com for bulk price"
    violations = DataQualityValidator.scan_for_pii(text)
    assert len(violations) > 0
    assert any("email address" in v for v in violations)


def test_detect_duplicates_finds_exact_matches():
    """Detects duplicate queries in the records list."""
    records = [
        {"sample_id": "1", "query": "toyota vitz 2015"},
        {"sample_id": "2", "query": "samsung galaxy s22"},
        {"sample_id": "3", "query": "toyota vitz 2015"},  # Duplicate
    ]
    duplicates = DataQualityValidator.detect_duplicates(records, field_name="query")
    assert len(duplicates) == 1
    assert "toyota vitz 2015" in duplicates


def test_detect_leakage_finds_cross_split_contamination():
    """Detects when an identical query exists in both train and evaluation splits."""
    records = [
        {"sample_id": "1", "query": "iphone 13 pro max", "split": "train"},
        {"sample_id": "2", "query": "leather sofa 3 seater", "split": "train"},
        {"sample_id": "3", "query": "iphone 13 pro max", "split": "test"},  # Leakage!
    ]
    leakages = DataQualityValidator.detect_leakage(records, text_field="query", split_field="split")
    assert len(leakages) == 1
    assert "iphone 13 pro max" in leakages
