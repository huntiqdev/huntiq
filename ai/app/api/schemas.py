"""Pydantic schemas for HUNTIQ AI conversational API and tool contracts."""

from typing import Any

from pydantic import BaseModel, Field


class ChatTurn(BaseModel):
    """Single dialogue turn in conversation history."""

    role: str
    content: str


class GuestChatRequest(BaseModel):
    """Incoming request payload for visitor AI chat."""

    message: str = Field(..., min_length=1, max_length=1000)
    locale: str = Field(default="en", max_length=10)
    session_id: str | None = Field(default=None, max_length=128)
    history: list[ChatTurn] = Field(default_factory=list)


class ListingSummary(BaseModel):
    """Authoritative, verified marketplace listing returned from Go PostgreSQL."""

    id: int
    title: str
    price: float
    currency: str = "ETB"
    listing_type: str = "sale"
    rental_price: float | None = None
    rental_duration_value: int | None = None
    rental_duration_unit: str | None = None
    condition: str = "used"
    seller_tier: str = "normal"
    city: str = "Addis Ababa"
    subcity: str | None = None
    detail_url: str
    image_url: str | None = None
    contact_action: str = "Buy / Contact Seller"
    availability_status: str = "Available"


class CategoryGroupResult(BaseModel):
    """Isolated category bucket in a multi-category or single-category search."""

    category_name: str
    category_slug: str | None = None
    search_term: str
    total_count: int = 0
    listings: list[ListingSummary] = Field(default_factory=list)
    error: str | None = None


class GuestChatResponse(BaseModel):
    """Context-aware response payload matching the Next.js frontend contract."""

    session_id: str
    greeting_prefix: str | None = None
    answer: str
    detected_locale: str = "en"
    intents: list[str] = Field(default_factory=list)
    category_groups: list[CategoryGroupResult] | None = None
    listings: list[ListingSummary] | None = None
    clarification: str | None = None
    latency_ms: float = 0.0
    model: str = "huntiq-intelligence-v2"


class InternalToolExecutionRequest(BaseModel):
    """Request schema sent from Python to Go to execute an authorized marketplace tool."""

    tool_name: str
    arguments: dict[str, Any] = Field(default_factory=dict)


class InternalToolExecutionResponse(BaseModel):
    """Response returned from Go backend after executing a database search."""

    success: bool
    tool_name: str
    total_matching: int = 0
    listings: list[ListingSummary] = Field(default_factory=list)
    error: str | None = None


# ---------------------------------------------------------------------------
# Go Backend Client Contract Schemas (Query Understanding & Product Analysis)
# ---------------------------------------------------------------------------


class UnderstandQueryParams(BaseModel):
    """Parameters sent from Go client to understand search queries."""

    query: str
    locale: str | None = "en"
    mode: str | None = "ai"
    preferred_language: str | None = None


class ExtractedPriceRange(BaseModel):
    min_price: float | None = None
    max_price: float | None = None
    currency: str = "ETB"


class ExtractedEntities(BaseModel):
    brand: str | None = None
    model: str | None = None
    price_range: ExtractedPriceRange = Field(default_factory=ExtractedPriceRange)
    condition: str | None = None
    location: str | None = None
    attributes: dict[str, Any] = Field(default_factory=dict)


class ReasoningStep(BaseModel):
    step: int
    title: str
    detail: str


class MarketplaceReasoningChain(BaseModel):
    intent_breakdown: str
    market_feasibility: str
    trade_off_analysis: str
    inspection_protocol: str
    steps: list[ReasoningStep] = Field(default_factory=list)


class QueryUnderstandingResponse(BaseModel):
    """Full validated response payload sent to Go backend."""

    normalized_query: str
    cleaned_search_term: str
    sort_by: str | None = None
    mode: str = "ai"
    primary_intent: str
    confidence: float
    suggested_category_ids: list[int] = Field(default_factory=list)
    extracted_entities: ExtractedEntities = Field(default_factory=ExtractedEntities)
    deep_parameters: dict[str, Any] = Field(default_factory=dict)
    localized_answer: str
    pii_detected: bool = False
    detected_locale: str = "en"
    execution_time_ms: float = 0.0
    is_safe: bool = True
    safety_flags: list[str] = Field(default_factory=list)
    reasoning_chain: MarketplaceReasoningChain | None = None


class ProductSummary(BaseModel):
    id: int
    title: str
    price: float
    seller_tier: str = "normal"
    currency: str = "ETB"
    condition: str | None = None


class PriceStats(BaseModel):
    min_price: float
    avg_price: float
    max_price: float
    currency: str = "ETB"


class ProductAnalysisRequest(BaseModel):
    query: str
    products: list[ProductSummary] = Field(default_factory=list)
    user_query: str | None = None
    locale: str = "en"


class ProductAnalysisResponse(BaseModel):
    total_analyzed: int
    total_matching: int
    premium_count: int
    gold_count: int
    normal_count: int
    price_stats: PriceStats
    top_premium_pick: ProductSummary | None = None
    top_gold_pick: ProductSummary | None = None
    market_verdict: str
    deal_ratings: dict[int, Any] = Field(default_factory=dict)


# ---------------------------------------------------------------------------
# Authenticated Chat Message Schemas (Go Backend SendChatMessage Contract)
# ---------------------------------------------------------------------------


class ChatFunctionCall(BaseModel):
    name: str
    arguments: dict[str, Any] = Field(default_factory=dict)


class ChatToolCall(BaseModel):
    id: str
    type: str = "function"
    function: ChatFunctionCall


class ChatLLMMessage(BaseModel):
    role: str
    content: str = ""
    name: str | None = None
    tool_call_id: str | None = None
    tool_calls: list[ChatToolCall] = Field(default_factory=list)


class AuthenticatedChatMessageRequest(BaseModel):
    conversation_id: str | None = None
    messages: list[ChatLLMMessage] = Field(default_factory=list)
    locale: str = "en"
    tools_enabled: bool = True


class AuthenticatedChatMessageResponse(BaseModel):
    content: str | None = None
    tool_calls: list[ChatToolCall] = Field(default_factory=list)
    finish_reason: str = "stop"
    model: str = "huntiq-intelligence-v2"
    latency_ms: float = 0.0
