package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

var DefaultParams = &types.ArgonParams{
	Memory:  64 * 1024, // 64 MB
	Time:    1,
	Threads: 4,
	KeyLen:  32,
	SaltLen: 16,
}

func generateSalt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var res byte = 0
	for i := range a {
		res |= a[i] ^ b[i] // Use i as index, not value
	}
	return res == 0
}

// AuthService provides authentication-related services
type AuthService struct {
	Logger *config.Logger
}

// HashPassword hashes a plain-text password and returns a string and possible error
func (a *AuthService) HashPassword(password string, p *types.ArgonParams) (string, error) {
	salt, err := generateSalt(p.SaltLen)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, p.KeyLen)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	// format: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
	params := fmt.Sprintf("m=%d,t=%d,p=%d", p.Memory, p.Time, p.Threads)
	encoded := fmt.Sprintf("$argon2id$v=19$%s$%s$%s", params, b64Salt, b64Hash)
	return encoded, nil
}

// ComparePasswordAndHash compares a plain-text password with a hashed password
// Returns true if they match, false otherwise + possible error
// Supports both bcrypt (legacy) and argon2 (new) password hashes
func (a *AuthService) ComparePasswordAndHash(password, encoded string) (bool, error) {
	// Check if it's an argon2 hash
	if strings.HasPrefix(encoded, "$argon2id$") {
		return a.compareArgon2Hash(password, encoded)
	}

	// Unknown hash format
	return false, fmt.Errorf("unsupported hash format: %s", encoded[:min(20, len(encoded))])
}

// compareArgon2Hash handles argon2 password comparison
func (a *AuthService) compareArgon2Hash(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("bad argon2 hash format: expected 6 parts, got %d", len(parts))
	}
	params := parts[3]
	saltB64 := parts[4]
	hashB64 := parts[5]
	var memory, time uint32
	var threads uint8

	for _, p := range strings.Split(params, ",") {
		kv := strings.Split(p, "=")
		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "m":
			v, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return false, fmt.Errorf("invalid memory parameter: %w", err)
			}
			memory = uint32(v)
		case "t":
			v, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return false, fmt.Errorf("invalid time parameter: %w", err)
			}
			time = uint32(v)
		case "p":
			v, err := strconv.ParseUint(kv[1], 10, 8)
			if err != nil {
				return false, fmt.Errorf("invalid threads parameter: %w", err)
			}
			threads = uint8(v)
		}
	}

	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, err
	}
	expected, err := base64.RawStdEncoding.DecodeString(hashB64)
	if err != nil {
		return false, err
	}
	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(expected)))
	return subtleCompare(hash, expected), nil
}

// GetRefreshTokenExpiration returns the expiration time for refresh tokens using configuration settings
func (a *AuthService) GetRefreshTokenExpiration() time.Time {
	cfg := config.Get()
	return time.Now().Add(cfg.Auth.RefreshTokenExpiry)
}

// GetAccessTokenExpiration returns the expiration time for access tokens using configuration settings
func (a *AuthService) GetAccessTokenExpiration() time.Time {
	cfg := config.Get()
	return time.Now().Add(cfg.Auth.AccessTokenExpiry)
}

// GenerateAccessToken generates a JWT access token for the given user
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
		Jti:   uuid.New(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   claims.Sub.String(),
		"email": claims.Email,
		"role":  claims.Role,
		"iat":   claims.Iat.Unix(),
		"exp":   claims.Exp.Unix(),
		"jti":   claims.Jti.String(),
	})
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken generates a JWT refresh token for the given user
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
		Jti:   uuid.New(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   claims.Sub.String(),
		"email": claims.Email,
		"role":  claims.Role,
		"iat":   claims.Iat.Unix(),
		"exp":   claims.Exp.Unix(),
		"jti":   claims.Jti.String(),
	})
	return token.SignedString([]byte(secret))
}

// ParseToken parses and validates a JWT token string and returns the claims
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

		jtiStr, ok := claims["jti"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid jti claim")
		}

		jti, err := uuid.Parse(jtiStr)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID in jti claim: %w", err)
		}

		return &types.AuthClaims{
			Sub:   sub,
			Email: email,
			Role:  role,
			Iat:   time.Unix(int64(iat), 0),
			Exp:   time.Unix(int64(exp), 0),
			Jti:   jti,
		}, nil
	}
	return nil, jwt.ErrInvalidKey
}

// Login authenticates a user and returns the user object if successful
func (a *AuthService) Login(authRequest *types.AuthRequest) (*types.User, error) {
	a.Logger.Info("Attempting login for user", "email", authRequest.Email)
	columns := []string{"id", "username", "email", "password_hash", "role"}
	query := Query().SetOperation("SELECT").SetTable("public.users").SetSelect(database.PrefixQuery("users", columns)).SetLimit(1)
	query.Where["public.users.email"] = authRequest.Email

	// Execute the query and get the user
	user, err := database.ExecuteQuery[types.User](query)
	if err != nil {
		return nil, err
	}

	if user.Single == nil {
		return nil, lib.ErrUserNotFound
	}

	isValid, err := a.ComparePasswordAndHash(authRequest.Password, user.Single.PasswordHash)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, lib.ErrInvalidCredentials
	}

	// Remove password hash before returning user object
	user.Single.PasswordHash = ""

	return user.Single, nil
}

// Register creates a new user account and returns the user object if successful
func (a *AuthService) Register(registerRequest *types.RegisterRequest) (*types.User, error) {
	// Check if user already exists
	query := Query().SetOperation("SELECT").SetTable("users").SetSelect([]string{"public.users.id"}).SetLimit(1)
	query.Where["public.users.email"] = registerRequest.Email

	existingUser, err := database.ExecuteQuery[types.User](query)
	if err == nil && existingUser.Single != nil {
		return nil, lib.ErrUserAlreadyExists
	}

	// Also check username
	query = Query().SetOperation("SELECT").SetTable("users").SetSelect([]string{"public.users.id"}).SetLimit(1)
	query.Where["public.users.username"] = registerRequest.Username

	existingUser, err = database.ExecuteQuery[types.User](query)
	if err == nil && existingUser.Single != nil {
		return nil, lib.ErrUsernameTaken
	}

	// Hash password
	hashedPassword, err := a.HashPassword(registerRequest.Password, DefaultParams)
	if err != nil {
		a.Logger.Error("Failed to hash password during registration", "error", err)
		return nil, lib.ErrHashingPassword
	}

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
	if err != nil || result.Single == nil {
		a.Logger.Error("Failed to create user during registration", "error", err)
		return nil, lib.ErrCreateUser
	}

	return result.Single, nil
}

// RefreshToken validates a refresh token and returns new JWT tokens with rotation for security
func (a *AuthService) RefreshToken(refreshTokenStr string) (*types.AuthResponse, error) {
	// Parse and validate refresh token
	claims, err := a.ParseToken(refreshTokenStr, false)
	if err != nil {
		a.Logger.Warn("Invalid refresh token during refresh attempt", "error", err)
		return nil, lib.ErrInvalidToken
	}

	// Check if token is expired
	if time.Now().After(claims.Exp) {
		return nil, lib.ErrExpiredToken
	}

	// Check if token is already blacklisted (detects token reuse/replay attacks)
	cacheService := CacheService{}
	blacklisted, err := cacheService.IsTokenBlacklisted(claims.Jti.String())
	if err != nil {
		a.Logger.Error("Failed to check token blacklist during refresh", "error", err, "jti", claims.Jti.String())
		return nil, lib.ErrValidatingToken
	}

	if blacklisted {
		a.Logger.Warn("Attempted reuse of blacklisted refresh token - possible replay attack",
			"jti", claims.Jti.String(),
			"user_id", claims.Sub,
			"user_email", claims.Email)
		return nil, lib.ErrInvalidToken
	}

	// Get user from database to ensure they still exist
	columns := []string{"id", "username", "email", "role"}
	query := Query().SetOperation("SELECT").SetTable("public.users").SetSelect(database.PrefixQuery("users", columns)).SetLimit(1)
	query.Where["public.users.id"] = claims.Sub.String()

	user, err := database.ExecuteQuery[types.User](query)
	if err != nil || user.Single == nil {
		return nil, lib.ErrUserNotFound
	}

	// SECURITY: Immediately blacklist the old refresh token to prevent reuse
	err = a.BlacklistToken(refreshTokenStr, false)
	if err != nil {
		a.Logger.Error("Failed to blacklist old refresh token during rotation",
			"error", err,
			"jti", claims.Jti.String(),
			"user_id", claims.Sub)
		// Don't fail the request - better to have working auth than strict security here
		// But log this as it could indicate Redis issues
	}

	// Generate new access token
	accessToken, err := a.GenerateAccessToken(user.Single)
	if err != nil {
		return nil, lib.ErrGeneratingToken
	}

	// Generate new refresh token (token rotation)
	newRefreshToken, err := a.GenerateRefreshToken(user.Single)
	if err != nil {
		return nil, lib.ErrGeneratingToken
	}

	return &types.AuthResponse{
		User:         user.Single,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// GetUserFromToken extracts the user information from a valid JWT access token
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
	columns := []string{"id", "username", "email", "password_hash", "role"}
	query := Query().SetOperation("SELECT").SetTable("public.users").SetSelect(database.PrefixQuery("users", columns)).SetLimit(1)
	query.Where["public.users.id"] = claims.Sub

	user, err := database.ExecuteQuery[types.User](query)
	if err != nil || user.Single == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user.Single, nil
}

func (a *AuthService) GetUserByID(userID uuid.UUID) (*types.User, error) {
	// Get user from database
	columns := []string{"id", "username", "email", "role"}
	query := Query().SetOperation("SELECT").SetTable("public.users").SetSelect(database.PrefixQuery("users", columns)).SetLimit(1)
	query.Where["public.users.id"] = userID

	user, err := database.ExecuteQuery[types.User](query)
	if err != nil || user.Single == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user.Single, nil
}

// BlacklistToken adds a token to the blacklist to prevent reuse
func (a *AuthService) BlacklistToken(tokenStr string, isAccessToken bool) error {
	claims, err := a.ParseToken(tokenStr, isAccessToken)
	if err != nil {
		return err
	}

	cacheService := CacheService{}

	// Store the token's JTI in Redis until it expires
	return cacheService.BlacklistToken(claims.Jti.String(), claims.Exp)
}
