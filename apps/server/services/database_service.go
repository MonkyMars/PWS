// Package services provides database and external service integrations for the PWS application.
// This package contains database connection management, data access layers, and external
// API integrations to support the core business logic of the application.
//
// The services package is designed to be modular and testable, with interfaces that
// allow for easy mocking and dependency injection in tests.
package services

import (
	"context"
	"sync"
	"time"

	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
)

var (
	dbCircuitBreaker *lib.DatabaseCircuitBreaker
	cbOnce           sync.Once
)

// initCircuitBreaker initializes the database circuit breaker
func initCircuitBreaker() {
	cbOnce.Do(func() {
		config := lib.CircuitBreakerConfig{
			MaxFailures:      5,
			Timeout:          30 * time.Second,
			MaxRequests:      3,
			SuccessThreshold: 2,
		}
		dbCircuitBreaker = lib.NewDatabaseCircuitBreaker("database", config)
	})
}

// GetCircuitBreaker returns the database circuit breaker instance
func GetCircuitBreaker() *lib.DatabaseCircuitBreaker {
	initCircuitBreaker()
	return dbCircuitBreaker
}

// Ping tests the database connection with circuit breaker protection
func Ping() error {
	cb := GetCircuitBreaker()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return cb.ExecuteQuery(ctx, func() error {
		db := database.GetInstance()
		return db.Health()
	})
}

// PingWithContext tests the database connection with circuit breaker protection and custom context
func PingWithContext(ctx context.Context) error {
	cb := GetCircuitBreaker()
	return cb.ExecuteQuery(ctx, func() error {
		db := database.GetInstance()
		return db.Health()
	})
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	return database.CloseInstance()
}

// ExecuteWithCircuitBreaker executes a database operation with circuit breaker protection
func ExecuteWithCircuitBreaker(ctx context.Context, operation func() error) error {
	cb := GetCircuitBreaker()
	return cb.ExecuteQuery(ctx, operation)
}

// GetCircuitBreakerStats returns statistics about the database circuit breaker
func GetCircuitBreakerStats() map[string]any {
	cb := GetCircuitBreaker()
	return cb.Stats()
}

// GetDatabaseConfig returns the current database configuration
func GetDatabaseConfig() types.DatabaseConfig {
	return database.GetConfig()
}

// Query returns a new QueryParams instance for building database queries
func Query() *types.QueryParams {
	return types.NewQuery()
}

type DatabaseServiceInterface interface {
	Ping() error
	PingWithContext(ctx context.Context) error
	CloseDatabase() error
	ExecuteWithCircuitBreaker(ctx context.Context, operation func() error) error
	GetCircuitBreakerStats() map[string]any
}
