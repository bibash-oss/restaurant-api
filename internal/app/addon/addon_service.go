package addon

import (
	"fmt"

	"github.com/google/uuid"
)

type AddonService struct {
	addonRepo *AddonRepository
}

func NewAddonService(addonRepo *AddonRepository) *AddonService {
	return &AddonService{
		addonRepo: addonRepo,
	}
}

func (s *AddonService) CreateAddon(req *CreateAddonRequest) (*Addon, error) {
	restaurantID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	existing, err := s.addonRepo.FindByRestaurantAndName(restaurantID, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("addon '%s' already exists for this restaurant", req.Name)
	}

	a := &Addon{
		RestaurantID: restaurantID,
		Name:         req.Name,
		Price:        req.Price,
		IsActive:     true,
	}

	if err := s.addonRepo.CreateAddon(a); err != nil {
		return nil, err
	}

	return a, nil
}

func (s *AddonService) GetAddonsByRestaurantID(restaurantIDStr string) ([]Addon, error) {
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}
	return s.addonRepo.GetAddonsByRestaurantID(restaurantID)
}

func (s *AddonService) GetAddonByID(idStr string) (*Addon, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid addon ID")
	}
	a, err := s.addonRepo.GetAddonByID(id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, fmt.Errorf("addon not found")
	}
	return a, nil
}

func (s *AddonService) UpdateAddon(idStr string, req *UpdateAddonRequest) (*Addon, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid addon ID")
	}
	a, err := s.addonRepo.GetAddonByID(id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, fmt.Errorf("addon not found")
	}

	if req.Name != nil && *req.Name != a.Name {
		existing, err := s.addonRepo.FindByRestaurantAndName(a.RestaurantID, *req.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, fmt.Errorf("addon '%s' already exists for this restaurant", *req.Name)
		}
		a.Name = *req.Name
	}
	if req.Price != nil {
		a.Price = *req.Price
	}
	if req.IsActive != nil {
		a.IsActive = *req.IsActive
	}

	if err := s.addonRepo.UpdateAddon(a); err != nil {
		return nil, err
	}

	return a, nil
}

func (s *AddonService) DeleteAddon(idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid addon ID")
	}
	return s.addonRepo.DeleteAddon(id)
}

func (s *AddonService) AssignAddonsToMenuItem(menuItemIDStr string, req *AssignAddonsToMenuItemRequest) error {
	menuItemID, err := uuid.Parse(menuItemIDStr)
	if err != nil {
		return fmt.Errorf("invalid menu item ID")
	}

	restaurantID, err := s.addonRepo.GetMenuItemRestaurantID(menuItemID)
	if err != nil {
		return err
	}
	if restaurantID == nil {
		return fmt.Errorf("menu item not found")
	}

	addonUUIDs := make([]uuid.UUID, 0, len(req.AddonIDs))
	for _, idStr := range req.AddonIDs {
		aid, err := uuid.Parse(idStr)
		if err != nil {
			return fmt.Errorf("invalid addon ID: %s", idStr)
		}
		addonUUIDs = append(addonUUIDs, aid)
	}

	if len(addonUUIDs) > 0 {
		addons, err := s.addonRepo.GetAddonsByIDs(addonUUIDs)
		if err != nil {
			return err
		}
		if len(addons) != len(addonUUIDs) {
			return fmt.Errorf("one or more addons not found")
		}
		// Ensure addons belong to the same restaurant as the menu item
		for _, a := range addons {
			if a.RestaurantID != *restaurantID {
				return fmt.Errorf("addon '%s' does not belong to the same restaurant as the menu item", a.Name)
			}
		}
	}

	return s.addonRepo.AssignAddonsToMenuItem(menuItemID, addonUUIDs)
}

func (s *AddonService) GetAddonsByMenuItemID(menuItemIDStr string) ([]Addon, error) {
	menuItemID, err := uuid.Parse(menuItemIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid menu item ID")
	}
	return s.addonRepo.GetAddonsByMenuItemID(menuItemID)
}
