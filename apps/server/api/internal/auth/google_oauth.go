package auth

import (
	"fmt"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/gofiber/fiber/v3"
)

// GoogleAuthURL handles getting the Google OAuth authorization URL
// GET /auth/google/url
func (ar *AuthRoutes) GoogleAuthURL(c fiber.Ctx) error {
	// Get user from auth middleware
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		msg := "Failed to get authenticated user claims for Google OAuth URL generation"
		return lib.HandleServiceError(c, err, msg)
	}

	// Generate OAuth URL
	authURL, err := ar.googleService.GenerateGoogleAuthURL(claims.Sub)
	if err != nil {
		msg := fmt.Sprintf("Failed to generate Google OAuth URL for user ID %s: %v", claims.Sub, err)
		return lib.HandleServiceError(c, err, msg)
	}

	return response.Success(c, authURL)
}

// GoogleAuthCallback handles the OAuth callback from Google
// GET /auth/google/callback
func (ar *AuthRoutes) GoogleAuthCallback(c fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	// Handle OAuth callback (this includes state validation and token exchange)
	redirectURL, err := ar.googleService.HandleGoogleCallback(state, code)
	if err != nil {
		msg := fmt.Sprintf("Failed to handle Google OAuth callback: %v", err)
		return lib.HandleServiceError(c, err, msg)
	}

	// Redirect to frontend success page
	return c.Redirect().To(redirectURL)
}

// GoogleAccessToken handles getting a fresh Google access token
// GET /auth/google/access-token
func (ar *AuthRoutes) GoogleAccessToken(c fiber.Ctx) error {
	// Get user from auth middleware
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		msg := "Failed to get authenticated user claims for Google access token"
		return lib.HandleServiceError(c, err, msg)
	}

	// Get access token
	tokenData, err := ar.googleService.GetGoogleAccessToken(claims.Sub)
	if err != nil {
		if err.Error() == "no linked Google account" {
			msg := fmt.Sprintf("User ID %s has no linked Google account", claims.Sub)
			return lib.HandleServiceError(c, lib.ErrNoLinkedAccount, msg)
		}
		msg := fmt.Sprintf("Failed to get Google access token for user ID %s: %v", claims.Sub, err)
		return lib.HandleServiceError(c, err, msg)
	}

	return response.Success(c, tokenData)
}

// GoogleUnlink handles unlinking a user's Google account
// DELETE /auth/google/unlink
func (ar *AuthRoutes) GoogleUnlink(c fiber.Ctx) error {
	// Get user from auth middleware
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		msg := "Failed to get authenticated user claims for Google account unlink"
		return lib.HandleServiceError(c, err, msg)
	}

	// Delete the user's refresh token
	err = ar.googleService.DeleteUserRefreshToken(claims.Sub)
	if err != nil {
		msg := fmt.Sprintf("Failed to unlink Google account for user ID %s: %v", claims.Sub, err)
		return lib.HandleServiceError(c, err, msg)
	}

	return response.Success(c, map[string]string{
		"message": "Google account unlinked successfully",
	})
}

// GoogleLinkStatus checks if user has linked their Google account
// GET /auth/google/status
func (ar *AuthRoutes) GoogleLinkStatus(c fiber.Ctx) error {
	// Get user from auth middleware
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		msg := "Failed to get authenticated user claims for Google link status check"
		return lib.HandleServiceError(c, err, msg)
	}

	// Check if user has a refresh token
	refreshToken, err := ar.googleService.LoadUserRefreshToken(claims.Sub)
	if err != nil {
		msg := fmt.Sprintf("Failed to check Google link status for user ID %s: %v", claims.Sub, err)
		return lib.HandleServiceError(c, err, msg)
	}

	isLinked := refreshToken != ""

	return response.Success(c, map[string]any{
		"linked":  isLinked,
		"user_id": claims.Sub,
	})
}
