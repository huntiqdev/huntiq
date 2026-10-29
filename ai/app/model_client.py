"""Minimal, robust Python client for interacting with the local Ollama instance."""

from typing import Any

import httpx
from pydantic import BaseModel

from .config import config
from .prompts import MARKETPLACE_TOOLS, build_system_prompt


class ChatMessage(BaseModel):
    role: str
    content: str
    tool_calls: list[dict[str, Any]] | None = None


class OllamaClient:
    """HTTP client for local Ollama daemon."""

    def __init__(self, base_url: str | None = None, timeout: float | None = None):
        self.base_url = (base_url or config.ollama_base_url).rstrip("/")
        self.timeout = timeout or config.request_timeout

    async def check_health(self) -> bool:
        """Check if local Ollama daemon is reachable."""
        try:
            async with httpx.AsyncClient(timeout=3.0) as client:
                res = await client.get(f"{self.base_url}/api/tags")
                return res.status_code == 200
        except Exception:  # noqa: BLE001
            return False

    async def list_models(self) -> list[str]:
        """Fetch list of locally installed model tags."""
        try:
            async with httpx.AsyncClient(timeout=5.0) as client:
                res = await client.get(f"{self.base_url}/api/tags")
                if res.status_code == 200:
                    data = res.json()
                    return [m.get("name", "") for m in data.get("models", [])]
        except Exception:  # noqa: BLE001
            return []
        return []

    async def chat(
        self,
        messages: list[dict[str, str]],
        model: str | None = None,
        locale: str = "en",
        tools: list[dict[str, Any]] | None = None,
    ) -> dict[str, Any]:
        """Send chat completion request to local Ollama."""
        target_model = model or config.huntiq_model
        system_prompt = build_system_prompt(locale)

        # Prepend system instruction if not already present
        full_messages = []
        if not any(m.get("role") == "system" for m in messages):
            full_messages.append({"role": "system", "content": system_prompt})
        full_messages.extend(messages)

        payload: dict[str, Any] = {
            "model": target_model,
            "messages": full_messages,
            "stream": False,
            "options": {
                "num_ctx": config.num_ctx,
            },
        }

        if tools is not None:
            payload["tools"] = tools
        else:
            payload["tools"] = MARKETPLACE_TOOLS

        async with httpx.AsyncClient(timeout=self.timeout) as client:
            res = await client.post(f"{self.base_url}/api/chat", json=payload)
            if res.status_code != 200:
                raise RuntimeError(
                    f"Ollama chat error (HTTP {res.status_code}): {res.text}"
                )
            return res.json()
