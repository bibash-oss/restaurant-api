package restaurant

type CreateRestaurantRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}
