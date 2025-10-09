// Package api provides the HTTP API layer for the PWS application.
// This package contains the main application setup, route configuration, and middleware
// integration. It serves as the orchestration layer that ties together configuration,
// logging, routing, and other core application components.
//
// The api package follows a modular architecture where concerns are separated into
// different sub-packages: config for application configuration, middleware for HTTP
// middleware, response for standardized API responses, and routes for endpoint handlers.
package api

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/api/routes"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/workers"
	"github.com/gofiber/fiber/v3"
)

// App initializes and starts the main application server.
// It loads configuration, sets up logging, creates the Fiber application with
// appropriate middleware, configures routes, and starts the HTTP server.
// Returns an error if the server fails to start or encounters a configuration issue.
func App() error {
	// Get centralized configuration
	cfg := config.Get()

	// Setup logger with centralized config
	logger := config.SetupLogger()

	// Create Fiber app with centralized config
	app := fiber.New(config.SetupFiber())

	// Add CORS middleware
	app.Use(middleware.SetupCORS())

	// Add logging middleware
	app.Use(logger.HTTPMiddleware())

	// Add health monitoring middleware
	app.Use(middleware.CreateHealthMiddleware())

	// Log server startup
	logger.ServerStart()

	// Setup routes
	SetupRoutes(app, logger)

	// Auto-discover routes for health monitoring
	workers.DiscoverRoutes(app)

	// Start audit logging system
	workers.StartAuditWorker()

	// Start health logging system
	workers.StartHealthLogWorker()

	// Log server ready
	logger.ServerReady()

	// Start server
	return app.Listen(cfg.GetServerAddress())
}

// SetupRoutes configures all application routes by delegating to specific route handlers.
// This function serves as the central point for route registration, making it easy to
// see all available routes and maintain route organization.
//
// Parameters:
//   - app: The Fiber application instance to register routes on
//   - logger: The centralized logger instance for route logging
func SetupRoutes(app *fiber.App, logger *config.Logger) {
	app.Get("/favicon.ico", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	router := routes.NewRouter()

	// Authentication routes
	router.SetupAuthRoutes(app)

	// Content routes
	router.SetupContentRoutes(app)

	// Health check
	router.SetupHealthRoutes(app)

	// Subject routes
	router.SetupSubjectRoutes(app)

	// Worker monitoring routes
	router.SetupWorkerRoutes(app)

	// Catch-all for undefined routes
	app.Use(func(c fiber.Ctx) error {
		return response.NotFound(c, "The requested resource was not found.")
	})
}
