"""Multilingual taxonomy mapping and synonym resolution for Ethiopian marketplace."""

# Mapping from multilingual terms across EN, AM, OM, TI to standard catalog root slugs
CATEGORY_SYNONYM_MAP: dict[str, list[str]] = {
    "phones-telecommunications": [
        "phone", "smartphone", "cellphone", "mobile",
        "ስልክ", "ስማርትፎን", "ሞባይል", "ተንቀሳቃሽ ስልክ",
        "bilbila", "moobaayilii",
        "ተንቀሳቓሲ ስልኪ", "ሞባይል"
    ],
    "computers-laptops-it-hardware": [
        "computer", "pc", "laptop", "desktop", "macbook",
        "ኮምፒውተር", "ላፕቶፕ", "ፒሲ", "ዴስክቶፕ",
        "kompiyitara", "laaptooppii",
        "ኮምፒተር", "ላፕቶፕ"
    ],
    "vehicles-transport-automotive": [
        "car", "vehicle", "automobile", "truck", "motorcycle", "van",
        "መኪና", "ተሽከርካሪ", "አውቶሞቢል", "ባጃጅ", "ሞተርሳይክል",
        "konkolaataa", "baajaajii", "mootorsayikilii",
        "መኪና", "ባጃጅ", "ተሽከርካሪ"
    ],
    "real-estate-property-housing": [
        "house", "apartment", "villa", "condominium", "studio",
        "ቤት", "አፓርታማ", "ቪላ", "ኮንዶሚኒየም",
        "mana", "appaartamaa", "viillaa",
        "ገዛ", "ኣፓርትመንት", "ቪላ"
    ],
    "land-plots-real-estate": [
        "land", "plot", "field", "site",
        "መሬት", "ቦታ", "የግንባታ ቦታ",
        "lafa", "bakka",
        "መሬት", "ቦታ"
    ],
    "agriculture-farming-livestock": [
        "tractor", "seeds", "fertilizer", "crop", "cattle", "sheep",
        "ትራክተር", "ዘር", "ማዳበሪያ", "ከብት", "በግ",
        "tiraaktara", "sanyii", "xaa'oo", "horii",
        "ትራክተር", "ዘርኢ", "መዳበሪያ", "ከብቲ"
    ],
}

# Transaction type mapping (Sale vs. Rent)
TRANSACTION_TYPE_MAP: dict[str, str] = {
    # Rent
    "rent": "rent",
    "rental": "rent",
    "lease": "rent",
    "ኪራይ": "rent",
    "ማከራየት": "rent",
    "kiraa": "rent",
    "kireessuu": "rent",
    "ክራይ": "rent",
    "ምክራይ": "rent",
    # Sale
    "buy": "sale",
    "purchase": "sale",
    "sale": "sale",
    "sell": "sale",
    "መግዛት": "sale",
    "ሽያጭ": "sale",
    "መሸጥ": "sale",
    "bituu": "sale",
    "gurgurtaa": "sale",
    "መሸጢ": "sale",
    "ምግዛእ": "sale",
}


def resolve_category_slug(term: str) -> str | None:
    """Resolve a multilingual keyword to a standard category slug."""
    clean = term.strip().lower()
    for slug, synonyms in CATEGORY_SYNONYM_MAP.items():
        if any(syn in clean for syn in synonyms):
            return slug
    return None


def resolve_transaction_type(term: str) -> str | None:
    """Resolve whether a term indicates rent or sale."""
    clean = term.strip().lower()
    return TRANSACTION_TYPE_MAP.get(clean)
