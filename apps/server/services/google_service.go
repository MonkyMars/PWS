package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleService struct{}

// getGoogleOAuthConfig returns the OAuth config using values from the centralized config
func getGoogleOAuthConfig() *oauth2.Config {
	cfg := config.Get()
	return &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		Scopes: []string{
			// minimal scope for picker & downloading file metadata
			"https://www.googleapis.com/auth/drive.readonly",
			// change permissions: "https://www.googleapis.com/auth/drive"
		},
		Endpoint:    google.Endpoint,
		RedirectURL: cfg.Google.RedirectURL,
	}
}

// generateState creates a CSRF state token
func (gs *GoogleService) generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// saveOAuthState saves the OAuth state mapped to user ID in cache with expiry
func (gs *GoogleService) saveOAuthState(userID uuid.UUID, state string) error {
	cacheService := &CacheService{}
	key := fmt.Sprintf("oauth_state:%s", state)

	// Store for 10 minutes (OAuth flow should complete quickly)
	return cacheService.Set(key, userID.String(), 10*time.Minute)
}

// getUserFromState retrieves and validates the user ID from OAuth state
func (gs *GoogleService) getUserFromState(state string) (uuid.UUID, error) {
	cacheService := &CacheService{}
	key := fmt.Sprintf("oauth_state:%s", state)

	userIDStr, err := cacheService.Get(key)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to retrieve state: %w", err)
	}
	if userIDStr == "" {
		return uuid.Nil, fmt.Errorf("invalid or expired state")
	}

	// Delete the state after use (one-time use)
	_ = cacheService.Delete(key)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in state: %w", err)
	}

	return userID, nil
}

// GET /api/auth/google/url
// Returns an OAuth URL to redirect the user to Google for consent.
func (gs *GoogleService) GoogleAuthURL(c fiber.Ctx) error {
	// You must be logged in with your app to link Google account:
	claimsInterface := c.Locals("claims")
	if claimsInterface == nil {
		return response.Unauthorized(c, "unauthenticated")
	}

	claims, ok := claimsInterface.(*types.AuthClaims)
	if claims == nil || !ok {
		return response.Unauthorized(c, "unauthenticated")
	}

	// create state and persist it server-side (or in a secure cookie) mapped to the user ID
	state, err := gs.generateState()
	if err != nil {
		return response.InternalServerError(c, "failed to generate state: "+err.Error())
	}

	// Store state -> user mapping in cache with expiry
	err = gs.saveOAuthState(claims.Sub, state)
	if err != nil {
		return response.InternalServerError(c, "failed to save OAuth state: "+err.Error())
	}

	// request offline access to get refresh_token. prompt=consent ensures refresh token is returned
	googleOAuthConfig := getGoogleOAuthConfig()
	authURL := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return response.Success(c, authURL)
}

// GET /api/auth/google/callback
// Handles the OAuth callback, exchanges code for tokens, stores refresh token server-side
func (gs *GoogleService) GoogleAuthCallback(c fiber.Ctx) error {
	ctx := context.Background()
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		return response.BadRequest(c, "state and code are required")
	}

	// Verify state maps to an authenticated user and is not expired
	userID, err := gs.getUserFromState(state)
	if err != nil {
		return response.BadRequest(c, "invalid or expired OAuth state: "+err.Error())
	}

	// Exchange the code for token
	googleOAuthConfig := getGoogleOAuthConfig()
	token, err := googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		log.Println("token exchange error:", err)
		return response.InternalServerError(c, "failed to exchange token: "+err.Error())
	}

	// token.RefreshToken will be non-empty when you correctly requested offline access and user consented
	if token.RefreshToken == "" {
		// Could be because the user previously granted permissions.
		// Consider using prompt=consent or handle reconsent flow.
		log.Println("no refresh token returned — might be previously granted")
	}

	// IMPORTANT: Save refresh token securely server-side (encrypt at rest)
	err = gs.SaveUserRefreshToken(userID, token.RefreshToken)
	if err != nil {
		log.Println("failed to save refresh token:", err)
		return response.InternalServerError(c, "failed to save token: "+err.Error())
	}

	// Redirect to frontend app (or render a success page)
	cfg := config.Get()
	return c.Redirect().To(cfg.FrontendURL + "/google-linked?success=1")
}

func (gs *GoogleService) GoogleAccessToken(c fiber.Ctx) error {
	ctx := context.Background()
	claimsInterface := c.Locals("claims")
	if claimsInterface == nil {
		return response.Unauthorized(c, "unauthenticated")
	}
	claims := claimsInterface.(*types.AuthClaims)
	userID := claims.Sub

	refreshToken, err := gs.LoadUserRefreshToken(userID)
	if err != nil || refreshToken == "" {
		return response.Unauthorized(c, "no linked Google account")
	}

	googleOAuthConfig := getGoogleOAuthConfig()
	ts := googleOAuthConfig.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	newToken, err := ts.Token()
	if err != nil {
		log.Println("token refresh error:", err)
		return response.InternalServerError(c, "failed to refresh token: "+err.Error())
	}

	// newToken.AccessToken is short-lived (~1 hour). Send to frontend to pass to Picker.
	return response.Success(c, map[string]any{
		"access_token": newToken.AccessToken,
		"expiry":       newToken.Expiry.Format(time.RFC3339),
		"token_type":   newToken.TokenType,
	})
}

func (gs *GoogleService) SaveUserRefreshToken(userID uuid.UUID, refreshToken string) error {
	query := Query().SetOperation("upsert").SetTable(lib.TableUserOAuthTokens).SetData(map[string]any{
		"user_id":       userID,
		"provider":      "google",
		"refresh_token": refreshToken,
	})

	_, err := database.ExecuteQuery[any](query)
	if err != nil {
		return err
	}
	return nil
}

func (gs *GoogleService) LoadUserRefreshToken(userID uuid.UUID) (string, error) {
	query := Query().SetOperation("select").SetTable(lib.TableUserOAuthTokens).SetLimit(1)
	query.Where["user_id"] = userID
	query.Where["provider"] = "google"

	data, err := database.ExecuteQuery[struct {
		RefreshToken string `json:"refresh_token"`
	}](query)
	if err != nil {
		return "", err
	}
	if data.Single == nil {
		return "", nil // no token found
	}
	return data.Single.RefreshToken, nil
}

func (gs *GoogleService) DeleteUserRefreshToken(userID uuid.UUID) error {
	query := Query().SetOperation("delete").SetTable(lib.TableUserOAuthTokens)
	query.Where["user_id"] = userID
	query.Where["provider"] = "google"

	_, err := database.ExecuteQuery[any](query)
	if err != nil {
		return err
	}
	return nil
}
