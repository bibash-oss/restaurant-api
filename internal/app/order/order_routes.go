package order

import (
	"kitchen-api/internal/app/addon"
	menuitem "kitchen-api/internal/app/menu_item"
	"kitchen-api/internal/app/table"

	"github.com/gin-gonic/gin"
	"kitchen-api/internal/database"
)

func RegisterOrderRoutes(r *gin.Engine, db *database.OrmDb) {
	orderRepo := NewOrderRepository(db)
	tableRepo := table.NewTableRepository(db)
	itemRepo := menuitem.NewMenuItemRepository(db)
	addonRepo := addon.NewAddonRepository(db)

	service := NewOrderService(orderRepo, tableRepo, itemRepo, addonRepo)
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

	tableOrderGroup := r.Group("/tables/:id/orders")
	{
		tableOrderGroup.GET("", controller.GetOrdersByTable)
	}
}
