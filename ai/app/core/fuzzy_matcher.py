"""
HUNTIQ Fast Fuzzy Entity & Brand / Location Matcher.
Uses RapidFuzz to resolve typos, phonetic variants, and noisy inputs
in Ethiopian multilingual marketplace queries.
"""

from rapidfuzz import fuzz, process

# Canonical canonical location map (standardized -> list of variants)
CANONICAL_LOCATIONS = {
    "Addis Ababa": [
        "addis ababa",
        "adis ababa",
        "addis abeba",
        "finfinnee",
        "finfinne",
        "ኣዲስ ኣበባ",
        "አዲስ አበባ",
    ],
    "Bole": ["bole", "boolee", "bolee", "ቦሌ", "ቦለ"],
    "CMC": ["cmc", "c.m.c", "ሲኤምሲ", "ሲ ኤም ሲ"],
    "Sarbet": ["sarbet", "sar bet", "ሳርቤት", "ሳር ቤት"],
    "Megenagna": ["megenagna", "megenanya", "መገናኛ"],
    "Piassa": ["piassa", "piyassa", "piasa", "ፒያሳ"],
    "Mexico": ["mexico", "meksiko", "ሜክሲኮ"],
    "Kazanchis": ["kazanchis", "kasanchis", "ካዛንቺስ"],
    "Hawassa": ["hawassa", "hawasa", "hawaasaa", "ሀዋሳ", "ሓዋሳ"],
    "Adama": ["adama", "adaamaa", "nazret", "nazareth", "አዳማ", "ናዝሬት"],
    "Bahir Dar": ["bahir dar", "bahirdar", "ባህር ዳር", "ባሕር ዳር"],
    "Mekelle": ["mekelle", "mekele", "መቀሌ", "መቐለ"],
    "Dire Dawa": ["dire dawa", "diredawa", "ድሬዳዋ", "ድሬ ዳዋ"],
    "Jimma": ["jimma", "jima", "jimmaa", "ጅማ"],
    "Gondar": ["gondar", "gonder", "ጎንደር"],
    "Bishoftu": ["bishoftu", "debre zeit", "ቢሾፍቱ", "ደብረዘይት"],
    "Jemo": ["jemo", "ጀሞ"],
    "Gerji": ["gerji", "ገርጂ"],
    "Summit": ["summit", "ሳሚት", "ሰሚት"],
    "Ayat": ["ayat", "አያት", "ኣያት"],
}

# Canonical brands
CANONICAL_BRANDS = {
    "Samsung": ["samsung", "samson", "samsong", "ሳምሰንግ"],
    "Apple": ["apple", "iphone", "ipon", "aifon", "አፕል", "አይፎን"],
    "Toyota": ["toyota", "tayota", "toyta", "ቶዮታ"],
    "Suzuki": ["suzuki", "suzki", "ሱዙኪ"],
    "Dell": ["dell", "del", "ዴል"],
    "HP": ["hp", "hewlett packard", "ኤችፒ", "ኤች ፒ"],
    "Lenovo": ["lenovo", "ሊኖቮ", "ሌኖቮ"],
    "MacBook": ["macbook", "mac book", "ማክቡክ"],
    "Tecno": ["tecno", "tekno", "ቴክኖ"],
    "Infinix": ["infinix", "ኢንፊኒክስ"],
    "Xiaomi": ["xiaomi", "redmi", "ሬድሚ", "ሻኦሚ"],
    "LG": ["lg", "ኤልጂ"],
    "Sony": ["sony", "ሶኒ"],
    "Hisense": ["hisense", "ሃይሴንስ", "ሀይሴንስ"],
}


class FuzzyEntityMatcher:
    def __init__(self):
        # Flatten locations for fast fuzzy search
        self.location_lookup = {}
        for canonical, variants in CANONICAL_LOCATIONS.items():
            for v in variants:
                self.location_lookup[v.lower()] = canonical

        # Flatten brands for fast fuzzy search
        self.brand_lookup = {}
        for canonical, variants in CANONICAL_BRANDS.items():
            for v in variants:
                self.brand_lookup[v.lower()] = canonical

        self.loc_choices = list(self.location_lookup.keys())
        self.brand_choices = list(self.brand_lookup.keys())

    def match_location(
        self, word: str, score_cutoff: float = 85.0
    ) -> tuple[str, float] | None:
        """
        Fuzzy matches a candidate location word to canonical Ethiopian location.
        Returns (canonical_name, score) or None.
        """
        w = word.strip().lower()
        if not w or len(w) < 3:
            return None

        # Direct match
        if w in self.location_lookup:
            return self.location_lookup[w], 100.0

        # Fuzzy ratio match
        res = process.extractOne(
            w, self.loc_choices, scorer=fuzz.ratio, score_cutoff=score_cutoff
        )
        if res:
            matched_variant, score, _ = res
            return self.location_lookup[matched_variant], float(score)

        return None

    def match_brand(
        self, word: str, score_cutoff: float = 85.0
    ) -> tuple[str, float] | None:
        """
        Fuzzy matches a candidate word to canonical electronics/vehicle brand.
        Returns (canonical_name, score) or None.
        """
        w = word.strip().lower()
        if not w or len(w) < 2:
            return None

        # Direct match
        if w in self.brand_lookup:
            return self.brand_lookup[w], 100.0

        # Fuzzy ratio match
        res = process.extractOne(
            w, self.brand_choices, scorer=fuzz.ratio, score_cutoff=score_cutoff
        )
        if res:
            matched_variant, score, _ = res
            return self.brand_lookup[matched_variant], float(score)

        return None


fuzzy_matcher = FuzzyEntityMatcher()
