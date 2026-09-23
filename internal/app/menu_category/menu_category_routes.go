package menucategory

import (
	"github.com/gin-gonic/gin"
	"kitchen-api/internal/database"
)

func RegisterMenuCategoryRoutes(r *gin.Engine, db *database.OrmDb) {
	repo := NewMenuCategoryRepository(db)
	service := NewMenuCategoryService(repo)
	controller := NewMenuCategoryController(service)

	categoryGroup := r.Group("/menu-categories")
	{
		categoryGroup.POST("", controller.CreateCategory)
		categoryGroup.GET("/:id", controller.GetCategoryByID)
		categoryGroup.PATCH("/:id", controller.UpdateCategory)
		categoryGroup.DELETE("/:id", controller.DeleteCategory)
	}

	restaurantCategoryGroup := r.Group("/restaurants/:id/menu-categories")
	{
		restaurantCategoryGroup.GET("", controller.GetCategoriesByRestaurant)
	}
}
