package adapter

import (
	"sync"
	"time"
)

type CircuitState int

const (
	CircuitStateClosed CircuitState = iota
	CircuitStateOpen
	CircuitStateHalfOpen
)

type CircuitBreakerConfig struct {
	FailureThreshold int
	SuccessThreshold int
	Timeout          time.Duration
}

type CircuitBreaker struct {
	mu              sync.RWMutex
	state           CircuitState
	failureCount    int
	successCount    int
	lastFailureTime time.Time

	failureThreshold int
	successThreshold int
	timeout          time.Duration
}

func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	failureThreshold := cfg.FailureThreshold
	if failureThreshold <= 0 {
		failureThreshold = 5
	}

	successThreshold := cfg.SuccessThreshold
	if successThreshold <= 0 {
		successThreshold = 2
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &CircuitBreaker{
		state:            CircuitStateClosed,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		timeout:          timeout,
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitStateClosed:
		return true

	case CircuitStateOpen:
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = CircuitStateHalfOpen
			cb.successCount = 0
			return true
		}
		return false

	case CircuitStateHalfOpen:
		return true
	}

	return false
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitStateClosed:
		cb.failureCount = 0

	case CircuitStateHalfOpen:
		cb.successCount++
		if cb.successCount >= cb.successThreshold {
			cb.state = CircuitStateClosed
			cb.failureCount = 0
			cb.successCount = 0
		}
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.lastFailureTime = time.Now()

	switch cb.state {
	case CircuitStateClosed:
		cb.failureCount++
		if cb.failureCount >= cb.failureThreshold {
			cb.state = CircuitStateOpen
		}

	case CircuitStateHalfOpen:
		cb.state = CircuitStateOpen
		cb.successCount = 0
	}
}

func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = CircuitStateClosed
	cb.failureCount = 0
	cb.successCount = 0
}
