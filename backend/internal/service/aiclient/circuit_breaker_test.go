package aiclient

import (
	"testing"
	"time"
)

func TestCircuitBreaker_StateTransitions(t *testing.T) {
	cb := NewCircuitBreaker(3, 50*time.Millisecond)

	// Initially CLOSED
	if cb.State() != StateClosed {
		t.Fatalf("expected state CLOSED, got %s", cb.State())
	}
	if !cb.Allow() {
		t.Fatalf("expected Allow to be true in CLOSED state")
	}

	// 2 failures should not trip
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.State() != StateClosed {
		t.Fatalf("expected state CLOSED after 2 failures, got %s", cb.State())
	}

	// 3rd failure should trip to OPEN
	cb.RecordFailure()
	if cb.State() != StateOpen {
		t.Fatalf("expected state OPEN after 3 failures, got %s", cb.State())
	}

	// Requests should be denied while OPEN
	if cb.Allow() {
		t.Fatalf("expected Allow to be false in OPEN state before cooldown")
	}

	// Wait for cooldown
	time.Sleep(60 * time.Millisecond)

	// Now should transition to HALF_OPEN and allow probe
	if !cb.Allow() {
		t.Fatalf("expected Allow to be true after cooldown (HALF_OPEN)")
	}
	if cb.State() != StateHalfOpen {
		t.Fatalf("expected state HALF_OPEN, got %s", cb.State())
	}

	// Probe succeeds -> reset to CLOSED
	cb.RecordSuccess()
	if cb.State() != StateClosed {
		t.Fatalf("expected state CLOSED after probe success, got %s", cb.State())
	}
}

func TestCircuitBreaker_ProbeFailureRetrips(t *testing.T) {
	cb := NewCircuitBreaker(2, 20*time.Millisecond)

	cb.RecordFailure()
	cb.RecordFailure()
	if cb.State() != StateOpen {
		t.Fatalf("expected state OPEN")
	}

	time.Sleep(30 * time.Millisecond)
	if !cb.Allow() {
		t.Fatalf("expected Allow to be true for probe")
	}
	if cb.State() != StateHalfOpen {
		t.Fatalf("expected state HALF_OPEN")
	}

	// Probe fails -> should immediately re-trip to OPEN
	cb.RecordFailure()
	if cb.State() != StateOpen {
		t.Fatalf("expected state OPEN after probe failure, got %s", cb.State())
	}
}
