from typing import Any

from pydantic import BaseModel, Field


class LivenessResponse(BaseModel):
    """Schema returned by liveness probe endpoint."""

    status: str = Field(default="live", description="Process liveness state")
    service: str = Field(default="huntiq-ai", description="Service identifier")
    version: str = Field(default="0.1.0", description="Service version")


class ReadinessResponse(BaseModel):
    """Schema returned by readiness probe endpoint."""

    status: str = Field(default="ready", description="Service readiness state")
    service: str = Field(default="huntiq-ai", description="Service identifier")
    version: str = Field(default="0.1.0", description="Service version")
    environment: str = Field(..., description="Current deployment environment")
    models_loaded: bool = Field(default=False, description="Whether ML models are loaded in memory")
    details: dict[str, Any] = Field(
        default_factory=dict,
        description="Non-sensitive operational telemetry details",
    )
