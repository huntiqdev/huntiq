package ai

import (
	"time"
)

// Supported Language Codes (Strict 4-Language Policy)
const (
	LangEN = "en" // English
	LangAM = "am" // Amharic
	LangTI = "ti" // Tigrinya
	LangOR = "or" // Afaan Oromo
)

// IsSupportedLanguage checks if the language code is one of the four official HUNTIQ languages.
func IsSupportedLanguage(lang string) bool {
	switch lang {
	case LangEN, LangAM, LangTI, LangOR:
		return true
	default:
		return false
	}
}

// Intent Classification Types
type DomainIntent string

const (
	IntentMarketplaceDiscovery DomainIntent = "MARKETPLACE_DISCOVERY"
	IntentListingSearch        DomainIntent = "LISTING_SEARCH"
	IntentPolicyHelp           DomainIntent = "POLICY_HELP"
	IntentAccountHelp          DomainIntent = "ACCOUNT_HELP"
	IntentAmbiguous            DomainIntent = "AMBIGUOUS"
	IntentOutOfDomain          DomainIntent = "OUT_OF_DOMAIN"
	IntentAdversarialBypass    DomainIntent = "ADVERSARIAL_BYPASS"
)

// Audit Event Types
const (
	AuditEventQuery             = "QUERY"
	AuditEventOutOfDomain       = "OUT_OF_DOMAIN"
	AuditEventAdversarialBlocked= "PROMPT_INJECTION_BLOCKED"
	AuditEventToolCalled        = "TOOL_CALLED"
	AuditEventProviderError     = "PROVIDER_ERROR"
	AuditEventRateLimited       = "RATE_LIMITED"
	AuditEventUnsupportedLang   = "UNSUPPORTED_LANGUAGE"
)

// ChatRequest is the client-facing payload for conversational marketplace discovery.
type ChatRequest struct {
	SessionID string `json:"session_id"`           // Unique client conversation UUID
	Message   string `json:"message"`              // User query / input text
	Language  string `json:"language,omitempty"`   // Preferred language (en, am, ti, or)
	Category  string `json:"category,omitempty"`   // Optional category context
	Location  string `json:"location,omitempty"`   // Optional location context
}

// ChatResponse is the server response returned to the client.
type ChatResponse struct {
	SessionID      string              `json:"session_id"`
	Language       string              `json:"language"`
	Reply          string              `json:"reply"`                 // Conversational answer
	Intent         DomainIntent        `json:"intent"`                // Classified domain intent
	Listings       []ListingCardDTO    `json:"listings,omitempty"`    // Real matching listing cards
	SuggestedQuery []string            `json:"suggested_queries,omitempty"` // Helpful follow-ups
	IsOutOfDomain  bool                `json:"is_out_of_domain"`
	FallbackUsed   bool                `json:"fallback_used"`
	TokensUsed     int                 `json:"tokens_used,omitempty"`
	LatencyMs      int64               `json:"latency_ms"`
}

// ListingCardDTO represents an authorized real marketplace listing shown beneath AI responses.
type ListingCardDTO struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Price          float64   `json:"price"`
	Currency       string    `json:"currency"`
	CategoryName   string    `json:"category_name"`
	LocationName   string    `json:"location_name"`
	PrimaryImage   string    `json:"primary_image,omitempty"`
	ImageURLs      []string  `json:"image_urls,omitempty"`
	URL            string    `json:"url"`
	Condition      string    `json:"condition,omitempty"`
	ListingType    string    `json:"listing_type,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// Conversation represents a persistent multi-turn conversation session.
type Conversation struct {
	ID           int64     `json:"id"`
	SessionID    string    `json:"session_id"`
	UserID       *int      `json:"user_id,omitempty"`
	Title        string    `json:"title"`
	LanguageCode string    `json:"language_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsArchived   bool      `json:"is_archived"`
}

// Message represents an individual turn in a conversation.
type Message struct {
	ID             int64        `json:"id"`
	ConversationID int64        `json:"conversation_id"`
	Role           string       `json:"role"` // 'user', 'assistant', 'system', 'tool'
	Content        string       `json:"content"`
	Intent         DomainIntent `json:"intent"`
	ToolCalls      interface{}  `json:"tool_calls,omitempty"`
	ToolResults    interface{}  `json:"tool_results,omitempty"`
	Metadata       interface{}  `json:"metadata,omitempty"`
	TokenCount     int          `json:"token_count"`
	LatencyMs      int          `json:"latency_ms"`
	CreatedAt      time.Time    `json:"created_at"`
}

// AuditEvent represents a logged security, domain, or operational event.
type AuditEvent struct {
	ID         int64                  `json:"id"`
	SessionID  string                 `json:"session_id,omitempty"`
	UserID     *int                   `json:"user_id,omitempty"`
	EventType  string                 `json:"event_type"`
	Details    map[string]interface{} `json:"details"`
	IPAddress  string                 `json:"ip_address,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// AIHealthStatus describes the operational state of HUNTIQ AI.
type AIHealthStatus struct {
	Status           string    `json:"status"` // "operational", "degraded", "disabled"
	Provider         string    `json:"provider"`
	Model            string    `json:"model"`
	Enabled          bool      `json:"enabled"`
	LatencyMs        int64     `json:"latency_ms"`
	LastChecked      time.Time `json:"last_checked"`
	SupportedLanguages []string `json:"supported_languages"`
}

// AIMetricsReport holds administrative metrics for monitoring.
type AIMetricsReport struct {
	TotalQueries       int64   `json:"total_queries"`
	OutOfDomainQueries int64   `json:"out_of_domain_queries"`
	AdversarialBlocks  int64   `json:"adversarial_blocks"`
	ProviderErrors     int64   `json:"provider_errors"`
	TotalTokens        int64   `json:"total_tokens"`
	EstimatedCostUSD   float64 `json:"estimated_cost_usd"`
	AverageLatencyMs   float64 `json:"average_latency_ms"`
	ActiveSessions24h  int64   `json:"active_sessions_24h"`
}
