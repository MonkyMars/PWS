package middleware

import (
	"strings"
	"time"

	"github.com/MonkyMars/PWS/workers"
	"github.com/gofiber/fiber/v3"
)

// CreateHealthMiddleware returns a middleware that tracks route metrics
func (mw *Middleware) CreateHealthMiddleware() fiber.Handler {
	manager := workers.GetGlobalManager()
	if manager == nil {
		// Return no-op middleware if health logging is disabled
		return func(c fiber.Ctx) error {
			return c.Next()
		}
	}

	return func(c fiber.Ctx) error {
		start := time.Now()

		// Extract service name from path
		serviceName := extractBasePath(c.Path())
		if serviceName == "" {
			return c.Next()
		}

		// Record metrics using the worker manager
		err := c.Next()

		// Record metrics
		latency := time.Since(start)
		statusCode := c.Response().StatusCode()
		manager.RecordHealthMetric(serviceName, statusCode, latency)

		return err
	}
}

// extractBasePath extracts the base path from a route path
func extractBasePath(path string) string {
	// Remove leading slash and split by slash
	if len(path) <= 1 {
		return ""
	}

	trimmed := path[1:] // Remove leading slash
	segments := strings.Split(trimmed, "/")
	if len(segments) == 0 {
		return ""
	}

	basePath := segments[0]

	// Skip health and system routes
	if basePath == "health" || basePath == "metrics" || basePath == "logs" {
		return ""
	}

	// Handle parameterized routes (remove :param)
	if strings.HasPrefix(basePath, ":") {
		return ""
	}

	return basePath
}
