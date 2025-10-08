// Package lib provides utility components for the PWS application.
// This file implements a circuit breaker pattern to protect against cascading failures
// when database operations are failing consistently.
package lib

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// CircuitState represents the current state of the circuit breaker
type CircuitState int

const (
	// StateClosed indicates the circuit is closed and requests are allowed
	StateClosed CircuitState = iota
	// StateOpen indicates the circuit is open and requests are rejected
	StateOpen
	// StateHalfOpen indicates the circuit is testing if the service has recovered
	StateHalfOpen
)

// String returns the string representation of the circuit state
func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreakerConfig holds configuration for the circuit breaker
type CircuitBreakerConfig struct {
	// MaxFailures is the maximum number of failures before opening the circuit
	MaxFailures int64
	// Timeout is how long to wait before attempting to close the circuit
	Timeout time.Duration
	// MaxRequests is the maximum number of requests allowed in half-open state
	MaxRequests int64
	// SuccessThreshold is the number of successful requests needed to close the circuit
	SuccessThreshold int64
}

// DefaultCircuitBreakerConfig returns a default configuration
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxFailures:      5,
		Timeout:          30 * time.Second,
		MaxRequests:      3,
		SuccessThreshold: 2,
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config          CircuitBreakerConfig
	state           atomic.Int32
	failures        atomic.Int64
	requests        atomic.Int64
	successes       atomic.Int64
	lastFailureTime atomic.Int64
	lastStateChange atomic.Int64
	mu              sync.RWMutex
	onStateChange   func(from, to CircuitState)
}

// CircuitBreakerError represents an error from the circuit breaker
type CircuitBreakerError struct {
	State   CircuitState
	Message string
}

func (e *CircuitBreakerError) Error() string {
	return fmt.Sprintf("circuit breaker %s: %s", e.State.String(), e.Message)
}

// IsCircuitBreakerError checks if an error is a circuit breaker error
func IsCircuitBreakerError(err error) bool {
	var cbe *CircuitBreakerError
	return errors.As(err, &cbe)
}

// NewCircuitBreaker creates a new circuit breaker with the given configuration
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	if config.MaxFailures <= 0 {
		config.MaxFailures = DefaultCircuitBreakerConfig().MaxFailures
	}
	if config.Timeout <= 0 {
		config.Timeout = DefaultCircuitBreakerConfig().Timeout
	}
	if config.MaxRequests <= 0 {
		config.MaxRequests = DefaultCircuitBreakerConfig().MaxRequests
	}
	if config.SuccessThreshold <= 0 {
		config.SuccessThreshold = DefaultCircuitBreakerConfig().SuccessThreshold
	}

	cb := &CircuitBreaker{
		config: config,
	}
	cb.state.Store(int32(StateClosed))
	cb.lastStateChange.Store(time.Now().UnixNano())

	return cb
}

// SetOnStateChange sets a callback function that is called when the circuit state changes
func (cb *CircuitBreaker) SetOnStateChange(fn func(from, to CircuitState)) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.onStateChange = fn
}

// Execute executes the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	// Check if we can execute the request
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	// Execute the function with timeout context
	var fnErr error
	done := make(chan struct{})

	go func() {
		defer close(done)
		fnErr = fn()
	}()

	select {
	case <-done:
		// Function completed, handle the result
		cb.afterRequest(fnErr)
		return fnErr
	case <-ctx.Done():
		// Context cancelled or timed out
		cb.afterRequest(ctx.Err())
		return ctx.Err()
	}
}

// beforeRequest checks if the request can be executed based on circuit state
func (cb *CircuitBreaker) beforeRequest() error {
	state := cb.getState()
	now := time.Now()

	switch state {
	case StateClosed:
		return nil
	case StateOpen:
		// Check if enough time has passed to transition to half-open
		lastFailure := time.Unix(0, cb.lastFailureTime.Load())
		if now.Sub(lastFailure) >= cb.config.Timeout {
			cb.setState(StateHalfOpen)
			cb.requests.Store(0)
			cb.successes.Store(0)
			return nil
		}
		return &CircuitBreakerError{
			State:   StateOpen,
			Message: "circuit breaker is open",
		}
	case StateHalfOpen:
		// Check if we've exceeded the maximum requests in half-open state
		if cb.requests.Load() >= cb.config.MaxRequests {
			return &CircuitBreakerError{
				State:   StateHalfOpen,
				Message: "too many requests in half-open state",
			}
		}
		cb.requests.Add(1)
		return nil
	default:
		return &CircuitBreakerError{
			State:   state,
			Message: "unknown circuit breaker state",
		}
	}
}

// afterRequest handles the result of the request execution
func (cb *CircuitBreaker) afterRequest(err error) {
	state := cb.getState()

	if err != nil {
		cb.recordFailure()

		switch state {
		case StateClosed:
			if cb.failures.Load() >= cb.config.MaxFailures {
				cb.setState(StateOpen)
			}
		case StateHalfOpen:
			cb.setState(StateOpen)
		}
	} else {
		cb.recordSuccess()

		switch state {
		case StateHalfOpen:
			if cb.successes.Load() >= cb.config.SuccessThreshold {
				cb.setState(StateClosed)
				cb.reset()
			}
		case StateClosed:
			// Reset failure count on successful request
			cb.failures.Store(0)
		}
	}
}

// recordFailure records a failure
func (cb *CircuitBreaker) recordFailure() {
	cb.failures.Add(1)
	cb.lastFailureTime.Store(time.Now().UnixNano())
}

// recordSuccess records a success
func (cb *CircuitBreaker) recordSuccess() {
	cb.successes.Add(1)
}

// reset resets the circuit breaker counters
func (cb *CircuitBreaker) reset() {
	cb.failures.Store(0)
	cb.requests.Store(0)
	cb.successes.Store(0)
}

// getState returns the current circuit state
func (cb *CircuitBreaker) getState() CircuitState {
	return CircuitState(cb.state.Load())
}

// setState changes the circuit state and calls the state change callback
func (cb *CircuitBreaker) setState(newState CircuitState) {
	oldState := CircuitState(cb.state.Swap(int32(newState)))
	cb.lastStateChange.Store(time.Now().UnixNano())

	// Call state change callback if set
	if cb.onStateChange != nil && oldState != newState {
		go cb.onStateChange(oldState, newState)
	}
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() CircuitState {
	return cb.getState()
}

// Failures returns the current failure count
func (cb *CircuitBreaker) Failures() int64 {
	return cb.failures.Load()
}

// Requests returns the current request count (relevant in half-open state)
func (cb *CircuitBreaker) Requests() int64 {
	return cb.requests.Load()
}

// Successes returns the current success count (relevant in half-open state)
func (cb *CircuitBreaker) Successes() int64 {
	return cb.successes.Load()
}

// LastFailureTime returns the time of the last failure
func (cb *CircuitBreaker) LastFailureTime() time.Time {
	return time.Unix(0, cb.lastFailureTime.Load())
}

// LastStateChange returns the time of the last state change
func (cb *CircuitBreaker) LastStateChange() time.Time {
	return time.Unix(0, cb.lastStateChange.Load())
}

// Stats returns comprehensive statistics about the circuit breaker
func (cb *CircuitBreaker) Stats() map[string]any {
	return map[string]any{
		"state":             cb.State().String(),
		"failures":          cb.Failures(),
		"requests":          cb.Requests(),
		"successes":         cb.Successes(),
		"last_failure_time": cb.LastFailureTime(),
		"last_state_change": cb.LastStateChange(),
		"max_failures":      cb.config.MaxFailures,
		"timeout":           cb.config.Timeout.String(),
		"max_requests":      cb.config.MaxRequests,
		"success_threshold": cb.config.SuccessThreshold,
	}
}

// ForceOpen forces the circuit breaker to open state
func (cb *CircuitBreaker) ForceOpen() {
	cb.setState(StateOpen)
	cb.lastFailureTime.Store(time.Now().UnixNano())
}

// ForceClose forces the circuit breaker to closed state and resets counters
func (cb *CircuitBreaker) ForceClose() {
	cb.setState(StateClosed)
	cb.reset()
}

// ForceClosed is an alias for ForceClose for clarity
func (cb *CircuitBreaker) ForceClosed() {
	cb.ForceClose()
}

// IsOpen returns true if the circuit breaker is in open state
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.getState() == StateOpen
}

// IsClosed returns true if the circuit breaker is in closed state
func (cb *CircuitBreaker) IsClosed() bool {
	return cb.getState() == StateClosed
}

// IsHalfOpen returns true if the circuit breaker is in half-open state
func (cb *CircuitBreaker) IsHalfOpen() bool {
	return cb.getState() == StateHalfOpen
}

// DatabaseCircuitBreaker is a specialized circuit breaker for database operations
type DatabaseCircuitBreaker struct {
	*CircuitBreaker
	name string
}

// NewDatabaseCircuitBreaker creates a new circuit breaker specifically for database operations
func NewDatabaseCircuitBreaker(name string, config CircuitBreakerConfig) *DatabaseCircuitBreaker {
	cb := NewCircuitBreaker(config)

	dbCb := &DatabaseCircuitBreaker{
		CircuitBreaker: cb,
		name:           name,
	}

	// Set up logging for state changes
	cb.SetOnStateChange(func(from, to CircuitState) {
		// Note: We avoid importing config/logger here to prevent circular dependencies
		// The actual logging should be set up by the caller
	})

	return dbCb
}

// Name returns the name of the database circuit breaker
func (dcb *DatabaseCircuitBreaker) Name() string {
	return dcb.name
}

// ExecuteQuery executes a database query with circuit breaker protection
func (dcb *DatabaseCircuitBreaker) ExecuteQuery(ctx context.Context, queryFn func() error) error {
	return dcb.Execute(ctx, queryFn)
}
