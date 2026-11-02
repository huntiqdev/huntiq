"""Restricted HTTP client communicating with Go backend's authorized marketplace tools."""

from typing import Any

import httpx

from ..api.schemas import (
    InternalToolExecutionRequest,
    InternalToolExecutionResponse,
    ListingSummary,
)
from ..config import config


class GoMarketplaceClient:
    """Restricted marketplace tool client for internal service-to-service execution."""

    def __init__(self, base_url: str | None = None, secret_key: str | None = None):
        self.base_url = (base_url or config.go_backend_url).rstrip("/")
        self.secret_key = secret_key or config.api_secret_key
        self.endpoint = f"{self.base_url}/api/internal/ai/tools/execute"

    async def execute_tool(
        self, tool_name: str, arguments: dict[str, Any]
    ) -> InternalToolExecutionResponse:
        """Call Go backend internal tool endpoint with strict authentication and timeout."""
        if tool_name not in config.allowed_tools:
            return InternalToolExecutionResponse(
                success=False,
                tool_name=tool_name,
                error=f"Unauthorized or unsupported tool '{tool_name}'",
            )

        payload = InternalToolExecutionRequest(
            tool_name=tool_name, arguments=arguments
        ).model_dump()

        headers = {
            "X-HUNTIQ-AI-SECRET": self.secret_key,
            "Content-Type": "application/json",
            "Accept": "application/json",
        }

        try:
            async with httpx.AsyncClient(timeout=3.0) as client:
                res = await client.post(self.endpoint, json=payload, headers=headers)
                if res.status_code == 200:
                    data = res.json()
                    listings = [
                        ListingSummary(**item) for item in data.get("listings", [])
                    ]
                    return InternalToolExecutionResponse(
                        success=True,
                        tool_name=tool_name,
                        total_matching=data.get("total_matching", len(listings)),
                        listings=listings,
                    )
                return InternalToolExecutionResponse(
                    success=False,
                    tool_name=tool_name,
                    error=f"Go backend returned status {res.status_code}",
                )
        except Exception as e:  # noqa: BLE001
            # Graceful failure when Go backend is offline or starting up
            return InternalToolExecutionResponse(
                success=False,
                tool_name=tool_name,
                error=f"Go backend connection failed: {e!s}",
            )

    async def search_marketplace(
        self,
        search_term: str,
        category_slug: str | None = None,
        min_price: float | None = None,
        max_price: float | None = None,
        location: str | None = None,
        condition: str | None = None,
        listing_type: str | None = None,
        limit: int = 5,
    ) -> list[ListingSummary]:
        """Execute search_marketplace_listings tool against Go backend."""
        args: dict[str, Any] = {"search_term": search_term, "limit": limit}
        if category_slug:
            args["category_slug"] = category_slug
        if min_price is not None:
            args["min_price"] = min_price
        if max_price is not None:
            args["max_price"] = max_price
        if location:
            args["location"] = location
        if condition:
            args["condition"] = condition
        if listing_type:
            args["listing_type"] = listing_type

        resp = await self.execute_tool("search_marketplace_listings", args)
        if resp.success:
            return resp.listings
        return []


# Global singleton instance
go_client = GoMarketplaceClient()
