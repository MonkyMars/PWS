package middleware

import (
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/services"
)

// Middleware handles HTTP routing for content-related endpoints.
// It follows clean architecture principles by depending on interfaces rather than concrete implementations.
// This makes the code more testable and maintainable.
type Middleware struct {
	authService  services.AuthServiceInterface
	cacheService *services.CacheService
	logger       *config.Logger
}

// NewMiddleware creates a Middleware instance with default dependencies.
// This is a convenience constructor for production use where you want to use
// the default implementations of all services.
func NewMiddleware() *Middleware {
	return &Middleware{
		authService:  services.NewAuthService(),
		cacheService: services.NewCacheService(),
		logger:       config.SetupLogger(),
	}
}
