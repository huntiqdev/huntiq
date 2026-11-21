package ai

type ChatMessage struct {
	Role  string `json:"role"`
	Parts []struct {
		Text string `json:"text"`
	} `json:"parts,omitempty"`
}

type ClassifyRequest struct {
	Message  string        `json:"message"`
	Language string        `json:"language,omitempty"`
	History  []ChatMessage `json:"history,omitempty"`
}

type ClassifyResponse struct {
	IsAllowed        bool     `json:"is_allowed"`
	Intent           string   `json:"intent"`
	Path             string   `json:"path"`
	Language         string   `json:"language"`
	SearchTerm       string   `json:"search_term,omitempty"`
	MinPrice         *float64 `json:"min_price,omitempty"`
	MaxPrice         *float64 `json:"max_price,omitempty"`
	InstantReply     string   `json:"instant_reply,omitempty"`
	SuggestedQueries []string `json:"suggested_queries,omitempty"`
	Reason           string   `json:"reason,omitempty"`
}

type ListingCard struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Price        float64  `json:"price"`
	Currency     string   `json:"currency"`
	CategoryName string   `json:"category_name"`
	LocationName string   `json:"location_name"`
	Condition    string   `json:"condition"`
	ListingType  string   `json:"listing_type"`
	ImageURL     string   `json:"image_url,omitempty"`
	ViewCount    int      `json:"view_count,omitempty"`
}

type PythonChatRequest struct {
	Message      string        `json:"message"`
	Language     string        `json:"language,omitempty"`
	History      []ChatMessage `json:"history,omitempty"`
	Listings     []ListingCard `json:"listings,omitempty"`
	TrendSummary string        `json:"trend_summary,omitempty"`
}

type ChatRequest struct {
	Message  string        `json:"message"`
	Language string        `json:"language,omitempty"`
	History  []ChatMessage `json:"history,omitempty"`
}

type ChatResponse struct {
	Reply            string        `json:"reply"`
	Language         string        `json:"language"`
	Intent           string        `json:"intent"`
	Path             string        `json:"path"`
	IsOutOfDomain    bool          `json:"is_out_of_domain"`
	FallbackUsed     bool          `json:"fallback_used"`
	TokensUsed       int           `json:"tokens_used"`
	LatencyMs        int64         `json:"latency_ms"`
	SuggestedQueries []string      `json:"suggested_queries,omitempty"`
	Listings         []ListingCard `json:"listings,omitempty"`
	Error            string        `json:"error,omitempty"`
}
