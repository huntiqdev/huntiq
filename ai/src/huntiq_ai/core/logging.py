import json
import logging
import sys
from datetime import UTC, datetime
from typing import Any

from huntiq_ai.core.config import get_settings

SENSITIVE_KEYS = {
    "authorization",
    "x-huntiq-ai-secret",
    "password",
    "secret",
    "token",
    "api_key",
    "cv",
    "phone",
}


class JSONFormatter(logging.Formatter):
    """Formats log records as single-line JSON with redacted sensitive data."""

    def format(self, record: logging.LogRecord) -> str:
        log_entry: dict[str, Any] = {
            "timestamp": datetime.now(UTC).isoformat(),
            "level": record.levelname,
            "logger": record.name,
            "message": record.getMessage(),
        }

        # Include correlation ID if attached to record
        if hasattr(record, "correlation_id") and record.correlation_id:
            log_entry["correlation_id"] = record.correlation_id

        # Include exception trace if present
        if record.exc_info:
            log_entry["exception"] = self.formatException(record.exc_info)

        return json.dumps(log_entry)


def sanitize_dict(data: dict[str, Any]) -> dict[str, Any]:
    """Recursively redacts sensitive keys from log dictionaries."""
    sanitized: dict[str, Any] = {}
    for k, v in data.items():
        if k.lower() in SENSITIVE_KEYS:
            sanitized[k] = "[REDACTED]"
        elif isinstance(v, dict):
            sanitized[k] = sanitize_dict(v)
        else:
            sanitized[k] = v
    return sanitized


def setup_logging() -> None:
    """Configures root and application loggers based on settings."""
    settings = get_settings()
    level = getattr(logging, settings.log_level.upper(), logging.INFO)

    root_logger = logging.getLogger()
    root_logger.setLevel(level)

    # Remove existing handlers to avoid duplicates
    for handler in list(root_logger.handlers):
        root_logger.removeHandler(handler)

    handler = logging.StreamHandler(sys.stdout)
    handler.setLevel(level)

    if settings.log_format == "json":
        handler.setFormatter(JSONFormatter())
    else:
        handler.setFormatter(
            logging.Formatter(
                fmt="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
                datefmt="%Y-%m-%d %H:%M:%S",
            )
        )

    root_logger.addHandler(handler)


def get_logger(name: str) -> logging.Logger:
    """Returns a named logger configured for the platform."""
    return logging.getLogger(name)
