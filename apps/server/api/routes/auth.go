package routes

import (
	"github.com/MonkyMars/PWS/api/internal"
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupAuthRoutes configures authentication-related routes
func SetupAuthRoutes(app *fiber.App) {
	// Create auth route group
	auth := app.Group("/auth")

	// Authentication endpoints
	auth.Post("/login", internal.Login)
	auth.Post("/register", internal.Register)
	auth.Post("/refresh", internal.RefreshToken)
	auth.Post("/logout", internal.Logout)

	// Protected route to get current user info
	auth.Get("/me", middleware.AuthMiddleware(), internal.Me)
}
