package lib

import "errors"

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserNotFound         = errors.New("user not found")
	ErrFailedToRefreshToken = errors.New("failed to refresh token")
	ErrFailedToDeleteToken  = errors.New("failed to delete token")
)

const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleStudent = "student"
)

const (
	TableUsers           = "users"
	TableFiles           = "files"
	TableSubjects        = "subjects"
	TableUserOAuthTokens = "user_oauth_tokens"
	TableUserSubjects    = "user_subjects"
)
