package auth

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// GoogleAuthURL handles getting the Google OAuth authorization URL
// GET /auth/google/url
func (ar *AuthRoutes) GoogleAuthURL(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Get user from auth middleware
	claimsInterface := c.Locals("claims")
	if claimsInterface == nil {
		logger.Error("No claims found in context for Google OAuth URL")
		return response.Unauthorized(c, "unauthenticated")
	}

	claims, ok := claimsInterface.(*types.AuthClaims)
	if !ok || claims == nil {
		logger.Error("Invalid claims type in context")
		return response.Unauthorized(c, "unauthenticated")
	}

	// Generate OAuth URL
	authURL, err := ar.googleService.GenerateGoogleAuthURL(claims.Sub)
	if err != nil {
		logger.AuditError("Failed to generate Google OAuth URL",
			"user_id", claims.Sub,
			"error", err)
		return response.InternalServerError(c, "failed to generate OAuth URL")
	}

	return response.Success(c, authURL)
}

// GoogleAuthCallback handles the OAuth callback from Google
// GET /auth/google/callback
func (ar *AuthRoutes) GoogleAuthCallback(c fiber.Ctx) error {
	logger := config.SetupLogger()

	state := c.Query("state")
	code := c.Query("code")

	// Handle OAuth callback (this includes state validation and token exchange)
	redirectURL, err := ar.googleService.HandleGoogleCallback(state, code)
	if err != nil {
		logger.AuditError("Failed to handle Google OAuth callback",
			"state", state,
			"code_present", code != "",
			"error", err)

		// Return error response for invalid callback
		return response.BadRequest(c, "OAuth callback failed: "+err.Error())
	}

	// Redirect to frontend success page
	return c.Redirect().To(redirectURL)
}

// GoogleAccessToken handles getting a fresh Google access token
// GET /auth/google/access-token
func (ar *AuthRoutes) GoogleAccessToken(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Get user from auth middleware
	claimsInterface := c.Locals("claims")
	if claimsInterface == nil {
		logger.Error("No claims found in context for Google access token")
		return response.Unauthorized(c, "unauthenticated")
	}

	claims, ok := claimsInterface.(*types.AuthClaims)
	if !ok || claims == nil {
		logger.Error("Invalid claims type in context")
		return response.Unauthorized(c, "unauthenticated")
	}

	// Get access token
	tokenData, err := ar.googleService.GetGoogleAccessToken(claims.Sub)
	if err != nil {
		logger.AuditError("Failed to get Google access token",
			"user_id", claims.Sub,
			"error", err)

		if err.Error() == "no linked Google account" {
			return response.Unauthorized(c, "no linked Google account")
		}
		return response.InternalServerError(c, "failed to refresh token")
	}

	return response.Success(c, tokenData)
}

// GoogleUnlink handles unlinking a user's Google account
// DELETE /auth/google/unlink
func (ar *AuthRoutes) GoogleUnlink(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Get user from auth middleware
	claimsInterface := c.Locals("claims")
	if claimsInterface == nil {
		logger.Error("No claims found in context for Google unlink")
		return response.Unauthorized(c, "unauthenticated")
	}

	// Delete the user's refresh token
	claims := claimsInterface.(*types.AuthClaims)
	err := ar.googleService.DeleteUserRefreshToken(claims.Sub)
	if err != nil {
		logger.AuditError("Failed to unlink Google account",
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
func (ar *AuthRoutes) GoogleLinkStatus(c fiber.Ctx) error {
	logger := config.SetupLogger()

	// Get user from auth middleware
	claimsInterface := c.Locals("claims")
	if claimsInterface == nil {
		logger.Error("No claims found in context for Google status check")
		return response.Unauthorized(c, "unauthenticated")
	}

	// Check if user has a refresh token
	claims := claimsInterface.(*types.AuthClaims)
	refreshToken, err := ar.googleService.LoadUserRefreshToken(claims.Sub)
	if err != nil {
		logger.AuditError("Failed to check Google link status",
			"user_id", claims.Sub,
			"error", err)
		return response.InternalServerError(c, "failed to check Google link status")
	}

	isLinked := refreshToken != ""

	return response.Success(c, map[string]any{
		"linked":  isLinked,
		"user_id": claims.Sub,
	})
}
