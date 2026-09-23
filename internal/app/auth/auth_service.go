package auth

import (
	"errors"
	"fmt"

	"kitchen-api/internal/app/user"
	"kitchen-api/internal/helper"
	"strings"

	"github.com/google/uuid"
	"kitchen-api/internal/config"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *user.UserRepository
}

func NewAuthService(userRepo *user.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	u, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("invalid email or password")
	}

	if !u.IsActive {
		return nil, errors.New("account is inactive, please contact your administrator")
	}

	// Verify password: first try bcrypt comparison
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		// Fallback check for any existing plain-text password from previous tests
		if u.Password == req.Password {
			// Upgrade plain password to bcrypt hash in background
			if hashed, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); hashErr == nil {
				u.Password = string(hashed)
				_ = s.userRepo.UpdateUser(u)
			}
		} else {
			return nil, errors.New("invalid email or password")
		}
	}

	token, expiresIn, err := helper.GenerateToken(u.ID, u.Email, u.Role, u.RestaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate authentication token: %w", err)
	}

	return &LoginResponse{
		Token:     token,
		ExpiresIn: expiresIn,
		Role:      u.Role,
	}, nil
}

func (s *AuthService) GetProfile(userID uuid.UUID) (*ProfileResponse, error) {
	u, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	baseURL := config.Default().GetString("app.url")
	baseURL = strings.TrimRight(baseURL, "/")

	var qrURL string
	if u.RestaurantID != nil {
		qrURL = fmt.Sprintf("%s/restaurant/%s", baseURL, u.RestaurantID.String())
	} else if u.Restaurant != nil && u.Restaurant.ID != uuid.Nil {
		qrURL = fmt.Sprintf("%s/restaurant/%s", baseURL, u.Restaurant.ID.String())
	}

	return &ProfileResponse{
		Users: &u,
		QrUrl: qrURL,
	}, nil
}
