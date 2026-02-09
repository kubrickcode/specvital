package reliability

import (
	"testing"
	"time"
)

func TestCircuitBreaker_InitialState(t *testing.T) {
	cb := NewCircuitBreaker(DefaultPhase1CircuitConfig())

	if cb.State() != CircuitClosed {
		t.Errorf("expected initial state to be Closed, got %v", cb.State())
	}
}

func TestCircuitBreaker_AllowWhenClosed(t *testing.T) {
	cb := NewCircuitBreaker(DefaultPhase1CircuitConfig())

	for i := 0; i < 10; i++ {
		if !cb.Allow() {
			t.Errorf("expected Allow() to return true when circuit is closed")
		}
	}
}

func TestCircuitBreaker_TripsAfterConsecutiveFailures(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 3,
		ResetTimeout:     1 * time.Second,
		HalfOpenMaxCalls: 1,
	}
	cb := NewCircuitBreaker(config)

	// Record failures up to threshold
	for i := 0; i < 3; i++ {
		if !cb.Allow() {
			t.Errorf("expected Allow() to return true before threshold")
		}
		cb.RecordFailure()
	}

	// Circuit should now be open
	if cb.State() != CircuitOpen {
		t.Errorf("expected state to be Open after %d failures, got %v", 3, cb.State())
	}

	// Allow should return false when open
	if cb.Allow() {
		t.Error("expected Allow() to return false when circuit is open")
	}
}

func TestCircuitBreaker_SuccessResetsFailureCount(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 3,
		ResetTimeout:     1 * time.Second,
		HalfOpenMaxCalls: 1,
	}
	cb := NewCircuitBreaker(config)

	// Record 2 failures (below threshold)
	cb.RecordFailure()
	cb.RecordFailure()

	// Success should reset
	cb.RecordSuccess()

	// Record 2 more failures
	cb.RecordFailure()
	cb.RecordFailure()

	// Should still be closed (total 4 failures but reset in between)
	if cb.State() != CircuitClosed {
		t.Errorf("expected state to be Closed after success reset, got %v", cb.State())
	}
}

func TestCircuitBreaker_TransitionsToHalfOpenAfterTimeout(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 2,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
	cb := NewCircuitBreaker(config)

	// Trip the circuit
	cb.RecordFailure()
	cb.RecordFailure()

	if cb.State() != CircuitOpen {
		t.Errorf("expected state to be Open, got %v", cb.State())
	}

	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)

	// Allow should transition to half-open
	if !cb.Allow() {
		t.Error("expected Allow() to return true after timeout")
	}

	if cb.State() != CircuitHalfOpen {
		t.Errorf("expected state to be HalfOpen, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenSuccessCloses(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 2,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
	cb := NewCircuitBreaker(config)

	// Trip the circuit
	cb.RecordFailure()
	cb.RecordFailure()

	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)

	// Transition to half-open
	cb.Allow()

	// Success in half-open should close
	cb.RecordSuccess()

	if cb.State() != CircuitClosed {
		t.Errorf("expected state to be Closed after half-open success, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenFailureOpens(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 2,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
	cb := NewCircuitBreaker(config)

	// Trip the circuit
	cb.RecordFailure()
	cb.RecordFailure()

	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)

	// Transition to half-open
	cb.Allow()

	// Failure in half-open should open
	cb.RecordFailure()

	if cb.State() != CircuitOpen {
		t.Errorf("expected state to be Open after half-open failure, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenLimitsCallCount(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 2,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
	cb := NewCircuitBreaker(config)

	// Trip the circuit
	cb.RecordFailure()
	cb.RecordFailure()

	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)

	// First call should succeed (transitions to half-open)
	if !cb.Allow() {
		t.Error("expected first Allow() to return true")
	}

	// Second call should fail (half-open limit reached)
	if cb.Allow() {
		t.Error("expected second Allow() to return false in half-open")
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 2,
		ResetTimeout:     1 * time.Second,
		HalfOpenMaxCalls: 1,
	}
	cb := NewCircuitBreaker(config)

	// Trip the circuit
	cb.RecordFailure()
	cb.RecordFailure()

	if cb.State() != CircuitOpen {
		t.Errorf("expected state to be Open, got %v", cb.State())
	}

	// Reset
	cb.Reset()

	if cb.State() != CircuitClosed {
		t.Errorf("expected state to be Closed after reset, got %v", cb.State())
	}

	if !cb.Allow() {
		t.Error("expected Allow() to return true after reset")
	}
}

func TestCircuitState_String(t *testing.T) {
	tests := []struct {
		state    CircuitState
		expected string
	}{
		{CircuitClosed, "closed"},
		{CircuitOpen, "open"},
		{CircuitHalfOpen, "half-open"},
		{CircuitState(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.expected {
			t.Errorf("CircuitState(%d).String() = %q, want %q", tt.state, got, tt.expected)
		}
	}
}
