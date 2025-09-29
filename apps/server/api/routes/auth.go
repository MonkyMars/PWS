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

	// Basic authentication endpoints
	auth.Post("/login", internal.Login)
	auth.Post("/register", internal.Register)
	auth.Post("/refresh", internal.RefreshToken)
	auth.Post("/logout", internal.Logout)

	// Protected route to get current user info
	auth.Get("/me", middleware.AuthMiddleware(), internal.Me)

	// Google OAuth endpoints
	googleAuth := auth.Group("/google")

	// Protected Google OAuth endpoints (require user to be logged in first)
	googleAuth.Get("/url", middleware.AuthMiddleware(), internal.GoogleAuthURL)
	googleAuth.Get("/status", middleware.AuthMiddleware(), internal.GoogleLinkStatus)
	googleAuth.Get("/access-token", middleware.AuthMiddleware(), internal.GoogleAccessToken)
	googleAuth.Delete("/unlink", middleware.AuthMiddleware(), internal.GoogleUnlink)

	// Public Google OAuth callback (no auth required as it validates state)
	googleAuth.Get("/callback", internal.GoogleAuthCallback)
}
