package order

import (
	"kitchen-api/internal/app/addon"
	menuitem "kitchen-api/internal/app/menu_item"

	"github.com/gin-gonic/gin"
	"kitchen-api/internal/database"
)

func RegisterOrderRoutes(r *gin.Engine, db *database.OrmDb) *OrderService {
	orderRepo := NewOrderRepository(db)
	itemRepo := menuitem.NewMenuItemRepository(db)
	addonRepo := addon.NewAddonRepository(db)

	service := NewOrderService(orderRepo, itemRepo, addonRepo)
	controller := NewOrderController(service)

	orderGroup := r.Group("/orders")
	{
		orderGroup.POST("", controller.CreateOrder)
		orderGroup.GET("/:id", controller.GetOrderByID)
		orderGroup.PATCH("/:id/status", controller.UpdateOrderStatus)
	}

	restaurantOrderGroup := r.Group("/restaurants/:id/orders")
	{
		restaurantOrderGroup.GET("", controller.GetOrdersByRestaurant)
	}

	return service
}
