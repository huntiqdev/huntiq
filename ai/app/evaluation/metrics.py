"""Evaluation metrics calculation for domain classification and tool extraction."""

from typing import Any


def calculate_precision_recall_f1(tp: int, fp: int, fn: int) -> dict[str, float]:
    """Calculate precision, recall, and F1 score with division-by-zero protection."""
    precision = tp / (tp + fp) if (tp + fp) > 0 else 0.0
    recall = tp / (tp + fn) if (tp + fn) > 0 else 0.0
    f1 = (
        (2 * precision * recall) / (precision + recall)
        if (precision + recall) > 0
        else 0.0
    )
    return {
        "precision": round(precision, 4),
        "recall": round(recall, 4),
        "f1": round(f1, 4),
    }


def validate_tool_call_schema(tool_call: dict[str, Any]) -> tuple[bool, str]:
    """Validate that a tool call structure matches HUNTIQ marketplace expectations."""
    func = tool_call.get("function")
    if not func or not isinstance(func, dict):
        return False, "Missing or invalid 'function' object in tool call"

    name = func.get("name")
    if not name or not isinstance(name, str):
        return False, "Missing or invalid 'name' in tool function"

    args = func.get("arguments")
    if not isinstance(args, dict):
        return False, "'arguments' must be a valid dictionary"

    if name == "search_marketplace_listings":
        if "searches" not in args or not isinstance(args["searches"], list):
            return False, "'search_marketplace_listings' requires 'searches' list"
        for s in args["searches"]:
            if not isinstance(s, dict) or "search_term" not in s:
                return False, "Each search item must specify 'search_term'"

    return True, "Valid"
