// Package middleware provides HTTP middleware functions for the PWS application.
// This package contains reusable middleware components for cross-cutting concerns
// such as CORS handling, authentication, logging, rate limiting, and request validation.
//
// Middleware functions in this package follow the Fiber middleware pattern and can be
// easily integrated into the application's HTTP request processing pipeline.
package middleware

import (
	"github.com/MonkyMars/PWS/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func (mw *Middleware) SetupCORS() fiber.Handler {
	cfg := config.Get()
	return cors.New(cors.Config{
		AllowOrigins:     cfg.Cors.AllowOrigins,
		AllowMethods:     cfg.Cors.AllowMethods,
		AllowHeaders:     cfg.Cors.AllowHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
	})
}
