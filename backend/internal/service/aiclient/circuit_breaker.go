package aiclient

import (
	"sync"
	"time"
)

// CircuitBreakerState represents the operational state of the circuit breaker.
type CircuitBreakerState string

const (
	StateClosed   CircuitBreakerState = "CLOSED"
	StateOpen     CircuitBreakerState = "OPEN"
	StateHalfOpen CircuitBreakerState = "HALF_OPEN"
)

// CircuitBreaker protects the application from cascading failures when an external
// or worker AI service experiences latency spikes, network partition, or downtime.
type CircuitBreaker struct {
	mu              sync.RWMutex
	state           CircuitBreakerState
	failureCount    int
	threshold       int
	cooldown        time.Duration
	lastStateChange time.Time
}

// NewCircuitBreaker creates a circuit breaker with configurable failure threshold and cooldown.
func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	if threshold <= 0 {
		threshold = 5
	}
	if cooldown <= 0 {
		cooldown = 30 * time.Second
	}
	return &CircuitBreaker{
		state:           StateClosed,
		threshold:       threshold,
		cooldown:        cooldown,
		lastStateChange: time.Now(),
	}
}

// Allow reports whether a new outgoing request is permitted to proceed.
func (cb *CircuitBreaker) Allow() bool {
	if cb == nil {
		return true
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if now.Sub(cb.lastStateChange) >= cb.cooldown {
			cb.state = StateHalfOpen
			cb.lastStateChange = now
			return true // Allow single probe request
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return true
	}
}

// RecordSuccess records a successful call, resetting failures and closing the circuit.
func (cb *CircuitBreaker) RecordSuccess() {
	if cb == nil {
		return
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount = 0
	if cb.state != StateClosed {
		cb.state = StateClosed
		cb.lastStateChange = time.Now()
	}
}

// RecordFailure records a failed call. Trips the circuit breaker to OPEN when threshold is met.
func (cb *CircuitBreaker) RecordFailure() {
	if cb == nil {
		return
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	if cb.state == StateHalfOpen || cb.failureCount >= cb.threshold {
		cb.state = StateOpen
		cb.lastStateChange = time.Now()
	}
}

// State returns the current operational state of the circuit breaker.
func (cb *CircuitBreaker) State() CircuitBreakerState {
	if cb == nil {
		return StateClosed
	}

	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}
