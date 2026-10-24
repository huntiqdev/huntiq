import json
from pathlib import Path

import pytest

DATA_PLATFORM_ROOT = Path(__file__).parent.parent


@pytest.fixture
def sample_manifest_path() -> Path:
    fixtures = DATA_PLATFORM_ROOT / "ml" / "datasets" / "synthetic_fixtures"
    return fixtures / "sample_manifest.json"


@pytest.fixture
def sample_dataset_path() -> Path:
    fixtures = DATA_PLATFORM_ROOT / "ml" / "datasets" / "synthetic_fixtures"
    return fixtures / "sample_query_intent_dataset.json"


@pytest.fixture
def sample_manifest(sample_manifest_path: Path) -> dict:
    with open(sample_manifest_path, encoding="utf-8") as f:
        return json.load(f)


@pytest.fixture
def sample_dataset(sample_dataset_path: Path) -> list[dict]:
    with open(sample_dataset_path, encoding="utf-8") as f:
        return json.load(f)
