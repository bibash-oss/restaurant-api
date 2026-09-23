package addon

import (
	"github.com/gin-gonic/gin"
	"github.com/ktmbeestech/yanshi/database"
)

func RegisterAddonRoutes(r *gin.Engine, db *database.OrmDb) {
	addonRepo := NewAddonRepository(db)
	service := NewAddonService(addonRepo)
	controller := NewAddonController(service)

	addonGroup := r.Group("/addons")
	{
		addonGroup.POST("", controller.CreateAddon)
		addonGroup.GET("/:id", controller.GetAddonByID)
		addonGroup.PATCH("/:id", controller.UpdateAddon)
		addonGroup.DELETE("/:id", controller.DeleteAddon)
	}

	restaurantAddonGroup := r.Group("/restaurants/:id/addons")
	{
		restaurantAddonGroup.GET("", controller.GetAddonsByRestaurant)
	}

	menuItemAddonGroup := r.Group("/menu-items/:id/addons")
	{
		menuItemAddonGroup.POST("", controller.AssignAddonsToMenuItem)
		menuItemAddonGroup.GET("", controller.GetAddonsByMenuItem)
	}
}
