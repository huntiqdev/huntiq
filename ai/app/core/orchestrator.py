"""Unified conversational AI orchestrator for HUNTIQ."""

import json
import time
import uuid
from typing import Any

from ..api.schemas import (
    AuthenticatedChatMessageRequest,
    AuthenticatedChatMessageResponse,
    CategoryGroupResult,
    ChatFunctionCall,
    ChatToolCall,
    GuestChatRequest,
    GuestChatResponse,
    ListingSummary,
)
from ..integrations.go_marketplace_client import go_client
from .context_manager import session_manager
from .intent import IntentType, classify_message
from .policy import get_out_of_domain_redirection
from .response_grounding import compose_grounded_response


class ConversationalOrchestrator:
    """Executes the complete natural language interpretation, planning, and grounding flow."""

    async def process_chat(self, req: GuestChatRequest) -> GuestChatResponse:
        t0 = time.time()
        session_id = req.session_id or str(uuid.uuid4())
        locale = req.locale or "en"

        # 1. Resolve bounded temporary context
        _ctx = session_manager.get_or_create(session_id, locale)

        # 2. Multidimensional intent classification and entity slot extraction
        plan = classify_message(req.message, locale)
        effective_locale = plan.detected_language or locale

        # 3. Apply conversational follow-up merging if ongoing dialogue
        effective_searches = session_manager.apply_contextual_followup(
            session_id, req.message, plan.searches
        )

        # 3.5. Safety & Scam Guard Check
        if (
            IntentType.SAFETY_VIOLATION in plan.intents
            or IntentType.SCAM_ATTEMPT in plan.intents
        ):
            answer = (
                plan.safety_warning or "Please follow HUNTIQ safe marketplace policies."
            )
            session_manager.record_turn(session_id, req.message, answer, None)
            latency = (time.time() - t0) * 1000.0
            return GuestChatResponse(
                session_id=session_id,
                answer=answer,
                detected_locale=effective_locale,
                intents=[i.value for i in plan.intents],
                latency_ms=round(latency, 2),
                model="huntiq-guard:safety-moderation-v1",
            )

        # 4. Out-of-Domain Guard
        if IntentType.OUT_OF_DOMAIN in plan.intents:
            answer = get_out_of_domain_redirection(effective_locale)
            session_manager.record_turn(session_id, req.message, answer, None)
            latency = (time.time() - t0) * 1000.0
            return GuestChatResponse(
                session_id=session_id,
                answer=answer,
                detected_locale=effective_locale,
                intents=[i.value for i in plan.intents],
                latency_ms=round(latency, 2),
                model="huntiq-guard:domain-policy-v2",
            )

        # 5. Platform FAQ Check
        if IntentType.PLATFORM_FAQ in plan.intents and plan.platform_faq_response:
            answer = plan.platform_faq_response
            session_manager.record_turn(session_id, req.message, answer, None)
            latency = (time.time() - t0) * 1000.0
            return GuestChatResponse(
                session_id=session_id,
                answer=answer,
                detected_locale=effective_locale,
                intents=[i.value for i in plan.intents],
                latency_ms=round(latency, 2),
                model="huntiq-kb:platform-policy-v2",
            )

        # 6. Ambiguity Clarification Check
        if plan.requires_clarification and plan.clarification_question:
            clarification = plan.clarification_question
            session_manager.record_turn(
                session_id,
                req.message,
                clarification,
                searches=effective_searches,
                unresolved_slot=plan.unresolved_slot,
            )
            latency = (time.time() - t0) * 1000.0
            return GuestChatResponse(
                session_id=session_id,
                answer=clarification,
                detected_locale=effective_locale,
                intents=[i.value for i in plan.intents],
                clarification=clarification,
                latency_ms=round(latency, 2),
                model="huntiq-assistant:clarification-v2",
            )

        # 7. Conversational / Social / Identity Inquiries (Zero DB search)
        if not plan.requires_database_search and not effective_searches:
            ans_parts = []
            if plan.social_greeting:
                ans_parts.append(plan.social_greeting)
            if plan.identity_response:
                ans_parts.append(plan.identity_response)
            elif plan.social_farewell:
                ans_parts.append(plan.social_farewell)

            ans = " ".join(ans_parts).strip()
            if not ans:
                ans = (
                    plan.social_greeting
                    or "Hello! How can I help with your HUNTIQ marketplace search today?"
                )

            session_manager.record_turn(session_id, req.message, ans, None)
            latency = (time.time() - t0) * 1000.0
            return GuestChatResponse(
                session_id=session_id,
                greeting_prefix=plan.social_greeting,
                answer=ans,
                detected_locale=effective_locale,
                intents=[i.value for i in plan.intents],
                latency_ms=round(latency, 2),
                model="huntiq-assistant:conversational-v2",
            )

        # 8. Authorized Marketplace Search via Go Backend
        category_groups: list[CategoryGroupResult] = []
        flat_listings: list[ListingSummary] = []

        for s in effective_searches:
            listings = await go_client.search_marketplace(
                search_term=s.search_term,
                category_slug=s.category_slug,
                min_price=s.min_price,
                max_price=s.max_price,
                location=s.location,
                condition=s.condition,
                listing_type=s.listing_type,
                limit=5,
            )
            category_groups.append(
                CategoryGroupResult(
                    category_name=s.category_name,
                    category_slug=s.category_slug,
                    search_term=s.search_term,
                    total_count=len(listings),
                    listings=listings,
                )
            )
            flat_listings.extend(listings)

        # 9. Grounded Summary Synthesis
        answer = compose_grounded_response(
            query=req.message,
            greeting_prefix=plan.social_greeting,
            groups=category_groups,
            locale=effective_locale,
        )

        session_manager.record_turn(session_id, req.message, answer, effective_searches)
        latency = (time.time() - t0) * 1000.0

        return GuestChatResponse(
            session_id=session_id,
            greeting_prefix=plan.social_greeting,
            answer=answer,
            detected_locale=effective_locale,
            intents=[i.value for i in plan.intents],
            category_groups=category_groups if len(category_groups) > 0 else None,
            listings=flat_listings if len(flat_listings) > 0 else None,
            latency_ms=round(latency, 2),
            model="huntiq-intelligence-v2",
        )

    async def process_chat_message(
        self, req: AuthenticatedChatMessageRequest
    ) -> AuthenticatedChatMessageResponse:
        """
        Process authenticated multi-turn chat with tool calling integration.
        Called by Go backend aiclient SendChatMessage.
        """
        t0 = time.time()
        locale = req.locale or "en"
        messages = req.messages

        if not messages:
            msg = "Hello! How can I assist you with the HUNTIQ marketplace today?"
            return AuthenticatedChatMessageResponse(
                content=msg,
                tool_calls=[],
                finish_reason="stop",
                model="huntiq-intelligence-v2",
                latency_ms=(time.time() - t0) * 1000.0,
            )

        last_msg = messages[-1]

        # Case A: If last message is a TOOL result, synthesize final response
        if last_msg.role == "tool":
            listings_data = []
            try:
                tool_data = json.loads(last_msg.content)
                listings_data = tool_data.get("listings", [])
            except (json.JSONDecodeError, TypeError, KeyError):
                listings_data = []

            if not listings_data:
                ans = (
                    "I searched our verified marketplace records, but couldn't find matching active listings at this time. "
                    "You can try adjusting your search terms, price range, or city."
                )
            else:
                lines = [
                    f"I found {len(listings_data)} verified listings matching your request:"
                ]
                for item in listings_data[:5]:
                    title = item.get("title", "Marketplace item")
                    price = item.get("price", 0)
                    city = item.get("city", "Addis Ababa")
                    tier = item.get("seller_tier", "normal")
                    tier_str = (
                        f" [{tier.capitalize()} Seller]"
                        if tier in ("premium", "gold")
                        else ""
                    )
                    lines.append(f"• {title} - {price:,.0f} ETB ({city}){tier_str}")
                lines.append(
                    "\nFeel free to ask for contact details, price comparisons, or safety recommendations!"
                )
                ans = "\n".join(lines)

            latency = (time.time() - t0) * 1000.0
            return AuthenticatedChatMessageResponse(
                content=ans,
                tool_calls=[],
                finish_reason="stop",
                model="huntiq-intelligence-v2",
                latency_ms=round(latency, 2),
            )

        # Case B: Last message is from USER
        user_query = last_msg.content.strip()
        plan = classify_message(user_query, locale)
        effective_locale = plan.detected_language or locale

        # Safety & Scam Guard Check
        if (
            IntentType.SAFETY_VIOLATION in plan.intents
            or IntentType.SCAM_ATTEMPT in plan.intents
        ):
            answer = (
                plan.safety_warning
                or "Security Warning: Please follow HUNTIQ safe marketplace policies."
            )
            return AuthenticatedChatMessageResponse(
                content=answer,
                tool_calls=[],
                finish_reason="stop",
                model="huntiq-guard:safety-moderation-v1",
                latency_ms=(time.time() - t0) * 1000.0,
            )

        # Out-of-Domain Guard
        if IntentType.OUT_OF_DOMAIN in plan.intents:
            answer = get_out_of_domain_redirection(effective_locale)
            return AuthenticatedChatMessageResponse(
                content=answer,
                tool_calls=[],
                finish_reason="stop",
                model="huntiq-guard:domain-policy-v2",
                latency_ms=(time.time() - t0) * 1000.0,
            )

        # Platform FAQ Check
        if IntentType.PLATFORM_FAQ in plan.intents and plan.platform_faq_response:
            return AuthenticatedChatMessageResponse(
                content=plan.platform_faq_response,
                tool_calls=[],
                finish_reason="stop",
                model="huntiq-kb:platform-policy-v2",
                latency_ms=(time.time() - t0) * 1000.0,
            )

        # Ambiguity Clarification Check
        if plan.requires_clarification and plan.clarification_question:
            return AuthenticatedChatMessageResponse(
                content=plan.clarification_question,
                tool_calls=[],
                finish_reason="stop",
                model="huntiq-assistant:clarification-v2",
                latency_ms=(time.time() - t0) * 1000.0,
            )

        # Conversational / Social
        if not plan.requires_database_search and not plan.searches:
            ans_parts = []
            if plan.social_greeting:
                ans_parts.append(plan.social_greeting)
            if plan.identity_response:
                ans_parts.append(plan.identity_response)
            elif plan.social_farewell:
                ans_parts.append(plan.social_farewell)

            ans = " ".join(ans_parts).strip()
            if not ans:
                ans = (
                    plan.social_greeting
                    or "Hello! How can I assist you with your HUNTIQ marketplace search today?"
                )

            return AuthenticatedChatMessageResponse(
                content=ans,
                tool_calls=[],
                finish_reason="stop",
                model="huntiq-assistant:conversational-v2",
                latency_ms=(time.time() - t0) * 1000.0,
            )

        # Marketplace Search Intent: Tool Calling vs Direct Execution
        if req.tools_enabled and plan.searches:
            s0 = plan.searches[0]
            call_id = f"call_{uuid.uuid4().hex[:12]}"
            arguments: dict[str, Any] = {"search_term": s0.search_term}
            if s0.min_price is not None:
                arguments["min_price"] = float(s0.min_price)
            if s0.max_price is not None:
                arguments["max_price"] = float(s0.max_price)
            if s0.condition:
                arguments["condition"] = s0.condition
            if s0.location:
                arguments["location"] = s0.location
            if s0.listing_type:
                arguments["listing_type"] = s0.listing_type
            arguments["limit"] = 5

            tool_call = ChatToolCall(
                id=call_id,
                type="function",
                function=ChatFunctionCall(
                    name="search_marketplace_listings",
                    arguments=arguments,
                ),
            )
            return AuthenticatedChatMessageResponse(
                content=None,
                tool_calls=[tool_call],
                finish_reason="tool_calls",
                model="huntiq-intelligence-v2",
                latency_ms=(time.time() - t0) * 1000.0,
            )

        # If tools not enabled, process via guest chat logic
        guest_req = GuestChatRequest(
            message=user_query,
            locale=locale,
            session_id=req.conversation_id,
        )
        res = await self.process_chat(guest_req)
        return AuthenticatedChatMessageResponse(
            content=res.answer,
            tool_calls=[],
            finish_reason="stop",
            model="huntiq-intelligence-v2",
            latency_ms=(time.time() - t0) * 1000.0,
        )


# Global singleton instance
orchestrator = ConversationalOrchestrator()
