// Package services provides database and external service integrations for the PWS application.
// This package contains database connection management, data access layers, and external
// API integrations to support the core business logic of the application.
//
// The services package is designed to be modular and testable, with interfaces that
// allow for easy mocking and dependency injection in tests.
package services

import (
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
)

// Ping tests the database connection
func Ping() error {
	db := database.GetInstance()
	return db.Health()
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	return database.CloseInstance()
}

// GetDatabaseStats returns connection pool statistics and logs them
func GetDatabaseStats() {
	cfg := config.Get()
	db := database.GetInstance()
	stats := db.GetStats()

	// Use centralized logger if available, otherwise fall back to standard log
	if cfg.IsDevelopment() {
		logger := config.SetupLogger()
		logger.Performance("database_pool_stats", 0)
		logger.Info("Database Pool Statistics",
			"total_connections", stats.TotalConns,
			"idle_connections", stats.IdleConns,
			"stale_connections", stats.StaleConns,
			"hits", stats.Hits,
			"misses", stats.Misses,
			"timeouts", stats.Timeouts,
		)
	}
}

// GetDatabaseConfig returns the current database configuration
func GetDatabaseConfig() config.DatabaseConfig {
	return database.GetConfig()
}
