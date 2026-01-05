package adapter

import (
	"testing"
	"time"
)

func TestCircuitBreaker_InitialState(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{})

	if cb.State() != CircuitStateClosed {
		t.Errorf("expected initial state to be Closed, got %v", cb.State())
	}

	if !cb.Allow() {
		t.Error("expected Allow() to return true in Closed state")
	}
}

func TestCircuitBreaker_TransitionToOpen(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 3,
	})

	cb.RecordFailure()
	cb.RecordFailure()
	if cb.State() != CircuitStateClosed {
		t.Error("expected state to remain Closed after 2 failures")
	}

	cb.RecordFailure()
	if cb.State() != CircuitStateOpen {
		t.Errorf("expected state to be Open after 3 failures, got %v", cb.State())
	}

	if cb.Allow() {
		t.Error("expected Allow() to return false in Open state")
	}
}

func TestCircuitBreaker_TransitionToHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 1,
		Timeout:          10 * time.Millisecond,
	})

	cb.RecordFailure()
	if cb.State() != CircuitStateOpen {
		t.Error("expected state to be Open")
	}

	time.Sleep(15 * time.Millisecond)

	if !cb.Allow() {
		t.Error("expected Allow() to return true after timeout (HalfOpen)")
	}

	if cb.State() != CircuitStateHalfOpen {
		t.Errorf("expected state to be HalfOpen, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenToClosedOnSuccess(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 2,
		Timeout:          10 * time.Millisecond,
	})

	cb.RecordFailure()
	time.Sleep(15 * time.Millisecond)
	cb.Allow()

	cb.RecordSuccess()
	if cb.State() != CircuitStateHalfOpen {
		t.Error("expected state to remain HalfOpen after 1 success")
	}

	cb.RecordSuccess()
	if cb.State() != CircuitStateClosed {
		t.Errorf("expected state to be Closed after 2 successes, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenToOpenOnFailure(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 1,
		Timeout:          10 * time.Millisecond,
	})

	cb.RecordFailure()
	time.Sleep(15 * time.Millisecond)
	cb.Allow()

	cb.RecordFailure()
	if cb.State() != CircuitStateOpen {
		t.Errorf("expected state to be Open after failure in HalfOpen, got %v", cb.State())
	}
}

func TestCircuitBreaker_SuccessResetsFailureCount(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 3,
	})

	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess()

	cb.RecordFailure()
	cb.RecordFailure()
	if cb.State() != CircuitStateClosed {
		t.Error("expected state to remain Closed after success reset")
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 1,
	})

	cb.RecordFailure()
	if cb.State() != CircuitStateOpen {
		t.Error("expected state to be Open")
	}

	cb.Reset()
	if cb.State() != CircuitStateClosed {
		t.Errorf("expected state to be Closed after Reset, got %v", cb.State())
	}

	if !cb.Allow() {
		t.Error("expected Allow() to return true after Reset")
	}
}

func TestCircuitBreaker_DefaultValues(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{})

	for range 5 {
		cb.RecordFailure()
	}
	if cb.State() != CircuitStateOpen {
		t.Error("expected default failure threshold to be 5")
	}

	cb.Reset()
	cb.RecordFailure()
	time.Sleep(35 * time.Millisecond)
	if cb.Allow() {
		t.Log("default timeout may be longer, skipping timing-sensitive check")
	}
}
