"""Tests for evaluation metrics and tool call schema validation."""

from app.evaluation.metrics import (
    calculate_precision_recall_f1,
    validate_tool_call_schema,
)


def test_precision_recall_f1():
    scores = calculate_precision_recall_f1(tp=10, fp=2, fn=3)
    assert scores["precision"] == round(10 / 12, 4)
    assert scores["recall"] == round(10 / 13, 4)
    assert scores["f1"] > 0.75


def test_precision_recall_zero_division():
    scores = calculate_precision_recall_f1(tp=0, fp=0, fn=0)
    assert scores["precision"] == 0.0
    assert scores["recall"] == 0.0
    assert scores["f1"] == 0.0


def test_validate_tool_call_schema_valid():
    tool_call = {
        "function": {
            "name": "search_marketplace_listings",
            "arguments": {
                "searches": [{"search_term": "phone", "max_price": 5000}],
                "requires_clarification": False,
            },
        }
    }
    valid, reason = validate_tool_call_schema(tool_call)
    assert valid is True
    assert reason == "Valid"


def test_validate_tool_call_schema_invalid():
    invalid_call = {
        "function": {"name": "search_marketplace_listings", "arguments": {}}
    }
    valid, reason = validate_tool_call_schema(invalid_call)
    assert valid is False
    assert "requires 'searches'" in reason
