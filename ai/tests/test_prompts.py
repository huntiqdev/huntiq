"""Unit tests for prompts, multilingual support, and tool schemas."""

from ai.app.prompts import (
    BASE_SYSTEM_INSTRUCTION,
    LANGUAGE_GUIDELINES,
    MARKETPLACE_TOOLS,
    build_system_prompt,
)


def test_build_system_prompt_all_languages():
    for lang in ["en", "am", "om", "ti"]:
        prompt = build_system_prompt(locale=lang)
        assert BASE_SYSTEM_INSTRUCTION in prompt
        assert LANGUAGE_GUIDELINES[lang] in prompt


def test_marketplace_tools_schema_validity():
    assert len(MARKETPLACE_TOOLS) >= 2
    tool_names = [t["function"]["name"] for t in MARKETPLACE_TOOLS]
    assert "search_marketplace_listings" in tool_names
    assert "inspect_seller_reputation" in tool_names

    # Check search tool parameters
    search_tool = next(
        t
        for t in MARKETPLACE_TOOLS
        if t["function"]["name"] == "search_marketplace_listings"
    )
    props = search_tool["function"]["parameters"]["properties"]
    assert "query" in props
    assert "min_price" in props
    assert "location" in props
    assert "query" in search_tool["function"]["parameters"]["required"]
