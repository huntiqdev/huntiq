from functools import lru_cache
from typing import Literal

from pydantic import Field, field_validator
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Central typed settings for the HUNTIQ AI microservice."""

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
        extra="ignore",
    )

    # Environment
    environment: Literal["development", "staging", "production", "test"] = Field(
        default="development",
        description="Current runtime environment",
    )

    # Server Bind
    host: str = Field(default="127.0.0.1", description="Bind IP address")
    port: int = Field(default=8000, description="Listening TCP port")

    # Logging
    log_level: str = Field(default="INFO", description="Log level: DEBUG, INFO, WARNING, ERROR")
    log_format: Literal["json", "text"] = Field(
        default="json", description="Log format: json or text"
    )

    # Security
    huntiq_ai_secret_key: str = Field(
        default="dev_secret_change_in_production",
        description="Internal service authentication pre-shared secret key",
    )
    max_payload_size_bytes: int = Field(
        default=65536,
        description="Maximum accepted request body size in bytes (default 64KB)",
    )
    request_timeout_seconds: float = Field(
        default=5.0,
        description="Maximum request processing timeout in seconds",
    )

    # Model Storage
    model_dir: str = Field(
        default="./models",
        description="Directory containing approved model weights",
    )

    @field_validator("huntiq_ai_secret_key")
    @classmethod
    def validate_secret_key(cls, v: str, info) -> str:
        """Fail fast if running in production with weak or default secret."""
        # Get environment value if present in values dict
        env = info.data.get("environment")
        if env == "production":
            if v == "dev_secret_change_in_production" or len(v.strip()) < 32:
                raise ValueError(
                    "Production environment requires a strong HUNTIQ_AI_SECRET_KEY "
                    "with at least 32 characters."
                )
        return v


@lru_cache(maxsize=1)
def get_settings() -> Settings:
    """Returns the cached singleton Settings instance."""
    return Settings()
