package auth

import (
	"kitchen-api/internal/app/user"
	"kitchen-api/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/ktmbeestech/yanshi/database"
)

func RegisterAuthRoutes(r *gin.Engine, db *database.OrmDb) {
	userRepo := user.NewUserRepository(db)
	service := NewAuthService(userRepo)
	controller := NewAuthController(service)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", controller.Login)
		authGroup.GET("/me", middlewares.AuthMiddleware(), controller.Me)
	}

	// Also support /users/login as an alias
	r.POST("/users/login", controller.Login)
}
