package internal

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// GoogleAuthURL handles getting the Google OAuth authorization URL
// GET /auth/google/url
func GoogleAuthURL(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Initialize Google service
	googleService := &services.GoogleService{}

	// Get OAuth URL (this includes user authentication check)
	err := googleService.GoogleAuthURL(c)
	if err != nil {
		logger.Error("Failed to generate Google OAuth URL", "error", err)
		return err // GoogleService already returns proper response
	}

	return nil // GoogleService already sent response
}

// GoogleAuthCallback handles the OAuth callback from Google
// GET /auth/google/callback
func GoogleAuthCallback(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Initialize Google service
	googleService := &services.GoogleService{}

	// Handle OAuth callback (this includes state validation and token exchange)
	err := googleService.GoogleAuthCallback(c)
	if err != nil {
		logger.Error("Failed to handle Google OAuth callback",
			"state", c.Query("state"),
			"code_present", c.Query("code") != "",
			"error", err)
		return err // GoogleService already returns proper response
	}

	return nil // GoogleService already sent response or redirect
}

// GoogleAccessToken handles getting a fresh Google access token
// GET /auth/google/access-token
func GoogleAccessToken(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Initialize Google service
	googleService := &services.GoogleService{}

	// Get access token (this includes user authentication check and refresh token validation)
	err := googleService.GoogleAccessToken(c)
	if err != nil {
		logger.Error("Failed to get Google access token", "error", err)
		return err // GoogleService already returns proper response
	}

	return nil // GoogleService already sent response
}

// GoogleUnlink handles unlinking a user's Google account
// DELETE /auth/google/unlink
func GoogleUnlink(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Get user from auth middleware
	claimsInterface := c.Locals("claims")
	if claimsInterface == nil {
		logger.Error("No claims found in context for Google unlink")
		return response.Unauthorized(c, "unauthenticated")
	}

	// Initialize Google service
	googleService := &services.GoogleService{}

	// Delete the user's refresh token
	claims := claimsInterface.(*types.AuthClaims)
	err := googleService.DeleteUserRefreshToken(claims.Sub)
	if err != nil {
		logger.Error("Failed to unlink Google account",
			"user_id", claims.Sub,
			"error", err)
		return response.InternalServerError(c, "failed to unlink Google account")
	}

	logger.Info("Google account unlinked successfully", "user_id", claims.Sub)
	return response.Success(c, map[string]string{
		"message": "Google account unlinked successfully",
	})
}

// GoogleLinkStatus checks if user has linked their Google account
// GET /auth/google/status
func GoogleLinkStatus(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Get user from auth middleware
	claimsInterface := c.Locals("claims")
	if claimsInterface == nil {
		logger.Error("No claims found in context for Google status check")
		return response.Unauthorized(c, "unauthenticated")
	}

	// Initialize Google service
	googleService := &services.GoogleService{}

	// Check if user has a refresh token
	claims := claimsInterface.(*types.AuthClaims)
	refreshToken, err := googleService.LoadUserRefreshToken(claims.Sub)
	if err != nil {
		logger.Error("Failed to check Google link status",
			"user_id", claims.Sub,
			"error", err)
		return response.InternalServerError(c, "failed to check Google link status")
	}

	isLinked := refreshToken != ""

	return response.Success(c, map[string]interface{}{
		"linked":  isLinked,
		"user_id": claims.Sub,
	})
}
