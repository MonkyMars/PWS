package lib

import "errors"

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrHashingPassword    = errors.New("error hashing password")
	ErrCreateUser         = errors.New("error creating user")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("expired token")
	ErrGeneratingToken    = errors.New("error generating token")
	ErrValidatingToken    = errors.New("error validating token")
)

const (
	TableAuditLogs = "audit_logs"
)
