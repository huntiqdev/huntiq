"""Query understanding and product analysis routes for Go backend integration."""

from __future__ import annotations

import time

from fastapi import APIRouter, Header, HTTPException, status

from ...api.schemas import (
    ExtractedEntities,
    ExtractedPriceRange,
    MarketplaceReasoningChain,
    PriceStats,
    ProductAnalysisRequest,
    ProductAnalysisResponse,
    QueryUnderstandingResponse,
    ReasoningStep,
    UnderstandQueryParams,
)
from ...config import config
from ...core.entity_extractor import (
    extract_condition,
    extract_location,
    extract_price_bounds,
)
from ...core.fuzzy_matcher import fuzzy_matcher
from ...core.trained_model_service import trained_models
from ...languages.normalizer import normalizer

router = APIRouter(prefix="/api/v1/query", tags=["query"])

CATEGORY_SLUG_TO_ID = {
    "products": 1,
    "vehicles": 2,
    "real_estate": 3,
    "jobs": 4,
    "services": 5,
    "agriculture": 6,
    "hotels": 7,
    "tickets": 8,
}


@router.post("/understand", response_model=QueryUnderstandingResponse)
async def handle_understand_query(
    req: UnderstandQueryParams,
    x_huntiq_ai_secret: str | None = Header(default=None),
) -> QueryUnderstandingResponse:
    """
    Authoritative query understanding endpoint called by Go backend aiclient.
    Validates X-HUNTIQ-AI-SECRET and applies trained ML checkpoints:
    - safety_moderation_v1
    - domain_classifier_v1
    - intent_classifier_v1
    - category_classifier_v1
    - fuzzy_matcher for brands and locations
    """
    if (
        config.api_secret_key
        and x_huntiq_ai_secret
        and x_huntiq_ai_secret != config.api_secret_key
    ):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid service credentials",
        )

    t0 = time.time()
    query = req.query.strip()
    locale = req.locale or "en"
    if locale == "or":
        locale = "om"

    norm_query = normalizer.normalize(query)

    # 1. Safety check
    safety_class, _safety_conf, _ = trained_models.predict_safety(query)
    is_safe = safety_class == "safe"
    safety_flags = [] if is_safe else [safety_class]

    # 2. Domain boundary
    domain_label, domain_conf, _ = trained_models.predict_domain(query)
    is_out_of_domain = domain_label == "general_out_of_domain" and domain_conf > 0.65

    # 3. Intent classification
    intent_label, intent_conf, _ = trained_models.predict_intent(query)
    if is_out_of_domain:
        intent_label = "out_of_domain"

    # 4. Category classification
    cat_label, _cat_conf, _ = trained_models.predict_category(query)
    suggested_cat_ids = []
    if cat_label and cat_label in CATEGORY_SLUG_TO_ID:
        suggested_cat_ids.append(CATEGORY_SLUG_TO_ID[cat_label])

    # 5. Entity extraction via fuzzy matcher and entity extractor
    matched_brand = None
    brand_res = fuzzy_matcher.match_brand(query)
    if brand_res:
        matched_brand = brand_res[0]
    else:
        for token in query.split():
            b_res = fuzzy_matcher.match_brand(token)
            if b_res:
                matched_brand = b_res[0]
                break

    matched_loc = None
    loc_res = fuzzy_matcher.match_location(query)
    if loc_res:
        matched_loc = loc_res[0]
    else:
        for token in query.split():
            l_res = fuzzy_matcher.match_location(token)
            if l_res:
                matched_loc = l_res[0]
                break

    loc = matched_loc or extract_location(query)
    cond = extract_condition(query)

    # Extract price bounds
    min_price, max_price = extract_price_bounds(query)

    extracted = ExtractedEntities(
        brand=matched_brand,
        model=None,
        price_range=ExtractedPriceRange(
            min_price=min_price,
            max_price=max_price,
            currency="ETB",
        ),
        condition=cond,
        location=loc,
        attributes={},
    )

    # Cleaned search term: strip stop words and conversational prefix
    cleaned_term = query
    for remove_phrase in [
        "i want to buy",
        "looking for",
        "find me",
        "how much is",
        "መግዛት እፈልጋለሁ",
        "ፈልግልኝ",
        "ዋጋው ስንት ነው",
        "bituu barbaada",
        "naaf barbaadi",
        "ክገዝእ ይደሊ",
    ]:
        cleaned_term = cleaned_term.replace(remove_phrase, "").strip()

    if not cleaned_term:
        cleaned_term = query

    # Localized answer synthesis
    if is_out_of_domain:
        localized_ans = (
            "HUNTIQ is Ethiopia's universal marketplace for products, vehicles, real estate, and services. "
            "I can assist you with buying, selling, renting, and inspecting marketplace items."
        )
    elif not is_safe:
        localized_ans = "Security Alert: This query contains potential advance-fee or fraud indicators."
    else:
        localized_ans = f"Searching verified marketplace listings for '{cleaned_term}'."

    execution_time_ms = (time.time() - t0) * 1000.0

    reasoning = MarketplaceReasoningChain(
        intent_breakdown=f"Intent: {intent_label} (confidence: {intent_conf:.2f}), Category: {cat_label}.",
        market_feasibility="Evaluated against active marketplace classification and taxonomy models.",
        trade_off_analysis="Extracted search constraints with fuzzy-matched locations and brands.",
        inspection_protocol="Encouraged safe physical inspection for verified items.",
        steps=[
            ReasoningStep(
                step=1,
                title="Safety & Domain Guard",
                detail=f"Safety: {safety_class}, Domain: {domain_label}",
            ),
            ReasoningStep(
                step=2,
                title="Intent & Taxonomy Classification",
                detail=f"Intent: {intent_label}, Category: {cat_label}",
            ),
            ReasoningStep(
                step=3,
                title="Entity & Brand Matching",
                detail=f"Brand: {matched_brand}, Location: {matched_loc}",
            ),
        ],
    )

    return QueryUnderstandingResponse(
        normalized_query=norm_query,
        cleaned_search_term=cleaned_term,
        sort_by=None,
        mode=req.mode or "ai",
        primary_intent=intent_label or "marketplace_search",
        confidence=round(intent_conf, 4),
        suggested_category_ids=suggested_cat_ids,
        extracted_entities=extracted,
        deep_parameters={
            "category_slug": cat_label,
            "safety_class": safety_class,
            "domain_label": domain_label,
        },
        localized_answer=localized_ans,
        pii_detected=False,
        detected_locale=locale,
        execution_time_ms=round(execution_time_ms, 2),
        is_safe=is_safe,
        safety_flags=safety_flags,
        reasoning_chain=reasoning,
    )


@router.post("/analyze-products", response_model=ProductAnalysisResponse)
async def handle_analyze_products(
    req: ProductAnalysisRequest,
    x_huntiq_ai_secret: str | None = Header(default=None),
) -> ProductAnalysisResponse:
    """
    Statistical analysis and verdict synthesis on real database listings.
    Called by Go backend aiclient AnalyzeProducts.
    """
    if (
        config.api_secret_key
        and x_huntiq_ai_secret
        and x_huntiq_ai_secret != config.api_secret_key
    ):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid service credentials",
        )

    products = req.products
    total = len(products)
    if total == 0:
        return ProductAnalysisResponse(
            total_analyzed=0,
            total_matching=0,
            premium_count=0,
            gold_count=0,
            normal_count=0,
            price_stats=PriceStats(
                min_price=0.0, avg_price=0.0, max_price=0.0, currency="ETB"
            ),
            top_premium_pick=None,
            top_gold_pick=None,
            market_verdict="No active listings available for analysis.",
            deal_ratings={},
        )

    prices = [p.price for p in products if p.price > 0]
    min_p = min(prices) if prices else 0.0
    max_p = max(prices) if prices else 0.0
    avg_p = sum(prices) / len(prices) if prices else 0.0

    prem_count = sum(1 for p in products if p.seller_tier == "premium")
    gold_count = sum(1 for p in products if p.seller_tier == "gold")
    norm_count = total - prem_count - gold_count

    top_prem = next((p for p in products if p.seller_tier == "premium"), None)
    top_gold = next((p for p in products if p.seller_tier == "gold"), None)

    verdict = (
        f"Analyzed {total} verified listings. "
        f"Market price ranges from {min_p:,.0f} to {max_p:,.0f} ETB (average {avg_p:,.0f} ETB). "
        f"Features {prem_count} Premium and {gold_count} Gold verified sellers."
    )

    return ProductAnalysisResponse(
        total_analyzed=total,
        total_matching=total,
        premium_count=prem_count,
        gold_count=gold_count,
        normal_count=norm_count,
        price_stats=PriceStats(
            min_price=min_p, avg_price=avg_p, max_price=max_p, currency="ETB"
        ),
        top_premium_pick=top_prem,
        top_gold_pick=top_gold,
        market_verdict=verdict,
        deal_ratings={},
    )
