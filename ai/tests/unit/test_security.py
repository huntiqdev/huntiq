import pytest
from fastapi.testclient import TestClient

from huntiq_ai.core.errors import AuthenticationError
from huntiq_ai.core.security import verify_internal_secret


@pytest.mark.asyncio
async def test_verify_internal_secret_success():
    """Valid secret header passes authentication."""
    valid_secret = "dev_secret_change_in_production"
    result = await verify_internal_secret(
        request=None,  # type: ignore
        x_huntiq_ai_secret=valid_secret,
    )
    assert result == valid_secret


@pytest.mark.asyncio
async def test_verify_internal_secret_missing_header():
    """Missing secret header raises AuthenticationError."""
    with pytest.raises(AuthenticationError, match="Missing required"):
        await verify_internal_secret(
            request=None,  # type: ignore
            x_huntiq_ai_secret=None,
        )


@pytest.mark.asyncio
async def test_verify_internal_secret_invalid_secret():
    """Invalid secret header raises AuthenticationError."""
    with pytest.raises(AuthenticationError, match="Invalid internal service authentication secret"):
        await verify_internal_secret(
            request=None,  # type: ignore
            x_huntiq_ai_secret="invalid_token_999",
        )


def test_payload_too_large_rejected(client: TestClient):
    """Requests exceeding max_payload_size_bytes must be rejected with HTTP 413."""
    # Test settings configure 1024 bytes limit
    oversized_data = "x" * 2048
    response = client.post(
        "/health/live",
        content=oversized_data,
        headers={"Content-Length": "2048"},
    )
    assert response.status_code == 413
    body = response.json()
    assert body["error"]["code"] == "PAYLOAD_TOO_LARGE"
