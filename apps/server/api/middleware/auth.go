package middleware

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		logger := config.SetupLogger()

		token := c.Cookies(lib.AccessTokenCookieName)

		if token == "" {
			logger.Error("No access token found in cookies")
			return response.Unauthorized(c, "Missing access token")
		}

		authService := services.NewAuthService()

		claims, err := authService.ParseToken(token, true)
		if err != nil {
			logger.AuditError("Failed to parse access token", "error", err)
			return response.Unauthorized(c, "Invalid or expired access token")
		}

		// Initialize Cache service
		cacheService := services.NewCacheService()

		// Check if token is blacklisted with graceful Redis failure handling
		blacklisted, err := cacheService.IsTokenBlacklisted(claims.Jti.String())
		if err != nil {
			logger.AuditError("Redis blacklist check failed, denying request for security", "error", err, "jti", claims.Jti.String())
			// Do not return faulty Redis errors to the client, let the request through if Redis is down
		} else if blacklisted {
			// SECURITY: This could indicate a token reuse attack
			logger.Warn("Blacklisted token access attempt detected",
				"jti", claims.Jti.String(),
				"user_id", claims.Sub,
				"user_email", claims.Email,
				"client_ip", c.IP(),
				"user_agent", c.Get("User-Agent"))

			// TODO: Invalidate all tokens for this user as a precaution

			return response.Unauthorized(c, "Token has been revoked")
		}

		// Store user claims in context locals for downstream handlers
		c.Locals("claims", claims)

		return c.Next()
	}
}

func AdminMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		logger := config.SetupLogger()

		token := c.Cookies(lib.AccessTokenCookieName)

		if token == "" {
			logger.Error("No access token found in cookies")
			return response.Unauthorized(c, "Missing access token")
		}

		authService := services.NewAuthService()

		claims, err := authService.ParseToken(token, true)
		if err != nil {
			logger.AuditError("Failed to parse access token", "error", err)
			return response.Unauthorized(c, "Invalid or expired access token")
		}

		// Initialize Cache service
		cacheService := services.NewCacheService()

		// Check if token is blacklisted with graceful Redis failure handling
		blacklisted, err := cacheService.IsTokenBlacklisted(claims.Jti.String())
		if err != nil {
			logger.AuditError("Redis blacklist check failed, denying request for security", "error", err, "jti", claims.Jti.String())
			// Do not return faulty Redis errors to the client, let the request through if Redis is down
		} else if blacklisted {
			// SECURITY: This could indicate a token reuse attack
			logger.Warn("Blacklisted token access attempt detected",
				"jti", claims.Jti.String(),
				"user_id", claims.Sub,
				"user_email", claims.Email,
				"client_ip", c.IP(),
				"user_agent", c.Get("User-Agent"))

			// TODO: Invalidate all tokens for this user as a precaution

			return response.Unauthorized(c, "Token has been revoked")
		}

		if claims.Role != lib.RoleAdmin {
			logger.Warn("Unauthorized admin access attempt detected",
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
