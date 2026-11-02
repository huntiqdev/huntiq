"""Tests for FastAPI HTTP endpoints using ASGITransport."""

import httpx
import pytest

from app.main import app


@pytest.mark.asyncio
async def test_health_endpoint():
    """Verify GET /health returns 200 and healthy status."""
    transport = httpx.ASGITransport(app=app)
    async with httpx.AsyncClient(transport=transport, base_url="http://test") as client:
        resp = await client.get("/health")
        assert resp.status_code == 200
        data = resp.json()
        assert data["status"] == "ok"
        assert data["service"] == "huntiq-ai"


@pytest.mark.asyncio
async def test_ready_endpoint():
    """Verify GET /ready returns 200."""
    transport = httpx.ASGITransport(app=app)
    async with httpx.AsyncClient(transport=transport, base_url="http://test") as client:
        resp = await client.get("/ready")
        assert resp.status_code == 200
        data = resp.json()
        assert data["ready"] is True


@pytest.mark.asyncio
async def test_api_v1_chat_greeting():
    """Verify POST /api/v1/chat processes greetings without DB searches."""
    transport = httpx.ASGITransport(app=app)
    async with httpx.AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {"message": "Hello! What can you do?", "locale": "en"}
        resp = await client.post("/api/v1/chat", json=payload)
        assert resp.status_code == 200
        data = resp.json()
        assert "answer" in data
        assert data["detected_locale"] == "en"
        assert data["session_id"] is not None


@pytest.mark.asyncio
async def test_api_v1_chat_out_of_domain():
    """Verify POST /api/v1/chat politely redirects out-of-domain requests."""
    transport = httpx.ASGITransport(app=app)
    async with httpx.AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {"message": "Write a python function for fibonacci", "locale": "en"}
        resp = await client.post("/api/v1/chat", json=payload)
        assert resp.status_code == 200
        data = resp.json()
        assert "out_of_domain" in data["intents"]
        assert "marketplace" in data["answer"].lower()


@pytest.mark.asyncio
async def test_api_v1_query_understand_endpoint():
    """Verify POST /api/v1/query/understand endpoint provides full taxonomy and entity extraction."""
    from app.config import config

    transport = httpx.ASGITransport(app=app)
    headers = {"X-HUNTIQ-AI-SECRET": config.api_secret_key}
    async with httpx.AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {
            "query": "Toyota Corolla Bole under 500000 birr",
            "locale": "en",
            "mode": "ai",
        }
        resp = await client.post(
            "/api/v1/query/understand", json=payload, headers=headers
        )
        assert resp.status_code == 200
        data = resp.json()
        assert data["is_safe"] is True
        assert "normalized_query" in data
        assert data["primary_intent"] in ("buy_search", "marketplace_search")
        assert 2 in data["suggested_category_ids"]  # Vehicles
        assert data["extracted_entities"]["brand"] == "Toyota"
        assert data["extracted_entities"]["price_range"]["max_price"] == 500000.0


@pytest.mark.asyncio
async def test_api_v1_query_analyze_products_endpoint():
    """Verify POST /api/v1/query/analyze-products calculates statistical market metrics."""
    from app.config import config

    transport = httpx.ASGITransport(app=app)
    headers = {"X-HUNTIQ-AI-SECRET": config.api_secret_key}
    async with httpx.AsyncClient(transport=transport, base_url="http://test") as client:
        payload = {
            "query": "laptop",
            "products": [
                {
                    "id": 1,
                    "title": "Dell XPS 13",
                    "price": 45000.0,
                    "seller_tier": "premium",
                },
                {
                    "id": 2,
                    "title": "HP Envy 15",
                    "price": 35000.0,
                    "seller_tier": "normal",
                },
            ],
        }
        resp = await client.post(
            "/api/v1/query/analyze-products", json=payload, headers=headers
        )
        assert resp.status_code == 200
        data = resp.json()
        assert data["total_analyzed"] == 2
        assert data["premium_count"] == 1
        assert data["price_stats"]["min_price"] == 35000.0
        assert data["price_stats"]["max_price"] == 45000.0
        assert data["price_stats"]["avg_price"] == 40000.0


@pytest.mark.asyncio
async def test_api_v1_chat_message_tool_calling_and_synthesis():
    """Verify POST /api/v1/chat/message returns tool calls for search and synthesizes tool responses."""
    from app.config import config

    transport = httpx.ASGITransport(app=app)
    headers = {"X-HUNTIQ-AI-SECRET": config.api_secret_key}
    async with httpx.AsyncClient(transport=transport, base_url="http://test") as client:
        # 1. Search message generates tool call
        search_payload = {
            "messages": [
                {"role": "user", "content": "I want to buy a used iPhone 13 in Addis"}
            ],
            "locale": "en",
            "tools_enabled": True,
        }
        resp = await client.post(
            "/api/v1/chat/message", json=search_payload, headers=headers
        )
        assert resp.status_code == 200
        data = resp.json()
        assert data["finish_reason"] == "tool_calls"
        assert len(data["tool_calls"]) > 0
        call = data["tool_calls"][0]
        assert call["function"]["name"] == "search_marketplace_listings"
        assert "search_term" in call["function"]["arguments"]

        # 2. Tool response synthesizes final grounded answer
        tool_payload = {
            "messages": [
                {"role": "user", "content": "I want to buy a used iPhone 13 in Addis"},
                {
                    "role": "assistant",
                    "content": "",
                    "tool_calls": [call],
                },
                {
                    "role": "tool",
                    "tool_call_id": call["id"],
                    "content": '{"total_matches": 1, "listings": [{"id": 10, "title": "iPhone 13 Pro 128GB", "price": 65000, "city": "Addis Ababa", "seller_tier": "premium"}]}',
                },
            ],
            "locale": "en",
            "tools_enabled": True,
        }
        resp2 = await client.post(
            "/api/v1/chat/message", json=tool_payload, headers=headers
        )
        assert resp2.status_code == 200
        data2 = resp2.json()
        assert data2["finish_reason"] == "stop"
        assert data2["content"] is not None
        assert "iPhone 13 Pro" in data2["content"]
        assert "65,000" in data2["content"]
