"""Common shared schemas, enums, and models for the HUNTIQ AI platform."""

from __future__ import annotations

from datetime import datetime
from enum import StrEnum
from typing import Any

from pydantic import BaseModel, Field


class Locale(StrEnum):
    """Canonical supported language locales across HUNTIQ services."""

    EN = "en"  # English
    AM = "am"  # Amharic (አማርኛ)
    OM = "om"  # Afaan Oromoo
    TI = "ti"  # Tigrinya (ትግርኛ)


class QueryIntent(StrEnum):
    """Canonical search query intent categories defined by contract."""

    BUY_PRODUCT = "buy_product"
    RENT_PROPERTY = "rent_property"
    FIND_SERVICE = "find_service"
    JOB_SEARCH = "job_search"
    VEHICLE_SEARCH = "vehicle_search"
    PRICE_CHECK = "price_check"
    AMBIGUOUS = "ambiguous"


class ErrorDetail(BaseModel):
    """Structured error details conforming to RFC 7807 problem representation."""

    code: str = Field(..., description="Machine-readable application error code")
    message: str = Field(..., description="Human-readable explanation of the error")
    correlation_id: str = Field(..., description="Distributed tracing correlation identifier")
    timestamp: datetime = Field(..., description="UTC timestamp of error occurrence")
    details: dict[str, Any] | None = Field(
        default=None,
        description="Optional diagnostic details or field-level validation errors",
    )


class ErrorResponse(BaseModel):
    """Top-level error response envelope."""

    error: ErrorDetail = Field(..., description="Encapsulated error details")


class ResponseMetadata(BaseModel):
    """Standard response metadata for tracking, latency, and telemetry."""

    correlation_id: str = Field(..., description="Correlation ID for distributed tracing")
    latency_ms: float | None = Field(
        default=None,
        description="Elapsed inference compute time in milliseconds",
    )
