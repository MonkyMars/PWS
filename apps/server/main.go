package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MonkyMars/PWS/api"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/workers"

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

	// Initialize worker manager with dependency injection
	workerManager := workers.NewWorkerManager(cfg, logger)

	// Initialize audit logging with the new manager
	initializeAuditLogging(workerManager)

	// Start all background workers
	if err := workerManager.Start(); err != nil {
		logger.Error("Failed to start worker manager", "error", err)
		log.Fatalf("Worker manager start error: %v", err)
	}

	// Minimal config info in development mode
	if cfg.IsDevelopment() {
		log.Printf("Development mode - %s:%s", cfg.AppName, cfg.Port)
	}

	logger.Info("Starting application",
		"app_name", cfg.AppName,
		"environment", cfg.Environment,
		"port", cfg.Port,
	)

	// Initialize database connection with centralized config
	err = database.Initialize()
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
	err = services.NewCacheService().Ping()
	if err != nil {
		logger.AuditError("Redis connection error", "error", err)
		log.Fatalf("Redis connection error: %v", err)
	}

	// Setup graceful shutdown with coordinated worker manager
	setupGracefulShutdown(logger, workerManager)

	// Ensure database and Redis connections are closed on exit
	defer func() {
		logger.Shutdown("application_exit")

		// Stop worker manager with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := workerManager.Stop(shutdownCtx); err != nil {
			logger.Error("Worker manager shutdown error", "error", err)
		}

		if err := services.CloseDatabase(); err != nil {
			logger.DatabaseError("close", err)
		}
		if err := services.CloseRedisConnection(); err != nil {
			logger.AuditError("Redis close error", "error", err)
		}
	}()

	// Start the API server
	err = api.App()
	if err != nil {
		logger.ServerError(err)
		// Fatal here to ensure the application exits if the server fails to start
		log.Fatal(err)
	}
}

// setupGracefulShutdown sets up signal handling for graceful application shutdown
func setupGracefulShutdown(logger *config.Logger, workerManager *workers.WorkerManager) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Shutdown("signal_received")

		// Stop worker manager with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := workerManager.Stop(shutdownCtx); err != nil {
			logger.Error("Worker manager shutdown error", "error", err)
		}

		// Close database connection
		if err := services.CloseDatabase(); err != nil {
			logger.DatabaseError("shutdown_close", err)
		}

		// Close Redis connection
		if err := services.CloseRedisConnection(); err != nil {
			logger.AuditError("Redis shutdown close error", "error", err)
		}

		os.Exit(0)
	}()
}

func initializeAuditLogging(workerManager *workers.WorkerManager) {
	// Wire up the audit logging function to avoid circular dependencies
	config.SetAuditLogFunc(workerManager.AddAuditLog)
}
