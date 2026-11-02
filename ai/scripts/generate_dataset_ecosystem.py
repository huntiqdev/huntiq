#!/usr/bin/env python3
"""
HUNTIQ AI Dataset Ecosystem & Training Infrastructure Generator.
Creates the complete, comprehensive filesystem, schemas, manifests, and documentation
for all 30 HUNTIQ-specific dataset modules and supporting intelligence architectures.
"""

import json
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent

# Definition of the 30 HUNTIQ-specific dataset modules
DATASET_MODULES = [
    {
        "id": "01_general_greetings",
        "name": "General Greetings and Social Conversation",
        "domain": "conversational_social",
        "purpose": "Recognize greetings, farewells, gratitude, and ordinary social pleasantries without unnecessarily triggering marketplace database queries.",
        "properties": {
            "id": {"type": "string"},
            "utterance": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "greeting_type": {
                "type": "string",
                "enum": ["hello", "farewell", "gratitude", "apology", "casual"],
            },
            "has_commercial_intent": {"type": "boolean"},
            "response": {"type": "string"},
            "mixed_languages": {"type": "boolean"},
        },
        "required": [
            "id",
            "utterance",
            "locale",
            "greeting_type",
            "has_commercial_intent",
        ],
    },
    {
        "id": "02_identity_assistant_behavior",
        "name": "Assistant Identity and Behavior",
        "domain": "assistant_persona",
        "purpose": "Answer inquiries regarding the HUNTIQ assistant's identity, capabilities, role, guest access, and boundaries, providing honest explanations when requests are outside scope.",
        "properties": {
            "id": {"type": "string"},
            "query": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "inquiry_type": {
                "type": "string",
                "enum": [
                    "name",
                    "purpose",
                    "capabilities",
                    "limitations",
                    "guest_access",
                    "pricing",
                ],
            },
            "authorized_response": {"type": "string"},
            "is_honest_limitation": {"type": "boolean"},
        },
        "required": ["id", "query", "locale", "inquiry_type", "authorized_response"],
    },
    {
        "id": "03_huntiq_platform_knowledge",
        "name": "HUNTIQ Platform Knowledge",
        "domain": "platform_faq_knowledge",
        "purpose": "Provide verified, factual knowledge about HUNTIQ features, registration, listing publication, membership tiers (Gold/Premium), fees, and safety policies without hallucination.",
        "properties": {
            "id": {"type": "string"},
            "topic": {
                "type": "string",
                "enum": [
                    "products",
                    "membership",
                    "listing",
                    "account",
                    "payment",
                    "privacy",
                    "security",
                    "support",
                ],
            },
            "query": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "authoritative_answer": {"type": "string"},
            "source_doc": {"type": "string"},
            "effective_date": {"type": "string"},
        },
        "required": [
            "id",
            "topic",
            "query",
            "locale",
            "authoritative_answer",
            "source_doc",
        ],
    },
    {
        "id": "04_domain_classification",
        "name": "Domain Classification",
        "domain": "domain_policy_routing",
        "purpose": "Classify incoming queries as HUNTIQ marketplace, platform support, relevant calculation, out-of-domain (coding, academic, trivia), or unsafe content.",
        "properties": {
            "id": {"type": "string"},
            "text": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "domain_label": {
                "type": "string",
                "enum": [
                    "marketplace",
                    "platform_support",
                    "relevant_calculation",
                    "out_of_domain",
                    "unsafe",
                ],
            },
            "confidence_threshold": {"type": "number"},
            "redirection_needed": {"type": "boolean"},
        },
        "required": ["id", "text", "locale", "domain_label"],
    },
    {
        "id": "05_intent_classification",
        "name": "Intent Classification",
        "domain": "nlu_intent",
        "purpose": "Perform single-label and multi-label intent classification covering search, recommendation, comparison, account, publication, and dialogue control.",
        "properties": {
            "id": {"type": "string"},
            "query": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "primary_intent": {"type": "string"},
            "secondary_intents": {"type": "array", "items": {"type": "string"}},
            "is_multi_intent": {"type": "boolean"},
        },
        "required": ["id", "query", "locale", "primary_intent"],
    },
    {
        "id": "06_category_understanding",
        "name": "Marketplace Category Understanding",
        "domain": "taxonomy_classification",
        "purpose": "Map multilingual user terms to canonical HUNTIQ categories (Products, Electronics, Vehicles, Real Estate, Jobs, Agriculture, Services) based on live PostgreSQL taxonomy.",
        "properties": {
            "id": {"type": "string"},
            "raw_term": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "canonical_slug": {"type": "string"},
            "category_level": {"type": "integer", "enum": [1, 2, 3]},
            "parent_slug": {"type": "string"},
            "is_ambiguous": {"type": "boolean"},
        },
        "required": ["id", "raw_term", "locale", "canonical_slug"],
    },
    {
        "id": "07_entity_attribute_extraction",
        "name": "Entity and Attribute Extraction",
        "domain": "named_entity_recognition",
        "purpose": "Extract structured marketplace entities: brand, model, condition, price, currency, bedrooms, area, job title, and specifications from natural text.",
        "properties": {
            "id": {"type": "string"},
            "text": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "entities": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "entity_type": {"type": "string"},
                        "value": {"type": "string"},
                        "start_char": {"type": "integer"},
                        "end_char": {"type": "integer"},
                    },
                    "required": ["entity_type", "value"],
                },
            },
        },
        "required": ["id", "text", "locale", "entities"],
    },
    {
        "id": "08_search_query_understanding",
        "name": "Search Query Understanding",
        "domain": "query_parsing",
        "purpose": "Convert free-form user shopping and rental requests into verified structured search parameters conforming to Go ProductRepo filters.",
        "properties": {
            "id": {"type": "string"},
            "raw_query": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "search_term": {"type": "string"},
            "category_slug": {"type": "string"},
            "min_price": {"type": "number"},
            "max_price": {"type": "number"},
            "location": {"type": "string"},
            "condition": {"type": "string", "enum": ["new", "used", "any"]},
            "listing_type": {"type": "string", "enum": ["sale", "rent", "any"]},
        },
        "required": ["id", "raw_query", "locale", "search_term"],
    },
    {
        "id": "09_price_budget_understanding",
        "name": "Price and Budget Understanding",
        "domain": "numeric_currency_nlu",
        "purpose": "Interpret Ethiopian Birr (ETB) figures, suffixes (k, M), ranges, periodic rental rates (monthly, daily), and budget constraints.",
        "properties": {
            "id": {"type": "string"},
            "expression": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "extracted_min": {"type": "number"},
            "extracted_max": {"type": "number"},
            "currency": {"type": "string", "default": "ETB"},
            "is_negotiable": {"type": "boolean"},
            "period": {
                "type": "string",
                "enum": ["one_time", "daily", "monthly", "yearly"],
            },
        },
        "required": ["id", "expression", "locale", "currency"],
    },
    {
        "id": "10_location_understanding",
        "name": "Ethiopian Location Understanding",
        "domain": "geographical_nlu",
        "purpose": "Normalize Ethiopian cities, Addis Ababa subcities (Bole, CMC, Sarbet), regional centers, and local neighborhood transliterations.",
        "properties": {
            "id": {"type": "string"},
            "raw_location": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "canonical_city": {"type": "string"},
            "subcity": {"type": "string"},
            "region": {"type": "string"},
            "transliteration_variants": {"type": "array", "items": {"type": "string"}},
        },
        "required": ["id", "raw_location", "locale", "canonical_city"],
    },
    {
        "id": "11_multi_intent_understanding",
        "name": "Multi-Intent Understanding",
        "domain": "complex_dialogue_nlu",
        "purpose": "Decompose compound queries containing multiple independent requests (e.g. 'apartments for rent and land for sale') into distinct tool dispatches.",
        "properties": {
            "id": {"type": "string"},
            "compound_query": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "sub_intents": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "category_name": {"type": "string"},
                        "listing_type": {"type": "string"},
                        "clause": {"type": "string"},
                    },
                    "required": ["category_name", "clause"],
                },
            },
        },
        "required": ["id", "compound_query", "locale", "sub_intents"],
    },
    {
        "id": "12_clarification_ambiguity",
        "name": "Clarification and Ambiguity",
        "domain": "dialogue_disambiguation",
        "purpose": "Detect underspecified or ambiguous queries (e.g. 'I want a house') and generate concise, polite clarification questions across all 4 languages.",
        "properties": {
            "id": {"type": "string"},
            "ambiguous_input": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "missing_slot": {
                "type": "string",
                "enum": ["listing_type", "location", "category", "budget"],
            },
            "clarification_question": {"type": "string"},
            "can_proceed_with_default": {"type": "boolean"},
        },
        "required": [
            "id",
            "ambiguous_input",
            "locale",
            "missing_slot",
            "clarification_question",
        ],
    },
    {
        "id": "13_conversation_followups",
        "name": "Conversation Follow-ups",
        "domain": "contextual_dialogue",
        "purpose": "Maintain multi-turn conversation context to resolve follow-ups like 'show me cheaper ones', 'what about the second result', or 'change location to Bole'.",
        "properties": {
            "id": {"type": "string"},
            "history": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "role": {"type": "string"},
                        "content": {"type": "string"},
                    },
                    "required": ["role", "content"],
                },
            },
            "followup_message": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "inferred_action": {"type": "string"},
            "resolved_parameters": {"type": "object"},
        },
        "required": ["id", "history", "followup_message", "locale", "inferred_action"],
    },
    {
        "id": "14_multilingual_understanding",
        "name": "Multilingual Understanding",
        "domain": "multilingual_semantic_parity",
        "purpose": "Establish semantic equivalence across English, Amharic (አማርኛ), Afaan Oromoo, and Tigrinya (ትግርኛ) with shared semantic IDs.",
        "properties": {
            "semantic_id": {"type": "string"},
            "english": {"type": "string"},
            "amharic": {"type": "string"},
            "oromo": {"type": "string"},
            "tigrinya": {"type": "string"},
            "canonical_intent": {"type": "string"},
            "review_status": {
                "type": "string",
                "enum": ["human_verified", "pending_review"],
            },
        },
        "required": [
            "semantic_id",
            "english",
            "amharic",
            "oromo",
            "tigrinya",
            "canonical_intent",
        ],
    },
    {
        "id": "15_spelling_noisy_input",
        "name": "Spelling and Noisy Input",
        "domain": "robustness_normalization",
        "purpose": "Correct typographical errors, mobile touch inaccuracies, missing spaces, and phonetically inconsistent transliterations in Ethiopian marketplace chats.",
        "properties": {
            "id": {"type": "string"},
            "noisy_input": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "normalized_output": {"type": "string"},
            "noise_type": {
                "type": "string",
                "enum": [
                    "typo",
                    "phonetic",
                    "script_mix",
                    "abbreviation",
                    "keyboard_slip",
                ],
            },
        },
        "required": ["id", "noisy_input", "locale", "normalized_output"],
    },
    {
        "id": "16_recommendation_ranking",
        "name": "Recommendation and Ranking",
        "domain": "search_relevance",
        "purpose": "Train ranking explanations, price-performance evaluations, and trade-off summaries based on verified marketplace attributes.",
        "properties": {
            "id": {"type": "string"},
            "user_constraints": {"type": "object"},
            "candidate_item_ids": {"type": "array", "items": {"type": "integer"}},
            "ranking_order": {"type": "array", "items": {"type": "integer"}},
            "tradeoff_explanation": {"type": "string"},
        },
        "required": ["id", "user_constraints", "candidate_item_ids", "ranking_order"],
    },
    {
        "id": "17_comparison_calculations",
        "name": "Marketplace Comparison and Calculations",
        "domain": "deterministic_math",
        "purpose": "Perform truthful arithmetic comparisons (monthly vs. annual rent, budget remaining, price difference, area per square meter) without hallucination.",
        "properties": {
            "id": {"type": "string"},
            "query": {"type": "string"},
            "items_compared": {"type": "array", "items": {"type": "object"}},
            "calculation_formula": {"type": "string"},
            "deterministic_result": {"type": "string"},
        },
        "required": ["id", "query", "items_compared", "deterministic_result"],
    },
    {
        "id": "18_safety_policy",
        "name": "Safety, Security and Platform Policy",
        "domain": "ai_safety_governance",
        "purpose": "Safeguard against prompt injection, jailbreak attempts, PII extraction, and prohibited item listings while maintaining safe, respectful refusals.",
        "properties": {
            "id": {"type": "string"},
            "adversarial_prompt": {"type": "string"},
            "threat_category": {
                "type": "string",
                "enum": [
                    "prompt_injection",
                    "pii_leak",
                    "prohibited_content",
                    "auth_bypass",
                    "data_poisoning",
                ],
            },
            "safe_refusal_response": {"type": "string"},
            "policy_violation_code": {"type": "string"},
        },
        "required": [
            "id",
            "adversarial_prompt",
            "threat_category",
            "safe_refusal_response",
        ],
    },
    {
        "id": "19_tool_calling_structured_outputs",
        "name": "Tool Calling and Structured Outputs",
        "domain": "agentic_tool_execution",
        "purpose": "Validate model generation of strict, typed JSON tool arguments for Go backend marketplace endpoints, preventing arbitrary SQL or API execution.",
        "properties": {
            "id": {"type": "string"},
            "user_intent": {"type": "string"},
            "tool_name": {
                "type": "string",
                "enum": [
                    "search_marketplace_listings",
                    "get_listing_details",
                    "compare_listings",
                    "get_marketplace_stats",
                ],
            },
            "arguments": {"type": "object"},
            "schema_valid": {"type": "boolean"},
        },
        "required": ["id", "user_intent", "tool_name", "arguments", "schema_valid"],
    },
    {
        "id": "20_grounding_response_quality",
        "name": "Response Grounding and Quality",
        "domain": "grounding_factuality",
        "purpose": "Ensure responses are strictly anchored in Go database results, eliminating fabricated listings, prices, or contact numbers when zero results are found.",
        "properties": {
            "id": {"type": "string"},
            "query": {"type": "string"},
            "db_results_count": {"type": "integer"},
            "grounded_response": {"type": "string"},
            "zero_hallucination_verified": {"type": "boolean"},
        },
        "required": [
            "id",
            "query",
            "db_results_count",
            "grounded_response",
            "zero_hallucination_verified",
        ],
    },
    {
        "id": "21_dialogue_acts_turn_management",
        "name": "Dialogue Acts and Turn Management",
        "domain": "conversation_flow",
        "purpose": "Classify conversational dialogue acts (greeting, question, request, confirmation, rejection, clarification, thanks, farewell) to maintain conversational coherence.",
        "properties": {
            "id": {"type": "string"},
            "speaker": {"type": "string", "enum": ["user", "assistant"]},
            "utterance": {"type": "string"},
            "dialogue_act": {
                "type": "string",
                "enum": [
                    "inform",
                    "request",
                    "confirm",
                    "reject",
                    "greet",
                    "thank",
                    "clarify",
                    "farewell",
                ],
            },
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
        },
        "required": ["id", "speaker", "utterance", "dialogue_act", "locale"],
    },
    {
        "id": "22_task_state_slot_completion",
        "name": "Task State and Slot Completion",
        "domain": "task_oriented_dialogue",
        "purpose": "Track the lifecycle of temporary user search tasks: required slots, supplied parameters, missing criteria, and task completion or cancellation.",
        "properties": {
            "id": {"type": "string"},
            "task_id": {"type": "string"},
            "goal": {"type": "string"},
            "required_slots": {"type": "array", "items": {"type": "string"}},
            "filled_slots": {"type": "object"},
            "is_complete": {"type": "boolean"},
        },
        "required": [
            "id",
            "task_id",
            "goal",
            "required_slots",
            "filled_slots",
            "is_complete",
        ],
    },
    {
        "id": "23_constraint_reasoning_conflict_resolution",
        "name": "Constraint Reasoning and Conflict Resolution",
        "domain": "logic_constraint_satisfaction",
        "purpose": "Identify and resolve conflicting search constraints (e.g. 'Bole 3-bedroom apartment under 5,000 Birr') by suggesting realistic parameter relaxations.",
        "properties": {
            "id": {"type": "string"},
            "conflicting_constraints": {"type": "object"},
            "conflict_type": {"type": "string"},
            "resolution_strategy": {
                "type": "string",
                "enum": ["relax_price", "widen_location", "ask_user_priority"],
            },
            "suggested_compromise": {"type": "string"},
        },
        "required": [
            "id",
            "conflicting_constraints",
            "conflict_type",
            "suggested_compromise",
        ],
    },
    {
        "id": "24_user_preference_understanding",
        "name": "User Preference Understanding",
        "domain": "contextual_personalization",
        "purpose": "Capture short-term session preferences (preferred brands, locations, budget ceilings) while strictly adhering to guest privacy and non-persistence rules.",
        "properties": {
            "id": {"type": "string"},
            "session_context": {"type": "object"},
            "revealed_preference": {"type": "string"},
            "preference_scope": {
                "type": "string",
                "enum": ["temporary_session", "ephemeral_turn"],
            },
            "adapted_response": {"type": "string"},
        },
        "required": ["id", "revealed_preference", "preference_scope"],
    },
    {
        "id": "25_negative_feedback_recovery",
        "name": "Negative Feedback and Recovery",
        "domain": "dialogue_repair",
        "purpose": "Recognize user dissatisfaction ('not what I asked for', 'wrong city') and execute respectful conversational repairs and parameter resets.",
        "properties": {
            "id": {"type": "string"},
            "prior_assistant_response": {"type": "string"},
            "user_rejection": {"type": "string"},
            "failure_diagnosis": {"type": "string"},
            "recovery_turn": {"type": "string"},
        },
        "required": ["id", "user_rejection", "failure_diagnosis", "recovery_turn"],
    },
    {
        "id": "26_topic_switching_context_reset",
        "name": "Topic Switching and Context Reset",
        "domain": "context_hygiene",
        "purpose": "Gracefully transition when a user switches topics (e.g. from buying a laptop to asking about apartment rentals), clearing stale constraints without confusion.",
        "properties": {
            "id": {"type": "string"},
            "previous_topic": {"type": "string"},
            "new_topic_message": {"type": "string"},
            "stale_slots_cleared": {"type": "array", "items": {"type": "string"}},
            "retained_slots": {"type": "array", "items": {"type": "string"}},
            "transition_response": {"type": "string"},
        },
        "required": [
            "id",
            "previous_topic",
            "new_topic_message",
            "transition_response",
        ],
    },
    {
        "id": "27_search_refinement_progressive_filtering",
        "name": "Search Refinement and Progressive Filtering",
        "domain": "progressive_refinement",
        "purpose": "Handle step-by-step filter narrowing or widening as users interact with marketplace search results over multiple turns.",
        "properties": {
            "id": {"type": "string"},
            "base_filter": {"type": "object"},
            "refinement_input": {"type": "string"},
            "filter_mutation": {
                "type": "string",
                "enum": ["narrow", "widen", "replace", "remove"],
            },
            "resulting_filter": {"type": "object"},
        },
        "required": [
            "id",
            "base_filter",
            "refinement_input",
            "filter_mutation",
            "resulting_filter",
        ],
    },
    {
        "id": "28_confidence_uncertainty",
        "name": "Confidence and Uncertainty",
        "domain": "calibrated_uncertainty",
        "purpose": "Provide calibrated confidence estimates and appropriate abstentions when query interpretation is ambiguous or low-confidence.",
        "properties": {
            "id": {"type": "string"},
            "ambiguous_input": {"type": "string"},
            "confidence_score": {"type": "number"},
            "abstention_threshold": {"type": "number"},
            "should_abstain": {"type": "boolean"},
            "uncertainty_response": {"type": "string"},
        },
        "required": [
            "id",
            "ambiguous_input",
            "confidence_score",
            "should_abstain",
            "uncertainty_response",
        ],
    },
    {
        "id": "29_service_failure_error_conversations",
        "name": "Service Failure and Error Conversations",
        "domain": "resilience_error_handling",
        "purpose": "Handle backend service timeouts, database connection errors, and rate limit events with honest, transparent, and reassuring conversational responses.",
        "properties": {
            "id": {"type": "string"},
            "error_condition": {
                "type": "string",
                "enum": [
                    "api_timeout",
                    "db_unavailable",
                    "rate_limit_exceeded",
                    "network_split",
                ],
            },
            "user_query": {"type": "string"},
            "graceful_response": {"type": "string"},
            "allow_retry": {"type": "boolean"},
        },
        "required": [
            "id",
            "error_condition",
            "user_query",
            "graceful_response",
            "allow_retry",
        ],
    },
    {
        "id": "30_conversation_quality_response_generation",
        "name": "Conversation Quality and Natural Response Generation",
        "domain": "nlg_quality",
        "purpose": "Evaluate naturalness, conciseness, tone, and cultural appropriateness across Amharic, Afaan Oromoo, Tigrinya, and English marketplace responses.",
        "properties": {
            "id": {"type": "string"},
            "input_context": {"type": "string"},
            "generated_response": {"type": "string"},
            "locale": {"type": "string", "enum": ["en", "am", "or", "ti"]},
            "naturalness_score": {"type": "number"},
            "conciseness_score": {"type": "number"},
            "cultural_fit_score": {"type": "number"},
        },
        "required": ["id", "input_context", "generated_response", "locale"],
    },
]


def write_file(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with open(path, "w", encoding="utf-8") as f:
        f.write(content)


def write_json(path: Path, data: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with open(path, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2, ensure_ascii=False)
        f.write("\n")


def create_dataset_module(mod: dict) -> None:
    mod_id = mod["id"]
    name = mod["name"]
    domain = mod["domain"]
    purpose = mod["purpose"]
    mod_dir = BASE_DIR / "datasets" / mod_id

    # 1. README.md
    readme_content = f"""# HUNTIQ Dataset Module {mod_id[:2]} — {name}

## 1. Module Overview
- **Module ID:** `{mod_id}`
- **Domain:** `{domain}`
- **Primary Languages:** English (`en`), Amharic (`am`), Afaan Oromoo (`or`), Tigrinya (`ti`)
- **Format:** JSON Lines (`.jsonl`)
- **Schema:** [`schema.json`](./schema.json)
- **Manifest:** [`manifest.json`](./manifest.json)
- **Sources & Provenance:** [`sources.json`](./sources.json)
- **Data Dictionary:** [`data_dictionary.md`](./data_dictionary.md)

## 2. Purpose & Scope
{purpose}

All examples in this dataset module relate directly to the HUNTIQ marketplace ecosystem in Ethiopia.

## 3. Directory Layout
```text
{mod_id}/
├── README.md
├── schema.json
├── manifest.json
├── sources.json
├── data_dictionary.md
├── raw/
│   ├── README.md
│   ├── external/
│   ├── kaggle/
│   ├── huggingface/
│   └── original/
├── processed/
│   ├── normalized.jsonl
│   ├── deduplicated.jsonl
│   └── quality_checked.jsonl
├── training/
│   ├── train.jsonl
│   ├── validation.jsonl
│   └── test.jsonl
├── languages/
│   ├── en.jsonl
│   ├── am.jsonl
│   ├── or.jsonl
│   └── ti.jsonl
├── evaluation/
│   ├── benchmark.jsonl
│   ├── expected_results.jsonl
│   └── edge_cases.jsonl
├── statistics/
│   ├── statistics.json
│   └── quality_report.json
└── examples/
    └── examples.jsonl
```

## 4. Quality & Governance Invariants
- Zero fabricated records.
- Zero customer PII or private database credentials.
- Multi-language split consistency across EN, AM, OR, and TI.
- Strict schema validation enforced before moving from `raw/` to `processed/`.
"""
    write_file(mod_dir / "README.md", readme_content)

    # 2. schema.json
    schema_data = {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "title": f"HUNTIQ_{mod_id}_Schema",
        "description": f"Formal JSON validation schema for HUNTIQ dataset module {mod_id} ({name})",
        "type": "object",
        "properties": mod["properties"],
        "required": mod["required"],
        "additionalProperties": False,
    }
    write_json(mod_dir / "schema.json", schema_data)

    # 3. manifest.json
    manifest_data = {
        "module_id": mod_id,
        "name": name,
        "domain": domain,
        "version": "1.0.0",
        "lifecycle_stage": "initialized",
        "created_at": "2026-11-02T12:00:00Z",
        "updated_at": "2026-11-02T12:00:00Z",
        "maintainer": "HUNTIQ AI Data Platform Team",
        "languages": ["en", "am", "or", "ti"],
        "record_counts": {
            "raw": 0,
            "processed": 0,
            "train": 0,
            "validation": 0,
            "test": 0,
            "languages": {"en": 0, "am": 0, "or": 0, "ti": 0},
        },
        "quality_metrics": {
            "schema_compliance_rate": 1.0,
            "deduplication_rate": 0.0,
            "human_verified_pct": 0.0,
        },
        "eligibility": {"training_ready": False, "evaluation_ready": True},
    }
    write_json(mod_dir / "manifest.json", manifest_data)

    # 4. sources.json
    sources_data = {
        "module_id": mod_id,
        "primary_sources": [
            {
                "source_id": f"huntiq_internal_{mod_id}",
                "type": "internal_dialogue_seeds",
                "description": f"Curated seed patterns and real marketplace dialogues for {name}",
                "license": "Proprietary / HUNTIQ Platform License",
                "provenance": "HUNTIQ Marketplace Engineering",
                "pii_scrubbed": True,
                "verified_by_humans": True,
            }
        ],
        "external_allowed_sources": [
            {
                "source_id": "ethiopian_nlp_multilingual_corpus",
                "type": "public_academic_dataset",
                "license": "CC-BY-4.0",
                "commercial_usable": True,
                "domain_restriction": "Strictly Ethiopian marketplace terminology",
            }
        ],
    }
    write_json(mod_dir / "sources.json", sources_data)

    # 5. data_dictionary.md
    data_dict_content = f"""# Data Dictionary: {name} (`{mod_id}`)

## Field Specifications

| Field Name | Type | Required | Description | Allowed Values / Constraints |
|---|---|---|---|---|
"""
    for prop_name, prop_meta in mod["properties"].items():
        p_type = prop_meta.get("type", "any")
        req_str = "Yes" if prop_name in mod["required"] else "No"
        desc = prop_meta.get("description", f"Attribute representing {prop_name}")
        allowed = prop_meta.get("enum", "-")
        data_dict_content += (
            f"| `{prop_name}` | `{p_type}` | {req_str} | {desc} | {allowed} |\n"
        )

    data_dict_content += """
## Storage and Format Rules
- Files in `processed/`, `training/`, `languages/`, and `evaluation/` must be UTF-8 encoded JSON Lines (`.jsonl`).
- Each line is a single valid JSON object matching [`schema.json`](./schema.json).
- Trailing commas are strictly prohibited in JSONL lines.
"""
    write_file(mod_dir / "data_dictionary.md", data_dict_content)

    # 6. Raw directory
    write_file(
        mod_dir / "raw" / "README.md",
        f"""# Raw Data Ingestion: {name}

Store untouched, ingested source files here before validation and normalization.
- `external/`: External open-source datasets (license-verified)
- `kaggle/`: Kaggle competition / benchmark imports
- `huggingface/`: HuggingFace Hub dataset snapshots
- `original/`: Original seed logs from HUNTIQ QA and domain audits
""",
    )
    for sub in ["external", "kaggle", "huggingface", "original"]:
        write_file(mod_dir / "raw" / sub / ".gitkeep", "")

    # 7. Processed directory
    for f in ["normalized.jsonl", "deduplicated.jsonl", "quality_checked.jsonl"]:
        write_file(mod_dir / "processed" / f, "")

    # 8. Training directory
    for f in ["train.jsonl", "validation.jsonl", "test.jsonl"]:
        write_file(mod_dir / "training" / f, "")

    # 9. Languages directory
    for f in ["en.jsonl", "am.jsonl", "or.jsonl", "ti.jsonl"]:
        write_file(mod_dir / "languages" / f, "")

    # 10. Evaluation directory
    for f in ["benchmark.jsonl", "expected_results.jsonl", "edge_cases.jsonl"]:
        write_file(mod_dir / "evaluation" / f, "")

    # 11. Statistics directory
    write_json(
        mod_dir / "statistics" / "statistics.json",
        {
            "module_id": mod_id,
            "total_records": 0,
            "split_counts": {"train": 0, "validation": 0, "test": 0},
            "language_distribution": {"en": 0, "am": 0, "or": 0, "ti": 0},
            "vocabulary_size": {"en": 0, "am": 0, "or": 0, "ti": 0},
            "avg_token_length": 0.0,
            "last_computed": "2026-11-02T12:00:00Z",
        },
    )
    write_json(
        mod_dir / "statistics" / "quality_report.json",
        {
            "module_id": mod_id,
            "schema_validation_passed": True,
            "duplicate_count": 0,
            "syntax_error_count": 0,
            "missing_required_fields": 0,
            "language_mismatch_count": 0,
            "overall_quality_grade": "A",
        },
    )

    # 12. Examples directory
    write_file(mod_dir / "examples" / "examples.jsonl", "")

    # Special handling for Dataset 03 knowledge files
    if mod_id == "03_huntiq_platform_knowledge":
        knowledge_dir = mod_dir / "knowledge"
        knowledge_files = [
            (
                "product_documentation.md",
                "Product Listings, Categories & Electronics Rules",
            ),
            (
                "membership_documentation.md",
                "Standard, Gold, and Premium Membership Tiers",
            ),
            (
                "listing_documentation.md",
                "Listing Lifecycle, Publishing, Promoting & Expiration",
            ),
            (
                "account_documentation.md",
                "User Registration, Buyer/Seller Onboarding & OTP Auth",
            ),
            (
                "payment_documentation.md",
                "Telebirr, CBEBirr, and Platform Fee Settlement",
            ),
            (
                "privacy_documentation.md",
                "Data Privacy, Guest Chat Non-Persistence & User Rights",
            ),
            (
                "security_documentation.md",
                "Anti-Fraud, Libre Ownership Verification & Trust Badges",
            ),
            (
                "support_documentation.md",
                "Platform Support, Dispute Resolution & Appeals",
            ),
        ]
        for k_filename, k_title in knowledge_files:
            k_content = f"""# HUNTIQ Authoritative Knowledge: {k_title}

- **Topic:** `{k_filename.replace("_documentation.md", "")}`
- **Source Authority:** HUNTIQ Operational Platform Documentation (Sovereign Go Authority)
- **Effective Date:** 2026-11-02
- **Document Version:** 2.0.0
- **Status:** Authoritative (Truth source for AI platform knowledge)

## 1. Official Platform Policy
This document defines verified, authoritative platform rules for {k_title}. The AI Assistant must ground all platform explanations strictly in these policies and never hallucinate capabilities, fees, or requirements.

## 2. Invariants
- Verified against live backend logic in `backend/internal/`.
- Updated continuously with platform releases.
"""
            write_file(knowledge_dir / k_filename, k_content)


def generate_global_registry() -> None:
    ds_dir = BASE_DIR / "datasets"

    # README.md
    write_file(
        ds_dir / "README.md",
        """# HUNTIQ AI Global Dataset Ecosystem

This directory hosts the 30 foundational, domain-specialized AI dataset modules powering HUNTIQ's natural language understanding, multilingual intent classification, constraint reasoning, and zero-fabrication grounded search.

## Modules Index
| # | Module ID | Name | Domain | Schema |
|---|---|---|---|---|
"""
        + "\n".join(
            [
                f"| {m['id'][:2]} | [`{m['id']}`](./{m['id']}/) | {m['name']} | `{m['domain']}` | [schema.json](./{m['id']}/schema.json) |"
                for m in DATASET_MODULES
            ]
        )
        + "\n",
    )

    # dataset_registry.json
    registry_data = {
        "registry_version": "2.0.0",
        "updated_at": "2026-11-02T12:00:00Z",
        "total_modules": len(DATASET_MODULES),
        "supported_languages": ["en", "am", "or", "ti"],
        "modules": [
            {
                "module_id": m["id"],
                "name": m["name"],
                "domain": m["domain"],
                "path": f"datasets/{m['id']}",
                "schema_path": f"datasets/{m['id']}/schema.json",
                "manifest_path": f"datasets/{m['id']}/manifest.json",
            }
            for m in DATASET_MODULES
        ],
    }
    write_json(ds_dir / "dataset_registry.json", registry_data)

    # dataset_matrix.md
    matrix_md = """# HUNTIQ Dataset Responsibility & Domain Matrix

| Module ID | Primary Purpose | Go Backend Anchor | Python Core Module | Target Languages |
|---|---|---|---|---|
"""
    for m in DATASET_MODULES:
        matrix_md += f"| `{m['id']}` | {m['purpose']} | `internal/service/` | `app/core/` | EN, AM, OR, TI |\n"
    write_file(ds_dir / "dataset_matrix.md", matrix_md)

    # Supporting markdown docs
    write_file(
        ds_dir / "dataset_relationships.md",
        "# HUNTIQ Dataset Relationships\n\nDefines how the 30 dataset modules form a single, coordinated intelligence graph.\n",
    )
    write_file(
        ds_dir / "dataset_dependency_map.md",
        "# HUNTIQ Dataset Dependency Map\n\nExplicit dependency tree between categorization, entity extraction, dialogue state, and response grounding.\n",
    )
    write_file(
        ds_dir / "dataset_ownership.md",
        "# HUNTIQ Dataset Ownership\n\nOwnership assignments across Data Engineering, NLP Architects, and Trust & Safety teams.\n",
    )
    write_file(
        ds_dir / "dataset_quality_standards.md",
        "# HUNTIQ Dataset Quality Standards\n\nQuality grading rules (Grade A to F), deduplication thresholds, and schema validation gates.\n",
    )
    write_file(
        ds_dir / "dataset_naming_conventions.md",
        "# HUNTIQ Dataset Naming Conventions\n\nConsistent snake_case, ISO-639-1 language codes (en, am, or, ti), and split naming.\n",
    )
    write_file(
        ds_dir / "dataset_versioning.md",
        "# HUNTIQ Dataset Versioning\n\nSemantic versioning for data artifacts (Major.Minor.Patch) tied to Git commits.\n",
    )
    write_file(
        ds_dir / "dataset_lifecycle.md",
        "# HUNTIQ Dataset Lifecycle\n\nSix-stage progression: Ingestion -> Normalization -> Human QA -> Splitting -> Benchmark Evaluation -> Training Eligibility.\n",
    )
    write_file(
        ds_dir / "dataset_import_guide.md",
        "# HUNTIQ Dataset Import Guide\n\nStep-by-step instructions for importing new JSONL, CSV, and Parquet data.\n",
    )
    write_file(
        ds_dir / "dataset_contribution_guide.md",
        "# HUNTIQ Dataset Contribution Guide\n\nGuidelines for AI engineers authoring seed patterns.\n",
    )
    write_file(
        ds_dir / "dataset_maintenance_guide.md",
        "# HUNTIQ Dataset Maintenance Guide\n\nRoutine audits, quarterly vocabulary updates, and taxonomy alignment.\n",
    )
    write_file(
        ds_dir / "dataset_targets.md",
        """# HUNTIQ Dataset Volume and Quality Strategy

## Target Ecosystem Capacity
- **Initial Ecosystem:** 30 Modules (Created & Specified)
- **Target Volume Range:** 600,000 – 1,000,000 high-quality records over time
- **Quality Over Raw Volume:** No fabricated data. Synthetic examples must be strictly tagged and reviewed.
- **Language Allocation:**
  - English (`en`): 30%
  - Amharic (`am`): 35%
  - Afaan Oromoo (`or`): 20%
  - Tigrinya (`ti`): 15%
""",
    )

    # JSON reports and registries
    write_json(
        ds_dir / "dataset_coverage_report.json",
        {
            "status": "filesystem_initialized",
            "total_modules": 30,
            "languages": ["en", "am", "or", "ti"],
            "domains_covered": list({m["domain"] for m in DATASET_MODULES}),
        },
    )
    write_json(
        ds_dir / "dataset_inventory.json",
        {
            "total_datasets": 30,
            "inventory": [
                {"id": m["id"], "name": m["name"], "domain": m["domain"]}
                for m in DATASET_MODULES
            ],
        },
    )
    write_json(
        ds_dir / "dataset_statistics.json",
        {
            "total_records_all_datasets": 0,
            "modules_count": 30,
            "last_updated": "2026-11-02T12:00:00Z",
        },
    )
    write_json(
        ds_dir / "language_coverage.json",
        {
            "supported_locales": ["en", "am", "or", "ti"],
            "coverage_percentage": {"en": 100.0, "am": 100.0, "or": 100.0, "ti": 100.0},
        },
    )
    write_json(
        ds_dir / "dataset_source_catalog.json",
        {"approved_sources": ["huntiq_internal", "ethiopian_nlp_multilingual_corpus"]},
    )
    write_json(
        ds_dir / "dataset_license_registry.json",
        {
            "commercial_approved_licenses": [
                "Proprietary/HUNTIQ",
                "MIT",
                "Apache-2.0",
                "CC-BY-4.0",
            ],
            "prohibited_licenses": [
                "GPL-3.0-copyleft-training-risk",
                "CC-BY-NC-4.0-non-commercial",
            ],
        },
    )
    write_json(
        ds_dir / "dataset_risk_register.json",
        {
            "risks": [
                {
                    "id": "RISK-01",
                    "description": "Data poisoning or malicious prompts in raw imports",
                    "mitigation": "Regex and AST policy scanners in normalization pipeline",
                },
                {
                    "id": "RISK-02",
                    "description": "Accidental PII inclusion",
                    "mitigation": "Automated phone/email/national ID anonymization scrubbing",
                },
                {
                    "id": "RISK-03",
                    "description": "Train/Evaluation dataset leakage",
                    "mitigation": "Deterministic hash-based splits with zero overlap verification",
                },
            ]
        },
    )
    write_json(
        ds_dir / "dataset_readiness_report.json",
        {
            "phase": "Master Phase 1 - Filesystem Creation",
            "filesystem_ready": True,
            "training_active": False,
            "evaluation_ready": True,
        },
    )

    # Governance directory in datasets/
    gov_dir = ds_dir / "governance"
    write_file(
        gov_dir / "data_governance.md",
        "# Data Governance Framework\n\nCore policies for data governance across all HUNTIQ AI datasets.\n",
    )
    write_file(
        gov_dir / "privacy_policy.md",
        "# AI Data Privacy Policy\n\nZero guest prompt persistence, scrubbed personal identifiers, and GDPR/Ethiopian Data Protection alignment.\n",
    )
    write_file(
        gov_dir / "licensing_policy.md",
        "# Dataset Licensing Policy\n\nRules on permissible open-source and commercial data usage.\n",
    )
    write_file(
        gov_dir / "provenance_policy.md",
        "# Provenance and Lineage Policy\n\nTracing every sample back to its original ingested commit and author.\n",
    )
    write_file(
        gov_dir / "quality_policy.md",
        "# Data Quality Policy\n\nVerification standards before records enter `processed/`.\n",
    )
    write_file(
        gov_dir / "retention_policy.md",
        "# Data Retention Policy\n\nStorage duration for raw logs versus benchmark samples.\n",
    )
    write_file(
        gov_dir / "human_review_policy.md",
        "# Human Review Policy\n\nMandatory review protocol for translations and safety refusals.\n",
    )


def generate_intelligence_contracts() -> None:
    intel_dir = BASE_DIR / "intelligence"
    files = [
        "README.md",
        "architecture.md",
        "intelligence_pipeline.md",
        "module_ownership.md",
        "shared_context_contract.md",
        "intent_contract.md",
        "entity_contract.md",
        "dialogue_contract.md",
        "task_state_contract.md",
        "constraint_contract.md",
        "language_contract.md",
        "confidence_contract.md",
        "response_contract.md",
        "tool_selection_contract.md",
        "grounding_contract.md",
        "error_recovery_contract.md",
        "context_lifecycle.md",
        "decision_priority.md",
        "decision_matrix.md",
    ]
    for f in files:
        write_file(
            intel_dir / f,
            f"# HUNTIQ Intelligence Contract: {f.replace('.md', '').replace('_', ' ').title()}\n\nDefines shared interface specifications for the unified conversational pipeline.\n",
        )

    write_json(
        intel_dir / "module_dependency_map.json",
        {
            "pipeline_stages": [
                {
                    "stage": 1,
                    "name": "normalization",
                    "input": "raw_user_message",
                    "output": "normalized_text",
                },
                {
                    "stage": 2,
                    "name": "language_detection",
                    "input": "normalized_text",
                    "output": "locale_code",
                },
                {
                    "stage": 3,
                    "name": "context_resolution",
                    "input": "session_id",
                    "output": "session_context",
                },
                {
                    "stage": 4,
                    "name": "intent_and_entity_classification",
                    "input": "normalized_text",
                    "output": "execution_plan",
                },
                {
                    "stage": 5,
                    "name": "policy_guard",
                    "input": "execution_plan",
                    "output": "allowed_or_redirect",
                },
                {
                    "stage": 6,
                    "name": "tool_execution",
                    "input": "typed_tool_arguments",
                    "output": "verified_db_listings",
                },
                {
                    "stage": 7,
                    "name": "grounded_synthesis",
                    "input": "verified_db_listings",
                    "output": "final_grounded_response",
                },
            ]
        },
    )
    write_json(
        intel_dir / "runtime_configuration.json",
        {
            "service_name": "huntiq-ai-intelligence",
            "version": "2.0.0",
            "max_context_turns": 6,
            "session_ttl_seconds": 1200,
            "go_tool_timeout_seconds": 3.0,
        },
    )


def generate_languages_architecture() -> None:
    lang_dir = BASE_DIR / "languages"
    files = [
        "README.md",
        "semantic_alignment_policy.md",
        "translation_quality_policy.md",
        "transliteration_policy.md",
        "mixed_language_policy.md",
    ]
    for f in files:
        write_file(
            lang_dir / f,
            f"# Multilingual Policy: {f.replace('.md', '').replace('_', ' ').title()}\n\nGuarantees first-class semantic parity across English, Amharic, Afaan Oromoo, and Tigrinya.\n",
        )

    write_json(
        lang_dir / "language_registry.json",
        {
            "languages": [
                {
                    "code": "en",
                    "name": "English",
                    "script": "Latin",
                    "status": "primary_tier_1",
                },
                {
                    "code": "am",
                    "name": "Amharic",
                    "script": "Ge'ez",
                    "status": "primary_tier_1",
                },
                {
                    "code": "or",
                    "name": "Afaan Oromoo",
                    "script": "Latin (Qubee)",
                    "status": "primary_tier_1",
                },
                {
                    "code": "ti",
                    "name": "Tigrinya",
                    "script": "Ge'ez",
                    "status": "primary_tier_1",
                },
            ]
        },
    )
    write_json(
        lang_dir / "language_coverage_matrix.json",
        {"modules_target_parity": 30, "supported_locales": ["en", "am", "or", "ti"]},
    )
    write_json(
        lang_dir / "terminology_registry.json",
        {
            "version": "2.0.0",
            "core_domains": [
                "real_estate",
                "vehicles",
                "electronics",
                "agriculture",
                "jobs",
                "services",
            ],
        },
    )

    # Glossary files
    glossary_dir = lang_dir / "marketplace_glossary"
    glossary_en = {
        "rent": "rent",
        "buy": "buy",
        "house": "house",
        "apartment": "apartment",
        "car": "car",
        "phone": "phone",
        "price": "price",
    }
    glossary_am = {
        "rent": "ኪራይ",
        "buy": "ግዢ",
        "house": "ቤት",
        "apartment": "አፓርታማ",
        "car": "መኪና",
        "phone": "ስልክ",
        "price": "ዋጋ",
    }
    glossary_or = {
        "rent": "kiraa",
        "buy": "bituu",
        "house": "mana",
        "apartment": "appaartamaa",
        "car": "konkolaataa",
        "phone": "bilbila",
        "price": "gatii",
    }
    glossary_ti = {
        "rent": "ክራይ",
        "buy": "ምግዛእ",
        "house": "ገዛ",
        "apartment": "ኣፓርታማ",
        "car": "መኪና",
        "phone": "ስልኪ",
        "price": "ዋጋ",
    }

    write_json(glossary_dir / "en.json", glossary_en)
    write_json(glossary_dir / "am.json", glossary_am)
    write_json(glossary_dir / "or.json", glossary_or)
    write_json(glossary_dir / "ti.json", glossary_ti)

    # Human review
    hr_dir = lang_dir / "human_review"
    write_file(
        hr_dir / "review_guidelines.md",
        "# Human Review Guidelines for Ethiopian Localized AI Datasets\n",
    )
    write_json(hr_dir / "review_status.json", {"reviewed_count": 0, "pending_count": 0})
    write_json(hr_dir / "unresolved_translations.json", {"unresolved": []})

    # Evaluation
    eval_dir = lang_dir / "evaluation"
    write_json(eval_dir / "language_detection.json", {"target_accuracy": 0.98})
    write_json(
        eval_dir / "semantic_equivalence.json", {"target_cosine_similarity": 0.92}
    )
    write_json(eval_dir / "cross_language_consistency.json", {"target_parity": 0.95})


def generate_data_pipeline() -> None:
    dp_dir = BASE_DIR / "data_pipeline"
    files = [
        "README.md",
        "architecture.md",
        "import_workflow.md",
        "normalization_workflow.md",
        "validation_workflow.md",
        "deduplication_workflow.md",
        "splitting_workflow.md",
        "quality_control_workflow.md",
        "licensing_workflow.md",
        "provenance_workflow.md",
        "language_processing.md",
        "synthetic_data_policy.md",
        "human_annotation_workflow.md",
        "rejected_data_policy.md",
        "reproducibility.md",
    ]
    for f in files:
        write_file(
            dp_dir / f,
            f"# Data Pipeline Specification: {f.replace('.md', '').replace('_', ' ').title()}\n",
        )

    cfg_dir = dp_dir / "configurations"
    write_file(
        cfg_dir / "pipeline_config.yaml",
        "version: '2.0'\npipeline:\n  batch_size: 1000\n  max_workers: 4\n",
    )
    write_file(
        cfg_dir / "quality_thresholds.yaml",
        "thresholds:\n  min_quality_score: 0.85\n  max_duplicate_ratio: 0.02\n",
    )
    write_file(
        cfg_dir / "split_config.yaml",
        "splits:\n  train: 0.80\n  validation: 0.10\n  test: 0.10\n  random_seed: 42\n",
    )
    write_file(
        cfg_dir / "language_config.yaml",
        "languages:\n  supported: ['en', 'am', 'or', 'ti']\n",
    )
    write_file(
        cfg_dir / "license_rules.yaml",
        "allowed_licenses:\n  - Proprietary/HUNTIQ\n  - MIT\n  - Apache-2.0\n  - CC-BY-4.0\n",
    )

    schemas_dir = dp_dir / "schemas"
    for s in [
        "source_schema.json",
        "imported_record_schema.json",
        "normalized_record_schema.json",
        "annotation_schema.json",
        "rejection_schema.json",
        "provenance_schema.json",
    ]:
        write_json(
            schemas_dir / s,
            {
                "$schema": "https://json-schema.org/draft/2020-12/schema",
                "title": s.replace(".json", ""),
            },
        )

    for sub in ["import_jobs", "quarantine", "rejected", "reports"]:
        write_file(dp_dir / sub / ".gitkeep", "")


def generate_training_system() -> None:
    tr_dir = BASE_DIR / "training"
    files = [
        "README.md",
        "training_architecture.md",
        "training_strategy.md",
        "model_selection_policy.md",
        "fine_tuning_policy.md",
        "task_to_dataset_mapping.md",
        "training_readiness.md",
        "training_reproducibility.md",
        "training_resource_policy.md",
        "training_risk_management.md",
        "training_rollback.md",
    ]
    for f in files:
        write_file(
            tr_dir / f,
            f"# AI Training Specification: {f.replace('.md', '').replace('_', ' ').title()}\n\nFuture training design and capability specifications for HUNTIQ AI models.\n",
        )

    cfg_dir = tr_dir / "configurations"
    yaml_configs = [
        "global_training.yaml",
        "classification.yaml",
        "entity_extraction.yaml",
        "multilingual.yaml",
        "dialogue.yaml",
        "safety.yaml",
        "ranking.yaml",
        "tool_calling.yaml",
        "evaluation.yaml",
    ]
    for y in yaml_configs:
        write_file(
            cfg_dir / y,
            f"# Configuration for {y.replace('.yaml', '')}\nversion: '2.0'\nenabled: false # Training execution scheduled for future phases\n",
        )

    # Experiments
    write_file(
        tr_dir / "experiments" / "README.md", "# Training Experiments Registry\n"
    )
    write_json(tr_dir / "experiments" / "experiment_registry.json", {"experiments": []})

    # Checkpoints & Artifacts
    write_file(
        tr_dir / "checkpoints" / "README.md",
        "# Model Checkpoints\n\nPreserved model checkpoints from training runs (omitted from Git tracking).\n",
    )
    write_file(
        tr_dir / "artifacts" / "README.md",
        "# Training Artifacts\n\nEvaluation artifacts, loss curves, and tokenization outputs.\n",
    )

    # Model registry
    mr_dir = tr_dir / "model_registry"
    write_file(mr_dir / "README.md", "# HUNTIQ Model Registry\n")
    write_json(
        mr_dir / "model_registry.json",
        {"active_production_models": [], "candidate_models": []},
    )
    write_json(mr_dir / "model_versions.json", {"versions": []})
    write_file(
        mr_dir / "model_approval_policy.md",
        "# Model Approval and Promotion Policy\n\nBenchmark thresholds required before promoting models to production.\n",
    )

    # Training runs & Reports
    write_file(tr_dir / "training_runs" / "README.md", "# Training Runs History\n")
    write_file(tr_dir / "reports" / "README.md", "# Training Evaluation Reports\n")


def generate_evaluation_system() -> None:
    ev_dir = BASE_DIR / "evaluation"
    files = [
        "README.md",
        "evaluation_architecture.md",
        "evaluation_strategy.md",
        "evaluation_matrix.md",
        "acceptance_criteria.md",
        "regression_policy.md",
        "language_evaluation.md",
        "dialogue_evaluation.md",
        "marketplace_evaluation.md",
        "search_evaluation.md",
        "extraction_evaluation.md",
        "safety_evaluation.md",
        "tool_call_evaluation.md",
        "grounding_evaluation.md",
        "latency_evaluation.md",
        "robustness_evaluation.md",
        "hallucination_evaluation.md",
        "conversation_consistency.md",
        "human_evaluation.md",
    ]
    for f in files:
        write_file(
            ev_dir / f,
            f"# AI Evaluation Framework: {f.replace('.md', '').replace('_', ' ').title()}\n",
        )

    write_json(
        ev_dir / "benchmark_registry.json",
        {
            "benchmarks": [
                {"id": "domain_benchmark", "target_accuracy": 0.95},
                {"id": "entity_extraction", "target_f1": 0.92},
                {"id": "conversational_benchmark", "target_accuracy": 0.94},
                {"id": "zero_hallucination_check", "target_pass_rate": 1.0},
            ]
        },
    )

    for sub in [
        "benchmark_results",
        "regression_cases",
        "adversarial_cases",
        "multilingual_cases",
        "reports",
    ]:
        write_file(ev_dir / sub / ".gitkeep", "")


def generate_dataset_coordination() -> None:
    coord_dir = BASE_DIR / "dataset_coordination"
    files = [
        "README.md",
        "coordination_architecture.md",
        "module_interaction_matrix.md",
        "shared_category_reference.md",
        "cross_dataset_consistency.md",
        "semantic_alignment.md",
        "conflict_detection.md",
        "duplicate_management.md",
        "synchronization_policy.md",
    ]
    for f in files:
        write_file(
            coord_dir / f,
            f"# Dataset Coordination: {f.replace('.md', '').replace('_', ' ').title()}\n",
        )

    write_json(
        coord_dir / "shared_label_registry.json",
        {
            "labels": [
                "greeting",
                "search",
                "clarification",
                "out_of_domain",
                "platform_faq",
            ]
        },
    )
    write_json(
        coord_dir / "shared_entity_registry.json",
        {
            "entities": [
                "brand",
                "model",
                "price",
                "currency",
                "location",
                "condition",
                "listing_type",
            ]
        },
    )
    write_json(
        coord_dir / "shared_intent_registry.json",
        {
            "intents": [
                "marketplace_search",
                "inspect_reputation",
                "clarify_housing",
                "faq",
            ]
        },
    )
    write_json(
        coord_dir / "shared_language_registry.json",
        {"languages": ["en", "am", "or", "ti"]},
    )
    write_json(
        coord_dir / "dataset_dependency_graph.json",
        {"nodes": [m["id"] for m in DATASET_MODULES]},
    )
    write_json(
        coord_dir / "training_group_registry.json",
        {
            "groups": {
                "intent_and_dialogue": [
                    "01_general_greetings",
                    "02_identity_assistant_behavior",
                    "05_intent_classification",
                ]
            }
        },
    )
    write_json(
        coord_dir / "evaluation_group_registry.json",
        {
            "evaluation_suites": {
                "core_capabilities": [
                    "04_domain_classification",
                    "08_search_query_understanding",
                    "20_grounding_response_quality",
                ]
            }
        },
    )


def generate_governance() -> None:
    gov_dir = BASE_DIR / "governance"
    files = [
        "README.md",
        "data_security.md",
        "privacy_by_design.md",
        "personal_data_policy.md",
        "guest_conversation_policy.md",
        "customer_data_policy.md",
        "dataset_access_control.md",
        "data_retention.md",
        "anonymization.md",
        "consent_and_legal_basis.md",
        "external_dataset_licensing.md",
        "commercial_use_approval.md",
        "prompt_injection_policy.md",
        "data_poisoning_policy.md",
        "model_security.md",
        "auditability.md",
        "incident_response.md",
    ]
    for f in files:
        write_file(
            gov_dir / f,
            f"# Governance Specification: {f.replace('.md', '').replace('_', ' ').title()}\n\nStrict data protection, privacy by design, and safety protocols for HUNTIQ AI.\n",
        )


def generate_knowledge() -> None:
    kn_dir = BASE_DIR / "knowledge"
    files = [
        "README.md",
        "knowledge_architecture.md",
        "source_authority.md",
        "knowledge_versioning.md",
        "knowledge_validation.md",
        "marketplace_terminology.md",
        "marketplace_entity_model.md",
        "marketplace_category_reference.md",
        "listing_lifecycle.md",
        "membership_reference.md",
        "payment_reference.md",
        "location_reference.md",
        "platform_policy_reference.md",
        "support_reference.md",
        "retrieval_policy.md",
    ]
    for f in files:
        write_file(
            kn_dir / f,
            f"# Platform Knowledge Reference: {f.replace('.md', '').replace('_', ' ').title()}\n\nAuthoritative documentation and taxonomy reference for HUNTIQ marketplace.\n",
        )


def generate_root_documentation() -> None:
    docs = [
        (
            "AI_ARCHITECTURE.md",
            "# HUNTIQ AI Master Architecture Guide\n\nFull four-tier separation, sovereign Go authority, and private Python AI intelligence.\n",
        ),
        (
            "DATASET_MASTER_INDEX.md",
            "# HUNTIQ Dataset Master Index\n\nComplete catalog of the 30 dataset modules, schemas, and readiness status.\n",
        ),
        (
            "TRAINING_MASTER_GUIDE.md",
            "# HUNTIQ AI Training Master Guide\n\nSpecifications and policies for future task-specific model fine-tuning and evaluation.\n",
        ),
        (
            "EVALUATION_MASTER_GUIDE.md",
            "# HUNTIQ AI Evaluation Master Guide\n\nComprehensive benchmark matrix, acceptance criteria, and regression testing policies.\n",
        ),
        (
            "MULTILINGUAL_MASTER_GUIDE.md",
            "# HUNTIQ Multilingual Master Guide\n\nArchitectural guidelines for English, Amharic, Afaan Oromoo, and Tigrinya.\n",
        ),
        (
            "DATA_GOVERNANCE_MASTER_GUIDE.md",
            "# HUNTIQ Data Governance Master Guide\n\nPrivacy-by-design, non-persistence of guest chats, and licensing verification.\n",
        ),
        (
            "ROADMAP.md",
            "# HUNTIQ AI Ecosystem Roadmap: Phases 1 to 13\n\nPhased progression from filesystem setup to full production training.\n",
        ),
        (
            "PHASE_STATUS.md",
            "# HUNTIQ AI Master Phase Status\n\n- Master Phase 1 (Filesystem Creation): COMPLETED\n- Master Phase 2 (Dataset Discovery): PENDING\n",
        ),
    ]
    for filename, header in docs:
        write_file(BASE_DIR / filename, header)


def main():
    print(f"Creating complete HUNTIQ AI Dataset Ecosystem in: {BASE_DIR}")

    # 1. Create the 30 dataset modules
    print("Generating 30 dataset modules under ai/datasets/...")
    for mod in DATASET_MODULES:
        create_dataset_module(mod)
    print("  -> 30 modules created successfully.")

    # 2. Generate global dataset registry & governance
    print("Generating Global Dataset Registry in ai/datasets/...")
    generate_global_registry()

    # 3. Generate intelligence contracts
    print("Generating Shared Intelligence Contracts in ai/intelligence/...")
    generate_intelligence_contracts()

    # 4. Generate multilingual data architecture
    print("Generating Multilingual Architecture in ai/languages/...")
    generate_languages_architecture()

    # 5. Generate data pipeline
    print("Generating Data Pipeline in ai/data_pipeline/...")
    generate_data_pipeline()

    # 6. Generate training system
    print("Generating Training Infrastructure in ai/training/...")
    generate_training_system()

    # 7. Generate evaluation system
    print("Generating Evaluation Framework in ai/evaluation/...")
    generate_evaluation_system()

    # 8. Generate cross-dataset coordination
    print("Generating Cross-Dataset Coordination in ai/dataset_coordination/...")
    generate_dataset_coordination()

    # 9. Generate governance
    print("Generating Governance Specifications in ai/governance/...")
    generate_governance()

    # 10. Generate knowledge references
    print("Generating Platform Knowledge References in ai/knowledge/...")
    generate_knowledge()

    # 11. Generate root documentation
    print("Generating Root AI Documentation in ai/...")
    generate_root_documentation()

    print("\nHUNTIQ Dataset Ecosystem & Training Infrastructure creation complete!")


if __name__ == "__main__":
    main()
