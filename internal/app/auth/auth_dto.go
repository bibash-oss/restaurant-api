package auth

import (
	"kitchen-api/internal/app/user"
	"kitchen-api/internal/enums"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string         `json:"token"`
	ExpiresIn int64          `json:"expiresIn"`
	Role      enums.UserRole `json:"role"`
}

type ProfileResponse struct {
	*user.Users
	QrUrl string `json:"qrUrl"`
}
