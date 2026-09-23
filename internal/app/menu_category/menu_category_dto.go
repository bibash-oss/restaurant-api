package menucategory

type CreateMenuCategoryRequest struct {
	RestaurantID string `json:"restaurantId" binding:"required"`
	Name         string `json:"name" binding:"required"`
}

type UpdateMenuCategoryRequest struct {
	Name *string `json:"name" binding:"required"`
}
