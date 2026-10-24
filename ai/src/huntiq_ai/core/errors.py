from datetime import UTC, datetime
from typing import Any

from fastapi import Request, status
from fastapi.responses import JSONResponse


class AIPlatformError(Exception):
    """Base exception for all HUNTIQ AI platform errors."""

    def __init__(
        self,
        message: str,
        code: str = "INTERNAL_ERROR",
        status_code: int = status.HTTP_500_INTERNAL_SERVER_ERROR,
        details: dict[str, Any] | None = None,
    ):
        super().__init__(message)
        self.message = message
        self.code = code
        self.status_code = status_code
        self.details = details or {}


class AuthenticationError(AIPlatformError):
    """Raised when internal service authentication fails or is missing."""

    def __init__(self, message: str = "Missing or invalid internal service authentication secret"):
        super().__init__(
            message=message,
            code="UNAUTHORIZED",
            status_code=status.HTTP_401_UNAUTHORIZED,
        )


class PayloadTooLargeError(AIPlatformError):
    """Raised when request payload exceeds configured maximum size."""

    def __init__(self, message: str = "Request payload exceeds maximum allowed size"):
        super().__init__(
            message=message,
            code="PAYLOAD_TOO_LARGE",
            status_code=status.HTTP_413_REQUEST_ENTITY_TOO_LARGE,
        )


class ModelUnavailableError(AIPlatformError):
    """Raised when a requested model is not loaded or unavailable."""

    def __init__(self, model_name: str):
        super().__init__(
            message=f"Model '{model_name}' is not loaded or currently unavailable",
            code="MODEL_UNAVAILABLE",
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            details={"model": model_name},
        )


def format_error_response(
    request: Request,
    message: str,
    code: str,
    status_code: int,
    details: dict[str, Any] | None = None,
) -> JSONResponse:
    """Creates a consistent RFC 7807 compliant JSON error response."""
    correlation_id = getattr(request.state, "correlation_id", "unknown")
    content: dict[str, Any] = {
        "error": {
            "code": code,
            "message": message,
            "correlation_id": correlation_id,
            "timestamp": datetime.now(UTC).isoformat(),
        }
    }
    if details:
        content["error"]["details"] = details

    return JSONResponse(status_code=status_code, content=content)


async def ai_platform_exception_handler(request: Request, exc: AIPlatformError) -> JSONResponse:
    """FastAPI exception handler for AIPlatformError."""
    return format_error_response(
        request=request,
        message=exc.message,
        code=exc.code,
        status_code=exc.status_code,
        details=exc.details,
    )
