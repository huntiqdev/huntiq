"""Response composition and zero-fabrication grounding engine for HUNTIQ marketplace."""

from ..api.schemas import CategoryGroupResult

EMPTY_RESULTS_MESSAGES = {
    "en": "I searched HUNTIQ's verified database but found no approved listings matching your request right now. New items are published daily, so feel free to check back shortly or explore related categories.",
    "am": "የHUNTIQ የቀጥታ ዳታቤዝ ተፈትሾ በአሁኑ ሰዓት ከፍለጋዎ ጋር የሚዛመድ ንቁ ዝርዝር አልተገኘም። አዳዲስ እቃዎች በየቀኑ ስለሚለጠፉ እባክዎ ከጥቂት ጊዜ በኋላ እንደገና ይመልከቱ።",
    "om": "Kuusaa gabaa HUNTIQ sakatta'uun oomisha barbaaddan yeroo ammaa hin arganne. Oomishaaleen haaraan guyyaa guyyaan waan maxxanfamaniif xiqqoo turanii irra deebi'aa ilaalaa.",
    "ti": "ኣብ HUNTIQ ተፈቲሹ ኣብዚ እዋን ምስ ድላይኩም ዝሰማማዕ ንጡፍ ዝርዝር ኣይተረኽበን። ሓደሽቲ ፍርያት መዓልታዊ ስለዝልጠፉ ድሒርኩም ደጊምኩም ፈትኑ።",
}

CATEGORY_FOUND_HEADERS = {
    "en": "Here are {count} verified listings for {term} on HUNTIQ:",
    "am": "በHUNTIQ ላይ ለ{term} የተገኙ {count} የተረጋገጡ ዝርዝሮች እነሆ:",
    "om": "HUNTIQ irratti {term} oomishaalee mirkanaa'an {count} kunooti:",
    "ti": "ኣብ HUNTIQ ን{term} ዝተረኽቡ {count} ዝተረጋገጹ ዝርዝራት እንሆ:",
}

MULTI_CATEGORY_HEADERS = {
    "en": "Here are the live marketplace results across your requested categories:",
    "am": "በጠየቋቸው ምድቦች የተገኙ የቀጥታ የገበያ ውጤቶች እነሆ:",
    "om": "Ramaddiiwwan isin gaafattan keessatti bu'aawwan gabaa qophaa'an kunooti:",
    "ti": "ኣብ ዝሓተትኩምዎም ምድባት ዝተረኽቡ ቀጥታ ናይ ዕዳጋ ውጽኢታት እንሆ:",
}


def compose_grounded_response(
    query: str,
    greeting_prefix: str | None,
    groups: list[CategoryGroupResult],
    locale: str = "en",
) -> str:
    """
    Produce an honest, grounded natural language summary of verified tool outputs.
    Ensures zero hallucinated prices, URLs, or inventory.
    """
    loc = (locale or "en").lower()
    if loc == "or":
        loc = "om"

    parts = []
    if greeting_prefix:
        parts.append(greeting_prefix)

    total_listings = sum(len(g.listings) for g in groups)

    if total_listings == 0:
        empty_msg = EMPTY_RESULTS_MESSAGES.get(loc, EMPTY_RESULTS_MESSAGES["en"])
        parts.append(empty_msg)
        return " ".join(parts).strip()

    if len(groups) == 1:
        grp = groups[0]
        term = grp.search_term or grp.category_name
        hdr_template = CATEGORY_FOUND_HEADERS.get(loc, CATEGORY_FOUND_HEADERS["en"])
        header = hdr_template.format(count=len(grp.listings), term=term)
        parts.append(header)
    else:
        multi_hdr = MULTI_CATEGORY_HEADERS.get(loc, MULTI_CATEGORY_HEADERS["en"])
        parts.append(multi_hdr)

    return " ".join(parts).strip()
