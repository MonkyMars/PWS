package middleware

import (
	"fmt"
	"slices"

	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
)

func (mw *Middleware) AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies(lib.AccessTokenCookieName)

		if token == "" {
			msg := "No access token found in cookies during authentication middleware"
			return lib.HandleServiceError(c, lib.ErrInvalidToken, msg)
		}

		claims, err := mw.authService.ParseToken(token, true)
		if err != nil {
			msg := fmt.Sprintf("Failed to parse access token in authentication middleware: %v", err)
			return lib.HandleServiceError(c, err, msg)
		}

		// Initialize Cache service
		cacheService := services.NewCacheService()

		// Check if token is blacklisted with graceful Redis failure handling
		blacklisted, err := cacheService.IsTokenBlacklisted(claims.Jti)
		if err != nil {
			lib.HandleServiceWarning(c, "Redis blacklist check failed in auth middleware", "error", err, "jti", claims.Jti.String())
			// Do not return faulty Redis errors to the client, let the request through if Redis is down
		} else if blacklisted {
			// SECURITY: This could indicate a token reuse attack
			msg := fmt.Sprintf("Blacklisted token access attempt - jti: %s, user_id: %s, user_email: %s, client_ip: %s, user_agent: %s",
				claims.Jti.String(), claims.Sub, claims.Email, c.IP(), c.Get("User-Agent"))
			// TODO: Invalidate all tokens for this user as a precaution
			return lib.HandleServiceError(c, lib.ErrTokenRevoked, msg)
		}

		// Store user claims in context locals for downstream handlers
		c.Locals("claims", claims)

		return c.Next()
	}
}

func (mw *Middleware) AdminMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies(lib.AccessTokenCookieName)

		if token == "" {
			msg := "No access token found in cookies during admin middleware authentication"
			return lib.HandleServiceError(c, lib.ErrInvalidToken, msg)
		}

		claims, err := mw.authService.ParseToken(token, true)
		if err != nil {
			msg := fmt.Sprintf("Failed to parse access token in admin middleware: %v", err)
			return lib.HandleServiceError(c, err, msg)
		}

		// Check if token is blacklisted with graceful Redis failure handling
		blacklisted, err := mw.cacheService.IsTokenBlacklisted(claims.Jti)
		if err != nil {
			lib.HandleServiceWarning(c, "Redis blacklist check failed in admin middleware", "error", err, "jti", claims.Jti.String())
			// Do not return faulty Redis errors to the client, let the request through if Redis is down
		} else if blacklisted {
			// SECURITY: This could indicate a token reuse attack
			msg := fmt.Sprintf("Blacklisted token access attempt in admin middleware - jti: %s, user_id: %s, user_email: %s, client_ip: %s, user_agent: %s",
				claims.Jti.String(), claims.Sub, claims.Email, c.IP(), c.Get("User-Agent"))
			// TODO: Invalidate all tokens for this user as a precaution
			return lib.HandleServiceError(c, lib.ErrTokenRevoked, msg)
		}

		if claims.Role != lib.RoleAdmin {
			msg := fmt.Sprintf("Unauthorized admin access attempt - user_id: %s, user_email: %s, user_role: %s, client_ip: %s, user_agent: %s",
				claims.Sub, claims.Email, claims.Role, c.IP(), c.Get("User-Agent"))
			return lib.HandleServiceError(c, lib.ErrInsufficientPermissions, msg)
		}

		// Store user claims in context locals for downstream handlers
		c.Locals("claims", claims)

		return c.Next()
	}
}

func (mw *Middleware) RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, err := lib.GetValidatedClaims(c)
		if err != nil {
			return lib.HandleServiceError(c, err, "Failed to get validated claims in RoleMiddleware")
		}

		isAllowed := slices.Contains(allowedRoles, claims.Role)

		if !isAllowed {
			msg := fmt.Sprintf("Insufficient permissions. User with role '%s' tried to access a route that requires one of '%v'", claims.Role, allowedRoles)
			return lib.HandleServiceError(c, lib.ErrInsufficientPermissions, msg)
		}

		return c.Next()
	}
}
