"""Benchmark dataset validation tests for HUNTIQ AI evaluation."""

import json
from pathlib import Path

DATASETS_DIR = Path(__file__).resolve().parent.parent / "datasets" / "evaluation"


def test_domain_benchmark_dataset():
    path = DATASETS_DIR / "domain_benchmark.jsonl"
    assert path.exists(), f"Missing dataset: {path}"

    lines = path.read_text(encoding="utf-8").strip().splitlines()
    assert len(lines) >= 10, f"Expected at least 10 domain samples, got {len(lines)}"

    valid_domains = {"out_of_domain", "platform_faq", "marketplace", "social_greeting", "identity"}
    for idx, line in enumerate(lines):
        data = json.loads(line)
        assert "id" in data, f"Line {idx} missing 'id'"
        assert "query" in data, f"Line {idx} missing 'query'"
        assert data["expected_domain"] in valid_domains, f"Invalid domain {data['expected_domain']}"


def test_entity_extraction_dataset():
    path = DATASETS_DIR / "entity_extraction.jsonl"
    assert path.exists(), f"Missing dataset: {path}"

    lines = path.read_text(encoding="utf-8").strip().splitlines()
    assert len(lines) >= 5, f"Expected at least 5 entity extraction samples, got {len(lines)}"

    for idx, line in enumerate(lines):
        data = json.loads(line)
        assert "id" in data
        assert "query" in data
        assert "expected_searches" in data
        assert isinstance(data["expected_searches"], list)


def test_conversational_benchmark_dataset():
    path = DATASETS_DIR / "conversational_benchmark.jsonl"
    assert path.exists(), f"Missing dataset: {path}"

    lines = path.read_text(encoding="utf-8").strip().splitlines()
    assert len(lines) >= 3, f"Expected at least 3 conversational samples, got {len(lines)}"

    for idx, line in enumerate(lines):
        data = json.loads(line)
        assert "id" in data
        assert "query" in data
        assert "expected_action" in data
