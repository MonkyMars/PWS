package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// Login handles user authentication and returns JWT tokens
func Login(c fiber.Ctx) error {
	logger := config.SetupLogger()

	var authRequest types.AuthRequest
	if err := c.Bind().Body(&authRequest); err != nil {
		logger.Error("Failed to parse login request", "error", err)
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate email
	if strings.TrimSpace(authRequest.Email) == "" {
		return response.SendValidationError(c, []types.ValidationError{
			{
				Field:   "email",
				Message: "Email is required",
				Value:   authRequest.Email,
			},
		})
	}

	// Validate password
	if strings.TrimSpace(authRequest.Password) == "" {
		return response.SendValidationError(c, []types.ValidationError{
			{
				Field:   "password",
				Message: "Password is required",
				Value:   authRequest.Password,
			},
		})
	}

	// Initialize auth service
	authService := &services.AuthService{Logger: logger}

	// Initialize cookie service
	cookieService := &services.CookieService{}

	// Attempt login
	user, err := authService.Login(&authRequest)
	if err != nil {
		logger.Error("Login failed", "email", authRequest.Email, "error", err)

		if errors.Is(err, lib.ErrInvalidCredentials) {
			return response.Unauthorized(c, "Invalid email or password")
		}

		return response.InternalServerError(c, "An error occurred during login")
	}

	// Generate tokens
	accessToken, err := authService.GenerateAccessToken(user)
	if err != nil {
		logger.AuditError("Failed to generate access token", "user_id", user.Id, "error", err)
		return response.InternalServerError(c, "Failed to generate access token")
	}

	refreshToken, err := authService.GenerateRefreshToken(user)
	if err != nil {
		logger.AuditError("Failed to generate refresh token", "user_id", user.Id, "error", err)
		return response.InternalServerError(c, "Failed to generate refresh token")
	}

	cookieService.SetAuthCookies(c, accessToken, refreshToken)

	return response.Success(c, user)
}

// Register handles user registration and returns JWT tokens
func Register(c fiber.Ctx) error {
	logger := config.SetupLogger()

	var registerRequest types.RegisterRequest
	if err := c.Bind().Body(&registerRequest); err != nil {
		logger.Error("Failed to parse register request", "error", err)
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate username
	if strings.TrimSpace(registerRequest.Username) == "" {
		return response.SendValidationError(c, []types.ValidationError{
			{
				Field:   "username",
				Message: "Username is required",
				Value:   registerRequest.Username,
			},
		})
	}

	// Validate email
	if strings.TrimSpace(registerRequest.Email) == "" {
		return response.SendValidationError(c, []types.ValidationError{
			{
				Field:   "email",
				Message: "Email is required",
				Value:   registerRequest.Email,
			},
		})
	}

	// Validate password
	if strings.TrimSpace(registerRequest.Password) == "" {
		return response.SendValidationError(c, []types.ValidationError{
			{
				Field:   "password",
				Message: "Password is required",
				Value:   registerRequest.Password,
			},
		})
	}

	// Basic password validation
	if len(registerRequest.Password) < 6 {
		return response.SendValidationError(c, []types.ValidationError{
			{
				Field:   "password",
				Message: "Password must be at least 6 characters long",
				Value:   registerRequest.Password,
			},
		})
	}

	// Initialize auth service
	authService := &services.AuthService{Logger: logger}
	// Initialize cookie service
	cookieService := &services.CookieService{}

	// Attempt registration
	user, err := authService.Register(&registerRequest)
	if err != nil {
		logger.Error("Registration failed", "email", registerRequest.Email, "username", registerRequest.Username, "error", err)

		if errors.Is(err, lib.ErrUserAlreadyExists) {
			return response.Conflict(c, "User with this email or username already exists")
		}

		return response.InternalServerError(c, "An error occurred during registration")
	}

	// Generate tokens for the new user
	accessToken, err := authService.GenerateAccessToken(user)
	if err != nil {
		logger.AuditError("Failed to generate access token", "user_id", user.Id, "error", err)
		return response.InternalServerError(c, "Failed to generate access token")
	}

	refreshToken, err := authService.GenerateRefreshToken(user)
	if err != nil {
		logger.AuditError("Failed to generate refresh token", "user_id", user.Id, "error", err)
		return response.InternalServerError(c, "Failed to generate refresh token")
	}

	cookieService.SetAuthCookies(c, accessToken, refreshToken)

	return response.Success(c, user)
}

// RefreshToken handles token refresh using refresh tokens
func RefreshToken(c fiber.Ctx) error {
	logger := config.SetupLogger()

	token := c.Cookies(lib.RefreshTokenCookieName)

	// Initialize auth service
	authService := &services.AuthService{Logger: logger}

	// Initialize cookie service for setting new cookies
	cookieService := &services.CookieService{}

	// Refresh tokens with rotation
	authResponse, err := authService.RefreshToken(token)
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

	// Set new rotated tokens in secure cookies
	cookieService.SetAuthCookies(c, authResponse.AccessToken, authResponse.RefreshToken)

	return response.Success(c, authResponse)
}

// Me returns the current authenticated user's information
func Me(c fiber.Ctx) error {
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

	// Initialize auth service
	authService := &services.AuthService{Logger: logger}

	// Fetch user info
	user, err := authService.GetUserByID(claims.Sub)
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
func Logout(c fiber.Ctx) error {
	logger := config.SetupLogger()

	accessToken := c.Cookies(lib.AccessTokenCookieName)
	refreshToken := c.Cookies(lib.RefreshTokenCookieName)

	// Initialize services
	authService := &services.AuthService{Logger: logger}
	cookieService := &services.CookieService{}

	// Blacklist access token if present
	if strings.TrimSpace(accessToken) != "" {
		// Validate and blacklist access token
		_, err := authService.GetUserFromToken(accessToken)
		if err != nil {
			logger.Warn("Invalid access token during logout, clearing anyway", "error", err)
		} else {
			// Token is valid, blacklist it
			if err := authService.BlacklistToken(accessToken, true); err != nil {
				logger.AuditError("Failed to blacklist access token", "error", err)
				// Don't return error, continue with logout process
			}
		}
	}

	// Process refresh token if present
	if strings.TrimSpace(refreshToken) != "" {
		// Try to blacklist refresh token (may be invalid, but that's okay)
		if err := authService.BlacklistToken(refreshToken, false); err != nil {
			logger.Warn("Failed to blacklist refresh token, may already be invalid", "error", err)
			// Don't return error, continue with logout process
		}
	}

	// Always clear auth cookies regardless of token validity
	cookieService.ClearAuthCookies(c)

	return response.Success(c, types.LogoutResponse{
		Message: "Logged out successfully",
	})
}
