package restaurant

import (
	"fmt"
	"strings"
	"time"
)

type RestaurantService struct {
	restaurantRepo *RestaurantRepository
}

func NewRestaurantService(restaurantRepo *RestaurantRepository) *RestaurantService {
	return &RestaurantService{restaurantRepo: restaurantRepo}
}

func (service *RestaurantService) CreateRestaurant(userReq *CreateRestaurantRequest) (*Restaurants, error) {
	existingUser, err := service.restaurantRepo.FindByName(userReq.Name)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("Restaurant With This Name Already Exists")
	}
	slug := strings.ToLower(strings.ReplaceAll(userReq.Name, " ", "-"))
	slug = slug + "-" + time.Now().Format("20060102")
	restaurant := &Restaurants{
		Name:     userReq.Name,
		Slug:     slug,
		Address:  &userReq.Address,
		IsActive: true,
	}

	if err := service.restaurantRepo.CreateRestaurant(restaurant); err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (service *RestaurantService) GetAllRestaurants() ([]Restaurants, error) {
	return service.restaurantRepo.GetRestaurants()
}

func (service *RestaurantService) GetRestaurantByName(name string) (*Restaurants, error) {
	user, err := service.restaurantRepo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("Restaurant not found")
	}
	return user, nil
}
