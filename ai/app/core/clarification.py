"""Ambiguity detection and localized clarification question generator."""

from .entity_extractor import ExtractedSearch

CLARIFICATION_HOUSING = {
    "en": "Are you looking to rent or buy a property on HUNTIQ?",
    "am": "ቤት ወይም አፓርታማ ለመከራየት ወይስ ለመግዛት ይፈልጋሉ?",
    "om": "Mana kireeffachuu moo bituu barbaaddu?",
    "ti": "ገዛ ንምክራይ ድዩ ወይስ ንምግዛእ ትደልዩ?",
}

CLARIFICATION_LOCATION = {
    "en": "Which city or subcity in Ethiopia would you like to search in (e.g., Bole, CMC, Hawassa)?",
    "am": "በየትኛው ከተማ ወይም ክፍለ ከተማ መፈለግ ይፈልጋሉ (ለምሳሌ፡ ቦሌ፣ ሲኤምሲ፣ ሀዋሳ)?",
    "om": "Magaalaa ykn kutaa magaalaa kam keessatti barbaaduu feetu (fkn: Boolee, CMC, Hawaasaa)?",
    "ti": "ኣየናይ ከተማ ወይ ክፍለ ከተማ ክትደልዩ ትደልዩ (ንኣብነት፡ ቦሌ፣ ሲኤምሲ፣ ሓዋሳ)?",
}


def check_clarification_needed(
    searches: list[ExtractedSearch], raw_query: str, locale: str
) -> tuple[bool, str | None, str | None]:
    """
    Determine if the query lacks essential disambiguating slots.
    Returns: (requires_clarification, clarification_question, unresolved_slot)
    """
    loc = (locale or "en").lower()
    if loc == "or":
        loc = "om"

    lower = raw_query.lower()

    for s in searches:
        # Housing ambiguity: if category is Real Estate and user did not specify rent vs. sale
        if (
            s.category_name in ("Real Estate & Housing", "Property Rentals")
            and not s.listing_type
        ):
            # Check if raw text had rent or buy keywords
            has_rent = any(k in lower for k in ["rent", "ኪራይ", "kiraa", "ክራይ"])
            has_buy = any(
                k in lower for k in ["buy", "sale", "ግዢ", "ሽያጭ", "bituu", "ምግዛእ"]
            )
            if not has_rent and not has_buy:
                question = CLARIFICATION_HOUSING.get(loc, CLARIFICATION_HOUSING["en"])
                return True, question, "listing_type"

    return False, None, None
