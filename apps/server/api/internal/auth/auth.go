package auth

import (
	"fmt"
	"strings"

	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// Login handles user authentication and returns JWT tokens
func (ar *AuthRoutes) Login(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Get validated request from context (validation middleware has already processed it)
	authRequest, err := middleware.GetValidatedRequest[types.AuthRequest](c)
	if err != nil {
		logger.Error("Failed to get validated request", "error", err)
		return lib.HandleValidationError(c, err, "request")
	}

	// Attempt login using injected service
	user, err := ar.authService.Login(authRequest)
	if err != nil {
		logger.AuditError("Login failed", "email", authRequest.Email, "error", err.Error())
		return lib.HandleAuthError(c, err, "login")
	}

	// Generate tokens using injected service
	accessToken, err := ar.authService.GenerateAccessToken(user)
	if err != nil {
		logger.AuditError("Failed to generate access token", "user_id", user.Id, "error", err)
		return lib.HandleServiceError(c, lib.ErrTokenGeneration)
	}

	refreshToken, err := ar.authService.GenerateRefreshToken(user)
	if err != nil {
		logger.AuditError("Failed to generate refresh token", "user_id", user.Id, "error", err)
		return lib.HandleServiceError(c, lib.ErrTokenGeneration)
	}

	ar.cookieService.SetAuthCookies(c, accessToken, refreshToken)

	return response.Success(c, user)
}

// Register handles user registration and returns JWT tokens
func (ar *AuthRoutes) Register(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Get validated request from context (validation middleware has already processed it)
	registerRequest, err := middleware.GetValidatedRequest[types.RegisterRequest](c)
	if err != nil {
		logger.Error("Failed to get validated request", "error", err)
		return lib.HandleValidationError(c, err, "request")
	}

	// Attempt registration using injected service
	user, err := ar.authService.Register(registerRequest)
	if err != nil {
		logger.AuditError("Registration failed", "email", registerRequest.Email, "username", registerRequest.Username, "error", err.Error())
		return lib.HandleServiceError(c, err)
	}

	// Generate tokens for the new user using injected service
	accessToken, err := ar.authService.GenerateAccessToken(user)
	if err != nil {
		logger.AuditError("Failed to generate access token", "user_id", user.Id, "error", err)
		return response.InternalServerError(c, "Failed to generate access token")
	}

	refreshToken, err := ar.authService.GenerateRefreshToken(user)
	if err != nil {
		logger.AuditError("Failed to generate refresh token", "user_id", user.Id, "error", err)
		return response.InternalServerError(c, "Failed to generate refresh token")
	}

	ar.cookieService.SetAuthCookies(c, accessToken, refreshToken)

	return response.Success(c, user)
}

// RefreshToken handles token refresh using refresh tokens
func (ar *AuthRoutes) RefreshToken(c fiber.Ctx) error {
	logger := config.SetupLogger()

	token := c.Cookies(lib.RefreshTokenCookieName)

	// Refresh tokens with rotation using injected service
	authResponse, err := ar.authService.RefreshToken(token)
	if err != nil {
		logger.Error("Token refresh failed", "error", err)

		// Check if this might be a token reuse attack
		if strings.Contains(err.Error(), "revoked") || strings.Contains(err.Error(), "blacklisted") {
			logger.Warn("Possible token reuse attack detected during refresh",
				"client_ip", c.IP(),
				"user_agent", c.Get("User-Agent"),
				"error", err.Error())
		}

		return response.Unauthorized(c, "Invalid or expired refresh token")
	}

	// Set new rotated tokens in secure cookies using injected service
	ar.cookieService.SetAuthCookies(c, authResponse.AccessToken, authResponse.RefreshToken)

	return response.Success(c, authResponse)
}

// Me returns the current authenticated user's information
func (ar *AuthRoutes) Me(c fiber.Ctx) error {
	logger := config.SetupLogger()

	claimsInterface := c.Locals("claims")

	if claimsInterface == nil {
		logger.Error("No claims found in context")
		return response.Unauthorized(c, "Unauthorized")
	}

	claims, ok := claimsInterface.(*types.AuthClaims)
	if !ok {
		logger.AuditError("Invalid claims type in context", "type", fmt.Sprintf("%T", claimsInterface))
		return response.Unauthorized(c, "Unauthorized")
	}

	// Fetch user info using injected service
	user, err := ar.authService.GetUserByID(claims.Sub)
	if err != nil {
		logger.AuditError("Failed to retrieve user info", "user_id", claims.Sub, "error", err)
		return response.InternalServerError(c, "Failed to retrieve user info")
	}

	if user == nil {
		logger.Error("User not found", "user_id", claims.Sub)
		return response.NotFound(c, "User not found")
	}

	return response.Success(c, user)
}

// Logout handles user logout with graceful handling of missing/invalid tokens
func (ar *AuthRoutes) Logout(c fiber.Ctx) error {
	logger := config.SetupLogger()

	accessToken := c.Cookies(lib.AccessTokenCookieName)
	refreshToken := c.Cookies(lib.RefreshTokenCookieName)

	// Blacklist access token if present using injected service
	if strings.TrimSpace(accessToken) != "" {
		// Validate and blacklist access token
		_, err := ar.authService.GetUserFromToken(accessToken)
		if err != nil {
			logger.Warn("Invalid access token during logout, clearing anyway", "error", err)
		} else {
			// Token is valid, blacklist it
			if err := ar.authService.BlacklistToken(accessToken, true); err != nil {
				logger.AuditError("Failed to blacklist access token", "error", err)
				// Don't return error, continue with logout process
			}
		}
	}

	// Process refresh token if present using injected service
	if strings.TrimSpace(refreshToken) != "" {
		// Try to blacklist refresh token (may be invalid, but that's okay)
		if err := ar.authService.BlacklistToken(refreshToken, false); err != nil {
			logger.Warn("Failed to blacklist refresh token, may already be invalid", "error", err)
			// Don't return error, continue with logout process
		}
	}

	// Always clear auth cookies regardless of token validity using injected service
	ar.cookieService.ClearAuthCookies(c)

	return response.Success(c, types.LogoutResponse{
		Message: "Logged out successfully",
	})
}
