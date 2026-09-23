package user

import (
	"github.com/gin-gonic/gin"
	"kitchen-api/internal/database"
)

func RegisterUserRoutes(r *gin.Engine, db *database.OrmDb) {
	repo := NewUserRepository(db)
	service := NewUserService(repo)
	controller := NewUserController(service)

	userGroup := r.Group("/users")
	{
		userGroup.POST("", controller.CreateUser)
		userGroup.GET("", controller.GetAllUsers)
	}
}
