package auth

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// AuthRoutes handles HTTP routing for auth-related endpoints.
// It follows clean architecture principles by depending on interfaces rather than concrete implementations.
// This makes the code more testable and maintainable.
type AuthRoutes struct {
	authService   services.AuthServiceInterface
	cookieService services.CookieServiceInterface
	googleService services.GoogleServiceInterface
	middleware    *middleware.Middleware
}

// NewAuthRoutesWithDefaults creates an AuthRoutes instance with default dependencies.
// This is a convenience constructor for production use where you want to use
// the default implementations of all services.
func NewAuthRoutesWithDefaults() *AuthRoutes {
	return &AuthRoutes{
		authService:   services.NewAuthService(),
		cookieService: services.NewCookieService(),
		googleService: services.NewGoogleService(),
		middleware:    middleware.NewMiddleware(),
	}
}

// RegisterRoutes registers all auth-related routes with the Fiber application.
// This method organizes routes logically and follows RESTful conventions.
// It groups related functionality and applies appropriate middleware.
func (ar *AuthRoutes) RegisterRoutes(app *fiber.App) {
	// Auth API group - handles user authentication and management
	auth := app.Group("/auth")
	ar.registerAuthRoutes(auth)

	google := app.Group("/google")
	ar.registerOAuthRoutes(google)
}

// registerAuthRoutes sets up all auth-related endpoints with proper middleware and handlers
func (ar *AuthRoutes) registerAuthRoutes(router fiber.Router) {
	// Public auth endpoints with validation middleware
	router.Post("/login",
		middleware.ValidateRequest[types.AuthRequest](middleware.AuthRequestValidation),
		ar.Login,
	)
	router.Post("/register",
		middleware.ValidateRequest[types.RegisterRequest](middleware.RegisterRequestValidation),
		ar.Register,
	)
	router.Post("/refresh", ar.RefreshToken)

	// Authenticated endpoints (require valid access token)
	protected := router.Group("/", ar.middleware.AuthMiddleware())
	protected.Get("/me", ar.Me)
	protected.Post("/logout", ar.Logout)
}

func (ar *AuthRoutes) registerOAuthRoutes(router fiber.Router) {
	// OAuth endpoints:
	// Keep public routes public
	router.Get("/callback", ar.GoogleAuthCallback)
	// Protected routes
	protected := router.Group("/", ar.middleware.AuthMiddleware())
	protected.Delete("/unlink", ar.GoogleUnlink)
	protected.Get("/status", ar.GoogleLinkStatus)
	router.Get("/url", ar.GoogleAuthURL)
	router.Get("/access-token", ar.GoogleAccessToken)
}
