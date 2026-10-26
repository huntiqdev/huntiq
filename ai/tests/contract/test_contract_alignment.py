import json
from pathlib import Path

from fastapi.testclient import TestClient

from huntiq_ai.schemas.common import ErrorResponse, Locale, QueryIntent

DATA_CONTRACTS_DIR = Path(__file__).parents[2] / "data-platform" / "data-contracts" / "schemas"
TAXONOMY_SCHEMAS_DIR = Path(__file__).parents[2] / "data-platform" / "taxonomy" / "schemas"


def test_query_intent_contract_locale_alignment():
    """Ensure Python Locale enum matches the JSON Schema query_intent_contract."""
    contract_path = DATA_CONTRACTS_DIR / "query_intent_contract.json"
    assert contract_path.exists(), f"Contract file not found: {contract_path}"

    with open(contract_path, encoding="utf-8") as f:
        contract = json.load(f)

    allowed_locales = set(contract["properties"]["locale"]["enum"])
    defined_locales = {loc.value for loc in Locale}

    assert defined_locales == allowed_locales, (
        f"Mismatch between Python Locale ({defined_locales}) "
        f"and JSON Schema contract ({allowed_locales})"
    )


def test_query_intent_contract_intent_alignment():
    """Ensure Python QueryIntent enum matches the JSON Schema query_intent_contract."""
    contract_path = DATA_CONTRACTS_DIR / "query_intent_contract.json"
    with open(contract_path, encoding="utf-8") as f:
        contract = json.load(f)

    allowed_intents = set(contract["properties"]["intent"]["enum"])
    defined_intents = {intent.value for intent in QueryIntent}

    assert defined_intents == allowed_intents, (
        f"Mismatch between Python QueryIntent ({defined_intents}) "
        f"and JSON Schema contract ({allowed_intents})"
    )


def test_taxonomy_export_schema_translation_locales():
    """Ensure taxonomy export translation keys correspond to valid non-English locales."""
    schema_path = TAXONOMY_SCHEMAS_DIR / "taxonomy_export.schema.json"
    assert schema_path.exists(), f"Taxonomy schema not found: {schema_path}"

    with open(schema_path, encoding="utf-8") as f:
        schema = json.load(f)

    # In taxonomy_export.schema.json, translations maps locale codes (am, om, ti)
    # Verify description or schema documents these exact Ethiopian languages
    translations_desc = schema["properties"]["categories"]["items"]["properties"]["translations"][
        "description"
    ]
    for loc in [Locale.AM, Locale.OM, Locale.TI]:
        assert loc.value in translations_desc


def test_error_response_validates_live_error_payload(client: TestClient):
    """Verify that runtime RFC 7807 error responses validate against ErrorResponse."""
    oversized = "a" * 1500
    response = client.post(
        "/health/live",
        content=oversized,
        headers={"Content-Length": "1500"},
    )
    assert response.status_code == 413
    error_model = ErrorResponse.model_validate(response.json())
    assert error_model.error.code == "PAYLOAD_TOO_LARGE"
    assert error_model.error.correlation_id is not None
    assert error_model.error.timestamp is not None
