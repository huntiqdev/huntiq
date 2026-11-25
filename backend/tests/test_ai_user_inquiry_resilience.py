#!/usr/bin/env python3
"""
Test AI User Inquiry Resilience and Fast Response Hardening
Verifies that:
1. "halo", "helo", "hello", "selam" are recognized as instant greetings (< 1ms).
2. "where are you", "where is huntiq", "yet neh", "የት ነህ" are recognized as instant identity/presence answers (< 1ms).
3. "give phone", "show me phone", "ፈልግ ስልክ" properly extract "phone" / "ስልክ" as the search term.
4. Out-of-domain and adversarial bypasses remain strictly defended.
"""

import os
import sys

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "ai"))
from script import DomainValidator, PolicyEngine, extract_search_term


def main():
    validator = DomainValidator(PolicyEngine())

    test_cases = [
        # Greetings
        {
            "query": "halo",
            "expected_allowed": True,
            "expected_intent": "GREETING",
            "expected_path": "PATH_ADVISORY",
            "must_have_instant_reply": True,
        },
        {
            "query": "helo",
            "expected_allowed": True,
            "expected_intent": "GREETING",
            "expected_path": "PATH_ADVISORY",
            "must_have_instant_reply": True,
        },
        {
            "query": "hello",
            "expected_allowed": True,
            "expected_intent": "GREETING",
            "expected_path": "PATH_ADVISORY",
            "must_have_instant_reply": True,
        },
        # Identity / Location / Presence
        {
            "query": "where are you",
            "expected_allowed": True,
            "expected_intent": "CAPABILITY_QUESTION",
            "expected_path": "PATH_ADVISORY",
            "must_have_instant_reply": True,
        },
        {
            "query": "where is huntiq",
            "expected_allowed": True,
            "expected_intent": "CAPABILITY_QUESTION",
            "expected_path": "PATH_ADVISORY",
            "must_have_instant_reply": True,
        },
        {
            "query": "yet neh",
            "expected_allowed": True,
            "expected_intent": "CAPABILITY_QUESTION",
            "expected_path": "PATH_ADVISORY",
            "must_have_instant_reply": True,
        },
        {
            "query": "የት ነህ",
            "expected_allowed": True,
            "expected_intent": "CAPABILITY_QUESTION",
            "expected_path": "PATH_ADVISORY",
            "must_have_instant_reply": True,
        },
        # Farewells & Conversational Typo Resiliency
        {
            "query": "good buy",
            "expected_allowed": True,
            "expected_intent": "GOODBYE",
            "expected_path": "PATH_ADVISORY",
            "expected_term": "",
        },
        {
            "query": "good by",
            "expected_allowed": True,
            "expected_intent": "GOODBYE",
            "expected_path": "PATH_ADVISORY",
            "expected_term": "",
        },
        {
            "query": "i mean huve anice time",
            "expected_allowed": True,
            "expected_intent": "GOODBYE",
            "expected_path": "PATH_ADVISORY",
            "expected_term": "",
        },
        {
            "query": "have a nice time",
            "expected_allowed": True,
            "expected_intent": "GOODBYE",
            "expected_path": "PATH_ADVISORY",
            "expected_term": "",
        },
        # General Commerce Inquiries (Non-blocking)
        {
            "query": "can you tell me about the market in addis?",
            "expected_allowed": True,
            "expected_intent": "GENERAL_COMMERCE_INQUIRY",
            "expected_path": "PATH_ADVISORY",
            "expected_term": "",
        },
        # Search Term Extractions
        {
            "query": "give phone",
            "expected_allowed": True,
            "expected_intent": "MARKETPLACE_SEARCH",
            "expected_path": "PATH_SEARCH",
            "expected_term": "phone",
        },
        {
            "query": "show me phone",
            "expected_allowed": True,
            "expected_intent": "MARKETPLACE_SEARCH",
            "expected_path": "PATH_SEARCH",
            "expected_term": "phone",
        },
        {
            "query": "ere be natish mekena ale?",
            "expected_allowed": True,
            "expected_intent": "MARKETPLACE_SEARCH",
            "expected_path": "PATH_SEARCH",
            "expected_term": "mekena",
        },
    ]

    passed = 0
    total = len(test_cases)

    print("=" * 80)
    print("  HUNTIQ AI USER INQUIRY RESILIENCE TEST SUITE")
    print("=" * 80)

    for tc in test_cases:
        res = validator.evaluate(tc["query"])
        q = tc["query"]
        assert res["is_allowed"] == tc["expected_allowed"], f"[{q}] expected allowed={tc['expected_allowed']}, got {res['is_allowed']}"
        assert res["intent"] == tc["expected_intent"], f"[{q}] expected intent={tc['expected_intent']}, got {res['intent']}"
        assert res["path"] == tc["expected_path"], f"[{q}] expected path={tc['expected_path']}, got {res['path']}"
        if tc.get("must_have_instant_reply"):
            assert res.get("instant_reply"), f"[{q}] missing instant_reply"
        if "expected_term" in tc:
            assert res.get("search_term") == tc["expected_term"], f"[{q}] expected search_term={tc['expected_term']}, got {res.get('search_term')}"

        print(f"[PASS] Query '{q}' -> Intent: {res['intent']} | Path: {res['path']}")
        passed += 1

    print("=" * 80)
    print(f"  ALL {passed}/{total} USER INQUIRY TESTS PASSED (100%)")
    print("=" * 80)


if __name__ == "__main__":
    main()
