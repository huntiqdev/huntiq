"""Languages package for Ethiopian multilingual processing."""

from .normalizer import (
    contains_geez,
    detect_language,
    normalize_geez_text,
    normalize_oromo_text,
)
from .terminology import (
    CATEGORY_SYNONYM_MAP,
    TRANSACTION_TYPE_MAP,
    resolve_category_slug,
    resolve_transaction_type,
)

__all__ = [
    "CATEGORY_SYNONYM_MAP",
    "TRANSACTION_TYPE_MAP",
    "contains_geez",
    "detect_language",
    "normalize_geez_text",
    "normalize_oromo_text",
    "resolve_category_slug",
    "resolve_transaction_type",
]
