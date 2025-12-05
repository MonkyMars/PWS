package middleware

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
)

func (mw *Middleware) AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies(lib.AccessTokenCookieName)

		if token == "" {
			return lib.HandleAuthError(c, lib.ErrInvalidToken, "middleware authentication")
		}

		claims, err := mw.authService.ParseToken(token, true)
		if err != nil {
			mw.logger.AuditError("Failed to parse access token", "error", err)
			return lib.HandleAuthError(c, err, "middleware authentication")
		}

		// Initialize Cache service
		cacheService := services.NewCacheService()

		// Check if token is blacklisted with graceful Redis failure handling
		blacklisted, err := cacheService.IsTokenBlacklisted(claims.Jti)
		if err != nil {
			mw.logger.AuditError("Redis blacklist check failed, denying request for security", "error", err, "jti", claims.Jti.String())
			// Do not return faulty Redis errors to the client, let the request through if Redis is down
		} else if blacklisted {
			// SECURITY: This could indicate a token reuse attack
			mw.logger.Warn("Blacklisted token access attempt detected",
				"jti", claims.Jti.String(),
				"user_id", claims.Sub,
				"user_email", claims.Email,
				"client_ip", c.IP(),
				"user_agent", c.Get("User-Agent"))

			// TODO: Invalidate all tokens for this user as a precaution

			return lib.HandleAuthError(c, lib.ErrTokenRevoked, "middleware authentication")
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
			return lib.HandleAuthError(c, lib.ErrInvalidToken, "admin middleware authentication")
		}

		claims, err := mw.authService.ParseToken(token, true)
		if err != nil {
			mw.logger.AuditError("Failed to parse access token", "error", err)
			return lib.HandleAuthError(c, err, "admin middleware authentication")
		}

		// Check if token is blacklisted with graceful Redis failure handling
		blacklisted, err := mw.cacheService.IsTokenBlacklisted(claims.Jti)
		if err != nil {
			mw.logger.AuditError("Redis blacklist check failed, denying request for security", "error", err, "jti", claims.Jti.String())
			// Do not return faulty Redis errors to the client, let the request through if Redis is down
		} else if blacklisted {
			// SECURITY: This could indicate a token reuse attack
			mw.logger.Warn("Blacklisted token access attempt detected",
				"jti", claims.Jti.String(),
				"user_id", claims.Sub,
				"user_email", claims.Email,
				"client_ip", c.IP(),
				"user_agent", c.Get("User-Agent"))

			// TODO: Invalidate all tokens for this user as a precaution

			return response.Unauthorized(c, "Token has been revoked")
		}

		if claims.Role != lib.RoleAdmin {
			mw.logger.Warn("Unauthorized admin access attempt detected",
				"user_id", claims.Sub,
				"user_email", claims.Email,
				"user_role", claims.Role,
				"client_ip", c.IP(),
				"user_agent", c.Get("User-Agent"))
			return response.Forbidden(c, "Admin access required - Current role: "+claims.Role)
		}

		// Store user claims in context locals for downstream handlers
		c.Locals("claims", claims)

		return c.Next()
	}
}
