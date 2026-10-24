import pytest
from fastapi.testclient import TestClient

from huntiq_ai.core.config import Settings, get_settings
from huntiq_ai.main import create_app


@pytest.fixture
def test_settings() -> Settings:
    """Returns safe settings for test executions."""
    return Settings(
        environment="test",
        host="127.0.0.1",
        port=8000,
        log_level="DEBUG",
        log_format="text",
        huntiq_ai_secret_key="test_secret_32_characters_long_for_safe_testing",
        max_payload_size_bytes=1024,  # Small 1KB limit for fast testing
        request_timeout_seconds=2.0,
        model_dir="./test_models",
    )


@pytest.fixture
def client(test_settings: Settings) -> TestClient:
    """FastAPI TestClient with overridden test settings."""
    app = create_app(test_settings)
    app.dependency_overrides[get_settings] = lambda: test_settings
    with TestClient(app) as test_client:
        yield test_client
