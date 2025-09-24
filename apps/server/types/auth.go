package types

import (
	"time"

	"github.com/google/uuid"
)

type ArgonParams struct {
	Memory  uint32
	Time    uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

type AuthClaims struct {
	Sub   uuid.UUID `json:"sub"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
	Iat   time.Time `json:"iat"`
	Exp   time.Time `json:"exp"`
	Jti   uuid.UUID `json:"jti"`
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

type SecurityBreachResponse struct {
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type TokenFamilyBreach struct {
	UserID          uuid.UUID `json:"user_id"`
	SuspiciousJTI   string    `json:"suspicious_jti"`
	DetectedAt      time.Time `json:"detected_at"`
	ActionTaken     string    `json:"action_taken"`
	RecommendedStep string    `json:"recommended_step"`
}

type User struct {
	Id           uuid.UUID `json:"id" pg:"id,pk,type:uuid,default:gen_random_uuid()"`
	Username     string    `json:"username" pg:"username,unique,notnull"`
	Email        string    `json:"email" pg:"email,unique,notnull"`
	PasswordHash string    `json:"-" pg:"password_hash,notnull"`
	Role         string    `json:"role" pg:"role,notnull,default:'student'"`
}

type UserOAuthToken struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Provider     string    `json:"provider"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
