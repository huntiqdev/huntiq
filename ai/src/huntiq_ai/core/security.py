import secrets

from fastapi import Header, Request

from huntiq_ai.core.config import get_settings
from huntiq_ai.core.errors import AuthenticationError


async def verify_internal_secret(
    request: Request,
    x_huntiq_ai_secret: str | None = Header(default=None, alias="X-HUNTIQ-AI-SECRET"),
) -> str:
    """Verifies that incoming request contains the valid pre-shared inter-service secret.

    Uses constant-time comparison to prevent timing side-channel attacks.
    """
    if not x_huntiq_ai_secret:
        raise AuthenticationError("Missing required 'X-HUNTIQ-AI-SECRET' authentication header")

    settings = get_settings()
    expected_secret = settings.huntiq_ai_secret_key

    # Constant-time comparison
    if not secrets.compare_digest(x_huntiq_ai_secret.strip(), expected_secret.strip()):
        raise AuthenticationError("Invalid internal service authentication secret")

    return x_huntiq_ai_secret
