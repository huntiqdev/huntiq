"""Unit tests for Ollama model client."""

import asyncio
from unittest.mock import AsyncMock, MagicMock, patch

from ai.app.model_client import OllamaClient


def test_ollama_client_health_check_success():
    async def run():
        client = OllamaClient(base_url="http://127.0.0.1:11434")
        mock_response = MagicMock(status_code=200)
        with patch("httpx.AsyncClient.get", new_callable=AsyncMock) as mock_get:
            mock_get.return_value = mock_response
            is_healthy = await client.check_health()
            assert is_healthy is True

    asyncio.run(run())


def test_ollama_client_health_check_failure():
    async def run():
        client = OllamaClient(base_url="http://127.0.0.1:11434")
        with patch(
            "httpx.AsyncClient.get",
            new_callable=AsyncMock,
            side_effect=Exception("Connection refused"),
        ):
            is_healthy = await client.check_health()
            assert is_healthy is False

    asyncio.run(run())


def test_ollama_client_list_models():
    async def run():
        client = OllamaClient(base_url="http://127.0.0.1:11434")
        mock_response = MagicMock(status_code=200)
        mock_response.json.return_value = {
            "models": [{"name": "qwen3:4b"}, {"name": "huntiq-assistant:latest"}]
        }
        with patch("httpx.AsyncClient.get", new_callable=AsyncMock) as mock_get:
            mock_get.return_value = mock_response
            models = await client.list_models()
            assert "qwen3:4b" in models
            assert "huntiq-assistant:latest" in models

    asyncio.run(run())


def test_ollama_client_chat_payload_structure():
    async def run():
        client = OllamaClient(base_url="http://127.0.0.1:11434")
        mock_response = MagicMock(status_code=200)
        mock_response.json.return_value = {
            "message": {"role": "assistant", "content": "Hello, I am HUNTIQ assistant."}
        }
        with patch("httpx.AsyncClient.post", new_callable=AsyncMock) as mock_post:
            mock_post.return_value = mock_response

            res = await client.chat(
                messages=[{"role": "user", "content": "Hello"}],
                model="huntiq-assistant:latest",
                locale="en",
            )
            assert res["message"]["content"] == "Hello, I am HUNTIQ assistant."

            # Verify system prompt was injected and tools attached
            mock_post.assert_called_once()
            _, kwargs = mock_post.call_args
            json_payload = kwargs.get("json", {})
            assert json_payload["model"] == "huntiq-assistant:latest"
            assert len(json_payload["messages"]) >= 2
            assert json_payload["messages"][0]["role"] == "system"
            assert "tools" in json_payload

    asyncio.run(run())
