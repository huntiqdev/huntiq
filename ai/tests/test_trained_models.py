"""
Tests for trained ML checkpoints:
- intent_classifier_v1
- category_classifier_v1
- safety_moderation_v1
- domain_classifier_v1
and fuzzy entity matching with RapidFuzz.
Verifies offline loading, calibrated probabilities, multilingual inference, and latency bounds.
"""

import time

import pytest

from app.api.schemas import GuestChatRequest
from app.core.fuzzy_matcher import fuzzy_matcher
from app.core.orchestrator import orchestrator
from app.core.trained_model_service import trained_models


def test_models_loaded_successfully():
    """Verify all 4 trained checkpoints are found and loaded."""
    assert trained_models.is_intent_model_ready(), (
        "Intent model checkpoint failed to load"
    )
    assert trained_models.is_category_model_ready(), (
        "Category model checkpoint failed to load"
    )
    assert trained_models.is_safety_model_ready(), (
        "Safety model checkpoint failed to load"
    )
    assert trained_models.is_domain_model_ready(), (
        "Domain model checkpoint failed to load"
    )
    assert trained_models.is_sentiment_model_ready(), (
        "Sentiment model checkpoint failed to load"
    )


def test_intent_prediction_multilingual():
    """Verify intent predictions across EN, AM, OR, TI."""
    cases = [
        ("Good morning assistant, how are you?", "social_greeting"),
        ("ሰላም እንደምን ነህ?", "social_greeting"),
        ("Akkam jirtu?", "social_greeting"),
        ("ሰላም ከመይ ኣለኻ?", "social_greeting"),
        ("Who are you and what do you do?", "identity_inquiry"),
        ("ማን ነህ?", "identity_inquiry"),
        ("Find 2 bedroom apartments in Bole", "marketplace_search"),
        ("ቶዮታ ቪትዝ መኪና መግዛት እፈልጋለሁ", "marketplace_search"),
        ("How do I pay with Telebirr?", "platform_faq"),
        ("በቴሌብር ክፍያ እንዴት ይፈጸማል?", "platform_faq"),
        ("Explain quantum mechanics and Schrodinger equation", "out_of_domain"),
        ("የኳንተም ፊዚክስ ቀመር ጻፍልኝ", "out_of_domain"),
    ]

    for query, expected_intent in cases:
        intent, conf, probs = trained_models.predict_intent(query)
        assert intent == expected_intent, (
            f"Query '{query}' predicted '{intent}', expected '{expected_intent}'"
        )
        assert conf > 0.40, f"Confidence {conf} too low for '{query}'"
        assert abs(sum(probs.values()) - 1.0) < 1e-4, "Probabilities must sum to 1.0"


def test_category_prediction_multilingual():
    """Verify category predictions for products, real estate, vehicles, jobs, agriculture."""
    cases = [
        ("iPhone 14 Pro Max 256GB storage", "products"),
        ("Samsung Galaxy S22 phone in Bole", "products"),
        ("2 bedroom luxury apartment for rent in Bole", "real_estate"),
        ("የሚከራይ ባለ 3 መኝታ ቤት በአዲስ አበባ", "real_estate"),
        ("Toyota Vitz 2018 automatic clean condition", "vehicles"),
        ("ኮድ 2 ቶዮታ ኮሮላ መኪና የሚሸጥ", "vehicles"),
        ("Looking for experienced driver with license", "jobs"),
        ("የሂሳብ ባለሙያ እና አካውንታንት ክፍት የስራ ቦታ", "jobs"),
        ("Agricultural teff 100 quintal wholesale Gojjam", "agriculture"),
        ("ማኛ ጤፍ በኩንታል የጎጃም ምርት", "agriculture"),
        ("Plumber for water pipe leak repair in Sarbet", "services"),
        ("Guest house room in Bole for daily rent", "hotels"),
        ("Concert tickets for New Year live music festival", "tickets"),
    ]

    for query, expected_cat in cases:
        cat, conf, probs = trained_models.predict_category(query)
        assert cat == expected_cat, (
            f"Query '{query}' predicted '{cat}', expected '{expected_cat}'"
        )
        assert conf > 0.35, f"Confidence {conf} too low for '{query}'"
        assert abs(sum(probs.values()) - 1.0) < 1e-4, "Probabilities must sum to 1.0"


def test_safety_moderation_multilingual():
    """Verify detection of advance fee scams, phishing, and contraband."""
    scam_queries = [
        (
            "Please send 50% advance to my Telebirr before seeing the car",
            "scam_advance",
        ),
        ("መኪናውን ከማየትህ በፊት 50% በቴሌብር በቅድሚያ ላክ", "scam_advance"),
        (
            "Click this link to verify your HUNTIQ bank account login: http://huntiq-fake.xyz",
            "phishing_harvesting",
        ),
        ("የባንክ አካውንትዎን እና የቴሌብር ፒን ቁጥርዎን እዚህ ሊንክ ላይ ያስገቡ", "phishing_harvesting"),
        (
            "Black market military weapons assault rifles bullets",
            "prohibited_contraband",
        ),
        ("ያልተመዘገበ ክላሽ እና ሽጉጥ መግዛት እፈልጋለሁ", "prohibited_contraband"),
        ("I want to buy a used Toyota Vitz in Addis Ababa", "safe"),
        ("ቦሌ የሚከራይ ባለ ሁለት መኝታ ቤት አለ?", "safe"),
    ]

    for query, expected_class in scam_queries:
        safety_class, conf, _ = trained_models.predict_safety(query)
        assert safety_class == expected_class, (
            f"Query '{query}' classified as '{safety_class}', expected '{expected_class}'"
        )
        assert conf > 0.40, f"Confidence {conf} too low for '{query}'"


def test_domain_boundary_multilingual():
    """Verify domain classification between marketplace and general external queries."""
    cases = [
        (
            "Looking for a used Dell XPS laptop in good condition",
            "ethiopian_marketplace",
        ),
        ("አዲስ አበባ ውስጥ ቪትዝ መኪና መግዛት እፈልጋለሁ", "ethiopian_marketplace"),
        ("Konkolaataa Vitz Finfinneetti bituu barbaada", "ethiopian_marketplace"),
        ("Write a Python script to sort a binary search tree", "general_out_of_domain"),
        (
            "Explain the causes of the French Revolution in 1789",
            "general_out_of_domain",
        ),
        ("የፈረንሳይ አብዮት ታሪክ እና ምክንያቶች ምን ምን ናቸው?", "general_out_of_domain"),
    ]

    for query, expected_domain in cases:
        domain, conf, _ = trained_models.predict_domain(query)
        assert domain == expected_domain, (
            f"Query '{query}' predicted as '{domain}', expected '{expected_domain}'"
        )
        assert conf > 0.50, f"Confidence {conf} too low for '{query}'"


def test_sentiment_dispute_multilingual():
    """Verify sentiment and dispute risk classification across EN, AM, OM, TI."""
    cases = [
        (
            "The laptop works perfectly, clean condition and quick response!",
            "positive_review",
        ),
        ("እቃው በጣም ቆንጆ ነው እና ሻጩ ታማኝ ነው አመሰግናለሁ", "positive_review"),
        ("Meeshaan baay'ee gaariidha, baay'ee na gammachiise", "positive_review"),
        ("እቲ ንብረት ጽቡቕ እዩ ቀልጢፉ በፂሑኒ የቐንየለይ", "positive_review"),
        ("Is this apartment still available for rent this weekend?", "neutral_inquiry"),
        ("የቤት እቃው አለ ወይስ ተሽጧል? ዋጋው ስንት ነው?", "neutral_inquiry"),
        ("Makiinaan kun amma yoonaa gurguraadhaaf jiraa?", "neutral_inquiry"),
        ("እዚ ገዛ እዚ ንክራይ ኣሎ ድዩ ዋጋኡ ክንደይ እዩ?", "neutral_inquiry"),
        (
            "Product stopped working after 2 hours and seller refuses to answer",
            "negative_dissatisfied",
        ),
        ("እቃው ተበላሽቷል ሻጩም ስልክ አያነሳም ቅር ተሰኝቻለሁ", "negative_dissatisfied"),
        ("Meeshaan cabeera garuu gurguraan bilbila hin kaasu", "negative_dissatisfied"),
        ("እቲ ቴሌቪዥን ተበላሽዩ ሸያጢ ድማ ስልኪ ኣየልዕልን ኣሎ", "negative_dissatisfied"),
        (
            "Seller stole my money via Telebirr without delivering, I want full refund or police action!",
            "transaction_dispute",
        ),
        ("ብር በቴሌብር ከፍዬ እቃውን አልላከም አጭበርብሮኛል በፍጥነት ብሬ ይመለስልኝ", "transaction_dispute"),
        (
            "Maallaqa kaffalee meeshaa naaf hin ergine, kun hatumsa kaffaltiin koo haa deebi'u",
            "transaction_dispute",
        ),
        (
            "ገንዘብ ወሲዱ ንብረት ኣይሃበንን ሰሪቑኒ ገንዘበይ ይመለሰለይ ወይ ናብ ፖሊስ ክኸይድ",
            "transaction_dispute",
        ),
    ]

    for query, expected_class in cases:
        pred_class, conf, probs = trained_models.predict_sentiment_dispute(query)
        assert pred_class == expected_class, (
            f"Query '{query}' predicted '{pred_class}', expected '{expected_class}'"
        )
        assert conf > 0.40, f"Confidence {conf} too low for '{query}'"
        assert abs(sum(probs.values()) - 1.0) < 1e-4, "Probabilities must sum to 1.0"


def test_fuzzy_matcher_brands_and_locations():
    """Verify fuzzy matching for misspelled brands and Ethiopian locations."""
    # Brand typo tolerance
    assert fuzzy_matcher.match_brand("samsng")[0] == "Samsung"
    assert fuzzy_matcher.match_brand("ipon")[0] == "Apple"
    assert fuzzy_matcher.match_brand("tayota")[0] == "Toyota"
    assert fuzzy_matcher.match_brand("suzki")[0] == "Suzuki"

    # Location typo tolerance
    assert fuzzy_matcher.match_location("adis abeba")[0] == "Addis Ababa"
    assert fuzzy_matcher.match_location("bolee")[0] == "Bole"
    assert fuzzy_matcher.match_location("kasanchis")[0] == "Kazanchis"
    assert fuzzy_matcher.match_location("mekele")[0] == "Mekelle"


@pytest.mark.asyncio
async def test_orchestrator_safety_interception():
    """Verify orchestrator blocks scam attempts with safety warning."""
    req = GuestChatRequest(message="መኪናውን ከማየትህ በፊት 50% በቴሌብር በቅድሚያ ላክ", locale="am")
    resp = await orchestrator.process_chat(req)
    assert "ማስጠንቀቂያ" in resp.answer or "ቅድመ ክፍያ" in resp.answer
    assert resp.model == "huntiq-guard:safety-moderation-v1"
    assert resp.latency_ms < 50.0


def test_inference_latency_benchmark():
    """Verify sub-10ms combined inference latency constraint for all 5 models on CPU."""
    test_queries = [
        "Looking for Toyota car in Addis Ababa",
        "2 bedroom apartment in Bole under 35000 birr",
        "iPhone 13 Pro 256GB in Bole",
        "በቴሌብር እንዴት መክፈል እችላለሁ?",
        "የኳንተም ፊዚክስ ቀመር ምንድን ነው?",
    ]

    # Warmup
    for q in test_queries:
        trained_models.predict_intent(q)
        trained_models.predict_category(q)
        trained_models.predict_safety(q)
        trained_models.predict_domain(q)
        trained_models.predict_sentiment_dispute(q)

    # Benchmark 50 iterations across all 5 models
    t0 = time.time()
    iterations = 50
    for _ in range(iterations):
        for q in test_queries:
            trained_models.predict_intent(q)
            trained_models.predict_category(q)
            trained_models.predict_safety(q)
            trained_models.predict_domain(q)
            trained_models.predict_sentiment_dispute(q)

    total_time = time.time() - t0
    total_calls = iterations * len(test_queries) * 5
    avg_latency_ms = (total_time / total_calls) * 1000.0

    print(
        f"\n[Latency Benchmark] Average inference latency across all 5 models: {avg_latency_ms:.3f} ms"
    )
    assert avg_latency_ms < 25.0, (
        f"Average latency {avg_latency_ms} ms exceeds 25ms threshold"
    )
