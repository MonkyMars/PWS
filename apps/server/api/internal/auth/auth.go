package auth

import (
	"strings"

	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// Login handles user authentication and returns JWT tokens
func (ar *AuthRoutes) Login(c fiber.Ctx) error {
	// Get validated request from context (validation middleware has already processed it)
	authRequest, err := middleware.GetValidatedRequest[types.AuthRequest](c)
	if err != nil {
		msg := fmt.Sprintf("Failed to get validated request: %v", err)
		return lib.HandleServiceError(c, err, msg)
	}

	// Attempt login using injected service
	user, err := ar.authService.Login(authRequest)
	if err != nil {
		msg := fmt.Sprintf("Login failed for email %s: %v", authRequest.Email, err)
		return lib.HandleServiceError(c, err, msg)
	}

	// Generate tokens using injected service
	accessToken, err := ar.authService.GenerateAccessToken(user)
	if err != nil {
		msg := fmt.Sprintf("Failed to generate access token for user ID %s: %v", user.Id, err)
		return lib.HandleServiceError(c, err, msg)
	}

	refreshToken, err := ar.authService.GenerateRefreshToken(user)
	if err != nil {
		msg := fmt.Sprintf("Failed to generate refresh token for user ID %s: %v", user.Id, err)
		return lib.HandleServiceError(c, err, msg)
	}

	ar.cookieService.SetAuthCookies(c, accessToken, refreshToken)

	return response.Success(c, user)
}

// Register handles user registration and returns JWT tokens
func (ar *AuthRoutes) Register(c fiber.Ctx) error {
	// Get validated request from context (validation middleware has already processed it)
	registerRequest, err := middleware.GetValidatedRequest[types.RegisterRequest](c)
	if err != nil {
		msg := fmt.Sprintf("Failed to get validated register request: %v", err)
		return lib.HandleServiceError(c, lib.ErrInvalidRequest, msg)
	}

	// Attempt registration using injected service
	user, err := ar.authService.Register(registerRequest)
	if err != nil {
		msg := fmt.Sprintf("Registration failed for email %s, username %s: %v", registerRequest.Email, registerRequest.Username, err)
		return lib.HandleServiceError(c, err, msg)
	}

	// Generate tokens for the new user using injected service
	accessToken, err := ar.authService.GenerateAccessToken(user)
	if err != nil {
		msg := fmt.Sprintf("Failed to generate access token for user ID %s: %v", user.Id, err)
		return lib.HandleServiceError(c, err, msg)
	}

	refreshToken, err := ar.authService.GenerateRefreshToken(user)
	if err != nil {
		msg := fmt.Sprintf("Failed to generate refresh token for user ID %s: %v", user.Id, err)
		return lib.HandleServiceError(c, err, msg)
	}

	ar.cookieService.SetAuthCookies(c, accessToken, refreshToken)

	return response.Success(c, user)
}

// RefreshToken handles token refresh using refresh tokens
func (ar *AuthRoutes) RefreshToken(c fiber.Ctx) error {
	token := c.Cookies(lib.RefreshTokenCookieName)

	// Refresh tokens with rotation using injected service
	authResponse, err := ar.authService.RefreshToken(token)
	if err != nil {
		// Check if this might be a token reuse attack
		if strings.Contains(err.Error(), "revoked") || strings.Contains(err.Error(), "blacklisted") {
			msg := fmt.Sprintf("Possible token reuse attack detected - client_ip: %s, user_agent: %s, error: %v", c.IP(), c.Get("User-Agent"), err)
			return lib.HandleServiceError(c, lib.ErrTokenReuse, msg)
		}
		msg := fmt.Sprintf("Token refresh failed: %v", err)
		return lib.HandleServiceError(c, err, msg)
	}

	// Set new rotated tokens in secure cookies using injected service
	ar.cookieService.SetAuthCookies(c, authResponse.AccessToken, authResponse.RefreshToken)

	return response.Success(c, authResponse)
}

// Me returns the current authenticated user's information
func (ar *AuthRoutes) Me(c fiber.Ctx) error {
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		msg := "Failed to get authenticated user claims from context"
		return lib.HandleServiceError(c, err, msg)
	}

	// Fetch user info using injected service
	user, err := ar.authService.GetUserByID(claims.Id)
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve user info for user ID %s: %v", claims.Sub, err)
		return lib.HandleServiceError(c, err, msg)
	}

	if user == nil {
		msg := fmt.Sprintf("User not found for user ID %s", claims.Sub)
		return lib.HandleServiceError(c, lib.ErrUserNotFound, msg)
	}

	return response.Success(c, user)
}

// Logout handles user logout with graceful handling of missing/invalid tokens
func (ar *AuthRoutes) Logout(c fiber.Ctx) error {
	// Extract values from context before spawning goroutine to avoid race conditions
	accessToken := c.Cookies(lib.AccessTokenCookieName)
	refreshToken := c.Cookies(lib.RefreshTokenCookieName)
	user := lib.GetUserFromContext(c)

	// Blacklist access token if present using injected service
	if strings.TrimSpace(accessToken) != "" {
		// Validate and blacklist access token
		_, err := ar.authService.GetUserFromToken(accessToken)
		if err != nil {
			lib.HandleServiceWarning(c, "Invalid access token during logout, clearing anyway", "error", err)
		} else {
			// Token is valid, blacklist it
			if err := ar.authService.BlacklistToken(accessToken, true); err != nil {
				lib.HandleServiceWarning(c, "Failed to blacklist access token during logout", "error", err)
				// Don't return error, continue with logout process
			}
		}

	// Process refresh token if present using injected service
	if strings.TrimSpace(refreshToken) != "" {
		// Try to blacklist refresh token (may be invalid, but that's okay)
		if err := ar.authService.BlacklistToken(refreshToken, false); err != nil {
			lib.HandleServiceWarning(c, "Failed to blacklist refresh token, may already be invalid", "error", err)
			// Don't return error, continue with logout process
		}

		// Clear user from cache if user exists
		if user != nil {
			err := ar.authService.ClearUserCache(user.Id)
			if err != nil {
				ar.logger.Warn("Failed to clear user cache during logout", "user_id", user.Id, "error", err)
			}
		}
	}()

	// Always clear auth cookies regardless of token validity using injected service
	ar.cookieService.ClearAuthCookies(c)

	return response.Success(c, types.LogoutResponse{
		Message: "Logged out successfully",
	})
}
