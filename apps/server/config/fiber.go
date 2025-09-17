package config

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

// SetupFiber creates and returns a Fiber configuration based on centralized application settings.
// This function configures the Fiber web framework with appropriate defaults and
// application-specific settings for optimal performance and functionality.
//
// The configuration includes:
//   - Case-sensitive routing for precise URL matching
//   - Strict routing to prevent trailing slash ambiguity
//   - Custom JSON encoder/decoder for improved performance
//   - Application-specific headers and naming
//   - Environment-based error handling
//
// Returns a Fiber configuration struct ready to be used when creating a new Fiber app.
func SetupFiber() fiber.Config {
	cfg := Get()

	return fiber.Config{
		CaseSensitive:    true,
		StrictRouting:    true,
		AppName:          cfg.AppName,
		JSONEncoder:      json.Marshal,
		JSONDecoder:      json.Unmarshal,
		ServerHeader:     cfg.AppName,
		ReadTimeout:      cfg.Server.ReadTimeout,
		WriteTimeout:     cfg.Server.WriteTimeout,
		IdleTimeout:      cfg.Server.IdleTimeout,
		ErrorHandler:     setupErrorHandler(cfg),
		DisableKeepalive: false,
	}
}

// setupErrorHandler creates a custom error handler based on environment
func setupErrorHandler(cfg *Config) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		// Default to 500 server error
		code := fiber.StatusInternalServerError

		// Check if it's a Fiber error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		// In development, return detailed error information
		if cfg.IsDevelopment() {
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
				"code":    code,
			})
		}

		// In production, return generic error messages
		var message string
		switch code {
		case fiber.StatusNotFound:
			message = "Resource not found"
		case fiber.StatusUnauthorized:
			message = "Unauthorized access"
		case fiber.StatusForbidden:
			message = "Access forbidden"
		case fiber.StatusBadRequest:
			message = "Bad request"
		default:
			message = "Internal server error"
		}

		return c.Status(code).JSON(fiber.Map{
			"error":   true,
			"message": message,
			"code":    code,
		})
	}
}
