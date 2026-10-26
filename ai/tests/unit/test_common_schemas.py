from datetime import UTC, datetime

import pytest
from pydantic import ValidationError

from huntiq_ai.schemas.common import (
    ErrorDetail,
    ErrorResponse,
    Locale,
    QueryIntent,
    ResponseMetadata,
)


def test_locale_enum_members():
    """Verify all 4 Ethiopian marketplace language locales are supported."""
    assert [locale.value for locale in Locale] == ["en", "am", "om", "ti"]
    assert Locale.EN == "en"
    assert Locale.AM == "am"
    assert Locale.OM == "om"
    assert Locale.TI == "ti"


def test_query_intent_enum_members():
    """Verify all 7 marketplace query intents are represented."""
    expected_intents = {
        "buy_product",
        "rent_property",
        "find_service",
        "job_search",
        "vehicle_search",
        "price_check",
        "ambiguous",
    }
    assert {intent.value for intent in QueryIntent} == expected_intents


def test_error_detail_model_validation():
    """Validate RFC 7807 error detail parsing."""
    now = datetime.now(UTC)
    detail = ErrorDetail(
        code="UNAUTHORIZED",
        message="Invalid secret",
        correlation_id="cid-123",
        timestamp=now,
        details={"hint": "Check X-HUNTIQ-AI-SECRET"},
    )
    assert detail.code == "UNAUTHORIZED"
    assert detail.correlation_id == "cid-123"
    assert detail.details == {"hint": "Check X-HUNTIQ-AI-SECRET"}


def test_error_response_envelope():
    """Validate full error response envelope."""
    payload = {
        "error": {
            "code": "PAYLOAD_TOO_LARGE",
            "message": "Exceeded 64KB",
            "correlation_id": "cid-456",
            "timestamp": "2026-10-26T12:00:00Z",
        }
    }
    response = ErrorResponse.model_validate(payload)
    assert response.error.code == "PAYLOAD_TOO_LARGE"
    assert response.error.correlation_id == "cid-456"
    assert response.error.details is None


def test_error_response_missing_fields_fails():
    """Missing required fields in error response must raise ValidationError."""
    with pytest.raises(ValidationError):
        ErrorResponse.model_validate({"error": {"code": "BAD"}})


def test_response_metadata():
    """Validate response telemetry metadata."""
    meta = ResponseMetadata(correlation_id="cid-789", latency_ms=12.5)
    assert meta.correlation_id == "cid-789"
    assert meta.latency_ms == 12.5
