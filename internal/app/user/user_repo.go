package user

import (
	"errors"

	"github.com/google/uuid"
	"kitchen-api/internal/database"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *database.OrmDb
}

func NewUserRepository(db *database.OrmDb) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(user *Users) error {
	return repo.db.OrmInstance.Create(user).Error
}

func (repo *UserRepository) GetUsers() ([]Users, error) {
	var users []Users
	return users, repo.db.OrmInstance.Preload("Restaurant").Find(&users).Error
}

func (repo *UserRepository) GetUserById(id uuid.UUID) (Users, error) {
	var user Users
	return user, repo.db.OrmInstance.Preload("Restaurant").First(&user, id).Error
}

func (repo *UserRepository) FindByEmail(email string) (*Users, error) {
	var user Users

	err := repo.db.OrmInstance.
		Preload("Restaurant").
		Where("email = ?", email).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, err
}

func (repo *UserRepository) UpdateUser(user *Users) error {
	return repo.db.OrmInstance.Save(user).Error
}

func (repo *UserRepository) DeleteUser(id uuid.UUID) error {
	return repo.db.OrmInstance.Delete(&Users{}, id).Error
}
