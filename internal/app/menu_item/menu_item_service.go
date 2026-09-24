package menuitem

import (
	"fmt"
	"kitchen-api/internal/enums"

	"github.com/google/uuid"
)

type MenuItemService struct {
	itemRepo *MenuItemRepository
}

func NewMenuItemService(itemRepo *MenuItemRepository) *MenuItemService {
	return &MenuItemService{itemRepo: itemRepo}
}

func (s *MenuItemService) CreateMenuItem(req *CreateMenuItemRequest) (*MenuItem, error) {
	restaurantID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}

	menuType := enums.MenuTypeKitchen
	if req.MenuType != nil && *req.MenuType != "" {
		menuType = *req.MenuType
	}

	item := &MenuItem{
		RestaurantID: restaurantID,
		CategoryID:   categoryID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		ImageURL:     req.ImageURL,
		MenuType:     menuType,
		IsActive:     true,
	}

	if err := s.itemRepo.CreateMenuItem(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *MenuItemService) GetMenuItemsByRestaurantID(restaurantIDStr string, menuType string) ([]MenuItem, error) {
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}
	return s.itemRepo.GetMenuItemsByRestaurantID(restaurantID, menuType)
}

func (s *MenuItemService) GetMenuItemsByCategoryID(categoryIDStr string) ([]MenuItem, error) {
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}
	return s.itemRepo.GetMenuItemsByCategoryID(categoryID)
}

func (s *MenuItemService) GetMenuItemByID(idStr string) (*MenuItem, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid menu item ID")
	}
	item, err := s.itemRepo.GetMenuItemByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("menu item not found")
	}
	return item, nil
}

func (s *MenuItemService) UpdateMenuItem(idStr string, req *UpdateMenuItemRequest) (*MenuItem, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid menu item ID")
	}
	item, err := s.itemRepo.GetMenuItemByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("menu item not found")
	}

	if req.CategoryID != nil {
		catID, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID")
		}
		item.CategoryID = catID
	}
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = req.Description
	}
	if req.Price != nil {
		item.Price = *req.Price
	}
	if req.ImageURL != nil {
		item.ImageURL = req.ImageURL
	}
	if req.MenuType != nil && *req.MenuType != "" {
		item.MenuType = *req.MenuType
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}

	if err := s.itemRepo.UpdateMenuItem(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *MenuItemService) DeleteMenuItem(idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid menu item ID")
	}
	return s.itemRepo.DeleteMenuItem(id)
}
