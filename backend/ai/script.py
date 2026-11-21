#!/usr/bin/env python3
"""
HUNTIQ AI — Master Engine (v2.0)
Single-file marketplace intelligence, contextual domain validation,
four-language understanding, fast social routing, and Google Gemini API integration.
Uses ONLY Python Standard Library (no third-party dependencies).
"""

import argparse
import http.server
import json
import logging
import os
import re
import socketserver
import sys
import threading
import time
import urllib.error
import urllib.request
from typing import Any, Dict, List, Optional, Tuple

logging.basicConfig(
    level=logging.INFO,
    format="[HUNTIQ-AI] %(asctime)s [%(levelname)s] %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
logger = logging.getLogger("huntiq_ai")

# ══════════════════════════════════════════════════════════════════════════════
# 1. CENTRAL POLICY LOADER (v2.0)
# ══════════════════════════════════════════════════════════════════════════════

DEFAULT_POLICY = {
    "policy_version": "2.0",
    "domain_identity": {
        "name": "HUNTIQ AI",
        "creator": "HUNTIQ",
        "purpose": "Marketplace intelligence, discovery, and commercial advisory",
        "scope": "HUNTIQ marketplace and closely related commercial activities",
    },
    "supported_languages": ["en", "am", "ti", "or"],
    "allowed_domains": [
        "marketplace_discovery",
        "product_search",
        "product_comparison",
        "product_purchasing_advice",
        "technology_purchasing_advice",
        "pricing_and_valuation",
        "property_and_rentals",
        "vehicles_and_transport",
        "employment_and_jobs",
        "freelance_and_gigs",
        "professional_services",
        "business_opportunities_and_tenders",
        "events_and_tickets",
        "parking_and_facilities",
        "broker_business_model",
        "marketplace_economics",
        "pricing_strategy",
        "negotiation_tactics",
        "supply_and_demand",
        "customer_trust_and_safety",
        "listing_optimization",
        "seller_performance_and_growth",
        "buyer_advisory_and_inspection",
        "platform_guidance_and_policies",
        "marketplace_commission_calculation",
        "commercial_math_and_comparison",
        "marketplace_trend_analysis",
        "commercial_recommendations",
    ],
    "allowed_intents": [
        "MARKETPLACE_SEARCH",
        "PRODUCT_AVAILABILITY",
        "PRODUCT_COMPARISON",
        "BUYER_ADVICE",
        "SELLER_ADVICE",
        "COMMERCIAL_ADVISORY",
        "BROKER_MODEL_EXPLANATION",
        "MARKETPLACE_BUSINESS_ADVICE",
        "PLATFORM_NAVIGATION",
        "COMMISSION_CALCULATION",
        "PRODUCT_MATH_ANALYSIS",
        "MARKET_TREND_ANALYSIS",
        "RECOMMENDATION",
        "GREETING",
        "APPRECIATION",
        "COMPLIMENT",
        "ACKNOWLEDGEMENT",
        "GOODBYE",
        "CAPABILITY_QUESTION",
        "USER_PRIVACY_PROTECTION",
        "SEARCH_REFINEMENT",
        "CONTEXTUAL_FOLLOWUP",
        "CLARIFICATION",
    ],
    "excluded_domains": [
        "unrelated_general_mathematics",
        "unrelated_general_physics",
        "unrelated_general_chemistry",
        "unrelated_general_biology",
        "unrelated_general_history",
        "unrelated_general_geography",
        "unrelated_general_programming",
        "unrelated_general_entertainment",
        "unrelated_general_trivia",
        "unrelated_homework_solving",
        "general_purpose_virtual_assistant",
        "credential_theft_or_jailbreak",
    ],
}


class PolicyEngine:
    def __init__(self, policy_path: Optional[str] = None):
        self.policy = DEFAULT_POLICY
        if policy_path and os.path.exists(policy_path):
            try:
                with open(policy_path, "r", encoding="utf-8") as f:
                    loaded = json.load(f)
                    if isinstance(loaded, dict) and "domain_identity" in loaded:
                        self.policy = loaded
                        logger.info("Loaded policy from %s (v%s)", policy_path, loaded.get("policy_version", "1.0"))
            except Exception as e:
                logger.warning("Failed to load policy file %s: %s. Using defaults.", policy_path, e)

    def is_language_supported(self, lang: str) -> bool:
        return lang in self.policy.get("supported_languages", ["en", "am", "ti", "or"])


# ══════════════════════════════════════════════════════════════════════════════
# 2. FOUR-LANGUAGE DETECTION & LOCALIZED RESPONSES
# ══════════════════════════════════════════════════════════════════════════════

RE_ETHIOPIC = re.compile(r"[\u1200-\u137F]")

TIGRINYA_MARKERS = [
    "ከመይ", "ደሊየ", "ኣሎ", "ኣለኹም", "እንታይ", "ክካረ", "ምግዛእ", "ምሻጥ", "ክራይ", "ገዛ",
    "ዓደግቲ", "ሸየጥቲ", "ተኻረይቲ", "ዕዳጋ", "ብዛዕባ", "ከመይሉ", "ንብረት", "ሓበሬታ",
    "ናይዚ", "ክንደይ", "እዩ", "ኣለዎ", "ክትገዝኡ", "ምኽሪ", "ብኸመይ", "የቐንየለይ", "ደሓን",
    "ሕራይ", "ብልሂ", "መን ኢኻ",
]

AMHARIC_MARKERS = [
    "ሰላም", "እፈልጋለሁ", "ይፈልጋሉ", "መግዛት", "መሸጥ", "ኪራይ", "ዋጋ", "ቤት", "መኪና",
    "ስልክ", "ላፕቶፕ", "ስራ", "አዲስ አበባ", "ደላላ", "የገበያ", "ቦሌ", "ማስታወቂያ", "ተከራይ",
    "ገዢ", "ሻጭ", "ባለሙያ", "አሳይ", "ፈልግ", "ምን", "እንዴት", "ዋጋው", "አመሰግናለሁ", "ጎበዝ",
    "እሺ", "ደህና", "ማን ነህ", "ብልህ",
]

OROMO_MARKERS = [
    "akkam", "barbaada", "barbaaddu", "biyya", "maaloo", "gatii", "mana", "konkolaataa",
    "hojii", "dilaala", "gabaa", "bittaa", "gurgurtaa", "kireeffachuu", "kireessuu", "ogeessa",
    "galatoomi", "galatoomaa", "nagaatti", "tole", "qaruutee", "ati eenyu",
]


def detect_language(text: str, requested_lang: str = "") -> str:
    if requested_lang in ("en", "am", "ti", "or"):
        return requested_lang

    clean = text.strip()
    if not clean:
        return "en"

    # Ge'ez Unicode script check (Amharic vs Tigrinya)
    if RE_ETHIOPIC.search(clean):
        lower = clean.lower()
        ti_score = sum(1 for m in TIGRINYA_MARKERS if m in lower)
        am_score = sum(1 for m in AMHARIC_MARKERS if m in lower)
        if ti_score > am_score:
            return "ti"
        return "am"

    lower = clean.lower()
    for marker in OROMO_MARKERS:
        if re.search(r"\b" + re.escape(marker) + r"\b", lower):
            return "or"

    return "en"


def get_localized_rejection(lang: str, reason: str = "") -> str:
    if lang == "am":
        return (
            "እኔ HUNTIQ AI ነኝ፤ የተዘጋጀሁት በHUNTIQ ነው። "
            "በHUNTIQ የገበያ ቦታ ላይ እቃዎችን፣ ቤቶችን፣ መኪኖችን፣ የሥራ እና የሙያ አገልግሎቶችን ለማግኘት "
            "ወይም የገበያ እና የንግድ ምክር ለመጠየቅ እባክዎ ከHUNTIQ ጋር የተያያዘ ጥያቄ ያቅርቡ።"
        )
    elif lang == "ti":
        return (
            "ኣነ HUNTIQ AI እየ፤ ብHUNTIQ ዝተዳለኹ። "
            "ኣብ HUNTIQ ዕዳጋ ንብረታት፣ ኣባይቲ፣ መካይን፣ ስራሕን ናይ ሞያ ኣገልግሎታትን ንምርካብ "
            "ወይ ናይ ዕዳጋን ንግድን ምኽሪ ንምሕታት በጃኹም ምስ HUNTIQ ዝተኣሳሰረ ሕቶ ሕተቱ።"
        )
    elif lang == "or":
        return (
            "Ani HUNTIQ AI dha; kan qophaa'e HUNTIQ dhaan. "
            "Gabaa HUNTIQ irratti meeshaalee, manneen, konkolaattota, hojii fi tajaajiloota ogeessaa barbaaduuf "
            "ykn gorsa daldalaa argachuuf maaloo gaaffii HUNTIQ waliin walqabatu gaafadhaa."
        )
    return (
        "I am HUNTIQ AI, made by HUNTIQ. I assist exclusively with HUNTIQ marketplace listings, "
        "products, properties, rentals, vehicles, jobs, services, pricing strategies, and platform guidance. "
        "Please ask something related to HUNTIQ or marketplace commerce."
    )


def get_localized_clarification(lang: str) -> str:
    if lang == "am":
        return "እባክዎ ጥያቄዎን በዝርዝር ያብራሩልኝ — በHUNTIQ ላይ ምን አይነት እቃ ወይም የገበያ መረጃ መፈለግ ይፈልጋሉ?"
    elif lang == "ti":
        return "በጃኹም ሕቶኹም ብዝርዝር ግለጹለይ — ኣብ HUNTIQ እንታይ ዓይነት ንብረት ወይ ናይ ዕዳጋ ሓበሬታ ክትደልዩ ትደልዩ?"
    elif lang == "or":
        return "Maaloo gaaffii keessan bal'inaan ibsaa — HUNTIQ irratti meeshaa akkamii ykn odeeffannoo gabaa maalii barbaaddu?"
    return "Could you please clarify your question? What marketplace listing, product, property, service, or business topic can I assist you with on HUNTIQ?"


def get_fast_greeting(lang: str) -> str:
    if lang == "am":
        return "ሰላም! ዛሬ በHUNTIQ ምን መፈለግ ይፈልጋሉ ወይም ምን እንድረዳዎ ይፈልጋሉ?"
    elif lang == "ti":
        return "ከመይ! ሎሚ ኣብ HUNTIQ እንታይ ክትደልዩ ትደልዩ ወይ እንታይ ክሕግዘኩም ትደልዩ?"
    elif lang == "or":
        return "Akkam! Har'a HUNTIQ irratti maal barbaaduu barbaaddu ykn maaliin isin gargaaruu?"
    return "Hello! How can I help you discover marketplace listings or assist your commerce on HUNTIQ today?"


def get_fast_compliment(lang: str) -> str:
    if lang == "am":
        return "አመሰግናለሁ! በHUNTIQ የገበያ ቦታ ላይ ምርጥ እቃዎችን፣ ቤቶችን፣ መኪኖችን እና የንግድ ምክሮችን እንዲያገኙ ለመርዳት ሁሌም ዝግጁ ነኝ።"
    elif lang == "ti":
        return "የቐንየለይ! ኣብ HUNTIQ ዕዳጋ ዝበለጹ ንብረታት፣ ኣባይቲ፣ መካይን ከምኡ'ውን ናይ ንግዲ ምኽሪ ንምርካብ ንምሕጋዝኩም ኩሉ ግዜ ድሉው እየ።"
    elif lang == "or":
        return "Galatoomaa! Gabaa HUNTIQ irratti meeshaalee, manneen, konkolaattota fi gorsa daldalaa akka argattan isin gargaaruuf yeroo hundumaa qophiidha."
    return "Thank you for the kind words! I am here to help you navigate HUNTIQ's marketplace, find great deals, and provide commercial advice whenever you need it."


def get_fast_appreciation(lang: str) -> str:
    if lang == "am":
        return "ምንም አይደል! ስላመሰገኑኝ ደስ ብሎኛል። በHUNTIQ ላይ ሌላ ማንኛውም የፍለጋ ወይም የገበያ ምክር ከፈለጉ በደስታ እረዳዎታለሁ።"
    elif lang == "ti":
        return "ገንዘብኩም ይኹን! ብጣዕሚ የቐንየለይ። ኣብ HUNTIQ ካልእ ናይ ፍለሻ ወይ ናይ ዕዳጋ ምኽሪ እንተደሊኹም ብሓጎስ ክሕግዘኩም እየ።"
    elif lang == "or":
        return "Homaa miti! Galatoomaa. HUNTIQ irratti barbaacha ykn gorsa daldalaa dabalataa yoo barbaaddan gammachuudhaan isin gargaara."
    return "You are very welcome! I appreciate your kind words. Let me know if you need any other assistance with HUNTIQ listings or marketplace advisory."


def get_fast_acknowledgement(lang: str) -> str:
    if lang == "am":
        return "እሺ! በHUNTIQ ላይ እቃዎችን፣ ቤቶችን ወይም አገልግሎቶችን መፈለግ ወይም የንግድ ምክር መጠየቅ ሲፈልጉ በማንኛውም ጊዜ ይጠይቁኝ።"
    elif lang == "ti":
        return "ሕራይ! ኣብ HUNTIQ ንብረታት ክትደልዩ ወይ ናይ ዕዳጋ ምኽሪ ክትሓቱ ምስ ትደልዩ ኣብ ዝኾነ ግዜ ሕተቱኒ።"
    elif lang == "or":
        return "Tole! HUNTIQ irratti meeshaalee barbaaduu ykn gorsa daldalaa argachuu yoo barbaaddan yeroo kamiyyuu na gaafadhaa."
    return "Understood! Feel free to ask whenever you want to search listings or explore marketplace insights on HUNTIQ."


def get_fast_capability(lang: str) -> str:
    if lang == "am":
        return (
            "እኔ በHUNTIQ የተዘጋጀሁት HUNTIQ AI ነኝ። በኢትዮጵያ ውስጥ የHUNTIQ የደላላ የገበያ ቦታ ረዳትዎ ነኝ። "
            "እቃዎችን፣ ቤቶችን፣ መኪኖችን፣ የስራ እና የሙያ አገልግሎቶችን እንዲፈልጉ፣ ዋጋና የደላላ ኮሚሽን እንዲያሰሉ፣ "
            "እንዲሁም ለገዢዎች፣ ሻጮችና አከራዮች ተግባራዊ የገበያ ምክር እንድሰጥ ተዘጋጅቻለሁ።"
        )
    elif lang == "ti":
        return (
            "ኣነ ብHUNTIQ ዝተዳለኹ HUNTIQ AI እየ። ኣብ ኢትዮጵያ ናይ HUNTIQ ደላላ ዕዳጋ ሓጋዚኹም እየ። "
            "ንብረታት፣ ኣባይቲ፣ መካይን፣ ስራሕን ናይ ሞያ ኣገልግሎታትን ክትደልዩ፣ ዋጋን ኮሚሽን ደላላን ክትሕስቡ፣ "
            "ከምኡ'ውን ንዓደግቲ፣ ሸየጥትን መኻረይትን ተግባራዊ ምኽሪ ክህብ ተዳልየ ኣለኹ።"
        )
    elif lang == "or":
        return (
            "Ani HUNTIQ AI dha, kan qophaa'e HUNTIQ dhaan. Itoophiyaa keessatti gabaa dilaala HUNTIQ irratti gargaaraa keessan. "
            "Meeshaalee, manneen, konkolaattota, hojii fi tajaajiloota ogeessaa barbaaduu, gatii fi komishinii dilaalaa shallaguu, "
            "akkasumas bittattoota, gurgurtoota fi kireessitootaaf gorsa daldalaa kennuu nan danda'a."
        )
    return (
        "I am HUNTIQ AI, created by HUNTIQ. I am your specialized broker marketplace assistant in Ethiopia. "
        "I can help you search products, properties, vehicles, jobs, and professional services, compare options, "
        "evaluate pricing and commissions, and provide practical advice for buyers, sellers, and landlords."
    )


def get_fast_goodbye(lang: str) -> str:
    if lang == "am":
        return "ደህና ሁኑ! HUNTIQ ን ስለተጠቀሙ አመሰግናለሁ። መልካም ቀን ይሁንልዎ፤ የገበያ እርዳታ ሲያስፈልግዎ በማንኛውም ጊዜ ይመለሱ!"
    elif lang == "ti":
        return "ደሓን ኩኑ! ንHUNTIQ ስለዝተጠቐምኩም የቐንየለይ። ብሩኽ መዓልቲ ይኹነልኩም፤ ናይ ዕዳጋ ሓገዝ ምስ ዘድልየኩም ኣብ ዝኾነ ግዜ ተመለሱ!"
    elif lang == "or":
        return "Nagaatti! HUNTIQ fayyadamuu keessaniif galatoomaa. Guyyaa gaarii qabaadhaa; gargaarsi gabaa yoo isin barbaachise yeroo kamiyyuu deebi'aa!"
    return "Goodbye! Thank you for using HUNTIQ. Have a wonderful day, and return whenever you need marketplace assistance!"


def get_fast_privacy(lang: str) -> str:
    if lang == "am":
        return (
            "HUNTIQ የተጠቃሚዎችን የግል መረጃ እና ደህንነት በጥብቅ ይጠብቃል። "
            "የተመዘገቡ ተጠቃሚዎችን ዝርዝር፣ የግል አካውንቶችን ወይም ማንነቶችን አልገልጽም። "
            "ነገር ግን አጠቃላይ የገበያ መረጃዎችን፣ በይፋ የወጡ ማስታወቂያዎችን ወይም የፕላትፎርሙን አሰራር ማሳየት እችላለሁ።"
        )
    elif lang == "ti":
        return (
            "HUNTIQ ናይ ተጠቀምቲ ብሕታዊ ሓበሬታን ድሕንነትን ብጽኑዕ ይሕሉ። "
            "ዝተመዝገቡ ተጠቀምቲ ዝርዝር ወይ ብሕታዊ መንነቶም ኣይገልጽን። "
            "ይኹን እምበር ሓፈሻዊ ናይ ዕዳጋ ሓበሬታ፣ ብግህዶ ዝወጹ መወዓውዒታት ወይ ናይ ፕላትፎርም ኣሰራርሓ ከርእየኩም እኽእል እየ።"
        )
    elif lang == "or":
        return (
            "HUNTIQ icciitii fi nageenya fayyadamtootaa cimsee eega. "
            "Tarree fayyadamtoota galmaa'anii, herrega dhuunfaa ykn eenyummaa isaanii hin ibsu. "
            "Haa ta'u malee, odeeffannoo gabaa waliigalaa, beeksisa ifatti ba'e ykn hojimaata waltajjichaa isinitti agarsiisuu nan danda'a."
        )
    return (
        "HUNTIQ strictly protects user privacy and account security. I do not disclose personal registration details, "
        "user accounts, or private identities. I can, however, provide aggregate marketplace statistics, platform guidance, "
        "or help you find publicly listed items and services."
    )


# ══════════════════════════════════════════════════════════════════════════════
# 3. REGEX DEFINITIONS & CONTEXTUAL PATTERNS
# ══════════════════════════════════════════════════════════════════════════════

# Adversarial prompt injections, secret theft, jailbreaks
RE_JAILBREAK = re.compile(
    r"(?i)(ignore\s+(all\s+)?previous\s+instructions|disregard\s+(all\s+)?prior|"
    r"you\s+are\s+now\s+(an\s+)?unrestricted|pretend\s+to\s+be|dan\s+mode|jailbreak|"
    r"reveal\s+(the\s+)?system\s+prompt|reveal\s+.*secret|developer\s+mode|"
    r"show\s+me\s+(backend\s+)?credentials|database\s+passwords?|\.env\s+contents?|bypass\s+security)"
)

# User privacy protection (Phase 15): Refuse inquiries regarding private user accounts / registrations
RE_USER_PRIVACY = re.compile(
    r"(?i)(who\s+(?:registered|signed\s*up|joined|created\s+an?\s+account)\s+(?:yesterday|today|recently|this\s+week)|"
    r"list\s+(?:all\s+)?(?:registered\s+)?(?:users|accounts|members|buyers|sellers)|"
    r"show\s+(?:me\s+)?(?:all\s+)?(?:user\s+|account\s+|member\s+)*(?:emails|phone\s+numbers|passwords|credentials|identities|personal\s+data)|"
    r"who\s+are\s+the\s+(?:registered\s+)?users|"
    r"(?:who|ማን|መን|eenyutu)\s+.*(?:registered|ተመዘገበ|ተመዝጊቡ|galmaa'e)|"
    r"(?:ተጠቃሚዎችን|የተጠቃሚ)\s+(?:አሳይ|ዝርዝር|መረጃ))"
)

# Social interaction patterns (Phase 6)
RE_GREETING = re.compile(
    r"(?i)^(?:hello|hi|hey|good\s+(?:morning|afternoon|evening)|selam|ሰላም|ከመይ|akkam)\b[!?.]*$"
)

RE_COMPLIMENT = re.compile(
    r"(?i)(?:^|\b)("
    r"you\s+are\s+(?:clever|smart|helpful|awesome|great|amazing|the\s+best|genius|good)|"
    r"you're\s+(?:clever|smart|helpful|awesome|great|amazing|the\s+best|genius|good)|"
    r"smart\s+bot|clever\s+ai|good\s+job|nice\s+work|impressive|"
    r"ጎበዝ|ብልህ|ጥሩ\s+ስራ|ብልሂ|qaruutee|baay'ee\s+gaarii"
    r")(?:\s+\w+)?(?:[!?.\s]*)$"
)

RE_APPRECIATION_THANKS = re.compile(
    r"(?i)\b("
    r"thank\s+you(\s+so\s+much)?|thanks(\s+a\s+lot)?|i\s+appreciate\s+you|much\s+appreciated|thx|ty|"
    r"አመሰግናለሁ|እናመሰግናለን|የቐንየለይ|galatoomi|galatoomaa"
    r")\b[!?.]*$"
)

RE_ACKNOWLEDGEMENT = re.compile(
    r"(?i)^(ok|okay|alright|got\s+it|understood|cool|sure|fine|nice|i\s+see|እሺ|ሕራይ|tole)\b[!?.]*$"
)

RE_CAPABILITY_OR_IDENTITY = re.compile(
    r"(?i)\b("
    r"who\s+are\s+you|what\s+can\s+you\s+do|who\s+created\s+you|who\s+made\s+you|"
    r"what\s+is\s+your\s+name|what\s+is\s+huntiq\s*ai|tell\s+me\s+about\s+yourself|"
    r"አንተ\s+ማን\s+ነህ|ማን\s+ነሽ|ማን\s+ነህ|ምን\s+መስራት\s+ትችላለህ|ምን\s+ትሰራለህ|"
    r"መን\s+ኢኻ|እንታይ\s+ክትገብር\s+ትኽእል|ati\s+eenyu|maal\s+gochuu\s+dandeessa"
    r")\b[!?.]*$"
)

RE_GOODBYE = re.compile(
    r"(?i)\b("
    r"goodbye|bye|bye\s+bye|see\s+you(\s+later)?|ciao|take\s+care|"
    r"ደህና\s+ሁን|ደህና\s+ሁኚ|ደህና\s+ሁኑ|ቻው|ደሓን\s+ኩን|ደሓን\s+ኩኑ|nagaatti"
    r")\b[!?.]*$"
)

# Unrelated pure sciences, academic physics, chemistry, biology
RE_PURE_SCIENCE = re.compile(
    r"(?i)(e\s*=\s*mc\s*[\^²2]|quantum\s+mechanics|quantum\s+physics|thermodynamics|"
    r"photosynthesis|mitochondria|chemical\s+reaction|avogadro|periodic\s+table|"
    r"speed\s+of\s+light|black\s+hole|derivative\s+of|integral\s+of|"
    r"quadratic\s+equation|organic\s+chemistry|dna\s+replication)"
)

# Pure ungrounded math (e.g. "Solve 25 x 37", "25 * 37", "sqrt(144)")
RE_PURE_MATH = re.compile(
    r"(?i)(^(?:solve|calculate|what\s+is|compute)?\s*\(?\s*\d+(?:\.\d+)?\s*[\+\-\*\/×÷\^]\s*\d+(?:\.\d+)?\s*\)?$|"
    r"^(?:solve|calculate|what\s+is|compute)\s+\d+\s*(?:times|multiplied\s+by|plus|minus|divided\s+by)\s*\d+$|"
    r"(?:sqrt\s*\(\d+\)|solve\s+for\s+x))"
)

# Unrelated general programming tutorials (hello world, leetcode, binary tree)
RE_GENERIC_CODING = re.compile(
    r"(?i)(write\s+a\s+(python|javascript|golang|c\+\+|java)\s+hello\s*world|"
    r"create\s+a\s+binary\s+tree|implement\s+(quick\s*sort|merge\s*sort|bubble\s*sort)|"
    r"leetcode|hackerrank|reverse\s+a\s+linked\s*list)"
)

# Unrelated history, astronomy, celebrity trivia, entertainment
RE_GENERIC_TRIVIA = re.compile(
    r"(?i)(who\s+won\s+.*world\s+cup|who\s+won\s+.*historical\s+war|capital\s+of\s+[a-z]+|"
    r"who\s+is\s+the\s+president\s+of\s+[a-z]+|tell\s+me\s+a\s+joke|write\s+an\s+essay\s+about|"
    r"horoscope|celebrity\s+gossip|distance\s+to\s+mars|constellation)"
)

# Practical marketplace math exception (discounts, commissions, power consumption, price per m²)
RE_MARKETPLACE_MATH_EXCEPTION = re.compile(
    r"(?i)(discount|percentage|%|commission|broker\s+fee|profit|margin|depreciation|"
    r"fuel\s+consumption|engine\s+capacity|cc\b|kilowatt|watt|power\s+consumption|"
    r"electricity\s+consumption|price\s+per\s+square|square\s+meter|m2|sqm|cost\s+per)"
)

# Commercial & marketplace keyword signals
MARKETPLACE_SIGNALS = [
    # English Commercial & General Items
    "buy", "sell", "rent", "price", "cost", "house", "apartment", "car", "vehicle",
    "truck", "phone", "laptop", "pc", "computer", "hardware", "software", "ram", "ssd",
    "processor", "core", "cpu", "screen", "monitor", "display", "specs", "specification",
    "tools", "equipment", "job", "hire", "service", "electrician", "plumber", "mechanic",
    "tender", "gig", "contract", "listing", "product", "huntiq", "account", "register",
    "seller", "buyer", "commission", "negotiate", "discount", "demand", "supply",
    "wholesale", "retail", "store", "deposit", "lease", "tenant", "landlord", "customer",
    "trust", "sales", "revenue", "inspection", "pricing", "commercial", "property",
    "toyota", "hyundai", "lenovo", "dell", "hp", "macbook", "iphone", "samsung", "thinkpad",
    "addis ababa", "bole", "hawassa", "adama", "dire dawa", "bahir dar", "mekelle",
    # Amharic
    "መግዛት", "መሸጥ", "ኪራይ", "ዋጋ", "ቤት", "መኪና", "ስልክ", "ላፕቶፕ", "ኮምፒውተር", "ስራ",
    "አዲስ አበባ", "ደላላ", "የገበያ", "ቦሌ", "ማስታወቂያ", "ተከራይ", "ገዢ", "ሻጭ", "ባለሙያ",
    "ኮሚሽን", "ንግድ", "ምርመራ", "መግለጫ", "ኮምፒተር",
    # Tigrinya
    "ምግዛእ", "ምሻጥ", "ክራይ", "ዋጋ", "ገዛ", "መኪና", "ኮምፒተር", "ስራሕ", "ዓደግቲ", "ሸየጥቲ",
    "ተኻረይቲ", "ዕዳጋ", "ኮሚሽን", "ደላላ",
    # Oromo
    "bittaa", "gurgurtaa", "kireeffachuu", "kireessuu", "gatii", "mana", "konkolaataa",
    "hojii", "dilaala", "komishinii", "kompiitara",
]

RE_PRICE_UNDER = re.compile(r"(?i)(?:under|below|less\s+than|maximum|max)\s+([0-9]+(?:,[0-9]+)*(?:\.[0-9]+)?)\s*(?:etb|birr|k)?")
RE_PRICE_ABOVE = re.compile(r"(?i)(?:above|over|more\s+than|minimum|min)\s+([0-9]+(?:,[0-9]+)*(?:\.[0-9]+)?)\s*(?:etb|birr|k)?")

# Explicit advisory, purchasing recommendations, property valuation, comparison patterns
RE_EXPLICIT_ADVISORY = re.compile(
    r"(?i)\b("
    r"how (can|should|do|does|to|would)|"
    r"what (is|are|makes|should|to check|about|spec|specs)|"
    r"which (is|are|pc|laptop|computer|phone|car|apartment|vehicle|item|brand|model|option|one)|"
    r"why (does|is|are)|"
    r"recommend|suggest|advice|advise|suitable\s+for|better\s+for|best\s+for|"
    r"tips?|guide|strategy|pricing strategy|commission model|"
    r"increase (sales|trust|views)|customer trust|inspect before|what to check|"
    r"calculate|compare|explain how|how to advertise|how to attract|best way to|"
    r"is\s+.*(enough|sufficient|good|better)|should\s+i\s+(buy|rent|choose|get)"
    r")\b|"
    r"(እንዴት|ምን|ስለ|መምሪያ|ምክር|የቱ|የትኛው|ይምረጡ|akkamitti|gorsa|kamtu|ከመይ|እንታይ|ኣየናይ)"
)


def extract_price_filters(text: str) -> Tuple[Optional[float], Optional[float]]:
    clean = text.lower()
    min_p, max_p = None, None

    m_under = RE_PRICE_UNDER.search(clean)
    if m_under:
        val_str = m_under.group(1).replace(",", "")
        try:
            val = float(val_str)
            if "k" in m_under.group(0).lower():
                val *= 1000
            max_p = val
        except ValueError:
            pass

    m_above = RE_PRICE_ABOVE.search(clean)
    if m_above:
        val_str = m_above.group(1).replace(",", "")
        try:
            val = float(val_str)
            if "k" in m_above.group(0).lower():
                val *= 1000
            min_p = val
        except ValueError:
            pass

    return min_p, max_p


def extract_search_term(text: str) -> str:
    clean = text.lower()
    fillers = [
        "find me", "find", "show me", "show", "search for", "search", "looking for",
        "i want", "i need", "can you show", "available", "please", "around", "near",
        "in", "for sale", "for rent", "buy", "sell", "cheap", "good", "best",
        "እፈልጋለሁ", "አሳየኝ", "ፈልግ", "barbaada", "ደሊየ",
    ]
    term = clean
    for f in fillers:
        term = term.replace(f, " ")
    term = RE_PRICE_UNDER.sub(" ", term)
    term = RE_PRICE_ABOVE.sub(" ", term)
    term = re.sub(r"\s+", " ", term).strip()
    return term


def is_contextual_followup(text: str, history: Optional[List[Dict[str, Any]]]) -> bool:
    if not history or len(history) == 0:
        return False
    lower = text.lower().strip()
    words = lower.split()
    if len(words) > 8:
        return False
    followup_patterns = [
        r"^(and|what\s+about|how\s+about|which\s+(one|is)|is\s+it|are\s+they|how\s+much|why|tell\s+me\s+more)",
        r"(cheaper|better|more|less|difference|compare|recommended|suitable|enough)",
        r"^(እንስ|ስለ|የትኛው|ስንት|ለምን|ዋጋውስ|ከመይ|እንታይ|akkam)",
    ]
    for pat in followup_patterns:
        if re.search(pat, lower):
            return True
    return False


# ══════════════════════════════════════════════════════════════════════════════
# 4. CONTEXTUAL DETERMINISTIC DOMAIN VALIDATION ENGINE
# ══════════════════════════════════════════════════════════════════════════════

class DomainValidator:
    def __init__(self, policy: PolicyEngine):
        self.policy = policy

    def evaluate(self, message: str, requested_lang: str = "", history: Optional[List[Dict[str, Any]]] = None) -> Dict[str, Any]:
        text = message.strip()
        lang = detect_language(text, requested_lang)

        # 1. Security Check: Adversarial Bypass & Jailbreak
        if RE_JAILBREAK.search(text):
            return {
                "is_allowed": False,
                "intent": "ADVERSARIAL_BYPASS",
                "path": "PATH_BLOCKED",
                "language": lang,
                "reason": "adversarial_prompt_injection",
                "instant_reply": "HUNTIQ AI does not provide private user information, backend credentials, or bypass security rules.",
            }

        # 2. User Privacy Protection Check (Phase 15): Refuse private account / registration data
        if RE_USER_PRIVACY.search(text):
            return {
                "is_allowed": False,
                "intent": "USER_PRIVACY_PROTECTION",
                "path": "PATH_BLOCKED",
                "language": lang,
                "reason": "privacy_policy_restriction",
                "instant_reply": get_fast_privacy(lang),
                "suggested_queries": [
                    "What categories are active on HUNTIQ?",
                    "How do I create a seller account?",
                    "How are broker transactions secured?",
                ],
            }

        # 3. Fast Social Responses (< 1ms instant deterministic response, 0 Gemini calls, 0 DB queries)
        if RE_GREETING.match(text):
            return {
                "is_allowed": True,
                "intent": "GREETING",
                "path": "PATH_ADVISORY",
                "language": lang,
                "instant_reply": get_fast_greeting(lang),
                "suggested_queries": [
                    "Find laptops under 50,000 ETB",
                    "Apartments for rent in Bole",
                    "How does the HUNTIQ broker model work?",
                ],
            }

        if RE_COMPLIMENT.search(text):
            return {
                "is_allowed": True,
                "intent": "COMPLIMENT",
                "path": "PATH_ADVISORY",
                "language": lang,
                "instant_reply": get_fast_compliment(lang),
                "suggested_queries": [
                    "Which PC is suitable for a software engineer?",
                    "How do I price my apartment for rent?",
                    "What should I inspect before buying a used car?",
                ],
            }

        if RE_APPRECIATION_THANKS.search(text):
            return {
                "is_allowed": True,
                "intent": "APPRECIATION",
                "path": "PATH_ADVISORY",
                "language": lang,
                "instant_reply": get_fast_appreciation(lang),
                "suggested_queries": [
                    "How can a seller increase sales on HUNTIQ?",
                    "Calculate 10% commission on 100,000 ETB",
                    "Find Toyota cars in Addis Ababa",
                ],
            }

        if RE_ACKNOWLEDGEMENT.match(text):
            return {
                "is_allowed": True,
                "intent": "ACKNOWLEDGEMENT",
                "path": "PATH_ADVISORY",
                "language": lang,
                "instant_reply": get_fast_acknowledgement(lang),
                "suggested_queries": [
                    "Find houses for sale in Addis Ababa",
                    "What to check when buying a used phone?",
                ],
            }

        if RE_CAPABILITY_OR_IDENTITY.search(text):
            return {
                "is_allowed": True,
                "intent": "CAPABILITY_QUESTION",
                "path": "PATH_ADVISORY",
                "language": lang,
                "instant_reply": get_fast_capability(lang),
                "suggested_queries": [
                    "Find laptops under 50,000 ETB",
                    "How does broker commission work?",
                    "Tips for renting an apartment safely",
                ],
            }

        if RE_GOODBYE.search(text):
            return {
                "is_allowed": True,
                "intent": "GOODBYE",
                "path": "PATH_ADVISORY",
                "language": lang,
                "instant_reply": get_fast_goodbye(lang),
            }

        # 4. Contextual Math & Utility Exception: Allow percentage, discount, commission, engine capacity, power consumption
        has_math_exception = bool(RE_MARKETPLACE_MATH_EXCEPTION.search(text))

        # 5. Out-of-Domain Detection: Pure Math, Academic Sciences, Generic Coding, Trivia
        if not has_math_exception:
            if RE_PURE_MATH.search(text):
                return {
                    "is_allowed": False,
                    "intent": "OUT_OF_DOMAIN",
                    "path": "PATH_BLOCKED",
                    "language": lang,
                    "reason": "unrelated_general_mathematics",
                    "instant_reply": get_localized_rejection(lang, "unrelated_general_mathematics"),
                }
            if RE_PURE_SCIENCE.search(text):
                return {
                    "is_allowed": False,
                    "intent": "OUT_OF_DOMAIN",
                    "path": "PATH_BLOCKED",
                    "language": lang,
                    "reason": "unrelated_general_science",
                    "instant_reply": get_localized_rejection(lang, "unrelated_general_science"),
                }
            if RE_GENERIC_CODING.search(text):
                return {
                    "is_allowed": False,
                    "intent": "OUT_OF_DOMAIN",
                    "path": "PATH_BLOCKED",
                    "language": lang,
                    "reason": "unrelated_general_programming",
                    "instant_reply": get_localized_rejection(lang, "unrelated_general_programming"),
                }
            if RE_GENERIC_TRIVIA.search(text):
                return {
                    "is_allowed": False,
                    "intent": "OUT_OF_DOMAIN",
                    "path": "PATH_BLOCKED",
                    "language": lang,
                    "reason": "unrelated_general_trivia",
                    "instant_reply": get_localized_rejection(lang, "unrelated_general_trivia"),
                }

        # 6. Check In-Domain Marketplace Signals & Contextual Continuity
        lower = text.lower()
        has_marketplace_signal = any(sig in lower for sig in MARKETPLACE_SIGNALS)
        is_followup = is_contextual_followup(text, history)

        # 7. Ambiguity check: very short ambiguous queries without signal or context
        if not has_marketplace_signal and not has_math_exception and not is_followup:
            if len(text.split()) <= 2 and not RE_ETHIOPIC.search(text):
                return {
                    "is_allowed": False,
                    "intent": "AMBIGUOUS",
                    "path": "PATH_BLOCKED",
                    "language": lang,
                    "reason": "ambiguous_input",
                    "instant_reply": get_localized_clarification(lang),
                }
            return {
                "is_allowed": False,
                "intent": "OUT_OF_DOMAIN",
                "path": "PATH_BLOCKED",
                "language": lang,
                "reason": "unrelated_user_intent",
                "instant_reply": get_localized_rejection(lang, "unrelated_user_intent"),
            }

        # 8. Extraction of Search Filters
        min_p, max_p = extract_price_filters(text)
        search_term = extract_search_term(text)

        # Distinguish Search vs Advisory
        is_math_advisory = has_math_exception and bool(re.search(r"(?i)\b(calculate|compute|what\s+is|how\s+much|solve|formula|compare)\b", text))
        is_explicit_advisory = bool(RE_EXPLICIT_ADVISORY.search(text))
        is_advisory = is_explicit_advisory or is_math_advisory or is_followup

        if is_advisory:
            path = "PATH_ADVISORY"
            search_term = ""
            if is_math_advisory or (has_math_exception and not is_explicit_advisory):
                if any(k in lower for k in ["commission", "percentage", "discount"]):
                    intent = "COMMISSION_CALCULATION"
                else:
                    intent = "PRODUCT_MATH_ANALYSIS"
            elif any(k in lower for k in ["broker", "commission"]):
                intent = "BROKER_MODEL_EXPLANATION"
            elif any(k in lower for k in ["sell", "landlord", "advertise", "listing", "converting", "pricing"]):
                intent = "SELLER_ADVICE"
            elif any(k in lower for k in ["compare", "which", "difference", "versus", "vs"]):
                intent = "PRODUCT_COMPARISON"
            elif any(k in lower for k in ["recommend", "suggest", "suitable"]):
                intent = "RECOMMENDATION"
            elif any(k in lower for k in ["buy", "rent", "tenant", "inspect"]):
                intent = "BUYER_ADVICE"
            elif is_followup:
                intent = "CONTEXTUAL_FOLLOWUP"
            else:
                intent = "MARKETPLACE_BUSINESS_ADVICE"
        else:
            path = "PATH_SEARCH"
            intent = "MARKETPLACE_SEARCH"

        return {
            "is_allowed": True,
            "intent": intent,
            "path": path,
            "language": lang,
            "search_term": search_term,
            "min_price": min_p,
            "max_price": max_p,
            "suggested_queries": [
                "Sort by lowest price",
                "Filter by verified location in Addis Ababa",
                "What should I inspect before making a deposit?",
            ],
        }


# ══════════════════════════════════════════════════════════════════════════════
# 5. GEMINI API CLIENT (Standard Library HTTPS)
# ══════════════════════════════════════════════════════════════════════════════

class GeminiClient:
    def __init__(self, api_key: str, model: str = "gemini-3.5-flash-lite", timeout: int = 20):
        self.api_key = api_key
        self.model = model or "gemini-3.5-flash-lite"
        self.timeout = timeout
        self.base_url = "https://generativelanguage.googleapis.com/v1beta"

    def _sanitize_error(self, err_msg: str) -> str:
        if self.api_key and self.api_key in err_msg:
            return err_msg.replace(self.api_key, "[REDACTED_API_KEY]")
        return err_msg

    def generate(self, system_instruction: str, history: List[Dict[str, Any]], prompt: str) -> Tuple[str, int, Optional[str]]:
        if not self.api_key:
            return "", 0, "GEMINI_API_KEY is not configured on server"

        url = f"{self.base_url}/models/{self.model}:generateContent?key={self.api_key}"

        contents = []
        if history:
            for turn in history[-4:]:
                role = turn.get("role", "user")
                if role not in ("user", "model"):
                    role = "user"
                parts = turn.get("parts", [])
                contents.append({"role": role, "parts": parts})

        contents.append({"role": "user", "parts": [{"text": prompt}]})

        payload = {
            "contents": contents,
            "systemInstruction": {"parts": [{"text": system_instruction}]},
            "generationConfig": {
                "temperature": 0.2,
                "maxOutputTokens": 1024,
            },
        }

        req_bytes = json.dumps(payload).encode("utf-8")
        req = urllib.request.Request(
            url,
            data=req_bytes,
            headers={"Content-Type": "application/json"},
            method="POST",
        )

        # Retry transient HTTP 429/503 errors up to 2 times
        for attempt in range(1, 3):
            try:
                with urllib.request.urlopen(req, timeout=self.timeout) as resp:
                    resp_bytes = resp.read()
                    data = json.loads(resp_bytes.decode("utf-8"))

                    candidates = data.get("candidates", [])
                    if not candidates:
                        return "", 0, "No response candidates returned by Gemini"

                    content = candidates[0].get("content", {})
                    parts = content.get("parts", [])
                    reply_text = "".join(p.get("text", "") for p in parts).strip()
                    tokens_used = data.get("usageMetadata", {}).get("totalTokenCount", 0)
                    return reply_text, tokens_used, None

            except urllib.error.HTTPError as e:
                err_body = ""
                try:
                    err_body = e.read().decode("utf-8")
                except Exception:
                    pass
                sanitized = self._sanitize_error(f"HTTP {e.code}: {e.reason} - {err_body}")

                if e.code in (429, 503) and attempt < 2:
                    logger.warning("Transient Gemini error (%s). Retrying in 1s...", sanitized)
                    time.sleep(1.0)
                    continue

                return "", 0, sanitized

            except urllib.error.URLError as e:
                sanitized = self._sanitize_error(f"Network error contacting Gemini: {e.reason}")
                return "", 0, sanitized

            except Exception as e:
                sanitized = self._sanitize_error(str(e))
                return "", 0, sanitized

        return "", 0, "Failed to connect to Gemini after retries"


# ══════════════════════════════════════════════════════════════════════════════
# 6. SYSTEM PROMPT BUILDER
# ══════════════════════════════════════════════════════════════════════════════

def build_system_instruction(lang: str, path: str, trend_summary: str = "") -> str:
    base = (
        "You are HUNTIQ AI — Made by HUNTIQ. You are an authoritative marketplace intelligence and broker platform advisor.\n"
        "HUNTIQ is Ethiopia's premier broker marketplace connecting buyers, sellers, renters, and service providers.\n"
        "Operational Principles:\n"
        "1. Platform Architecture: HUNTIQ operates a broker marketplace model (not an inventory-holding retailer). Revenue is earned via broker commissions, seller subscriptions, and featured listings.\n"
        "2. Mode A (Database Authority): When matching database listings are provided below, summarize them factually without inventing fake prices, products, or seller contact numbers.\n"
        "3. Mode B (Marketplace & Commercial Advisory): When the user asks for purchasing advice, specifications (e.g. PC for software engineers, laptop RAM, vehicle inspection, apartment pricing, landlord advertising, broker commissions), provide thorough, expert, actionable guidance regardless of whether matching database listings exist.\n"
        "4. Mode C (Market Intelligence): When trend signals are provided, use them to describe marketplace demand accurately.\n"
        "5. Statistical Honesty: Distinguish views from sales, searches from purchases, and listing count from product quality. Aggregate marketplace trends reflect viewer interest and search activity, not verified transaction volume.\n"
        "6. Architecture Transparency: You are augmented with live database listings and aggregate platform signals; never claim model retraining or weight updates.\n"
        "7. Facts vs Estimates: Always distinguish real verified database prices from general market estimates or typical price ranges.\n"
        "8. No Redundant Introductions: Get straight to the answer without repetitive phrases like 'I am HUNTIQ AI, your dedicated assistant...'. Keep formatting clear and professional."
    )

    if trend_summary:
        base += f"\n\n[VERIFIED AGGREGATE MARKETPLACE SIGNALS:\n{trend_summary}\nUse these aggregate signals to provide relevant contextual guidance without claiming automatic model retraining.]"

    if lang == "am":
        base += "\nRespond primarily in clear, natural Amharic (አማርኛ). Preserve brand names and technical specs accurately."
    elif lang == "ti":
        base += "\nRespond primarily in clear, natural Tigrinya (ትግርኛ). Preserve brand names and technical specs accurately."
    elif lang == "or":
        base += "\nRespond primarily in clear, natural Afaan Oromo. Preserve brand names and technical specs accurately."
    else:
        base += "\nRespond in clear, professional English."

    return base


# ══════════════════════════════════════════════════════════════════════════════
# 7. HTTP SERVER & REQUEST HANDLER
# ══════════════════════════════════════════════════════════════════════════════

class AIRequestHandler(http.server.BaseHTTPRequestHandler):
    auth_token: str = ""
    policy_engine: PolicyEngine = None
    validator: DomainValidator = None
    gemini_client: Optional[GeminiClient] = None

    def log_message(self, format: str, *args: Any) -> None:
        pass

    def _send_json(self, status_code: int, data: Dict[str, Any]) -> None:
        payload = json.dumps(data).encode("utf-8")
        self.send_response(status_code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(payload)))
        self.send_header("Connection", "close")
        self.end_headers()
        self.wfile.write(payload)

    def _check_auth(self) -> bool:
        if not self.auth_token:
            return True
        token = self.headers.get("X-Internal-Token", "")
        return token == self.auth_token

    def do_GET(self) -> None:
        if self.path == "/health":
            self._send_json(200, {
                "status": "operational",
                "service": "HUNTIQ AI Python Engine",
                "version": "2.0",
                "languages": ["en", "am", "ti", "or"],
            })
            return

        self._send_json(404, {"error": "endpoint not found"})

    def do_POST(self) -> None:
        if not self._check_auth():
            self._send_json(401, {"error": "unauthorized: missing or invalid internal auth token"})
            return

        content_length = int(self.headers.get("Content-Length", 0))
        if content_length > 10 * 1024 * 1024:
            self._send_json(413, {"error": "payload exceeds 10MB limit"})
            return

        try:
            body = self.rfile.read(content_length)
            data = json.loads(body.decode("utf-8"))
        except Exception as e:
            self._send_json(400, {"error": f"invalid JSON payload: {str(e)}"})
            return

        if self.path == "/classify":
            self.handle_classify(data)
        elif self.path == "/chat":
            self.handle_chat(data)
        elif self.path == "/shutdown":
            self.close_connection = True
            self._send_json(200, {"status": "shutting down"})
            def do_shutdown():
                time.sleep(0.05)
                os._exit(0)
            threading.Thread(target=do_shutdown, daemon=True).start()
        else:
            self._send_json(404, {"error": "endpoint not found"})

    def handle_classify(self, data: Dict[str, Any]) -> None:
        message = data.get("message", "")
        lang = data.get("language", "")
        history = data.get("history", [])

        res = self.validator.evaluate(message, lang, history)
        self._send_json(200, res)

    def handle_chat(self, data: Dict[str, Any]) -> None:
        message = data.get("message", "")
        lang = data.get("language", "")
        history = data.get("history", [])
        listings = data.get("listings", [])
        trend_summary = data.get("trend_summary", "")

        t0 = time.time()
        # 1. Local Domain Validation (Zero Gemini Calls for Rejected Queries)
        evaluation = self.validator.evaluate(message, lang, history)

        if not evaluation.get("is_allowed", False):
            latency_ms = int((time.time() - t0) * 1000)
            self._send_json(200, {
                "reply": evaluation.get("instant_reply", ""),
                "language": evaluation.get("language", "en"),
                "intent": evaluation.get("intent", "OUT_OF_DOMAIN"),
                "path": evaluation.get("path", "PATH_BLOCKED"),
                "is_out_of_domain": True,
                "fallback_used": False,
                "tokens_used": 0,
                "latency_ms": latency_ms,
                "suggested_queries": evaluation.get("suggested_queries", []),
            })
            return

        # 2. Instant Reply (Fast Social Interactions < 1ms, 0 Gemini Calls, 0 DB Queries)
        if evaluation.get("instant_reply") and evaluation.get("intent") in (
            "GREETING", "COMPLIMENT", "APPRECIATION", "ACKNOWLEDGEMENT", "GOODBYE", "CAPABILITY_QUESTION"
        ):
            latency_ms = int((time.time() - t0) * 1000)
            self._send_json(200, {
                "reply": evaluation["instant_reply"],
                "language": evaluation.get("language", "en"),
                "intent": evaluation.get("intent"),
                "path": "PATH_ADVISORY",
                "is_out_of_domain": False,
                "fallback_used": False,
                "tokens_used": 0,
                "latency_ms": latency_ms,
                "suggested_queries": evaluation.get("suggested_queries", []),
            })
            return

        # 3. Construct Augmented Prompt with Real Database Context (Authoritative Source)
        resolved_lang = evaluation.get("language", "en")
        path = evaluation.get("path", "PATH_SEARCH")
        system_instruction = build_system_instruction(resolved_lang, path, trend_summary)

        augmented_prompt = message
        if listings:
            augmented_prompt += f"\n\n[AUTHORIZED DATABASE CONTEXT: Found {len(listings)} matching verified listings on HUNTIQ:\n"
            for c in listings[:6]:
                cid = c.get("id")
                title = c.get("title", "")
                price = c.get("price", 0)
                curr = c.get("currency", "ETB")
                cat = c.get("category_name", "")
                loc = c.get("location_name", "")
                augmented_prompt += f"- #{cid}: {title} | Price: {price:.2f} {curr} | Category: {cat} | Location: {loc}\n"
            augmented_prompt += "Concisely explain these listings to the user and offer helpful marketplace recommendations.]"
        elif path == "PATH_SEARCH" or (evaluation.get("search_term") and not listings):
            # Phase 7: Strict Zero-Result Handling
            search_label = evaluation.get("search_term") or message
            augmented_prompt += (
                f"\n\n[DATABASE CONTEXT: 0 matching listings found in database for '{search_label}'.\n"
                "STRICT INSTRUCTIONS:\n"
                "1. State factually and clearly that no listings currently match this exact search on HUNTIQ.\n"
                "2. DO NOT fabricate or hallucinate specific listings, inventory items, or fake seller names.\n"
                "3. Provide helpful purchasing advice, recommended specifications, or what to look for when shopping for this item.\n"
                "4. If mentioning typical price ranges, explicitly label them as general market estimates, not verified listings.\n"
                "5. Suggest broader search terms or related categories on HUNTIQ.]"
            )

        # 4. Invoke Google Gemini via Standard-Library HTTPS Client
        reply_text, tokens_used, err = self.gemini_client.generate(system_instruction, history, augmented_prompt)
        latency_ms = int((time.time() - t0) * 1000)

        if err:
            logger.error("Gemini call error: %s", err)
            fallback = get_localized_rejection(resolved_lang, "service_unavailable")
            if listings:
                fallback = f"I retrieved {len(listings)} verified listings from HUNTIQ for you below:"
            self._send_json(200, {
                "reply": fallback,
                "language": resolved_lang,
                "intent": evaluation.get("intent", "MARKETPLACE_SEARCH"),
                "path": path,
                "is_out_of_domain": False,
                "fallback_used": True,
                "tokens_used": 0,
                "latency_ms": latency_ms,
                "error": err,
            })
            return

        self._send_json(200, {
            "reply": reply_text,
            "language": resolved_lang,
            "intent": evaluation.get("intent", "MARKETPLACE_SEARCH"),
            "path": path,
            "is_out_of_domain": False,
            "fallback_used": False,
            "tokens_used": tokens_used,
            "latency_ms": latency_ms,
            "suggested_queries": evaluation.get("suggested_queries", []),
        })


# ══════════════════════════════════════════════════════════════════════════════
# 8. SERVER INITIALIZATION & MAIN ENTRY
# ══════════════════════════════════════════════════════════════════════════════

class ThreadedHTTPServer(socketserver.ThreadingMixIn, http.server.HTTPServer):
    daemon_threads = True


def run_server(port: int, auth_token: str, policy_path: str, gemini_key: str, gemini_model: str) -> None:
    policy_engine = PolicyEngine(policy_path)
    validator = DomainValidator(policy_engine)
    gemini_client = GeminiClient(api_key=gemini_key, model=gemini_model)

    AIRequestHandler.auth_token = auth_token
    AIRequestHandler.policy_engine = policy_engine
    AIRequestHandler.validator = validator
    AIRequestHandler.gemini_client = gemini_client

    server = ThreadedHTTPServer(("127.0.0.1", port), AIRequestHandler)
    actual_port = server.server_address[1]
    logger.info("HUNTIQ AI Engine listening on http://127.0.0.1:%d", actual_port)
    print(f"HUNTIQ_AI_PORT={actual_port}", flush=True)

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        logger.info("Keyboard interrupt received. Shutting down...")
    finally:
        server.server_close()
        logger.info("HUNTIQ AI Engine stopped cleanly.")


def main():
    parser = argparse.ArgumentParser(description="HUNTIQ AI Standalone Engine")
    parser.add_argument("--port", type=int, default=0, help="Loopback port to bind (default 0 for random free port)")
    parser.add_argument("--auth-token", type=str, default="", help="Internal authorization token")
    parser.add_argument("--policy", type=str, default="backend/ai/policy.json", help="Path to policy.json")
    parser.add_argument("--gemini-key", type=str, default="", help="Google Gemini API key")
    parser.add_argument("--gemini-model", type=str, default="gemini-3.5-flash-lite", help="Gemini Model")
    args = parser.parse_args()

    key = args.gemini_key or os.environ.get("GEMINI_API_KEY", "")
    model = args.gemini_model or os.environ.get("GEMINI_MODEL", "gemini-3.5-flash-lite")

    run_server(args.port, args.auth_token, args.policy, key, model)


if __name__ == "__main__":
    main()
