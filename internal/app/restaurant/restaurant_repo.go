package restaurant

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ktmbeestech/yanshi/database"
	"gorm.io/gorm"
)

type RestaurantRepository struct {
	db *database.OrmDb
}

func NewRestaurantRepository(db *database.OrmDb) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

func (repo *RestaurantRepository) CreateRestaurant(Restaurant *Restaurants) error {
	return repo.db.OrmInstance.Create(Restaurant).Error
}

func (repo *RestaurantRepository) GetRestaurants() ([]Restaurants, error) {
	var Restaurants []Restaurants
	return Restaurants, repo.db.OrmInstance.Find(&Restaurants).Error
}

func (repo *RestaurantRepository) GetRestaurantById(id uuid.UUID) (Restaurants, error) {
	var Restaurant Restaurants
	return Restaurant, repo.db.OrmInstance.First(&Restaurant, id).Error
}

func (repo *RestaurantRepository) FindByName(name string) (*Restaurants, error) {
	var Restaurant Restaurants

	err := repo.db.OrmInstance.
		Where("name = ?", name).
		First(&Restaurant).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &Restaurant, err
}

func (repo *RestaurantRepository) UpdateRestaurant(Restaurant *Restaurants) error {
	return repo.db.OrmInstance.Save(Restaurant).Error
}

func (repo *RestaurantRepository) DeleteRestaurant(id uuid.UUID) error {
	return repo.db.OrmInstance.Delete(&Restaurants{}, id).Error
}
