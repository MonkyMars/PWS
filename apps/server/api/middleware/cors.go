// Package middleware provides HTTP middleware functions for the PWS application.
// This package contains reusable middleware components for cross-cutting concerns
// such as CORS handling, authentication, logging, rate limiting, and request validation.
//
// Middleware functions in this package follow the Fiber middleware pattern and can be
// easily integrated into the application's HTTP request processing pipeline.
package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func SetupCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})
}
