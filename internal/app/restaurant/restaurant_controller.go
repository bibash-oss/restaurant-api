package restaurant

import "github.com/gin-gonic/gin"

type RestaurantController struct {
	restaurantService *RestaurantService
}

func NewRestaurantController(restaurantService *RestaurantService) *RestaurantController {
	return &RestaurantController{restaurantService: restaurantService}
}

func (controller *RestaurantController) CreateRestaurant(c *gin.Context) {
	var restaurantReq CreateRestaurantRequest
	if err := c.ShouldBindJSON(&restaurantReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	createdRestaurant, err := controller.restaurantService.CreateRestaurant(&restaurantReq)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{
		"message": "Restaurant created successfully",
		"data":    createdRestaurant,
		"success": true,
	})
}

func (controller *RestaurantController) GetAllRestaurants(c *gin.Context) {
	restaurants, err := controller.restaurantService.GetAllRestaurants()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": restaurants, "success": true, "message": "Restaurants Fetched Successfully"})
}
