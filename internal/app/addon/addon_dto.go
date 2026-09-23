package addon

type CreateAddonRequest struct {
	RestaurantID string  `json:"restaurantId" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Price        float64 `json:"price" binding:"required,gte=0"`
}

type UpdateAddonRequest struct {
	Name     *string  `json:"name"`
	Price    *float64 `json:"price" binding:"omitempty,gte=0"`
	IsActive *bool    `json:"isActive"`
}

type AssignAddonsToMenuItemRequest struct {
	AddonIDs []string `json:"addonIds" binding:"required"`
}
