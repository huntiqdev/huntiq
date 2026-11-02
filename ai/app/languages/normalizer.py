"""Text normalization and script preprocessing for Ethiopian languages.

Supports English, Amharic (አማርኛ), Afaan Oromoo, and Tigrinya (ትግርኛ).
"""

import re
import unicodedata

# Ge'ez Homophone Normalization Mapping (Amharic standard phonetic canonicalization)
GEEZ_HOMOPHONE_MAP = {
    # Ha variants -> 1st order ሀ
    "ሐ": "ሀ",
    "ኀ": "ሀ",
    "ሃ": "ሀ",
    "ሁ": "ሁ",
    "ሑ": "ሁ",
    "ኁ": "ሁ",
    "ሂ": "ሂ",
    "ሒ": "ሂ",
    "ኺ": "ሂ",
    "ሄ": "ሄ",
    "ሔ": "ሄ",
    "ኼ": "ሄ",
    "ህ": "ህ",
    "ሕ": "ህ",
    "ኽ": "ህ",
    "ሆ": "ሆ",
    "ሖ": "ሆ",
    "ኆ": "ሆ",
    # Se variants -> ሰ
    "ሠ": "ሰ",
    "ሡ": "ሱ",
    "ሢ": "ሲ",
    "ሣ": "ሳ",
    "ሤ": "ሴ",
    "ሥ": "ስ",
    "ሦ": "ሶ",
    # Tse variants -> ጸ
    "ፀ": "ጸ",
    "ፁ": "ጹ",
    "ፂ": "ጺ",
    "ፃ": "ጻ",
    "ፄ": "ጼ",
    "ፅ": "ጽ",
    "ፆ": "ጾ",
    # A variants -> አ
    "ዐ": "አ",
    "ዑ": "ኡ",
    "ዒ": "ኢ",
    "ዓ": "ኣ",
    "ዔ": "ኤ",
    "ዕ": "እ",
    "ዖ": "ኦ",
}

# Ge'ez Punctuation to ASCII Mapping
GEEZ_PUNCT_MAP = {
    "፡": " ",  # Word separator
    "።": ".",  # Full stop
    "፣": ",",  # Comma
    "፤": ";",  # Semicolon
    "፥": ":",  # Colon
    "፦": ":",  # Preface colon
    "፧": "?",  # Question mark
    "፨": " ",  # Paragraph separator
}

GEEZ_RANGE = re.compile(r"[\u1200-\u137F\u1380-\u139F\u2D80-\u2DDF]")


def contains_geez(text: str) -> bool:
    """Check if the string contains any Ge'ez script characters."""
    return bool(GEEZ_RANGE.search(text))


def normalize_geez_text(text: str) -> str:
    """Canonicalize Ge'ez characters, unify homophones, and normalize punctuation."""
    if not text:
        return ""

    # Normalize unicode composition
    normalized = unicodedata.normalize("NFKC", text)

    # Replace Ge'ez punctuation
    for punct, replacement in GEEZ_PUNCT_MAP.items():
        normalized = normalized.replace(punct, replacement)

    # Normalize homophones
    result = []
    for char in normalized:
        result.append(GEEZ_HOMOPHONE_MAP.get(char, char))

    # Collapse multiple whitespaces
    return re.sub(r"\s+", " ", "".join(result)).strip()


def normalize_oromo_text(text: str) -> str:
    """Normalize Afaan Oromoo orthography: unified apostrophes and whitespace."""
    if not text:
        return ""

    # Normalize glottal stop / hudhaa apostrophe variants
    text = re.sub(r"[`'ʻʼ’]", "'", text)
    return re.sub(r"\s+", " ", text).strip()


def detect_language(text: str, hint_locale: str = "en") -> str:
    """Heuristic language detection for HUNTIQ supported locales."""
    clean = text.strip()
    if not clean:
        return hint_locale

    if contains_geez(clean):
        # Tigrinya specific vocabulary heuristics
        tigrinya_markers = ["እንታይ", "ኣበይ", "ከመይ", "ኣለኩም", "ክንደይ", "ድልየት", "ናይ"]
        for marker in tigrinya_markers:
            if marker in clean:
                return "ti"
        return "am"

    lower = clean.lower()
    # Oromo specific vocabulary heuristics
    oromo_markers = [
        "akkam",
        "maal",
        "meeqa",
        "eessatti",
        "barbaada",
        "jira",
        "biyya",
        "mana",
        "konkolaataa",
        "bilbila",
        "qarshii",
    ]
    words = set(re.findall(r"\b\w+\b", lower))
    if any(marker in words for marker in oromo_markers):
        return "om"

    return "en"


def normalize_text(text: str) -> str:
    """Universal text normalizer for Ge'ez, Oromo, and English."""
    if not text:
        return ""
    if contains_geez(text):
        return normalize_geez_text(text)
    return normalize_oromo_text(text)


class TextNormalizer:
    @staticmethod
    def normalize(text: str) -> str:
        return normalize_text(text)


normalizer = TextNormalizer()
