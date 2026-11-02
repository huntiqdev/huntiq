"""Thread-safe, bounded in-memory session manager for temporary AI dialogue context."""

import threading
import time

from pydantic import BaseModel, Field

from .entity_extractor import ExtractedSearch

DEFAULT_SESSION_TTL_SECONDS = 20 * 60  # 20 minutes
MAX_SESSIONS = 10000
MAX_TURNS_PER_SESSION = 6


class ConversationTurn(BaseModel):
    role: str
    content: str
    timestamp: float = Field(default_factory=time.time)


class ConversationContext(BaseModel):
    session_id: str
    active_language: str = "en"
    recent_searches: list[ExtractedSearch] = Field(default_factory=list)
    last_max_price: float | None = None
    last_min_price: float | None = None
    last_location: str | None = None
    last_category: str | None = None
    last_brand: str | None = None
    last_listing_type: str | None = None
    last_condition: str | None = None
    unresolved_slot: str | None = None
    turns: list[ConversationTurn] = Field(default_factory=list)
    last_active_at: float = Field(default_factory=time.time)
    created_at: float = Field(default_factory=time.time)


class SessionManager:
    """Manages sliding-TTL, bounded in-memory conversation contexts."""

    def __init__(
        self,
        ttl_seconds: float = DEFAULT_SESSION_TTL_SECONDS,
        max_sessions: int = MAX_SESSIONS,
    ):
        self._lock = threading.Lock()
        self._sessions: dict[str, ConversationContext] = {}
        self._ttl = ttl_seconds
        self._max_sessions = max_sessions

    def get_or_create(
        self, session_id: str, language: str = "en"
    ) -> ConversationContext:
        """Retrieve existing active session or initialize a fresh bounded context."""
        now = time.time()
        with self._lock:
            # Periodic opportunistic cleanup if table is large
            if len(self._sessions) > self._max_sessions:
                self._evict_expired_locked(now)

            ctx = self._sessions.get(session_id)
            if ctx is not None:
                if now - ctx.last_active_at > self._ttl:
                    # Expired, reset
                    ctx = ConversationContext(
                        session_id=session_id,
                        active_language=language,
                        last_active_at=now,
                        created_at=now,
                    )
                    self._sessions[session_id] = ctx
                else:
                    ctx.last_active_at = now
                    if language:
                        ctx.active_language = language
                return ctx

            ctx = ConversationContext(
                session_id=session_id,
                active_language=language,
                last_active_at=now,
                created_at=now,
            )
            self._sessions[session_id] = ctx
            return ctx

    def record_turn(
        self,
        session_id: str,
        user_message: str,
        assistant_answer: str,
        searches: list[ExtractedSearch] | None = None,
        unresolved_slot: str | None = None,
    ) -> None:
        """Append user/assistant turn and update slots in conversation state."""
        now = time.time()
        with self._lock:
            ctx = self._sessions.get(session_id)
            if ctx is None:
                ctx = ConversationContext(
                    session_id=session_id,
                    last_active_at=now,
                    created_at=now,
                )
                self._sessions[session_id] = ctx

            ctx.last_active_at = now
            ctx.turns.append(
                ConversationTurn(role="user", content=user_message, timestamp=now)
            )
            ctx.turns.append(
                ConversationTurn(
                    role="assistant", content=assistant_answer, timestamp=now
                )
            )

            # Bound turns to MAX_TURNS_PER_SESSION
            if len(ctx.turns) > MAX_TURNS_PER_SESSION:
                ctx.turns = ctx.turns[-MAX_TURNS_PER_SESSION:]

            ctx.unresolved_slot = unresolved_slot

            if searches:
                ctx.recent_searches = searches[:4]
                latest = searches[0]
                if latest.category_name:
                    ctx.last_category = latest.category_name
                if latest.brand:
                    ctx.last_brand = latest.brand
                if latest.location:
                    ctx.last_location = latest.location
                if latest.max_price is not None:
                    ctx.last_max_price = latest.max_price
                if latest.min_price is not None:
                    ctx.last_min_price = latest.min_price
                if latest.listing_type:
                    ctx.last_listing_type = latest.listing_type
                if latest.condition:
                    ctx.last_condition = latest.condition

    def apply_contextual_followup(
        self, session_id: str, query: str, searches: list[ExtractedSearch]
    ) -> list[ExtractedSearch]:
        """
        Merge prior slot state into follow-up requests (e.g. 'Show me cheaper ones',
        'I want to rent, not buy', 'Find a phone near Bole').
        """
        lower = query.strip().lower()
        with self._lock:
            ctx = self._sessions.get(session_id)
            if ctx is None or not ctx.recent_searches:
                return searches

            # 1. Price adjustment ("cheaper", "less expensive", "ከዚህ ያነሰ")
            if any(
                k in lower
                for k in ["cheaper", "less expensive", "ያነሰ", "ዋጋው የቀነሰ", "gatii gadi"]
            ):
                base = ctx.recent_searches[0].model_copy()
                if ctx.last_max_price is not None and ctx.last_max_price > 0:
                    base.max_price = ctx.last_max_price * 0.8
                return [base]

            # 2. Transaction type correction ("I want to rent, not buy", "ለመከራየት")
            if any(
                k in lower
                for k in [
                    "want to rent",
                    "to rent",
                    "for rent",
                    "ለመከራየት",
                    "kiraa barbaada",
                ]
            ):
                base = ctx.recent_searches[0].model_copy()
                base.listing_type = "rent"
                return [base]

            if any(
                k in lower
                for k in ["want to buy", "for sale", "ለመግዛት", "bituu barbaada"]
            ):
                base = ctx.recent_searches[0].model_copy()
                base.listing_type = "sale"
                return [base]

            # 3. Location refinement ("in Bole", "near Bole", "ቦሌ ውስጥ")
            if any(k in lower for k in ["near bole", "in bole", "around bole", "ቦሌ"]):
                if not searches:
                    base = ctx.recent_searches[0].model_copy()
                    base.location = "Bole"
                    return [base]
                for s in searches:
                    s.location = "Bole"
                return searches

            # 4. Brand refinement ("what about samsung?", "ሳምሰንግስ?")
            if "samsung" in lower or "ሳምሰንግ" in lower:
                base = ctx.recent_searches[0].model_copy()
                base.brand = "Samsung"
                base.search_term = "Samsung"
                return [base]

            # 5. Slot inheritance: if searches lack location or listing_type, inherit from recent context
            for s in searches:
                if not s.location and ctx.last_location:
                    s.location = ctx.last_location
                if not s.listing_type and ctx.last_listing_type:
                    s.listing_type = ctx.last_listing_type

            return searches

    def _evict_expired_locked(self, now: float) -> None:
        """Remove sessions exceeding the TTL."""
        expired_keys = [
            sid
            for sid, ctx in self._sessions.items()
            if now - ctx.last_active_at > self._ttl
        ]
        for sid in expired_keys:
            del self._sessions[sid]


# Global singleton instance
session_manager = SessionManager()
