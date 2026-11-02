"""Configuration settings for HUNTIQ Ollama integration and evaluation."""

import os

from pydantic import BaseModel, Field


class AIConfig(BaseModel):
    """Configuration schema for HUNTIQ AI local inference."""

    ollama_base_url: str = Field(
        default_factory=lambda: os.getenv("OLLAMA_BASE_URL", "http://127.0.0.1:11434")
    )
    base_model: str = Field(
        default_factory=lambda: os.getenv("OLLAMA_BASE_MODEL", "qwen3:4b")
    )
    huntiq_model: str = Field(
        default_factory=lambda: os.getenv(
            "OLLAMA_HUNTIQ_MODEL", "huntiq-assistant:latest"
        )
    )
    request_timeout: float = Field(
        default_factory=lambda: float(os.getenv("AI_REQUEST_TIMEOUT", "60.0"))
    )
    num_ctx: int = Field(
        default_factory=lambda: int(os.getenv("AI_CONTEXT_LENGTH", "4096"))
    )
    allowed_tools: list[str] = Field(
        default_factory=lambda: [
            "search_marketplace_listings",
            "inspect_seller_reputation",
            "get_listing_details",
            "compare_listings",
            "get_marketplace_stats",
        ]
    )
    supported_languages: list[str] = Field(
        default_factory=lambda: ["en", "am", "om", "ti"]
    )
    api_secret_key: str = Field(
        default_factory=lambda: os.getenv(
            "HUNTIQ_AI_SECRET",
            os.getenv("AI_SECRET_KEY", "huntiq-local-ai-secret-key-2026"),
        )
    )
    go_backend_url: str = Field(
        default_factory=lambda: os.getenv("GO_BACKEND_URL", "http://127.0.0.1:8080")
    )
    host: str = Field(default_factory=lambda: os.getenv("HOST", "127.0.0.1"))
    port: int = Field(default_factory=lambda: int(os.getenv("PORT", "8000")))


# Singleton instance for import
config = AIConfig()
