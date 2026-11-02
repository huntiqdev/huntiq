"""Tests for Python AI intent classifier, policy, context manager, and orchestrator."""

import pytest

from app.api.schemas import GuestChatRequest
from app.core.context_manager import SessionManager
from app.core.entity_extractor import (
    extract_location,
    extract_price_bounds,
    extract_searches,
)
from app.core.intent import IntentType, classify_message
from app.core.orchestrator import orchestrator
from app.core.policy import check_platform_faq, is_out_of_domain


def test_intent_greetings_multilingual():
    """Verify social greetings in all four supported languages."""
    cases = [
        ("Hello, how are you?", "en"),
        ("ሰላም እንደምን ነህ?", "am"),
        ("Akkam jirtu?", "om"),
        ("ሰላም ከመይ ኣለኻ?", "ti"),
    ]
    for text, loc in cases:
        plan = classify_message(text, loc)
        assert IntentType.SOCIAL_GREETING in plan.intents
        assert not plan.requires_database_search
        assert plan.social_greeting is not None


def test_intent_identity_multilingual():
    """Verify identity inquiries in all four languages."""
    cases = [
        ("Who are you and what do you do?", "en"),
        ("ማን ነህ? ምንስ ትሰራለህ?", "am"),
        ("Eenyu ati?", "om"),
        ("መን ኢኻ?", "ti"),
    ]
    for text, loc in cases:
        plan = classify_message(text, loc)
        assert IntentType.IDENTITY_INQUIRY in plan.intents
        assert not plan.requires_database_search


def test_out_of_domain_detection():
    """Verify coding and academic requests are caught."""
    assert is_out_of_domain("Write a Python script to reverse a string")
    assert is_out_of_domain("What is the quadratic formula?")
    assert is_out_of_domain("የኳንተም ፊዚክስ ቀመር ምንድን ነው?")
    # Legitimate marketplace queries with numbers should NOT be out of domain
    assert not is_out_of_domain("Find a phone under 20,000 birr in Addis Ababa")
    assert not is_out_of_domain("Toyota Vitz price and specs")


def test_platform_faq_detection():
    """Verify HUNTIQ FAQs are answered."""
    is_faq, ans = check_platform_faq("How do I publish a listing on HUNTIQ?", "en")
    assert is_faq
    assert "Post Listing" in ans

    is_faq_am, ans_am = check_platform_faq("በHUNTIQ ላይ ዕቃዎችን እንዴት መሸጥ እችላለሁ?", "am")
    assert is_faq_am
    assert "ማስታወቂያ ልጠፍ" in ans_am


def test_entity_extraction_prices_and_locations():
    """Verify extraction of prices with currencies and locations."""
    # Under pattern
    _min_p, max_p = extract_price_bounds(
        "Find a Samsung phone under 25,000 birr in Bole"
    )
    assert max_p == 25000.0
    loc = extract_location("Find a Samsung phone under 25,000 birr in Bole")
    assert loc == "Bole"

    # Suffixes
    _min_p2, max_p2 = extract_price_bounds("budget of 20k ETB")
    assert max_p2 == 20000.0

    # Amharic
    _min_am, max_am = extract_price_bounds("ከ 5000 ብር በታች ስልክ")
    assert max_am == 5000.0


def test_multi_category_extraction():
    """Verify decomposition of multi-intent queries."""
    searches = extract_searches(
        "Find apartments for rent and land for sale in Addis Ababa"
    )
    assert len(searches) >= 2
    cat_names = [s.category_name for s in searches]
    assert "Property Rentals" in cat_names or "Real Estate & Housing" in cat_names
    assert "Land & Plots" in cat_names


def test_housing_ambiguity_clarification():
    """Verify that unspecified housing requests prompt for rent vs. buy."""
    plan = classify_message("I am looking for a house in Addis Ababa", "en")
    assert plan.requires_clarification
    assert plan.unresolved_slot == "listing_type"
    assert "rent or buy" in plan.clarification_question.lower()


def test_context_followup_cheaper():
    """Verify contextual adjustment when user asks for cheaper alternative."""
    sm = SessionManager()
    session_id = "test-session-1"

    searches = extract_searches("Samsung phone under 20000 birr")
    sm.record_turn(
        session_id, "Samsung phone under 20000 birr", "Here are 3 phones", searches
    )

    # User says "Show me cheaper ones"
    updated = sm.apply_contextual_followup(session_id, "Show me cheaper ones", [])
    assert len(updated) == 1
    assert updated[0].max_price == 16000.0  # 20000 * 0.8


@pytest.mark.asyncio
async def test_orchestrator_conversational_turn():
    """Verify full end-to-end orchestrator pipeline on greetings."""
    req = GuestChatRequest(message="Hello HUNTIQ", locale="en")
    resp = await orchestrator.process_chat(req)
    assert resp.session_id is not None
    assert "Hello" in resp.answer or "HUNTIQ" in resp.answer
    assert resp.detected_locale == "en"
    assert "social_greeting" in resp.intents
    assert resp.listings is None  # Zero DB search for greeting!
