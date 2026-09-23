package menuitem

type CreateMenuItemRequest struct {
	RestaurantID string  `json:"restaurantId" binding:"required"`
	CategoryID   string  `json:"categoryId" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Description  *string `json:"description"`
	Price        float64 `json:"price" binding:"required,gt=0"`
	ImageURL     *string `json:"imageUrl"`
}

type UpdateMenuItemRequest struct {
	CategoryID  *string  `json:"categoryId"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	ImageURL    *string  `json:"imageUrl"`
	IsActive    *bool    `json:"isActive"`
}
