"""Authoritative Multilingual Intent Classifier for HUNTIQ AI Platform."""

import re
from enum import Enum

from pydantic import BaseModel, Field

from ..languages.normalizer import (
    detect_language,
    normalize_geez_text,
    normalize_oromo_text,
)
from .clarification import check_clarification_needed
from .entity_extractor import ExtractedSearch, extract_searches
from .policy import check_platform_faq, get_safety_warning, is_out_of_domain
from .trained_model_service import trained_models


class IntentType(str, Enum):
    SOCIAL_GREETING = "social_greeting"
    SOCIAL_FAREWELL = "social_farewell"
    SOCIAL_GRATITUDE = "social_gratitude"
    IDENTITY_INQUIRY = "identity_inquiry"
    CAPABILITY_INQUIRY = "capability_inquiry"
    PLATFORM_FAQ = "platform_faq"
    OUT_OF_DOMAIN = "out_of_domain"
    MARKETPLACE_SEARCH = "marketplace_search"
    CLARIFICATION_NEEDED = "clarification_needed"
    GENERAL_KNOWLEDGE = "general_knowledge"
    SAFETY_VIOLATION = "safety_violation"
    SCAM_ATTEMPT = "scam_attempt"


# Social greetings across 4 languages
GREETING_PATTERNS = [
    r"(?i)\b(hello|hi|hey|good morning|good afternoon|good evening)\b",
    r"(?i)(ሰላም|እንደምን|ጤና ይስጥልኝ|እንደምን አደርክ|እንደምን ዋልክ)",
    r"(?i)\b(akkam|akkam jirtu|nagaa|ashamaa)\b",
    r"(?i)(ሰላም|ከመይ ኣለኻ|ከመይ ሓዲርኩም|ከመይ ዊዒልኩም)",
]

# Social farewells
FAREWELL_PATTERNS = [
    r"(?i)\b(bye|goodbye|see you|farewell)\b",
    r"(?i)(ቻው|ደህና ሁን|ደህና ሁኚ|ደህና ሁኑ)",
    r"(?i)\b(nagaatti|booda wal argina)\b",
    r"(?i)(ደሓን ኩን|ደሓን ኩኑ|ቻው)",
]

# Social gratitude
GRATITUDE_PATTERNS = [
    r"(?i)\b(thanks|thank you|much appreciated)\b",
    r"(?i)(አመሰግናለሁ|እናመሰግናለን|እግዚአብሔር ይስጥልኝ)",
    r"(?i)\b(galatoomaa|galatoomi|baay'ee galatoomi)\b",
    r"(?i)(የቐንየለይ|ኣመስግነኩም)",
]

# Identity questions ("who are you", "what is your name")
IDENTITY_PATTERNS = [
    r"(?i)\b(who are you|what is your name|your name|what are you)\b",
    r"(?i)(ማን ነህ|ስምህ ማን ነው|ማን ነሽ|ስምሽ ማን ነው|አንተ ማን ነህ)",
    r"(?i)\b(eenyu ati|maqaan kee eenyu|maal ati)\b",
    r"(?i)(መን ኢኻ|መን ኢኺ|ስምካ መን እዩ)",
]

# Capability inquiries ("what can you do", "help me")
CAPABILITY_PATTERNS = [
    r"(?i)\b(what can you do|how can you help|what do you do)\b",
    r"(?i)(ምን ትሰራለህ|በምን ልትረዳኝ ትችላለህ|ምን ማድረግ ትችላለህ)",
    r"(?i)\b(maal hojjachuu dandeessa|akkamitti na gargaarta)\b",
    r"(?i)(እንታይ ክትገብር ትኽእል|ብምንታይ ክትሕግዘኒ ትኽእል)",
]

# Localized standard responses
GREETING_RESPONSES = {
    "en": "Hello! Welcome to HUNTIQ.",
    "am": "ሰላም! እንኳን ወደ HUNTIQ በደህና መጡ።",
    "om": "Akkam! Baga nagaan gara HUNTIQ dhuftan.",
    "ti": "ሰላም! ናብ HUNTIQ ብደሓን መጻእኹም።",
}

IDENTITY_RESPONSES = {
    "en": "I am HUNTIQ's AI Assistant, built specifically for Ethiopia's universal marketplace.",
    "am": "እኔ የHUNTIQ የገበያ ረዳት AI ነኝ፤ የኢትዮጵያን ገበያ ለማገልገል የተዘጋጀሁ።",
    "om": "Ani Gargaaraa AI HUNTIQ ti; gabaa Itoophiyaa tajaajiluuf kan qophaa'e.",
    "ti": "ኣነ ናይ HUNTIQ ናይ AI ረዳኢ እየ፤ ንዕዳጋ ኢትዮጵያ ንምግልጋል ዝተዳለዉ።",
}

FAREWELL_RESPONSES = {
    "en": "Goodbye! Let me know whenever you need anything else on HUNTIQ.",
    "am": "ደህና ሁኑ! በHUNTIQ ላይ ሌላ ማንኛውም ነገር ሲፈልጉ እዚህ አለሁ።",
    "om": "Nagaatti! Yeroo kamiyyuu HUNTIQ irratti waan barbaaddan yoo jiraate as jiras.",
    "ti": "ደሓን ኩኑ! ኣብ HUNTIQ ዝኾነ ነገር ክትደልዩ ከለኹም ኣብዚ ኣለኹ።",
}

GRATITUDE_RESPONSES = {
    "en": "You are very welcome! What else can I find for you on HUNTIQ today?",
    "am": "ምንም አይደል! ዛሬ በHUNTIQ ላይ ሌላ ምን ላሳይዎት?",
    "om": "Homaa miti! Har'a HUNTIQ irratti maal biraa isiniif barbaadu?",
    "ti": "ገንዘብኩም! ሎሚ ኣብ HUNTIQ ካልእ እንታይ ክረኽበልኩም?",
}


class IntentClassificationPlan(BaseModel):
    """Structured plan returned by the authoritative intent classifier."""

    detected_language: str
    intents: list[IntentType] = Field(default_factory=list)
    social_greeting: str | None = None
    social_farewell: str | None = None
    identity_response: str | None = None
    platform_faq_response: str | None = None
    safety_warning: str | None = None
    safety_class: str | None = None
    model_confidence: float | None = None
    requires_clarification: bool = False
    clarification_question: str | None = None
    unresolved_slot: str | None = None
    searches: list[ExtractedSearch] = Field(default_factory=list)
    requires_database_search: bool = False


def classify_message(message: str, locale_hint: str = "en") -> IntentClassificationPlan:
    """
    Authoritative intent classification and query understanding pipeline.
    Decomposes multi-intent messages and prepares execution plan.
    """
    clean_msg = message.strip()
    detected_lang = detect_language(clean_msg, locale_hint)
    if detected_lang == "or":
        detected_lang = "om"

    # Script normalization
    normalized_text = clean_msg
    if detected_lang in ("am", "ti"):
        normalized_text = normalize_geez_text(clean_msg)
    elif detected_lang == "om":
        normalized_text = normalize_oromo_text(clean_msg)

    intents: list[IntentType] = []
    social_greeting: str | None = None
    social_farewell: str | None = None
    identity_resp: str | None = None
    faq_resp: str | None = None

    # 0. Safety & Scam Guard Check via ML Model
    safety_class, safety_conf, _ = trained_models.predict_safety(normalized_text)
    if safety_class and safety_class != "safe" and safety_conf >= 0.70:
        warning_msg = get_safety_warning(safety_class, detected_lang)
        intent_type = (
            IntentType.SCAM_ATTEMPT
            if safety_class in ("scam_advance", "phishing_harvesting")
            else IntentType.SAFETY_VIOLATION
        )
        return IntentClassificationPlan(
            detected_language=detected_lang,
            intents=[intent_type],
            safety_warning=warning_msg,
            safety_class=safety_class,
            model_confidence=safety_conf,
            requires_database_search=False,
        )

    # Check if query matches in-platform conversational patterns (identity, greeting, farewell, capability)
    is_social_or_identity = any(
        re.search(pat, normalized_text)
        for pat_list in (
            GREETING_PATTERNS,
            FAREWELL_PATTERNS,
            GRATITUDE_PATTERNS,
            IDENTITY_PATTERNS,
            CAPABILITY_PATTERNS,
        )
        for pat in pat_list
    )

    # 1. Out-of-Domain Guard via Rule & ML Model (skipped for legitimate assistant identity/greeting)
    if is_out_of_domain(normalized_text) or (
        not is_social_or_identity
        and trained_models.predict_domain(normalized_text)[0] == "general_out_of_domain"
        and trained_models.predict_domain(normalized_text)[1] >= 0.88
    ):
        intents.append(IntentType.OUT_OF_DOMAIN)
        return IntentClassificationPlan(
            detected_language=detected_lang,
            intents=intents,
            requires_database_search=False,
        )

    # 2. Platform FAQ
    is_faq, faq_ans = check_platform_faq(normalized_text, detected_lang)
    if is_faq:
        intents.append(IntentType.PLATFORM_FAQ)
        faq_resp = faq_ans

    # 3. Social Greetings
    for pat in GREETING_PATTERNS:
        if re.search(pat, normalized_text):
            intents.append(IntentType.SOCIAL_GREETING)
            social_greeting = GREETING_RESPONSES.get(
                detected_lang, GREETING_RESPONSES["en"]
            )
            break

    # 4. Social Farewells
    for pat in FAREWELL_PATTERNS:
        if re.search(pat, normalized_text):
            intents.append(IntentType.SOCIAL_FAREWELL)
            social_farewell = FAREWELL_RESPONSES.get(
                detected_lang, FAREWELL_RESPONSES["en"]
            )
            break

    # 5. Social Gratitude
    for pat in GRATITUDE_PATTERNS:
        if re.search(pat, normalized_text):
            intents.append(IntentType.SOCIAL_GRATITUDE)
            break

    # 6. Identity Inquiries
    for pat in IDENTITY_PATTERNS:
        if re.search(pat, normalized_text):
            intents.append(IntentType.IDENTITY_INQUIRY)
            identity_resp = IDENTITY_RESPONSES.get(
                detected_lang, IDENTITY_RESPONSES["en"]
            )
            break

    # 7. Capability Inquiries
    for pat in CAPABILITY_PATTERNS:
        if re.search(pat, normalized_text):
            intents.append(IntentType.CAPABILITY_INQUIRY)
            break

    # 8. Marketplace Searches & Entity Extraction
    # Extract searches from normalized text
    searches = extract_searches(normalized_text)

    has_social_intent = any(
        i in intents
        for i in (
            IntentType.SOCIAL_GREETING,
            IntentType.IDENTITY_INQUIRY,
            IntentType.CAPABILITY_INQUIRY,
            IntentType.SOCIAL_FAREWELL,
            IntentType.SOCIAL_GRATITUDE,
            IntentType.PLATFORM_FAQ,
        )
    )

    # Only treat as marketplace search if:
    # 1. An explicit categorized search was found, OR
    # 2. No social/conversational intent is present and a search term was extracted
    is_genuine_market_search = searches and (
        any(s.category_name != "Marketplace" for s in searches) or not has_social_intent
    )

    if is_genuine_market_search:
        # Filter out generic 'Marketplace' catch-all if specific categories exist
        specific = [s for s in searches if s.category_name != "Marketplace"]
        if specific:
            searches = specific
        intents.append(IntentType.MARKETPLACE_SEARCH)
        requires_db = True
    else:
        searches = []
        requires_db = False

    # 9. Clarification Check
    req_clarification = False
    clarification_q = None
    unresolved_slot = None
    if searches:
        req_clarification, clarification_q, unresolved_slot = (
            check_clarification_needed(searches, normalized_text, detected_lang)
        )
        if req_clarification:
            intents.append(IntentType.CLARIFICATION_NEEDED)

    # 10. ML Intent Classifier fallback if no explicit rule matched
    if not intents:
        pred_intent, intent_conf, _ = trained_models.predict_intent(normalized_text)
        if pred_intent and intent_conf >= 0.55:
            if pred_intent == "social_greeting":
                intents.append(IntentType.SOCIAL_GREETING)
                social_greeting = GREETING_RESPONSES.get(
                    detected_lang, GREETING_RESPONSES["en"]
                )
            elif pred_intent == "social_farewell":
                intents.append(IntentType.SOCIAL_FAREWELL)
                social_farewell = FAREWELL_RESPONSES.get(
                    detected_lang, FAREWELL_RESPONSES["en"]
                )
            elif pred_intent == "social_gratitude":
                intents.append(IntentType.SOCIAL_GRATITUDE)
            elif pred_intent == "identity_inquiry":
                intents.append(IntentType.IDENTITY_INQUIRY)
                identity_resp = IDENTITY_RESPONSES.get(
                    detected_lang, IDENTITY_RESPONSES["en"]
                )
            elif pred_intent == "marketplace_search":
                intents.append(IntentType.MARKETPLACE_SEARCH)
                requires_db = True

    if not intents:
        intents.append(IntentType.GENERAL_KNOWLEDGE)

    return IntentClassificationPlan(
        detected_language=detected_lang,
        intents=intents,
        social_greeting=social_greeting,
        social_farewell=social_farewell,
        identity_response=identity_resp,
        platform_faq_response=faq_resp,
        requires_clarification=req_clarification,
        clarification_question=clarification_q,
        unresolved_slot=unresolved_slot,
        searches=searches,
        requires_database_search=requires_db,
    )
