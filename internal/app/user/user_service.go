package user

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *UserRepository
}

func NewUserService(userRepo *UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (service *UserService) CreateUser(userReq *CreateUserRequest) (*Users, error) {
	existingUser, err := service.userRepo.FindByEmail(userReq.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("User with this email already exists")
	}
	restaurantID, err := uuid.Parse(userReq.RestaurantID)
	if err != nil {
		return nil, fmt.Errorf("Invalid Restaurant ID")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &Users{
		Name:         userReq.Name,
		Email:        userReq.Email,
		Password:     string(hashedPassword),
		RestaurantID: &restaurantID,
		IsActive:     true,
	}

	if err := service.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) GetAllUsers() ([]Users, error) {
	return service.userRepo.GetUsers()
}

func (service *UserService) GetUserByEmail(email string) (*Users, error) {
	user, err := service.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("User not found")
	}
	return user, nil
}
