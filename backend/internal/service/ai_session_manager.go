package service

import (
	"sync"
	"time"
)

// ConversationTurn represents a single conversational exchange in memory.
type ConversationTurn struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ConversationContext represents the short-term, bounded conversational state.
// It is stored purely in-memory and never written to PostgreSQL.
type ConversationContext struct {
	SessionID      string                   `json:"session_id"`
	ActiveLanguage string                   `json:"active_language"`
	RecentSearches []ExtractedSearchRequest `json:"recent_searches"`
	LastMaxPrice   *float64                 `json:"last_max_price,omitempty"`
	LastMinPrice   *float64                 `json:"last_min_price,omitempty"`
	LastLocation   string                   `json:"last_location,omitempty"`
	LastCategory   string                   `json:"last_category,omitempty"`
	Turns          []ConversationTurn       `json:"turns"`
	LastActiveAt   time.Time                `json:"last_active_at"`
	CreatedAt      time.Time                `json:"created_at"`
}

// SessionManager manages thread-safe, bounded, in-memory dialogue sessions.
// Enforces a strict 20-minute sliding TTL, max 6 turns per session, and
// automatic eviction if active sessions exceed maxSessions (default 10,000).
type SessionManager struct {
	mu          sync.RWMutex
	sessions    map[string]*ConversationContext
	ttl         time.Duration
	maxSessions int
	stopChan    chan struct{}
}

// Default constants for session management
const (
	DefaultSessionTTL  = 20 * time.Minute
	DefaultMaxSessions = 10000
	MaxTurnsPerSession = 6
)

// NewSessionManager creates a new in-memory SessionManager and starts a background cleaner.
func NewSessionManager(ttl time.Duration, maxSessions int) *SessionManager {
	if ttl <= 0 {
		ttl = DefaultSessionTTL
	}
	if maxSessions <= 0 {
		maxSessions = DefaultMaxSessions
	}

	sm := &SessionManager{
		sessions:    make(map[string]*ConversationContext),
		ttl:         ttl,
		maxSessions: maxSessions,
		stopChan:    make(chan struct{}),
	}

	go sm.cleanupWorker(2 * time.Minute)

	return sm
}

// GetOrCreate retrieves an existing session or creates a new one with the given language.
func (sm *SessionManager) GetOrCreate(sessionID string, language string) *ConversationContext {
	if sessionID == "" {
		return nil
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	ctx, exists := sm.sessions[sessionID]
	if exists {
		// Verify if expired
		if now.Sub(ctx.LastActiveAt) > sm.ttl {
			// Expired, reset context
			ctx = &ConversationContext{
				SessionID:      sessionID,
				ActiveLanguage: language,
				RecentSearches: make([]ExtractedSearchRequest, 0, 4),
				Turns:          make([]ConversationTurn, 0, MaxTurnsPerSession),
				LastActiveAt:   now,
				CreatedAt:      now,
			}
			sm.sessions[sessionID] = ctx
			return ctx
		}
		ctx.LastActiveAt = now
		if language != "" && ctx.ActiveLanguage != language {
			ctx.ActiveLanguage = language
		}
		return ctx
	}

	// Evict oldest if capacity exceeded
	if len(sm.sessions) >= sm.maxSessions {
		sm.evictOldestLocked()
	}

	ctx = &ConversationContext{
		SessionID:      sessionID,
		ActiveLanguage: language,
		RecentSearches: make([]ExtractedSearchRequest, 0, 4),
		Turns:          make([]ConversationTurn, 0, MaxTurnsPerSession),
		LastActiveAt:   now,
		CreatedAt:      now,
	}
	sm.sessions[sessionID] = ctx
	return ctx
}

// Get safely retrieves a session snapshot if present and not expired.
func (sm *SessionManager) Get(sessionID string) (*ConversationContext, bool) {
	if sessionID == "" {
		return nil, false
	}

	sm.mu.RLock()
	defer sm.mu.RUnlock()

	ctx, exists := sm.sessions[sessionID]
	if !exists {
		return nil, false
	}

	if time.Since(ctx.LastActiveAt) > sm.ttl {
		return nil, false
	}

	return ctx, true
}

// RecordTurn appends user and assistant dialogue, updates search entities, and caps turns.
func (sm *SessionManager) RecordTurn(
	sessionID string,
	userMsg string,
	aiMsg string,
	plan *ConversationPlan,
	searches []ExtractedSearchRequest,
) {
	if sessionID == "" {
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	ctx, exists := sm.sessions[sessionID]
	if !exists {
		return
	}

	now := time.Now()
	ctx.LastActiveAt = now

	// Update turns
	if userMsg != "" {
		ctx.Turns = append(ctx.Turns, ConversationTurn{
			Role:      "user",
			Content:   userMsg,
			Timestamp: now,
		})
	}
	if aiMsg != "" {
		ctx.Turns = append(ctx.Turns, ConversationTurn{
			Role:      "assistant",
			Content:   aiMsg,
			Timestamp: now,
		})
	}

	// Keep only the most recent MaxTurnsPerSession turns
	if len(ctx.Turns) > MaxTurnsPerSession {
		ctx.Turns = ctx.Turns[len(ctx.Turns)-MaxTurnsPerSession:]
	}

	// Update active searches and entities
	if len(searches) > 0 {
		ctx.RecentSearches = append(ctx.RecentSearches, searches...)
		if len(ctx.RecentSearches) > 4 {
			ctx.RecentSearches = ctx.RecentSearches[len(ctx.RecentSearches)-4:]
		}

		last := searches[len(searches)-1]
		if last.CategoryName != "" {
			ctx.LastCategory = last.CategoryName
		}
		if last.MaxPrice != nil {
			ctx.LastMaxPrice = last.MaxPrice
		}
		if last.MinPrice != nil {
			ctx.LastMinPrice = last.MinPrice
		}
		if last.Location != "" {
			ctx.LastLocation = last.Location
		}
	}
}

// Clear removes a session from memory immediately.
func (sm *SessionManager) Clear(sessionID string) {
	if sessionID == "" {
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, sessionID)
}

// Len returns the current number of active in-memory sessions.
func (sm *SessionManager) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.sessions)
}

// Stop shuts down the background cleaner goroutine.
func (sm *SessionManager) Stop() {
	select {
	case <-sm.stopChan:
		// already stopped
	default:
		close(sm.stopChan)
	}
}

// evictOldestLocked deletes the session with the oldest LastActiveAt. Must be called with sm.mu held.
func (sm *SessionManager) evictOldestLocked() {
	var oldestID string
	var oldestTime time.Time

	first := true
	for id, s := range sm.sessions {
		if first || s.LastActiveAt.Before(oldestTime) {
			oldestID = id
			oldestTime = s.LastActiveAt
			first = false
		}
	}

	if oldestID != "" {
		delete(sm.sessions, oldestID)
	}
}

// cleanupWorker periodically purges expired sessions.
func (sm *SessionManager) cleanupWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-sm.stopChan:
			return
		case <-ticker.C:
			sm.purgeExpired()
		}
	}
}

func (sm *SessionManager) purgeExpired() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	for id, s := range sm.sessions {
		if now.Sub(s.LastActiveAt) > sm.ttl {
			delete(sm.sessions, id)
		}
	}
}
