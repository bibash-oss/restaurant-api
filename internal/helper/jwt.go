package helper

import (
	"errors"
	"time"

	"kitchen-api/internal/enums"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"kitchen-api/internal/config"
)

type JWTClaims struct {
	UserID       uuid.UUID      `json:"userId"`
	Email        string         `json:"email"`
	Role         enums.UserRole `json:"role"`
	RestaurantID *uuid.UUID     `json:"restaurantId,omitempty"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	secret := config.Default().GetString("jwt.secret")
	return []byte(secret)
}

func getExpirationDuration() time.Duration {
	hours := config.Default().GetInt("jwt.expirationHours")
	if hours <= 0 {
		hours = 72
	}
	return time.Duration(hours) * time.Hour
}

func GenerateToken(userID uuid.UUID, email string, role enums.UserRole, restaurantID *uuid.UUID) (string, int64, error) {
	duration := getExpirationDuration()
	expiresAt := time.Now().Add(duration)

	claims := JWTClaims{
		UserID:       userID,
		Email:        email,
		Role:         role,
		RestaurantID: restaurantID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(duration.Seconds()), nil
}

func ValidateToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
