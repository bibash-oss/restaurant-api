package table

type CreateTableRequest struct {
	RestaurantID string `json:"restaurantId" binding:"required"`
	Number       string `json:"number" binding:"required"`
}

type UpdateTableRequest struct {
	Number   *string `json:"number"`
	IsActive *bool   `json:"isActive"`
}
