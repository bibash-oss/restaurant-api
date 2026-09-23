package table

import (
	"github.com/gin-gonic/gin"
	"github.com/ktmbeestech/yanshi/database"
)

func RegisterTableRoutes(r *gin.Engine, db *database.OrmDb) {
	repo := NewTableRepository(db)
	service := NewTableService(repo)
	controller := NewTableController(service)

	tableGroup := r.Group("/tables")
	{
		tableGroup.POST("", controller.CreateTable)
		tableGroup.GET("/:id", controller.GetTableByID)
		tableGroup.PATCH("/:id", controller.UpdateTable)
		tableGroup.DELETE("/:id", controller.DeleteTable)
	}

	restaurantTableGroup := r.Group("/restaurants/:id/tables")
	{
		restaurantTableGroup.GET("", controller.GetTablesByRestaurant)
	}
}
