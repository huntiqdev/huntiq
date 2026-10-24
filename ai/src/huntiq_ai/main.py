import uuid
from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager

from fastapi import FastAPI, Request, Response
from starlette.middleware.base import BaseHTTPMiddleware, RequestResponseEndpoint

from huntiq_ai.api.router import api_router
from huntiq_ai.core.config import Settings, get_settings
from huntiq_ai.core.errors import (
    AIPlatformError,
    ai_platform_exception_handler,
    format_error_response,
)
from huntiq_ai.core.logging import get_logger, setup_logging

logger = get_logger("huntiq_ai.main")


class CorrelationIdMiddleware(BaseHTTPMiddleware):
    """Ensures every request has a correlation ID for distributed tracing."""

    async def dispatch(self, request: Request, call_next: RequestResponseEndpoint) -> Response:
        correlation_id = request.headers.get("X-Correlation-ID") or f"cid-{uuid.uuid4().hex[:12]}"
        request.state.correlation_id = correlation_id

        response = await call_next(request)
        response.headers["X-Correlation-ID"] = correlation_id
        return response


class PayloadLimitMiddleware(BaseHTTPMiddleware):
    """Enforces maximum incoming request payload size."""

    def __init__(self, app, max_size_bytes: int | None = None):
        super().__init__(app)
        self.max_size_bytes = max_size_bytes

    async def dispatch(self, request: Request, call_next: RequestResponseEndpoint) -> Response:
        limit = (
            self.max_size_bytes
            if self.max_size_bytes is not None
            else get_settings().max_payload_size_bytes
        )
        content_length = request.headers.get("content-length")
        if content_length and int(content_length) > limit:
            return format_error_response(
                request=request,
                message=(
                    f"Payload size {content_length} bytes exceeds "
                    f"maximum allowed limit of {limit} bytes"
                ),
                code="PAYLOAD_TOO_LARGE",
                status_code=413,
            )
        return await call_next(request)


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    """Application lifespan context manager for startup and shutdown."""
    setup_logging()
    settings = get_settings()
    logger.info(
        "Starting HUNTIQ AI microservice",
        extra={
            "environment": settings.environment,
            "port": settings.port,
        },
    )
    yield
    logger.info("Shutting down HUNTIQ AI microservice")


def create_app(settings: Settings | None = None) -> FastAPI:
    """Creates and configures the FastAPI application instance."""
    if settings is None:
        settings = get_settings()

    # Disable interactive docs in production to prevent schema leakage
    docs_url = "/docs" if settings.environment != "production" else None
    redoc_url = "/redoc" if settings.environment != "production" else None
    openapi_url = "/openapi.json" if settings.environment != "production" else None

    app = FastAPI(
        title="HUNTIQ AI Service",
        description="Internal prediction and classification worker for HUNTIQ marketplace",
        version="0.1.0",
        docs_url=docs_url,
        redoc_url=redoc_url,
        openapi_url=openapi_url,
        lifespan=lifespan,
    )

    # Middleware Registration: in Starlette, outermost runs first.
    # We add PayloadLimitMiddleware first, then CorrelationIdMiddleware second
    # so CorrelationIdMiddleware wraps PayloadLimitMiddleware.
    app.add_middleware(
        PayloadLimitMiddleware,
        max_size_bytes=settings.max_payload_size_bytes,
    )
    app.add_middleware(CorrelationIdMiddleware)

    # Exception Handlers
    app.add_exception_handler(AIPlatformError, ai_platform_exception_handler)

    # Router Registration
    app.include_router(api_router)

    return app


app = create_app()
