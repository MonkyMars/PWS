package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MonkyMars/PWS/api"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/services"

	"github.com/joho/godotenv"
)

/*
* main is the entry point of the application
* It initializes configuration, logging, database connections,
* starts the API server with graceful shutdown handling.
* It uses centralized configuration and logging throughout.
 */
func main() {
	// Load environment variables from .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file, proceeding with system environment variables")
	}

	// Load centralized configuration
	cfg := config.Load()

	// Setup centralized logger
	logger := config.SetupLogger()
	logger.ConfigLoaded()

	// Print configuration in development mode
	if cfg.IsDevelopment() {
		cfg.PrintConfig()
	}

	logger.Info("Starting application",
		"app_name", cfg.AppName,
		"environment", cfg.Environment,
		"port", cfg.Port,
	)

	// Initialize database connection with centralized config
	err = services.InitializeDatabase()
	if err != nil {
		logger.DatabaseError("initialization", err)
		log.Fatalf("Database initialization error: %v", err)
	}

	// Log successful database connection
	logger.DatabaseConnected()

	// Test database connection
	err = services.Ping()
	if err != nil {
		logger.DatabaseError("ping", err)
		log.Fatalf("Database connection error: %v", err)
	}

	// Initialize and test Redis connection
	cacheService := &services.CacheService{}
	err = cacheService.Ping()
	if err != nil {
		logger.Error("Redis connection error", "error", err)
		log.Fatalf("Redis connection error: %v", err)
	}
	logger.Info("Redis connection established successfully")

	// Setup graceful shutdown
	setupGracefulShutdown(logger)

	// Ensure database and Redis connections are closed on exit
	defer func() {
		logger.Shutdown("application_exit")
		if err := services.CloseDatabase(); err != nil {
			logger.DatabaseError("close", err)
		}
		if err := services.CloseRedisConnection(); err != nil {
			logger.Error("Redis close error", "error", err)
		}
	}()

	// Log database statistics in development mode
	if cfg.IsDevelopment() {
		services.GetDatabaseStats()
	}

	// Start the API server
	logger.ServerStart()
	err = api.App()
	if err != nil {
		logger.ServerError(err)
		log.Fatal(err)
	}
}

// setupGracefulShutdown sets up signal handling for graceful application shutdown
func setupGracefulShutdown(logger *config.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Shutdown("signal_received")

		// Close database connection
		if err := services.CloseDatabase(); err != nil {
			logger.DatabaseError("shutdown_close", err)
		}

		// Close Redis connection
		if err := services.CloseRedisConnection(); err != nil {
			logger.Error("Redis shutdown close error", "error", err)
		}

		logger.Info("Application shutdown complete")
		os.Exit(0)
	}()
}
