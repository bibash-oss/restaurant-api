package restaurant

import (
	"github.com/gin-gonic/gin"
	"kitchen-api/internal/database"
)

func RegisterRestaurantRoutes(r *gin.Engine, db *database.OrmDb) {
	repo := NewRestaurantRepository(db)
	service := NewRestaurantService(repo)
	controller := NewRestaurantController(service)

	restaurantGroup := r.Group("/restaurants")
	{
		restaurantGroup.POST("", controller.CreateRestaurant)
		restaurantGroup.GET("", controller.GetAllRestaurants)
	}
}
