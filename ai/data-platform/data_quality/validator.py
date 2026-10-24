"""HUNTIQ AI Data Platform Quality and Governance Validator.

Validates dataset manifests, scans for PII (including Ethiopian phone formats),
detects duplicate records, and prevents train/test leakage.
"""

from __future__ import annotations

import re
from dataclasses import asdict, dataclass, field
from datetime import UTC, datetime
from typing import Any


class ManifestValidationError(Exception):
    """Raised when a dataset manifest fails schema or checksum validation."""

    pass


@dataclass
class DataQualityReport:
    """Structured report returned from dataset quality validation."""

    valid: bool
    total_samples: int
    pii_violations: list[dict[str, str]] = field(default_factory=list)
    duplicate_violations: list[str] = field(default_factory=list)
    leakage_violations: list[str] = field(default_factory=list)
    manifest_errors: list[str] = field(default_factory=list)
    timestamp: str = field(default_factory=lambda: datetime.now(UTC).isoformat())

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


class DataQualityValidator:
    """Enforces quality gates, privacy rules, and split hygiene on datasets."""

    # Ethiopian Mobile Phone Regex: Supports +251 9..., 09..., +251 7..., 07...
    ETHIOPIAN_PHONE_PATTERN = re.compile(
        r"(?:(?:\+251\s?|0)[79]\d{2}[\s.-]?\d{3}[\s.-]?\d{3}|(?:\+251\s?|0)[79]\d{8})"
    )

    # Standard Email Regex
    EMAIL_PATTERN = re.compile(r"\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,7}\b")

    # Common credential or secret patterns in raw text
    SECRET_PATTERN = re.compile(
        r"(?:password|passwd|secret|api_key|bearer\s+[a-zA-Z0-9_\-\.]+)",
        re.IGNORECASE,
    )

    @classmethod
    def scan_for_pii(cls, text: str) -> list[str]:
        """Scans a text field for unscrubbed PII or secrets.

        Returns a list of violation descriptions (redacted for safe reporting).
        """
        violations: list[str] = []

        phone_matches = cls.ETHIOPIAN_PHONE_PATTERN.findall(text)
        if phone_matches:
            violations.append(f"Contains {len(phone_matches)} phone number pattern(s)")

        email_matches = cls.EMAIL_PATTERN.findall(text)
        if email_matches:
            violations.append(f"Contains {len(email_matches)} email address pattern(s)")

        secret_matches = cls.SECRET_PATTERN.findall(text)
        if secret_matches:
            violations.append(
                f"Contains {len(secret_matches)} potential secret/credential keyword(s)"
            )

        return violations

    @classmethod
    def validate_manifest(cls, manifest: dict[str, Any]) -> list[str]:
        """Validates manifest schema and split consistency.

        Returns a list of error strings (empty if valid).
        """
        errors: list[str] = []
        required_fields = [
            "manifest_id",
            "dataset_name",
            "version",
            "created_at",
            "source_export_id",
            "record_count",
            "sha256_checksum",
            "pii_scrubbed",
            "splits",
        ]

        for req in required_fields:
            if req not in manifest:
                errors.append(f"Missing required manifest field: '{req}'")

        if not manifest.get("pii_scrubbed", False):
            errors.append("Manifest pii_scrubbed must be true for approved datasets")

        splits = manifest.get("splits")
        if isinstance(splits, dict):
            train_cnt = splits.get("train", 0)
            val_cnt = splits.get("val", 0)
            test_cnt = splits.get("test", 0)
            total_splits = train_cnt + val_cnt + test_cnt
            record_count = manifest.get("record_count", 0)

            if total_splits != record_count:
                errors.append(
                    f"Splits sum ({total_splits}) does not match record_count ({record_count})"
                )
        else:
            errors.append(
                "Manifest 'splits' field must be an object with train, val, and test counts"
            )

        return errors

    @classmethod
    def detect_duplicates(
        cls, records: list[dict[str, Any]], field_name: str = "query"
    ) -> list[str]:
        """Identifies duplicate sample entries based on the target text field."""
        seen: set[str] = set()
        duplicates: list[str] = []

        for record in records:
            value = str(record.get(field_name, "")).strip().lower()
            if not value:
                continue
            if value in seen:
                duplicates.append(value)
            else:
                seen.add(value)

        return list(set(duplicates))

    @classmethod
    def detect_leakage(
        cls,
        records: list[dict[str, Any]],
        text_field: str = "query",
        split_field: str = "split",
    ) -> list[str]:
        """Detects samples appearing in both train and evaluation (val/test) splits."""
        train_texts: set[str] = set()
        eval_texts: set[str] = set()

        for rec in records:
            split_type = rec.get(split_field, "train")
            text = str(rec.get(text_field, "")).strip().lower()
            if not text:
                continue

            if split_type == "train":
                train_texts.add(text)
            elif split_type in ("val", "test", "evaluation"):
                eval_texts.add(text)

        overlap = train_texts.intersection(eval_texts)
        return list(overlap)

    @classmethod
    def validate_dataset(
        cls,
        records: list[dict[str, Any]],
        manifest: dict[str, Any] | None = None,
        text_field: str = "query",
        split_field: str = "split",
    ) -> DataQualityReport:
        """Runs full quality suite: manifest validation, PII, duplicates, and leakage."""
        manifest_errors = cls.validate_manifest(manifest) if manifest else []
        pii_violations: list[dict[str, str]] = []

        for idx, rec in enumerate(records):
            sample_id = rec.get("sample_id", f"index-{idx}")
            text_to_scan = " ".join(
                str(v) for k, v in rec.items() if isinstance(v, (str, int, float))
            )
            violations = cls.scan_for_pii(text_to_scan)
            for v in violations:
                pii_violations.append({"sample_id": sample_id, "violation": v})

        duplicates = cls.detect_duplicates(records, field_name=text_field)
        leakage = cls.detect_leakage(records, text_field=text_field, split_field=split_field)

        is_valid = (
            len(manifest_errors) == 0
            and len(pii_violations) == 0
            and len(duplicates) == 0
            and len(leakage) == 0
        )

        return DataQualityReport(
            valid=is_valid,
            total_samples=len(records),
            pii_violations=pii_violations,
            duplicate_violations=duplicates,
            leakage_violations=leakage,
            manifest_errors=manifest_errors,
        )


def validate_dataset_manifest(manifest_dict: dict[str, Any]) -> None:
    """Convenience helper that raises ManifestValidationError if invalid."""
    errors = DataQualityValidator.validate_manifest(manifest_dict)
    if errors:
        raise ManifestValidationError(f"Invalid manifest: {'; '.join(errors)}")
