package services

import (
	"fmt"
	"time"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func (a *AuthService) HashPassword(password string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash
}

func (a *AuthService) VerifyPassword(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

func (a *AuthService) GetRefreshTokenExpiration() time.Time {
	cfg := config.Get()
	return time.Now().Add(cfg.Auth.RefreshTokenExpiry)
}

func (a *AuthService) GetAccessTokenExpiration() time.Time {
	cfg := config.Get()
	return time.Now().Add(cfg.Auth.AccessTokenExpiry)
}

func (a *AuthService) GenerateAccessToken(user *types.User) (string, error) {
	secret := config.Get().Auth.AccessTokenSecret

	now := time.Now()
	exp := a.GetAccessTokenExpiration()

	claims := &types.AuthClaims{
		Sub:   user.Id,
		Email: user.Username,
		Role:  user.Role,
		Iat:   now,
		Exp:   exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   claims.Sub.String(),
		"email": claims.Email,
		"role":  claims.Role,
		"iat":   claims.Iat.Unix(),
		"exp":   claims.Exp.Unix(),
	})
	return token.SignedString([]byte(secret))
}

func (a *AuthService) GenerateRefreshToken(user *types.User) (string, error) {
	secret := config.Get().Auth.RefreshTokenSecret

	now := time.Now()
	exp := a.GetRefreshTokenExpiration()

	claims := &types.AuthClaims{
		Sub:   user.Id,
		Email: user.Username,
		Role:  user.Role,
		Iat:   now,
		Exp:   exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   claims.Sub.String(),
		"email": claims.Email,
		"role":  claims.Role,
		"iat":   claims.Iat.Unix(),
		"exp":   claims.Exp.Unix(),
	})
	return token.SignedString([]byte(secret))
}

func (a *AuthService) ParseToken(tokenStr string, isAccessToken bool) (*types.AuthClaims, error) {
	secret := config.Get().Auth.AccessTokenSecret
	if !isAccessToken {
		secret = config.Get().Auth.RefreshTokenSecret
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Safely extract and validate claims
		subStr, ok := claims["sub"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid sub claim")
		}

		sub, err := uuid.Parse(subStr)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID in sub claim: %w", err)
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid email claim")
		}

		role, ok := claims["role"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid role claim")
		}

		iat, ok := claims["iat"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid iat claim")
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid exp claim")
		}

		return &types.AuthClaims{
			Sub:   sub,
			Email: email,
			Role:  role,
			Iat:   time.Unix(int64(iat), 0),
			Exp:   time.Unix(int64(exp), 0),
		}, nil
	}
	return nil, jwt.ErrInvalidKey
}

func (a *AuthService) Login(authRequest *types.AuthRequest) (*types.User, error) {
	query := Query().SetOperation("SELECT").SetTable("users").SetSelect("id, email, username, password_hash, role")
	query.Where["email"] = authRequest.Email

	// Execute the query and get the user
	user, err := database.ExecuteQuery[types.User](query)
	if err != nil {
		return nil, err
	}

	if user.Single == nil {
		return nil, fmt.Errorf("user not found")
	}

	if !a.VerifyPassword(authRequest.Password, user.Single.PasswordHash) {
		return nil, bcrypt.ErrMismatchedHashAndPassword
	}

	return user.Single, nil
}

func (a *AuthService) Register(registerRequest *types.RegisterRequest) (*types.User, error) {
	// Check if user already exists
	query := Query().SetOperation("SELECT").SetTable("users").SetSelect("id")
	query.Where["email"] = registerRequest.Email

	existingUser, err := database.ExecuteQuery[types.User](query)
	if err == nil && existingUser.Single != nil {
		return nil, fmt.Errorf("user with email already exists")
	}

	// Also check username
	query = Query().SetOperation("SELECT").SetTable("users").SetSelect("id")
	query.Where["username"] = registerRequest.Username

	existingUser, err = database.ExecuteQuery[types.User](query)
	if err == nil && existingUser.Single != nil {
		return nil, fmt.Errorf("user with username already exists")
	}

	// Hash password
	hashedPassword := a.HashPassword(registerRequest.Password)

	// Create user
	newUserID := uuid.New()
	insertQuery := Query().SetOperation("INSERT").SetTable("users")
	insertQuery.Data = map[string]any{
		"id":            newUserID,
		"username":      registerRequest.Username,
		"email":         registerRequest.Email,
		"password_hash": hashedPassword,
		"role":          "student",
	}
	insertQuery.Returning = []string{"id", "username", "email", "role"}

	result, err := database.ExecuteQuery[types.User](insertQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if result.Single == nil {
		return nil, fmt.Errorf("failed to create user: no data returned")
	}

	return result.Single, nil
}

func (a *AuthService) RefreshToken(refreshTokenStr string) (*types.AuthResponse, error) {
	// Parse and validate refresh token
	claims, err := a.ParseToken(refreshTokenStr, false)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if token is expired
	if time.Now().After(claims.Exp) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Get user from database to ensure they still exist
	query := Query().SetOperation("SELECT").SetTable("users").SetSelect("id, email, username, role")
	query.Where["id"] = claims.Sub

	user, err := database.ExecuteQuery[types.User](query)
	if err != nil || user.Single == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new tokens
	accessToken, err := a.GenerateAccessToken(user.Single)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := a.GenerateRefreshToken(user.Single)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &types.AuthResponse{
		User:         user.Single,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *AuthService) GetUserFromToken(tokenStr string) (*types.User, error) {
	// Parse and validate access token
	claims, err := a.ParseToken(tokenStr, true)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	// Check if token is expired
	if time.Now().After(claims.Exp) {
		return nil, fmt.Errorf("access token expired")
	}

	// Get user from database
	query := Query().SetOperation("SELECT").SetTable("users").SetSelect("id, email, username, role")
	query.Where["id"] = claims.Sub

	user, err := database.ExecuteQuery[types.User](query)
	if err != nil || user.Single == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user.Single, nil
}
