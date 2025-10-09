package routes

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupAuthRoutes configures authentication-related routes
func (r *Router) SetupAuthRoutes(app *fiber.App) {
	// Create authGroup route group
	authGroup := app.Group("/auth")

	// Basic authentication endpoints
	authGroup.Post("/login", r.AuthRoutes.Login)
	authGroup.Post("/register", r.AuthRoutes.Register)
	authGroup.Post("/refresh", r.AuthRoutes.RefreshToken)
	authGroup.Post("/logout", r.AuthRoutes.Logout)

	// Protected route to get current user info
	authGroup.Get("/me", middleware.AuthMiddleware(), r.AuthRoutes.Me)

	// Google OAuth endpoints
	googleAuth := authGroup.Group("/google")

	// Protected Google OAuth endpoints (require user to be logged in first)
	googleAuth.Get("/url", middleware.AuthMiddleware(), r.AuthRoutes.GoogleAuthURL)
	googleAuth.Get("/status", middleware.AuthMiddleware(), r.AuthRoutes.GoogleLinkStatus)
	googleAuth.Get("/access-token", middleware.AuthMiddleware(), r.AuthRoutes.GoogleAccessToken)
	googleAuth.Delete("/unlink", middleware.AuthMiddleware(), r.AuthRoutes.GoogleUnlink)

	// Public Google OAuth callback (no auth required as it validates state)
	googleAuth.Get("/callback", r.AuthRoutes.GoogleAuthCallback)
}
