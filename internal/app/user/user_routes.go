package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ktmbeestech/yanshi/database"
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
