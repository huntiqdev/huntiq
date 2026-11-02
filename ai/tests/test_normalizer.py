"""Tests for Ethiopian language normalizer and language detection."""

from app.languages.normalizer import (
    contains_geez,
    detect_language,
    normalize_geez_text,
    normalize_oromo_text,
)
from app.languages.terminology import (
    resolve_category_slug,
    resolve_transaction_type,
)


def test_contains_geez():
    assert contains_geez("ሰላም") is True
    assert contains_geez("hello world") is False
    assert contains_geez("Hello ሰላም") is True


def test_normalize_geez_homophones_and_punct():
    # ሐመር ሀ variant canonicalization
    raw = "ሐመር፡መኪና።"
    normalized = normalize_geez_text(raw)
    assert "ሀ" in normalized or "ሐ" not in normalized
    assert "፡" not in normalized
    assert " " in normalized


def test_normalize_oromo():
    raw = "xaa`oo fi horii"
    normalized = normalize_oromo_text(raw)
    assert "`" not in normalized
    assert "'" in normalized


def test_detect_language():
    assert detect_language("Hello, I need a car") == "en"
    assert detect_language("ሰላም እንደምን ነህ") == "am"
    assert detect_language("Akkam jirtu bilbila barbaada") == "om"
    assert detect_language("እንታይ ዓይነት መኪና ኣለኩም") == "ti"


def test_resolve_category_slug():
    assert resolve_category_slug("smartphone") == "phones-telecommunications"
    assert resolve_category_slug("ስልክ") == "phones-telecommunications"
    assert resolve_category_slug("bilbila") == "phones-telecommunications"
    assert resolve_category_slug("መኪና") == "vehicles-transport-automotive"
    assert resolve_category_slug("mana") == "real-estate-property-housing"


def test_resolve_transaction_type():
    assert resolve_transaction_type("rent") == "rent"
    assert resolve_transaction_type("ኪራይ") == "rent"
    assert resolve_transaction_type("kiraa") == "rent"
    assert resolve_transaction_type("buy") == "sale"
    assert resolve_transaction_type("መግዛት") == "sale"
