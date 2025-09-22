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
	"github.com/MonkyMars/PWS/types"
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
	// Database pool statistics collection disabled to reduce log noise
	// Use monitoring tools or add temporary logging here if needed for debugging
}

// GetDatabaseConfig returns the current database configuration
func GetDatabaseConfig() config.DatabaseConfig {
	return database.GetConfig()
}

func Query() *types.QueryParams {
	return types.NewQuery()
}
