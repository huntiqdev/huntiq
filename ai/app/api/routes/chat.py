"""Chat API route handlers for HUNTIQ AI service."""

from fastapi import APIRouter, Header, HTTPException, status

from ...api.schemas import (
    AuthenticatedChatMessageRequest,
    AuthenticatedChatMessageResponse,
    GuestChatRequest,
    GuestChatResponse,
)
from ...config import config
from ...core.orchestrator import orchestrator

router = APIRouter(prefix="/api/v1/chat", tags=["chat"])


@router.post("", response_model=GuestChatResponse)
async def handle_guest_chat(
    req: GuestChatRequest, x_huntiq_ai_secret: str | None = Header(default=None)
) -> GuestChatResponse:
    """
    Public and internal conversational intelligence endpoint.
    If internal secret is configured, validates authorized calls from Go backend or Next.js proxy.
    """
    if (
        config.api_secret_key
        and x_huntiq_ai_secret
        and x_huntiq_ai_secret != config.api_secret_key
    ):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid service credentials",
        )

    return await orchestrator.process_chat(req)


@router.post("/message", response_model=AuthenticatedChatMessageResponse)
async def handle_authenticated_chat_message(
    req: AuthenticatedChatMessageRequest,
    x_huntiq_ai_secret: str | None = Header(default=None),
) -> AuthenticatedChatMessageResponse:
    """
    Authenticated conversational endpoint with tool calling.
    Called by Go backend aiclient SendChatMessage / sendPythonChat.
    """
    if (
        config.api_secret_key
        and x_huntiq_ai_secret
        and x_huntiq_ai_secret != config.api_secret_key
    ):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid service credentials",
        )

    return await orchestrator.process_chat_message(req)
