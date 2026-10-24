import pytest

from huntiq_ai.core.config import Settings


def test_default_settings_development():
    """Default settings in development should have safe default values."""
    settings = Settings(environment="development")
    assert settings.environment == "development"
    assert settings.host == "127.0.0.1"
    assert settings.port == 8000
    assert settings.max_payload_size_bytes == 65536


def test_production_settings_rejects_default_secret():
    """Production mode must reject default or short placeholder secret."""
    with pytest.raises(
        ValueError,
        match="Production environment requires a strong HUNTIQ_AI_SECRET_KEY",
    ):
        Settings(
            environment="production",
            huntiq_ai_secret_key="dev_secret_change_in_production",
        )


def test_production_settings_rejects_short_secret():
    """Production mode must reject secrets with fewer than 32 characters."""
    with pytest.raises(ValueError, match="at least 32 characters"):
        Settings(
            environment="production",
            huntiq_ai_secret_key="short_secret",
        )


def test_production_settings_accepts_strong_secret():
    """Production mode accepts strong 32+ character secret."""
    strong_secret = "a" * 48
    settings = Settings(
        environment="production",
        huntiq_ai_secret_key=strong_secret,
    )
    assert settings.huntiq_ai_secret_key == strong_secret
