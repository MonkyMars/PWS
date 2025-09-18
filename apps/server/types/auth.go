package types

import (
	"time"

	"github.com/google/uuid"
)

type AuthClaims struct {
	Sub uuid.UUID `json:"sub"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
	Iat	 time.Time     `json:"iat"`
	Exp  time.Time     `json:"exp"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id           uuid.UUID `json:"id" pg:"id,pk,type:uuid,default:gen_random_uuid()"`
	Username     string    `json:"username" pg:"username,unique,notnull"`
	PasswordHash []byte    `json:"-" pg:"password_hash,notnull"`
	Role         string    `json:"role" pg:"role,notnull,default:'student'"`
}
