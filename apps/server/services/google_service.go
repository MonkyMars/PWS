package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/google/uuid"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleService struct {
	logger *config.Logger
}

func NewGoogleService() *GoogleService {
	return &GoogleService{
		logger: config.SetupLogger(),
	}
}

// getGoogleOAuthConfig returns the OAuth config using values from the centralized config
func getGoogleOAuthConfig() *oauth2.Config {
	cfg := config.Get()
	return &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		Scopes: []string{
			// Full drive access for changing permissions and viewing.
			"https://www.googleapis.com/auth/drive",
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

// GenerateGoogleAuthURL generates an OAuth URL for the authenticated user
func (gs *GoogleService) GenerateGoogleAuthURL(userID uuid.UUID) (string, error) {
	// create state and persist it server-side (or in a secure cookie) mapped to the user ID
	state, err := gs.generateState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	// Store state -> user mapping in cache with expiry
	err = gs.saveOAuthState(userID, state)
	if err != nil {
		return "", fmt.Errorf("failed to save OAuth state: %w", err)
	}

	// request offline access to get refresh_token. prompt=consent ensures refresh token is returned
	googleOAuthConfig := getGoogleOAuthConfig()
	authURL := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return authURL, nil
}

// HandleGoogleCallback processes the OAuth callback and returns redirect URL
func (gs *GoogleService) HandleGoogleCallback(state, code string) (string, error) {
	ctx := context.Background()

	if state == "" || code == "" {
		return "", fmt.Errorf("state and code are required")
	}

	// Verify state maps to an authenticated user and is not expired
	userID, err := gs.getUserFromState(state)
	if err != nil {
		return "", fmt.Errorf("invalid or expired OAuth state: %w", err)
	}

	// Exchange the code for token
	googleOAuthConfig := getGoogleOAuthConfig()
	token, err := googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		log.Println("token exchange error:", err)
		return "", fmt.Errorf("failed to exchange token: %w", err)
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
		return "", fmt.Errorf("failed to save token: %w", err)
	}

	// Return redirect URL for frontend
	cfg := config.Get()
	return cfg.FrontendURL + "/dashboard", nil
}

// GetGoogleAccessToken gets a fresh access token for the user
func (gs *GoogleService) GetGoogleAccessToken(userID uuid.UUID) (map[string]any, error) {
	ctx := context.Background()

	refreshToken, err := gs.LoadUserRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load refresh token: %w", err)
	}

	if refreshToken == "" {
		return nil, fmt.Errorf("no linked Google account")
	}

	googleOAuthConfig := getGoogleOAuthConfig()
	ts := googleOAuthConfig.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	newToken, err := ts.Token()
	if err != nil {
		log.Printf("GetGoogleAccessToken: Token refresh error - %v", err)
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	log.Printf("GetGoogleAccessToken: Successfully refreshed token for user %s", userID.String())
	return map[string]any{
		"access_token": newToken.AccessToken,
		"expiry":       newToken.Expiry.Format(time.RFC3339),
		"token_type":   newToken.TokenType,
	}, nil
}

func (gs *GoogleService) SaveUserRefreshToken(userID uuid.UUID, refreshToken string) error {
	query := Query().SetOperation("insert").SetTable("user_oauth_tokens").SetData(map[string]any{
		"user_id":       userID,
		"provider":      "google",
		"refresh_token": refreshToken,
	})
	// On conflict (same user_id and provider), update the refresh_token and updated_at timestamp
	query.OnConflict = "(user_id, provider) DO UPDATE SET refresh_token = EXCLUDED.refresh_token, updated_at = NOW()"

	_, err := database.ExecuteQuery[types.UserOAuthToken](query)

	if err != nil {
		log.Printf("SaveUserRefreshToken: Failed to save token - %v", err)
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	log.Printf("SaveUserRefreshToken: Successfully saved token for user %s", userID.String())
	return nil
}

func (gs *GoogleService) LoadUserRefreshToken(userID uuid.UUID) (string, error) {
	query := Query().SetOperation("select").SetTable(lib.TableUserOAuthTokens).SetSelect(database.PrefixQuery(lib.TableUserOAuthTokens, []string{"refresh_token", "id"})).SetLimit(1)
	query.Where["user_oauth_tokens.user_id"] = userID
	query.Where["user_oauth_tokens.provider"] = "google"

	result, err := database.ExecuteQuery[types.GoogleRefreshTokenResponse](query)

	if err != nil {
		log.Printf("LoadUserRefreshToken: Failed to load token - %v", err)
		return "", lib.ErrFailedToRefreshToken
	}

	if len(result.Data) == 0 {
		return "", nil
	}

	return result.Single.RefreshToken, nil
}

func (gs *GoogleService) DeleteUserRefreshToken(userID uuid.UUID) error {
	query := Query().SetOperation("delete").SetTable(lib.TableUserOAuthTokens)
	query.Where["user_id"] = userID
	query.Where["provider"] = "google"

	_, err := database.ExecuteQuery[struct{}](query)

	if err != nil {
		return lib.ErrFailedToDeleteToken
	}

	return nil
}

func (gs *GoogleService) MakeFilePublic(userID uuid.UUID, fileID string) error {
	ctx := context.Background()

	// Get access token for this teacher
	tokenData, err := gs.GetGoogleAccessToken(userID)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Build oauth2 client with access token
	oauthCfg := getGoogleOAuthConfig()
	token := &oauth2.Token{AccessToken: tokenData["access_token"].(string)}
	client := oauthCfg.Client(ctx, token)

	// Create Drive service
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create drive client: %w", err)
	}

	// Add permission: anyone with link can read
	perm := &drive.Permission{
		Role: "reader",
		Type: "anyone",
	}

	_, err = srv.Permissions.Create(fileID, perm).Do()
	if err != nil {
		return fmt.Errorf("failed to set public permission: %w", err)
	}

	return nil
}

type GoogleServiceInterface interface {
	GenerateGoogleAuthURL(userID uuid.UUID) (string, error)
	HandleGoogleCallback(state, code string) (string, error)
	GetGoogleAccessToken(userID uuid.UUID) (map[string]any, error)
	SaveUserRefreshToken(userID uuid.UUID, refreshToken string) error
	LoadUserRefreshToken(userID uuid.UUID) (string, error)
	DeleteUserRefreshToken(userID uuid.UUID) error
	MakeFilePublic(userID uuid.UUID, fileID string) error
}
