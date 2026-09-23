package menuitem

import (
	"github.com/gin-gonic/gin"
	"kitchen-api/internal/database"
)

func RegisterMenuItemRoutes(r *gin.Engine, db *database.OrmDb) {
	repo := NewMenuItemRepository(db)
	service := NewMenuItemService(repo)
	controller := NewMenuItemController(service)

	itemGroup := r.Group("/menu-items")
	{
		itemGroup.POST("", controller.CreateMenuItem)
		itemGroup.GET("/:id", controller.GetMenuItemByID)
		itemGroup.PATCH("/:id", controller.UpdateMenuItem)
		itemGroup.DELETE("/:id", controller.DeleteMenuItem)
	}

	restaurantItemGroup := r.Group("/restaurants/:id/menu-items")
	{
		restaurantItemGroup.GET("", controller.GetMenuItemsByRestaurant)
	}

	categoryItemGroup := r.Group("/menu-categories/:id/menu-items")
	{
		categoryItemGroup.GET("", controller.GetMenuItemsByCategory)
	}
}
