package types

import (
	"time"

	"github.com/google/uuid"
)

type AuthClaims struct {
	Sub   uuid.UUID `json:"sub"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
	Iat   time.Time `json:"iat"`
	Exp   time.Time `json:"exp"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type User struct {
	Id           uuid.UUID `json:"id" pg:"id,pk,type:uuid,default:gen_random_uuid()"`
	Username     string    `json:"username" pg:"username,unique,notnull"`
	Email        string    `json:"email" pg:"email,unique,notnull"`
	PasswordHash []byte    `json:"-" pg:"password_hash,notnull"`
	Role         string    `json:"role" pg:"role,notnull,default:'student'"`
}
