"""Tests verifying that HUNTIQ AI operates entirely offline without external cloud dependencies."""

import asyncio
from unittest.mock import MagicMock

from ai.app.config import config
from ai.app.model_client import OllamaClient


def test_ollama_config_uses_local_address():
    """Verify that configuration points to local loopback interface."""
    assert (
        "127.0.0.1" in config.ollama_base_url or "localhost" in config.ollama_base_url
    )
    assert not any(
        cloud in config.ollama_base_url
        for cloud in ["api.openai.com", "anthropic.com", "googleapis.com"]
    )


def test_no_external_api_calls_made(monkeypatch):
    """Ensure that client only attempts to communicate with the local Ollama URL."""

    async def run():
        client = OllamaClient()

        called_urls = []

        async def fake_post(self, url, **kwargs):
            called_urls.append(str(url))
            mock_res = MagicMock(status_code=200)
            mock_res.json.return_value = {
                "message": {"role": "assistant", "content": "offline test ok"}
            }
            return mock_res

        import httpx

        monkeypatch.setattr(httpx.AsyncClient, "post", fake_post)

        await client.chat(
            messages=[{"role": "user", "content": "Testing offline capability"}]
        )

        assert len(called_urls) == 1
        assert called_urls[0].startswith("http://127.0.0.1:11434")

    asyncio.run(run())
