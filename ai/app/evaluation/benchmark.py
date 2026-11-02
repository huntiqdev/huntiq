"""Evaluation suite for assessing conversational, coding, multilingual, and marketplace accuracy."""

import asyncio
from typing import Any

import httpx

from ..model_client import OllamaClient

TEST_CASES = [
    {
        "id": "general_explanation",
        "category": "general",
        "locale": "en",
        "prompt": "Explain what an internal combustion engine is and how its four strokes work in simple terms.",
        "expected_keywords": ["stroke", "intake", "compression", "power", "exhaust"],
    },
    {
        "id": "coding_python",
        "category": "coding",
        "locale": "en",
        "prompt": "Write a clean Python function to reverse a string in-place without using slicing [::-1].",
        "expected_keywords": ["def", "while", "left", "right"],
    },
    {
        "id": "amharic_greeting",
        "category": "multilingual",
        "locale": "am",
        "prompt": "ሰላም እንደምን አደርክ? ዛሬ የቦሽ ድሪል መግዛት እፈልጋለሁ።",
        "expected_keywords": ["ሰላም", "ቦሽ", "ድሪል"],
    },
    {
        "id": "marketplace_tool_intent",
        "category": "marketplace",
        "locale": "en",
        "prompt": "I am looking to buy a used Toyota Vitz under 700,000 ETB in Addis Ababa.",
        "expect_tool": "search_marketplace_listings",
    },
]


async def run_evaluation(model_name: str = "qwen3:4b") -> dict[str, Any]:
    """Run evaluation cases against local Ollama model."""
    client = OllamaClient()
    healthy = await client.check_health()
    if not healthy:
        return {
            "status": "error",
            "message": "Local Ollama service is unreachable on 127.0.0.1:11434",
            "results": [],
        }

    results: list[dict[str, Any]] = []
    for test in TEST_CASES:
        messages = [{"role": "user", "content": test["prompt"]}]
        try:
            res = await client.chat(
                messages=messages,
                model=model_name,
                locale=test.get("locale", "en"),
            )
            message = res.get("message", {})
            content = message.get("content", "")
            tool_calls = message.get("tool_calls", [])

            passed = True
            details = []
            if "expected_keywords" in test:
                for kw in test["expected_keywords"]:
                    if kw.lower() not in content.lower():
                        passed = False
                        details.append(f"Missing keyword: {kw}")

            if "expect_tool" in test:
                tool_names = [
                    tc.get("function", {}).get("name")
                    for tc in tool_calls
                    if "function" in tc
                ]
                if test["expect_tool"] not in tool_names:
                    passed = False
                    details.append(
                        f"Expected tool call {test['expect_tool']}, got {tool_names}"
                    )

            results.append(
                {
                    "test_id": test["id"],
                    "category": test["category"],
                    "passed": passed,
                    "content_preview": content[:120],
                    "details": details,
                }
            )
        except (httpx.HTTPError, OSError, ValueError, KeyError, RuntimeError) as e:
            results.append(
                {
                    "test_id": test["id"],
                    "category": test["category"],
                    "passed": False,
                    "error": str(e),
                }
            )

    total = len(results)
    passed_count = sum(1 for r in results if r.get("passed"))
    return {
        "status": "completed",
        "model": model_name,
        "total": total,
        "passed": passed_count,
        "score_percent": round((passed_count / total) * 100, 2) if total > 0 else 0,
        "results": results,
    }


if __name__ == "__main__":
    summary = asyncio.run(run_evaluation())
    print(summary)
