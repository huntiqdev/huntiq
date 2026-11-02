"""Entity and search slot extraction for HUNTIQ multilingual queries."""

import re

from pydantic import BaseModel

from .fuzzy_matcher import fuzzy_matcher
from .trained_model_service import trained_models

# Monetary parsing regexes across EN, AM, OM, TI
RX_PRICE_UNDER = re.compile(
    r"(?i)\b(?:under|below|less than|max|maximum of|budget of|up to)\s*([\d,]+(?:\.\d+)?[kKmM]?)\s*(?:etb|birr|ብር|qarshii|ቅርሺ)?\b"
)
RX_PRICE_ABOVE = re.compile(
    r"(?i)\b(?:above|over|more than|min|minimum of|starting from|at least)\s*([\d,]+(?:\.\d+)?[kKmM]?)\s*(?:etb|birr|ብር|qarshii|ቅርሺ)?\b"
)
RX_PRICE_HAVE = re.compile(
    r"(?i)\b(?:i have|with|budget of|my budget is|around|have)\s*([\d,]+(?:\.\d+)?[kKmM]?)\s*(?:etb|birr|ብር|qarshii|ቅርሺ)?\b"
)
RX_PRICE_STANDALONE = re.compile(
    r"(?i)\b([\d,]+(?:\.\d+)?[kKmM]?)\s*(?:etb|birr|ብር|qarshii|ቅርሺ)\b"
)
RX_AM_MAX = re.compile(
    r"(?i)(?:ከ\s*)?([\d,]+(?:\.\d+)?[kKmM]?)\s*(?:ብር|etb)?\s*(?:በታች|ያነሰ)"
)
RX_AM_BUDGET = re.compile(r"(?i)([\d,]+(?:\.\d+)?[kKmM]?)\s*(?:ብር|etb)?\s*(?:አለኝ|ያለኝ)")
RX_OM_MAX = re.compile(
    r"(?i)(?:gatii\s*)?([\d,]+(?:\.\d+)?[kKmM]?)\s*(?:qarshii|etb)?\s*(?:gadi)"
)
RX_OM_BUDGET = re.compile(r"(?i)(?:qarshii\s*)?([\d,]+(?:\.\d+)?[kKmM]?)\s*(?:qaba)")
RX_TI_MAX = re.compile(
    r"(?i)(?:ትሕቲ|ትሕቲ\s*)?([\d,]+(?:\.\d+)?[kKmM]?)\s*(?:ቅርሺ|etb)?\s*(?:ትሕቲ)"
)
RX_TI_BUDGET = re.compile(r"(?i)(?:ቅርሺ\s*)?([\d,]+(?:\.\d+)?[kKmM]?)\s*(?:ኣለኒ)")

# Known Ethiopian locations (Latin & Ge'ez)
KNOWN_LOCATIONS: dict[str, str] = {
    "addis ababa": "Addis Ababa",
    "አዲስ አበባ": "Addis Ababa",
    "finfinnee": "Addis Ababa",
    "ኣዲስ ኣበባ": "Addis Ababa",
    "bole": "Bole",
    "ቦሌ": "Bole",
    "boolee": "Bole",
    "cmc": "CMC",
    "ሲኤምሲ": "CMC",
    "sarbet": "Sarbet",
    "ሳርቤት": "Sarbet",
    "megenagna": "Megenagna",
    "መገናኛ": "Megenagna",
    "piassa": "Piassa",
    "ፒያሳ": "Piassa",
    "mexico": "Mexico",
    "ሜክሲኮ": "Mexico",
    "kazanchis": "Kazanchis",
    "ካዛንቺስ": "Kazanchis",
    "hawassa": "Hawassa",
    "ሀዋሳ": "Hawassa",
    "hawaasaa": "Hawassa",
    "ሓዋሳ": "Hawassa",
    "adama": "Adama",
    "አዳማ": "Adama",
    "adaamaa": "Adama",
    "ናዝሬት": "Adama",
    "bahir dar": "Bahir Dar",
    "ባህር ዳር": "Bahir Dar",
    "ባሕር ዳር": "Bahir Dar",
    "mekelle": "Mekelle",
    "መቀሌ": "Mekelle",
    "መቐለ": "Mekelle",
    "dire dawa": "Dire Dawa",
    "ድሬዳዋ": "Dire Dawa",
    "ድሬ ዳዋ": "Dire Dawa",
    "jimma": "Jimma",
    "ጅማ": "Jimma",
    "jimmaa": "Jimma",
    "gondar": "Gondar",
    "ጎንደር": "Gondar",
    "bishoftu": "Bishoftu",
    "ቢሾፍቱ": "Bishoftu",
    "ደብረዘይት": "Bishoftu",
}

# Known brands
KNOWN_BRANDS: list[str] = [
    "samsung",
    "apple",
    "iphone",
    "toyota",
    "vitz",
    "corolla",
    "yaris",
    "suzuki",
    "dell",
    "hp",
    "lenovo",
    "macbook",
    "bosch",
    "makita",
    "devalt",
    "dewalt",
    "tecno",
    "infinix",
    "redmi",
    "xiaomi",
    "lg",
    "sony",
    "hisense",
]

# Canonical category definitions
CATEGORY_SPECS = [
    {
        "name": "Mobile Phones",
        "slug": "phones-telecommunications",
        "keywords": [
            "phone",
            "phones",
            "mobile",
            "mobiles",
            "smartphone",
            "smartphones",
            "iphone",
            "samsung",
            "redmi",
            "tecno",
            "infinix",
            "ስልክ",
            "ስልኮች",
            "ሞባይል",
            "ሞባይሎች",
            "አይፎን",
            "ሳምሰንግ",
            "bilbila",
            "bilbiloota",
            "moobaayilii",
            "ayifoonii",
            "ስልኪ",
            "ስልክታት",
            "ሞባይል",
            "ኣይፎን",
        ],
        "default_type": "sale",
    },
    {
        "name": "Computers & Laptops",
        "slug": "computers-laptops-it-hardware",
        "keywords": [
            "pc",
            "computer",
            "computers",
            "laptop",
            "laptops",
            "macbook",
            "desktop",
            "workstation",
            "dell",
            "hp",
            "lenovo",
            "ኮምፒውተር",
            "ኮምፒውተሮች",
            "ላፕቶፕ",
            "ላፕቶፖች",
            "ማክቡክ",
            "kompiitara",
            "kompiitaraa",
            "laaptooppii",
            "ኮምፒተር",
            "ላፕቶፕ",
        ],
        "default_type": "sale",
    },
    {
        "name": "Vehicles & Cars",
        "slug": "vehicles-transport-automotive",
        "keywords": [
            "car",
            "cars",
            "vehicle",
            "vehicles",
            "automobile",
            "vitz",
            "corolla",
            "yaris",
            "motorcycle",
            "bajaj",
            "van",
            "መኪና",
            "መኪኖች",
            "ተሽከርካሪ",
            "አውቶሞቢል",
            "ቪትዝ",
            "ኮሮላ",
            "ባጃጅ",
            "ሞተርሳይክል",
            "konkolaataa",
            "konkolaattota",
            "baajaajii",
            "መኪና",
            "ተሽከርካሪ",
            "ባጃጅ",
        ],
        "default_type": "sale",
    },
    {
        "name": "Property Rentals",
        "slug": "real-estate-property-housing",
        "keywords": [
            "apartments for rent",
            "apartment for rent",
            "houses for rent",
            "house for rent",
            "rent apartment",
            "rent apartments",
            "rent house",
            "rent houses",
            "condo for rent",
            "rental",
            "rentals",
            "የሚከራይ ቤት",
            "የሚከራዩ ቤቶች",
            "የሚከራይ አፓርታማ",
            "የቤት ኪራይ",
            "ኮንዶሚኒየም ኪራይ",
            "ኪራይ ቤት",
            "mana kireeffamu",
            "manneen kireeffaman",
            "mana kiraa",
            "appaartamaa kiraa",
            "ዝካረ ገዛ",
            "ዝካረ ኣፓርታማ",
            "ክራይ ገዛ",
        ],
        "default_type": "rent",
    },
    {
        "name": "Real Estate & Housing",
        "slug": "real-estate-property-housing",
        "keywords": [
            "house",
            "houses",
            "apartment",
            "apartments",
            "condo",
            "condominium",
            "villa",
            "studio",
            "ቤት",
            "ቤቶች",
            "አፓርታማ",
            "ኮንዶሚኒየም",
            "ቪላ",
            "mana",
            "manneen",
            "appaartamaa",
            "viillaa",
            "ገዛ",
            "ኣፓርትመንት",
            "ቪላ",
        ],
        "default_type": None,
    },
    {
        "name": "Land & Plots",
        "slug": "land-plots-real-estate",
        "keywords": [
            "land",
            "plot",
            "plots",
            "field",
            "site",
            "መሬት",
            "ቦታ",
            "የግንባታ ቦታ",
            "lafa",
            "bakka",
            "መሬት",
            "ቦታ",
        ],
        "default_type": "sale",
    },
    {
        "name": "Heavy Machinery & Tools",
        "slug": "industrial-machinery-equipment",
        "keywords": [
            "drill",
            "generator",
            "tractor",
            "pump",
            "welder",
            "machine",
            "tools",
            "bosch",
            "makita",
            "ድሪል",
            "ጄኔሬተር",
            "ትራክተር",
            "ፓምፕ",
            "ማሽን",
            "መሳሪያዎች",
            "ቦሽ",
            "tiraaktara",
            "jeneraetara",
            "meeshaa",
            "ድሪል",
            "ጄነሬተር",
            "ትራክተር",
            "ማሽን",
        ],
        "default_type": "sale",
    },
]


class ExtractedSearch(BaseModel):
    """Normalized structured search parameters for a single marketplace category."""

    category_name: str
    category_slug: str | None = None
    search_term: str
    brand: str | None = None
    listing_type: str | None = None  # "sale" or "rent"
    min_price: float | None = None
    max_price: float | None = None
    location: str | None = None
    condition: str | None = None  # "new" or "used"


def parse_numeric_price(val_str: str) -> float | None:
    """Parse numeric price strings with Birr/ETB and k/m suffixes (e.g. '20k' -> 20000)."""
    if not val_str:
        return None
    clean = val_str.strip().replace(",", "")
    multiplier = 1.0
    if clean.endswith(("k", "K")):
        multiplier = 1000.0
        clean = clean[:-1]
    elif clean.endswith(("m", "M")):
        multiplier = 1000000.0
        clean = clean[:-1]

    try:
        val = float(clean)
        return val * multiplier
    except ValueError:
        return None


def extract_price_bounds(query: str) -> tuple[float | None, float | None]:
    """Extract minimum and maximum price bounds from user query."""
    min_p: float | None = None
    max_p: float | None = None

    # Max patterns
    for rx in [RX_PRICE_UNDER, RX_AM_MAX, RX_OM_MAX, RX_TI_MAX]:
        m = rx.search(query)
        if m:
            val = parse_numeric_price(m.group(1))
            if val is not None and val > 0:
                max_p = val
                break

    # Min patterns
    m_above = RX_PRICE_ABOVE.search(query)
    if m_above:
        val = parse_numeric_price(m_above.group(1))
        if val is not None and val > 0:
            min_p = val

    # Budget / have patterns
    if max_p is None:
        for rx in [RX_PRICE_HAVE, RX_AM_BUDGET, RX_OM_BUDGET, RX_TI_BUDGET]:
            m = rx.search(query)
            if m:
                val = parse_numeric_price(m.group(1))
                if val is not None and val > 0:
                    max_p = val
                    break

    # Fallback standalone price
    if max_p is None and min_p is None:
        m_stand = RX_PRICE_STANDALONE.search(query)
        if m_stand:
            val = parse_numeric_price(m_stand.group(1))
            if val is not None and val > 0:
                max_p = val

    return min_p, max_p


def extract_location(query: str) -> str | None:
    """Extract known Ethiopian cities and subcities with exact and fuzzy matching."""
    lower = query.lower()
    for loc_key, canonical in KNOWN_LOCATIONS.items():
        # Match as word boundary or exact token
        pattern = r"(?i)(?:\b|(?<=\s))" + re.escape(loc_key) + r"(?:\b|(?=\s)|$)"
        if re.search(pattern, lower) or loc_key in lower:
            return canonical

    # Fuzzy fallback via RapidFuzz
    tokens = re.findall(r"\w+", lower)
    for tok in tokens:
        res = fuzzy_matcher.match_location(tok)
        if res and res[1] >= 85.0:
            return res[0]
    return None


def extract_condition(query: str) -> str | None:
    """Extract condition ('new' or 'used')."""
    lower = query.lower()
    if any(k in lower for k in ["new", "brand new", "አዲስ", "haaraa", "ሓድሽ"]):
        return "new"
    if any(k in lower for k in ["used", "second hand", "ያገለገለ", "moofaa", "ዝተጠቐመሉ"]):
        return "used"
    return None


def extract_listing_type(query: str) -> str | None:
    """Extract transaction mode ('rent' or 'sale')."""
    lower = query.lower()
    if any(
        k in lower
        for k in [
            "rent",
            "rental",
            "lease",
            "to rent",
            "for rent",
            "ኪራይ",
            "ማከራየት",
            "የሚከራይ",
            "ኪራይ ቤት",
            "kiraa",
            "kireeffamu",
            "kireessuu",
            "ክራይ",
            "ምክራይ",
            "ዝካረ",
        ]
    ):
        return "rent"
    if any(
        k in lower
        for k in [
            "buy",
            "purchase",
            "sale",
            "for sale",
            "መግዛት",
            "ሽያጭ",
            "የሚሸጥ",
            "መሸጥ",
            "bituu",
            "gurgurtaa",
            "gurguruu",
            "ምግዛእ",
            "መሸጢ",
        ]
    ):
        return "sale"
    return None


def extract_brand(query: str) -> str | None:
    """Extract recognized equipment or electronics brand with fuzzy support."""
    lower = query.lower()
    for b in KNOWN_BRANDS:
        if re.search(r"\b" + re.escape(b) + r"\b", lower):
            return b.capitalize()

    # Fuzzy fallback via RapidFuzz
    tokens = re.findall(r"\w+", lower)
    for tok in tokens:
        res = fuzzy_matcher.match_brand(tok)
        if res and res[1] >= 85.0:
            return res[0]
    return None


# Fallback mapping for ML predicted category classes
ML_CATEGORY_MAPPING = {
    "vehicles": {
        "name": "Vehicles & Cars",
        "slug": "vehicles-transport-automotive",
        "default_type": "sale",
    },
    "real_estate": {
        "name": "Real Estate & Housing",
        "slug": "real-estate-property-housing",
        "default_type": None,
    },
    "products": {
        "name": "Electronics & Products",
        "slug": "phones-telecommunications",
        "default_type": "sale",
    },
    "agriculture": {
        "name": "Agriculture & Food",
        "slug": "agriculture-food-commodities",
        "default_type": "sale",
    },
    "jobs": {
        "name": "Jobs & Employment",
        "slug": "jobs-employment",
        "default_type": "sale",
    },
    "services": {
        "name": "Professional Services",
        "slug": "professional-home-services",
        "default_type": "sale",
    },
    "hotels": {
        "name": "Travel & Hotels",
        "slug": "travel-hotels-tours",
        "default_type": "sale",
    },
    "tickets": {
        "name": "Event Tickets",
        "slug": "events-tickets",
        "default_type": "sale",
    },
}


def _extract_single_clause(
    clause: str, global_loc: str | None = None
) -> list[ExtractedSearch]:
    lower = clause.lower()
    min_p, max_p = extract_price_bounds(clause)
    loc = extract_location(clause) or global_loc
    cond = extract_condition(clause)
    tx_type = extract_listing_type(clause)
    brand = extract_brand(clause)

    matched_specs = []
    for spec in CATEGORY_SPECS:
        for kw in spec["keywords"]:
            if (
                re.search(
                    r"(?i)(?:\b|(?<=\s))" + re.escape(kw) + r"(?:\b|(?=\s)|$)", lower
                )
                or kw in lower
            ):
                # Avoid adding both Property Rentals and Real Estate if it's explicitly rent
                if spec["name"] == "Real Estate & Housing" and (
                    tx_type == "rent" or spec["default_type"] == "rent"
                ):
                    continue
                matched_specs.append((spec, kw))
                break

    # If rule matching didn't trigger, consult the trained category model
    if not matched_specs:
        predicted_cat, conf, _ = trained_models.predict_category(clause)
        if predicted_cat and conf >= 0.40 and predicted_cat in ML_CATEGORY_MAPPING:
            spec_info = ML_CATEGORY_MAPPING[predicted_cat]
            if predicted_cat == "real_estate" and tx_type == "rent":
                spec_info = {
                    "name": "Property Rentals",
                    "slug": "real-estate-property-housing",
                    "default_type": "rent",
                }
            clean_term = clause
            for filler in [
                "i need",
                "looking for",
                "find me",
                "show me",
                "where can i find",
                "እፈልጋለሁ",
                "ፈልግልኝ",
                "barbaada",
                "ይደሊ ኣለኹ",
                "find",
                "for sale",
                "for rent",
            ]:
                clean_term = re.sub(
                    re.escape(filler), "", clean_term, flags=re.IGNORECASE
                )
            clean_term = clean_term.strip() or predicted_cat
            matched_specs.append((spec_info, clean_term))

    if not matched_specs:
        search_term = clause
        for filler in [
            "i need",
            "looking for",
            "find me",
            "show me",
            "where can i find",
            "እፈልጋለሁ",
            "ፈልግልኝ",
            "barbaada",
            "ይደሊ ኣለኹ",
            "find",
            "for sale",
            "for rent",
        ]:
            search_term = re.sub(
                re.escape(filler), "", search_term, flags=re.IGNORECASE
            )
        search_term = search_term.strip()
        if search_term:
            return [
                ExtractedSearch(
                    category_name="Marketplace",
                    category_slug=None,
                    search_term=search_term,
                    brand=brand,
                    listing_type=tx_type,
                    min_price=min_p,
                    max_price=max_p,
                    location=loc,
                    condition=cond,
                )
            ]
        return []

    results = []
    seen = set()
    for spec, kw in matched_specs:
        if spec["name"] in seen:
            continue
        seen.add(spec["name"])
        eff_type = tx_type or spec["default_type"]
        search_kw = brand.lower() if brand else kw
        results.append(
            ExtractedSearch(
                category_name=spec["name"],
                category_slug=spec["slug"],
                search_term=search_kw,
                brand=brand,
                listing_type=eff_type,
                min_price=min_p,
                max_price=max_p,
                location=loc,
                condition=cond,
            )
        )
    return results


def extract_searches(query: str) -> list[ExtractedSearch]:
    """
    Decompose user message into one or more structured category searches.
    Handles multi-intent queries like 'Find apartments for rent and land for sale'.
    """
    clean_q = query.strip()
    if not clean_q:
        return []

    global_loc = extract_location(clean_q)

    # Split multi-intent queries by conjunctions if present
    conjunction_pattern = r"\s+(?:and|&|\+|እና|fi|ከምኡ'ውን)\s+"
    parts = re.split(conjunction_pattern, clean_q, flags=re.IGNORECASE)

    if len(parts) > 1:
        aggregated: list[ExtractedSearch] = []
        seen_names = set()
        for p in parts:
            part_searches = _extract_single_clause(p.strip(), global_loc)
            for ps in part_searches:
                if ps.category_name not in seen_names:
                    seen_names.add(ps.category_name)
                    aggregated.append(ps)
        if aggregated:
            return aggregated

    return _extract_single_clause(clean_q, global_loc)
