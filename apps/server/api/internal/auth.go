package internal

import (
	"errors"
	"strings"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
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
	authService := &services.AuthService{}

	// Attempt login
	user, err := authService.Login(&authRequest)
	if err != nil {
		logger.Error("Login failed", "email", authRequest.Email, "error", err)

		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return response.Unauthorized(c, "Invalid email or password")
		}

		return response.InternalServerError(c, "An error occurred during login")
	}

	// Generate tokens
	accessToken, err := authService.GenerateAccessToken(user)
	if err != nil {
		logger.Error("Failed to generate access token", "user_id", user.Id, "error", err)
		return response.InternalServerError(c, "Failed to generate access token")
	}

	refreshToken, err := authService.GenerateRefreshToken(user)
	if err != nil {
		logger.Error("Failed to generate refresh token", "user_id", user.Id, "error", err)
		return response.InternalServerError(c, "Failed to generate refresh token")
	}

	// Create response
	authResponse := types.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	logger.Info("User logged in successfully", "user_id", user.Id, "email", authRequest.Email)

	return response.Success(c, authResponse)
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
	authService := &services.AuthService{}

	// Attempt registration
	user, err := authService.Register(&registerRequest)
	if err != nil {
		logger.Error("Registration failed", "email", registerRequest.Email, "username", registerRequest.Username, "error", err)

		if strings.Contains(err.Error(), "already exists") {
			return response.Conflict(c, "User with the given email or username already exists")
		}

		return response.InternalServerError(c, "An error occurred during registration")
	}

	// Generate tokens for the new user
	accessToken, err := authService.GenerateAccessToken(user)
	if err != nil {
		logger.Error("Failed to generate access token", "user_id", user.Id, "error", err)
		return response.InternalServerError(c, "Failed to generate access token")
	}

	refreshToken, err := authService.GenerateRefreshToken(user)
	if err != nil {
		logger.Error("Failed to generate refresh token", "user_id", user.Id, "error", err)
		return response.InternalServerError(c, "Failed to generate refresh token")
	}

	// Create response
	authResponse := types.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	logger.Info("User registered successfully", "user_id", user.Id, "email", registerRequest.Email, "username", registerRequest.Username)

	return response.Success(c, authResponse)
}

// RefreshToken handles token refresh using refresh tokens
func RefreshToken(c fiber.Ctx) error {
	logger := config.SetupLogger()

	var refreshRequest types.RefreshTokenRequest

	if err := c.Bind().Body(&refreshRequest); err != nil {
		logger.Error("Failed to parse refresh request", "error", err)
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate refresh token
	if strings.TrimSpace(refreshRequest.RefreshToken) == "" {
		return response.SendValidationError(c, []types.ValidationError{
			{
				Field:   "refresh_token",
				Message: "Refresh token is required",
				Value:   refreshRequest.RefreshToken,
			},
		})
	}

	// Initialize auth service
	authService := &services.AuthService{}

	// Refresh tokens
	authResponse, err := authService.RefreshToken(refreshRequest.RefreshToken)
	if err != nil {
		logger.Error("Token refresh failed", "error", err)
		return response.Unauthorized(c, "Invalid or expired refresh token")
	}

	logger.Info("Token refreshed successfully", "user_id", authResponse.User.Id)

	return response.Success(c, authResponse)
}

// Me returns the current authenticated user's information
func Me(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return response.Unauthorized(c, "Authorization header required")
	}

	// Check Bearer token format
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return response.Unauthorized(c, "Invalid authorization format")
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if strings.TrimSpace(token) == "" {
		return response.Unauthorized(c, "Token is required")
	}

	// Initialize auth service
	authService := &services.AuthService{}

	// Get user from token
	user, err := authService.GetUserFromToken(token)
	if err != nil {
		logger.Error("Failed to get user from token", "error", err)
		return response.Unauthorized(c, "Invalid or expired token")
	}

	logger.Info("User info retrieved", "user_id", user.Id)

	return response.Success(c, user)
}

// Logout handles user logout (token invalidation would be implemented here)
func Logout(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return response.Unauthorized(c, "Authorization header required")
	}

	// Check Bearer token format
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return response.Unauthorized(c, "Invalid authorization format")
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if strings.TrimSpace(token) == "" {
		return response.Unauthorized(c, "Token is required")
	}

	// Initialize auth service
	authService := &services.AuthService{}

	// Validate token to ensure it's valid before logout
	_, err := authService.GetUserFromToken(token)
	if err != nil {
		logger.Error("Invalid token during logout", "error", err)
		return response.Unauthorized(c, "Invalid or expired token")
	}

	// TODO: In a production system, you would:
	// 1. Add the token to a blacklist/revocation list
	// 2. Store blacklisted tokens in Redis with expiration
	// 3. Check blacklist in authentication middleware

	logger.Info("User logged out successfully")

	return response.Success(c, types.LogoutResponse{
		Message: "Logged out successfully",
	})
}
