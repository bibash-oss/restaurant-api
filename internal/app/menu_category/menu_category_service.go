package menucategory

import (
	"fmt"

	"github.com/google/uuid"
)

type MenuCategoryService struct {
	categoryRepo *MenuCategoryRepository
}

func NewMenuCategoryService(categoryRepo *MenuCategoryRepository) *MenuCategoryService {
	return &MenuCategoryService{categoryRepo: categoryRepo}
}

func (s *MenuCategoryService) CreateCategory(req *CreateMenuCategoryRequest) (*MenuCategory, error) {
	restaurantID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	existing, err := s.categoryRepo.FindByRestaurantAndName(restaurantID, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("category '%s' already exists for this restaurant", req.Name)
	}

	c := &MenuCategory{
		RestaurantID: restaurantID,
		Name:         req.Name,
	}

	if err := s.categoryRepo.CreateCategory(c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *MenuCategoryService) GetCategoriesByRestaurantID(restaurantIDStr string) ([]MenuCategory, error) {
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}
	return s.categoryRepo.GetCategoriesByRestaurantID(restaurantID)
}

func (s *MenuCategoryService) GetCategoryByID(idStr string) (*MenuCategory, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}
	c, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("category not found")
	}
	return c, nil
}

func (s *MenuCategoryService) UpdateCategory(idStr string, req *UpdateMenuCategoryRequest) (*MenuCategory, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}
	c, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("category not found")
	}

	if req.Name != nil && *req.Name != c.Name {
		existing, err := s.categoryRepo.FindByRestaurantAndName(c.RestaurantID, *req.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, fmt.Errorf("category '%s' already exists for this restaurant", *req.Name)
		}
		c.Name = *req.Name
	}

	if err := s.categoryRepo.UpdateCategory(c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *MenuCategoryService) DeleteCategory(idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid category ID")
	}
	return s.categoryRepo.DeleteCategory(id)
}
