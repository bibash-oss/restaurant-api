package orderitem

import (
	"github.com/gin-gonic/gin"
	"github.com/ktmbeestech/yanshi/database"
)

func RegisterOrderItemRoutes(r *gin.Engine, db *database.OrmDb) {
	repo := NewOrderItemRepository(db)
	service := NewOrderItemService(repo)
	controller := NewOrderItemController(service)

	orderItemGroup := r.Group("/order-items")
	{
		orderItemGroup.GET("/:id", controller.GetOrderItemByID)
	}

	ordersGroup := r.Group("/orders/:id/items")
	{
		ordersGroup.GET("", controller.GetItemsByOrder)
	}
}
