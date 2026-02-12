package reliability

import (
	"errors"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker.
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case CircuitClosed:
		return "closed"
	case CircuitOpen:
		return "open"
	case CircuitHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// ErrCircuitOpen is returned when the circuit is open.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitBreakerConfig holds configuration for a circuit breaker.
type CircuitBreakerConfig struct {
	FailureThreshold int           // Consecutive failures to trip
	ResetTimeout     time.Duration // Time before attempting half-open
	HalfOpenMaxCalls int           // Max calls in half-open state
}

// DefaultPhase1CircuitConfig returns default config for Phase 1.
func DefaultPhase1CircuitConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 5,
		ResetTimeout:     30 * time.Second,
		HalfOpenMaxCalls: 1,
	}
}

// DefaultPhase2CircuitConfig returns default config for Phase 2.
func DefaultPhase2CircuitConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 3,
		ResetTimeout:     20 * time.Second,
		HalfOpenMaxCalls: 1,
	}
}

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	config CircuitBreakerConfig

	mu               sync.RWMutex
	state            CircuitState
	failures         int
	lastFailureTime  time.Time
	halfOpenCalls    int
	halfOpenFailures int
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  CircuitClosed,
	}
}

// State returns the current circuit state.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Allow checks if a request is allowed and records intention.
// Returns true if allowed, false if circuit is open.
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitClosed:
		return true

	case CircuitOpen:
		// Transition to half-open after reset timeout expires
		if time.Since(cb.lastFailureTime) >= cb.config.ResetTimeout {
			cb.state = CircuitHalfOpen
			cb.halfOpenCalls = 1 // This call counts as the first half-open probe
			cb.halfOpenFailures = 0
			return true
		}
		return false

	case CircuitHalfOpen:
		// Allow limited calls in half-open state to probe service health.
		// With HalfOpenMaxCalls=1 (default), only the transition call is allowed.
		// If it succeeds → close circuit, if it fails → reopen circuit.
		if cb.halfOpenCalls < cb.config.HalfOpenMaxCalls {
			cb.halfOpenCalls++
			return true
		}
		return false

	default:
		return false
	}
}

// RecordSuccess records a successful call.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitClosed:
		cb.failures = 0

	case CircuitHalfOpen:
		// Success in half-open: close the circuit
		cb.state = CircuitClosed
		cb.failures = 0
	}
}

// RecordFailure records a failed call.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.lastFailureTime = time.Now()

	switch cb.state {
	case CircuitClosed:
		cb.failures++
		if cb.failures >= cb.config.FailureThreshold {
			cb.state = CircuitOpen
		}

	case CircuitHalfOpen:
		// Failure in half-open: reopen the circuit
		cb.state = CircuitOpen
		cb.failures = 0
	}
}

// Reset resets the circuit breaker to closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = CircuitClosed
	cb.failures = 0
	cb.halfOpenCalls = 0
	cb.halfOpenFailures = 0
}
