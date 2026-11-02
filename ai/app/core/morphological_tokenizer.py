"""
HUNTIQ Multilingual Morphological & Letter-Level Tokenizer.
Provides deep subword, letter-by-letter, and morphological decomposition for:
- Amharic (AM)
- Tigrinya (TI)
- Afaan Oromoo (OM)
- English (EN)

Deconstructs words down to consonant roots, Fidel vowel orders, agglutinative affixes,
and character n-grams to capture the finest semantic and orthographic nuances.
"""

from __future__ import annotations

import re
from typing import Any

# Ge'ez Fidel Root Mapping: Map 1st-7th order fidels to base consonant root
GEEZ_CONSONANT_MAP: dict[str, str] = {
    # ሀ ሁ ሂ ሃ ሄ ህ ሆ
    "ሁ": "ሀ",
    "ሂ": "ሀ",
    "ሃ": "ሀ",
    "ሄ": "ሀ",
    "ህ": "ሀ",
    "ሆ": "ሀ",
    # ለ ሉ ሊ ላ ሌ ል ሎ ሏ
    "ሉ": "ለ",
    "ሊ": "ለ",
    "ላ": "ለ",
    "ሌ": "ለ",
    "ል": "ለ",
    "ሎ": "ለ",
    "ሏ": "ለ",
    # ሐ ሑ ሒ ሓ ሔ ሕ ሖ ሗ
    "ሑ": "ሐ",
    "ሒ": "ሐ",
    "ሓ": "ሐ",
    "ሔ": "ሐ",
    "ሕ": "ሐ",
    "ሖ": "ሐ",
    "ሗ": "ሐ",
    # መ ሙ ሚ ማ ሜ ም ሞ ሟ
    "ሙ": "መ",
    "ሚ": "መ",
    "ማ": "መ",
    "ሜ": "መ",
    "ም": "መ",
    "ሞ": "መ",
    "ሟ": "መ",
    # ሠ ሡ ሢ ሣ ሤ ሥ ሦ ሧ
    "ሡ": "ሠ",
    "ሢ": "ሠ",
    "ሣ": "ሠ",
    "ሤ": "ሠ",
    "ሥ": "ሠ",
    "ሦ": "ሠ",
    "ሧ": "ሠ",
    # ረ ሩ ሪ ራ ሬ ር ሮ ሯ
    "ሩ": "ረ",
    "ሪ": "ረ",
    "ራ": "ረ",
    "ሬ": "ረ",
    "ር": "ረ",
    "ሮ": "ረ",
    "ሯ": "ረ",
    # ሰ ሱ ሲ ሳ ሴ ስ ሶ ሷ
    "ሱ": "ሰ",
    "ሲ": "ሰ",
    "ሳ": "ሰ",
    "ሴ": "ሰ",
    "ስ": "ሰ",
    "ሶ": "ሰ",
    "ሷ": "ሰ",
    # ሸ ሹ ሺ ሻ ሼ ሽ ሾ ሿ
    "ሹ": "ሸ",
    "ሺ": "ሸ",
    "ሻ": "ሸ",
    "ሼ": "ሸ",
    "ሽ": "ሸ",
    "ሾ": "ሸ",
    "ሿ": "ሸ",
    # ቀ ቁ ቂ ቃ ቄ ቅ ቆ ቋ
    "ቁ": "ቀ",
    "ቂ": "ቀ",
    "ቃ": "ቀ",
    "ቄ": "ቀ",
    "ቅ": "ቀ",
    "ቆ": "ቀ",
    "ቋ": "ቀ",
    # በ ቡ ቢ ባ ቤ ብ ቦ ቧ
    "ቡ": "በ",
    "ቢ": "በ",
    "ባ": "በ",
    "ቤ": "በ",
    "ብ": "በ",
    "ቦ": "በ",
    "ቧ": "በ",
    # ተ ቱ ቲ ታ ቴ ት ቶ ቷ
    "ቱ": "ተ",
    "ቲ": "ተ",
    "ታ": "ተ",
    "ቴ": "ተ",
    "ት": "ተ",
    "ቶ": "ተ",
    "ቷ": "ተ",
    # ቸ ቹ ቺ ቻ ቼ ች ቾ ቿ
    "ቹ": "ቸ",
    "ቺ": "ቸ",
    "ቻ": "ቸ",
    "ቼ": "ቸ",
    "ች": "ቸ",
    "ቾ": "ቸ",
    "ቿ": "ቸ",
    # ኀ ኁ ኂ ኃ ኄ ኅ ኆ ኋ
    "ኁ": "ኀ",
    "ኂ": "ኀ",
    "ኃ": "ኀ",
    "ኄ": "ኀ",
    "ኅ": "ኀ",
    "ኆ": "ኀ",
    "ኋ": "ኀ",
    # ነ ኑ ኒ ና ኔ ን ኖ ኗ
    "ኑ": "ነ",
    "ኒ": "ነ",
    "ና": "ነ",
    "ኔ": "ነ",
    "ን": "ነ",
    "ኖ": "ነ",
    "ኗ": "ነ",
    # ኘ ኙ ኚ ኛ ኜ ኝ ኞ ኟ
    "ኙ": "ኘ",
    "ኚ": "ኘ",
    "ኛ": "ኘ",
    "ኜ": "ኘ",
    "ኝ": "ኘ",
    "ኞ": "ኘ",
    "ኟ": "ኘ",
    # አ ኡ ኢ ኣ ኤ እ ኦ ኧ
    "ኡ": "አ",
    "ኢ": "አ",
    "ኣ": "አ",
    "ኤ": "አ",
    "እ": "አ",
    "ኦ": "አ",
    "ኧ": "አ",
    # ከ ኩ ኪ ካ ኬ ክ ኮ ኳ
    "ኩ": "ከ",
    "ኪ": "ከ",
    "ካ": "ከ",
    "ኬ": "ከ",
    "ክ": "ከ",
    "ኮ": "ከ",
    "ኳ": "ከ",
    # ኸ ኹ ኺ ኻ ኼ ኽ ኾ ዃ
    "ኹ": "ኸ",
    "ኺ": "ኸ",
    "ኻ": "ኸ",
    "ኼ": "ኸ",
    "ኽ": "ኸ",
    "ኾ": "ኸ",
    "ዃ": "ኸ",
    # ወ ዉ ዊ ዋ ዌ ው ዎ
    "ዉ": "ወ",
    "ዊ": "ወ",
    "ዋ": "ወ",
    "ዌ": "ወ",
    "ው": "ወ",
    "ዎ": "ወ",
    # ዐ ዑ ዒ ዓ ዔ ዕ ዖ
    "ዑ": "ዐ",
    "ዒ": "ዐ",
    "ዓ": "ዐ",
    "ዔ": "ዐ",
    "ዕ": "ዐ",
    "ዖ": "ዐ",
    # ዘ ዙ ዚ ዛ ዜ ዝ ዞ ዟ
    "ዙ": "ዘ",
    "ዚ": "ዘ",
    "ዛ": "ዘ",
    "ዜ": "ዘ",
    "ዝ": "ዘ",
    "ዞ": "ዘ",
    "ዟ": "ዘ",
    # ዠ ዡ ዢ ዣ ዤ ዥ ዦ ዧ
    "ዡ": "ዠ",
    "ዢ": "ዠ",
    "ዣ": "ዠ",
    "ዤ": "ዠ",
    "ዥ": "ዠ",
    "ዦ": "ዠ",
    "ዧ": "ዠ",
    # የ ዩ ዪ ያ ዬ ይ ዮ ዯ
    "ዩ": "የ",
    "ዪ": "የ",
    "ያ": "የ",
    "ዬ": "የ",
    "ይ": "የ",
    "ዮ": "የ",
    "ዯ": "የ",
    # ደ ዱ ዲ ዳ ዴ ድ ዶ ዷ
    "ዱ": "ደ",
    "ዲ": "ደ",
    "ዳ": "ደ",
    "ዴ": "ደ",
    "ድ": "ደ",
    "ዶ": "ደ",
    "ዷ": "ደ",
    # ጀ ጁ ጂ ጃ ጄ ጅ ጆ ጇ
    "ጁ": "ጀ",
    "ጂ": "ጀ",
    "ጃ": "ጀ",
    "ጄ": "ጀ",
    "ጅ": "ጀ",
    "ጆ": "ጀ",
    "ጇ": "ጀ",
    # ገ ጉ ጊ ጋ ጌ ግ ጎ ጓ
    "ጉ": "ገ",
    "ጊ": "ገ",
    "ጋ": "ገ",
    "ጌ": "ገ",
    "ግ": "ገ",
    "ጎ": "ገ",
    "ጓ": "ገ",
    # ጠ ጡ ጢ ጣ ጤ ጥ ጦ ጧ
    "ጡ": "ጠ",
    "ጢ": "ጠ",
    "ጣ": "ጠ",
    "ጤ": "ጠ",
    "ጥ": "ጠ",
    "ጦ": "ጠ",
    "ጧ": "ጠ",
    # ጨ ጩ ጪ ጫ ጬ ጭ ጮ ጯ
    "ጩ": "ጨ",
    "ጪ": "ጨ",
    "ጫ": "ጨ",
    "ጬ": "ጨ",
    "ጭ": "ጨ",
    "ጮ": "ጨ",
    "ጯ": "ጨ",
    # ጰ ጱ ጲ ጳ ጴ ጵ ጶ ጷ
    "ጱ": "ጰ",
    "ጲ": "ጰ",
    "ጳ": "ጰ",
    "ጴ": "ጰ",
    "ጵ": "ጰ",
    "ጶ": "ጰ",
    "ጷ": "ጰ",
    # ጸ ጹ ጺ ጻ ጼ ጽ ጾ ጿ
    "ጹ": "ጸ",
    "ጺ": "ጸ",
    "ጻ": "ጸ",
    "ጼ": "ጸ",
    "ጽ": "ጸ",
    "ጾ": "ጸ",
    "ጿ": "ጸ",
    # ፀ ፁ ፂ ፃ ፄ ፅ ፆ ፇ
    "ፁ": "ፀ",
    "ፂ": "ፀ",
    "ፃ": "ፀ",
    "ፄ": "ፀ",
    "ፅ": "ፀ",
    "ፆ": "ፀ",
    "ፇ": "ፀ",
    # ፈ ፉ ፊ ፋ ፌ ፍ ፎ ፏ
    "ፉ": "ፈ",
    "ፊ": "ፈ",
    "ፋ": "ፈ",
    "ፌ": "ፈ",
    "ፍ": "ፈ",
    "ፎ": "ፈ",
    "ፏ": "ፈ",
    # ፐ ፑ ፒ ፓ ፔ ፕ ፖ ፗ
    "ፑ": "ፐ",
    "ፒ": "ፐ",
    "ፓ": "ፐ",
    "ፔ": "ፐ",
    "ፕ": "ፐ",
    "ፖ": "ፐ",
    "ፗ": "ፐ",
}

# Common agglutinative prefixes in Amharic & Tigrinya
GEEZ_PREFIXES: list[str] = [
    "እንደ",
    "የሚ",
    "የማ",
    "ስለ",
    "እስከ",
    "ወደ",
    "ያለ",
    "ሳይ",
    "በ",
    "ከ",
    "ለ",
    "የ",
    "ስ",
    "ዝ",
    "ክ",
    "ናይ",
]

# Common agglutinative suffixes in Amharic & Tigrinya
GEEZ_SUFFIXES: list[str] = [
    "ዎች",
    "ኦች",
    "ኣት",
    "ታት",
    "አቸው",
    "ኦቻቸው",
    "አችን",
    "ኦቻችን",
    "በት",
    "በትም",
    "ለት",
    "ለትም",
    "ናቸው",
    "አለች",
    "አሉ",
    "ም",
    "ና",
    "ኮ",
    "ሳ",
    "ኝ",
    "ህ",
    "ሽ",
    "ው",
    "ዋ",
    "ይ",
    "ታ",
]

# Common agglutinative suffixes in Afaan Oromoo
OROMO_SUFFIXES: list[str] = [
    "dhaaf",
    "ttii",
    "rratti",
    "irraa",
    "oota",
    "olii",
    "dha",
    "tti",
    "rraa",
    "f",
    "sa",
    "eera",
    "ame",
    "uu",
    "aa",
    "an",
]


class MorphologicalTokenizer:
    """
    Multilingual subword and morphological tokenizer for Ethiopian marketplace NLP.
    Understands every single letter, character n-gram, fidel root, and sentence discourse.
    """

    def __init__(self, min_ngram: int = 1, max_ngram: int = 5):
        self.min_ngram = min_ngram
        self.max_ngram = max_ngram

    @staticmethod
    def extract_letters(word: str) -> list[str]:
        """Decomposes a word into individual letters / unicode graphemes."""
        return [char for char in word if not char.isspace()]

    @staticmethod
    def extract_fidel_roots(word: str) -> list[str]:
        """Converts Ge'ez characters to their foundational 1st-order consonant roots."""
        return [GEEZ_CONSONANT_MAP.get(char, char) for char in word]

    @classmethod
    def strip_affixes_geez(cls, word: str) -> tuple[str, list[str]]:
        """
        Strips common Amharic / Tigrinya agglutinative prefixes and suffixes,
        returning the stripped root stem and detected affixes.
        """
        stem = word
        affixes_found: list[str] = []

        # Check prefixes
        for prefix in GEEZ_PREFIXES:
            if stem.startswith(prefix) and len(stem) > len(prefix) + 2:
                stem = stem[len(prefix) :]
                affixes_found.append(f"PRE:{prefix}")
                break

        # Check suffixes
        for suffix in GEEZ_SUFFIXES:
            if stem.endswith(suffix) and len(stem) > len(suffix) + 2:
                stem = stem[: -len(suffix)]
                affixes_found.append(f"SUF:{suffix}")
                break

        return stem, affixes_found

    @classmethod
    def strip_affixes_oromo(cls, word: str) -> tuple[str, list[str]]:
        """
        Strips common Afaan Oromoo agglutinative suffixes,
        returning the base stem and detected morphemes.
        """
        stem = word.lower()
        affixes_found: list[str] = []

        for suffix in OROMO_SUFFIXES:
            if stem.endswith(suffix) and len(stem) > len(suffix) + 3:
                stem = stem[: -len(suffix)]
                affixes_found.append(f"OR_SUF:{suffix}")
                break

        return stem, affixes_found

    def generate_character_ngrams(self, word: str) -> list[str]:
        """
        Generates character-level subword n-grams (from min_ngram to max_ngram),
        wrapping words with boundary markers for prefix/suffix position awareness.
        """
        padded = f"<{word}>"
        ngrams: list[str] = []
        n_len = len(padded)

        for n in range(self.min_ngram, min(self.max_ngram + 1, n_len + 1)):
            for i in range(n_len - n + 1):
                ngrams.append(padded[i : i + n])

        return ngrams

    def tokenize_document(self, text: str) -> dict[str, Any]:
        """
        Performs full hierarchical decomposition of a sentence or multi-sentence utterance:
        - Consecutive sentences
        - Word tokens
        - Morphological stems
        - Detected affixes
        - Individual letters
        - Consonant root radicals
        - Subword character n-grams
        """
        if not text:
            return {
                "sentences": [],
                "words": [],
                "stems": [],
                "affixes": [],
                "letters": [],
                "roots": [],
                "char_ngrams": [],
            }

        # 1. Split consecutive sentences (handles English '.', Amharic '።' and '?')
        sentences = [s.strip() for s in re.split(r"[.!?።፤\n]+", text) if s.strip()]
        if not sentences:
            sentences = [text.strip()]

        # 2. Extract words
        clean_text = re.sub(r"[^\w\s-]", " ", text)
        words = [w.strip() for w in clean_text.split() if w.strip()]

        stems: list[str] = []
        affixes: list[str] = []
        letters: list[str] = []
        roots: list[str] = []
        char_ngrams: list[str] = []

        for w in words:
            # Letters
            letters.extend(self.extract_letters(w))

            # Ge'ez Fidel roots
            roots.extend(self.extract_fidel_roots(w))

            # Affix stripping
            g_stem, g_affixes = self.strip_affixes_geez(w)
            o_stem, o_affixes = self.strip_affixes_oromo(w)

            stems.append(g_stem)
            stems.append(o_stem)
            affixes.extend(g_affixes)
            affixes.extend(o_affixes)

            # Character n-grams
            char_ngrams.extend(self.generate_character_ngrams(w))

        return {
            "sentences": sentences,
            "words": words,
            "stems": list(dict.fromkeys(stems)),
            "affixes": affixes,
            "letters": letters,
            "roots": roots,
            "char_ngrams": char_ngrams,
        }

    def to_enriched_feature_string(self, text: str) -> str:
        """
        Converts a text into a dense, multi-resolution feature document containing:
        original words, morphological stems, character n-grams, and consonant roots.
        This allows vectorizers to simultaneously learn word semantics, letter patterns,
        and morphological roots.
        """
        tokens = self.tokenize_document(text)
        parts = [
            text,  # Full original consecutive sentence
            " ".join(tokens["stems"]),  # Morphological stems
            " ".join(tokens["affixes"]),  # Affixes
            " ".join(tokens["roots"]),  # Ge'ez consonant roots
        ]
        return " ".join([p for p in parts if p.strip()])


# Global singleton instance
morph_tokenizer = MorphologicalTokenizer()
