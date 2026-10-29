package handlers

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/used-tools-marketplace/internal/service/aiclient"
)

func TestTierAuthorizationLogic(t *testing.T) {
	tests := []struct {
		name         string
		userTier     string
		requiredTier string
		isAdmin      bool
		expected     bool
	}{
		// Guest access
		{"Guest on normal product", "guest", "normal", false, true},
		{"Guest on empty required tier", "guest", "", false, true},
		{"Guest on gold product", "guest", "gold", false, false},
		{"Guest on premium product", "guest", "premium", false, false},

		// Registered normal buyer
		{"Registered buyer on normal product", "normal", "normal", false, true},
		{"Registered buyer on gold product", "normal", "gold", false, false},
		{"Registered buyer on premium product", "normal", "premium", false, false},

		// Gold member
		{"Gold member on normal product", "gold", "normal", false, true},
		{"Gold member on gold product", "gold", "gold", false, true},
		{"Gold member on premium product", "gold", "premium", false, false},

		// Premium member
		{"Premium member on normal product", "premium", "normal", false, true},
		{"Premium member on gold product", "premium", "gold", false, true},
		{"Premium member on premium product", "premium", "premium", false, true},

		// Admin overrides
		{"Admin on gold product without tier", "normal", "gold", true, true},
		{"Admin on premium product as guest", "guest", "premium", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTierAuthorized(tt.userTier, tt.requiredTier, tt.isAdmin)
			if got != tt.expected {
				t.Errorf("isTierAuthorized(%q, %q, %v) = %v; want %v",
					tt.userTier, tt.requiredTier, tt.isAdmin, got, tt.expected)
			}
		})
	}
}

type mockSearchLogger struct {
	mu          sync.Mutex
	calls       []searchLogEvent
	errToReturn error
	recorded    chan struct{}
}

func (m *mockSearchLogger) LogSearch(ctx context.Context, query string, resultCount int, categoryID *int, filters map[string]interface{}, userID *int, sessionID, ipAddress string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, searchLogEvent{
		query:       query,
		resultCount: resultCount,
		categoryID:  categoryID,
		filters:     filters,
		userID:      userID,
		sessionID:   sessionID,
		ipAddress:   ipAddress,
	})
	if m.recorded != nil {
		select {
		case m.recorded <- struct{}{}:
		default:
		}
	}
	return m.errToReturn
}

func TestProductHandler_SearchAnalytics_WorkerDispatchesEvents(t *testing.T) {
	mock := &mockSearchLogger{
		recorded: make(chan struct{}, 1),
	}
	h := &ProductHandler{}
	h.WithSearchLogger(mock)

	catID := 38
	uid := 42
	ev := searchLogEvent{
		query:       "Samsung Galaxy",
		resultCount: 5,
		categoryID:  &catID,
		filters:     map[string]interface{}{"condition": "new"},
		userID:      &uid,
		sessionID:   "sess-xyz-123",
		ipAddress:   "192.168.1.10",
	}

	h.searchChan <- ev

	select {
	case <-mock.recorded:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for search analytics event to be dispatched to logger")
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()
	if len(mock.calls) != 1 {
		t.Fatalf("expected 1 call to LogSearch, got %d", len(mock.calls))
	}
	call := mock.calls[0]
	if call.query != "Samsung Galaxy" {
		t.Errorf("expected query 'Samsung Galaxy', got %q", call.query)
	}
	if call.resultCount != 5 {
		t.Errorf("expected resultCount 5, got %d", call.resultCount)
	}
	if call.categoryID == nil || *call.categoryID != 38 {
		t.Errorf("expected categoryID 38, got %v", call.categoryID)
	}
	if call.userID == nil || *call.userID != 42 {
		t.Errorf("expected userID 42, got %v", call.userID)
	}
	if call.sessionID != "sess-xyz-123" {
		t.Errorf("expected sessionID 'sess-xyz-123', got %q", call.sessionID)
	}
	if call.ipAddress != "192.168.1.10" {
		t.Errorf("expected ipAddress '192.168.1.10', got %q", call.ipAddress)
	}
}

func TestProductHandler_SearchAnalytics_ZeroResultQueries(t *testing.T) {
	mock := &mockSearchLogger{
		recorded: make(chan struct{}, 1),
	}
	h := &ProductHandler{}
	h.WithSearchLogger(mock)

	ev := searchLogEvent{
		query:       "nonexistent query that finds nothing",
		resultCount: 0,
	}

	h.searchChan <- ev

	select {
	case <-mock.recorded:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for zero-result query to be logged")
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()
	if len(mock.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(mock.calls))
	}
	if mock.calls[0].resultCount != 0 {
		t.Errorf("expected 0 results, got %d", mock.calls[0].resultCount)
	}
}

func TestProductHandler_SearchAnalytics_LoggingFailureDoesNotPanic(t *testing.T) {
	mock := &mockSearchLogger{
		recorded:    make(chan struct{}, 1),
		errToReturn: errors.New("database connection refused"),
	}
	h := &ProductHandler{}
	h.WithSearchLogger(mock)

	ev := searchLogEvent{
		query:       "caterpillar excavator",
		resultCount: 2,
	}

	// Should safely log the error and continue without panic
	h.searchChan <- ev

	select {
	case <-mock.recorded:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for worker to process failing event")
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()
	if len(mock.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(mock.calls))
	}
}

func TestProductHandler_SearchAnalytics_BoundedChannelNonBlocking(t *testing.T) {
	h := &ProductHandler{
		searchChan: make(chan searchLogEvent, 2),
	}
	// Fill the buffer to capacity
	h.searchChan <- searchLogEvent{query: "test1"}
	h.searchChan <- searchLogEvent{query: "test2"}

	// Third enqueue must not block
	ev3 := searchLogEvent{query: "test3"}
	dropped := false
	select {
	case h.searchChan <- ev3:
	default:
		dropped = true
		atomic.AddUint64(&h.droppedCount, 1)
	}

	if !dropped {
		t.Errorf("expected 3rd event to drop when buffer is full")
	}
	if h.DroppedSearchEvents() != 1 {
		t.Errorf("expected DroppedSearchEvents() = 1, got %d", h.DroppedSearchEvents())
	}
}

func TestProductHandler_WithAIClient_Wiring(t *testing.T) {
	h := &ProductHandler{}
	if h.aiClient != nil {
		t.Fatalf("expected aiClient to be nil initially")
	}

	client := aiclient.NewClient("http://localhost:8000", "test-secret")
	h.WithAIClient(client)

	if h.aiClient == nil || !h.aiClient.IsEnabled() {
		t.Fatalf("expected aiClient to be wired and enabled")
	}
}

