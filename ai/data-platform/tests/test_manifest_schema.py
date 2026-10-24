import pytest

from data_quality.validator import (
    DataQualityValidator,
    ManifestValidationError,
    validate_dataset_manifest,
)


def test_sample_manifest_is_valid(sample_manifest: dict):
    """Synthetic sample manifest must pass all validation checks."""
    errors = DataQualityValidator.validate_manifest(sample_manifest)
    assert errors == []


def test_manifest_validation_helper_succeeds(sample_manifest: dict):
    """Helper function succeeds without raising an exception for valid manifest."""
    validate_dataset_manifest(sample_manifest)


def test_manifest_missing_required_field_fails():
    """Manifest missing critical fields must return errors."""
    incomplete_manifest = {
        "manifest_id": "mf-incomplete",
        "dataset_name": "test_dataset",
    }
    errors = DataQualityValidator.validate_manifest(incomplete_manifest)
    assert len(errors) > 0
    assert any("version" in err for err in errors)
    assert any("sha256_checksum" in err for err in errors)

    with pytest.raises(ManifestValidationError):
        validate_dataset_manifest(incomplete_manifest)


def test_manifest_split_mismatch_fails(sample_manifest: dict):
    """Manifest whose split sums do not equal record_count must fail."""
    corrupted_manifest = dict(sample_manifest)
    corrupted_manifest["record_count"] = 100  # Mismatch: splits total is 6
    errors = DataQualityValidator.validate_manifest(corrupted_manifest)
    assert any("does not match record_count" in err for err in errors)


def test_manifest_unscrubbed_pii_fails(sample_manifest: dict):
    """Manifest with pii_scrubbed=False must be rejected."""
    unscrubbed_manifest = dict(sample_manifest)
    unscrubbed_manifest["pii_scrubbed"] = False
    errors = DataQualityValidator.validate_manifest(unscrubbed_manifest)
    assert any("pii_scrubbed must be true" in err for err in errors)
