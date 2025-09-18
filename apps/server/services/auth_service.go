package services

import (
	"time"

	"github.com/MonkyMars/PWS/config"
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
		"sub":   claims.Sub,
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
		"sub":   claims.Sub,
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
		return &types.AuthClaims{
			Sub:   claims["sub"].(uuid.UUID),
			Email: claims["email"].(string),
			Role:  claims["role"].(string),
			Iat:   time.Unix(int64(claims["iat"].(float64)), 0),
			Exp:   time.Unix(int64(claims["exp"].(float64)), 0),
		}, nil
	}
	return nil, jwt.ErrInvalidKey
}

